package main

import (
	"fmt"
)

func fibonacci(limit int) chan int {
	c := make(chan int)
	a := 0
	b := 1
	go func() {
		for {
			if limit == 0 {
				close(c)
				return
			}
			c <- a
			a, b = b, a+b
			limit--
		}
	}()
	return c
}

func main() {
	for r := range fibonacci(20) {
		fmt.Printf("%v ", r)
	}
}
