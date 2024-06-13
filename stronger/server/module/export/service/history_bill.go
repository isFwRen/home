/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 4:37 下午
 */

package service

import (
	"errors"
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/model"
	"server/module/export/utils"
	model3 "server/module/load/model"
	model2 "server/module/pro_manager/model"
	model4 "server/module/report_management/model"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

// SaveBill 保存单的信息到历史库
func SaveBill(obj model.ResultDataBill) (err error) {
	bill := obj.Bill
	if bill.ID == "" {
		return errors.New("单据id为空")
	}
	db := global.ProDbMap[bill.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()

	allFields := make([]model3.ProjectField, 0)
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					allFields = append(allFields, fieldsArr[k]...)
				}
			}
		}
	}

	questionNum := 0
	for _, field := range allFields {
		questionNum += len(field.Issues)
		if field.IsChange == false {
			continue
		}
		err = tx.Model(&model3.ProjectField{}).Where("id = ?", field.ID).
			Select("result_value").
			Updates(model3.ProjectField{ResultValue: field.ResultValue}).Error
		if err != nil {
			global.GLog.Error(err.Error())
			tx.Rollback()
			return errors.New("保存出错啦")
		}
	}

	//CSB0108RC0011000
	//1、当案件导出校验的描述包含“清单项目金额”二字时，将案件录入状态改为“待审核”；
	//2、当案件导出校验的描述不包含“清单项目金额”二字时，将案件录入状态改为“已导出”；
	stage := 3
	if strings.Index(bill.WrongNote, "清单项目金额") != -1 {
		stage = 4
	}
	var projectBillOne model2.ProjectBill
	tx.Model(&model2.ProjectBill{}).Where("id = ?", bill.ID).Find(&projectBillOne)
	if projectBillOne.ExportStage == 0 {
		if stage == 3 || stage == 4 {
			err = tx.Model(&model2.ProjectBill{}).Where("id = ?", bill.ID).Updates(
				map[string]interface{}{
					"export_at":      time.Now(),
					"wrong_note":     bill.WrongNote,
					"is_auto_upload": false,
					"invoice_num":    len(obj.Invoice) - 1,
					//QualityUser:  bill.QualityUser,
					"question_num": questionNum,
					"stage":        stage,
					"export_stage": stage,
					"count_money":  bill.CountMoney,
					"bill_type":    bill.BillType,
				}).Error
		} else {
			err = tx.Model(&model2.ProjectBill{}).Where("id = ?", bill.ID).Updates(
				map[string]interface{}{
					"export_at":      time.Now(),
					"wrong_note":     bill.WrongNote,
					"is_auto_upload": false,
					"invoice_num":    len(obj.Invoice) - 1,
					//QualityUser:  bill.QualityUser,
					"question_num": questionNum,
					"stage":        stage,
					"count_money":  bill.CountMoney,
					"bill_type":    bill.BillType,
				}).Error
		}
	} else {
		err = tx.Model(&model2.ProjectBill{}).Where("id = ?", bill.ID).Updates(
			map[string]interface{}{
				"export_at":      time.Now(),
				"wrong_note":     bill.WrongNote,
				"is_auto_upload": false,
				"invoice_num":    len(obj.Invoice) - 1,
				//QualityUser:  bill.QualityUser,
				"question_num": questionNum,
				"stage":        stage,
				"count_money":  bill.CountMoney,
				"bill_type":    bill.BillType,
			}).Error
	}

	return tx.Commit().Error
}

func AutoExportFieldUpdate(obj model.ResultDataBill) (err error) {
	bill := obj.Bill
	if bill.ID == "" {
		return errors.New("单据id为空")
	}
	db := global.ProDbMap[bill.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	allFields := make([]model3.ProjectField, 0)
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					allFields = append(allFields, fieldsArr[k]...)
				}
			}
		}
	}

	for _, field := range allFields {
		if field.IsChange == false {
			continue
		}
		err = tx.Model(&model3.ProjectField{}).Where("id = ?", field.ID).
			Select("result_value").
			Updates(model3.ProjectField{ResultValue: field.ResultValue}).Error
		if err != nil {
			global.GLog.Error(err.Error())
			tx.Rollback()
			return errors.New("保存出错啦")
		}
	}

	//进入系统的单默认手动回传状态为“否”，案件第一次导出：
	//1.若有导出校验则更改手动回传状态为“是”
	//2.若无导出校验则不更改手动回传状态
	isAutoUpload := obj.Bill.IsAutoUpload
	if obj.Bill.WrongNote != "" || obj.Bill.Remark != "" {
		isAutoUpload = false
	}
	err = tx.Model(&model2.ProjectBill{}).Where("id = ?", bill.ID).Updates(
		map[string]interface{}{
			"is_auto_upload": isAutoUpload,
		}).Error

	return tx.Commit().Error
}

// GetFieldsByBillIdAndBlockCode 获取所有清单分块的字段
func GetFieldsByBillIdAndBlockCode(billId string, blockCode string, proCode string) (err error, fields []model3.ProjectField) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, fields
	}
	err = db.Where("block_id in (?) ",
		db.Table("project_blocks").Select("id").Where("bill_id = ? and code = ?", billId, blockCode),
	).Find(&fields).Error
	return err, fields
}

// GetBillBlocksAndFields 获取单据,分块和字段
func GetBillBlocksAndFields(billId string, proCode string, wg *sync.WaitGroup, bill *model2.ProjectBill, blocks *[]model3.ProjectBlock, fields *[]model3.ProjectField) {
	defer utils.DealErr(wg)
	db := global.ProDbMap[proCode]
	db.Model(&model2.ProjectBill{}).Where("id = ?", billId).First(&bill)
	db.Model(&model3.ProjectBlock{}).Where("bill_id = ?", billId).Find(&blocks)
	db.Model(&model3.ProjectField{}).Where("bill_id = ?", billId).Find(&fields)
}

// SaveBlocksAndFields 保存分块和字段到历史库
func SaveBlocksAndFields(proCode string, b []model3.ProjectBlock, f []model3.ProjectField) error {
	db, isOk := global.ProDbMap[proCode]
	if !isOk {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.GLog.Error(fmt.Sprintf("%v", r))
			tx.Rollback()
		}
	}()

	err := tx.Model(&model3.ProjectBlock{}).Create(&b).Error
	if err != nil {
		global.GLog.Error(fmt.Sprintf("%v", err))
		tx.Rollback()
		return err
	}

	err = tx.Model(&model3.ProjectField{}).CreateInBatches(&f, 1000).Error
	if err != nil {
		global.GLog.Error(fmt.Sprintf("%v", err))
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		global.GLog.Error(fmt.Sprintf("%v", err))
		tx.Rollback()
		return err
	}
	return nil
}

// InsertWrongs 插入错误查询数据
func InsertWrongs(proCode string, wrongs []model4.Wrong, billId string) (err error, rows int64) {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr, 0
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var w model4.Wrong
	err = tx.Model(&model4.Wrong{}).Where("bill_id = ?", billId).Delete(&w).Error
	if err != nil {
		tx.Rollback()
		return err, 0
	}
	if len(wrongs) > 0 {
		rows = tx.Model(&model4.Wrong{}).Create(wrongs).RowsAffected
		if err != nil {
			tx.Rollback()
			return err, 0
		}
	}

	return tx.Commit().Error, rows
}

// InsertWrongs 插入错误查询数据
func InsertOcrs(proCode string, ocrStatisticses []model4.OcrStatistics, billId string) (err error, rows int64) {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr, 0
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var w model4.OcrStatistics
	err = tx.Model(&model4.OcrStatistics{}).Where("bill_id = ?", billId).Delete(&w).Error
	if err != nil {
		tx.Rollback()
		return err, 0
	}
	if len(ocrStatisticses) > 0 {
		rows = tx.Model(&model4.OcrStatistics{}).Create(ocrStatisticses).RowsAffected
		if err != nil {
			tx.Rollback()
			return err, 0
		}
	}

	return tx.Commit().Error, rows
}

func GetCountBillNum(proCode string, bill_num string) (err error, num int64) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, num
	}
	err = db.Model(&model2.ProjectBill{}).Where("bill_num = ?", bill_num).Count(&num).Error
	return err, num
}

func GetCountValueByCode(proCode, id, code, value string) (err error, num int64) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, num
	}
	err = db.Model(&model3.ProjectField{}).Where("bill_id != ? and code = ? and result_value = ?", id, code, value).Count(&num).Error
	return err, num
}

// FetchInvoiceNum 查询发票号的单号——非当前单
func FetchInvoiceNum(proCode, billNum, invoiceNum string) (invoiceNums []model.InvoiceNum) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return invoiceNums
	}
	db.Model(&model.InvoiceNum{}).Where("num = ? and bill_num <> ?", invoiceNum, billNum).
		Find(&invoiceNums)
	return invoiceNums
}

// UpdateInvoiceNumByBillNum 更新收集当前单据的发票号
func UpdateInvoiceNumByBillNum(proCode, billNum string, invoiceNums []model.InvoiceNum) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	err := tx.Model(&model.InvoiceNum{}).Delete(&[]model.InvoiceNum{}, "bill_num  = ?", billNum).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(invoiceNums) > 0 {
		err = tx.Model(&model.InvoiceNum{}).Create(invoiceNums).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// UpdateHospitalCatalogue 收集不在常量库的数据field
func UpdateHospitalCatalogue(proCode string, billId string, hospitalCatalogues []model.HospitalCatalogue) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	err := tx.Model(&model.HospitalCatalogue{}).Delete(&[]model.HospitalCatalogue{}, "bill_id  = ?", billId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(hospitalCatalogues) > 0 {
		err = tx.Model(&model.HospitalCatalogue{}).Create(hospitalCatalogues).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// UpdateAgency 机构抽取
func UpdateAgency(proCode string, agency model.Agency) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	return db.Model(&model.Agency{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bill_id"}},
		UpdateAll: true,
	}).Create(&agency).Error
}
