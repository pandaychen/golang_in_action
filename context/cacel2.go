package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// gen generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the context once
	// they are done consuming generated integers not to leak
	// the internal goroutine started by gen.
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("recvs ctx.Done()...exit")
					return // returning not to leak the goroutine
				case dst <- n:
					fmt.Println("in subroutine,ctx=", ctx)
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel() // cancel when we are finished consuming integers
	fmt.Println("in main,ctx=", ctx)
	time.Sleep(1 * time.Second)
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
	cancel()
	time.Sleep(5 * time.Second)
}
