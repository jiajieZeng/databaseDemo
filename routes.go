package main

import (
	"databaseDemo/controller"

	"github.com/gin-gonic/gin"
)

func CollectRoutes(r *gin.Engine) *gin.Engine {

	// 注册
	r.POST("/register", controller.Register)
	// 登录
	r.POST("/login", controller.Login)

	// db.First
	r.GET("/queryfirst", controller.QueryFirst)

	// db.Raw
	r.GET("/raw", controller.RawSQL)
	return r

}
