# Codematic Multi-Provider Virtual Wallet & Payment System

## Overview
This project is a modular, scalable backend system for managing virtual wallets and facilitating payments across multiple fintech providers (e.g., Paystack, Flutterwave, Stripe). It is designed for multi-tenancy (each tenant is a business with its own users and wallets), provider-agnostic transaction handling, and extensibility for new payment providers.

**Key Features:**
- Multi-tenant wallet management (deposits, withdrawals, transfers, transaction history)
- Provider abstraction layer for easy integration of new payment providers
- Idempotent transaction endpoints
- Webhook processing and replay for failed events
- JWT-based authentication and tenant separation
- PostgreSQL persistence
- Dockerized for local and cloud deployment
- Swagger/OpenAPI documentation
- Extensible for background jobs, Redis caching, fraud detection, and more

---

## Architecture
- **Go Fiber**: High-performance web framework
- **Domain-Driven Design**: Clear separation of domain, infrastructure, and handler layers
- **Provider Abstraction**: Easily add new payment providers by implementing a common interface
- **Multi-Tenancy**: Tenant isolation at the DB and API level
- **Idempotency**: Ensured for all transaction-related endpoints
- **Extensible**: Add background jobs, Redis, audit logs, and more with minimal changes

```
/ cmd/main.go           # Entry point, server setup, Swagger
/ internal/
    config/             # App config, logger
    domain/             # Business logic (auth, user, wallet, provider abstractions)
    handler/            # HTTP handlers, route registration
    infrastructure/     # DB, cache, provider integrations
    middleware/         # Auth, rate limiting, etc.
    router/             # Fiber app and server setup
    shared/             # Common models, utils
/ docs/                 # Auto-generated Swagger docs
```

---

## Setup & Local Development

### Prerequisites
- Go 1.20+
- Docker & Docker Compose
- PostgreSQL

### 1. Clone & Configure
```sh
git clone <your-repo-url>
cd codematic-backend
cp .env.example .env # Edit as needed
```

### 2. Start Services
```sh
docker-compose up --build
```

### 3. Run Migrations
```sh
go run internal/infrastructure/db/migrate.go
```

### 4. Generate Swagger Docs (if editing endpoints)
```sh
swag init --generalInfo cmd/main.go --output ./docs
```

### 5. Start the Server
```sh
go run cmd/main.go
```

---

## API Documentation
- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Base Path**: `/api`
- **Auth**: JWT Bearer token in `Authorization` header

---

## Key Endpoints (Examples)
- `POST /api/auth/login` — User login
- `POST /api/wallet/deposit` — Deposit funds
- `POST /api/wallet/withdraw` — Withdraw funds
- `POST /api/wallet/transfer` — Transfer between wallets
- `GET  /api/wallet/transactions` — Transaction history
- `POST /api/webhook/provider` — Provider webhook endpoint
- `POST /api/webhook/replay` — Replay failed webhook events

See Swagger UI for full details and request/response schemas.

---

## Extending the System
- **Add a Provider**: Implement the provider interface in `internal/domain/provider/` and register it
- **Add a Tenant**: Insert a new tenant in the `tenants` table; all APIs are tenant-aware
- **Add Background Jobs**: Integrate with a job queue (e.g., RabbitMQ, BullMQ)
- **Enable Redis Caching**: Configure Redis in `internal/infrastructure/cache/`
- **Add Audit Logs**: Use middleware or service hooks to log sensitive actions

---

## Assumptions & Trade-offs
- Real payment provider calls are mocked for demo purposes
- Focus is on architecture, extensibility, and edge-case handling, not production completeness
- Security best practices (rate limiting, input validation, etc.) are demonstrated but not exhaustive

---

## Developer Guide
- **Code Style**: Idiomatic Go, DDD, clear separation of concerns
- **Testing**: Add unit/integration tests in each domain and handler package
- **Swagger**: Annotate handlers with Swag comments for API docs
- **CI/CD**: (Bonus) Add GitHub Actions for lint/test/build

---

## License
MIT 