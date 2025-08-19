package main

import (
	"fmt"
	"net/http"
	"os"
	"posts/handlers"
	"posts/internal/db"
	"posts/internal/db/migrations"
	"posts/repo"
	"time"
	"twitter/pkg/authjwt"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Connect()
	migrations.Run()

	router := mux.NewRouter()

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7

	jwtService := authjwt.NewJWTService(accessSecret, refreshSecret, accessTTL, refreshTTL)

	postRepo := repo.NewPost(db.DB)

	api := router.PathPrefix("/api").Subrouter()
	api.Use(authjwt.AuthMiddleware(jwtService))

	//приватные
	api.HandleFunc("/posts", handlers.GetPosts(postRepo)).Methods(http.MethodGet)

	fmt.Println("Service post run on 8082!")

	if err := http.ListenAndServe(":8082", router); err != nil {
		panic(err)
	}
}
