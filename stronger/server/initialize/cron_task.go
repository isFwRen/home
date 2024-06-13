package initialize

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
	"server/global"
)

func InitTaskCron() {
	log.Println("开启定时任务")
	// 例子:
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	global.GTaskCron = cron.New()

	_, err := cron.Parse(global.GProConf[global.GConfig.System.ProCode].RestartAt)
	if err != nil {
		global.GLog.Error("正在安全重启程序")
		return
	}
	//定时重启录入
	err = global.GTaskCron.AddFunc(global.GProConf[global.GConfig.System.ProCode].RestartAt, func() {
		global.GLog.Error("正在安全重启程序")
		os.Exit(0)
	})

	global.GTaskCron.Start()
	fmt.Println("定时任务开启成功")
}
