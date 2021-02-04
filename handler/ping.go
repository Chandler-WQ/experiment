package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context){
	ctx.String(http.StatusOK, "pong")
}
