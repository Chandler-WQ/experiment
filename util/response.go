package util

import (
	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/common"
)

//通用返回值格式方法
func SuccessResponse(ctx *gin.Context, tips string, data interface{}) *common.Response {
	ctx.Set("res_code", common.Success.Code)
	return &common.Response{Code: common.Success.Code, Message: tips, Data: data}
}

//通用返回值格式方法
func FailResponse(ctx *gin.Context, code int32, tips string, data interface{}) *common.Response {
	ctx.Set("res_code", code)
	return &common.Response{Code: code, Message: tips, Data: data}
}

//通用返回值格式方法
func ErrResponse(ctx *gin.Context, cloudError *common.Response) *common.Response {
	ctx.Set("res_code", cloudError.Code)
	return &common.Response{Code: cloudError.Code, Message: cloudError.Message, Data: nil}
}
