package repositories

import (
	"context"
	"database/sql"

	"github.com/mdmaceno/notificator/app/models"
	"github.com/mdmaceno/notificator/internal/db"
)

type MessageRepository struct {
	DB      *sql.DB
	Queries *db.Queries
}

func (mr MessageRepository) CreateMessage(message models.Message) error {
	tx, err := mr.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := mr.Queries.WithTx(tx)
	ctx := context.Background()

	m, err := qtx.CreateMessage(ctx, db.CreateMessageParams{
		Title:   message.Title,
		Body:    message.Body,
		Service: message.Service,
	})

	if err != nil {
		return err
	}

	for _, destination := range message.Destinations {
		_, err = qtx.CreateDestination(ctx, db.CreateDestinationParams{
			MessageID: m.ID,
			Receiver:  destination.Receiver,
		})
	}

	if err != nil {
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}
