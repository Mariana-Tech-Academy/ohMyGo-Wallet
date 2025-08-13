package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

//go get package github.com/go-redis/redis/v8


var Ctx = context.Background()
var Client *redis.Client
var Nil = redis.Nil

// connect to Redis client
func ConnectRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// test the connection
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}
