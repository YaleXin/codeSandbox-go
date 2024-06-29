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
5. 限制代码的读写权限
6. 验证 jwt 中的用户是在数据库中真实存在的（避免用户自己伪造）
7. 设置总的超时记录（不仅仅是单个输入用例）
8. 超时后，主动停止监控协程
9. 要求客户端使用 secretKey 计算 payload 哈希


*/

func main() {
	log.ConfigLog()
	routes.Starter()
}
