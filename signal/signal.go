package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Stops signals channel. This function exists
// in Go greater or equal to 1.1.
func signalStop(c chan<- os.Signal) {
	signal.Stop(c)
}

func signals() <-chan bool {
	quit := make(chan bool)

	go func() {
		signals := make(chan os.Signal)
		defer close(signals)

		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
		defer signalStop(signals)

		<-signals
		quit <- true
	}()

	return quit
}

func main() {

	s_quit := signals()

	<-s_quit

	fmt.Println("Recv Signals to Quit")

}
