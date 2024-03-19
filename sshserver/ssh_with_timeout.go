package main

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Conn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Conn) Read(b []byte) (int, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}

func main() {
	timeout := 3 * time.Second
	addr := "1.1.1.1:36000"
	config := ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
		Auth:            []ssh.AuthMethod{ssh.Password("xxxx")},
	}
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	timeoutConn := &Conn{conn, timeout, timeout}
	c, chans, reqs, err := ssh.NewClientConn(timeoutConn, addr, &config)
	if err != nil {
		panic(err)
	}

	client := ssh.NewClient(c, chans, reqs)
	defer client.Close()
}
