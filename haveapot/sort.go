package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{15, 20, 31, 4, 25, 10, 7, 8, 11, 2, 33, 1, 3}
	//sort.Ints(a)
	//a := []int{1, 2, 3, 4, 5}
	sort.Ints(a)
	fmt.Println(a)
	b := sort.Search(len(a), func(i int) bool { return a[i] >= 30 })
	fmt.Println(b)
	c := sort.Search(len(a), func(i int) bool { return a[i] >= 20 })
	fmt.Println(c)
	d := sort.Search(len(a), func(i int) bool { return a[i] >= 7 })
	fmt.Println(d)
	d = sort.Search(len(a), func(i int) bool { return a[i] <= 8 }) //WRONG!,Your Should not use <= in a asc array
	fmt.Println(d)
}
