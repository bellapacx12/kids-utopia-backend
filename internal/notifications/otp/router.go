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

func (r *Router) Send(to string, code string) error {

	if strings.Contains(to, "@") {
		return r.email.Send(to, code)
	}

	return r.sms.Send(to, code)
}