package sandbox

import (
	"codeSandbox/model/dto"
	"codeSandbox/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

const CODE_DIR_PREX string = "temp"

type SandBox struct {
	DockerInfo utils.DockerInfo
}

var DockerClient *client.Client

func initImages() {
	// 先看本地的镜像列表
	// 如果配置文件中指定的编程语言对应的镜像不在本地，则下载

	// 使用映射来模拟set，空结构体不会占用内存
	type empty struct{}
	localImageSet := make(map[string]empty)

	ctx := context.Background()
	images, _ := DockerClient.ImageList(ctx, types.ImageListOptions{All: true})
	for _, image := range images {
		localImageSet[image.RepoTags[0]] = empty{}
	}
	list := utils.Config.DockerInfoList

	shouldDownloadList := make([]string, 0, 0)
	for _, info := range list {
		imageName := info.ImageName
		// 检查元素是否存在
		if _, exists := localImageSet[imageName]; !exists {
			shouldDownloadList = append(shouldDownloadList, imageName)
		}
	}
	// 使用协程下载镜像，每个协程下载一个镜像
	var wait sync.WaitGroup
	wait.Add(len(shouldDownloadList))
	for _, image := range shouldDownloadList {
		log.Infof("start download image:%v", image)
		// 直接使用协程异步下载
		go func(imageName string) {
			reader, err := DockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
			if err != nil {
				log.Errorf("ImagePull %v fail:%v", imageName, err)
			}
			defer reader.Close()
			io.Copy(os.Stdout, reader)
			log.Infof("ImagePull %v success", imageName)
			// 协程执行完毕
			wait.Done()
		}(image)
	}
	// 等待子协程运行完毕
	wait.Wait()
}
func init() {
	docker, err := connectDocker()
	if err != nil {
		log.Panicf("init docker fail:%v", err)
		return
	}
	DockerClient = docker
	initImages()
}
func connectDocker() (cli *client.Client, err error) {
	dockerConfig := utils.Config.SandboxMachine

	cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost(fmt.Sprintf("tcp://%v:%v", dockerConfig.Host, dockerConfig.Port)))
	if err != nil {
		return nil, err
	}

	return cli, nil
}
func clearFile(codeFilename string) {
	err := os.Remove(codeFilename)
	if err != nil {
		log.Errorf("Remove code file fail:%v", err)
	}
	log.Debugf("clear file finish, codeFilename:%v", codeFilename)
}
func getOutputResponse(executeMessageArrayList []dto.ExecuteMessage) dto.ExecuteCodeResponse {
	return dto.ExecuteCodeResponse{}
}
func compileAndRun(language string, userCodeFile fs.File, inputList []string) []dto.ExecuteMessage {
	return nil
}

func (sandbox *SandBox) saveFile(code string) (fs.File, string) {
	// 不同的编程语言将会保存到不同的地方
	language := sandbox.DockerInfo.Language
	filename := sandbox.DockerInfo.Filename
	parentPath := CODE_DIR_PREX + string(filepath.Separator) + language
	// 限为 0666，表示为所有人都可以对该文件夹进行读写，且不存在时会自动创建。
	err := os.MkdirAll(parentPath, 0666)
	if err != nil {
		log.Errorf("MkdirAll %v fail:%v", parentPath, err)
		return nil, ""
	}
	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Errorf("Gennerate UUID fail: %v", err)
		return nil, ""
	}
	codeFilename := parentPath + string(filepath.Separator) + newUUID.String() + "_" + filename
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
	codeFile, codeFilename := sandbox.saveFile(code)
	// 2. 编译代码并执行代码
	language := executeCodeRequest.Language
	inputList := executeCodeRequest.InputList
	executeMessageArrayList := compileAndRun(language, codeFile, inputList)
	// 3. 整理输出信息
	executeCodeResponse := getOutputResponse(executeMessageArrayList)
	// 4. 文件清理
	defer clearFile(codeFilename)
	return executeCodeResponse
}
