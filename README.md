# codeSandbox
代码沙箱平台。

- 用户代码运行在容器层，增加安全性
- 基于配置 yaml 文件的设计思想，无需 修改代码即可添加新的编程语言
- 提供加密式开放 API 功能，使得用户借助 Access key 提交代码而无需维持登录信息。

## 技术栈
- Golang
- Gin
- Swagger
- GORM
- Redis
- Gomail
## 自动生成API文档
```shell
swag init --parseDependency --parseInternal --parseGoList=false --parseDepth=1
```
## 构建为Linux应用
```shell
go env -w GOARCH=amd64
go env -w GOOS=linux
go build -o codeSandbox
```
## 在目标平台上运行
```shell
export CodeSandboxConfigFileName=/path/to/config.yml
./codeSandbox
```

## 文件夹说明

```shell
.
|-- Dockerfile-java11
|-- LICENSE
|-- README.md
|-- api # controller层
|-- conf # 配置文件
|-- db # 数据库相关
|-- docs # swagger文档生成目录
|-- go.mod # 包管理相关
|-- go.sum # 包管理相关
|-- log # 日志相关
|-- main.go # 主启动文件
|-- model # 模型相关
|-- responses # 自定义响应体
|-- routes # 路由相关
|-- service # 服务层
|-- task # 自动任务
|-- temp # 代码存放中间目录
|-- test # 测试相关
`-- utils # 工具类相关
```

