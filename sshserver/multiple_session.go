package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
	mu     sync.Mutex
}

func (s *SSHClient) RunCommand(command string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("unable to create session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("unable to run command: %v", err)
	}

	return string(output), nil
}

func main() {
	// 读取私钥文件a
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// 解析私钥
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// 构建SSH客户端配置
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			// 使用私钥进行身份验证
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到SSH服务器
	client, err := ssh.Dial("tcp", "1.2.3.4:22", config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	sshClient := &SSHClient{client: client}

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			output, err := sshClient.RunCommand("echo 'Hello from goroutine " + fmt.Sprintf("%d", n) + "'")
			if err != nil {
				log.Printf("Error running command: %v", err)
			} else {
				log.Printf("Output: %s", output)
			}
		}(i)
	}

	wg.Wait()
}
