/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/7 13:49
 */

package service

import (
	"gorm.io/gorm"
	"server/global"
	"server/module/pro_manager/model"
	request2 "server/module/pro_manager/model/request"
	"time"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

// AddAnnouncement 新增公告
func AddAnnouncement(announcement model.Announcement) error {
	return global.GDb.Create(&announcement).Error
}

// GetAnnouncementByPage 获取公告管理分页
func GetAnnouncementByPage(announcementPageReq request2.AnnouncementPageReq) (err error, total int64, list []model.Announcement) {
	limit := announcementPageReq.PageSize
	offset := announcementPageReq.PageSize * (announcementPageReq.PageIndex - 1)
	db := global.GDb.Model(&model.Announcement{})
	if announcementPageReq.ProCode != "" {
		db = db.Where("pro_code = ? ", announcementPageReq.ProCode)
	}
	if announcementPageReq.Status > 0 {
		db = db.Where("status = ? ", announcementPageReq.Status)
	}
	if announcementPageReq.ReleaseType > 0 {
		db = db.Where("release_type = ? ", announcementPageReq.ReleaseType)
	}
	db = db.Where("created_at BETWEEN ? AND ?  ", announcementPageReq.StartTime, announcementPageReq.EndTime)
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("status asc,created_at desc").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

// ChangeStatusAnnouncementById 更改状态
func ChangeStatusAnnouncementById(announcementChangeStatusReq request2.AnnouncementChangeStatusReq, customClaims *model3.CustomClaims) int64 {
	obj := map[string]interface{}{
		"status": announcementChangeStatusReq.Status,
	}
	if announcementChangeStatusReq.Status == 2 {
		obj["release_date"] = time.Now()
		obj["release_user_code"] = customClaims.Code
		obj["release_user_name"] = customClaims.NickName
	}
	return global.GDb.Where("id in (?)", announcementChangeStatusReq.Ids).Updates(&model.Announcement{}).Updates(obj).RowsAffected
}

// UpdateAnnouncementById 更新公告信息
func UpdateAnnouncementById(announcement model.Announcement) error {
	return global.GDb.Where("id = ?", announcement.ID).Updates(&announcement).Updates(map[string]interface{}{
		"title":        announcement.Title,
		"release_type": announcement.ReleaseType,
		"pro_code":     announcement.ProCode,
		"content":      announcement.Content,
	}).Error
}

// GetHomePageAnnouncement 主页获取公告或动态
func GetHomePageAnnouncement(announcementPageHomeReq request2.AnnouncementPageHomeReq) (err error, total int64, list []model.Announcement) {
	limit := announcementPageHomeReq.PageSize
	offset := announcementPageHomeReq.PageSize * (announcementPageHomeReq.PageIndex - 1)
	db := global.GDb.Model(&model.Announcement{})
	db = db.Where("title like ? and status = 2 and release_type = ? ",
		"%"+announcementPageHomeReq.Title+"%",
		announcementPageHomeReq.ReleaseType)
	if !announcementPageHomeReq.StartTime.IsZero() && !announcementPageHomeReq.EndTime.IsZero() {
		db = db.Where("release_date BETWEEN ? AND ? ",
			announcementPageHomeReq.StartTime,
			announcementPageHomeReq.EndTime)
	}
	if announcementPageHomeReq.ProCode != "" {
		db.Where("pro_code = ? ", announcementPageHomeReq.ProCode)
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("created_at asc").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

// AddView 新增访问量
func AddView(id string) error {
	return global.GDb.Model(&model.Announcement{}).Where("id = ?", id).
		Update("visit_count", gorm.Expr("visit_count + ?", 1)).Error
}
