package service

import (
	"github.com/shopspring/decimal"
	"server/global"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/module/pro_manager/model/response"
	"server/utils"
	"time"
)

func GetQualityAnalysis(info request.QuaAnaRes) (err error, list interface{}, total int64) {

	StartTime, _ := time.ParseInLocation("2006-01-02", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02", info.EndTime, time.Local)

	if info.Types == "1" {
		var proCode []string
		err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).Count(&total).
			Pluck("DISTINCT(pro_code)", &proCode).Error
		if err != nil {
			return err, nil, 0
		}
		var item []response.ProAnalysis
		for _, v := range proCode {
			var it response.ProAnalysis
			var ps int64
			it.ProCode = v
			err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
				Where("pro_code = ? ", v).Count(&ps).Error
			if err != nil {
				return err, nil, 0
			}
			it.Nums = ps
			if decimal.NewFromFloat(float64(it.Nums)/float64(total)).RoundCeil(2).String() == "1" {
				it.Percentage = "100%"
			} else {
				it.Percentage = decimal.NewFromFloat(float64(it.Nums)/float64(total)).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%"
			}

			item = append(item, it)
		}
		return nil, item, total
	}

	if info.Types == "2" {
		var fields []string
		err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
			Count(&total).Pluck("DISTINCT(wrong_field_name)", &fields).Error
		if err != nil {
			return err, nil, 0
		}
		var item []response.FieldAnalysis
		for _, v := range fields {
			var it response.FieldAnalysis
			var fs int64
			it.FiledName = v
			err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
				Where("wrong_field_name = ? ", v).Count(&fs).Error
			if err != nil {
				return err, nil, 0
			}
			it.Nums = fs
			if decimal.NewFromFloat(float64(it.Nums)/float64(total)).RoundCeil(2).String() == "1" {
				it.Percentage = "100%"
			} else {
				it.Percentage = decimal.NewFromFloat(float64(it.Nums)/float64(total)).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%"
			}

			item = append(item, it)
		}
		return nil, item, total
	}

	if info.Types == "3" {
		var code []string
		var op0 []string
		var op1 []string
		var op2 []string
		var opq []string
		err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
			Pluck("op0_responsible_code", &op0).Error
		err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
			Pluck("op1_responsible_code", &op1).Error
		err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
			Pluck("op2_responsible_code", &op2).Error
		err = global.GDb.Model(&model.Quality{}).Where("feedback_date >= ? AND feedback_date <= ? ", StartTime, EndTime).
			Pluck("opq_responsible_code", &opq).Error

		if err != nil {
			return err, nil, 0
		}
		code = append(code, op0...)
		code = append(code, op1...)
		code = append(code, op2...)
		code = append(code, opq...)
		sum := 0
		m := make(map[string]int, 0)
		for _, v := range code {
			if v != "" {
				m[v]++
				sum++
			}
		}

		var item []response.PeopleAnalysis
		for k, v := range m {
			var it response.PeopleAnalysis
			it.People = k + global.UserCodeName[k]

			it.Nums = int64(v)
			if decimal.NewFromFloat(float64(it.Nums)/float64(sum)).RoundCeil(2).String() == "1" {
				it.Percentage = "100%"
			} else {
				it.Percentage = decimal.NewFromFloat(float64(it.Nums)/float64(sum)).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%"
			}

			item = append(item, it)
		}
		return nil, item, total
	}
	return nil, nil, 0
}

func ExportQualityAnalysis(startTime, endTime string) (error, string) {
	basicPath := "files/质量分析导出/"
	bookName := "理赔差错分析" + startTime + "~" + endTime + ".xlsx"

	proInfo := request.QuaAnaRes{
		Types:     "1",
		StartTime: startTime,
		EndTime:   endTime,
	}
	err, proList, _ := GetQualityAnalysis(proInfo)
	if err != nil {
		return err, ""
	}
	err = utils.ExportBigExcel2(basicPath, bookName, "按项目", proList, "Sheet1")
	if err != nil {
		return err, ""
	}

	fieldInfo := request.QuaAnaRes{
		Types:     "2",
		StartTime: startTime,
		EndTime:   endTime,
	}
	err, filedList, _ := GetQualityAnalysis(fieldInfo)
	if err != nil {
		return err, ""
	}
	err = utils.ExportBigExcel2(basicPath, bookName, "按字段", filedList, "Sheet2")
	if err != nil {
		return err, ""
	}

	PeopleInfo := request.QuaAnaRes{
		Types:     "3",
		StartTime: startTime,
		EndTime:   endTime,
	}
	err, peopleList, _ := GetQualityAnalysis(PeopleInfo)
	if err != nil {
		return err, ""
	}
	err = utils.ExportBigExcel2(basicPath, bookName, "按人员", peopleList, "Sheet3")
	if err != nil {
		return err, ""
	}
	return err, basicPath + bookName
}
