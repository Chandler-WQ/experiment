package util

import (
	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/common"
)

//通用返回值格式方法
func NewSuccessResponse(ctx *gin.Context, data interface{}, tips string) *common.Response {
	ctx.Set("res_code", common.Success.StatusCode)
	return &common.Response{StatusCode: common.Success.StatusCode, Message: tips, Data: data}
}

//通用返回值格式方法
func NewFailResponse(ctx *gin.Context, code int32, tips string,data interface{}) *common.Response {
	ctx.Set("res_code", code)
	return &common.Response{StatusCode: code, Message: tips, Data: data}
}

//通用返回值格式方法
func NewErrResponse(ctx *gin.Context, cloudError *common.Response) *common.Response {
	ctx.Set("res_code", cloudError.StatusCode)
	return &common.Response{StatusCode: cloudError.StatusCode, Message: cloudError.Message, Data: nil}
}
