package util

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//global var
var sequence int64 = 0
var lastTime int64 = -1

//every segment bit
var workerIdBits = 5
var datacenterIdBits = 5
var sequenceBits = 12

//every segment max number
var maxWorkerId int64 = -1 ^ (-1 << workerIdBits)
var maxDatacenterId int64 = -1 ^ (-1 << datacenterIdBits)
var maxSequence int64 = -1 ^ (-1 << sequenceBits)

//bit operation shift
var workerIdShift = sequenceBits
var datacenterShift = workerIdBits + sequenceBits
var timestampShift = datacenterIdBits + workerIdBits + sequenceBits

var timeSequ map[int64]int64

type Snowflake struct {
	datacenterId int64
	workerId     int64
	epoch        int64
	mt           *sync.Mutex
}

var IdGenClient *Snowflake

func initIdGen() error {
	var datacenterId int64 = 0
	var workerId int64 = 0
	var epoch int64 = 1596850974657
	var err error
	IdGenClient, err = newSnowflake(datacenterId, workerId, epoch)
	if err != nil {
		log.Error("init IdGenClient error %s", err.Error())
	}
	return nil
}
func newSnowflake(datacenterId int64, workerId int64, epoch int64) (*Snowflake, error) {
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenterId cant be greater than %d or less than 0", maxDatacenterId))
	}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("workerId cant be greater than %d or less than 0", maxWorkerId))
	}
	if epoch > getCurrentTime() {
		return nil, errors.New(fmt.Sprintf("epoch time cant be after now"))
	}
	sf := Snowflake{datacenterId, workerId, epoch, new(sync.Mutex)}
	return &sf, nil
}

func (sf *Snowflake) getUniqueId() int64 {
	sf.mt.Lock()
	defer sf.mt.Unlock()
	//get current time
	currentTime := getCurrentTime()
	//compute sequence
	if currentTime < lastTime { //occur clock back
		//panic or wait,wait is not the best way.can be optimized.
		currentTime = waitUntilNextTime(lastTime)
		sequence = 0
	} else if currentTime == lastTime { //at the same time(micro-second)
		sequence = (sequence + 1) & maxSequence
		if sequence == 0 { //overflow max num,wait next time
			currentTime = waitUntilNextTime(lastTime)
		}
	} else if currentTime > lastTime { //next time
		sequence = 0
		lastTime = currentTime
	}
	//generate id
	return (currentTime-sf.epoch)<<timestampShift | sf.datacenterId<<datacenterShift |
		sf.workerId<<workerIdShift | sequence
}

func waitUntilNextTime(lasttime int64) int64 {
	currentTime := getCurrentTime()
	for currentTime <= lasttime {
		time.Sleep(1 * time.Second / 1000) //sleep micro second
		currentTime = getCurrentTime()
	}
	return currentTime
}

func getCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6 //micro second
}

func GetId() int64 {
	if IdGenClient == nil {
		panic("the IdGenClient is nil")
	}
	return IdGenClient.getUniqueId()
}
