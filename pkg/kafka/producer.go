package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	client sarama.SyncProducer
}
func NewProducer(brokers []string) *Producer {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	client, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("Kafka init error:", err)
	}

	return &Producer{client: client}
}
func (p *Producer) Send(topic string, message []byte) error {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.client.SendMessage(msg)
	return err
}