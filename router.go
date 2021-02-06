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
		UrlConfig("/api/register", "Ping", POST, handler.Ping),
		UrlConfig("/api/login", "Login", POST, handler.Login),
		UrlSessionConfig("/api/logout", "Logout", GET, handler.Logout),
		UrlSessionConfig("/api/get/session/info", "GetSessionInfo", GET, handler.GetSessionInfo),
		UrlSessionConfig("/api/teach/offer", "CreatCourse", POST, handler.CreatCourse),
		UrlSessionConfig("/api/teach/:start_time", "MGetCourseInfo", GET, handler.MGetCourseInfo),
		UrlSessionConfig("/api/class/:course_id", "GetCourseAllInfo", GET, handler.GetCourseAllInfo),
		UrlSessionConfig("/api/course/update", "UpdateStudentCourse", POST, handler.UpdateStudentCourse),
	}
}
