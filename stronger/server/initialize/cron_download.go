package initialize

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
	"server/global"
)

func InitDownloadCron() {
	log.Println("开启下载定时任务")
	// 例子:
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	global.GDownloadCron = cron.New()

	//定时重启下载
	err = global.GDownloadCron.AddFunc("0 */30 * * * ?", func() {
		global.GLog.Error("正在安全重启程序")
		os.Exit(0)
	})

	//codes := []string{"B0103", "B0106", "B0110"}
	//for i, _ := range codes {
	//	code := &codes[i]
	//	//全局变量标识 标志执行完再进入下一轮
	//	var gTaskMark = 1
	//	err = global.GDownloadCron.AddFunc("@every 60s", func() {
	//		if gTaskMark == 1 {
	//			gTaskMark = 2
	//			//global.GLog.Info("", zap.Any(*code, "开始查询需要下载图片的单据"))
	//			err = schedule.DownloadImages(*code)
	//			if err != nil {
	//				global.GLog.Error("下载图片错误", zap.Error(err))
	//				return
	//			}
	//			gTaskMark = 1
	//		}
	//		//global.GLog.Info("", zap.Any(*code, "每5s轮训"))
	//	})
	//}

	global.GDownloadCron.Start()
	fmt.Println("下载定时任务开启成功")
}
