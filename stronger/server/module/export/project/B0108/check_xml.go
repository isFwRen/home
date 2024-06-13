package B0108

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"server/module/export/service"
	"server/module/load/model"
	proModel "server/module/pro_conf/model"
	_ "server/utils"
	utils2 "server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
	_ "github.com/wxnacy/wgo/arrays"
	_ "go.uber.org/zap"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	//global.GLog.Info("B0108:::CheckXml")
	obj := o.(FormatObj)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	// ----------------------------------code---------------------------------------------
	constMap := constDeal(obj.Bill.ProCode)
	wrongNote += CodeCheck(obj, xmlValue, constMap)
	wrongNote += XmlCheck(obj, xmlValue, constMap)

	// ----------------------------------xml---------------------------------------------
	fmt.Println(wrongNote)
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
	SaleChannel := obj.Bill.SaleChannel
	hyField := []string{}
	fc107s := []string{}
	fc041All := 0.0
	fc064All := 0.0
	//fc033Value := ""
	//fc059Value := ""
	isNull := true
	fc019 := ""
	fc019_is := false
	fc023 := ""
	fc021 := ""

	fc084Arr := []string{}
	fc035 := ""

	for _, field := range fields {

		//---------------
		if field.Code == "fc035" && field.FinalValue != "" && strings.Index(wrongNote, "就诊日期年份存在异常，请检查；") == -1 {
			if fc035 == "" {
				fc035 = Substr(field.FinalValue, 0, 4)
			} else if Substr(field.FinalValue, 0, 4) != fc035 {
				wrongNote += "就诊日期年份存在异常，请检查；"
			}
		}

		if field.Code == "fc084" {
			fc084Arr = append(fc084Arr, field.ResultValue)
		}

		if field.Code == "fc019" {
			fc019 = field.ResultValue
			fc019_is = true
		}

		if field.Code == "fc023" {
			fc023 = field.ResultValue
			fmt.Println(fc023)
		}

		if field.Code == "fc021" {
			fc021 = field.ResultValue
			fmt.Println(fc021)
		}

		if strings.Index(wrongNote, "诊断结论为“手术后对症治疗”，请确认；") == -1 {
			if RegIsMatch(field.Code, "^(fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$") && strings.Index(field.ResultValue, "手术后对症治疗") != -1 {
				wrongNote += "诊断结论为“手术后对症治疗”，请确认；"
			}
		}

		//CSB0108RC0098000
		//global.GLog.Info("CSB0108RC0098000")
		//当fc018、fc112、fc036、fc104、fc058、fc161、fc162、fc205、fc206、fc207、fc208、fc209
		//以上字段任一结果值没有转码，则出一条导出校验：疾病诊断录入有误，请修改；
		if strings.Index(wrongNote, "诊断结论录入内容不在代码表中，请修改；") == -1 {
			if RegIsMatch(field.Code, "^(fc018|fc112|fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209)$") && field.ResultValue != "" {
				if !HasKey(constMap["zhenDuanDaiMaMap"], field.ResultValue) {
					wrongNote += "诊断结论录入内容不在代码表中，请修改；"
				}
			}
		}
		//CSB0108RC0093000
		//fc034、fc055录入值不在《B0108_太平理赔_医院代码表》常量库时，出导出校验：就诊医院名称录入错误；
		//global.GLog.Info("CSB0108RC0093000")
		// if strings.Index(wrongNote, "就诊医院名称录入错误；") == -1 {
		// 	if RegIsMatch(field.Code, `^(fc034|fc055|fc226|fc227)$`) && field.ResultValue != "" {
		// 		if !HasKey(constMap["yiYuanDaiMaMap"], field.ResultValue) {
		// 			wrongNote = wrongNote + "就诊医院名称录入错误；"
		// 		}
		// 	}
		// }

		//CSB0108RC0126000
		//global.GLog.Info("CSB0108RC0126000")
		//"fc113、fc114、fc115、fc116、fc117、fc118、fc119、fc120、fc121、fc122、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc131、fc132以上字段的录入值必须为数据库中内容，否则出导出校验：发票名称大项[XXX]不为常量库内容，请修改；
		//（xxx为不为常量库内容的名称，如fc113录入值为氧气费，则导出校验为：发票名称大项[氧气费]不为常量库内容，请修改）"

		//CSB0108RC0076000
		//当fc003的结果值与fc020的结果值不一致时，出导出校验：实际门诊发票张数与申请表中【门诊发票张数】不一致，请核查（当不存在fc020或fc003字段时，不应进行该校验）
		////global.GLog.Info("CSB0108RC0076000")
		if strings.Index(wrongNote, "实际门诊发票张数与申请表中【门诊发票张数】不一致，请核查；") == -1 {
			if IsExist(fields, "fc003") && IsExist(fields, "fc020") {
				if GetFieldsFinal(fields, "fc003") != GetFieldsFinal(fields, "fc020") {
					wrongNote = wrongNote + "实际门诊发票张数与申请表中【门诊发票张数】不一致，请核查；"
				}
			}
		}

		//CSB0108RC0077000
		//global.GLog.Info("CSB0108RC0077000")
		//当fc004的结果值与fc022的结果值不一致时，出导出校验：实际住院发票张数与申请表中【住院发票张数】不一致，请核查（当不存在fc022或fc004字段时，不应进行该校验）
		////global.GLog.Info("CSB0108RC0077000")
		isExistFc004 := false
		isExistFc022 := false
		for _, field1 := range fields {
			if field1.Code == "fc004" {
				isExistFc004 = true
			}
			if field1.Code == "fc022" {
				isExistFc022 = true
			}
		}
		if isExistFc004 && isExistFc022 {
			if field.Code == "fc004" && GetFieldsFinal(fields, "fc004") != GetFieldsFinal(fields, "fc022") {
				wrongNote = wrongNote + "实际住院发票张数与申请表中【住院发票张数】不一致，请核查；"
			}
		}

		//CSB0108RC0091000
		//当fc034、fc055、fc226、fc227录入内容包含“毫州”时，出导出校验：医院名称亳（bo）州、毫（hao）州，请确认；
		//global.GLog.Info("CSB0108RC0091000")
		if RegIsMatch(field.Code, "fc034|fc055|fc226|fc227") {
			if strings.Index(wrongNote, "医院名称亳（bo）州、毫（hao）州，请确认；") == -1 {
				if strings.Index(field.ResultValue, "毫州") != -1 {
					wrongNote += "医院名称亳（bo）州、毫（hao）州，请确认；"
				}
			}
		}

		if RegIsMatch(field.Code, `^(fc107)$`) && field.FinalValue != "" && arrays.ContainsString(fc107s, field.FinalValue) == -1 {
			fc107s = append(fc107s, field.FinalValue)
		}
		//CSB0108RC0096000
		//当fc160结果值为0时，出导出校验：诊断结论内容为空，请确认修改；
		//global.GLog.Info("CSB0108RC0096000")
		if field.Code == "fc160" && field.ResultValue == "0" {
			wrongNote += "诊断结论内容为空，请确认修改；"
		}

		//CSB0108RC0099000
		//当fc036、fc104、fc058、fc161、fc162、fc205、fc206、fc207、fc208、fc209以上字段结果值均为空值时，出导出校验：诊断内容为空；
		////global.GLog.Info("CSB0108RC0099000")

		hyField = []string{"fc036", "fc104", "fc058", "fc161", "fc205", "fc206", "fc207", "fc208", "fc209", "fc239", "fc240", "fc241", "fc242", "fc243", "fc244", "fc245", "fc246", "fc247", "fc248"}
		for i := range hyField {
			if field.Code == hyField[i] && field.ResultValue != "" {
				isNull = false
			}
		}

		//CSB0108RC0107000
		////global.GLog.Info("CSB0108RC0107000")
		//将fc107与fc024结果值作对比（当存在多个fc107时，需将每个fc107的结果值相加求和再与fc024结果值作对比）， 不一致时出导出校验：
		//报销单统筹支付金额与申请书社保报销金额录入不一致，请核查（当不存在fc024或fc107字段时，不应进行该校验）

		//if strings.Index(wrongNote, "报销单统筹支付金额与申请书社保报销金额录入不一致，请核查；") == -1 {
		//	if IsExist(fields, "fc107") && IsExist(fields, "fc024") {
		//		difference := GetFieldsAllFinal(fields, "fc107") - ParseFloat(GetFieldsFinal(fields, "fc024"))
		//		if strconv.FormatFloat(difference, 'f', 2, 64) != "0.00" && strconv.FormatFloat(difference, 'f', 2, 64) != "-0.00" {
		//			wrongNote += "报销单统筹支付金额与申请书社保报销金额录入不一致，请核查；"
		//		}
		//	}
		//}

		//CSB0108RC0109000
		////global.GLog.Info("CSB0108RC0109000")
		//将fc108与fc025结果值作对比（当存在多个fc108时，需将每个fc108的结果值相加求和再与fc025结果值作对比），不一致出导出校验：
		//报销单其他扣除金额与申请书其他报销金额录入不一致，请核查（当不存在fc025或fc108字段时，不应进行该校验）
		if strings.Index(wrongNote, "报销单其他扣除金额与申请书其他报销金额录入不一致，请核查；") == -1 {
			if IsExist(fields, "fc108") && IsExist(fields, "fc025") {
				if GetFieldsAllFinal(fields, "fc108")-ParseFloat(GetFieldsFinal(fields, "fc025")) != 0 {
					wrongNote += "报销单其他扣除金额与申请书其他报销金额录入不一致，请核查；"
				}
			}
		}
		//CSB0108RC0101000
		////global.GLog.Info("CSB0108RC0101000")
		//当fc204录入值为0时，出导出校验：发票xxx有收取手术费，但手术数量为0，请核查；（xxx为同一发票属性下的fc032或fc054的结果值）
		////global.GLog.Info("CSB0108RC0101000")
		//if obj.Bill.Agency[:3] != "001" {
		//	if field.Code == "fc204" && field.ResultValue == "0" {
		//		wrongNote += "发票" + GetFieldsFinal(fields, "fc032") + "有收取手术费，但手术数量为0，请核查；"
		//	}
		//}

		//CSB0108RC0102000
		////global.GLog.Info("CSB0108RC0102000")
		//当fc165、fc196、fc197、fc198、fc199、fc204以上字段录入值包含?或？时，出一条导出校验：手术填写有误(包含?或？)，请核查；
		////global.GLog.Info("CSB0108RC0102000")
		if strings.Index(wrongNote, "手术填写有误(包含?或？)，请核查；") == -1 {
			if RegIsMatch(field.Code, "^(fc165|fc196|fc197|fc198|fc199|fc204)$") {
				if RegIsMatch(field.ResultValue, "\\?|？") {
					wrongNote += "手术填写有误(包含?或？)，请核查；"
				}
			}
		}

		//CSB0108RC0104000
		//"当左边字段有录入值时，校验对应右边字段是否有录入值，若无则出一条导出校验：手术填写有误(无录入值)，请核查；
		//fc164 fc165
		//fc200 fc196
		//fc201 fc197
		//fc202 fc198
		//fc203 fc199"
		//global.GLog.Info("CSB0108RC0104000")
		hyFieldArr := [][]string{
			{"fc164", "fc165"},
			{"fc200", "fc196"},
			{"fc201", "fc197"},
			{"fc202", "fc198"},
			{"fc203", "fc199"},
		}
		if strings.Index(wrongNote, "手术填写有误(无录入值)，请核查；") == -1 {
			for i := range hyFieldArr {
				if field.Code == hyFieldArr[i][0] && field.ResultValue != "" {
					if GetFieldsInput(fields, hyFieldArr[i][1]) == "" {
						wrongNote += "手术填写有误(无录入值)，请核查；"
					}
				}
			}
		}
		//CSB0108RC0120000
		//当fc176结果值为0时，出导出校验：科别不能为空，请修改；
		//global.GLog.Info("CSB0108RC0120000")
		//20230711 取消需求
		//if field.Code == "fc176" && field.FinalValue == "0" {
		//	wrongNote += "科别不能为空，请修改；"
		//}

		//CSB0108RC0086000
		//当fc036、fc104、fc058、fc161、fc162、fc205、fc206、fc207、fc208、fc209以上字段结果值为“1403A28.1”时，
		//对应字段出导出校验，请确认XXX录入内容是否有误，确认后修改；（xxx为对应字段名）
		//global.GLog.Info("CSB0108RC0086000")

		if RegIsMatch(field.Code, "^(fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$") && field.FinalValue == "1403A28.1" {
			if strings.Index(wrongNote, "请确认"+field.Code+"录入内容是否有误，确认后修改；") == -1 {
				wrongNote += "请确认" + field.Code + "录入内容是否有误，确认后修改；"
			}
		}

		//CSB0108RC0083000
		//将fc175结果值与当前系统日期作对比，当fc175的值小于当前系统日期时，出导出校验：
		//身份证过期；（如fc175的结果值为20210909，当前系统日期为20211019，则需出导出校验）
		//global.GLog.Info("CSB0108RC0083000")
		currentTime, _ := strconv.Atoi(time.Now().Format("20060102")) //获取当前时间 20220927
		if field.Code == "fc175" {
			fc175ValueInt, _ := strconv.Atoi(field.FinalValue)
			if fc175ValueInt < currentTime {
				wrongNote += "身份证截止日期不能小于当前系统日期，请检查；"
			}
		}

		//CSB0108RC0084000
		//当fc012录入内容为空值、0、6、7、8、10、11、12、13、14、15、19时，出导出校验：申请人类型录入内容不在选项中，请核查；
		//global.GLog.Info("CSB0108RC0084000")
		if strings.Index(wrongNote, "申请人类型录入内容不在选项中，请核查；") == -1 {
			if RegIsMatch(GetFieldsInput(fields, "fc012"), "^(0|6|7|8|10|11|12|13|14|15|19)$") || (GetFieldsExit(fields, "fc012") && GetFieldsInput(fields, "fc012") == "") {
				wrongNote += "申请人类型录入内容不在选项中，请核查；"
			}
		}

		//CSB0108RC0085000
		//当fc012录入值为A时，出导出校验：
		//事故者与申请人姓名一致时申请人类型选择“被保险人”，不一致时事故者为未成年人则选择“监护人”，成年则选择“其他”；
		//global.GLog.Info("CSB0108RC0085000")
		// 20230711 取消需求
		//if field.Code == "fc012" && field.ResultValue == "A" {
		//	wrongNote += "事故者与申请人姓名一致时申请人类型选择“被保险人”，不一致时事故者为未成年人则选择“监护人”，成年则选择“其他”；"
		//}

		//CSB0108RC0090000
		//当存在多个fc033或fc059字段，校验每个字段的结果值，不一致时出导出校验：病人姓名（发票）不一致，请确认后修改；（为空不校验）
		//global.GLog.Info("CSB0108RC0090000")
		//20230711 取消需求
		//if strings.Index(wrongNote, "病人姓名（发票）不一致，请确认后修改；") == -1 {
		//	if field.Code == "fc033" || field.Code == "fc059" {
		//		fc033Value = field.FinalValue
		//		fc059Value = field.FinalValue
		//		for _, index := range fields {
		//			if strings.Index(wrongNote, "病人姓名（发票）不一致，请确认后修改；") == -1 {
		//				if index.Code == "fc033" || index.Code == "fc059" {
		//					if fc033Value != index.FinalValue || fc059Value != index.FinalValue {
		//						wrongNote += "病人姓名（发票）不一致，请确认后修改；"
		//					}
		//				}
		//			}
		//		}
		//	}
		//}
		//CSB0108RC0080000
		//"将每个fc040、fc063的结果值相加之和与fc024结果值作对比，不一致出导出校验：发票统筹支付金额与申请书社保报销金额录入不一致，请核查（当不存在fc024字段时，不应进行该校验）
		//注：一个案件中可能存在一张或多张发票，即fc040、fc063可能存在一个或存在多个
		//如：
		//只有门诊发票则只会存在一个或多个fc040不存在fc063，需将每个fc040结果值相加之和与fc024结果值作对比；
		//只有住院发票则只会存在一个或多个fc063不存在fc040，需将每个fc063结果值相加之和与fc024结果值作对比；
		//门诊发票和住院发票同时存在时则会出现一个fc040和fc063或多个fc040多个fc063，需将每个fc040和fc063的结果值相加之和与fc024结果值作对比；"
		//global.GLog.Info("CSB0108RC0080000")
		if obj.Bill.Agency[:3] != "001" {
			if strings.Index(wrongNote, "发票统筹支付金额与申请书社保报销金额录入不一致，请核查；") == -1 {
				// fmt.Println("-----------发票统筹支付金额与申请书社保报销金额录入不一致，请核查；---------------")
				if IsExist(fields, "fc024") {
					// fmt.Println("-----------fc024fc024---------------")
					difference := (GetFieldsAllFinal(fields, "fc040") + GetFieldsAllFinal(fields, "fc063")) - GetFieldsAllFinal(fields, "fc024")
					if strconv.FormatFloat(difference, 'f', 2, 64) != "0.00" && strconv.FormatFloat(difference, 'f', 2, 64) != "-0.00" {
						wrongNote += "发票统筹支付金额与申请书社保报销金额录入不一致，请核查；"
					}
				}
			}
		}

		//CSB0108RC0081000
		//"将每个fc041、fc064的结果值相加之和与fc025结果值作对比，不一致出导出校验：发票其他扣除金额与申请书其他报销金额录入不一致，请核查（当不存在fc025字段时，不应进行该校验）
		//注：一个案件中可能存在一张或多张发票，即fc041、fc064可能存在一个或存在多个
		//如：
		//只有门诊发票则只会存在一个或多个fc041不存在fc064，需将每个fc040结果值相加之和与fc025结果值作对比；
		//只有住院发票则只会存在一个或多个fc064不存在fc041，需将每个fc063结果值相加之和与fc025结果值作对比；
		//门诊发票和住院发票同时存在时则会出现一个fc041和fc064或多个fc041多个fc064，需将每个fc041和fc064的结果值相加之和与fc025结果值作对比；"
		//global.GLog.Info("CSB0108RC0081000")
		if strings.Index(wrongNote, "发票其他扣除金额与申请书其他报销金额录入不一致，请核查；") == -1 {
			if IsExist(fields, "fc025") {
				if (GetFieldsAllFinal(fields, "fc041")+GetFieldsAllFinal(fields, "fc064"))-GetFieldsAllFinal(fields, "fc025") != 0 {
					wrongNote += "发票其他扣除金额与申请书其他报销金额录入不一致，请核查；"
				}
			}
		}

		//CSB0108RC0092000
		//fc034、fc055录入值为纯数字时，出导出校验：就诊医院名称录入错误；
		//global.GLog.Info("CSB0108RC0092000")
		if field.Code == "fc034" || field.Code == "fc055" {
			if RegIsMatch(field.ResultValue, "^[0-9]*$") {
				wrongNote += "就诊医院名称录入错误；"
			}
		}

		//CSB0108RC0097000
		//当fc160结果值不为数字格式时，出导出校验：fc160诊断结论数量录入内容有误，请确认后修改；
		//global.GLog.Info("CSB0108RC0097000")
		if field.Code == "fc160" && !RegIsMatch(field.FinalValue, "^[0-9]*$") {
			wrongNote += "fc160诊断结论数量录入内容有误，请确认后修改；"
		}

		//CSB0108RC0087000
		//当fc036、fc104、fc058、fc161、fc162、fc205、fc206、fc207、fc208、fc209
		//以上字段结果值为2402Z24.201、2402Z24.2、1413B33.803、1409A82.101、2402Z20.3、1409A82.9、1409A82.001、1409A82.901、2402Z20.301
		//的其中一项时，出导出校验，请确认XXX字段录入内容是否有误；（xxx为对应字段名）
		//global.GLog.Info("CSB0108RC0087000")
		fieldMatchResult, _ := regexp.MatchString("^(fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$", field.Code)
		if fieldMatchResult {
			matchFinalValue, _ := regexp.MatchString("^(2402Z24.201|2402Z24.2|1413B33.803|1409A82.101|2402Z20.3|1409A82.9|1409A82.001|1409A82.901|2402Z20.301)$", field.FinalValue)
			if matchFinalValue {
				if strings.Index(wrongNote, "请确认"+field.Code+"字段录入内容是否有误；") == -1 {
					wrongNote += "请确认" + field.Code + "字段录入内容是否有误；"
				}
			}
		}

		//CSB0108RC0100000
		//当fc113、fc114、fc115、fc116、fc117、fc118、fc119、fc120、fc121、fc122、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc131、fc132
		//以上字段的录入值其中之一包含“手术费”，但该案件无fc204字段时，出导出校验：有收取手术费，但无手术记录；
		//global.GLog.Info("CSB0108RC0100000")
		if RegIsMatch(field.Code, "^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$") {
			if strings.Index(wrongNote, "有收取手术费，但无手术记录；") == -1 {
				if field.ResultValue == "手术费" && !IsExist(fields, "fc204") {
					wrongNote += "有收取手术费，但无手术记录；"
				}
			}
		}

		//CSB0108RC0130000
		//当存在字段出问题件时，出导出校验：XXX为问题件，请确认；（xxx为字段名）
		//global.GLog.Info("CSB0108RC0130000")
		if field.Issues != nil {
			if strings.Index(wrongNote, field.Code+"为问题件，请确认；") == -1 {
				wrongNote += field.Code + "为问题件，请确认；"
			}
		}

		if field.Code == "fc041" {
			fc041All += ParseFloat(GetFieldsFinal(fields, "fc041"))
		}
		if field.Code == "fc064" {
			fc064All += ParseFloat(GetFieldsFinal(fields, "fc064"))
		}
		//编码 CSB0108RC0295000 当fc036、fc104、fc058、fc161、fc162、fc205、fc206、fc207、fc208、fc209、fc239、fc240、fc241、fc242、fc243、fc244、fc245、fc246、fc247、fc248以上字段录入值为B时，则出一条导出校验：诊断结论录入内容为B，请检查；
		if RegIsMatch(field.Code, "^(fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$") {
			if strings.Index(wrongNote, "诊断结论录入内容为B，请检查；") == -1 {
				if field.ResultValue == "B" && field.ResultValue != "" {
					wrongNote += "诊断结论录入内容为B，请检查；"
				}
			}
		}

		// CSB0108RC0311000
		// 清单循环分块中的字段：fc080、fc143、fc145、fc147、fc149、fc151、fc153、fc155录入值匹配到《B0108_太平理赔_住院发票大项代码表》中“大项名称”一列（第一列）时，出导出校验：【xxx】为大项，请检查；（xxx为以上字段的录入内容）
		if strings.Index(wrongNote, "【"+field.ResultValue+"】为大项，请检查；") == -1 {
			if RegIsMatch(field.Code, "^(fc080|fc143|fc145|fc147|fc149|fc151|fc153|fc155)$") && field.ResultValue != "" {
				if HasKey(constMap["faPiaoDaXiangMap3"], field.ResultValue) {
					wrongNote += "【" + field.ResultValue + "】为大项，请检查；"
				}
			}
		}
		//------------------------------------------------------------------------------------
		//CSB0108RC0326000 当fc084（模板类型字段，存在多个）录入内容包含99时，出导出校验：该案件存在多余报销单，请确认是否多合一案件；
		if field.Code == "fc084" {
			if strings.Index(field.ResultValue, "99") != -1 {
				if strings.Index(wrongNote, "该案件存在多余报销单，请确认是否多合一案件；") == -1 {
					wrongNote = wrongNote + "该案件存在多余报销单，请确认是否多合一案件；"
				}
			}
		}
	}
	//CSB0108RC0074000
	//global.GLog.Info("CSB0108RC0074000")
	//理赔案件，当fc084的值不包含5和6时，出导出校验，提示：事故发生日期请确认是否录入有误（注意死亡日期、意外发生日期、疾病首次就诊日期）
	isExist := true
	fc084_10 := false
	for i := range fc084Arr {
		if RegIsMatch(fc084Arr[i], "^(5|6)$") {
			isExist = false
		}
		if RegIsMatch(fc084Arr[i], "10") {
			fc084_10 = true
		}
	}

	if fc019_is && fc019 == "" {
		wrongNote += "申请理赔类别不能为空，请检查；"
	}

	//20230711 取消需求
	//if RegIsMatch(fc019, `15`) && fc023 == "" {
	//	wrongNote += "理赔类别与住院信息不符，请检查；"
	//}
	//20230711 取消需求
	//if RegIsMatch(fc019, `19`) && fc021 == "" {
	//	wrongNote += "理赔类别与门诊信息不符，请检查；"
	//}

	//20230711 取消需求
	//if RegIsMatch(fc019, `(14|12|13|17)`) && RegIsMatch(fc019, `11`) {
	//	wrongNote += "申请理赔类别不能同时存在身故和重疾/残疾/烧烫伤/特种疾病津贴，请检查;"
	//}

	if RegIsMatch(fc019, `(11|12|13|17)`) && !fc084_10 {
		wrongNote += "初审未切“10”，请确认是否为身故、重疾案件；"
	}
	if !RegIsMatch(fc019, `(11|12|13|17)`) && fc084_10 {
		wrongNote += "请确认是否为身故、重疾案件；"
	}

	ffMaps := [][]string{{"fc036", "fc228"}, {"fc104", "fc229"}, {"fc058", "fc230"}, {"fc161", "fc231"}, {"fc162", "fc232"}, {"fc205", "fc233"}, {"fc206", "fc234"}, {"fc207", "fc235"}, {"fc208", "fc236"}, {"fc209", "fc237"}, {"fc239", "fc249"},
		{"fc240", "fc250"},
		{"fc241", "fc251"},
		{"fc242", "fc252"},
		{"fc243", "fc253"},
		{"fc244", "fc254"},
		{"fc245", "fc255"},
		{"fc246", "fc256"},
		{"fc247", "fc257"},
		{"fc248", "fc258"}}
	for _, fmap := range ffMaps {
		f0 := GetFieldsInput(fields, fmap[0])
		f1 := GetFieldsInput(fields, fmap[1])
		if f0 == "B" && f1 == "" {
			wrongNote += "诊断结论为B，诊断详细名称为空，请修改;"
			break
		}
	}

	// fmt.Println("89898978667564648")
	// fmt.Println(fc084Arr)
	// fmt.Println(isExist)
	if SaleChannel == "理赔" && isExist {
		if strings.Index(wrongNote, "事故发生日期请确认是否录入有误（注意死亡日期、意外发生日期、疾病首次就诊日期）；") == -1 {
			wrongNote += "事故发生日期请确认是否录入有误（注意死亡日期、意外发生日期、疾病首次就诊日期）；"
		}
	}
	if strings.Index(wrongNote, "诊断内容为空") == -1 {
		if isNull {
			wrongNote += "疾病诊断不能为空;"
		}
	}

	//CSB0108RC0131000
	//当存在fc108时，将fc108与所有fc041、fc064结果值之和作对比（即所有fc041+fc064的值与fc108作对比），不一致出导出校验：
	//第三方报销金额与发票其他扣除金额录入不一致，请核查；
	//global.GLog.Info("CSB0108RC0131000")
	if strings.Index(wrongNote, "第三方报销金额与发票其他扣除金额录入不一致，请核查；") == -1 {
		if IsExist(fields, "fc108") {
			if ParseFloat(GetFieldsFinal(fields, "fc108"))-(GetFieldsAllFinal(fields, "fc041")+GetFieldsAllFinal(fields, "fc064")) != 0 {
				wrongNote += "第三方报销金额与发票其他扣除金额录入不一致，请核查；"
			}
		}
	}

	//CSB0108RC0079000
	//"1.将所有的fc039的结果值相加之和与fc021结果值作对比，不一致时出导出校验：
	//	实际门诊发票总金额与申请书门诊发票总金额不一致（当不存在fc039或fc021时，不执行该校验）
	//2.将所有的fc062的结果值相加之和与fc023结果值作对比，不一致时出导出校验：
	//	实际住院发票总金额与申请书住院发票总金额不一致（当不存在fc062或fc023时，不执行该校验）"
	//global.GLog.Info("CSB0108RC0079000")
	if obj.Bill.Agency[:3] != "001" {
		if strings.Index(wrongNote, "实际门诊发票总金额与申请书门诊发票总金额不一致") == -1 {
			if IsExist(fields, "fc039") && IsExist(fields, "fc021") {
				if GetFieldsAllFinal(fields, "fc039") != GetFieldsAllFinal(fields, "fc021") {
					wrongNote += "实际门诊发票总金额与申请书门诊发票总金额不一致；"
				}
			}
		}
		if strings.Index(wrongNote, "实际住院发票总金额与申请书住院发票总金额不一致；") == -1 {
			if IsExist(fields, "fc062") && IsExist(fields, "fc023") {
				if GetFieldsAllFinal(fields, "fc062") != GetFieldsAllFinal(fields, "fc023") {
					wrongNote += "实际住院发票总金额与申请书住院发票总金额不一致；"
				}
			}
		}
	}

	for _, fapiao := range obj.ClinicInfo {
		for _, field := range fapiao.Fields {
			if RegIsMatch(field.Code, `^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$`) && field.FinalValue != "" {
				if strings.Index(wrongNote, "发票名称大项["+field.ResultValue+"]不为常量库内容，请修改;") == -1 {
					if !HasKey(constMap["faPiaoDaXiangMap2"], field.ResultValue) {
						wrongNote += "发票名称大项[" + field.ResultValue + "]不为常量库内容，请修改;"
					}
				}
			}

			//CSB0108RC0101001
			//当fc204录入值为0时，出导出校验：发票xxx有收取手术费，但手术数量为0，请核查；（xxx为同一发票属性下的fc032或fc054的结果值）
			if obj.Bill.Agency[:3] != "001" && field.Code == "fc204" && field.ResultValue == "0" {
				fc032Value := GetFieldsFinal(fapiao.Fields, "fc032")
				fc054Value := GetFieldsFinal(fapiao.Fields, "fc054")
				if len(fc032Value) != 0 {
					mes := "发票" + fc032Value + "有收取手术费，但手术数量为0，请核查；"
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
					}
				} else {
					mes := "发票" + fc054Value + "有收取手术费，但手术数量为0，请核查；"
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
					}
				}
			}
		}
	}

	for _, fapiao := range obj.InpatientInfo {
		for _, field := range fapiao.Fields {
			if RegIsMatch(field.Code, `^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$`) && field.FinalValue != "" {
				if strings.Index(wrongNote, "发票名称大项["+field.ResultValue+"]不为常量库内容，请修改;") == -1 {
					if !HasKey(constMap["faPiaoDaXiangMap3"], field.ResultValue) {
						wrongNote += "发票名称大项[" + field.ResultValue + "]不为常量库内容，请修改;"
					}
				}
			}

			//CSB0108RC0101001
			//当fc204录入值为0时，出导出校验：发票xxx有收取手术费，但手术数量为0，请核查；（xxx为同一发票属性下的fc032或fc054的结果值）
			if obj.Bill.Agency[:3] != "001" && field.Code == "fc204" && field.ResultValue == "0" {
				fc032Value := GetFieldsFinal(fapiao.Fields, "fc032")
				fc054Value := GetFieldsFinal(fapiao.Fields, "fc054")
				if len(fc032Value) != 0 {
					mes := "发票" + fc032Value + "有收取手术费，但手术数量为0，请核查；"
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
					}
				} else {
					mes := "发票" + fc054Value + "有收取手术费，但手术数量为0，请核查；"
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
					}
				}
			}
		}
	}

	clinicInfo := obj.ClinicInfo                    //门诊发票
	inpatientInfo := obj.InpatientInfo              //住院发票
	fapiaos := append(clinicInfo, inpatientInfo...) //发票内容

	//无门诊发票发票时
	if clinicInfo[0].Fields == nil {
		for i, field := range clinicInfo[0].Fields {
			fmt.Println("field", field)
			fmt.Println("i=", i)
		}
		fapiaos = append(inpatientInfo, clinicInfo...)
	}
	fc032Map := make(map[string]string)
	fc054Map := make(map[string]string)
	var filedArr []string
	yiyuans := map[string]int{}
	for _, fapiao := range fapiaos {
		//qingDans := fapiao.QingDan //清单
		code := fapiao.Code //发票编码
		isExistHy := true
		for _, field := range fapiao.Fields {
			if field.Code == "fc226" || field.Code == "fc227" {
				_, isExits := yiyuans[field.ResultValue]
				if !isExits {
					yiyuans[field.ResultValue] = 0
				}
				yiyuans[field.ResultValue] += 1
			}

			//CSB0108RC0103000
			//global.GLog.Info("CSB0108RC0103000")
			//当fc165、fc196、fc197、fc198、fc199以上字段结果值无法转码时，出一条导出校验：手术填写有误(无法转码)，请核查；
			if strings.Index(wrongNote, "手术填写有误(无法转码)，请核查；") == -1 {
				if RegIsMatch(field.Code, "^(fc165|fc196|fc197|fc198|fc199)$") && field.ResultValue != "" {
					if !HasKey(constMap["shouShuDaiMaMap"], field.ResultValue) {
						wrongNote += "发票【" + code + "】手术填写有误(无法转码)，请核查；"
					}
				}
			}

			//CSB0108RC0093000
			//fc034、fc055录入值不在《B0108_太平理赔_医院代码表》常量库时，出导出校验：就诊医院名称录入错误；
			//global.GLog.Info("CSB0108RC0093000")
			if strings.Index(wrongNote, "发票【"+code+"】就诊医院名称录入错误；") == -1 {
				if RegIsMatch(field.Code, `^(fc034|fc055|fc226|fc227)$`) && field.ResultValue != "" {
					if !HasKey(constMap["yiYuanDaiMaMap"], field.ResultValue) {
						wrongNote += "发票【" + code + "】就诊医院名称录入错误；"
					}
				}
			}

			if field.Code == "fc062" {
				if ParseFloat(GetFieldsFinal(fapiao.Fields, "fc219")) > ParseFloat(field.FinalValue) {
					wrongNote += "发票" + code + "个人支付金额不可大于发票总金额，请核查；"
				}
				if ParseFloat(GetFieldsFinal(fapiao.Fields, "fc220")) > ParseFloat(field.FinalValue) {
					wrongNote += "发票" + code + "账户支付不可大于发票总金额，请核查;"
				}
				if ParseFloat(GetFieldsFinal(fapiao.Fields, "fc222")) > ParseFloat(field.FinalValue) {
					wrongNote += "发票" + code + "补缴金额不可大于发票总金额，请核查;"
				}

			}

			//CSB0108RC0114000
			//"fc055的录入值为B，且fc227为空时，出导出校验：发票xxx未填写其他医院名称；（xxx为该单发票号）
			//fc055的录入值不为B，且fc227不为空时，出导出校验：发票xxx，不应录入其他医院名称。（xxx为该单发票号）"
			//global.GLog.Info("CSB0108RC0114000")
			if field.Code == "fc055" && GetFieldsInput(fapiao.Fields, "fc055") == "B" && GetFieldsFinal(fapiao.Fields, "fc227") == "" {
				////global.GLog.Info("CSB0108RC0114000")
				wrongNote += "发票" + code + "未填写其他医院名称;"
			} else if field.Code == "fc055" && GetFieldsInput(fapiao.Fields, "fc055") != "B" && GetFieldsFinal(fapiao.Fields, "fc227") != "" {
				wrongNote += "发票" + code + "，不应录入其他医院名称;"
			}

			if RegIsMatch(field.Code, `^(fc034|fc226|fc055|fc227)$`) && RegIsMatch(field.ResultValue, `(宁海县第一医院|宁海县妇幼保健院)`) {
				wrongNote += "发票" + code + "为宁海县医院，请检查统筹和自费金额;"
			}

			//CSB0108RC0111000
			//"当fc140或fc142结果值为1时，校验同一发票属性下以下字段录入内容是否有包含“手术费”，无则出导出校验：发票大项中没有“手术费”，请确认；
			//fc113、fc114、fc115、fc116、fc117、fc118、fc119、fc120、fc121、fc122、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc131、fc132"
			//global.GLog.Info("CSB0108RC0111000")
			if GetFieldsFinal(fapiao.Fields, "fc140") == "1" || GetFieldsFinal(fapiao.Fields, "fc142") == "1" {
				if RegIsMatch(field.Code, "^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$") {
					if strings.Index(GetFieldsInput(fapiao.Fields, field.Code), "手术费") != -1 {
						isExistHy = false
					}
				}
			}
			//CSB0108RC0098000
			//global.GLog.Info("CSB0108RC0098000")
			//当fc018、fc112、fc036、fc104、fc058、fc161、fc162、fc205、fc206、fc207、fc208、fc209
			//以上字段任一结果值没有转码，则出一条导出校验：疾病诊断录入有误，请修改；
			if strings.Index(wrongNote, "诊断结论录入内容不在代码表中，请修改；") == -1 {
				if RegIsMatch(field.Code, "^(fc018|fc112|fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$") {
					if constMap["zhenDuanDaiMaMap"][field.ResultValue] != field.FinalValue {
						wrongNote += "诊断结论录入内容不在代码表中，请修改；"
					}
				}
			}
			//CSB0108RC0116000
			//当fc056或fc057的结果值晚于当前系统日期时，出导出校验：发票xxx的出入院日期晚于当前系统时间，请修改；（xxx为同一发票属性下的fc054的结果值）
			//global.GLog.Info("CSB0108RC0116000")
			currentTime, _ := strconv.Atoi(time.Now().Format("20060102")) //获取当前时间 20220927

			if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc054")+"的出入院日期晚于当前系统时间，请修改；") == -1 {
				if ParseFloat(GetFieldsFinal(fapiao.Fields, "fc056")) > float64(currentTime) || ParseFloat(GetFieldsFinal(fapiao.Fields, "fc057")) > float64(currentTime) {
					wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "的出入院日期晚于当前系统时间，请修改；"
				}
			}

			//CSB0108RC0094000
			//"fc034的录入值为B，且fc226为空时，出导出校验：发票xxx未填写其他医院名称；（xxx为该单发号）
			//fc034的录入值不为B，且fc226不为空时，出导出校验：发票xxx，不应录入其他医院名称。（xxx为该单发号）"
			//global.GLog.Info("CSB0108RC0094000")
			if strings.Index(wrongNote, "发票"+fapiao.Code+"未填写其他医院名称；") == -1 {
				if GetFieldsInput(fapiao.Fields, "fc034") == "B" && GetFieldsInput(fapiao.Fields, "fc226") == "" {
					wrongNote += "发票" + fapiao.Code + "未填写其他医院名称；"
				}
			}
			if strings.Index(wrongNote, "发票"+fapiao.Code+"，不应录入其他医院名称。") == -1 {
				if GetFieldsInput(fapiao.Fields, "fc034") != "B" && GetFieldsInput(fapiao.Fields, "fc226") != "" {
					wrongNote += "发票" + fapiao.Code + "，不应录入其他医院名称。"
				}
			}

			//CSB0108RC0095000
			//当fc035的结果值晚于当前系统日期时，出导出校验：发票xxx的就诊日期晚于当前系统时间，请修改；（xxx为同一发票属性下的fc032的结果值）
			//global.GLog.Info("CSB0108RC0095000")
			currentTime, _ = strconv.Atoi(time.Now().Format("20060102")) //获取当前时间 20220927
			if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc032")+"的就诊日期晚于当前系统时间，请修改；") == -1 {
				if ParseFloat(GetFieldsFinal(fapiao.Fields, "fc035")) > float64(currentTime) {
					wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc032") + "的就诊日期晚于当前系统时间，请修改；"
				}
			}

			//if field.Code == "fc035" {
			//	fc035ValueInt, _ := strconv.Atoi(field.ResultValue)
			//	if fc035ValueInt > currentTime {
			//		wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc032") + "的就诊日期晚于当前系统时间，请修改；"
			//	}
			//}
			//CSB0108RC0121000
			//在同一发票属性下，当fc063和fc065结果值相加之和大于或等于fc062结果值时，出导出校验：住院发票报销金额加自费金额大于发票总金额，请检查；
			//global.GLog.Info("CSB0108RC0121000")
			if strings.Index(wrongNote, "住院发票报销金额加自费金额大于发票总金额，请检查；") == -1 {
				if IsExist(fapiao.Fields, "fc063") && IsExist(fapiao.Fields, "fc065") && IsExist(fapiao.Fields, "fc062") {
					if (ParseFloat(GetFieldsFinal(fapiao.Fields, "fc063")) + ParseFloat(GetFieldsFinal(fapiao.Fields, "fc065"))) > ParseFloat(GetFieldsFinal(fapiao.Fields, "fc062")) {
						wrongNote += "住院发票报销金额加自费金额大于发票总金额，请检查；"
					}
				}
			}

			//CSB0108RC0105000
			//在同一发票属性下，当fc040和fc043结果值相加之和大于或等于fc039结果值时，出导出校验：门诊发票报销金额加自费金额大于发票总金额，请检查；
			if obj.Bill.Agency[:3] != "001" {
				if strings.Index(wrongNote, "门诊发票报销金额加自费金额大于发票总金额，请检查；") == -1 {
					if IsExist(fapiao.Fields, "fc040") && IsExist(fapiao.Fields, "fc043") && IsExist(fapiao.Fields, "fc039") {
						if ParseFloat(GetFieldsFinal(fapiao.Fields, "fc040"))+ParseFloat(GetFieldsFinal(fapiao.Fields, "fc043")) >= ParseFloat(GetFieldsFinal(fapiao.Fields, "fc039")) {
							wrongNote += "门诊发票报销金额加自费金额大于发票总金额，请检查；"
						}
					}
				}
			}

			//CSB0108RC0127000
			//当fc077结果值为0或为空时，出导出校验：住院发票[xxx]大项-床位费金额为0或为空，请核查；（xxx为同一发票属性下fc054结果值）
			//global.GLog.Info("CSB0108RC0127000")
			if field.Code == "fc077" && !RegIsMatch(obj.Bill.Agency, `^(00183000|00183010|00183002|00183012|00183300|00183310|00183301|00183311)$`) && RegIsMatch(field.FinalValue, "^(0|0.00| )$") {
				wrongNote += "住院发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "大项-床位费金额为0或为空，请核查；"
			}

			//CSB0108RC0123000
			//在同一发票属性下，fc062与fc063、fc064、fc219、fc220结果值的差不等于0时（fc062 - fc063 - fc064 - fc219 - fc220 != 0）
			//出导出校验：XX发票检查统筹、其他扣除、个人支付、账户支付是否有录入错误（备注：其中XX是fc054）
			//global.GLog.Info("CSB0108RC0123000")
			hyField = []string{"fc063", "fc064", "fc219", "fc220"}
			hyFieldFinalValue := 0.0
			for index := range hyField {
				hyFieldFinalValue += ParseFloat(GetFieldsFinal(fapiao.Fields, hyField[index]))
			}
			chae := ParseFloat(GetFieldsFinal(fapiao.Fields, "fc062")) - hyFieldFinalValue
			if field.Code == "fc062" && strconv.FormatFloat(chae, 'f', 2, 64) != "0.00" && strconv.FormatFloat(chae, 'f', 2, 64) != "-0.00" {
				wrongNote += GetFieldsFinal(fapiao.Fields, "fc054") + "发票检查统筹、其他扣除、个人支付、账户支付是否有录入错误；"
			}

			//CSB0108RC0115000
			//同一发票属性下，将“fc214与fc056”、“fc215与fc057”录入值作对比，出现不一致情况出一条导出校验：发票XXX，出入院日期不一致；（XXX为同一发票属性下，fc054的结果值）
			//global.GLog.Info("CSB0108RC0115000")
			if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc054")+"，出入院日期不一致；") == -1 {
				if IsExist(fapiao.HospitalizationDate, "fc214") && IsExist(fapiao.Fields, "fc056") && IsExist(fapiao.HospitalizationDate, "fc215") && IsExist(fapiao.Fields, "fc057") {
					if (GetFieldsInput(fapiao.HospitalizationDate, "fc214") != GetFieldsInput(fapiao.Fields, "fc056")) || (GetFieldsInput(fapiao.HospitalizationDate, "fc215") != GetFieldsInput(fapiao.Fields, "fc057")) {
						wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "，出入院日期不一致；"
					}
				}
			}

			//CSB0108RC0122000
			////global.GLog.Info("CSB0108RC0122000")
			//"fc066、fc067、fc068、fc069、fc070、fc071、fc072、fc073、fc074、fc075、fc076、fc077、fc078、fc079的结果值相加不等于fc062的值时，出一条导出校验：住院发票[xxx]明细金额与发票总金额不一致，差额为[XXX]，请核查
			//（第一个xxx为fc054结果值，第二个XXX为fc062-（fc066+fc067+fc068+fc069+fc070+fc071+fc072+fc073+fc074+fc075+fc076+fc077+fc078+fc079）的差额）（当明细金额相加大于fc062时，需在差额前面加""-""）"
			hyField = []string{"fc066", "fc067", "fc068", "fc069", "fc070", "fc071", "fc072", "fc073", "fc074", "fc075", "fc076", "fc077", "fc078", "fc079"}
			fieldsFinalAll := 0.0

			for i := range hyField {
				fieldsFinalAll += ParseFloat(GetFieldsFinal(fapiao.Fields, hyField[i]))
			}
			if IsExist(fapiao.Fields, "fc062") {

				if fieldsFinalAll != ParseFloat(GetFieldsFinal(fapiao.Fields, "fc062")) {
					difference := ParseFloat(GetFieldsFinal(fapiao.Fields, "fc062")) - fieldsFinalAll
					// fmt.Println("----differencedifferencedifference----", difference, strconv.FormatFloat(difference, 'f', 2, 64), strconv.FormatFloat(difference, 'f', 2, 64) != "0.00")
					if strconv.FormatFloat(difference, 'f', 2, 64) != "0.00" && strconv.FormatFloat(difference, 'f', 2, 64) != "-0.00" && strings.Index(wrongNote, "住院发票"+GetFieldsFinal(fapiao.Fields, "fc054")+"明细金额与发票总金额不一致，差额为"+strconv.FormatFloat(difference, 'f', 2, 64)+"，请核查；") == -1 {
						wrongNote += "住院发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "明细金额与发票总金额不一致，差额为" + strconv.FormatFloat(difference, 'f', 2, 64) + "，请核查；"
					}
				}
			}

			//CSB0108RC0106000
			//"在同一发票属性下，fc044、fc046、fc047、fc048、fc049、fc050、fc051、fc052、fc053、fc100、fc102的结果值相加不等于fc039的值时，出一条导出校验：门诊发票[xxx]明细金额与发票总金额不一致，差额为[XXX]，请核查
			//（第一个xxx为fc032的结果值，第二个XXX为fc039-（fc044+fc046+fc047+fc048+fc049+fc050+fc051+fc052+fc053+fc100+fc102）的差额）（当明细金额相加大于fc039时，需在差额前面加""-""）"
			//global.GLog.Info("CSB0108RC0106000")

			hyField = []string{"fc044", "fc046", "fc047", "fc048", "fc049", "fc050", "fc051", "fc052", "fc053", "fc100", "fc102"}
			fieldsFinalAll = 0.0

			for i := range hyField {
				fieldsFinalAll += ParseFloat(GetFieldsFinal(fapiao.Fields, hyField[i]))
			}
			//fmt.Println(ParseFloat(GetFieldsFinal(fields, "fc039")))
			if IsExist(fapiao.Fields, "fc039") {
				Difference := ParseFloat(GetFieldsFinal(fapiao.Fields, "fc039")) - fieldsFinalAll
				if strings.Index(wrongNote, "门诊发票"+GetFieldsFinal(fapiao.Fields, "fc032")+"明细金额与发票总金额不一致，差额为"+strconv.FormatFloat(Difference, 'f', 2, 64)+"，请核查；") == -1 {
					if strconv.FormatFloat(Difference, 'f', 2, 64) != "0.00" && strconv.FormatFloat(Difference, 'f', 2, 64) != "-0.00" && fieldsFinalAll != ParseFloat(GetFieldsFinal(fapiao.Fields, "fc039")) {
						wrongNote += "门诊发票" + GetFieldsFinal(fapiao.Fields, "fc032") + "明细金额与发票总金额不一致，差额为" + strconv.FormatFloat(Difference, 'f', 2, 64) + "，请核查；"
					}
				}
			}

			//CSB0108RC0124000
			//同一发票属性下，当存在fc107时，将fc107与fc063结果值作对比（当存在多个fc107时，需将每个fc107的结果值相加求和再与fc063结果值作对比），不一致出导出校验：
			//发票XXX报销单统筹支付金额与发票统筹支付金额录入不一致，请核查；（XXX为同一发票属性下，fc054的结果值）
			//global.GLog.Info("CSB0108RC0124000")

			if IsExist(fapiao.BaoXiaoDan, "fc107") && IsExist(fapiao.Fields, "fc054") {
				if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc054")+"报销单统筹支付金额与发票统筹支付金额录入不一致，请核查；") == -1 {
					if GetFieldsAllFinal(fapiao.BaoXiaoDan, "fc107")-ParseFloat(GetFieldsFinal(fapiao.Fields, "fc063")) != 0 {
						wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "报销单统筹支付金额与发票统筹支付金额录入不一致，请核查；"
					}
				}
			}

			//CSB0108RC0125000
			//同一发票属性下，当存在fc212时，将fc212与fc065结果值作对比（当存在多个fc212时，需将每个fc212的结果值相加求和再与fc065结果值作对比），不一致出导出校验：
			//发票XXX报销单自费金额与发票自费金额录入不一致，请核查；（XXX为同一发票属性下，fc054的结果值）（当不存在fc212字段时，不应进行该校验）
			//global.GLog.Info("CSB0108RC0125000")
			if IsExist(fapiao.BaoXiaoDan, "fc212") && IsExist(fapiao.Fields, "fc054") {
				if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc054")+"报销单自费金额与发票自费金额录入不一致，请核查；") == -1 {
					if GetFieldsAllFinal(fapiao.BaoXiaoDan, "fc212")-ParseFloat(GetFieldsFinal(fapiao.Fields, "fc065")) != 0 {
						wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "报销单自费金额与发票自费金额录入不一致，请核查；"
					}
				}
			}

			//CSB0108RC0108000
			//同一发票属性下，当存在fc107时，将fc107与fc040结果值作对比（当存在多个fc107时，需将每个fc107的结果值相加求和再与fc040结果值作对比），不一致出导出校验：
			//发票XXX报销单统筹支付金额与发票统筹支付金额录入不一致，请核查；（XXX为同一发票属性下，fc032的结果值）
			//global.GLog.Info("CSB0108RC0108000")

			if IsExist(fapiao.BaoXiaoDan, "fc107") && IsExist(fapiao.Fields, "fc032") {
				if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc032")+"报销单统筹支付金额与发票统筹支付金额录入不一致，请核查；") == -1 {

					if GetFieldsAllFinal(fapiao.BaoXiaoDan, "fc107")-ParseFloat(GetFieldsFinal(fapiao.Fields, "fc040")) != 0 {
						wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc032") + "报销单统筹支付金额与发票统筹支付金额录入不一致，请核查；"
					}
				}
			}

			//CSB0108RC0110000
			//同一发票属性下，当存在fc212时，将fc212与fc043结果值作对比（当存在多个fc212时，需将每个fc212的结果值相加求和再与fc043结果值作对比），不一致出导出校验：
			//发票XXX报销单自费金额与发票自费金额录入不一致，请核查；（XXX为同一发票属性下，fc032的结果值）（当不存在fc212字段时，不应进行该校验）
			//global.GLog.Info("CSB0108RC0110000")
			if IsExist(fapiao.BaoXiaoDan, "fc212") && IsExist(fapiao.Fields, "fc032") {
				if strings.Index(wrongNote, "发票"+GetFieldsFinal(fapiao.Fields, "fc032")+"报销单自费金额与发票自费金额录入不一致，请核查；") == -1 {
					if GetFieldsAllFinal(fapiao.BaoXiaoDan, "fc212")-ParseFloat(GetFieldsFinal(fapiao.Fields, "fc043")) != 0 {
						wrongNote += "发票" + GetFieldsFinal(fapiao.Fields, "fc032") + "报销单自费金额与发票自费金额录入不一致，请核查；"
					}
				}
			}

			//CSB0108RC0112000
			//"当以下字段录入内容包含“手术费”，且fc165结果值为空时，出导出校验：发票大项存在手术费，手术名称未填写，请确认；
			//fc113、fc114、fc115、fc116、fc117、fc118、fc119、fc120、fc121、fc122、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc131、fc132"
			//global.GLog.Info("CSB0108RC0112000")
			//修改导出校验提示：发票xxx发票大项存在手术费，手术名称未填写，请确认；（xxx为同一发票下对应的fc032或fc054的值）
			var fc032Result, fc054Result string
			if IsExist(fapiao.Fields, "fc032") {
				fc032Result = GetFieldsFinal(fapiao.Fields, "fc032")
				if RegIsMatch(field.Code, "^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$") {
					if strings.Index(wrongNote, "发票"+fc032Result+"发票大项存在手术费，手术名称未填写，请确认;") == -1 {
						if strings.Index(field.ResultValue, "手术费") != -1 && GetFieldsFinal(fields, "fc165") == "" {
							wrongNote += "发票" + fc032Result + "发票大项存在手术费，手术名称未填写，请确认;"
							fmt.Println("fc032Result=", fc032Result)
						}
					}
				}
			} else if IsExist(fapiao.Fields, "fc054") {
				fc054Result = GetFieldsFinal(fapiao.Fields, "fc054")
				if RegIsMatch(field.Code, "^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$") {
					if strings.Index(wrongNote, "发票"+fc054Result+"发票大项存在手术费，手术名称未填写，请确认;") == -1 {
						if strings.Index(field.ResultValue, "手术费") != -1 && GetFieldsFinal(fields, "fc165") == "" {
							wrongNote += "发票" + fc054Result + "发票大项存在手术费，手术名称未填写，请确认;"
							fmt.Println("fc054Result=", fc054Result)

						}
					}
				}
			}

			//CSB0108RC0129000
			//校验所有字段的录入值，当包含?或？时，出导出校验：案件中存在?号，请核实；
			//global.GLog.Info("CSB0108RC0129000")
			if IsExist(fields, "fc032") {
				//fc032 = GetFieldsFinal(fields, "fc032")
				//if strings.Index(wrongNote, "发票"+fc032val+"的"+hyFields2.Name+"存在?号，请核实；") == -1 {
				if RegIsMatch(field.ResultValue, "\\?|？") {
					filedArr = append(filedArr, field.Name)
				}
				//}
			} else if IsExist(fields, "fc054") {
				//fc054 = GetFieldsFinal(fields, "fc054")
				//if strings.Index(wrongNote, "发票"+fc032val+"的"+hyFields2.Name+"存在?号，请核实；") == -1 {
				if RegIsMatch(field.ResultValue, "\\?|？") {
					filedArr = append(filedArr, field.Name)
				}
				//}
			}

		}

		//项目编码:CSB0108RC0111001     20230711 - 修改导出校验提示：发票xxx没有手术费，请确认；（xxx为同一发票下对应的fc032或fc054的值）
		var fc032Val, fc054Val string
		if obj.Bill.Agency[:3] != "001" {
			if IsExist(fapiao.Fields, "fc032") {
				fc032Val = GetFieldsFinal(fapiao.Fields, "fc032")
				if IsExist(fapiao.Fields, "fc142") || IsExist(fapiao.Fields, "fc140") {
					if GetFieldsFinal(fapiao.Fields, "fc142") == "1" || GetFieldsFinal(fapiao.Fields, "fc140") == "1" {
						if isExistHy {
							if strings.Index(wrongNote, "发票"+fc032Val+"没有手术费，请确认") == -1 {
								wrongNote += "发票" + fc032Val + "没有手术费，请确认;"
							}
						}
					}
				}
			} else if IsExist(fapiao.Fields, "fc054") {
				fc054Val = GetFieldsFinal(fapiao.Fields, "fc054")
				if IsExist(fapiao.Fields, "fc142") || IsExist(fapiao.Fields, "fc140") {
					if GetFieldsFinal(fapiao.Fields, "fc142") == "1" || GetFieldsFinal(fapiao.Fields, "fc140") == "1" {
						if isExistHy {
							if strings.Index(wrongNote, "发票"+fc054Val+"没有手术费，请确认") == -1 {
								wrongNote += "发票" + fc054Val + "没有手术费，请确认;"
							}
						}
					}
				}
			}
		}

		//CSB0108RC0113000
		//"fc081、fc144、fc146、fc148、fc150、fc152、fc154、fc156以上字段为清单循环分块字段（即可能会存在多个的情况），
		//在同一发票属性下，当他们的结果值相加不等于fc039的结果值时，出一条导出校验：门诊发票[xxx]清单项目金额与发票总金额不一致，差额为[XXX]，请核查（当不存在fc081字段时，不应进行该校验）
		//（第一个xxx为fc032的结果值，第二个XXX为fc039-（所有fc081+fc144+fc146+fc148+fc150+fc152+fc154+fc156）的差额）（当金额相加大于fc039时，需在差额前面加""-""）"
		// fieldAllValue := 0.0
		// for _, qingDan := range fapiao.QingDan {
		// 	for _, qingDanFeild := range qingDan.Fields {
		// 		if IsExist(fapiao.Fields, "fc039") && IsExist(qingDan.Fields, "fc081") {
		// 			if RegIsMatch(qingDanFeild.Code, "^(fc081|fc144|fc146|fc148|fc150|fc152|fc154|fc156)$") {
		// 				fieldAllValue += ParseFloat(qingDanFeild.FinalValue)
		// 			}
		// 		}
		// 	}
		// }
		// // fmt.Println("-------------fieldAllValue---------------", fieldAllValue, GetFieldsFinal(fapiao.Fields, "fc039"))
		// difference := ParseFloat(GetFieldsFinal(fapiao.Fields, "fc039")) - fieldAllValue
		// differenceStr := strconv.FormatFloat(difference, 'f', 2, 64)
		// if IsExist(fapiao.Fields, "fc039") && IsExist(fapiao.QingDan[0].Fields, "fc081") {
		// 	if difference != 0.0 && differenceStr != "0.00" && differenceStr != "-0.00" {
		// 		if strings.Index(wrongNote, "门诊发票"+GetFieldsFinal(fapiao.Fields, "fc032")+"清单项目金额与发票总金额不一致，差额为"+differenceStr+"，请核查；") == -1 {
		// 			wrongNote += "门诊发票" + GetFieldsFinal(fapiao.Fields, "fc032") + "清单项目金额与发票总金额不一致，差额为" + differenceStr + "，请核查；"
		// 		}
		// 	}
		// }

		// //CSB0108RC0128000
		// //"fc081、fc144、fc146、fc148、fc150、fc152、fc154、fc156以上字段为清单循环分块字段（即可能会存在多个的情况），在同一发票属性下，
		// //当他们的结果值相加不等于fc062的结果值时，出一条导出校验：住院发票[xxx]清单项目金额与发票总金额不一致，差额为[XXX]，请核查（当不存在fc081字段时，不应进行该校验）
		// //（第一个xxx为fc054结果值，第二个XXX为fc062-（fc081+fc144+fc146+fc148+fc150+fc152+fc154+fc156）的差额）"
		// //global.GLog.Info("CSB0108RC0128000")
		// fieldAllValue = 0.0
		// for _, qingDan := range fapiao.QingDan {
		// 	for _, qingDanFeild := range qingDan.Fields {
		// 		if IsExist(fapiao.Fields, "fc062") && IsExist(qingDan.Fields, "fc081") {
		// 			if RegIsMatch(qingDanFeild.Code, "^(fc081|fc144|fc146|fc148|fc150|fc152|fc154|fc156)$") {
		// 				fieldAllValue += ParseFloat(qingDanFeild.FinalValue)
		// 			}
		// 		}
		// 	}
		// }
		// // fmt.Println("-------")
		// // fmt.Println(fieldAllValue)
		// // fmt.Println("=======")
		// // fmt.Println(GetFieldsFinal(fapiao.Fields, "fc062"))
		// difference = ParseFloat(GetFieldsFinal(fapiao.Fields, "fc062")) - fieldAllValue
		// differenceStr = strconv.FormatFloat(difference, 'f', 2, 64)
		// if fieldAllValue != ParseFloat(GetFieldsFinal(fapiao.Fields, "fc062")) {
		// 	if IsExist(fapiao.Fields, "fc062") && IsExist(fapiao.QingDan[0].Fields, "fc081") && differenceStr != "0.00" && differenceStr != "-0.00" {
		// 		if strings.Index(wrongNote, "住院发票"+GetFieldsFinal(fapiao.Fields, "fc054")+"清单项目金额与发票总金额不一致，差额为"+differenceStr+"，请核查；") == -1 {
		// 			wrongNote += "住院发票" + GetFieldsFinal(fapiao.Fields, "fc054") + "清单项目金额与发票总金额不一致，差额为" + differenceStr + "，请核查；"
		// 		}
		// 	}
		// }

	}

	aaa := ""
	bbb := ""
	fmt.Println("-------------yiyuans-----------------", yiyuans)
	for _, fapiao := range fapiaos {
		code := fapiao.Code //发票编码
		for _, field := range fapiao.Fields {
			if field.Code == "fc226" || field.Code == "fc227" {
				num, _ := yiyuans[field.ResultValue]
				for key, value := range yiyuans {
					if key != field.ResultValue {
						if num < value {
							aaa += "【" + code + "】、"
						} else if num == value {
							bbb = "其他医院名称与其他发票不一致，请检查；"
						}

					}
				}

			}
		}
	}
	if aaa != "" {
		wrongNote += "发票" + aaa + "其他医院名称与其他发票不一致，请检查；"
	}
	if bbb != "" {
		wrongNote += bbb
	}

	//CSB0108RC0089000
	//校验所有fc032和fc054的结果值是否存在一致的情况，一致则出导出校验：发票号xxx重复，请核查；
	//（xxx为重复的fc032或fc054的结果值）（为空不校验）
	//global.GLog.Info("CSB0108RC0089000")

	//fc032val := GetFieldsFinal(fields, "fc032")
	//fc054val := GetFieldsFinal(fields, "fc054")
	totalVal := 0
	for _, index := range fapiaos {
		for i, hyFields2 := range index.Fields {
			if i == len(index.Fields) {

			}

			if hyFields2.Code == "fc032" || hyFields2.Code == "fc054" {
				_, ok1 := fc032Map[hyFields2.FinalValue]
				if ok1 {
					if strings.Index(wrongNote, "发票号"+hyFields2.FinalValue+"重复，请核查；") == -1 {
						wrongNote += "发票号" + hyFields2.FinalValue + "重复，请核查；"
					}
				} else {
					fc032Map[hyFields2.FinalValue] = "1"
					_, ok2 := fc054Map[hyFields2.FinalValue]
					if ok2 {
						if strings.Index(wrongNote, "发票号"+hyFields2.FinalValue+"重复，请核查；") == -1 {
							wrongNote += "发票号" + hyFields2.FinalValue + "重复，请核查；"
						}
					} else {
						fc054Map[hyFields2.FinalValue] = "1"
					}
				}
			}

			//CSB0108RC0075000
			//global.GLog.Info("CSB0108RC0075000")
			//当fc217存在不为空的结果值时，与所有的fc035和fc056的结果值作对比，若存在fc217日期比较晚的情况，则出导出校验：首次就诊日期不能早于事故日期；
			//global.GLog.Info("CSB0108RC0075000")
			//修改导出校验提示：发票xxx就诊日期不能早于事故日期；（xxx为同一发票下对应的fc032或fc054的值）
			if hyFields2.Code == "fc032" || hyFields2.Code == "fc054" {
				if GetFieldsFinal(fields, "fc217") != "" {
					if IsExist(index.Fields, "fc035") {
						if strings.Index(wrongNote, "发票"+hyFields2.FinalValue+"首次就诊日期不能早于事故日期;") == -1 {
							if GetFieldsFinal(fields, "fc217") > GetFieldsFinal(index.Fields, "fc035") {
								wrongNote = wrongNote + "发票" + hyFields2.FinalValue + "首次就诊日期不能早于事故日期;"
							}
						}
					}
					if IsExist(index.Fields, "fc056") {
						if strings.Index(wrongNote, "发票"+hyFields2.FinalValue+"首次就诊日期不能早于事故日期;") == -1 {
							if GetFieldsFinal(fields, "fc217") > GetFieldsFinal(index.Fields, "fc056") {
								wrongNote = wrongNote + "发票" + hyFields2.FinalValue + "首次就诊日期不能早于事故日期;"
							}
						}
					}
				}
			}

			//global.GLog.Info("CSB0108RC0078000")
			//当存在fc033或fc059字段（或存在多个）时，将每个fc033、fc059字段的结果值与fc216结果值做对比，
			//不一致时出导出校验：病人姓名不一致，请确认后修改；
			//20230711 - 修改导出校验提示：发票xxx病人姓名与诊断姓名不一致，请确认后修改。（xxx为同一发票下对应的fc032或fc054的值）
			if hyFields2.Code == "fc032" || hyFields2.Code == "fc054" {
				if strings.Index(wrongNote, "发票"+hyFields2.FinalValue+"病人姓名不一致，请确认后修改;") == -1 {
					if IsExist(index.Fields, "fc033") && GetFieldsFinal(index.Fields, "fc033") != GetFieldsFinal(fields, "fc216") {
						wrongNote = wrongNote + "发票" + hyFields2.FinalValue + "病人姓名不一致，请确认后修改;"
					}
					if IsExist(index.Fields, "fc059") && GetFieldsFinal(index.Fields, "fc059") != GetFieldsFinal(fields, "fc216") {
						wrongNote = wrongNote + "发票" + hyFields2.FinalValue + "病人姓名不一致，请确认后修改;"
					}
				}
			}

			//CSB0108RC0330000 "同一发票下，以下字段的录入值包含两个或两个以上“床位费”时，出导出校验：xxx发票，请核实床位费；（xxx为同一发票下fc054的结果值）
			//fc113、fc114、fc115、fc116、fc117、fc118、fc119、fc120、fc121、fc122、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc131、fc132"
			if RegIsMatch(hyFields2.Code, "^(fc113|fc114|fc115|fc116|fc117|fc118|fc119|fc120|fc121|fc122|fc123|fc124|fc125|fc126|fc127|fc128|fc129|fc130|fc131|fc132)$") {
				if strings.Index(wrongNote, GetFieldsFinal(fields, "fc054")+"发票，请核实床位费；") == -1 {
					if strings.Index(hyFields2.ResultValue, "床位费") != -1 {
						totalVal++
						if totalVal >= 2 {
							wrongNote += GetFieldsFinal(fields, "fc054") + "发票，请核实床位费；"
						}
					}
				}
			}

		}
	}
	//编码CSB0108RC0129001 1.校验所有字段的录入值，当包含?或？时，出导出校验：发票【xxx】的【yyy】存在?号，请核实；20230713新增
	//2.如同一发票下多个字段包含?或？，则导出校验提示格式为：发票【xxx】的【yyy1、yyy2】存在?号，请核实；（xxx为发票号fc032或fc054的值，yyy为包含问号的字段名）
	//3.当最后一页字段包含?或？时，则出导出校验：最后一页存在?，请检查；20230713新增
	var clinicField []string
	var inpatientField []string
	for _, fieldsMap := range obj.ClinicInfo {
		_, fc032 := GetOneField(fieldsMap.Fields, "fc032", true)
		for _, field := range fieldsMap.QingDan {
			for _, qdField := range field.Fields {
				if RegIsMatch(qdField.ResultValue, `[?？]`) {
					clinicField = append(clinicField, qdField.Name)
				}
			}
		}
		for _, field := range fieldsMap.BaoXiaoDan {
			if RegIsMatch(field.ResultValue, `[?？]`) {
				clinicField = append(clinicField, field.Name)
			}
		}

		for _, field := range fieldsMap.HospitalizationDate {
			if RegIsMatch(field.ResultValue, `[?？]`) {
				clinicField = append(clinicField, field.Name)
			}
		}
		for j, field := range fieldsMap.Fields {
			if RegIsMatch(field.ResultValue, `[?？]`) {
				clinicField = append(clinicField, field.Name)
			}

			if j == len(fieldsMap.Fields)-1 {
				if len(clinicField) > 0 {
					join := strings.Join(clinicField, "、")
					mes := "发票" + fc032 + "的【" + join + "】存在?号，请核实;"
					if strings.Index(wrongNote, mes) == -1 {
						fmt.Println("--filedArr.=", clinicField)
						wrongNote += mes
						clinicField = make([]string, 0)
						fmt.Println("--filedArr null.=", clinicField)
					}
				}
			}

		}
	}

	for _, fieldsMap := range obj.InpatientInfo {
		_, fc054 := GetOneField(fieldsMap.Fields, "fc054", true)
		for _, field := range fieldsMap.QingDan {
			for _, qdField := range field.Fields {
				if RegIsMatch(qdField.ResultValue, `[?？]`) {
					inpatientField = append(inpatientField, qdField.Name)
				}
			}
		}
		for _, field := range fieldsMap.BaoXiaoDan {
			if RegIsMatch(field.ResultValue, `[?？]`) {
				inpatientField = append(inpatientField, field.Name)
			}
		}

		for _, field := range fieldsMap.HospitalizationDate {
			if RegIsMatch(field.ResultValue, `[?？]`) {
				inpatientField = append(inpatientField, field.Name)
			}
		}

		for j, field := range fieldsMap.Fields {
			if RegIsMatch(field.ResultValue, `[?？]`) {
				inpatientField = append(inpatientField, field.Name)
			}
			if j == len(fieldsMap.Fields)-1 {
				if len(inpatientField) > 0 {
					join := strings.Join(inpatientField, "、")
					mes := "发票" + fc054 + "的【" + join + "】存在?号，请核实;"
					if strings.Index(wrongNote, mes) == -1 {
						fmt.Println("--filedArr.=", inpatientField)
						wrongNote += mes
						inpatientField = make([]string, 0)
						fmt.Println("--filedArr null.=", inpatientField)
					}
				}
			}

		}

	}

	fmt.Println("111")
	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""

	//CSB0108RC0088000
	//当base_info大节点下的accident_id节点值为空时，出导出校验：事故原因为空，请检查；
	accidentType := utils2.GetNodeDecimalValue(xmlValue, "accident_type")
	hyField := GetNode(xmlValue, "base_info")
	for _, value := range hyField {
		if GetNodeValue(value, "accident_id") == "" {
			wrongNote += "事故原因为空，请检查；"
		}
		if RegIsMatch(GetNodeValue(value, "accident_id"), `[\u4e00-\u9fa5]`) && strings.Index(wrongNote, "诊断结论不能为中文；") == -1 {
			wrongNote += "诊断结论不能为中文；"
		}
	}

	//CSB0108RC0117000
	//当out_date的值早于in_date的值时，出导出校验：
	//XXX发票号的出院日期早于入院日期，请修改（xxx为与out_date在同一个inpatient_record大节点下的operation_no的值）
	hyField = GetNode(xmlValue, "inpatient_record")
	for _, value := range hyField {
		if GetNodeValue(value, "out_date") < GetNodeValue(value, "in_date") {
			wrongNote += GetNodeValue(value, "operation_no") + "发票号的出院日期早于入院日期，请修改；"
		}

		bed_amt_days := GetNodeValue(value, "bed_amt_days")
		if bed_amt_days != "" && ParseFloat(bed_amt_days) > 60 {
			wrongNote += GetNodeValue(value, "operation_no") + "住院天数大于60天，请确认；"
		}

		//CSB0108RC0118000
		//CSB0108RC0118002 取消需求   "在同一<inpatient_record>大节点下，当<bed>节点值除以<bed_amt_days>节点值的值大于80时，出导出校验：xxx发票，请核实床位费；（xxx为同一大节点下的<operation_no>的节点值）
		//当<bed_amt_days>为0时，将其视为1而后再用<bed>节点值除以<bed_amt_days>节点值
		//如<bed_amt_days>为0，<bed>为81，则用81/1=81>80出导出校验"

		//在同一<inpatient_record>大节点下，当<bed>节点值除以<bed_amt_days>节点值的值大于80时，出导出校验：
		//xxx发票，请核实床位费；（xxx为同一大节点下的<operation_no>的节点值）
		if bed_amt_days == "0" {
			bed_amt_days = "1"
		}
		//if ParseFloat(GetNodeValue(value, "bed"))/ParseFloat(bed_amt_days) > 80 {
		//	//wrongNote += GetNodeValue(value, "operation_no") + "发票，请核实床位费；"
		//}

		//CSB0108RC0119000
		//同一个<inpatient_record>下，将<diagnostic_amt_days>节点值与<bed_amt_days>节点值作对比，差额大于2时，出导出校验：发票XXX住院诊查费天数与床位费天数的差额为【xxx】天，请核查；
		//（XXX为同一个<inpatient_record>下的<operation_no>节点值，xxx为<diagnostic_amt_days>节点值与<bed_amt_days>节点值的差额）
		// fmt.Println("00000")
		chae := math.Abs(ParseFloat(GetNodeValue(value, "diagnostic_amt_days")) - ParseFloat(GetNodeValue(value, "bed_amt_days")))
		// fmt.Println(GetNodeValue(value, "diagnostic_amt_days"))
		// fmt.Println(GetNodeValue(value, "bed_amt_days"))
		// fmt.Println(chae)
		if chae > 2.0 {
			if ParseFloat(GetNodeValue(value, "diagnostic_amt_days")) < ParseFloat(GetNodeValue(value, "bed_amt_days")) {
				wrongNote += "发票" + GetNodeValue(value, "operation_no") + "住院诊查费天数与床位费天数的差为【" + strconv.Itoa(int(chae)) + "】天，请核查；"
			} else {
				wrongNote += "发票" + GetNodeValue(value, "operation_no") + "住院诊查费天数与床位费天数的差为【-" + strconv.Itoa(int(chae)) + "】天，请核查；"
			}
		}
	}

	//CSB0108RC0082000
	//当accident_type节点值不为2时，出导出校验：诊断为意外，请核实事故日期是否为发生日期。
	if !accidentType.Equal(decimal.NewFromInt(2)) {
		wrongNote += "诊断为意外，请核实事故日期是否为发生日期。"
	}

	if RegIsMatch(obj.Bill.Agency, `^(00083000|00083002|00083300|00083301)$`) {

		items := []string{"inpatient_record", "clinic_record"}
		for _, item := range items {
			records := GetNode(xmlValue, item)
			for _, record := range records {
				sum := 0.00
				sum_amount := GetNodeValue(record, "sum_amount")
				if sum_amount != "" {
					item_amounts := RegMatchAll(record, `<item_amount>.+<\/item_amount>`)
					for _, item_amount := range item_amounts {
						item_amount = GetNodeValue(item_amount, "item_amount")
						sum = utils2.SumFloat(sum, ParseFloat(item_amount), "+")
					}
					chae := utils2.SumFloat(ParseFloat(sum_amount), sum, "-")
					if chae != 0.00 && chae != -0.00 {
						if item == "inpatient_record" {
							operation_no := GetNodeValue(record, "operation_no")
							wrongNote += "住院发票[" + operation_no + "]清单项目金额与发票总金额不一致，差额为[" + ToString(chae) + "]，请核查；"
						} else {
							operation_no := GetNodeValue(record, "clinic_no")
							wrongNote += "门诊发票[" + operation_no + "]清单项目金额与发票总金额不一致，差额为[" + ToString(chae) + "]，请核查；"

						}
					}
				}
			}
		}
	}

	if obj.Bill.SaleChannel == "理赔" {
		nItems := []string{"myversion", "case_no", "return_flag"}
		wrongNote += CheckNotNill(xmlValue, "通用节点值不能为空，请核查；", nItems)
		nodeValues := GetNode(xmlValue, `accident_info`)
		for _, xValues := range nodeValues {
			nItems = []string{"insured_name", "accident_time", "accident_type"}
			wrongNote += CheckNotNill(xValues, "通用节点值不能为空，请核查；", nItems)
		}
		nodeValues = GetNode(xmlValue, `base_info`)
		for _, xValues := range nodeValues {
			nItems = []string{"name", "time", "accident_id", "apply_type", "social_security", "isCopy", "apply_certi_type", "apply_certi_code"}
			wrongNote += CheckNotNill(xValues, "通用节点值不能为空，请核查；", nItems)
		}
	} else if obj.Bill.SaleChannel == "秒赔" {
		nItems := []string{"myversion", "case_no", "return_flag"}
		wrongNote += CheckNotNill(xmlValue, "通用节点值不能为空，请核查；", nItems)
		nodeValues := GetNode(xmlValue, `accident_info`)
		for _, xValues := range nodeValues {
			nItems = []string{"accident_time", "accident_type"}
			wrongNote += CheckNotNill(xValues, "通用节点值不能为空，请核查；", nItems)
		}
		nodeValues = GetNode(xmlValue, `base_info`)
		for _, xValues := range nodeValues {
			nItems = []string{"accident_id", "social_security", "isCopy"}
			wrongNote += CheckNotNill(xValues, "通用节点值不能为空，请核查；", nItems)
		}
	}

	nodeValues := GetNode(xmlValue, `pay_info`)
	for _, xValues := range nodeValues {
		nodeValue := GetNodeValue(xValues, "pay_mode")
		if nodeValue == "3" {
			nItems := []string{"acco_name", "acco_certi_code", "bank_account"}
			wrongNote += CheckNotNill(xValues, "账户节点值不能为空，请核查；", nItems)
		}

	}

	fc084Arr := []string{}
	// fc173 := ""
	for _, field := range obj.Fields {
		if field.Code == "fc084" {
			fc084Arr = append(fc084Arr, field.ResultValue)
		}
		// if field.Code == "fc173" {
		// 	fc173 = field.ResultValue
		// }
	}
	if utils2.HasItem(fc084Arr, "10") {
		nodeValues := GetNode(xmlValue, `extension_info`)
		for _, xValues := range nodeValues {
			nItems := []string{"confirm_date", "accident_id"}
			wrongNote += CheckNotNill(xValues, "非医疗节点值不能为空，请核查；", nItems)
		}
	}
	// if fc173 == "1" {
	// 	nodeValues := GetNode(xmlValue, `agency_info`)
	// 	for _, xValues := range nodeValues {
	// 		nItems := []string{"agency_name", "agency_certi_type", "agency_certi_code", "agency_tel"}
	// 		wrongNote += CheckNotNill(xValues, "代办人节点值不能为空，请核查；", nItems)
	// 	}
	// }

	nodeValues = GetNode(xmlValue, `clinic_info`)
	//2023年11月24日16:22:56 xingqiyi
	//CSB0108RC0274001 CSB0108RC0275001
	//机构号为00183000、00183010、00183002、00183012、00183300、00183310的案件不执行该校验
	if !utils2.HasItem([]string{"00183000", "00183010", "00183002", "00183012", "00183300", "00183310", "00183301", "00183311"}, obj.Bill.Agency) {
		for _, xValues := range nodeValues {
			nodeValue := GetNodeValue(xValues, "record_counts")
			if nodeValue != "0" {
				nItems := []string{"seq", "hospital_code", "first_date", "end_date", "accident_status", "accident_id", "accident_name_desc", "sum_amount", "calc_amount", "deduct1", "deduct2", "deduct3", "west_medicine", "china_medicine", "herbal_medicine", "examination", "inspection", "laboratory", "special_inspection", "treatment", "surgery", "other", "material", "name"}
				xValues = RegReplace(xValues, `<fee_details>[\s\S]*?<\/fee_details>`, "")
				xValues = RegReplace(xValues, `<surgery_info>[\s\S]*?<\/surgery_info>`, "")
				wrongNote += CheckNotNill(xValues, "门诊发票节点值不能为空，请核查；", nItems)
			}

		}

		nodeValues = GetNode(xmlValue, `inpatient_info`)
		for _, xValues := range nodeValues {
			nodeValue := GetNodeValue(xValues, "record_counts")
			if nodeValue != "0" {
				nItems := []string{"seq", "hospital_code", "operation_type", "hospital_num", "in_date", "out_date", "bed_amt_days", "diagnostic_amt_days", "accident_status", "accident_id", "accident_name_desc", "doctor_name", "doctor_name2", "doctor_name3", "doctor_name4", "doctor_name5", "name", "sum_amount", "calc_amount", "deduct1", "deduct2", "deduct3", "west_medicine", "china_medicine", "herbal_medicine", "examination", "inspection", "laboratory", "special_inspection", "treatment", "surgery", "nursing", "blood_transfusion", "bed", "other", "material", "supplementary", "refund", "self_cash_payment", "account_payment", "prepay", "fee_count"}
				xValues = RegReplace(xValues, `<fee_details>[\s\S]*?<\/fee_details>`, "")
				xValues = RegReplace(xValues, `<surgery_info>[\s\S]*?<\/surgery_info>`, "")
				wrongNote += CheckNotNill(xValues, "住院发票节点值不能为空，请核查；", nItems)
			}

		}
	}

	nodeValues = GetNode(xmlValue, `surgery_info`)
	for _, xValues := range nodeValues {
		nodeValue := GetNodeValue(xValues, "record_counts")
		if nodeValue != "0" {
			nItems := []string{"seq", "surgery_date", "surgery_code"}
			wrongNote += CheckNotNill(xValues, "手术节点值不能为空，请核查；", nItems)
		}

	}

	if !RegIsMatch(obj.Bill.Agency, `^(00183000|00183010|00183002|00183012|00183300|00183310|00183301|00183311)$`) {
		nodeValues = GetNode(xmlValue, `fee_details`)
		for _, xValues := range nodeValues {
			nodeValue := GetNodeValue(xValues, "fee_count")
			if nodeValue != "0" {
				nItems := []string{"seq", "item_name", "item_amount"}
				wrongNote += CheckNotNill(xValues, "清单节点值不能为空，请核查；", nItems)
			}

		}
	}

	accident_time := GetNodeValue(xmlValue, "accident_time")
	nodeValues = RegMatchAll(xmlValue, `<(clinic_record|inpatient_record)>[\s\S]*?<\/(clinic_record|inpatient_record)>`)
	for _, xValues := range nodeValues {
		surgery_date := GetNodeValue(xValues, "surgery_date")
		clinic_no := GetNodeValue(xValues, "clinic_no")
		if clinic_no == "" {
			clinic_no = GetNodeValue(xValues, "operation_no")
		}
		if surgery_date != "" {
			a, _ := time.Parse("20060102", surgery_date)
			b, _ := time.Parse("20060102", accident_time)
			now := time.Now()
			if a.Before(b) {
				wrongNote += clinic_no + "发票手术日期不能早于事故日期；"
			}
			if a.After(now) {
				wrongNote += clinic_no + "发票手术日期不能晚于当前系统时间；"
			}
		}

		surgery := GetNodeValue(xValues, "surgery")
		record_counts := GetNodeValue(xValues, "record_counts")
		if surgery != "" {
			if surgery == "0.00" && record_counts != "0" {
				wrongNote += clinic_no + "未存在手术费，不应录入手术，请检查；"
			} else if surgery != "0.00" && record_counts == "0" {
				wrongNote += clinic_no + "存在手术费，但未录入手术，请检查；"
			}
		}

	}

	return wrongNote
}
func CheckNotNill(xmlValue, mess string, nItems []string) string {
	// nItems := []string{"myversion", "case_no", "return_flag"}
	mes := ""
	for _, node := range nItems {
		nodes := RegMatchAll(xmlValue, `<`+node+`(>.*?<\/`+node+`>| \/>)`)
		for _, nodeValue := range nodes {
			value := GetNodeValue(nodeValue, node)
			if value == "" {
				mes += node + mess
			}
		}
	}

	return mes
}

func GetNode(xmlValue, node string) []string {
	query := `<` + node + `>[\s\S]*?<\/` + node + `>`
	return RegMatchAll(xmlValue, query)
}

func GetFieldsInput(fields []model.ProjectField, code string) string {
	for _, field := range fields {
		if field.Code == code {
			return field.ResultValue
		}
	}
	return ""
}

func GetFieldsExit(fields []model.ProjectField, code string) bool {
	for _, field := range fields {
		if field.Code == code {
			return true
		}
	}
	return false
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

func GetFieldsAllFinal(fields []model.ProjectField, code string) float64 {
	fieldsValueAll := 0.0
	for _, field := range fields {
		if field.Code == code {
			fieldsValueAll = utils2.SumFloat(fieldsValueAll, ParseFloat(field.FinalValue), "+")
		}
	}
	return fieldsValueAll
}

func IsExist(fields []model.ProjectField, code string) bool {
	isExist := false
	for _, field := range fields {
		if field.Code == code {
			isExist = true
		}
	}
	return isExist
}

func IsRepeat(fields []model.ProjectField, code string) bool {
	isExist := false
	value := GetFieldsFinal(fields, code)
	nextValue := ""
	for _, field := range fields {
		if field.Code == code {
			nextValue = field.FinalValue
			if value == nextValue {
				isExist = true
			}
		}

	}
	return isExist
}
