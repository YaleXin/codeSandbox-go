package test

import (
	"codeSandbox/model/dto"
	"codeSandbox/service/sandbox"
	"codeSandbox/utils"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestExecuteCode(t *testing.T) {
	bytes, _ := os.ReadFile("C:\\Users\\Yalexin\\GolandProjects\\codeSandbox\\inner_test\\demo\\main.go")
	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	codeRequest := dto.ExecuteCodeRequest{
		Code:      string(bytes),
		Language:  "Go",
		InputList: []string{"1 9999999 \n"},
	}
	sandbox := sandbox.SandBox{
		DockerInfo: utils.DockerInfo{
			Language:       "Go",
			ImageName:      "golang:1.17",
			Filename:       "Main.go",
			CompileCmd:     "go build Main.go",
			RunCmd:         "./Main",
			ContainerCount: 1,
		},
	}
	runRes := sandbox.ExecuteCode(&codeRequest)
	logrus.Debugf("runRes:%v", runRes)

	for _, res := range runRes.ExecuteMessages {
		if res.MemoryCost == 0 {
			t.Errorf("res.MemoryCost is 0")
		}
	}

	t.Log("test finish")

}
