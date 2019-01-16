package main

//a simple interactive ssh client
//pandaychen

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func main() {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	client, err := ssh.Dial("tcp", "172.27.0.17:22", &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password("input your passwords")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),		//2019/01/16 09:32:07 dial error: ssh: must specify HostKeyCallback
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
	err = session.Shell()
	ce(err, "start shell")
	err = session.Wait()
	ce(err, "return")
}
