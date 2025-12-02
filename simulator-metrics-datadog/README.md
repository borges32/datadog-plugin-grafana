# Datadog API Simulator

This is a Datadog metrics API simulator for local testing.

## Installation

```bash
pip install -r requirements.txt
```

## Running

```bash
python app.py
```

The server will run on `http://localhost:5000`

## Credentials

- **DD-API-KEY**: `f22e8e0c4fcab646939943357ca7c201`
- **DD-APPLICATION-KEY**: `5469722d1b56bc1e652698267eb979c10b7f7216`

## Endpoints

### Health Check
```bash
curl http://localhost:5000/health
```

### Validate Credentials
```bash
curl http://localhost:5000/api/v1/validate \
  --header 'DD-API-KEY: f22e8e0c4fcab646939943357ca7c201' \
  --header 'DD-APPLICATION-KEY: 5469722d1b56bc1e652698267eb979c10b7f7216'
```

### Metrics Query
```bash
curl 'http://localhost:5000/api/v1/query?from=1764658800&to=1764662400&query=avg:processor.time{host:AH-CW-AP-104} by {host,instance}' \
  --header 'DD-API-KEY: f22e8e0c4fcab646939943357ca7c201' \
  --header 'DD-APPLICATION-KEY: 5469722d1b56bc1e652698267eb979c10b7f7216'
```

## Features

- ✅ Authentication with DD-API-KEY and DD-APPLICATION-KEY
- ✅ `/api/v1/validate` endpoint for health check
- ✅ `/api/v1/query` endpoint that simulates Datadog queries
- ✅ Generates simulated data with sinusoidal variation + noise
- ✅ Supports group by (e.g., `by {host,instance}`)
- ✅ Returns multiple series when group by is used
- ✅ Response format identical to Datadog

## Grafana Plugin Configuration

In the plugin ConfigEditor, use:
- **URL**: `http://localhost:5000/`
- **DD-API-KEY**: `f22e8e0c4fcab646939943357ca7c201`
- **DD-APPLICATION-KEY**: `5469722d1b56bc1e652698267eb979c10b7f7216`

## Simulated Data

The simulator generates data with:
- Base values between 2-10
- Sinusoidal variation to simulate patterns
- Random noise
- Occasional spikes (5% chance) between 20-50
- 20-second interval between points
