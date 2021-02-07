package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func CreateOrder(ctx *gin.Context) {
	req := model.ExperimentReserve{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsStudent(session.UserType) && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	err = db.Db.CreateExperimentReserveInfo(&model.ExperimentReserveInfo{
		ExperimentId: req.ExperimentId,
		UserId:       session.UserId,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Status:       -1,
	})

	if err != nil && strings.Contains(err.Error(), "RemainingSeat") {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "预约已满，请重选时间段"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}

func UpdateOrder(ctx *gin.Context) {
	req := model.UpdateExperimentReserve{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}

	experimentReserveInfo, err := db.Db.GetExperimentReserveInfo(req.ExperimentReserveId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	if experimentReserveInfo.UserId != session.UserId && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, "没有操作权限"))
		return
	}

	err = db.Db.UpdateExperimentReserveInfo(req.ExperimentReserveId, req.Status)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}

func MGetOrder(ctx *gin.Context) {
	startTimeStr := ctx.Params.ByName("start_time")
	if startTimeStr == "" {
		startTimeStr = "-1"
	}
	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "参数错误"))
		return
	}

	endTimeStr := ctx.Params.ByName("end_time")
	if endTimeStr == "" {
		endTimeStr = "0"
	}
	endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "参数错误"))
		return
	}
	if startTime > endTime {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "开始时间大于结束时间"))
		return
	}

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsStudent(session.UserType) && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	reserveInfoRsps, err := db.Db.MGetExperimentReserveInfo(session.UserId, startTime, endTime)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "查询为空"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, reserveInfoRsps))
}

func DeleteOrder(ctx *gin.Context) {
	req := model.DeleteExperimentReserve{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}

	experimentReserveInfo, err := db.Db.GetExperimentReserveInfo(req.ExperimentReserveId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	if experimentReserveInfo.UserId != session.UserId && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, "没有操作权限"))
		return
	}

	err = db.Db.DeleteExperimentReserveInfo(req.ExperimentReserveId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}
