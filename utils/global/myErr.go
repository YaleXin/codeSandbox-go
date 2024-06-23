package global

import "fmt"

// 自定义错误类型
type CustomError struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

// 实现error接口的Error方法
func (e *CustomError) Error() string {
	return fmt.Sprintf(`{"code": %d, "message": "%s"}`, e.ErrorCode, e.Message)
}
func (e *CustomError) Code() int {
	return e.ErrorCode
}
