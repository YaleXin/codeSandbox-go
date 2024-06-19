package service

import (
	dto "codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/service/sandbox"
	"codeSandbox/utils"
)

type SandboxService struct{}

func (sandboxService *SandboxService) GetSupportLanguages() []string {
	DockerInfoList := utils.Config.DockerInfoList
	var languages []string
	for _, info := range DockerInfoList {
		languages = append(languages, info.Language)
	}
	return languages
}

func (sandboxService *SandboxService) ExecuteCode(executeCodeRequest dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	// 找出该语言对应的 dockerinfo 对象
	language := executeCodeRequest.Language
	byLanguage := getDockerInfoByLanguage(language)
	box := sandbox.SandBox{
		DockerInfo: byLanguage,
	}

	// 获取每个执行用例的输出
	executeMessages := box.ExecuteCode(&executeCodeRequest)
	// 对执行用例脱敏
	executeCodeResponse := dto.ExecuteCodeResponse{}
	executeCodeResponse.ExecuteMessages = make([]vo.ExecuteMessageVO, 0, 0)
	for _, executeMessage := range executeMessages {
		executeCodeResponse.ExecuteMessages = append(executeCodeResponse.ExecuteMessages, executeMessage.ToVO())
	}
	return executeCodeResponse
}

func getDockerInfoByLanguage(codeLanguage string) utils.DockerInfo {
	DockerInfoList := utils.Config.DockerInfoList

	for _, info := range DockerInfoList {
		if info.Language == codeLanguage {
			return info
		}
	}
	return utils.DockerInfo{}
}
