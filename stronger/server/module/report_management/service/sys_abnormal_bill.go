package service

import (
	"fmt"
	"server/global"
	"server/module/pro_manager/const_data"
	m "server/module/pro_manager/model"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	u "server/utils"
	"strconv"
	"strings"
	"time"
)

func GetAbnormalBill(info request.AbnormalBillSearch) (err error, list interface{}, total int64) {
	abnormalBill := make([]model.AbnormalBill, 0)
	var bills []m.ProjectBill
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	ProCode := info.ProCode
	//连接数据库
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	db = db.Model(&m.ProjectBill{}).Where("created_at >= ? AND created_at <= ? AND status = ? AND stage = ? ", StartTime, EndTime, 3, 7)
	if info.Types != "" {
		Types, err := strconv.Atoi(info.Types)
		if err != nil {
			return err, nil, 0
		}
		db = db.Where("claim_type = ?", Types)
	}
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&bills).Error
		if err != nil {
			return err, nil, 0
		}
		var item model.AbnormalBill
		for _, v := range bills {
			item.CreateAt = v.CreatedAt.Format("2006-01-02 15:04:05")
			item.UploadAt = v.UploadAt.Format("2006-01-02 15:04:05")
			item.Agency = v.Agency
			item.AbnormalReason = strings.Split(u.RegReplace(v.DelRemarks, ".*：", ""), "-")[0]
			item.BillName = v.BillName
			item.Stage = v.Stage
			item.StageStr = const_data.BillStage[v.Stage]
			abnormalBill = append(abnormalBill, item)
		}
		return err, abnormalBill, total
	}
	return err, nil, total
}

func ExportAbnormalBill(info request.AbnormalBillSearch) (err error, path, name string) {
	err, list, total := GetAbnormalBill(info)
	if err != nil {
		return err, "", ""
	}
	fmt.Println("ExportAbnormalBill row :", total)
	s := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	e := strings.Replace(info.EndTime, " 23:59:59", "", -1)
	bookName := "异常件数据" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "异常件数据导出/" + info.ProCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcel(basicPath, bookName, "", list)
	return err, basicPath + bookName, bookName
}
