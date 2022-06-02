package main

import (
	"fmt"
)

type fibonacciChan chan int

func (f fibonacciChan) Next() *int {
	c, ok := <-f
	if !ok {
		return nil
	}
	return &c
}

func fibonacci(limit int) fibonacciChan {
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
	f := fibonacci(20)
	fmt.Printf("%v ", *f.Next())
	fmt.Printf("%v ", *f.Next())
	for r := range f {
		fmt.Printf("%v ", r)
	}
}
