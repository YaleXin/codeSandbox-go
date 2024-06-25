package db

import "codeSandbox/model"

type ExecutionDao struct {
}

// 找出没有被 “delete” 的
func (e *ExecutionDao) ListExecution() ([]model.Execution, error) {
	var executions []model.Execution
	find := dBClinet.Find(&executions)
	err := find.Error
	if err != nil {
		return nil, err
	}
	return executions, nil
}
func (e *ExecutionDao) ListExecutionByUserId(userId uint) ([]model.Execution, error) {
	var executions []model.Execution
	// 使用Preload预加载User关系
	result := dBClinet.Where("user_id = ?", userId).Find(&executions)
	if result.Error != nil {
		return nil, result.Error // 返回查询过程中可能遇到的错误
	}
	return executions, nil
}
