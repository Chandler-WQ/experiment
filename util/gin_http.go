package util

import (
	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/common"
)

var domain = "localhost"

func GetCookie(ctx *gin.Context, key string) (value string, err error) {
	return ctx.Cookie(key)
}

func SetCookie(ctx *gin.Context, name, key string) {
	ctx.SetCookie(name, key, common.SessionAge, "/", domain, false, false)
	return
}
