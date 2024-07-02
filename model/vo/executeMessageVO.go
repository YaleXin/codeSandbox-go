package vo

type ExecuteMessageVO struct {
	ExitCode     int8   `json:"exitCode"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage"` // 脱敏信息，将 docker 容器的错误信息做进一步处理
	TimeCost     int64  `json:"timeCost"`
	MemoryCost   uint64 `json:"memoryCost"`
}
