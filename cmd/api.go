package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/mdamaceno/notificator/app"
	"github.com/mdamaceno/notificator/config"
	"github.com/mdamaceno/notificator/internal/db"
)

func main() {
	dbconn, err := config.InitDB()

	defer dbconn.Close()

	if err != nil {
		log.Fatal(err)
	}

	q := db.New(dbconn)

	e := echo.New()
	routes := app.Routes{DB: dbconn, Queries: q, Echo: e}.Init()

	e.Logger.Fatal(routes.Start(":1323"))
}
