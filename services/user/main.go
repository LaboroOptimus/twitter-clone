package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"twitter/pkg/authjwt"
	"user/handlers"
	"user/internal/db"
	"user/internal/db/migrations"
	"user/internal/validation"
	"user/repo"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	// валидатор
	validation.Init()
	fmt.Println("User service running on :8080")

	// коннект к бд
	db.Connect()

	// миграции
	migrations.Run()

	// jwt
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7

	jwtService := authjwt.NewJWTService(accessSecret, refreshSecret, accessTTL, refreshTTL)

	userRepo := repo.NewUser(db.DB)
	refreshRepo := repo.NewRefresh(db.DB)

	loginDeps := handlers.LoginDeps{
		Users:   userRepo,
		Refresh: refreshRepo,
		JWT:     jwtService,
	}

	router := mux.NewRouter()

	// публичные
	router.HandleFunc("/register", handlers.Register(userRepo)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.Login(loginDeps)).Methods(http.MethodPost)
	router.HandleFunc("/refresh", handlers.RefreshToken(loginDeps)).Methods(http.MethodPost)

	api := router.PathPrefix("/api").Subrouter()
	api.Use(authjwt.AuthMiddleware(jwtService))

	// приватные
	api.HandleFunc("/me", handlers.GetProfile(userRepo)).Methods(http.MethodGet)
	api.HandleFunc("/me", handlers.UpdateProfile(userRepo)).Methods(http.MethodPatch)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
