# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Infrastructure
```bash
make dev_prep       # Start PostgreSQL, RabbitMQ, PgAdmin via Docker Compose
make dev_stop       # Stop infrastructure
make dev_clean      # Stop and remove infrastructure + volumes
make dbgui          # Open PgAdmin at localhost:8000 (password: xemmenu)
make mqgui          # Open RabbitMQ management at localhost:15672 (user/pass: xemmenu)
```

### Running Services
```bash
make run_order      # Run Order service (port 5000)
make run_kitchen    # Run Kitchen service (port 5001)
make run_payment    # Run Payment service (port 5002)

# Or from within a service directory:
cd services/order && make run
```

### Testing
```bash
# From within a service directory
go test ./...
go test -v ./...
go test -v ./internal/domain/...  # Run tests in a specific package
```

## Architecture

This is a **Saga Orchestration** demo: the Order service acts as the saga coordinator, while Kitchen and Payment services are participants. All inter-service communication is asynchronous via RabbitMQ.

### Repository Layout

```
services/
├── order/              # Saga orchestrator — owns the state machine
├── kitchen/            # Saga participant — manages tickets
├── payment/            # Saga participant — authorizes payments
├── shared/             # Shared library: saga framework, errors, logging
├── order_contract/     # Message contracts (commands/replies) for Order
├── kitchen_contract/   # Message contracts for Kitchen
└── payment_contract/   # Message contracts for Payment
deployments/dev/        # Docker Compose for local infrastructure
```

### Within Each Service

```
service/
├── cmd/server/main.go              # Entry point: wires DB, MQ, handlers
└── internal/
    ├── domain/                     # Entities and business rules
    ├── appservice/
    │   ├── create_*/               # Use-case logic per operation
    │   ├── proxy/                  # (Order only) generates commands to other services
    │   ├── port/                   # Repository interfaces
    │   ├── saga_command_handlers.go
    │   └── saga_state_machines.go  # (Order only) state machine registrations
    ├── adapter/
    │   ├── http/rest/              # Gin REST handlers
    │   ├── db/postgresql/          # GORM repository implementations
    │   └── pubsub/                 # RabbitMQ publisher/subscriber wrappers
    └── common/config/              # Environment-based configuration
```

### Saga Framework (`services/shared/saga/`)

The custom saga framework has three main components:

- **`Manager`** — Runs in the Order service. Executes state machine steps, sends commands, receives replies, and advances or compensates the saga.
- **`CommandHandler`** — Runs in Kitchen and Payment. Receives commands, calls business logic, sends back a Reply.
- **`StateMachine` / `StateMachineBuilder`** — Defines the ordered steps and their compensating actions.

### Create Order Saga Flow

```
Client → POST /order-service/api/orders
  → Order saved (PENDING)
  → Saga step 1: CreateTicket → Kitchen
      ← TicketAccepted reply
  → Saga step 2: AuthorizePayment → Payment
      ← PaymentAuthorized reply
  → Saga step 3: ApproveTicket → Kitchen
      ← TicketApproved reply
  → Saga step 4: ApproveOrder (local)
      → Order status = APPROVED

On failure at any step → compensating transactions run in reverse
```

### Messaging Contracts

Service contracts (`*_contract/` modules) define:
- Command structs sent to a service's `CommandChannel`
- Reply structs sent back to the requesting service's `ReplyChannel`
- Channel name constants (e.g., `KitchenServiceCommandChannel`)

Each service's `go.mod` uses `replace` directives to reference local contract and shared modules.

### Key Patterns

- **Domain validation** happens in entity constructors (e.g., `domain.NewOrder()`); app services trust the returned entity.
- **Idempotency**: `processed_messages` table deduplicates replayed MQ messages.
- **Error handling**: Use `apperr.AppErr` with typed codes (`BadRequest`, `UnprocessableEntity`). Wrap errors with `apperr.Wrap(err, "context")`.
- **Configuration**: `cleanenv`-based, loaded from environment variables. See `internal/common/config/` in each service for the full var list (prefixed `APP_`, `POSTGRES_`, `RABBIT_MQ_`).

### Database

Each service has its own PostgreSQL database (`order_service`, `kitchen_service`, `payment_service`). Schemas are initialized via `db/init/init.sql` in each service. The Order service schema includes `orders`, `sagas`, and `processed_messages` tables.
