package model

import (
	"gorm.io/gorm"
)

//实验课信息
type ExperimentCourse struct {
	gorm.Model
	Name         string
	ExperimentId int64
	StartTime    int64
	EndTime      int64
	Isopenpc     int32
	Resource     int32
	Desc         string
	TeacherId    int64
}

func (ExperimentCourse) TableName() string {
	return "experiment_course"
}

//实验课类型信息
type ExperimentType struct {
	gorm.Model
	ExperimentCourseId   int64
	ExperimentCourseType int64
}

func (ExperimentType) TableName() string {
	return "experiment_type"
}

//session信息
type Session struct {
	gorm.Model
	UserId     int64
	ExpireTime int64
	Data       string
}

func (Session) TableName() string {
	return "session"
}

//学生-学生-实验课关联信息
type StudentCourse struct {
	gorm.Model
	ExperimentCourseId int64
	TeacherId          int64
	StudentId          int64
	Score              int32
	ScoreEncry         string
	Status             int32
}

func (StudentCourse) TableName() string {
	return "student_course"
}

//用户信息
type UserInfo struct {
	gorm.Model
	Name         string
	Email        string
	Phone        string
	Password     string
	PassportName string
	College      string
	UserType     int32
	UserNumber   int64
}

func (UserInfo) TableName() string {
	return "user_info"
}

//实验室信息信息
type ExperimentInfo struct {
	gorm.Model
	Name      string
	Region    string
	TotalSeat int64
}

func (ExperimentInfo) TableName() string {
	return "experiment_info"
}

//实验室预约信息
type ExperimentSegmentInfo struct {
	gorm.Model
	ExperimentId  int64
	StartTime     int64
	EndTime       int64
	TotalSeat     int64
	RemainingSeat int64
}

func (ExperimentSegmentInfo) TableName() string {
	return "experiment_segment_info"
}

//实验室和学生关联占用的信息
type ExperimentReserveInfo struct {
	gorm.Model
	UserId       int64
	ExperimentId int64
	StartTime    int64
	EndTime      int64
	Status       int32
}

func (ExperimentReserveInfo) TableName() string {
	return "experiment_reserve_info"
}

//设备信息
type EquipmentInfo struct {
	gorm.Model
	Name           string
	Type           string
	Factory        string
	ExperimentName string
	EquipmentId    int64
}

func (EquipmentInfo) TableName() string {
	return "equipment_info"
}
