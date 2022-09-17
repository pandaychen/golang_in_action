package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func redisConnect() (rdb *redis.Client) {

	var (
		redisServer string
		port        string
		password    string
	)

	redisServer = os.Getenv("RedisUrl")
	if redisServer == "" {
		redisServer = "127.0.0.1"
	}
	port = os.Getenv("RedisPort")
	if port == "" {
		port = "6379"
	}
	password = os.Getenv("RedisPass")

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisServer + ":" + port,
		Password: password,
		DB:       0, // use default DB
	})

	return
}

func subMessage(channel string) {
	rdb := redisConnect()
	pubsub := rdb.Subscribe(context.Background(), channel)
	_, err := pubsub.Receive(context.Background())
	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func main() {
	channel := "hello"
	subMessage(channel)
}
