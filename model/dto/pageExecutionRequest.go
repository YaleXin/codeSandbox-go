package dto

// @Summary 执行记录请求参数
// @Description 执行记录请求参数
// @Accept json
// @Param pageSize body string true "每页大小"
// @Param pageNum body string true "页码数"
type PageExecutionRequest struct {
	PageSize int64 `bind:"required" json:"pageSize"` // 每页大小
	PageNum  int64 `bind:"required" json:"pageNum"`  // 页码数
}
