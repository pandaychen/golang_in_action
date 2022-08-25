package main

import (
	"fmt"
	"github.com/jrallison/go-workers"
	"time"
)

func myJob(message *workers.Msg) {
	fmt.Println("doing myjob start...")
	// do something with your message
	// message.Jid()
	// message.Args() is a wrapper around go-simplejson (http://godoc.org/github.com/bitly/go-simplejson)
	fmt.Println("doing myjob done...")
}

type myMiddleware struct{}

func (r *myMiddleware) Call(queue string, message *workers.Msg, next func() bool) (acknowledge bool) {
	// do something before each message is processed
	acknowledge = next()
	// do something after each message is processed
	return
}

func main() {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": "localhost:6379",
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})

	workers.Middleware.Append(&myMiddleware{})

	// pull messages from "myqueue" with concurrency of 10
	workers.Process("myqueue", myJob, 10)

	// pull messages from "myqueue2" with concurrency of 20
	workers.Process("myqueue2", myJob, 20)
	// Add a job to a queue
	workers.Enqueue("myqueue", "Add", []int{1, 2})
	res, err := workers.Enqueue("myqueue", "Add", []int{1, 2})
	fmt.Println(res, err)
	time.Sleep(5 * time.Second)

	// Add a job to a queue with retry
	workers.EnqueueWithOptions("myqueue2", "Add", []int{1, 2}, workers.EnqueueOptions{Retry: true})

	// stats will be available at http://localhost:8080/stats
	go workers.StatsServer(8081)

	// Blocks until process is told to exit via unix signal
	workers.Run()
}
