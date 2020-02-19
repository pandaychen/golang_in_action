package main

//多个goroutine作为writer，每个writer内部随机生成字符串写进去。唯一的reader读取数据并打印
//演示io.Pipe的用法

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func generate(writer *io.PipeWriter) {
	arr := make([]byte, 32)
	for {
		for i := 0; i < 32; i++ {
			arr[i] = byte(r.Uint32() >> 24)
		}
		n, err := writer.Write(arr)
		if nil != err {
			log.Fatal(err)
		} else {
			fmt.Println("write bytes:", n)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	rp, wp := io.Pipe()
	for i := 0; i < 20; i++ {
		go generate(wp)
	}
	time.Sleep(1 * time.Second)
	data := make([]byte, 64)
	for {
		n, err := rp.Read(data)
		if nil != err {
			log.Fatal(err)
		}
		if 0 != n {
			log.Println("main loop", n, string(data))
		}
		time.Sleep(1 * time.Second)
	}
}
