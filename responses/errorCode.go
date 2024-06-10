package responses

var (
	// OK
	// 通用成功
	OK = myResponse(200, "ok")
	// 通用错误
	Err = myResponse(500, "")

	// 服务级错误码
	ErrParam     = myResponse(10001, "参数有误")
	ErrSignParam = myResponse(10002, "签名参数有误")
)
