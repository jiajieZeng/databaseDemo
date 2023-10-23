package main

import (
	"databaseDemo/bootstrap"
	"databaseDemo/global"
)

func main() {
	// 初始化配置
	bootstrap.InitializeConfig()
	// 初始化日志
	global.App.Log = bootstrap.InitializeLog()

	global.App.Log.Info("log init success!")

	// 获取初始化的数据库
	//db := common.InitDB()
	// 延迟关闭数据库
	//defer db.Close()

	global.App.DB = bootstrap.InitializeDB()
	global.App.RDB = bootstrap.InitializeDBSQL()
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()

	// 创建一个默认的路由引擎
	//r := gin.Default()

	// 启动路由
	//CollectRoutes(r)
	// 在9090端口启动服务
	//panic(r.Run(":" + global.App.Config.App.Port))

	// 启动服务器
	bootstrap.RunServer()
}
