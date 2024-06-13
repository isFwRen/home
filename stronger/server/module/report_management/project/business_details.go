/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/12/2 16:50
 */

package project

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"server/global"
	l "server/module/load/model"
	"server/module/pro_manager/const_data"
	m "server/module/pro_manager/model"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	service2 "server/module/report_management/service"
	u "server/utils"
	"strconv"
	"strings"
	"time"
	"unicode"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

//报表内容说明：
//日期：匹配案件列表的日期，格式为：YYYY/MM/DD
//批次号：匹配案件列表的批次号
//销售渠道：匹配案件列表的销售渠道
//案件号：匹配案件列表的案件号
//影像数量：显示案件号对应的影像数量
//机构：匹配案件列表的机构
//医保类型：匹配案件列表的医保类型
//扫描时间：匹配案件列表的扫描时间，格式为：YYYY/MM/DD HH:MM:SS
//下载时间：匹配案件列表的日期+时间，格式为：YYYY/MM/DD HH:MM:SS
//导出时间：匹配案件列表的导出时间，格式为：YYYY/MM/DD HH:MM:SS
//初次回传时间：显示案件号第一次回传的时间，格式为：格式为：YYYY/MM/DD HH:MM:SS
//最新回传时间：匹配案件列表的回传时间，格式为：YYYY/MM/DD HH:MM:SS
//最晚回传时间：案件最晚的回传时间，格式为：YYYY-MM-DD hh:mm:ss，最晚回传时间计算方式：
//延迟时间：
//“初次回传时间”晚于“最晚回传时间”时，延迟时间=初次回传时间-最晚回传时间，格式为HH:MM:SS
//“初次回传时间”不晚于“最晚回传时间”，显示为空
//案件状态：匹配案件列表的案件状态
//是否超时：
//“初次回传时间”晚于“最晚回传时间”时，显示为“是”
//“初次回传时间”晚于“最晚回传时间”时，显示为“否”
//处理时长：处理时长=(初次回传时间-扫描时间)*24，格式为数值，保留2个小数点
//录入状态：匹配案件列表的录入状态
//理赔类型：匹配案件列表的理赔类型
//账单金额：匹配案件列表的账单金额
//发票数量：匹配案件列表的发票数量
//清单数量：计算案件清单的录入数量，计算方式为：不同影像中字段名称包含“清单所属发票”字样的个数，同一张影像中字段名称存在多个包含“清单所属发票”的字样时，属于1个；(例：初审时在A影像出现2个“清单所属发票”的字样，则清单数量为1)
//问题件数量：匹配案件列表的问题件数量
//疾病诊断：显示字段名称包含“疾病诊断”的结果数据，多个时，用分号“；”隔开
//录入字符数：显示案件录入的字符数，一个半角数字/字母/符号为一个字符，一个全角数字/字母/符号、中文为两个字符
//结算字符数：xml的字符数，一个半角数字/字母/符号为一个字符，一个全角数字/字母/符号、中文为两个字符(项目有特殊要求，对应项目则以最新的需求为准)

//结算金额：根据项目特殊定制
//时效考核要求：根据项目特殊定制

// 案件次数：计算同一案件号在30天内出现的次数，30天内出现多次时，第一次为1，第二次为2，以此类推；
// 初审进入时间：工序进入初审的时间，格式为：YYYY/MM/DD HH:MM:SS
// 初审结束时间：工序结束初审的时间，格式为：YYYY/MM/DD HH:MM:SS
// 初审处理时长：初审处理时长=(初审结束时间-初审进入时间)*24，格式为数值，保留2个小数点
// 一码进入时间：工序进入一码的时间，格式为：YYYY/MM/DD HH:MM:SS
// 一码结束时间：工序结束一码的时间，格式为：YYYY/MM/DD HH:MM:SS
// 一码处理时长：一码处理时长=(一码结束时间-一码进入时间)*24，格式为数值，保留2个小数点
// 二码进入时间：工序进入二码的时间，格式为：YYYY/MM/DD HH:MM:SS
// 二码结束时间：工序结束二码的时间，格式为：YYYY/MM/DD HH:MM:SS
// 二码处理时长：二码处理时长=(二码结束时间-二码进入时间)*24，格式为数值，保留2个小数点
// 问题件进入时间：工序进入问题件的时间，格式为：YYYY/MM/DD HH:MM:SS
// 问题件结束时间：工序结束问题件的时间，格式为：YYYY/MM/DD HH:MM:SS
// 问题件处理时长：问题件处理时长=(问题件结束时间-问题件进入时间)*24，格式为数值，保留2个小数点
// 质检进入时间：案件首次导出的时间，格式为：YYYY/MM/DD HH:MM:SS
// 质检结束时间：案件首次回传的时间，格式为：YYYY/MM/DD HH:MM:SS
// 质检处理时长：质检处理时长=(质检结束时间-质检进入时间)*24，格式为：YYYY/MM/DD HH:MM:SS
type billObj struct {
	m.ProjectBill
	ProjectBlock []l.ProjectBlock `json:"projectBlock" gorm:"foreignKey:BillID;references:ID;"`
	ProjectField []l.ProjectField `json:"projectField" gorm:"foreignKey:BillID;references:ID;"`
}

func (v billObj) TableName() string {
	return "project_bills"
}

func GetBusinessDetails(info request.BusinessDetailsSearch) (err error, list []model.BusinessDetailsExport, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	proCode := info.ProCode
	if len(info.StartTime) == 10 {
		info.StartTime = info.StartTime + " 00:00:00"
	}
	if len(info.EndTime) == 10 {
		info.EndTime = info.EndTime + " 23:59:59"
	}
	//连接数据库
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, list, total
	}
	businessDetails := make([]model.BusinessDetailsExport, 0)
	var bills []billObj
	//原来的可以用------------------
	//err = db.Where("scan_at >= ? AND scan_at <= ? AND (stage = 5 or stage = 7) ", info.StartTime, info.EndTime).
	//	//Select("id", "bill_name").
	//	Preload("ProjectBlock").
	//	Preload("ProjectField", func(db *gorm.DB) *gorm.DB {
	//		return db.Select("id,name,code,bill_id,result_value,result_input,final_value,final_input")
	//	}).
	//	Limit(limit).Offset(offset).Find(&bills).Error
	//原来的可以用------------------

	//////////////////////////年报
	if info.Type == 4 {
		//////////////////////////年报
		year, err := time.Parse("2006-01-02 15:04:05", info.StartTime)
		if err != nil {
			return err, list, total
		}
		yearStr := strconv.Itoa(year.Year())
		err = db.Where("scan_at >= ? AND scan_at <= ? AND (stage = 5 or stage = 7) ", yearStr+"-01-01 00:00:00", yearStr+"-12-31 23:59:59").
			//Select("id", "bill_name").
			Preload("ProjectBlock").
			Preload("ProjectField", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,name,code,bill_id,block_id,result_value,result_input,final_value,final_input")
			}).
			Order("created_at").Limit(limit).Offset(offset).Find(&bills).Error
		//////////////////////////年报
	} else {
		err = db.Where("scan_at >= ? AND scan_at <= ? AND (stage = 5 or stage = 7) ", info.StartTime, info.EndTime).
			//Select("id", "bill_name").
			Preload("ProjectBlock").
			Preload("ProjectField", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,name,code,bill_id,block_id,result_value,result_input,final_value,final_input")
			}).
			Order("created_at").Limit(limit).Offset(offset).Find(&bills).Error
	}
	//////////////////////////年报

	if err != nil {
		return err, list, total
	}
	db.Model(&m.ProjectBill{}).Where("scan_at >= ? AND scan_at <= ? AND (stage = 5 or stage = 7) ", info.StartTime, info.EndTime).Count(&total)
	var item model.BusinessDetailsExport
	for _, obj := range bills {
		bill := obj.ProjectBill
		blocks := obj.ProjectBlock
		fields := obj.ProjectField
		item.CreateAt = bill.CreatedAt.Format("2006-01-02") //日期
		item.BatchNum = bill.BatchNum                       //批次号
		item.SaleChannel = bill.SaleChannel                 //销售渠道
		item.BillNum = bill.BillNum                         //案件号
		item.PictureNumber = len(bill.Images)               //影像数量
		item.Agency = bill.Agency                           //机构
		item.InsuranceType = bill.InsuranceType             //医保类型
		//item.ScanAt = bill.ScanAt.Format("2006-01-02 15:04:05")         //扫描时间   修改为创建时间
		item.ScanAt = bill.CreatedAt.Format("2006-01-02 15:04:05")      //扫描时间
		item.DownloadAt = bill.DownloadAt.Format("2006-01-02 15:04:05") //下载时间
		item.ExportAt = bill.ExportAt.Format("2006-01-02 15:04:05")     //导出时间
		item.Stage = const_data.BillStage[bill.Stage]                   //录入状态
		item.ClaimType = const_data.BillClaimType[bill.ClaimType]       //理赔类型
		item.CountMoney = bill.CountMoney                               //账单金额
		item.InvoiceNum = bill.InvoiceNum                               //发票数量
		item.Status = const_data.BillStatus[bill.Status]                //案件状态
		item.QuestionNum = bill.QuestionNum                             //问题件数量
		item.QualityUserCode = bill.QualityUserCode                     //质检人员工号
		item.QualityUserName = bill.QualityUserName                     //质检人员姓名
		item.Remark = bill.Remark                                       //备注
		item.ExportStage = const_data.BillStage[bill.ExportStage]       //导出时候的状态
		FirstExportAt := bill.FirstExportAt.Format("2006-01-02 15:04:05")
		item.FirstExportAt = FirstExportAt                                    //初次导出时间
		item.FirstUploadAt = bill.UploadAt.Format("2006-01-02 15:04:05")      //初次回传时间
		item.LatestUploadAt = bill.LastUploadAt.Format("2006-01-02 15:04:05") //最新回传时间

		//err, backAtTheLatest := FindAndCalculateBackAtTheLatestAdapter(bill, proCode)
		err, backAtTheLatest := CalculateAllReturnTime(bill, proCode)
		if err != nil {
			global.GLog.Error("backAtTheLatest", zap.Error(err))
			backAtTheLatest = ""
		}
		item.AtTheLateUploadAt = backAtTheLatest //最晚回传时间
		if backAtTheLatest != "" {
			t1, _ := time.ParseInLocation("2006-01-02 15:04:05", backAtTheLatest, time.Local)
			//是否超时
			if t1.Before(bill.UploadAt) {
				item.IsTheTimeOut = "是"
			} else {
				item.IsTheTimeOut = "否"
			}
			if bill.UploadAt.Before(t1) {
				item.LateTime = "" //延迟时间
			} else {
				//item.LateTime = bill.UploadAt.Sub(t1).Round(time.Second).String() //延迟时间
				//delayTime := formatTimeDuration(bill.UploadAt.Sub(t1).Round(time.Second))
				format := ObtainDuration(t1, bill.UploadAt) //延迟时间
				item.LateTime = format                      //延迟时间
			}

		} else {
			item.IsTheTimeOut = ""
			item.LateTime = "0.0.0" //延迟时间
		}
		//item.WorkTime = bill.LastUploadAt.Sub(bill.ScanAt).Round(time.Second).String()                     //处理时长
		hanldTime := formatTimeDuration(bill.LastUploadAt.Sub(bill.ScanAt).Round(time.Second))
		fmt.Println(hanldTime)
		//item.WorkTime = formatTimeDuration(bill.LastUploadAt.Sub(bill.ScanAt).Round(time.Second))
		//WorkTime := ConvertFormatTime(result) //处理时长
		WorkTime := ObtainDuration(bill.ScanAt, bill.LastUploadAt)                                         //处理时长
		item.WorkTime = WorkTime                                                                           //处理时长
		item.DurationRange = GetDurationRange(item.WorkTime)                                               //时长范围
		item.DiseaseDiagnosis = GetJiBingFiled(fields)                                                     //疾病诊断
		item.SettlementMoney = CalculateMoney(bill, bill.LastUploadAt.Sub(bill.ScanAt).Round(time.Second)) //结算金额?

		//item.RequirementOfAging = B0108.CalculateRequirementOfAging(bill) //时效考核要求
		err, requirementOfAging := CalculateRequirementOfAgingAdapter(bill, proCode) //时效考核要求
		if err != nil {
			global.GLog.Error("requirementOfAging", zap.Error(err))
			requirementOfAging = ""
		}
		item.RequirementOfAging = requirementOfAging
		item.ListingNum = CalculateListingNum(blocks) //清单数量
		//item.TheNumberOfCase = CalculateTheNumberOfCase(bill)  //案件次数
		item.FieldCharacter = CalculateWriteCharacter(fields)  //录入字符数
		item.SettlementCharacter = CalculateXmlCharacter(bill) //结算字符数

		//var blocks []l.ProjectBlock
		//err = db.Model(&l.ProjectBlock{}).Where("bill_id = ? ", bill.ID).Find(&blocks).Error
		//if err != nil {
		//	continue
		//}

		op0BlockTime := make([]time.Time, 0)
		op1BlockTime := make([]time.Time, 0)
		op2BlockTime := make([]time.Time, 0)
		opQBlockTime := make([]time.Time, 0)
		for _, block := range blocks {
			op0BlockTime = append(op0BlockTime, block.Op0ApplyAt, block.Op0SubmitAt)
			op1BlockTime = append(op1BlockTime, block.Op1ApplyAt, block.Op1SubmitAt)
			op2BlockTime = append(op2BlockTime, block.Op2ApplyAt, block.Op2SubmitAt)
			opQBlockTime = append(opQBlockTime, block.OpqApplyAt, block.OpqSubmitAt)
		}
		op0MaxTime, op0MinTime := Compare(op0BlockTime)
		op1MaxTime, op1MinTime := Compare(op1BlockTime)
		op2MaxTime, op2MinTime := Compare(op2BlockTime)
		opQMaxTime, opQMinTime := Compare(opQBlockTime)
		item.Op0Entry = op0MinTime.Format("2006-01-02 15:04:05") //初审进入时间
		item.Op0End = op0MaxTime.Format("2006-01-02 15:04:05")   //初审结束时间
		preliminaryTime := Calculate(op0MinTime, op0MaxTime)
		fmt.Println("preliminaryTime=", preliminaryTime)
		Op0EntryTime := ObtainDuration(op0MinTime, op0MaxTime) //初审处理时间
		item.Op0WorkTime = Op0EntryTime                        //初审处理时间

		item.Op1Entry = op1MinTime.Format("2006-01-02 15:04:05") //一码进入时间
		item.Op1End = op1MaxTime.Format("2006-01-02 15:04:05")   //一码结束时间
		oneCode := Calculate(op1MinTime, op1MaxTime)
		fmt.Println("oneCode=", oneCode)
		Op1WorkTime := ObtainDuration(op1MinTime, op1MaxTime) //一码处理时间
		item.Op1WorkTime = Op1WorkTime                        //一码处理时间

		item.Op2Entry = op2MinTime.Format("2006-01-02 15:04:05") //二码进入时间
		item.Op2End = op2MaxTime.Format("2006-01-02 15:04:05")   //二码结束时间
		calculate := Calculate(op2MinTime, op2MaxTime)
		fmt.Println("calculate=", calculate)
		Op2WorkTime := ObtainDuration(op2MinTime, op2MaxTime) //二码处理时间
		item.Op2WorkTime = Op2WorkTime                        //二码处理时间

		item.OpQEntry = opQMinTime.Format("2006-01-02 15:04:05") //问题件进入时间
		item.OpQEnd = opQMaxTime.Format("2006-01-02 15:04:05")   //问题件结束时间
		problemTime := Calculate(opQMinTime, opQMaxTime)
		fmt.Println("problemTime=", problemTime)
		OpQWorkTime := ObtainDuration(op2MinTime, op2MaxTime) //问题件处理时间
		item.OpQWorkTime = OpQWorkTime                        //问题件处理时间
		if bill.BillType == 1 {
			item.BillType = "门诊" //单据类型
		} else if bill.BillType == 2 {
			item.BillType = "住院" //单据类型
		} else if bill.BillType == 3 {
			item.BillType = "门诊,住院" //单据类型
		}

		//---------------------------------------------------- 一二码处理时长
		var createArray []string //创建时间
		var mergeAtArray []string
		for _, block := range blocks {
			if block.Name != "未定义" && block.Name != "" {
				newTime := block.CreatedAt.Format("2006-01-02 15:04:05")
				if newTime != "0001-01-01 00:00:00" && newTime != "" {
					createArray = append(createArray, newTime)
				}
				if item.Op1End != "" {
					mergeAtArray = append(mergeAtArray, item.Op1End)
				}
				if item.Op2End != "" {
					mergeAtArray = append(mergeAtArray, item.Op2End)
				}
			}
		}
		duration := GetOp1AndOp2Duration(createArray, mergeAtArray)
		item.Op1AndOp2Duration = duration //一二码处理时长
		//----------------------------------------------------- 一二码处理时长

		//申请书录入字符、发票录入字符、清单录入字符、其他录入字符
		item.ApplicationChar, item.InvoiceChar, item.ListChar = FetchFieldsByName(obj)
		item.OtherChar = item.FieldCharacter - item.ApplicationChar - item.InvoiceChar - item.ListChar
		finalMap := CountBillTypeOutpatient(fields, bill)
		item.BillTypeOS = finalMap["BillTypeOS"]                     //票据类型（门诊）
		item.BillTypeIH = finalMap["BillTypeIH"]                     //票据类型（住院）
		item.ElectronBillCount = finalMap["ElectronBillCount"]       //电子票据数量
		item.NonElectronBillCount = finalMap["NonElectronBillCount"] //非电子票据数量
		item.InvoiceSum = finalMap["InvoiceSum"]                     //发票张数汇总
		businessDetails = append(businessDetails, item)
	}
	return nil, businessDetails, total
}

func ExportBusinessDetails(info request.BusinessDetailsSearch) (err error, path, name string) {
	proCode := info.ProCode
	total := service2.CountBill(info)
	var businessAllDetails []model.BusinessDetailsExport
	batchSize := 100
	for offset := 0; offset < int(total); offset += batchSize {
		info.PageInfo.PageSize = batchSize
		info.PageInfo.PageIndex = offset/batchSize + 1
		err, businessDetails, _ := GetBusinessDetails(info)
		if err != nil {
			return err, "", ""
		}
		businessAllDetails = append(businessAllDetails, businessDetails...)
	}
	//s := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	//e := strings.Replace(info.EndTime, " 00:00:00", "", -1)
	bookName := proCode + "业务明细表" + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "业务明细表导出/" + proCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcel(basicPath, bookName, "", businessAllDetails)
	return err, basicPath + bookName, bookName
}

//func FindAndCalculateBackAtTheLatest(bill m.ProjectBill, proCode string) string {
//	backAtTheLatest, _, _ := B0108.CalculateBackTimeAndTimeRemaining(bill)
//	return backAtTheLatest
//}

func Calculate(min, max time.Time) string {
	//return max.Sub(min).Round(time.Second).String()
	return formatTimeDuration(max.Sub(min).Round(time.Second))
}

func Compare(times []time.Time) (maxTime, minTime time.Time) {
	isNotInitTime := true
	for i := 0; i < len(times); i++ {
		T := times[i]
		if T.Format("2006-01-02") != "0001-01-01" {
			if isNotInitTime {
				maxTime = T
				minTime = T
				isNotInitTime = false
			}
		} else {
			continue
		}
		if T.Before(minTime) && T.Format("2006-01-02") != "0001-01-01" {
			minTime = T
		}
		if T.After(maxTime) {
			maxTime = T
		}
	}
	return maxTime, minTime
}

func CalculateXmlCharacter(bill m.ProjectBill) int {
	xmlFile := global.GConfig.LocalUpload.FilePath + bill.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/%v.xml",
			bill.CreatedAt.Year(), int(bill.CreatedAt.Month()),
			bill.CreatedAt.Day(), bill.BillNum)
	//xmlFile := "D:/stronger/server/files/B0118/upload_xml/2022/4/13/532022010000106.xml"
	global.GLog.Info("xml file:::" + xmlFile)

	//拿到项目回传单据xml字符串
	data, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		global.GLog.Error(xmlFile+" File reading error", zap.Error(err))
		return 0
	}

	//GKB格式
	if bill.ProCode == "B0108" {
		data, err = ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(data), simplifiedchinese.GBK.NewDecoder()))
		if err != nil {
			global.GLog.Error(xmlFile+" File reading error1", zap.Error(err))
			return 0
		}
	}
	return CalculateXmlValue(string(data))
}

func CalculateXmlValue(xml string) int {
	character := 0
	reg := regexp.MustCompile("<[\\s\\S]*?>[\\s\\S]*?</[\\s\\S]*?>")
	reg1 := regexp.MustCompile("<[\\s\\S]*?>")
	reg2 := regexp.MustCompile("</[\\s\\S]*?>")
	//reg3 := regexp.MustCompile("[^\\x00-\\xff]")
	arr := reg.FindAllString(xml, -1)
	for _, v := range arr {
		//fmt.Println("1", v)
		v = reg2.ReplaceAllString(reg1.ReplaceAllString(v, ""), "")
		//fmt.Println("2", v)
		v = strings.Replace(v, "\n\t\t", "", -1)
		v = strings.ReplaceAll(v, "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		v = strings.ReplaceAll(v, " ", "")
		//fmt.Println("3", v)
		//global.GLog.Info("v", zap.Any("", v))
		//一个半角数字/字母/符号为一个字符，一个全角数字/字母/符号、中文为两个字符(项目有特殊要求，对应项目则以最新的需求为准)
		if v == "" {
			continue
		}
		reg3 := regexp.MustCompile(`·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、`)
		for _, k := range v {
			if unicode.In(k, unicode.Han, unicode.Greek) ||
				reg3.MatchString(string(k)) ||
				strings.IndexRune("ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ", k) != -1 {
				character += 2
			} else {
				character += 1
			}
		}
		//fmt.Println("4", character)
	}
	return character
}

func CalculateWriteCharacter(fields []l.ProjectField) int {
	Character := 0
	//db := global.ProDbMap[bill.ProCode]
	//if db == nil {
	//	global.GLog.Error("CalculateWriteCharacter : "+global.ProDbErr.Error(), zap.Error(errors.New("CalculateWriteCharacter : "+global.ProDbErr.Error())))
	//	return 0
	//}
	//var blocks []l.ProjectBlock
	//err := db.Model(&l.ProjectBlock{}).Order("created_at desc").Where("bill_id = ? ", bill.ID).Find(&blocks).Error
	//if err != nil {
	//	global.GLog.Error(err.Error(), zap.Error(err))
	//	return 0
	//}
	//for _, v := range blocks {
	//var fields []l.ProjectField
	//err := db.Model(&l.ProjectField{}).Order("created_at desc").Where("bill_id = ? ", bill.ID).Find(&fields).Error
	//if err != nil {
	//	global.GLog.Error(err.Error(), zap.Error(err))
	//	return 0
	//}
	for _, f := range fields {
		if f.ResultInput != "no" {
			Character += Difference(f.ResultValue, "")
			//fmt.Println("2", f.Name+f.ResultValue, Character)
		}
	}
	//}
	return Character
}

func CalculateTheNumberOfCase(bill m.ProjectBill) int {
	db := global.ProDbMap[bill.ProCode]
	if db == nil {
		global.GLog.Error("CalculateTheNumberOfCase : "+global.ProDbErr.Error(), zap.Error(errors.New("CalculateTheNumberOfCase : +"+global.ProDbErr.Error())))
		return 0
	}
	var total int64
	err := db.Model(&m.ProjectBill{}).Order("created_at desc").Where("bill_num = ? ", bill.BillNum).Count(&total).Error
	if err != nil {
		global.GLog.Error(err.Error(), zap.Error(err))
		return 0
	}
	return int(total)
}

func CalculateListingNum(blocks []l.ProjectBlock) int {
	ListingNum := 0
	//db := global.ProDbMap[bill.ProCode]
	//if db == nil {
	//	global.GLog.Error("CalculateListingNum : "+global.ProDbErr.Error(), zap.Error(errors.New("CalculateListingNum : +"+global.ProDbErr.Error())))
	//	return 0
	//}
	//var blocks []l.ProjectBlock
	//err := db.Model(&l.ProjectBlock{}).Order("created_at desc").Where("bill_id = ? ", bill.ID).Find(&blocks).Error
	//if err != nil {
	//	global.GLog.Error(err.Error(), zap.Error(err))
	//	return 0
	//}
	for _, v := range blocks {
		if v.IsLoop && strings.Index(v.Name, "清单") != -1 {
			ListingNum++
		}
	}
	return ListingNum
}

func Difference(input, result string) (wrong int) {
	if len(input) > len(result) {
		if result == "" {
			wrong = GetWrongSumVersionTwo(input, result, len(result))
		} else {
			wrong = GetWrongSumVersionTwo(input, result, len(result)) + GetWrongSumVersionTwo(input[len(result):], "", 0)
		}
	} else if len(input) == len(result) {
		wrong = GetWrongSumVersionTwo(input, result, len(input))
	} else if len(input) < len(result) {
		wrong = GetWrongSumVersionTwo(input, result, len(input)) + GetWrongSumVersionTwo("", result[len(input):], 0)
	}
	return wrong
}

func GetWrongSumVersionTwo(input, result string, length int) (wrong int) {
	if input == "" {
		return A(result)
	}
	w := 0
	//匹配中英文?/？
	//reg2 := regexp.MustCompile("[?|？]")
	if result != "" {
		if strings.Index(input, "?") != -1 || strings.Index(input, "？") != -1 {
			return AHasQuestMa(result)
		} else {
			for i := 0; i < length; i++ {
				if input[i] != result[i] {
					w += A(string(input[i]))
				}
			}
		}
	} else {
		return A(input)
	}
	return w
}

func A(str string) int {
	fieldCharacter := 0
	reg2 := regexp.MustCompile("^[?|？]$")
	for _, rr := range str {
		if reg2.MatchString(string(rr)) {
			continue
		}
		if unicode.Is(unicode.Han, rr) {
			fieldCharacter = fieldCharacter + 2
		} else {
			fieldCharacter = fieldCharacter + 1
		}
	}
	return fieldCharacter
}

func AHasQuestMa(str string) int {
	fieldCharacter := 0
	reg2 := regexp.MustCompile("^[?|？]$")
	for _, rr := range str {
		if reg2.MatchString(string(rr)) {
			continue
		}
		if unicode.Is(unicode.Han, rr) {
			fieldCharacter = fieldCharacter + 2
		} else {
			fieldCharacter = fieldCharacter + 1
		}
	}
	return fieldCharacter
}

func GetJiBingFiled(f []l.ProjectField) (str string) {

	//db := global.ProDbMap[proCode]
	//if db == nil {
	//	fmt.Println("GetJiBingFiled-3", global.ProDbErr)
	//}
	//
	//var f []l.ProjectField
	//err := db.Model(&l.ProjectField{}).Where("bill_id = ? ", billId).Find(&f).Error
	//if err != nil {
	//	fmt.Println("GetJiBingFiled-4", err)
	//	return ""
	//}
	for _, v := range f {
		if strings.Index(v.Name, "疾病诊断") != -1 && strings.Index(v.Name, "疾病诊断代码") == -1 && v.ResultValue != "" {
			str += v.ResultValue + ";"
		}
	}
	return str
}

func CalculateMoney(bill m.ProjectBill, t time.Duration) string {
	types := bill.InsuranceType
	if types == "补录" {
		return "2.85"
	}
	if types == "" {
		if t <= time.Hour {
			return "4.28"
		}
		if t <= 2*time.Hour && time.Hour < t {
			return "4.09"
		}
		if t > 2*time.Hour {
			return "2"
		}
	}
	return ""
}

func formatTimeDuration(duration time.Duration) string {
	// 格式化 Duration
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	durationString := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	//fmt.Println(durationString) // 输出：01:30:10
	return durationString
}

// 格式化时间 输出 1.2.3(不要秒)
func ConvertFormat(timeStr string) string {
	matched, _ := regexp.MatchString(`:`, timeStr)
	if !matched {
		return "0"
	}
	splitTime := strings.Split(timeStr, ":")
	formattedTime := splitTime[0] + "." + splitTime[1]
	if len(splitTime) < 3 {
		formattedTime = splitTime[0]
		return formattedTime
	} else {
		formattedTime = splitTime[0] + "." + splitTime[1]
		return formattedTime
	}
	return "0"
}
func GetDurationRange(timeRange string) string {
	float, _ := strconv.ParseFloat(timeRange, 64)
	if float > 0 && float <= 0.5 {
		return "0-0.5小时"
	} else if float > 0.5 && float <= 1 {
		return "0.5-1小时"
	} else if float > 1 && float <= 2 {
		return "1-2小时"
	} else if float > 2 && float <= 3 {
		return "2-3小时"
	} else if float > 3 {
		return "3小时以上"
	}
	return ""
}

// GetOp1AndOp2Duration 返回一二码处理时长
func GetOp1AndOp2Duration(createArray, mergeAtArray []string) string {
	//创建时间取最早的
	var firstTime time.Time
	//1.2码结束取最晚
	var endTime time.Time
	if len(createArray) > 0 && len(mergeAtArray) > 0 {
		if len(createArray) > 1 {
			firstTime, _ = time.Parse("2006-01-02 15:04:05", createArray[0])
			for i := 1; i < len(createArray); i++ {
				parse, _ := time.Parse("2006-01-02 15:04:05", createArray[i])
				if parse.Before(firstTime) {
					firstTime = parse
				}
			}
		} else if len(createArray) == 1 {
			firstTime, _ = time.Parse("2006-01-02 15:04:05", createArray[0])
		}

		if len(mergeAtArray) > 1 {
			endTime, _ = time.Parse("2006-01-02 15:04:05", mergeAtArray[0])
			for i := 1; i < len(mergeAtArray); i++ {
				parse, _ := time.Parse("2006-01-02 15:04:05", mergeAtArray[i])
				if parse.After(endTime) {
					endTime = parse
				}
			}
		} else if len(mergeAtArray) == 1 {
			endTime, _ = time.Parse("2006-01-02 15:04:05", mergeAtArray[0])
		}
		rounds := endTime.Sub(firstTime).Round(time.Second)
		hours := rounds.Hours()
		hours = math.Round(hours*100) / 100
		sprintf := fmt.Sprintf("%.2f", hours)
		return sprintf
	} else {
		return "0.00"
	}
}

// 计算时长 ObtainDuration 大放后面小放前面
func ObtainDuration(beginTime, endTime time.Time) string {
	round := endTime.Sub(beginTime).Round(time.Second)
	hours := round.Hours()
	hours = math.Round(hours*100) / 100
	sprintf := fmt.Sprintf("%.2f", hours)
	return sprintf
}

// FetchFieldsByName 获取对应分块的字段
func FetchFieldsByName(obj billObj) (applicationChar, invoiceChar, listChar int) {
	applicationBlockIdArr := make([]string, 0)
	invoiceBlockIdArr := make([]string, 0)
	listBlockIdArr := make([]string, 0)
	for _, block := range obj.ProjectBlock {
		if strings.Index(block.Name, "申请书") != -1 {
			applicationBlockIdArr = append(applicationBlockIdArr, block.ID)
		}
		if strings.Index(block.Name, "发票") != -1 {
			invoiceBlockIdArr = append(invoiceBlockIdArr, block.ID)
		}
		if strings.Index(block.Name, "清单") != -1 {
			listBlockIdArr = append(listBlockIdArr, block.ID)
		}
	}
	for _, field := range obj.ProjectField {
		if strings.Index(strings.Join(applicationBlockIdArr, ","), field.BlockID) != -1 && field.ResultInput != "no" && field.ResultValue != "" {
			applicationChar += A(field.ResultValue)
		}
		if strings.Index(strings.Join(invoiceBlockIdArr, ","), field.BlockID) != -1 && field.ResultInput != "no" && field.ResultValue != "" {
			invoiceChar += A(field.ResultValue)
		}
		if strings.Index(strings.Join(listBlockIdArr, ","), field.BlockID) != -1 && field.ResultInput != "no" && field.ResultValue != "" {
			listChar += A(field.ResultValue)
		}
	}
	return applicationChar, invoiceChar, listChar
}

// CountBillTypeOutpatient
// B0103 B0106
// 1.票据类型（门诊）：统计fc009字段录入值为“1”的数量 fc009Outpatient
// 2.票据类型（住院）：统计fc009字段录入值为“2”的数量 fc009Hospitalization
// 3.电子票据数量：统计fc090字段录入值为“1”的数量 electronicCount
// 4.非电子票据数量：【发票张数汇总】-【电子票据数量】的值 notElectronicCount
// B0110
// 1.票据类型（门诊）：统计fc009字段录入值为“1”的数量 fc009Outpatient
// 2.票据类型（住院）：统计fc009字段录入值为“2”的数量 fc009Hospitalization
// 3.电子票据数量：统计fc066字段录入值为“1”的数量 electronicCount
// 4.非电子票据数量：【发票张数汇总】-【电子票据数量】的值 notElectronicCount
// B0118
// 1.票据类型（门诊）：统计fc003字段录入值为“1”的数量
// 2.票据类型（住院）：统计fc003字段录入值为“2”的数量
// 3.电子票据数量：统计fc277字段录入值为“1”的数量
// 4.非电子票据数量：【发票张数汇总】-【电子票据数量】的值
// 5.发票张数汇总：统计fc003字段的数量
func CountBillTypeOutpatient(fields []l.ProjectField, bill m.ProjectBill) (list map[string]int) {
	billTypeOS := 0           //票据类型(门诊) 009 1门诊  2住院  fc090  1电子  2非电子
	billTypeIH := 0           //票据类型(住院)
	electronBillCount := 0    //电子票据数量
	invoiceSum := 0           //发票张数汇总
	nonElectronBillCount := 0 //非电子票据数量
	finalMap := make(map[string]int)
	iMap := map[string][]string{
		"B0103": {"fc009", "fc277"},
		"B0106": {"fc009", "fc277"},
		"B0110": {"fc009", "fc066"},
		"B0118": {"fc003", "fc277"},
	}
	codeArr, isExist := iMap[bill.ProCode]
	if isExist {
		for _, field := range fields {
			if field.Code == codeArr[0] {
				if field.ResultValue == "1" {
					billTypeOS = billTypeOS + 1
				} else if field.ResultValue == "2" {
					billTypeIH = billTypeIH + 1
				}
				invoiceSum = invoiceSum + 1
			}
			if field.Code == codeArr[1] && field.ResultValue == "1" {
				electronBillCount = electronBillCount + 1
			}
		}
		nonElectronBillCount = invoiceSum - electronBillCount
	} else {
		if bill.ProCode == "B0108" {
			for _, field := range fields {
				if field.Code == "fc096" {
					billTypeOS++
					invoiceSum++
					nonElectronBillCount++
				}
				if field.Code == "fc097" {
					billTypeIH++
					invoiceSum++
					nonElectronBillCount++
				}
			}
		}
	}

	finalMap["BillTypeOS"] = billTypeOS
	finalMap["BillTypeIH"] = billTypeIH
	finalMap["ElectronBillCount"] = electronBillCount
	finalMap["NonElectronBillCount"] = nonElectronBillCount
	finalMap["InvoiceSum"] = invoiceSum
	return finalMap
}
