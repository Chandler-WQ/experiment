package service

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/db"
)

func InitTimeRun() {
	SegmentCreatefunc := func() {
		defer func() {
			if s := recover(); s != nil {
				log.Errorf("[SegmentCreatefunc]panic: %s", s)
			}
		}()
		_ = SegmentCreate()
	}

	SegmentCreatefunc()

	go func() {
		for {
			time.Sleep(5 * time.Hour)
			SegmentCreatefunc()
		}
	}()
}

func SegmentCreate() error {
	timeNow := time.Now().Unix()
	timeStart := timeNow - timeNow%(3600*24*2)
	experimentInfos, err := db.Db.MgetExperimentInfo(-1, -1)

	if err != nil {
		return err
	}

	for i := 0; i < len(experimentInfos); i++ {
		for j := 0; j < 96; j++ {
			experimentSegmentInfo := model.ExperimentSegmentInfo{
				ExperimentId:  int64(experimentInfos[i].ID),
				StartTime:     timeStart + int64(j*30*60),
				EndTime:       timeStart + int64(j*30*60) + 30*60,
				TotalSeat:     experimentInfos[i].TotalSeat,
				RemainingSeat: experimentInfos[i].TotalSeat,
			}
			err := db.Db.GetOrCreateExperimentSegmentInfo(&experimentSegmentInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
