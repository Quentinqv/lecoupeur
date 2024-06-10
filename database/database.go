package database

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Create the client variable
var client *redis.Client
var ctx = context.Background()

func Connect() {
	redisHost := ""
	if os.Getenv("REDIS_HOST") == "" {
		redisHost = "localhost"
	} else {
		redisHost = os.Getenv("REDIS_HOST")
	}

	redisPort := ""
	if os.Getenv("REDIS_PORT") == "" {
		redisPort = "6379"
	} else {
		redisPort = os.Getenv("REDIS_PORT")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
		DB:       0,
	})

	// Check that the connection is working
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func StoreURL(url, id string) {
	// Store the URL in the database with a TTL of 10 minutes
	err := client.Set(ctx, id, url, 10*time.Minute).Err()
	if err != nil {
		panic(err)
	}
}

func GetURL(id string) (string, error) {
	// Retrieve the URL from the database
	url, err := client.Get(ctx, id).Result()
	if err != nil {
		return "", err
	}
	return url, nil
}
