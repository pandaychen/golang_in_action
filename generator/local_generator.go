package main

import (
	"fmt"
	"time"
)

type AutoInc struct {
	start, step int
	queue       chan int
	running     bool
}

func New(start, step int) (ai *AutoInc) {
	ai = &AutoInc{
		start:   start,
		step:    step,
		running: true,
		queue:   make(chan int, 4),
	}
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	var i = 0
	for {
		select {
		case ai.queue <- i:
			i++
		}
	}
}

func (ai *AutoInc) Id() int {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	ai.running = false
	close(ai.queue)
}

func main() {
	ai := New(0, 2)
	defer ai.Close()
	for {
		id := ai.Id()
		fmt.Print(id)
		time.Sleep(1 * time.Second)
	}
}
