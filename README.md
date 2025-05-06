# go-location-finder

Go microservice for ZIP code and weather lookup, instrumented with OpenTelemetry and exporting traces to Zipkin.

---

## Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [How to Run](#how-to-run)
- [How to Test](#how-to-test)
- [Observability](#observability)
- [Environment Variables](#environment-variables)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

---

## Architecture

- **go-location-finder**: Looks up ZIP code and weather.
- **go-location-validator**: Validates ZIP code and calls the finder.
- **otel-collector**: OpenTelemetry Collector for traces.
- **zipkin-all-in-one**: UI for distributed traces.

---

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## How to Run

1. **Clone the repository:**

   ```sh
   git clone https://github.com/your-user/go-location-finder.git
   cd go-location-finder
   ```

2. **Set up environment variables:**

   - Create a `.env` file inside the `go-location-finder` folder:
     ```
     API_KEY=your_weatherapi_key
     ```

3. **Start the services:**

   ```sh
   docker-compose up --build
   ```

   Wait until all containers are "Up".

---

## How to Test

### Using the validator service

```sh
curl -X POST http://localhost:9000/validate \
  -H "Content-Type: application/json" \
  -d '{"cep":"01001000"}'
```

### Using the finder service directly

```sh
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01001000"}'
```

---

## Observability

- **Zipkin UI:** [http://localhost:9411](http://localhost:9411)

  1. Open the link above.
  2. Select the service name in "Service Name" (e.g., `go-location-finder` or `microservice-go-location-validator`).
  3. Click "Find Traces" to view the generated traces.

---

## Environment Variables

See examples in `docker-compose.yaml`:

- `OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317`
- `OTEL_SERVICE_NAME=go-location-finder` (or `go-location-validator`)
- `HTTP_PORT=:8080` or `:9000`
- `API_KEY=...` (your WeatherAPI key)

---

## FAQ

**1. I get a `connection refused` error for the OTEL Collector. What should I do?**  
Wait for the `otel-collector` service to be ready before starting the microservices, or restart the microservices after the collector is running.

**2. How do I make calls between microservices?**  
Use the service name as the hostname (e.g., `http://go-location-finder:8080/cep`) inside the containers.

**3. How do I view the traces?**  
Open Zipkin at [http://localhost:9411](http://localhost:9411) and search by service name.
