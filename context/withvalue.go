package main

import (
	"context"
	"fmt"
	"time"
)

var key string = "ctxkey"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	//generate child ctx
	valueCtx := context.WithValue(ctx, key, "childctxvalue1")
	go watch(valueCtx)

	//generate another child ctx
	valueCtx2 := context.WithValue(ctx, key, "childctxvalue2")
	go watch(valueCtx2)

	//generate a grandson ctx
	valueCtx3, _ := context.WithCancel(valueCtx2)
	go watch(valueCtx3)

	//generate a grandson	ctx
	valueCtx4 := context.WithValue(valueCtx3, key, "childctxvalue4")
	go watch(valueCtx4)

	time.Sleep(4 * time.Second)
	fmt.Println("Call Parent Ctx.cancel() to notify all sub-routine Done(Exit)....")
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(2 * time.Second)
}

func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			//取出值
			fmt.Println(ctx.Value(key), "[exit]groutine recv ctx.Done(),exit..")
			return
		default:
			//取出值
			fmt.Println(ctx.Value(key), ctx, "groutines works...")
			time.Sleep(1 * time.Second)
		}
	}
}
