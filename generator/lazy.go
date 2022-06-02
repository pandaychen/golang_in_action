package main

import (
	"fmt"
)

var resume chan int

func integers() chan int {
	yield := make(chan int)
	count := 0
	go func() {
		for {	
			//blocked when full
			yield <- count
			count++
		}
	}()
	return yield
}

func generateInteger() int {
	return <-resume
}

func main() {
	resume = integers()
	for i := 0; i < 10; i++ {
		fmt.Println(generateInteger()) //=> 0
	}
}
