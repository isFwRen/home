/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/18 10:56
 */

package service

import (
	"server/global"
	"server/module/msg_manager/model"
	"time"
)

//GetCustomerNoticeByPage 获取客户通知列表
func GetCustomerNoticeByPage(search model.CustomerNoticeSearchReq) (err error, total int64, list []model.CustomerNotice) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, 0, list
	}
	limit := search.PageSize
	offset := search.PageSize * (search.PageIndex - 1)
	db = db.Model(&model.CustomerNotice{}).Where("send_time BETWEEN ? AND ? ", search.StartTime, search.EndTime)
	if search.MsgType > 0 {
		db = db.Where("msg_type = ? ", search.MsgType)
	}
	if search.Status > 0 {
		db = db.Where("status = ? ", search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

//GetCustomerNoticeById 获取客户通知
func GetCustomerNoticeById(search model.CustomerNotice) (err error, item model.CustomerNotice) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, item
	}
	err = db.Model(&model.CustomerNotice{}).Where("id = ? ", search.ID).First(&item).Error
	return err, item
}

//Reply 回复客户
func Reply(customerNotice model.CustomerNotice) (err error) {
	db := global.ProDbMap[customerNotice.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	err = db.Model(&model.CustomerNotice{}).Where("id = ?", customerNotice.ID).Updates(
		map[string]interface{}{
			"is_reply":       customerNotice.IsReply,
			"expect_num":     customerNotice.ExpectNum,
			"reply_time":     time.Now(),
			"status":         2,
			"deal_user_code": customerNotice.DealUserCode,
			"deal_user_name": customerNotice.DealUserName,
		}).Error
	return err
}

//GetCustomerNoticeByStatus 统计有多少条通知函2
func GetCustomerNoticeByStatus(proCode string) (err error, total int64) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, 0
	}
	err = db.Model(&model.CustomerNotice{}).Where("status = 1").Count(&total).Error
	return err, total
}
