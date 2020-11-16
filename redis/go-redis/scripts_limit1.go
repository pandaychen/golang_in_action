package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"sync"
)

func createScript() *redis.Script {
	script := redis.NewScript(`
	-- 获取调用脚本时传入的第一个key值（用作限流的 key）
local key = KEYS[1]
-- 获取调用脚本时传入的第一个参数值（限流大小）
local limit = tonumber(ARGV[1])

local ttl = tonumber(ARGV[2])

-- 获取当前流量大小
local curentLimit = tonumber(redis.call('get', key) or "0")

-- 是否超出限流
if curentLimit + 1 > limit then
    -- 返回(拒绝)
    return 0
else
    -- 没有超出 value + 1
    redis.call('INCRBY', key, 1)
    -- 设置过期时间
    redis.call('EXPIRE', key, ttl)
    -- 返回(放行)
    return 1
end

	`)
	return script
}

func evalScript(client *redis.Client, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	script := createScript()
	sha, err := script.Load(client.Context(), client).Result()
	if err != nil {
		log.Fatalln(err)
	}
	ret := client.EvalSha(client.Context(), sha, []string{
		"limiterkey",
	}, 1, 60)
	if result, err := ret.Result(); err != nil {
		log.Fatalf("Execute Redis fail: %v", err.Error())
	} else {
		fmt.Printf("id: %d, result: %d\n", id, result)
	}
}

func main() {
	var wg sync.WaitGroup
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go evalScript(client, i, &wg)
	}
	wg.Wait()

}
