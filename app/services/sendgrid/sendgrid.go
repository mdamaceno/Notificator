package sendgrid

import (
	"log"
	"net/smtp"
	"os"
)

type SendgridService struct{}

func (s SendgridService) Send(emails []string, title string, body string) []error {
	from := os.Getenv("EMAIL_FROM")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")

	errList := []error{}

	for _, email := range emails {
		to := []string{email}
		message := []byte("To: " + email + "\r\n" +
			"From: " + from + "\r\n" +
			"Subject: " + title + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(host+":"+port, nil, from, to, message)

		if err != nil {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		return errList
	}

	return nil
}
