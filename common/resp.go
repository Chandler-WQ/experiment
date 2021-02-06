package common

import (
	"fmt"
)

var (
	Success       = New(0, "success")
	ParaErr       = New(1, "参数格式错误或缺少对应参数")
	DBErr         = New(2, "数据库错误，请重试或并联系管理员")
	SerErr        = New(3, "服务错误，请联系管理员")
	SessionErr    = New(4, "请先登录")
	PermissionErr = New(5, "Permission denied")
)

type Response struct {
	Code    int32       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func New(code int32, message string) *Response {
	err := &Response{
		Code:    code,
		Message: message,
		Data:    fmt.Sprintf("%d:%s", code, message),
	}
	return err
}
