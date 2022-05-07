package main

import (
	"fmt"
)

func main() {
	var array []*int

	var (
		a = 1
		b = 2
		c = 3
		d = 4

		pa = &a
		pb = &b
		pc = &c
		pd = &d
	)
	array = append(array, &a)
	array = append(array, &b)
	array = append(array, &c)
	fmt.Println(array, pa, pb, pc, pd)

	//交换两个指针变量的值，不影响array
	pa, pd = pd, pa
	fmt.Println("after swap:")
	fmt.Println(array, pa, pb, pc, pd)
}
