package dto

type ExecuteCodeRequest struct {
	Code      string
	Language  string
	InputList []string
}
