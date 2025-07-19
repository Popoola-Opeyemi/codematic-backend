package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codematic/internal/config"
	"codematic/internal/domain/tenants"
	"codematic/internal/infrastructure/db"
	"codematic/internal/infrastructure/events/kafka"

	"go.uber.org/zap"
)

const (
	groupID = "webhook-processor-group"
)

type DepositEvent struct {
	TenantID  string                 `json:"tenant_id"`
	WalletID  string                 `json:"wallet_id"`
	Amount    string                 `json:"amount"`
	Provider  string                 `json:"provider"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadAppConfig()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if cfg.KAFKA_BROKER_URL == "" {
		log.Fatal("KAFKA_BROKER_URL env var required")
	}
	dbConn := db.InitDB(cfg, logger)
	defer dbConn.Close()

	tenantsRepo := tenants.NewRepository(dbConn.Queries, dbConn.Pool)

	log.Println("Webhook processor started, listening for deposit events...")

	err := kafka.Subscribe(ctx, cfg.KAFKA_BROKER_URL, kafka.WalletDepositSuccessTopic, groupID, func(_ []byte, value []byte) {
		var event DepositEvent
		if err := json.Unmarshal(value, &event); err != nil {
			log.Printf("invalid deposit event: %v", err)
			return
		}
		dbTenant, err := tenantsRepo.GetTenantByID(ctx, event.TenantID)
		if err != nil {
			log.Printf("failed to get tenant: %v", err)
			return
		}
		if dbTenant.WebhookUrl == "" {
			log.Printf("tenant %s has no webhook_url set", event.TenantID)
			return
		}
		payload, _ := json.Marshal(event)
		resp, err := http.Post(dbTenant.WebhookUrl, "application/json", bytes.NewReader(payload))
		if err != nil {
			log.Printf("webhook POST failed: %v", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			log.Printf("webhook POST returned status %d", resp.StatusCode)
		}
	})
	if err != nil {
		log.Fatalf("failed to subscribe to kafka: %v", err)
	}

	<-ctx.Done()
	log.Println("shutting down webhook processor")
}
