package otp

import (
	"strings"
)

type Router struct {
	email Sender
	sms   Sender
}

func NewRouter(email Sender, sms Sender) *Router {
	return &Router{
		email: email,
		sms:   sms,
	}
}

func (r *Router) Send(to string, subject string, content string) error {

	if strings.Contains(to, "@") {
		return r.email.Send(to, subject, content)
	}

	return r.sms.Send(to, subject, content)
}