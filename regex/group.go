package main

import (
	"fmt"
	"regexp"
)

func main() {

	re := regexp.MustCompile("a.")

	ipre := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	fmt.Println(re.FindAllString("paranormal", -1))
	fmt.Println(re.FindAllString("paranormal", 2))
	fmt.Println(re.FindAllString("graal", -1))
	fmt.Println(re.FindAllString("none", -1))

	fmt.Println(len(ipre.FindAllString(" 1.2.3.4 ,1.2.3.1,11.", -1)))

}
