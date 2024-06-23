package global

const (
	SUCCESS              = 0
	PARAMS_ERROR         = 40000
	DATA_REPEAT_ERROR    = 40001
	NOT_LOGIN_ERROR      = 40100
	NO_AUTH_ERROR        = 40101
	LACK_AUTH_ERROR      = 40102
	BAN_ERROR            = 40103
	NOT_FOUND_ERROR      = 40400
	NOT_FOUND_USER_ERROR = 40401
	PWD_ERROR            = 40402
	FORBIDDEN_ERROR      = 40300
	SYSTEM_ERROR         = 50000
	OPERATION_ERROR      = 50001
	API_REQUEST_ERROR    = 50020
)

var codemsg = map[int]string{
	SUCCESS:              "OK",
	PARAMS_ERROR:         "请求参数错误",
	DATA_REPEAT_ERROR:    "数据重复",
	NOT_LOGIN_ERROR:      "未登录",
	NO_AUTH_ERROR:        "无权限",
	LACK_AUTH_ERROR:      "权限不足",
	BAN_ERROR:            "已经禁用该用户",
	NOT_FOUND_ERROR:      "请求数据不存在",
	NOT_FOUND_USER_ERROR: "用户不存在",
	PWD_ERROR:            "密码错误",
	FORBIDDEN_ERROR:      "禁止访问",
	SYSTEM_ERROR:         "系统内部异常",
	OPERATION_ERROR:      "操作失败",
	API_REQUEST_ERROR:    "外部API调用失败",
}

func GetErrMsg(code int) string {
	return codemsg[code]
}
