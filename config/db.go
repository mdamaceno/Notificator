package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	url := "postgresql://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	return sql.Open("postgres", url)
}
