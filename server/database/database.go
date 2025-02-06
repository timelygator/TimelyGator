package database

import (
	"timelygator/server/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("timelygator.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	DB.AutoMigrate(&models.Task{})

	return nil
}
