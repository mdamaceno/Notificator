package app

import (
	"github.com/labstack/echo"
	"github.com/mdmaceno/notificator/app/controllers"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, DB *gorm.DB) *echo.Echo {
	messageController := controllers.MessageController{}
	api := e.Group("/api")

	api.POST("/message", messageController.Create)

	return e
}
