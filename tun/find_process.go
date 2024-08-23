package main

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	// 替换为您要查找的源IP和源端口
	srcIP := "192.168.1.1"
	srcPort := uint32(61549)

	// 获取所有网络连接
	conns, err := net.Connections("all")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 查找与给定源IP和源端口匹配的连接
	var matchingConn *net.ConnectionStat
	for _, conn := range conns {
		if conn.Laddr.IP == srcIP && conn.Laddr.Port == srcPort {
			matchingConn = &conn
			break
		}
	}

	if matchingConn == nil {
		fmt.Println("No connection found for the given source IP and port.")
		return
	}

	// 获取进程信息
	p, err := process.NewProcess(matchingConn.Pid)
	if err != nil {
		fmt.Printf("Error getting process for PID %d: %v\n", matchingConn.Pid, err)
		return
	}

	// 获取进程的绝对路径
	exe, err := p.Exe()
	if err != nil {
		fmt.Printf("Error getting exe for process %d: %v\n", p.Pid, err)
		return
	}

	// 获取进程名字
	name, err := p.Name()
	if err != nil {
		fmt.Printf("Error getting name for process %d: %v\n", p.Pid, err)
		return
	}

	// 输出进程的绝对路径和名字
	fmt.Printf("Process %d (%s):\n", p.Pid, name)
	fmt.Println("Path:", exe)
}
