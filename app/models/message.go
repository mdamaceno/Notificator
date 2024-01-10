package models

import (
	"time"

	"github.com/google/uuid"
)

type Payload map[string]interface{}

type Message struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Service   string    `gorm:"type:varchar(255);not null"`
	Payload   Payload   `gorm:"type:varchar(255);not null;serializer:json"`
	SendAt    time.Time `gorm:"type:timestamp;"`
	CreatedAt time.Time `gorm:"not null;"`
	UpdatedAt time.Time `gorm:"not null;"`
}
