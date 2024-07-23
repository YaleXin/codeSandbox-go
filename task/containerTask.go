package task

import (
	"codeSandbox/service/sandboxDockerServices"
	log "github.com/sirupsen/logrus"
)

func ResetCodeSandbox() {
	log.Infof("Start to reset codeSandbox...")
	box := sandboxDockerServices.SandBox{}
	box.ResetCodeSandbox()
	log.Infof("Reset codeSandbox finish...")
}
