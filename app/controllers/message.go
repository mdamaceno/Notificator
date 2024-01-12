package controllers

import (
	"net/http"

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

	if len(messageParams.Receivers) == 0 {
		return ctx.JSON(http.StatusUnprocessableEntity, _response.NewAPIErrorResponse(_response.NOT_ENOUGH_RECEIVERS, nil))
	}

	message, err := models.NewMessage(messageParams)

	if err != nil {
		response := _response.NewAPIErrorResponse(_response.INTERNAL_SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	go message.Send()

	if err := c.DB.Create(&message).Error; err != nil {
		response := _response.NewAPIErrorResponse(_response.INTERNAL_SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	return ctx.JSON(http.StatusNoContent, nil)
}
