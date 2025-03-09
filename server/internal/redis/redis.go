package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Username: os.Getenv("REDIS_USERNAME"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	fmt.Println("Connected to Redis")
}