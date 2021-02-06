package service

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/pb"
	"github.com/Chandler-WQ/experiment/util"
)

func GetSessionFromCtx(ctx *gin.Context) (*pb.Session, error) {
	sessionCtx, exists := ctx.Get(common.SessionInfo)
	if !exists {
		ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
		return nil, errors.New("[GetSessionInfo]sessionInfo is empty")
	}

	session, ok := sessionCtx.(*pb.Session)
	if !ok {
		return nil, errors.New("[GetSessionInfo]the type error")
	}
	return session, nil
}
