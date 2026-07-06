package kafka

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(client *Client) *Consumer {
	return &Consumer{
		client: client.client,
	}
}

func (c *Consumer) Poll(ctx context.Context) ([]*kgo.Record, error) {
	fetches := c.client.PollFetches(ctx)

	if err := fetches.Err(); err != nil {
		log.Printf("❌ PollFetches error: %v", err)
		return nil, err
	}

	var records []*kgo.Record
	count := 0

	fetches.EachRecord(func(r *kgo.Record) {
		count++


		records = append(records, r)
	})

	log.Printf("📦 POLLED RECORDS COUNT: %d", count)

	return records, nil
}

func (c *Consumer) Commit(ctx context.Context, records ...*kgo.Record) error {
	return c.client.CommitRecords(ctx, records...)
}