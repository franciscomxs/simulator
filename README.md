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

Compute a loan amortization schedule, optionally with a grace period.

**Request**

```json
{
  "amount": 10000,
  "rate": 0.02,
  "term": 12,
  "system": "PRICE",
  "grace_period": 3
}
```

| Field          | Type    | Required | Description                                                                                                  |
|----------------|---------|----------|--------------------------------------------------------------------------------------------------------------|
| `amount`       | decimal | yes      | Loan principal (> 0)                                                                                         |
| `rate`         | decimal | yes      | Monthly interest rate (0 to < 1)                                                                             |
| `term`         | integer | yes      | Number of paying installments (1 to 600)                                                                     |
| `system`       | string  | yes      | Amortization system: `PRICE`, `SAC`                                                                          |
| `grace_period` | integer | no       | Months before amortization begins (default `0`, must be `>= 0`, `grace_period + term <= 600`)                |

During the grace period no payment is charged; interest is capitalized into the balance at the contracted `rate`. After the grace period the PMT is computed from the capitalized balance (`adjusted_amount = amount * (1 + rate)^grace_period`) over the original `term`.

**Response — 200**

```json
{
  "amount": 10000,
  "rate": 0.02,
  "term": 12,
  "system": "PRICE",
  "grace_period": 3,
  "adjusted_amount": 10612.08,
  "total_duration": 15,
  "total_amount": 12041.99,
  "installments": [
    { "number": 1, "type": "GRACE",   "payment": 0,       "principal": 0,      "interest": 200.00, "balance": 10200.00 },
    { "number": 2, "type": "GRACE",   "payment": 0,       "principal": 0,      "interest": 204.00, "balance": 10404.00 },
    { "number": 3, "type": "GRACE",   "payment": 0,       "principal": 0,      "interest": 208.08, "balance": 10612.08 },
    { "number": 4, "type": "PAYMENT", "payment": 1003.47, "principal": 791.23, "interest": 212.24, "balance":  9820.85 }
  ]
}
```

| Field             | Description                                                       |
|-------------------|-------------------------------------------------------------------|
| `grace_period`    | Grace months echoed from the request (`0` when omitted)           |
| `adjusted_amount` | Balance after grace capitalization; equals `amount` when grace=`0`|
| `total_duration`  | Total operation length in months (`grace_period + term`)          |
| `installments[].type`    | `GRACE` or `PAYMENT`                                       |
| `installments[].balance` | Outstanding balance after this installment                 |

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
