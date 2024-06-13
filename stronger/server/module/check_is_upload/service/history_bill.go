/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/20 11:11
 */

package service

import (
	"server/global"
	"server/module/pro_manager/model"
	"time"
)

// GetCheckIsUploadHistoryBills 获取需要检查是否接收成功的单
func GetCheckIsUploadHistoryBills(proCode string, startTime time.Time) (err error, bills []model.ProjectBill) {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr, bills
	}
	err = db.Model(&model.ProjectBill{}).Select("id,bill_name,bill_num,batch_num,stage,sale_channel,created_at").
		Where("stage = 5 and created_at BETWEEN ? and ?", startTime, time.Now()).
		Order("created_at desc").Limit(500).Find(&bills).Error

	return err, bills
}

// UpdateBillStage 更新单据状态
func UpdateBillStage(proCode string, bill model.ProjectBill) error {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr
	}
	return db.Model(&model.ProjectBill{}).Where("id = ?", bill.ID).Updates(map[string]interface{}{
		"stage": 7,
	}).Error
}
