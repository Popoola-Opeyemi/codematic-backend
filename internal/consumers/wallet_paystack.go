package consumers

import (
	"codematic/internal/domain/wallet"
	"codematic/internal/infrastructure/events/kafka"
	"context"

	"go.uber.org/zap"
)

const (
	walletGroupID = "wallet-paystack-consumer-group"
)

func StartWalletPaystackConsumer(
	ctx context.Context,
	broker string,
	walletService wallet.Service,
	logger *zap.Logger,
) {
	go func() {
		err := kafka.Subscribe(
			ctx,
			broker,
			kafka.PaystackWalletEventTopic,
			walletGroupID,
			func(key, value []byte) {
				walletService.HandlePaystackKafkaEvent(ctx, key, value)
			},
		)
		if err != nil {
			logger.Sugar().Errorf("Failed to subscribe to Paystack wallet events: %v", err)
		}
	}()
}
