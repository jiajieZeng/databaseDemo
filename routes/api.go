package routes

import (
	"databaseDemo/app/controller"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	router.POST("/register", controller.Register)

	router.POST("/login", controller.Login)

	router.POST("/queryfirst", controller.QueryFirst)

	router.POST("/raw", controller.RawSQL)

	router.POST("/checkin", controller.CheckIn)

	router.POST("/hash", controller.HashData)

	router.POST("/Zset", controller.Zset)

	router.POST("/TxBegin", controller.TxBegin)

	router.POST("/TxCommit", controller.TxCommit)

	router.POST("/TxRaw", controller.TxRaw)

	router.POST("/TxRollback", controller.TxRollBack)

	router.POST("/person-by-id", controller.GetPersonByID)

	router.POST("/belongings-by-id", controller.GetBelongingsByID)

	router.POST("/infos-by-id", controller.GetInfosByID)

	router.POST("/all-by-id", controller.GetAllByID)

	router.POST("/all-same-size", controller.GetAllWearSameSize)

	router.POST("/all-same-business", controller.GetAllRunSameBusiness)

	router.POST("/limitter", controller.Ex02)
}
