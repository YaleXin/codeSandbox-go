package service_v1

import (
	dto "codeSandbox/model/dto"
)

func GetSupportLanguages() []string {
	res := []string{"DEBUG a", "DEBUG b"}
	return res
}

func ExecuteCode(executeCodeRequest dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	e := dto.ExecuteCodeResponse{}
	return e
}
