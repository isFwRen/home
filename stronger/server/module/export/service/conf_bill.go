/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/24 3:27 下午
 */

package service

import (
	"server/global"
	"server/module/export/utils"
	"server/module/pro_conf/model"
	"sync"

	"gorm.io/gorm"
)

//GetFieldConf 获取字段配置
func GetFieldConf(proCode string, wg *sync.WaitGroup, sysProFieldList *[]model.SysProField) {
	defer utils.DealErr(wg)

	proId := global.ProCodeId[proCode]
	global.GLog.Info("项目code为：：：" + proCode)
	_ = global.GDb.Where("pro_id = ?", proId).Preload("SysIssues", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).Find(&sysProFieldList).Error
	//return err, sysProFieldList
}

func GetProFieldCheckConf(proCode string) (error, []model.SysProField) {
	sysProFieldList := []model.SysProField{}
	proId := global.ProCodeId[proCode]
	global.GLog.Info("项目code为：：：" + proCode)
	err := global.GDb.Where("pro_id = ?", proId).Preload("SysProFieldChecks", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).Find(&sysProFieldList).Error
	return err, sysProFieldList
	//return err, sysProFieldList
}

//GetExportConf 获取导出配置
func GetExportConf(proCode string) (err error, export model.SysExport) {
	proId := global.ProCodeId[proCode]
	global.GLog.Info("项目code为：：：" + proCode)
	err = global.GDb.Where("pro_id = ?", proId).First(&export).Error
	return err, export
}

//GetInspectConf 获取审核配置
func GetInspectConf(proCode string) (err error, i []model.SysInspection) {
	proId := global.ProCodeId[proCode]
	global.GLog.Info("项目code为：：：" + proCode)
	err = global.GDb.Where("pro_id = ?", proId).Find(&i).Error
	return err, i
}
