# Saga Orchestration in Go

A working microservices demo implementing the **Saga Orchestration** pattern with a reusable, framework-level saga engine in Go.

`Go` | `PostgreSQL` | `RabbitMQ` | `Docker Compose`

## Overview

This project demonstrates distributed transactions across three microservices — **Order** (saga orchestrator), **Kitchen**, and **Payment** (saga participants) — coordinated through asynchronous messaging via RabbitMQ, with each service owning its own PostgreSQL database.

Built after reading [_Microservices Patterns_](https://www.manning.com/books/microservices-patterns) by Chris Richardson. The saga framework in `services/shared/saga/` is inspired by [Eventuate Tram Sagas](https://github.com/eventuate-tram/eventuate-tram-sagas) (Java), reimplemented from scratch in Go with clean interfaces so it can be extracted and reused in other projects.

## Defining a Saga

A saga is defined as a sequence of forward steps and compensating actions using a fluent builder:

```go
saga.StateMachineBuilder().
    For("CreateOrder").
    WithCompensation(m.rejectOrder).       // undo order if any later step fails
    InvokeParticipant(m.createTicket).     // step 1: Kitchen creates a ticket
    WithCompensation(m.rejectTicket).      // undo ticket if any later step fails
    InvokeParticipant(m.authorizePayment). // step 2: Payment authorizes charge
    InvokeParticipant(m.approveTicket).    // step 3: Kitchen approves ticket
    InvokeParticipant(m.approveOrder).     // step 4: Order marks as APPROVED
    Build()
```

Builder methods:

- **`For(sagaType)`** — name of the saga
- **`InvokeParticipant(fn)`** — forward step that sends a command to a service
- **`WithCompensation(fn)`** — compensating action, executed in reverse if a later step fails
- **`OnReply(replyType, fn)`** — process a specific reply before advancing
- **`Build()`** — returns the state machine

## Create Order Saga Flow

```

                                          Ticket rejected /               Ticket creation failed /
  send CreateTicket  ┌─────────────────┐  send RejectOrder  ┌─────────────────┐  send RejectOrder  ┌─────────────────┐
 ●──────────────────►│   Creating      ├───────────────────►│   Rejecting     ├───────────────────►│     Order       ├──────►◉
                     │   ticket        │                    │   order         |                    │     rejected    |
                     └────────┬────────┘                    └─────────────────┘                    └─────────────────┘
                              │                                      ▲
           Ticket accepted /  │                                      │
       send AuthorizePayment  |                                      │
                              │                                      │
                              ▼           Payment failed /           │
                     ┌─────────────────┐  send RejectTicket ┌─────────────────┐
                     │  Authorizing    ├───────────────────►│   Rejecting     │
                     │  payment        |                    │   ticket        │
                     └────────┬────────┘                    └─────────────────┘
                              │
        Payment authorized /  │
          send ApproveTicket  │
                              │
                              ▼
                     ┌─────────────────┐
                     │   Approving     │
                     │   ticket        │
                     └────────┬────────┘
                              │
          Ticket approved /   │
          send ApproveOrder   │
                              │
                              ▼
                     ┌─────────────────┐
                     │   Approving     │
                     │   order         │
                     └────────┬────────┘
                              │
           Order approved     │
                              │
                              ▼
                     ┌─────────────────┐
                     │     Order       │
                     │     approved    │
                     └────────┬────────┘
                              │
                              ▼
                              ◉
```

Commands and replies flow asynchronously through RabbitMQ queues. Each service listens on its own command channel and replies back to the orchestrator. The saga **Manager** in the Order service tracks progress in a `sagas` table and advances or compensates based on each reply.

## Saga Framework

The reusable framework lives in `services/shared/saga/` and has three components:

**Manager** — the orchestrator runtime. Registers state machines, creates saga instances, sends commands, listens for replies, and advances or rolls back the saga. Runs in the orchestrator service (Order).

**CommandHandler** — the participant runtime. Receives commands on a channel, calls your handler function, and sends a reply back. Includes idempotency via a `processed_messages` table. Runs in participant services (Kitchen, Payment).

**StateMachineBuilder** — the fluent DSL shown above for defining saga step sequences.

### Port Interfaces

The framework depends only on interfaces defined in `port.go`, making it broker-agnostic and database-agnostic:

```go
type Producer interface {
    Send(channel string, message msg.Message) error
}

type Consumer interface {
    Consume(channel string) (dChan <-chan msg.Delivery, close func() error, err error)
}

type Repo interface {
    CreateSaga(*Saga) error
    UpdateSaga(*Saga) error
    FindSagaByID(id string) *Saga
    MessageRepo
    BeginTransaction() Transaction
}
```

This demo wires these to RabbitMQ and PostgreSQL, but any implementation of these interfaces works.

## Project Structure

```
services/
├── order/              # Saga orchestrator — owns the state machine
├── kitchen/            # Saga participant — manages tickets
├── payment/            # Saga participant — authorizes payments
├── shared/             # Saga framework, errors, logging
├── order_contract/     # Message contracts (commands/replies) for Order
├── kitchen_contract/   # Message contracts for Kitchen
└── payment_contract/   # Message contracts for Payment
deployments/dev/        # Docker Compose for local infrastructure
```

Each service follows hexagonal architecture:

```
services/order/
├── cmd/server/main.go              # Entry point, wiring
└── internal/
    ├── domain/                     # Entities, business rules
    ├── appservice/                 # Use cases, saga definitions, proxies
    │   ├── create_order/           # CreateOrder saga + service
    │   └── proxy/                  # Generates commands to other services
    ├── adapter/
    │   ├── http/rest/              # Gin REST handlers
    │   ├── db/postgresql/          # GORM repository implementations
    │   └── pubsub/                 # RabbitMQ adapters
    └── common/config/              # Environment-based configuration
```

Domain has no infrastructure imports. Adapters implement port interfaces.

## Getting Started

**Prerequisites:** Go 1.16+, Docker & Docker Compose

```bash
# 1. Start infrastructure (PostgreSQL, RabbitMQ, pgAdmin)
make dev_prep

# 2. Run each service in a separate terminal
make run_order      # port 5000
make run_kitchen    # port 5001
make run_payment    # port 5002
```

**Create an order:**

```bash
curl -X POST http://localhost:5000/order-service/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "vendor": "Pizza Place",
    "location": "123 Main St",
    "line_items": [
      { "id": "margherita", "quantity": 2, "note": "extra cheese" }
    ]
  }'
```

**Check order status:**

```bash
curl http://localhost:5000/order-service/api/orders
```

The order starts as `PENDING`. Once the saga completes, it becomes `APPROVED`. If any step fails, compensating transactions run and the order becomes `REJECTED`.

**Admin UIs:**

```bash
make dbgui    # pgAdmin at localhost:8000
make mqgui    # RabbitMQ management at localhost:15672
```

**Teardown:**

```bash
make dev_clean
```

## References

- [_Microservices Patterns_](https://www.manning.com/books/microservices-patterns) by Chris Richardson
- [Saga pattern](https://microservices.io/patterns/data/saga.html) — microservices.io
- [Eventuate Tram Sagas](https://github.com/eventuate-tram/eventuate-tram-sagas) — the Java framework that inspired this implementation
