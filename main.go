package main

import (
	"codeSandbox/db"
	"codeSandbox/routes"
	"codeSandbox/utils/log"
)

/*

TODO
1. - [ ] 完善 docker 客户端初始化方式（生产模式不使用tcp）
2. - [ ] 重要接口防抖
3. - [ ] key 个数限制 输入用例个数限制
*/

func main() {
	db.InitRedis()
	log.ConfigLog()
	routes.Starter()

}
