package main

//a simple interactive ssh client
//pandaychen

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

var ce = func(err error, msg string) {
	if err != nil {
		log.Fatalf("%s error: %v", msg, err)
	}
}

func main() {
	client, err := ssh.Dial("tcp", "127.0.0.1:2222", &ssh.ClientConfig{
		User:            "foo",
		Auth:            []ssh.AuthMethod{ssh.Password("bar")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //2019/01/16 09:32:07 dial error: ssh: must specify HostKeyCallback
	})
	ce(err, "dial")

	go StartSession(client)
	go StartSession(client)

	select {}
}

func StartSession(client *ssh.Client) {
	session, err := client.NewSession()
	ce(err, "new session")
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	fmt.Println(session.Run("/usr/sbin/ifconfig eth1"))
}
