package email

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESSender struct {
	client *ses.Client
	from   string
}

func NewSESSender(region, from string) (*SESSender, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &SESSender{
		client: ses.NewFromConfig(cfg),
		from:   from,
	}, nil
}

func (s *SESSender) Send(to string,  code string) error {
	log.Println("📧 SES SENDING TO:", to)
log.Println("📧 FROM:", s.from)

	subject := "Your OTP Code"
	body := "Your KIDS UTOPIA OTP is: " + code

	_, err := s.client.SendEmail(context.TODO(), &ses.SendEmailInput{
		Source: aws.String(s.from),
		Destination: &sesTypes.Destination{
			ToAddresses: []string{to},
		},
		Message: &sesTypes.Message{
			Subject: &sesTypes.Content{
				Data: aws.String(subject),
			},
			Body: &sesTypes.Body{
				Text: &sesTypes.Content{
					Data: aws.String(body),
				},
			},
		},
	})

	return err
}