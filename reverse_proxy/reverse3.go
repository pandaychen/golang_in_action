package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// 地址重写实例
	// http://127.0.0.1:8888/test?id=1  =》 http://127.0.0.1:8081/reverse/test?id=1

	rs1 := "http://127.0.0.1:8081/reverse"
	targetUrl, err := url.Parse(rs1)
	if err != nil {
		log.Fatal("err")
	}
	proxy := NewSingleHostReverseProxy(targetUrl)
	log.Println("Reverse proxy server serve at : 127.0.0.1:8888")
	if err := http.ListenAndServe(":8888", proxy); err != nil {
		log.Fatal("Start server failed,err:", err)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	// 自定义ModifyResponse
	modifyResp := func(resp *http.Response) error {
		var oldData, newData []byte
		oldData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// 根据不同状态码修改返回内容
		if resp.StatusCode == 200 {
			newData = []byte("[INFO] " + string(oldData))

		} else {
			newData = []byte("[ERROR] " + string(oldData))
		}

		// 修改返回内容及ContentLength
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(newData))
		resp.ContentLength = int64(len(newData))
		resp.Header.Set("Content-Length", fmt.Sprint(len(newData)))
		return nil
	}
	// 传入自定义的ModifyResponse
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyResp}
}
