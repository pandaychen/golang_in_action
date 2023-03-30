package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

}

func batchHSet() {
	pipeline := client.Pipeline()
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("hkey%d", i)
		pipeline.HSet(ctx, key, map[string]interface{}{key: key})

	}
	_, err := pipeline.Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func batchGet() {
	pipeline := client.Pipeline()
	ctx := context.Background()
	result := make([]*redis.StringCmd, 0)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		result = append(result, pipeline.Get(ctx, key))
	}

	res, err := pipeline.Exec(ctx)
	fmt.Println(err, res)
	for _, r := range result {
		v, err := r.Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("no data")
				continue
			}

		}
		fmt.Println(v)
	}
}

func batchHGet() {
	pipeline := client.Pipeline()
	ctx := context.Background()
	result := make([]*redis.StringStringMapCmd, 0)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("hkey%d", i)
		result = append(result, pipeline.HGetAll(ctx, key))
	}
	res, err := pipeline.Exec(ctx)
	fmt.Println(err, res)
	for _, r := range result {
		v, err := r.Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(v)
	}
}

func main() {
	//batchSet()
	//batchHGet()
	batchGet()
}
