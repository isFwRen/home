package initialize

import (
	"errors"
	"fmt"
	"log"
	"os"
	"server/global"
	"server/module/msg_manager/model"
	"server/module/msg_manager/service"
	"server/module/schedule"
	"time"

	"github.com/robfig/cron"
	"go.uber.org/zap"
)

func InitCron() {
	log.Println("开启定时任务")
	// 例子:
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	global.GCron = cron.New()

	//定时重启管理
	err = global.GCron.AddFunc("0 0 0 * * ?", func() {
		global.GLog.Error("正在安全重启程序")
		os.Exit(0)
	})

	//测试
	err = global.GCron.AddFunc("0 */5 * * *", schedule.BusinessAnalyze)
	err = global.GCron.AddFunc("0 */10 * * *", schedule.BusinessUploadAnalyzeInTheMorning)
	err = global.GCron.AddFunc("0 */15 * * *", schedule.WrongAnalysis)
	err = global.GCron.AddFunc("0 0 1 * * ?", schedule.WrongAnalysis)
	err = global.GCron.AddFunc("0 0 1 * * ?", func() {
		go schedule.BackupsAndDel()
	})
	// go schedule.BackupsAndDel()
	if err != nil {
		global.GLog.Error("WrongAnalysis时间有误")
		//return
	}
	//来量/回传统计
	//err = global.GCron.AddFunc("0 30 8,9,10,11,12,13,14,15,16,17,18 * * *", BusinessAnalyze)
	//回传分析
	//err = global.GCron.AddFunc("0 0 0 * * *", BusinessUploadAnalyzeInTheMorning)
	//错误分析
	//err = global.GCron.AddFunc("0 0 0 * * *", WrongAnalysis)

	//定时重启录入
	//for _, project := range global.GProConf {
	//	code := project.Code
	//	restartAt := project.RestartAt
	//	port := project.BackEndPort
	//	s, err := time.Parse("15:04:05", restartAt)
	//	if err != nil {
	//		global.GLog.Error("项目" + code + "重启时间有误:::" + restartAt)
	//		continue
	//	}
	//	if port == 0 {
	//		continue
	//	}
	//	global.GLog.Info(code + "	" + strconv.Itoa(s.Second()) + " " + strconv.Itoa(s.Minute()) + " " + strconv.Itoa(s.Hour()) + " * * ? " + strconv.Itoa(port))
	//	err = global.GCron.AddFunc(strconv.Itoa(s.Second())+" "+strconv.Itoa(s.Minute())+" "+strconv.Itoa(s.Hour())+" * * ?", func() {
	//		schedule.RestartTaskByCode(code, strconv.Itoa(port))
	//	})
	//	if err != nil {
	//		global.GLog.Error("项目" + code + "定时有误:::" + restartAt)
	//		continue
	//	}
	//}

	//定时通知
	err = GroupNotice()
	if err != nil {
		global.GLog.Error("项目定时通知错误:::" + err.Error())
		//return
	}

	//定时计算单据是否超时
	err = global.GCron.AddFunc("0 */1 * * *", schedule.CalculateTimeout)
	if err != nil {
		global.GLog.Error("定时计算单据是否超时错误:::" + err.Error())
		//return
	}

	//定时计算日报
	err = global.GCron.AddFunc("0 */1 * * * ?", schedule.CalculateDayReport)
	//err = global.GCron.AddFunc("0 * 14 * * ?", schedule.CalculateDayReport)
	if err != nil {
		global.GLog.Error("定时计算日报错误:::" + err.Error())
		//return
	}

	//定时计算字符数
	err = global.GCron.AddFunc("0 1 * * * ?", schedule.CalculateDayCharSum)
	//err = global.GCron.AddFunc("0 * 10 * * ?", schedule.CalculateDayCharSum)
	if err != nil {
		global.GLog.Error("定时计算日报错误:::" + err.Error())
		return
	}

	//err = global.GCron.AddFunc("0 0 1 * * ?", schedule.UpdateHospitalAndCatalogue)
	err = global.GCron.AddFunc("0 */5 * * * ?", schedule.UpdateHospitalAndCatalogue)
	if err != nil {
		global.GLog.Error("定时更新国寿常量错误:::" + err.Error())
		//return
	}

	//CSB0108RC0006000 B0108
	//当时效表中出现剩余时间为负数的案件时，将状态为“录入中”的案件单号放到紧急单池中
	//2023年10月11日11:57:34 修改需求
	//秒赔案件：
	//1.当机构号倒数第二位为“0”，且案件剩余时间≤10分钟时，将状态为“录入中”的案件单号放到优先单池中；
	//2.当机构号倒数第二位为“1”，且案件剩余时间为负数时，将状态为“录入中”的案件单号放到优先单池中；
	//理赔案件：
	//1.当机构号倒数第二位为“0”，且案件剩余时间为负数时，将状态为“录入中”的案件单号放到优先单池中；
	//2.当机构号倒数第二位为“1”，且案件剩余时间≤-10分钟时，将状态为“录入中”的案件单号放到优先单池中；
	err = global.GCron.AddFunc("0 */1 * * * ?", schedule.UpdateTimeoutStage)
	if err != nil {
		global.GLog.Error("录入中超时单据放优先单池:::" + err.Error())
		//return
	}
	//时效简报
	err = global.GCron.AddFunc("*/5 * * * * ?", schedule.GetTimeLinessBriefingCron)
	if err != nil {
		global.GLog.Error("获取简报错误:::" + err.Error())
		//return
	}

	global.GCron.Start()
	fmt.Println("定时任务开启成功")
}

func GroupNotice() error {
	//定时通知
	global.GLog.Warn("定时通知", zap.Any("星期", int(time.Now().Weekday())))
	err, noticeTwos, noticeOnes := service.GetGroupNotice()
	if err != nil {
		return err
	}

	//构造所有发送时间点
	for _, one := range noticeOnes {
		startTime, err := time.Parse("15:04:05", one.StartTime)
		if err != nil {
			global.GLog.Error("项目" + one.ProCode + "定时开始时间有误:::" + one.StartTime)
			return err
		}
		endTime, err := time.Parse("15:04:05", one.EndTime)
		if err != nil {
			global.GLog.Error("项目" + one.ProCode + "定时结束时间有误:::" + one.EndTime)
			return err
		}
		if startTime.After(endTime) {
			return errors.New("开始时间比结束时间大,id:::" + one.ID + ",ProCode:::" + one.ProCode)
		}
		newTime := startTime
		for {
			if newTime.After(endTime) {
				break
			}
			var newOne = model.GroupNoticeTwo{
				ProCode:   one.ProCode,
				SendTime:  newTime.Format("15:04:05"),
				DayOfWeek: one.DayOfWeek,
				GroupId:   one.GroupId,
			}
			noticeTwos = append(noticeTwos, newOne)
			newTime = newTime.Add(time.Duration(one.Interval) * time.Minute)
		}
	}

	//新建定时器
	for _, notice := range noticeTwos {
		sendTime, err := time.Parse("15:04:05", notice.SendTime)
		if err != nil {
			global.GLog.Error("项目" + notice.ProCode + "定时时间有误:::" + notice.SendTime)
			return err
		}
		spec := fmt.Sprintf("%d %d %d * * %d", sendTime.Second(), sendTime.Minute(), sendTime.Hour(), notice.DayOfWeek)
		global.GLog.Info(spec + "   " + notice.ProCode)
		err = global.GCron.AddFunc(spec, func() {
			schedule.SendMsg(notice)
		})
		if err != nil {
			global.GLog.Error("项目" + notice.ProCode + "定时有误:::" + notice.SendTime)
			return err
		}
	}
	return err
}
