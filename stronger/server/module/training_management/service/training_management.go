package service

import (
	"server/global"
	"server/module/training_management/model"
	"server/module/training_management/model/request"
)

// PageTrainingManagement 分页查询培训设置
func PageTrainingManagement(reqs request.ReqTraining) (err error, list []map[string]interface{}, total int64) {
	limit := reqs.PageSize
	offset := reqs.PageSize * (reqs.PageIndex - 1)
	db := global.GDb.Model(&model.TrainingManagement{})
	if reqs.ProjectCode != "" {
		db.Where("project_code = ?", reqs.ProjectCode)
	}
	if reqs.IsAt != "" {
		db.Where("entry_start_at = ?", reqs.IsAt)
	}
	if reqs.UserName != "" {
		db.Where("user_name like ?", "%"+reqs.UserName+"%")
	}
	if reqs.UserCode != "" {
		db.Where("user_code like ?", "%"+reqs.UserCode+"%")
	}
	if reqs.AuditStatus != "" {
		db.Where("audit_status = ?", reqs.AuditStatus)
	}

	err = db.Count(&total).Limit(limit).Offset(offset).Order("created_at desc").Scan(&list).Error
	return err, list, total
}

// InfoTrainingManagement 根据ID查询
func InfoTrainingManagement(id string) (err error, info model.TrainingManagement) {
	err = global.GDb.Model(&model.TrainingManagement{}).Where("id = ?", id).Find(&info).Error
	if err != nil {
		return err, info
	}
	return err, info
}

// EditTrainingManagement 审核
func EditTrainingManagement(reqId request.ReqAuditStatus) (err error) {
	err = global.GDb.Model(&model.TrainingManagement{}).Where("id in (?)", reqId.Ids).Updates(map[string]interface{}{
		"audit_status": reqId.AuditStatus,
	}).Error
	return err
}

// ExportTrainingManagementInfo 查询要导出数据
func ExportTrainingManagementInfo(ids []string) (list []model.TrainingManagement, err error) {
	if len(ids) == 0 {
		err = global.GDb.Model(&model.TrainingManagement{}).Find(&list).Error
		return list, err
	}
	err = global.GDb.Model(&model.TrainingManagement{}).Where("id in ?", ids).Find(&list).Error
	return list, err

}
