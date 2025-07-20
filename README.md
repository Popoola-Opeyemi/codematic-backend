# Codematic Backend ⚡️

## Overview

**Codematic** is a modular, multi-tenant backend platform for wallet, payment, and webhook management. It is designed for fintech and SaaS platforms that require:
- Multi-tenant wallet management (each business/tenant has its own users and wallets)
- Provider-agnostic payment processing (integrate Paystack, Flutterwave, etc.)
- Secure authentication and role-based access (platform admin, tenant admin, user)
- Idempotent transaction endpoints and robust webhook handling
- Extensibility for new payment providers, background jobs, and analytics

Codematic is built in Go, fully Dockerized, and leverages PostgreSQL, Redis, and Kafka for persistence, caching, and event streaming. It provides a clean, well-documented API and is ready for local development or cloud deployment.

## Architecture

- **Go Fiber**: High-performance web framework for REST APIs
- **Domain-Driven Design**: Clear separation of domain logic, infrastructure, and handlers
- **Provider Abstraction**: Easily add new payment providers by implementing a common interface
- **Multi-Tenancy**: Tenant isolation at the DB and API level
- **Idempotency**: Ensured for all transaction-related endpoints
- **Event-Driven**: Kafka for event streaming and async processing
- **Extensible**: Add background jobs, Redis caching, audit logs, and more with minimal changes
- **Swagger/OpenAPI**: Auto-generated API documentation
- **Dockerized**: For local and production-ready deployments

## 🔧 Project Structure

```
.
├── cmd/                        # Main entry points (app, migrations)
│   └── migrate/                # Migration runner
├── internal/                   # Main application code
│   ├── app/                    # Application bootstrap
│   ├── config/                 # Configuration and logging
│   ├── consumers/              # Kafka consumers
│   ├── domain/                 # Domain logic and interfaces (auth, tenants, wallet, etc.)
│   ├── handler/                # HTTP route handlers
│   ├── infrastructure/         # External services
│   │   ├── cache/              # Redis and cache providers
│   │   ├── db/                 # Database, migrations, SQLC, queries
│   │   │   ├── migrations/     # SQL migration files
│   │   │   ├── query/          # Raw SQL queries
│   │   │   └── sqlc/           # SQLC generated Go code
│   │   ├── events/             # Event bus, Kafka, dispatcher
│   │   ├── search/             # Search integrations
│   │   └── ws/                 # WebSocket support
│   ├── middleware/             # Fiber middlewares
│   ├── router/                 # Fiber server and route bootstrap
│   ├── scheduler/              # Background job scheduler and jobs
│   │   └── jobs/               # Individual job implementations
│   └── shared/                 # Shared models and utilities
│       ├── model/              # Common models (error, jwt, response, etc.)
│       └── utils/              # Utility functions
│   └── thirdparty/             # Third-party payment provider clients
│       ├── baseclient/         # Base HTTP client
│       ├── flutterwave/        # Flutterwave integration
│       └── paystack/           # Paystack integration
├── docs/                       # Documentation (Swagger/OpenAPI)
├── logs/                       # App logs
├── postman/                    # Postman collections for API testing
│   └── collections/
├── tmp/                        # Temporary files
├── go.mod / go.sum             # Go dependencies
├── dockerfile                  # Dockerfile for the Go application
├── docker-compose.yml          # Multi-service Docker configuration
├── prometheus.yml              # Prometheus monitoring config
├── sqlc.yaml                   # SQLC config for generating DB code
└── README.md                   # You are here
```

## 🚀 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/)
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) for typed SQL queries
- (Optional) [Air](https://github.com/cosmtrek/air) for live reload during development

### Environment Setup

Set your environment variables in a `.env` file (see `docker-compose.yml` for required variables: `POSTGRES_DSN`, `REDIS_ADDR`, etc.).

Before running migrations, set your Goose database connection string:

```bash
export GOOSE_DBSTRING="postgres://user:password@localhost:5433/dbname?sslmode=disable"
```

### 🐳 Run with Docker Compose

To build and start the project:

```bash
docker-compose up --build codematic-dev
```

To start everything fresh:

```bash
docker-compose down -v --remove-orphans
# then
docker-compose up --build codematic-dev
```

This will rebuild all containers and clear old volumes.

#### Rebuilding and Restarting Services

```bash
# Stop all containers
docker-compose down
# Rebuild all images
docker-compose build
# Start all services in detached mode
docker-compose up -d
```

To rebuild and restart a specific service:

```bash
docker-compose build codematic-dev
docker-compose up -d codematic-dev
```

To view logs for a specific service:

```bash
docker-compose logs -f codematic-dev
```

After startup, the API will be available at:

- `http://localhost:9082/api`
- Health: `http://localhost:9082/api/health`
- Prometheus metrics: `http://localhost:9082/metrics`

## 🗄️ Database Schema

The core schema includes multi-tenant support and user roles:

```sql
CREATE TABLE "tenants" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "slug" VARCHAR UNIQUE NOT NULL,
  "webhook_url" TEXT DEFAULT '' NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY,
  "tenant_id" UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  "email" VARCHAR UNIQUE NOT NULL,
  "phone" VARCHAR,
  "password_hash" VARCHAR NOT NULL,
  "role" VARCHAR DEFAULT 'USER' CHECK (
    role IN ('PLATFORM_ADMIN', 'TENANT_ADMIN', 'USER')
  ),
  "is_active" BOOLEAN DEFAULT true,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);
```
Seed data includes example tenants and users (see migrations for details).
## Database Diagram

You can view the database schema here: [dbdiagram.io](https://dbdiagram.io/d/6808ff771ca52373f5105101) 



## 📚 API Documentation (Swagger/OpenAPI)

The API is documented using Swagger/OpenAPI. See [`docs/swagger.yaml`](docs/swagger.yaml) for the full spec.

### Main Endpoints

- `POST /api/auth/login` — Login for tenant users (regular or admin)
- `POST /api/auth/admin` — Login for platform admins
- `POST /api/auth/signup` — Register a new user
- `GET /api/auth/me` — Get current authenticated user
- `GET /api/tenant` — List all tenants
- `POST /api/tenant/create` — Create a new tenant
- `GET /api/tenant/{id}` — Get tenant by ID
- `PUT /api/tenant/{id}` — Update tenant
- `DELETE /api/tenant/{id}` — Delete tenant
- `GET /api/tenant/slug/{slug}` — Get tenant by slug
- `GET /api/wallet/{wallet_id}/balance` — Get wallet balance
- `GET /api/wallet/{wallet_id}/transactions` — Get wallet transactions
- `POST /api/wallet/transfer` — Transfer funds between wallets
- `POST /api/wallet/withdraw` — Withdraw funds from a wallet
- `POST /api/webhook/{provider}` — Handle provider webhook

For full request/response models and error codes, see the Swagger file or run the service and visit `/docs` if enabled.

## 🛠️ Development Guide

### Migrations

Create a new migration:

```bash
goose -dir ./internal/infrastructure/db/migrations create <migration_name> sql
```

Apply migrations:

```bash
goose -dir ./internal/infrastructure/db/migrations up
```

### SQLC Code Generation

After updating SQL files, regenerate Go code:

```bash
sqlc generate
```

### Running Locally (without Docker)

1. Start PostgreSQL, Redis, and Kafka locally (see `docker-compose.yml` for config).
2. Set environment variables as in `.env`.
3. Run migrations with Goose.
4. Start the server:

```bash
go run cmd/main.go
```

### Testing

Tests follow Go's `*_test.go` convention. Place tests in the relevant package directories.

## 📦 Tech Stack

- **Go (Fiber v2)** — Web framework
- **PostgreSQL + SQLC** — Typed SQL queries and database
- **Redis** — Caching
- **Kafka** — Event streaming
- **Docker** — Containerization
- **Prometheus/Grafana** — Monitoring
- **Zap** — Structured logging
- **JWT** — Authentication
- **gocron** — Job scheduling

## 📜 License

MIT License — Free to use, modify, and contribute. 

