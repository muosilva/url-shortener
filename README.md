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

These credentials are for local development only.

## Example Request

```bash
curl -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

## Load Test

Run a small k6 load test:

```bash
make load-test
```

Default load:

```text
2 iterations/second for 1 minute
```

Each iteration creates one short URL and then requests the generated short URL without following the external redirect.

Increase load gradually:

```bash
make load-test K6_RATE=5 K6_DURATION=2m K6_MAX_VUS=12
```

## Useful Commands

```bash
make test
make load-test
make logs
make ps
make compose-config
```

## Português

Este projeto é um encurtador de URLs feito com Go, PostgreSQL, Redis e uma stack local de observabilidade.

Ele foi pensado para estudar:

- métricas com Prometheus;
- logs com Loki;
- coleta com Grafana Alloy;
- dashboards no Grafana;
- teste de carga com k6;
- análise de gargalos usando gráficos.

### Tecnologias

- API em Go
- PostgreSQL para persistência
- Redis para cache
- Prometheus para métricas
- Loki para logs
- Grafana Alloy para coleta
- Grafana para visualização
- Air para live reload no desenvolvimento local
- k6 para teste de carga

### Requisitos

- Docker
- Docker Compose
- Make

### Como rodar

Suba todo o ambiente de desenvolvimento com um único comando:

```bash
make dev
```

Esse comando sobe a API com Air dentro do Docker, além de PostgreSQL, Redis, Prometheus, Loki, Alloy e Grafana.

Para parar os containers:

```bash
make down
```

### URLs locais

- API: http://localhost:8080
- Health check: http://localhost:8080/healthz
- Métricas: http://localhost:8080/metrics
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100
- Alloy UI: http://localhost:12345

Login do Grafana:

```text
admin / admin
```

Essas credenciais são apenas para desenvolvimento local.

### Criar uma URL curta

```bash
curl -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

### Rodar teste de carga

Execute um teste pequeno e seguro com k6:

```bash
make load-test
```

Por padrão, o teste roda com:

```text
2 iterações por segundo durante 1 minuto
```

Cada iteração cria uma URL curta e depois acessa a URL gerada sem seguir o redirect externo.

Para aumentar a carga gradualmente:

```bash
make load-test K6_RATE=5 K6_DURATION=2m K6_MAX_VUS=12
```

O valor de `K6_RATE` controla quantas iterações por segundo o k6 tenta executar.

### Observabilidade

Depois de rodar `make dev`, acesse o Grafana em http://localhost:3000.

Use o Prometheus como fonte de dados para visualizar métricas como:

- requests por segundo;
- taxa de erro 5xx;
- duração das requests;
- p95 de criação de URL;
- duração de operações no PostgreSQL;
- duração de operações no Redis;
- cache hit ratio;
- goroutines;
- CPU e memória da API;
- requests em andamento.

Use o Loki como fonte de dados para visualizar os logs da aplicação.

### Comandos úteis

```bash
make test
make load-test
make logs
make ps
make compose-config
```
