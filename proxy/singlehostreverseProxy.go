package main

// 测试NewSingleHostReverseProxy的使用

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	remote, err := url.Parse("http://google.com")
	if err != nil {
		panic(err)
	}

	remote2, err := url.Parse("http://www.163.com/aaa")
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy2 := httputil.NewSingleHostReverseProxy(remote2)

	var maper map[string]*httputil.ReverseProxy = make(map[string]*httputil.ReverseProxy)
	fmt.Println(remote.String(), remote2.String())
	maper[remote.String()] = proxy
	maper[remote2.String()] = proxy2

	// 设置前置路由及其对应的代理处理对象
	//http.HandleFunc("/", handler(proxy))
	//http.HandleFunc("/test", handler(proxy2))
	http.HandleFunc("/", handler(maper))
	http.HandleFunc("/test", handler(maper))
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		panic(err)
	}
}

func randMapKey(m map[string]*httputil.ReverseProxy) *httputil.ReverseProxy {
	mapKeys := make([]string, 0, len(m)) // pre-allocate exact size
	for key := range m {
		mapKeys = append(mapKeys, key)
	}
	return m[mapKeys[rand.Intn(len(mapKeys))]]
}

// 在handler方法中，1.选择处理前置路由的处理对象 2. 设置后置的路由（是否修改）及各种http头部属性，体现在w上 3. 对修改后的的w，调用httputil.ReverseProxy的ServeHTTP方法p.ServeHTTP(w, r)完成整个代理过程
// 以上3步即完成对proxy的前后路由绑定过程
func handler(maper map[string]*httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	//p := maper["http://www.163.com/aaa"]
	p := randMapKey(maper)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-Ben", "Rad")
		p.ServeHTTP(w, r)
	}
}
