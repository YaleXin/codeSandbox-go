package adminServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/service/mailServices"
	"codeSandbox/utils/global"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminService struct {
}

var AdminServiceInstance AdminService

var executionDao db.ExecutionDao
var userDao db.UserDao

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
	user := model.User{
		Model: gorm.Model{
			ID: userId,
		},
		Audit: true,
	}
	_, err := userDao.UpdateUserById(&user)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	// 发送邮件通知用户
	_, err = userDao.GetUserById(&user, userId)
	if err != nil {
		log.Warnf("After AuditUserByUserId query user fail:%v", err)
	} else {
		mailServiceInstance := &mailServices.MailServiceInstance
		// 使用协程，避免阻塞
		go mailServiceInstance.SendToUser(&user)
	}
	return global.SUCCESS
}

func (a *AdminService) RejectUserByUserId(userId uint) int {
	return global.SUCCESS
}

func (a *AdminService) BanUserByUserId(userId uint) int {
	user := model.User{
		Model: gorm.Model{
			ID: userId,
		},
		Ban: true,
	}
	_, err := userDao.UpdateUserById(&user)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	return global.SUCCESS
}

func (a *AdminService) CancelUserByUserId(userId uint) int {
	user := model.User{
		Model: gorm.Model{
			ID: userId,
		},
		Ban: false,
	}
	_, err := userDao.UpdateUserBanById(&user)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	return global.SUCCESS
}

func (a *AdminService) PageUser(page *dto.PageExecutionRequest) (int, *vo.PageDataVO) {
	// 每页大小应该是正整数，避免除零错误
	if page.PageSize == 0 {
		return global.PARAMS_ERROR, nil
	}
	users, total, err := userDao.PageUser(page.PageSize, page.PageNum)
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	// 封装页码和执行记录
	dataVO := vo.PageDataVO{
		Data:  getUserListVO(users),
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

func getUserListVO(users []model.User) []vo.UserDetailVO {
	userVOS := make([]vo.UserDetailVO, 0, len(users))
	// 每个都要封装
	for _, user := range users {
		detailVO := vo.UserDetailVO{}
		getUserDetailVO(&user, &detailVO)
		userVOS = append(userVOS, detailVO)
	}
	return userVOS
}

func getUserDetailVO(user *model.User, userVO *vo.UserDetailVO) {
	userVO.Id = user.ID
	userVO.Username = user.Username
	userVO.Role = user.Role
	userVO.Email = user.Email
	userVO.CurrentUsage = user.CurrentUsage
	userVO.MonthLimit = user.MonthLimit
	userVO.CreateAt = user.CreatedAt
	userVO.Ban = user.Ban
	userVO.Audit = user.Audit
}
