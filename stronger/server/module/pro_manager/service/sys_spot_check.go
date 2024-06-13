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

func UpdateSysSpotCheckById(configAging model.SysSpotCheck, id string) (err error) {
	err = global.GDb.Where("id = ?", id).First(&model.SysSpotCheck{}).Save(&configAging).Error
	return err
}

func InsertSysSpotCheck(agingConfig model.SysSpotCheck) (err error) {
	err = global.GDb.Model(&model.SysSpotCheck{}).Create(&agingConfig).Error
	return err
}

func SelectSysSpotCheckByID(id string) (error, model.SysSpotCheck) {
	var configsRes model.SysSpotCheck
	err := global.GDb.Where("id = ?", id).First(&configsRes).Error
	return err, configsRes
}

func GetSysSpotCheckByPage(announcementPageReq request2.SysSpotCheckQuery) (err error, total int64, list []model.SysSpotCheck) {
	limit := announcementPageReq.PageSize
	offset := announcementPageReq.PageSize * (announcementPageReq.PageIndex - 1)
	db := global.GDb.Model(&model.SysSpotCheck{})
	if announcementPageReq.ProCode != "" {
		db = db.Where("pro_code = ? ", announcementPageReq.ProCode)
	}
	if announcementPageReq.Status > 0 {
		db = db.Where("status = ? ", announcementPageReq.Status)
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

func GetSysSpotCheckByType(typeValue int) (err error, list []model.SysSpotCheck) {
	db := global.GDb.Model(&model.SysSpotCheck{})
	db = db.Where("type = ? ", typeValue)
	err = db.Find(&list).Error
	return err, list
}
