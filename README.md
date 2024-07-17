# codeSandbox
代码沙箱平台
## 技术栈
- Golang
- Gin
- Swagger
## 自动生成API文档
```shell
swag init --parseDependency --parseInternal --parseGoList=false --parseDepth=1
```
## 构建为Linux应用
```shell
go env -w GOOS=linux
go build -o codeSandbox
```
## 在目标平台上运行
```shell
export CodeSandboxConfigFileName=/path/to/config.yml
./codeSandbox
```
## 切换回 Windows
```shell
go env -w GOOS=windows
```