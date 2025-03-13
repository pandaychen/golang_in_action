package main

import (
    "fmt"
    "github.com/shirou/gopsutil/v3/process"
)

func main() {
    pid := 3021027  // 替换为目标进程 PID

    // 获取进程对象
    p, err := process.NewProcess(int32(pid))
    if err != nil {
        fmt.Printf("Failed to get process: %v\n", err)
        return
    }

    // 获取内存信息
    memInfo, err := p.MemoryInfo()
    if err != nil {
        fmt.Printf("Failed to get memory info: %v\n", err)
        return
    }

    // 输出 RES（RSS）的值（单位：字节）
    fmt.Printf("Process %d RES: %d KB\n", pid, memInfo.RSS/1024)
}
