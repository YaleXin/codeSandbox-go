package test

import (
	"codeSandbox/service/sandbox"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
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
