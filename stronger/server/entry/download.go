package main

import (
	"fmt"
	core2 "server/core"
	"server/global"
	"server/initialize"
	"server/module"
	"server/module/download/project/B0102"
	B0108 "server/module/download/project/B0108"
	"server/module/download/project/B0113"
	"server/module/download/project/B0114"
	"server/module/download/project/B0116"
	B0118 "server/module/download/project/B0118"
	"server/module/download/project/B0121"
	"server/module/download/project/B0122"
	"server/module/download/project/guoshou"
	"server/utils"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

//go:generate swag init --parseDependency --parseInternal -o ../docs

// @title stronger doc
// @version 1.0.1
// @description This is stronger Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @BasePath /
func main() {
	/* 这是我的第一个简单的程序 */
	fmt.Println("start download")
	// global.GLog.Info("eeeeeeee")
	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "download"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	// initialize.GormPostgreSql()
	module.Base()
	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()
	//开启定时任务
	initialize.InitDownloadCron()
	defer global.GDownloadCron.Stop()
	if utils.RegIsMatch("^(B0106-B0103-B0110)$", global.GConfig.System.ProCode) {
		router := initialize.RoutersDownload()
		core2.RunServer(router, "nobug")
	} else {
		downloadProcess(global.GConfig.System.ProCode)
		timeAfterTrigger := time.After(30 * time.Second)
		<-timeAfterTrigger
	}
	fmt.Println("程序结束")
	// os.Exit(0)
}

func downloadProcess(proCode string) {
	switch proCode {
	case "B0118":
		B0118.Process()
	case "B0108":
		B0108.Process()
	case "B0114":
		B0114.Process()
	case "B0113":
		B0113.Process()
	case "B0121":
		B0121.Process()
	case "B0122":
		B0122.Process()
	case "B0116":
		B0116.Process()
	case "B0102":
		B0102.Process()
	case "B0103":
		guoshou.DownloadImages("B0103")
	case "B0106":
		guoshou.DownloadImages("B0106")
	case "B0110":
		guoshou.DownloadImages("B0110")
	default:
		fmt.Println("不存在的下载进程项目编号:", proCode)
	}

}
