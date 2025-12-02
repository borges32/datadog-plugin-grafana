# Integração com Datadog - Plugin Grafana

## Visão Geral

Este plugin permite integração do Grafana com a API de Métricas do Datadog.

## Configuração

### 1. Instalação do Mage

O build do backend Go requer o Mage. Para instalar:

```bash
go install github.com/magefile/mage@latest
```

### 2. Build do Plugin

Para compilar o plugin completo (frontend TypeScript + backend Go):

```bash
npm run build
```

Isso executará:
- `npm run build:frontend` - Compila o código TypeScript com webpack
- `npm run build:backend` - Compila os executáveis Go para múltiplas plataformas

### 3. Configuração no Grafana

Após instalar o plugin no Grafana, configure os seguintes campos:

#### URL
- URL da API do Datadog (exemplo: `https://api.us3.datadoghq.com/`)
- Varie de acordo com a região do seu Datadog (us3, us5, eu, etc.)

#### DD-API-KEY
- Chave de API do Datadog
- Campo seguro (criptografado no backend)

#### DD-APPLICATION-KEY
- Chave de Aplicação do Datadog
- Campo seguro (criptografado no backend)

## Uso

### Query Editor

No painel de query do Grafana, você pode inserir queries do Datadog:

**Exemplo de Query:**
```
avg:processor.time{host:AH-CW-AP-104} by {host,instance}
```

### Timeframe

O plugin automaticamente usa o timeframe selecionado no dashboard do Grafana. Os valores de `from` e `to` são convertidos para Unix timestamps (Epoch) e enviados para a API do Datadog.

## Estrutura de Retorno

O plugin processa a resposta da API do Datadog e converte para o formato de Data Frames do Grafana:

- **Series**: Cada série retornada pelo Datadog vira um frame
- **Pointlist**: Array de pontos onde [0] é o timestamp (ms) e [1] é o valor
- **Tag Set**: Labels extraídos e aplicados às métricas
- **Scope**: Escopo da métrica aplicado

## Arquivos Modificados

### Frontend (TypeScript)
- `src/types.ts` - Tipos atualizados para Datadog
- `src/components/ConfigEditor.tsx` - UI de configuração
- `src/components/QueryEditor.tsx` - UI de query

### Backend (Go)
- `pkg/plugin/datasource.go` - Implementação da integração com Datadog
- `pkg/main.go` - Entrada do plugin

### Build
- `package.json` - Scripts de build atualizados

## Exemplo de Resposta da API Datadog

```json
{
  "status": "ok",
  "series": [
    {
      "metric": "processor.time",
      "display_name": "processor.time",
      "tag_set": ["host:AH-CW-AP-104", "instance:3"],
      "pointlist": [
        [1764658800000.0, 8.64499773957933],
        [1764658840000.0, 6.66704217626698]
      ]
    }
  ]
}
```

## Health Check

O plugin implementa health check que valida:
- Configuração das credenciais (API Key e Application Key)
- Conectividade com a API do Datadog
- Validade das credenciais

## Executáveis Gerados

O build gera executáveis para múltiplas plataformas:
- `gpx_datadog_plugin_grafana_linux_amd64`
- `gpx_datadog_plugin_grafana_linux_arm64`
- `gpx_datadog_plugin_grafana_linux_arm`
- `gpx_datadog_plugin_grafana_darwin_amd64`
- `gpx_datadog_plugin_grafana_darwin_arm64`
- `gpx_datadog_plugin_grafana_windows_amd64.exe`
