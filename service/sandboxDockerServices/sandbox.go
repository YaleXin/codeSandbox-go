package sandboxDockerServices

import (
	"codeSandbox/model/dto"
	utilsType "codeSandbox/utils"
	filesUtils "codeSandbox/utils/files"
	"context"
	"fmt"
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
const ERR_MSG_TIME_OUT string = "TIMEOUT"

// 每个执行用例执行最大的时间，单位为秒
const RUN_CODE_TIME_OUT = 5 * time.Second

// 超时外开销
const TIMEOUT_OVERHEAD = time.Second * 2

// timemout 命令执行指定命令后的发生超时的返回码
const TIMEOUT_CMD_EXITCODE = 124

func clearFile(codeFilename string) {

	dir := filepath.Dir(codeFilename)
	err := os.RemoveAll(dir)
	if err != nil {
		log.Errorf("Remove dir %v fail:%v", dir, err)
	}
	log.Debugf("clear file finish, dir :%v", dir)
}
func getOutputMessage(executeMessageArrayList []dto.ExecuteMessage) []dto.ExecuteMessage {
	executeMessages := make([]dto.ExecuteMessage, 0, 0)
	for _, executeMessage := range executeMessageArrayList {
		executeMessages = append(executeMessages, executeMessage)
	}
	return executeMessages
}

func copyFileToContainer(containerId, userCodeFilePath, uuid string) bool {
	//======== 容器中先创建文件夹，然后本地打包文件上传至文件夹
	ctx := context.Background()
	sourceFiles := []string{userCodeFilePath}
	tarFilePath := "main.tar"
	destFilePath := WORDING_DIR + string(filepath.Separator) + uuid

	message := runCmdByContainer(containerId, []string{"mkdir", "-p", uuid}, "", "", "mkdir", "")
	if message.ExitCode != utilsType.EXIT_CODE_OK {
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

	// 规定该新目录只能给该用户读写（root除外）
	newUserName := strings.Replace(uuid, "-", "", -1)
	message = runCmdByContainer(containerId, []string{"useradd", newUserName, "-m"}, "", "", "useradd", "")
	if message.ExitCode != utilsType.EXIT_CODE_OK {
		return false
	}
	// 更改新建的目录归属权为新建的用户
	cmds := []string{"chown", "-R", fmt.Sprintf("%v:%v", newUserName, newUserName), uuid}
	message = runCmdByContainer(containerId, cmds, "", "", "chown", "")
	if message.ExitCode != utilsType.EXIT_CODE_OK {
		return false
	}
	cmds = []string{"chmod", "-R", "700", uuid}
	message = runCmdByContainer(containerId, cmds, "", "", "chmod", "")
	if message.ExitCode != utilsType.EXIT_CODE_OK {
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
	newUserName := strings.Replace(uuid, "-", "", -1)
	copyStatus := copyFileToContainer(containerId, userCodeFilePath, uuid)
	if !copyStatus {
		return []dto.ExecuteMessage{{
			ExitCode:     utilsType.EXIT_CODE_BASE_ERROR,
			ErrorMessage: "System error",
		}}
	}
	// 对容器中刚刚创建的目录和用户执行删除操作
	defer clearContainerFileAndUser(containerId, uuid)

	//====== 编译文件
	compileCmd := dockerInfo.CompileCmd
	cmdSplit := strings.Split(compileCmd, " ")
	// Linux系统下，路径分隔符必然为 /
	workDir := WORDING_DIR + "/" + uuid
	compileRes := runCmdByContainer(containerId, cmdSplit, workDir, "", "compile", newUserName)
	log.Infof("compileRes:%v", compileRes)
	if compileRes.ExitCode != utilsType.EXIT_CODE_OK {
		compileRes.ExitCode = utilsType.EXIT_CODE_COMPILE_ERROR
		compileRes.ErrorMessage = "Compile fail"
		return []dto.ExecuteMessage{compileRes}
	}

	//====== 运行代码
	messages := runCode(containerId, dockerInfo, inputList, workDir, newUserName)

	return messages
}

func clearContainerFileAndUser(containerId string, uuid string) {
	// 删除用户和对应的家目录
	newUserName := strings.Replace(uuid, "-", "", -1)
	cmds := []string{"userdel", "-r", newUserName}
	message := runCmdByContainer(containerId, cmds, "", "", "userdel", "")
	if message.ExitCode != utilsType.EXIT_CODE_OK {
		log.Errorf("userdel %v fail", newUserName)
	}
	// 删除新建的用于存放代码的目录
	cmds = []string{"rm", "-rf", uuid}
	message = runCmdByContainer(containerId, cmds, "", "", "rm -rf uuid", "")
	if message.ExitCode != utilsType.EXIT_CODE_OK {
		log.Errorf("rm -rf %v fail", uuid)
	}
}

func runCode(containerId string, dockerInfo utilsType.DockerInfo, inputList []string, workDir string, user string) []dto.ExecuteMessage {
	messages := make([]dto.ExecuteMessage, 0, 0)
	runCmd := dockerInfo.RunCmd
	runSplit := strings.Split(runCmd, " ")
	for _, inputStr := range inputList {
		runRes := runCmdByContainer(containerId, runSplit, workDir, inputStr, "run", user)
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

func (sandbox *SandBox) ExecuteCode(executeCodeRequest *dto.ExecuteCodeRequest) []dto.ExecuteMessage {
	// 1. 保存用户代码为文件
	code := executeCodeRequest.Code
	_, codeFilePath := sandbox.saveFile(code)
	// 4. 文件清理
	defer clearFile(codeFilePath)
	// 2. 编译代码并执行代码
	language := executeCodeRequest.Language
	inputList := executeCodeRequest.InputList
	executeMessageArrayList := sandbox.compileAndRun(language, codeFilePath, inputList)
	// 3. 整理输出信息
	executeMessages := getOutputMessage(executeMessageArrayList)
	return executeMessages
}
