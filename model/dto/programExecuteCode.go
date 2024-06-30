package dto

type ProgramExecuteCodeRequest struct {
	Payload   string `bind:"required" json:"payload"`   // 使用 secret 加密后的数据
	PublicKey string `bind:"required" json:"publicKey"` // 公钥
	Signature string `bind:"required" json:"signature"` // 对数据的签名，防止
}
