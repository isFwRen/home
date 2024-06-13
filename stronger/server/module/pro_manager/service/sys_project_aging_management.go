package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/global"
	"server/module/pro_manager/const_data"
	billmap "server/module/pro_manager/const_data"
	M "server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/module/pro_manager/project/B0108"
	"server/module/pro_manager/project/B0113"
	"server/module/pro_manager/project/B0118"
	"strings"
	"time"

	"go.uber.org/zap"
)

func GetProjectAgingAll(info request.ProjectAgingManagementSearchAll) (err error, list interface{}, total int64) {
	var ProjectAgingManagement []M.ProjectAgingManagement
	errMsg := ""
	arr := strings.Split(info.ProCode, ",")
	for _, proCode := range arr {
		infos := request.ProjectAgingManagementSearch{
			ProCode:    proCode,
			StartTime:  info.StartTime,
			EndTime:    info.EndTime,
			CaseNumber: info.CaseNumber,
			Agency:     info.Agency,
			CaseStatus: info.CaseStatus,
			Stage:      info.Stage,
			PageInfo:   info.PageInfo,
		}
		err, p, t := GetProjectAging(infos)
		if err != nil {
			errMsg += err.Error() + ";"
			continue
		}
		total += t
		ProjectAgingManagement = append(ProjectAgingManagement, p...)
	}
	return errors.New(errMsg), ProjectAgingManagement, total
}

func GetProjectAging(info request.ProjectAgingManagementSearch) (err error, list []M.ProjectAgingManagement, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	var bills []M.ProjectBill
	//连接数据库
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	db = db.Model(&M.ProjectBill{})
	if info.StartTime != "" {
		StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
		db = db.Where("scan_at >= ? ", StartTime)
	}
	if info.EndTime != "" {
		EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
		db = db.Where("scan_at <= ? ", EndTime)
	}
	if info.CaseNumber != "" {
		db = db.Where("bill_num LIKE ? ", "%"+info.CaseNumber+"%")
	}
	if info.Agency != "" {
		db = db.Where("agency = ? ", info.Agency)
	}
	if info.CaseStatus != "" {
		db = db.Where("status = ? ", info.CaseStatus)
	}
	if info.Stage != 0 {
		db = db.Where("stage = ? ", info.Stage)
	} else {
		//if info.ProCode == "B0108" {
		//CSB0108RC0014000 时效管理：显示所有录入状态不为“已接收”的案件
		//2023年05月08日20:07:01
		//改为不为已接收不为已回传
		//db = db.Where("stage <> ? ", 7)
		//} else {
		//}
		//2023年08月21日16:52:14 文雅提的需求
		if info.ProCode == "B0108" {
			db = db.Where("stage <> ? ", 7)
		} else {
			db = db.Where("stage <> ? and stage <> ? ", 5, 7)
		}
	}

	var orderBys [][]string
	err = json.Unmarshal([]byte(info.PageInfo.OrderBy), &orderBys)
	global.GLog.Error("", zap.Error(err))
	if err == nil {
		for _, orderBy := range orderBys {
			global.GLog.Info("", zap.Any("", orderBy))
			fmt.Println("========const_data.OrderBy[orderBy[0]]======", const_data.OrderBy[orderBy[0]])
			//db = db.Order(orderBy[0] + " " + orderBy[1])
		}
	}

	var ProjectAgingManagement []M.ProjectAgingManagement
	var projectAgingManagementItem M.ProjectAgingManagement
	err = db.Count(&total).Error
	if total == 0 {
		return nil, make([]M.ProjectAgingManagement, 0), 0
	}
	if err != nil {
		return err, make([]M.ProjectAgingManagement, 0), 0
	}
	//err = db.Order("id desc").Limit(limit).Offset(offset).Find(&bills).Error
	//err = db.Order("scan_at desc").Limit(limit).Offset(offset).Find(&bills).Error
	if len(orderBys) > 0 {
		if orderBys[0][1] == "asc" {
			err = db.Order("scan_at asc").Limit(limit).Offset(offset).Find(&bills).Error

		} else if orderBys[0][1] == "desc" {
			err = db.Order("scan_at desc").Limit(limit).Offset(offset).Find(&bills).Error

		}
	} else {
		err = db.Order("scan_at desc").Limit(limit).Offset(offset).Find(&bills).Error
	}
	if err != nil {
		return err, make([]M.ProjectAgingManagement, 0), 0
	}
	for _, bill := range bills {
		projectAgingManagementItem.ProCode = info.ProCode
		projectAgingManagementItem.CreatedAt = bill.CreatedAt.Format("2006-01-02 15:04:05") //扫案件创建时间
		projectAgingManagementItem.ScanAt = bill.ScanAt.Format("2006-01-02 15:04:05")       //扫描时间
		projectAgingManagementItem.Agency = bill.Agency
		projectAgingManagementItem.Stage = billmap.BillStage[bill.Stage]
		projectAgingManagementItem.Status = billmap.BillStatus[bill.Status]
		projectAgingManagementItem.CaseNumber = bill.BillNum
		//global.GLog.Info("", zap.Any(info.ProCode, bill))

		backAtTheLatest := ""
		timeRemaining := ""
		second := 0.0
		switch info.ProCode {
		case "B0108":
			backAtTheLatest, timeRemaining, second = B0108.CalculateBackTimeAndRemainder(bill)
		case "B0113":
			zipNameArr := strings.Split(bill.BillName, "_")
			if len(zipNameArr) != 6 {
				global.GLog.Error("案件号压缩包名字有误", zap.Any("", bill.BillName))
			}
			if zipNameArr[len(zipNameArr)-1] == "1" {
				//CSB0113RC0081000
				backAtTheLatest, timeRemaining, second = B0113.CalculateBackTimeAndRemainder(bill)
			} else {
				err, backAtTheLatest, timeRemaining, second = B0118.CalculateBackTimeAndRemainder(bill, total, info.ProCode)
				if err != nil {
					global.GLog.Error("", zap.Error(err))
				}
			}
		default:
			err, backAtTheLatest, timeRemaining, second = B0118.CalculateBackTimeAndRemainder(bill, total, info.ProCode)
			if err != nil {
				global.GLog.Error("", zap.Error(err))
			}
		}
		projectAgingManagementItem.BackAtTheLatest = backAtTheLatest
		projectAgingManagementItem.TimeRemaining = timeRemaining
		projectAgingManagementItem.Second = second
		projectAgingManagementItem.ClaimType = bill.ClaimType
		ProjectAgingManagement = append(ProjectAgingManagement, projectAgingManagementItem)
	}
	fmt.Println("GetProjectAging", len(ProjectAgingManagement))
	return err, ProjectAgingManagement, total
}
