package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-16653.crce179.ap-south-1-1.ec2.redns.redis-cloud.com:16653",
		Username: "default",
		Password: "QGbByRLjrpPlEylxaEQop2KwOKftKnpD",
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