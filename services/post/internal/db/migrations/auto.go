package migrations

import (
	"log"
	"posts/internal/db"
	"posts/internal/db/models"
)

func Run() {
	err := db.DB.AutoMigrate(&models.Post{}, &models.Image{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
