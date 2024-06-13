package service

import (
	"server/global"
	"server/module/pro_conf/model"
	sysModel "server/module/sys_base/model"
)

func GetSysProTempByProIdAndName(proId string, name string) (error, model.SysProTemplate) {
	var sysProTemplate model.SysProTemplate
	err := global.GDb.Model(&model.SysProTemplate{}).Where("pro_id = ? AND name = ?", proId, name).First(&sysProTemplate).Error
	//err = global.G_DB.Where("id = ?", id).Preload("SysProTempB").First(&reSysProTemplate).Error
	return err, sysProTemplate
}

func GetSysProject(code string) (error, model.SysProject) {
	var sysProject model.SysProject
	err := global.GDb.Model(&model.SysProject{}).Where("code = ?", code).First(&sysProject).Error
	return err, sysProject
}

func GetSysWhiteLists(proCode string) (error, []sysModel.SysWhiteList) {
	var sysWhiteList []sysModel.SysWhiteList
	err := global.ProDbMap[proCode].Model(&sysModel.SysWhiteList{}).Find(&sysWhiteList).Error
	return err, sysWhiteList
}

func GetSysProTempBlockByTempId(tempId string) (error, []model.SysProTempB) {
	db := global.GDb.Model(&model.SysProTempB{}).Where("pro_temp_id = ?", tempId)
	var sysProTempBList []model.SysProTempB
	err := db.Find(&sysProTempBList).Error

	return err, sysProTempBList
}

func GetSysFields(ProId string) (error, []model.SysProField) {
	// 创建db
	db := global.GDb.Model(&model.SysProField{}).Where("pro_id = ?", ProId)
	var sysProFieldList []model.SysProField

	err := db.Find(&sysProFieldList).Error
	return err, sysProFieldList
}

func GetTempBFRelationByBId(bId string) (err error, tempBFRelationList []model.TempBFRelation) {
	err = global.GDb.Model(&model.TempBFRelation{}).
		Where("temp_b_id = ? ", bId).
		Order("updated_at asc").
		Find(&tempBFRelationList).
		Error
	return err, tempBFRelationList
}

func GetBlockRelationsByBId(bId string) (err error, tempBlockRelationList []model.TempBlockRelation) {

	var tempBlockRelation []model.TempBlockRelation
	err = global.GDb.Model(&model.TempBlockRelation{}).Where("temp_b_id = ?", bId).Find(&tempBlockRelation).Error
	return err, tempBlockRelation
}
