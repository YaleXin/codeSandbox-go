package service_v1

import (
	dto "codeSandbox/model/dto"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var DockerInfoList dto.DockerInfoList

func init() {
	var file []byte
	var err error
	// 通过环境变量来判断是否使用默认配置文件，方便开发
	if filename, ok := os.LookupEnv("CodeSandboxDockerConfigFileName"); ok {
		file, err = os.ReadFile(filename)
	} else {
		file, err = os.ReadFile("./conf/machine.yml")
	}
	if err != nil {
		panic(fmt.Sprintf("DockerConfig file read fail: %v1", err))
	}
	err = yaml.Unmarshal(file, &DockerInfoList)
	if err != nil {
		panic(fmt.Sprintf("DockerConfig parse fail: %v1", err))
	}
}

func GetSupportLanguages() []string {
	var languages []string
	for _, info := range DockerInfoList.DockerInfoList {
		languages = append(languages, info.Language)
	}
	return languages
}

func ExecuteCode(executeCodeRequest dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	e := dto.ExecuteCodeResponse{}
	return e
}
