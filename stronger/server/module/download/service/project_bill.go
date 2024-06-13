package service

import (
	"server/global"
	model3 "server/module/download/model"
	model2 "server/module/msg_manager/model"
	"server/module/pro_manager/model"
)

func InsertBill(proCode string, agingConfig model.ProjectBill) (err error, configInter model.ProjectBill) {
	agingConfig.IsAutoUpload = true
	err = global.ProDbMap[proCode].Create(&agingConfig).Error
	return err, agingConfig
}

func InsertBills(proCode string, agingConfigs []model.ProjectBill) (err error, configInter []model.ProjectBill) {
	// agingConfig.IsAutoUpload = true
	err = global.ProDbMap[proCode].Create(&agingConfigs).Error
	return err, agingConfigs
}

func DelBill(proCode string, id string) (err error) {
	var configAging model.ProjectBill
	err = global.ProDbMap[proCode].Where("id = ?", id).Delete(&configAging).Error
	return err
}

func UpdateBill(proCode string, configAging model.ProjectBill, id string) (err error) {
	err = global.ProDbMap[proCode].Where("id = ?", id).First(&model.ProjectBill{}).Save(&configAging).Error
	return err
}

func SelectBill(proCode string, name string) (err error, configs []model.ProjectBill) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBill{})
	if name != "" {
		db = db.Where("bill_name = ?", name)
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

func CountBillByBatchNum(proCode string, name string) (err error, total int64) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, 0
	}
	db = db.Model(&model.ProjectBill{})
	if name != "" {
		db = db.Where("batch_num = ?", name)
	}
	err = db.Count(&total).Error
	return err, total
}

// InsertCustomerNotice 新增客户通知2
func InsertCustomerNotice(proCode string, customerNotice model2.CustomerNotice) (err error) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	err = db.Create(&customerNotice).Error
	return err
}

// GetCustomerNoticeByStatus 统计有多少条通知函2
func GetCustomerNoticeByStatus(proCode string) (err error, total int64) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, 0
	}
	err = db.Model(&model2.CustomerNotice{}).Where("status = 1").Count(&total).Error
	err = db.Count(&total).Error
	return err, total
}

// CountBillByName 统计有多少条通知函2待回复
func CountBillByName(proCode string, name string) (err error, total int64) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, 0
	}
	db = db.Model(&model2.CustomerNotice{})
	if name != "" {
		db = db.Where("file_name = ?", name)
	}
	err = db.Count(&total).Error
	return err, total
}

// FindBillByStage 获取存在已接收的单
func FindBillByStage(proCode string, name string) (err error, bill model.ProjectBill) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, bill
	}
	db = db.Model(&model.ProjectBill{})
	if name != "" {
		db = db.Where("bill_num = ? and stage = 7", name)
	}
	err = db.First(&bill).Error
	return err, bill
}

// UpdateBillPackCode 更新单据通知函信息
func UpdateBillPackCode(bill model.ProjectBill) (err error) {

	db := global.ProDbMap[bill.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	return db.Model(&model.ProjectBill{}).Where("id = ? and stage = 7", bill.ID).Updates(map[string]interface{}{
		"pack_code":  bill.PackCode,
		"wrong_note": bill.WrongNote,
		"stage":      3,
	}).Error
}

// InsertUpdateConstLogs 插入多条更新常量log
func InsertUpdateConstLogs(log []model3.UpdateConstLog) error {
	return global.GDb.Model(&model3.UpdateConstLog{}).Create(&log).Error
}

// InsertUpdateConstLog 插入1条更新常量log
func InsertUpdateConstLog(log model3.UpdateConstLog) error {
	return global.GDb.Model(&model3.UpdateConstLog{}).Create(&log).Error
}

// UpdateBillStage 更新待下载文件单据流程状态为待加载
func UpdateBillStage(bill model.ProjectBill) (err error) {
	db := global.ProDbMap[bill.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	return db.Model(&model.ProjectBill{}).Where("id = ? and stage = 8", bill.ID).Updates(map[string]interface{}{
		"stage":      1,
		"images":     bill.Images,
		"status":     bill.Status,
		"wrong_note": bill.WrongNote,
	}).Error
}

// FetchBillByStage 获取存在待下载的单
func FetchBillByStage(proCode string) (err error, bills []model.ProjectBill) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, bills
	}
	err = db.Model(&model.ProjectBill{}).
		Where("stage = 8").
		Order("created_at ASC").
		Limit(10).
		Find(&bills).Error
	return err, bills
}
