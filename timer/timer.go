package main

import (
	"fmt"
	"time"
)

func main() {

	d := time.Duration(time.Second * 2)
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			fmt.Println("timer timeout triggered...")
			// need reset
			t.Reset(time.Second * 2)
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("choose default")
		}
	}
}
