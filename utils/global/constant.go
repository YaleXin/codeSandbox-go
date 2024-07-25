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
	CONTAINER_MAX_MEMORY = 1024 * 1024 * 1024 // 每个容器最大允许使用 1 GB 内存
	CONTAINER_MAX_CPU    = 1                  // 每个容器最大允许使用 1 个 CPU
)
const (
	USER_VALIDITY_PERIOD    = 180 // 用户体验时间（单位为天）
	CAPTCHA_LEN             = 5   // 验证码长度
	INPUT_LIST_MAX_LEN      = 5   //输入用例最大数目
	MAINTAIN_KEY_MAX_LEN    = 5   //每个用户最多可以有的 key
	EXECUTION_PAGE_MAX_SIZE = 100 // 用户查询执行记录每页最大数
	USER_INIT_MONTH_LIMIT   = 500 // 每个用户初始的每月调用额度
)
const (
	APP_MODE_DEV  = "dev"
	APP_MODE_PROD = "prod"
)
const (
	REGISTER_URL_PREFIX = "https://code.yalexin.top/api/v1/user/check?token="
)
