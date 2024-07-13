package main

import (
	"codeSandbox/db"
	"codeSandbox/routes"
	"codeSandbox/utils/log"
)

/*

TODO
1. - [ ] 完善 docker 客户端初始化方式（生产模式不使用tcp）
4. bug: 使用验证码的token来提交代码，提示用户不存在
*/

func main() {
	db.InitRedis()
	log.ConfigLog()
	routes.Starter()

}
