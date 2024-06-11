package test

import (
	"codeSandbox/model/dto"
	sandbox "codeSandbox/service/sandbox"
	"testing"
)

func TestExecuteCode(t *testing.T) {
	codeRequest := dto.ExecuteCodeRequest{
		Code:      "package main\n\nimport (\n\t\"codeSandbox/routes\"\n\t\"codeSandbox/utils/log\"\n)\n\nfunc main() {\n\tlog.ConfigLog()\n\troutes.Starter()\n}\n",
		Language:  "Go",
		InputList: []string{"1", "2", "2"},
	}
	sandbox := sandbox.SandBox{
		DockerInfo: dto.DockerInfo{
			Language:   "Go",
			ImageName:  "go:1.21",
			Filename:   "Main.go",
			CompileCmd: "go build Main.go",
			RunCmd:     "./main",
		},
	}
	sandbox.ExecuteCode(&codeRequest)
	t.Log("test finish")
}
