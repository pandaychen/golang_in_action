package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
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

func pubMessage(channel, msg string) error {
	rdb := redisConnect()
	redisErr := rdb.Publish(context.Background(), channel, msg)
	return redisErr.Err()
}

func main() {
	channel := "hello"
	msgList := []string{"hello", "world"}

	for {
		for _, msg := range msgList {
			err := pubMessage(channel, msg)
			fmt.Println("send", msg, channel, err)
		}
		time.Sleep(2 * time.Second)
	}

}
