# Datadog Plugin for Grafana

A Grafana datasource plugin that integrates with the Datadog Metrics API, allowing you to query and visualize Datadog metrics directly in Grafana dashboards.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Repository Structure](#repository-structure)
- [Quick Start](#quick-start)
- [Plugin Installation](#plugin-installation)
- [Configuration](#configuration)
- [Development](#development)
- [API Simulator](#api-simulator)
- [Testing](#testing)

## 🎯 Overview

This project consists of two main components:

1. **Grafana Plugin** - A backend datasource plugin for Grafana that queries Datadog metrics
2. **API Simulator** - A Python-based Datadog API simulator for local development and testing

## ✨ Features

- ✅ Query Datadog metrics directly from Grafana
- ✅ Support for Datadog query syntax (e.g., `avg:metric{tag} by {dimension}`)
- ✅ Automatic timeframe conversion from Grafana dashboards
- ✅ Multiple series support with proper labeling
- ✅ Alerting support
- ✅ Secure credential storage (API keys encrypted in backend)
- ✅ Health check for credential validation
- ✅ Cross-platform backend executables (Linux, macOS, Windows)

## 📁 Repository Structure

```
datadog-plugin-grafana/
├── plugin/
│   └── opensource-datadog-plugin-grafana-datasource/
│       ├── src/                      # Frontend TypeScript/React code
│       │   ├── components/
│       │   │   ├── ConfigEditor.tsx  # Datasource configuration UI
│       │   │   └── QueryEditor.tsx   # Query builder UI
│       │   ├── types.ts              # TypeScript type definitions
│       │   └── datasource.ts         # Frontend datasource logic
│       ├── pkg/                      # Backend Go code
│       │   ├── main.go              # Plugin entry point
│       │   └── plugin/
│       │       └── datasource.go    # Datadog API integration
│       ├── dist/                     # Compiled plugin (generated)
│       ├── package.json             # NPM dependencies and scripts
│       ├── go.mod                   # Go dependencies
│       └── Magefile.go              # Build configuration
│
└── simulator-metrics-datadog/       # Datadog API simulator
    ├── app.py                       # Flask application
    ├── requirements.txt             # Python dependencies
    ├── test.sh                      # Test script
    └── README.md                    # Simulator documentation
```

## 🚀 Quick Start

### Prerequisites

- **Node.js** >= 18
- **Go** >= 1.19
- **Python** >= 3.8 (for simulator)
- **Docker** (optional, for Grafana instance)
- **Mage** (Go build tool)

### Install Mage

```bash
go install github.com/magefile/mage@latest
```

### Build the Plugin

```bash
cd plugin/opensource-datadog-plugin-grafana-datasource
npm install
npm run build
```

This will:
1. Compile TypeScript/React frontend → `dist/module.js`
2. Compile Go backend for all platforms → `dist/gpx_datadog_plugin_grafana_*`

### Run Grafana with Docker

```bash
cd plugin/opensource-datadog-plugin-grafana-datasource
docker compose up -d
```

Access Grafana at `http://localhost:3000` (admin/admin)

## 🔧 Plugin Installation

### Option 1: Local Development (Docker Compose)

The plugin is automatically mounted when using the provided `docker-compose.yaml`.

### Option 2: Manual Installation

1. Build the plugin (see Quick Start)
2. Copy the `dist/` folder to your Grafana plugins directory:
   ```bash
   cp -r dist /var/lib/grafana/plugins/opensource-datadog-datasource
   ```
3. Restart Grafana
4. Enable the plugin in Grafana settings (if required)

## ⚙️ Configuration

### Add Datasource in Grafana

1. Go to **Configuration** → **Data Sources** → **Add data source**
2. Search for "Datadog Plugin Grafana"
3. Configure the following fields:

#### URL
The Datadog API endpoint for your region:
- US1: `https://api.datadoghq.com/`
- US3: `https://api.us3.datadoghq.com/`
- US5: `https://api.us5.datadoghq.com/`
- EU: `https://api.datadoghq.eu/`
- Local simulator: `http://localhost:5000/`

#### DD-API-KEY
Your Datadog API key (stored securely, encrypted in backend)

#### DD-APPLICATION-KEY
Your Datadog Application key (stored securely, encrypted in backend)

4. Click **Save & Test** to verify the connection

### Query Syntax

In the Query Editor, use standard Datadog query syntax:

```
avg:processor.time{host:myhost} by {host,instance}
```

**Examples:**

```
# Average CPU usage for a specific host
avg:system.cpu.user{host:web-server-01}

# Memory usage grouped by host
avg:system.mem.used{env:production} by {host}

# Network traffic with aggregation
sum:system.net.bytes_rcvd{*} by {host}
```

The time range is automatically taken from the Grafana dashboard.

## 💻 Development

### Frontend Development (Hot Reload)

```bash
cd plugin/opensource-datadog-plugin-grafana-datasource
npm run dev
```

This watches for TypeScript changes and rebuilds automatically.

### Backend Development

After modifying Go code:

```bash
npm run build:backend
```

Or build for specific platform:

```bash
mage -v build:linux    # Linux only
```

### Project Scripts

```bash
npm run build           # Build frontend + backend
npm run build:frontend  # Build only TypeScript/React
npm run build:backend   # Build only Go backend
npm run dev            # Watch mode for frontend
npm run lint           # Run ESLint
npm run test           # Run Jest tests
```

## 🧪 API Simulator

For local development without a Datadog account, use the included API simulator.

### Start the Simulator

```bash
cd simulator-metrics-datadog
pip install -r requirements.txt
python app.py
```

The simulator runs on `http://localhost:5000`

### Test Credentials

- **DD-API-KEY**: `f22e8e0c4fcab646939943357ca7c201`
- **DD-APPLICATION-KEY**: `5469722d1b56bc1e652698267eb979c10b7f7216`

### Run Tests

```bash
cd simulator-metrics-datadog
./test.sh
```

### Features

- ✅ Simulates `/api/v1/query` endpoint
- ✅ Authentication validation
- ✅ Generates realistic metric data with patterns
- ✅ Supports `group by` for multiple series
- ✅ Compatible response format

See [simulator-metrics-datadog/README.md](simulator-metrics-datadog/README.md) for more details.

## 🧪 Testing

### Test with Real Datadog API

Configure the plugin with your actual Datadog credentials and query real metrics.

### Test with Simulator

1. Start the simulator (see above)
2. Configure plugin with:
   - URL: `http://localhost:5000/`
   - Use test credentials above
3. Create a panel and query: `avg:processor.time{host:AH-CW-AP-104} by {host,instance}`

### Expected Results

You should see 4 series in the graph:
- `host:AH-CW-AP-104,instance:0`
- `host:AH-CW-AP-104,instance:1`
- `host:AH-CW-AP-104,instance:2`
- `host:AH-CW-AP-104,instance:3`

## 📊 Alerting

The plugin supports Grafana alerting. You can:

1. Create alert rules using Datadog metrics
2. Set thresholds and conditions
3. Configure notification channels

The backend evaluates queries and provides data for alert evaluation.

## 🔨 Build Outputs

After running `npm run build`, the following executables are generated in `dist/`:

- `gpx_datadog_plugin_grafana_linux_amd64` (Linux 64-bit)
- `gpx_datadog_plugin_grafana_linux_arm64` (Linux ARM64)
- `gpx_datadog_plugin_grafana_linux_arm` (Linux ARM)
- `gpx_datadog_plugin_grafana_darwin_amd64` (macOS Intel)
- `gpx_datadog_plugin_grafana_darwin_arm64` (macOS Apple Silicon)
- `gpx_datadog_plugin_grafana_windows_amd64.exe` (Windows 64-bit)

## 📝 License

Apache-2.0

## 👤 Author

Opensource

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📚 Additional Resources

- [Grafana Plugin Documentation](https://grafana.com/docs/grafana/latest/developers/plugins/)
- [Datadog API Documentation](https://docs.datadoghq.com/api/latest/metrics/)
- [Grafana Plugin SDK for Go](https://github.com/grafana/grafana-plugin-sdk-go)
