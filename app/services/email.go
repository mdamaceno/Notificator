package services

type Email interface {
	Send(receivers []string, title string, body string) []error
}
