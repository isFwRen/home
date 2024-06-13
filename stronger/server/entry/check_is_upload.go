/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/20 09:22
 */

package main

import (
	"go.uber.org/zap"
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/check_is_upload"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

func main() {

	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "check_is_upload"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()

	module.Base()

	err := check_is_upload.CheckIsUpload(global.GConfig.System.ProCode)
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
	}
	timeAfterTrigger := time.After(15 * time.Second)
	<-timeAfterTrigger
	global.GLog.Info("检查一轮结束")
	//os.Exit(0)
	//routerUpload := initialize.RoutersUpload()
	//core.RunServer(routerUpload, "keyboard")
}
