package utils

type ServerConfig struct {
	Server         Server         `yaml:"Server"`
	SandboxMachine SandboxMachine `yaml:"SandboxMachine"`
	DockerInfoList []DockerInfo   `yaml:"DockerInfoList"`
}

// 每种编程语言对应的配置信息
type DockerInfo struct {
	// 编程语言
	Language string `yaml:"Language"`
	// 对应的镜像
	ImageName string `yaml:"ImageName"`
	// 容器数量
	ContainerCount int `yaml:"ContainerCount"`
	// 保存的代码文件名
	Filename string `yaml:"Filename"`
	// 编译命令
	CompileCmd string `yaml:"CompileCmd"`
	// 运行命令
	RunCmd string `yaml:"RunCmd"`
}

type SandboxMachine struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}
type Server struct {
	AppMode  string   `yaml:"AppMode"`
	Host     string   `yaml:"Host"`
	Port     string   `yaml:"Port"`
	JwtKey   string   `yaml:"JwtKey"`
	Database Database `yaml:"Database"`
	Oss      Oss      `yaml:"Oss"`
	Push     Push     `yaml:"Push,omitempty"`
}
type Database struct {
	Type     string `yaml:"Type"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	PassWord string `yaml:"PassWord"`
	Name     string `yaml:"Name"`
	Redis    Redis  `yaml:"Redis"`
}
type Redis struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Password string `yaml:"Password"`
	Db       int    `yaml:"Db"`
}
type Oss struct {
	Name  string `yaml:"Name"`
	Qiniu struct {
		AccessKey string `yaml:"AccessKey,omitempty"`
		SecretKey string `yaml:"SecretKey,omitempty"`
		Bucket    string `yaml:"Bucket,omitempty"`
		Sever     string `yaml:"Sever,omitempty"`
	} `yaml:"qiniu"`
	Aliyun struct {
		AccessKeyID     string `yaml:"AccessKeyId,omitempty"`
		AccessKeySecret string `yaml:"AccessKeySecret,omitempty"`
		Endpoint        string `yaml:"Endpoint,omitempty"`
		BucketName      string `yaml:"BucketName,omitempty"`
	} `yaml:"aliyun"`
}
type Push struct {
	Enable string `yaml:"Enable,omitempty"`
	WxPush WxPush `yaml:"WxPush,omitempty"`
	Email  Email  `yaml:"Email,omitempty"`
}
type WxPush struct {
	CorpId  string `yaml:"CorpId,omitempty"`
	Agentid string `yaml:"Agentid,omitempty"`
	Secret  string `yaml:"Secret,omitempty"`
}
type Email struct {
	To       string `yaml:"To,omitempty"`
	Password string `yaml:"Password,omitempty"`
	From     string `yaml:"From,omitempty"`
	Host     string `yaml:"Host,omitempty"`
	Port     string `yaml:"Port,omitempty"`
}
