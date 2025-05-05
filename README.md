# go-location-finder

A microservices project in Go for querying ZIP code and weather information, fully instrumented for observability (OpenTelemetry, Jaeger, Prometheus, Zipkin).

---

## Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Running with Docker Compose](#running-with-docker-compose)
- [Testing the Microservices](#testing-the-microservices)
- [Observability](#observability)
- [Environment Variables](#environment-variables)
- [FAQ](#faq)

---

## Architecture

- **go-location-finder**: Service that queries ZIP code and weather information.
- **go-location-validator**: Service that validates ZIP codes and calls `go-location-finder`.
- **otel-collector**: OpenTelemetry Collector for traces and metrics.
- **jaeger-all-in-one**: UI for distributed tracing.
- **zipkin-all-in-one**: Alternative UI for traces.
- **prometheus**: Collects metrics from the services.

---

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## Running with Docker Compose

1. **Clone the repository:**

   ```sh
   git clone https://github.com/your-user/go-location-finder.git
   cd go-location-finder
   ```

2. **Start all services:**

   ```sh
   docker-compose up --build
   ```

   > Wait until all containers show status `Up`.

3. **Check running containers:**
   ```sh
   docker-compose ps
   ```

---

## Testing the Microservices

### 1. **From your host terminal:**

- **Validate a ZIP code:**

  ```sh
  curl -X POST http://localhost:9000/validate \
    -H "Content-Type: application/json" \
    -d '{"cep": "12222-050"}'
  ```

### 2. **Using `.http` file (VS Code REST Client):**

In your `request.http` file:

```http
POST http://localhost:9000/validate
Content-Type: application/json

{
  "cep": "12222-050"
}
```

---

## Observability

- **Jaeger UI:** [http://localhost:16686](http://localhost:16686)
- **Zipkin UI:** [http://localhost:9411](http://localhost:9411)
- **Prometheus:** [http://localhost:9090](http://localhost:9090)

---

## Environment Variables

See examples in `docker-compose.yaml`:

- `OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317`
- `OTEL_SERVICE_NAME=go-location-finder` (or `go-location-validator`)
- `HTTP_PORT=:8080` or `:9000`
- `API_KEY=...` (for weather API access)

---

## FAQ

**1. I get a `connection refused` error for the OTEL Collector. What should I do?**  
Wait for the `otel-collector` service to be ready before starting the microservices, or restart the microservices after the collector is running.

**2. How do I make calls between microservices?**  
Use the service name as the hostname (e.g., `http://go-location-finder:8080/cep`) inside the containers.

**3. How do I view traces?**  
Open Jaeger UI at [http://localhost:16686](http://localhost:16686).

---

## Contributing

Pull requests are welcome! Please open an issue to discuss any changes.

---

## License

MIT
