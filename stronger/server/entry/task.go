package main

import (
	"fmt"
	core2 "server/core"
	"server/global"
	"server/initialize"
	"server/module"
	practice "server/module/practice/service"
	"server/module/task"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

// @title stronger doc
// @version 1.0.1
// @description This is stronger Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	fmt.Println("Hello, World task!")
	// global.GLog.Info("eeeeeeee")
	core.InitConfig()
	core2.InitConfig()
	task.TaskProCode = global.GConfig.System.ProCode + "_task"
	global.GConfig.System.Process = "task"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	// core.

	// initialize.GormPostgreSql()
	module.Base()

	//开启定时任务
	initialize.InitTaskCron()
	defer global.GTaskCron.Stop()

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()
	task.Init()
	go task.DoCheckTask()
	go task.TaskList()
	go practice.CheckPracticeUser()
	fmt.Println("---------------------11111-------------------------")

	routerUpload := initialize.RoutersTask()
	core2.RunServer(routerUpload, "keyboard")

	// process("B0118")
	// fmt.Println("DoCheckTaskDoCheckTaskDoCheckTask:")

	// err, block := service.GetTaskBlock("B0118", "op1", []string{})
	// fmt.Println("errerrerrerrerrerr222:", err, block)
	// if err == nil {
	// 	// fmt.Println("fffffffffff:")
	// 	err, fields := service.SelectOpFieldsByBlockID("B0118", block.ID)
	// 	fmt.Println("fieldsfieldsfieldsfields:", err, fields)
	// }

	// timeAfterTrigger := time.After(10 * time.Second)
	// <-timeAfterTrigger
	// fmt.Println("程序结束")

}
