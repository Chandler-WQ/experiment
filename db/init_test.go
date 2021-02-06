package db

import (
	"testing"
	"time"

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
	err := Db.CreateCourse(&model.ExperimentCourse{
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
	ExperimentCourse, err := Db.MGetCourses(1, 1)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(ExperimentCourse))
	assert.NotNil(t, ExperimentCourse)
	log.Infof("the ExperimentInfo is %+v", ExperimentCourse)
}

func TestDbProxy_MGetCourseAllInfo(t *testing.T) {
	courseAllInfo, err := Db.MGetCourseAllInfo(1)
	assert.Nil(t, err)
	assert.NotNil(t, courseAllInfo)
	log.Infof("the ExperimentInfo is %+v", courseAllInfo.ExperimentCourse)
	log.Infof("the CourseTeacherInfo is %+v", courseAllInfo.CourseTeacherInfo)
	log.Infof("the CourseStudentInfos[0] is %+v", courseAllInfo.CourseStudentInfos[0])
}

func TestDbProxy_CreateStudentsCourse(t *testing.T) {
	err := Db.CreateStudentsCourse([]int64{3}, 2, 1)
	assert.Nil(t, err)
}

func TestDbProxy_UpdateStudentsCourse(t *testing.T) {
	err := Db.UpdateStudentsCourse(&model.StudentCourse{
		ExperimentCourseId: 1,
		StudentId:          1,
		Score:              100,
		ScoreEncry:         "xxkahsd",
		Status:             1,
	})
	assert.Nil(t, err)
}

func TestDbProxy_CreateExperimentInfo(t *testing.T) {
	err := Db.CreateExperimentInfo(&model.ExperimentInfo{
		Name:      "物理实验室",
		Region:    "物理楼118",
		TotalSeat: 210,
	})
	assert.Nil(t, err)
}

func TestDbProxy_UpdateExperimentInfo(t *testing.T) {
	err := Db.UpdateExperimentInfo(2, &model.ExperimentInfo{
		TotalSeat: 100,
	})
	assert.Nil(t, err)
}

func TestDbProxy_GetExperimentInfo(t *testing.T) {
	experimentInfo, err := Db.GetExperimentInfo(2)
	assert.Nil(t, err)
	log.Infof("the experimentInfo is %+v", experimentInfo)
}

func TestDbProxy_MgetExperimentInfo(t *testing.T) {
	experimentInfos, err := Db.MgetExperimentInfo(1, 2)
	assert.Nil(t, err)
	log.Infof("the experimentInfo is %+v,the len is %v", experimentInfos[0], len(experimentInfos))
}

func TestDbProxy_CreateExperimentReserveInfo(t *testing.T) {
	timeNow := time.Now().Unix()

	err := Db.CreateExperimentReserveInfo(&model.ExperimentReserveInfo{
		UserId:       1,
		ExperimentId: 2,
		StartTime:    timeNow,
		EndTime:      timeNow + 30*61*3,
		Status:       1,
	})
	assert.Nil(t, err)
}

func TestDbProxy_DeleteExperimentReserveInfo(t *testing.T) {
	err := Db.DeleteExperimentReserveInfo(2)
	assert.Nil(t, err)
}

func TestDbProxy_MGetExperimentReserveInfo(t *testing.T) {
	timeNow := time.Now().Unix()
	experimentReserveInfos, userInfo, err := Db.MGetExperimentReserveInfo(1, timeNow-30*61, timeNow+30*61*20)
	assert.Nil(t, err)
	log.Infof("userInfo is %v", userInfo)
	log.Infof("experimentReserveInfos is %v", experimentReserveInfos[0])
}

func TestDbProxy_GetExperimentSegmentInfo(t *testing.T) {
	timeNow := time.Now().Unix()
	experimentSegmentInfo, err := Db.GetExperimentSegmentInfo(1, timeNow-30*61, timeNow+30*61*20)
	assert.Nil(t, err)
	log.Infof("experimentSegmentInfo is %v", experimentSegmentInfo)
}

func TestDbProxy_MGetExperimentSegmentInfo(t *testing.T) {
	timeNow := time.Now().Unix()
	experimentSegmentInfo, err := Db.MGetExperimentSegmentInfo(timeNow-30*61, timeNow+30*61*20)
	assert.Nil(t, err)
	log.Infof("experimentSegmentInfo is %v", experimentSegmentInfo[1])
}
