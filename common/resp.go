package common

import (
	"fmt"
)

var (
	Success                = New(0, "success")
	ParamError             = New(1, "param error")
	DBError                = New(2, "db error，please try again later")
	ServerError            = New(3, "server error")
)

type Response struct {
	StatusCode int32       `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

const (
	nilString = "nil"
)

func New(code int32, message string) *Response {
	err := &Response{
		StatusCode: code,
		Message:    message,
		Data:       fmt.Sprintf("%d:%s", code, message),
	}
	return err
}
func (e *Response) ErrNo() int32 {
	if e == nil {
		return 0
	}
	return e.StatusCode
}
func (e *Response) ErrTips() string {
	if e == nil {
		return nilString
	}
	return e.Message
}