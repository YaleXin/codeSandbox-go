package dto

type ExecuteCodeRequest struct {
	Code     string
	Language string
	// 运行用例组，每一个元素代表一个用例，例如 1 2\n
	InputList []string
}
