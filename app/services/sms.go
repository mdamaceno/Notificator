package services

var ServiceID = struct {
	Email string
	SMS   string
}{
	Email: "email",
	SMS:   "sms",
}

type SMS interface {
	Send(receivers []string, message string) []error
}
