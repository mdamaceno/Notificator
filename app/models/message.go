package models

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdmaceno/notificator/app/services"
	"github.com/mdmaceno/notificator/internal/helpers"
)

type Email interface {
	Send(receivers []string, title string, body string) []error
}

type SMS interface {
	Send(receivers []string, message string) []error
}

var MessageType = struct {
	Email string
	SMS   string
}{
	Email: "email",
	SMS:   "sms",
}

var service = struct {
	Email Email
	SMS   SMS
}{
	Email: services.SendgridService{},
	SMS:   services.TwilioSMSService{},
}

type Message struct {
	ID        uuid.UUID
	Service   string
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time

	Destinations []Destination
}

type MessageReceivers []string

type IncomingMessage struct {
	Service   []string         `json:"service" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Body      string           `json:"body" validate:"required"`
	Receivers MessageReceivers `json:"receivers" validate:"required"`
}

func NewMessage(im *IncomingMessage) (Message, error) {
	messageId := uuid.New()

	if im == nil {
		return Message{}, errors.New("message params is nil")
	}

	message := Message{
		ID:      messageId,
		Service: strings.Join(im.Service, ","),
		Title:   im.Title,
		Body:    im.Body,
	}

	destinations := make([]Destination, len(im.Receivers))
	for i, receiver := range im.Receivers {
		destinations[i] = Destination{
			MessageId: messageId,
			Receiver:  receiver,
		}
	}

	message.Destinations = destinations

	return message, nil
}

func (m Message) FilterEmails() []string {
	emails := make([]string, 0)

	for _, destination := range m.Destinations {
		err := helpers.Validate.Var(destination.Receiver, "email")
		if err == nil {
			emails = append(emails, destination.Receiver)
		}
	}

	return emails
}

func (m Message) FilterPhoneNumbers() []string {
	phoneNumbers := make([]string, 0)

	for _, destination := range m.Destinations {
		err := helpers.Validate.Var(destination.Receiver, "e164")
		if err == nil {
			phoneNumbers = append(phoneNumbers, destination.Receiver)
		}
	}

	return phoneNumbers
}

func (m Message) hasService(id string) bool {
	s := strings.Split(m.Service, ",")
	for _, v := range s {
		if v == id {
			return true
		}
	}

	return false
}

func (m Message) Send() []error {
	var errList []error

	if m.hasService(MessageType.Email) {
		emails := m.FilterEmails()
		emailErr := service.Email.Send(emails, m.Title, m.Body)

		for _, err := range emailErr {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	if m.hasService(MessageType.SMS) {
		phoneNumbers := m.FilterPhoneNumbers()
		smsErr := service.SMS.Send(phoneNumbers, m.Body)

		for _, err := range smsErr {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	return errList
}
