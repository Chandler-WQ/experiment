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
