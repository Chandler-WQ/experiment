package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Chandler-WQ/experiment/handler"
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

func PostConfig(url string, desc string,method string, handlerFunc gin.HandlerFunc) *RouterConfig {
	return &RouterConfig{
		URL:    url,
		Method: method,
		HandlerFunc: []gin.HandlerFunc{
			handlerFunc,
		},
	}
}

func GetRouter() []*RouterConfig {
	return []*RouterConfig{
		PostConfig("/ping","ping",GET,handler.Ping),
	}
}