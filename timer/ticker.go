package main

import (
	"fmt"
	"time"
)

func main() {

	d := time.Duration(time.Second * 2)

	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			fmt.Println("timeout...")
		}
	}
}
