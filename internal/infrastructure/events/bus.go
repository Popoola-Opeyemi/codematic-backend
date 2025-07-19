package events

import "context"

type EventBus interface {
	Publish(ctx context.Context, topic string, key string, value []byte) error
}
