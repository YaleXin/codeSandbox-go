package service_v1

import (
	dto "codeSandbox/model/dto"
	"codeSandbox/utils"
)

func GetSupportLanguages() []string {
	DockerInfoList := utils.Config.DockerInfoList
	var languages []string
	for _, info := range DockerInfoList {
		languages = append(languages, info.Language)
	}
	return languages
}

func ExecuteCode(executeCodeRequest dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	e := dto.ExecuteCodeResponse{}
	return e
}
