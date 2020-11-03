package main

import (
	"context"
	"fmt"
	"time"
)

type HookFunc func(ctx context.Context, call string, customName, funcName string) func(err error)

type WrapperClass struct {
	addr     string
	hooks    []HookFunc
	hookname []string
}

func (c *WrapperClass) invokeHook(ctx context.Context, call string, customName string) func(error) {
	finishHooks := make([]func(error), 0, len(c.hooks))
	for index, fn := range c.hooks {
		fmt.Printf("\n[invokeHook]range every fn:%v\n", fn)
		// 将实际的参数，绑定到中间件上，同时fn的结果添加到finishHooks中
		finishHooks = append(finishHooks, fn(ctx, call, customName, c.hookname[index]))
	}
	// 返回一个函数
	return func(err error) {
		// 遍历 finishHooks，依次运行各个中间件
		for _, fn := range finishHooks {
			fmt.Printf("\n[invokeHook]range fn result:%v,err=%v\n", fn, err)
			fn(err)
		}
	}
}

// AddHook add hook function.
func (c *WrapperClass) AddHook(hookFn HookFunc) {
	c.hooks = append(c.hooks, hookFn)
}

func NewWrapperClass(addr string) *WrapperClass {
	wc := &WrapperClass{
		addr:  addr,
		hooks: make([]HookFunc, 0),
	}
	wc.AddHook(NewSlowLogHook(250 * time.Millisecond))
	wc.hookname = append(wc.hookname, "NewSlowLogHook")
	wc.AddHook(NewTestLogHook())
	wc.hookname = append(wc.hookname, "AddHook")
	return wc
}

//
// NewSlowLogHook log slow operation.
func NewSlowLogHook(threshold time.Duration) HookFunc {
	return func(ctx context.Context, call string, customName, funcName string) func(err error) {
		start := time.Now()
		return func(error) {
			duration := time.Since(start)
			if duration < threshold {
				//return
			}
			fmt.Printf("slow log test1: %s %s time: %s[funcname:%s]\n", customName, call, duration, funcName)
		}
	}
}

func NewTestLogHook() HookFunc {
	return func(ctx context.Context, call string, customName, funcName string) func(err error) {
		return func(error) {
			fmt.Printf("print log test2: %s %s[func:%s]\n", customName, call, funcName)
		}
	}
}

func (c *WrapperClass) DoHookAction(ctx context.Context) (err error) {
	fmt.Println("[DoHookAction]start")
	finishHook := c.invokeHook(ctx, "hookstruct", "GET")
	fmt.Println("[DoHookAction]do real origin lib jobs...")
	time.Sleep(1 * time.Second)
	finishHook(err)
	fmt.Println("[DoHookAction]end")
	return
}

func main() {
	wc := NewWrapperClass("127.0.0.1")
	wc.DoHookAction(context.Background())
}
