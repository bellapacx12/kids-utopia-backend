package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers []string
	Username string
	Password string
	CAFile   string
	Topic string
	GroupID string
}

type Client struct {
	client *kgo.Client
}

func New(cfg Config) (*Client, error) {
	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    caPool,
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),

		// 🔥 CRITICAL FIX
		kgo.ConsumeTopics(cfg.Topic),

		kgo.ConsumerGroup(cfg.GroupID),
        kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.DialTLSConfig(tlsConfig),

		kgo.SASL(
			plain.Auth{
				User: cfg.Username,
				Pass: cfg.Password,
			}.AsMechanism(),
		),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx)
}

func (c *Client) Close() {
	c.client.Close()
}