package otp

type Sender interface {
	Send(to string, code string,  content string) error
}