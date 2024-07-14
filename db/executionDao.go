package db

import (
	"codeSandbox/model"
	"fmt"
)

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

func (e *ExecutionDao) PageExecutionByUserId(userId uint, pageSize, pageNumber int64) ([]model.Execution, int64, error) {
	var executions []model.Execution

	// 使用Preload预加载User关系， 查询总数
	var count int64
	result := dBClinet.Where("user_id = ?", userId).Find(&executions).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error // 返回查询过程中可能遇到的错误
	}
	// 跳过的记录数
	offset := (pageNumber - 1) * pageSize
	result = dBClinet.Where("user_id = ?", userId).Offset(int(offset)).Limit(int(pageSize)).Find(&executions)
	//result = dBClinet.Where("user_id = ?", userId).Find(&executions).Offset(int(offset)).Limit(int(pageSize))
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return executions, count, nil
}

func (e *ExecutionDao) PageExecution(pageSize, pageNumber int64) ([]model.Execution, int64, error) {
	var executions []model.Execution

	var count int64
	result := dBClinet.Find(&executions).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error // 返回查询过程中可能遇到的错误
	}
	// 跳过的记录数
	offset := (pageNumber - 1) * pageSize
	// 顺便把 User 信息也查出来
	result = dBClinet.Preload("User").Offset(int(offset)).Limit(int(pageSize)).Find(&executions)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return executions, count, nil
}

func (e *ExecutionDao) ExecutionAdd(execution *model.Execution) (int64, error) {
	create := dBClinet.Create(execution)
	err := create.Error
	if err != nil {
		return 0, err
	}
	affected := create.RowsAffected
	return affected, nil
}

func (e *ExecutionDao) ExecutionUpdateById(execution *model.Execution) error {
	// gorm 中，没有主键时会调用 create ，这是不符合我们的
	if execution.ID == uint(0) {
		return fmt.Errorf("execution.id == 0")
	}
	// 根据 `struct` 更新属性，只会更新非零值的字段
	updates := dBClinet.Updates(execution)
	err := updates.Error
	if err != nil {
		return err
	}
	return nil
}
