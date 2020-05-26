package main

import (
	"container/ring"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	limitCount  int        = 10 // 6s限频
	limitBucket int        = 6  // 滑动窗口个数
	curCount    int32      = 0  // 记录限频数量
	head        *ring.Ring      // 环形队列（链表）
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "0.0.0.0:9091") //获取一个tcpAddr
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	checkError(err)
	defer listener.Close()
	// 初始化滑动窗口
	head = ring.New(limitBucket)
	for i := 0; i < limitBucket; i++ {
		head.Value = 0
		head = head.Next()
	}
	// 启动执行器
	go func() {
		timer := time.NewTicker(time.Second * 1)
		for range timer.C { // 定时每隔1秒刷新一次滑动窗口数据
			subCount := int32(0 - head.Value.(int))
			newCount := atomic.AddInt32(&curCount, subCount)

			arr := [6]int{}
			for i := 0; i < limitBucket; i++ { // 这里是为了方便打印
				arr[i] = head.Value.(int)
				head = head.Next()
			}
			fmt.Println("move subCount,newCount,arr", subCount, newCount, arr)
			head.Value = 0
			head = head.Next()
		}
	}()

	for {
		conn, err := listener.Accept() // 在此处阻塞，每次来一个请求才往下运行handle函数
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(&conn) // 起一个单独的协程处理，有多少个请求，就起多少个协程，协程之间共享同一个全局变量limiting，对其进行原子操作。
	}
}

func handle(conn *net.Conn) {
	defer (*conn).Close()
	n := atomic.AddInt32(&curCount, 1)
	//fmt.Println("handler n:", n)
	if n > int32(limitCount) { // 超出限频
		atomic.AddInt32(&curCount, -1) // add 1 by atomic，业务处理完毕，放回令牌
		(*conn).Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\nError, too many request, please try again."))
	} else {
		mu := sync.Mutex{}
		mu.Lock()
		pos := head.Prev()
		val := pos.Value.(int)
		val++
		pos.Value = val
		mu.Unlock()
		time.Sleep(1 * time.Second)                                             // 假设我们的应用处理业务用了1s的时间
		(*conn).Write([]byte("HTTP/1.1 200 OK\r\n\r\nI can change the world!")) // 业务处理结束后，回复200成功。
	}
}

// 异常报错的处理
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
