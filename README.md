# Codematic Backend ⚡️

## Overview

**Codematic** is a modular, multi-tenant backend platform for wallet, payment, and webhook management. It is designed for fintech and SaaS platforms that require:
- Multi-tenant wallet management (each business/tenant has its own users and wallets)
- Provider-agnostic payment processing (integrate Paystack, Flutterwave, etc.)
- Secure authentication and role-based access (platform admin, tenant admin, user)
- Idempotent transaction endpoints and robust webhook handling
- Extensibility for new payment providers, background jobs, and analytics

Codematic is built in Go, fully Dockerized, and leverages PostgreSQL, Redis, and Kafka for persistence, caching, and event streaming. It provides a clean, well-documented API and is ready for local development or cloud deployment.

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
- `GET /api/transactions/{id}` — Get a single transaction (access controlled)
- `GET /api/transactions` — List transactions (with filters and access control)

#### Transaction Endpoints

- `GET /api/transactions/{id}`
  - Get a single transaction by ID.
  - **Access Control:**
    - **User:** Can only get their own transactions.
    - **Tenant Admin:** Can get transactions for their tenant.
    - **Admin:** Can get any transaction.

- `GET /api/transactions`
  - List transactions with optional filters (`status`, `limit`, `offset`).
  - **Access Control:**
    - **User:** Only their own transactions.
    - **Tenant Admin:** All transactions for their tenant, or filter by status.
    - **Admin:** All transactions, or filter by status.

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

# Developer Usage Guide

Welcome to the Codematic Backend! This section provides a concise guide for setup, development workflow, and key conventions.

---

## 🚀 Getting Started

### 1. Clone the Repository
```sh
git clone <your-repo-url>
cd codematic-backend
```

### 2. Environment Setup
- **Go version:** Ensure Go 1.20+ is installed.
- **Dependencies:**  
  ```sh
  go mod download
  ```

### 3. Configuration
- Copy and edit your environment config as needed:
  ```sh
  cp .env.example .env
  # Edit .env with your DB, cache, and other secrets
  ```

### 4. Database
- **Migrations:**  
  ```sh
  go run cmd/migrate/main.go
  ```
- **SQLC:**  
  All DB queries are managed via [sqlc](https://sqlc.dev/).  
  To regenerate Go code from SQL:
  ```sh
  sqlc generate
  ```

### 5. Running the App
```sh
go run cmd/main.go
```
Or use Docker Compose:
```sh
docker-compose up --build
```

---

## 🗂️ Project Structure

- `cmd/` – Entrypoints (main, migrations)
- `internal/`
  - `app/` – Bootstrap logic
  - `config/` – Configuration & logging
  - `domain/` – Business logic (auth, tenants, wallet, etc.)
  - `handler/` – HTTP handlers (Fiber)
  - `infrastructure/` – DB, cache, events, third-party integrations
  - `middleware/` – Fiber middleware
  - `router/` – HTTP server & routes
  - `shared/` – Common models & utilities

---

## 🧑‍💻 Development Workflow

- **HTTP Handlers:**  
  Follow the style in `internal/handler/wallet.go` for clean, consistent handlers.
- **Database Access:**  
  Use the repository pattern and sqlc-generated code (see `internal/infrastructure/db/sqlc/`).
- **Kafka Topics:**  
  Define all topics centrally in `internal/infrastructure/events/kafka/topics.go`.
- **Authentication:**  
  - `/auth/login` for tenant users (including tenant admins)
  - `/auth/admin` for platform admins  
  Check user roles after login.

---

## 🛠️ Useful Commands

- **Lint:**  
  ```sh
  golangci-lint run
  ```
- **Generate SQLC code:**  
  ```sh
  sqlc generate
  ```

---

## 📚 Documentation

- API docs: See `docs/swagger.yaml` or `docs/swagger.json`
- Postman collection: `postman/collections/`

---

## 🤝 Contributing

1. Fork & branch from `main`
2. Follow code style and conventions
3. Open a PR with a clear description

---

For more details, see the [README.md](./README.md)

--- 

