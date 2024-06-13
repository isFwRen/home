package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"os"
	"server/global"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

func CreateQualityTable() {
	fmt.Println(1111111)
	err := global.GDb.Migrator().CreateTable(&model.Quality{})
	//err = global.GDb.Migrator().CreateTable(&model.SysCorrected{})
	//err = db.Migrator().CreateTable(&model.OutputStatisticsSummary{}, &model.Op1{}, &model.Op2{}, &model.Op0{}, &model.OpQ{}, &model.OutputStatistics{})
	fmt.Println(err)
}

func GetQualityManagement(info request.QualityRequest) (err error, list interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.GDb.Model(&model.Quality{})
	proCode := strings.Split(info.ProCode, ",")

	if len(proCode) > 1 {
		db = db.Where("pro_code in ? ", proCode)
	} else {
		db = db.Where("pro_code = ? ", info.ProCode)
	}

	if info.BillName != "" {
		db = db.Where("bill_name LIKE ? ", "%"+info.BillName+"%")
	}
	if info.WrongFieldName != "" {
		db = db.Where("wrong_field_name LIKE ? ", "%"+info.WrongFieldName+"%")
	}

	if info.ResponsibleName != "" {
		db = db.Where("op0_responsible_name LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("op1_responsible_name LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("op2_responsible_name LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("opq_responsible_name LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("op0_responsible_code LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("op1_responsible_code LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("op2_responsible_code LIKE ? ", "%"+info.ResponsibleName+"%").
			Or("opq_responsible_code LIKE ? ", "%"+info.ResponsibleName+"%")
	}

	//StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	//EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	if info.Month != "" {
		db = db.Where("month = ? ", info.Month)
	}

	var lr []model.Quality
	err = db.Count(&total).Error
	if err != nil {
		return nil, nil, total
	}
	err = db.Limit(limit).Offset(offset).Find(&lr).Error
	if err != nil {
		return nil, nil, total
	}

	return nil, lr, total
}

func AddQualityManagement(info request.QualitiesAeRequest) error {
	var total int64
	err := global.GDb.Debug().Model(&model.Quality{}).Where("pro_code = ? AND bill_name = ? AND wrong_field_name = ? ", info.ProCode, info.BillName, info.WrongFieldName).
		Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("存在重复数据")
	}
	var newData model.Quality
	err = copier.Copy(&newData, &info)
	if err != nil {
		return err
	}
	feedbackDate, _ := time.ParseInLocation("2006-01-02", info.FeedbackDate, time.Local)
	newData.FeedbackDate = feedbackDate
	entryDate, _ := time.ParseInLocation("2006-01-02", info.EntryDate, time.Local)
	newData.EntryDate = entryDate
	err = global.GDb.Debug().Model(&model.Quality{}).Create(&newData).Error
	return err
}

func EditQualityManagement(info request.QualitiesAeRequest) error {
	Update := make(map[string]interface{})
	utils.CheckRequired(info, Update)
	//arr := pq.StringArray{"aaaa"}
	err := global.GDb.Debug().Model(&model.Quality{}).Where("id = ? ", info.Id).Updates(&Update).Error
	return err
}

func UploadQualityManagement(dst map[string]string) error {
	failMsg := ""
	for fn, path := range dst {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			failMsg += path + "不存在;"
			continue
		}
		xlsx, err := excelize.OpenFile(path + "/" + fn)
		if err != nil {
			global.GLog.Error("open excel error:" + err.Error())
			continue
		}

		rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))

		if !CheckQualityTable(rows[0]) {
			failMsg += path + "表头不符合规则;"
		}

		var proQuality []model.Quality

		for i, row := range rows {
			var proQualityitem model.Quality
			if i == 0 {
				continue
			}
			fmt.Println("aaa", row)
			proQualityitem.Month = row[0]
			proQualityitem.ProCode = row[1]
			proQualityitem.BillName = row[2]
			proQualityitem.FeedbackDate = excelDateToDate(row[3])
			proQualityitem.EntryDate = excelDateToDate(row[4])
			proQualityitem.WrongFieldName = row[5]
			proQualityitem.Right = row[6]
			proQualityitem.Wrong = row[7]
			proQualityitem.Op0ResponsibleCode = row[8]
			proQualityitem.Op0ResponsibleName = row[9]
			proQualityitem.Op1ResponsibleCode = row[10]
			proQualityitem.Op1ResponsibleName = row[11]
			proQualityitem.Op2ResponsibleCode = row[12]
			proQualityitem.Op2ResponsibleName = row[13]
			proQualityitem.OpqResponsibleCode = row[14]
			proQualityitem.OpqResponsibleName = row[15]

			_, b := xlsx.GetPicture("Sheet1", "Q"+strconv.Itoa(i+1))
			if err := ioutil.WriteFile(path+"/"+proQualityitem.ProCode+"-"+proQualityitem.BillName+"-"+proQualityitem.WrongFieldName+".png", b, 0644); err != nil {
				fmt.Println(err)
			}
			proQualityitem.ImagePath = append(proQualityitem.ImagePath, path+"/"+proQualityitem.ProCode+"-"+proQualityitem.BillName+"-"+proQualityitem.WrongFieldName+".png")
			proQuality = append(proQuality, proQualityitem)
		}
		//err = db.Model(&model.Quality{}).Create(proQuality).Error
		//err = global.GDb.Model(&model.Quality{}).Create(&proQuality).Error
		//fmt.Println(err)
		err = global.GDb.Debug().Model(&model.Quality{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pro_code"}, {Name: "bill_name"}, {Name: "wrong_field_name"}},
			UpdateAll: true,
		}).Create(&proQuality).Error
		if err != nil {
			failMsg += err.Error() + ";"
			continue
		}
	}
	return errors.New(failMsg)
}

func ExportQualityData(proCode []string, month string) (error, string) {
	db := global.GDb.Model(&model.Quality{})
	for _, v := range proCode {
		if v == "all" {
			proCode = []string{"all"}
		}
	}
	var total int64
	db = db.Where("month = ? ", month).Count(&total)
	fmt.Println("total", total)
	var proQual []model.Quality
	for _, v := range proCode {
		if v != "all" {
			db = db.Where("pro_code = ? ", v)
			var proQualitem []model.Quality
			var total int64
			err := db.Find(&proQualitem).Count(&total).Error
			if err != nil {
				return err, ""
			}
			proQual = append(proQual, proQualitem...)
		} else {
			var proQualitem []model.Quality
			var total int64
			err := db.Find(&proQualitem).Count(&total).Error
			if err != nil {
				return err, ""
			}
			proQual = append(proQual, proQualitem...)
		}
	}
	basicPath := global.GConfig.LocalUpload.FilePath + "质量管理导出/"
	bookName := "理赔差错汇总" + strings.Replace(month, "-", "", -1) + ".xlsx"
	err := utils.ExportBigExcel3(basicPath, bookName, "", proQual)

	//插入图片
	xlsx, err := excelize.OpenFile(basicPath + bookName)
	if err != nil {
		return err, basicPath + bookName
	}
	for i, v := range proQual {
		fmt.Println(v.ImagePath[0])
		_, err := os.Stat(v.ImagePath[0])
		if os.IsNotExist(err) {
			continue
		}
		file, err := ioutil.ReadFile(v.ImagePath[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Q" + strconv.Itoa(i+2))
		xlsx.SetCellValue("Sheet1", "R"+strconv.Itoa(i+2), " ")
		if err := xlsx.AddPictureFromBytes("Sheet1", "R"+strconv.Itoa(i+2), "", v.ProCode+"-"+v.WrongFieldName, ".png", file); err != nil {
			fmt.Println(err)
		}
	}
	if err := xlsx.SaveAs(basicPath + bookName); err != nil {
		fmt.Println(err)
	}
	return err, basicPath + bookName
}

func DeleteQualityData(ids []string) (err error) {
	var Quality model.Quality
	for _, id := range ids {
		err = global.GDb.Where("id = ?", id).Delete(&Quality).Error
		if err != nil {
			return err
		}
	}
	return err
}

func CheckQualityTable(row []string) bool {
	topName := []string{"月份", "项目编码", "案件号", "反馈日期", "录入日期", "错误字段", "正确值", "错误值", "初审责任人工号", "初审责任人姓名", "一码责任人工号", "一码责任人姓名", "二码责任人工号", "二码责任人姓名", "问题件责任人工号", "问题件责任人姓名"}
	for i, a := range topName {
		for j, b := range row {
			if i == j {
				if a != b {
					return false
				}
			}
		}
	}
	return true
}

func excelDateToDate(excelDate string) time.Time {
	excelTime := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	var days, _ = strconv.Atoi(excelDate)
	return excelTime.Add(time.Second * time.Duration(days*86400))
}
