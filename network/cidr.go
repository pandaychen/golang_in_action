package main

import (
	"fmt"
	"net"
)

func main() {

	A := "10.12.221.0/24"
	B := "10.12.221.123"
	C := "10.12.221.10/29"

	ipA, ipnetA, _ := net.ParseCIDR(A)
	ipB := net.ParseIP(B)

	ipC, _, _ := net.ParseCIDR(C)

	fmt.Println("Network address A: ", A)
	fmt.Println("IP address      B: ", B)
	fmt.Println("ipA              : ", ipA)
	fmt.Println("ipnetA           : ", ipnetA)

	fmt.Printf("\nDoes A (%s) contain: B (%s)?\n", ipnetA, ipB)

	if ipnetA.Contains(ipB) {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}

	if ipnetA.Contains(ipC) {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}

}
