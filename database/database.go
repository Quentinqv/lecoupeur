package database

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Create the client variable
var client *redis.Client
var ctx = context.Background()

// Struct to store the URL with statistics of number of times it was accessed
type URL struct {
	URL     string `json:"url"`
	Counter int    `json:"counter"`
}

func Connect() *redis.Client {
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

	return client
}

func FlushAll() {
	client.FlushAll(ctx)
}

func StoreURL(url, id string) {
	// Store the URL in the database with a TTL of 10 minutes
	urlStruct := URL{URL: url, Counter: 0}
	urlString, err := json.Marshal(urlStruct)
	if err != nil {
		panic(err)
	}

	err = client.Set(ctx, id, urlString, 10*time.Minute).Err()
	if err != nil {
		panic(err)
	}
}

func GetURL(id string) (string, error) {
	// Retrieve the URL from the database
	urlStruct, err := client.Get(ctx, id).Result()
	if err != nil {
		return "", err
	}

	// Increment the counter and return the URL
	var url URL
	err = json.NewDecoder(strings.NewReader(urlStruct)).Decode(&url)
	if err != nil {
		return "", err
	}

	url.Counter++
	urlString, err := json.Marshal(url)
	if err != nil {
		return "", err
	}
	err = client.Set(ctx, id, urlString, 10*time.Minute).Err()
	if err != nil {
		return "", err
	}

	return url.URL, nil
}

func GetAllURLs() (map[string]URL, error) {
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	urls := make(map[string]URL)
	for _, key := range keys {
		urlStruct, err := client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var url URL
		err = json.NewDecoder(strings.NewReader(urlStruct)).Decode(&url)
		if err != nil {
			return nil, err
		}
		urls[key] = url
	}

	return urls, nil
}
