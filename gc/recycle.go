package main

import (
    "fmt"
    "runtime"
    "time"
)

type A struct {
    Map map[int]int
}

func main() {
    var a A
    done := make(chan struct{})

    m1 := make(map[int]int)
    m1[1] = 1
    a.Map = m1

    // 监控原 map 是否被回收
    runtime.SetFinalizer(&m1, func(m *map[int]int) {
        fmt.Println("m1's map is finalized")
        close(done)
    })

    m2 := make(map[int]int)
    m2[2] = 2
    a.Map = m2 // 切断与原 map 的关联

    // 显式触发 GC（仅用于测试）
    runtime.GC()
//	fmt.Println(m1)
    select {
    case <-done:
        fmt.Println("m1's map was recycled")
	//fmt.Println(m1)
    case <-time.After(time.Second):
        fmt.Println("m1's map not recycled")
    }
}
