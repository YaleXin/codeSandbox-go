package main

import (
	"codeSandbox/routes"
	"codeSandbox/utils/log"
)

/*

TODO
1. - [ ] 完善 docker 客户端初始化方式（生产模式不使用tcp）

*/

func main() {
	log.ConfigLog()
	routes.Starter()
}
