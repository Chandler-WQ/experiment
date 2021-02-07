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
	Isopenpc     int32   `json:"isopenpc"`
	Resource     int32   `json:"resource"`
	Desc         string  `json:"desc" binding:"required"`
	ClassType    []int64 `json:"type"`
	StudentIds   []int64 `json:"student_ids"`
}

type UpdateStudentCourse struct {
	StudentId int64 `json:"student_id" binding:"required"`
	ClassId   int64 `json:"class_id" binding:"required"`
	Score     int32 `json:"score"`
	Status    int32 `json:"status"`
}

type Equipment struct {
	Name           string `json:"equipname" binding:"required"`
	Type           string `json:"type" binding:"required"`
	Factory        string `json:"factory" binding:"required"`
	ExperimentName string `json:"lab" binding:"required"`
	ExperimentId   int64  `json:"experiment_id"`
	Sum            int64  `json:"sum" binding:"required"`
}

type UpdateEquipment struct {
	EquipmentId int64 `json:"equipment_id" binding:"required"`
	Sum         int64 `json:"sum" binding:"required"`
}

type DeleteEquipment struct {
	EquipmentId int64 `json:"equipment_id" binding:"required"`
}

type Experiment struct {
	Name      string `json:"name" binding:"required"`
	Region    string `json:"region" binding:"required"`
	TotalSeat int64  `json:"total_seat" binding:"required"`
}

type UpdateExperiment struct {
	ExperimentId int64 `json:"experiment_id" binding:"required"`
	TotalSeat    int64 `json:"total_seat" binding:"required"`
}

type DeleteExperiment struct {
	ExperimentId int64 `json:"experiment_id" binding:"required"`
}

type ExperimentReserve struct {
	ExperimentId int64 `json:"experiment_id" binding:"required"`
	StartTime    int64 `json:"start_time" binding:"required"`
	EndTime      int64 `json:"end_time" binding:"required"`
}

type UpdateExperimentReserve struct {
	ExperimentReserveId int64 `json:"experiment_reserve_id" binding:"required"`
	Status              int64 `json:"status" binding:"required"`
}

type DeleteExperimentReserve struct {
	ExperimentReserveId int64 `json:"experiment_reserve_id" binding:"required"`
}
