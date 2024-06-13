package B0114

import (
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/service"
	proModel "server/module/pro_conf/model"
	sUtils "server/utils"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
	_ "github.com/wxnacy/wgo/arrays"
	_ "go.uber.org/zap"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	global.GLog.Info("B0114:::CheckXml")
	obj := o.(FormatObj)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	// ----------------------------------code---------------------------------------------
	constMap := constDeal(obj.Bill.ProCode)
	wrongNote += CodeCheck(obj, xmlValue, constMap)
	wrongNote += XmlCheck(obj, xmlValue, constMap)

	// ----------------------------------xml---------------------------------------------

	global.GLog.Error("wrongNote：：：" + wrongNote)
	fmt.Println("------------------------wrongNote-----------------------:", wrongNote)
	return err, wrongNote
}

func CheckWrongNote(pro, xmlValue string, obj FormatObj) (error, string) {
	// fmt.Println("---------CheckWrongNoteCheckWrongNote------------")
	wrongNote := ""

	err, fieldCheckConfs := service.GetProFieldCheckConf(pro)
	fmt.Println("--------fieldCheckConfs-err------------", err)
	if err != nil {
		return err, wrongNote
	}
	fieldCheckConfMap := make(map[string][]proModel.SysProFieldCheck)
	// fmt.Println("---------fieldCheckConfMap------------", len(fieldCheckConfMap))
	for _, fieldCheckConf := range fieldCheckConfs {
		fieldCheckConfMap[fieldCheckConf.Code] = fieldCheckConf.SysProFieldChecks
	}
	eleLen := reflect.ValueOf(obj).NumField()
	for j := 0; j < eleLen; j++ {
		if reflect.TypeOf(obj).Field(j).Name != "Bill" && reflect.TypeOf(obj).Field(j).Name != "Fields" {
			//每张发票每种类型的字段
			fmt.Println("---------------------------", reflect.TypeOf(obj).Field(j).Name)
			fieldsMaps := reflect.ValueOf(obj).Field(j).Interface().([]FieldsMap)
			for _, fieldsMap := range fieldsMaps {
				if fieldsMap.Code == "" {
					continue
				}
				for _, field := range fieldsMap.Fields {
					items, isExit := fieldCheckConfMap[field.Code]
					// fmt.Println("---------items------------", items)
					if isExit {
						for _, item := range items {
							fffs := strings.Split(item.Value, ";")
							// fmt.Println("---------fffs------------", fffs)
							for _, fff := range fffs {
								mess := "账单号:" + fieldsMap.Code + item.Mark + ";"
								if strings.Index(wrongNote, mess) != -1 {
									continue
								}
								if item.CheckType == "1" {
									if field.ResultValue == fff {
										wrongNote += mess
									}
								} else if item.CheckType == "2" {
									if strings.Index(field.ResultValue, fff) != -1 {
										wrongNote += mess
									}
								} else if item.CheckType == "3" {
									if strings.Index(field.ResultValue, fff) == -1 {
										wrongNote += mess
									}
								}
							}

						}
					}
				}
			}
		}
	}
	// for
	return err, wrongNote

}

func CodeCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	fields := obj.Fields
	wrongNote := ""
	fc092Arr := []string{}
	for _, MedicalBillInfo := range obj.MedicalBillInfo {
		_, fc126 := GetOneField(MedicalBillInfo.Fields, "fc126", false)
		fc092s := GetSameField(MedicalBillInfo.Fields, "fc092", false)
		_, fc104 := GetOneField(MedicalBillInfo.Fields, "fc104", true)
		for _, fc092 := range fc092s {
			if fc092 != "" {
				if sUtils.HasItem(fc092Arr, fc092) {
					wrongNote += fc104 + "存在多张报销单，请检查报销金额/其他商保金额;"
				} else {
					fc092Arr = append(fc092Arr, fc092)
				}
			}
		}
		if fc126 != "A" {
			sum := 0.0
			isExits := false
			for _, field := range MedicalBillInfo.Fields {
				// fmt.Println("------------------------Code-----------------------:", field.Code, field.ResultValue)
				if RegIsMatch(field.Code, `^(fc268|fc269|fc270)$`) {
					value := field.ResultValue
					isExits = true
					if !RegIsMatch(value, `^(A|0|)$`) && strings.Index(value, "?") == -1 {
						sum += ParseFloat(value)
					}
				}
			}
			// fmt.Println("------------------------fc126-----------------------:", fc126, sum)
			if isExits && ParseFloat(fc126) != sum {
				wrongNote += "发票" + fc104 + "的发票报销金额不等于报销总金额，请检查;"
			}
		}
		for _, QingDan := range MedicalBillInfo.QingDan {
			for _, field := range QingDan.Fields {
				if RegIsMatch(field.Code, `^(fc257|fc258|fc259|fc260|fc261|fc262|fc263|fc264)$`) {
					if strings.Index(field.ResultValue, "?") != -1 {
						wrongNote += fc104 + "账单" + field.Name + "异常，请修改"
					}
				}
			}

		}

	}
	for _, field := range fields {
		if field.Code == "fc014" && len(field.Issues) > 0 {
			wrongNote += "意外细节录入内容不在常量库中，请检查;"
		}
		if field.Code == "fc015" && len(field.Issues) > 0 {
			wrongNote += "损伤外部原因录入内容不在常量库中，请检查;"
		}

		if strings.Index(wrongNote, "案件中存在?号，请核实；") == -1 {
			if !RegIsMatch(field.Code, "^(fc108|fc106|fc104|fc111|fc112)$") && RegIsMatch(field.ResultValue, "\\?|？") {
				wrongNote += "案件中存在?号，请核实；"
			}
		}
	}
	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""
	causeCodes := sUtils.GetNodeData(obj.Bill.OtherInfo, "causeCode")
	// fmt.Println("----------------------------------------------------------------------------:")
	NaNs := RegMatchAll(xmlValue, `>.*NaN.*?<\/`)
	if len(NaNs) > 0 {
		wrongNote += "录入值错误导致存在“NaN”，请修改;"
	}

	if sUtils.HasItem(causeCodes, "21") || sUtils.HasItem(causeCodes, "22") {
		surgeryCodes := RegMatchAll(xmlValue, `<surgeryCode>.+?<\/surgeryCode>`)
		operationMedicalAmounts := RegMatchAll(xmlValue, `<operationMedicalAmount>(0|0.00|)<\/operationMedicalAmount>`)
		if len(surgeryCodes) > 0 && len(operationMedicalAmounts) > 0 {
			wrongNote += "有手术代码但缺少手术津贴内容，请确认；"
		}
	}
	surgeryCodes := RegMatchAll(xmlValue, `<surgeryCode>.*?<\/surgeryCode>`)
	// menzhen := map[string]string{}
	// zhuyuan := map[string]string{}
	cValue := []string{}
	for _, surgeryCode := range surgeryCodes {
		bValue := GetNodeValue(surgeryCode, "surgeryCode")
		if bValue != "" {
			if arrays.ContainsString(cValue, bValue) != -1 {
				wrongNote += "手术术士编码" + bValue + "重复，请检查；"
			}
			cValue = append(cValue, bValue)
		}
	}

	cValue = []string{}
	cacheXmlValue := RegReplace(xmlValue, `<medicalDefuctInfos>[\s\S]*?<\/medicalDefuctInfos>`, "")
	billNos := RegMatchAll(cacheXmlValue, `<billNo>.*?<\/billNo>`)
	for _, billNo := range billNos {
		bValue := GetNodeValue(billNo, "billNo")
		if bValue != "" {
			if arrays.ContainsString(cValue, bValue) != -1 {
				wrongNote += "发票号" + bValue + "重复，请检查；"
			}
			cValue = append(cValue, bValue)
		}
	}

	medicalBillInfos := RegMatchAll(xmlValue, `<medicalBillInfo>[\s\S]*?<\/medicalBillInfo>`)
	// fmt.Println("------------------------medicalBillInfos-----------------------:", len(medicalBillInfos))
	for _, medicalBillInfo := range medicalBillInfos {
		billNo := GetNodeValue(medicalBillInfo, "billNo")
		billType := GetNodeValue(medicalBillInfo, "billType")
		therapyType := GetNodeValue(medicalBillInfo, "therapyType")
		if billType == "06" && therapyType == "1" {
			wrongNote += "发票号" + billNo + "账单类型为“新农合住院收据”，治疗类型不能为门诊；"
		}
		inHospitalDate := GetNodeValue(medicalBillInfo, "inHospitalDate")
		outHospitalDate := GetNodeValue(medicalBillInfo, "outHospitalDate")
		a, _ := time.Parse("2006/01/02", inHospitalDate)
		b, _ := time.Parse("2006/01/02", outHospitalDate)
		if a.After(b) {
			wrongNote += "发票号" + billNo + "入院日期晚于出院日期;"
		}
		// if therapyType == "2" {
		// 	medicalBillDetail := RegMatchAll(medicalBillInfo, `<medicalBillDetails>[\s\S]*?<\/medicalBillDetails>`)[0]
		// 	if GetNodeValue(medicalBillDetail, "count") == "1" {
		// 		wrongNote += "票据号" + billNo + "，医疗费用类型不可汇总录入，请确认后修改;"
		// 	}
		// }
		// medicalBillAmount := GetNodeValue(medicalBillInfo, "medicalBillAmount")
		medicalBillAmounts := RegMatchAll(medicalBillInfo, `<medicalBillAmount>.*?<\/medicalBillAmount>`)
		medicalBillAmount := 0.0
		for _, medicalBillInfo := range medicalBillAmounts {
			medicalBillAmount += ParseFloat(GetNodeValue(medicalBillInfo, "medicalBillAmount"))
		}
		billAmount := GetNodeValue(medicalBillInfo, "billAmount")
		abateAmount := GetNodeValue(medicalBillInfo, "abateAmount")
		// fmt.Println("------------------------billAmount-----------------------:", billAmount, medicalBillAmount, abateAmount)
		if billAmount != "" && abateAmount != "" {
			chae := SumFloat(ParseFloat(billAmount), medicalBillAmount, "-")
			if chae != 0 && chae != 0.0 {
				wrongNote += "发票号" + billNo + "医疗费用金额总金额不等于费用金额，差额" + ToString(chae) + ";"
			}
			if ParseFloat(abateAmount) > ParseFloat(billAmount) {
				wrongNote += "发票号" + billNo + "扣减金额（社保支付+新农合支付+其他商保支付+其他扣减费用+自费金额+自付2）大于费用金额;"
			}
			if billAmount == abateAmount {
				wrongNote += "发票号" + billNo + "发票总金额与扣减金额一致，请确认;"

			}
		}
		// fmt.Println("--------------------!!!!!!!!!!!---------------------:")

	}

	return wrongNote
}

func SumFloat(a1, a2 float64, ff string) float64 {
	b1 := decimal.NewFromFloat(a1)
	b2 := decimal.NewFromFloat(a2)
	// b1.StringFixed(2)
	if ff == "+" {
		return ParseFloat(b1.Add(b2).StringFixed(2))
	}
	if ff == "-" {
		return ParseFloat(b1.Sub(b2).StringFixed(2))
	}

	if ff == "*" {
		return ParseFloat(b1.Mul(b2).StringFixed(2))
	}

	return 0.00
}
