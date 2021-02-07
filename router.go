package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/handler"
	"github.com/Chandler-WQ/experiment/middleware"
)

var (
	GET  = "GET"
	POST = "POST"
)

type RouterConfig struct {
	URL         string
	Method      string // 支持的方法,例如GET,POST
	HandlerFunc []gin.HandlerFunc
}

func setupRouter(r *gin.Engine) {
	for _, config := range GetRouter() {
		r.Handle(config.Method, config.URL, config.HandlerFunc...)
	}
}

func UrlConfig(url string, desc string, method string, handlerFunc gin.HandlerFunc) *RouterConfig {
	return &RouterConfig{
		URL:    url,
		Method: method,
		HandlerFunc: []gin.HandlerFunc{
			handlerFunc,
		},
	}
}

func UrlSessionConfig(url string, desc string, method string, handlerFunc gin.HandlerFunc) *RouterConfig {
	return &RouterConfig{
		URL:    url,
		Method: method,
		HandlerFunc: []gin.HandlerFunc{
			middleware.Session(),
			handlerFunc,
		},
	}
}

func GetRouter() []*RouterConfig {
	return []*RouterConfig{
		UrlConfig("/ping", "Ping", GET, handler.Ping),

		UrlConfig("/api/register", "Register", POST, handler.Register),
		UrlConfig("/api/login", "Login", POST, handler.Login),
		UrlSessionConfig("/api/logout", "Logout", GET, handler.Logout),
		UrlSessionConfig("/api/get/session/info", "GetSessionInfo", GET, handler.GetSessionInfo),

		UrlSessionConfig("/api/teach/offer", "CreatCourse", POST, handler.CreatCourse),
		UrlSessionConfig("/api/teach/:start_time/:end_time", "MGetCourseInfo", GET, handler.MGetCourseInfo),
		UrlSessionConfig("/api/class/:course_id", "GetCourseAllInfo", GET, handler.GetCourseAllInfo),
		UrlSessionConfig("/api/course/update", "UpdateStudentCourse", POST, handler.UpdateStudentCourse),

		UrlSessionConfig("/api/equip/create", "CreateEquip", POST, handler.CreateEquip),
		UrlSessionConfig("/api/equip/update", "UpdateEquip", POST, handler.UpdateEquip),
		UrlSessionConfig("/api/equip/mget/:offset/:limit", "MGetEquips", GET, handler.MGetEquips),
		UrlSessionConfig("/api/equip/get/:id", "GetEquip", GET, handler.GetEquip),
		UrlSessionConfig("/api/equip/delete", "DeleteEquips", POST, handler.DeleteEquips),

		UrlSessionConfig("/api/experiment/create", "CreateExperiment", POST, handler.CreateExperiment),
		UrlSessionConfig("/api/experiment/update", "UpdateExperiment", POST, handler.UpdateExperiment),
		UrlSessionConfig("/api/experiment/mget/:offset/:limit", "MGetExperiments", GET, handler.MGetExperiments),
		UrlSessionConfig("/api/experiment/get/:id", "GetExperiment", GET, handler.GetExperiment),
		UrlSessionConfig("/api/experiment/delete", "DeleteExperiments", POST, handler.DeleteExperiments),

		UrlSessionConfig("/api/order/create", "CreateOrder", POST, handler.CreateOrder),
		UrlSessionConfig("/api/order/update", "UpdateOrder", POST, handler.UpdateOrder),
		UrlSessionConfig("/api/order/mget/:start_time/:end_time", "MGetOrder", GET, handler.MGetOrder),
		UrlSessionConfig("/api/order/delete", "DeleteOrder", POST, handler.DeleteOrder),
	}
}
