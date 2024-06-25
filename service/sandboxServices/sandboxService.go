package sandboxServices

import (
	dto "codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/service/cryptoServices"
	"codeSandbox/service/keypairService"
	"codeSandbox/service/sandboxDockerServices"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"encoding/json"
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

func (sandboxService *SandboxService) ProgramExecuteCode(programExecuteCodeRequest *dto.ProgramExecuteCodeRequest) (int, *dto.ExecuteCodeResponse) {
	// 1. 根据公钥，找出私钥
	// 2. 使用私钥解密出 json 字符串
	// 3. 字符串还原成 dto.ExecuteCodeRequest
	// 4. 调用 ExecuteCode()

	// 1.
	keypairServiceInstance := &keypairService.KeyPairServiceInstance
	code, keyPair := keypairServiceInstance.GetKeyPairByPublicKey(programExecuteCodeRequest.PublicKey)
	if code != global.SUCCESS {
		return code, nil
	}
	cryptoInstance := &cryptoServices.CryptoServiceInstance
	privateKeyBase64 := keyPair.SecretKey
	encryptedBase64 := programExecuteCodeRequest.Payload
	// 2.
	decryptWithPrivateKeyBase64, err := cryptoInstance.DecryptWithPrivateKeyBase64(privateKeyBase64, encryptedBase64)
	if err != nil {
		return global.KEY_PAIR_ERROR, nil
	}
	// 3.

	var executeCodeRequest dto.ExecuteCodeRequest
	err = json.Unmarshal([]byte(decryptWithPrivateKeyBase64), &executeCodeRequest)
	if err != nil {
		return global.PARAMS_ERROR, nil
	}
	executeCode := sandboxService.ExecuteCode(executeCodeRequest)
	return global.SUCCESS, &executeCode
}

func (sandboxService *SandboxService) ExecuteCode(executeCodeRequest dto.ExecuteCodeRequest) dto.ExecuteCodeResponse {
	// 找出该语言对应的 dockerinfo 对象
	language := executeCodeRequest.Language
	byLanguage := getDockerInfoByLanguage(language)
	box := sandboxDockerServices.SandBox{
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
