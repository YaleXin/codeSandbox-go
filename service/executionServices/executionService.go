package executionServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/utils/global"
	"encoding/json"
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
func getExecutionVO(e *model.Execution) *vo.ExecutionVO {
	// 在数据库中，输入列表是被序列化后的，因此要先还原
	var inputList []string
	var outputList []string
	json.Unmarshal([]byte(e.InputList), &inputList)
	json.Unmarshal([]byte(e.OutputList), &outputList)
	executionVO := vo.ExecutionVO{
		Id:            e.ID,
		CreateAt:      e.CreatedAt,
		Code:          e.Code,
		Language:      e.Language,
		MaxMemoryCost: e.MaxMemoryCost,
		MaxTimeCost:   e.MaxTimeCost,
		InputList:     inputList,
		OutputList:    outputList,
		Status:        e.Status,
		KeyPairId:     e.KeyPairId,
	}
	return &executionVO
}
func getExecutionListVO(eList []model.Execution) []vo.ExecutionVO {
	executionVOS := make([]vo.ExecutionVO, 0, len(eList))
	// 每个都要封装
	for _, execution := range eList {
		executionVOS = append(executionVOS, *getExecutionVO(&execution))
	}
	return executionVOS
}
func (e *ExecutionService) PageUserExecution(user *model.User, page *dto.PageExecutionRequest) (int, *vo.PageDataVO) {
	// 每页大小应该是正整数，避免除零错误
	if page.PageSize == 0 {
		return global.PARAMS_ERROR, nil
	}
	// 不能超过预设页大小
	if page.PageSize > global.EXECUTION_PAGE_MAX_SIZE {
		return global.PARAMS_ERROR, nil
	}
	executionListByUserId, total, err := executionDao.PageExecutionByUserId(user.ID, page.PageSize, page.PageNum)
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	// 封装页码和执行记录
	dataVO := vo.PageDataVO{
		Data:  getExecutionListVO(executionListByUserId),
		Total: total,
	}
	// 总的页数，整除后向上取整
	if total%page.PageSize != 0 {
		dataVO.PageCount = total/page.PageSize + 1
	} else {
		dataVO.PageCount = total / page.PageSize
	}
	return global.SUCCESS, &dataVO
}
