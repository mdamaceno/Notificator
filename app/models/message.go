package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mdmaceno/notificator/app/_validation"
	"github.com/mdmaceno/notificator/app/services"
)

var Services = struct {
	Email string
	SMS   string
}{
	Email: "email",
	SMS:   "sms",
}

type Message struct {
	Id        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Service   string     `gorm:"type:varchar(255);not null"`
	Title     string     `gorm:"type:varchar(255);not null"`
	Body      string     `gorm:"type:text;not null"`
	SendAt    *time.Time `gorm:"type:timestamp;"`
	CreatedAt time.Time  `gorm:"not null;"`
	UpdatedAt time.Time  `gorm:"not null;"`

	Destinations []Destination
}

type MessageReceivers []string

type IncomingMessage struct {
	Service   string           `json:"service" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Body      string           `json:"body" validate:"required"`
	SendAt    string           `json:"send_at" validate:"datetime,omitempty"`
	Receivers MessageReceivers `json:"receivers" validate:"required"`
}

func NewMessage(im *IncomingMessage) (Message, error) {
	messageId := uuid.New()

	if im == nil {
		return Message{}, errors.New("message params is nil")
	}

	message := Message{
		Id:      messageId,
		Service: im.Service,
		Title:   im.Title,
		Body:    im.Body,
	}

	if im.SendAt != "" {
		sendAt, err := time.Parse(time.RFC3339, im.SendAt)

		if err != nil {
			return Message{}, err
		}

		message.SendAt = &sendAt
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
		err := _validation.Validate.Var(destination.Receiver, "email")
		if err == nil {
			emails = append(emails, destination.Receiver)
		}
	}

	return emails
}

func (m Message) Send() error {
	emails := m.FilterEmails()

	if len(emails) == 0 {
		return errors.New("no email to send")
	}

	errList := services.SendEmail(emails, m.Title, m.Body)

	if len(errList) > 0 {
		return errList[0]
	}

	return nil
}
