package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func CreateExperiment(ctx *gin.Context) {
	req := model.Experiment{}
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
	if !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	err = db.Db.CreateExperimentInfo(&model.ExperimentInfo{
		Name:      req.Name,
		Region:    req.Region,
		TotalSeat: req.TotalSeat,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	go func() {
		defer func() {
			if s := recover(); s != nil {
				log.Errorf("[SegmentCreatefunc]panic: %s", s)
			}
		}()
		_ = service.SegmentCreate()
	}()
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))

}

func UpdateExperiment(ctx *gin.Context) {
	req := model.UpdateExperiment{}
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
	if !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	err = db.Db.UpdateExperimentInfo(req.ExperimentId, &model.ExperimentInfo{
		TotalSeat: req.TotalSeat,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}

func MGetExperiments(ctx *gin.Context) {
	offsetStr := ctx.Params.ByName("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, errA := strconv.ParseInt(offsetStr, 10, 64)
	limitStr := ctx.Params.ByName("limit")
	if limitStr == "" {
		limitStr = "100"
	}
	limit, errB := strconv.ParseInt(limitStr, 10, 64)
	if errA != nil || errB != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}
	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	ExperimentInfos, err := db.Db.MgetExperimentInfo(int(offset), int(limit))
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	if len(ExperimentInfos) == 0 {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "查询数据为空"))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, ExperimentInfos))
}

func GetExperiment(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}
	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	ExperimentInfos, err := db.Db.GetExperimentInfo(id)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "查询数据为空"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, ExperimentInfos))
}

func DeleteExperiments(ctx *gin.Context) {
	req := model.DeleteExperiment{}
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
	if !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	err = db.Db.DeleteExperimentInfo(req.ExperimentId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}
