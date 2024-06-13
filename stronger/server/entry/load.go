package main

import (
	"fmt"
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/load"
	"server/module/load/service"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
	// "server/module/load/"
)

func main() {
	fmt.Println("Hello, World!")
	// global.GLog.Info("eeeeeeee")
	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "load"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	// core.

	// initialize.GormPostgreSql()
	module.Base()
	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()
	loadProcess(global.GConfig.System.ProCode)
	timeAfterTrigger := time.After(30 * time.Second)
	<-timeAfterTrigger
	fmt.Println("程序结束")

}

// type process interface {
// 	Scan()
// }

func loadProcess(proCode string) {
	err, projectBills := service.SelectBillsByStage(proCode, 1)
	fmt.Println("SelectBillsByStage：", err, projectBills)
	// projectBills := B0001.Scan()
	// var err error
	for ii, projectBill := range projectBills {
		fmt.Println("projectBill", ii, projectBill)
		err = load.ProLoadFunc(proCode, projectBill)
		// err, projectBill = B0001.FetchBill(projectBill)
		// if err != nil {
		// 	fmt.Println("errerrerrerr:", err)
		// }
	}

}
