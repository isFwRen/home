package service

import (
	"server/global"
	model3 "server/module/pro_manager/model"
	model2 "server/module/report_management/model"
	"server/module/report_management/model/request"
	"server/module/sys_base/model"
)

//func GetBusinessDetails(info request.GetBusinessDetailsSearch) (err error, list interface{}, total int64) {
//	switch info.ProCode {
//	case "B0118":
//		err, list, total = B0118.GetBusinessDetails(info)
//	case "B0108":
//		err, list, total = B0108.GetBusinessDetails(info)
//	default:
//		fmt.Println("123")
//	}
//	return
//}

// GetCharSumReportByPage 分页获取字符统计
func GetCharSumReportByPage(baseTimePageCode model.BaseTimePageCode) (err error, list []model2.CharSum, row int64) {
	db := global.GDb.Model(&model2.CharSum{}).Where("sum_date BETWEEN ? AND ?", baseTimePageCode.StartTime, baseTimePageCode.EndTime)
	_, ok := global.GProConf[baseTimePageCode.ProCode]
	if !ok {
		return global.ProDbErr, list, row
	}
	db = db.Where("pro_code = ?", baseTimePageCode.ProCode)
	db.Count(&row)
	if baseTimePageCode.PageSize > 0 && baseTimePageCode.PageIndex > -1 {
		limit := baseTimePageCode.PageSize
		offset := baseTimePageCode.PageSize * (baseTimePageCode.PageIndex - 1)
		db = db.Order("id desc").Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return err, list, row
}

// CountBill 统计单据
func CountBill(info request.BusinessDetailsSearch) (total int64) {
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return 0
	}
	db.Model(&model3.ProjectBill{}).
		Where("scan_at >= ? AND scan_at <= ? AND (stage = 5 or stage = 7) ", info.StartTime, info.EndTime).
		Count(&total)
	return total
}
