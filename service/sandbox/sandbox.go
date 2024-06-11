package sandbox

import (
	"codeSandbox/model/dto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
)

const CODE_DIR_PREX string = "temp"

type SandBox struct {
	DockerInfo dto.DockerInfo
}

func clearFile(codeFile fs.File) {

}
func getOutputResponse(executeMessageArrayList []dto.ExecuteMessage) dto.ExecuteCodeResponse {
	return dto.ExecuteCodeResponse{}
}
func compileAndRun(language string, userCodeFile fs.File, inputList []string) []dto.ExecuteMessage {
	return nil
}

func (sandbox *SandBox) saveFile(code string) fs.File {
	// 不同的编程语言将会保存到不同的地方
	language := sandbox.DockerInfo.Language
	filename := sandbox.DockerInfo.Filename
	parentPath := CODE_DIR_PREX + string(filepath.Separator) + language
	// 限为 0666，表示为所有人都可以对该文件夹进行读写，且不存在时会自动创建。
	err := os.MkdirAll(parentPath, 0666)
	if err != nil {
		log.Errorf("MkdirAll %v fail:%v", parentPath, err)
		return nil
	}
	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Errorf("Gennerate UUID fail: %v", err)
		return nil
	}
	codeFilename := parentPath + string(filepath.Separator) + newUUID.String() + "_" + filename
	file, err := os.OpenFile(codeFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Errorf("OpenFile %v fail: %v", codeFilename, err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Close fail: %v", err)
		}
	}(file)
	_, err = file.WriteString(code)
	if err != nil {
		log.Errorf("WriteString fail: %v", err)
		return nil
	}
	return file
}

func (sandbox *SandBox) ExecuteCode(executeCodeRequest *dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	// 1. 保存用户代码为文件
	code := executeCodeRequest.Code
	codeFile := sandbox.saveFile(code)
	// 2. 编译代码并执行代码
	language := executeCodeRequest.Language
	inputList := executeCodeRequest.InputList
	executeMessageArrayList := compileAndRun(language, codeFile, inputList)
	// 3. 整理输出信息
	executeCodeResponse := getOutputResponse(executeMessageArrayList)
	// 4. 文件清理
	clearFile(codeFile)
	return executeCodeResponse
}
