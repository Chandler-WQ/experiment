package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

func Init(r *gin.Engine) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	db.MustInitDb()
	service.InitTimeRun()
	util.Init()
	setupRouter(r)
}

func main() {
	r := gin.Default()
	Init(r)

	err := r.Run(":8080")

	if err != nil {
		log.Errorf("gin run err %s", err.Error())
	}
}
