package sandbox

import (
	"codeSandbox/model/dto"
	utilsType "codeSandbox/utils"
	filesUtils "codeSandbox/utils/files"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const CODE_LOCAL_DIR_PREX string = "temp"

// 每个执行用例执行最大的时间，单位为秒
const RUN_CODE_TIME_OUT = 5 * time.Second

func clearFile(codeFilename string) {
	err := os.Remove(codeFilename)
	if err != nil {
		log.Errorf("Remove code file fail:%v", err)
	}
	log.Debugf("clear file finish, codeFilename:%v", codeFilename)
}
func getOutputResponse(executeMessageArrayList []dto.ExecuteMessage) dto.ExecuteCodeResponse {
	response := dto.ExecuteCodeResponse{}
	for _, executeMessage := range executeMessageArrayList {
		response.ExecuteMessages = append(response.ExecuteMessages, executeMessage)
	}
	return response
}

func copyFileToContainer(containerId, userCodeFilePath, uuid string) bool {
	//======== 容器中先创建文件夹，然后本地打包文件上传至文件夹
	ctx := context.Background()
	sourceFiles := []string{userCodeFilePath}
	tarFilePath := "main.tar"
	destFilePath := WORDING_DIR + string(filepath.Separator) + uuid

	message := runCmdByContainer(containerId, []string{"mkdir", "-p", uuid}, "", "", "mkdir")
	if message.ExitCode == EXIT_CODE_ERROR {
		return false
	}
	// 将代码文件打包为 main.tar
	err := filesUtils.CreateTarArchiveFiles(sourceFiles, tarFilePath)
	if err != nil {
		log.Errorf("create tar file fail: %v", err)
	}
	srcFile, err := os.Open(tarFilePath)
	// 先close 再删除
	defer os.Remove(tarFilePath)
	defer srcFile.Close()
	err = DockerClient.CopyToContainer(ctx, containerId, destFilePath, srcFile, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if err != nil {
		log.Errorf("copy to containerId:%v fail:%v", containerId, err)
		return false
	}
	return true
}

// 将本地文件保存到容器中，并编译运行
func (sandbox *SandBox) compileAndRun(language string, userCodeFilePath string, inputList []string) []dto.ExecuteMessage {
	dockerInfo := sandbox.DockerInfo
	// 有多个容器可以选择时，随机抽一个进行使用
	// 设置随机数种子，通常使用时间作为种子以获得更好的随机性
	rand.Seed(time.Now().UnixNano())
	count := dockerInfo.ContainerCount
	selectIdx := rand.Intn(count)
	log.Infof("selectIdx:%v", selectIdx)
	containerId := getContainerId(dockerInfo, selectIdx)

	//======== 复制文件（先打包文件，再复制到容器中）
	// 对于 temp\Go\81b6f397-a185-4ef2-b3c4-908c3ad4d20c\Main.go uuid = 81b6f397-a185-4ef2-b3c4-908c3ad4d20c
	uuid := filepath.Base(filepath.Dir(userCodeFilePath))
	copyStatus := copyFileToContainer(containerId, userCodeFilePath, uuid)
	if !copyStatus {
		return []dto.ExecuteMessage{{
			ExitCode:     EXIT_CODE_ERROR,
			ErrorMessage: "System error",
		}}
	}

	//====== 编译文件
	compileCmd := dockerInfo.CompileCmd
	cmdSplit := strings.Split(compileCmd, " ")
	// Linux系统下，路径分隔符必然为 /
	workDir := WORDING_DIR + "/" + uuid
	compileRes := runCmdByContainer(containerId, cmdSplit, workDir, "", "compile")
	log.Infof("compileRes:%v", compileRes)
	if compileRes.ExitCode == EXIT_CODE_ERROR {
		compileRes.ErrorMessage = "Compile fail"
		return []dto.ExecuteMessage{compileRes}
	}

	//====== 运行代码
	messages := runCode(containerId, dockerInfo, inputList, workDir)

	return messages
}

func runCode(containerId string, dockerInfo utilsType.DockerInfo, inputList []string, workDir string) []dto.ExecuteMessage {
	messages := make([]dto.ExecuteMessage, 0, 0)
	runCmd := dockerInfo.RunCmd
	runSplit := strings.Split(runCmd, " ")
	for _, inputStr := range inputList {
		runRes := runCmdByContainer(containerId, runSplit, workDir, inputStr, "run")
		messages = append(messages, runRes)
	}
	return messages
}
func (sandbox *SandBox) saveFile(code string) (fs.File, string) {
	// 不同的编程语言将会保存到不同的地方
	language := sandbox.DockerInfo.Language
	filename := sandbox.DockerInfo.Filename
	newUUID, err := uuid.NewRandom()
	// 例如父级路径为 temp/Go/uuid/
	parentPath := CODE_LOCAL_DIR_PREX + string(filepath.Separator) + language + string(filepath.Separator) + newUUID.String()
	// 限为 0666，表示为所有人都可以对该文件夹进行读写，且不存在时会自动创建。
	err = os.MkdirAll(parentPath, 0666)
	if err != nil {
		log.Errorf("MkdirAll %v fail:%v", parentPath, err)
		return nil, ""
	}
	if err != nil {
		log.Errorf("Gennerate UUID fail: %v", err)
		return nil, ""
	}
	codeFilename := parentPath + string(filepath.Separator) + filename
	// O_WRONLY 以只写的模式打开文件, O_CREATE 如果文件不存在则创建文件
	file, err := os.OpenFile(codeFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Errorf("OpenFile %v fail: %v", codeFilename, err)
		return nil, ""
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
		return nil, ""
	}
	log.Debugf("save file finish, file:%v, codeFilename:%v", file, codeFilename)
	return file, codeFilename
}

func (sandbox *SandBox) ExecuteCode(executeCodeRequest *dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	// 1. 保存用户代码为文件
	code := executeCodeRequest.Code
	_, codeFilePath := sandbox.saveFile(code)
	// 4. 文件清理（暂时不清理了，留存）
	// defer clearFile(codeFilename)
	// 2. 编译代码并执行代码
	language := executeCodeRequest.Language
	inputList := executeCodeRequest.InputList
	executeMessageArrayList := sandbox.compileAndRun(language, codeFilePath, inputList)
	// 3. 整理输出信息
	executeCodeResponse := getOutputResponse(executeMessageArrayList)
	return executeCodeResponse
}
