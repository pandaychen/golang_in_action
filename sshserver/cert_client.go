package main

//a simple interactive ssh client
//pandaychen

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var addrName = flag.String("addr", "127.0.0.1:22", "addr")

func main() {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	flag.Parse()

	// parse the user's private key:
	pvtKeyBts, err := ioutil.ReadFile("./login.key")
	if err != nil { /* handle it */
	}

	signer, err := ssh.ParsePrivateKey(pvtKeyBts)
	if err != nil { /* handle it */
	}

	// parse the user's certificate:
	certBts, err := ioutil.ReadFile("./login-cert.pub")
	if err != nil { /* handle it */
	}

	cert, _, _, _, err := ssh.ParseAuthorizedKey(certBts)
	if err != nil { /* handle it */
	}

	// create a signer using both the certificate and the private key:
	certSigner, err := ssh.NewCertSigner(cert.(*ssh.Certificate), signer)
	if err != nil { /* handle it */
	}

	client, err := ssh.Dial("tcp", *addrName, &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(certSigner),
		},
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
