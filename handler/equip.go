package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func CreateEquip(ctx *gin.Context) {
	req := model.Equipment{
		ExperimentId: -1,
	}
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

	err = db.Db.CreateEquipmentInfo(&model.EquipmentInfo{
		Name:           req.Name,
		ExperimentId:   req.ExperimentId,
		ExperimentName: req.ExperimentName,
		Factory:        req.Factory,
		Type:           req.Type,
		Sum:            req.Sum,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))

}

func UpdateEquip(ctx *gin.Context) {
	req := model.UpdateEquipment{}
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

	err = db.Db.UpdateEquipmentInfo(req.EquipmentId, &model.EquipmentInfo{
		Sum: req.Sum,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}

func MGetEquips(ctx *gin.Context) {
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

	equipmentInfos, err := db.Db.MgetEquipmentInfo(int(offset), int(limit))
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "查询数据为空"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, equipmentInfos))
}

func GetEquip(ctx *gin.Context) {
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

	equipmentInfos, err := db.Db.GetEquipmentInfo(id)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "查询数据为空"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, equipmentInfos))
}

func DeleteEquips(ctx *gin.Context) {
	req := model.DeleteEquipment{}
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

	err = db.Db.DeleteEquipmentInfo(req.EquipmentId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}
