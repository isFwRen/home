/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/7 13:49
 */

package service

import (
	"fmt"
	"server/global"
	loadModel "server/module/load/model"
	"server/module/pro_manager/model"
	request2 "server/module/pro_manager/model/request"
	"time"

	"github.com/shopspring/decimal"
	//  "time"
)

func GetSysSpotCheckStatisticByPage(announcementPageReq request2.SysSpotCheckStatisticQuery) (err error, total int64, list []model.SysSpotCheckStatistic) {
	limit := announcementPageReq.PageSize
	offset := announcementPageReq.PageSize * (announcementPageReq.PageIndex - 1)
	db := global.GDb.Model(&model.SysSpotCheckStatistic{})
	if announcementPageReq.ProCode != "" {
		db = db.Where("pro_code = ? ", announcementPageReq.ProCode)
	}
	if announcementPageReq.Type > 0 {
		db = db.Where("type = ? ", announcementPageReq.Type)
	}
	db = db.Where("created_at BETWEEN ? AND ?  ", announcementPageReq.StartTime, announcementPageReq.EndTime)
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("created_at asc").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

func GetStaticByType(announcementPageReq request2.SysSpotCheckStatisticQuery) (err error, total int64, list []model.SysSpotCheckStatistic) {
	err, sum := CountProBlock(announcementPageReq.ProCode, announcementPageReq.StartTime, announcementPageReq.EndTime)
	if err != nil {
		return
	}
	checkList := []model.SysSpotCheckData{}
	if announcementPageReq.Type == 1 {
		err, total, checkList = CountCheckData1(announcementPageReq)
	} else {
		err, total, checkList = CountCheckData2(announcementPageReq)
	}

	for _, checkData := range checkList {
		data := model.SysSpotCheckStatistic{}
		data.Sum = sum
		data.ProCode = announcementPageReq.ProCode
		data.CheckSum = data.CheckSum + checkData.Num
		data.WrongNum = data.CheckSum + checkData.WrongNum
		if announcementPageReq.Type == 2 {
			data.SubmitDay = checkData.SubmitDay
		}
		if data.CheckSum != 0 {
			data.WrongRatio = decimal.NewFromInt(int64(data.WrongNum)).Div(decimal.NewFromInt(int64(data.CheckSum)))
		}
		if data.Sum != 0 {
			data.WrongRatio = decimal.NewFromInt(int64(data.CheckSum)).Div(decimal.NewFromInt(int64(data.Sum)))
		}
		list = append(list, data)
	}
	return err, total, list
}

func CountCheckData1(announcementPageReq request2.SysSpotCheckStatisticQuery) (err error, total int64, list []model.SysSpotCheckData) {
	db := global.GDb.Model(&model.SysSpotCheckData{})
	db = db.Where("type = ?", 2)
	db = db.Where("created_at BETWEEN ? AND ?  ", announcementPageReq.StartTime, announcementPageReq.EndTime)
	err = db.Select("pro_code, sum(num) as num, sum(done_num) as done_num, sum(undone_num) as undone_num, sum(wrong_num) as wrong_num").Group("pro_code").Having("pro_code = ?", announcementPageReq.ProCode).Find(&list).Error
	fmt.Println("--------------------total-----list------------------", err, list)
	return err, total, list
}

func CountCheckData2(announcementPageReq request2.SysSpotCheckStatisticQuery) (err error, total int64, list []model.SysSpotCheckData) {
	db := global.GDb.Model(&model.SysSpotCheckData{})
	db = db.Where("type = ?", 2)
	db = db.Where("created_at BETWEEN ? AND ?  ", announcementPageReq.StartTime, announcementPageReq.EndTime)
	err = db.Select("submit_day, pro_code, sum(num) as num, sum(done_num) as done_num, sum(undone_num) as undone_num, sum(wrong_num) as wrong_num").Group("pro_code, submit_day").Having("pro_code = ?", announcementPageReq.ProCode).Find(&list).Error
	fmt.Println("--------------------total-----list------------------", err, list)
	return err, total, list
}

func CountProBlock(proCode string, startTime time.Time, endTime time.Time) (err error, configs int64) {
	var total int64 = 0
	db := global.ProDbMap[proCode].Model(&loadModel.ProjectBlock{})
	db = db.Where("created_at BETWEEN ? AND ?  ", startTime, endTime)
	err = db.Count(&total).Error
	fmt.Println("--------------------total-----err------------------", total, err)
	return err, total
}
