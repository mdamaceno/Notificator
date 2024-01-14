package models

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageModel(t *testing.T) {
	t.Run("NewMessage", func(t *testing.T) {
		t.Run("should return message when message params is not nil", func(t *testing.T) {
			im := &IncomingMessage{
				Service:   []string{Services.Email},
				Title:     "title",
				Body:      "body",
				SendAt:    "2021-01-01T00:00:00Z",
				Receivers: MessageReceivers{"john@email.com", "doe@email.com"},
			}

			message, err := NewMessage(im)
			timeParsed, _ := time.Parse(time.RFC3339, im.SendAt)

			assert.Nil(t, err)
			assert.Equal(t, im.Service, strings.Split(message.Service, ","))
			assert.Equal(t, im.Title, message.Title)
			assert.Equal(t, im.Body, message.Body)
			assert.Equal(t, &timeParsed, message.SendAt)
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

	t.Run("Send", func(t *testing.T) {
		t.Run("should return nothing when email is sent", func(t *testing.T) {
			message := Message{
				Title: "title",
				Body:  "body hello world",
				Destinations: []Destination{
					{Receiver: "johndoe@email.com"},
				},
			}

			err := message.Send()

			assert.Empty(t, err)
		})
	})
}
