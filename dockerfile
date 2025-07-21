# syntax=docker/dockerfile:1.4

##############################
# STAGE 1: Common base builder
##############################
FROM golang:1.24-alpine AS base

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOBIN=/go/bin

WORKDIR /app

RUN apk add --no-cache git tzdata ca-certificates

COPY go.mod go.sum ./
RUN go mod download

# Install sqlc (unset GOOS and GOARCH for correct platform binary)
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    env -u GOOS -u GOARCH go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY . .

# Generate SQLC code
RUN /go/bin/sqlc generate

##############################
# STAGE 2: Dev image with Air and Swag
##############################
FROM base AS dev

# Install Air
RUN apk add --no-cache curl && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s && \
    mv ./bin/air /usr/local/bin/ && \
    rm -rf ./bin

# Install swag CLI (Swagger/OpenAPI generator)
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    env -u GOOS -u GOARCH go install github.com/swaggo/swag/cmd/swag@latest && \
    mv /go/bin/swag /usr/local/bin/swag

WORKDIR /app

CMD ["air"]

##############################
# STAGE 3: Builder for production
##############################
FROM base AS builder

# Build metadata
ARG GIT_HASH
ARG BUILD_TIME

# Build the Go binary from cmd
RUN go build -ldflags="-X 'main.gitHash=${GIT_HASH}' -X 'main.buildTime=${BUILD_TIME}'" -o codematic ./cmd

##############################
# STAGE 4: Migration binary builder
##############################
FROM builder AS migrate-builder
RUN go build -o migrate ./cmd/migrate

##############################
# STAGE 5: Final minimal production image
##############################
FROM alpine AS prod

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copy Go binaries
COPY --from=builder /app/codematic /app/codematic
COPY --from=migrate-builder /app/migrate /app/migrate

# Copy SSL certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy migration files
COPY --from=builder /app/internal/infrastructure/db/migrations /app/internal/infrastructure/db/migrations

# Create logs directory with correct ownership
RUN adduser -D appuser && \
    mkdir -p /app/logs && \
    chown appuser:appuser /app/logs

USER appuser

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "/app/migrate && /app/codematic"]
