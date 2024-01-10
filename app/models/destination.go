package models

import (
	"time"

	"github.com/google/uuid"
)

type Destination struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MessageId uuid.UUID `gorm:"type:uuid;not null"`
	Receiver  string    `gorm:"type:varchar(255);not null;"`
	CreatedAt time.Time `gorm:"not null;"`
	UpdatedAt time.Time `gorm:"not null;"`

	Message Message
}
