package service

import (
	"server/global"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"time"
)

// GetIsUploadUnNoticeBills 获取需要客户未接收成功已回传的单
func GetIsUploadUnNoticeBills(proCode string, startTime, endTime time.Time, billNums []string) (err error, bills []model.ProjectBill) {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr, bills
	}
	err = db.Model(&model.ProjectBill{}).Select("id,bill_name,bill_num,batch_num,stage,sale_channel,last_upload_at").
		Where("stage = 5 and last_upload_at BETWEEN ? and ? and batch_num in (?)", startTime, endTime, billNums).
		Order("created_at desc").Limit(500).Find(&bills).Error

	return err, bills
}

// FetchFtpMonitorConf 根据proCode获取ftp监控配置
func FetchFtpMonitorConf(proCode string) (err error, info model2.SysFtpMonitor) {
	err = global.GDb.Model(&model2.SysFtpMonitor{}).
		Where("pro_code = ?", proCode).
		First(&info).Error
	return err, info
}
