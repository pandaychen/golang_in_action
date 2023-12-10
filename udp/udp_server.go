package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 创建一个UDP地址
	serverAddr, err := net.ResolveUDPAddr("udp", ":18080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// 创建一个UDP连接
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		// 读取客户端发送的数据
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		// 输出接收到的数据
		fmt.Printf("Received %s from %s\n", string(buf[:n]), addr)

		// 将接收到的数据发送回客户端（echo）
		_, err = conn.WriteToUDP(buf[:n], addr)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			continue
		}
	}
}
