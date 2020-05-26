package main

import (
	"fmt"
	"github.com/juju/ratelimit"
	"sync"
	"time"
)

// RateLimitingPlugin can limit connecting per unit time
type RateLimitingPlugin struct {
	FillInterval time.Duration
	Capacity     int64
	bucket       *ratelimit.Bucket
}

// NewRateLimitingPlugin creates a new RateLimitingPlugin
func NewRateLimitingPlugin(fillInterval time.Duration, capacity int64) *RateLimitingPlugin {
	tb := ratelimit.NewBucket(fillInterval, capacity)

	return &RateLimitingPlugin{
		FillInterval: fillInterval,
		Capacity:     capacity,
		bucket:       tb}
}

// HandleConnAccept can limit connecting rate
func (plugin *RateLimitingPlugin) HandleConnAccept() bool {
	//return conn, plugin.bucket.TakeAvailable(1) > 0
	return plugin.bucket.TakeAvailable(1) > 0
}

func (plugin *RateLimitingPlugin) Consumer(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ret := plugin.HandleConnAccept()
		if ret {
			fmt.Println("get token..run")
		} else {
			fmt.Println("Not get token..quit")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func main() {
	wg := sync.WaitGroup{}
	buc := NewRateLimitingPlugin(200*time.Millisecond, 5)
	wg.Add(1)
	go buc.Consumer(&wg)
	wg.Wait()
}
