package B0103

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	"server/module/export/service"
	"server/module/load/model"
	model3 "server/module/load/model"
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
	global.GLog.Info("B0103:::CheckXml")
	obj := o.(FormatObj)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	// ----------------------------------code---------------------------------------------
	constMap := InitConst(obj.Bill.ProCode, obj.Bill.Agency)
	wrongNote += XmlCheck(obj, xmlValue, constMap)
	wrongNote += CodeCheck(obj, xmlValue, constMap)

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
	// sanfang := constMap["sanfang"]
	// fmt.Println("--------------sanfang----------------------", sanfang)
	//jibin := constMap["jibin"]
	//特殊的常量-全部医疗目录
	//constSpecialMap := constSpecialDeal(obj.Bill.ProCode)
	var fc005Arr []string
	var fc091Arr []string
	for _, invoice := range obj.Invoice {
		// _, fc005 := GetOneField(invoice.Fields, "fc005", true)
		_, f005 := GetOneField(invoice.Fields, "fc005", true)
		var resultNames []string
		var filedArr []string
		fmt.Println("----f005-----", f005)
		fmt.Println("----invoice.QingDan-----", len(invoice.QingDan))
		for _, qingDan := range invoice.QingDan {
			for i, field := range qingDan.Fields {
				//_, fc008 := GetOneField(invoice.Fields, "fc008", true)
				//bianCode, ok := constSpecialMap["shuJuKuBianMaMap"][obj.Bill.Agency]
				muLu, total := utils.FetchConst(obj.Bill.ProCode, "B0103_广西贵州国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": obj.Bill.Agency})
				fmt.Println("==============muLu=======", muLu)
				//if ok && len(bianCode) > 3 {
				if total > 0 {
					yiLiaoMuLu := "医疗目录" + muLu
					//yiLiaoMuLu := "医疗目录" + bianCode[3]
					otherInfo := make(map[string]interface{})
					err := json.Unmarshal([]byte(obj.Bill.OtherInfo), &otherInfo)
					if err != nil {
						global.GLog.Error("CSB0103RC0056000", zap.Error(err))
					}
					//yiLiaoMuLuConst, _ := constSpecialMap[yiLiaoMuLu]
					if RegIsMatch(field.Code, `^(fc012|fc025|fc026|fc027|fc028|fc029|fc030|fc031)$`) {
						yiLiaoMuLu = "B0103_广西贵州国寿理赔_" + yiLiaoMuLu
						_, total1 := utils.FetchConst(obj.Bill.ProCode, yiLiaoMuLu, "医疗项目编码", map[string]string{"项目名称": field.ResultValue})
						//_, ok := yiLiaoMuLuConst[field.ResultValue]
						//if !ok && field.ResultValue != "" {

						if total1 == 0 && field.ResultValue != "" {
							resultNames = append(resultNames, "【"+field.ResultValue+"】")
							//wrongNote += field.Name + "【" + field.ResultValue + "】不为常量库内容，请修改;"
						}
					}
				}
				if len(resultNames) > 0 {
					if i == len(qingDan.Fields)-1 {
						join := strings.Join(resultNames, "、")
						mes := "发票" + f005 + "项目名称" + join + "不为常量库内容，请修改;"
						if strings.Index(wrongNote, mes) == -1 {
							wrongNote += mes
						}
					}
				}

				//需求编码 CSB0103RC0095000 清单
				if RegIsMatch(field.ResultValue, `[?？]`) {
					filedArr = append(filedArr, field.Name)
				}
			}

			//CSB0103RC0125000
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
			var cSB0103RC0125000Msg []string
			for _, item := range sFields {
				_, name := GetOneField(qingDan.Fields, item[0], false)
				_, money := GetOneField(qingDan.Fields, item[1], false)
				_, number := GetOneField(qingDan.Fields, item[2], false)
				if !(len(name) == 0 && len(money) == 0 && len(number) == 0) &&
					!(len(name) > 0 && len(money) > 0 && len(number) > 0) {
					for _, value := range item {
						_, val := GetOneField(qingDan.Fields, value, true)
						if len(val) == 0 {
							cSB0103RC0125000Msg = append(cSB0103RC0125000Msg, value)
						}
					}
				}
			}
			if len(cSB0103RC0125000Msg) > 0 {
				mes := "【" + f005 + "】发票【" + strings.Join(cSB0103RC0125000Msg, "、") + "】字段位为空；"
				if strings.Index(wrongNote, mes) == -1 {
					wrongNote += mes
				}
			}
		}

		_, fc010 := GetOneField(invoice.Fields, "fc010", false)
		_, fc005 := GetOneField(invoice.Fields, "fc005", false)
		fc016 := ""
		fc062 := ""
		fc063 := ""
		if len(invoice.ThirdBaoXiaoDan) == 0 {
			sum := 0.0
			_, fc062 = GetOneField(invoice.Fields, "fc062", true)
			for _, field := range invoice.BaoXiaoDan {
				if RegIsMatch(field.Code, `^(fc015|fc016|fc063|fc065|fc067|fc069|fc071|fc073)$`) && field.ResultValue != "" && field.ResultValue != "A" {
					sum = SumFloat(sum, ParseFloat(field.ResultValue), "+")
				}
				if field.Code == "fc016" {
					fc016 = field.ResultValue
				}
				// 编码 CSB0103RC0070001 取消 第二点
				//if fc062 != "" && RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) {
				//	if HasKey(constMap["sanfangAll"], field.ResultValue) == false && field.ResultValue != "" {
				//		wrongNote += "发票号[" + fc005 + "]第三方出具单位录入有误，请检查;"
				//	}
				//}
			}
			if sum != ParseFloat(fc010) && fc010 != "" && fc016 != "" {
				wrongNote += "发票号[" + fc005 + "]，发票报销金额与结算单报销金额加第三方报销金额之和不一致，请确认;"
			}
		} else if len(invoice.BaoXiaoDan) != 0 {
			sum := 0.0
			for _, field := range invoice.BaoXiaoDan {
				if RegIsMatch(field.Code, `^(fc015|fc016)$`) && field.ResultValue != "" && field.ResultValue != "A" {
					sum = SumFloat(sum, ParseFloat(field.ResultValue), "+")
				}
				if field.Code == "fc016" {
					fc016 = field.ResultValue
				}
				if field.Code == "fc062" {
					fc062 = field.ResultValue
				}
				if field.Code == "fc063" {
					fc063 = field.ResultValue
				}
				// 编码 ：CSB0103RC0096000
				//校验所有fc062、fc064、fc066、fc068、fc070、fc072的录入值，不为常量表《B0103_广西贵州国寿理赔_第三方出具单位》中的“第三方出具单位名称”（第三列）的内容时，出导出校验：发票【xxx】第三方出具单位录入有误，请检查；
				//（字段录入值为A、为空时不执行该校验，xxx是对应出提示的发票中fc005的值；）
				// CSB0103RC0096001
				if strings.Index(wrongNote, "发票号【"+fc005+"】第三方出具单位录入有误，请检查；") == -1 {
					if RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) && field.ResultValue != "" && field.ResultValue != "A" {
						_, total := utils.FetchConst(obj.Bill.ProCode, "B0103_广西贵州国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": field.ResultValue})
						//isBool := HasKey(constMap["sanfangAll"], field.ResultValue)

						if total == 0 {
							wrongNote += "发票号【" + fc005 + "】第三方出具单位录入有误，请检查；"
						}
					}
				}

			}
			for _, field := range invoice.ThirdBaoXiaoDan {

				if RegIsMatch(field.Code, `^(fc063)$`) && field.ResultValue != "" && field.ResultValue != "A" {
					sum = SumFloat(sum, ParseFloat(field.ResultValue), "+")
				}
				//编码 CSB0103RC0070001  取消 第一点
				//if field.Code == "fc062" && (fc062 == "" && fc063 == "") {
				//isBool := HasKey(constMap["sanfang"], field.FinalValue)
				//	if isBool == false && field.ResultValue != "" {
				//		wrongNote += "发票号[" + fc005 + "]第三方出具单位录入有误，请检查;"
				//	}
				//}
				//if field.Code == "fc062" && field.ResultValue != "" && fc062 != "" {
				//	wrongNote += "发票[" + fc005 + "]重复录入第三方支付数据，请检查;"
				//}
				//if field.Code == "fc063" && field.ResultValue != "" && fc063 != "" {
				//	wrongNote += "发票[" + fc005 + "]重复录入第三方支付数据，请检查;"
				//}
				if field.Code == "fc062" && field.ResultValue != "" && fc062 != "" && fc062 != field.ResultValue {
					wrongNote += "发票[" + fc005 + "]重复录入第三方支付数据，请检查;"
				}
				if field.Code == "fc063" && field.ResultValue != "" && fc063 != "" && fc063 != field.ResultValue {
					wrongNote += "发票[" + fc005 + "]重复录入第三方支付数据，请检查;"
				}

			}

			if sum != ParseFloat(fc010) && fc010 != "" && fc016 != "" {
				wrongNote += "发票号[" + fc005 + "]，发票报销金额与结算单报销金额加第三方报销金额之和不一致，请确认;"
			}
		}
		// 编码 ：CSB0103RC0096000
		//校验所有fc062、fc064、fc066、fc068、fc070、fc072的录入值，不为常量表《B0103_广西贵州国寿理赔_第三方出具单位》中的“第三方出具单位名称”（第三列）的内容时，出导出校验：发票【xxx】第三方出具单位录入有误，请检查；
		//（字段录入值为A、为空时不执行该校验，xxx是对应出提示的发票中fc005的值；）

		if len(invoice.ThirdBaoXiaoDan) != 0 {
			for _, field := range invoice.ThirdBaoXiaoDan {
				if strings.Index(wrongNote, "发票号【"+fc005+"】第三方出具单位录入有误，请检查；") == -1 {
					if RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) && field.ResultValue != "" && field.ResultValue != "A" {
						_, total := utils.FetchConst(obj.Bill.ProCode, "B0103_广西贵州国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": field.ResultValue})
						//isBool := HasKey(constMap["sanfangAll"], field.ResultValue)

						if total == 0 {
							mes := "发票号【" + fc005 + "】第三方出具单位录入有误，请检查；"
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += "发票号【" + fc005 + "】第三方出具单位录入有误，请检查；"
							}
						}
					}
				}
			}
		}

		for _, field := range fields {
			if RegIsMatch(field.Code, `^(fc007|fc074|fc075|fc076|fc077|fc078|fc079|fc080|fc081|fc082)$`) {
				if field.ResultValue == "妊娠合并低蛋白血症" {
					wrongNote += "“妊娠合并低蛋白血症”错误诊断请改为“低(白)蛋白血(症)”;"
				} else if field.ResultValue == "EB病毒遗传缺陷反应后的免疫缺陷" {
					wrongNote += "“EB病毒遗传缺陷反应后的免疫缺陷”错误诊断请改为“ＥＢ病毒感染”;"
				} else if field.ResultValue == "唇裂" {
					wrongNote += "“唇裂”错误诊断，请确认;"

				}
				//matched, _ := regexp.MatchString(`^(A|B|\?|？|)$`, field.ResultValue)
				//if field.ResultValue != "" && !matched {
				//	//_, total := utils.FetchConst(obj.Bill.ProCode, "B0103_广西贵州国寿理赔_ICD10疾病编码", "疾病代码", map[string]string{"疾病名称": field.ResultValue})
				//	//if !RegIsMatch(field.ResultValue, `^(|A|B|\?)$`) && !HasKey(jibin, field.ResultValue) {
				//	//	wrongNote += "发票号[" + fc005 + "]，核实疾病代码录入是否准确;"
				//	//}
				//	//mes := "发票号[" + fc005 + "]，核实疾病代码录入是否准确;"
				//	//if strings.Index(wrongNote, mes) == -1 {
				//	//	if total == 0 {
				//	//		wrongNote += "发票号[" + fc005 + "]，核实疾病代码录入是否准确;"
				//	//	}
				//	//}
				//}
			}

			if field.Code == "fc007" {
				if field.ResultValue == "" {
					wrongNote += "发票号[" + fc005 + "]，疾病诊断不能为空;"
				}
			}

		}

		//编码：CSB0103RC0095000 校验所有字段的录入值，当包含?或？时，出导出校验：发票【xxx】的【yyy】存在?号，请核实；
		//如同一发票下多个字段包含?或？，则导出校验提示格式为：发票【xxx】的【yyy1、yyy2】存在?号，请核实；（xxx为发票号fc005的值，yyy为包含问号的字段名）
		//日期：20230711
		// 取消需求 导出校验 最后一页存在?，请检查;
		//if j == len(obj.Invoice)-1 {
		//	for _, field := range invoice.Fields {
		//		if RegIsMatch(field.ResultValue, `[?？]`) {
		//			if strings.Index(wrongNote, "最后一页存在?，请检查;") == -1 {
		//				wrongNote += "最后一页存在?，请检查;"
		//			}
		//		}
		//	}
		//}

		for i, field := range invoice.Fields {

			if RegIsMatch(field.ResultValue, `[?？]`) {
				filedArr = append(filedArr, field.Name)
			}

			if i == len(invoice.Fields)-1 {
				if len(filedArr) > 0 {
					join := strings.Join(filedArr, "、")
					mes := "发票" + f005 + "的【" + join + "】存在?号，请核实;"
					if strings.Index(wrongNote, mes) == -1 {
						fmt.Println("--filedArr.=", filedArr)
						wrongNote += mes
						filedArr = make([]string, 0)
						fmt.Println("--filedArr null.=", filedArr)
					}
				}
			}

			//CSB0103RC0108000
			//同一发票下，当fc090录入值为1时，校验fc092、fc093、fc094、fc095、fc096是否存在结果值，任意一个字段结果值为空时，出导出校验：电子发票五要素未录入齐全，请确认是否为电子发票；
			sFiled := []string{"fc092", "fc005", "fc094", "fc095", "fc096"}
			if field.Code == "fc090" {
				_, fc091Val := GetOneField(invoice.Fields, "fc091", false)
				if fc091Val != "3" && field.ResultValue == "1" {
					for k := 0; k < len(sFiled); k++ {
						_, val := GetOneField(invoice.Fields, sFiled[k], true)
						if val == "" {
							mes := "【" + f005 + "】" + "电子发票五要素未录入齐全，请确认是否为电子发票；"
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}
					}
				}
			}
			//CSB0103RC0109000
			//"同一发票下，当fc091录入值为1时，（xxx为fc005的结果值）
			//1、校验fc092的结果值是否为12位数，否则出导出校验：【XX】发票增值税票据代码字段应为12位数，请检查；
			//2、校验fc093的结果值是否为8位数，否则出导出校验：【XX】发票增值税票据号码字段应为8位数，请检查；
			//3、校验fc095的结果值是否为6位数，否则出导出校验：【XX】发票校验码字段应为6位数，请检查；"
			//CSB0103RC0110000
			//"同一发票下，当fc091录入值为2时，（xxx为fc005的结果值）
			//1、校验fc092的结果值是否为8位数，否则出导出校验：XX】发票财政部票据代码字段应为8位数，请检查；
			//2、校验fc093的结果值是否为10位数，否则出导出校验：【XX】发票财政部票据号码字段应为10位数，请检查；"
			sFiledTo := [][]string{{"fc092", "12", "发票增值税票据代码字段应为12位数，请检查；"}, {"fc005", "8", "发票增值税票据号码字段应为8位数，请检查；"}, {"fc095", "6", "发票校验码字段应为6位数，请检查；"}}
			sFiledThree := [][]string{{"fc092", "8", "发票财政部票据代码字段应为8位数，请检查；"}, {"fc005", "10", "发票财政部票据号码字段应为10位数，请检查；"}}
			if field.Code == "fc091" {
				if field.ResultValue == "1" {
					for k := 0; k < len(sFiledTo); k++ {
						_, val := GetOneField(invoice.Fields, sFiledTo[k][0], true)
						atoi, _ := strconv.Atoi(sFiledTo[k][1])
						if val != "" && len(val) != atoi {
							mes := "【" + f005 + "】" + sFiledTo[k][2]
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}
					}

					// _, val := GetOneField(invoice.Fields, "fc097", true)
					// if val != "01" {
					// 	mes := "【" + fc005 + "】" + "发票查验结果不通过，请确认五要素是否录入正确；"
					// 	if strings.Index(wrongNote, mes) == -1 {
					// 		wrongNote += mes
					// 	}
					// }
				} else if field.ResultValue == "2" {
					for k := 0; k < len(sFiledThree); k++ {
						_, val := GetOneField(invoice.Fields, sFiledThree[k][0], true)
						atoi, _ := strconv.Atoi(sFiledThree[k][1])
						if val != "" && len(val) != atoi {
							mes := "【" + f005 + "】" + sFiledThree[k][2]
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}
					}
				}
			}

			//CSB0103RC0114001
			if field.Code == "fc090" && field.ResultValue == "1" {
				_, val := GetOneField(invoice.Fields, "fc097", true)
				_, fc091Val := GetOneField(invoice.Fields, "fc091", false)
				_, fc094Val := GetOneField(invoice.Fields, "fc094", true)
				fc094Parse, _ := time.Parse("2006-01-02", fc094Val)
				subTime := time.Now().Sub(fc094Parse)
				oneYear, _ := time.ParseDuration("8760h")
				if fc091Val == "2" {
					if subTime < oneYear && val != "01" {
						fc005Arr = append(fc005Arr, f005)
					}
				} else {
					if val != "01" {
						fc005Arr = append(fc005Arr, f005)
					}
				}
			}

			//CSB0103RC0122000
			if field.Code == "fc091" && field.ResultValue == "3" {
				_, val := GetOneField(invoice.Fields, "fc005", false)
				if len(val) < 20 {
					fc091Arr = append(fc091Arr, f005)
				}
			}

			//fc007|fc074|fc075|fc076|fc077|fc078|fc079|fc080|fc081|fc082 录入值不为<B0103_广西贵州国寿理赔_ICD10疾病编码>的疾病名称(第一列)且不为空值,A,B,?时出导出校验发票号[XXX]，核实疾病代码录入是否准确;
			if RegIsMatch(field.Code, `^(fc007|fc074|fc075|fc076|fc077|fc078|fc079|fc080|fc081|fc082)$`) {
				matched, _ := regexp.MatchString(`^(A|B|\?|？|)$`, field.ResultValue)
				if field.ResultValue != "" && !matched {
					_, total := utils.FetchConst(obj.Bill.ProCode, "B0103_广西贵州国寿理赔_ICD10疾病编码", "疾病代码", map[string]string{"疾病名称": field.ResultValue})
					//if !RegIsMatch(field.ResultValue, `^(|A|B|\?)$`) && !HasKey(jibin, field.ResultValue) {
					//	wrongNote += "发票号[" + fc005 + "]，核实疾病代码录入是否准确;"
					//}
					mes := "发票号[" + fc005 + "]，核实疾病代码录入是否准确;"
					if strings.Index(wrongNote, mes) == -1 {
						if total == 0 {
							wrongNote += "发票号[" + fc005 + "]，核实疾病代码录入是否准确;"
						}
					}
				}
			}
		}
	}

	if len(fc005Arr) > 0 {
		join := strings.Join(fc005Arr, "、")
		mes := "【" + join + "】" + "发票查验结果不通过，请确认五要素是否录入正确；"
		if strings.Index(wrongNote, mes) == -1 {
			wrongNote += mes
		}
	}

	if len(fc091Arr) > 0 {
		join := strings.Join(fc091Arr, "、")
		mes := "【" + join + "】" + "发票票据号不为20位数，请检查；"
		if strings.Index(wrongNote, mes) == -1 {
			wrongNote += mes
		}
	}

	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""
	invoice := obj.Invoice
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
	fc056Issue := false //问题件

	items := []string{}
	//rcptNos := RegMatchAll(xmlValue, `<rcptNo>.*?<\/rcptNo>`)
	//for _, rcptNo := range rcptNos {
	//	//if !RegIsMatch(rcptNo, `^\d+$`) {
	//	//	wrongNote += "发票号[" + rcptNo + "]票据号录入格式错误，请核实;"
	//	//}
	//	if arrays.Contains(items, rcptNo) != -1 {
	//		wrongNote += "发票号[" + rcptNo + "]重复，请检查;"
	//	} else {
	//		items = append(items, rcptNo)
	//	}
	//}

	rcptInfoLists := RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	var endDateLast time.Time
	lianxu := map[string]string{}
	filename := true
	imageType := false
	for _, rcptInfoList := range rcptInfoLists {
		rcptNo := GetNodeValue(rcptInfoList, "rcptNo")
		endDate := GetNodeValue(rcptInfoList, "endDate")
		beginDate := GetNodeValue(rcptInfoList, "beginDate")
		indexRcptNo := strings.Index(rcptInfoList, "</rcptNo>")
		// CSB0103RC0057001
		if arrays.Contains(items, rcptNo) != -1 {
			mes := "发票号[" + rcptNo + "]重复，请检查;"
			if strings.Index(wrongNote, mes) == -1 {
				wrongNote += "发票号[" + rcptNo + "]重复，请检查;"
			}
		} else {
			items = append(items, rcptNo)
		}
		a, _ := time.Parse("2006-01-02", endDate)
		b, _ := time.Parse("2006-01-02", beginDate)
		if a.After(endDateLast) {
			endDateLast = a
		}
		if a.Before(b) {
			wrongNote += "发票号[" + rcptNo + "]出院时间早于入院时间，请检查；"
		}

		expenMode := GetNodeValue(rcptInfoList, "expenMode")
		d := a.Sub(b)
		aa := d.Hours() / 24
		if (expenMode == "0020" && aa != 0) || (expenMode == "0040" && aa <= 0) || time.Now().Before(a) || time.Now().Before(b) {
			wrongNote += "发票号[" + rcptNo + "]的出入院日期异常，请核实；"
		}

		socialPayAmnt := GetNodeValue(rcptInfoList, "socialPayAmnt")
		additionalPayAmnt := GetNodeValue(rcptInfoList, "additionalPayAmnt")

		cmsImageInfoList := otherInfo.CmsImageInfoList
		for _, list := range cmsImageInfoList {
			if list.ImageType == "3032" {
				imageType = true
				if ParseFloat(socialPayAmnt) > 0 || ParseFloat(additionalPayAmnt) > 0 {
					filename = false
				}
			}
		}
		//if ParseFloat(socialPayAmnt) > 0 || ParseFloat(additionalPayAmnt) > 0 {
		//	filename = false
		//}
		rcptAmnt := GetNodeValue(rcptInfoList, "rcptAmnt")
		lianxu[rcptNo] = rcptAmnt
		if ParseFloat(socialPayAmnt) >= ParseFloat(rcptAmnt) && indexRcptNo != -1 {
			wrongNote += "发票号[" + rcptNo + "]报销金额大于或等于发票总金额，请检查;"
		}
		//CSB0103RC0066000 判断<socialPayAmnt>为非数字或小于0时(空值不进行校验)，出导出校验：“发票号[xxx]社保报销金额为无效数值，请核实；”，xxx是同一个rcptInfoList中的rcptNo值
		if socialPayAmnt != "" && !RegIsMatch(socialPayAmnt, `^[\d\.]+$`) || ParseFloat(socialPayAmnt) < 0 {
			wrongNote += "发票号[" + rcptNo + "]社保报销金额为无效数值，请核实；"
		}

		siHealthType := GetNodeValue(rcptInfoList, "siHealthType")
		accountPayAmnt := GetNodeValue(rcptInfoList, "accountPayAmnt")
		if siHealthType == "00" && (!RegIsMatch(socialPayAmnt, `^(|0\.00|0)$`) || !RegIsMatch(accountPayAmnt, `^(|0\.00|0)$`)) {
			wrongNote += "发票号[" + rcptNo + "]的医保身份或社保支付金额有误，请确认;"
			fc056Issue = true
		}
		for _, invoices := range invoice {
			for _, field := range invoices.Fields {
				if field.Code == "fc056" {
					if fc056Issue == true {
						issueMap := model3.Issue{
							Message: "医保身份未填写内容",
						}
						field.Issues = append(field.Issues, issueMap)
					}
				}
			}
		}

		rcptLists := RegMatchAll(rcptInfoList, `<(rcptList|errorRcptList)>[\s\S]*?<\/(rcptList|errorRcptList)>`)
		sum := 0.0
		difference := ""
		for ii, rcptList := range rcptLists {
			// listName := GetNodeValue(rcptList, "listName")
			// if RegIsMatch(listName, `(第三方报销|全自费费用|先行自付费用|社保统筹范围内费用|不予支付部分|目录外费用|超限价自付费用|公务员报销|床位费补助|企补支付|大病报销)`) {
			// 	wrongNote += "发票号" + rcptNo + "报销单录入错误；"
			// }
			// if strings.Index(rcptList, "<errorRcptList>") != -1 && RegIsMatch(listName, `^(6501000000014689|6527000000014708|6532000000014698|6541000000014699|6542000000014698|6590000000014708)$`) {
			// 	continue wrongNote += "[" + price + "]为空，请检查;"
			// }
			quantity := GetNodeValue(rcptList, "quantity")
			listName := GetNodeValue(rcptList, "listName")
			price := GetNodeValue(rcptList, "price")
			priceIndex := strings.Index(rcptList, "price")       //price
			quantityIndex := strings.Index(rcptList, "quantity") //quantity
			if quantity == "" && quantityIndex != -1 {
				if strings.Index(wrongNote, "发票号["+rcptNo+"]数量/金额为空，请检查;") == -1 {
					wrongNote += "发票号[" + rcptNo + "]数量/金额为空，请检查;"
				}
			}
			if price == "" && priceIndex != -1 {
				if strings.Index(wrongNote, "发票号["+rcptNo+"]数量/金额为空，请检查;") == -1 {
					wrongNote += "发票号[" + rcptNo + "]数量/金额为空，请检查;"
				}
			}

			if ParseFloat(price) > ParseFloat(rcptAmnt) {
				wrongNote += "发票号[" + rcptNo + "]费用明细大于发票总金额;"
			}
			if quantity != "" && price != "" && listName != "补差金额" {
				total := SumFloat(ParseFloat(quantity), ParseFloat(price), "*")
				sum = sum + total
			}
			if len(rcptLists)-1 == ii {
				price = GetNodeValue(rcptList, "price")
				fmt.Println("len(rcptLists)-1 == ii = ", price, difference)
				difference = price
			}
		}

		chae := SumFloat(ParseFloat(rcptAmnt), sum, "-")
		// expenMode := GetNodeValue(rcptInfoList, "expenMode")
		//float, _ := strconv.ParseFloat(difference, 64)
		if sum != 0.0 && (chae >= 2 || chae <= -2) {
			wrongNote += "发票号[" + rcptNo + "]清单项目金额与发票金额不一致，差额:" + strconv.Itoa(int(chae)) + "，请检查；"
		}
	}

	wmes := ""
	for key, value := range lianxu {
		if RegIsMatch(key, `^\d+$`) {
			nKey, _ := strconv.Atoi(key)
			nKey = nKey + 1
			rcptAmnt, ok := lianxu[strconv.Itoa(nKey)]
			if ok && value == rcptAmnt {
				rcptNo := "发票号[" + key + "]"
				if strings.Index(wmes, rcptNo) == -1 {
					wmes += rcptNo
				}
				rcptNo = "发票号[" + strconv.Itoa(nKey) + "]"
				if strings.Index(wmes, rcptNo) == -1 {
					wmes += rcptNo
				}
			}
		}
	}
	if wmes != "" {
		wrongNote += wmes + "连续且总金额相同，请确认是否为连续发票;"
	}
	if filename && imageType {
		wrongNote += wmes + "请确认报销金额、统筹支付、附加支付是否录入正确;"
	}

	dd := time.Now().Sub(endDateLast)
	if (dd.Hours() / 24) > 730 {
		wrongNote += "超索赔期限，请核实出院日期是否准确;"
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
