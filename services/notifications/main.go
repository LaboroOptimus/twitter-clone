package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"notifications/internal/db"
	"twitter/pkg/redisx"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	_ = godotenv.Load()

	redisAddr := getenv("REDIS_ADDR", "redis:6379")
	stream := getenv("REDIS_NOTIF_STREAM", "events:notifications")
	group := getenv("REDIS_GROUP", "notif-group")
	consumer := getenv("CONSUMER_NAME", hostnameFallback("notif-c1"))
	port := getenv("PORT", "8083")
	tlsCert := os.Getenv("TLS_CERT_FILE")
	tlsKey := os.Getenv("TLS_KEY_FILE")

	db.Connect()
	sqlDB := db.DB
	_ = sqlDB

	// --- Redis client через redisx ---
	ctx := context.Background()
	rdb := redisx.NewClient(redisx.Config{Addr: redisAddr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	log.Printf("connected to redis: %s", redisAddr)

	if err := rdb.XGroupCreateMkStream(ctx, stream, group, "$").Err(); err != nil {
		if !isBusyGroupErr(err) {
			log.Fatalf("XGroupCreateMkStream failed: %v", err)
		}
	}

	consumerCtx, cancelConsumer := context.WithCancel(context.Background())
	go func() {
		backoff := time.Second
		for {
			select {
			case <-consumerCtx.Done():
				log.Println("[consumer] stopping...")
				return
			default:
			}

			res, err := rdb.XReadGroup(consumerCtx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  []string{stream, ">"},
				Count:    100,
				Block:    5 * time.Second,
			}).Result()

			if err != nil && !errors.Is(err, redis.Nil) {
				log.Printf("[consumer] read error: %v", err)
				time.Sleep(backoff)
				if backoff < 5*time.Second {
					backoff *= 2
				}
				continue
			}
			backoff = time.Second

			for _, str := range res {
				for _, msg := range str.Messages {
					// Разбор события (пока просто логируем)
					event := map[string]any{
						"id":       msg.ID,
						"type":     asString(msg.Values["type"]),
						"user_id":  asString(msg.Values["user_id"]),
						"actor_id": asString(msg.Values["actor_id"]),
						"post_id":  asString(msg.Values["post_id"]),
					}
					b, _ := json.Marshal(event)
					log.Printf("[event] %s", string(b))

					// TODO: здесь вызов сервиса: запись в PG, инкремент счётчика в Redis

					// Подтверждаем обработку
					if err := rdb.XAck(consumerCtx, stream, group, msg.ID).Err(); err != nil {
						log.Printf("[consumer] XAck error: %v", err)
					}
				}
			}
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)
	router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := rdb.Ping(r.Context()).Err(); err != nil {
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	}).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("notifications-service listening on :%s (TLS:%t)", port, tlsCert != "" && tlsKey != "")
		var err error
		if tlsCert != "" && tlsKey != "" {
			err = srv.ListenAndServeTLS(tlsCert, tlsKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")

	cancelConsumer()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func hostnameFallback(def string) string {
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	return def
}
func isBusyGroupErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "BUSYGROUP")
}
func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
