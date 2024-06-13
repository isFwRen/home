package service

import (
	"gorm.io/gorm/clause"
	"server/global"
	"server/module/pro_conf/model"
)

// FetchFtpMonitorConf 根据proCode获取ftp监控配置
func FetchFtpMonitorConf(proCode string) (err error, info []model.SysFtpMonitor) {
	err = global.GDb.Model(&model.SysFtpMonitor{}).
		Where("pro_code = ?", proCode).
		Find(&info).Error
	return err, info
}

// EditFtpMonitorConf 编辑ftp监控配置
func EditFtpMonitorConf(editReq model.SysFtpMonitor) (row int64) {
	assignmentColumns := []string{"frequency", "wrong_msg", "desc", "pro_code"}
	if editReq.ID == "" {
		assignmentColumns = append(assignmentColumns, "created_code", "created_name")
	}
	return global.GDb.Model(&model.SysFtpMonitor{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(assignmentColumns),
	}).Create(&editReq).RowsAffected
}
