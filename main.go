package main

import (
	"github.com/labstack/echo"
	"github.com/mdmaceno/notificator/app"
	"github.com/mdmaceno/notificator/config"
)

func main() {
	env := config.Envs()
	DB := config.InitDB(env)

	e := echo.New()
	routes := app.InitRoutes(e, DB)

	e.Logger.Fatal(routes.Start(":1323"))
}
