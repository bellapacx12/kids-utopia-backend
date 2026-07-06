package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}
func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client.client,
	}
}
func (p *Producer) Publish(
	ctx context.Context,
	topic string,
	key string,
	eventBytes []byte,
) error {

	log.Printf("📦 PRODUCER SEND topic=%s key=%s bytes=%d",
		topic, key, len(eventBytes),
	)

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: eventBytes, // ✅ already JSON
	}

	err := p.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	log.Printf("✅ PRODUCED SUCCESS topic=%s offset=sent", topic)

	return nil
}