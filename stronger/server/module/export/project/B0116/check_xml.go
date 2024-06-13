package B0116

import (
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	"server/module/export/service"
	"server/module/load/model"
	proModel "server/module/pro_conf/model"
	"server/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	_ "github.com/wxnacy/wgo/arrays"
	_ "go.uber.org/zap"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	global.GLog.Info("B0116:::CheckXml")
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
	wrongNote := ""
	//意外出险细节
	accidental := constMap["accidental"]
	//损伤外部原因
	damages := constMap["damages"]
	//所有字段
	fields := obj.Fields
	// 存储fc045
	var billCodes []string
	for _, fieldsMap := range obj.MedicalBillInfo {
		var sumNumber float64
		saveField := []string{}
		for _, field := range fieldsMap.Fields {
			//CSB0116RC0119000
			//符合以下条件时，出导出校验：XXX账单类型错误，请确认（XXX为该发票中fc045的结果值）
			//1.当同一发票中fc201录入值为“1”，且该发票的fc111录入值不为“1、3、6、8、10”中任意一个时；
			//2.当同一发票中fc201录入值为“2”，且该发票的fc111录入值不为“2、4、5、7、9、11”中任意一个时；
			//（会存在多组fc201和fc111的情况，根据字段fc088判断是否为同一发票）
			//1.当同一发票中fc201录入值为“1”，且该发票的fc111录入值不为“1、3、6、8、10”中任意一个时；
			_, fc045Val := GetOneField(fieldsMap.Fields, "fc045", false)
			_, fc111ResultVal := GetOneField(fieldsMap.Fields, "fc111", false)
			if RegIsMatch(field.Code, `^(fc201)$`) && field.ResultValue == "1" && utils.RegIsMatch(`^(1|3|6|8|10)$`, fc111ResultVal) {
				if strings.Index(wrongNote, fc045Val+"账单类型错误，请确认；") == -1 {
					wrongNote += fc045Val + "账单类型错误，请确认；"
				}
			}
			//2.当同一发票中fc201录入值为“2”，且该发票的fc111录入值不为“2、4、5、7、9、11”中任意一个时；
			if RegIsMatch(field.Code, `^(fc201)$`) && field.ResultValue == "2" && utils.RegIsMatch(`^(2|4|5|7|9|11)$`, fc111ResultVal) {
				if strings.Index(wrongNote, fc045Val+"账单类型错误，请确认；") == -1 {
					wrongNote += fc045Val + "账单类型错误，请确认；"
				}
			}

			_, fc008Val := GetOneField(fields, "fc008", true)
			_, fc009Val := GetOneField(fields, "fc009", true)
			//CSB0116RC0123000 1.当fc008和fc009结果值均为空时，出提示：证件有效期和证件是否长期有效至少录一个；当单证类型为“1”时，不进行该校验
			if fc008Val == "" && fc009Val == "" && obj.Bill.InsuranceType != "1" {
				if strings.Index(wrongNote, "证件有效期和证件是否长期有效至少录一个；") == -1 {
					wrongNote += "证件有效期和证件是否长期有效至少录一个；"
				}
			}
			//CSB0116RC0123000 2.当fc008和fc009结果值同时不为空时，出提示：证件有效期和证件是否长期有效不能同时录入； 当单证类型为“1”时，不进行该校验
			if fc008Val != "" && fc009Val != "" && obj.Bill.InsuranceType != "1" {
				if strings.Index(wrongNote, "证件有效期和证件是否长期有效不能同时录入；") == -1 {
					wrongNote += "证件有效期和证件是否长期有效不能同时录入；"
				}
			}
			//CSB0116RC0132000 fc045的录入值重复时，出导出校验提示：发票号：XXX重复，请检查；（XXX为重复的fc045）
			if field.Code == "fc045" {
				if field.ResultValue != "" && (field.ResultValue != "A" || !RegIsMatch(field.ResultValue, "\\?|？")) {
					billCodes = append(billCodes, field.ResultValue)
				}
			}
			// CSB0116RC0134000
			//同一发票下，（xxx为fc045的值）
			//1.当fc048、fc049、fc050、fc051的结果值之和大于fc046的结果值时，出导出校验提示：xxx发票其他扣减费用、新农合、医保支付金额、其他商保支付之和不得大于账单发生金额；
			//2.当fc060、fc061、fc062、fc059、fc050的结果值之和大于fc046的结果值时，出导出校验提示：xxx发票自付一、自付二、自费、医疗保险基金支付金额、其他商保支付金额之和不得大于账单发生金额；
			//6.当fc063、fc072、fc073、fc074的结果值之和大于fc059的结果值时，出导出校验提示：xxx发票门诊大额支付、退休人员补充医疗保险支付、残疾军人补助支付、公务员医疗补助支付之和，不得大于医疗保险基金支付金额；
			sFields := [][]string{
				{"fc048", "fc049", "fc050", "fc051", "发票其他扣减费用、新农合、医保支付金额、其他商保支付之和不得大于账单发生金额；"},
				{"fc060", "fc061", "fc062", "fc059", "fc050", "发票自付一、自付二、自费、医疗保险基金支付金额、其他商保支付金额之和不得大于账单发生金额；"},
				{"fc063", "fc072", "fc073", "fc074", "发票门诊大额支付、退休人员补充医疗保险支付、残疾军人补助支付、公务员医疗补助支付之和，不得大于医疗保险基金支付金额；"},
			}
			_, fc046Val := GetOneField(fieldsMap.Fields, "fc046", true)
			//var fc046Int int
			for i := 0; i < len(sFields); i++ {
				codeNumber := 0.00
				for ii, item := range sFields[i] {
					if field.Code == item && field.FinalValue != "" && ii < len(sFields[i])-1 {
						atoi, _ := strconv.ParseFloat(field.FinalValue, 64)
						codeNumber = codeNumber + atoi
					}
				}
				if fc046Val != "" {
					atoi, _ := strconv.ParseFloat(fc046Val, 64)
					if codeNumber > atoi {
						if strings.Index(wrongNote, fc045Val+sFields[i][len(sFields[i])-1]) == -1 {
							wrongNote += fc045Val + sFields[i][len(sFields[i])-1]
						}
					}
				}

			}
			//3.当fc058的结果值大于fc046的结果值时，出导出校验提示：xxx发票医疗保险范围内金额不能大于账单发生金额；
			//4.当fc077的结果值大于fc046的结果值时，出导出校验提示：xxx发票筹基金本次支付不能大于账单发生金额；
			//5.当fc080的结果值大于fc046的结果值时，出导出校验提示：xxx发票大额互助资金本次支付不能大于账单发生金额；
			//7.当fc077的结果值大于fc082的结果值时，出导出校验提示：xxx发票统筹基金本次支付不能大于统筹基金年度内累计支付；
			//8.当fc080的结果值大于fc083的结果值时，出导出校验提示：xxx发票大额互助资金本次支付不能大于大额互助资金累计支付金额；
			// CSB0116RC0140000 同一发票下，当fc051的结果值大于fc046时，出导出校验提示：xxx发票其他扣减费用大于账单金额，请检查；（xxx为fc045的值）
			sField := [][]string{
				{"fc058", "fc046", "发票医疗保险范围内金额不能大于账单发生金额；"},
				{"fc077", "fc046", "发票筹基金本次支付不能大于账单发生金额；"},
				{"fc080", "fc046", "发票大额互助资金本次支付不能大于账单发生金额；"},
				{"fc077", "fc082", "发票统筹基金本次支付不能大于统筹基金年度内累计支付；"},
				{"fc080", "fc083", "发票大额互助资金本次支付不能大于大额互助资金累计支付金额；"},
				{"fc051", "fc046", "发票其他扣减费用大于账单金额，请检查；"},
			}
			for i := 0; i < len(sField); i++ {
				for ii, item := range sField[i] {
					if ii == 0 {
						if field.Code == item {
							_, vals := GetOneField(fieldsMap.Fields, sField[i][1], true)
							if vals != "" {
								atoi := ParseFloat(vals)
								fiedNumber := ParseFloat(field.FinalValue)
								if fiedNumber > atoi {
									if strings.Index(wrongNote, fc045Val+sField[i][2]) == -1 {
										wrongNote += fc045Val + sField[i][2]
									}
								}
							}
						}
					}
				}
			}
			//CSB0116RC0135000
			//同一发票下，当fc155、fc156、fc157、fc158、fc159、fc160、fc161、fc162、fc163、fc164、fc165、fc166、fc167、fc168、fc169、fc170、fc171、fc172、fc173、fc174的结果值之和大于fc046的结果值时
			//出导出校验提示：xxx发票所有扣除总金额之和不能大于账单发生金额；（xxx为fc045的值）
			fieldCode := []string{"fc155", "fc156", "fc157", "fc158", "fc159", "fc160", "fc161", "fc162", "fc163", "fc164", "fc165", "fc166", "fc167", "fc168", "fc169", "fc170", "fc171", "fc172", "fc173", "fc174"}
			var number float64
			for _, code := range fieldCode {
				if field.Code == code {
					if field.FinalValue != "" {
						fiedNumber := ParseFloat(field.FinalValue)
						number += fiedNumber
					}
				}
			}
			if fc046Val != "" {
				atoi := ParseFloat(fc046Val)
				if number > atoi {
					if strings.Index(wrongNote, fc045Val+"发票所有扣除总金额之和不能大于账单发生金额；") == -1 {
						wrongNote += fc045Val + "发票所有扣除总金额之和不能大于账单发生金额；"
					}
				}
			}

			//CSB0116RC0136000 同一发票下，当fc111录入值为5或8或10或11时，对应的报销单中的fc050、fc060、fc061、fc062、fc077、fc080的结果值之和不等于fc046的结果值时
			//出导出校验提示：xxx发票自付一、自付二、自费、大额互助基金本次支付、统筹基金本次支付、其他商保支付金额之和必须等于账单金额；（xxx为fc045的值）
			fieldCode = []string{"fc050", "fc060", "fc061", "fc062", "fc077", "fc080"}
			var sum int
			for _, code := range fieldCode {
				if field.Code == code {

					if field.FinalValue != "" {
						fiedNumber, _ := strconv.Atoi(field.FinalValue)
						sum += fiedNumber
					}
				}
			}
			if fc111ResultVal != "" && utils.RegIsMatch(`^(5|8|10|11)$`, fc111ResultVal) {
				fiedNumber, _ := strconv.Atoi(fc111ResultVal)
				if sum != fiedNumber {
					if strings.Index(wrongNote, fc045Val+"发票自付一、自付二、自费、大额互助基金本次支付、统筹基金本次支付、其他商保支付金额之和必须等于账单金额；") == -1 {
						wrongNote += fc045Val + "发票自付一、自付二、自费、大额互助基金本次支付、统筹基金本次支付、其他商保支付金额之和必须等于账单金额；"
					}
				}
			}

			//CSB0116RC0137000
			//同一发票下，当fc048和fc049和fc050的结果值之和为0.00时（结果值为“?”时当作0.00计算），出导出校验提示：发票号XXX报销金额为0,请检查；（xxx为fc045的结果值）
			codes := []string{"fc048", "fc049", "fc050"}
			for _, code := range codes {
				if RegIsMatch(field.Code, code) && IsExist(fieldsMap.Fields, code) {
					if field.FinalValue != "" {
						if field.FinalValue == "?" || field.FinalValue == "？" {
							field.FinalValue = "0.00"
						}
						fiedNumber := ParseFloat(field.FinalValue)
						sumNumber = sumNumber + fiedNumber
					}
				}
			}
			sprintf := fmt.Sprintf("%.2f", sumNumber)
			if (IsExist(fieldsMap.Fields, "fc048") || IsExist(fieldsMap.Fields, "fc049") || IsExist(fieldsMap.Fields, "fc050")) && RegIsMatch(`^(0|0.00)$`, sprintf) {
				if strings.Index(wrongNote, "发票号"+fc045Val+"报销金额为0,请检查；") == -1 {
					wrongNote += "发票号" + fc045Val + "报销金额为0,请检查；"
				}
			}
			// CSB0116RC0138000 当同一张发票的fc046的结果值等于该发票的“fc048和fc049和fc050的结果值之和”时（结果值为“?”时当作0.00计算
			//出导出校验提示：发票号XXX报销金额与总金额一致,请检查；（xxx为fc045的结果值）
			var atois float64
			if fc046Val != "" {
				atois = ParseFloat(fc046Val)
			}
			fc046float := fmt.Sprintf("%.2f", atois)
			if fc046float == sprintf {

				if strings.Index(wrongNote, "发票号"+fc045Val+"报销金额与总金额一致,请检查；") == -1 {
					wrongNote += "发票号" + fc045Val + "报销金额与总金额一致,请检查；"
				}
			}
			//CSB0116RC0139000 同一发票下，当fc175、fc176、fc177、fc178、fc179、fc180、fc181、fc182录入值包含“?”时
			//出导出校验提示：xxxx票据yyyy异常，请修改；（xxx为fc045的值，yyyy为对应的fc175、fc176、fc177、fc178、fc179、fc180、fc181、fc182字段名）
			if utils.RegIsMatch(`^(fc175|fc176|fc177|fc178|fc179|fc180|fc181|fc182)$`, field.Code) {
				if field.ResultValue != "" {
					matched, _ := regexp.MatchString(`[\\?|\\？]`, field.ResultValue)
					if matched {
						fmt.Println("出校验字段 ======", field.Code)
						saveField = append(saveField, field.Code)
					}
				}
			}
			join := strings.Join(saveField, "、")
			if strings.Index(wrongNote, fc045Val+"票据"+join+"异常，请修改；") == -1 {
				wrongNote += fc045Val + "票据" + join + "异常，请修改；"
			}

			//CSB0116RC0141000 同一发票下，当fc054的结果值早于fc017的结果值时，出导出校验提示：xxx发票门诊发生日期不能早于出险日期，请检查；（xxx为fc045的值）
			//CSB0116RC0142000 同一发票下，当fc055的结果值早于fc017的结果值时，出导出校验提示：xxx发票入院日期不能早于出险日期，请检查；（xxx为fc045的值）
			fieldDates := [][]string{
				{"fc054", "fc017", "发票门诊发生日期不能早于出险日期，请检查；"},
				{"fc055", "fc017", "发票入院日期不能早于出险日期，请检查；"},
			}
			for i := 0; i < len(fieldDates); i++ {
				for ii, fieldDate := range fieldDates[i] {
					if ii == 0 {
						_, vals := GetOneField(fieldsMap.Fields, fieldDates[i][1], true)
						if field.Code == fieldDate && (field.FinalValue != "" && vals != "") {
							parseField, _ := time.Parse("2006-01-02", field.FinalValue)
							parseFc017, _ := time.Parse("2006-01-02", fieldDates[i][1])
							if parseField.Before(parseFc017) {
								if strings.Index(wrongNote, fc045Val+fieldDates[i][2]) == -1 {
									wrongNote += fc045Val + fieldDates[i][2]
								}
							}
						}
					}
				}
			}

			//CSB0116RC0143000 同一发票下，当fc055的结果值晚于fc056的结果值时，出导出校验提示：xxx发票入院日期不能晚于出院日期，请检查；（xxx为fc045的值）
			_, fc055Val := GetOneField(fieldsMap.Fields, "fc055", true)
			_, fc056Val := GetOneField(fieldsMap.Fields, "fc056", true)
			if fc055Val != "" && fc056Val != "" {
				parsefc055Val, _ := time.Parse("2006-01-02", fc055Val)
				parsefc056Val, _ := time.Parse("2006-01-02", fc056Val)
				if parsefc055Val.After(parsefc056Val) {
					if strings.Index(wrongNote, fc045Val+"发票入院日期不能晚于出院日期，请检查；") == -1 {
						wrongNote += fc045Val + "发票入院日期不能晚于出院日期，请检查；"
					}
				}
			}

		}
	}

	var fc003Vals []string
	var fc003One string
	var fc185Vals []string
	var fc185One string
	var fc027s []string
	var fc109s []string
	var fc110s []string
	isFool := false
	for _, field := range fields {
		//CSB0116RC0118000
		//1.当fc013录入值为“1或2”且fc086的录入值不存在“3”时（fc086为循环字段），导出校验出提示：出险类型为医疗但缺少账单信息，请检查；
		//2.当fc013录入值为“9”且“不存在fc019或fc019的结果值为空时”，导出校验出提示：出险类型为身故但缺少身故日期，请检查；
		//3.当fc013录入值为“3或4”且fc086的录入值不存在“12”时（fc086为循环字段），导出校验出提示：出险类型为伤残/全残但缺少相应信息，请检查；
		//4.当fc013录入值为“5”且fc086的录入值不存在“14”时（fc086为循环字段），导出校验出提示：出险类型为失能但缺少相应信息，请检查；
		//5.当fc013录入值为“6或7”且fc086的录入值不存在“11”时（fc086为循环字段），导出校验出提示：出险类型为重疾/轻症但缺少相应信息，请检查；
		//6.当fc013录入值为“8”且fc086的录入值不存在“13”时（fc086为循环字段），导出校验出提示：出险类型为特种病但缺少相应信息，请检查；
		if field.Code == "fc013" {
			//1.当fc013录入值为“1或2”且fc086的录入值不存在“3”时（fc086为循环字段），导出校验出提示：出险类型为医疗但缺少账单信息，请检查；
			_, fc086ResultVal := GetOneField(fields, "fc086", false)
			fmt.Println("fc086ResultVal===========", fc086ResultVal)
			fmt.Println("==========field.ResultValue==", field.ResultValue)
			if utils.RegIsMatch(`^(1|2)$`, field.ResultValue) && (fc086ResultVal != "" && !utils.RegIsMatch(`^(3)$`, fc086ResultVal)) {
				if strings.Index(wrongNote, "出险类型为医疗但缺少账单信息，请检查；") == -1 {
					wrongNote += "出险类型为医疗但缺少账单信息，请检查；"
				}
			}
			//2.当fc013录入值为“9”且“不存在fc019或fc019的结果值为空时”，导出校验出提示：出险类型为身故但缺少身故日期，请检查；
			_, fc019Val := GetOneField(fields, "fc019", true)
			if utils.RegIsMatch(`^(9)$`, field.ResultValue) && (!IsExist(fields, "fc019") || fc019Val == "") {
				if strings.Index(wrongNote, "出险类型为身故但缺少身故日期，请检查；") == -1 {
					wrongNote += "出险类型为身故但缺少身故日期，请检查；"
				}
			}
			//3.当fc013录入值为“3或4”且fc086的录入值不存在“12”时（fc086为循环字段），导出校验出提示：出险类型为伤残/全残但缺少相应信息，请检查；
			if utils.RegIsMatch(`^(3|4)$`, field.ResultValue) && (fc086ResultVal != "" && !utils.RegIsMatch(`^[12]$`, fc086ResultVal)) {
				if strings.Index(wrongNote, "出险类型为伤残/全残但缺少相应信息，请检查；") == -1 {
					wrongNote += "出险类型为伤残/全残但缺少相应信息，请检查；"
				}
			}
			//4.当fc013录入值为“5”且fc086的录入值不存在“14”时（fc086为循环字段），导出校验出提示：出险类型为失能但缺少相应信息，请检查；
			if utils.RegIsMatch(`^(5)$`, field.ResultValue) && (fc086ResultVal != "" && !utils.RegIsMatch(`^[14]$`, fc086ResultVal)) {
				if strings.Index(wrongNote, "出险类型为失能但缺少相应信息，请检查；") == -1 {
					wrongNote += "出险类型为失能但缺少相应信息，请检查；"
				}
			}
			//5.当fc013录入值为“6或7”且fc086的录入值不存在“11”时（fc086为循环字段），导出校验出提示：出险类型为重疾/轻症但缺少相应信息，请检查；
			if utils.RegIsMatch(`^(6|7)$`, field.ResultValue) && (fc086ResultVal != "" && !utils.RegIsMatch(`^[11]$`, fc086ResultVal)) {
				if strings.Index(wrongNote, "出险类型为重疾/轻症但缺少相应信息，请检查；") == -1 {
					wrongNote += "出险类型为重疾/轻症但缺少相应信息，请检查；"
				}
			}
			//6.当fc013录入值为“8”且fc086的录入值不存在“13”时（fc086为循环字段），导出校验出提示：出险类型为特种病但缺少相应信息，请检查；
			if utils.RegIsMatch(`^(6|7)$`, field.ResultValue) && (fc086ResultVal != "" && !utils.RegIsMatch(`^[13]$`, fc086ResultVal)) {
				if strings.Index(wrongNote, "出险类型为特种病但缺少相应信息，请检查；") == -1 {
					wrongNote += "出险类型为特种病但缺少相应信息，请检查；"
				}
			}

		}

		//CSB0116RC0120000 校验所有的fc003，当存在fc003的结果值不一致或全部fc003的结果值均为空时（不存在fc003时不进行校验），导出校验出提示：出险人姓名异常，请检查；
		if field.Code == "fc003" {
			fc003Vals = append(fc003Vals, field.FinalValue)
		}
		_, fc005Val := GetOneField(fields, "fc005", true)
		_, fc017Val := GetOneField(fields, "fc017", true)
		_, fc008Val := GetOneField(fields, "fc008", true)
		_, fc018Val := GetOneField(fields, "fc018", true)
		fc005Time, _ := time.Parse("2006-01-02", fc005Val)
		fc017Time, _ := time.Parse("2006-01-02", fc017Val)
		fc008Time, _ := time.Parse("2006-01-02", fc008Val)
		fc018Time, _ := time.Parse("2006-01-02", fc018Val)

		//CSB0116RC0121000 fc005结果值晚于fc017结果值时（fc005结果值为空时不进行校验），导出校验出提示：出生日期不能晚于出险日期；
		if fc005Time.After(fc017Time) {
			if strings.Index(wrongNote, "出生日期不能晚于出险日期；") == -1 {
				wrongNote += "出生日期不能晚于出险日期；"
			}
		}
		//CSB0116RC0122000 fc008结果值早于fc017结果值时（fc008结果值为空时不进行校验），导出校验出提示：证件有效期不能早于出险日期；
		if fc008Time.Before(fc017Time) {
			if strings.Index(wrongNote, "证件有效期不能早于出险日期；") == -1 {
				wrongNote += "证件有效期不能早于出险日期；"
			}
		}
		// CSB0116RC0124000 fc018结果值早于fc017结果值时（fc018结果值为空时不进行校验），导出校验出提示：认定日期不能早于出险日期；
		if fc018Time.Before(fc017Time) {
			if strings.Index(wrongNote, "认定日期不能早于出险日期；") == -1 {
				wrongNote += "认定日期不能早于出险日期；"
			}
		}
		// CSB0116RC0125000 结果数据中fc015的录入值匹配《B0116_华夏人寿团险理赔_意外出险细节》的“中文名”（第二列），当fc015的录入值不在常量库时（录入值为空时不校验），导出校验出提示：意外细节录入内容不在常量库中，请检查；
		if field.Code == "fc015" {
			if !HasKey(accidental, field.ResultValue) {
				if strings.Index(wrongNote, "意外细节录入内容不在常量库中，请检查；") == -1 {
					wrongNote += "意外细节录入内容不在常量库中，请检查；"
				}
			}
		}
		//CSB0116RC0126000 结果数据中fc016的录入值匹配《B0116_华夏人寿团险理赔_损伤外部原因》的“中文名”（第二列），当fc016的录入值不在常量库时（录入值为空时不校验），导出校验出提示：损伤外部原因录入内容不在常量库中，请检查；
		if field.Code == "fc016" {
			if !HasKey(damages, field.ResultValue) {
				if strings.Index(wrongNote, "损伤外部原因录入内容不在常量库中，请检查；") == -1 {
					wrongNote += "损伤外部原因录入内容不在常量库中，请检查；"
				}
			}
		}
		//CSB0116RC0128000
		//fc022或fc023的结果值的前两位数不等于fc021的结果值时，或fc023的结果值的前四位数不等于fc022的结果值时，或fc023和fc022和fc021的结果值均为空时，导出校验出提示：省市区不一致，请核对；
		//fc195或fc194的结果值的前两位数不等于fc196的结果值时，或fc194的结果值的前四位数不等于fc195的结果值时，或fc194和fc195和fc196的结果值均为空时，导出校验出提示：省市区不一致，请核对；
		//fc211或fc210的结果值的前两位数不等于fc212的结果值时，或fc210的结果值的前四位数不等于fc211的结果值时，或fc210和fc211和fc212的结果值均为空时，导出校验出提示：省市区不一致，请核对；
		//当单证类型为“1”时，不进行该校验
		sFields := [][]string{
			{"fc022", "fc023", "fc021", "省市区不一致，请核对；"},
			{"fc195", "fc194", "fc196", "省市区不一致，请核对；"},
			{"fc211", "fc210", "fc212", "省市区不一致，请核对；"},
		}
		for _, sField := range sFields {
			_, oneRowVal := GetOneField(fields, sField[0], true)
			_, toRowVal := GetOneField(fields, sField[1], true)
			_, thRowVal := GetOneField(fields, sField[2], true)

			if oneRowVal == "" && toRowVal == "" && thRowVal == "" && obj.Bill.InsuranceType != "1" {
				if strings.Index(wrongNote, sField[3]) == -1 {
					wrongNote += sField[3]
				}
			}
			if oneRowVal != "" && toRowVal != "" && thRowVal != "" {
				//前两位
				oneRowValAfterTo := oneRowVal[:2]
				toRowValAfterTo := toRowVal[:2]
				toRowValAfterFour := toRowVal[:4]
				if ((oneRowValAfterTo != thRowVal || toRowValAfterTo != thRowVal) || toRowValAfterFour != oneRowVal) && obj.Bill.InsuranceType != "1" {
					if strings.Index(wrongNote, sField[3]) == -1 {
						wrongNote += sField[3]
					}
				}
			}

		}

		// CSB0116RC0129000 校验所有的fc185，当存在fc185的结果值不一致或全部fc185的结果值均为空时（不存在fc185时不进行校验），导出校验出提示：领取人姓名异常，请检查；
		//当单证类型为“1”时，不进行该校验
		if field.Code == "fc185" {
			fc185Vals = append(fc185Vals, field.FinalValue)
		}
		//CSB0116RC0130000 fc044,fc096,fc099的结果值存在重复时，导出校验出提示：伤病诊断不能重复，请检查；（结果值为空时不进行校验）
		_, fc044Val := GetOneField(fields, "fc044", true)
		_, fc096Val := GetOneField(fields, "fc096", true)
		_, fc099Val := GetOneField(fields, "fc099", true)
		if fc044Val != "" && fc096Val != "" && fc099Val != "" {
			if fc044Val == fc096Val || fc044Val == fc099Val || fc096Val == fc099Val {
				if strings.Index(wrongNote, "伤病诊断不能重复，请检查；") == -1 {
					wrongNote += "伤病诊断不能重复，请检查；"
				}
			}
		}
		//CSB0116RC0131000 fc041，fc102，fc105的结果值存在重复时，导出校验出提示：手术术式不能重复，请检查；（结果值为空时不进行校验）
		_, fc041Val := GetOneField(fields, "fc041", true)
		_, fc102Val := GetOneField(fields, "fc102", true)
		_, fc105Val := GetOneField(fields, "fc105", true)
		if fc041Val != "" && fc102Val != "" && fc105Val != "" {
			if fc041Val == fc102Val || fc041Val == fc105Val || fc102Val == fc105Val {
				if strings.Index(wrongNote, "手术术式不能重复，请检查；") == -1 {
					wrongNote += "手术术式不能重复，请检查；"
				}
			}
		}

		// CSB0116RC0144000 同一案件存在多个fc027时，对所有fc027的结果值进行校验（结果值为空的不进行校验），当fc027的结果值重复时，导出校验出提示：同一案件残疾代码不可重复；
		if field.Code == "fc027" {
			if field.FinalValue != "" {
				fc027s = append(fc027s, field.FinalValue)
			}
		}
		// CSB0116RC0145000 fc032结果值早于fc017结果值时（fc032结果值为空时不进行校验），导出校验出提示：鉴定日期不能早于出险日期；
		// CSB0116RC0146000 fc036结果值早于fc017结果值时（fc036结果值为空时不进行校验），导出校验出提示：确诊日期不能早于出险日期；
		// CSB0116RC0149000 fc217结果值早于fc216结果值时（fc217结果值为空时不进行校验），导出校验出提示：失能截止日期不能早于失能开始日期；
		// CSB0116RC0150000 fc221结果值早于fc220结果值时（fc221结果值为空时不进行校验），导出校验出提示：护理截止日期不能早于护理开始日期；
		sField := [][]string{
			{"fc032", "fc017", "鉴定日期不能早于出险日期；"},
			{"fc036", "fc017", "确诊日期不能早于出险日期；"},
			{"fc217", "fc216", "失能截止日期不能早于失能开始日期；"},
			{"fc221", "fc220", "护理截止日期不能早于护理开始日期；"},
		}
		for i := 0; i < len(sField); i++ {
			for ii, item := range sField[i] {
				if ii == 0 {
					if field.Code == item {
						_, vals := GetOneField(fields, sField[i][1], true)
						if vals != "" && field.FinalValue != "" {
							parse, _ := time.Parse("2006-01-02", vals)
							itemParse, _ := time.Parse("2006-01-02", field.FinalValue)
							if itemParse.Before(parse) {
								if strings.Index(wrongNote, sField[i][2]) == -1 {
									wrongNote += sField[i][2]
								}
							}
						}
					}
				}
			}
		}
		// CSB0116RC0147000 同一案件存在多个fc109时，对所有fc109的结果值进行校验（结果值为空的不进行校验），当fc109的结果值重复时，导出校验出提示：同一案件特种病类型不可重复；
		if field.Code == "fc109" {
			if field.FinalValue != "" {
				fc109s = append(fc109s, field.FinalValue)
			}
		}

		if field.Code == "fc110" {
			if field.FinalValue != "" {
				fc110s = append(fc110s, field.FinalValue)
			}
		}

	}
	//CSB0116RC0120000 校验所有的fc003，当存在fc003的结果值不一致或全部fc003的结果值均为空时（不存在fc003时不进行校验），导出校验出提示：出险人姓名异常，请检查；
	if len(fc003Vals) > 0 {
		fc003One = fc003Vals[0]
		for i := 0; i < len(fc003Vals); i++ {
			if fc003Vals[i] != "" {
				isFool = true
			}
		}
		for i := 1; i < len(fc003Vals); i++ {
			if fc003Vals[i] != fc003One || !isFool {
				if strings.Index(wrongNote, "出险人姓名异常，请检查；") == -1 {
					wrongNote += "出险人姓名异常，请检查；"
				}
			}
		}
	}
	// CSB0116RC0129000 校验所有的fc185，当存在fc185的结果值不一致或全部fc185的结果值均为空时（不存在fc185时不进行校验），导出校验出提示：领取人姓名异常，请检查；
	//当单证类型为“1”时，不进行该校验
	notEmpty := true
	if len(fc185Vals) > 0 {
		fc185One = fc185Vals[0]
		for _, item := range fc185Vals {
			if item != "" {
				notEmpty = false
			}
		}
		for i := 1; i < len(fc185Vals); i++ {
			if (fc185One != fc185Vals[i] || notEmpty) && obj.Bill.InsuranceType != "1" {
				if strings.Index(wrongNote, "领取人姓名异常，请检查；") == -1 {
					wrongNote += "领取人姓名异常，请检查；"
				}
			}
		}
	}
	// CSB0116RC0132000 fc045的录入值重复时，出导出校验提示：发票号：XXX重复，请检查；（XXX为重复的fc045）
	if len(billCodes) > 0 {
		fc405Rest := billCodes[0]
		if len(billCodes) > 1 {
			for i := 1; i < len(billCodes); i++ {
				if fc405Rest == billCodes[i] {
					if strings.Index(wrongNote, "发票号："+billCodes[i]+"重复，请检查；") == -1 {
						wrongNote += "发票号：" + billCodes[i] + "重复，请检查；"
					}
				}
			}
		}
	}
	// CSB0116RC0133000 fc045的录入值出现连号时（如第三个fc045录入值为123，第七个fc045录入值为122；fc045录入值为A或包含?时不进行校验），出导出校验提示：发票号XXX、YYY为连号
	//注意检查是否为同一张发票，为同一张发票时需要重加载；（XXX为数值较小的发票号，YYY为数值较大的发票号，如果出现多张连号的，如111、112、113、114，则提示为“发票号111、112、113、114为连号”）
	//存储连号  12 13 14 15 16
	var serialNumber []string
	if len(billCodes) > 0 {
		sort.Strings(billCodes)
		count := 0
		for i := 0; i < len(billCodes); i++ {
			if i+1 <= len(billCodes)-1 {
				atoi := 0
				atoiNext := 0
				if billCodes[i] != "" {
					atoi, _ = strconv.Atoi(billCodes[i])
					atoiNext, _ = strconv.Atoi(billCodes[i+1])
					count = atoi + 1
				}
				if count == atoiNext {
					//练号保存在切片
					serialNumber = append(serialNumber, billCodes[i], billCodes[i+1])
				}
			}
		}
	}
	seenMap := make(map[string]bool)
	storageSlice := []string{}
	if len(serialNumber) > 0 {
		for _, item := range serialNumber {
			if ok := seenMap[item]; !ok {
				seenMap[item] = true
				storageSlice = append(storageSlice, item)
			}
		}
		join := strings.Join(storageSlice, "、")
		if strings.Index(wrongNote, "发票号"+join+"为连号；") == -1 {
			wrongNote += "发票号" + join + "为连号；"
		}
	}

	// CSB0116RC0144000 同一案件存在多个fc027时，对所有fc027的结果值进行校验（结果值为空的不进行校验），当fc027的结果值重复时，导出校验出提示：同一案件残疾代码不可重复；
	if len(fc027s) > 0 {
		fc027One := fc027s[0]
		if len(fc027s) > 1 {
			for i := 1; i < len(fc027s); i++ {
				if fc027One == fc027s[i] {
					if strings.Index(wrongNote, "同一案件残疾代码不可重复；") == -1 {
						wrongNote += "同一案件残疾代码不可重复；"
					}
				}
			}
		}
	}
	// CSB0116RC0147000 同一案件存在多个fc109时，对所有fc109的结果值进行校验（结果值为空的不进行校验），当fc109的结果值重复时，导出校验出提示：同一案件特种病类型不可重复；
	if len(fc109s) > 0 {
		fc109One := fc109s[0]
		if len(fc109s) > 1 {
			for i := 1; i < len(fc109s); i++ {
				if fc109One == fc109s[i] {
					if strings.Index(wrongNote, "同一案件特种病类型不可重复；") == -1 {
						wrongNote += "同一案件特种病类型不可重复；"
					}
				}
			}
		}
	}
	// CSB0116RC0148000 同一案件存在多个fc110时，对所有fc110的结果值进行校验（结果值为空的不进行校验），当fc110的结果值重复时，导出校验出提示：同一案件失能类型不可重复；
	if len(fc110s) > 0 {
		fc110One := fc110s[0]
		if len(fc110s) > 1 {
			for i := 1; i < len(fc110s); i++ {
				if fc110One == fc110s[i] {
					if strings.Index(wrongNote, "同一案件失能类型不可重复；") == -1 {
						wrongNote += "同一案件失能类型不可重复；"
					}
				}
			}
		}
	}
	return wrongNote
}

func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""
	//CSB0116RC0127000 所有的hospitalCode和unitName中（hospitalCode、unitName可能有多个且所在节点位置不一样），存在任一一个值为空时，出导出校验，提示：医疗机构异常，请检查；
	//（例如XML中有多个hospitalCode的值为空，该提示也只出一个）
	hospNode := []string{"hospitalCode", "unitName"}
	for _, node := range hospNode {
		items := RegMatchAll(xmlValue, `<(`+node+`)>[\s\S]*?<\/(`+node+`)>`)
		for _, item := range items {
			nodeVal := GetNodeValue(item, node)
			if nodeVal == "" {
				if strings.Index(wrongNote, "医疗机构异常，请检查；") == -1 {
					wrongNote += "医疗机构异常，请检查；"
				}
			}

		}
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
func IsExist(fields []model.ProjectField, code string) bool {
	isExist := false
	for _, field := range fields {
		if field.Code == code {
			isExist = true
		}
	}
	return isExist
}
