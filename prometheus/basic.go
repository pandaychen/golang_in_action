package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 创建一个计数器指标，用于记录请求的总数
var requestCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests.",
})

// 初始化Prometheus的指标收集器，并注册我们的计数器指标
func initPrometheus() {
	prometheus.MustRegister(requestCount)
}

func initPrometheusVal() {
	requestCount.Add(0)
	/*
		# HELP http_requests_total Total number of HTTP requests.
		# TYPE http_requests_total counter
		http_requests_total 0
	*/
}

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 初始化Prometheus指标收集器
	initPrometheus()
	initPrometheusVal()
	// 定义一个路由，用于暴露Prometheus指标
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 定义一个路由，用于处理正常的HTTP请求，并增加计数器的值
	r.GET("/", func(c *gin.Context) {
		requestCount.Add(1) // 每次请求时增加计数器的值
		c.String(http.StatusOK, "Hello, world!")
	})

	// 启动Gin服务器
	r.Run(":8088")
}
