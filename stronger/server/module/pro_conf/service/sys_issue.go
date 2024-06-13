/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/2/1 13:46
 */

package service

import (
	"go.uber.org/zap"
	"server/global"
	"server/module/pro_conf/model"
)

//GetIssueListByFieldId 根据字段id获取问题件配置
func GetIssueListByFieldId(fId string) (err error, list []model.SysIssue) {
	err = global.GDb.Model(&model.SysIssue{}).Where("f_id = ? ", fId).
		Order("created_at").Limit(1000).Find(&list).Error
	return err, list
}

//EditIssuesByFieldId 根据字段id编辑问题件
func EditIssuesByFieldId(fId string, list []model.SysIssue) error {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("插入历史库失败")
		}
	}()
	//删除
	dRow := tx.Model(&model.SysIssue{}).Where("f_id = ?", fId).Delete(model.SysIssue{}).RowsAffected
	global.GLog.Info("删除数量数::", zap.Int64("dRow", dRow))

	//新增
	aRow := tx.Model(&model.SysIssue{}).Create(list).RowsAffected
	global.GLog.Info("新增数量数::", zap.Int64("aRow", aRow))
	return tx.Commit().Error
}
