package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/monitor"
	"server/module/monitor/service"
	"server/utils"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

func main() {
	fmt.Println("start monitor")
	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "monitor"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog
	module.Base()

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()

	//获取监控配置
	err, conf := service.FetchFtpMonitorConf(global.GConfig.System.ProCode)
	if err != nil {
		global.GLog.Error("monitor-"+global.GConfig.System.ProCode, zap.Error(err))
		return
	}
	err = monitor.AdapterMonitor(conf)
	if err != nil {
		global.GLog.Error("monitor-"+global.GConfig.System.ProCode, zap.Error(err))
		robot := utils.NewRobot("b72bc04ad782bdeea40828d064df3869fd49c4121ed0bbcbb1ec8295508c1b01", "SEC467698ef65ac33aa5e2f04e2c12eb9b5a55325446df8be33594927b91412ab30")
		err = robot.SendTextMessage("监控报错\n"+err.Error(), []string{}, true)
	}

	timeAfterTrigger := time.After(time.Duration(conf.Frequency) * time.Minute)
	<-timeAfterTrigger
	global.GLog.Info("一轮监控结束")
	os.Exit(0)
}
