package main

/*
	封装了ss调用的read /write函数和其他网络接口调用
*/

import (
	"net"
)

const (
	AddrMask byte = 0xf
)

//ss的Conn 结构，封装了
type Conn struct {
	net.Conn        //原生net.Conn
	readBuf  []byte //
	writeBuf []byte
}

func t_proxy_conn(remote,local net.Conn) error{
	return nil
}

func main(){
	var cconn *Conn = new(Conn)
	var native_conn net.Conn

	t_proxy_conn(*cconn,native_conn)
}
