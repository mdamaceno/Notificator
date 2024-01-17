package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdmaceno/notificator/app/models"
	"github.com/mdmaceno/notificator/app/repositories"
	"github.com/mdmaceno/notificator/app/services"
	"github.com/mdmaceno/notificator/internal/db"
	"github.com/mdmaceno/notificator/internal/helpers"
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
		return ctx.JSON(http.StatusUnprocessableEntity, helpers.NewAPIErrorResponse(helpers.INVALID_REQUEST, nil))
	}

	if err := helpers.Validate.Struct(messageParams); err != nil {
		mapErrors := helpers.MapValidationErrors(err)
		return ctx.JSON(http.StatusUnprocessableEntity, helpers.NewAPIErrorResponse(helpers.INVALID_REQUEST, mapErrors))
	}

	if len(messageParams.Receivers) == 0 {
		return ctx.JSON(http.StatusUnprocessableEntity, helpers.NewAPIErrorResponse(helpers.NOT_ENOUGH_RECEIVERS, nil))
	}

	message, err := models.NewMessage(messageParams)

	if err != nil {
		response := helpers.NewAPIErrorResponse(helpers.INTERNAL_SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	message.Sender = sender

	errList := message.Send()

	for _, err := range errList {
		log.Println(err)
	}

	err = repositories.MessageRepository{DB: c.DB, Queries: c.Queries}.CreateMessage(message)

	if err != nil {
		log.Println(err)
	}

	return ctx.JSON(http.StatusNoContent, nil)
}
