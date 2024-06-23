package global

func GetErrMsg(code int) string {
	return codemsg[code]
}

var codemsg = map[int]string{
	SUCCESS:           "OK",
	PARAMS_ERROR:      "请求参数错误",
	DATA_REPEAT:       "数据重复",
	NOT_LOGIN_ERROR:   "未登录",
	NO_AUTH_ERROR:     "无权限",
	LACK_AUTH_ERROR:   "权限不足",
	BAN_ERROR:         "已经禁用该用户",
	NOT_FOUND_ERROR:   "请求数据不存在",
	FORBIDDEN_ERROR:   "禁止访问",
	SYSTEM_ERROR:      "系统内部异常",
	OPERATION_ERROR:   "操作失败",
	API_REQUEST_ERROR: "外部API调用失败",
}
