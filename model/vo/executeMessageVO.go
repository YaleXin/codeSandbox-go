package vo

type ExecuteMessageVO struct {
	ExitCode int8
	Message  string
	// 脱敏信息，将 docker 容器的错误信息做进一步处理
	ErrorMessage string
	TimeCost     int64
	MemoryCost   uint64
}
