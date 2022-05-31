package main

import (
	"context"
	"fmt"
	"time"
)

func hardWork(job interface{}) error {
	time.Sleep(time.Second * 10)
	return nil
}

func requestBlockedWork(ctx context.Context, job interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		done <- hardWork(job)
	}()
	select {
	case err := <-done:
		return err
	case p := <-panicChan:
		panic(p)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	now := time.Now()
	requestBlockedWork(context.Background(), "any")
	fmt.Println("elapsed:", time.Since(now))
}
