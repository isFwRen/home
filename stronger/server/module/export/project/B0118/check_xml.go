/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/26 6:01 下午
 */

package B0118

import (
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/service"
	"server/module/load/model"
	proModel "server/module/pro_conf/model"
	utils2 "server/utils"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/wxnacy/wgo/arrays"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	global.GLog.Info("B0118:::CheckXml")

	obj := o.(FormatObj)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	// ----------------------------------code---------------------------------------------
	constMap := constDeal(obj.Bill.ProCode)
	wrongNote += CodeCheck(obj, xmlValue, constMap)
	wrongNote += XmlCheck(obj, xmlValue, constMap)

	// ----------------------------------xml---------------------------------------------

	global.GLog.Error("wrongNote：：：" + wrongNote)
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
	// constMap["yiYuanDaiMaMap"] := map[string]string{}
	fc062s := []string{}
	fc149 := "null"
	ccffs := [][]string{
		{"fc084", "fc154", "fc092", "fc162", "fc172"},
		{"fc085", "fc155", "fc093", "fc163", "fc173"},
		{"fc086", "fc156", "fc094", "fc164", "fc174"},
		{"fc087", "fc157", "fc095", "fc165", "fc175"},
		{"fc088", "fc158", "fc096", "fc166", "fc176"},
		{"fc089", "fc159", "fc097", "fc167", "fc177"},
		{"fc090", "fc160", "fc098", "fc168", "fc178"},
		{"fc091", "fc161", "fc099", "fc169", "fc179"},
	}

	for _, field := range fields {
		if obj.Bill.ClaimType == 1 {
			if (field.Code == "fc180" && field.ResultValue == "徐州") || (field.Code == "fc054" && strings.Index(field.ResultValue, "徐州") != -1) {
				wrongNote += "江苏徐州市的结算单模版，社保自费请按单录入“先行支付” + “自费”;"
			}
		}
		if RegIsMatch(field.Code, `^(fc191|fc191|fc192|fc193|fc194|fc195|fc196|fc197|fc198|fc199|fc200)$`) && field.FinalValue != "" {
			if !HasKey(constMap["jiBingZhenDuanDaiMaMap"], field.FinalValue) {
				wrongNote = wrongNote + field.Name + "疾病代码录入有误，请检查；"
			} else if RegIsMatch(field.FinalValue, `\..+$`) {
				ffinput := Substr(field.FinalValue, 0, strings.Index(field.FinalValue, ".")+2)
				if !HasKey(constMap["jiBingZhenDuanDaiMaMap"], ffinput) {
					wrongNote = wrongNote + field.Name + "疾病代码录入有误，请检查；"
				}
			}
		}
		if RegIsMatch(field.Code, `^(fc053|fc182|fc183|fc184|fc185|fc186|fc187|fc188|fc189|fc190)$`) && field.FinalValue != "" && !HasKey(constMap["jiBingZhenDuanDaiMaMap"], field.FinalValue) {
			wrongNote = wrongNote + field.Name + "疾病诊断错误，请检查；"
		}
		if RegIsMatch(field.Code, `^(fc053)$`) && field.FinalValue == "" {
			wrongNote = wrongNote + field.Name + "疾病诊断不能为空，请检查；"
		}
		if RegIsMatch(field.Code, `^(fc062)$`) && field.FinalValue != "" && arrays.ContainsString(fc062s, field.FinalValue) == -1 {
			fc062s = append(fc062s, field.FinalValue)
		}
		if RegIsMatch(field.Code, `^(fc149)$`) {
			fc149 = field.FinalValue
		}

		// if RegIsMatch(field.Code, `^(fc084|fc154|fc092|fc162|fc172|fc085|fc155|fc093|fc163|fc173|fc086|fc156|fc094|fc164|fc174|fc087|fc157|fc095|fc165|fc175|fc088|fc158|fc096|fc166|fc176|fc089|fc159|fc097|fc167|fc177|fc090|fc160|fc098|fc168|fc178|fc091|fc161|fc099|fc169|fc179)$`) && (strings.Index(field.ResultValue, "?") != -1 || strings.Index(field.ResultValue, "？") != -1) {
		// 	// fmt.Println("---------------清单内容不能录入问号，请检查------------------:", field.ResultValue)
		// 	wrongNote += "清单内容不能录入问号，请检查；"
		// }
		//fc103 := ""
		//if field.Code == "fc103" {
		//	fc103 = field.FinalValue
		//}
		//if RegIsMatch(field.Code, `^(fc154|fc155|fc156|fc157|fc158|fc159|fc160|fc161)$`) && field.ResultValue != "" && (ParseFloat(field.ResultValue) > 1 || ParseFloat(field.FinalValue) > 1) {
		//	wrongNote += fc103 + "账单号" + field.Name + "的清单项目类型录入错误，请检查；"
		//}

		//if RegIsMatch(field.Code, `^(fc154|fc155|fc156|fc157|fc158|fc159|fc160|fc161)$`) && field.FinalValue != "" {
		//	if !RegIsMatch(field.FinalValue, `^(1|2|3|4|5|6)$`) {
		//		wrongNote += fc103 + "账单号" + field.Name + "的清单项目类型录入错误，请检查;"
		//	}
		//}
		// if field.Code == "fc180" && (obj.Bill.ClaimType == 5 || obj.Bill.ClaimType == 6) && field.ResultValue == "F" && strings.Index(wrongNote, "fc180出险地点市为F，请修改；") == -1 {
		// 	wrongNote += "fc180出险地点市为F，请修改；"
		// }
		if field.Code == "fc180" && field.ResultValue == "佛山" && strings.Index(wrongNote, "案件为佛山案件，请确认报销金额是否等于总金额-个人支付；") == -1 {
			wrongNote += "案件为佛山案件，请确认报销金额是否等于总金额-个人支付；"
		}

	}

	// siLists := obj.SiList
	// for _, siList := range siLists {
	// 	fields := siList.Fields
	// 	for _, field := range fields {
	// 		if field.Code == "fc067" {
	// 			fmt.Println("--------------fc067------------------:", field.Code, field.ResultValue)
	// 			wrongNote += siList.Code + "账单号的社保自费为0，请确认;"
	// 		}
	// 	}
	// }

	hospitalList := obj.HospitalList
	opList := obj.OpList
	fapiaos := append(hospitalList, opList...)
	fc002Maps := make(map[string][]string)
	for _, fapiao := range fapiaos {
		qingDans := fapiao.QingDan
		code := fapiao.Code
		num14 := 0
		names := ""
		arrs := []string{}
		fc052 := ""
		fc206 := ""
		fc059 := ""
		fc052Len := 0
		for _, field := range fapiao.Fields {
			if field.Code == "fc052" {
				fc052 = field.ResultValue
				fc052Len++
			}
			if field.Code == "fc059" {
				fc059 = field.ResultValue
			}
			if field.Code == "fc206" {
				fc206 = field.ResultValue
			}
			if field.Code == "fc180" && (obj.Bill.ClaimType == 5 || obj.Bill.ClaimType == 6) && field.ResultValue == "F" && strings.Index(wrongNote, fapiao.Code+"账单号的fc180出险地点市为F，请修改；") == -1 {
				wrongNote += fapiao.Code + "账单号的fc180出险地点市为F，请修改；"
			}
		}

		fc002Map, isExits := fc002Maps[code]
		if !isExits {
			fc002Maps[code] = []string{}
		}
		fc002Maps[code] = append(fc002Map, fc059)

		if fc052 != fc206 && fc052Len > 0 {
			// wrongNote += "发票" + code + "报销金额与结算单报销金额不一致，请确认；"
		}

		//当一张单子中，左、右列字段同时存在时，需对比两个分块中的录入值是否一致：
		//同一行的两个字段视为一组对比数据，当一组字段的录入值均为空或录入值一致或fc205录入为A,fc206、fc207、fc208录入为空时，则不出导出校验，如一组字段有录入值且录入内容不一致时，需出导出校验：xxx账单号，发票报销内容与结算单报销内容录入不一致，请修改；（如fc205录入1，fc101录入2，则需出导出校验）
		//bc001  bc003
		fArr := [][]string{
			{"fc205", "fc101"},
			{"fc206", "fc052"},
			{"fc207", "fc051"},
			{"fc208", "fc079"},
		}
		flag := false
		for _, f := range fArr {
			is0, f0Val := GetFieldsInputExits(fapiao.Fields, f[0])
			is1, f1Val := GetFieldsInputExits(fapiao.Fields, f[1])
			if f0Val == "" && f1Val == "" {
				continue
			}
			if !is0 || !is1 {
				continue
			}
			if f0Val != f1Val {
				flag = true
			}
		}

		fc205Val := GetFieldsInput(fapiao.Fields, fArr[0][0])
		fc206Val := GetFieldsInput(fapiao.Fields, fArr[1][0])
		fc207Val := GetFieldsInput(fapiao.Fields, fArr[2][0])
		fc208Val := GetFieldsInput(fapiao.Fields, fArr[3][0])
		if fc205Val == "A" && fc206Val == "" && fc207Val == "" && fc208Val == "" {
			flag = false
		}

		if flag {
			msg := "发票" + code + "报销金额与结算单报销金额不一致，请确认；"
			if strings.Index(wrongNote, msg) == -1 {
				wrongNote += msg
			}
		}

		fc067 := ""
		fc008 := fapiao.Money
		sum := decimal.NewFromFloat(0.0)
		fmt.Println("---------------money------------------:", fapiao.Money)
		for _, field := range fapiao.Fields {
			if field.Code == "fc067" {
				// fmt.Println("---------------fc067------------------:", field.ResultValue)
				fc067 = field.ResultValue
				//2023年09月07日15:54:40
				//CSB0118RC0283001
				fc180Val := GetFieldsInput(fapiao.Fields, "fc180")
				if fc067 == "0" && fc180Val != "徐州" {
					wrongNote += code + "账单号的社保自费为0，请确认;"
				}
			}
			// if field.Code == "fc008" {
			// 	fmt.Println("---------------fc008fc008fc008------------------:", field.FinalInput)
			// }
			if field.Code == "fc008" && field.FinalInput != "" {
				fc008 = ParseFloat(field.FinalInput)
			}
			if RegIsMatch(field.Code, `^(fc010|fc012|fc014|fc016|fc018|fc020|fc022|fc024|fc026|fc028|fc030|fc032|fc034|fc036|fc038|fc040|fc042|fc044|fc046|fc048)$`) && field.FinalValue != "" {
				a1, _ := decimal.NewFromString(field.FinalValue)
				sum = a1.Add(sum)
				// a2.Tan().BigFloat().Float64()
				// fff, _ := decimal.NewFromFloat(sum).Add()
				// sum += ParseFloat(field.FinalValue)
				// fmt.Println("---------------sum------------------:", sum, field.FinalValue, ParseFloat(field.FinalValue), fmt.Sprintf("%.8f", field.FinalValue))
				// sum =decimal.Sum(sum, decimal.NewFromString(field.FinalValue))
			}
		}
		fc153 := 0.0
		fc153s := []string{}

		isfc092 := false
		fc092All := 0.0
		fc153All := 0.0
		fc010All := 0.0

		//在同一发票中(根据fc059发票属性、fc060清单所属属性判断是否属于同一张发票)对下列分组中左边字段录入内容进行校验,
		//左边字段中录入内容包含"走廊"、"床位"、"人间"、"病房"、"陪护床"字样时，对对应的右边字段的录入内容进行合计，
		//若以上字样的项目名称中包含"取暖"或"空调"字样时，该项目名称不作统计；
		//fc084   fc092
		//fc085   fc093
		//fc086   fc094
		//fc087   fc095
		//fc088   fc096
		//fc089   fc097
		//fc090   fc098
		//fc091   fc099
		//1.对fc152的录入内容进行校验，若fc152录入内容为“床位费”，其对应的fc153的录入内容与上面分组合计结果数据不一致时，
		//出导出校验：发票XXX（XXX为发票号）床位费明细金额与床位费大项金额不一致，请修改。
		//2.上面分组合计结果数据与下面分组中左边字段录入内容为"床位费"的字段所对应的右边字段的录入内容不一致时，
		//出导出校验：发票XXX（XXX为发票号）床位费明细金额与床位费大项金额不一致，请修改。
		fieldsArr := [][]string{
			{"fc009", "fc010"},
			{"fc011", "fc012"},
			{"fc013", "fc014"},
			{"fc015", "fc016"},
			{"fc017", "fc018"},
			{"fc019", "fc020"},
			{"fc021", "fc022"},
			{"fc023", "fc024"},
			{"fc025", "fc026"},
			{"fc027", "fc028"},
			{"fc029", "fc030"},
			{"fc031", "fc032"},
			{"fc033", "fc034"},
			{"fc035", "fc036"},
			{"fc037", "fc038"},
			{"fc039", "fc040"},
			{"fc041", "fc042"},
			{"fc043", "fc044"},
			{"fc045", "fc046"},
			{"fc047", "fc048"},
		}
		// fc153isZero := false
		for _, qingDan := range qingDans {
			fc171 := ""
			// fc067 := ""
			isOK := false
			for _, field := range qingDan.Fields {
				if field.Code == "fc153" && field.FinalValue != "" {
					cache := field.FinalValue + "_" + GetFieldsFinal(qingDan.Fields, "fc152")
					if field.FinalValue == "0" {
						fmt.Println("---------------fc153------------------:", field.FinalValue)
						// fc153isZero = true
					}
					if arrays.Contains(fc153s, cache) == -1 {
						// fc153 = fc153 + ParseFloat(field.FinalValue)
						fc153 = SumFloat(fc153, ParseFloat(field.FinalValue), "+")
					}
					fc153s = append(fc153s, cache)
				}
				if field.Code == "fc152" && field.ResultValue == "床位费" {
					fc153 := GetFieldsInput(qingDan.Fields, "fc153")
					fc153All += ParseFloat(fc153)
				}
				for _, ccff := range ccffs {
					if field.Code == ccff[0] {
						c3 := GetFieldsInput(qingDan.Fields, ccff[2])
						if RegIsMatch(field.ResultValue, `(走廊|床位|人间|病房|陪护床)`) && !RegIsMatch(field.ResultValue, `(取暖|空调)`) {
							fmt.Println("---------------------------------:---------------------------------:", field.ResultValue, c3)
							isfc092 = true
							fc092All += ParseFloat(c3)
						}
					}
				}
				if RegIsMatch(field.Code, `^(fc084|fc154|fc092|fc162|fc172|fc085|fc155|fc093|fc163|fc173|fc086|fc156|fc094|fc164|fc174|fc087|fc157|fc095|fc165|fc175|fc088|fc158|fc096|fc166|fc176|fc089|fc159|fc097|fc167|fc177|fc090|fc160|fc098|fc168|fc178|fc091|fc161|fc099|fc169|fc179)$`) {
					// fmt.Println("---------------录入重复，请确认；------------------:", field.Code, field.ResultValue)
					if strings.Index(field.ResultValue, "?") != -1 || strings.Index(field.ResultValue, "？") != -1 {
						// fmt.Println("---------------清单内容不能录入问号，请检查------------------:", field.ResultValue)
						wrongNote += fapiao.Code + "账单号的清单内容不能录入问号，请检查；"
					}
					for _, ccff := range ccffs {
						if field.Code == ccff[1] {
							c1 := GetFieldsInputIndex(qingDan.Fields, ccff[0], field.BlockIndex)
							c2 := field.ResultValue
							c3 := GetFieldsInputIndex(qingDan.Fields, ccff[2], field.BlockIndex)
							c4 := GetFieldsInputIndex(qingDan.Fields, ccff[3], field.BlockIndex)
							c5 := GetFieldsInputIndex(qingDan.Fields, ccff[4], field.BlockIndex)
							// fmt.Println("---------------！！！！！------------------:", field.BlockIndex)
							// fmt.Println("---------------录入重复，请确认；------------------:", field.Code, c1, c2, c3, c4, c5)

							//if RegIsMatch(c1, `(走廊|床位|人间|病房|陪护床)`) && !RegIsMatch(c1, `(取暖|空调)`) {
							//	fc092All += ParseFloat(c3)
							//}

							if c1 != "" || c2 != "" || c3 != "" || c4 != "" || c5 != "" {
								ccvalue := c1 + "_" + c2 + "_" + c3 + "_" + c4 + "_" + c5
								// && field.ResultValue != ""
								if arrays.Contains(arrs, ccvalue) != -1 {
									// wrongNote += code + "账单号" + c1 + "录入重复，请确认；"
								}
								arrs = append(arrs, ccvalue)
							}
							if RegIsMatch(c2, `^(1)$`) {
								if !(c1 != "" && c2 != "" && c3 != "" && c4 != "" && c5 == "") {
									wrongNote += fapiao.Code + "账单号的清单组合内容录入有误，请检查；"
								}
							} else if RegIsMatch(c2, `^(2|3)$`) {
								if !(c1 != "" && c2 != "" && c3 != "" && c4 == "" && c5 == "") {
									wrongNote += fapiao.Code + "账单号的清单组合内容录入有误，请检查；"
								}
							} else if RegIsMatch(c2, `^(4)$`) {
								if !(c1 != "" && c2 != "" && c3 != "" && c4 == "" && c5 != "") {
									wrongNote += fapiao.Code + "账单号的清单组合内容录入有误，请检查；"
								}
							}

						}
						// if arrays.Contains(ccff, field.Code) != -1 {
						// 	_, OK := arrs[qq]
						// 	// fmt.Println("---------------arr------------------:", OK, qq, arrs[qq])
						// 	if OK {
						// 		if arrays.Contains(arrs[qq], field.ResultValue) != -1 && field.ResultValue != "" {
						// 			wrongNote += code + "账单号" + field.ResultValue + "录入重复，请确认；"
						// 		}
						// 		arrs[qq] = append(arrs[qq], field.ResultValue)
						// 		// fmt.Println("---------------arr22222------------------:", arrs[qq])
						// 	} else {
						// 		arrs[qq] = []string{field.ResultValue}
						// 	}
						// 	break
						// }
					}
				}
				if RegIsMatch(field.Code, `^(fc154|fc155|fc156|fc157|fc158|fc159|fc160|fc161)$`) {
					if RegIsMatch(field.ResultValue, `^(1|2|4)$`) {
						isOK = true
					}
					if RegIsMatch(field.ResultValue, `^(1|4)$`) {
						num14++
					} else if field.ResultValue != "" {
						names += field.Name + ","
					}
				}

				if RegIsMatch(field.Code, `^(fc154|fc155|fc156|fc157|fc158|fc159|fc160|fc161)$`) && field.ResultValue != "" {
					if !RegIsMatch(field.ResultValue, `^(1|2|3|4|5|6)$`) {
						wrongNote += code + "账单号的" + field.Name + "录入错误，请检查;"
					}
				}

				if RegIsMatch(field.Code, `^(fc162|fc163|fc164|fc165|fc166|fc167|fc168|fc169)$`) && (ParseFloat(field.ResultValue) > 1 || ParseFloat(field.FinalValue) > 1) {
					wrongNote += code + "账单号的" + field.Name + "最大值不能超过1，请检查"
				}

				if field.Code == "fc171" && field.FinalValue != "" {
					fc171 = field.FinalValue
				}

			}
			if isOK && fc171 != "" && fc067 != "" && fc171 != fc067 {
				chae := ParseFloat(fc067) - ParseFloat(fc171)
				wrongNote += code + "账单号的清单自费金额不等于结算单自费金额，差额" + ToString(chae) + "，请检查;"
			}
			// if fc067 == "0" {
			// 	wrongNote += code + "账单号的社保自费为0，请确认;"
			// }
			if num14 != 0 && names != "" {
				msg := code + "账单号的" + names + "清单类型录入不一致，请检查;"
				if strings.Index(wrongNote, msg) == -1 {
					wrongNote += code + "账单号的" + names + "清单类型录入不一致，请检查;"
				}
			}
		}

		hasChuangWeiFei := false
		for _, field := range fapiao.Fields {
			for _, fCodes := range fieldsArr {
				if field.Code == fCodes[0] && field.ResultValue == "床位费" {
					fc010 := GetFieldsInput(fapiao.Fields, fCodes[1])
					fc010All += ParseFloat(fc010)
					hasChuangWeiFei = true
				}
			}
		}
		global.GLog.Info("fc153All", zap.Any("", fc153All))
		global.GLog.Info("fc092All", zap.Any("", fc092All))
		global.GLog.Info("fc010All", zap.Any("", fc010All))
		fc056s := GetFieldsInputs(fapiao.Fields, "fc056")
		fc152Val := ""
		fc153Val := ""
		if fapiao.QingDan != nil && len(fapiao.QingDan) > 0 {
			fc152Val = GetFieldsInput(fapiao.QingDan[0].Fields, "fc152")
			fc153Val = GetFieldsInput(fapiao.QingDan[0].Fields, "fc153")
		}
		fff := []string{}
		for _, field := range fapiao.Fields {
			if RegIsMatch(field.Code, `^(fc009|fc011|fc013|fc015|fc017|fc019|fc021|fc023|fc025|fc027|fc029|fc031|fc033|fc035|fc037|fc039|fc041|fc043|fc045|fc047)$`) {
				fff = append(fff, field.FinalValue)
			}
		}
		for _, QingDan := range fapiao.QingDan {
			fc152 := getOneValue(QingDan.Fields, "fc152")
			fc171 := getOneValue(QingDan.Fields, "fc171")
			if arrays.Contains(fff, fc152) == -1 && fc171 != "" && fc171 != "0" && fc171 != "0.00" {
				wrongNote += "账单号" + code + "清单大项" + getOneResultValue(QingDan.Fields, "fc152") + "不在发票明细中，请核实并修改。"
			}
		}
		global.GLog.Info("fc056s", zap.Any("", fc056s))
		global.GLog.Info("fc152Val", zap.Any("", fc152Val))
		global.GLog.Info("fc153Val", zap.Any("", fc153Val))
		//2022年09月15日09:12:54 不存在bc007
		// if isfc092 && fc153All != fc092All && !utils2.HasItem(fc056s, "8") {
		// 	wrongNote += "发票" + code + "床位费明细金额与床位费大项金额不一致，请修改。"
		// }
		if isfc092 && fc010All != fc092All && utils2.HasItem(fc056s, "8") && hasChuangWeiFei {
			wrongNote += "发票" + code + "床位费明细金额与床位费大项金额不一致，请修改。"
		}
		chae := 0.0
		// if len(fc153s) > 0 && !fc153isZero {
		// 	// chae = fc008 - fc153
		// 	chae = SumFloat(fc008, fc153, "-")
		// }
		// if fc153 == 0 {
		// chae = fc008 - sum
		a1 := decimal.NewFromFloat(fc008)
		chae = ParseFloat(fmt.Sprint(a1.Sub(sum)))
		// }
		fmt.Println("--------------chaechaechae------------------:", fc008, fc153, chae)
		msg := "小于"
		if chae < 0 {
			msg = "大于"
		}
		if chae != 0.0 {
			wrongNote += "账单号" + code + "的发票明细金额" + msg + "总金额，差额" + ToString(chae) + ";"
		}
	}

	for key, value := range fc002Maps {
		if len(value) > 1 {
			fc059s := ""
			for _, fc059 := range value {
				fc059s += fc059 + "、"
			}
			wrongNote += "属性" + fc059s + "账单号[" + key + "]重复，请检查;"
		}
	}

	for _, fc062 := range fc062s {
		if fc062 != fc149 {

		}
	}
	if (fc149 == "null" && len(fc062s) > 1) || (fc149 != "null" && utils2.ItemHasNotInArr(fc062s, fc149)) {
		wrongNote = wrongNote + "发票姓名录入不一致,请检查；"
	}
	return wrongNote
}

func SumFloat(a1, a2 float64, ff string) float64 {
	b1 := decimal.NewFromFloat(a1)
	b2 := decimal.NewFromFloat(a2)
	if ff == "+" {
		return ParseFloat(fmt.Sprint(b1.Add(b2)))
	}
	if ff == "-" {
		return ParseFloat(fmt.Sprint(b1.Sub(b2)))
	}

	return 0.0
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""
	billNumbers := RegMatchAll(xmlValue, `<billNumber>.*?<\/billNumber>`)
	// menzhen := map[string]string{}
	// zhuyuan := map[string]string{}
	cValue := []string{}
	for _, billNumber := range billNumbers {
		bValue := GetNodeValue(billNumber, "billNumber")
		if bValue != "" {
			if arrays.ContainsString(cValue, bValue) != -1 {
				wrongNote += "账单号" + bValue + "重复，请检查；"
			}
			cValue = append(cValue, bValue)
		}
	}

	items := RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	for _, item := range items {
		billNumber := GetNodeValue(item, "billNumber")
		billStartDate := GetNodeValue(item, "billStartDate")
		billEndDate := GetNodeValue(item, "billEndDate")
		if billEndDate < billStartDate {
			wrongNote += billNumber + "账单号的出院日期早于入院日期，请修改;"
		}
		feeLists := RegMatchAll(item, `<feeList>[\s\S]*?<\/feeList>`)
		chongfu := []string{}
		isZore := true
		for _, feeList := range feeLists {
			feeTypeCode := GetNodeValue(feeList, "feeTypeCode")
			initialFee := GetNodeValue(feeList, "initialFee")
			selfFee := GetNodeValue(feeList, "selfFee")
			yyy := feeTypeCode
			if strings.Index(item, "<opList>") != -1 && HasKey(constMap["menZhenFeiYongLeiXingDaiMaMap"], feeTypeCode) {
				yyy = constMap["menZhenFeiYongLeiXingDaiMaMap"][feeTypeCode]
			} else if strings.Index(item, "<hospitalList>") != -1 && HasKey(constMap["zhuYuanFeiYongDaiMaMap"], feeTypeCode) {
				yyy = constMap["zhuYuanFeiYongDaiMaMap"][feeTypeCode]
			}
			if isZore && selfFee != "" && initialFee != "" && strings.Index(initialFee, "-") == -1 && ParseFloat(selfFee) > ParseFloat(initialFee) {
				wrongNote += billNumber + "账单号的" + yyy + "发票大项的自费金额大于总金额，请修改;"
			}
			// fmt.Println("---------------feeTypeCode------------------:", feeTypeCode)
			if !RegIsMatch(feeTypeCode, `^\d+$`) {
				wrongNote += "账单号[" + billNumber + "]的发票明细费用类型应该录入数据库中的内容，请检查;"
			}
			fmt.Println(chongfu)
			if arrays.ContainsString(chongfu, billNumber+"_"+feeTypeCode+"_"+initialFee) != -1 {
				wrongNote += billNumber + "账单号的发票大项内容存在重复," + yyy + "重复，请检查;"
			}
			chongfu = append(chongfu, billNumber+"_"+feeTypeCode+"_"+initialFee)
		}
	}

	deal_node := []string{"opList", "hospitalList"}
	for _, node_1_name := range deal_node {
		node_1_arr := GetNode(xmlValue, node_1_name)
		for _, node_1 := range node_1_arr {
			billNumber_val := GetNodeValue(node_1, "billNumber")
			feeList_node_arr := GetNode(node_1, "feeList")
			for _, feeList := range feeList_node_arr {
				initialFee_val := GetNodeDecimalValue(feeList, "initialFee")
				feeTypeCode_val := GetNodeValue(feeList, "feeTypeCode")
				selfFee_val := GetNodeDecimalValue(feeList, "selfFee")
				mess := billNumber_val + "账单号的" + feeTypeCode_val + "大项金额为负数，请修改；"
				messSelfFee := billNumber_val + "账单号的" + feeTypeCode_val + "大项的自费金额为负数，请修改；"
				fmt.Println("initialFee_val", initialFee_val)
				fmt.Println("selfFee_val", selfFee_val)
				if initialFee_val.LessThan(decimal.Zero) && strings.Index(wrongNote, mess) == -1 {
					wrongNote += mess
				}
				if selfFee_val.LessThan(decimal.Zero) && strings.Index(wrongNote, messSelfFee) == -1 {
					wrongNote += messSelfFee
				}
			}
		}
	}

	siList_node_arr := GetNode(xmlValue, "siList")
	for _, siList := range siList_node_arr {
		inSIAmount_val := GetNodeDecimalValue(siList, "inSIAmount")
		obillnumber_val := GetNodeValue(siList, "obillnumber")
		fee := GetNodeDecimalValue(siList, "fee")
		mess := obillnumber_val + "账单号的范围内金额为负数，请检查该发票的自费金额是否大于总金额；"
		messLessThanFee := "发票" + obillnumber_val + "范围内金额比报销金额小，请检查扣费内容是否有误；"
		fmt.Println("inSIAmount_val", inSIAmount_val)
		fmt.Println("fee", fee)
		if inSIAmount_val.LessThan(decimal.Zero) && strings.Index(wrongNote, mess) == -1 {
			wrongNote += mess
		}
		if inSIAmount_val.LessThan(fee) && strings.Index(wrongNote, messLessThanFee) == -1 {
			wrongNote += messLessThanFee
		}
	}

	hospitalLists := RegMatchAll(xmlValue, `<(hospitalList)>[\s\S]*?<\/(hospitalList)>`)
	//调整 CSB0118RC0279001 （当errorDesc节点值包含xxx执行单病种录入规则，且xxx与hospitalList节点下billNumber一致，则不执行该校验）
	errorDescList := RegMatchAll(xmlValue, `<(errorDesc)>[\s\S]*?<\/(errorDesc)>`)
	isErrorDesc := true
	for _, hospitalList := range hospitalLists {
		billNumber := GetNodeValue(hospitalList, "billNumber")
		for _, errorDesc := range errorDescList {
			errorDescVal := GetNodeValue(errorDesc, "errorDesc")
			if strings.Index(errorDescVal, billNumber) != -1 {
				if strings.Index(errorDescVal, "执行单病种录入规则") != -1 {
					isErrorDesc = false
				}
			}
		}
		if strings.Index(hospitalList, "<feeTypeCode>204</feeTypeCode>") == -1 && isErrorDesc {
			wrongNote += billNumber + "账单号的住院发票没有床位费大项，请检查；"
		} else {
			feeLists := RegMatchAll(xmlValue, `<feeList>[\s\S]*?<\/feeList>`)
			for _, feeList := range feeLists {
				feeTypeCode := GetNodeValue(feeList, "feeTypeCode")
				if feeTypeCode == "204" && GetNodeValue(feeList, "initialFee") == "0" {
					wrongNote += billNumber + "账单号的住院发票床位费大项金额为0，请检查；"
					break
				}
			}
		}
	}

	//CSB0118RC0301000
	//当xml文件中不存在opList、hospitalList大节点时，出导出校验：不存在发票情况，请确认初审是否录入有误，（异常件不执行该检验）
	if obj.Bill.Status != 3 {
		arr := GetNode(xmlValue, "opList")
		arr1 := GetNode(xmlValue, "hospitalList")
		if len(arr)+len(arr1) < 1 {
			wrongNote += "不存在发票情况，请确认初审是否录入有误；"
		}
	}

	//CSB0118RC0348000
	//fc277录入值为1时
	//1、fc279 节点“billCode”的值为空时，出导出校验“票据代码不能为空”。
	//2、fc281 节点“billCheckCode”的值为空时，出导出校验“校验码不能为空”。
	//CSB0118RC0343000 05/06 调整新添加
	//当fc278录入值为3且fc279和fc281录入值都为空时，不出该校验
	for _, opItem := range obj.OpList {
		fcMap := make(map[string]string)
		fcList := []string{"fc277", "fc278", "fc279", "fc281"}
		for _, code := range fcList {
			fcMap[code] = ""
		}
		currentFc002ResultValue := ""
		for _, opField := range opItem.Fields {
			_, ok := fcMap[opField.Code]
			if ok {
				fcMap[opField.Code] = opField.ResultValue
			}
			if opField.Code == "fc002" {
				currentFc002ResultValue = opField.ResultValue
			}
		}
		if !(fcMap["fc278"] == "3" && fcMap["fc279"] == "" && fcMap["fc281"] == "") {
			for _, opField := range opItem.Fields {
				if opField.Code == "fc277" && opField.ResultValue == "1" {
					opList := utils2.GetNodeData(xmlValue, "opList")
					for _, opItem := range opList {
						billNumberItem := utils2.GetNodeData(opItem, `billNumber`)
						billCodeItem := utils2.GetNodeData(opItem, `billCode`)
						billCheckCodeItem := utils2.GetNodeData(opItem, `billCheckCode`)
						billNumberItemStr := ""
						if len(billNumberItem) == 1 && billNumberItem[0] == currentFc002ResultValue {
							billNumberItemStr = billNumberItem[0]
							if len(billCodeItem) == 1 && billCodeItem[0] == "" {
								wrongNote += "【" + billNumberItemStr + "】 票据代码不能为空；"
							}
							if len(billCheckCodeItem) == 1 && billCheckCodeItem[0] == "" {
								wrongNote += "【" + billNumberItemStr + "】 校验码不能为空；"
							}
						}

					}
				}
			}
		}

	}
	//for _, fieldsMap := range obj.Fields {
	//fmt.Println(fmt.Sprintf("----==%v==--", fieldsMap))
	//if fieldsMap.Code == "fc277" && fieldsMap.ResultValue == "1" {
	//	opList := utils2.GetNodeData(xmlValue, "opList")
	//	for _, opItem := range opList {
	//		billNumberItem := utils2.GetNodeData(opItem, `billNumber`)
	//		billCodeItem := utils2.GetNodeData(opItem, `billCode`)
	//		billCheckCodeItem := utils2.GetNodeData(opItem, `billCheckCode`)
	//		billNumberItemStr := ""
	//		if len(billNumberItem) == 1 {
	//			billNumberItemStr = billNumberItem[0]
	//			if len(billCodeItem) == 1 && billCodeItem[0] == "" {
	//				wrongNote += "【" + billNumberItemStr + "】 票据代码不能为空；"
	//			}
	//			if len(billCheckCodeItem) == 1 && billCheckCodeItem[0] == "" {
	//				wrongNote += "【" + billNumberItemStr + "】 校验码不能为空；"
	//			}
	//		}
	//
	//	}
	//}
	//}
	//}

	return wrongNote
}

func GetNode(xmlValue, node string) []string {
	query := `<` + node + `>[\s\S]*?<\/` + node + `>`
	return RegMatchAll(xmlValue, query)
}

func GetFieldsInput(fields []model.ProjectField, code string) string {
	for _, field := range fields {
		if field.Code == code {
			// fmt.Println("---------------！！！！！------------------:", field.Code, field.BlockIndex)
			return field.ResultValue
		}
	}
	return ""
}

func GetFieldsInputExits(fields []model.ProjectField, code string) (bool, string) {
	for _, field := range fields {
		if field.Code == code {
			// fmt.Println("---------------！！！！！------------------:", field.Code, field.BlockIndex)
			return true, field.ResultValue
		}
	}
	return false, ""
}

func GetFieldsFinal(fields []model.ProjectField, code string) string {
	for _, field := range fields {
		if field.Code == code {
			return field.FinalValue
		}
	}
	return ""
}

func GetFieldsInputs(fields []model.ProjectField, code string) []string {
	var result []string
	for _, field := range fields {
		if field.Code == code {
			result = append(result, field.ResultValue)
		}
	}
	return result
}

func GetFieldsInputIndex(fields []model.ProjectField, code string, index int) string {
	for _, field := range fields {
		if field.Code == code && field.BlockIndex == index {
			// fmt.Println("---------------！！！！！------------------:", field.Code, field.BlockIndex)
			return field.ResultValue
		}
	}
	return ""
}
