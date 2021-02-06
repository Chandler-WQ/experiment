package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func Session() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := util.GetCookie(ctx, common.SessionKey)
		if err != nil {
			ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SessionErr.Code, common.SessionErr.Message, nil))
			ctx.Abort()
			return
		}

		session, err := service.GetSession(token)
		if err != nil {
			ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SessionErr.Code, common.SessionErr.Message, err.Error()))
			ctx.Abort()
			return
		}
		ctx.Set(common.SessionInfo, session)
		ctx.Next()
	}
}
