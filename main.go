package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/mdmaceno/notificator/app"
	"github.com/mdmaceno/notificator/internal/db"
)

func main() {
	url := "postgresql://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	dbconn, err := sql.Open("postgres", url)

	defer dbconn.Close()

	if err != nil {
		log.Fatal(err)
	}

	q := db.New(dbconn)

	e := echo.New()
	routes := app.Routes{DB: dbconn, Queries: q, Echo: e}.Init()

	e.Logger.Fatal(routes.Start(":1323"))
}
