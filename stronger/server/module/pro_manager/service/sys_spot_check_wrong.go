/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/7 13:49
 */

package service

import (
	"server/global"
	"server/module/pro_manager/model"

	//  "time"
	request2 "server/module/pro_manager/model/request"
)

func UpdateSysSpotCheckWrong(configAging model.SysSpotCheckWrong, id string) (err error) {
	err = global.GDb.Where("id = ?", id).First(&model.SysSpotCheckWrong{}).Save(&configAging).Error
	return err
}

func InsertSysSpotCheckWrong(agingConfig model.SysSpotCheckWrong) (err error, configInter model.SysSpotCheckWrong) {
	err = global.GDb.Model(&model.SysSpotCheckWrong{}).Create(&agingConfig).Error
	return err, agingConfig
}

func SelectSysSpotCheckWrongByID(id string) (error, model.SysSpotCheckWrong) {
	var configsRes model.SysSpotCheckWrong
	err := global.GDb.Where("id = ?", id).First(&configsRes).Error
	return err, configsRes
}

func GetSysSpotCheckWrongByPage(announcementPageReq request2.SysSpotCheckWrongQuery) (err error, total int64, list []model.SysSpotCheckWrong) {
	limit := announcementPageReq.PageSize
	offset := announcementPageReq.PageSize * (announcementPageReq.PageIndex - 1)
	db := global.GDb.Model(&model.SysSpotCheckWrong{})
	if announcementPageReq.ProCode != "" {
		db = db.Where("pro_code = ? ", announcementPageReq.ProCode)
	}
	if announcementPageReq.CreatedCode != "" {
		db = db.Where("created_code = ? ", announcementPageReq.CreatedCode)
	}
	if announcementPageReq.CreatedName != "" {
		db = db.Where("created_name = ? ", announcementPageReq.CreatedName)
	}
	if announcementPageReq.FieldName != "" {
		db = db.Where("field_name = ? ", announcementPageReq.FieldName)
	}
	if announcementPageReq.Type != "" {
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
