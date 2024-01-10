package config

import (
	"log"

	"github.com/mdmaceno/notificator/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(a *AppConf) *gorm.DB {
	url := "host=" + a.DB_HOST + " user=" + a.DB_USER + " password=" + a.DB_PASS + " dbname=" + a.DB_NAME + " port=" +
		a.DB_PORT + " sslmode=disable TimeZone=UTC"

	orm, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established")

	db.RunMigrations(orm)

	log.Println("Database migrations ran successfully")

	return orm
}
