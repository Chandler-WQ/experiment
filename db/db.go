package db

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/util"
)

var _ DbController = (*dbProxy)(nil)

type DbController interface {

	//开设实验课，实验成绩相关的接口
	CreateCourse(course *model.ExperimentCourse, courseType []int64) error        //开设实验课
	MGetCourses(teacherId int64, startTime int64) ([]*model.CourseInfoRsp, error) //批量查询某个老师的实验课
	MGetCourseAllInfo(courseId int64) (*model.CourseAllInfo, error)               //批量查询某个实验课的所有学生关联的基本信息和老师的基本信息
	CreateStudentsCourse(studentIds []int64, teacherId, courseId int64) error     //插入学生与课程的信息
	UpdateStudentsCourse(studentCourse *model.StudentCourse) error                //更新学生某一个课程的信息

	//实验室管理
	CreateExperimentInfo(experimentInfo *model.ExperimentInfo) error                     // 创建实验室
	UpdateExperimentInfo(experimentId int64, experimentInfo *model.ExperimentInfo) error //更新实验室信息
	GetExperimentInfo(experimentId int64) (*model.ExperimentInfo, error)                 //查询实验室信息
	MgetExperimentInfo(offest, limit int) ([]model.ExperimentInfo, error)                //查询实验室信息

	//实验室学生预约占用管理
	GetOrCreateExperimentSegmentInfo(experimentSegmentInfos *model.ExperimentSegmentInfo) error
	CreateExperimentReserveInfo(experimentReserveInfo *model.ExperimentReserveInfo) error //预约实验室，并且同步更新segment和预约表
	UpdateExperimentReserveInfo(experimentReserveId, status int64) error
	GetExperimentSegmentInfo(experimentId int64, starTime int64, endTime int64) (*model.ExperimentSegmentInfo, error)        //获得某个实验室某个时间段的信息
	MGetExperimentSegmentInfo(starTime, endTime int64) ([]model.ExperimentSegmentInfo, error)                                //批量获取实验室某个时间段的信息
	MGetExperimentReserveInfo(userId int64, starTime, endTime int64) ([]model.ExperimentReserveInfo, *model.UserInfo, error) //获取某个人某个时间端的预约信息
	DeleteExperimentReserveInfo(experimentReserveId int64) error                                                             //删除预约，并且segment表进行相应的减一

	//设备管理
	CreateEquipmentInfo(equipmentInfo *model.EquipmentInfo) error                    // 创建设备
	UpdateEquipmentInfo(equipmentId int64, equipmentInfo *model.EquipmentInfo) error //更新实验室信息
	GetEquipmentInfo(equipmentId int64) (*model.EquipmentInfo, error)                //查询设备
	MgetEquipmentInfo(offest, limit int) ([]model.EquipmentInfo, error)              //批量查询设备

	//session管理
	CreateSession(session *model.Session) error         //创建session信息
	UpdateSession(session *model.Session) error         //更新session信息，失效时间
	DeleteSession(sessionId, userId int64) error        //删除session信息
	GetSession(sessionId int64) (*model.Session, error) //查询session信息
	MGetSession(userId int64) ([]model.Session, error)  //查询用户的所有session

	CreateUserInfo(userInfo *model.UserInfo) error
}

func (dbProxy *dbProxy) checkDb() {
	if dbProxy == nil || dbProxy.DbConn == nil {
		log.Error("the db is nil")
		panic("the db is nil")
	}
}

func (dbProxy *dbProxy) CreateCourse(course *model.ExperimentCourse, courseType []int64) error {
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

	_, startTime, endTime := util.SegmentTime(course.StartTime, course.EndTime)
	err = tx.Debug().Exec("update experiment_segment_info set remaining_seat = remaining_seat - ? where start_time >= ? and experiment_id = ? and end_time <= ?",
		course.StudentSum, startTime, course.ExperimentId, endTime).Error

	err = tx.Commit().Error
	if err != nil {
		log.Errorf("[CreateCourse]the db error ,the error is %v", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) MGetCourses(teacherId int64, startTime int64) ([]*model.CourseInfoRsp, error) {
	dbProxy.checkDb()
	var experimentCourses []*model.ExperimentCourse
	var err error
	if startTime > 0 {
		err = dbProxy.DbConn.Debug().Model(&model.ExperimentCourse{}).Where("start_time > ? and teacher_id = ?", startTime, teacherId).Find(&experimentCourses).Error
	} else {
		err = dbProxy.DbConn.Debug().Model(&model.ExperimentCourse{}).Where("teacher_id = ? ", teacherId).Find(&experimentCourses).Error
	}
	if err != nil {
		log.Errorf("[MGetCourses]the db error ,the error is %s", err.Error())
		return nil, err
	}
	courseInfoRsp := make([]*model.CourseInfoRsp, 0, len(experimentCourses))

	ExperimentIds := make([]int64, len(experimentCourses))
	for i := 0; i < len(experimentCourses); i++ {
		ExperimentIds[i] = experimentCourses[i].ExperimentId
		courseInfoRsp = append(courseInfoRsp, &model.CourseInfoRsp{
			ExperimentCourse: experimentCourses[i],
		})
	}

	var experimentInfos []*model.ExperimentInfo
	err = dbProxy.DbConn.Model(&model.ExperimentInfo{}).Where("id in ?", ExperimentIds).Find(&experimentInfos).Error
	if err != nil {
		log.Errorf("[MGetCourses]the db error ,the error is %s", err.Error())
		return nil, err
	}
	experimentInfoMap := make(map[int64]*model.ExperimentInfo, len(experimentInfos))
	for i := 0; i < len(experimentInfos); i++ {
		experimentInfoMap[int64(experimentInfos[i].ID)] = experimentInfos[i]
	}

	for i := 0; i < len(courseInfoRsp); i++ {
		courseInfoRsp[i].ExperimentInfo = experimentInfoMap[courseInfoRsp[i].ExperimentCourse.ExperimentId]
	}
	return courseInfoRsp, err
}

func (dbProxy *dbProxy) MGetCourseAllInfo(courseId int64) (*model.CourseAllInfo, error) {
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
	var experimentInfo model.ExperimentInfo
	err = dbProxy.DbConn.Model(&model.ExperimentInfo{}).Where("id = ?", experimentCourse.ExperimentId).First(&experimentInfo).Error
	if err != nil {
		log.Errorf("[MGetStudents]the db error ,the error is %s", err.Error())
		return nil, err
	}
	return &model.CourseAllInfo{
		ExperimentCourse:   &experimentCourse,
		ExperimentInfo:     &experimentInfo,
		CourseTeacherInfo:  &courseTeacherInfo,
		CourseStudentInfos: courseStudentInfos,
	}, nil
}

func (dbProxy *dbProxy) CreateStudentsCourse(studentIds []int64, teacherId, courseId int64) error {
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

func (dbProxy *dbProxy) UpdateStudentsCourse(studentCourse *model.StudentCourse) error {
	dbProxy.checkDb()
	tx := dbProxy.DbConn.Begin()
	if studentCourse.ExperimentCourseId <= 0 || studentCourse.StudentId <= 0 {
		return errors.New("the parameter error ,missing ExperimentCourseId or StudentId")
	}

	var err error
	if studentCourse.Score != -1 || studentCourse.ScoreEncry != "" {
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

func (dbProxy *dbProxy) CreateExperimentInfo(experimentInfo *model.ExperimentInfo) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.ExperimentInfo{}).Create(experimentInfo).Error
	if err != nil {
		log.Errorf("[CreateExperimentInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) UpdateExperimentInfo(experimentId int64, experimentInfo *model.ExperimentInfo) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.ExperimentInfo{}).Where("id = ?", experimentId).Updates(experimentInfo).Error
	if err != nil {
		log.Errorf("[UpdateExperimentInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) GetExperimentInfo(experimentId int64) (*model.ExperimentInfo, error) {
	dbProxy.checkDb()
	var experimentInfo model.ExperimentInfo
	err := dbProxy.DbConn.Model(&model.ExperimentInfo{}).Where("id = ?", experimentId).First(&experimentInfo).Error
	if err != nil {
		log.Errorf("[GetExperimentInfo]the db error ,the error is %s", err.Error())
	}
	return &experimentInfo, err
}

func (dbProxy *dbProxy) MgetExperimentInfo(offest, limit int) ([]model.ExperimentInfo, error) {
	dbProxy.checkDb()
	var experimentInfos []model.ExperimentInfo
	var err error
	if offest == -1 && limit == -1 {
		err = dbProxy.DbConn.Model(&model.ExperimentInfo{}).Find(&experimentInfos).Error
	} else {
		err = dbProxy.DbConn.Model(&model.ExperimentInfo{}).Offset(offest).Limit(limit).Find(&experimentInfos).Error
	}

	if err != nil {
		log.Errorf("[GetExperimentInfo]the db error ,the error is %s", err.Error())
	}
	return experimentInfos, err
}

func (dbProxy *dbProxy) GetOrCreateExperimentSegmentInfo(experimentSegmentInfo *model.ExperimentSegmentInfo) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.ExperimentSegmentInfo{}).
		Where("start_time = ? and experiment_id = ?", experimentSegmentInfo.StartTime, experimentSegmentInfo.ExperimentId).
		FirstOrCreate(experimentSegmentInfo).Error
	if err != nil {
		log.Errorf("[CreateExperimentSegmentInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) CreateExperimentReserveInfo(experimentReserveInfo *model.ExperimentReserveInfo) error {
	dbProxy.checkDb()
	tx := dbProxy.DbConn.Begin()
	err := tx.Model(&model.ExperimentReserveInfo{}).Create(&experimentReserveInfo).Error
	if err != nil {
		log.Errorf("[CreateExperimentReserveInfo]the db error ,the error is %s", err.Error())
		tx.Rollback()
		return err
	}
	sum, startTime, endTime := util.SegmentTime(experimentReserveInfo.StartTime, experimentReserveInfo.EndTime)
	for i := 0; i < int(sum); i++ {
		var experimentSegmentInfo model.ExperimentSegmentInfo
		start := startTime + int64(i*1800)
		err := dbProxy.DbConn.Model(&model.ExperimentSegmentInfo{}).
			Where("start_time = ? and experiment_id = ?", start, experimentReserveInfo.ExperimentId).
			First(&experimentSegmentInfo).Error
		if err != nil {
			log.Errorf("[CreateExperimentReserveInfo]the db error ,the error is %s", err.Error())
			tx.Rollback()
			return err
		}
		if experimentSegmentInfo.RemainingSeat <= 0 {
			return errors.New("the RemainingSeat is 0")
		}
	}

	err = tx.Exec("update  experiment_segment_info set remaining_seat = remaining_seat - ? where start_time >= ? and experiment_id = ? and end_time <= ?",
		1, startTime, experimentReserveInfo.ExperimentId, endTime).Error
	if err != nil {
		log.Errorf("[CreateExperimentReserveInfo]the db error ,the error is %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Errorf("[CreateExperimentReserveInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) UpdateExperimentReserveInfo(experimentReserveId, status int64) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.ExperimentReserveInfo{}).Where("id = ? ", experimentReserveId).Update("status = ?", status).Error
	if err != nil {
		log.Errorf("[UpdateExperimentReserveInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) GetExperimentSegmentInfo(experimentId, starTime, endTime int64) (*model.ExperimentSegmentInfo, error) {
	dbProxy.checkDb()
	_, start, end := util.SegmentTime(starTime, endTime)
	var experimentSegmentInfo model.ExperimentSegmentInfo
	err := dbProxy.DbConn.Debug().Model(&model.ExperimentSegmentInfo{}).
		Select("experiment_id,max(total_seat) as total_seat ,min(remaining_seat) as remaining_seat").
		Where("start_time >= ? and end_time <= ? and experiment_id = ? ", start, end, experimentId).
		Group("experiment_id").
		Find(&experimentSegmentInfo).Error
	if err != nil {
		log.Errorf("[GetExperimentSegmentInfo]the db error ,the error is %s", err.Error())
		return nil, err
	}

	return &model.ExperimentSegmentInfo{
		ExperimentId:  experimentSegmentInfo.ExperimentId,
		StartTime:     starTime,
		EndTime:       endTime,
		TotalSeat:     experimentSegmentInfo.TotalSeat,
		RemainingSeat: experimentSegmentInfo.RemainingSeat,
	}, nil
}

func (dbProxy *dbProxy) MGetExperimentSegmentInfo(starTime, endTime int64) ([]model.ExperimentSegmentInfo, error) {
	dbProxy.checkDb()
	var experimentSegmentInfos []model.ExperimentSegmentInfo
	_, start, end := util.SegmentTime(starTime, endTime)
	err := dbProxy.DbConn.Debug().Model(&model.ExperimentSegmentInfo{}).
		Select("experiment_id,max(total_seat) as total_seat ,min(remaining_seat) as remaining_seat").
		Where("start_time >= ? and end_time <= ?", start, end).
		Group("experiment_id").
		Find(&experimentSegmentInfos).Error
	if err != nil {
		log.Errorf("[MGetExperimentSegmentInfo]the db error ,the error is %s", err.Error())
	}
	for i := 0; i < len(experimentSegmentInfos); i++ {
		experimentSegmentInfos[i].StartTime = starTime
		experimentSegmentInfos[i].EndTime = endTime
	}
	return experimentSegmentInfos, err
}

func (dbProxy *dbProxy) MGetExperimentReserveInfo(userId int64, starTime int64, endTime int64) ([]model.ExperimentReserveInfo, *model.UserInfo, error) {
	dbProxy.checkDb()
	var experimentReserveInfos []model.ExperimentReserveInfo
	err := dbProxy.DbConn.Model(&model.ExperimentReserveInfo{}).
		Debug().Where("user_id = ? and start_time >= ? and end_time <= ?", userId, starTime, endTime).
		Find(&experimentReserveInfos).Error
	if err != nil {
		log.Errorf("[MGetExperimentReserveInfo]the db error ,the error is %s", err.Error())
		return nil, nil, err
	}

	var userInfo model.UserInfo
	err = dbProxy.DbConn.Model(&model.UserInfo{}).Where("id = ?", userId).First(&userInfo).Error
	if err != nil {
		log.Errorf("[MGetExperimentReserveInfo]the db error ,the error is %s", err.Error())
	}
	return experimentReserveInfos, &userInfo, err
}

func (dbProxy *dbProxy) DeleteExperimentReserveInfo(experimentReserveId int64) error {
	dbProxy.checkDb()
	var experimentReserveInfoTem model.ExperimentReserveInfo
	err := dbProxy.DbConn.Model(&model.ExperimentReserveInfo{}).Where("id = ?", experimentReserveId).First(&experimentReserveInfoTem).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		log.Errorf("[DeleteExperimentReserveInfo]the db error ,the error is %s", err.Error())
		return err
	}

	tx := dbProxy.DbConn.Begin()
	_, startTime, endTime := util.SegmentTime(experimentReserveInfoTem.StartTime, experimentReserveInfoTem.EndTime)
	err = tx.Model(&model.ExperimentReserveInfo{}).Where("id = ?", experimentReserveId).Delete(&model.ExperimentReserveInfo{}).Error
	if err != nil {
		log.Errorf("[DeleteExperimentReserveInfo]the db error ,the error is %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Exec("update  experiment_segment_info set remaining_seat = remaining_seat + ? where start_time >= ? and experiment_id = ? and end_time <= ?",
		1, startTime, experimentReserveInfoTem.ExperimentId, endTime).Error
	if err != nil {
		log.Errorf("[CreateExperimentReserveInfo]the db error ,the error is %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Errorf("[CreateExperimentReserveInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) CreateEquipmentInfo(equipmentInfo *model.EquipmentInfo) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.EquipmentInfo{}).Create(equipmentInfo).Error
	if err != nil {
		log.Errorf("[CreateEquipmentInfo]the db error ,the error is %s", err.Error())
	}
	return err

}

func (dbProxy *dbProxy) UpdateEquipmentInfo(equipmentId int64, equipmentInfo *model.EquipmentInfo) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.EquipmentInfo{}).Where("id = ?", equipmentId).Updates(equipmentInfo).Error
	if err != nil {
		log.Errorf("[UpdateEquipmentInfo]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) GetEquipmentInfo(equipmentId int64) (*model.EquipmentInfo, error) {
	dbProxy.checkDb()
	var equipmentInfo model.EquipmentInfo
	err := dbProxy.DbConn.Model(&model.EquipmentInfo{}).Where("id = ?", equipmentId).First(&equipmentInfo).Error
	if err != nil {
		log.Errorf("[GetEquipmentInfo]the db error ,the error is %s", err.Error())
	}
	return &equipmentInfo, err
}

func (dbProxy *dbProxy) MgetEquipmentInfo(offest, limit int) ([]model.EquipmentInfo, error) {
	dbProxy.checkDb()
	var equipmentInfos []model.EquipmentInfo

	if limit == 0 {
		limit = 100
	}
	err := dbProxy.DbConn.Model(&model.EquipmentInfo{}).Offset(offest).Limit(limit).Find(&equipmentInfos).Error
	if err != nil {
		log.Errorf("[MgetEquipmentInfo]the db error ,the error is %s", err.Error())
	}
	return equipmentInfos, err
}

func (dbProxy *dbProxy) CreateSession(session *model.Session) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Debug().Model(&model.Session{}).Create(session).Error
	if err != nil {
		log.Errorf("[CreateSession]the db error ,the error is %s", err.Error())
	}
	return err
}

func (dbProxy *dbProxy) UpdateSession(session *model.Session) error {
	dbProxy.checkDb()
	if session.SessionId != 0 {
		err := dbProxy.DbConn.Model(&model.Session{}).Where("session_id = ?", session.SessionId).
			Update("expire_time = ?", session.ExpireTime).Error
		if err != nil {
			log.Errorf("[UpdateSession]the db error ,the error is %s", err.Error())
		}
		return err
	}

	if session.UserId != 0 {
		err := dbProxy.DbConn.Model(&model.Session{}).Where("user_id = ?", session.UserId).
			Update("expire_time = ?", session.ExpireTime).Error
		if err != nil {
			log.Errorf("[UpdateSession]the db error ,the error is %s", err.Error())
		}
		return err
	}

	return errors.New("[UpdateSession]the userId is 0 and data is null")

}

func (dbProxy *dbProxy) DeleteSession(sessionId, userId int64) error {
	dbProxy.checkDb()
	if sessionId != 0 {
		err := dbProxy.DbConn.Model(&model.Session{}).Where("session_id = ?", sessionId).Delete(&model.Session{}).Error
		if err != nil {
			log.Errorf("[DeleteSession]the db error ,the error is %s", err.Error())
		}
		return err
	}

	if userId != 0 {
		err := dbProxy.DbConn.Model(&model.Session{}).Where("user_id = ?", userId).Delete(&model.Session{}).Error
		if err != nil {
			log.Errorf("[DeleteSession]the db error ,the error is %s", err.Error())
		}
		return err
	}

	return errors.New("[DeleteSession]the userId is 0 and data is null")

}

func (dbProxy *dbProxy) GetSession(sessionId int64) (*model.Session, error) {
	dbProxy.checkDb()
	var session model.Session
	err := dbProxy.DbConn.Model(&model.Session{}).Where("session_id = ?", sessionId).First(&session).Error
	if err != nil {
		log.Errorf("[GetSession]the db error ,the error is %s", err.Error())
	}
	return &session, err
}

func (dbProxy *dbProxy) MGetSession(userId int64) ([]model.Session, error) {
	dbProxy.checkDb()
	var sessions []model.Session
	err := dbProxy.DbConn.Model(&model.Session{}).Where("user_id = ?", userId).Find(&sessions).Error
	if err != nil {
		log.Errorf("[MGetSession]the db error ,the error is %s", err.Error())
	}
	return sessions, err
}

func (dbProxy *dbProxy) CreateUserInfo(userInfo *model.UserInfo) error {
	dbProxy.checkDb()
	err := dbProxy.DbConn.Model(&model.UserInfo{}).Create(userInfo).Error
	if err != nil {
		log.Errorf("[MGetSession]the db error ,the error is %s", err.Error())
	}

	return err
}

func (dbProxy *dbProxy) GetUserInfo(passportName, password string) (*model.UserInfo, error) {
	dbProxy.checkDb()
	var userInfo model.UserInfo
	err := dbProxy.DbConn.Model(&model.UserInfo{}).Where("passport_name = ? and password = ?", passportName, password).First(&userInfo).Error
	if err != nil {
		log.Errorf("[MGetSession]the db error ,the error is %s", err.Error())
	}
	return &userInfo, err
}
