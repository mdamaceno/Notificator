package models

type Email interface {
	Send(receivers []string, title string, body string) []error
}

type SMS interface {
	Send(receivers []string, message string) []error
}

type Whatsapp interface {
	Send(receivers []string, message string) []error
}

type Sender struct {
	Email    Email
	SMS      SMS
	Whatsapp Whatsapp
}
