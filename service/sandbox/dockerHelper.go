package sandbox

import (
	"codeSandbox/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

var DockerClient *client.Client

type SandBox struct {
	DockerInfo utils.DockerInfo
}

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
