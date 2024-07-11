package main

import (
	"fmt"

	"github.com/miekg/dns"
)

func main() {
	// 设置 DNS 服务器地址和端口
	server := "8.8.8.8:53"

	// 设置查询的域名和类型
	domain := "www.google.com."
	qType := dns.TypeA

	// 创建一个 DNS 查询消息
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qType)
	m.RecursionDesired = true

	// 创建一个 DNS 客户端
	c := new(dns.Client)

	// 发送 DNS 查询请求
	r, _, err := c.Exchange(m, server)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 检查响应状态
	if r.Rcode != dns.RcodeSuccess {
		fmt.Printf("Failed to get an answer: %s\n", dns.RcodeToString[r.Rcode])
		return
	}

	// 打印查询结果
	for _, ans := range r.Answer {
		Arecord, ok := ans.(*dns.A)
		if ok {
			fmt.Printf("The A record of %s is: %s\n", domain, Arecord.A.String())
		}
	}
}
