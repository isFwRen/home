package service

import (
	"server/global"
	eModel "server/module/export/model"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

// 获取扣费明细信息
func GetDeductionDetails(DeductionDetailSearch model.DeductionDetailSearch) (err error, total int64, list []eModel.DeductionDetail) {
	limit := DeductionDetailSearch.PageSize
	offset := DeductionDetailSearch.PageSize * (DeductionDetailSearch.PageIndex - 1)
	var projectBills []eModel.DeductionDetail

	db := global.ProDbMap[DeductionDetailSearch.ProCode]
	if db == nil {
		return global.ProDbErr, 0, projectBills
	}

	db = db.Model(&eModel.DeductionDetail{}).Where("date BETWEEN ? AND ? ", DeductionDetailSearch.TimeStart, DeductionDetailSearch.TimeEnd)

	if DeductionDetailSearch.BillCode != "" {
		db.Where("bill_num LIKE ? ", "%"+DeductionDetailSearch.BillCode+"%")
	}

	if DeductionDetailSearch.Name != "" {
		db.Where("inventory_name LIKE ? ", "%"+DeductionDetailSearch.Name+"%")
	}

	err = db.Count(&total).Error
	if limit >= 0 {
		err = db.Order("date").Limit(limit).Offset(offset).Find(&projectBills).Error
	} else {
		err = db.Order("date").Find(&projectBills).Error
	}

	return err, total, projectBills
	//获取history单据信息
	// err = db.Model(&eModel.DeductionDetail{}).Where("id = ?", reqParam.ID).Find(&b).Error
	// return err, b
}

func ExportDeductionDetailExcel(DeductionDetailSearch model.DeductionDetailSearch) (err error, path, name string) {
	DeductionDetailSearch.PageSize = -1
	err, _, projectBills := GetDeductionDetails(DeductionDetailSearch)
	if err != nil {
		return err, "", ""
	}

	s := strings.Replace(DeductionDetailSearch.TimeStart.Format("2006-01-02"), " 00:00:00", "", -1)
	e := strings.Replace(DeductionDetailSearch.TimeEnd.Format("2006-01-02"), " 23:59:59", "", -1)
	bookName := "扣费明细" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "扣费明细/" + DeductionDetailSearch.ProCode + "/"
	// 本地
	//basicPath := "./"
	err = utils.ExportBigExcel(basicPath, bookName, "", projectBills)
	return err, basicPath, bookName
}
