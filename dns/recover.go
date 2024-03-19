package main

import (
	"fmt"

	"github.com/miekg/dns"
)

type dnsHandler struct{}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// 使用 defer 和 recover 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recovered from panic: %v\n", err)
		}
	}()

	// 故意触发 panic
	panic("An error occurred in ServeDNS")

	// 处理 DNS 请求（这里的代码实际上不会被执行，因为前面触发了 panic）
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	w.WriteMsg(m)
}

func main() {
	server := &dns.Server{Addr: ":5353", Net: "udp"}
	dns.HandleFunc(".", (&dnsHandler{}).ServeDNS)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
