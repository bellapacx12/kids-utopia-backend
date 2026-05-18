package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	client sarama.Consumer
}
func NewConsumer(brokers []string) *Consumer {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal("Kafka consumer error:", err)
	}

	return &Consumer{client: client}
}
func (c *Consumer) Consume(topic string, handler func([]byte)) {

	partitionConsumer, err := c.client.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range partitionConsumer.Messages() {
		handler(msg.Value)
	}
}