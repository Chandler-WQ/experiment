package model

type Register struct {
	Name     string `json:"name" binding:"required"`
	Identity string `json:"identity" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Login struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateExperiment struct {
	ExperimentId int64   `json:"experiment_id" binding:"required"`
	ClassName    string  `json:"class_name" binding:"required"`
	StartTime    int64   `json:"start_time" binding:"required"`
	EndTime      int64   `json:"end_time" binding:"required"`
	Isopenpc     int32   `json:"isopenpc" binding:"required"`
	Resource     int32   `json:"resource" binding:"required"`
	Desc         string  `json:"desc" binding:"required"`
	ClassType    []int64 `json:"type"`
	StudentSum   int32   `json:"student_sum" binding:"required"`
}

type UpdateStudentCourse struct {
	StudentId int64 `json:"student_id" binding:"required"`
	ClassId   int64 `json:"class_id" binding:"required"`
	Score     int32 `json:"score"`
	Status    int32 `json:"status"`
}
