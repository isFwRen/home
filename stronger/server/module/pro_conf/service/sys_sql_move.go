package service

import (
	"fmt"
	"server/global"
	"server/module/pro_conf/model"

	"gorm.io/gorm"
)

func MoveData(proCode, mtype, TemplateId string) (err error) {
	// proId := global.ProCodeId[proCode]
	// err = global.GDb.Model(&model.SysIssue{}).Where("f_id = ? ", fId).
	// 	Order("created_at").Limit(1000).Find(&list).Error
	switch mtype {
	case "field":
		return MoveField(proCode)
	case "export":
		return MoveExport(proCode)
	case "template":
		return MoveTemplate(proCode, TemplateId)
	}
	return err
}

func MoveField(proCode string) (err error) {
	proId := global.ProCodeId[proCode]
	fmt.Println("----------MoveField---------", proCode, proId)
	var sysProFieldList []model.SysProField
	db := global.GDb.Model(&model.SysProField{}).Where("pro_id = ?", proId).
		Preload("SysIssues").
		Preload("SysProFieldChecks")
	err = db.Find(&sysProFieldList).Error
	fmt.Println("----------sysProFieldList---------", len(sysProFieldList))
	return global.GDbMove.Transaction(func(tx *gorm.DB) error {
		ids := []string{}
		if err := tx.Model(&model.SysProField{}).Select("id").Where("pro_id = ?", proId).Find(&ids).Error; err != nil {
			return err
		}
		if len(ids) > 0 {
			if err = tx.Where("f_id in ?", ids).Delete(&model.SysIssue{}).Error; err != nil {
				return err
			}
			if err = tx.Where("f_id in ?", ids).Delete(&model.SysProFieldCheck{}).Error; err != nil {
				return err
			}
			if err := tx.Where("pro_id = ?", proId).Delete(&model.SysProField{}).
				Preload("SysIssues").
				Preload("SysProFieldChecks").Error; err != nil {
				return err
			}
		}

		for _, sysProFieldList := range sysProFieldList {
			if err = tx.Model(&model.SysProField{}).Create(&sysProFieldList).
				Preload("SysIssues").
				Preload("SysProFieldChecks").Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func MoveExport(proCode string) (err error) {
	proId := global.ProCodeId[proCode]
	fmt.Println("----------MoveExport---------", proCode, proId)
	var SysExportList []model.SysExport
	db := global.GDb.Model(&model.SysExport{}).Where("pro_id = ?", proId).
		Preload("SysExportNodes")
	err = db.Find(&SysExportList).Error
	fmt.Println("----------sysProFieldList---------", len(SysExportList))
	fmt.Println("----------SysExportNodes---------", len(SysExportList[0].SysExportNodes))

	return global.GDbMove.Transaction(func(tx *gorm.DB) error {
		// aaa := model.SysExport{}
		// aaa.ProId = proId  .Select("SysExportNodes") .Where("pro_id = ?", proId)
		ids := []string{}
		if err := tx.Model(&model.SysExport{}).Select("id").Where("pro_id = ?", proId).Find(&ids).Error; err != nil {
			return err
		}
		fmt.Println("----------ids-----------", ids)
		if len(ids) > 0 {
			if err = tx.Where("export_id in ?", ids).Delete(&model.SysExportNode{}).Error; err != nil {
				return err
			}
			if err = tx.Where("pro_id = ?", proId).Delete(&model.SysExport{}).Error; err != nil {
				return err
			}
		}
		// subQuery := tx.Select("id").Where("pro_id = ?", "proId%").Table("users")

		fmt.Println("----------err---------", err)
		for _, sysExport := range SysExportList {
			fmt.Println("----------sysExport---------", len(sysExport.SysExportNodes))
			if err = tx.Model(&model.SysExport{}).Create(&sysExport).
				Preload("SysExportNodes").Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func MoveTemplate(proCode, TemplateId string) (err error) {
	proId := global.ProCodeId[proCode]
	fmt.Println("----------MoveTemplate---------", proCode, proId)
	// var SysExportList []model.SysExport
	var reSysProTemplate model.SysProTemplate
	err = global.GDb.Where("name = ? and pro_id = ?", TemplateId, proId).
		Preload("SysProTempBs").
		Preload("SysProTempBs.TempBFRelations").
		Preload("SysProTempBs.TempBlockRelations").
		First(&reSysProTemplate).Error
	// db := global.GDb.Model(&model.SysExport{}).Where("pro_id = ?", proId).
	// 	Preload("SysExportNodes")
	// err = db.Find(&SysExportList).Error
	// fmt.Println("----------sysProFieldList---------", len(SysExportList))

	return global.GDbMove.Transaction(func(tx *gorm.DB) error {
		ids := []string{}
		if err := tx.Model(&model.SysProTemplate{}).Select("id").Where("name = ? and pro_id = ?", TemplateId, proId).Find(&ids).Error; err != nil {
			return err
		}
		fmt.Println("----------ids-----------", ids)
		if len(ids) > 0 {
			bids := []string{}
			if err := tx.Model(&model.SysProTempB{}).Select("id").Where("pro_temp_id in ?", ids).Find(&bids).Error; err != nil {
				return err
			}
			if len(bids) > 0 {
				if err = tx.Where("temp_b_id in ?", bids).Delete(&model.TempBFRelation{}).Error; err != nil {
					return err
				}
				if err = tx.Where("temp_b_id in ?", bids).Delete(&model.TempBlockRelation{}).Error; err != nil {
					return err
				}
			}
			if err = tx.Where("pro_temp_id in ?", ids).Delete(&model.SysProTempB{}).Error; err != nil {
				return err
			}
			if err := tx.Where("name = ? and pro_id = ?", TemplateId, proId).Delete(&model.SysProTemplate{}).
				Preload("SysProTempBs").Error; err != nil {
				return err
			}
		}

		// for _, sysExport := range SysExportList {
		if err = tx.Model(&model.SysProTemplate{}).Create(&reSysProTemplate).Error; err != nil {
			return err
		}
		// }
		return nil
	})
}
