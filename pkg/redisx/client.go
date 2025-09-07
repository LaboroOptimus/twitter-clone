package redisx

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr        string 
	Password    string
	DB          int
	DialTimeout time.Duration
}

func NewClient(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        cfg.Addr,
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: ifZero(cfg.DialTimeout, 5*time.Second),
	})
}

func ifZero[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}
