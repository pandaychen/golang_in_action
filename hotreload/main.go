package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"log"
)

var configPath string

func main() {
	configPath="./config.toml"
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
			ReloadConfig()
			log.Println("Reloaded config succ")
		}
	}()

	http.HandleFunc("/", showcfg) 
	err := http.ListenAndServe(":2345", nil) 
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


func showcfg(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", Config().DB.Server)
}
