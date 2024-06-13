package service

import (
	"fmt"
	"math"
	"server/global"
	"server/module/report_management/model"
	u "server/utils"
	"strconv"
	"strings"
)

func GetBusinessDownloadAnalysis(info model.BusinessAnalysisSearch) (err error, list interface{}, total int64) {
	var BusinessAnalyze []model.BusinessAnalysis
	if info.IsCheckAll {
		//all := make([]model.AllBusinessAnalysis, 0)
		var all []model.AllBusinessAnalysis
		var item model.AllBusinessAnalysis
		for dbCode, db := range global.ProDbMap {
			if strings.Index(dbCode, "_task") == -1 {
				db = db.Model(&model.BusinessAnalysis{})
				db = db.Where("types = ?", info.Types)
				if info.StartTime != "" {
					db = db.Where("create_at >= ?", info.StartTime)
				}
				if info.EndTime != "" {
					db = db.Where("create_at <= ?", info.EndTime)
				}
				err = db.Order("created_at asc").Find(&BusinessAnalyze).Error
				if err != nil {
					return err, nil, 0
				}
				for _, v := range BusinessAnalyze {
					item.AHalfPastEight += v.AHalfPastEight
					item.AHalfPastNine += v.AHalfPastNine
					item.AHalfPastTen += v.AHalfPastTen
					item.AHalfPastEleven += v.AHalfPastEleven
					item.AHalfPastTwelve += v.AHalfPastTwelve
					item.AHalfPastThirteen += v.AHalfPastThirteen
					item.AHalfPastFourteen += v.AHalfPastFourteen
					item.AHalfPastFifteen += v.AHalfPastFifteen
					item.AHalfPastSixteen += v.AHalfPastSixteen
					item.AHalfPastSeventeen += v.AHalfPastSeventeen
					item.AHalfPastEighteen += v.AHalfPastEighteen
					item.BeforeZeroHour += v.BeforeZeroHour
				}
				item.CreateAt = info.StartTime + " - " + info.EndTime
				item.ProCode = dbCode
				all = append(all, item)
			}
		}
		var calculateBusinessAnalyzeSum model.AllBusinessAnalysis
		for _, v := range all {
			calculateBusinessAnalyzeSum.ProCode = "合计"
			calculateBusinessAnalyzeSum.CreateAt = "合计"
			calculateBusinessAnalyzeSum.AHalfPastEight += v.AHalfPastEight
			calculateBusinessAnalyzeSum.AHalfPastNine += v.AHalfPastNine
			calculateBusinessAnalyzeSum.AHalfPastTen += v.AHalfPastTen
			calculateBusinessAnalyzeSum.AHalfPastEleven += v.AHalfPastEleven
			calculateBusinessAnalyzeSum.AHalfPastTwelve += v.AHalfPastTwelve
			calculateBusinessAnalyzeSum.AHalfPastThirteen += v.AHalfPastThirteen
			calculateBusinessAnalyzeSum.AHalfPastFourteen += v.AHalfPastFourteen
			calculateBusinessAnalyzeSum.AHalfPastFifteen += v.AHalfPastFifteen
			calculateBusinessAnalyzeSum.AHalfPastSixteen += v.AHalfPastSixteen
			calculateBusinessAnalyzeSum.AHalfPastSeventeen += v.AHalfPastSeventeen
			calculateBusinessAnalyzeSum.AHalfPastEighteen += v.AHalfPastEighteen
			calculateBusinessAnalyzeSum.BeforeZeroHour += v.BeforeZeroHour
		}
		all = append(all, calculateBusinessAnalyzeSum)
		return err, all, int64(len(all))
	} else {
		//连接数据库
		db := global.ProDbMap[info.ProCode]
		if db == nil {
			return global.ProDbErr, nil, 0
		}
		db = db.Model(&model.BusinessAnalysis{})
		db = db.Where("types = ?", info.Types)
		if info.StartTime != "" {
			db = db.Where("create_at >= ?", info.StartTime)
		}
		if info.EndTime != "" {
			db = db.Where("create_at <= ?", info.EndTime)
		}
		err = db.Count(&total).Error
		err = db.Order("created_at asc").Find(&BusinessAnalyze).Error
		if err != nil {
			return err, nil, 0
		}
		var calculateBusinessAnalyzeSum model.BusinessAnalysis
		for _, v := range BusinessAnalyze {
			calculateBusinessAnalyzeSum.CreateAt = "合计"
			calculateBusinessAnalyzeSum.AHalfPastEight += v.AHalfPastEight
			calculateBusinessAnalyzeSum.AHalfPastNine += v.AHalfPastNine
			calculateBusinessAnalyzeSum.AHalfPastTen += v.AHalfPastTen
			calculateBusinessAnalyzeSum.AHalfPastEleven += v.AHalfPastEleven
			calculateBusinessAnalyzeSum.AHalfPastTwelve += v.AHalfPastTwelve
			calculateBusinessAnalyzeSum.AHalfPastThirteen += v.AHalfPastThirteen
			calculateBusinessAnalyzeSum.AHalfPastFourteen += v.AHalfPastFourteen
			calculateBusinessAnalyzeSum.AHalfPastFifteen += v.AHalfPastFifteen
			calculateBusinessAnalyzeSum.AHalfPastSixteen += v.AHalfPastSixteen
			calculateBusinessAnalyzeSum.AHalfPastSeventeen += v.AHalfPastSeventeen
			calculateBusinessAnalyzeSum.AHalfPastEighteen += v.AHalfPastEighteen
			calculateBusinessAnalyzeSum.BeforeZeroHour += v.BeforeZeroHour
		}
		BusinessAnalyze = append(BusinessAnalyze, calculateBusinessAnalyzeSum)
		return err, BusinessAnalyze, total
	}
}

func GetBusinessUploadAnalysis(info model.BusinessAnalysisSearch) (err error, list interface{}, total int64) {
	var BusinessAnalyze []model.BusinessAnalysis
	var BusinessUploadAnalyze []model.BusinessUploadAnalysis
	if info.ProCode == "全部" {
		//全部
		if info.IsCheckAll {
			//按整体
			//all := make([]model.AllBusinessUploadAnalysis, 0)
			var all []model.AllBusinessUploadAnalysis
			var item model.AllBusinessUploadAnalysis
			for dbCode, db := range global.ProDbMap {
				if strings.Index(dbCode, "_task") == -1 {
					db = db.Model(&model.BusinessUploadAnalysis{})
					if info.StartTime != "" {
						db = db.Where("create_at >= ?", info.StartTime)
					}
					if info.EndTime != "" {
						db = db.Where("create_at <= ?", info.EndTime)
					}
					err = db.Order("created_at asc").Find(&BusinessUploadAnalyze).Error
					for _, v := range BusinessUploadAnalyze {
						item.VolumeOfBusiness = v.VolumeOfBusiness
						item.OneHours = v.OneHours
						item.TwoHours = v.TwoHours
						item.ThreeHours = v.ThreeHours
						item.MoreThanThreeHours = v.MoreThanThreeHours
						item.TheAverageTime = v.TheAverageTime
					}
					item.TheAverageTime = decimals(item.TheAverageTime / float64(len(BusinessUploadAnalyze)))
					item.OneHoursRate = percent(float64(item.OneHours), float64(item.VolumeOfBusiness))
					item.TwoHoursRate = percent(float64(item.TwoHours), float64(item.VolumeOfBusiness))
					item.ThreeHoursRate = percent(float64(item.ThreeHours), float64(item.VolumeOfBusiness))
					item.MoreThanThreeHoursRate = percent(float64(item.MoreThanThreeHours), float64(item.VolumeOfBusiness))
					item.CreateAt = info.StartTime + " - " + info.EndTime
					item.ProCode = dbCode
					all = append(all, item)
				}
			}
			var calculateBusinessAnalyzeSum model.AllBusinessUploadAnalysis
			for _, v := range all {
				calculateBusinessAnalyzeSum.ProCode = "合计"
				calculateBusinessAnalyzeSum.CreateAt = "合计"
				calculateBusinessAnalyzeSum.VolumeOfBusiness += v.VolumeOfBusiness
				calculateBusinessAnalyzeSum.TheAverageTime += v.TheAverageTime
				calculateBusinessAnalyzeSum.OneHours += v.OneHours
				calculateBusinessAnalyzeSum.TwoHours += v.TwoHours
				calculateBusinessAnalyzeSum.ThreeHours += v.ThreeHours
				calculateBusinessAnalyzeSum.MoreThanThreeHours += v.MoreThanThreeHours
			}
			calculateBusinessAnalyzeSum.TheAverageTime = math.Round(calculateBusinessAnalyzeSum.TheAverageTime*100) / 100
			//calculateBusinessAnalyzeSum.TheAverageTime = decimals(calculateBusinessAnalyzeSum.TheAverageTime / float64(len(all)))
			calculateBusinessAnalyzeSum.OneHoursRate = percent(float64(calculateBusinessAnalyzeSum.OneHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			calculateBusinessAnalyzeSum.TwoHoursRate = percent(float64(calculateBusinessAnalyzeSum.TwoHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			calculateBusinessAnalyzeSum.ThreeHoursRate = percent(float64(calculateBusinessAnalyzeSum.ThreeHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			calculateBusinessAnalyzeSum.MoreThanThreeHoursRate = percent(float64(calculateBusinessAnalyzeSum.MoreThanThreeHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			all = append(all, calculateBusinessAnalyzeSum)
			return err, all, int64(len(all))
		} else {
			//按明细
			//all := make([]model.AllBusinessAnalysis, 0)
			var all []model.AllBusinessAnalysis
			var item model.AllBusinessAnalysis
			for dbCode, db := range global.ProDbMap {
				if strings.Index(dbCode, "_task") == -1 {
					db = db.Model(&model.BusinessAnalysis{})
					db = db.Where("types = ?", info.Types)
					if info.StartTime != "" {
						db = db.Where("create_at >= ?", info.StartTime)
					}
					if info.EndTime != "" {
						db = db.Where("create_at <= ?", info.EndTime)
					}
					err = db.Order("created_at asc").Find(&BusinessAnalyze).Error
					if err != nil {
						return err, nil, 0
					}
					fmt.Println("BusinessAnalyze", BusinessAnalyze)
					for _, v := range BusinessAnalyze {
						item.AHalfPastEight += v.AHalfPastEight
						item.AHalfPastNine += v.AHalfPastNine
						item.AHalfPastTen += v.AHalfPastTen
						item.AHalfPastEleven += v.AHalfPastEleven
						item.AHalfPastTwelve += v.AHalfPastTwelve
						item.AHalfPastThirteen += v.AHalfPastThirteen
						item.AHalfPastFourteen += v.AHalfPastFourteen
						item.AHalfPastFifteen += v.AHalfPastFifteen
						item.AHalfPastSixteen += v.AHalfPastSixteen
						item.AHalfPastSeventeen += v.AHalfPastSeventeen
						item.AHalfPastEighteen += v.AHalfPastEighteen
						item.BeforeZeroHour += v.BeforeZeroHour
					}
					item.CreateAt = info.StartTime + " - " + info.EndTime
					item.ProCode = dbCode
					all = append(all, item)
				}
			}
			var calculateBusinessAnalyzeSum model.AllBusinessAnalysis
			for _, v := range all {
				calculateBusinessAnalyzeSum.ProCode = "合计"
				calculateBusinessAnalyzeSum.CreateAt = "合计"
				calculateBusinessAnalyzeSum.AHalfPastEight += v.AHalfPastEight
				calculateBusinessAnalyzeSum.AHalfPastNine += v.AHalfPastNine
				calculateBusinessAnalyzeSum.AHalfPastTen += v.AHalfPastTen
				calculateBusinessAnalyzeSum.AHalfPastEleven += v.AHalfPastEleven
				calculateBusinessAnalyzeSum.AHalfPastTwelve += v.AHalfPastTwelve
				calculateBusinessAnalyzeSum.AHalfPastThirteen += v.AHalfPastThirteen
				calculateBusinessAnalyzeSum.AHalfPastFourteen += v.AHalfPastFourteen
				calculateBusinessAnalyzeSum.AHalfPastFifteen += v.AHalfPastFifteen
				calculateBusinessAnalyzeSum.AHalfPastSixteen += v.AHalfPastSixteen
				calculateBusinessAnalyzeSum.AHalfPastSeventeen += v.AHalfPastSeventeen
				calculateBusinessAnalyzeSum.AHalfPastEighteen += v.AHalfPastEighteen
				calculateBusinessAnalyzeSum.BeforeZeroHour += v.BeforeZeroHour
			}
			all = append(all, calculateBusinessAnalyzeSum)
			return err, all, int64(len(all))
		}
	} else {
		//项目
		if info.IsCheckAll {
			//按整体
			//连接数据库
			db := global.ProDbMap[info.ProCode]
			if db == nil {
				return global.ProDbErr, nil, 0
			}
			db = db.Model(&model.BusinessUploadAnalysis{})
			if info.StartTime != "" {
				db = db.Where("create_at >= ?", info.StartTime)
			}
			if info.EndTime != "" {
				db = db.Where("create_at <= ?", info.EndTime)
			}
			err = db.Count(&total).Error
			err = db.Order("created_at asc").Find(&BusinessUploadAnalyze).Error
			if err != nil {
				return err, nil, 0
			}
			var calculateBusinessAnalyzeSum model.BusinessUploadAnalysis
			for i, v := range BusinessUploadAnalyze {
				BusinessUploadAnalyze[i].OneHoursRate = percent(float64(v.OneHours), float64(v.VolumeOfBusiness))
				BusinessUploadAnalyze[i].TwoHoursRate = percent(float64(v.TwoHours), float64(v.VolumeOfBusiness))
				BusinessUploadAnalyze[i].ThreeHoursRate = percent(float64(v.ThreeHours), float64(v.VolumeOfBusiness))
				BusinessUploadAnalyze[i].MoreThanThreeHoursRate = percent(float64(v.MoreThanThreeHours), float64(v.VolumeOfBusiness))

				calculateBusinessAnalyzeSum.CreateAt = "合计"
				calculateBusinessAnalyzeSum.VolumeOfBusiness += v.VolumeOfBusiness
				calculateBusinessAnalyzeSum.TheAverageTime += v.TheAverageTime
				calculateBusinessAnalyzeSum.OneHours += v.OneHours
				calculateBusinessAnalyzeSum.TwoHours += v.TwoHours
				calculateBusinessAnalyzeSum.ThreeHours += v.ThreeHours
				calculateBusinessAnalyzeSum.MoreThanThreeHours += v.MoreThanThreeHours
			}
			calculateBusinessAnalyzeSum.TheAverageTime = decimals(calculateBusinessAnalyzeSum.TheAverageTime / float64(len(BusinessUploadAnalyze)))
			calculateBusinessAnalyzeSum.OneHoursRate = percent(float64(calculateBusinessAnalyzeSum.OneHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			calculateBusinessAnalyzeSum.TwoHoursRate = percent(float64(calculateBusinessAnalyzeSum.TwoHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			calculateBusinessAnalyzeSum.ThreeHoursRate = percent(float64(calculateBusinessAnalyzeSum.ThreeHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			calculateBusinessAnalyzeSum.MoreThanThreeHoursRate = percent(float64(calculateBusinessAnalyzeSum.MoreThanThreeHours), float64(calculateBusinessAnalyzeSum.VolumeOfBusiness))
			BusinessUploadAnalyze = append(BusinessUploadAnalyze, calculateBusinessAnalyzeSum)
			return err, BusinessUploadAnalyze, total
		} else {
			//按明细
			//连接数据库
			db := global.ProDbMap[info.ProCode]
			if db == nil {
				return global.ProDbErr, nil, 0
			}
			db = db.Model(&model.BusinessAnalysis{})
			db = db.Where("types = ?", info.Types)
			if info.StartTime != "" {
				db = db.Where("create_at >= ?", info.StartTime)
			}
			if info.EndTime != "" {
				db = db.Where("create_at <= ?", info.EndTime)
			}
			err = db.Count(&total).Error
			err = db.Order("created_at asc").Find(&BusinessAnalyze).Error
			if err != nil {
				return err, nil, 0
			}
			var BusinessAnalyzeRe []model.BusinessAnalysis
			var item model.BusinessAnalysis
			for _, v := range BusinessAnalyze {
				item.CreateAt = v.CreateAt
				item.AHalfPastEight = v.AHalfPastEight
				item.AHalfPastNine = v.AHalfPastNine
				item.AHalfPastTen = v.AHalfPastTen
				item.AHalfPastEleven = v.AHalfPastEleven
				item.AHalfPastTwelve = v.AHalfPastTwelve
				item.AHalfPastThirteen = v.AHalfPastThirteen
				item.AHalfPastFourteen = v.AHalfPastFourteen
				item.AHalfPastFifteen = v.AHalfPastFifteen
				item.AHalfPastSixteen = v.AHalfPastSixteen
				item.AHalfPastSeventeen = v.AHalfPastSeventeen
				item.AHalfPastEighteen = v.AHalfPastEighteen
				item.BeforeZeroHour = v.BeforeZeroHour
				BusinessAnalyzeRe = append(BusinessAnalyzeRe, item)
			}
			var calculateBusinessAnalyzeSum model.BusinessAnalysis
			for _, v := range BusinessAnalyze {
				calculateBusinessAnalyzeSum.CreateAt = "合计"
				calculateBusinessAnalyzeSum.AHalfPastEight += v.AHalfPastEight
				calculateBusinessAnalyzeSum.AHalfPastNine += v.AHalfPastNine
				calculateBusinessAnalyzeSum.AHalfPastTen += v.AHalfPastTen
				calculateBusinessAnalyzeSum.AHalfPastEleven += v.AHalfPastEleven
				calculateBusinessAnalyzeSum.AHalfPastTwelve += v.AHalfPastTwelve
				calculateBusinessAnalyzeSum.AHalfPastThirteen += v.AHalfPastThirteen
				calculateBusinessAnalyzeSum.AHalfPastFourteen += v.AHalfPastFourteen
				calculateBusinessAnalyzeSum.AHalfPastFifteen += v.AHalfPastFifteen
				calculateBusinessAnalyzeSum.AHalfPastSixteen += v.AHalfPastSixteen
				calculateBusinessAnalyzeSum.AHalfPastSeventeen += v.AHalfPastSeventeen
				calculateBusinessAnalyzeSum.AHalfPastEighteen += v.AHalfPastEighteen
				calculateBusinessAnalyzeSum.BeforeZeroHour += v.BeforeZeroHour
			}
			BusinessAnalyzeRe = append(BusinessAnalyzeRe, calculateBusinessAnalyzeSum)
			return err, BusinessAnalyzeRe, total
		}
	}
}

func ExportBusinessDownloadAnalysis(info model.BusinessAnalysisSearch) (err error, path, name string) {
	err, analysis, _ := GetBusinessDownloadAnalysis(info)
	if err != nil {
		return err, "", ""
	}
	s := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	e := strings.Replace(info.EndTime, " 00:00:00", "", -1)
	bookName := info.ProCode + "-" + "来量分析" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "来量分析导出/" + info.ProCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcel(basicPath, bookName, "", analysis)
	return err, basicPath + bookName, bookName
}

func ExportBusinessUploadAnalysis(info model.BusinessAnalysisSearch) (err error, path, name string) {
	err, analysis, _ := GetBusinessUploadAnalysis(info)
	if err != nil {
		return err, "", ""
	}
	s := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	e := strings.Replace(info.EndTime, " 00:00:00", "", -1)
	bookName := info.ProCode + "-" + "回传分析" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "回传分析导出/" + info.ProCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcel(basicPath, bookName, "", analysis)
	return err, basicPath + bookName, bookName
}

func percent(num1, num2 float64) string {
	if num2 == 0 {
		return "0.00%"
	}
	if num1/num2 == 1 {
		return "100%"
	}
	return fmt.Sprintf("%.2f", (num1/num2)*100) + "%"
}

func decimals(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
