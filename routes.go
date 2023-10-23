package main

import (
	"databaseDemo/app/controller"
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
	r.POST("/raw", controller.RawSQL)
	return r

}
