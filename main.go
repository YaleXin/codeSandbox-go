package main

import (
	"codeSandbox/routes"
	"codeSandbox/utils/log"
)

/*

TODO
1. 完善 docker 客户端初始化方式（生产模式不使用tcp）
2. 开放平台接口 accesskey 和 secretkey 设计
3. 用户功能
4. 接口访问记录入库


*/

func main() {
	log.ConfigLog()
	routes.Starter()
}
