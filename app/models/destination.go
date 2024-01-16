package models

import (
	"time"

	"github.com/google/uuid"
)

type Destination struct {
	Id        uuid.UUID
	MessageId uuid.UUID
	Receiver  string
	CreatedAt time.Time
	UpdatedAt time.Time

	Message Message
}
