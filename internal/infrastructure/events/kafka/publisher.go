package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Broker string
}

func NewKafkaProducer(broker string) *KafkaProducer {
	return &KafkaProducer{Broker: broker}
}

func (kp *KafkaProducer) Publish(ctx context.Context, topic, key string, value []byte) error {
	writer := NewWriter(kp.Broker, topic)
	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("Kafka publish error: %v", err)
	}
	return err
}
