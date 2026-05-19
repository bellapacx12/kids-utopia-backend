package email

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	apiKey string
	from   string
}

func NewSendGrid(apiKey, from string) *SendGrid {
	return &SendGrid{
		apiKey: apiKey,
		from:   from,
	}
}
func (s *SendGrid) Send(to string, subject string, html string) error {

	from := mail.NewEmail("Kids Utopia", s.from)
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, subject, toEmail, "", html)

	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}