package service

import (
	"fmt"
	"server/global"
	model2 "server/module/download/model"
	"server/module/load/model"
	"time"

	"github.com/shopspring/decimal"
	// sumModel "server/module/report_management/model"
	// "strings"
)

func SelectBlockByUser(proCode string, code string, ratio decimal.Decimal) (err error, configs []string) {
	var total int64 = 0
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	if code != "" {
		db = db.Where("op0_code = ? or op1_code = ? or op2_code = ? or opq_code = ?", code, code, code, code)
	}
	// db = db.Where("created_at BETWEEN ? AND ?  ", announcementPageReq.StartTime, announcementPageReq.EndTime)
	err = db.Count(&total).Error
	fmt.Println("--------------------total-----err------------------", total, err)
	aaa := decimal.NewFromInt(total)
	ratio = ratio.Div(decimal.NewFromInt(100))
	limit := int(aaa.Mul(ratio).IntPart())
	// fmt.Println("--------------------limit-----------------------", limit)
	var configsRes []string
	err = db.Select("id").Limit(limit).Find(&configsRes).Error
	return err, configsRes
}

func SelectBlockByBlock(proCode string, code string, ratio decimal.Decimal) (err error, configs []string) {
	var total int64 = 0
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	if code != "" {
		db = db.Where("code = ?", code)
	}
	// db = db.Where("created_at BETWEEN ? AND ?  ", announcementPageReq.StartTime, announcementPageReq.EndTime)
	err = db.Count(&total).Error
	fmt.Println("--------------------total-----err------------------", total, err)
	aaa := decimal.NewFromInt(total)
	ratio = ratio.Div(decimal.NewFromInt(100))
	limit := int(aaa.Mul(ratio).IntPart())
	// fmt.Println("--------------------limit-----------------------", limit)
	var configsRes []string
	err = db.Select("id").Limit(limit).Find(&configsRes).Error
	return err, configsRes
}

func CountProBlock(proCode string) (err error, configs int64) {
	var total int64 = 0
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	// db = db.Where("created_at BETWEEN ? AND ?  ", StartTime, EndTime)
	err = db.Count(&total).Error
	fmt.Println("--------------------total-----err------------------", total, err)
	return err, total
}

func FetchLogsYesterday() (list []model2.UpdateConstLog, err error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err = global.GDb.Model(model2.UpdateConstLog{}).
		Order("created_at").
		Where("to_char(created_at, 'YYYY-MM-DD') = ? and is_updated = false", yesterday).
		Find(&list).Error
	return list, err
}
