/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/16 14:18
 */

package service

import (
	"go.uber.org/zap"
	"server/global"
	"server/module/export/model"
	model3 "server/module/pro_manager/model"
	"server/module/report_management/model/request"
	model2 "server/module/sys_base/model"
)

// PageNewHospitalAndCatalogue 分页获取列表
func PageNewHospitalAndCatalogue(search request.NewHospitalAndCatalogueSearch) (err error, total int64, list []model.HospitalCatalogue) {
	limit := search.PageSize
	offset := search.PageSize * (search.PageIndex - 1)
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, 0, list
	}
	db = db.Model(&model.HospitalCatalogue{})
	if search.Type > 0 {
		db = db.Where("type = ?", search.Type)
	}
	err = db.Count(&total).Error
	err = db.Order("created_at").
		Where("created_at BETWEEN ? AND ? ", search.StartTime, search.EndTime).
		Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

// ExportNewHospitalAndCatalogue 导出列表
func ExportNewHospitalAndCatalogue(search request.NewHospitalAndCatalogueExportSearch) (err error, list []model.HospitalCatalogue) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, list
	}
	db = db.Model(&model.HospitalCatalogue{})
	if search.Type > 0 {
		db = db.Where("type = ?", search.Type)
	}
	err = db.Order("created_at").
		Where("created_at BETWEEN ? AND ? ", search.StartTime, search.EndTime).
		Limit(100000).Find(&list).Error
	return err, list
}

// PageExtractAgency 分页获取列表
func PageExtractAgency(search model2.BaseTimePageCode) (err error, total int64, list []model.Agency) {
	limit := search.PageSize
	offset := search.PageSize * (search.PageIndex - 1)
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, 0, list
	}
	db = db.Model(&model.Agency{})
	err = db.Count(&total).Error
	err = db.Order("created_at").
		Where("created_at BETWEEN ? AND ? ", search.StartTime, search.EndTime).
		Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

// ExportExtractAgency 导出列表
func ExportExtractAgency(search model2.BaseTimeRangeWithCode) (err error, list []model.Agency) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, list
	}
	db = db.Model(&model.Agency{})
	err = db.Order("created_at").
		Where("created_at BETWEEN ? AND ? ", search.StartTime, search.EndTime).
		Limit(100000).Find(&list).Error
	return err, list
}

// FetchBills 获取单据
func FetchBills(search model2.BaseTimeRangeWithCode) (err error, list []model3.ProjectBill) {
	db := global.ProDbMap[search.ProCode]
	if db == nil {
		return global.ProDbErr, list
	}
	var countBill int64
	batchSize := 100
	db.Model(&model3.ProjectBill{}).
		Where("(stage = 5 or stage = 3 or stage = 7) and created_at BETWEEN ? AND ? ", search.StartTime, search.EndTime).
		Count(&countBill)
	for offset := 0; offset < int(countBill); offset += batchSize {
		var bills []model3.ProjectBill
		err = db.Where("(stage = 5 or stage = 3 or stage = 7) and created_at BETWEEN ? AND ? ", search.StartTime, search.EndTime).
			Select("id,bill_name,bill_num,pro_code,stage,download_path,download_at,batch_num,agency,created_at").
			Order("created_at").Limit(batchSize).Offset(offset).Find(&bills).Error
		if err != nil {
			global.GLog.Error("计算字符数Find bills", zap.Error(global.ProDbErr))
			return err, list
		}
		list = append(list, bills...)
	}
	return err, list
}
