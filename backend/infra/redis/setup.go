package redisclient

import (
	"context"
	"fmt"
	"os"
	"sync"

	logger "gotemplate/utility/logger"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	once        sync.Once
	ctx         = context.Background()
	RedisClient *redis.Client
)

// Initialize sets up the Redis client (singleton).
func Initialize() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			logger.Logger.Error("No .env file found or unable to load it:", err)
		}
		addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
		password := os.Getenv("REDIS_PASSWORD")
		logger.Logger.Infof("Connecting to redis on addr: %s", addr)

		RedisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})

		// Test connection
		if err := RedisClient.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
		fmt.Println("Redis connected successfully")
	})
}

// CloseRedis closes the Redis connection gracefully.
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
