package models

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdamaceno/notificator/internal/helpers"
)

var MessageType = struct {
	Email    string
	SMS      string
	Whatsapp string
}{
	Email:    "email",
	SMS:      "sms",
	Whatsapp: "whatsapp",
}

type Message struct {
	ID        uuid.UUID
	Service   string
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time

	Destinations []Destination
	Sender       Sender
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
	var phoneNumbers []string

	for _, destination := range m.Destinations {
		err := helpers.Validate.Var(destination.Receiver, "e164")
		if err == nil {
			phoneNumbers = append(phoneNumbers, destination.Receiver)
		}
	}

	return phoneNumbers
}

func (m Message) hasMessageType(id string) bool {
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

	if m.hasMessageType(MessageType.Email) {
		emails := m.FilterEmails()
		emailErr := m.Sender.Email.Send(emails, m.Title, m.Body)

		for _, err := range emailErr {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	if m.hasMessageType(MessageType.SMS) {
		phoneNumbers := m.FilterPhoneNumbers()
		smsErr := m.Sender.SMS.Send(phoneNumbers, m.Body)

		for _, err := range smsErr {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	if m.hasMessageType(MessageType.Whatsapp) {
		phoneNumbers := m.FilterPhoneNumbers()
		waErr := m.Sender.Whatsapp.Send(phoneNumbers, m.Body)

		for _, err := range waErr {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	return errList
}

func (m Message) FromJSON(body []byte) (Message, error) {
	var im IncomingMessage

	err := json.Unmarshal(body, &im)
	if err != nil {
		return Message{}, err
	}

	message, err := NewMessage(&im)

	if err != nil {
		return Message{}, err
	}

	return message, nil
}
