package migrations

import (
	"log"
	"user/internal/db"
	"user/internal/db/models"
)

func Run() {
	err := db.DB.AutoMigrate(&models.User{}, &models.Refresh{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
