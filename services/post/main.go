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
	"twitter/pkg/redisx"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Connect()
	migrations.Run()

	router := mux.NewRouter()
	v := validator.New()
	rdb := redisx.NewClient(redisx.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	producer := redisx.NewProducer(rdb)

	notifStream := os.Getenv("REDIS_NOTIF_STREAM")
	if notifStream == "" {
		notifStream = "events:notifications"
	}

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7

	jwtService := authjwt.NewJWTService(accessSecret, refreshSecret, accessTTL, refreshTTL)

	postRepo := repo.NewPost(db.DB)
	likeRepo := repo.NewLike(db.DB)

	var postDeps = &handlers.PostsDeps{
		Validator: v,
		PostRepo:  postRepo,
	}

	var likeDeps = &handlers.LikesDeps{
		LikeRepo:    likeRepo,
		PostRepo:    postRepo,
		Validator:   v,
		Producer:    producer,
		NotifStream: notifStream,
	}

	router.Use(authjwt.AuthMiddleware(jwtService))

	router.HandleFunc("/posts", handlers.GetPosts(postDeps)).Methods(http.MethodGet)

	router.HandleFunc("/post", handlers.CreatePost(postDeps)).Methods(http.MethodPost)
	router.HandleFunc("/post", handlers.UpdatePost(postDeps)).Methods(http.MethodPatch)
	router.HandleFunc("/post", handlers.DeletePost(postDeps)).Methods(http.MethodDelete)
	router.HandleFunc("/post/{id:[0-9]+}", handlers.GetPostById(postDeps)).Methods(http.MethodGet)

	router.HandleFunc("/like", handlers.LikePost(likeDeps)).Methods(http.MethodPost)
	router.HandleFunc("/unlike", handlers.UnlikePost(likeDeps)).Methods(http.MethodPost)

	fmt.Println("Service post run on 8082!")

	if err := http.ListenAndServe(":8082", router); err != nil {
		panic(err)
	}
}
