package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"twitter/pkg/authjwt"
	"twitter/pkg/redisx"
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

	rdb := redisx.NewClient(redisx.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	producer := redisx.NewProducer(rdb)

	notifStream := os.Getenv("REDIS_NOTIF_STREAM")
	if notifStream == "" {
		notifStream = "events:notifications"
	}

	// jwt
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7

	jwtService := authjwt.NewJWTService(accessSecret, refreshSecret, accessTTL, refreshTTL)

	userRepo := repo.NewUser(db.DB)
	refreshRepo := repo.NewRefresh(db.DB)
	followRepo := repo.NewFollow(db.DB)

	loginDeps := handlers.LoginDeps{
		Users:   userRepo,
		Refresh: refreshRepo,
		JWT:     jwtService,
	}

	followDeps := &handlers.FollowDeps{
		Follow: followRepo,
		Producer: producer,
		NotifStream: notifStream,
	}

	userDeps := handlers.UserDeps{
		Users: userRepo,
	}

	router := mux.NewRouter()

	// публичные
	router.HandleFunc("/register", handlers.Register(userRepo)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.Login(loginDeps)).Methods(http.MethodPost)
	router.HandleFunc("/refresh", handlers.RefreshToken(loginDeps)).Methods(http.MethodPost)

	private := router.NewRoute().Subrouter()
	private.Use(authjwt.AuthMiddleware(jwtService))

	// приватные
	private.HandleFunc("/me", handlers.GetProfile(userDeps)).Methods(http.MethodGet)
	private.HandleFunc("/me", handlers.UpdateProfile(userDeps)).Methods(http.MethodPatch)
	private.HandleFunc("/user", handlers.GetUser(userDeps)).Methods(http.MethodGet)

	private.HandleFunc("/follow/{id:[0-9]+}", handlers.FollowUser(followDeps)).Methods(http.MethodPost)
	private.HandleFunc("/unfollow/{id:[0-9]+}", handlers.UnfollowUser(followDeps)).Methods(http.MethodDelete)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
