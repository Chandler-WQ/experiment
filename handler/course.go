package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func CreatCourse(ctx *gin.Context) {
	req := model.CreateExperiment{
		Isopenpc:   0,
		Resource:   0,
		StudentIds: []int64{},
	}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, req))
		return
	}

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsTeacher(session.UserType) && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	if req.StartTime > req.EndTime || req.StartTime < time.Now().Unix() || !util.IsHalfTime(req.StartTime) || !util.IsHalfTime(req.EndTime) || req.StartTime > time.Now().Unix()+common.Week {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "时间错误"))
		return
	}

	experimentSegmentInfo, err := db.Db.GetExperimentSegmentInfo(req.ExperimentId, req.StartTime, req.EndTime)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "实验室id不存在"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	if experimentSegmentInfo.RemainingSeat < int64(len(req.StudentIds)) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "实验室剩余座位数量不够"))
		return
	}

	experimentCourse := model.ExperimentCourse{
		ExperimentId: req.ExperimentId,
		Name:         req.ClassName,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Isopenpc:     req.Isopenpc,
		Resource:     req.Resource,
		Desc:         req.Desc,
		TeacherId:    session.UserId,
		StudentSum:   int32(len(req.StudentIds)),
	}

	err = db.Db.CreateCourse(&experimentCourse, req.ClassType, req.StudentIds)

	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}

func MGetCourseInfo(ctx *gin.Context) {
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
	if !util.IsTeacher(session.UserType) && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	CourseInfoRsp, err := db.Db.MGetCourses(session.UserId, startTime, endTime)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "查询为空"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, CourseInfoRsp))
}

func GetCourseAllInfo(ctx *gin.Context) {
	courseIdStr := ctx.Params.ByName("course_id")
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "参数错误"))
		return
	}

	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsTeacher(session.UserType) && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}

	courseAllInfo, err := db.Db.MGetCourseAllInfo(courseId)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, "课程id不存在"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, courseAllInfo))
}

func UpdateStudentCourse(ctx *gin.Context) {
	req := model.UpdateStudentCourse{
		Score:  -1,
		Status: -1,
	}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}
	if req.Score == -1 && req.Status == -1 {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.ParaErr.Code, common.ParaErr.Message, nil))
		return
	}
	session, err := service.GetSessionFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.SerErr.Code, common.SerErr.Message, err.Error()))
		return
	}
	if !util.IsTeacher(session.UserType) && !util.IsAdmin(session.UserType) {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.PermissionErr.Code, common.PermissionErr.Message, nil))
		return
	}
	//todo:需要做权限隔离

	err = db.Db.UpdateStudentsCourse(&model.StudentCourse{
		StudentId:          req.StudentId,
		ExperimentCourseId: req.ClassId,
		Score:              req.Score,
		Status:             req.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailResponse(ctx, common.DBErr.Code, common.DBErr.Message, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse(ctx, common.SUCCESS, nil))
}
