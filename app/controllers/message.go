package controllers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdmaceno/notificator/app/_response"
	"github.com/mdmaceno/notificator/app/_validation"
	"github.com/mdmaceno/notificator/app/models"
	"github.com/mdmaceno/notificator/app/repositories"
	"github.com/mdmaceno/notificator/internal/db"
)

type MessageController struct {
	DB      *sql.DB
	Queries *db.Queries
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

	if len(messageParams.Receivers) == 0 {
		return ctx.JSON(http.StatusUnprocessableEntity, _response.NewAPIErrorResponse(_response.NOT_ENOUGH_RECEIVERS, nil))
	}

	message, err := models.NewMessage(messageParams)

	if err != nil {
		response := _response.NewAPIErrorResponse(_response.INTERNAL_SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	go message.Send()

	repositories.MessageRepository{DB: c.DB, Queries: c.Queries}.CreateMessage(message)

	return ctx.JSON(http.StatusNoContent, nil)
}
