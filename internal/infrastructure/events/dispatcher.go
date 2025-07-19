package events

import (
	"context"
	"log"
)

var bus EventBus

func Init(b EventBus) {
	bus = b
}

func Dispatch(ctx context.Context, topic string, key string, payload []byte) {
	if bus == nil {
		log.Println("event dispatcher not initialized")
		return
	}
	if err := bus.Publish(ctx, topic, key, payload); err != nil {
		log.Printf("failed to publish event: %v\n", err)
	}
}
