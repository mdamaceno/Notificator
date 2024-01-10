package models

import (
	"time"

	"github.com/google/uuid"
)

type Receiver []string

type Destination struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MessageId uuid.UUID `gorm:"type:uuid;not null"`
	Receiver  Receiver  `gorm:"type:varchar(255);not null;serializer:json"`
	CreatedAt time.Time `gorm:"not null;"`
	UpdatedAt time.Time `gorm:"not null;"`
}
