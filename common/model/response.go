package model

type CourseAllInfo struct {
	ExperimentCourse   *ExperimentCourse
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
