/**
 * @Author: xingqiyi
 * @Description:导出入口
 * @Date: 2021/11/12 9:55 上午
 */

package main

import (
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/export"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"

	"go.uber.org/zap"
)

func main() {
	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "export"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()

	module.Base()
	exportProcessFunc()

	//os.Exit(0)
	//routerUpload := initialize.RoutersExport()
	//core.RunServer(routerUpload, "keyboard")
}

func exportProcessFunc() {
	err := export.Process(global.GConfig.System.ProCode)
	if err != nil {
		global.GLog.Error("导出错误", zap.Error(err))
	}

	global.GLog.Info("导出一轮结束")
	timeAfterTrigger := time.After(30 * time.Second)
	<-timeAfterTrigger
	exportProcessFunc()
}
