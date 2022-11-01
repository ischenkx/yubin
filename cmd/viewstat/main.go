package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port  int
	Redis struct {
		Channel  string
		URL      string
		Password string
		Username string
		DB       int
	}
}

func getConfig() (config Config) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("failed to get the port:", err)
	}
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Println("failed to get the db:", err)
	}

	config.Port = port
	config.Redis.DB = db
	config.Redis.URL = os.Getenv("REDIS_URL")
	config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	config.Redis.Username = os.Getenv("REDIS_USERNAME")
	config.Redis.Channel = os.Getenv("REDIS_CHANNEL")
	return
}

func main() {
	config := getConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.URL,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		info := parts[len(parts)-1]
		cmd := redisClient.Publish(context.Background(), config.Redis.Channel, info)
		if cmd.Err() != nil {
			log.Println("failed to publish:", cmd.Err())
		}
	})

	addr := fmt.Sprintf(":%d", config.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("failed to start an http server:", err)
	}
}
