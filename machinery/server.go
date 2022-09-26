package main

import (
	"context"
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func main() {

	cnf, err := config.NewFromYaml("./config.yml", false)
	if err != nil {
		log.Println("config failed", err)
		return
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		log.Println("start server failed", err)
		return
	}

	// 注册任务
	err = server.RegisterTask("sum", Sum)
	if err != nil {
		log.Println("reg task failed", err)
		return
	}
	err = server.RegisterTask("call", CallBack)
	if err != nil {
		log.Println("reg task failed", err)
		return
	}

	//task signature
	signature1 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		RetryTimeout: 100,
		RetryCount:   3,
		RoutingKey:   "high_queue",
	}

	signature2 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		RetryTimeout: 100,
		RetryCount:   3,
		RoutingKey:   "high_queue",
	}

	signature3 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		RetryTimeout: 100,
		RetryCount:   3,
		RoutingKey:   "high_queue",
	}

	//// group
	group, err := tasks.NewGroup(signature1, signature2, signature3)
	if err != nil {
		log.Println("add group failed", err)
		return
	}

	asyncResults, err := server.SendGroupWithContext(context.Background(), group, 0)
	if err != nil {
		log.Println(err)
		return
	}
	for _, asyncResult := range asyncResults {
		results, err := asyncResult.Get(1)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf(
			"%v  %v \n",
			asyncResult.Signature.Args[0].Value,
			tasks.HumanReadableResults(results),
		)
	}

}

func Sum(args []int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

// Multiply ...
func CallBack(args ...int64) (int64, error) {
	sum := int64(1)
	for _, arg := range args {
		sum *= arg
	}
	return sum, nil
}
