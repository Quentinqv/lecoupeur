package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Create the client variable
var client *redis.Client
var ctx = context.Background()

func Connect() {
	// Connect to the database
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
