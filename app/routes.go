package app

import (
	"database/sql"

	"github.com/labstack/echo"
	"github.com/mdamaceno/notificator/app/controllers"
	"github.com/mdamaceno/notificator/internal/db"
)

type Routes struct {
	DB      *sql.DB
	Queries *db.Queries
	Echo    *echo.Echo
}

func (r Routes) Init() *echo.Echo {
	messageController := controllers.MessageController{DB: r.DB, Queries: r.Queries}
	api := r.Echo.Group("/api")

	api.POST("/message", messageController.Create)

	return r.Echo
}
