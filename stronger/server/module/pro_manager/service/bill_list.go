/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:47 下午
 */

package service

import (
	"encoding/json"
	"fmt"
	"server/global"
	model2 "server/module/load/model"
	model3 "server/module/pro_conf/model"
	pf "server/module/pro_conf/model"
	"server/module/pro_manager/const_data"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/response"
	pro118 "server/module/pro_manager/project/B0118"
	util "server/utils"
	"strconv"
	"time"
	model4 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetBillByPage 根据一堆过滤条件获取案件列表
func GetBillByPage(billListSearch model.BillListSearch) (err error, total int64, list interface{}) {
	limit := billListSearch.PageSize
	offset := billListSearch.PageSize * (billListSearch.PageIndex - 1)
	var projectBills []model.ProjectBill

	db := global.ProDbMap[billListSearch.ProCode]
	if db == nil {
		return global.ProDbErr, 0, projectBills
	}
	if billListSearch.BillCode != "" {
		db = db.Where("bill_name like ? ", "%"+billListSearch.BillCode+"%")
	}
	if billListSearch.Status != -1 {
		db = db.Where("status = ?", billListSearch.Status)
	}
	if billListSearch.SaleChannel != "" {
		db = db.Where("sale_channel like ?", "%"+billListSearch.SaleChannel+"%")
	}
	if billListSearch.BatchNum != "" {
		db = db.Where("batch_num like ?", "%"+billListSearch.BatchNum+"%")
	}
	if billListSearch.Agency != "" {
		db = db.Where("agency like ?", "%"+billListSearch.Agency+"%")
	}
	if billListSearch.ClaimType != -1 {
		db = db.Where("claim_type = ?", billListSearch.ClaimType)
	}
	if billListSearch.InsuranceType != "" {
		db = db.Where("insurance_type like ? ", billListSearch.InsuranceType+"%")
	}
	if billListSearch.StickLevel != -1 {
		db = db.Where("stick_level = ?", billListSearch.StickLevel)
	}
	if billListSearch.IsQuestion != -1 {
		if billListSearch.IsQuestion == 1 {
			db = db.Where("question_num > ? ", 0)
		} else {
			db = db.Where("question_num <= ? ", 0)
		}
	}
	if billListSearch.InvoiceNum != -1 {
		db = db.Where("invoice_num = ? ", billListSearch.InvoiceNum)
	}
	if billListSearch.QualityUser != "" {
		db = db.Where("CONCAT(quality_user_code,quality_user_name) like ? ", "%"+billListSearch.QualityUser+"%")
	}
	if billListSearch.Stage != "" {
		db = db.Where("stage = ? ", billListSearch.Stage)
	}

	if billListSearch.MaxCountMoney != 0 {
		db = db.Where("count_money BETWEEN ? AND ? ", billListSearch.MinCountMoney, billListSearch.MaxCountMoney)
	} else {
		db = db.Where("count_money >= ? ", billListSearch.MinCountMoney)
	}

	db = db.Model(&model.ProjectBill{}).Where("created_at BETWEEN ? AND ? ", billListSearch.TimeStart, billListSearch.TimeEnd)

	err = db.Count(&total).Error

	if billListSearch.OrderBy != "" {
		var orderBys [][]string
		err = json.Unmarshal([]byte(billListSearch.OrderBy), &orderBys)
		global.GLog.Error("", zap.Error(err))
		if err == nil {
			for _, orderBy := range orderBys {
				global.GLog.Info("", zap.Any("", orderBy))
				db = db.Order(const_data.OrderBy[orderBy[0]] + " " + orderBy[1])
			}
		}
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&projectBills).Error
	return err, total, projectBills
}

// UpdateBillObjInForceExport 强制导出更新单
func UpdateBillObjInForceExport(reqParam model.DelByIdAndProCode, billObj model.BillObj) error {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("插入历史库失败")
		}
	}()
	//新增block
	bRow := tx.Model(&model2.ProjectBlock{}).Create(billObj.ProjectBlockList).RowsAffected
	global.GLog.Info("新增分块数::", zap.Int64("row", bRow))

	fRow := tx.Model(&model2.ProjectField{}).Create(billObj.ProjectFieldList).RowsAffected
	global.GLog.Info("新增字段数::", zap.Int64("row", fRow))
	billObj.ProjectBill.ID = reqParam.ID
	//更新整个bill数据
	billRow := tx.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).Updates(map[string]interface{}{
		"status":         1,
		"stage":          4,
		"is_auto_upload": false,
		"wrong_note":     billObj.ProjectBill.WrongNote,
		//billObj.ProjectBill
	}).RowsAffected
	global.GLog.Info("更新单据数::", zap.Int64("row", billRow))
	return tx.Commit().Error
}

// UpdateBillObjInDel 删除单据的更新单信息
func UpdateBillObjInDel(reqParam model.DelByIdAndProCode, billObj model.BillObj, customClaims *model4.CustomClaims) error {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("插入历史库失败")
		}
	}()
	//新增block
	bRow := tx.Model(&model2.ProjectBlock{}).Create(billObj.ProjectBlockList).RowsAffected
	global.GLog.Info("新增分块数::", zap.Int64("row", bRow))

	fRow := tx.Model(&model2.ProjectField{}).Create(billObj.ProjectFieldList).RowsAffected
	global.GLog.Info("新增字段数::", zap.Int64("row", fRow))
	billObj.ProjectBill.ID = reqParam.ID
	var timeZero time.Time
	//更新整个bill数据
	//2.理赔类型 ，导出时间，回传时间，理赔类型，账单金额，发票数量，问题件，质检人
	billRow := tx.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).Updates(map[string]interface{}{
		"status":            4,
		"pre_status":        gorm.Expr("status"),
		"is_auto_upload":    false,
		"wrong_note":        "",
		"del_remarks":       "（" + customClaims.Code + customClaims.NickName + time.Now().Format("2006-01-02 15:04:05") + "）：" + reqParam.DelRemarks,
		"export_at":         timeZero,
		"upload_at":         timeZero,
		"claim_type":        0,
		"count_money":       0,
		"invoice_num":       0,
		"question_num":      0,
		"quality_user_code": "",
		"quality_user_name": "",
		"export_stage":      0,
	}).RowsAffected
	global.GLog.Info("更新单据数::", zap.Int64("row", billRow))
	return tx.Commit().Error
}

// GetBillInfo 获取history单据信息
func GetBillInfo(reqParam model.DelByIdAndProCode) (err error, b model.ProjectBill) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, b
	}
	//获取history单据信息
	err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).First(&b).Error
	return err, b
}

// DelByIdsAndProCode 批量删除
func DelByIdsAndProCode(reqParam model.DelByIdAndProCode) (e error, rows int64) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, 0
	}

	//更新状态
	tx := db.Begin().Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(fmt.Sprintf("DelByIdsAndProCode→→→→→%v", r))
		}
	}()
	rowsAffected := tx.UpdateColumn("pre_status", gorm.Expr("status")).RowsAffected
	affected := tx.UpdateColumns(model.ProjectBill{
		Status:     4,
		DelRemarks: reqParam.DelRemarks,
		WrongNote:  reqParam.DelRemarks,
	}).RowsAffected
	if rowsAffected != affected {
		panic("更新行数不一样")
	}
	tx.Commit()
	return nil, rowsAffected
}

// RecoverBill 恢复单据
func RecoverBill(reqParam model.ProCodeAndId) (err error) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(fmt.Sprintf("DelByIdsAndProCode→→→→→%v", r))
		}
	}()

	err = tx.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).
		UpdateColumn("status", gorm.Expr("pre_status")).
		Error

	rmB := tx.Where("bill_id = ?", reqParam.ID).Delete(&model2.ProjectBlock{}).RowsAffected
	global.GLog.Info("删除的分块条数有", zap.Int64("", rmB))
	rmF := tx.Where("bill_id = ?", reqParam.ID).Delete(&model2.ProjectField{}).RowsAffected
	global.GLog.Info("删除的字段条数有", zap.Int64("", rmF))

	return tx.Commit().Error
}

// GetBillObj 获取整个单据
func GetBillObj(reqParam model.ProCodeAndId) (err error, obj model.BillObj) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, obj
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).
		First(&obj.ProjectBill).Error

	err = db.Model(&model2.ProjectBlock{}).Where("bill_id = ?", reqParam.ID).
		Find(&obj.ProjectBlockList).Error
	err = db.Model(&model2.ProjectField{}).Where("bill_id = ?", reqParam.ID).
		Find(&obj.ProjectFieldList).Error
	return err, obj
}

// GetProBillById 查询项目单据
func GetProBillById(reqParam model.ProCodeAndId) (err error, bill model.ProjectBill) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, bill
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).First(&bill).
		Error
	return err, bill
}

// UploadState 更新单据状态
func UploadState(reqParam model.ProCodeAndId, stage int) (err error) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).Update("stage", stage).Error
	return err
}

// UploadRemark 更新备注
func UploadRemark(reqParam model.Remark) (err error) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	err = db.Model(&model.ProjectBill{}).
		Where("id = ? and edit_version = ?", reqParam.ID, reqParam.EditVersion).
		Updates(map[string]interface{}{
			"edit_version": reqParam.EditVersion + 1,
			"remark":       reqParam.Remark,
		}).Error
	return err
}

// SetUploadType 更新回传方式
func SetUploadType(reqParam model.AutoUpload) (err error, b model.ProjectBill) {
	db := global.ProDbMap[reqParam.ProCode]
	dbCache := global.ProDbMap[reqParam.ProCode+"_task"]
	if db == nil || dbCache == nil {
		return global.ProDbErr, b
	}
	//err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).First(&b).Error
	err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).Update("is_auto_upload", reqParam.IsAutoUpload).Error
	err = dbCache.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).Update("is_auto_upload", reqParam.IsAutoUpload).Error
	return err, b
}

// UpdateField 更新字段
func UpdateField(p model.EditBillResultDataManyFields) (err error) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Model(&model2.ProjectField{}).Begin()
	for _, field := range p.Fields {
		err = tx.Where("id = ?", field.ID).
			Select("result_value", "field_input").
			Updates(model2.ProjectField{ResultValue: field.ResultValue, ResultInput: field.ResultInput}).Error
		if err != nil {
			global.GLog.Error(err.Error())
			tx.Rollback()
		}
	}
	return tx.Commit().Error
}

// GetField 获取字段信息
func GetField(p model.EditBillResultData) (err error, f model2.ProjectField) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr, f
	}
	err = db.Model(&model2.ProjectField{}).
		Where("id = ?", p.FieldId).
		First(&f).Error
	return err, f
}

// GetBillAndBlocks 获取单据和分开
func GetBillAndBlocks(p model.EditBillResultDataManyFields) (err error, bill model.ProjectBill, blocks []model2.ProjectBlock, fieldsLen int64) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr, bill, blocks, 0
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", p.BillId).First(&bill).Error
	if err != nil {
		return err, bill, blocks, 0
	}
	err = db.Model(&model2.ProjectBlock{}).Where("bill_id = ?", p.BillId).Find(&blocks).Error
	if err != nil {
		return err, bill, blocks, 0
	}
	db.Model(&model2.ProjectField{}).Where("bill_id = ?", p.BillId).Count(&fieldsLen)
	return err, bill, blocks, fieldsLen
}

// GetBill 获取单据
func GetBill(p model.EditBillResultDataManyFields) (err error, bill model.ProjectBill) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr, bill
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", p.BillId).First(&bill).Error
	return err, bill
}

// GetBlockAndFields 获取字段和分块
func GetBlockAndFields(p model.EditBillResultDataManyFields) (err error, blocks []model2.ProjectBlock, fields []model2.ProjectField, fieldsLen int64) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr, blocks, fields, 0
	}
	err = db.Model(&model2.ProjectBlock{}).Where("bill_id = ?", p.BillId).Find(&blocks).Error
	if err != nil {
		return err, blocks, fields, 0
	}
	if p.EditType == 1 {
		err = db.Model(&model2.ProjectField{}).Where("bill_id = ?", p.BillId).Order("block_index asc").Order("field_index asc").Find(&fields).Error
	}
	db.Model(&model2.ProjectField{}).Where("bill_id = ?", p.BillId).Count(&fieldsLen)
	return err, blocks, fields, fieldsLen
}

// GetFieldConf 获取字段配置
func GetFieldConf(proCode string) (err error, sysProFields []model3.SysProField) {
	proId := global.ProCodeId[proCode]
	if proId == "" {
		return global.NotPro, sysProFields
	}
	err = global.GDb.Model(&model3.SysProField{}).Where("pro_id = ?", proId).Preload("SysIssues", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).Find(&sysProFields).Error
	return err, sysProFields
}

// GetFieldInfo 获取字段、分块和字段配置信息
func GetFieldInfo(p model.ProCodeAndId) (err error, obj model.FieldObj) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr, obj
	}

	proId := global.ProCodeId[p.ProCode]
	if proId == "" {
		return global.NotPro, obj
	}
	err = db.Model(&model2.ProjectField{}).Where("id = ?", p.ID).First(&obj.Field).Error
	if err != nil {
		return err, obj
	}
	err = db.Model(&model2.ProjectBlock{}).Where("id = ?", obj.Field.BlockID).First(&obj.Block).Error
	if err != nil {
		return err, obj
	}
	err = global.GDb.Model(&model3.SysProField{}).Where("code = ? and pro_id = ? ", obj.Field.Code, proId).First(&obj.FieldConf).Error
	return err, obj
}

// UpdateFeedback 修改反馈值
func UpdateFeedback(p model.EditFeedbackVal) (err error) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	err = db.Model(&model2.ProjectField{}).Where("id = ?", p.FieldId).
		Updates(model2.ProjectField{ResultValue: p.FieldValue, FeedbackDate: p.EditDate}).Error
	if err != nil {
		return err
	}
	return err
}

// SetPractice 设置为练习单
func SetPractice(p model.SetPracticeForm) (err error) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(fmt.Sprintf("SetPractice→→→→→%v", r))
		}
	}()

	err = tx.Model(&model2.ProjectBlock{}).Where("id in (?)", p.BlockIds).
		Updates(map[string]interface{}{
			"is_practice": p.IsPractice,
		}).Error
	err = tx.Model(&model2.ProjectField{}).Where("block_id in (?)", p.BlockIds).
		Updates(map[string]interface{}{
			"is_practice": p.IsPractice,
		}).Error
	return tx.Commit().Error
}

// GetLog 获取单据的修改日志
func GetLog(p model.ProCodeAndId) (err error, logs []model.ResultDataLog) {
	db := global.ProDbMap[p.ProCode]
	if db == nil {
		return global.ProDbErr, logs
	}
	err = db.Model(&model.ResultDataLog{}).Order("created_at desc").Where("bill_id = ?", p.ID).
		Find(&logs).Error
	return err, logs
}

// GetQualitiesByPro 获取项目质检配置
func GetQualitiesByPro(pro string) (err error, s []model3.SysQuality) {
	proId := global.ProCodeId[pro]
	if proId == "" {
		return global.NotPro, s
	}
	err = global.GDb.Model(&model3.SysQuality{}).Where("pro_id = ?", proId).Order("belong_type desc").
		Find(&s).Error
	return err, s
}

//type QQQ struct {
//	ID   int64  `json:"id,string"`
//	Name string `json:"name"`
//}
//func QQ() (err error, s QQQ) {
//
//	//err = global.GDb.Model(&QQQ{}).Create(QQQ{
//	//	920628794505560064,
//	//}).Error
//	err = global.GDb.First(&s).Error
//	return err, s
//}

// GetTaskById 获取任务的整个单据
func GetTaskById(reqParam model.DelByIdAndProCode) (err error, obj model.BillObj) {
	db := global.ProDbMap[reqParam.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, obj
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).
		First(&obj.ProjectBill).Error

	err = db.Model(&model2.ProjectBlock{}).Where("bill_id = ?", reqParam.ID).
		Find(&obj.ProjectBlockList).Error
	err = db.Model(&model2.ProjectField{}).Where("bill_id = ?", reqParam.ID).
		Find(&obj.ProjectFieldList).Error
	return err, obj
}

// DelTaskById 删除任务
func DelTaskById(reqParam model.DelByIdAndProCode) (e error, row int64) {
	db := global.ProDbMap[reqParam.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, 0
	}

	//更新状态
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("删除task失败")
		}
	}()
	//return db.Transaction(func(tx *gorm.DB) error {
	//	err := tx.Where("id = ?", reqParam.ID).Delete(&model.ProjectBill{}).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//	err = tx.Where("bill_id = ?", reqParam.ID).Delete(&model2.ProjectBlock{}).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//	err = tx.Where("bill_id = ?", reqParam.ID).Delete(&model2.ProjectField{}).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//	return nil
	//})
	rowsAffected := tx.Where("id = ?", reqParam.ID).Delete(&model.ProjectBill{}).RowsAffected
	rowsB := tx.Where("bill_id = ?", reqParam.ID).Delete(&model2.ProjectBlock{}).RowsAffected
	global.GLog.Info("删除的分块条数有", zap.Int64("", rowsB))
	rowsF := tx.Where("bill_id = ?", reqParam.ID).Delete(&model2.ProjectField{}).RowsAffected
	global.GLog.Info("删除的字段条数有", zap.Int64("", rowsF))
	err := tx.Commit().Error
	return err, rowsAffected
}

// InsertTaskBillObj 插入整张单
func InsertTaskBillObj(reqParam model.ProCodeAndId, billObj model.BillObj) error {
	db := global.ProDbMap[reqParam.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("新增task失败")
		}
	}()
	//新增block
	bRow := tx.Model(&model2.ProjectBlock{}).Create(billObj.ProjectBlockList).RowsAffected
	global.GLog.Info("新增分块数::", zap.Int64("row", bRow))

	fRow := tx.Model(&model2.ProjectField{}).Create(billObj.ProjectFieldList).RowsAffected
	global.GLog.Info("新增字段数::", zap.Int64("row", fRow))

	//更新整个bill数据
	billRow := tx.Model(&model.ProjectBill{}).Create(billObj.ProjectBill).RowsAffected
	global.GLog.Info("更新单据数::", zap.Int64("row", billRow))
	return tx.Commit().Error
}

// GetBlockById 根据分块id获取分块，单据和关联分块
func GetBlockById(reqParam model.ProCodeAndId) (err error, blocks []response.BlockImg, bill model.ProjectBill) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, blocks, bill
	}
	var block response.BlockImg
	var blocks1 []response.BlockImg
	err = db.Model(&model2.ProjectBlock{}).Where("id = ?", reqParam.ID).First(&block).Error
	if err != nil {
		return err, blocks, bill
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", block.BillID).First(&bill).Error
	if err != nil {
		return err, blocks, bill
	}
	err = db.Model(&model2.ProjectBlock{}).Where("code in (?)", block.LinkBCode).Find(&blocks1).Error
	if err != nil {
		return err, blocks, bill
	}
	blocks = append(append(blocks, blocks1...), block)
	return nil, blocks, bill
}

// GetFieldById 根据id获取字段
func GetFieldById(reqParam model.ProCodeAndId) (err error, fields []model2.ProjectField) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, fields
	}
	var field model2.ProjectField
	err = db.Model(&model2.ProjectField{}).Where("id = ?", reqParam.FieldId).First(&field).Error
	if err != nil {
		return err, fields
	}
	err = db.Model(&model2.ProjectField{}).Where("block_id = ? and block_index = ?", field.BlockID, field.BlockIndex).Find(&fields).Error
	if err != nil {
		return err, fields
	}
	return nil, fields
}

func GetUploadPath(reqParam model.ProCodeAndId) (err error, uploadPaths model3.SysProDownloadPaths) {
	id, ok := global.ProCodeId[reqParam.ProCode]
	if !ok {
		return global.NotPro, uploadPaths
	}
	err = global.GDb.Where("pro_id = ? and is_upload = true", id).First(&uploadPaths).Error
	return err, uploadPaths
}

// UpdateQualityUser 更新单据质检人
func UpdateQualityUser(customClaims *model4.CustomClaims, proCode, billId string, stage int, exportAt time.Time) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	err := db.Model(&model.ProjectBill{}).Where("id = ?", billId).Updates(map[string]interface{}{
		"quality_user_code": customClaims.Code,
		"quality_user_name": customClaims.Name,
		"export_stage":      stage,
		"first_export_at":   exportAt,
	}).Error

	return err
}

// GetTimeLinessBriefing 时效简报
func GetTimeLinessBriefing(proCode string) (err error, list map[string]interface{}) {
	var notReturn []string   // 超时30分钟未回传
	var in5Minutes []string  //时效剩余0-5分钟以内
	var in15Minutes []string //时效剩余5-15分钟以内
	var in20Minutes []string //时效剩余15-20分钟以内
	return err, list
	//找到待加载 录入中 未定义的单
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, list
	}
	var loadedBill []model.ProjectBill // 待加载 录入中 未定义
	// 获取合同时效
	var configs pf.SysProjectConfigAgingContract
	//contractDb := global.GDb.Model(&pf.SysProjectConfigAgingContract{}).Order("created_at ASC")
	////////////////////////循环  找到理赔类型相同的单和合同时效
	claimTypeNumber := []int{0, 3, 4, 5, 6, 7}
	for _, number := range claimTypeNumber {
		// 查看单 是否是节假日
		//循环单 查看每个单是不是在节假日 和 非节假日  区分开
		if number == 0 {
			//找到 该项目 未定义的时效
			contractDb := global.GDb.Model(&pf.SysProjectConfigAgingContract{}).Order("created_at ASC")
			contractDb.Where("code = ? AND claim_type = 0", proCode).Find(&configs)
			err = db.Model(&model.ProjectBill{}).Where("created_at > '2023-09-01 00:00:00'  AND created_at < '2023-09-25 23:59:59' AND stage = 1 OR stage = 2 AND claim_type = 0 AND status <> 4 ").Find(&loadedBill).Error
		} else {
			contractDb := global.GDb.Model(&pf.SysProjectConfigAgingContract{}).Order("created_at ASC")
			contractDb.Where("code = ? AND claim_type = ?", proCode, number).Find(&configs)
			err = db.Model(&model.ProjectBill{}).Where("created_at > '2023-09-01 00:00:00'  AND created_at < '2023-09-25 23:59:59' AND stage = 3 OR stage = 4 AND claim_type = ? AND status <> 4 ", number).Find(&loadedBill).Error
		}
		for _, bill := range loadedBill {
			//循环这些单  先计算最晚回传时间 然后最晚回传时间 - 当前时间  找到时效区间
			//节假日 agingHoliday
			var agingHoliday pf.SysProjectConfigAgingHoliday
			err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", bill.CreatedAt.Format("20061")).Find(&agingHoliday).Error
			if err != nil {
				return err, list
			}
			//根据单的日期来查看是否节假日 今天日期
			num, _ := strconv.Atoi(bill.CreatedAt.Format("02"))
			//创建时间
			billCreatedAt := bill.CreatedAt.Format("15:04:05")
			//转换成map
			_, m := pro118.C(agingHoliday)
			//判断单是时效内还是时效外
			_, ist, count := pro118.AgeingInner(configs, billCreatedAt)

			if configs.ContractStartTime == "" || configs.ContractEndTime == "" || billCreatedAt == "" {
				continue
			}
			//这张单的最晚回传时间
			var backAtTheLatest string
			if !m[num] {
				//工作日
				// 开始计算最晚回传时间
				// 查看是否在时效内 工作日
				//1.在时效内工作日
				if ist && count == 1 {
					requirementsTime := configs.RequirementsTime + "m" //考核要求时间
					duration, _ := time.ParseDuration(requirementsTime)
					latestTime := bill.CreatedAt.Add(duration) // 时效内最晚回传时间
					backAtTheLatest = latestTime.Format("2006-01-02 15:04:05")

				} else if ist && count == 2 {
					//时效外-----------------------
					err, backAtTheLatest = pro118.AgeingOutTimeDefault(m, num, bill.CreatedAt, configs, bill)
				}
			} else if m[num] {
				//节假日
				//节假日 - 时效外
				_, backAtTheLatest = pro118.AgeingOutTimeDefault(m, num, bill.CreatedAt, configs, bill)
			}
			//计算时间时效简报剩余时间  最晚回传时间 -  当前时间
			timeSub, _ := time.Parse("2006-01-02 15:04:05", backAtTheLatest)
			//createAt, _ := time.Parse("2006-01-02 15:04:05", billCreatedAt)
			format := time.Now().Format("2006-01-02 15:04:05")
			parse, _ := time.Parse("2006-01-02 15:04:05", format)
			briefingTime := timeSub.Sub(parse)
			minutes := briefingTime.Minutes()
			if minutes < -30 {
				notExist := util.IsRepeateds(notReturn, bill.BillNum)
				if notExist {
					notReturn = append(notReturn, bill.BillNum)
				}
			} else if minutes > 0 && minutes <= 5 {
				notExist := util.IsRepeateds(in5Minutes, bill.BillNum)
				if notExist {
					in5Minutes = append(in5Minutes, bill.BillNum)
				}
			} else if minutes > 5 && minutes <= 15 {
				notExist := util.IsRepeateds(in15Minutes, bill.BillNum)
				if notExist {
					in15Minutes = append(in15Minutes, bill.BillNum)
				}
			} else if minutes > 15 && minutes <= 20 {
				notExist := util.IsRepeateds(in20Minutes, bill.BillNum)
				if notExist {
					in20Minutes = append(in20Minutes, bill.BillNum)
				}
			}
		}
	}

	////////////////////////循环

	//返回的数据
	contractBillMap := map[string]interface{}{
		"not30ReturnMap":    notReturn,
		"not30ReturnMapLen": len(notReturn),
		"in5MinutesMap":     in5Minutes,
		"in5MinutesMapLen":  len(in5Minutes),
		"in15MinutesMap":    in15Minutes,
		"in15MinutesMapLen": len(in15Minutes),
		"in20MinutesMap":    in20Minutes,
		"in20MinutesMapLen": len(in20Minutes),
	}

	return err, contractBillMap
}

// 首次先把项目编码存入redis
// 之后查询出来数据 在把数据存入redis
func SetTimeLinessBriefing(proCode string) (err error) {
	setBriefingCode := global.GRedis.Set("ageing_briefing:code", proCode, 3600*time.Second)
	fmt.Println("=============set", setBriefingCode)
	return err
}
