package main

import (
	core2 "server/core"
	"server/global"
	"server/initialize"
	"server/module"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

// @BasePath /api/
// @title stronger doc
// @version 1.0.1
// @description This is stronger Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @securityDefinitions.apikey UserID
// @in header
// @name x-user-id
// @securityDefinitions.apikey ProCode
// @in header
// @name pro-code
func main() {
	//初始化运行配置
	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "common"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog
	module.Base()

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()
	//开启定时任务
	initialize.InitCron()
	defer global.GCron.Stop()

	router := initialize.Routers()
	core2.RunServer(router, "nobug")
}
