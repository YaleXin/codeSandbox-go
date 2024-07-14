package adminServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/utils/global"
	"encoding/json"
)

type AdminService struct {
}

var AdminServiceInstance AdminService

var executionDao db.ExecutionDao

func (a *AdminService) PageExecution(page *dto.PageExecutionRequest) (int, *vo.PageDataVO) {
	// 每页大小应该是正整数，避免除零错误
	if page.PageSize == 0 {
		return global.PARAMS_ERROR, nil
	}
	executions, total, err := executionDao.PageExecution(page.PageSize, page.PageNum)
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	// 封装页码和执行记录
	dataVO := vo.PageDataVO{
		Data:  getExecutionListVO(executions),
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
func getExecutionListVO(eList []model.Execution) []vo.AdminExecutionVO {
	executionVOS := make([]vo.AdminExecutionVO, 0, len(eList))
	// 每个都要封装
	for _, execution := range eList {
		executionVOS = append(executionVOS, *getExecutionVO(&execution))
	}
	return executionVOS
}
func getExecutionVO(e *model.Execution) *vo.AdminExecutionVO {
	// 在数据库中，输入列表是被序列化后的，因此要先还原
	var inputList []string
	var outputList []string
	json.Unmarshal([]byte(e.InputList), &inputList)
	json.Unmarshal([]byte(e.OutputList), &outputList)
	executionVO := vo.AdminExecutionVO{
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
		UserId:        e.UserId,
		Username:      e.User.Username,
	}
	return &executionVO
}
func (a *AdminService) AuditUserByUserId(userId uint) int {
	return global.SUCCESS
}

func (a *AdminService) RejectUserByUserId(userId uint) int {
	return global.SUCCESS
}

func (a *AdminService) BanUserByUserId(userId uint) int {
	return global.SUCCESS
}
