package sandbox

import (
	"bytes"
	"codeSandbox/model/dto"
	"codeSandbox/utils"
	"context"
	"encoding/json"
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
	"time"
)

const CONTAINER_PREFIX string = "codeSandbox"
const WORDING_DIR string = "/codeSandbox"

// 代码沙箱执行过程的退出码
const (
	// 正常退出
	EXIT_CODE_OK = iota
	// 异常退出
	EXIT_CODE_ERROR
)

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
func runCmdByContainer(containerId string, cmd []string, workDir string, input string, tag string) dto.ExecuteMessage {
	message := dto.ExecuteMessage{}
	message.ExitCode = EXIT_CODE_ERROR
	ctx := context.Background()
	// 创建执行命令实例
	execConfig := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Tty:          false,
		Cmd:          cmd,
	}
	if workDir != "" {
		execConfig.WorkingDir = workDir
	}
	resp, err := DockerClient.ContainerExecCreate(ctx, containerId, execConfig)
	if err != nil {
		errMsg := fmt.Sprintf("ContainerExecCreate fail:%v", err)
		log.Panicf(errMsg)
		message.ErrorMessage = errMsg
		return message
	}

	// 开始时间打点
	startT := time.Now()

	// 开始监控容器（内存信息）
	memCn := make(chan uint64)
	//TODO 关闭
	//defer close(memCn)
	done := make(chan struct{})
	defer close(done)
	monitorReady := make(chan struct{})
	defer close(monitorReady)
	go monitorContainerStats(containerId, done, memCn, monitorReady, tag)

	// 等待监控程序就绪
	<-monitorReady

	// 启动执行命令并连接到输入输出流
	execID := resp.ID
	execAttachResp, err := DockerClient.ContainerExecAttach(ctx, execID, types.ExecStartCheck{})
	if err != nil {
		errMsg := fmt.Sprintf("ContainerExecAttach fail:%v{}", err)
		log.Panicf(errMsg)
		message.ErrorMessage = errMsg
		return message
	}
	defer execAttachResp.Close()

	hijackedResp := execAttachResp.Conn
	defer hijackedResp.Close()

	// 向输入流中写入数据
	if input != "" {
		log.Debugf("start to write data to stdin")
		write, err := hijackedResp.Write([]byte(input))
		if err != nil {
			errMsg := fmt.Sprintf("Write fail:%v{}", err)
			log.Panicf(errMsg)
		}
		log.Debugf("write:%v bytes finish", write)
	}

	// 创建一个bytes.Buffer实例用于接收输出
	var buf bytes.Buffer
	chDone := make(chan struct{})
	defer close(chDone)
	go func() {
		// 将 hijackedResp 中的数据复制到buf中
		_, err = io.Copy(&buf, hijackedResp)
		log.Debugf("read data from stdout finish...")
		chDone <- struct{}{}
	}()
	mainCn := make(chan struct{})
	go func() {
	Loop1:
		for {
			select {
			// 在规定时间内完成
			case <-chDone:
				// 关闭监控并获取最大内存消耗
				done <- struct{}{}
				memCost, readStatus := <-memCn
				tc := time.Since(startT)
				resultStr := buf.String()
				message.ExitCode = EXIT_CODE_OK
				message.Message = resultStr
				message.TimeCost = tc.Milliseconds()
				message.MemoryCost = memCost
				log.Debugf("before return message = %v, get memCost:%v readStatus:%v on %v", message, memCost, readStatus, tag)
				break Loop1
				// 超时完成
			case <-time.After(RUN_CODE_TIME_OUT):
				// 关闭监控并获取最大内存消耗
				done <- struct{}{}
				memCost, _ := <-memCn
				tc := time.Since(startT)
				message.ExitCode = EXIT_CODE_ERROR
				message.ErrorMessage = "Timout"
				message.TimeCost = tc.Milliseconds()
				message.MemoryCost = memCost
				break Loop1
			}
		}
		mainCn <- struct{}{} // 告诉主协程可以退出了
	}()
	<-mainCn
	return message
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

// 监控指定容器的内存信息（建议协程方式调用），直至主协程发来 done 信号，并将运行期间使用的最大内存返回
func monitorContainerStats(containerID string, done chan struct{}, memCn chan uint64, monitorReady chan struct{}, flag string) {
	ctx := context.Background()

	// 实时获取容器统计信息
	statsReader, err := DockerClient.ContainerStats(ctx, containerID, true)
	if err != nil {
		fmt.Errorf("Failed to start stats stream: %v", err)
	}
	defer statsReader.Body.Close()
	decoder := json.NewDecoder(statsReader.Body)
	var initMemoryUsage uint64 = 0
	var maxMemoryUsage uint64 = 0
	monitorReady <- struct{}{}
Loop:
	for {
		select {
		// 等待主协程发来的 done 信号
		case <-done:
			log.Debugf("return main coroutine initMemoryUsage:%v, maxMemoryUsage:%v, usage:%v on %v", initMemoryUsage, maxMemoryUsage, maxMemoryUsage-initMemoryUsage, flag)
			// 将运行期间占用的最大内存返回给主协程
			memCn <- maxMemoryUsage - initMemoryUsage
			break Loop
		default:
			var stat types.StatsJSON
			err := decoder.Decode(&stat)
			if err != nil {
				if err == io.EOF { // 连接被Docker关闭或意外中断
					log.Errorf("Connection closed or error receiving stats:", err)
					return
				}
				log.Errorf("Error decoding stats:", err)
				continue
			}

			// 处理统计信息
			memUsage := stat.MemoryStats.Usage
			//log.Debugf("memUsage : %v", memUsage)
			if initMemoryUsage == 0 {
				initMemoryUsage = memUsage
			}
			maxMemoryUsage = max(maxMemoryUsage, memUsage)
		}
	}

}
