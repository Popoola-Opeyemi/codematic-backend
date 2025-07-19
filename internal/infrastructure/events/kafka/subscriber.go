package kafka

import (
	"context"
	"log"
)

func StartTokenPriceConsumer(broker string) {
	reader := NewReader(broker, TopicTokenPriceUpdated, "dexly-alert-engine")

	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("Kafka consumer error:", err)
				continue
			}
			go handleTokenPriceUpdate(msg.Key, msg.Value)
		}
	}()
}

func handleTokenPriceUpdate(key, value []byte) {
	log.Printf("Received token price update: %s", value)
}
