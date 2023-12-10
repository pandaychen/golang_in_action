package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	dest string
)

func init() {
	flag.StringVar(&dest, "dest", "x.x.x.x:18080", "dest udp addr")
}

func main() {
	// 创建一个UDP地址
	serverAddr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// 创建一个UDP连接
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		// 向服务器发送数据
		_, err := conn.Write([]byte("Hello, server!"))
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			continue
		}

		// 读取服务器返回的数据
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		// 输出接收到的数据
		fmt.Printf("Received %s from server\n", string(buf[:n]))

		// 等待1秒后再发送下一条消息
		time.Sleep(1 * time.Second)
	}
}
