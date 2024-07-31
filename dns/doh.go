package main

import (
        "encoding/base64"
        "fmt"
        "github.com/miekg/dns"
        "io/ioutil"
        "net/http"
        "os"
)

func main() {
       query := dns.Msg{}
       query.SetQuestion("www.taobao.com.", dns.TypeA)
       msg, _ := query.Pack()
       b64 := base64.RawURLEncoding.EncodeToString(msg)
       resp, err := http.Get("https://dns.alidns.com/dns-query?dns=" + b64)
       if err != nil {
            fmt.Printf("Send query error, err:%v\n", err)
            os.Exit(1)
       }
       defer resp.Body.Close()
       bodyBytes, _ := ioutil.ReadAll(resp.Body)
       response := dns.Msg{}
       response.Unpack(bodyBytes)
       fmt.Printf("Dns answer is :%v\n", response.String())
}
