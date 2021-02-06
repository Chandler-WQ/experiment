package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/service"
	"github.com/Chandler-WQ/experiment/util"
)

//// 通过字典模拟 DB
//var db = make(map[string]string)
//
//func setupRouter1(r *gin.Engine) {
//	// 初始化 Gin 框架默认实例，该实例包含了路由、中间件以及配置信息
//
//
//	// Ping 测试路由
//	r.GET("/ping", func(c *gin.Context) {
//		c.String(http.StatusOK, "pong")
//	})
//
//	// 获取用户数据路由
//	r.GET("/user/:name", func(c *gin.Context) {
//		user := c.Params.ByName("name")
//		value, ok := db[user]
//		if ok {
//			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
//		}
//	})
//
//	// 需要 HTTP 基本授权认证的子路由群组设置
//	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
//		"foo":  "bar", // 用户名:foo 密码:bar
//		"manu": "123", // 用户名:manu 密码:123
//	}))
//
//	// 保存用户信息路由
//	authorized.POST("admin", func(c *gin.Context) {
//		user := c.MustGet(gin.AuthUserKey).(string)
//
//		// 解析并验证 JSON 格式请求数据
//		var json struct {
//			Value string `json:"value" binding:"required"`
//		}
//
//		if c.Bind(&json) == nil {
//			db[user] = json.Value
//			c.JSON(http.StatusOK, gin.H{"status": "ok"})
//		}
//	})
//}

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
		log.Error("gin run err %s", err.Error())
	}
}
