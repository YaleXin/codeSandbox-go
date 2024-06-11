package utils

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	Config ServerConfig
)

// 全局唯一一个 config 变量
func init() {
	var file []byte
	var err error
	// 通过环境变量来判断是否使用默认配置文件，方便开发
	if filename, ok := os.LookupEnv("CodeSandboxConfigFileName"); ok {
		file, err = os.ReadFile(filename)
	} else {
		file, err = os.ReadFile("./conf/config.yml")
	}
	if err != nil {
		log.Panicf("config file read fail : %v", err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Panicf("config file parse fail : %v", err)
	}
}
