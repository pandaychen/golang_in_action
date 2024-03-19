package main

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
	"net"
	"sync"
	"time"
)

var (
	limiterCache = cache.New(5*time.Minute, 10*time.Minute)
	rateLimit    = rate.Limit(5) // 每秒最多5个请求
	rateLimitMu  sync.RWMutex    // 保护rateLimit的读写锁
)

func SetRateLimit(newRateLimit rate.Limit) {
	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()
	rateLimit = newRateLimit
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	// 获取客户端的源IP
	clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())

	// 检查客户端的IP是否在限速缓存中
	value, found := limiterCache.Get(clientIP)
	if !found {
		// 如果不存在，则创建一个新的令牌桶限速器，并将其添加到缓存中
		rateLimitMu.RLock()
		limiter := rate.NewLimiter(rateLimit, int(rateLimit))
		rateLimitMu.RUnlock()
		value = limiter
		limiterCache.Set(clientIP, value, cache.DefaultExpiration)
	}

	// 尝试从令牌桶中获取一个令牌，如果没有令牌，则返回错误
	limiter := value.(*rate.Limiter)
	if !limiter.Allow() {
		w.WriteMsg(makeRcode(r, dns.RcodeRefused))
		return
	}

	// 处理DNS请求
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	w.WriteMsg(m)
}

func makeRcode(req *dns.Msg, rcode int) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetRcode(req, rcode)
	return resp
}

func main() {
	server := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", handleRequest)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
