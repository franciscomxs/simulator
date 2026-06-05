# CLAUDE.md — Simulator Project

## Language: Go

This project is written in **Go**. All code, tooling, and conventions follow Go idioms and the standard Go toolchain (`go build`, `go test`, `go mod`).

## Infrastructure: Docker Only

**All project operations must use Docker.** Never run `go`, databases, or any tooling directly on the host machine.

- Run the application: `docker compose up`
- Run all tests: `docker compose run --rm app go test ./...`
- Run a specific package: `docker compose run --rm app go test ./internal/domain/...`
- Add a dependency: `docker compose run --rm app go get <module>`
- Tidy modules: `docker compose run --rm app go mod tidy`
- Build: `docker compose run --rm app go build ./...`
- Access a shell: `docker compose run --rm app sh`

Do not run `go`, `psql`, or any other tool directly on the host. If a required service (database, cache, queue) is not yet in `docker-compose.yml`, add it there — do not install it locally.

## Architecture: Ports and Adapters (Hexagonal)

All code must follow the Ports and Adapters architecture. The core domain must be completely isolated from infrastructure concerns.

### Package Structure

```
cmd/
├── lambda/
│   └── main.go
└── web/
    └── main.go

internal/
├── domain/
│   ├── order.go                   # Entity
│   ├── payment.go                 # Entity
│   ├── pricing_service.go         # Domain service
│   ├── errors.go
│   └── value_objects/
│       ├── order_id.go
│       └── money.go
│
├── application/
│   ├── ports/
│   │   ├── order_repository.go      # Output port
│   │   ├── payment_gateway.go       # Output port
│   │   ├── event_publisher.go       # Output port
│   │   └── create_order_usecase.go  # Input port
│   │
│   └── usecases/
│       ├── create_order.go
│       ├── cancel_order.go
│       └── pay_order.go
│
├── adapters/
│   ├── inbound/
│   │   ├── http/
│   │   │   ├── dto/
│   │   │   │   ├── create_order_request.go
│   │   │   │   └── create_order_response.go
│   │   │   ├── order_handler.go
│   │   │   └── router.go
│   │   │
│   │   ├── grpc/
│   │   │   └── order_service.go
│   │   │
│   │   └── consumer/
│   │       └── order_created.go
│   │
│   └── outbound/
│       ├── postgres/
│       │   ├── order_repository.go
│       │   └── models/
│       │       └── order_record.go
│       │
│       ├── memory/
│       │   └── order_repository.go
│       │
│       ├── stripe/
│       │   └── payment_gateway.go
│       │
│       └── kafka/
│           └── event_publisher.go
│
└── infrastructure/
    ├── config/
    │   └── config.go
    │
    ├── database/
    │   └── postgres.go
    │
    └── container/
        └── dependencies.go

test/
├── integration/
│   ├── postgres_order_repository_test.go
│   └── stripe_payment_gateway_test.go
│
└── e2e/
    └── create_order_test.go
```

Fluxo

```
HTTP Request
      │
      ▼
Inbound Adapter
      │
      ▼
Use Case
      │
      ▼
Output Port
      │
      ▼
Outbound Adapter
      │
      ▼
Postgres / RabbitMQ
```

### Rules

1. **Domain has zero external dependencies.** No database drivers, HTTP clients, frameworks, or infrastructure concerns inside `internal/domain/`. The domain layer may import only the Go standard library and other domain packages. The application layer may import the Go standard library, domain packages, and interfaces defined in `application/ports/`.
2. **Depend on abstractions.** Use cases depend on the interfaces defined in `application/ports/`, never on concrete adapter types.
3. **Adapters implement ports.** Each outbound adapter implements one or more interfaces defined in `application/ports/`. Go's implicit interface satisfaction applies — no `implements` keyword needed, but the contract must be explicit. Implementing multiple interfaces is acceptable when they represent different facets of the same resource (e.g. `OrderReader` and `OrderWriter` both implemented by `PostgresOrderRepository`); it is not acceptable when they group unrelated responsibilities — in that case, create separate adapters.
4. **Inbound adapters call use cases.** HTTP handlers and consumers translate HTTP/message input into use-case calls and translate results back into HTTP/message output. No business logic inside adapters.
5. **`main.go` wires everything.** The composition root is the only place that instantiates concrete adapters and injects them. No `init()` side effects, no global state.
6. **Tests target layer boundaries.** Unit tests cover domain and use cases using in-memory fakes. Integration tests cover adapters against real infrastructure running in Docker.
7. **Keep domain models isolated from transport concerns.** Domain entities and value objects must never be exposed directly through HTTP, gRPC, messaging, database, or third-party APIs. Adapters are responsible for translating between transport-specific DTOs and domain models. Changes to external contracts must not require changes to the domain layer.

### Naming Convention

| Concept | Example name |
|---|---|
| Input port (interface) | `CreateOrderUseCase` |
| Output port (interface) | `OrderRepository`, `PaymentGateway` |
| Use case struct | `createOrderUseCase` (unexported) |
| Inbound adapter | `OrderHTTPHandler`, `OrderConsumer` |
| Outbound adapter | `PostgresOrderRepository`, `StripePaymentGateway` |
| Domain entity | `Order`, `Payment` |
| Domain service | `PricingService` |

## Testing

### Non-negotiable rules

1. **Tests are written before implementation (TDD).** For every new behavior, write a failing test first, then write only the code needed to make it pass. Implementation must never precede its test.
2. **Let the failing test drive the design.** Before writing any production code, run `go test ./...` inside Docker and confirm the new test fails for the right reason. A test that passes without implementation is not testing anything.
3. **Never delete tests.** Removing a test is forbidden unless the code it covers has been deleted entirely. If a test is failing because the behavior changed, update the test to match the new intended behavior — do not delete it.
4. **All tests must pass before any commit or merge.** Fix failing tests before moving on to anything else.
5. **Never skip tests.** Do not call `t.Skip()` to silence a failing test. If a test cannot pass yet, fix the underlying code or the test.
6. **Never comment out tests.** A commented-out test is a deleted test.

### TDD cycle (Red → Green → Refactor)

Follow this cycle for every unit of behavior, without exception:

1. **Red** — write the smallest possible `Test*` function describing the next desired behavior. Run `docker compose run --rm app go test ./...` and confirm it fails. If it does not fail, the test is wrong.
2. **Green** — write the minimum production code to make the failing test pass. No more. Do not generalize or anticipate future cases.
3. **Refactor** — with tests green, improve structure (naming, duplication, abstractions) without changing behavior. Run the suite after every change.

Repeat. Each cycle should take minutes. If a cycle is taking too long, the step is too large — split it.

### What to test and how

| Layer | Test type | How |
|---|---|---|
| `internal/domain/` | Unit | Plain `go test`, no fakes needed — pure functions and value objects |
| `internal/application/usecases/` | Unit | Inject in-memory fakes that implement the port interfaces |
| `internal/adapters/outbound/` | Integration | Real infrastructure via Docker; use `TestMain` to set up/tear down |
| `internal/adapters/inbound/` | Integration | `httptest.NewRecorder` + real use case wired in |
| End-to-end | E2E | Full stack via `docker compose`; exercise the HTTP API end to end |

### Test as design tool

Tests written before implementation expose design problems early:

- If a test is hard to set up, the struct under test has too many dependencies — split it.
- If the test needs to reach into unexported fields, the abstraction boundary is wrong — redesign the interface.
- If you cannot test a use case without a real database, the layer is not properly isolated — fix the architecture first.

Never suppress design feedback by making tests more complex. A painful test is a signal, not a problem to work around.

### In-memory adapters

For every output port interface, maintain an in-memory implementation under adapters/memory.

```
internal/
├── application/
│   └── ports/
│       └── order_repository.go      # OrderRepository interface (port)
├── adapters/
│   ├── postgres/
│   │   └── order_repository.go      # PostgreSQL implementation
│   └── memory/
│       └── order_repository.go      # In-memory implementation (tests/dev)
```

These adapters implement the full interface using in-memory maps or slices. They serve as the canonical test double for unit tests.

### Test file structure

Co-locate unit tests with the package they test (`order_test.go` next to `order.go`). Keep integration and E2E tests under a top-level `test/` directory:

```
internal/
  domain/
    order.go
    order_test.go          # unit test, same package (white-box) or _test package (black-box)
  application/
    usecases/
      create_order.go
      create_order_test.go # unit test using fakes

test/
  integration/
    postgres_order_repository_test.go
  e2e/
    create_order_test.go
```

### Running tests (all inside Docker)

```bash
# All tests
docker compose run --rm app go test ./...

# Unit tests only (no Docker services needed beyond the app container)
docker compose run --rm app go test ./internal/...

# Integration tests
docker compose run --rm app go test ./test/integration/...

# Single package
docker compose run --rm app go test ./internal/domain/

# With coverage
docker compose run --rm app go test -coverprofile=coverage.out ./...
docker compose run --rm app go tool cover -func=coverage.out
```

### Coverage

- `internal/domain/` must maintain 100% statement coverage — every business rule must be exercised by a test.
- `internal/application/usecases/` should maintain high coverage, focusing on business behavior rather than coverage metrics.
- Adapters must have integration tests covering the happy path and primary error cases.
- Do not use `//nolint` or coverage exclusion tricks to paper over untested code. If a branch is genuinely unreachable, remove it.
