package service

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
	"time"
)

// GetBillInfo 获取history单据信息
func GetDelBills(billListSearch model.BillListSearch) (err error, total int64, list []model.ProjectDelBill) {
	limit := billListSearch.PageSize
	offset := billListSearch.PageSize * (billListSearch.PageIndex - 1)
	var projectBills []model.ProjectDelBill

	db := global.ProDbMap[billListSearch.ProCode]
	if db == nil {
		return global.ProDbErr, 0, projectBills
	}

	rojectCache := global.GProConf[billListSearch.ProCode]
	timeEnd := time.Now().Add(-24 * time.Duration(rojectCache.SaveDate) * time.Hour)
	if billListSearch.TimeEnd.After(timeEnd) {
		billListSearch.TimeEnd = timeEnd
	}

	db = db.Model(&model.ProjectDelBill{}).Where("scan_at BETWEEN ? AND ? ", billListSearch.TimeStart, billListSearch.TimeEnd)

	err = db.Count(&total).Error
	if limit >= 0 {
		err = db.Order("scan_at").Limit(limit).Offset(offset).Find(&projectBills).Error
	} else {
		err = db.Order("scan_at").Find(&projectBills).Error
	}

	return err, total, projectBills
	//获取history单据信息
	// err = db.Model(&model.ProjectDelBill{}).Where("id = ?", reqParam.ID).Find(&b).Error
	// return err, b
}

func ExportExcel(billListSearch model.BillListSearch) (err error, path, name string) {
	billListSearch.PageSize = -1
	err, _, projectBills := GetDelBills(billListSearch)
	if err != nil {
		return err, "", ""
	}

	projectBillExcels := []model.ProjectDelBillExcel{}
	for ii, bill := range projectBills {
		billExcel := model.ProjectDelBillExcel{}
		billExcel.Idx = ii + 1
		billExcel.BillNum = bill.BillNum
		billExcel.Size = utils.ToString(bill.Size) + "KB"
		billExcel.DelAt = bill.DelAt
		billExcel.ScanAt = bill.ScanAt
		billExcel.Remarks = bill.Remarks
		billExcel.Stage = bill.Stage
		billExcel.Name = bill.Name
		billExcel.Describe = bill.Describe
		billExcel.Image = bill.Image
		projectBillExcels = append(projectBillExcels, billExcel)
	}

	s := strings.Replace(billListSearch.TimeStart.Format("2006-01-02"), " 00:00:00", "", -1)
	e := strings.Replace(billListSearch.TimeEnd.Format("2006-01-02"), " 23:59:59", "", -1)
	bookName := "销毁报告" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "销毁报告/" + billListSearch.ProCode + "/"
	// 本地
	//basicPath := "./"
	err = utils.ExportBigExcel(basicPath, bookName, "", projectBillExcels)
	return err, basicPath, bookName
}

// GetBillInfo 获取history单据信息
func SaveDelBills(bill model.ProjectBill, stage string) (err error) {
	db := global.ProDbMap[bill.ProCode]
	if db == nil {
		return global.ProDbErr
	}

	rojectCache := global.GProConf[bill.ProCode]

	projectDelBill := model.ProjectDelBill{}
	projectDelBill.BillName = bill.BillName
	projectDelBill.BillNum = bill.BillNum
	projectDelBill.ScanAt = bill.ScanAt
	// projectDelBill.DelAt = bill.CreatedAt
	projectDelBill.Name = "系统"
	projectDelBill.Stage = stage
	projectDelBill.Remarks = bill.SaleChannel

	if rojectCache.SaveDate != 0 {
		projectDelBill.DelAt = bill.ScanAt.Add(24 * time.Duration(rojectCache.SaveDate) * time.Hour)
	}

	if stage == "下载文件" {
		path := global.GConfig.LocalUpload.FilePath + strings.Replace(bill.DownloadPath, "files/", "", 1)
		xmlCmd := "ls " + path + " |grep xml"
		err, stdout, _ := project.ShellOut(xmlCmd)
		fmt.Println("---------------xmlCmd-----------------------", xmlCmd, err, stdout)
		if err == nil {
			projectDelBill.Describe = stdout
		}

		projectDelBill.Image = bill.BatchNum + ".zip"
		zipCmd := "du -s " + path + bill.BatchNum + ".zip"
		err, stdout, _ = project.ShellOut(zipCmd)
		fmt.Println("---------------下载文件-----------------------", zipCmd, err, stdout)
		if err == nil {
			lineReg := regexp.MustCompile(`([\d\.]+)\s+\S+`)
			arr := lineReg.FindStringSubmatch(stdout)
			fmt.Println("---------------arr-----------------------", arr)
			if len(arr) > 1 {
				projectDelBill.Size = utils.ParseFloat(arr[1])
			}

		} else {
			zipCmd := "du -s " + strings.Replace(path, bill.BatchNum+"/", "", 1) + bill.BatchNum + ".zip"
			err, stdout, _ := project.ShellOut(zipCmd)
			fmt.Println("---------------下载文件2-----------------------", zipCmd, err, stdout)
			if err == nil {
				lineReg := regexp.MustCompile(`([\d\.]+)\s+\S+`)
				arr := lineReg.FindStringSubmatch(stdout)
				fmt.Println("---------------arr-----------------------", arr)
				if len(arr) > 1 {
					projectDelBill.Size = utils.ParseFloat(arr[1])
				}
			} else {
				aaa := strings.Replace(path, bill.BatchNum+"/", "", 1)
				aaa = strings.Replace(aaa, bill.BillNum+"/", "", 1)
				zipCmd = "du -s " + aaa + bill.BatchNum + ".zip"
				err, stdout, _ := project.ShellOut(zipCmd)
				fmt.Println("---------------下载文件2-----------------------", zipCmd, err, stdout)
				if err == nil {
					lineReg := regexp.MustCompile(`([\d\.]+)\s+\S+`)
					arr := lineReg.FindStringSubmatch(stdout)
					fmt.Println("---------------arr-----------------------", arr)
					if len(arr) > 1 {
						projectDelBill.Size = utils.ParseFloat(arr[1])
					}
				}
			}
		}
	} else {
		path := global.GConfig.LocalUpload.FilePath + bill.ProCode + "/upload_xml/" +
			fmt.Sprintf("%v/%v/%v/",
				bill.CreatedAt.Year(), int(bill.CreatedAt.Month()),
				bill.CreatedAt.Day())
		zipCmd := "du -s " + path + bill.BillNum + ".xml"
		err, stdout, _ := project.ShellOut(zipCmd)
		fmt.Println("---------------结果文件-----------------------", zipCmd, err, stdout)
		if err == nil {
			lineReg := regexp.MustCompile(`([\d\.]+)\s+\S+`)
			arr := lineReg.FindStringSubmatch(stdout)
			fmt.Println("---------------arr-----------------------", arr)
			if len(arr) > 1 {
				projectDelBill.Size = utils.ParseFloat(arr[1])
			}
		}
		projectDelBill.Image = bill.BillNum + ".xml"
		projectDelBill.Describe = bill.BillNum + ".xml"
	}

	err = db.Model(&model.ProjectDelBill{}).Where("bill_name = ? and stage = ?", bill.BillName, stage).FirstOrCreate(&projectDelBill).Error

	//获取history单据信息
	// err = db.Model(&model.ProjectDelBill{}).Where("id = ?", reqParam.ID).Find(&b).Error
	return err
}

func HistoryDelList(proCode string) (err error, total int64) {
	// limit := billListSearch.PageSize
	// offset := billListSearch.PageSize * (billListSearch.PageIndex - 1)
	var projectBills []model.ProjectBill

	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, 0
	}

	db = db.Model(&model.ProjectBill{})

	err = db.Count(&total).Error
	if err != nil {
		return err, 0
	}

	offset := 0
	for offset*100 < int(total) {
		err = db.Order("created_at asc").Limit(100).Offset(offset * 100).Find(&projectBills).Error
		if err != nil {
			break
		}
		for _, projectBill := range projectBills {
			SaveDelBills(projectBill, "下载文件")
			SaveDelBills(projectBill, "结果数据")
		}
		offset = offset + 1
	}

	// if limit >= 0 {
	// 	err = db.Order("created_at asc").Limit(limit).Offset(offset).Find(&projectBills).Error
	// }

	//获取history单据信息
	return err, total
}
