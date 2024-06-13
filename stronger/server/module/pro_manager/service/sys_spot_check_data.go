/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/7 13:49
 */

package service

import (
	"server/global"
	"server/module/pro_manager/model"
	request2 "server/module/pro_manager/model/request"
	//  "time"
)

func UpdateSysSpotCheckDataById(configAging model.SysSpotCheckData, id string) (err error) {
	err = global.GDb.Where("id = ?", id).First(&model.SysSpotCheckData{}).Save(&configAging).Error
	return err
}

func InsertSysSpotCheckData(agingConfig model.SysSpotCheckData) (err error) {
	err = global.GDb.Model(&model.SysSpotCheckData{}).Create(&agingConfig).Error
	return err
}

func SelectSysSpotCheckDataByID(id string) (error, model.SysSpotCheckData) {
	var configsRes model.SysSpotCheckData
	err := global.GDb.Where("id = ?", id).First(&configsRes).Error
	return err, configsRes
}

func GetSysSpotCheckDataByPage(announcementPageReq request2.SysSpotCheckDataQuery) (err error, total int64, list []model.SysSpotCheckData) {
	limit := announcementPageReq.PageSize
	offset := announcementPageReq.PageSize * (announcementPageReq.PageIndex - 1)
	db := global.GDb.Model(&model.SysSpotCheckData{})
	if announcementPageReq.ProCode != "" {
		db = db.Where("pro_code = ? ", announcementPageReq.ProCode)
	}
	if announcementPageReq.Code > 0 {
		db = db.Where("code = ? ", announcementPageReq.Code)
	}
	if announcementPageReq.Name > 0 {
		db = db.Where("name = ? ", announcementPageReq.Name)
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
