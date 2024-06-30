package main

import (
	"codeSandbox/routes"
	"codeSandbox/utils/log"
)

/*

TODO
1. - [ ] 完善 docker 客户端初始化方式（生产模式不使用tcp）
2. - [x] 开放平台接口 accesskey 和 secretkey 设计
3. - [x] 用户功能
4. - [x] 接口访问记录入库
5. - [x] 限制代码的读写权限
6. 设置总的超时记录（不仅仅是单个输入用例）
7. - [x] 容器中的命令也要有超时控制（即宿主机上面的协程有超时控制，同时容器运行示例也要有超时控制）
8. - [x] 要求客户端使用 secretKey 计算 payload 哈希


*/

func main() {
	log.ConfigLog()
	routes.Starter()
}
