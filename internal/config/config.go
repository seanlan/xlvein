package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"math/rand"
	"os"
	"time"
)

var C *Config

func setup() {
	initLogging()
}

func init() {
	var err error
	var s Config
	// 初始化随机种子
	rand.Seed(time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
	// 初始化配置
	path := "conf.yaml"
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, &s)
	if err != nil {
		log.Fatal(err)
	}
	C = &s
	setup()
}
