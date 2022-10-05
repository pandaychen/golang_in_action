package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/memberlist"
	// "net"
	"os"
	"strconv"
	"time"
)

var (
	bindPort = flag.Int("port", 8001, "gossip port")
)

func main() {
	flag.Parse()
	hostname, _ := os.Hostname()
	config := memberlist.DefaultLocalConfig()
	config.Name = hostname + "-" + strconv.Itoa(*bindPort)
	// config := memberlist.DefaultLocalConfig()
	config.BindPort = *bindPort
	config.AdvertisePort = *bindPort

	fmt.Println("config.DisableTcpPings", config.DisableTcpPings)
	fmt.Println("config.IndirectChecks", config.IndirectChecks)
	fmt.Println("config.RetransmitMult", config.RetransmitMult)

	fmt.Println("config.PushPullInterval", config.PushPullInterval)

	fmt.Println("config.ProbeInterval", config.ProbeInterval)

	fmt.Println("config.GossipInterval", config.GossipInterval)
	fmt.Println("config.GossipNodes", config.GossipNodes)

	fmt.Println("config.BindPort", config.BindPort)

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist:" + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	// 配置种子节点
	_, err = list.Join([]string{"127.0.0.1:8001", "127.0.0.1:8002"})
	fmt.Println("err", err)

	if err != nil {
		panic("Failed to join cluster:" + err.Error())
	}

	// Ask for members of the cluster
	for {
		fmt.Println("-------------start--------------")
		for _, member := range list.Members() {
			fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
		}
		fmt.Println("-------------end--------------")
		time.Sleep(time.Second * 3)

	}

}
