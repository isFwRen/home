package service

import (
	"gorm.io/gorm"
	"server/global"
	model2 "server/module/load/model"
	"server/module/pro_manager/model"
	"time"
)

func DelBill(proCode string, id string) (err error) {
	var configAging model.ProjectBill
	err = global.ProDbMap[proCode].Where("id = ?", id).Delete(&configAging).Error
	return err
}

func UpdateBill(proCode string, configAging model.ProjectBill, id string) (err error) {
	err = global.ProDbMap[proCode].Where("id = ?", id).First(&model.ProjectBill{}).Save(&configAging).Error
	return err
}

func InsertBill(proCode string, agingConfig model.ProjectBill) (err error, configInter model.ProjectBill) {
	err = global.ProDbMap[proCode].Model(&model.ProjectBill{}).Create(&agingConfig).Error
	return err, agingConfig
}

func SelectBillByID(proCode string, id string) (error, model.ProjectBill) {
	var configsRes model.ProjectBill
	err := global.ProDbMap[proCode].Where("id = ?", id).First(&configsRes).Error
	return err, configsRes
}

func SelectBillsByStage(proCode string, stage int) (err error, configs []model.ProjectBill) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBill{})
	if stage != 0 {
		db = db.Where("stage = ?", stage)
	}
	var configsRes []model.ProjectBill
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func CountBill(proCode string, name string) (err error, total int64) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBill{})
	if name != "" {
		db = db.Where("bill_name = ?", name)
	}
	err = db.Count(&total).Error
	return err, total
}

//UpdateClaimType 更新理赔类型
func UpdateClaimType(proCode string, billId string, claimType int) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	return db.Model(&model.ProjectBill{}).Where("id = ?", billId).Update("claim_type", claimType).Error
}

//UpdateClaimTypeStickLevel 更新理赔类型
func UpdateClaimTypeStickLevel(proCode string, billId string, claimType, stickLevel int) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	return db.Model(&model.ProjectBill{}).Where("id = ?", billId).Updates(map[string]interface{}{
		"claim_type":  claimType,
		"stick_level": stickLevel,
	}).Error
}

//DelBlocksFieldsAndUpdateBillStageByBillID 删除分块字段并更新单据状态
func DelBlocksFieldsAndUpdateBillStageByBillID(dbStr string, projectBill model.ProjectBill) (err error) {
	db := global.ProDbMap[dbStr]
	if db == nil {
		return global.ProDbErr
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("bill_id = ?", projectBill.ID).Delete(&model2.ProjectField{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("bill_id = ?", projectBill.ID).Delete(&model2.ProjectBlock{}).Error
		if err != nil {
			return err
		}
		err = tx.Model(&model.ProjectBill{}).Where("id = ?", projectBill.ID).Updates(map[string]interface{}{
			"stage":       projectBill.Stage,
			"status":      projectBill.Status,
			"updated_at":  time.Now(),
			"pictures":    projectBill.Pictures,
			"images_type": projectBill.ImagesType,
		}).Error
		return err
	})
	return err
}

//DelBlocksFieldsAndInsertBillDetail 删除分块和字段并插入新的单的信息
func DelBlocksFieldsAndInsertBillDetail(dbStr string, proFields []model2.ProjectField, projectBlock model2.ProjectBlock, projectBill model.ProjectBill) (err error) {
	db := global.ProDbMap[dbStr]
	if db == nil {
		return global.ProDbErr
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("bill_id = ?", projectBill.ID).Delete(&model2.ProjectField{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("bill_id = ?", projectBill.ID).Delete(&model2.ProjectBlock{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("id = ?", projectBill.ID).Delete(&model.ProjectBill{}).Error
		if err != nil {
			return err
		}

		err = tx.Model(model2.ProjectField{}).Create(&proFields).Error
		if err != nil {
			return err
		}
		err = tx.Model(model2.ProjectBlock{}).Create(&projectBlock).Error
		if err != nil {
			return err
		}
		err = tx.Model(model.ProjectBill{}).Create(&projectBill).Error
		if err != nil {
			return err
		}
		return err
	})
	return err
}
