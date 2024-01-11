package models

import (
	"time"

	"github.com/google/uuid"
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
	Payload   Payload    `gorm:"type:jsonb;not null;serializer:json"`
	SendAt    *time.Time `gorm:"type:timestamp;"`
	CreatedAt time.Time  `gorm:"not null;"`
	UpdatedAt time.Time  `gorm:"not null;"`
}

type MessageReceiver []string

type IncomingMessage struct {
	Service   string          `json:"service" validate:"required"`
	Payload   IncomingPayload `json:"payload" validate:"required"`
	SendAt    string          `json:"send_at" validate:"datetime,omitempty"`
	Receivers MessageReceiver `json:"receivers" validate:"required"`
}
