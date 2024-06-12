package test

import (
	"codeSandbox/model/dto"
	"codeSandbox/service/sandbox"
	"codeSandbox/utils"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestExecuteCode(t *testing.T) {

	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	codeRequest := dto.ExecuteCodeRequest{
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a int\n\tvar b int\n\tscanf, err := fmt.Scanf(\"%d%d\", &a, &b)\n\tif err != nil {\n\t\tfmt.Println(scanf, err)\n\t}\n\tfmt.Println(a + b)\n}",
		Language:  "Go",
		InputList: []string{"8", "2"},
	}
	sandbox := sandbox.SandBox{
		DockerInfo: utils.DockerInfo{
			Language:       "Go",
			ImageName:      "golang:1.17",
			Filename:       "Main.go",
			CompileCmd:     "go build Main.go",
			RunCmd:         "./main",
			ContainerCount: 3,
		},
	}
	runRes := sandbox.ExecuteCode(&codeRequest)
	logrus.Debugf("runRes:%v", runRes)
	t.Log("test finish")

}
