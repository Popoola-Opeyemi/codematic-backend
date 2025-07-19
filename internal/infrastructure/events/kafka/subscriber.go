package kafka

import (
	"context"
	"log"
)

// Subscribe starts a consumer for the given topic and groupID, calling handler for each message.
func Subscribe(ctx context.Context, broker, topic, groupID string, handler func(key, value []byte)) error {
	reader := NewReader(broker, topic, groupID)
	go func() {
		defer reader.Close()
		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // context cancelled
				}
				log.Printf("Kafka read error: %v", err)
				continue
			}
			handler(m.Key, m.Value)
		}
	}()
	return nil
}
