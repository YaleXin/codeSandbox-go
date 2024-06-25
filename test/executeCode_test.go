package test

import (
	"bytes"
	"codeSandbox/model/dto"
	"codeSandbox/service/cryptoServices"
	"codeSandbox/service/sandboxDockerServices"
	"codeSandbox/utils"
	"encoding/json"
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
	for _, executeMsg := range response {
		assert.Equal(t, executeMsg.ExitCode, utils.EXIT_CODE_OK)
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
	for _, executeMsg := range response {
		assert.Equal(t, executeMsg.ExitCode, utils.EXIT_CODE_TIME_OUT)
		assert.Equal(t, executeMsg.ErrorMessage, sandboxDockerServices.ERR_MSG_TIME_OUT)
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
	for _, executeMsg := range response {
		assert.Equal(t, executeMsg.ExitCode, utils.EXIT_CODE_OK)
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
	for _, executeMsg := range response {
		assert.Equal(t, executeMsg.ExitCode, utils.EXIT_CODE_RUNTIME_ERROR)
		//assert.NotEqual(t, executeMsg.MemoryCost, uint64(0))
	}

}

func execode(code string, inputList []string) []dto.ExecuteMessage {
	codeRequest := dto.ExecuteCodeRequest{
		Code:      code,
		Language:  "Go",
		InputList: inputList,
	}
	sandbox := sandboxDockerServices.SandBox{
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

func TestProgramExecuteCode(t *testing.T) {
	//code := "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tvar a int\n\tvar b int\n\tscanf, err := fmt.Scanf(\"%d%d\", &a, &b)\n\tif err != nil {\n\t\tfmt.Println(scanf, err)\n\t}\n\tsum := a + b\n\tsliceNums := make([]int, 0, 0)\n\tfor i := 0; i < sum; i++ {\n\t\tsliceNums = append(sliceNums, i)\n\t}\n\tfmt.Println(sum)\n}\n"
	code := "gggggg"
	inputList := []string{"1 99999\n"}
	codeRequest := dto.ExecuteCodeRequest{
		Code:      code,
		Language:  "Go",
		InputList: inputList,
	}
	origidata, err := json.Marshal(codeRequest)
	if err != nil {
		t.Fatalf("json.Marshal %v", err)
	}
	publicKeyBase64 := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUE2M0UydVFUVVhFZTYzQXV1bEpPRQpOeFlJU0QvcElpamswWk9QTDRSWUFjMWhBRTlQelR0blBHUWJzMFRNZk1BRWI1WFM0bUNXYlgveHJyQVo1dlc4CnBieEorNXNMTFJjOEY0aXh0QUlOY2pxYTI2Mkh2R2JQOFNCbzFwdW54NWt3Z3M5b0tLN3M4R1h1ejZhT01STXEKSEJaYXhwVGtsdEo5c2NyTTlQUFhVSUFScEVpZHBqNDBiRU0rcE1nTGNQSDA5U1F6VE1WbjZ0RG9Fd05WRDNydwo0WWtJNWxYK2YwZi9WMFNVT3NrbUFvbk1aMGtnUVpuNDIwNWt6SFBvSXpGSEFTbmNhbG1vcGNRVk9NWnp1ZWZlCnh0ZGlrbmtlb1ZpUXZ0TVlsM0N0VXhibEMxUnJGQk1qZ243WVliUStSVjhKb01IOG8zQ2IrdTNCRU1IeXNWdE4KTFFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="

	var cryptoService cryptoServices.CryptoServiceInterface = new(cryptoServices.CryptoService)
	base64, err := cryptoService.EncryptWithPublicKeyBase64(publicKeyBase64, string(origidata))
	if err != nil {
		t.Fatalf("EncryptWithPublicKeyBase64 %v", err)
	}
	t.Logf("base64:%v", base64)
}
