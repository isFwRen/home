package service

import (
	"server/global"
	"server/module/export/model"

	"gorm.io/gorm"
)

func InsertDeductionDetails(proCode, billId string, items []model.DeductionDetail) error {

	db := global.ProDbMap[proCode]
	if db == nil {
		panic("没有该项目的连接")
	}

	return global.ProDbMap[proCode].Transaction(func(tx *gorm.DB) error {
		err := tx.Where("bill_id = ?", billId).Delete(&model.DeductionDetail{}).Error
		if err != nil {
			return err
		}
		if len(items) > 0 {
			return tx.Save(&items).Error
		}
		return err
	})
	// err := db.Order("created_at desc").Where("stage = 6").Find(&obj).Error

}
