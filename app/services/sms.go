package services

type SMS interface {
	Send(receivers []string, message string) []error
}
