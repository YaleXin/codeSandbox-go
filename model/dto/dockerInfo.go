package dto

type DockerInfoList struct {
	DockerInfoList []DockerInfo `yaml:"DockerInfoList"`
}

// 每种编程语言对应的配置信息
type DockerInfo struct {
	// 编程语言
	Language string `yaml:"Language"`
	// 对应的镜像
	ImageName string `yaml:"ImageName"`
	// 保存的代码文件名
	Filename string `yaml:"Filename"`
	// 编译命令
	CompileCmd string `yaml:"CompileCmd"`
	// 运行命令
	RunCmd string `yaml:"RunCmd"`
}
