package db

import (
	"github.com/mdmaceno/notificator/app/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&models.Message{},
		&models.Destination{},
	)
}
