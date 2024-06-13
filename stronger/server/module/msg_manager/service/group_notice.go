/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/25 14:53
 */

package service

import (
	"errors"
	"gorm.io/gorm/clause"
	"server/global"
	"server/module/msg_manager/model"
	model2 "server/module/msg_manager/model/request"
)

// AddGroupNotices 新增一组固定通知
func AddGroupNotices(groupNoticeAddReq model2.GroupNoticeAddReq) (err error) {
	if groupNoticeAddReq.Type == 1 {
		if len(groupNoticeAddReq.Ones) < 1 {
			return errors.New("新增的数据为空")
		}
		var maxBlock int
		global.GDb.Raw("SELECT max(block) as value FROM group_notice_ones WHERE group_id = ? ", groupNoticeAddReq.Ones[0].GroupId).
			Scan(&maxBlock)
		var dingtalkGroup model.DingtalkGroup
		global.GDb.Model(&model.DingtalkGroup{}).Where("id = ?", groupNoticeAddReq.Ones[0].GroupId).First(&dingtalkGroup)
		for i, _ := range groupNoticeAddReq.Ones {
			groupNoticeAddReq.Ones[i].Block = maxBlock + 1
			groupNoticeAddReq.Ones[i].ProCode = dingtalkGroup.ProCode
		}
		err = global.GDb.Model(&model.GroupNoticeOne{}).Create(&groupNoticeAddReq.Ones).Error
	} else {
		if len(groupNoticeAddReq.Twos) < 1 {
			return errors.New("新增的数据为空")
		}
		var dingtalkGroup model.DingtalkGroup
		global.GDb.Model(&model.DingtalkGroup{}).Where("id = ?", groupNoticeAddReq.Twos[0].GroupId).First(&dingtalkGroup)
		var maxBlock int
		global.GDb.Raw("SELECT max(block) as value FROM group_notice_twos WHERE group_id = ? ", groupNoticeAddReq.Twos[0].GroupId).
			Scan(&maxBlock)
		for i, _ := range groupNoticeAddReq.Twos {
			groupNoticeAddReq.Twos[i].Block = maxBlock + 1
			groupNoticeAddReq.Twos[i].ProCode = dingtalkGroup.ProCode
		}
		err = global.GDb.Model(&model.GroupNoticeTwo{}).Create(&groupNoticeAddReq.Twos).Error
	}

	return err
}

// GetGroupNoticeByGroupId 获取项目的模板
func GetGroupNoticeByGroupId(groupNoticeReq model2.GroupNoticeReq) (err error, ones []model.GroupNoticeOne, twos []model.GroupNoticeTwo) {
	err = global.GDb.Model(&model.GroupNoticeOne{}).Where("group_id = ?", groupNoticeReq.GroupId).
		Find(&ones).Order("created_at").Error
	if err != nil {
		return err, ones, twos
	}
	err = global.GDb.Model(&model.GroupNoticeTwo{}).Where("group_id = ?", groupNoticeReq.GroupId).
		Find(&twos).Order("created_at").Error
	return err, ones, twos
}

// ReGroupNotice 重置
func ReGroupNotice(groupNoticeReq model2.GroupNoticeReq) (rows int64) {
	if groupNoticeReq.Type == 1 {
		var groupNoticeOne model.GroupNoticeOne
		rows = global.GDb.Model(&model.GroupNoticeOne{}).Where("group_id = ?", groupNoticeReq.GroupId).
			Delete(&groupNoticeOne).RowsAffected
	} else {
		var groupNoticeTwo model.GroupNoticeTwo
		rows = global.GDb.Model(&model.GroupNoticeTwo{}).Where("group_id = ?", groupNoticeReq.GroupId).
			Delete(&groupNoticeTwo).RowsAffected
	}
	return rows
}

// EditGroupNotices 编辑
func EditGroupNotices(groupNoticeAddReq model2.GroupNoticeAddReq) (err error) {
	//for i, _ := range ones {
	//	if ones[i].ID == "" {
	//		ones[i].ID ==
	//	}
	//}
	if groupNoticeAddReq.Type == 1 {
		err = global.GDb.Debug().Model(&model.GroupNoticeOne{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
			//DoUpdates: clause.AssignmentColumns([]string{"start_time", "end_time",
			//	"interval", "day_of_week", "block", "group_id"}),
		}).Create(&groupNoticeAddReq.Ones).Error
	} else {
		err = global.GDb.Debug().Model(&model.GroupNoticeTwo{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
			//DoUpdates: clause.AssignmentColumns([]string{"send_time"}),
		}).Create(&groupNoticeAddReq.Twos).Error
	}
	return err
}

// GetGroupNotice 获取定时发送
func GetGroupNotice() (err error, twos []model.GroupNoticeTwo, ones []model.GroupNoticeOne) {
	err = global.GDb.Model(&model.GroupNoticeTwo{}).Limit(1000).
		Find(&twos).Error
	if err != nil {
		return err, twos, ones
	}
	err = global.GDb.Model(&model.GroupNoticeOne{}).Limit(1000).
		Find(&ones).Error
	return err, twos, ones
}
