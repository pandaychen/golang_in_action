package main

//test dail a domain

import (
	"flag"
	"fmt"
	"net"
	"os"
)
func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", "www.baidu.com:80")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	done := make(chan string)
	go handleWrite(conn, done)
	go handleRead(conn, done)
	fmt.Println(<-done)
	fmt.Println(<-done)
}
func handleWrite(conn net.Conn, done chan string) {
	for i := 1; i > 0; i-- {
		_, e := conn.Write([]byte("GET / HTTP/1.1\r\nHost: www.baidu.com\r\n\r\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	done <- "Sent"
}
func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println(string(buf[:reqLen-1]))
	done <- "Read"
}
