package model

// API 统一响应格式
type API struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// OK 响应成功
func OK(data interface{}) *API {
	return &API{Success: true, Data: data}
}

// Fail 响应错误
func Fail(message string) *API {
	return &API{Success: false, Message: message}
}
