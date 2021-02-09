package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func Register(ctx *gin.Context) {
	var req model.Register
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, req))
		return
	}
	userType, err := util.IdentityToCode(req.Identity)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, err.Error(), req))
		return
	}
	userinfo := model.UserInfo{
		Name:         req.Name,
		PassportName: req.UserName,
		Password:     req.Password,
		UserType:     userType,
	}

	err = db.Db.CreateUserInfo(&userinfo)
	//代表用户名已经被注册过
	if err != nil && strings.Contains(err.Error(), "Duplicate") {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, "用户名已经被占用", nil))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, ""))
}

func Login(ctx *gin.Context) {
	var req model.Login
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, req))
		return
	}

	userInfo, err := db.Db.GetUserInfo(req.UserName, req.Password)
	//代表用户名已经被注册过
	if err != nil && err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, "用户名或密码错误", nil))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	session := userInfo.ToSession()
	session.SessionId = util.GetId()
	session.ExpireTime = time.Now().Unix() + common.SessionAge
	token, err := service.CreateSession(&session)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, err.Error(), nil))
		return
	}
	util.SetCookie(ctx, common.SessionKey, token)
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, userInfo))
}

func Logout(ctx *gin.Context) {

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}

	err = db.Db.DeleteSession(session.SessionId, session.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	util.SetCookie(ctx, common.SessionKey, "")
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}

func GetSessionInfo(ctx *gin.Context) {
	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, nil))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, session))
}
