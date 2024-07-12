package responses

import "codeSandbox/utils/global"

var (
	// OK
	// 通用成功
	OK = myResponse(200, "ok")
	// 通用错误
	Err = myResponse(500, "")

	// TODO 把每种错误都定义一个 myResponse
	// 服务级错误码
	ErrParam     = myResponse(10001, "参数有误")
	ErrSignParam = myResponse(10002, "签名参数有误")

	ErrTooManyRequests = myResponse(global.TOO_MANY_REQUESTS_ERROR, global.GetErrMsg(global.TOO_MANY_REQUESTS_ERROR))
)
