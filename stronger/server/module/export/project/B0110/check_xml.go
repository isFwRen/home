package B0110

import (
	"encoding/json"
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/service"
	"server/module/load/model"
	proModel "server/module/pro_conf/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/wxnacy/wgo/arrays"
	_ "github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	global.GLog.Info("B0110:::CheckXml")
	obj := o.(FormatObj)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	// ----------------------------------code---------------------------------------------
	constMap := InitConst(obj.Bill.ProCode)
	wrongNote += XmlCheck(obj, xmlValue, constMap)
	wrongNote += CodeCheck(obj, xmlValue, constMap)

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
	//sanfang := constMap["sanfang"]
	//ICD10Code := constMap["ICD10Code"]
	// fmt.Println("--------------sanfang----------------------", sanfang)
	var fc091Arr []string
	var fc005Arr []string
	for _, invoice := range obj.Invoice {
		_, fc005 := GetOneField(invoice.Fields, "fc005", true)
		fc020s := []string{}
		fc061s := []string{}
		var filedArr []string

		for _, QinDan := range invoice.QingDan {
			for _, field := range QinDan.Fields {
				//需求编码 CSB0110RC0091000 清单
				if RegIsMatch(field.ResultValue, `[?？]`) {
					filedArr = append(filedArr, field.Name)
				}
				fmt.Println("QinDan.field=", field)
			}

			//CSB0110RC0115000
			sFields := [][]string{
				//名称 金额 数量
				{"fc012", "fc014", "fc013"},
				{"fc025", "fc039", "fc032"},
				{"fc026", "fc040", "fc033"},
				{"fc027", "fc041", "fc034"},
				{"fc028", "fc042", "fc035"},
				{"fc029", "fc043", "fc036"},
				{"fc030", "fc044", "fc037"},
				{"fc031", "fc045", "fc038"},
			}
			var cSB0110RC0115000Msg []string
			for _, item := range sFields {
				_, name := GetOneField(QinDan.Fields, item[0], false)
				_, money := GetOneField(QinDan.Fields, item[1], false)
				_, number := GetOneField(QinDan.Fields, item[2], false)
				if !(len(name) == 0 && len(money) == 0 && len(number) == 0) &&
					!(len(name) > 0 && len(money) > 0 && len(number) > 0) {
					for _, value := range item {
						_, val := GetOneField(QinDan.Fields, value, true)
						if len(val) == 0 {
							cSB0110RC0115000Msg = append(cSB0110RC0115000Msg, value)
						}
					}
				}
			}
			if len(cSB0110RC0115000Msg) > 0 {
				mes := "【" + fc005 + "】发票【" + strings.Join(cSB0110RC0115000Msg, "、") + "】字段位为空；"
				if strings.Index(wrongNote, mes) == -1 {
					wrongNote += mes
				}
			}
		}

		//_, fc005nc := GetOneField(invoice.Fields, "fc005", true)
		//需求编码 CSB0110RC0091000
		for _, field := range invoice.Fields {
			fmt.Println("field=", field)
			if field.Code == "fc020" && field.ResultValue != "" {
				if arrays.Contains(fc020s, field.ResultValue) == -1 {
					fc020s = append(fc020s, field.ResultValue)
				} else {
					wrongNote += "发票号[" + fc005 + "]同一发票对应多张报销单；"
				}
			}
			if field.Code == "fc061" && field.ResultValue != "" {
				// fmt.Println("--------------fc061----------------------", field.ResultValue)
				v, total := utils.FetchConst(obj.Bill.ProCode, "B0110_新疆国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": field.ResultValue})
				fmt.Println("==============total =======", total)
				fmt.Println("============== v =======", v)

				if total == 0 && strings.Index(wrongNote, "第三方出具单位录入有误，请检查;") == -1 {
					wrongNote += "第三方出具单位录入有误，请检查;"
				}
				if arrays.Contains(fc061s, field.ResultValue) == -1 {
					fc061s = append(fc061s, field.ResultValue)
				} else {
					wrongNote += "发票号[" + fc005 + "]重复录入第三方支付公司，请修改；"
				}
			}

			//CSB0110RC0104000
			//同一发票下，当fc066录入值为1时，校验fc068、fc069、fc070、fc071、fc072是否存在结果值，任意一个字段结果值为空时，出导出校验：电子发票五要素未录入齐全，请确认是否为电子发票；
			sFiled := []string{"fc068", "fc069", "fc070", "fc071", "fc072"}
			if field.Code == "fc066" {
				_, fc091Val := GetOneField(invoice.Fields, "fc091", false)
				if fc091Val != "3" && field.ResultValue == "1" {
					for k := 0; k < len(sFiled); k++ {
						_, val := GetOneField(invoice.Fields, sFiled[k], true)
						if val == "" {
							mes := "【" + fc005 + "】" + "电子发票五要素未录入齐全，请确认是否为电子发票；"
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}

					}
				}
			}
			//CSB0110RC0105000
			//"同一发票下，当fc067录入值为1时，（xxx为fc005的结果值）
			//1、校验fc068的结果值是否为12位数，否则出导出校验：【XX】发票增值税票据代码字段应为12位数，请检查；
			//2、校验fc069的结果值是否为8位数，否则出导出校验：【XX】发票增值税票据号码字段应为8位数，请检查；
			//3、校验fc071的结果值是否为6位数，否则出导出校验：【XX】发票校验码字段应为6位数，请检查；
			//4、fc073的结果值不为“01”时，出导出校验，【XX】发票查验结果不通过，请确认五要素是否录入正确；"
			//CSB0110RC0106000
			//"同一发票下，当fc067录入值为2时，（xxx为fc005的结果值）
			//1、校验fc068的结果值是否为8位数，否则出导出校验：XX】发票财政部票据代码字段应为8位数，请检查；
			//2、校验fc069的结果值是否为10位数，否则出导出校验：【XX】发票财政部票据号码字段应为10位数，请检查；"
			//3、fc073的结果值不为“01”时，出导出校验，【XX】发票查验结果不通过，请确认五要素是否录入正确；
			sFiledTo := [][]string{{"fc068", "12", "发票增值税票据代码字段应为12位数，请检查；"}, {"fc069", "8", "发票增值税票据号码字段应为8位数，请检查；"}, {"fc071", "6", "发票校验码字段应为6位数，请检查；"}}
			sFiledThree := [][]string{{"fc068", "8", "发票财政部票据代码字段应为8位数，请检查；"}, {"fc069", "10", "发票财政部票据号码字段应为10位数，请检查；"}}
			if field.Code == "fc067" {
				if field.ResultValue == "1" {
					for k := 0; k < len(sFiledTo); k++ {
						_, val := GetOneField(invoice.Fields, sFiledTo[k][0], true)
						atoi, _ := strconv.Atoi(sFiledTo[k][1])
						if val != "" && len(val) != atoi {
							mes := "【" + fc005 + "】" + sFiledTo[k][2]
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}
					}
					//_, val097 := GetOneField(invoice.Fields, "fc073", true)
					//if val097 != "" && val097 != "01" {
					//	mes := "【" + fc005 + "】发票查验结果不通过，请确认五要素是否录入正确；"
					//	if strings.Index(wrongNote, mes) == -1 {
					//		wrongNote += mes
					//	}
					//}

				} else if field.ResultValue == "2" {
					_, fc070Val := GetOneField(invoice.Fields, "fc070", true)
					fc070Parse, _ := time.Parse("2006-01-02", fc070Val)
					subTime := time.Now().Sub(fc070Parse)
					oneYear, _ := time.ParseDuration("8760h")
					for k := 0; k < len(sFiledThree); k++ {
						_, val := GetOneField(invoice.Fields, sFiledThree[k][0], true)
						atoi, _ := strconv.Atoi(sFiledThree[k][1])
						if val != "" && len(val) != atoi {
							mes := "【" + fc005 + "】" + sFiledThree[k][2]
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}
					}
					if subTime < oneYear {
						_, val := GetOneField(invoice.Fields, "fc073", true)
						if val != "01" {
							fc005Arr = append(fc005Arr, fc005)
						}
					}
				} else if field.ResultValue == "3" {
					_, val := GetOneField(invoice.Fields, "fc005", false)
					if len(val) < 20 {
						fc091Arr = append(fc091Arr, fc005)
					}
				}
			}

		}

		// _, fc061 := GetOneField(invoice.Fields, "fc061", false)
		// if !HasKey(sanfang, fc061) && strings.Index(wrongNote, "第三方出具单位录入有误，请检查;") == -1 {
		// 	wrongNote += "第三方出具单位录入有误，请检查;"
		// }

		_, fc010 := GetOneField(invoice.Fields, "fc010", true)
		isFc024, fc024 := GetOneField(invoice.Fields, "fc024", true)

		if isFc024 && fc010 != fc024 {
			wrongNote += "发票[" + fc005 + "]与结算单的报销金额不一致，请检查;"
		}
		_, fc017 := GetOneField(invoice.Fields, "fc017", false)
		isFc024, fc024 = GetOneField(invoice.Fields, "fc024", false)
		if isFc024 && ParseFloat(fc017) > ParseFloat(fc024) {
			wrongNote += fc005 + "其他报销金额必须小于或等于报销金额，请修改;"
		}

		_, fc066 := GetOneField(invoice.Fields, "fc066", true)
		_, fc073 := GetOneField(invoice.Fields, "fc073", true)
		if fc066 == "Y" && fc073 == "" {
			wrongNote += fc005 + "发票,发票查验结果为空,请检查;"
		}
	}

	for _, field := range obj.Fields {
		if field.Code == "fc007" {
			if field.ResultValue == "" {
				wrongNote += "疾病诊断不能为空；"
			}
			//编码 CSB0110RC0092000 当fc007录入值不为常量表《B0110_新疆国寿理赔_ICD10疾病编码》中“疾病名称”（第一列）一列内容时，出导出校验：疾病诊断不在常量库内，请检查；
			_, total := utils.FetchConst(obj.Bill.ProCode, "B0110_新疆国寿理赔_ICD10疾病编码", "疾病名称", map[string]string{"疾病名称": field.ResultValue})
			if total == 0 && field.ResultValue != "" {
				if strings.Index(wrongNote, "疾病诊断不在常量库内，请检查；") == -1 {
					wrongNote += "疾病诊断不在常量库内，请检查；"
				}
			}
		}
	}

	if len(fc005Arr) > 0 {
		mes := "【" + strings.Join(fc005Arr, "、") + "】" + "发票查验结果不通过，请确认五要素是否录入正确；"
		if strings.Index(wrongNote, mes) == -1 {
			wrongNote += mes
		}
	}

	if len(fc091Arr) > 0 {
		mes := "【" + strings.Join(fc091Arr, "、") + "】" + "发票票据号不为20位数，请检查；"
		if strings.Index(wrongNote, mes) == -1 {
			wrongNote += mes
		}
	}

	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""

	var otherInfo OtherInfo
	err := json.Unmarshal([]byte(obj.Bill.OtherInfo), &otherInfo)
	if err != nil {
		global.GLog.Error("otherInfo", zap.Error(err))
	}

	bpoSendRemark := otherInfo.BpoSendRemark
	if bpoSendRemark != "" {
		bpoSendRemark = RegReplace(bpoSendRemark, `【.*】`, "")
		bpoSendRemark = RegReplace(bpoSendRemark, `(；|，|；)`, "")
		if bpoSendRemark != "" {
			wrongNote += bpoSendRemark + ";"
		}
	}

	//items := []string{}
	//rcptNos := RegMatchAll(xmlValue, `<rcptNo>.*?<\/rcptNo>`)
	//for _, rcptNo := range rcptNos {
	//	if arrays.Contains(items, rcptNo) != -1 {
	//		wrongNote += "发票号[" + rcptNo + "]重复，请检查；"
	//	} else {
	//		items = append(items, rcptNo)
	//	}
	//}

	// if strings.Index(otherInfo["bpoSendRemark"].(string), "全民") != -1 {

	rcptInfoLists := RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	rcptNoName := []string{}
	for _, rcptInfoList := range rcptInfoLists {
		rcptNo := GetNodeValue(rcptInfoList, "rcptNo")
		endDate := GetNodeValue(rcptInfoList, "endDate")
		beginDate := GetNodeValue(rcptInfoList, "beginDate")
		a, _ := time.Parse("2006-01-02", endDate)
		b, _ := time.Parse("2006-01-02", beginDate)
		if a.Before(b) {
			wrongNote += "发票号[" + rcptNo + "]出院时间早于入院时间，请检查；"
		}
		socialPayAmnt := GetNodeValue(rcptInfoList, "socialPayAmnt")
		rcptAmnt := GetNodeValue(rcptInfoList, "rcptAmnt")
		fmt.Println("==================rcptNo==========", rcptNo)
		fmt.Println("===========ParseFloat(socialPayAmnt)========", ParseFloat(socialPayAmnt))
		fmt.Println("===========ParseFloat(rcptAmnt)========", ParseFloat(rcptAmnt))
		if ParseFloat(socialPayAmnt) >= ParseFloat(rcptAmnt) {
			if rcptNo != "" {
				wrongNote += "发票号[" + rcptNo + "]报销金额大于或等于发票总金额，请检查；"
			}
		}
		rcptLists := RegMatchAll(rcptInfoList, `<(rcptList|errorRcptList)>[\s\S]*?<\/(rcptList|errorRcptList)>`)
		sum := 0.0
		for _, rcptList := range rcptLists {
			listName := GetNodeValue(rcptList, "listName")
			if RegIsMatch(listName, `(第三方报销|全自费费用|先行自付费用|社保统筹范围内费用|不予支付部分|目录外费用|超限价自付费用|公务员报销|床位费补助|企补支付|大病报销)`) {
				wrongNote += "发票号" + rcptNo + "报销单录入错误；"
			}
			if strings.Index(rcptList, "<errorRcptList>") != -1 && RegIsMatch(listName, `^(6501000000014689|6527000000014708|6532000000014698|6541000000014699|6542000000014698|6590000000014708)$`) {
				continue
			}
			quantity := GetNodeValue(rcptList, "quantity")
			price := GetNodeValue(rcptList, "price")
			if quantity != "" && price != "" {
				total := SumFloat(ParseFloat(quantity), ParseFloat(price), "*")
				sum = SumFloat(sum, total, "+")
			}
		}
		fmt.Println("--------------rcptAmnt----------------------", rcptAmnt, sum)
		chae := SumFloat(ParseFloat(rcptAmnt), sum, "-")
		// expenMode := GetNodeValue(rcptInfoList, "expenMode")
		if sum != 0.0 && (chae >= 0.1 || chae <= -0.1) {
			wrongNote += "发票号[" + rcptNo + "]清单项目金额与发票金额不一致，差额:" + ToString(chae) + "，请检查；"
		}
		//CSB0110RC0050001 每个案件可能包含多个rcptInfoList，将每个rcptInfoList大节点下的第一个<rcptNo>进行对比，
		//如果有重复值，将自动回传改为手动回传，并在保单列表界面案件下出蓝色提示：提示为“发票号[xxx]重复，请检查；”（xxx是rcptNo值）
		if arrays.Contains(rcptNoName, rcptNo) != -1 {
			if strings.Index(wrongNote, "发票号["+rcptNo+"]重复，请检查；") == -1 {
				wrongNote += "发票号[" + rcptNo + "]重复，请检查；"
			}
		} else {
			rcptNoName = append(rcptNoName, rcptNo)
		}
	}

	return wrongNote
}

func CheckFieldHasIssue(fields []model.ProjectField, code string) bool {
	for _, field := range fields {
		if field.Code == code {
			if len(field.Issues) > 0 {
				return true
			}
		}
	}
	return false
}
