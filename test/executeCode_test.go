package test

import (
	"bytes"
	"codeSandbox/model/dto"
	"codeSandbox/service/sandbox"
	"codeSandbox/utils"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
	"testing"
)

// CustomFormatter 自定义日志格式化器
type CustomFormatter struct{}

// Format 实现Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 定义时间格式，包含毫秒
	timestamp := entry.Time.Format("2024-06-18 15:04:05.000")

	// 构建自定义的日志格式
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	fmt.Fprintf(b, "%s [%s] %s\n", timestamp, entry.Level.String(), entry.Message)

	// 添加额外的fields，如果有的话
	if len(entry.Data) > 0 {
		for k, v := range entry.Data {
			fmt.Fprintf(b, "    %s=%v\n", k, v)
		}
	}

	return b.Bytes(), nil
}

// 有睡眠
func TestExecuteLowCode(t *testing.T) {
	code := "package main\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nfunc main() {\n\tvar a int\n\tvar b int\n\tscanf, err := fmt.Scanf(\"%d%d\", &a, &b)\n\tif err != nil {\n\t\tfmt.Println(scanf, err)\n\t}\n\tsum := a + b\n\tsliceNums := make([]int, 0, 0)\n\tfor i := 0; i < sum; i++ {\n\t\tsliceNums = append(sliceNums, i)\n\t}\n\ttime.Sleep(3 * time.Second)\n\tfmt.Println(sum)\n}\n"
	inputList := []string{"1 99999\n"}
	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&CustomFormatter{})
	response := execode(code, inputList)
	logrus.Debugf("runRes:%v", response)
	for _, executeMsg := range response.ExecuteMessages {
		assert.Equal(t, executeMsg.ExitCode, sandbox.EXIT_CODE_OK)
		assert.NotEqual(t, executeMsg.MemoryCost, uint64(0))
	}

}

func TestExecuteTimeoutCode(t *testing.T) {
	code := "package main\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nfunc main() {\n\tvar a int\n\tvar b int\n\tscanf, err := fmt.Scanf(\"%d%d\", &a, &b)\n\tif err != nil {\n\t\tfmt.Println(scanf, err)\n\t}\n\tsum := a + b\n\tsliceNums := make([]int, 0, 0)\n\tfor i := 0; i < sum; i++ {\n\t\tsliceNums = append(sliceNums, i)\n\t}\n\ttime.Sleep(10 * time.Second)\n\tfmt.Println(sum)\n}\n"
	inputList := []string{"1 99999\n"}
	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&CustomFormatter{})
	response := execode(code, inputList)
	logrus.Debugf("runRes:%v", response)
	for _, executeMsg := range response.ExecuteMessages {
		assert.Equal(t, executeMsg.ExitCode, sandbox.EXIT_CODE_ERROR)
		assert.Equal(t, executeMsg.ErrorMessage, sandbox.ERR_MSG_TIME_OUT)
		assert.Equal(t, executeMsg.Message, "")
	}

}

// 无睡眠
func TestExecuteFastCode(t *testing.T) {
	code := "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tvar a int\n\tvar b int\n\tscanf, err := fmt.Scanf(\"%d%d\", &a, &b)\n\tif err != nil {\n\t\tfmt.Println(scanf, err)\n\t}\n\tsum := a + b\n\tsliceNums := make([]int, 0, 0)\n\tfor i := 0; i < sum; i++ {\n\t\tsliceNums = append(sliceNums, i)\n\t}\n\tfmt.Println(sum)\n}\n"
	inputList := []string{"1 99999\n"}
	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&CustomFormatter{})
	response := execode(code, inputList)
	logrus.Debugf("runRes:%v", response)
	for _, executeMsg := range response.ExecuteMessages {
		assert.Equal(t, executeMsg.ExitCode, sandbox.EXIT_CODE_OK)
		//assert.NotEqual(t, executeMsg.MemoryCost, uint64(0))
	}

}

// 除 0 代码
func TestExecuteDiveZeroCode(t *testing.T) {
	code := "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tvar a int\n\tvar b int\n\tfmt.Scanf(\"%d%d\", &a, &b)\n\tsum := a + b\n\tfmt.Println(a / b)\n\tsliceNums := make([]int, 0, 0)\n\tfor i := 0; i < sum; i++ {\n\t\tsliceNums = append(sliceNums, i)\n\t}\n\tfmt.Println(sum)\n}\n"
	inputList := []string{"bbbb aaa\n"}
	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&CustomFormatter{})
	response := execode(code, inputList)
	logrus.Debugf("runRes:%v", response)
	for _, executeMsg := range response.ExecuteMessages {
		assert.Equal(t, executeMsg.ExitCode, sandbox.EXIT_CODE_ERROR)
		//assert.NotEqual(t, executeMsg.MemoryCost, uint64(0))
	}

}

func execode(code string, inputList []string) dto.ExecuteCodeResponse {
	codeRequest := dto.ExecuteCodeRequest{
		Code:      code,
		Language:  "Go",
		InputList: inputList,
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

	return runRes
}
