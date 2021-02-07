package model

type CourseAllInfo struct {
	ExperimentCourse   *ExperimentCourse
	ExperimentInfo     *ExperimentInfo
	CourseStudentInfos []CourseStudentInfo
	CourseTeacherInfo  *CourseTeacherInfo
}

type CourseTeacherInfo struct {
	Id         int64
	Name       string
	College    string
	UserNumber int64
}

type CourseStudentInfo struct {
	Id         int64
	Score      int32
	ScoreEncry string
	Status     int32
	Name       string
	College    string
	UserNumber int64
}

type CourseInfoRsp struct {
	ExperimentCourse *ExperimentCourse
	ExperimentInfo   *ExperimentInfo
}

type ReserveInfoRsp struct {
	ExperimentReserveInfo *ExperimentReserveInfo
	ExperimentInfo        *ExperimentInfo
}
