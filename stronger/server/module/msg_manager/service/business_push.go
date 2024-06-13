/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/18 10:56
 */

package service

import (
	"go.uber.org/zap"
	"server/global"
	"server/module/msg_manager/model"
)

// GetBusinessPushByPage 获取业务通知列表
func GetBusinessPushByPage(search model.BusinessPushSearchReq, uid string) (err error, total int64, list []model.BusinessPushSend) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, 0, list
	}
	limit := search.PageSize
	offset := search.PageSize * (search.PageIndex - 1)
	db = db.Model(&model.BusinessPushSend{}).Joins("BusinessPush").
		Where("pro_code = ? and business_push_sends.user_id = ? and \"BusinessPush\".\"created_at\" BETWEEN ? AND ? ",
			search.ProCode, uid, search.StartTime, search.EndTime)
	if search.MsgType > 0 {
		db = db.Where("type = ? ", search.MsgType)
	}
	err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return err, total, list
	}
	err = db.Offset(-1).Count(&total).Error
	return err, total, list
}

// Read 将推送小新标志为已读
//
// Parameters:
//
//	search (model.BusinessPushSendReadReq): id数组和项目编码
//
// Returns:
//
//	int64: 更新条数
//
// Description:
//
//	直接更新数据库business_push_send的is_read字段为true
func Read(search model.BusinessPushSendReadReq) (row int64) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		global.GLog.Error("", zap.Error(global.ProDbErr))
		return 0
	}
	return db.Model(model.BusinessPushSend{}).Where("id in (?)", search.Ids).
		Updates(map[string]interface{}{
			"is_read": true,
		}).RowsAffected
}
