/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/12 9:55 上午
 */

package main

import (
	"go.uber.org/zap"
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/upload"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

func main() {

	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "upload"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()

	module.Base()

	err := upload.AutomaticUpload(global.GConfig.System.ProCode)
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
	}
	timeAfterTrigger := time.After(30 * time.Second)
	<-timeAfterTrigger
	global.GLog.Info("导出一轮结束")
	//os.Exit(0)
	//routerUpload := initialize.RoutersUpload()
	//core.RunServer(routerUpload, "keyboard")
}
