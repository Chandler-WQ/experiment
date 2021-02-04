package db

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Chandler-WQ/experiment/common/model"
)

type DbController interface {

	//开设实验课，实验成绩相关的接口
	CreateCourse(ctx *gin.Context, course *model.ExperimentCourse, courseType []int64) error                                 //开设实验课
	MGetCourses(ctx *gin.Context, teacherId int64, startTime int64) ([]model.ExperimentCourse, *model.ExperimentInfo, error) //批量查询某个老师的实验课
	MGetCourseAllInfo(ctx *gin.Context, courseId int64) (*model.CourseAllInfo, error)                                        //批量查询某个实验课的所有学生关联的基本信息和老师的基本信息
	CreateStudentsCourse(ctx *gin.Context, studentIds []int64, teacherId, courseId int64) error                              //插入学生与课程的信息
	UpdateStudentsCourse(ctx *gin.Context, studentCourse *model.StudentCourse) error                                         //更新学生某一个课程的信息

	//实验室管理
	CreateExperimentInfo(ctx *gin.Context, experimentInfo *model.ExperimentInfo) error
	UpdateExperimentInfo(ctx *gin.Context, experimentInfo *model.ExperimentInfo) error
	GetExperimentInfo(ctx *gin.Context, experimentId int64) (*model.ExperimentInfo, error)
	MgetExperimentInfo(ctx *gin.Context) ([]model.ExperimentInfo, error)

	//实验室学生预约占用管理
	MGetExperimentSegmentInfo(ctx *gin.Context, starTime int64, endTime int64) ([]model.ExperimentSegmentInfo, error)
	GetExperimentSegmentInfo(ctx *gin.Context, experimentId int64, starTime int64, endTime int64) ([]model.ExperimentSegmentInfo, error)
	CreateExperimentReserveInfo(ctx *gin.Context, experimentReserveInfo *model.ExperimentReserveInfo) error
	DeleteExperimentReserveInfo(ctx *gin.Context, experimentReserveInfo *model.ExperimentReserveInfo) error
	MGetExperimentReserveInfo(ctx *gin.Context, userId int64) ([]model.ExperimentReserveInfo, *model.UserInfo, error)

	//设备管理
	CreateEquipmentInfo(ctx *gin.Context, equipmentInfo *model.EquipmentInfo) error
	UpdateEquipmentInfo(ctx *gin.Context, equipmentInfo *model.EquipmentInfo) error
	GetEquipmentInfo(ctx *gin.Context, experimentId int64) (*model.EquipmentInfo, error)
	MgetEquipmentInfo(ctx *gin.Context) ([]model.EquipmentInfo, error)
}

func (dbProxy *dbProxy) checkDb() {
	if dbProxy == nil || dbProxy.DbConn == nil {
		log.Error("the db is nil")
		panic("the db is nil")
	}
}

func (dbProxy *dbProxy) CreateCourse(ctx *gin.Context, course *model.ExperimentCourse, courseType []int64) error {
	dbProxy.checkDb()
	tx := dbProxy.DbConn.Begin()
	err := tx.Create(&course).Error
	if err != nil {
		log.Errorf("[CreateCourse]the db error ,the error is %v", err.Error())
		tx.Rollback()
		return err
	}

	for _, i := range courseType {
		err := tx.Create(&model.ExperimentType{
			ExperimentCourseId:   int64(course.ID),
			ExperimentCourseType: i,
		}).Error
		if err != nil {
			log.Errorf("[CreateCourse]the db error ,the error is %v", err.Error())
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		log.Errorf("[CreateCourse]the db error ,the error is %v", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) MGetCourses(ctx *gin.Context, teacherId int64, startTime int64) ([]model.ExperimentCourse, *model.ExperimentInfo, error) {
	dbProxy.checkDb()
	var experimentCourses []model.ExperimentCourse
	var err error
	if startTime > 0 {
		err = dbProxy.DbConn.Model(&model.ExperimentCourse{}).Where("startTime > ? and teacher_id = ?", startTime, teacherId).Find(&experimentCourses).Error
	} else {
		err = dbProxy.DbConn.Model(&model.ExperimentCourse{}).Where("teacher_id = ? ", teacherId).Find(&experimentCourses).Error
	}
	if err != nil {
		log.Errorf("[MGetCourses]the db error ,the error is %s", err.Error())
		return nil, nil, err
	}

	var experimentInfo model.ExperimentInfo
	err = dbProxy.DbConn.Model(&model.ExperimentInfo{}).Where("id = ?", experimentCourses[0].ExperimentId).First(&experimentInfo).Error
	if err != nil {
		log.Errorf("[MGetCourses]the db error ,the error is %s", err.Error())
		return experimentCourses, nil, err
	}
	return experimentCourses, &experimentInfo, err
}

func (dbProxy *dbProxy) MGetCourseAllInfo(ctx *gin.Context, courseId int64) (*model.CourseAllInfo, error) {
	dbProxy.checkDb()
	var experimentCourse model.ExperimentCourse
	err := dbProxy.DbConn.Model(&model.ExperimentCourse{}).Where("id = ?", courseId).First(&experimentCourse).Error
	if err != nil {
		log.Errorf("[MGetStudents]the db error ,the error is %s", err.Error())
		return nil, err
	}
	courseStudentInfos := []model.CourseStudentInfo{}
	err = dbProxy.DbConn.Debug().Table("student_course").
		Select("user_info.id as id ,user_info.name as name,user_info.college as college ,user_info.user_number as user_number,"+
			"student_course.status as status ,student_course.score as score,student_course.score_encry as score_encry").
		Joins("left join user_info on user_info.id = student_course.student_id").
		Where("student_course.experiment_course_id = ?", experimentCourse.ID).
		Find(&courseStudentInfos).Error
	if err != nil {
		log.Errorf("[MGetStudents]the db error ,the error is %s", err.Error())
		return nil, err
	}

	courseTeacherInfo := model.CourseTeacherInfo{}
	err = dbProxy.DbConn.Debug().Table("experiment_course").
		Select("user_info.id as id ,user_info.name as name,user_info.college as college ,user_info.user_number as user_number").
		Joins("left join user_info on user_info.id = experiment_course.teacher_id").
		Where("experiment_course.id = ?", courseId).
		First(&courseTeacherInfo).Error

	if err != nil {
		log.Errorf("[MGetStudents]the db error ,the error is %s", err.Error())
		return nil, err
	}
	return &model.CourseAllInfo{
		ExperimentCourse:   &experimentCourse,
		CourseTeacherInfo:  &courseTeacherInfo,
		CourseStudentInfos: courseStudentInfos,
	}, nil
}

func (dbProxy *dbProxy) CreateStudentsCourse(ctx *gin.Context, studentIds []int64, teacherId, courseId int64) error {
	dbProxy.checkDb()
	length := len(studentIds)
	studentCourses := make([]model.StudentCourse, length)
	for i := 0; i < length; i++ {
		studentCourses[i].StudentId = studentIds[i]
		studentCourses[i].TeacherId = teacherId
		studentCourses[i].ExperimentCourseId = courseId
		studentCourses[i].Score = -1
	}
	err := dbProxy.DbConn.Create(&studentCourses).Error
	if err != nil {
		log.Errorf("[CreateStudentsCourse]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) UpdateStudentsCourse(ctx *gin.Context, studentCourse *model.StudentCourse) error {
	dbProxy.checkDb()
	tx := dbProxy.DbConn.Begin()
	if studentCourse.ExperimentCourseId <= 0 || studentCourse.StudentId <= 0 {
		return errors.New("the parameter error ,missing ExperimentCourseId or StudentId")
	}

	var err error
	if studentCourse.Score != -1 && studentCourse.ScoreEncry != "" {
		err = tx.Model(&model.StudentCourse{}).Where("experiment_course_id = ? And student_id = ?", studentCourse.ExperimentCourseId, studentCourse.StudentId).
			Updates(model.StudentCourse{
				ScoreEncry: studentCourse.ScoreEncry,
				Score:      studentCourse.Score,
			}).Error
		if err != nil {
			log.Errorf("[CreateStudentsCourse]the db error ,the error is %s", err.Error())
			tx.Rollback()
			return err
		}
	}
	if studentCourse.Status != -1 {
		err = tx.Model(&model.StudentCourse{}).Where("experiment_course_id = ? And student_id = ?", studentCourse.ExperimentCourseId, studentCourse.StudentId).
			Update("status", studentCourse.Status).Error
		if err != nil {
			log.Errorf("[CreateStudentsCourse]the db error ,the error is %s", err.Error())
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		log.Errorf("[CreateStudentsCourse]the db error ,the error is %s", err.Error())
	}
	return err
}
