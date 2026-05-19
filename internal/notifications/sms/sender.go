package sms

import "fmt"

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) Send(to string, code string, content string) error {
	// DEV implementation (replace later with Twilio / Africa's Talking)
	fmt.Printf("[SMS OTP] to=%s code=%s\n", to, code)
	return nil
}