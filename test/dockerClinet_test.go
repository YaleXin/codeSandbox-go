package test

import (
	"codeSandbox/service/sandbox"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"
)

func TestDockerClineInit(t *testing.T) {
	// 在测试开始前设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	client := sandbox.DockerClient
	logrus.Debugf("docker client:%v", client)

	err := getImageList(client)
	// docker client:&{http tcp://192.168.254.1:2375 tcp 192.168.254.1:2375  0xc00002f980 1.45 <nil> map[] false true false <nil> 0xc0001057c0}
	if err != nil {
		logrus.Errorf("GetImageList:%v", err)
	}

}

func TestContainerStart(t *testing.T) {
	ctx := context.Background()
	cli := sandbox.DockerClient

	imageName := "golang:1.17"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "aaa")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

}

func getImageList(cli *client.Client) error {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return err
	}
	//打印结果
	for _, image := range images {
		logrus.Debugf("%v:%v", image.RepoTags, image.ID)
	}
	return nil
}

func TestRunCmd(t *testing.T) {
	//logrus.SetLevel(logrus.DebugLevel)
	//containerId := "7634b7988b"
	//runCmdByContainer := sandbox.runCmdByContainer(containerId, []string{"ls", "-l"}, "/tmp")
	//logrus.Debugf("runCmdByContainer: %v", runCmdByContainer)
}
