package dto

type ProgramExecuteCodeRequest struct {
	// 使用 secret 加密后的数据
	Payload string `bind:"required" json:"payload"`
	// 公钥
	PublicKey string `bind:"required" json:"publicKey"`
}
