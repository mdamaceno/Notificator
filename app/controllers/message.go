package controllers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdamaceno/notificator/app/models"
	"github.com/mdamaceno/notificator/app/repositories"
	"github.com/mdamaceno/notificator/app/services"
	"github.com/mdamaceno/notificator/internal/db"
	"github.com/mdamaceno/notificator/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	deliveryCount int = 0
)

type MessageController struct {
	DB      *sql.DB
	Queries *db.Queries
}

var sender models.Sender = models.Sender{
	Email:    services.SendgridService{},
	SMS:      services.TwilioSMSService{},
	Whatsapp: services.TwilioWhatsappService{},
}

func (c MessageController) Create(ctx echo.Context) error {
	messageParams := new(models.IncomingMessage)

	if err := ctx.Bind(&messageParams); err != nil {
		helpers.ErrLog.Printf("Error binding message: %v", err)
		return ctx.JSON(http.StatusUnprocessableEntity, helpers.NewAPIErrorResponse(helpers.INVALID_REQUEST, nil))
	}

	if err := helpers.Validate.Struct(messageParams); err != nil {
		helpers.ErrLog.Printf("Error validating message: %v", err)
		mapErrors := helpers.MapValidationErrors(err)
		return ctx.JSON(http.StatusUnprocessableEntity, helpers.NewAPIErrorResponse(helpers.INVALID_REQUEST, mapErrors))
	}

	if len(messageParams.Receivers) == 0 {
		return ctx.JSON(http.StatusUnprocessableEntity, helpers.NewAPIErrorResponse(helpers.NOT_ENOUGH_RECEIVERS, nil))
	}

	message, err := models.NewMessage(messageParams)

	if err != nil {
		helpers.ErrLog.Printf("Error creating message: %v", err)
		response := helpers.NewAPIErrorResponse(helpers.INTERNAL_SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	message.Sender = sender

	errList := message.Send()

	for _, err := range errList {
		helpers.ErrLog.Printf("Error sending message: %v", err)
	}

	err = repositories.MessageRepository{DB: c.DB, Queries: c.Queries}.CreateMessage(message)

	if err != nil {
		helpers.ErrLog.Printf("Error saving message: %v", err)
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (c MessageController) Consume(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		helpers.Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()
	defer c.DB.Close()

	for d := range deliveries {
		deliveryCount++
		helpers.Log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		d.Ack(false)

		message, err := new(models.Message).FromJSON(d.Body)

		if err != nil {
			helpers.ErrLog.Printf("Error parsing message: %v", err)
		}

		message.Sender = sender

		errList := message.Send()

		for _, err := range errList {
			helpers.ErrLog.Printf("Error sending message: %v", err)
		}

		err = repositories.MessageRepository{DB: c.DB, Queries: c.Queries}.CreateMessage(message)

		if err != nil {
			helpers.ErrLog.Printf("Error saving message: %v", err)
		}
	}
}
