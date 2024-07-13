package vo

// 验证码返回体，需要用验证码的地方要请求该数据，例如注册和登录
type CaptchaVO struct {
	ImageBase64 string `json:"imageBase64"` // 验证码 base64
	Token       string `json:"token"`       // 登录和注册要带上
}
