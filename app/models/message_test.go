package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEmailService struct{}
type MockSMSService struct{}
type MockWhatsappService struct{}

func (s MockEmailService) Send(receivers []string, title string, body string) []error {
	return nil
}

func (s MockSMSService) Send(receivers []string, message string) []error {
	return nil
}

func (s MockWhatsappService) Send(receivers []string, message string) []error {
	return nil
}

func TestMessageModel(t *testing.T) {
	t.Run("NewMessage", func(t *testing.T) {
		t.Run("should return message when message params is not nil", func(t *testing.T) {
			im := &IncomingMessage{
				Service:   []string{MessageType.Email},
				Title:     "title",
				Body:      "body",
				Receivers: MessageReceivers{"john@email.com", "doe@email.com"},
			}

			message, err := NewMessage(im)

			assert.Nil(t, err)
			assert.Equal(t, im.Service, strings.Split(message.Service, ","))
			assert.Equal(t, im.Title, message.Title)
			assert.Equal(t, im.Body, message.Body)
			for i, destination := range message.Destinations {
				assert.Equal(t, im.Receivers[i], destination.Receiver)
			}
		})

		t.Run("should return error when message params is nil", func(t *testing.T) {
			_, err := NewMessage(nil)

			assert.NotNil(t, err)
		})
	})

	t.Run("FilterEmails", func(t *testing.T) {
		t.Run("should return emails when message has destinations", func(t *testing.T) {
			message := Message{
				Destinations: []Destination{
					{Receiver: "john@email.com"},
					{Receiver: "+628123456789"},
				},
			}

			emails := message.FilterEmails()

			assert.Equal(t, 1, len(emails))
			assert.Equal(t, message.Destinations[0].Receiver, emails[0])
		})
	})

	t.Run("FilterPhoneNumbers", func(t *testing.T) {
		t.Run("should return phone numbers when message has destinations", func(t *testing.T) {
			message := Message{
				Destinations: []Destination{
					{Receiver: "john@email.com"},
					{Receiver: "+628123456789"},
					{Receiver: "+628123456780"},
				},
			}

			phoneNumbers := message.FilterPhoneNumbers()

			assert.Equal(t, 2, len(phoneNumbers))
			assert.Equal(t, message.Destinations[1].Receiver, phoneNumbers[0])
			assert.Equal(t, message.Destinations[2].Receiver, phoneNumbers[1])
		})
	})

	t.Run("Send", func(t *testing.T) {
		mockEmailService := &MockEmailService{}
		mockSMSService := &MockSMSService{}
		mockWhatsappService := &MockWhatsappService{}

		message := Message{
			Title: "title",
			Body:  "body hello world",
			Destinations: []Destination{
				{Receiver: "johndoe@email.com"},
			},
		}

		t.Run("should call email service when service contains email", func(t *testing.T) {
			message.Service = strings.Join([]string{MessageType.Email}, ",")
			message.Sender = Sender{
				Email: mockEmailService,
			}

			err := message.Send()

			assert.Empty(t, err)
		})

		t.Run("should call sms service when service contains sms", func(t *testing.T) {
			message.Service = strings.Join([]string{MessageType.SMS}, ",")
			message.Sender = Sender{
				SMS: mockSMSService,
			}

			err := message.Send()

			assert.Empty(t, err)
		})

		t.Run("should call whatsapp service when service contains whatsapp", func(t *testing.T) {
			message.Service = strings.Join([]string{MessageType.Whatsapp}, ",")
			message.Sender = Sender{
				Whatsapp: mockWhatsappService,
			}

			err := message.Send()

			assert.Empty(t, err)
		})
	})
}
