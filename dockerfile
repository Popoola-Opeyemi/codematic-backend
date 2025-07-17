# syntax=docker/dockerfile:1.4

##############################
# STAGE 1: Common base builder
##############################
FROM golang:1.24-alpine AS base

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk add --no-cache git tzdata ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

##############################
# STAGE 2: Dev image with Air
##############################
FROM base AS dev

# Install Air
RUN apk add --no-cache curl && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s && \
    mv ./bin/air /usr/local/bin/ && \
    rm -rf ./bin


# Set workdir to /app to match docker-compose volume mount
WORKDIR /app

# Air will rebuild and restart the server
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
# STAGE 4: Final minimal production image
##############################
FROM scratch AS prod

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

WORKDIR /app

COPY --from=builder /app/codematic /app/codematic-backend

USER 1000
EXPOSE 8080
ENTRYPOINT ["./codematic"]
