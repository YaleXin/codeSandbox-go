package dto

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"` // 旧密码
	NewPassword string `json:"newPassword"` // 新密码
}
