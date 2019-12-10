package main

import (
    "fmt"
    _ "go.uber.org/automaxprocs"
    "runtime"
)

func main() {
    // Your application logic here.
    fmt.Println("real GOMAXPROCS", runtime.GOMAXPROCS(-1))
}
