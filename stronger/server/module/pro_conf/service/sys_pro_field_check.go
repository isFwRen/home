package service

import (
	"server/global"
	"server/module/pro_conf/model"

	"go.uber.org/zap"
)

//GetFieldCheckListByFieldId 根据字段id获取问题件配置
func GetFieldCheckListByFieldId(fId string) (err error, list []model.SysProFieldCheck) {
	err = global.GDb.Model(&model.SysProFieldCheck{}).Where("f_id = ? ", fId).
		Order("created_at").Limit(1000).Find(&list).Error
	return err, list
}

//EditIssuesByFieldId 根据字段id编辑问题件
func EditFieldCheckByFieldId(fId string, list []model.SysProFieldCheck) error {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("插入历史库失败")
		}
	}()
	//删除
	dRow := tx.Model(&model.SysProFieldCheck{}).Where("f_id = ?", fId).Delete(model.SysProFieldCheck{}).RowsAffected
	global.GLog.Info("删除数量数::", zap.Int64("dRow", dRow))

	//新增
	aRow := tx.Model(&model.SysProFieldCheck{}).Create(list).RowsAffected
	global.GLog.Info("新增数量数::", zap.Int64("aRow", aRow))
	return tx.Commit().Error
}
