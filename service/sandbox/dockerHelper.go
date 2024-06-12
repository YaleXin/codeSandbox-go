package sandbox

import (
	"bytes"
	"codeSandbox/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const CONTAINER_PREFIX string = "codeSandbox"
const WORDING_DIR string = "/codeSandbox"

var DockerClient *client.Client

type SandBox struct {
	DockerInfo utils.DockerInfo
}

// 将编程语言和下标组合成容器名字
func getContainerName(language string, idx int) string {
	return CONTAINER_PREFIX + "_" + language + "_" + strconv.Itoa(idx)
}

// 根据编程语言信息和下标获取容器id，当容器不存在时，创建容器
func getContainerId(dockerInfo utils.DockerInfo, idx int) string {
	ctx := context.Background()
	containerName := getContainerName(dockerInfo.Language, idx)

	// 列出所有容器，过滤出指定名称的容器
	containers, err := DockerClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		log.Errorf("Failed to list containers: %v", err)
	}

	containerId := ""
	// 检查是否有匹配的容器
	if len(containers) == 0 {
		log.Debugf("Container '%s' not found.\n", containerName)

	} else if len(containers) > 1 {
		for _, cntn := range containers {
			names := cntn.Names
			// 由于返回来的名字中，会在前面自动添加 "/" 因此要先把它去掉
			if names[0][1:] == containerName {
				containerId = cntn.ID
				break
			}
		}
	}
	if containerId == "" {
		containerId = createContainer(&dockerInfo, containerName)
	}
	return containerId
}

// 当 workDir 为空字符串，即 ""，则不设置 WorkingDir
func runCmdByContainer(containerId string, cmd []string, workDir string) string {
	ctx := context.Background()
	// 创建执行命令实例
	execConfig := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  false,
		Tty:          false,
		Cmd:          cmd,
	}
	if workDir != "" {
		execConfig.WorkingDir = workDir
	}
	resp, err := DockerClient.ContainerExecCreate(ctx, containerId, execConfig)
	if err != nil {
		log.Panicf("ContainerExecCreate fail:%v", err)
	}

	// 启动执行命令并连接到输出流
	execID := resp.ID
	respAttach, err := DockerClient.ContainerExecAttach(context.Background(), execID, types.ExecStartCheck{})
	if err != nil {
		log.Panicf("ContainerExecAttach fail: %v", err)
	}
	defer respAttach.Close()
	var stdoutBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		io.Copy(&stdoutBuf, respAttach.Reader)
	}()

	// 等待goroutines完成
	wg.Wait()
	// 将读取的内容转换为字符串
	stdoutStr := stdoutBuf.String()
	log.Debugf("STDOUT:", stdoutStr)

	return stdoutStr
}

func createContainer(dockerInfo *utils.DockerInfo, containerName string) string {
	ctx := context.Background()
	resp, err := DockerClient.ContainerCreate(ctx, &container.Config{
		Image:           dockerInfo.ImageName,
		AttachStdin:     true,
		AttachStdout:    true,
		AttachStderr:    true,
		Tty:             true,
		NetworkDisabled: true,
		WorkingDir:      WORDING_DIR,
	}, nil, nil, nil, containerName)
	var containerId string
	// 判断容器名字是否被占用，被占用则直接启动
	if err != nil {
		errorStr := err.Error()
		if strings.Contains(errorStr, "already in use") {
			// 提取出容器 id
			// 编译正则表达式，用于匹配容器ID
			re := regexp.MustCompile(`container "([^"]+)".*`)
			matches := re.FindStringSubmatch(errorStr)

			// 检查是否有匹配项
			if len(matches) > 1 {

				containerId = matches[1]
			} else {
				log.Error("No container ID found in the error message.")
			}
		}
		log.Errorf("create container %v fail:%v", containerName, err)
	} else {
		containerId = resp.ID
	}
	err = DockerClient.ContainerStart(ctx, containerId, container.StartOptions{})
	if err != nil {
		log.Errorf("start container %v %v fail: %v", containerName, containerId, err)
	}

	inspect, err := DockerClient.ContainerInspect(ctx, containerId)
	if err != nil {
		log.Errorf("get ContainerInspect fail %v", err)
	}
	log.Debugf("container: %q,status: %q", inspect.ID[:10], inspect.State.Status)
	return containerId
}
func initContainer(dockerInfoList *[]utils.DockerInfo) {
	languageListLen := len(*dockerInfoList)
	var wait sync.WaitGroup
	// 每种编程语言使用一个协程来创建容器
	wait.Add(languageListLen)
	for _, dockerInfo := range *dockerInfoList {
		log.Debugf("handle %v ...", dockerInfo)
		go func(info utils.DockerInfo) {
			log.Debugf("go coroutine get %v", info)
			// 创建指定数量的容器
			for i := 0; i < info.ContainerCount; i++ {
				containerName := getContainerName(info.Language, i)
				createContainer(&info, containerName)
			}
			wait.Done()
			log.Debugf("handle %v finish, create %v containers", info, info.ContainerCount)
		}(dockerInfo)
	}
	wait.Wait()
	log.Infof("init container success")
}
func initImagesAndContainer() {
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
	log.Infof("images prepare success, start to init container")
	// 必须将所需的镜像下载完后才能初始化容器
	go initContainer(&list)
}
func init() {
	docker, err := connectDocker()
	if err != nil {
		log.Panicf("init docker fail:%v", err)
		return
	}
	DockerClient = docker
	go initImagesAndContainer()
}
func connectDocker() (cli *client.Client, err error) {
	dockerConfig := utils.Config.SandboxMachine

	cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost(fmt.Sprintf("tcp://%v:%v", dockerConfig.Host, dockerConfig.Port)))
	if err != nil {
		return nil, err
	}

	return cli, nil
}
