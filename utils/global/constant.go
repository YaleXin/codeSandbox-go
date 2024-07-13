package global

const (
	NORMAL_USER_ROLE = int(10)
	ADMIN_USER_ROLE  = int(1)
)
const (
	// 执行记录状态
	EXECUTION_STATUS_RUNNING = iota
	EXECUTION_STATUS_NOMAL_EXIT
	EXECUTION_STATUS_ERROR_EXIT
)
const (
	// 用户体验时间（单位为天）
	USER_VALIDITY_PERIOD = 7
	CAPTCHA_LEN          = 5
)
