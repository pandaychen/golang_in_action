package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			return new(Person)
		},
	}
}

func main() {
	initPool()

	p := pool.Get().(*Person)
	p.Name = "first"
	pool.Put(p)
	fmt.Println("Get from pool:", pool.Get().(*Person))
	fmt.Println("Pool is empty", pool.Get().(*Person))
}
