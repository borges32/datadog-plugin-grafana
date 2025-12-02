#!/bin/bash
# Script de teste do simulador da API do Datadog

API_KEY="f22e8e0c4fcab646939943357ca7c201"
APP_KEY="5469722d1b56bc1e652698267eb979c10b7f7216"
BASE_URL="http://localhost:5000"

echo "=================================================="
echo "Testes do Simulador da API do Datadog"
echo "=================================================="
echo ""

echo "1. Health Check (sem autenticação)"
echo "--------------------------------------------------"
curl -s "${BASE_URL}/health" | jq '.'
echo ""
echo ""

echo "2. Validar Credenciais"
echo "--------------------------------------------------"
curl -s "${BASE_URL}/api/v1/validate" \
  --header "DD-API-KEY: ${API_KEY}" \
  --header "DD-APPLICATION-KEY: ${APP_KEY}" | jq '.'
echo ""
echo ""

echo "3. Testar com credenciais inválidas"
echo "--------------------------------------------------"
curl -s "${BASE_URL}/api/v1/validate" \
  --header "DD-API-KEY: invalid" \
  --header "DD-APPLICATION-KEY: invalid" | jq '.'
echo ""
echo ""

echo "4. Query de Métricas - Single Series"
echo "--------------------------------------------------"
FROM=$(date -d '1 hour ago' +%s)
TO=$(date +%s)

curl -s "${BASE_URL}/api/v1/query?from=${FROM}&to=${TO}&query=avg:processor.time{host:AH-CW-AP-104}" \
  --header "DD-API-KEY: ${API_KEY}" \
  --header "DD-APPLICATION-KEY: ${APP_KEY}" | jq '.series | length'
echo "Séries retornadas: 1"
echo ""
echo ""

echo "5. Query de Métricas - Multiple Series (com group by)"
echo "--------------------------------------------------"
curl -s "${BASE_URL}/api/v1/query?from=${FROM}&to=${TO}&query=avg:processor.time{host:AH-CW-AP-104}%20by%20{host,instance}" \
  --header "DD-API-KEY: ${API_KEY}" \
  --header "DD-APPLICATION-KEY: ${APP_KEY}" | jq '{
    status: .status,
    query: .query,
    series_count: (.series | length),
    first_series: .series[0] | {
      scope: .scope,
      metric: .metric,
      points: (.pointlist | length)
    }
  }'
echo ""
echo ""

echo "6. Query de Métricas - Resposta Completa (primeira série)"
echo "--------------------------------------------------"
curl -s "${BASE_URL}/api/v1/query?from=${FROM}&to=${TO}&query=avg:processor.time{host:AH-CW-AP-104}%20by%20{host,instance}" \
  --header "DD-API-KEY: ${API_KEY}" \
  --header "DD-APPLICATION-KEY: ${APP_KEY}" | jq '.series[0] | {
    metric: .metric,
    scope: .scope,
    tag_set: .tag_set,
    interval: .interval,
    sample_points: .pointlist[0:3]
  }'
echo ""
echo ""

echo "=================================================="
echo "Testes concluídos!"
echo "=================================================="
