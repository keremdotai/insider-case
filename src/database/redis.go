package database

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// Define Redis client and context
var ctx = context.Background()
var RedisDB *redis.Client

// Connect to the Redis database
func ConnectRedisDB() {
	// Get environment variables for Redis connection
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")

	RedisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port), // Redis connection address
		Username: username,                         // Username for Redis connection
		Password: password,                         // Password for Redis connection
		DB:       0,                                // Default DB
	})

	// Check the connection status
	_, err := RedisDB.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
}

// Append message ID to the Redis list
func AppendMessageID(listKey string, messageID string) error {
	err := RedisDB.RPush(ctx, listKey, messageID).Err()
	if err != nil {
		return fmt.Errorf("failed to append message id to the Redis list: %v", err)
	}
	return nil
}

// Retrieve all list items from Redis
func RetrieveListItems(listKey string) ([]string, error) {
	items, err := RedisDB.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve Redis list items: %v", err)
	}
	return items, nil
}

// Close the connection to the Redis database
func CloseRedisDBConnection() error {
	err := RedisDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close the Redis connection: %v", err)
	}
	return nil
}
