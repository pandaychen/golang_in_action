package main

import (
	"fmt"
	"time"
	"github.com/go-redis/redis"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		PoolSize:     10,
		ReadTimeout:  time.Millisecond * time.Duration(100),
		WriteTimeout: time.Millisecond * time.Duration(100),
		IdleTimeout:  time.Second * time.Duration(60),
		//add MinIdleConns
		MinIdleConns: 3,
	})

	_, err := Client.Ping().Result()
	if err != nil {
		panic("init redis error")
	} else {
		fmt.Println("init redis ok")
	}
}

func get(key string) (string, bool) {
	r, err := Client.Get(key).Result()
	if err != nil {
		return "", false
	}

	return r, true
}

func set(key string, val string, expTime int32) {
	Client.Set(key, val, time.Duration(expTime)*time.Second)
}

func main() {
	for {
		for i := 0; i < 10000; i++ {
			set("name", "x", 100)
			s, b := get("name")
			fmt.Println(s, b)
		}

		time.Sleep(1 * time.Second)
	}
}
