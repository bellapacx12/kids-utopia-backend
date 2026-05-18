package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Client struct {
	client *sqs.Client
	queue  string
}
func New(queueURL string, region string) (*Client, error) {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: sqs.NewFromConfig(cfg),
		queue:  queueURL,
	}, nil
}
func (c *Client) Send(message string) error {

	_, err := c.client.SendMessage(
		context.TODO(),
		&sqs.SendMessageInput{
			QueueUrl:    &c.queue,
			MessageBody: &message,
		},
	)

	return err
}
func (c *Client) Receive() ([]types.Message, error) {

	out, err := c.client.ReceiveMessage(
		context.TODO(),
		&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queue),
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     10,
		},
	)

	if err != nil {
		return nil, err
	}

	return out.Messages, nil
}
func (c *Client) Delete(receipt string) error {

	_, err := c.client.DeleteMessage(
		context.TODO(),
		&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(c.queue),
			ReceiptHandle: aws.String(receipt),
		},
	)

	return err
}