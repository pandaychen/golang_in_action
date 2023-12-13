package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	sync.RWMutex
	Name    string
	Version float64
}

var appConfig *AppConfig

func main() {
	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	appConfig = &AppConfig{}
	updateConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		updateConfig()
	})

	// 模拟多个并发读取配置的 goroutine
	for i := 0; i < 10; i++ {
		go func() {
			for {
				printConfig()
			}
		}()
	}

	// 阻塞主 goroutine，以便其他 goroutine 可以继续运行
	select {}
}

func updateConfig() {
	appConfig.Lock()
	defer appConfig.Unlock()

	appConfig.Name = viper.GetString("app.name")
	appConfig.Version = viper.GetFloat64("app.version")
}

func printConfig() {
	appConfig.RLock()
	defer appConfig.RUnlock()

	fmt.Printf("App Name: %s, Version: %f\n", appConfig.Name, appConfig.Version)
}
