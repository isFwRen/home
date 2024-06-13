package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/module/export"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	"server/module/report_management/model/response"
	u "server/utils"
	"strings"
	"time"
	model2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
)

var IsOperationLogMap = map[string]string{
	"0": "待审核",
	"1": "通过",
	"2": "不通过",
}

func GetIncorrectList(info model.WrongSearch) (err error, list interface{}, total int64) {
	var wrongs []model.Wrong
	var wrongGet []response.GetWrong
	var wrongGetItem response.GetWrong
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	ProCode := info.ProCode
	//连接数据库
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	err = db.Model(&model.Wrong{}).Find(&wrongs).Error
	db = db.Model(&model.Wrong{})
	if info.Code != "" {
		db = db.Where("code LIKE ?", "%"+info.Code+"%")
	}
	if info.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+info.NickName+"%")
	}
	if info.FieldName != "" {
		db = db.Where("field_name LIKE ?", "%"+info.FieldName+"%")
	}
	if info.Op != "" {
		db = db.Where("op = ?", info.Op)
	}
	if info.Complaint != "" {
		db = db.Where("is_complain = ?", info.Complaint)
	}
	if info.Confirm != "" {
		db = db.Where("is_wrong_confirm = ?", info.Confirm)
	}
	if !info.IsAudit {
		//待审核 = 非已审核 && 提交申诉
		db = db.Where("is_audit = ? and is_complain = true", info.IsAudit)
	}
	if info.StartTime != "" {
		StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
		db = db.Where("submit_day >= ?", StartTime)
	}
	if info.EndTime != "" {
		EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
		db = db.Where("submit_day <= ?", EndTime)
	}
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&wrongs).Error
		for _, wrong := range wrongs {
			u.StructAssign(&wrongGetItem, &wrong)
			wrongGetItem.ID = wrong.ID
			wrongGetItem.IsOperationLog = wrong.IsOperationLog
			wrongGetItem.IsAudit = wrong.IsAudit
			wrongGet = append(wrongGet, wrongGetItem)
		}
		return err, wrongGet, total
	}
	return err, nil, total
}

func GetExportIncorrectList(info model.WrongExport) (err error, wrongs []model.WrongExportResp) {
	//连接数据库
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr, wrongs
	}
	err = db.Model(&model.Wrong{}).
		Select("submit_day,code,nick_name,bill_name,bill_num,agency,field_name,wrong,\"right\"").
		Where("submit_day BETWEEN ? AND ?", info.StartTime, info.EndTime).
		Find(&wrongs).Error
	return err, wrongs
}

func GetIncorrectTaskList(info model.WrongSearch, uid string) (err error, list interface{}, total int64) {
	var wrongs []model.Wrong
	var wrongGet []response.GetWrong
	var wrongGetItem response.GetWrong
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	ProCode := info.ProCode

	var user model2.SysUser
	err = global.GUserDb.Model(&model2.SysUser{}).Where("id = ? ", uid).Find(&user).Error
	if err != nil {
		return err, nil, 0
	}

	//连接数据库
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	err = db.Model(&model.Wrong{}).Find(&wrongs).Error
	db = db.Model(&model.Wrong{})
	db = db.Where("code = ? AND nick_name = ? ", user.Code, user.Name)
	if info.FieldName != "" {
		db = db.Where("field_name LIKE ?", "%"+info.FieldName+"%")
	}
	if info.Op != "" {
		db = db.Where("op = ?", info.Op)
	}
	if info.Complaint != "" {
		db = db.Where("is_complain = ?", info.Complaint)
	}
	if info.StartTime != "" {
		StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
		db = db.Where("submit_day >= ?", StartTime)
	}
	if info.EndTime != "" {
		EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
		db = db.Where("submit_day <= ?", EndTime)
	}
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&wrongs).Error
		for _, wrong := range wrongs {
			u.StructAssign(&wrongGetItem, &wrong)
			wrongGetItem.ID = wrong.ID
			wrongGetItem.IsOperationLog = IsOperationLogMap[wrong.IsOperationLog]
			wrongGetItem.IsAudit = wrong.IsAudit
			wrongGet = append(wrongGet, wrongGetItem)
		}
		return err, wrongGet, total
	}
	return err, nil, total
}

func ComplainConfirm(confirm request.Complain) (list interface{}, err error) {
	var wrong model.Wrong
	//连接数据库
	db := global.ProDbMap[confirm.ProCode]
	if db == nil {
		return nil, global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//超过3天就不让申诉
	err = tx.Model(&model.Wrong{}).Where("id = ? ", confirm.Id).Find(&wrong).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if time.Now().Sub(wrong.SubmitDay).Hours() > 72 {
		tx.Rollback()
		return nil, errors.New("只允许在三天之内对该错误进行申诉")
	}
	if wrong.Wrong == "?" || wrong.Wrong == "？" {
		tx.Rollback()
		return nil, errors.New("纯问号内容，不计入正确率统计，无需进行申诉")
	}
	err = tx.Model(&model.Wrong{}).Where("id = ? ", confirm.Id).Updates(map[string]interface{}{
		"is_complain":      confirm.ComplainConfirm,
		"is_operation_log": "0",
	}).Find(&wrong).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return wrong, err
}

func ComplainConfirmTask(confirm request.ComplainTask) (list interface{}, err error) {
	var wrong model.Wrong
	//连接数据库
	db := global.ProDbMap[confirm.ProCode]
	if db == nil {
		return nil, global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, v := range confirm.List {
		//超过3天就不让申诉
		err = tx.Model(&model.Wrong{}).Where("id = ? ", v.Id).Find(&wrong).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if time.Now().Sub(wrong.SubmitDay).Hours() > 72 {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("日期：%s, 工号：%s ,姓名：%s, 单号：%s, 机构号：%s, 字段编号：%s, 字段：%s, 只允许在三天之内对该错误进行申诉", wrong.SubmitDay.Format("2006-01-02"),
				wrong.Code, wrong.NickName, wrong.BillName, wrong.Agency, wrong.FieldCode, wrong.FieldName))
		}
	}

	for _, v := range confirm.List {
		err = tx.Model(&model.Wrong{}).Where("id = ? ", v.Id).Updates(map[string]interface{}{
			"is_complain":      confirm.ComplainConfirm,
			"is_operation_log": "0",
		}).Find(&wrong).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return wrong, err
}

func WrongConfirm(confirm request.WrongConfirmArray) (list interface{}, err error) {

	var lists []model.Wrong
	//连接数据库
	db := global.ProDbMap[confirm.ProCode]
	if db == nil {
		return nil, global.ProDbErr
	}

	dayCodeMap := []string{}

	for _, v := range confirm.List {
		//判断有没有申诉
		var wrong model.Wrong
		err = db.Model(&model.Wrong{}).Where("id = ? ", v.Id).Find(&wrong).Error
		if err != nil {
			return nil, err
		}
		if !wrong.IsComplain {
			return nil, errors.New(fmt.Sprintf("日期: %s, 工号: %s, 错误内容: %s, 没有申诉", wrong.SubmitDay.Format("2006-01-02"), wrong.Code, wrong.Wrong))
		}
		// is_wrong_confirm := map[string]bool{
		// 	"is_wrong_confirm": true,
		// }
		OperationOp := make([]string, 0)
		if confirm.WrongConfirm == "1" {
			wrong.IsWrongConfirm = true
			wrong.IsOperationLog = "1"
			wrong.IsAudit = true
			wrong.Right = wrong.Wrong
			OperationOp = []string{"IsWrongConfirm", "IsOperationLog", "IsAudit", "Right"}
		} else {
			wrong.IsWrongConfirm = false
			wrong.IsOperationLog = "2"
			wrong.IsAudit = true
			wrong.IsComplain = false
			OperationOp = []string{"IsWrongConfirm", "IsOperationLog", "IsAudit", "IsComplain"}
		}
		// if confirm.WrongConfirm == "2" {
		// 	is_wrong_confirm = map[string]bool{
		// 		"is_wrong_confirm": false,
		// 	}
		// }
		err = db.Select(OperationOp).Model(&model.Wrong{}).Where("id = ? ", wrong.ID).Updates(&wrong).Error
		if err != nil {
			return nil, err
		}

		day := wrong.SubmitDay.Format("2006-01-02")
		dayCode := day + "_" + wrong.Code
		if arrays.Contains(dayCodeMap, dayCode) == -1 {
			dayCodeMap = append(dayCodeMap, dayCode)
		}

		//判断有没有操作过差错审核(也就是通过或者不通过按钮)
		//IsOperationLog:  true:表明上次修改差错审核为通过   false:表明上次修改差错审核为不通过
		// if confirm.WrongConfirm == "1" {
		// 	confirm.Right = wrong.Wrong
		// }
		// if confirm.WrongConfirm == "1" {
		// 	msg, err := ReWrongSum(wrong, confirm.ProCode, confirm.Right, true)
		// 	if err != nil {
		// 		return nil, errors.New(msg + "," + err.Error())
		// 	}
		// } else if confirm.WrongConfirm == "2" {
		// msg, err := ReWrongSum(wrong, confirm.ProCode, confirm.Right, false)
		// 	if err != nil {
		// 		return nil, errors.New(msg + "," + err.Error())
		// 	}
		// } else {
		// 	return nil, errors.New("审核失败")
		// }

		// wrong.Right = confirm.Right
		// if confirm.WrongConfirm == "1" {
		// 	wrong.IsWrongConfirm = true
		// } else if confirm.WrongConfirm == "2" {
		// 	wrong.IsWrongConfirm = false
		// }

		// var newWrong model.Wrong
		// err = db.Model(&model.Wrong{}).Where("id = ? ", v.Id).Find(&newWrong).Error
		// if err != nil {
		// 	return nil, err
		// }

		// lists = append(lists, newWrong)
	}

	err, tx := export.EffectiveCharacter(dayCodeMap, confirm.ProCode)

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return lists, err
	}

	err = export.ExpenseAccountSum(dayCodeMap, confirm.ProCode)
	if err != nil {
		return lists, err
	}

	return lists, err
}

func IncorrectAnalysis(info model.WrongSearch) (err error, list interface{}, total int64) {
	var wrongAnalysis []model.WrongAnalysis
	var wrongAnalysisGet []response.GetWrongAnalysis
	var wrongAnalysisGetItem response.GetWrongAnalysis
	//limit := info.PageInfo.PageSize
	//offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	ProCode := info.ProCode
	//连接数据库
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	db = db.Model(&model.WrongAnalysis{})
	if info.Code != "" {
		db = db.Where("code LIKE ? ", "%"+info.Code+"%")
	}
	if info.NickName != "" {
		db = db.Where("nick_name LIKE ? ", "%"+info.NickName+"%")
	}
	if info.StartTime != "" {
		StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
		db = db.Where("statistical_time >= ? ", StartTime)
	}
	if info.EndTime != "" {
		EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
		db = db.Where("statistical_time <= ? ", EndTime)
	}
	err = db.Count(&total).Error
	if total > 0 {
		//err = db.Order("created_at asc").Limit(limit).Offset(offset).Find(&wrongAnalysis).Error
		err = db.Order("created_at asc").Find(&wrongAnalysis).Error
		for _, wrong := range wrongAnalysis {
			u.StructAssign(&wrongAnalysisGetItem, &wrong)
			wrongAnalysisGet = append(wrongAnalysisGet, wrongAnalysisGetItem)
		}
		return err, wrongAnalysisGet, total
	}
	return err, nil, total
}

func OcrAnalysis(info model.WrongSearch) (err error, list interface{}, total int64) {
	var wrong []model.Wrong
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	ProCode := info.ProCode
	//连接数据库
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	db = db.Model(&model.Wrong{})
	if info.FieldName != "" {
		db = db.Where("field_name = ? ", info.FieldName)
	}
	if info.Op != "" {
		db = db.Where("field_name = ? ", info.Op)
	}
	if info.StartTime != "" {
		StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
		db = db.Where("submit_day >= ? ", StartTime)
	}
	if info.EndTime != "" {
		EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
		db = db.Where("submit_day <= ? ", EndTime)
	}
	db = db.Where("id_ocr = ? ", true)
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&wrong).Error
		ocrWrong := make([]model.OcrAnalysis, len(wrong))
		for i, v := range wrong {
			ocrWrong[i].StatisticalTime = v.SubmitDay
			ocrWrong[i].CaseNumber = v.CaseNumber
			ocrWrong[i].Agency = v.Agency
			ocrWrong[i].FieldName = v.FieldName
			ocrWrong[i].Wrong = v.Wrong
			ocrWrong[i].Right = v.Right
		}
		return err, ocrWrong, total
	}
	return err, nil, total
}

func ExportIncorrectAnalysis(info request.IncorrectAnalysisExport) (err error, path, name string) {
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	//查询明细
	//连接数据库
	var total int64
	ProCode := info.ProCode
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, "", ""
	}
	//直接查！
	db = db.Model(&model.WrongAnalysis{}).Where("statistical_time >= ? AND statistical_time <= ? ", StartTime, EndTime)
	var WrongAnalysisReturn []model.WrongAnalysis
	err = db.Find(&WrongAnalysisReturn).Count(&total).Error
	var WrongAnalysisExport []model.WrongAnalysisExport
	var WrongAnalysisExportItem model.WrongAnalysisExport
	for _, v := range WrongAnalysisReturn {
		WrongAnalysisExportItem.StatisticalTime = v.StatisticalTime
		WrongAnalysisExportItem.Code = v.Code
		WrongAnalysisExportItem.NickName = v.NickName
		WrongAnalysisExportItem.WrongNumber = v.WrongNumber
		WrongAnalysisExportItem.TheNumberOfComplaints = v.TheNumberOfComplaints
		WrongAnalysisExportItem.NonPassingQuantity = v.NonPassingQuantity
		WrongAnalysisExportItem.ThroughTheNumber = v.ThroughTheNumber

		WrongAnalysisExportItem.TheComplaintRate = C(v.TheComplaintRate)
		WrongAnalysisExportItem.ThePassRate = C(v.ThePassRate)
		WrongAnalysisExportItem.UnqualifiedRate = C(v.UnqualifiedRate)

		WrongAnalysisExport = append(WrongAnalysisExport, WrongAnalysisExportItem)
	}
	if err != nil {
		return err, "", ""
	}
	s := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	e := strings.Replace(info.EndTime, " 23:59:59", "", -1)
	bookName := "错误分析" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "错误分析导出/" + ProCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcel(basicPath, bookName, "", WrongAnalysisExport)
	return err, basicPath + bookName, bookName
}

func C(a float64) string {
	if a > 0.001 {
		r, _ := decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(100)).Float64()
		return fmt.Sprintf("%.2f", r) + "%"
	} else {
		return decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%"
	}
}
