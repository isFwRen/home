/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/28 15:32
 */

package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"server/global"
	model3 "server/module/homepage/model"
	request2 "server/module/homepage/model/request"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	model4 "server/module/report_management/model"
	"server/module/report_management/model/request"
	"server/utils"
	"sort"
	"strconv"
	"time"
)

var dateFormat = map[int]string{
	1: "to_char(scan_at,'YYYY-MM-DD') as count_date",                                          //日
	2: "concat_ws('-',extract (year from scan_at),extract (week from scan_at)) as count_date", //周
	3: "to_char(scan_at,'YYYY-MM') as count_date",                                             //月
	4: "to_char(scan_at,'YYYY') as count_date",                                                //年
}

// GetProReport 获取项目报表
func GetProReport(reportReq request.BusinessReportReq, reportType int) (err error, obj map[string]interface{}) {
	obj = make(map[string]interface{}, 0)
	proConfList := map[string]model2.SysProjectCache{}
	if reportReq.ProCode == "all" {
		proConfList = global.GProConf
	} else if _, ok := global.GProConf[reportReq.ProCode]; ok {
		proConfList = map[string]model2.SysProjectCache{
			reportReq.ProCode: global.GProConf[reportReq.ProCode],
		}
	} else {
		return errors.New("项目参数错误proCode[all或者B0118等]"), obj
	}

	for _, proConf := range proConfList {
		if global.ProDbMap[proConf.Code] == nil {
			continue
		}
		var list = make([]map[string]interface{}, 0)
		switch reportType {
		case 1:
			err, list = GetBusinessReportByProCode(reportReq, proConf.Code)
		case 2:
			err, list = GetAgingReportByProCode(reportReq, proConf.Code)
		case 3:
			err, list = GetDealTimeReportByProCode(reportReq, proConf.Code)
		}
		if err != nil {
			return err, obj
		}
		//构造没有数据的日期
		obj[proConf.Code] = formatList(reportReq, list)
	}
	return nil, obj
}

// GetBusinessReportByProCode 获取单个项目业务量报表
func GetBusinessReportByProCode(businessReportReq request.BusinessReportReq, proCode string) (err error, list []map[string]interface{}) {
	err = global.ProDbMap[proCode].Model(&model.ProjectBill{}).
		Raw("SELECT "+dateFormat[businessReportReq.Type]+",count(1) FROM project_bills where (stage = 5 or stage = 7)  and scan_at BETWEEN ? AND ? GROUP BY count_date",
			businessReportReq.StartTime, businessReportReq.EndTime).Scan(&list).Error
	return err, list
}

// GetAgingReportByProCode 获取单个项目的时效报表
func GetAgingReportByProCode(agingReportReq request.BusinessReportReq, proCode string) (err error, list []map[string]interface{}) {
	//err = global.ProDbMap[proCode].Model(&model.ProjectBill{}).
	//	Raw("SELECT a.count_date,CAST(a.count/b.count as DECIMAL(18,4)) as count from(SELECT "+dateFormat[agingReportReq.Type]+",count(1)::numeric as count FROM project_bills WHERE  (stage = 5 or stage = 7) and is_timeout = false and scan_at BETWEEN ? AND ? GROUP BY count_date)a LEFT JOIN(SELECT "+dateFormat[agingReportReq.Type]+",count(1)::numeric as count FROM project_bills WHERE  stage = 5 and scan_at BETWEEN ? AND ? GROUP BY count_date)b  on a.count_date = b.count_date",
	//		agingReportReq.StartTime, agingReportReq.EndTime, agingReportReq.StartTime, agingReportReq.EndTime).Scan(&list).Error
	//日报根据业务明细表来  获取单个项目的业务明细表
	var info request.BusinessDetailsSearch
	info.ProCode = agingReportReq.ProCode
	if agingReportReq.ProCode == "all" {
		info.ProCode = proCode
	}
	info.StartTime = agingReportReq.StartTime.Format("2006-01-02")
	info.EndTime = agingReportReq.EndTime.Format("2006-01-02")
	info.Type = agingReportReq.Type

	//处理量太大先注释 要此功能需重写
	//err, reportList, _ := report.GetBusinessDetails(info)
	//用日期进行分类
	//reportMarshal, err := json.Marshal(reportList)
	//var reportMap []map[string]interface{}
	//err = json.Unmarshal([]byte(reportMarshal), &reportMap)
	var reportDetails []model4.BusinessDetailsExport
	//err = mapstructure.Decode(reportList, &reportDetails)
	if err != nil {
		err.Error()
	}
	err, resultMap := FindReportSameData(reportDetails, agingReportReq)
	// TODO 补充周报
	if agingReportReq.Type == 2 {

	}
	return err, resultMap
}

// FindReportSameData  找到相同日期的数据
func FindReportSameData(reportDetails []model4.BusinessDetailsExport, agingReportReq request.BusinessReportReq) (err error, resultMap []map[string]interface{}) {
	//记录已经循环过的数据
	seenDate := make(map[string]bool)
	//记录日期返回的数据
	result := make([]map[string]interface{}, 0)
	finalMap := make(map[string]bool)
	//便利过的年
	finalYearMap := make(map[string]bool)
	// 封装周数据
	weekFinalMap := make([]map[string]interface{}, 0)
	//封装年数据
	yearFinalMap := make([]map[string]interface{}, 0)
	//相同周
	weekMap := make(map[string]bool, 0)
	list := make([]map[string]interface{}, 0)
	for i := 0; i < len(reportDetails); i++ {
		if agingReportReq.Type == 2 {
			fmt.Println("=================我是周-------------------------")
			//////////////////////////////// 周 ///////////////////////////
			parse, err := time.Parse("2006-01-02", reportDetails[i].CreateAt)
			if err != nil {
				return err, nil
			}
			year := parse.Year()
			firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
			days := int(parse.Sub(firstDayOfYear).Hours() / 24)
			weekNumber := days/7 + 2
			dateFinal := strconv.Itoa(year) + "-" + strconv.Itoa(weekNumber)
			reportSameWeekList := make([]model4.BusinessDetailsExport, 0)
			if !weekMap[reportDetails[i].CreateAt] {
				for j := 0; j < len(reportDetails); j++ {
					if reportDetails[i].CreateAt == reportDetails[j].CreateAt {
						reportSameWeekList = append(reportSameWeekList, reportDetails[i])
						seenDate[reportDetails[i].CreateAt] = true
					}
				}
			}
			var outTimeWeekCase int
			for _, reportCase := range reportSameWeekList {
				if reportCase.IsTheTimeOut != "" && reportCase.IsTheTimeOut == "是" {
					outTimeWeekCase++
				}
			}
			if len(reportSameWeekList) > 0 {
				// 计算时效保障率
				var timeGuaranteeRate float64
				timeGuaranteeRate = (float64(len(reportSameWeekList)) - float64(outTimeWeekCase)) / float64(len(reportSameWeekList))
				lastTo := strconv.FormatFloat(timeGuaranteeRate, 'f', 2, 64)
				for _, item := range reportSameWeekList {
					countDate := item.CreateAt
					timeGuaranteeRateCount := lastTo
					if !finalMap[countDate] {
						weekFinalMap = append(weekFinalMap, map[string]interface{}{
							"count":      timeGuaranteeRateCount,
							"count_date": dateFinal,
						})
						finalMap[countDate] = true
					}
				}
			}
			resultMap = formatListReturnMap(agingReportReq, list)
			for ii, finalKey := range weekFinalMap {
				for j, m := range resultMap {
					if finalKey["count_date"] == m["countDate"] {
						resultMap[j]["count"] = weekFinalMap[ii]["count"]
						resultMap[j]["count_date"] = weekFinalMap[ii]["count_date"]
					}
				}
			}
			return err, resultMap
			//////////////////////////////// 周 ///////////////////////////
		} else if agingReportReq.Type == 4 {
			//reportSameYearList := make([]model4.BusinessDetailsExport, 0)
			//////////////////////////////// 年 ///////////////////////////
			fmt.Println("=================我是年-------------------------")

			outTimeYearCase := 0
			for _, reportCase := range reportDetails {
				if reportCase.IsTheTimeOut != "" && reportCase.IsTheTimeOut == "是" {
					outTimeYearCase++
				}
			}
			// 计算时效保障率
			var timeGuaranteeRateYear float64
			timeGuaranteeRateYear = (float64(len(reportDetails)) - float64(outTimeYearCase)) / float64(len(reportDetails))
			for _, item := range reportDetails {
				parse, err := time.Parse("2006-01-02", item.CreateAt)
				if err != nil {
					return err, nil
				}
				countDate := strconv.Itoa(parse.Year())
				timeGuaranteeRateCount := timeGuaranteeRateYear
				rateYear := strconv.FormatFloat(timeGuaranteeRateCount, 'f', 2, 64)
				if !finalYearMap[countDate] {
					yearFinalMap = append(yearFinalMap, map[string]interface{}{
						"count":      rateYear,
						"count_date": countDate,
					})
					finalYearMap[countDate] = true
				}
			}
			return err, yearFinalMap
			//////////////////////////////// 年 ///////////////////////////
		}
		//找到相同时间的
		reportSameList := make([]model4.BusinessDetailsExport, 0)
		if !seenDate[reportDetails[i].CreateAt] {
			for j := 0; j < len(reportDetails); j++ {
				if reportDetails[i].CreateAt == reportDetails[j].CreateAt {
					reportSameList = append(reportSameList, reportDetails[i])
					seenDate[reportDetails[i].CreateAt] = true

				}
			}
		}
		// 把相同时间的数据进行计算  查看是否超时 统计时效保障率  时效保障率 = （总案件量 - 超时案件量）/ 总案件量  超时案件量=是否超时为”是“的案件
		//统计超时案件
		outTimeCase := 0
		for _, reportCase := range reportSameList {
			if reportCase.IsTheTimeOut != "" && reportCase.IsTheTimeOut == "是" {
				outTimeCase++
			}
		}
		if len(reportSameList) > 0 {
			// 计算时效保障率
			timeGuaranteeRate := (len(reportSameList) - outTimeCase) / len(reportSameList)
			for _, item := range reportSameList {
				countDate := item.CreateAt
				timeGuaranteeRateCount := timeGuaranteeRate
				if !finalMap[countDate] {
					result = append(result, map[string]interface{}{
						"count":      timeGuaranteeRateCount,
						"count_date": countDate,
					})
					finalMap[countDate] = true
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i]["count_date"].(string) < result[j]["count_date"].(string)
	})
	return err, result
}

// GetDealTimeReportByProCode 获取单个项目的处理时间报表
func GetDealTimeReportByProCode(dealTimeReportReq request.BusinessReportReq, proCode string) (err error, list []map[string]interface{}) {
	err = global.ProDbMap[proCode].Model(&model.ProjectBill{}).
		Raw("SELECT "+dateFormat[dealTimeReportReq.Type]+",date_part('epoch',sum(last_upload_at - download_at)/count(1)) as count FROM project_bills WHERE  (stage = 5 or stage = 7) and scan_at BETWEEN ? AND ? GROUP BY count_date",
			//Raw("SELECT "+dateFormat+",sum(export_at - download_at)/count(1) as count FROM project_bills WHERE  stage = 3 and scan_at BETWEEN ? AND ? GROUP BY count_date",
			dealTimeReportReq.StartTime, dealTimeReportReq.EndTime).Scan(&list).Error
	return err, list
}

// formatList 构造没有数据的日期
func formatList(businessReportReq request.BusinessReportReq, list []map[string]interface{}) interface{} {
	var businessReports = make([]interface{}, 0)
	if businessReportReq.Type == 2 {
		weekday := int(businessReportReq.StartTime.Weekday() - 1)
		endStartWeekday := int(businessReportReq.EndTime.Weekday())
		if weekday == -1 {
			weekday = 6
		}
		//获取当前周开始,结束时间
		businessReportReq.StartTime = businessReportReq.StartTime.Add(-24 * time.Duration(weekday) * 60 * time.Minute)
		businessReportReq.EndTime = businessReportReq.EndTime.Add(24 * time.Duration(7-endStartWeekday) * 60 * time.Minute)
		//fmt.Println("1111")
	}
	newDate := businessReportReq.StartTime
	for {
		if newDate.After(businessReportReq.EndTime) {
			return businessReports
		}
		//日
		hasDay := 1
		dStr := newDate.Format("2006-01-02")
		switch businessReportReq.Type {
		case 2:
			//周
			hasDay = 7
			y, w := newDate.ISOWeek()
			global.GLog.Info("当前周", zap.Any(strconv.Itoa(y), w))
			dStr = fmt.Sprintf("%d-%d", y, w)
		case 3:
			//月
			hasDay = utils.GetYearMonthToDay(newDate.Year(), int(newDate.Month()))
			dStr = newDate.Format("2006-01")
		case 4:
			//年
			hasDay = 365
			if newDate.Year()%4 == 0 {
				hasDay = 366
			}
			dStr = newDate.Format("2006")
		}
		newBusinessReport := map[string]interface{}{
			"countDate": dStr,
			"count":     0,
		}
		for _, businessReport := range list {
			if businessReport["count_date"] == dStr {
				newBusinessReport["count"] = businessReport["count"]
				continue
			}
		}
		businessReports = append(businessReports, newBusinessReport)

		newDate = newDate.Add(time.Duration(hasDay) * 24 * 60 * time.Minute)
		//global.GLog.Info("dddd", zap.Any("ddd", newDate))

	}
}

// formatList 构造没有数据的日期
func formatListReturnMap(businessReportReq request.BusinessReportReq, list []map[string]interface{}) (finalMaps []map[string]interface{}) {
	var businessReports = make([]map[string]interface{}, 0)
	if businessReportReq.Type == 2 {
		weekday := int(businessReportReq.StartTime.Weekday() - 1)
		endStartWeekday := int(businessReportReq.EndTime.Weekday())
		if weekday == -1 {
			weekday = 6
		}
		//获取当前周开始,结束时间
		businessReportReq.StartTime = businessReportReq.StartTime.Add(-24 * time.Duration(weekday) * 60 * time.Minute)
		businessReportReq.EndTime = businessReportReq.EndTime.Add(24 * time.Duration(7-endStartWeekday) * 60 * time.Minute)

	}
	newDate := businessReportReq.StartTime
	for {
		if newDate.After(businessReportReq.EndTime) {
			return businessReports
		}
		//日
		hasDay := 1
		dStr := newDate.Format("2006-01-02")
		switch businessReportReq.Type {
		case 2:
			//周
			hasDay = 7
			y, w := newDate.ISOWeek()
			global.GLog.Info("当前周", zap.Any(strconv.Itoa(y), w))
			dStr = fmt.Sprintf("%d-%d", y, w)
		case 3:
			//月
			hasDay = utils.GetYearMonthToDay(newDate.Year(), int(newDate.Month()))
			dStr = newDate.Format("2006-01")
		case 4:
			//年
			hasDay = 365
			if newDate.Year()%4 == 0 {
				hasDay = 366
			}
			dStr = newDate.Format("2006")
		}

		newBusinessReport := map[string]interface{}{
			"countDate": dStr,
			"count":     0,
		}
		for _, businessReport := range list {
			if businessReport["count_date"] == dStr {
				newBusinessReport["count"] = businessReport["count"]
				continue
			}
		}
		businessReports = append(businessReports, newBusinessReport)
		newDate = newDate.Add(time.Duration(hasDay) * 24 * 60 * time.Minute)
		//global.GLog.Info("dddd", zap.Any("ddd", newDate))

	}
}

// SetOtherReportInfo 设置项目日报其他信息
func SetOtherReportInfo(proReportOtherInfoReq request2.ProReportOtherInfoReq) (err error) {
	proReportOtherInfo := proReportOtherInfoReq.ProReportOtherInfo
	proReportOtherInfo.ReportDate = time.Date(proReportOtherInfo.ReportDate.Year(),
		proReportOtherInfo.ReportDate.Month(),
		proReportOtherInfo.ReportDate.Day(),
		0, 0, 0, 0, time.UTC)
	for i, _ := range proReportOtherInfoReq.ProReport {
		proReportOtherInfoReq.ProReport[i].ReportDate = proReportOtherInfo.ReportDate
	}
	if err != nil {
		return err
	}
	tx := global.GDb.Begin()
	err = tx.Model(&model3.ProReportOtherInfo{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "report_date"}},
		UpdateAll: true,
	}).Create(&proReportOtherInfo).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&model3.ProReport{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "report_date"}, {Name: "pro_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"predict_value"}),
	}).Create(&proReportOtherInfoReq.ProReport).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
