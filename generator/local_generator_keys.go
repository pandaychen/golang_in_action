package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"

	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func GetRsaKeyPair() (string, string, error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	publicKey, err := ssh.NewPublicKey(&private.PublicKey)
	if err != nil {
		return "", "", err
	}

	privateKey := pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(private),
		Type:  "RSA PRIVATE KEY",
	})

	return string(privateKey), string(ssh.MarshalAuthorizedKey(publicKey)), nil
}

type Keys struct {
	Pk string
	Sk string
}

type AutoInc struct {
	start, step int
	queue       chan Keys
	running     bool
}

func New(start, step int) (ai *AutoInc) {
	ai = &AutoInc{
		start:   start,
		step:    step,
		running: true,
		queue:   make(chan Keys, 10),
	}
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	var i = 0
	for {
		sk, pk, _ := GetRsaKeyPair()
		skeyitem := Keys{
			Pk: pk,
			Sk: sk,
		}
		select {
		case ai.queue <- skeyitem:
			i++
		}
	}
}

func (ai *AutoInc) Id() *Keys {
	select {
	case id := <-ai.queue:
		return &id
	default:
		return nil
	}
}

func (ai *AutoInc) Close() {
	ai.running = false
	close(ai.queue)
}

func main() {
	ai := New(0, 2)
	defer ai.Close()

	time.Sleep(1 * time.Second)
	for {
		id := ai.Id()
		fmt.Print(id)
		time.Sleep(1 * time.Second)
	}
}
