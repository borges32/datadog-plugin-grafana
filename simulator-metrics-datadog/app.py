#!/usr/bin/env python3
"""
Simulador da API de Métricas do Datadog
Simula o endpoint /api/v1/query com autenticação
"""

from flask import Flask, request, jsonify
import time
import random
import math
from urllib.parse import unquote

app = Flask(__name__)

# Credenciais fake para autenticação
VALID_API_KEY = "f22e8e0c4fcab646939943357ca7c201"
VALID_APP_KEY = "5469722d1b56bc1e652698267eb979c10b7f7216"

def check_auth():
    """Verifica se as credenciais estão corretas"""
    api_key = request.headers.get('DD-API-KEY')
    app_key = request.headers.get('DD-APPLICATION-KEY')
    
    if not api_key or not app_key:
        return False, {"errors": ["API key and application key are required"]}
    
    if api_key != VALID_API_KEY:
        return False, {"errors": ["Invalid API key"]}
    
    if app_key != VALID_APP_KEY:
        return False, {"errors": ["Invalid application key"]}
    
    return True, None

def generate_pointlist(from_ts, to_ts, interval=20):
    """Gera dados simulados de métrica"""
    points = []
    current = from_ts * 1000  # Converter para milliseconds
    to_ms = to_ts * 1000
    
    base_value = random.uniform(2, 10)
    
    while current <= to_ms:
        # Gera valor com variação senoidal + ruído
        t = (current - from_ts * 1000) / 1000
        value = base_value + math.sin(t / 100) * 3 + random.uniform(-2, 2)
        
        # Adiciona picos aleatórios
        if random.random() < 0.05:
            value += random.uniform(20, 50)
        
        points.append([float(current), max(0, value)])
        current += interval * 1000  # Intervalo em milliseconds
    
    return points

def parse_query(query_str):
    """Parse básico da query do Datadog"""
    # Exemplo: avg:processor.time{host:AH-CW-AP-104} by {host,instance}
    parts = query_str.split('{')
    
    if len(parts) < 2:
        return None, None, []
    
    # Extrair métrica (ex: avg:processor.time)
    metric_part = parts[0].strip()
    if ':' in metric_part:
        aggregation, metric = metric_part.split(':', 1)
    else:
        aggregation = 'avg'
        metric = metric_part
    
    # Extrair tags (ex: host:AH-CW-AP-104)
    tags_part = parts[1].split('}')[0]
    tags = []
    for tag in tags_part.split(','):
        tag = tag.strip()
        if ':' in tag:
            tags.append(tag)
    
    # Extrair group by
    group_by = []
    if 'by' in query_str and '{' in query_str.split('by')[1]:
        group_by_part = query_str.split('by')[1].split('{')[1].split('}')[0]
        group_by = [g.strip() for g in group_by_part.split(',')]
    
    return metric, aggregation, tags, group_by

@app.route('/api/v1/validate', methods=['GET'])
def validate():
    """Endpoint para health check"""
    is_valid, error = check_auth()
    
    if not is_valid:
        return jsonify(error), 403
    
    return jsonify({"valid": True}), 200

@app.route('/api/v1/query', methods=['GET'])
def query_metrics():
    """Endpoint principal de query de métricas"""
    
    # Verificar autenticação
    is_valid, error = check_auth()
    if not is_valid:
        return jsonify(error), 403
    
    # Obter parâmetros
    from_ts = request.args.get('from', type=int)
    to_ts = request.args.get('to', type=int)
    query = request.args.get('query', '')
    
    if not from_ts or not to_ts or not query:
        return jsonify({
            "errors": ["Missing required parameters: from, to, query"]
        }), 400
    
    # Decodificar query se necessário
    query = unquote(query)
    
    # Parse da query
    metric, aggregation, tags, group_by = parse_query(query)
    
    if not metric:
        return jsonify({
            "errors": ["Invalid query format"]
        }), 400
    
    # Gerar séries baseadas no group by
    series_list = []
    
    # Se tem group by, gera múltiplas séries
    if 'instance' in group_by or 'instance' in str(group_by):
        instances = ['0', '1', '2', '3']
    else:
        instances = ['0']
    
    for instance in instances:
        # Construir tag_set
        tag_set = []
        scope_parts = []
        
        for tag in tags:
            if ':' in tag:
                key, value = tag.split(':', 1)
                tag_set.append(f"{key}:{value}")
                scope_parts.append(f"{key}:{value}")
        
        # Adicionar instance ao tag_set se estiver no group by
        if len(instances) > 1:
            tag_set.append(f"instance:{instance}")
            scope_parts.append(f"instance:{instance}")
        
        scope = ','.join(scope_parts)
        
        # Gerar pointlist
        pointlist = generate_pointlist(from_ts, to_ts)
        
        series = {
            "aggr": aggregation,
            "attributes": {},
            "display_name": metric,
            "end": to_ts * 1000,
            "expression": f"{aggregation}:{metric}{{{','.join(tag_set)}}}",
            "interval": 20,
            "length": len(pointlist),
            "metric": metric,
            "pointlist": pointlist,
            "query_index": 0,
            "scope": scope,
            "start": from_ts * 1000,
            "tag_set": tag_set,
            "unit": None
        }
        
        series_list.append(series)
    
    response = {
        "status": "ok",
        "res_type": "time_series",
        "resp_version": 1,
        "query": query,
        "from_date": from_ts * 1000,
        "to_date": to_ts * 1000,
        "series": series_list,
        "values": [],
        "times": [],
        "message": "",
        "group_by": group_by
    }
    
    return jsonify(response), 200

@app.route('/health', methods=['GET'])
def health():
    """Health check sem autenticação"""
    return jsonify({"status": "healthy"}), 200

@app.errorhandler(404)
def not_found(error):
    return jsonify({"errors": ["Endpoint not found"]}), 404

@app.errorhandler(500)
def internal_error(error):
    return jsonify({"errors": ["Internal server error"]}), 500

if __name__ == '__main__':
    print("=" * 60)
    print("Simulador da API do Datadog")
    print("=" * 60)
    print(f"API Key: {VALID_API_KEY}")
    print(f"Application Key: {VALID_APP_KEY}")
    print("=" * 60)
    print("Endpoints disponíveis:")
    print("  GET /health - Health check")
    print("  GET /api/v1/validate - Validar credenciais")
    print("  GET /api/v1/query - Query de métricas")
    print("=" * 60)
    print("\nServidor rodando em http://0.0.0.0:5000")
    print("\nExemplo de uso:")
    print("curl 'http://localhost:5000/api/v1/query?from=1764658800&to=1764662400&query=avg:processor.time{host:AH-CW-AP-104} by {host,instance}' \\")
    print(f"  --header 'DD-API-KEY: {VALID_API_KEY}' \\")
    print(f"  --header 'DD-APPLICATION-KEY: {VALID_APP_KEY}'")
    print("=" * 60)
    
    app.run(host='0.0.0.0', port=5000, debug=True)
