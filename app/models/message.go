package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdmaceno/notificator/app/services"
	"github.com/mdmaceno/notificator/app/services/twilio"
	"github.com/mdmaceno/notificator/internal/helpers"
)

var Services = struct {
	Email string
	SMS   string
}{
	Email: "email",
	SMS:   "sms",
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

func (m Message) Send() error {
	emails := m.FilterEmails()

	for _, destination := range m.Destinations {
		err := helpers.Validate.Var(destination.Receiver, "e164")
		if err == nil {
			phoneNumbers = append(phoneNumbers, destination.Receiver)
		}
	}

	errList := services.SendEmail(emails, m.Title, m.Body)

	if len(errList) > 0 {
		return errList[0]
	}

	return nil
}
