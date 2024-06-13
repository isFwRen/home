package schedule

import (
	"database/sql"
	"errors"
	"fmt"
	"server/global"
	model4 "server/module/homepage/model"
	l "server/module/load/model"
	"server/module/pro_manager/model"
	"server/module/pro_manager/project/B0118"
	model3 "server/module/report_management/model"
	"server/module/report_management/project"
	model2 "server/module/schedule/model"
	"server/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PeriodOfTimeInit() map[string]string {
	PeriodOfTime := make(map[string]string, 0)
	PeriodOfTime["a_half_past_eight"] = "00:00:00-08:30:00"
	PeriodOfTime["a_half_past_nine"] = "08:30:01-09:30:00"
	PeriodOfTime["a_half_past_ten"] = "09:30:01-10:30:00"
	PeriodOfTime["a_half_past_eleven"] = "10:30:01-11:30:00"
	PeriodOfTime["a_half_past_twelve"] = "11:30:01-12:30:00"
	PeriodOfTime["a_half_past_thirteen"] = "12:30:01-13:30:00"
	PeriodOfTime["a_half_past_fourteen"] = "13:30:01-14:30:00"
	PeriodOfTime["a_half_past_fifteen"] = "14:30:01-15:30:00"
	PeriodOfTime["a_half_past_sixteen"] = "15:30:01-16:30:00"
	PeriodOfTime["a_half_past_seventeen"] = "16:30:01-17:30:00"
	PeriodOfTime["a_half_past_eighteen"] = "17:30:01-18:30:00"
	PeriodOfTime["before_zero_hour"] = "18:30:01-23:59:59"
	return PeriodOfTime
}

func calculateHistoryDbConnection() (num int) {
	for code, _ := range global.ProDbMap {
		if strings.Index(code, "_task") == -1 {
			num++
		}
	}
	return num
}

func BusinessAnalyze() {
	PeriodOfTime := PeriodOfTimeInit()
	fmt.Println("start cron, 每隔30分钟统计一下业务量-下载/回传")
	fmt.Println("global.ProDbMap", global.ProDbMap)
	wg := sync.WaitGroup{}
	for dbCode, v := range global.ProDbMap {
		if strings.Index(dbCode, "_task") == -1 {
			wg.Add(1)
			go func(v *gorm.DB) {
				defer wg.Done()
				//第一次判断有没有今天的数据
				var t int64
				err := v.Model(&model2.BusinessAnalysis{}).Where("create_at = ? AND types = 'download'", time.Now().Format("2006-01-02")).Count(&t).Error
				if err != nil {
					global.GLog.Error("BusinessAnalysis function", zap.Any("错误❎", err.Error()))
					return
				}
				if t == 0 {
					NewDay := &model2.BusinessAnalysis{CreateAt: time.Now().Format("2006-01-02"), AHalfPastEight: 0, AHalfPastNine: 0, AHalfPastTen: 0, AHalfPastEleven: 0, AHalfPastTwelve: 0, AHalfPastThirteen: 0, AHalfPastFourteen: 0, AHalfPastFifteen: 0, AHalfPastSixteen: 0, AHalfPastSeventeen: 0, AHalfPastEighteen: 0, BeforeZeroHour: 0, Types: "download"}
					err = v.Model(&model2.BusinessAnalysis{}).Create(&NewDay).Error
					if err != nil {
						global.GLog.Error("BusinessAnalysisInTheMorning function", zap.Any("错误❎", err.Error()))
						return
					}
				}

				err = v.Model(&model2.BusinessAnalysis{}).Where("create_at = ? AND types = 'upload'", time.Now().Format("2006-01-02")).Count(&t).Error
				if err != nil {
					global.GLog.Error("BusinessAnalysis function", zap.Any("错误❎", err.Error()))
					return
				}
				if t == 0 {
					NewDay := &model2.BusinessAnalysis{CreateAt: time.Now().Format("2006-01-02"), AHalfPastEight: 0, AHalfPastNine: 0, AHalfPastTen: 0, AHalfPastEleven: 0, AHalfPastTwelve: 0, AHalfPastThirteen: 0, AHalfPastFourteen: 0, AHalfPastFifteen: 0, AHalfPastSixteen: 0, AHalfPastSeventeen: 0, AHalfPastEighteen: 0, BeforeZeroHour: 0, Types: "upload"}
					err = v.Model(&model2.BusinessAnalysis{}).Create(&NewDay).Error
					if err != nil {
						global.GLog.Error("BusinessAnalysisInTheMorning function", zap.Any("错误❎", err.Error()))
						return
					}
				}

				//开始统计
				for i, k := range PeriodOfTime {
					var total int64
					arr := strings.Split(k, "-")
					StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" "+arr[0], time.Local)
					EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" "+arr[1], time.Local)
					//下载
					err = v.Model(&model.ProjectBill{}).Where("scan_at >= ? AND scan_at <= ? ", StartTime, EndTime).Count(&total).Error
					if err != nil {
						global.GLog.Error("BusinessAnalysis function", zap.Any("错误❎", err.Error()))
						return
					}
					err = v.Model(&model2.BusinessAnalysis{}).Where("create_at = ? AND types = 'download'", time.Now().Format("2006-01-02")).Update(i, total).Error
					if err != nil {
						global.GLog.Error("BusinessAnalysis function", zap.Any("错误❎", err.Error()))
						return
					}
					//回传
					var bi []model.ProjectBill
					err = v.Model(&model.ProjectBill{}).Where("upload_at >= ? AND upload_at <= ? ", StartTime, EndTime).Find(&bi).Count(&total).Error
					if err != nil {
						global.GLog.Error("BusinessAnalysis function", zap.Any("错误❎", err.Error()))
						return
					}
					err = v.Model(&model2.BusinessAnalysis{}).Where("create_at = ? AND types = 'upload'", time.Now().Format("2006-01-02")).Update(i, total).Error
					if err != nil {
						global.GLog.Error("BusinessAnalysis function", zap.Any("错误❎", err.Error()))
						return
					}
				}
			}(v)
		}
	}
	wg.Wait()
}

// BusinessUploadAnalyzeInTheMorning 凌晨统计回传量
func BusinessUploadAnalyzeInTheMorning() {
	wg := sync.WaitGroup{}
	wg.Add(calculateHistoryDbConnection())
	for dbCode, v := range global.ProDbMap {
		if strings.Index(dbCode, "_task") == -1 {
			fmt.Println("dbCode", dbCode)
			go func(v *gorm.DB) {
				var total int64
				var bills []model.ProjectBill
				var BusinessUpload model2.BusinessUploadAnalysis
				BeforeDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
				NowDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 59, time.Now().Location())
				err := v.Model(&model.ProjectBill{}).Where("scan_at >= ? AND scan_at <= ? AND (stage = '5' OR  stage = '7')", BeforeDay, NowDay).Count(&total).Error
				if err != nil {
					global.GLog.Error("BusinessUploadAnalyzeInTheMorning function", zap.Any("错误❎", err.Error()))
					return
				}
				err = v.Order("id desc").Where("scan_at >= ? AND scan_at <= ? AND (stage = '5' OR stage = '7')", BeforeDay, NowDay).Find(&bills).Error
				if err != nil {
					global.GLog.Error("BusinessUploadAnalyzeInTheMorning function", zap.Any("错误❎", err.Error()))
					return
				}
				totalTime := 0.0
				for _, bill := range bills {
					BusinessUpload.VolumeOfBusiness += 1
					diff := bill.UploadAt.Sub(bill.ScanAt).Hours()
					totalTime += bill.UploadAt.Sub(bill.ScanAt).Seconds()
					if diff >= 0 && diff <= 1 {
						BusinessUpload.OneHours += 1
					} else if diff > 1 && diff <= 2 {
						BusinessUpload.TwoHours += 1
					} else if diff > 2 && diff <= 3 {
						BusinessUpload.ThreeHours += 1
					} else {
						BusinessUpload.MoreThanThreeHours += 1
					}
				}

				//算比例
				//NaN inf
				if BusinessUpload.VolumeOfBusiness == 0 {
					BusinessUpload.OneHoursRate = 0
					BusinessUpload.TwoHoursRate = 0
					BusinessUpload.ThreeHoursRate = 0
					BusinessUpload.MoreThanThreeHoursRate = 0
					BusinessUpload.TheAverageTime = 0
				} else {
					BusinessUpload.OneHoursRate = float64(BusinessUpload.OneHours) / float64(BusinessUpload.VolumeOfBusiness)
					BusinessUpload.TwoHoursRate = float64(BusinessUpload.TwoHours) / float64(BusinessUpload.VolumeOfBusiness)
					BusinessUpload.ThreeHoursRate = float64(BusinessUpload.ThreeHours) / float64(BusinessUpload.VolumeOfBusiness)
					BusinessUpload.MoreThanThreeHoursRate = float64(BusinessUpload.MoreThanThreeHours) / float64(BusinessUpload.VolumeOfBusiness)
					BusinessUpload.TheAverageTime = Decimal(totalTime / 3600 / float64(BusinessUpload.VolumeOfBusiness))
				}
				//保存进数据库
				var t int64
				err = v.Model(&model2.BusinessUploadAnalysis{}).Where("create_at = ? ", time.Now().Format("2006-01-02")).Count(&t).Error
				if err != nil {
					global.GLog.Error("BusinessUploadAnalyzeInTheMorning function", zap.Any("错误❎", err.Error()))
					return
				}
				if t == 0 {
					BusinessUpload.CreateAt = time.Now().Format("2006-01-02")
					BusinessUpload.TheAverageTime = 0
					err = v.Model(&model2.BusinessUploadAnalysis{}).Create(&BusinessUpload).Error
					if err != nil {
						global.GLog.Error("BusinessUploadAnalyzeInTheMorning function", zap.Any("错误❎", err.Error()))
						return
					}
				} else {
					err = v.Model(&model2.BusinessUploadAnalysis{}).Where("create_at = ? ", time.Now().Format("2006-01-02")).Updates(map[string]interface{}{
						"volume_of_business":         BusinessUpload.VolumeOfBusiness,
						"the_average_time":           BusinessUpload.TheAverageTime,
						"one_hours":                  BusinessUpload.OneHours,
						"two_hours":                  BusinessUpload.TwoHours,
						"three_hours":                BusinessUpload.ThreeHours,
						"more_than_three_hours":      BusinessUpload.MoreThanThreeHours,
						"one_hours_rate":             BusinessUpload.OneHoursRate,
						"two_hours_rate":             BusinessUpload.TwoHoursRate,
						"three_hours_rate":           BusinessUpload.ThreeHoursRate,
						"more_than_three_hours_rate": BusinessUpload.MoreThanThreeHoursRate,
					}).Error
					if err != nil {
						global.GLog.Error("BusinessUploadAnalyzeInTheMorning function", zap.Any("错误❎", err.Error()))
						return
					}
				}

				wg.Done()
			}(v)
		}
	}
	wg.Wait()
}

// WrongAnalysis 错误分析
func WrongAnalysis() {
	wg := sync.WaitGroup{}
	wg.Add(calculateHistoryDbConnection())
	for dbCode, v := range global.ProDbMap {
		if strings.Index(dbCode, "_task") == -1 {
			go func(v *gorm.DB) {
				//dbCode 会给污染
				//fmt.Println("dbCode", dbCode)

				//k: code|submitDay v: model3.WrongAnalysis
				wrongAnalysisMap := make(map[string]*model3.WrongAnalysis, 0)

				threeDaysAgo := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-3, 0, 0, 0, 0, time.Now().Location())
				nowDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 0, time.Now().Location())
				var total int64
				err := v.Model(&model3.Wrong{}).Where("submit_day >= ? AND submit_day <= ? AND code <> '0'", threeDaysAgo, nowDay).Count(&total).Error
				if err != nil {
					global.GLog.Error("WrongAnalysis function", zap.Any("错误❎", dbCode+":"+err.Error()))
					return
				}

				var wrongs []model3.Wrong
				err = v.Model(&model3.Wrong{}).Order("id desc").Where("submit_day >= ? AND submit_day <= ? AND code <> '0'", threeDaysAgo, nowDay).Find(&wrongs).Error
				if err != nil {
					global.GLog.Error("WrongAnalysis function", zap.Any("错误❎", dbCode+":"+err.Error()))
					return
				}

				for _, wrong := range wrongs {
					if _, ok := wrongAnalysisMap[wrong.Code+"|"+wrong.SubmitDay.Format("2006-01-02")]; !ok {
						wrongAnalysisMap[wrong.Code+"|"+wrong.SubmitDay.Format("2006-01-02")] = &model3.WrongAnalysis{}
						EditWrongAnalysis(wrongAnalysisMap, wrong, wrong.Code+"|"+wrong.SubmitDay.Format("2006-01-02"))
					} else {
						EditWrongAnalysis(wrongAnalysisMap, wrong, wrong.Code+"|"+wrong.SubmitDay.Format("2006-01-02"))
					}
				}

				for i, m := range wrongAnalysisMap {
					arr := strings.Split(i, "|")
					var t int64
					err = v.Model(&model3.WrongAnalysis{}).Where("to_char(statistical_time,'YYYY-MM-DD') = ? AND code = ? ", arr[1], arr[0]).Count(&t).Error
					if err != nil {
						global.GLog.Error("WrongAnalysis function", zap.Any("错误❎", dbCode+":"+err.Error()))
						return
					}
					if t == 0 {
						err = v.Model(&model3.WrongAnalysis{}).Create(&m).Error
						if err != nil {
							global.GLog.Error("WrongAnalysis function", zap.Any("错误❎", dbCode+":"+err.Error()))
							return
						}
					} else {
						err = v.Model(&model3.WrongAnalysis{}).Where("to_char(statistical_time,'YYYY-MM-DD') = ? AND code = ? ", arr[1], arr[0]).Updates(map[string]interface{}{
							"wrong_number":             m.WrongNumber,
							"the_number_of_complaints": m.TheNumberOfComplaints,
							"the_complaint_rate":       m.TheComplaintRate,
							"through_the_number":       m.ThroughTheNumber,
							"the_pass_rate":            m.ThePassRate,
							"non_passing_quantity":     m.NonPassingQuantity,
							"unqualified_rate":         m.UnqualifiedRate,
						}).Error
						if err != nil {
							global.GLog.Error("WrongAnalysis function", zap.Any("错误❎", dbCode+":"+err.Error()))
							return
						}
					}
				}
				fmt.Println("WrongAnalysis finish!")
			}(v)

		}
	}
	wg.Wait()
}

func EditWrongAnalysis(M map[string]*model3.WrongAnalysis, w model3.Wrong, key string) {
	//申诉数量：录入错误人员已申诉的错误数
	//申诉率：申诉率=申诉数量/错误数量*100%
	//通过数量：录入错误人员已申诉且审核通过的错误数；
	//通过率：通过率=通过数量/申诉数量*100%
	//不通过数量：录入错误人员已申诉且审核不通过的错误数
	//不通过率：不通过率=不通过数量/申诉数量*100%
	M[key].StatisticalTime = w.SubmitDay
	M[key].Code = w.Code
	M[key].NickName = w.NickName
	M[key].WrongNumber += 1
	if w.IsComplain {
		M[key].TheNumberOfComplaints += 1
		if w.IsWrongConfirm {
			M[key].ThroughTheNumber += 1
		}
	} else {
		if w.IsAudit && !w.IsWrongConfirm {
			M[key].TheNumberOfComplaints += 1
			M[key].NonPassingQuantity += 1
		}
	}
	//防止分母为0
	if M[key].TheNumberOfComplaints == 0 {
		M[key].ThePassRate = 0
		M[key].UnqualifiedRate = 0
	} else {
		M[key].ThePassRate = M[key].ThroughTheNumber / M[key].TheNumberOfComplaints
		M[key].UnqualifiedRate = M[key].NonPassingQuantity / M[key].TheNumberOfComplaints
	}
	M[key].TheComplaintRate = M[key].TheNumberOfComplaints / M[key].WrongNumber
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// CalculateTimeout 定时计算单据是否超时
func CalculateTimeout() {
	global.GLog.Info("定时计算单据是否超时", zap.Any("时间", time.Now()))
	for _, db := range global.ProDbMap {
		var bills []model.ProjectBill
		db.Model(&model.ProjectBill{}).
			Where("(stage = 5 or stage = 7) and (to_char(deadline_upload_time, 'YYYY-MM-DD') = '0001-01-01' or deadline_upload_time is null ) and scan_at BETWEEN ? AND ? ",
				time.Now().Add(-48*60*time.Minute), time.Now()).Limit(100).Find(&bills)
		global.GLog.Info("定时计算单据是否超时", zap.Any("bills", len(bills)))
		for _, bill := range bills {
			switch bill.ProCode {
			case "B0118":
				err, deadLineUploadTime, _, second := B0118.CalculateBackTimeAndTimeRemaining(bill, 1, bill.ProCode)
				if err != nil {
					global.GLog.Error("定时计算单据是否超时--锦浩的方法", zap.Error(err))
					continue
				}
				isTimeout := false
				global.GLog.Info("定时计算单据是否超时", zap.Any("second", second))
				if second < 0 {
					isTimeout = true
				}
				//parse, err := time.Parse("2006-01-02 15:04:05", deadLineUploadTime)
				//if err != nil {
				//	global.GLog.Error("定时计算单据是否超时", zap.Error(err))
				//	continue
				//}
				db.Model(&model.ProjectBill{}).Where("id = ?", bill.ID).Updates(&map[string]interface{}{
					"deadline_upload_time": deadLineUploadTime,
					"is_timeout":           isTimeout,
				})
			}
		}
	}
}

// CalculateDayReport 计算日报
func CalculateDayReport() {
	calTime := time.Now().Add(-24 * 60 * time.Minute)
	global.GLog.Info("计算日报", zap.Any("时间", calTime))
	reportDate := time.Date(calTime.Year(), calTime.Month(), calTime.Day(), 0, 0, 0, 0, time.UTC)
	reportDayStr := reportDate.Format("2006-01-02")
	day := calTime.Day()
	mDay := utils.GetYearMonthToDay(reportDate.Year(), int(reportDate.Month()))
	timePercent, _ := decimal.NewFromInt(int64(day)).DivRound(decimal.NewFromInt(int64(mDay)), 4).Float64()
	for _, pro := range global.GProConf {
		proDb := global.ProDbMap[pro.Code]
		if proDb == nil {
			global.GLog.Error("计算日报", zap.Error(global.ProDbErr))
			continue
		}
		var predictValue float64
		err := global.GDb.Model(&model4.ProReport{}).Select("predict_value").
			Where("to_char(report_date,'YYYY-MM-DD') = ? and pro_code = ?", reportDayStr, pro.Code).Scan(&predictValue).Error
		if err != nil {
			global.GLog.Error("计算日报predictValue", zap.Error(global.ProDbErr))
			continue
		}

		err, monthCount := getCount(proDb, 2, reportDate, false)
		if err != nil {
			global.GLog.Error("计算日报monthCount", zap.Error(global.ProDbErr))
			continue
		}

		err, dayCount := getCount(proDb, 1, reportDate, false)
		if err != nil {
			global.GLog.Error("计算日报dayCount", zap.Error(global.ProDbErr))
			continue
		}

		err, monthTimeoutCount := getCount(proDb, 2, reportDate, true)
		if err != nil {
			global.GLog.Error("计算日报monthCount", zap.Error(global.ProDbErr))
			continue
		}

		err, dayTimeoutCount := getCount(proDb, 1, reportDate, true)
		if err != nil {
			global.GLog.Error("计算日报dayCount", zap.Error(global.ProDbErr))
			continue
		}
		monthAgingPercent := 0.0
		if monthCount != 0 {
			monthAgingPercent, _ = decimal.NewFromInt(monthCount).Sub(decimal.NewFromInt(monthTimeoutCount)).DivRound(decimal.NewFromInt(monthCount), 4).Float64()
		}
		//var finishValue float64
		//err = proDb.Raw("SELECT coalesce(sum(count_money),0) FROM project_bills WHERE to_char(scan_at, 'YYYY-MM-DD') = ?", reportDayStr).Scan(&finishValue).Error
		//if err != nil {
		//	global.GLog.Error("计算日报finishValue", zap.Error(global.ProDbErr))
		//	continue
		//}
		finishValue := CalculateMoney(proDb, reportDate, 1)

		err, monthErrorCount := getQualitiesCount(2, reportDate, pro.Code)
		if err != nil {
			global.GLog.Error("计算日报monthErrorCount", zap.Error(global.ProDbErr))
			continue
		}

		err, dayErrorCount := getQualitiesCount(1, reportDate, pro.Code)
		if err != nil {
			global.GLog.Error("计算日报dayErrorCount", zap.Error(global.ProDbErr))
			continue
		}

		finishPercent := 0.0
		if predictValue != 0 {
			finishPercent, _ = decimal.NewFromFloat(finishValue).DivRound(decimal.NewFromFloat(predictValue), 4).Float64()
		}

		monthRightPercent := 0.0
		if monthCount != 0 {
			rightCount := decimal.NewFromInt(monthCount).Sub(decimal.NewFromInt(monthErrorCount))
			monthRightPercent, _ = rightCount.DivRound(decimal.NewFromInt(monthCount), 4).Float64()
		}
		var proReport model4.ProReport
		proReport.ProCode = pro.Code                    //项目编码
		proReport.FinishValue = finishValue             //实际完成
		proReport.TimePercent = timePercent             //时间比例
		proReport.FinishPercent = finishPercent         //完成比例
		proReport.MonthCount = monthCount               //月业务总量
		proReport.DayCount = dayCount                   //日业务量
		proReport.MonthAgingPercent = monthAgingPercent //月时效保障率
		proReport.MonthTimeoutCount = monthTimeoutCount //月超时数量
		proReport.DayTimeoutCount = dayTimeoutCount     //日超时数量
		proReport.MonthRightPercent = monthRightPercent //月质量准确率
		proReport.MonthErrorCount = monthErrorCount     //月差错数量
		proReport.DayErrorCount = dayErrorCount         //日差错数量
		proReport.ReportDate = reportDate               //报表日期

		err = global.GDb.Model(&model4.ProReport{}).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "report_date"}, {Name: "pro_code"}},
			DoUpdates: clause.AssignmentColumns([]string{"finish_value", "time_percent",
				"finish_percent", "month_count", "day_count", "month_aging_percent", "month_timeout_count",
				"day_timeout_count", "month_right_percent", "month_error_count", "day_error_count"}),
		}).Create(&proReport).Error
		if err != nil {
			global.GLog.Error("计算日报CreateOrUpdate", zap.Error(global.ProDbErr))
			continue
		}
	}
}

// getCount 获取数量
func getCount(proDb *gorm.DB, dateType int, reportDate time.Time, isCalTimeout bool) (err error, dayCount int64) {
	reportDateStr := reportDate.Format("2006-01-02")
	formatStr := "YYYY-MM-DD"
	if dateType == 2 {
		reportDateStr = reportDate.Format("2006-01")
		formatStr = "YYYY-MM"
	}
	proDb = proDb.Model(&model.ProjectBill{}).
		Where("to_char(scan_at,'"+formatStr+"') = ?  and stage = 5", reportDateStr)
	if isCalTimeout {
		proDb = proDb.Where("is_timeout = true")
	}
	err = proDb.Count(&dayCount).Error
	return err, dayCount
}

// getQualitiesCount 获取质检数量
func getQualitiesCount(dateType int, reportDate time.Time, proCode string) (err error, dayCount int64) {
	reportDateStr := reportDate.Format("2006-01-02")
	formatStr := "YYYY-MM-DD"
	if dateType == 2 {
		reportDateStr = reportDate.Format("2006-01")
		formatStr = "YYYY-MM"
	}
	err = global.GDb.Model(&model.Quality{}).Select("count(distinct(bill_name))").
		Where("to_char(entry_date,'"+formatStr+"') = ? and pro_code = ?", reportDateStr, proCode).Count(&dayCount).Error
	return err, dayCount
}

type billObj struct {
	model.ProjectBill
	ProjectField []l.ProjectField `json:"projectField" gorm:"foreignKey:BillID;references:ID;"`
}

func (v billObj) TableName() string {
	return "project_bills"
}

// CalculateDayCharSum 定时计算字符数
func CalculateDayCharSum() {
	calTime := time.Now().Add(-24 * 60 * time.Minute)
	//calTime := time.Now().Add(-35 * 24 * 60 * time.Minute)
	global.GLog.Info("计算字符数", zap.Any("时间", calTime))
	calDate := time.Date(calTime.Year(), calTime.Month(), calTime.Day(), 0, 0, 0, 0, time.UTC)
	calDayStr := calDate.Format("2006-01-02")
	for _, pro := range global.GProConf {
		proDb := global.ProDbMap[pro.Code]
		if proDb == nil {
			global.GLog.Error("计算字符数", zap.Error(global.ProDbErr))
			continue
		}
		//var bills []model.ProjectBill
		var allBills []billObj
		var countBill int64
		batchSize := 100

		//计算员工字符========================================================================================
		var staffCharSum *int64
		var staffCharSumInt64 sql.NullInt64
		tmp := int64(0)
		staffCharSum = &tmp
		//err := proDb.Table("output_statistics_summaries").
		err := global.GDb.Table("output_statistics_summaries").
			Select("SUM(mary) AS mary_sum").
			Where("to_char(submit_time,'YYYY-MM-DD') = ? and pro_code = ?", calDayStr, pro.Code).
			Scan(&staffCharSumInt64).Error
		if err != nil {
			global.GLog.Error("查询员工字符数 err ", zap.Error(global.ProDbErr))
			continue
		}
		if staffCharSumInt64.Valid {
			*staffCharSum = staffCharSumInt64.Int64
		} else {
			*staffCharSum = int64(0)
		}
		//==================================================================================================

		proDb.Model(&model.ProjectBill{}).Where("(stage = 5 or stage = 7) and to_char(scan_at,'YYYY-MM-DD') = ?", calDayStr).Count(&countBill)
		for offset := 0; offset < int(countBill); offset += batchSize {
			var bills []billObj
			err := proDb.Where("(stage = 5 or stage = 7) and to_char(scan_at,'YYYY-MM-DD') = ?", calDayStr).
				//Select("id,bill_name,bill_num,pro_code,stage,download_path,download_at,batch_num").
				Preload("ProjectField").
				Order("created_at").Limit(batchSize).Offset(offset).Find(&bills).Error
			if err != nil {
				global.GLog.Error("计算字符数Find bills", zap.Error(global.ProDbErr))
				return
			}
			allBills = append(allBills, bills...)
		}
		//countBill := int64(len(bills))
		global.GLog.Info("计算字符数", zap.Any("bills", countBill))
		if countBill < 1 {
			global.GLog.Error("计算字符数", zap.Error(errors.New(calDayStr+"没单")))
			continue
		}
		charSum := 0
		charInputSum := 0
		charPercent := 0.0
		claimTypeCharSumMap := make(map[int]*BillsScaleSum)      //理赔类型（有，无，混合）->结算字符
		claimTypeCharInputSumMap := make(map[int]*BillsScaleSum) //理赔类型（有，无，混合）->录入字符
		//const_data.BillClaimType
		for _, obj := range allBills {
			bill := obj.ProjectBill
			fields := obj.ProjectField
			//var fields []l.ProjectField
			//err = proDb.Model(&l.ProjectField{}).Order("created_at desc").Where("bill_id = ? ", bill.ID).Find(&fields).Error
			//if err != nil {
			//	global.GLog.Error(err.Error(), zap.Error(err))
			//}
			//switch bill.ProCode {
			//case "B0118":

			//charSum += project.CalculateXmlCharacter(bill) //结算字符
			//charInputSum += project.CalculateWriteCharacter(fields) //录入字符
			CSTmp := project.CalculateXmlCharacter(bill)
			CISTmp := project.CalculateWriteCharacter(fields)
			charSum += CSTmp       //结算字符
			charInputSum += CISTmp //录入字符

			//-----------分别统计 理赔类型-字符数:
			_, ok := claimTypeCharSumMap[bill.ClaimType]
			if ok {
				claimTypeCharSumMap[bill.ClaimType].UpdateValue(CSTmp)
			} else {
				b := BillsScaleSum{
					Num:     1,
					CharSum: CSTmp,
				}
				claimTypeCharSumMap[bill.ClaimType] = &b
			}
			_, ok1 := claimTypeCharInputSumMap[bill.ClaimType]
			if ok1 {
				claimTypeCharInputSumMap[bill.ClaimType].UpdateValue(CISTmp)
			} else {
				b := BillsScaleSum{
					Num:     1,
					CharSum: CISTmp,
				}
				claimTypeCharInputSumMap[bill.ClaimType] = &b
			}
			//default:
			//	global.GLog.Error("计算字符数", zap.Error(errors.New("没有该项目"+bill.ProCode)))
			//}
		}

		var settleStaffPercent float64 //结算字符数与员工字符数的比例
		//业务明细
		//billDetail := project.
		//err = proDb.Where("scan_at >= ? AND scan_at <= ? AND (stage = 5 or stage = 7) ", info.StartTime, info.EndTime).
		//	//Select("id", "bill_name").
		//	Preload("ProjectBlock").
		//	Preload("ProjectField", func(db *gorm.DB) *gorm.DB {
		//		return db.Select("id,name,code,bill_id,result_value,result_input,final_value,final_input")
		//	}).Find(&bills).Error

		averageCharCount, _ := decimal.NewFromInt(int64(charSum)).DivRound(decimal.NewFromInt(countBill), 4).Float64()
		averageCharInputCount, _ := decimal.NewFromInt(int64(charInputSum)).DivRound(decimal.NewFromInt(countBill), 4).Float64()
		averageStaffInputCount, _ := decimal.NewFromInt(*staffCharSum).DivRound(decimal.NewFromInt(countBill), 4).Float64()

		//比值
		if charSum > 0 {
			charPercent, _ = decimal.NewFromInt(int64(charInputSum)).DivRound(decimal.NewFromInt(int64(charSum)), 4).Float64()
			settleStaffPercent, _ = decimal.NewFromInt(*staffCharSum).DivRound(decimal.NewFromInt(int64(charSum)), 4).Float64()
		}
		claimsAccHasHis := 0.0
		claimsAccNoneHis := 0.0
		claimsAccMixHis := 0.0
		claimsAicHasHis := 0.0
		claimsAicNoneHis := 0.0
		claimsAicMixHis := 0.0
		if value, ok := claimTypeCharSumMap[5]; ok {
			claimsAccHasHis = value.GetPercent()
		}
		if value, ok := claimTypeCharSumMap[4]; ok {
			claimsAccNoneHis = value.GetPercent()
		}
		if value, ok := claimTypeCharSumMap[6]; ok {
			claimsAccMixHis = value.GetPercent()
		}
		if value, ok := claimTypeCharInputSumMap[5]; ok {
			claimsAicHasHis = value.GetPercent()
		}
		if value, ok := claimTypeCharInputSumMap[4]; ok {
			claimsAicNoneHis = value.GetPercent()
		}
		if value, ok := claimTypeCharInputSumMap[6]; ok {
			claimsAicMixHis = value.GetPercent()
		}

		var charSumObj = model3.CharSum{
			SumDate:               calDate,
			ProCode:               pro.Code,
			BillCount:             countBill,
			CharCount:             charSum,
			InputCharCount:        charInputSum,
			AverageCharCount:      averageCharCount,
			AverageInputCharCount: averageCharInputCount,
			CharPercent:           charPercent,

			//=========new add
			StaffInputCount:        int(*staffCharSum),
			SettleStaffPercent:     settleStaffPercent,
			AverageStaffInputCount: averageStaffInputCount,

			ClaimsAccHasHis:  claimsAccHasHis,
			ClaimsAccNoneHis: claimsAccNoneHis,
			ClaimsAccMixHis:  claimsAccMixHis,
			ClaimsAicHasHis:  claimsAicHasHis,
			ClaimsAicNoneHis: claimsAicNoneHis,
			ClaimsAicMixHis:  claimsAicMixHis,
		}
		//4: "无报销",     5: "有报销",     6: "混合型",
		err = global.GDb.Model(&model3.CharSum{}).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "sum_date"}, {Name: "pro_code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"bill_count", "char_count", "input_char_count",
				"average_char_count", "average_input_char_count", "char_percent",
				"staff_input_count", "average_staff_input_count", "settle_staff_percent",
				"claims_acc_has_his", "claims_acc_none_his", "claims_acc_mix_his",
				"claims_aic_has_his", "claims_aic_none_his", "claims_aic_mix_his"}),
		}).Create(&charSumObj).Error
	}

}

type BillsScaleSum struct {
	Num     int //案件量
	CharSum int //字符数
}

func (b *BillsScaleSum) UpdateValue(charNum int) {
	b.Num++
	b.CharSum += charNum
}

// 计算比值
func (b *BillsScaleSum) GetPercent() float64 {
	if b.Num > 0 {
		resFloat, _ := decimal.NewFromInt(int64(b.CharSum)).DivRound(decimal.NewFromInt(int64(b.Num)), 4).Float64()
		return resFloat
	}
	return 0.0
}

// CalculateMoney 计算实际完成金额
func CalculateMoney(proDB *gorm.DB, reportDate time.Time, dateType int) float64 {
	reportDateStr := reportDate.Format("2006-01-02")
	formatStr := "YYYY-MM-DD"
	if dateType == 2 {
		reportDateStr = reportDate.Format("2006-01")
		formatStr = "YYYY-MM"
	}
	var bills []model.ProjectBill
	totalMoney := 0.0
	err := proDB.Model(&model.ProjectBill{}).
		Select("scan_at,last_upload_at,insurance_type").
		Where("to_char(scan_at,'"+formatStr+"') = ?", reportDateStr).
		Find(&bills).Error
	if err != nil {
		global.GLog.Error("计算实际完成金额", zap.Error(err))
		return 0
	}
	for _, bill := range bills {
		t := bill.LastUploadAt.Sub(bill.ScanAt).Round(time.Second)
		types := bill.InsuranceType
		if types == "补录" {
			totalMoney += 2.85
		}
		if types == "" {
			if t <= time.Hour {
				totalMoney += 4.28
			}
			if t <= 2*time.Hour && time.Hour < t {
				totalMoney += 4.09
			}
			if t > 2*time.Hour {
				totalMoney += 2
			}
		}
	}
	total, _ := decimal.NewFromFloat(totalMoney).RoundFloor(4).Float64()
	return total
}
