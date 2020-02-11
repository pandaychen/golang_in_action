package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	RR_NGINX = 0
	RR_LVS   = 1
)

// 基于权重rr算法的接口
type RoundRobin interface {
	Next() interface{} //返回pick的节点
	Add(node interface{}, weight int)
	RemoveAll()
	Reset()
}

//算法实现工厂类
func NewWeightedRR(rtype int) RoundRobin {
	if rtype == RR_NGINX {
		return &WNGINX{}
	} else if rtype == RR_LVS {
		return &WLVS{}
	}
	return nil
}

//单个节点结构
type WeightNginx struct {
	Node            interface{}
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func (ww *WeightNginx) fail() {
	ww.EffectiveWeight -= ww.Weight
	if ww.EffectiveWeight < 0 {
		ww.EffectiveWeight = 0
	}
}

//nginx算法实现类（RoundRobin抽象结构的实例化）
type WNGINX struct {
	//need a lock
	nodes []*WeightNginx
	n     int
}

//增加权重节点
func (w *WNGINX) Add(node interface{}, weight int) {
	weighted := &WeightNginx{
		Node:            node,
		Weight:          weight,
		EffectiveWeight: weight}
	w.nodes = append(w.nodes, weighted)
	w.n++
}

func (w *WNGINX) RemoveAll() {
	w.nodes = w.nodes[:0] //slice
	w.n = 0
}

//下次轮询事件
func (w *WNGINX) Next() interface{} {
	if w.n == 0 {
		return nil
	}
	if w.n == 1 {
		return w.nodes[0].Node
	}

	//根据wrr算法选择下次节点
	return nextWeightedRRNode(w.nodes).Node
}

func nextWeightedRRNode(nodes []*WeightNginx) (best *WeightNginx) {
	total := 0

	for i := 0; i < len(nodes); i++ {
		w := nodes[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight //CurrentWeight：当前权重(变化的)，EffectiveWeight：初始化权重
		total += w.EffectiveWeight           //计算权重和
		if w.EffectiveWeight < w.Weight {
			//本节点权重值慢慢恢复
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	if best == nil {
		return nil
	}
	//选中节点减去权重和
	best.CurrentWeight -= total
	return best
}

func (w *WNGINX) Reset() {
	for _, s := range w.nodes {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0 //CurrentWeight初始值为0
	}
}

//LVS算法单个节点结构
type WeightLvs struct {
	Node   interface{}
	Weight int
}

//lvs算法实现类
type WLVS struct {
	nodes []*WeightLvs
	n     int
	gcd   int //通用的权重因子
	maxW  int //最大权重
	i     int //被选择的次数
	cw    int //当前的权重值
}

//下次轮询事件
func (w *WLVS) Next() interface{} {
	if w.n == 0 {
		return nil
	}

	if w.n == 1 {
		return w.nodes[0].Node
	}

	for {
		w.i = (w.i + 1) % w.n
		if w.i == 0 {
			w.cw = w.cw - w.gcd
			if w.cw <= 0 {
				w.cw = w.maxW
				if w.cw == 0 {
					return nil
				}
			}
		}
		if w.nodes[w.i].Weight >= w.cw {
			return w.nodes[w.i].Node
		}
	}
}

//增加权重节点
func (w *WLVS) Add(node interface{}, weight int) {
	weighted := &WeightLvs{Node: node, Weight: weight}
	if weight > 0 {
		if w.gcd == 0 {
			w.gcd = weight
			w.maxW = weight
			w.i = -1
			w.cw = 0
		} else {
			w.gcd = gcd(w.gcd, weight)
			if w.maxW < weight {
				w.maxW = weight
			}
		}
	}
	w.nodes = append(w.nodes, weighted)
	w.n++
}

func gcd(x, y int) int {
	var t int
	for {
		t = (x % y)
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}
func (w *WLVS) RemoveAll() {
	w.nodes = w.nodes[:0]
	w.n = 0
	w.gcd = 0
	w.maxW = 0
	w.i = -1
	w.cw = 0
}
func (w *WLVS) Reset() {
	w.i = -1
	w.cw = 0
}

var grr_obj = NewWeightedRR(RR_NGINX)

type handle struct {
	addrs []string
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//choose a backend,then proxy it...
	addr := grr_obj.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func StartServer() {
	//backend addr list
	h := &handle{}
	h.addrs = []string{"127.0.0.1:8080", "127.0.0.1:8081"}

	w := 1
	for _, e := range h.addrs {
		grr_obj.Add(e, w)
		w++
	}
	err := http.ListenAndServe(":12345", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	StartServer()
}
