package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout     int = 90
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},

		Timeout: 20 * time.Second,
	}
	return client
}

func main() {
	var endPoint string = "http://127.0.0.1:12345/post"

	for {
		req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer([]byte("Post this data")))
		if err != nil {
			log.Fatalf("Error Occured. %+v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// use httpClient to send request
		response, err := httpClient.Do(req)
		if err != nil && response == nil {
			log.Fatalf("Error sending request to API endpoint. %+v", err)
		} else {
			// Close the connection to reuse it
			defer response.Body.Close()

			// Let's check if the work actually is done
			// We have seen inconsistencies even when we get 200 OK response
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("Couldn't parse response body. %+v", err)
			}

			log.Println("Response Body:", string(body))
		}
		time.Sleep(1 * time.Second)
	}

}
