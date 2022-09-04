package main

import (
        "sync"
        "github.com/BurntSushi/toml"
        "path/filepath"
)

type tomlConfig struct {
        Title string
        DB database `toml:"database"`
        Servers map[string]server
        Clients clients
}


type database struct {
        Server string
        Ports []int
        ConnMax int `toml:"connection_max"`
        Enabled bool
}

type server struct {
        IP string
        DC string
}

type clients struct {
        Data [][]interface{}
        Hosts []string
}

var (
        cfg * tomlConfig
        once sync.Once
        cfgLock = new(sync.RWMutex)
)

func Config() *tomlConfig {
        once.Do(ReloadConfig)
        cfgLock.RLock()
        defer cfgLock.RUnlock()
        return cfg
}

func ReloadConfig() {
        filePath, err := filepath.Abs(configPath)
        if err != nil {
                panic(err)
        }
        config := new(tomlConfig)
        if _ , err := toml.DecodeFile(filePath, config); err != nil {
                panic(err)
        }
        cfgLock.Lock()
        defer cfgLock.Unlock()
        cfg = config
}
