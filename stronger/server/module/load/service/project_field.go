package service

import (
	"server/global"
	"server/module/load/model"
)

func InsertField(proCode string, agingConfig model.ProjectField) (err error, configInter model.ProjectField) {
	err = global.ProDbMap[proCode].Create(&agingConfig).Error
	return err, agingConfig
}

func InsertFields(proCode string, agingConfig []model.ProjectField) (err error) {
	err = global.ProDbMap[proCode].Create(&agingConfig).Error
	return err
}

func UpdateField(proCode string, configAging model.ProjectField, id string) (err error) {
	err = global.ProDbMap[proCode].Where("id = ?", id).Updates(configAging).Error
	return err
}

func DelFieldsByBillID(proCode string, bill_id string) (err error) {
	var configAging model.ProjectField
	err = global.ProDbMap[proCode].Where("bill_id = ?", bill_id).Delete(&configAging).Error
	return err
}

func SelectFieldsByBlockID(proCode string, block_id string) (err error, configs []model.ProjectField) {
	db := global.ProDbMap[proCode].Model(&model.ProjectField{})
	if block_id != "" {
		db = db.Where("block_id = ?", block_id)
	}
	db.Order("created_at asc").Order("block_index asc").Order("field_index asc")
	var configsRes []model.ProjectField
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func SelectCropFields(proCode string, block_id string, name string) (err error, configs []model.ProjectField) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})

	db = db.Where("block_id = ? AND name = ?", block_id, name)

	var configsRes []model.ProjectField
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func SelectOpFieldsByBlockID(proCode string, block_id string) (err error, configs []model.ProjectField) {
	db := global.ProDbMap[proCode].Model(&model.ProjectField{})
	// db.Select("id, created_at, updated_at, name, code, bill_id, block_id, block_index, field_index, op1_value, op1_input, op2_value, op2_input, op_q_value, op_q_input, op_d_value, op_d_input, result_value, result_input, final_value, final_input, right_value")
	if block_id != "" {
		db = db.Where("block_id = ?", block_id)
	}
	// db.Group("block_index")
	db.Order("updated_at asc").Order("block_index asc").Order("field_index asc")
	var configsRes []model.ProjectField
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func SelectBillFields(proCode string, bill_id string, block_index int, code string) (err error, configs model.ProjectField) {
	db := global.ProDbMap[proCode].Model(&model.ProjectField{})
	// if block_id != "" {
	db = db.Where("bill_id = ? AND code = ?", bill_id, code)
	if block_index != -1 {
		db = db.Where("block_index = ? ", block_index)
	}
	// db.Order("block_index asc").Order("field_index asc")
	var configsRes model.ProjectField
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func GetZeroFieldsValue(proCode string, block_id string, value string) (err error, configs []model.ProjectField) {
	db := global.ProDbMap[proCode].Model(&model.ProjectField{})
	// if block_id != "" {
	// db = db.Where("bill_id = ? AND code = ? AND  = ?", bill_id, "", )
	// }
	// db.Order("block_index asc").Order("field_index asc")

	var configsRes []model.ProjectField
	err = db.Raw("SELECT * FROM project_fields WHERE block_id = ? AND block_index = ( SELECT block_index FROM project_fields WHERE block_id = ? AND code = 'fc059' AND result_value = ? )", block_id, block_id, value).Scan(&configsRes).Error
	// err = db.Find(&configsRes).Error
	return err, configsRes
}

func GetZeroFieldsValueByCode(proCode string, block_id string, code, value string) (err error, configs []model.ProjectField) {
	db := global.ProDbMap[proCode].Model(&model.ProjectField{})
	// if block_id != "" {
	// db = db.Where("bill_id = ? AND code = ? AND  = ?", bill_id, "", )
	// }
	// db.Order("block_index asc").Order("field_index asc")

	var configsRes []model.ProjectField
	err = db.Raw("SELECT * FROM project_fields WHERE block_id = ? AND block_index = ( SELECT block_index FROM project_fields WHERE block_id = ? AND code = ? AND result_value = ? )", block_id, block_id, code, value).Scan(&configsRes).Error
	// err = db.Find(&configsRes).Error
	return err, configsRes
}
