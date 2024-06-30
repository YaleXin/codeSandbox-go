package sandboxServices

import (
	"codeSandbox/model"
	dto "codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/service/cryptoServices"
	"codeSandbox/service/executionServices"
	"codeSandbox/service/keypairService"
	"codeSandbox/service/sandboxDockerServices"
	"codeSandbox/service/userServices"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"codeSandbox/utils/tool"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
func checkProgramExecuteCodeRequest(programExecuteCodeRequest *dto.ProgramExecuteCodeRequest, privateKeyBase64 string) bool {
	if tool.IsStructEmpty(programExecuteCodeRequest) {
		return false
	}
	// 判断签名
	return verifySignature(programExecuteCodeRequest, privateKeyBase64)
}

// verifySignature 判断签名
func verifySignature(programExecuteCodeRequest *dto.ProgramExecuteCodeRequest, privateKeyBase64 string) bool {
	return tool.MD5Str(programExecuteCodeRequest.Payload+privateKeyBase64) == programExecuteCodeRequest.Signature
}
func (sandboxService *SandboxService) ProgramExecuteCode(c *gin.Context, programExecuteCodeRequest *dto.ProgramExecuteCodeRequest) (int, *dto.ExecuteCodeResponse) {
	// 1. 根据公钥，找出私钥
	// 1.5 判断签名
	// 2. 使用私钥解密出 json 字符串
	// 3. 字符串还原成 dto.ExecuteCodeRequest
	// 4. 调用 sandboxDockerServices.ExecuteCode()

	// 1.
	keypairServiceInstance := &keypairService.KeyPairServiceInstance
	code, keyPair := keypairServiceInstance.GetKeyPairByPublicKey(programExecuteCodeRequest.PublicKey)
	if code != global.SUCCESS {
		return code, nil
	}
	cryptoInstance := &cryptoServices.CryptoServiceInstance
	privateKeyBase64 := keyPair.SecretKey
	encryptedBase64 := programExecuteCodeRequest.Payload
	// 1.5
	if !checkProgramExecuteCodeRequest(programExecuteCodeRequest, privateKeyBase64) {
		return global.PARAMS_ERROR, nil
	}
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
	// 4.
	code, executeCode := sandboxService.ExecuteCode(c, executeCodeRequest, keyPair)
	return code, &executeCode
}

// 要么是用户提供 accesskey 的方式调用， 要么是用户登录后调用
func (sandboxService *SandboxService) ExecuteCode(c *gin.Context, executeCodeRequest dto.ExecuteCodeRequest, keyPair *model.KeyPair) (int, dto.ExecuteCodeResponse) {
	// 找出该语言对应的 dockerinfo 对象
	language := executeCodeRequest.Language
	byLanguage := getDockerInfoByLanguage(language)
	box := sandboxDockerServices.SandBox{
		DockerInfo: byLanguage,
	}
	// 先添加一条执行记录到数据库中
	var execution model.Execution
	code := addExecutionRecord(c, &executeCodeRequest, &execution, keyPair)
	if code != global.SUCCESS {
		return code, dto.ExecuteCodeResponse{}
	}
	// 获取每个执行用例的输出
	executeMessages := box.ExecuteCode(&executeCodeRequest)
	// 更新数据库中的执行记录
	updateExecutionRecord(&execution, executeMessages)
	// 对执行用例脱敏
	executeCodeResponse := dto.ExecuteCodeResponse{}
	executeCodeResponse.ExecuteMessages = make([]vo.ExecuteMessageVO, 0, 0)
	for _, executeMessage := range executeMessages {
		executeCodeResponse.ExecuteMessages = append(executeCodeResponse.ExecuteMessages, executeMessage.ToVO())
	}
	return global.SUCCESS, executeCodeResponse
}

func updateExecutionRecord(execution *model.Execution, messages []dto.ExecuteMessage) int {
	// 只有所有的输出用例都是正常，该执行记录状态才为正常
	isNormalExit := true
	maxMemoryCost := uint64(0)
	maxTimeCost := int64(0)
	outputList := make([]string, 0, len(messages))
	for _, msg := range messages {
		if msg.ExitCode != utils.EXIT_CODE_OK {
			isNormalExit = false
			outputList = append(outputList, msg.ErrorMessage)
		} else {
			outputList = append(outputList, msg.Message)
		}
		maxTimeCost = max(maxTimeCost, msg.TimeCost)
		maxMemoryCost = max(maxMemoryCost, msg.MemoryCost)
	}
	// 将输出转为 json 字符串
	outputListBytes, _ := json.Marshal(outputList)
	execution.OutputList = string(outputListBytes)
	if isNormalExit {
		execution.Status = global.EXECUTION_STATUS_NOMAL_EXIT
	} else {
		execution.Status = global.EXECUTION_STATUS_ERROR_EXIT
	}
	execution.MaxTimeCost = maxTimeCost
	execution.MaxMemoryCost = maxMemoryCost

	instance := &executionServices.ExecutionServiceInstance
	updateExecutionStatus := instance.UpdateExecution(execution)
	return updateExecutionStatus
}

func addExecutionRecord(c *gin.Context, executeCodeRequest *dto.ExecuteCodeRequest, execution *model.Execution, keyPair *model.KeyPair) int {
	// 如果不是程序方式调用，则说明是登录后在前端调用
	if keyPair == nil {
		_, user := userServices.UserServiceInstance.GetLoginUser(c)
		execution.User = *user
		execution.UserId = user.ID
	} else {
		// 反之如果是程序方式调用，其会提供 accessKey，前面的逻辑会查出对应的用户
		user := keyPair.User
		execution.User = user
		execution.UserId = keyPair.UserId
		execution.KeyPairId = keyPair.ID
	}
	execution.Status = global.EXECUTION_STATUS_RUNNING
	execution.Language = executeCodeRequest.Language
	execution.Code = executeCodeRequest.Code
	inputListJsonBytes, _ := json.Marshal(executeCodeRequest.InputList)
	execution.InputList = string(inputListJsonBytes)
	// 判断有没有额度调用
	code := userServices.UserServiceInstance.CheckAndUpdateUserUsage(execution.UserId)
	if code != global.SUCCESS {
		return code
	}
	instance := &executionServices.ExecutionServiceInstance
	addExecutionStatus := instance.AddExecution(execution)
	return addExecutionStatus

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
