package main

//a simple interactive ssh client
//pandaychen

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	client, err := ssh.Dial("tcp", "127.0.0.1:2222", &ssh.ClientConfig{
		User: "foo",
		//Auth:            []ssh.AuthMethod{ssh.Password("input your passwords")},
		Auth:            []ssh.AuthMethod{ssh.Password("bar")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //2019/01/16 09:32:07 dial error: ssh: must specify HostKeyCallback
	})
	ce(err, "dial")
	session, err := client.NewSession()
	ce(err, "new session")
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	termFD := int(os.Stdin.Fd())
	w, h, _ := terminal.GetSize(termFD)
	termState, _ := terminal.MakeRaw(termFD)
	defer terminal.Restore(termFD, termState)
	err = session.RequestPty("xterm-256color", h, w, modes)
	ce(err, "request pty")
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(30))
			session.SendRequest("keepalive@openssh.com", true, nil)

		}
	}()
	err = session.Shell()
	ce(err, "start shell")
	session.Wait()
	fmt.Println("Exit.")
}
