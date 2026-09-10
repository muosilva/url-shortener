# URL Shortener

URL shortener built with Go, PostgreSQL, Redis, and a local observability stack.

## Stack

- Go API
- PostgreSQL for persistence
- Redis for cache
- Prometheus for metrics
- Loki for logs
- Grafana Alloy for collection
- Grafana for dashboards
- Air for live reload during local development

## Requirements

- Docker
- Docker Compose
- Make

## Run

Start the full development environment:

```bash
make dev
```

This starts the API with Air inside Docker, plus PostgreSQL, Redis, Prometheus, Loki, Alloy, and Grafana.

Stop everything:

```bash
make down
```

## Local URLs

- API: http://localhost:8080
- Health check: http://localhost:8080/healthz
- Metrics: http://localhost:8080/metrics
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100
- Alloy UI: http://localhost:12345

Grafana login:

```text
admin / admin
```

## Example Request

```bash
curl -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

## Useful Commands

```bash
make test
make logs
make ps
make compose-config
```
