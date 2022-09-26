package main

import (
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
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

	worker := server.NewCustomQueueWorker("useless", 1, "high_queue")
	err = worker.Launch()
	if err != nil {
		log.Println("start worker error", err)
		return
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
