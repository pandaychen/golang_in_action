package main

import (
	"fmt"
	"strconv"
	"github.com/manifoldco/promptui"
)

func main() {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
			"Saturday", "Sunday"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
	
	days:=make([]string,0)
	for i:=0;i<25;i++{
		days=append(days,strconv.Itoa(i))
	}

	if result == "Tuesday" {
		prompt := promptui.Select{
			Label: "Select Hour",
			Items: days,
		}

		_, result, _ := prompt.Run()
		fmt.Printf("You choose day %q\n", result)
	}
}
