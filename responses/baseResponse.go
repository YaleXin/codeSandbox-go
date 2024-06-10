package responses

import "encoding/json"

// 通用返回结构体和函数

type Response struct {
	// 错误码
	Code int `json:"code"`
	// 错误描述
	Msg string `json:"msg"`
	// 返回数据
	Data interface{} `json:"data"`
}

// 自定义响应信息
func (res *Response) WithMsg(message string) Response {
	return Response{
		Code: res.Code,
		Msg:  message,
		Data: res.Data,
	}
}

// 追加响应数据
func (res *Response) WithData(data interface{}) Response {
	return Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *Response) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: res.Code,
		Msg:  res.Msg,
		Data: res.Data,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// 构造函数
func myResponse(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
