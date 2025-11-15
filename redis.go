package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
	