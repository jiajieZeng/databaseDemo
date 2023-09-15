package main

import (
	"databaseDemo/app/common"
	"databaseDemo/bootstrap"
	"databaseDemo/global"
	"github.com/gin-gonic/gin"
)

func main() {

	//获取初始化的数据库
	db := common.InitDB()
	//延迟关闭数据库
	defer db.Close()

	bootstrap.InitializeConfig()
	//创建一个默认的路由引擎
	r := gin.Default()

	//启动路由
	CollectRoutes(r)

	//在9090端口启动服务
	//panic(r.Run(":9090"))
	panic(r.Run(":" + global.App.Config.App.Port))
}
