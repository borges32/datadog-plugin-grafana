package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(settings.JSONData, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSONData: %w", err)
	}

	apiURL, _ := jsonData["url"].(string)
	if apiURL == "" {
		apiURL = "https://api.datadoghq.com/"
	}

	apiKey := settings.DecryptedSecureJSONData["apiKey"]
	applicationKey := settings.DecryptedSecureJSONData["applicationKey"]

	return &Datasource{
		settings:       settings,
		apiURL:         apiURL,
		apiKey:         apiKey,
		applicationKey: applicationKey,
	}, nil
}

// Datasource is a Datadog datasource which can respond to data queries and reports its health.
type Datasource struct {
	settings       backend.DataSourceInstanceSettings
	apiURL         string
	apiKey         string
	applicationKey string
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	Query string `json:"query"`
}

type datadogResponse struct {
	Status      string          `json:"status"`
	Series      []datadogSeries `json:"series"`
	Message     string          `json:"message"`
	ErrorType   string          `json:"error_type"`
	ErrorDetail string          `json:"error"`
}

type datadogSeries struct {
	Metric      string      `json:"metric"`
	DisplayName string      `json:"display_name"`
	TagSet      []string    `json:"tag_set"`
	Pointlist   [][]float64 `json:"pointlist"`
	Scope       string      `json:"scope"`
	Expression  string      `json:"expression"`
}

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	if qm.Query == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "query is required")
	}

	// Convert time range to Unix timestamps (Epoch)
	from := query.TimeRange.From.Unix()
	to := query.TimeRange.To.Unix()

	// Build Datadog API URL
	baseURL := d.apiURL
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}
	apiEndpoint := baseURL + "api/v1/query"

	// Build query parameters
	params := url.Values{}
	params.Add("from", strconv.FormatInt(from, 10))
	params.Add("to", strconv.FormatInt(to, 10))
	params.Add("query", qm.Query)

	fullURL := apiEndpoint + "?" + params.Encode()

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("error creating request: %v", err.Error()))
	}

	// Add Datadog headers
	httpReq.Header.Set("DD-API-KEY", d.apiKey)
	httpReq.Header.Set("DD-APPLICATION-KEY", d.applicationKey)

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("error executing request: %v", err.Error()))
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("error reading response: %v", err.Error()))
	}

	// Check HTTP status
	if httpResp.StatusCode != http.StatusOK {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("Datadog API error (status %d): %s", httpResp.StatusCode, string(body)))
	}

	// Parse Datadog response
	var ddResp datadogResponse
	err = json.Unmarshal(body, &ddResp)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("error parsing response: %v", err.Error()))
	}

	if ddResp.Status != "ok" {
		errMsg := ddResp.Message
		if ddResp.ErrorDetail != "" {
			errMsg = ddResp.ErrorDetail
		}
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("Datadog query error: %s", errMsg))
	}

	// Convert Datadog series to Grafana frames
	for _, series := range ddResp.Series {
		// Use scope as the series name (e.g., "host:AH-CW-AP-104,instance:3")
		seriesName := series.Scope

		frame := data.NewFrame("")

		// Build time and value arrays from pointlist
		times := make([]time.Time, 0, len(series.Pointlist))
		values := make([]float64, 0, len(series.Pointlist))

		for _, point := range series.Pointlist {
			if len(point) >= 2 {
				// point[0] is timestamp in milliseconds, point[1] is value
				timestamp := time.Unix(0, int64(point[0])*int64(time.Millisecond))
				times = append(times, timestamp)
				values = append(values, point[1])
			}
		}

		// Add fields to frame - use metric name for field, scope is in frame name
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, times),
			data.NewField(seriesName, nil, values),
		)

		// Add the frame to the response
		response.Frames = append(response.Frames, frame)
	}

	if len(response.Frames) == 0 {
		log.DefaultLogger.Warn("No series returned from Datadog", "query", qm.Query)
	}

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	// Validate configuration
	if d.apiKey == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "DD-API-KEY is not configured",
		}, nil
	}

	if d.applicationKey == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "DD-APPLICATION-KEY is not configured",
		}, nil
	}

	// Try a simple API call to validate credentials
	baseURL := d.apiURL
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}
	testURL := baseURL + "api/v1/validate"

	httpReq, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Error creating request: %v", err),
		}, nil
	}

	httpReq.Header.Set("DD-API-KEY", d.apiKey)
	httpReq.Header.Set("DD-APPLICATION-KEY", d.applicationKey)

	client := &http.Client{Timeout: 10 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Error connecting to Datadog: %v", err),
		}, nil
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusOK {
		message = "Successfully connected to Datadog API"
	} else if httpResp.StatusCode == http.StatusForbidden || httpResp.StatusCode == http.StatusUnauthorized {
		status = backend.HealthStatusError
		message = "Invalid Datadog API credentials"
	} else {
		status = backend.HealthStatusError
		message = fmt.Sprintf("Datadog API returned status %d", httpResp.StatusCode)
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
