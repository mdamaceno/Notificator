package controllers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/mdmaceno/notificator/app/_response"
	"github.com/mdmaceno/notificator/app/_validation"
	"github.com/mdmaceno/notificator/app/models"
	"gorm.io/gorm"
)

type MessageController struct {
	DB *gorm.DB
}

func (c MessageController) Create(ctx echo.Context) error {
	messageParams := new(models.IncomingMessage)

	if err := ctx.Bind(&messageParams); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, _response.NewAPIErrorResponse(_response.INVALID_REQUEST, nil))
	}

	if err := _validation.Validate.Struct(messageParams); err != nil {
		mapErrors := _validation.MapValidationErrors(err)
		return ctx.JSON(http.StatusUnprocessableEntity, _response.NewAPIErrorResponse(_response.INVALID_REQUEST, mapErrors))
	}

	receivers := models.MessageReceiver(messageParams.Receivers)

	if messageParams.Service == "email" {
		for _, receiver := range receivers {
			if err := _validation.Validate.Var(receiver, "email"); err != nil {
				mapErrors := _validation.MapValidationErrors(err)
				return ctx.JSON(http.StatusUnprocessableEntity, _response.NewAPIErrorResponse(_response.INVALID_REQUEST, mapErrors))
			}
		}
	}

	messageId := uuid.New()

	message := models.Message{
		Id:      messageId,
		Service: messageParams.Service,
		Payload: models.Payload{
			Title: messageParams.Payload.Title,
			Body:  messageParams.Payload.Body,
		},
	}

	if messageParams.SendAt != "" {
		sendAt, err := time.Parse(time.RFC3339, messageParams.SendAt)

		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, _response.NewAPIErrorResponse(_response.INVALID_REQUEST, err.Error()))
		}

		message.SendAt = &sendAt
	}

	destinations := make([]models.Destination, len(receivers))
	for i, receiver := range receivers {
		destinations[i] = models.Destination{
			MessageId: messageId,
			Receiver:  receiver,
		}
	}

	c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&message).Error; err != nil {
			response := _response.NewAPIErrorResponse(_response.INTERNAL_SERVER_ERROR, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response)
		}

		if err := tx.Create(&destinations).Error; err != nil {
			response := _response.NewAPIErrorResponse(_response.INTERNAL_SERVER_ERROR, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response)
		}

		return nil
	})

	return ctx.JSON(http.StatusNoContent, nil)
}
