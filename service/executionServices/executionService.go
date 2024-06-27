package executionServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/utils/global"
)

type ExecutionService struct{}

var ExecutionServiceInstance ExecutionService
var executionDao db.ExecutionDao

func (e *ExecutionService) AddExecution(execution *model.Execution) int {
	_, err := executionDao.ExecutionAdd(execution)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	return global.SUCCESS
}

func (e *ExecutionService) UpdateExecution(execution *model.Execution) int {
	err := executionDao.ExecutionUpdateById(execution)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	return global.SUCCESS
}
