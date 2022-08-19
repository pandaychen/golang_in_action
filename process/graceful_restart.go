package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/pandaychen/goes-wrapper/process"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %v!\n", os.Getpid())
}

func startServer(addr string, ln net.Listener) *http.Server {
	http.HandleFunc("/hello", handler)

	httpServer := &http.Server{
		Addr: addr,
	}
	go httpServer.Serve(ln)

	return httpServer
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":12345", "Address to listen on.")

	// Create (or import) a net.Listener and start a goroutine that runs
	// a HTTP server on that net.Listener.
	pl := process.ProcListener{
		Addr: addr,
	}
	ln, err := pl.RecreateListener(addr)
	if err != nil {
		fmt.Printf("Unable to create or import a listener: %v.\n", err)
		os.Exit(1)
	}
	server := startServer(addr, ln)

	// Wait for signals to either fork or quit.
	err = process.WaitSignals(addr, ln, server)
	if err != nil {
		fmt.Printf("Exiting: %v\n", err)
		return
	}
}
