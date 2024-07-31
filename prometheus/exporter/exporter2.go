package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}

func init() {
	prometheus.MustRegister(NewNodeCollector())
}

var reqCount int32
var hostname string

type NodeCollector struct {
	requestDesc    *prometheus.Desc //Counter
	nodeMetrics    nodeStatsMetrics //混合方式
	goroutinesDesc *prometheus.Desc //Gauge
	threadsDesc    *prometheus.Desc //Gauge
	summaryDesc    *prometheus.Desc //summary
	histogramDesc  *prometheus.Desc //histogram
	mutex          sync.Mutex
}

// 混合方式数据结构
type nodeStatsMetrics []struct {
	desc    *prometheus.Desc
	eval    func(*mem.VirtualMemoryStat) float64
	valType prometheus.ValueType
}

// 初始化采集器
func NewNodeCollector() prometheus.Collector {
	host, _ := host.Info()
	hostname = host.Hostname
	return &NodeCollector{
		requestDesc: prometheus.NewDesc(
			"total_request_count",
			"请求数",
			[]string{"DYNAMIC_HOST_NAME"}, //动态标签名称
			prometheus.Labels{"STATIC_LABEL1": "静态值可以放在这里", "HOST_NAME": hostname}),
		nodeMetrics: nodeStatsMetrics{
			{
				desc: prometheus.NewDesc(
					"total_mem",
					"内存总量",
					nil, nil),
				valType: prometheus.GaugeValue,
				eval:    func(ms *mem.VirtualMemoryStat) float64 { return float64(ms.Total) / 1e9 },
			},
			{
				desc: prometheus.NewDesc(
					"free_mem",
					"内存空闲",
					nil, nil),
				valType: prometheus.GaugeValue,
				eval:    func(ms *mem.VirtualMemoryStat) float64 { return float64(ms.Free) / 1e9 },
			},
		},
		goroutinesDesc: prometheus.NewDesc(
			"goroutines_num",
			"协程数.",
			nil, nil),
		threadsDesc: prometheus.NewDesc(
			"threads_num",
			"线程数",
			nil, nil),
		summaryDesc: prometheus.NewDesc(
			"summary_http_request_duration_seconds",
			"summary类型",
			[]string{"code", "method"},
			prometheus.Labels{"owner": "example"},
		),
		histogramDesc: prometheus.NewDesc(
			"histogram_http_request_duration_seconds",
			"histogram类型",
			[]string{"code", "method"},
			prometheus.Labels{"owner": "example"},
		),
	}
}

// Describe returns all descriptions of the collector.
// 实现采集器Describe接口
func (n *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.requestDesc
	for _, metric := range n.nodeMetrics {
		ch <- metric.desc
	}
	ch <- n.goroutinesDesc
	ch <- n.threadsDesc
	ch <- n.summaryDesc
	ch <- n.histogramDesc
}

// Collect returns the current state of all metrics of the collector.
// 实现采集器Collect接口,真正采集动作
func (n *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	ch <- prometheus.MustNewConstMetric(n.requestDesc, prometheus.CounterValue, 0, hostname)
	vm, _ := mem.VirtualMemory()
	for _, metric := range n.nodeMetrics {
		ch <- prometheus.MustNewConstMetric(metric.desc, metric.valType, metric.eval(vm))
	}

	ch <- prometheus.MustNewConstMetric(n.goroutinesDesc, prometheus.GaugeValue, float64(runtime.NumGoroutine()))

	num, _ := runtime.ThreadCreateProfile(nil)
	ch <- prometheus.MustNewConstMetric(n.threadsDesc, prometheus.GaugeValue, float64(num))

	//模拟数据
	ch <- prometheus.MustNewConstSummary(
		n.summaryDesc,
		4711, 403.34,
		map[float64]float64{0.5: 42.3, 0.9: 323.3},
		"200", "get",
	)

	//模拟数据
	ch <- prometheus.MustNewConstHistogram(
		n.histogramDesc,
		4711, 403.34,
		map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
		"200", "get",
	)
	n.mutex.Unlock()
}
