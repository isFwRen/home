package B0113

import (
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/service"
	proModel "server/module/pro_conf/model"
	"strconv"
	"strings"
	"time"

	"github.com/wxnacy/wgo/arrays"
	_ "github.com/wxnacy/wgo/arrays"
	_ "go.uber.org/zap"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	global.GLog.Info("B0113:::CheckXml")
	obj := o.(FormatObj)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	// ----------------------------------code---------------------------------------------
	// constMap := constDeal(obj.Bill.ProCode)
	wrongNote += CodeCheck(obj, xmlValue, map[string]map[string]string{})
	wrongNote += XmlCheck(obj, xmlValue, map[string]map[string]string{})

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
	wrongNote := ""
	// _, fc059 := GetOneField(obj.Fields, "fc059", false)
	// _, fc282 := GetOneField(obj.Fields, "fc282", false)
	// if fc059 != fc282 {
	// 	wrongNote += "申请书中银行账户名与证件银行账户名不一致,请检查;"
	// }

	for _, field := range obj.Fields {
		if strings.Index(wrongNote, "案件中存在?号，请核实；") == -1 {
			if !RegIsMatch(field.Code, "^(fc065)$") && RegIsMatch(field.ResultValue, "\\?|？") {
				wrongNote += "案件中存在?号，请核实；"
			}
		}
	}

	cfields := map[string][]string{
		"fc090": {"fc082", "fc234"},
		"fc091": {"fc083", "fc235"},
		"fc092": {"fc084", "fc236"},
		"fc093": {"fc085", "fc237"},
		"fc094": {"fc086", "fc238"},
		"fc095": {"fc087", "fc239"},
		"fc096": {"fc088", "fc240"},
		"fc097": {"fc089", "fc241"},
	}

	for _, MedicalBillInfo := range obj.Invoice {
		_, fc062 := GetOneField(MedicalBillInfo.Fields, "fc062", true)
		_, fc063 := GetOneField(MedicalBillInfo.Fields, "fc063", true)
		if fc062 == "B" {
			_, fc068 := GetOneField(MedicalBillInfo.Fields, "fc068", true)
			_, fc069 := GetOneField(MedicalBillInfo.Fields, "fc069", true)
			if fc068 == fc069 {
				wrongNote += fc063 + "票据出入院日期一致，请确认是否正确;"
			}
		}
		_, fc062 = GetOneField(MedicalBillInfo.Fields, "fc062", false)
		if fc062 == "2" {
			_, fc068 := GetOneField(MedicalBillInfo.Fields, "fc068", true)
			_, fc069 := GetOneField(MedicalBillInfo.Fields, "fc069", true)
			a, _ := time.Parse("2006-01-02", fc068)
			b, _ := time.Parse("2006-01-02", fc069)
			if a.After(b) {
				wrongNote += fc063 + "入院日期晚于出院日期，请修改;"
			}

		}
		cValue := []string{}
		isWenhao := false
		_, fc332 := GetOneField(MedicalBillInfo.Fields, "fc332", false)
		fc332_names := ""
		for _, field := range MedicalBillInfo.Fields {
			if RegIsMatch(field.Code, `^(fc224|fc225|fc226|fc227|fc228|fc229|fc229|fc230|fc231|fc232|fc233|fc252|fc253|fc254|fc255|fc256)$`) && strings.Index(field.ResultValue, "-") != -1 {
				wrongNote += fc063 + "票据的" + field.Name + "不允许为负数，请修改;"
			}
			if fc332 == "1" {
				if RegIsMatch(field.Code, `^(fc333|fc334|fc335|fc336)$`) && field.FinalValue == "" {
					fc332_names += field.Name + "、"
				}
			} else if fc332 == "2" {
				if RegIsMatch(field.Code, `^(fc333|fc334|fc335|fc337|fc072)$`) && field.FinalValue == "" {
					fc332_names += field.Name + "、"
				}
			}
		}
		if fc332_names != "" {
			wrongNote += "发票[" + fc063 + "]" + fc332_names + "为空,请检查;"
		}
		for _, QingDan := range MedicalBillInfo.QingDan {
			aa := true
			for _, field := range QingDan.Fields {
				if RegIsMatch(field.Code, `^(fc082|fc234|fc090|fc083|fc235|fc091|fc084|fc236|fc092|fc085|fc237|fc093|fc086|fc238|fc094|fc087|fc239|fc095|fc088|fc240|fc096|fc089|fc241|fc097)$`) {
					if strings.Index(field.ResultValue, "?") != -1 && strings.Index(field.ResultValue, fc063+"票据的项目名称，项目明细分类，项目金额未填写或模糊;") == -1 {
						wrongNote += fc063 + "票据的项目名称，项目明细分类，项目金额未填写或模糊;"
					}

				}
				if RegIsMatch(field.Code, `^(fc260|fc261|fc262|fc263|fc264|fc265|fc266|fc267)$`) {
					cValue = append(cValue, field.ResultValue)
					if field.ResultValue == "" || field.ResultValue == "2" || field.ResultValue == "5" {
						aa = false
					}
					if strings.Index(field.ResultValue, "?") != -1 {
						isWenhao = true
					}
				}

				if RegIsMatch(field.Code, `^(fc098|fc099|fc100|fc101|fc102|fc103|fc104|fc105)$`) {
					if field.FinalValue != "1.00" {
						aa = false
					}
				}

				cCode, is := cfields[field.Code]
				if is {
					_, fValue0 := GetOneField(QingDan.Fields, field.Code, false)
					_, fValue := GetOneField(QingDan.Fields, cCode[0], false)
					_, fValue1 := GetOneField(QingDan.Fields, cCode[1], false)
					if (fValue0 != "" || fValue != "" || fValue1 != "") && (fValue0 == "" || fValue == "" || fValue1 == "") && strings.Index(wrongNote, fc063+"票据的项目名称，项目明细分类，项目金额未填写或模糊;") == -1 {
						wrongNote += fc063 + "票据的项目名称，项目明细分类，项目金额未填写或模糊;"
					}
					if field.ResultValue == "0" {
						if fValue != "" && fValue != "A" && strings.Index(fValue, "?") == -1 {
							wrongNote += "发票号[" + fc063 + "]" + fValue + "清单项目金额不能为0，请检查;"
						}
					}

				}
			}
			if aa {
				wrongNote += "发票" + fc063 + "的清单分块为全自费项目，请确认;"
			}
		}
		if arrays.Contains(cValue, "5") != -1 {
			if arrays.Contains(cValue, "1") != -1 || arrays.Contains(cValue, "2") != -1 || arrays.Contains(cValue, "3") != -1 || arrays.Contains(cValue, "4") != -1 {
				wrongNote += fc063 + "票据，请确认医保类型是否正确;"
			}
		}
		if isWenhao {
			wrongNote += fc063 + "票据医保类型异常，请修改;"
		}
	}

	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""

	c, ok := global.GProConf[obj.Bill.ProCode].ConstTable["B0113_百年理赔_百年理赔费用项目代码表"]
	feiyong := make(map[string]string, 0)
	if ok {
		fmt.Println("------------------------c-----------------------:", c[0])
		for _, arr := range c {
			feiyong[strings.TrimSpace(arr[0])] = arr[0]
		}
	}

	State := GetNodeValue(xmlValue, "State")
	if RegIsMatch(State, `^(1|2|3|4|5|6)$`) {
		wrongNote += "请检查影像上的索赔类型是否勾选“医疗费用”，如有勾选，则可直接返回；如未勾选，则把[案件录入状态]字段修改为0后返回；"
	}

	_, num := service.GetCountBillNum("B0113", obj.Bill.BillNum)
	if num > 1 && State != "0" {
		wrongNote += "案件重复推送，请确认案件录入状态是否正确；"
	}

	// FeeItemCodes := RegMatchAll(xmlValue, `<FeeItemCode>.*?<\/FeeItemCode>`)
	// for _, FeeItemCode := range FeeItemCodes {
	// 	_, isExits := feiyong[FeeItemCode]
	// 	if FeeItemCode == "" || !isExits {
	// 		wrongNote += "项目名称" + FeeItemCode + "不在常量库中，FeeItemCode不能为空，项目编码-发票错误;"
	// 	}
	// }

	lLFeeMainSchemas := RegMatchAll(xmlValue, `<LLFeeMainSchema>[\s\S]*?<\/LLFeeMainSchema>`)
	MainFeeNoMap := map[string]int{}
	for _, lLFeeMainSchema := range lLFeeMainSchemas {
		lLCaseReceiptSchemas := RegMatchAll(lLFeeMainSchema, `<LLCaseReceiptSchema>[\s\S]*?<\/LLCaseReceiptSchema>`)
		FeeSum := 0.00
		FeeAmnt := GetNodeValue(lLFeeMainSchema, "FeeAmnt")
		MainFeeNo := GetNodeValue(lLFeeMainSchema, "MainFeeNo")
		num, isOK := MainFeeNoMap[MainFeeNo]
		if isOK {
			MainFeeNoMap[MainFeeNo] = num + 1
			if RegIsMatch(wrongNote, `发票号`+MainFeeNo+`重复录入\d+次，请检查;`) {
				wrongNote = RegReplace(wrongNote, `发票号`+MainFeeNo+`重复录入\d+次，请检查;`, `发票号`+MainFeeNo+`重复录入`+strconv.Itoa(num+1)+`次，请检查;`)
			} else {
				wrongNote += "发票号" + MainFeeNo + "重复录入1次，请检查;"
			}
		} else {
			MainFeeNoMap[MainFeeNo] = 1
		}

		cccs := [][]string{{"FeeAmnt", "账单总金额"}, {"FeeRealAmnt", "调整后的账单总金额* 默认等于账单总金额"}, {"SelfAmnt", "自费总金额"}, {"SelfItemAmnt", "自费项目总金额"}, {"SelfDrugAmnt", "自费药品总金额"}}
		for _, ccc := range cccs {
			cc := GetNodeValue(lLFeeMainSchema, ccc[0])
			if cc == "" {
				wrongNote += "[" + MainFeeNo + "]发票[" + ccc[1] + "]为空，请检查；"
			}
		}

		FeeType := GetNodeValue(lLFeeMainSchema, "FeeType")

		qq := 0.00
		fWrong := ""
		for _, lLCaseReceiptSchema := range lLCaseReceiptSchemas {

			FeeItemCode := GetNodeValue(lLCaseReceiptSchema, "FeeItemCode")
			_, isExits := feiyong[FeeItemCode]
			if FeeItemCode == "" || !isExits {
				if fWrong != "" {
					fWrong += ","
				}

				// wrongNote += "项目名称" + FeeItemCode + "不在常量库中，FeeItemCode不能为空，项目编码-发票错误;"
				if FeeItemCode == "" {
					FeeItemName := GetNodeValue(lLCaseReceiptSchema, "FeeItemName")
					fWrong += FeeItemName + " "
				} else {
					fWrong += FeeItemCode + " "
				}
			}
			Fee := GetNodeValue(lLCaseReceiptSchema, "Fee")
			if !RegIsMatch(FeeItemCode, `^(CS001|CZ001|CZ002)$`) {
				FeeSum = SumFloat(FeeSum, ParseFloat(Fee), "+")
			}

			if RegIsMatch(FeeItemCode, `^(CS001|CZ001)$`) {
				qq = SumFloat(qq, ParseFloat(Fee), "+")
			}

			if RegIsMatch(Fee, `^(0.00|B|\?)$`) {
				if FeeItemCode == "CS001" {
					wrongNote += MainFeeNo + "票据中报销金额为0，请检查;"
				}
				if FeeItemCode == "CZ001" {
					wrongNote += MainFeeNo + "票据中自费金额为0，请检查;"
				}

			}

		}
		if fWrong != "" {
			wrongNote += "发票号[" + MainFeeNo + "]项目名称" + fWrong + "不在常量库中，请检查;"
		}
		fmt.Println("-------------------FeeType-------------------", FeeType, qq, FeeAmnt)
		if FeeType == "A" || FeeType == "B" {

			if qq > ParseFloat(FeeAmnt) {
				wrongNote += MainFeeNo + "票据中的报销金额与自费总金额相加大于账单总金额，请检查;"
			}
		}

		if FeeType == "B" && qq == ParseFloat(FeeAmnt) {
			wrongNote += MainFeeNo + "票据中的报销金额与自费总金额相加等于账单总金额，请检查;"

		}

		chae := SumFloat(ParseFloat(FeeAmnt), FeeSum, "-")
		if chae != 0.00 {
			wrongNote += "账单号" + MainFeeNo + "明细金额与账单总金额不一致，差额" + ToString(chae) + ";"
		}
	}

	return wrongNote
}
