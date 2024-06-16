package dto

type ExecuteMessage struct {
	ExitCode     int8
	Message      string
	ErrorMessage string
	TimeCost     int64
	MemoryCost   uint64
}
