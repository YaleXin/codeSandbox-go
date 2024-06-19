package dto

import (
	"codeSandbox/model/vo"
	"codeSandbox/utils"
)

const ERR_MSG_RUNTIME_ERROR string = "RUNTIME ERROR"

type ExecuteMessage struct {
	ExitCode     int8
	Message      string
	ErrorMessage string
	TimeCost     int64
	MemoryCost   uint64
}

func (executeMsg *ExecuteMessage) ToVO() vo.ExecuteMessageVO {
	messageVO := vo.ExecuteMessageVO{
		ExitCode:     executeMsg.ExitCode,
		Message:      executeMsg.Message,
		ErrorMessage: executeMsg.ErrorMessage,
		TimeCost:     executeMsg.TimeCost,
		MemoryCost:   executeMsg.MemoryCost,
	}
	// 信息脱敏
	if messageVO.ExitCode != utils.EXIT_CODE_OK {
		messageVO.ErrorMessage = utils.EXIT_ERROR_MESSAGE[messageVO.ExitCode]
	}
	// 将 ErrorMessage 脱敏
	return messageVO
}
