package db

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"

	"github.com/Chandler-WQ/experiment/common/model"
)

func TestMain(m *testing.M) {
	MustInitDb()
	m.Run()
}

func TestFirst(t *testing.T) {
	//Db.DbConn.AutoMigrate(&model.ExperimentReserveInfo{})
	//Db.DbConn.AutoMigrate(&model.ExperimentSegmentInfo{})
	//Db.DbConn.AutoMigrate(&model.EquipmentInfo{})
	//Db.DbConn.AutoMigrate(&model.ExperimentInfo{})
	//Db.DbConn.AutoMigrate(&model.ExperimentCourse{})
	//Db.DbConn.AutoMigrate(&model.StudentCourse{})
	//Db.DbConn.AutoMigrate(&model.TeacherCourse{})
	//timeNow := time.Now().Unix()
	//a := timeNow % 3600
	//b := timeNow - a + 3600
	//Db.DbConn.Create(&model.ExperimentCourse{
	//	Name:         "电工实验",
	//	ExperimentId: 1,
	//	StartTime:    b,
	//	EndTime:      b + 3600*2,
	//	Isopenpc:     1,
	//	Resource:     1,
	//	Desc:         "无",
	//	TeacherId:    1,
	//})
	//Db.DbConn.Create(&model.ExperimentInfo{
	//	Name:      "电工实验",
	//	Region:    "实验楼304",
	//	TotalSeat: 100,
	//})

	//
	//Db.DbConn.Create(&model.ExperimentReserveInfo{
	//	UserId:       1,
	//	ExperimentId: 1,
	//	StartTime:    b,
	//	EndTime:      b + 3600*2,
	//	Status:       1,
	//})
	//Db.DbConn.Create(&model.EquipmentInfo{
	//	Name:           "烧杯",
	//	Type:           "玻璃制品",
	//	Factory:        "未知",
	//	ExperimentName: "电工实验",
	//	EquipmentId:    1,
	//})
	//Db.DbConn.Create(&model.ExperimentSegmentInfo{
	//	ExperimentId:  1,
	//	StartTime:     b,
	//	EndTime:       b + 3600*2,
	//	TotalSeat:     100,
	//	RemainingSeat: 99,
	//})
	//Db.DbConn.Create(&model.ExperimentType{
	//	ExperimentCourseId:   1,
	//	ExperimentCourseType: 2,
	//})
	//Db.DbConn.Create(&model.TeacherCourse{
	//	ExperimentCourseId: 1,
	//	TeacherId:          1,
	//})
	//Db.DbConn.Create(&model.StudentCourse{
	//	ExperimentCourseId: 1,
	//	TeacherId:          1,
	//	StudentId:          1,
	//	Score:              99,
	//})
	//Db.DbConn.Create(&model.UserInfo{
	//	Name:         "张三",
	//	Email:        "123456@qq.com",
	//	PassportName: "zhangsan",
	//	Phone:        "12345678911",
	//	Password:     "teststt",
	//	College:      "计算机院",
	//	UserType:     1,
	//	UserNumber:   2,
	//})
}

func TestDbProxy_CreateCourse(t *testing.T) {
	err := Db.CreateCourse(&gin.Context{}, &model.ExperimentCourse{
		Name:         "物理实验",
		ExperimentId: 1,
		StartTime:    1612447200,
		EndTime:      1612454400,
		Isopenpc:     0,
		Resource:     1,
		Desc:         "无",
		TeacherId:    1,
	}, []int64{1})
	assert.Nil(t, err)
}

func TestDbProxy_MGetCourses(t *testing.T) {
	ExperimentCourse, ExperimentInfo, err := Db.MGetCourses(&gin.Context{}, 1)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(ExperimentCourse))
	assert.NotNil(t, ExperimentInfo)
	log.Infof("the ExperimentInfo is %+v", ExperimentInfo)
}

func TestDbProxy_MGetCourseAllInfo(t *testing.T) {
	courseAllInfo, err := Db.MGetCourseAllInfo(&gin.Context{}, 1)
	assert.Nil(t, err)
	assert.NotNil(t, courseAllInfo)
	log.Infof("the ExperimentInfo is %+v", courseAllInfo.ExperimentCourse)
	log.Infof("the CourseTeacherInfo is %+v", courseAllInfo.CourseTeacherInfo)
	log.Infof("the CourseStudentInfos[0] is %+v", courseAllInfo.CourseStudentInfos[0])
}

func TestDbProxy_CreateStudentsCourse(t *testing.T) {
	err := Db.CreateStudentsCourse(&gin.Context{}, []int64{3}, 2, 1)
	assert.Nil(t, err)
}

func TestDbProxy_UpdateStudentsCourse(t *testing.T) {
	err := Db.UpdateStudentsCourse(&gin.Context{}, &model.StudentCourse{
		ExperimentCourseId: 1,
		StudentId:          1,
		Score:              100,
		ScoreEncry:         "xxkahsd",
		Status:             1,
	})
	assert.Nil(t, err)
}
