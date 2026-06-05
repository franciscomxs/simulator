# Simulator

Financial simulation API in Go. Supports two simulation types:

- **Loan simulation** — full amortization schedule for PRICE (French) and SAC (Constant Amortization) systems.
- **Investment simulation** — compound growth projection with monthly contributions and a yearly timeline.

Runs as a traditional HTTP server or as an AWS Lambda + API Gateway function — both using the same domain, use cases, and adapters.

---

## Requirements

- Docker
- Docker Compose

No local Go installation required. All commands run inside containers.

---

## Running

```bash
# Build the image
docker compose build

# Start the web server (http://localhost:8080)
docker compose up
```

---

## API

### POST /loan/simulate

Compute a loan amortization schedule.

**Request**

```json
{
  "amount": 10000,
  "rate": 0.02,
  "term": 12,
  "system": "PRICE"
}
```

| Field    | Type    | Description                         |
|----------|---------|-------------------------------------|
| `amount` | decimal | Loan principal (> 0)                |
| `rate`   | decimal | Monthly interest rate (0 to < 1)    |
| `term`   | integer | Number of installments (1 to 600)   |
| `system` | string  | Amortization system: `PRICE`, `SAC` |

**Response — 200**

```json
{
  "amount": 10000,
  "rate": 0.02,
  "term": 12,
  "system": "PRICE",
  "total_amount": 11347.20,
  "installments": [
    { "number": 1, "payment": 945.60, "principal": 745.60, "interest": 200.00 },
    { "number": 2, "payment": 945.60, "principal": 760.51, "interest": 185.09 }
  ]
}
```

---

### POST /investment/simulate

Simulate compound investment growth with monthly contributions.

**Request**

```json
{
  "initial_amount": 0,
  "monthly_contribution": 1000,
  "rate": 0.01,
  "term": 10
}
```

| Field                  | Type    | Description                              |
|------------------------|---------|------------------------------------------|
| `initial_amount`       | decimal | Initial investment amount (>= 0)         |
| `monthly_contribution` | decimal | Monthly contribution amount (>= 0)       |
| `rate`                 | decimal | Monthly interest rate (>= 0)             |
| `term`                 | integer | Investment horizon in years (> 0)        |

Formula per month: `balance = balance × (1 + rate) + monthly_contribution`

**Response — 200**

```json
{
  "initial_amount": 0,
  "monthly_contribution": 1000,
  "rate": 0.01,
  "term": 10,
  "final_amount": 230038.69,
  "timeline": [
    { "year": 1, "contributions": 12000, "balance": 12682.50 },
    { "year": 2, "contributions": 24000, "balance": 26973.46 }
  ]
}
```

**Response — 400**

```json
{ "error": "invalid simulation parameters" }
```

---

### GET /health

Health check. Returns `{"status":"ok"}`.

### GET /docs

Swagger UI — interactive API documentation.

### GET /openapi.yaml

Raw OpenAPI 3.0 specification.

---

## Amortization Systems

### PRICE (French)

Fixed payment amount. Interest portion decreases and principal portion increases with each installment.

```
PMT = PV × i / (1 − (1 + i)^−n)
```

### SAC (Constant Amortization)

Constant principal amortization. Interest is calculated on the outstanding balance, so the payment amount decreases over time.

```
Principal amortization = PV / n
Interest              = OutstandingBalance × i
Payment               = PrincipalAmortization + Interest
```

---

## Testing

```bash
# Run all tests
docker compose run --rm app go test ./...

# Unit tests only
docker compose run --rm app go test ./internal/...

# Integration / E2E
docker compose run --rm app go test ./test/...

# With coverage
docker compose run --rm app go test -coverprofile=coverage.out ./...
docker compose run --rm app go tool cover -func=coverage.out
```

---

## Project Structure

```
cmd/
  web/          # HTTP server entry point
  lambda/       # AWS Lambda entry point

internal/
  domain/       # Business rules — pure Go, no external dependencies
  application/
    ports/      # Input and output port interfaces
    usecases/   # Use case implementations
  adapters/
    inbound/
      http/     # HTTP handler and router (package httphandler)
      lambda/   # Lambda handler (package lambdahandler)
      dto/      # Shared request/response types
    outbound/
      memory/   # In-memory adapter (demonstrates port pattern)
  infrastructure/
    config/     # Environment configuration
    container/  # Dependency wiring

test/
  e2e/          # End-to-end tests against the full HTTP stack
```

### Dependency rule

```
cmd → infrastructure → adapters → application → domain
```

The domain layer imports only the Go standard library. The application layer imports only the domain and its own port interfaces. Adapters and infrastructure are the only layers allowed to import external packages.

---

## Architecture

This project follows the **Ports and Adapters (Hexagonal)** architecture. See [AGENTS.md](AGENTS.md) for the full set of architectural and testing conventions.
