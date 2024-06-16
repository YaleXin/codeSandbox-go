package dto

// @Summary 执行代码请求参数
// @Description 用户需要提交的代码执行请求参数
// @Tags Code Execution
// @Accept json
// @Param code body string true "待执行的代码"
// @Param language body string true "代码语言"
// @Param inputList body []string true "运行用例数组，每个元素代表一个用例的输入"
// @Router /execute-code [post]
type ExecuteCodeRequest struct {
	Code     string `bind:"required" json:"code"`
	Language string `bind:"required" json:"language"`
	// 运行用例组，每一个元素代表一个用例，例如 1 2\n
	InputList []string `bind:"required" json:"inputList"`
}
