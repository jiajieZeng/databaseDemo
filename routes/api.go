package routes

import (
	"databaseDemo/app/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SetApiGroupRoutes 定义 api 分组路由
func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/test", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "success")
	})

	// 注册
	router.POST("/register", controller.Register)
	// 登录
	router.POST("/login", controller.Login)

	// db.First
	router.GET("/queryfirst", controller.QueryFirst)

	// db.Raw
	router.GET("/raw", controller.RawSQL)

}
