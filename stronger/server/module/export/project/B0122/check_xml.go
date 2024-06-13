package B0122

import (
	"fmt"
	"math"
	"regexp"
	"server/module/load/model"
	"strconv"

	"server/utils"
	"strings"

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

	// err, fieldCheckConfs := service.GetProFieldCheckConf(pro)
	// fmt.Println("--------fieldCheckConfs-err------------", err)
	// if err != nil {
	// 	return err, wrongNote
	// }
	// fieldCheckConfMap := make(map[string][]proModel.SysProFieldCheck)
	// // fmt.Println("---------fieldCheckConfMap------------", len(fieldCheckConfMap))
	// for _, fieldCheckConf := range fieldCheckConfs {
	// 	fieldCheckConfMap[fieldCheckConf.Code] = fieldCheckConf.SysProFieldChecks
	// }
	// eleLen := reflect.ValueOf(obj).NumField()
	// for j := 0; j < eleLen; j++ {
	// 	if reflect.TypeOf(obj).Field(j).Name != "Bill" && reflect.TypeOf(obj).Field(j).Name != "Fields" {
	// 		//每张发票每种类型的字段
	// 		fmt.Println("---------------------------", reflect.TypeOf(obj).Field(j).Name)
	// 		yiYuanObjs := reflect.ValueOf(obj).Field(j).Interface().([]YiYuanObj)
	// 		for _, yiYuanObj := range yiYuanObjs {
	// 			yiYuaneleLen := reflect.ValueOf(yiYuanObj).NumField()
	// 			for y := 0; y < yiYuaneleLen; y++ {
	// 				if reflect.TypeOf(yiYuaneleLen).Field(y).Name != "Name" && reflect.TypeOf(yiYuaneleLen).Field(y).Name != "Fields" {
	// 					fieldsMaps := reflect.ValueOf(yiYuaneleLen).Field(y).Interface().([]FieldsMap)
	// 					for _, fieldsMap := range fieldsMaps {
	// 						if fieldsMap.Code == "" {
	// 							continue
	// 						}
	// 						for _, field := range fieldsMap.Fields {
	// 							items, isExit := fieldCheckConfMap[field.Code]
	// 							// fmt.Println("---------items------------", items)
	// 							if isExit {
	// 								for _, item := range items {
	// 									fffs := strings.Split(item.Value, ";")
	// 									// fmt.Println("---------fffs------------", fffs)
	// 									for _, fff := range fffs {
	// 										mess := "账单号:" + fieldsMap.Code + item.Mark + ";"
	// 										if strings.Index(wrongNote, mess) != -1 {
	// 											continue
	// 										}
	// 										if item.CheckType == "1" {
	// 											if field.ResultValue == fff {
	// 												wrongNote += mess
	// 											}
	// 										} else if item.CheckType == "2" {
	// 											if strings.Index(field.ResultValue, fff) != -1 {
	// 												wrongNote += mess
	// 											}
	// 										} else if item.CheckType == "3" {
	// 											if strings.Index(field.ResultValue, fff) == -1 {
	// 												wrongNote += mess
	// 											}
	// 										}
	// 									}

	// 								}
	// 							}
	// 						}
	// 					}
	// 				}
	// 			}

	// 		}

	// 	}
	// }
	// // for
	return nil, wrongNote

}

func CodeCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	// fields := obj.Fields
	wrongNote := ""
	var filedArr []string
	fc079Map := map[string]int{}
	fc037Arrs := []string{}
	fc037FArrs := []string{}
	//CSB0122RC0230001 调整新增字段：
	//fc497 fc507
	//fc498 fc508
	//fc499 fc509
	//fc500 fc510
	//fc501 fc511
	//fc502 fc512
	//fc503 fc513
	//fc504 fc514
	//fc505 fc515
	//fc506 fc516
	codes := [][]string{
		{"fc164", "fc174"},
		{"fc165", "fc175"},
		{"fc166", "fc176"},
		{"fc167", "fc177"},
		{"fc168", "fc178"},
		{"fc169", "fc179"},
		{"fc170", "fc180"},
		{"fc171", "fc181"},
		{"fc172", "fc182"},
		{"fc173", "fc183"},
		{"fc497", "fc507"},
		{"fc498", "fc508"},
		{"fc499", "fc509"},
		{"fc500", "fc510"},
		{"fc501", "fc511"},
		{"fc502", "fc512"},
		{"fc503", "fc513"},
		{"fc504", "fc514"},
		{"fc505", "fc515"},
		{"fc506", "fc516"},
	}
	for i, medical := range obj.Medical {
		allInvoice := medical.Hospital
		allInvoice = append(allInvoice, medical.Clinc...)
		allInvoice = append(allInvoice, medical.Pharmacy...)
		allInvoice = append(allInvoice, medical.Nonmedical...)
		allInvoice = append(allInvoice, medical.Injury...)
		allInvoice = append(allInvoice, medical.Operation...)
		for _, invoice := range allInvoice {
			for _, codef := range codes {
				isExist, fcaaa := GetOneField(invoice.Fields, codef[0], true)
				if isExist && !utils.RegIsMatch(`^(|W53\.951|W54\.951|W55\.951|W55\.952|W56\.851|W56\.852|W57\.951|W58\.951|W59\.951|W60\.951|W64\.951|W64\.952|W64\.953|W64\.954)$`, fcaaa) {
					isExist, fcbbb := GetOneField(invoice.Fields, codef[1], true)
					if isExist && fcbbb == "" {
						wrongNote += "【" + codef[1] + "】伤病补充说明不可为空，请检查；"
					}
				}
			}

			_, fc079 := GetOneField(invoice.Fields, "fc079", true)

			if fc079 != "" {
				_, isExist := fc079Map[fc079]
				if isExist {
					fc079Map[fc079]++
				} else {
					fc079Map[fc079] = 1
				}
			}

			_, fc135 := GetOneField(invoice.Fields, "fc135", true)
			_, fc491 := GetOneField(invoice.Fields, "fc491", true)
			if fc135 == "1" && fc491 == "" {
				wrongNote += fc079 + "发票,发票查验结果为空,请检查;"
			}
			if fc135 == "1" && (fc491 == "2" || fc491 == "3") {
				wrongNote += fc079 + "发票查验结果不通过，请确认五要素是否录入正确；"
			}

			allFields := invoice.Fields
			mes := ""
			//fc080 := GetFieldValue(invoice.Fields, "fc080", true, -1)
			sum2 := 0.0
			sum1 := 0.0
			issum1 := false
			// isCheck := true
			//fc432 := false
			for _, field := range invoice.Fields {
				if field.Code == "fc037" {
					fc037Arrs = append(fc037Arrs, field.ResultValue)
					fc037FArrs = append(fc037FArrs, field.FinalValue)
				}
				if RegIsMatch(field.Code, `^(fc432|fc433|fc434|fc435|fc436|fc437|fc438|fc439|fc440|fc441|fc442|fc443|fc444|fc445|fc446|fc447|fc448|fc449|fc450|fc451|fc452|fc453|fc454|fc455|fc456|fc457|fc458|fc459|fc460|fc461)$`) && field.FinalValue != "" {
					sum2 = utils.SumFloat(sum2, utils.ParseFloat(field.FinalValue), "+")
				}

				if RegIsMatch(field.Code, `^(fc492|fc493|fc494|fc495)$`) && field.FinalValue != "" {
					if field.ResultValue != "A" {
						issum1 = true
					}
					sum1 = utils.SumFloat(sum1, utils.ParseFloat(field.FinalValue), "+")
				}

				if fc135 == "1" && fc491 != "" && RegIsMatch(field.Code, `^(fc079|fc487|fc488|fc489|fc490)$`) && field.FinalValue == "" && strings.Index(wrongNote, invoice.Code+"电子发票五要素未录入齐全，请确认是否为电子发票；") == -1 {
					wrongNote += invoice.Code + "电子发票五要素未录入齐全，请确认是否为电子发票；"
				}
				//CSB0122RC0233001  2024-04-28 调整 新增字段 fc497|fc498|fc499|fc500|fc501|c502|fc503|fc504|fc505|fc506
				if RegIsMatch(field.Code, `^(fc164|fc165|fc166|fc167|fc168|fc169|fc170|fc171|fc172|fc173|fc497|fc498|fc499|fc500|fc501|c502|fc503|fc504|fc505|fc506)$`) && field.FinalValue != "" && strings.Index(wrongNote, invoice.Code+"伤病代码不在常量表中，请检查；") == -1 {
					_, isExit := constMap["jiBingDaiMaMap"][field.ResultValue]
					if !isExit {
						wrongNote += invoice.Code + "伤病代码不在常量表中，请检查；"
					}
				}
			}

			isExist, fc338 := GetOneField(invoice.Fields, "fc338", true)
			_, fc280 := GetOneField(invoice.Fields, "fc280", true)
			_, fc080 := GetOneField(invoice.Fields, "fc080", true)
			if isExist {
				aa := utils.SumFloat(utils.ParseFloat(fc338), utils.ParseFloat(fc280), "+")
				if aa >= utils.ParseFloat(fc080) {
					wrongNote += invoice.Code + "报销金额加自费金额大于发票总金额，请检查；"
				}
				cache := utils.SumFloat(utils.ParseFloat(fc338), sum1, "-")
				if issum1 && cache != 0.00 && cache != -0.00 {
					wrongNote += invoice.Code + "发票报销单自费金额与总自费金额差[" + utils.ToString(cache) + "]，请核查；"
				}
			}
			//cache := utils.SumFloat(utils.ParseFloat(fc080), sum2, "-")
			//if fc432 && cache != 0.00 && cache != -0.00 {
			//	wrongNote += "发票【" + invoice.Code + "】明细金额与发票总金额不一致，差额为[" + utils.ToString(cache) + "]，请核查；"
			//}

			//20230810 取消
			//sum := 0.0
			for _, qingDan := range invoice.QingDan {
				for _, field := range qingDan.Fields {
					if RegIsMatch(field.Code, `^(fc142|fc324|fc325|fc326|fc327|fc328|fc329|fc330)$`) && field.FinalValue != "" {
						//20230810 取消
						//sum = utils.SumFloat(sum, utils.ParseFloat(field.FinalValue), "+")
					}

				}
				allFields = append(allFields, qingDan.Fields...)
			}
			//20230810 取消
			//cache = utils.SumFloat(utils.ParseFloat(fc080), sum, "-")
			//if len(invoice.QingDan) > 0 && cache != 0.00 && cache != -0.00 {
			//	wrongNote += "发票【" + invoice.Code + "】清单项目金额与发票总金额不一致，差额为[" + utils.ToString(cache) + "]，请核查；"
			//}
			for jj, field := range allFields {
				//编码:CSB0122RC0148000
				//1.校验所有字段的录入值，当包含?或？时，出导出校验：发票【xxx】的【yyy】存在?号，请核实；
				//2.如同一发票下多个字段包含?或？，则导出校验提示格式为：发票【xxx】的【yyy1、yyy2】存在?号，请核实；（xxx为发票号fc079的值，yyy为包含问号的字段名）

				if IsExist(allFields, "fc079") {
					if RegIsMatch(field.ResultValue, "\\?|？") {
						fmt.Println("--i != len(obj.Medical)-1进入", i, len(obj.Medical))
						filedArr = append(filedArr, field.Name)
					}
				}

				if RegIsMatch(field.ResultValue, `(\?|\？)`) {
					mes += field.Name + "、"
				}

				if jj == len(allFields)-1 {
					if len(filedArr) > 0 {
						join := strings.Join(filedArr, "、")
						message := "发票" + fc079 + "的【" + join + "】存在?号，请核实;"
						if strings.Index(wrongNote, message) == -1 {
							wrongNote += message
							filedArr = make([]string, 0)
						}
					}
				}

			}
			if mes != "" {
				wrongNote += "发票【" + invoice.Code + "】的【" + mes + "】存在?号，请核实；"
			}

			// 编码 CSB0122RC0172000 时间 20230809  当fc214结果值不为空时，fc216结果值不能为空，空则出导出校验：意外发生地不能为空
			_, fc214Val := GetOneField(invoice.Fields, "fc214", true)
			if fc214Val != "" {
				_, fc216Val := GetOneField(invoice.Fields, "fc216", true)
				if fc216Val == "" {
					wrongNote += "意外发生地不能为空；"
				}
			}

			//(1). 编码 CSB0122RC0162000 时间 20230809  在同一发票属性下，
			//fc432、fc433、fc434、fc435、fc436、fc437、fc438、fc439、fc440、fc441、
			//fc442、fc443、fc444、fc445、fc446、fc447、fc448、fc449、fc450、fc451、
			//fc452、fc453、fc454、fc455、fc456、fc457、fc458、fc459、fc460、fc461
			//的录入值相加不等于fc080的值时，出导出校验：“xxx发票明细金额与总金额不一致，差额为XXX（当明细金额相加大于fc080时，需在差额前面加"-"），请确认并修改；
			//说明：
			//1、上述字段包含有问号时，不进行校验；        2、第一个xxx为同一发票下的fc079的值；
			//3、第二个XXX为fc080-（fc432+fc433+fc434+fc435+fc436+fc437+fc438+fc439+fc440+fc441+fc442+fc443+fc444+fc445+fc446+fc447+fc448+fc449+fc450+fc451+fc452+fc453+fc454+fc455+fc456+fc457+fc458+fc459+fc460+fc461）的差额（字段没有录入值时，按0计算）；

			//(2) CSB0122RC0166000 时间 20230810 在同一发票属性下，将
			//fc402、fc403、fc404、fc405、fc406、fc407、fc408、fc409、fc410、fc411、fc412、fc413、fc414、fc415、fc416、fc417、fc418、fc419、fc420、fc421、fc422、fc423、fc424、fc425、fc426、fc427、fc428、fc429、fc430、fc431
			//的录入值与所有fc161录入值作对比，存在不一致的情况时，出导出校验：xxx发票下【yyy1，yyy2，……】发票大项名称与清单大项名称存在不一致的情况，请核查；
			//（第一个xxx为同一发票fc079的值，yyy为不一致的大项内容，多个不一致时合成一条导出校验）
			sField := []string{"fc432", "fc433", "fc434", "fc435", "fc436", "fc437", "fc438", "fc439", "fc440", "fc441", "fc442", "fc443", "fc444", "fc445",
				"fc446", "fc447", "fc448", "fc449", "fc450", "fc451", "fc452", "fc453", "fc454", "fc455", "fc456", "fc457", "fc458", "fc459", "fc460", "fc461"}
			sFieldName := []string{"fc402", "fc403", "fc404", "fc405", "fc406", "fc407", "fc408", "fc409", "fc410", "fc411", "fc412", "fc413", "fc414", "fc415", "fc416", "fc417", "fc418", "fc419", "fc420", "fc421", "fc422", "fc423", "fc424", "fc425", "fc426", "fc427", "fc428", "fc429", "fc430", "fc431"}
			//sField字段存在在校验
			sFieldIsExist := false
			fc080Input := GetFieldsInput(invoice.Fields, "fc080")
			filedCount := 0.00
			var fieldResultValue []string
			for _, field := range sField {
				exist := IsExist(invoice.Fields, field)
				if exist {
					sFieldIsExist = true
				}
				fieldInput := GetFieldsInput(invoice.Fields, field)
				if RegIsMatch(fieldInput, "[\\?|\\？]") {
					continue
				}
				if fieldInput == "" {
					fieldInput = "0"
					float, _ := strconv.ParseFloat(fieldInput, 64)
					filedCount += float
				}
				sFieldFloat, _ := strconv.ParseFloat(fieldInput, 64)
				filedCount += sFieldFloat
			}
			fc080float, _ := strconv.ParseFloat(fc080Input, 64)
			CountRound := math.Round(filedCount*100) / 100
			fc080SubFiledCount := fc080float - CountRound
			if fc080Input != utils.ToString(filedCount) {
				if sFieldIsExist && fc080SubFiledCount != 0 && fc080SubFiledCount != -0 && fc080SubFiledCount != 0.00 && fc080SubFiledCount != -0.00 {
					if strings.Index(wrongNote, GetFieldsFinal(invoice.Fields, "fc079")+"发票明细金额与总金额不一致，差额为"+utils.ToString(fc080SubFiledCount)+"；") == -1 {
						wrongNote += GetFieldsFinal(invoice.Fields, "fc079") + "发票明细金额与总金额不一致，差额为" + utils.ToString(fc080SubFiledCount) + "；"
					}
				}
			}
			//CSB0122RC0166000
			for _, item := range sFieldName {
				fieldInput, _ := GetFieldsInputAndName(invoice.Fields, item)
				exist := IsExist(invoice.Fields, item)
				if exist && fieldInput != "" {
					fieldResultValue = append(fieldResultValue, fieldInput)
				}
			}

			var finalResul []string
			//是否有床位费
			isBedFee := false
			//记录是否有fc161
			isExistFc161 := false
			for _, QingDan := range invoice.QingDan {
				for _, field := range QingDan.Fields {
					//（2）CSB0122RC0166000
					if field.Code == "fc161" {
						isExistFc161 = true
						if field.ResultValue != "" {
							finalResul = append(finalResul, field.ResultValue)
						}
					}
				}
			}
			//存储fc161录入值
			fc161ResultMap := make(map[string]bool)
			for _, item := range finalResul {
				fc161ResultMap[item] = true
				if item == "床位费" {
					isBedFee = true
				}
			}
			// 编码 CSB0122RC0165000  时间 20230810  同一发票下，fc279录入值为2时，校验所有fc161录入内容是否存在床位费，否则出导出校验：xxx住院发票没有床位费，请确认是否漏切或录入有误；（xxx为同一发票下fc079的值）
			_, fc279Input := GetOneField(invoice.Fields, "fc279", false)
			if fc279Input != "" && fc279Input == "2" && isExistFc161 {
				if !isBedFee {
					message := GetFieldsFinal(invoice.Fields, "fc079") + "住院发票没有床位费，请确认是否漏切或录入有误；"
					if strings.Index(wrongNote, message) == -1 {
						wrongNote += message
					}
				}
			}

			//存储发票大项录入值
			sFieldNameResultMap := make(map[string]bool)
			for _, item := range fieldResultValue {
				sFieldNameResultMap[item] = true
			}
			//找出不一致 存储发票大项值
			var fieldNotInvoiceItem []string
			for _, item := range fieldResultValue {
				_, fieldInputItem := fc161ResultMap[item]
				if !fieldInputItem {
					fieldNotInvoiceItem = append(fieldNotInvoiceItem, item)
				}
			}
			//找出不一致 存储清单值
			var inventoryNotItem []string
			for _, item := range finalResul {
				_, isExist := sFieldNameResultMap[item]
				if !isExist && len(fieldNotInvoiceItem) > 0 {
					inventoryNotItem = append(inventoryNotItem, item)
				}
			}
			// CSB0122RC0166000
			if len(fieldNotInvoiceItem) > 0 && isExistFc161 {
				join := strings.Join(fieldNotInvoiceItem, "、")
				message := GetFieldsFinal(invoice.Fields, "fc079") + "发票下【" + join + "】发票大项名称与清单大项名称存在不一致的情况，请核查；"
				if strings.Index(wrongNote, message) == -1 {
					wrongNote += message
					fieldNotInvoiceItem = make([]string, 0)
				}
			}
			//CSB0122RC0166000
			if len(inventoryNotItem) > 0 && isExistFc161 {
				join := strings.Join(inventoryNotItem, "、")
				message := GetFieldsFinal(invoice.Fields, "fc079") + "发票下【" + join + "】清单大项名称与发票大项名称存在不一致的情况，请核查；"
				if strings.Index(wrongNote, message) == -1 {
					wrongNote += message
					inventoryNotItem = make([]string, 0)
				}
			}
			//编码 CSB0122RC0167000 同一发票下
			//fc083、fc084、fc085、fc087、fc089、fc090、fc091、fc093、fc094、fc096、fc098、fc099、fc101、fc103、fc105、fc107、fc109、fc111、fc113、fc114、fc116、fc118、fc119、fc120、fc121、fc122
			//fc124、fc126、fc128、fc130、fc131以上字段结果值为负数时；出导出校验：xxx发票yyy字段结果值为负数，请核查；（xxx为同一发票下fc079的结果值，yyy为字段编码）
			Fields := []string{"fc083", "fc084", "fc085", "fc087", "fc089", "fc090", "fc091", "fc093", "fc094", "fc096", "fc098", "fc099", "fc101", "fc103", "fc105", "fc107", "fc109", "fc111", "fc113", "fc114", "fc116", "fc118", "fc119", "fc120", "fc121", "fc122", "fc124", "fc126", "fc128", "fc130", "fc131"}
			for _, field := range Fields {
				_, item := GetOneField(invoice.Fields, field, true)
				if item != "" {
					float, _ := strconv.ParseFloat(item, 64)
					if float < 0 {
						wrongNote += GetFieldsFinal(invoice.Fields, "fc079") + "发票" + field + "字段结果值为负数，请核查；"
					}
				}
			}

			// CSB0122RC0209000
			//在同一发票属性下，
			//fc479、fc480、fc481、fc482、fc483、fc162、fc163的录入值相加不等于fc080的值时，出导出校验：“xxx非医疗发票明细金额与总金额不一致，差额为XXX（当明细金额相加大于fc080时，需在差额前面加"-"），请确认并修改；
			//说明：
			//1、上述字段包含有问号时，不进行校验；
			//2、第一个xxx为同一发票下的fc079的值；
			//3、第二个XXX为fc080-（fc479+fc480+fc481+fc482+fc483+fc162+fc163）的差额（字段没有录入值时，按0计算）；
			balanceFields := []string{"fc479", "fc480", "fc481", "fc482", "fc483", "fc162", "fc163"}
			count := 0.00
			isExistField := false
			_, fc080ResultValue := GetOneField(invoice.Fields, "fc080", false)
			for _, field := range balanceFields {
				_, item := GetOneField(invoice.Fields, field, false)
				exist := IsExist(invoice.Fields, field)
				if RegIsMatch(item, `[?|？]`) {
					continue
				}
				if exist {
					if item == "" {
						item = "0"
					}
					isExistField = true
					count += ParseFloat(item)
				}
			}
			if fc080ResultValue == "" {
				fc080ResultValue = "0.00"
			}
			balanceCount := ParseFloat(fc080ResultValue) - count
			countToDecimal := fmt.Sprintf("%.2f", count)
			fc080ToDecimal := fmt.Sprintf("%.2f", fc080ResultValue)
			if isExistField && (countToDecimal != fc080ToDecimal) && (balanceCount != 0.00 || balanceCount != 0) {
				message := GetFieldsFinal(invoice.Fields, "fc079") + "非医疗发票明细金额与总金额不一致，差额为" + strconv.FormatFloat(float64(balanceCount), 'f', 2, 64) + "，请确认并修改；"
				if strings.Index(wrongNote, message) == -1 {
					wrongNote += message
				}
			}
			//	//CSB0122RC0235000 同一发票下，当fc280、fc130两个字段同时存在时，校验两个字段的录入内容是否一致，不一致时出导出校验，【XX】发票和报销单统筹金额不一致，请检查；（xx为同一发票的fc079的值）
			exit280 := IsExist(invoice.Fields, "fc280")
			exit130 := IsExist(invoice.Fields, "fc130")
			if exit280 && exit130 {
				_, fc280ResultValue := GetOneField(invoice.Fields, "fc280", false)
				_, fc130ResultValue := GetOneField(invoice.Fields, "fc130", false)
				if fc280ResultValue != fc130ResultValue {
					message := GetFieldsFinal(invoice.Fields, "fc079") + "发票和报销单统筹金额不一致，请检查；"
					if strings.Index(wrongNote, message) == -1 {
						wrongNote += message
					}
				}
			}
			//CSB0122RC0238000 同一发票下，当fc279录入值为2时，fc081、fc082两个字段的结果值不能为空，存在为空情况时出导出校验：【XX】发票出入院日期为空，请检查；
			_, fc279Result := GetOneField(invoice.Fields, "fc279", false)
			if fc279Result == "2" {
				_, fc081Result := GetOneField(invoice.Fields, "fc081", false)
				_, fc082Result := GetOneField(invoice.Fields, "fc082", false)
				if fc081Result == "" || fc082Result != "" {
					message := GetFieldsFinal(invoice.Fields, "fc079") + "发票出入院日期为空，请检查；"
					if strings.Index(wrongNote, message) == -1 {
						wrongNote += message
					}
				}

			}

		}
	}

	for key, value := range fc079Map {
		if value > 1 {
			wrongNote += "【" + key + "】发票重复" + strconv.Itoa(value) + "次，请检查；"
		}
	}
	for _, field := range obj.Fields {

		if RegIsMatch(field.Code, "^(fc474|fc475)$") {
			dizhiQu := constMap["dizhiQu"]
			_, isExist := dizhiQu[field.FinalValue]
			if !isExist && strings.Index(wrongNote, "出险地点省市区不在地址库中，请检查；") == -1 {
				wrongNote += "出险地点省市区不在地址库中，请检查；"
			}
		}

		if RegIsMatch(field.Code, "^(fc040)$") && len(fc037Arrs) > 0 && !utils.HasItem(fc037Arrs, field.ResultValue) && strings.Index(wrongNote, "诊断书对应医院名称不在发票中，请检查；") == -1 {
			wrongNote += "诊断书对应医院名称不在发票中，请检查；"
		}

		if RegIsMatch(field.Code, "^(fc040)$") && len(fc037FArrs) > 0 && strings.Index(field.FinalValue, "_") != -1 && !utils.HasItem(fc037FArrs, field.FinalValue) && strings.Index(wrongNote, "住院发票出入院日期与诊断出入院日期无法对应，请检查；") == -1 {
			wrongNote += "住院发票出入院日期与诊断出入院日期无法对应，请检查；"
		}

		// 编码 CSB0122RC0164000 时间 20230809 需求：fc003、fc004、fc013、fc014、fc027、fc028、fc054、fc055、fc056、fc081、fc082、fc189、fc190、fc191、fc192、fc193、fc197、fc198、fc210、
		//fc212、fc213、fc217、fc225、fc226、fc235、fc236、fc245、fc246、fc255、fc256、fc278、fc202、fc203、fc019、fc029、fc227、fc237、fc247、fc257
		//结果值不为YYYYMMDD格式，则出导出校验：xxx字段日期录入格式错误，请检查；（xxx为字段编码，录入内容为A、空或包含?时不校验）
		if RegIsMatch(field.Code, "^(fc003|fc004|fc013|fc014|fc027|fc028|fc054|fc055|fc056|fc081|fc082|fc189|fc190|fc191|fc192|fc193|fc197|fc198|fc210|fc212|fc213|fc217|fc225|fc226|fc235|fc236|fc245|fc246|fc255|fc256|fc278|fc202|fc203|fc019|fc029|fc227|fc237|fc247|fc257)$") {
			if RegIsMatch(field.ResultValue, `[?|？|A|]`) || field.ResultValue == "" {
				continue
			}
			matched, _ := regexp.MatchString(`^[0-9]{8}$`, field.FinalValue)
			if !matched {
				if strings.Index(wrongNote, field.Code+"字段日期录入格式错误，请检查；") == -1 {
					wrongNote += field.Code + "字段日期录入格式错误，请检查；"
				}
			}
		}

		//编码 CSB0122RC0180000
		//当左边字段录入值与对应右边字段录入值不一致时（字段唯一），出一条导出校验：开户银行或银行账号不一致，请检查；
		//账号：fc021  fc264
		//银行：fc022  fc263
		sFileds := [][]string{
			{"fc021", "fc264"},
			{"fc022", "fc263"},
		}
		for _, item := range sFileds {
			_, Oneitem := GetOneField(obj.Fields, item[1], false)
			if field.Code == item[0] {
				if (Oneitem != "" && field.ResultValue != "") && (Oneitem != field.ResultValue) {
					if strings.Index(wrongNote, "开户银行或银行账号不一致，请检查；") == -1 {
						wrongNote += "开户银行或银行账号不一致，请检查；"
					}
				}
			}
		}
		//CSB0122RC0201000 当单证类型为“重大疾病”，fc211结果值为空时，出导出校验：该案件为重大疾病案件，重疾原因不能为空，请参考病历、病理报告资料进行补充
		if field.Code == "fc211" && obj.Bill.InsuranceType == "重大疾病" {
			if field.FinalValue == "" {
				if strings.Index(wrongNote, "该案件为重大疾病案件，重疾原因不能为空，请参考病历、病理报告资料进行补充；") == -1 {
					wrongNote += "该案件为重大疾病案件，重疾原因不能为空，请参考病历、病理报告资料进行补充；"
				}
			}
		}
		//CSB0122RC0236000 当fc017结果值为空时，出导出校验：受款人姓名不可为空，请修改；
		_, fc017Val := GetOneField(obj.Fields, "fc017", true)
		if fc017Val == "" {
			if strings.Index(wrongNote, "受款人姓名不可为空，请修改；") == -1 {
				wrongNote += "受款人姓名不可为空，请修改；"
			}
		}
		//CSB0122RC0237000 校验每个fc042、fc043的结果值，与fc001作对比，存在不一致时出导出校验：事故人姓名与发票姓名不一致，请检查；
		_, fc001val := GetOneField(obj.Fields, "fc001", true)
		if RegIsMatch(field.Code, "^(fc042|fc043)$") {
			if fc001val != field.FinalValue {
				if strings.Index(wrongNote, "事故人姓名与发票姓名不一致，请检查；") == -1 {
					wrongNote += "事故人姓名与发票姓名不一致，请检查；"
				}
			}
		}

	}

	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""
	//编码 CSB0122RC0149000 时间 20230809 同一个bill_record下，当存在total_amount节点时，将所有total_amount节点值进行合计，合计后与bill_amount节点值进行对比，
	//不一致时出一条导出校验：发票[xxx]清单项目金额与发票总金额不一致，差额为[XXX]，请核查； （第一个xxx为bill_no节点值，第二个xxx为差额，差额为0.00或-0.00时无需出导出校验）
	// 编码 CSB0122RC0151000 时间 20230809 同一个bill_record下，当存在total_price节点时，将所有total_price节点值进行合计，合计后与【mealsfee节点值和ambulance节点值 drugfee_return、carfee、lodgingfee、body_return、other 节点值相加后得到的值】进行对比，
	//不一致时出一条导出校验：发票[xxx]清单项目金额与发票总金额不一致，差额为[XXX]，请核查；（第一个xxx为bill_no节点值，第二个xxx为差额，差额为0.00或-0.00时无需出导出校验）
	billRecords := RegMatchAll(xmlValue, `<bill_records>[\s\S]*?<\/bill_records>`)
	for _, record := range billRecords {
		totalPricesCount := 0.00
		totalAmountCount := 0.00
		mealsfeAdd := 0.00
		bill_amounts := RegMatchAll(record, `<bill_amount>[\s\S]*?<\/bill_amount>`)
		totalAmounts := RegMatchAll(record, `<total_amount>[\s\S]*?<\/total_amount>`)
		totalPrices := RegMatchAll(record, `<total_price>[\s\S]*?<\/total_price>`)
		mealsfees := RegMatchAll(record, `<mealsfee>[\s\S]*?<\/mealsfee>`)
		ambulances := RegMatchAll(record, `<ambulance>[\s\S]*?<\/ambulance>`)
		drugfeeReturn := RegMatchAll(record, `<drugfee_return>[\s\S]*?<\/drugfee_return>`)
		carfee := RegMatchAll(record, `<carfee>[\s\S]*?<\/carfee>`)
		lodgingfee := RegMatchAll(record, `<lodgingfee>[\s\S]*?<\/lodgingfee>`)
		bodyReturn := RegMatchAll(record, `<body_return>[\s\S]*?<\/body_return>`)
		other := RegMatchAll(record, `<other>[\s\S]*?<\/other>`)

		if len(totalAmounts) > 0 {
			for _, totalAmount := range totalAmounts {
				totalAmountVal := GetNodeValue(totalAmount, "total_amount")
				float, _ := strconv.ParseFloat(totalAmountVal, 64)
				totalAmountCount = totalAmountCount + float
			}
		}
		if len(totalPrices) > 0 {
			for _, totalPrice := range totalPrices {
				totalPriceVal := GetNodeValue(totalPrice, "total_price")
				float, _ := strconv.ParseFloat(totalPriceVal, 64)
				totalPricesCount = totalPricesCount + float
			}
		}
		if len(mealsfees) > 0 && len(ambulances) > 0 {
			mealsfeeVal := GetNodeValue(mealsfees[0], "mealsfee")
			ambulanceVal := GetNodeValue(ambulances[0], "ambulance")
			drugFeeReturnVal := GetNodeValue(drugfeeReturn[0], "drugfee_return")
			carFeeVal := GetNodeValue(carfee[0], "carfee")
			lodgingFeeVal := GetNodeValue(lodgingfee[0], "lodgingfee")
			bodyReturnVal := GetNodeValue(bodyReturn[0], "body_return")
			otherVal := GetNodeValue(other[0], "other")
			add, _ := ConvertAdd(mealsfeeVal, ambulanceVal)
			addTo, _ := ConvertAdd(drugFeeReturnVal, carFeeVal)
			number := add + addTo
			addThree, _ := ConvertAdd(lodgingFeeVal, bodyReturnVal)
			numberTow := number + addThree
			otherFloat, _ := strconv.ParseFloat(otherVal, 64)
			val := numberTow + otherFloat
			mealsfeAdd = val
		}
		billNo := GetNodeValue(record, "bill_no")
		chae := 0.00
		if len(bill_amounts) > 0 && len(totalAmounts) > 0 {
			billAmount := GetNodeValue(bill_amounts[0], "bill_amount")
			totalAmountCountFloat := ParseFloat(fmt.Sprintf("%.2f", totalAmountCount))
			billAmountFloat, _ := strconv.ParseFloat(billAmount, 64)
			chae = billAmountFloat - totalAmountCountFloat
			if billAmount != strconv.FormatFloat(float64(totalAmountCountFloat), 'f', 2, 64) {
				if chae != 0.00 && chae != -0.00 && chae != 0 && chae != -0 {
					if strings.Index(wrongNote, "发票["+billNo+"]清单项目金额与发票总金额不一致，差额为["+strconv.FormatFloat(float64(chae), 'f', 2, 64)+"]，请核查；") == -1 {
						wrongNote += "发票[" + billNo + "]清单项目金额与发票总金额不一致，差额为[" + strconv.FormatFloat(float64(chae), 'f', 2, 64) + "]，请核查；"
					}
				}
			}
		}
		chaePrices := totalPricesCount - mealsfeAdd
		if totalPricesCount != chaePrices && len(totalPrices) > 0 {
			if chaePrices != 0.00 || chaePrices != -0.00 {
				if strings.Index(wrongNote, "发票["+billNo+"]清单项目金额与发票总金额不一致，差额为["+strconv.FormatFloat(float64(chaePrices), 'f', 2, 64)+"]，请核查；") == -1 {
					wrongNote += "发票[" + billNo + "]清单项目金额与发票总金额不一致，差额为[" + strconv.FormatFloat(float64(chaePrices), 'f', 2, 64) + "]，请核查；"
				}
			}
		}

	}

	insurance_relate := GetNodeValue(xmlValue, "insurance_relate")
	accident_name := GetNodeValue(xmlValue, "accident_name")
	payee := GetNodeValue(xmlValue, "payee")
	apply_name := GetNodeValue(xmlValue, "apply_name")
	if insurance_relate == "I" {
		if accident_name != payee {
			wrongNote += "关系为本人，受款人和事故人姓名不一致，请检查；"
		}
	} else {
		if accident_name == apply_name {
			wrongNote += "关系不为本人，申请人姓名和事故人姓名一致，请检查；"
		}
	}

	return wrongNote
}

func ConvertAdd(strOne, strTo string) (float64, error) {
	float, err := strconv.ParseFloat(strOne, 64)
	floatTo, err := strconv.ParseFloat(strTo, 64)
	countFloat := float + floatTo
	return countFloat, err
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
func GetFieldsInputAndName(fields []model.ProjectField, code string) (str, name string) {
	for _, field := range fields {
		if field.Code == code {
			return field.ResultValue, field.Name
		}
	}
	return "", ""
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
func ParseFloat(value string) float64 {
	val, _ := strconv.ParseFloat(value, 64)
	return val
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
func GetNodeValue(xmlValue, nodeName string) string {
	beginNode := strings.Index(xmlValue, "<"+nodeName+">") + len(nodeName) + 2
	endNode := strings.Index(xmlValue, "</"+nodeName+">")
	sValue := ""
	if beginNode != -1 && endNode != -1 {
		sValue = Substr(xmlValue, beginNode, endNode)
	}
	return sValue
}
func Substr(str string, start, end int) string {
	if start == -1 || end == -1 {
		return ""
	}
	return string(str[start:end])
}
