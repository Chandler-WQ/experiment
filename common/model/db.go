package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common/pb"
)

type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

//实验课信息
type ExperimentCourse struct {
	Model
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
	Model
	ExperimentCourseId   int64
	ExperimentCourseType int64
}

func (ExperimentType) TableName() string {
	return "experiment_type"
}

//session信息
type Session struct {
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	SessionId  int64
	UserId     int64
	ExpireTime int64
	Data       string
}

func (Session) TableName() string {
	return "session"
}

//学生-学生-实验课关联信息
type StudentCourse struct {
	Model
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
	Model
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
	Model
	Name      string
	Region    string
	TotalSeat int64
}

func (ExperimentInfo) TableName() string {
	return "experiment_info"
}

//实验室预约信息
type ExperimentSegmentInfo struct {
	Model
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
	Model
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
	Model
	Name           string
	Type           string
	Factory        string
	ExperimentName string
	EquipmentId    int64
}

func (EquipmentInfo) TableName() string {
	return "equipment_info"
}

func (userInfo *UserInfo) ToSession() pb.Session {
	return pb.Session{
		UserNumber:   userInfo.UserNumber,
		UserId:       int64(userInfo.ID),
		Name:         userInfo.Name,
		PassportName: userInfo.PassportName,
		College:      userInfo.College,
		UserType:     int64(userInfo.UserType),
	}

}
