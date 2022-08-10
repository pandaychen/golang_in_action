package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	MyCounter prometheus.Counter
)

// init 注册指标
func init() {
	// 1.定义指标（类型，名字，帮助信息）
	MyCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_counter_total",
		Help: "自定义counter",
	})
	// 2.注册指标
	prometheus.MustRegister(MyCounter)
}

// Sayhello
func Sayhello(w http.ResponseWriter, r *http.Request) {
	// 接口请求量递增
	MyCounter.Inc()
	fmt.Fprintf(w, "Hello Wrold!")
}

func main() {

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/counter", Sayhello)
	http.ListenAndServe(":8080", nil)
}
