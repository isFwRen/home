/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/9/2 10:30
 */

package B0108

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"server/global"
	model2 "server/module/export/model"
	"server/module/export/utils"
	utils3 "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	utils2 "server/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
)

//模板类型字段
//1-全国版申请书（影像上方有申请人声明及授权）；
//2-反欺诈新版申请书（影像上方直接为基本信息栏）；
//3-第三方报销；
//5-门诊发票；
//6-住院发票；
//7-清单；
//8-报销单；
//9-诊断书；
//10-死亡证明书/火化证明书/户籍注销证明书/残疾/重疾诊断/特殊疾病；
//11-证件信息；
//12-存折/银行卡；
//13-手术（发票上印有手术费的需要切手术）；
//14-住院日期（出院小结/出院记录/出院证明/出院证，无以上四类影像则切空白处）；
//16-住院诊查费天数（费用清单）
//5-门诊发票；
//6-住院发票；
//7-清单；
//8-报销单；
//13-手术（发票上印有手术费的需要切手术）；
//14-住院日期（出院小结/出院记录/出院证明/出院证，无以上四类影像则切空白处）；
//16-住院诊查费天数（费用清单）

//初审关系字段
//fc096门诊票据属性-模板
//fc097住院票据属性-模板
//fc089清单所属发票
//fc091报销单所属发票
//fc101手术所属发票
//fc213住院日期所属发票
//fc225住院诊查费所属发票

//非初审分块
//发票 	bc002，bc003
//清单 	bc004
//报销单 	bc009
//手术 	bc013
//住院日期 bc001
//住院诊查费天数 无非初审

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:             model2.TypeCode{FieldCode: []string{"fc096", "fc097"}, BlockCode: []string{"bc002", "bc003"}},
	QingDanFieldCode:             model2.TypeCode{FieldCode: []string{"fc089"}, BlockCode: []string{"bc004"}},
	BaoXiaoDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc091"}, BlockCode: []string{"bc009"}},
	OperationFieldCode:           model2.TypeCode{FieldCode: []string{"fc101"}, BlockCode: []string{"bc013"}},
	HospitalizationDateFieldCode: model2.TypeCode{FieldCode: []string{"fc213"}, BlockCode: []string{"bc001"}},
	HospitalizationFeeFieldCode:  model2.TypeCode{FieldCode: []string{"fc225"}, BlockCode: []string{""}},
	OtherTempType:                map[string]string{"1": "1", "2": "2", "3": "3", "9": "9", "10": "10", "11": "11", "12": "12", "15": "15", "99": "99"},
	TempTypeField:                "fc084",
	InvoiceNumField:              []string{"fc032", "fc054"},
	MoneyField:                   []string{"fc039", "fc062"},
	//InvoiceTypeField: "fc003",
}

// ResultData B0108
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {
	obj = utils.RelationDealManyInvoiceType(bill, blocks, fieldMap, relationConf)
	//常量
	constMap := constDeal(bill.ProCode)

	fieldLocationMap := make(map[string][][]int)
	n := 0
	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票
		invoiceMap := obj.Invoice[i]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					for l := 0; l < len(fields); l++ {
						//fields[l].Issues = nil
						//i:发票index j:发票结构体字段index k:字段二维数组x的index l:字段二维数组y的index
						fieldLocationMap[fields[l].Code] = append(fieldLocationMap[fields[l].Code], []int{i, j, k, l, fields[l].BlockIndex})
						n++
					}
				}
			}
		}
	}
	fmt.Println(n)

	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票的所有关系信息
		invoiceMap := obj.Invoice[i]
		global.GLog.Info("B108", zap.Any("invoiceMap", invoiceMap.Id))
	}

	blockMap := make(map[string]string)
	for i := 0; i < len(blocks); i++ {
		//同一发票的所有关系信息
		block := blocks[i]
		blockMap[block.Code] = block.Model.ID
	}
	//CSB0108RC0022000
	// 校验fc084的录入值（fc084为循环字段），统计fc084录入5的数量，将数量赋值给fc003
	//CSB0108RC0023000
	// 校验fc084的录入值（fc084为循环字段），统计fc084录入6的数量，将数量赋值给fc004
	global.GLog.Info("CSB0108RC0023000")
	fc084Locs := fieldLocationMap["fc084"]
	fc084Count5 := 0
	fc084Count6 := 0
	for _, loc := range fc084Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				if fieldsArr[loc[2]][loc[3]].ResultValue == "5" {
					fc084Count5++
				} else if fieldsArr[loc[2]][loc[3]].ResultValue == "6" {
					fc084Count6++
				}
			}
		}
	}
	fc003Locs := fieldLocationMap["fc003"]
	for _, loc := range fc003Locs {
		obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = strconv.Itoa(fc084Count5)
	}
	fc004Locs := fieldLocationMap["fc004"]
	for _, loc := range fc004Locs {
		obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = strconv.Itoa(fc084Count6)
	}
	//CSB0108RC0024000
	// 当fc003的结果值与fc020的结果值不一致时，fc003出问题件：实际门诊发票张数与申请表中【门诊发票张数】不一致，请核查
	global.GLog.Info("CSB0108RC0024000")
	fc003Locs = fieldLocationMap["fc003"]
	fc020Locs := fieldLocationMap["fc020"]
	if len(fc003Locs) > 0 && len(fc020Locs) > 0 {
		fc003Obj := obj.Invoice[fc003Locs[0][0]].Invoice[fc003Locs[0][2]][fc003Locs[0][3]]
		fc020Obj := obj.Invoice[fc020Locs[0][0]].Invoice[fc020Locs[0][2]][fc020Locs[0][3]]
		if fc003Obj.FinalValue != fc020Obj.FinalValue {
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "实际门诊发票张数与申请表中【门诊发票张数】不一致，请核查",
			}
			obj.Invoice[fc003Locs[0][0]].Invoice[fc003Locs[0][2]][fc003Locs[0][3]].Issues = append(fc003Obj.Issues, issue)
		}
	}
	//CSB0108RC0025000
	// 当fc004的结果值与fc022的结果值不一致时，fc004出问题件：实际住院发票张数与申请表中【住院发票张数】不一致，请核查
	global.GLog.Info("CSB0108RC0025000")
	fc004Locs = fieldLocationMap["fc004"]
	fc022Locs := fieldLocationMap["fc022"]
	if len(fc004Locs) > 0 && len(fc022Locs) > 0 {
		fc004Obj := obj.Invoice[fc004Locs[0][0]].Invoice[fc004Locs[0][2]][fc004Locs[0][3]]
		fc022Obj := obj.Invoice[fc022Locs[0][0]].Invoice[fc022Locs[0][2]][fc022Locs[0][3]]
		if fc004Obj.FinalValue != fc022Obj.FinalValue {
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "实际住院发票张数与申请表中【住院发票张数】不一致，请核查",
			}
			obj.Invoice[fc004Locs[0][0]].Invoice[fc004Locs[0][2]][fc004Locs[0][3]].Issues = append(fc004Obj.Issues, issue)
		}
	}
	//CSB0108RC0026000
	// 将每个fc041、fc064的结果值相加之和与fc025结果值作对比，不一致时fc025出问题件：发票其他扣除金额与申请书其他报销金额不一致
	global.GLog.Info("CSB0108RC0026000")
	fc041Locs := fieldLocationMap["fc041"]
	valueSum := decimal.NewFromFloat32(0.0)
	for _, loc := range fc041Locs {
		fc041Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		d, err := decimal.NewFromString(fc041Obj.FinalValue)
		if err != nil {
			continue
		}
		valueSum = valueSum.Add(d)
	}
	fc064Locs := fieldLocationMap["fc064"]
	for _, loc := range fc064Locs {
		fc064Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		d, err := decimal.NewFromString(fc064Obj.FinalValue)
		if err != nil {
			continue
		}
		valueSum = valueSum.Add(d)
	}
	fmt.Println(valueSum.String())
	fc025Locs := fieldLocationMap["fc025"]
	if len(fc025Locs) > 0 {
		fc025Obj := obj.Invoice[fc025Locs[0][0]].Invoice[fc025Locs[0][2]][fc025Locs[0][3]]
		if ParseFloat(fc025Obj.FinalValue) != ParseFloat(valueSum.String()) {
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "发票其他扣除金额与申请书其他报销金额不一致",
			}
			obj.Invoice[fc025Locs[0][0]].Invoice[fc025Locs[0][2]][fc025Locs[0][3]].Issues = append(fc025Obj.Issues, issue)
		}
	}

	//CSB0108RC0027000
	// 当fc110录入值为1时，将fc010结果值赋值给fc009
	global.GLog.Info("CSB0108RC0027000")
	fc110Locs := fieldLocationMap["fc110"]
	if len(fc110Locs) > 0 {
		fc110Obj := obj.Invoice[fc110Locs[0][0]].Invoice[fc110Locs[0][2]][fc110Locs[0][3]]
		if fc110Obj.ResultValue == "1" {
			fc010Locs := fieldLocationMap["fc010"]
			fc009Locs := fieldLocationMap["fc009"]
			if len(fc010Locs) > 0 && len(fc009Locs) > 0 {
				obj.Invoice[fc009Locs[0][0]].Invoice[fc009Locs[0][2]][fc009Locs[0][3]].FinalValue = obj.Invoice[fc010Locs[0][0]].Invoice[fc010Locs[0][2]][fc010Locs[0][3]].FinalValue
			}
		}
	}
	/*
		CSB0108RC0028000
		fc017、fc035、fc056、fc057、fc099、fc158、fc164、fc192、
		fc193、fc194、fc195、fc175、fc200、fc201、fc202、fc203、
		fc214、fc215、fc217
		当以上字段录入值为6位数时，在其结果值前面增加20（如字段录入值为190101，则对应结果值为20190101）
	*/
	global.GLog.Info("CSB0108RC0028000")
	fieldCodeS := [20]string{"fc017", "fc035", "fc056", "fc057", "fc099",
		"fc158", "fc164", "fc192", "fc193", "fc194", "fc195", "fc175",
		"fc200", "fc201", "fc202", "fc203", "fc214", "fc215", "fc217", "fc265"}
	for _, fieldCode := range fieldCodeS {
		fieldLocs := fieldLocationMap[fieldCode]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					if len(fieldsArr[loc[2]][loc[3]].ResultValue) == 6 {
						fieldsArr[loc[2]][loc[3]].FinalValue = "20" + fieldsArr[loc[2]][loc[3]].FinalValue
					}
				}
			}
		}
	}

	//CSB0108RC0029000
	// fc099录入值为A或B或?或E时，结果值默认为当前日期；
	global.GLog.Info("CSB0108RC0029000")
	fc099Locs := fieldLocationMap["fc099"]
	for _, loc := range fc099Locs {
		fc099Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		compile, err := regexp.Compile("\\?|？|A|B|E")
		if err != nil {
			global.GLog.Error("CSB0108RC0029000" + err.Error())
			continue
		}
		if compile.MatchString(fc099Obj.ResultValue) {
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = time.Now().Format("20060102")
		}
	}

	//CSB0108RC0030000
	// 将fc019结果值中的“/”转换成英文状态下的逗号“,”
	global.GLog.Info("CSB0108RC0030000")
	fc019Locs := fieldLocationMap["fc019"]
	for _, loc := range fc019Locs {
		field019Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = strings.Replace(field019Obj.FinalValue, "/", ",", -1)
	}
	/*
		CSB0108RC0031000
		（当fc012录入值为A时，不执行该校验）
		当fc012录入值存在多个且包含2时，结果值转为2；当fc012录入值存在多个且不包含2时，取第一个数字作为fc012结果值
		如：fc012录入值为1/2/3/5，则结果值为2；
		fc012录入值为1/3/4，则结果值为1；
	*/
	global.GLog.Info("CSB0108RC0031000")
	fc012Locs := fieldLocationMap["fc012"]
	for _, loc := range fc012Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				if fieldsArr[loc[2]][loc[3]].ResultValue != "A" {
					strs := strings.Split(fieldsArr[loc[2]][loc[3]].ResultValue, "/")
					flag := false
					for _, str := range strs {
						if str == "2" {
							flag = true
							break
						}
					}
					if flag {
						fieldsArr[loc[2]][loc[3]].FinalValue = "2"
					} else {
						fieldsArr[loc[2]][loc[3]].FinalValue = strs[0]
					}
				}
			}
		}
	}

	/*
		CSB0108RC0038000
		fc018、fc112、fc036、fc104、fc058、fc161、
		fc162、fc205、fc206、fc207、fc208、fc209
		录入值为《B0108_太平理赔_诊断代码表》中“ACCIDENT_NAME”（第二列），
		结果值根据常量表转换成对应的“ACCIDENT_ID”（第一列）值最长的那条代码
	*/
	global.GLog.Info("CSB0108RC0038000")
	codes := []string{"fc018", "fc112", "fc036", "fc104", "fc058", "fc161", "fc162", "fc205", "fc206", "fc207", "fc208", "fc209", "fc239", "fc240", "fc241", "fc242", "fc243", "fc244", "fc245", "fc246", "fc247", "fc248"}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					fieldsArr[loc[2]][loc[3]].FinalValue = constMap["zhenDuanDaiMaMap"][fieldsArr[loc[2]][loc[3]].ResultValue]
				}
			}
		}
	}

	//CSB0108RC0032000
	// 当fc157录入值为2时，将fc036的结果值赋值给fc018
	global.GLog.Info("CSB0108RC0032000")
	fc157Locs := fieldLocationMap["fc157"]
	if len(fc157Locs) > 0 {
		fc157Obj := obj.Invoice[fc157Locs[0][0]].Invoice[fc157Locs[0][2]][fc157Locs[0][3]]
		if fc157Obj.ResultValue == "2" {
			fc036Locs := fieldLocationMap["fc036"]
			fc018Locs := fieldLocationMap["fc018"]
			if len(fc036Locs) > 0 && len(fc018Locs) > 0 {
				obj.Invoice[fc018Locs[0][0]].Invoice[fc018Locs[0][2]][fc018Locs[0][3]].FinalValue = obj.Invoice[fc036Locs[0][0]].Invoice[fc036Locs[0][2]][fc036Locs[0][3]].FinalValue
			}
		}
	}

	//fc035Locs := fieldLocationMap["fc035"]
	//fc035CountMap := map[string]int{}
	//for _, loc := range fc035Locs {
	//	invoiceMap := obj.Invoice[loc[0]]
	//	eleLen := reflect.ValueOf(invoiceMap).NumField()
	//	if eleLen > 0 {
	//		if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
	//			//每张发票每种类型的字段
	//			fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
	//			fc032Loc := utils2.GetFieldLoc(fieldLocationMap["fc032"], loc[0], loc[1], -1, -1)
	//			if len(fc032Loc) > 0 {
	//				fc032Obj := obj.Invoice[fc032Loc[0][0]].Invoice[fc032Loc[0][2]][fc032Loc[0][3]]
	//				count, ok := fc035CountMap[fieldsArr[loc[2]][loc[3]].FinalValue]
	//				sId := "01"
	//				if ok {
	//					count = count + 1
	//				} else {
	//					count = 1
	//				}
	//				fc035CountMap[fieldsArr[loc[2]][loc[3]].FinalValue] = count
	//				if count < 10 {
	//					sId = "0" + strconv.Itoa(count)
	//				} else {
	//					sId = strconv.Itoa(count)
	//				}
	//				compile, err := regexp.Compile(`^(\\?|？|A)$`)
	//				if err != nil {
	//					global.GLog.Error("CSB0108RC0033000" + err.Error())
	//					continue
	//				}
	//				if compile.MatchString(fc032Obj.ResultValue) {
	//					obj.Invoice[fc032Loc[0][0]].Invoice[fc032Loc[0][2]][fc032Loc[0][3]].FinalValue = fieldsArr[loc[2]][loc[3]].FinalValue + sId
	//					issue := model3.Issue{
	//						Type:    "",
	//						Code:    "",
	//						Message: "门诊发票：第[" + (fieldsArr[loc[2]][loc[3]].FinalValue + sId) + "]门诊号模糊不清",
	//					}
	//					obj.Invoice[fc032Loc[0][0]].Invoice[fc032Loc[0][2]][fc032Loc[0][3]].Issues = append(fc032Obj.Issues, issue)
	//				}
	//			}
	//		}
	//	}
	//}

	//CSB0108RC0034000
	// 当fc033结果值与fc009结果值不一致，fc033出问题件：病人姓名与申请表事故者不一致
	global.GLog.Info("CSB0108RC0034000")
	fc009Locs := fieldLocationMap["fc009"]
	if len(fc009Locs) > 0 {
		fc009Obj := obj.Invoice[fc009Locs[0][0]].Invoice[fc009Locs[0][2]][fc009Locs[0][3]]
		fc033Locs := fieldLocationMap["fc033"]
		for _, loc := range fc033Locs {
			fc033Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
			if fc033Obj.FinalValue != fc009Obj.FinalValue {
				issue := model3.Issue{
					Type:    "",
					Code:    "",
					Message: "病人姓名与申请表事故者不一致",
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc033Obj.Issues, issue)
			}
		}
	}

	//CSB0108RC0035000
	// fc034、fc055录入值为《B0108_太平理赔_医院代码表》中“医院名称”（第二列），结果值根据常量表转换成对应的“医院编码”（第一列）
	global.GLog.Info("CSB0108RC0035000")
	codes = []string{"fc034", "fc055"}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					fieldsArr[loc[2]][loc[3]].FinalValue = constMap["yiYuanDaiMaMap"][fieldsArr[loc[2]][loc[3]].ResultValue]
				}
			}
		}
	}
	//CSB0108RC0036000
	// 当fc034录入值为B时，出问题件：医院名称为空或模糊
	global.GLog.Info("CSB0108RC0036000")
	fc034Locs := fieldLocationMap["fc034"]
	for _, loc := range fc034Locs {
		fc034Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		if fc034Obj.ResultValue == "B" {
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "医院名称为空或模糊",
			}
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc034Obj.Issues, issue)
		}
	}
	/*
		CSB0108RC0037000
		当fc035录入值为A或?或？时，fc035结果值默认为当前日期并出问题件：门诊发票：第[xxx]就诊时间模糊不清（xxx为fc032结果值）
		如：20201226
	*/
	global.GLog.Info("CSB0108RC0037000")
	fc035Locs := fieldLocationMap["fc035"]
	for _, loc := range fc035Locs {
		fc035Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		// compile, err := regexp.Compile(`^(\\?|？|A)$`)
		// if err != nil {
		// 	global.GLog.Error("CSB0108RC0037000" + err.Error())
		// 	continue
		// }
		// fmt.Println("-------------11111111111111111111111111111111111111-----------field-------------------", fc035Obj.ResultValue, RegIsMatch(`^(\\?|？|A)$`, fc035Obj.ResultValue))
		if RegIsMatch(`^(\\?|？|A)$`, fc035Obj.ResultValue) {
			fc032Locs := utils2.GetFieldLoc(fieldLocationMap["fc032"], loc[0], -1, loc[2], -1)
			fc032Val := obj.Invoice[loc[0]].Invoice[loc[2]][fc032Locs[0][3]].FinalValue
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "门诊发票：第[" + fc032Val + "]就诊时间模糊不清",
			}
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = time.Now().Format("20060102")
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc035Obj.Issues, issue)
		}
	}

	/*
		CSB0108RC0039000
		当左边字段录入值不为“B”且不为空时，对应右边字段赋值为“NA”
		fc036、fc228
		fc104、fc229
		fc058、fc230
		fc161、fc231
		fc162、fc232
		fc205、fc233
		fc206、fc234
		fc207、fc235
		fc208、fc236
		fc209、fc237
	*/
	// global.GLog.Info("CSB0108RC0039000")
	// iCodes := [][]string{
	// 	{"fc036", "fc228"},
	// 	{"fc104", "fc229"},
	// 	{"fc058", "fc230"},
	// 	{"fc161", "fc231"},
	// 	{"fc162", "fc232"},
	// 	{"fc205", "fc233"},
	// 	{"fc206", "fc234"},
	// 	{"fc207", "fc235"},
	// 	{"fc208", "fc236"},
	// 	{"fc209", "fc237"},
	// 	{"fc239", "fc249"},
	// 	{"fc240", "fc250"},
	// 	{"fc241", "fc251"},
	// 	{"fc242", "fc252"},
	// 	{"fc243", "fc253"},
	// 	{"fc244", "fc254"},
	// 	{"fc245", "fc255"},
	// 	{"fc246", "fc256"},
	// 	{"fc247", "fc257"},
	// 	{"fc248", "fc258"},
	// }
	// for _, iCode := range iCodes {
	// 	fieldLLocs := fieldLocationMap[iCode[0]]
	// 	if len(fieldLLocs) > 0 {
	// 		fieldLObj := obj.Invoice[fieldLLocs[0][0]].Invoice[fieldLLocs[0][2]][fieldLLocs[0][3]]
	// 		if fieldLObj.ResultValue != "" && fieldLObj.ResultValue != "B" {
	// 			fieldRLocs := fieldLocationMap[iCode[1]]
	// 			if len(fieldRLocs) > 0 {
	// 				obj.Invoice[fieldRLocs[0][0]].Invoice[fieldRLocs[0][2]][fieldRLocs[0][3]].FinalValue = "NA"
	// 			}
	// 		}
	// 	}

	// }
	//CSB0108RC0040000
	// 当以下字段录入值包含?或？时，对应字段结果值默认为“其他”，并出问题件：诊断结论模糊；
	// fc228、fc229、fc230、fc231、fc232、fc233、fc234、fc235、fc236、fc237
	global.GLog.Info("CSB0108RC0040000")
	codes = []string{"fc228", "fc229", "fc230", "fc231", "fc232", "fc233", "fc234", "fc235", "fc236", "fc237", "fc249", "fc250", "fc251", "fc252", "fc253", "fc254", "fc255", "fc256", "fc257", "fc258"}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					compile, err := regexp.Compile("\\?|？")
					if err != nil {
						global.GLog.Error("CSB0108RC0040000" + err.Error())
						continue
					}
					if compile.MatchString(fieldsArr[loc[2]][loc[3]].ResultValue) {
						fieldsArr[loc[2]][loc[3]].FinalValue = "其他"
						issue := model3.Issue{
							Type:    "",
							Code:    "",
							Message: "诊断结论模糊；",
						}
						fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
					}
				}
			}
		}
	}
	/*
		CSB0108RC0041000
		fc106、fc165、fc196、fc197、fc198、fc199录入值为《B0108_太平理赔_手术代码表》中“手术名称”（第二列），
		结果值根据常量表转换成对应的“手术代码”（第一列）
	*/
	global.GLog.Info("CSB0108RC0041000")
	codes = []string{"fc165", "fc196", "fc197", "fc198", "fc199", "fc266"}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					fieldsArr[loc[2]][loc[3]].FinalValue = constMap["shouShuDaiMaMap"][fieldsArr[loc[2]][loc[3]].ResultValue]
				}
			}
		}
	}

	/*
		CSB0108RC0044000
		MB002-bc002为“门诊发票”分块，同一分块下：
		1、当第一列字段的结果值为第三列的某一内容时，则将第一列对应的第二列字段的结果值赋值给与第三列内容所对应的第四列字段结果值处；
		（如，当第一列字段的fc113结果值为“诊查费”则将对应的第二列字段fc133的结果值赋值给“诊查费”所对应的字段fc046的结果值处）
		2、当第一列字段存在多个第三列相同结果值时，则将第一列字段相同结果值所对应的第二列字段的结果值相加求和，再赋值给与第三列内容所对应的第四列字段结果值处；
		（如存在多个“检查费”，fc113和fc114的结果值均为“检查费”，则将fc133、fc134两个字段的结果值相加求和所得的值，赋值给“检查费”所对应的字段fc047的结果值处）

		fc113、fc133    西药费      fc044
		fc114、fc134    中成药费    fc100
		fc115、fc135    中草药费    fc102
		fc116、fc136    诊查费      fc046
		fc117、fc137    检查费      fc047
		fc118、fc005    化验费      fc048
		fc119、fc006    特殊检查费  fc049
		fc120、fc007    治疗费      fc050
		fc121、fc008    手术费      fc051
		fc122、fc037    材料费      fc052
		fc123、fc038    其他费      fc053
		fc124、fc060
		fc125、fc061
		fc126、fc086
		fc127、fc087
		fc128、fc088
		fc129、fc092
		fc130、fc093
		fc131、fc098
		fc132、fc109

		CSB0108RC0054000
		MB002-bc003为“住院发票”分块，同一分块下：
		1、当第一列字段的结果值为第三列的某一内容时，则将第一列对应的第二列字段的结果值赋值给与第三列内容所对应的第四列字段结果值处；
		（如，当第一列字段的fc113结果值为“诊查费”则将对应的第二列字段fc133的结果值赋值给“诊查费”所对应的字段fc069的结果值处）
		2、当第一列字段存在多个第三列相同结果值时，则将第一列字段相同结果值所对应的第二列字段的结果值相加求和，再赋值给与第三列内容所对应的第四列字段结果值处；
		（如存在多个“检查费”，fc113和fc114的结果值均为“检查费”，则将fc133、fc134两个字段的结果值相加求和所得的值，赋值给“检查费”所对应的字段fc070的结果值处）

		fc113、fc133    西药费      fc066
		fc114、fc134    中成药费    fc067
		fc115、fc135    中草药费    fc068
		fc116、fc136    诊查费      fc069
		fc117、fc137    检查费      fc070
		fc118、fc005    化验费      fc071
		fc119、fc006    特殊检查费  fc072
		fc120、fc007    治疗费      fc073
		fc121、fc008    手术费      fc076
		fc122、fc037    材料费      fc078
		fc123、fc038    其他费      fc079
		fc124、fc060    输血费      fc074
		fc125、fc061    护理费      fc075
		fc126、fc086    床位费      fc077
		fc127、fc087
		fc128、fc088
		fc129、fc092
		fc130、fc093
		fc131、fc098
		fc132、fc109
	*/
	global.GLog.Info("CSB0108RC0044000")
	global.GLog.Info("CSB0108RC0054000")
	iCodes := [][]string{
		{"fc113", "fc133"},
		{"fc114", "fc134"},
		{"fc115", "fc135"},
		{"fc116", "fc136"},
		{"fc117", "fc137"},
		{"fc118", "fc005"},
		{"fc119", "fc006"},
		{"fc120", "fc007"},
		{"fc121", "fc008"},
		{"fc122", "fc037"},
		{"fc123", "fc038"},
		{"fc124", "fc060"},
		{"fc125", "fc061"},
		{"fc126", "fc086"},
		{"fc127", "fc087"},
		{"fc128", "fc088"},
		{"fc129", "fc092"},
		{"fc130", "fc093"},
		{"fc131", "fc098"},
		{"fc132", "fc109"},
	}
	iMap1 := map[string]string{
		"西药费":   "fc044",
		"中成药费":  "fc100",
		"中草药费":  "fc102",
		"诊查费":   "fc046",
		"检查费":   "fc047",
		"化验费":   "fc048",
		"特殊检查费": "fc049",
		"治疗费":   "fc050",
		"手术费":   "fc051",
		"材料费":   "fc052",
		"其他费":   "fc053",
	}
	iMap2 := map[string]string{
		"西药费":   "fc066",
		"中成药费":  "fc067",
		"中草药费":  "fc068",
		"诊查费":   "fc069",
		"检查费":   "fc070",
		"化验费":   "fc071",
		"特殊检查费": "fc072",
		"治疗费":   "fc073",
		"手术费":   "fc076",
		"材料费":   "fc078",
		"其他费":   "fc079",
		"输血费":   "fc074",
		"护理费":   "fc075",
		"床位费":   "fc077",
	}
	for _, iCode := range iCodes {
		fieldLLocs := fieldLocationMap[iCode[0]]
		for _, loc := range fieldLLocs {
			fieldLObj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
			s, ok := iMap1[fieldLObj.FinalValue]
			if obj.Invoice[loc[0]].InvoiceType == "门诊" {
				s, ok = iMap1[fieldLObj.FinalValue]
			} else if obj.Invoice[loc[0]].InvoiceType == "住院" {
				s, ok = iMap2[fieldLObj.FinalValue]
			}
			if ok {
				fieldRLocs := utils2.GetFieldLoc(fieldLocationMap[iCode[1]], loc[0], loc[1], loc[2], -1)
				fieldRObj := obj.Invoice[fieldRLocs[0][0]].Invoice[fieldRLocs[0][2]][fieldRLocs[0][3]]
				fieldLoc := utils2.GetFieldLoc(fieldLocationMap[s], loc[0], loc[1], loc[2], -1)
				fieldObj := obj.Invoice[fieldLoc[0][0]].Invoice[fieldLoc[0][2]][fieldLoc[0][3]]
				findValue, _ := decimal.NewFromString(fieldObj.FinalValue)
				d, _ := decimal.NewFromString(fieldRObj.FinalValue)
				obj.Invoice[fieldLoc[0][0]].Invoice[fieldLoc[0][2]][fieldLoc[0][3]].FinalValue = findValue.Add(d).String()
				finalValueFloat, err := strconv.ParseFloat(findValue.Add(d).String(), 64)
				if err != nil {
					continue
				}
				value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", finalValueFloat), 64)
				obj.Invoice[fieldLoc[0][0]].Invoice[fieldLoc[0][2]][fieldLoc[0][3]].FinalValue = strconv.FormatFloat(value, 'f', 2, 64)
			}
		}
	}
	//CSB0108RC0045000
	// fc188为循环分块MB002-bc004中的字段，fc188字段的结果值默认为该分块对应的同一发票属性的fc032结果值或fc054的结果值
	global.GLog.Info("CSB0108RC0045000")
	fc188Locs := fieldLocationMap["fc188"]
	for _, loc := range fc188Locs {
		fieldLoc := utils2.GetFieldLoc(fieldLocationMap["fc032"], loc[0], -1, -1, -1)
		if len(fieldLoc) == 0 {
			fieldLoc = utils2.GetFieldLoc(fieldLocationMap["fc054"], loc[0], -1, -1, -1)
		}
		if len(fieldLoc) > 0 {
			fieldObj := obj.Invoice[loc[0]].Invoice[fieldLoc[0][2]][fieldLoc[0][3]]
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					fieldsArr[loc[2]][loc[3]].FinalValue = fieldObj.FinalValue
				}
			}
		}
	}

	/*
		CSB0108RC0046000
		第一列字段和第二列字段为一一对应关系，两者必须同时有值或同时为空，当其中一个字段有值，
		对应字段为空时，空的字段出问题件：项目名称与项目金额不匹配，请确认修改；
		项目名称，项目金额
		fc080，fc081
		fc143，fc144
		fc145，fc146
		fc147，fc148
		fc149，fc150
		fc151，fc152
		fc153，fc154
		fc155，fc156

		如：fc080结果值为“氯化钠”，fc081结果值为空，则fc081出问题件；
		fc080结果值为空，fc081结果值为“10.00”，则fc080出问题件；"
	*/
	global.GLog.Info("CSB0108RC0046000")
	iCodes = [][]string{
		{"fc080", "fc081"},
		{"fc143", "fc144"},
		{"fc145", "fc146"},
		{"fc147", "fc148"},
		{"fc149", "fc150"},
		{"fc151", "fc152"},
		{"fc153", "fc154"},
		{"fc155", "fc156"},
	}
	for _, iCode := range iCodes {
		fieldLLocs := fieldLocationMap[iCode[0]]
		fieldRLocs := fieldLocationMap[iCode[1]]
		if len(fieldLLocs) > 0 && len(fieldRLocs) > 0 && len(fieldLLocs) == len(fieldRLocs) {
			for i := 0; i < len(fieldLLocs); i++ {
				fieldLObj := obj.Invoice[fieldLLocs[i][0]].QingDan[fieldLLocs[i][2]][fieldLLocs[i][3]]
				fieldRObj := obj.Invoice[fieldRLocs[i][0]].QingDan[fieldRLocs[i][2]][fieldRLocs[i][3]]
				if fieldLObj.FinalValue == "" && fieldRObj.FinalValue != "" {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "项目名称与项目金额不匹配，请确认修改；",
					}
					obj.Invoice[fieldLLocs[i][0]].QingDan[fieldLLocs[i][2]][fieldLLocs[i][3]].Issues = append(fieldLObj.Issues, issue)
				} else if fieldLObj.FinalValue != "" && fieldRObj.FinalValue == "" {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "项目名称与项目金额不匹配，请确认修改；",
					}
					obj.Invoice[fieldRLocs[i][0]].QingDan[fieldRLocs[i][2]][fieldRLocs[i][3]].Issues = append(fieldRObj.Issues, issue)
				}
			}
		}
	}

	//fc056Locs := fieldLocationMap["fc056"]
	//fc056CountMap := map[string]int{}
	//for _, loc := range fc056Locs {
	//	invoiceMap := obj.Invoice[loc[0]]
	//	eleLen := reflect.ValueOf(invoiceMap).NumField()
	//	if eleLen > 0 {
	//		if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
	//			//每张发票每种类型的字段
	//			fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
	//			fc054Loc := utils2.GetFieldLoc(fieldLocationMap["fc054"], loc[0], loc[1], -1, -1)
	//			if len(fc054Loc) > 0 {
	//				fc054Obj := obj.Invoice[fc054Loc[0][0]].Invoice[fc054Loc[0][2]][fc054Loc[0][3]]
	//				count, ok := fc056CountMap[fieldsArr[loc[2]][loc[3]].FinalValue]
	//				sId := "01"
	//				if ok {
	//					count = count + 1
	//				} else {
	//					count = 1
	//				}
	//				fc056CountMap[fieldsArr[loc[2]][loc[3]].FinalValue] = count
	//				if count < 10 {
	//					sId = "0" + strconv.Itoa(count)
	//				} else {
	//					sId = strconv.Itoa(count)
	//				}
	//				compile, err := regexp.Compile(`^(\\?|？|A)$`)
	//				if err != nil {
	//					global.GLog.Error("CSB0108RC0047000" + err.Error())
	//					continue
	//				}
	//				if compile.MatchString(fc054Obj.ResultValue) {
	//					obj.Invoice[fc054Loc[0][0]].Invoice[fc054Loc[0][2]][fc054Loc[0][3]].FinalValue = fieldsArr[loc[2]][loc[3]].FinalValue + sId
	//					issue := model3.Issue{
	//						Type:    "",
	//						Code:    "",
	//						Message: "住院发票：第[" + (fieldsArr[loc[2]][loc[3]].FinalValue + sId) + "]住院号模糊不清",
	//					}
	//					obj.Invoice[fc054Loc[0][0]].Invoice[fc054Loc[0][2]][fc054Loc[0][3]].Issues = append(fc054Obj.Issues, issue)
	//				}
	//			}
	//		}
	//	}
	//}

	//CSB0108RC0048000
	// 当fc055录入值为B时，出问题件：医院名称为空或模糊
	global.GLog.Info("CSB0108RC0048000")
	fc055Locs := fieldLocationMap["fc055"]
	for _, loc := range fc055Locs {
		fc055Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		if fc055Obj.ResultValue == "B" {
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "医院名称为空或模糊",
			}
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc055Obj.Issues, issue)
		}
	}
	//CSB0108RC0049000
	// 当fc056录入值为A或?或？时，fc056结果值默认为当前日期并出问题件：住院发票：第[xxx]就诊时间模糊不清（xxx为fc054结果值）
	// 如：20201226
	global.GLog.Info("CSB0108RC0049000")
	fc056Locs := fieldLocationMap["fc056"]
	for _, loc := range fc056Locs {
		fc056Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		// compile, err := regexp.Compile(`^(\\?|？|A)$`)
		// if err != nil {
		// 	global.GLog.Error("CSB0108RC0049000" + err.Error())
		// 	continue
		// }
		if RegIsMatch(`^(\\?|？|A)$`, fc056Obj.ResultValue) {
			fc054Locs := utils2.GetFieldLoc(fieldLocationMap["fc054"], loc[0], -1, loc[2], -1)
			fc054Val := obj.Invoice[loc[0]].Invoice[loc[2]][fc054Locs[0][3]].FinalValue
			issue := model3.Issue{
				Type:    "",
				Code:    "",
				Message: "住院发票：第[" + fc054Val + "]就诊时间模糊不清",
			}
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = time.Now().Format("20060102")
			obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc056Obj.Issues, issue)
		}
	}
	//CSB0108RC0050000
	// 同一发票属性下，校验fc214与fc056的结果值，不一致时，取最晚时间赋值给fc056
	// （当fc214录入值为A、?、？、空值时，不执行该校验）
	global.GLog.Info("CSB0108RC0050000")
	fc214Locs := fieldLocationMap["fc214"]
	for _, loc := range fc214Locs {
		fc214Obj := obj.Invoice[loc[0]].HospitalizationDate[loc[2]][loc[3]]
		compile, err := regexp.Compile("\\?|？|A")
		if err != nil {
			global.GLog.Error("CSB0108RC0050000" + err.Error())
			continue
		}
		if !(compile.MatchString(fc214Obj.ResultValue) || fc214Obj.ResultValue == "") {
			global.GLog.Info("", zap.Any("", fieldLocationMap["fc056"]))
			fc056Loc := utils2.GetFieldLoc(fieldLocationMap["fc056"], loc[0], -1, -1, -1)
			if len(fc056Loc) <= 0 {
				continue
			}
			if len(fc056Loc[0]) <= 0 {
				continue
			}
			invoiceMap := obj.Invoice[fc056Loc[0][0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(fc056Loc[0][1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(fc056Loc[0][1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(fc056Loc[0][1]).Interface().([][]model3.ProjectField)
					if fc214Obj.FinalValue != fieldsArr[fc056Loc[0][2]][fc056Loc[0][3]].FinalValue {
						if fc214Obj.FinalValue > fieldsArr[fc056Loc[0][2]][fc056Loc[0][3]].FinalValue {
							fieldsArr[fc056Loc[0][2]][fc056Loc[0][3]].FinalValue = fc214Obj.FinalValue
						}
					}
				}
			}
		}
	}
	//CSB0108RC0051000
	// 同一发票属性下，校验fc215与fc057的结果值，不一致时，取最早时间赋值给fc057
	// （当fc215录入值为A、?、？、空值时，不执行该校验）
	global.GLog.Info("CSB0108RC0051000")
	fc215Locs := fieldLocationMap["fc215"]
	for _, loc := range fc215Locs {
		fc215Obj := obj.Invoice[loc[0]].HospitalizationDate[loc[2]][loc[3]]
		compile, err := regexp.Compile("\\?|？|A")
		if err != nil {
			global.GLog.Error("CSB0108RC0051000" + err.Error())
			continue
		}
		if !(compile.MatchString(fc215Obj.ResultValue) || fc215Obj.ResultValue == "") {
			fc057Loc := utils2.GetFieldLoc(fieldLocationMap["fc057"], loc[0], -1, -1, -1)
			invoiceMap := obj.Invoice[fc057Loc[0][0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(fc057Loc[0][1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(fc057Loc[0][1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(fc057Loc[0][1]).Interface().([][]model3.ProjectField)
					if fc215Obj.FinalValue != fieldsArr[fc057Loc[0][2]][fc057Loc[0][3]].FinalValue {
						if fc215Obj.FinalValue < fieldsArr[fc057Loc[0][2]][fc057Loc[0][3]].FinalValue {
							fieldsArr[fc057Loc[0][2]][fc057Loc[0][3]].FinalValue = fc215Obj.FinalValue
						}
					}
				}
			}
		}
	}
	//CSB0108RC0052000
	// 当fc224的录入值不为整数时，结果值四舍五入取整数。
	global.GLog.Info("CSB0108RC0052000")
	fc224Locs := fieldLocationMap["fc224"]
	for _, loc := range fc224Locs {
		fc224Obj := obj.Invoice[loc[0]].HospitalizationFee[loc[2]][loc[3]]
		_, err := strconv.Atoi(fc224Obj.ResultValue)
		if err != nil {
			float, err := strconv.ParseFloat(fc224Obj.FinalValue, 64)
			if err != nil {
				global.GLog.Error("CSB0108RC0052000" + err.Error())
				continue
			}
			obj.Invoice[loc[0]].HospitalizationFee[loc[2]][loc[3]].FinalValue = strconv.Itoa(int(math.Floor(float + 0.5)))
		}
	}
	//CSB0108RC0053000
	// 当fc059结果值与fc009结果值不一致，fc059出问题件：病人姓名与申请表事故者不一致
	global.GLog.Info("CSB0108RC0053000")
	fc059Locs := fieldLocationMap["fc059"]
	fc009Locs = fieldLocationMap["fc009"]
	for _, loc := range fc059Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				if len(fc009Locs) > 0 {
					fc009Obj := obj.Invoice[fc009Locs[0][0]].Invoice[fc009Locs[0][2]][fc009Locs[0][3]]
					if fieldsArr[loc[2]][loc[3]].FinalValue != fc009Obj.FinalValue {
						issue := model3.Issue{
							Type:    "",
							Code:    "",
							Message: "病人姓名与申请表事故者不一致",
						}
						fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
					}
				}
			}
		}
	}
	codeMaps := map[string][]string{
		"fc043": {"fc001", "fc014", "fc045"},
		"fc065": {"fc002", "fc015", "fc082"},
		"fc212": {"fc013", "fc016", "fc083"},
		"fc040": {"fc090", "fc111"},
		"fc063": {"fc105", "fc138"},
		"fc107": {"fc106", "fc139"},
	}
	for key, codes := range codeMaps {
		// fmt.Println("-------------key---------------", key)
		for _, keyLoc := range fieldLocationMap[key] {
			sum := 0.0
			for _, code := range codes {
				codeLocs := utils2.GetFieldLoc(fieldLocationMap[code], keyLoc[0], -1, -1, -1)
				if len(codeLocs) > 0 {
					codeVal := utils2.GetFieldValueByLoc(obj, codeLocs[0], true)
					// fmt.Println("-------------codeVal---------------", code, codeVal)
					if utils2.RegIsMatch(`^[0-9\.]+$`, codeVal) {
						sum = utils2.SumFloat(sum, utils2.ParseFloat(codeVal), "+")
						// fmt.Println("-------------+++++++---------------", sum)
					}
				}

			}
			// fmt.Println("-------------sum---------------", utils2.ToString(sum))
			utils2.SetOnlyOneFinalValue(obj, keyLoc, utils2.ToString(sum))
		}
	}

	//CSB0108RC0042000
	// 将fc107与fc024结果值作对比（当存在多个fc107时，需将每个fc107的结果值相加求和再与fc024结果值作对比），
	// 不一致时第一个fc107出问题件：报销单统筹支付金额与申请书社保报销金额不一致
	global.GLog.Info("CSB0108RC0042000")
	fc107Locs := fieldLocationMap["fc107"]
	fc107Sum := decimal.NewFromFloat32(0.0)
	for _, loc := range fc107Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				d, err := decimal.NewFromString(fieldsArr[loc[2]][loc[3]].FinalValue)
				if err != nil {
					continue
				}
				fc107Sum = fc107Sum.Add(d)
			}
		}
	}
	fc024Locs := fieldLocationMap["fc024"]
	if len(fc024Locs) > 0 {
		fc024Obj := obj.Invoice[fc024Locs[0][0]].Invoice[fc024Locs[0][2]][fc024Locs[0][3]]
		if ParseFloat(fc024Obj.FinalValue) != ParseFloat(fc107Sum.String()) {
			if len(fc107Locs) > 0 {
				invoiceMap := obj.Invoice[fc107Locs[0][0]]
				eleLen := reflect.ValueOf(invoiceMap).NumField()
				if eleLen > 0 {
					if reflect.ValueOf(invoiceMap).Field(fc107Locs[0][1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(fc107Locs[0][1]).Kind() != reflect.Float64 {
						//每张发票每种类型的字段
						fieldsArr := reflect.ValueOf(invoiceMap).Field(fc107Locs[0][1]).Interface().([][]model3.ProjectField)
						issue := model3.Issue{
							Type:    "",
							Code:    "",
							Message: "报销单统筹支付金额与申请书社保报销金额不一致",
						}
						fieldsArr[fc107Locs[0][2]][fc107Locs[0][3]].Issues = append(fieldsArr[fc107Locs[0][2]][fc107Locs[0][3]].Issues, issue)
					}
				}
			}
		}
	}

	//CSB0108RC0043000
	// 将fc108与fc025结果值作对比（当存在多个fc108时，需将每个fc108的结果值相加求和再与fc025结果值作对比），
	// 不一致时第一个fc108出问题件：报销单其他扣除金额与申请书其他报销金额不一致
	global.GLog.Info("CSB0108RC0043000")
	fc108Locs := fieldLocationMap["fc108"]
	fc108Sum := decimal.NewFromFloat32(0.0)
	for _, loc := range fc108Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				d, err := decimal.NewFromString(fieldsArr[loc[2]][loc[3]].FinalValue)
				if err != nil {
					continue
				}
				fc108Sum = fc108Sum.Add(d)
			}
		}
	}
	fc025Locs = fieldLocationMap["fc025"]
	if len(fc025Locs) > 0 {
		fc025Obj := obj.Invoice[fc025Locs[0][0]].Invoice[fc025Locs[0][2]][fc025Locs[0][3]]
		fc025Value, err := strconv.ParseFloat(fc025Obj.FinalValue, 64)
		fc108Value, err1 := strconv.ParseFloat(fc108Sum.String(), 64)
		if err == nil && err1 == nil && fc025Value != fc108Value {
			if len(fc108Locs) > 0 {
				invoiceMap := obj.Invoice[fc108Locs[0][0]]
				eleLen := reflect.ValueOf(invoiceMap).NumField()
				if eleLen > 0 {
					if reflect.ValueOf(invoiceMap).Field(fc108Locs[0][1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(fc108Locs[0][1]).Kind() != reflect.Float64 {
						//每张发票每种类型的字段
						fieldsArr := reflect.ValueOf(invoiceMap).Field(fc108Locs[0][1]).Interface().([][]model3.ProjectField)
						issue := model3.Issue{
							Type:    "",
							Code:    "",
							Message: "报销单其他扣除金额与申请书其他报销金额不一致",
						}
						fieldsArr[fc108Locs[0][2]][fc108Locs[0][3]].Issues = append(fieldsArr[fc108Locs[0][2]][fc108Locs[0][3]].Issues, issue)
					}
				}
			}
		}
	}

	/*
		CSB0108RC0055000
		以下字段结果值保留两位小数（如：18则为18.00,18.5则为18.50，-18则为-18.00，0则为0.00）
		fc005、fc006、fc007、fc008、fc021、fc023、fc024、fc025、
		fc037、fc038、fc039、fc040、fc041、fc043、fc060、fc061、
		fc062、fc063、fc064、fc065、fc081、fc086、fc087、fc088、
		fc092、fc093、fc098、fc107、fc108、fc109、fc133、fc134、
		fc135、fc136、fc137、fc144、fc146、fc148、fc150、fc152、
		fc154、fc156、fc219、fc220、fc221、fc222、fc223
	*/
	global.GLog.Info("CSB0108RC0055000")
	codes = []string{
		"fc005", "fc006", "fc007", "fc008", "fc021", "fc023", "fc024", "fc025",
		"fc037", "fc038", "fc039", "fc040", "fc041", "fc043", "fc060", "fc061",
		"fc062", "fc063", "fc064", "fc065", "fc081", "fc086", "fc087", "fc088",
		"fc092", "fc093", "fc098", "fc107", "fc108", "fc109", "fc133", "fc134",
		"fc135", "fc136", "fc137", "fc144", "fc146", "fc148", "fc150", "fc152",
		"fc154", "fc156", "fc219", "fc220", "fc221", "fc222", "fc223"}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					finalValue, err := strconv.ParseFloat(fieldsArr[loc[2]][loc[3]].FinalValue, 64)
					if err != nil {
						//global.GLog.Error("CSB0108RC0055000" + err.Error())
						continue
					}
					value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", finalValue), 64)
					fieldsArr[loc[2]][loc[3]].FinalValue = strconv.FormatFloat(value, 'f', 2, 64)
				}
			}
		}
	}
	//CSB0108RC0056000
	// 当fc077结果值为0.00或为空时，fc077出问题件：发票、清单、结算单均无床位费或模糊、缺失
	global.GLog.Info("CSB0108RC0056000")
	if !RegIsMatch(obj.Bill.Agency, `^(00183000|00183010|00183002|00183012|00183300|00183310|00183301|00183311)$`) {
		fc077Locs := fieldLocationMap["fc077"]
		for _, loc := range fc077Locs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					if fieldsArr[loc[2]][loc[3]].FinalValue == "0.00" || fieldsArr[loc[2]][loc[3]].FinalValue == "" {
						issue := model3.Issue{
							Type:    "",
							Code:    "",
							Message: "发票、清单、结算单均无床位费或模糊、缺失",
						}
						fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
					}
				}
			}
		}
	}

	//CSB0108RC0058000
	// 当fc010与fc028结果值不一致时，fc028出问题件：账户姓名与申请人姓名不一致
	global.GLog.Info("CSB0108RC0058000")
	fc010Locs := fieldLocationMap["fc010"]
	fc028Locs := fieldLocationMap["fc028"]
	if len(fc010Locs) > 0 && len(fc028Locs) > 0 && len(fc010Locs) == len(fc028Locs) {
		for i := 0; i < len(fc010Locs); i++ {
			fc010Obj := obj.Invoice[fc010Locs[i][0]].Invoice[fc010Locs[i][2]][fc010Locs[i][3]]
			fc028Obj := obj.Invoice[fc028Locs[i][0]].Invoice[fc028Locs[i][2]][fc028Locs[i][3]]
			if fc010Obj.FinalValue != fc028Obj.FinalValue {
				issue := model3.Issue{
					Type:    "",
					Code:    "",
					Message: "账户姓名与申请人姓名不一致",
				}
				obj.Invoice[fc028Locs[i][0]].Invoice[fc028Locs[i][2]][fc028Locs[i][3]].Issues = append(fc028Obj.Issues, issue)
			}
		}
	}
	//CSB0108RC0059000
	// 当fc028录入值为A、为空时，将fc010结果值赋值给fc028
	global.GLog.Info("CSB0108RC0059000")
	fc028Locs = fieldLocationMap["fc028"]
	fc010Locs = fieldLocationMap["fc010"]
	for _, loc := range fc028Locs {
		fc028Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		if fc028Obj.ResultValue == "A" || fc028Obj.ResultValue == "" {
			if len(fc010Locs) > 0 {
				fc010Obj := obj.Invoice[fc010Locs[0][0]].Invoice[fc010Locs[0][2]][fc010Locs[0][3]]
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = fc010Obj.FinalValue
			}
		}
	}
	//CSB0108RC0060000
	// 当fc011与fc029结果值不一致时，fc029出问题件：账户证件号码与申请人证件号码不一致
	global.GLog.Info("CSB0108RC0060000")
	fc011Locs := fieldLocationMap["fc011"]
	fc029Locs := fieldLocationMap["fc029"]
	if len(fc011Locs) > 0 && len(fc029Locs) > 0 && len(fc011Locs) == len(fc029Locs) {
		for i := 0; i < len(fc011Locs); i++ {
			fc011Obj := obj.Invoice[fc011Locs[i][0]].Invoice[fc011Locs[i][2]][fc011Locs[i][3]]
			fc029Obj := obj.Invoice[fc029Locs[i][0]].Invoice[fc029Locs[i][2]][fc029Locs[i][3]]
			if fc011Obj.FinalValue != fc029Obj.FinalValue {
				issue := model3.Issue{
					Type:    "",
					Code:    "",
					Message: "账户证件号码与申请人证件号码不一致",
				}
				obj.Invoice[fc029Locs[i][0]].Invoice[fc029Locs[i][2]][fc029Locs[i][3]].Issues = append(fc029Obj.Issues, issue)
			}
		}
	}
	//CSB0108RC0061000
	// 当fc029录入值为A、为空时，将fc011结果值赋值给fc029
	global.GLog.Info("CSB0108RC0061000")
	fc029Locs = fieldLocationMap["fc029"]
	fc011Locs = fieldLocationMap["fc011"]
	for _, loc := range fc029Locs {
		fc029Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		if fc029Obj.ResultValue == "A" || fc029Obj.ResultValue == "" {
			if len(fc011Locs) > 0 {
				fc011Obj := obj.Invoice[fc011Locs[0][0]].Invoice[fc011Locs[0][2]][fc011Locs[0][3]]
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = fc011Obj.FinalValue
			}
		}
	}
	/*
		CSB0108RC0062000
		fc030、fc190录入值为《B0108_太平理赔_银行代码表》中“银行名称”（第二列），
		结果值根据常量表转换成对应的“银行代码”（第一列）
	*/
	global.GLog.Info("CSB0108RC0062000")
	codes = []string{"fc030", "fc190"}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					fieldsArr[loc[2]][loc[3]].FinalValue = constMap["yinHangDaiMaMap"][fieldsArr[loc[2]][loc[3]].ResultValue]
				}
			}
		}
	}
	//CSB0108RC0063000
	//当fc190与fc030结果值不一致时，fc030出问题件：账户所属银行不一致
	global.GLog.Info("CSB0108RC0063000")
	fc190Locs := fieldLocationMap["fc190"]
	fc030Locs := fieldLocationMap["fc030"]
	if len(fc190Locs) > 0 && len(fc030Locs) > 0 && len(fc190Locs) == len(fc030Locs) {
		for i := 0; i < len(fc190Locs); i++ {
			fc190Obj := obj.Invoice[fc190Locs[i][0]].Invoice[fc190Locs[i][2]][fc190Locs[i][3]]
			fc030Obj := obj.Invoice[fc030Locs[i][0]].Invoice[fc030Locs[i][2]][fc030Locs[i][3]]
			if fc190Obj.FinalValue != fc030Obj.FinalValue {
				issue := model3.Issue{
					Type:    "",
					Code:    "",
					Message: "账户所属银行不一致",
				}
				obj.Invoice[fc030Locs[i][0]].Invoice[fc030Locs[i][2]][fc030Locs[i][3]].Issues = append(fc030Obj.Issues, issue)
			}
		}
	}
	//CSB0108RC0064000
	// 当fc191与fc031结果值不一致时，fc031出问题件：银行账号不一致
	global.GLog.Info("CSB0108RC0064000")
	fc191Locs := fieldLocationMap["fc191"]
	fc031Locs := fieldLocationMap["fc031"]
	if len(fc191Locs) > 0 && len(fc031Locs) > 0 && len(fc191Locs) == len(fc031Locs) {
		for i := 0; i < len(fc191Locs); i++ {
			fc191Obj := obj.Invoice[fc191Locs[i][0]].Invoice[fc191Locs[i][2]][fc191Locs[i][3]]
			fc031Obj := obj.Invoice[fc031Locs[i][0]].Invoice[fc031Locs[i][2]][fc031Locs[i][3]]
			if fc191Obj.FinalValue != fc031Obj.FinalValue {
				issue := model3.Issue{
					Type:    "",
					Code:    "",
					Message: "银行账号不一致",
				}
				obj.Invoice[fc031Locs[i][0]].Invoice[fc031Locs[i][2]][fc031Locs[i][3]].Issues = append(fc031Obj.Issues, issue)
			}
		}
	}
	//CSB0108RC0065000
	// 当fc169录入值是B或包含?时，出问题件：代办人/代理人姓名全部模糊或部分模糊
	global.GLog.Info("CSB0108RC0065000")
	fc169Locs := fieldLocationMap["fc169"]
	for _, loc := range fc169Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				compile, err := regexp.Compile("\\?|？|B")
				if err != nil {
					global.GLog.Error("CSB0108RC0065000" + err.Error())
					continue
				}
				if compile.MatchString(fieldsArr[loc[2]][loc[3]].ResultValue) {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "代办人/代理人姓名全部模糊或部分模糊",
					}
					fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
				}
			}
		}
	}
	//CSB0108RC0066000
	// 当fc170录入值是B时，出问题件：代办人/代理人证件类型全部模糊或部分模糊
	global.GLog.Info("CSB0108RC0066000")
	fc170Locs := fieldLocationMap["fc170"]
	for _, loc := range fc170Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				if fieldsArr[loc[2]][loc[3]].ResultValue == "B" {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "代办人/代理人证件类型全部模糊或部分模糊",
					}
					fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
				}
			}
		}
	}
	//CSB0108RC0067000
	// 当fc171结果值为15位数字或18位数字时，fc170结果值默认为1并清空问题件
	global.GLog.Info("CSB0108RC0067000")
	fc171Locs := fieldLocationMap["fc171"]
	for _, loc := range fc171Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				if len(fieldsArr[loc[2]][loc[3]].FinalValue) == 15 || len(fieldsArr[loc[2]][loc[3]].FinalValue) == 18 {
					fc170Loc := utils2.GetFieldLoc(fieldLocationMap["fc170"], loc[0], -1, -1, -1)
					if len(fc170Loc) > 0 {
						obj.Invoice[fc170Loc[0][0]].Invoice[fc170Loc[0][2]][fc170Loc[0][3]].FinalValue = "1"
						obj.Invoice[fc170Loc[0][0]].Invoice[fc170Loc[0][2]][fc170Loc[0][3]].Issues = nil
					}
				}
			}
		}
	}
	//CSB0108RC0068000
	// 当fc171录入值是B或包含?时，出问题件：代办人/代理人有效证件号全部模糊或部分模糊
	global.GLog.Info("CSB0108RC0068000")
	fc171Locs = fieldLocationMap["fc171"]
	for _, loc := range fc171Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				compile, err := regexp.Compile("\\?|？|B")
				if err != nil {
					global.GLog.Error("CSB0108RC0068000" + err.Error())
					continue
				}
				if compile.MatchString(fieldsArr[loc[2]][loc[3]].ResultValue) {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "代办人/代理人有效证件号全部模糊或部分模糊",
					}
					fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
				}
			}
		}
	}
	//CSB0108RC0069000
	// 当fc172录入值是B或包含?时，出问题件：代办人/代理人联系电话全部模糊或部分模糊
	global.GLog.Info("CSB0108RC0069000")
	fc172Locs := fieldLocationMap["fc172"]
	for _, loc := range fc172Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				compile, err := regexp.Compile("\\?|？|B")
				if err != nil {
					global.GLog.Error("CSB0108RC0069000" + err.Error())
					continue
				}
				if compile.MatchString(fieldsArr[loc[2]][loc[3]].ResultValue) {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "代办人/代理人联系电话全部模糊或部分模糊",
					}
					fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
				}
			}
		}
	}
	//CSB0108RC0070000
	// 校验所有字段的结果值，当包含?或？时，将?、？转为^
	//CSB0108RC0241001 2023年10月24日10:06:17
	//校验fc080、fc143、fc145、fc147、fc149、fc151、fc153、fc155、fc169、fc171、fc172、fc009、fc010、fc028、fc033、fc059、fc216字段的结果值，当包含?或？时，将?、？转为^
	global.GLog.Info("CSB0108RC0070000")
	myCodes := "fc080、fc143、fc145、fc147、fc149、fc151、fc153、fc155、fc169、fc171、fc172、fc009、fc010、fc028、fc033、fc059、fc216"
	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票
		invoiceMap := obj.Invoice[i]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					for l := 0; l < len(fields); l++ {
						compile, err := regexp.Compile("\\?|？")
						if err != nil {
							global.GLog.Error("CSB0108RC0069000" + err.Error())
							continue
						}
						if strings.Index(myCodes, fields[l].Code) == -1 {
							continue
						}
						if compile.MatchString(fields[l].FinalValue) {
							fields[l].FinalValue = strings.Replace(fields[l].FinalValue, "?", "^", -1)
							fields[l].FinalValue = strings.Replace(fields[l].FinalValue, "？", "^", -1)
						}
					}
				}
			}
		}
	}
	/*
		CSB0108RC0071000
		销售渠道为“秒赔”的案件，清空以下所有字段的结果值及问题件
		fc009、fc019、fc099、fc010、fc011、fc095、fc012、fc027、
		fc028、fc029、fc030、fc031、fc158、fc112、fc169、fc170、
		fc171、fc172、fc173、fc110、fc190、fc191、fc174、fc175
	*/
	global.GLog.Info("CSB0108RC0071000")
	codes = []string{
		"fc009", "fc019", "fc099", "fc010", "fc011", "fc095", "fc012", "fc027",
		"fc028", "fc029", "fc030", "fc031", "fc158", "fc112", "fc169", "fc170",
		"fc171", "fc172", "fc173", "fc110", "fc190", "fc191", "fc174", "fc175",
	}
	if bill.SaleChannel == "秒赔" {
		for _, code := range codes {
			fieldLocs := fieldLocationMap[code]
			for _, loc := range fieldLocs {
				invoiceMap := obj.Invoice[loc[0]]
				eleLen := reflect.ValueOf(invoiceMap).NumField()
				if eleLen > 0 {
					if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
						//每张发票每种类型的字段
						fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
						fieldsArr[loc[2]][loc[3]].FinalValue = ""
						fieldsArr[loc[2]][loc[3]].Issues = nil
					}
				}
			}
		}
	}
	//CSB0108RC0072000
	// 销售渠道为“秒赔”的案件，清空fc033、fc059的问题件
	global.GLog.Info("CSB0108RC0072000")
	codes = []string{"fc033", "fc059"}
	if bill.SaleChannel == "秒赔" {
		for _, code := range codes {
			fieldLocs := fieldLocationMap[code]
			for _, loc := range fieldLocs {
				invoiceMap := obj.Invoice[loc[0]]
				eleLen := reflect.ValueOf(invoiceMap).NumField()
				if eleLen > 0 {
					if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
						//每张发票每种类型的字段
						fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
						fieldsArr[loc[2]][loc[3]].Issues = nil
					}
				}
			}
		}
	}
	/*
		CSB0108RC0073000
		当清空fc084（fc084为模板类型字段）的录入值时，清空对应生成的分块下的所有字段的结果值及问题件
	*/
	//global.GLog.Info("CSB0108RC0073000")
	//for i := 0; i < len(obj.Invoice); i++ {
	//	//同一发票
	//	invoiceMap := obj.Invoice[i]
	//	eleLen := reflect.ValueOf(invoiceMap).NumField()
	//	for j := 0; j < eleLen; j++ {
	//		if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
	//			//每张发票每种类型的字段
	//			fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
	//			for k := 0; k < len(fieldsArr); k++ {
	//				fields := fieldsArr[k]
	//				fc084Loc := utils2.GetFieldLoc(fieldLocationMap["fc084"], i, j, k, -1)
	//				fc084Obj := fieldsArr[fc084Loc[0][2]][fc084Loc[0][3]]
	//				if fc084Obj.ResultValue == "" {
	//					for l := 0; l < len(fields); l++ {
	//						if fields[l].BlockID != fc084Obj.BlockID && fields[l].BlockIndex != fc084Obj.BlockIndex {
	//							fields[l].FinalValue = ""
	//							fields[l].Issues = nil
	//						}
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	//CSB0108RC0057000
	// 当fc027录入值为1或A时，清空fc028、fc029、fc030、fc031的问题件。
	global.GLog.Info("CSB0108RC0057000")
	fc027Locs := fieldLocationMap["fc027"]
	codes = []string{"fc028", "fc029", "fc030", "fc031"}
	for _, loc := range fc027Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				if fieldsArr[loc[2]][loc[3]].ResultValue == "1" || fieldsArr[loc[2]][loc[3]].ResultValue == "A" {
					for _, code := range codes {
						fieldLocs := fieldLocationMap[code]
						for _, fieldLoc := range fieldLocs {
							obj.Invoice[fieldLoc[0]].Invoice[fieldLoc[2]][fieldLoc[3]].Issues = nil
						}
					}
				}
			}
		}
	}

	//CSB0108RC0245000
	//"fc012录入值为A时，
	//1.当fc009和fc010结果值一致时，fc012结果值默认为2
	//2.当fc009和fc010结果值不一致时，校验fc238的身份证第7-14位出生日期，
	//	①与当前时间作对比，如为未成年时，fc012结果值默认为5
	//	②与当前时间作对比，如为已成年时，fc012结果值默认为9"
	for _, fc012Loc := range fieldLocationMap["fc012"] {
		fc012Val := utils2.GetFieldValueByLoc(obj, fc012Loc, false)
		fc009Locs = utils2.GetFieldLoc(fieldLocationMap["fc009"], fc012Loc[0], -1, -1, -1)
		fc010Locs = utils2.GetFieldLoc(fieldLocationMap["fc010"], fc012Loc[0], -1, -1, -1)
		if utils2.RegIsMatch(`^(A)$`, fc012Val) && len(fc009Locs) == 1 && len(fc010Locs) == 1 {
			fc009Val := utils2.GetFieldValueByLoc(obj, fc009Locs[0], true)
			fc010Val := utils2.GetFieldValueByLoc(obj, fc010Locs[0], true)
			if fc009Val == fc010Val {
				utils2.SetOnlyOneFinalValue(obj, fc012Loc, "2")
			} else {
				fc238Locs := utils2.GetFieldLoc(fieldLocationMap["fc238"], fc012Loc[0], -1, -1, -1)
				if len(fc238Locs) != 1 {
					continue
				}
				fc238Val := utils2.GetFieldValueByLoc(obj, fc238Locs[0], true)
				if !IDCard(fc238Val) {
					continue
				}
				birthday, err := time.Parse("20060102", fc238Val[6:14])
				if err != nil {
					global.GLog.Error("CSB0108RC0245000", zap.Error(err))
					continue
				}
				if birthday.AddDate(18, 0, 0).Before(time.Now()) {
					utils2.SetOnlyOneFinalValue(obj, fc012Loc, "9")
				} else {
					utils2.SetOnlyOneFinalValue(obj, fc012Loc, "5")
				}
			}
		}
	}

	/*
		CSB0108RC0233000
		当fc189的结果值有内容时，屏蔽所有的问题件，增加该字段的问题件，每个分号（;或；）
		之间作为一个问题件描述，没有分号则视为一个问题件；(如：“申请书未填写；”则视为1个问题件；“申请书未填写；发票模糊”则视为2个问题件)
		(将该代码放在所有需求的最后面。)
	*/
	fc189Locs := fieldLocationMap["fc189"]
	for _, loc := range fc189Locs {
		fc189Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		if fc189Obj.FinalValue != "" {
			// for i := 0; i < len(obj.Invoice); i++ {
			// 	//同一发票
			// 	invoiceMap := obj.Invoice[i]
			// 	eleLen := reflect.ValueOf(invoiceMap).NumField()
			// 	for j := 0; j < eleLen; j++ {
			// 		if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
			// 			//每张发票每种类型的字段
			// 			fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
			// 			for k := 0; k < len(fieldsArr); k++ {
			// 				fields := fieldsArr[k]
			// 				for l := 0; l < len(fields); l++ {
			// 					fields[l].Issues = nil
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			issuesList := strings.Split(strings.ReplaceAll(fc189Obj.FinalValue, ";", "；"), "；")
			for _, s := range issuesList {
				issue := model3.Issue{
					Type:    "",
					Code:    "",
					Message: s,
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues, issue)
			}
		}
	}

	/*
		CSB0108RC0047000
		fc054录入值为A或?或？时，fc054结果值默认为“fc056结果值+序号”，并出问题件：住院发票：第[xxx]住院号模糊不清（xxx为fc054结果值）
		如：2020122601，如同一就诊日期（即存在多个fc056且结果值相同）有多张发票，则序号自动顺延，如：2020122602；

		2023年06月27日14:04:33
		fc054录入值为A或?或？时，fc054结果值默认为“fc056结果值+fc062结果值（去掉小数点.）”
	*/
	// global.GLog.Info("CSB0108RC0047000")
	fc054Num := 0
	for _, fc054Loc := range fieldLocationMap["fc054"] {
		fc054Val := utils2.GetFieldValueByLoc(obj, fc054Loc, false)
		if regexp.MustCompile(`^(A|？|\?)$`).MatchString(fc054Val) {
			fc054Num++
			utils2.SetOnlyOneFinalValue(obj, fc054Loc, "99999"+strconv.Itoa(fc054Num))
		}
	}
	// for _, fc054Loc := range fieldLocationMap["fc054"] {
	// 	fc054Val := utils2.GetFieldValueByLoc(obj, fc054Loc, false)
	// 	fc056Locs := utils2.GetFieldLoc(fieldLocationMap["fc056"], fc054Loc[0], -1, -1, -1)
	// 	fc062Locs := utils2.GetFieldLoc(fieldLocationMap["fc062"], fc054Loc[0], -1, -1, -1)
	// 	if regexp.MustCompile(`^(A|？|\?)$`).MatchString(fc054Val) && len(fc056Locs) == 1 && len(fc062Locs) == 1 {
	// 		fc056Val := utils2.GetFieldValueByLoc(obj, fc056Locs[0], true)
	// 		fc062Val := utils2.GetFieldValueByLoc(obj, fc062Locs[0], true)
	// 		fc062Val = utils2.RegReplace(fc062Val, `0+$`, "")
	// 		fc062Val = strings.Replace(fc062Val, ".", "", -1)
	// 		utils2.SetOnlyOneFinalValue(obj, fc054Loc, fc056Val+fc062Val)
	// 	}
	// }

	/*
		CSB0108RC0033000
		fc032录入值为A或?或？时，fc032结果值默认为“fc035结果值+序号”，并出问题件：门诊发票：第[xxx]门诊号模糊不清（xxx为fc032结果值）
		如：2020122601，如同一就诊日期（即存在多个fc035且结果值相同）有多张发票，则序号自动顺延，如：2020122602；

		2023年06月27日13:57:43
		fc032录入值为A或?或？时，fc032结果值默认为“fc035结果值+fc039结果值（去掉小数点.）”，
	*/
	global.GLog.Info("CSB0108RC0033000")
	fc032Num := 0
	for _, fc032Loc := range fieldLocationMap["fc032"] {
		fc032Val := utils2.GetFieldValueByLoc(obj, fc032Loc, false)
		if regexp.MustCompile(`^(A|？|\?)$`).MatchString(fc032Val) {
			fc032Num++
			utils2.SetOnlyOneFinalValue(obj, fc032Loc, "88888"+strconv.Itoa(fc032Num))
		}
	}
	// for _, fc032Loc := range fieldLocationMap["fc032"] {
	// 	fc032Val := utils2.GetFieldValueByLoc(obj, fc032Loc, false)
	// 	fc035Locs := utils2.GetFieldLoc(fieldLocationMap["fc035"], fc032Loc[0], -1, -1, -1)
	// 	fc039Locs := utils2.GetFieldLoc(fieldLocationMap["fc039"], fc032Loc[0], -1, -1, -1)
	// 	if regexp.MustCompile(`^(A|？|\?)$`).MatchString(fc032Val) && len(fc035Locs) == 1 && len(fc039Locs) == 1 {
	// 		fc035Val := utils2.GetFieldValueByLoc(obj, fc035Locs[0], true)
	// 		fc039Val := utils2.GetFieldValueByLoc(obj, fc039Locs[0], true)
	// 		fc039Val = utils2.RegReplace(fc039Val, `0+$`, "")
	// 		fc039Val = strings.Replace(fc039Val, ".", "", -1)
	// 		utils2.SetOnlyOneFinalValue(obj, fc032Loc, fc035Val+fc039Val)
	// 	}
	// }

	var wrong []string
	// 编码CSB0108RC0129001 导出校验 最后一页字段录入存在?？ 就导出校验  "最后一页存在?，请检查;" 20230713新增
	for p, invoiceMap := range obj.Invoice {
		for _, fields := range invoiceMap.Invoice {
			if p == len(obj.Invoice)-1 {
				for _, field := range fields {
					if strings.Index(obj.Bill.WrongNote, "最后一页存在?，请检查;") == -1 {
						if RegIsMatch(field.ResultValue, `[?？]`) {
							wrong = append(wrong, "最后一页存在?，请检查;")
						}
					}
				}
			}
		}
	}

	fc189Locs = fieldLocationMap["fc189"]
	for _, loc := range fc189Locs {
		fc189Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		if fc189Obj.ResultValue == "A" {
			for i := 0; i < len(obj.Invoice); i++ {
				//同一发票
				invoiceMap := obj.Invoice[i]
				eleLen := reflect.ValueOf(invoiceMap).NumField()
				for j := 0; j < eleLen; j++ {
					if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
						//每张发票每种类型的字段
						fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
						for k := 0; k < len(fieldsArr); k++ {
							fields := fieldsArr[k]
							for l := 0; l < len(fields); l++ {
								fields[l].Issues = nil
							}
						}
					}
				}
			}
		}
	}

	if len(wrong) > 1 {
		obj.Bill.WrongNote += "最后一页存在?，请检查;"
	}

	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票
		invoiceMap := obj.Invoice[i]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					for l := 0; l < len(fields); l++ {
						if !RegIsMatch(fields[l].Code, `^(fc228|fc229|fc230|fc231|fc232|fc233|fc234|fc235|fc236|fc237|fc249|fc250|fc251|fc252|fc253|fc254|fc255|fc256|fc257|fc258fc226|fc227)$`) {
							fields[l].FinalValue = RegReplace(fields[l].FinalValue, ` +`, "")
						}
					}
				}
			}
		}
	}
	// CSB0108RC0298000
	//1.统计所有fc039的结果值，将金额合计赋值给fc259，不存在fc039时，fc259结果值默认0
	//2.统计所有fc062的结果值，将金额合计赋值给fc260 不存在fc062时 ，fc260结果值默认0
	fieldCodes := [][]string{
		{"fc039", "fc259"},
		{"fc062", "fc260"},
	}
	for _, code := range fieldCodes {
		var exist bool
		for _, invoiceMap := range obj.Invoice {
			for _, fields := range invoiceMap.Invoice {
				exist = IsExist(fields, code[0])
			}
		}
		var sun float64
		for _, loc := range fieldLocationMap[code[0]] {
			val := utils2.GetFieldValueByLoc(obj, loc, true)
			if val != "" {
				sun = sun + ParseFloat(val)
			}
		}
		fieldRLocs := fieldLocationMap[code[1]]
		if len(fieldRLocs) > 0 && !exist {
			obj.Invoice[fieldRLocs[0][0]].Invoice[fieldRLocs[0][2]][fieldRLocs[0][3]].FinalValue = "0.00"
		}
		if len(fieldRLocs) > 0 {
			obj.Invoice[fieldRLocs[0][0]].Invoice[fieldRLocs[0][2]][fieldRLocs[0][3]].FinalValue = fmt.Sprintf("%.2f", sun)
		}
	}
	//CSB0108RC0300000 当fc217（字段唯一）结果值为空时，对比所有fc035、fc056、fc158不为空的结果值，将最早日期赋值给fc217（当不存在以上字段时不执行该校验）
	codes = []string{"fc035", "fc056", "fc158"}
	fc214LocArr := fieldLocationMap["fc217"]
	var fc217Val string
	if len(fc214LocArr) > 0 {
		fc217Val = utils2.GetFieldValueByLoc(obj, fc214LocArr[0], true)
	}
	if fc217Val == "" {
		isStorageVal := make([]string, 0)
		for _, code := range codes {
			exist := false
			for _, invoiceMap := range obj.Invoice {
				for _, fields := range invoiceMap.Invoice {
					isExist := IsExist(fields, code)
					if isExist {
						exist = true
					}
				}
			}
			for _, loc := range fieldLocationMap[code] {
				val := utils2.GetFieldValueByLoc(obj, loc, true)
				if exist && val != "" {
					isStorageVal = append(isStorageVal, val)
				}
			}
		}
		sort.Slice(isStorageVal, func(i, j int) bool {
			return isStorageVal[i] < isStorageVal[j]
		})
		if len(isStorageVal) > 0 {
			if len(fc214LocArr) > 0 {
				obj.Invoice[fc214LocArr[0][0]].Invoice[fc214LocArr[0][2]][fc214LocArr[0][3]].FinalValue = isStorageVal[0]
			}
		}
	}

	//CSB0108RC0318000
	fieldCodes = [][]string{
		{"fc036", "fc228"}, {"fc104", "fc229"}, {"fc058", "fc230"}, {"fc161", "fc231"}, {"fc162", "fc232"}, {"fc205", "fc233"}, {"fc206", "fc234"}, {"fc207", "fc235"}, {"fc208", "fc236"}, {"fc209", "fc237"}, {"fc239", "fc249"}, {"fc240", "fc250"}, {"fc241", "fc251"}, {"fc242", "fc252"}, {"fc243", "fc253"}, {"fc244", "fc254"}, {"fc245", "fc255"}, {"fc246", "fc256"}, {"fc247", "fc257"}, {"fc248", "fc258"},
	}
	cacheArrs := []string{}
	for ii, invoice := range obj.Invoice {
		// fmt.Println("----invoice.ZhenDuan-------", len(invoice.ZhenDuan))
		for jj, zhenDuan := range invoice.Invoice {
			for kk, field := range zhenDuan {
				for _, fcode := range fieldCodes {
					if fcode[0] == field.Code && field.ResultValue != "" && field.ResultValue != "B" {
						if field.ResultValue != "B" {
							obj.Invoice[ii].Invoice[jj] = setValue(obj.Invoice[ii].Invoice[jj], fcode[1], "NA")
						}
						// fmt.Println("----cacheArrscacheArrscacheArrs-------", field.Code, field.ResultValue)
						if arrays.Contains(cacheArrs, field.ResultValue) == -1 {
							cacheArrs = append(cacheArrs, field.ResultValue)
						} else {
							obj.Invoice[ii].Invoice[jj][kk].FinalValue = ""
							obj.Invoice[ii].Invoice[jj] = setValue(obj.Invoice[ii].Invoice[jj], fcode[1], "")
						}
					}
				}
			}
		}
	}

	fieldCodes = [][]string{{"fc261", "fc040"}, {"fc262", "fc063"}}
	for _, codes := range fieldCodes {
		// fmt.Println("-------------key---------------", key)
		if len(fieldLocationMap[codes[0]]) > 0 {
			sum := 0.0
			for _, keyLoc := range fieldLocationMap[codes[1]] {
				codeLocs := utils2.GetFieldLoc(fieldLocationMap[codes[1]], keyLoc[0], -1, -1, -1)
				if len(codeLocs) > 0 {
					codeVal := utils2.GetFieldValueByLoc(obj, codeLocs[0], true)
					// fmt.Println("-------------codeVal---------------", code, codeVal)
					if utils2.RegIsMatch(`^[0-9\.]+$`, codeVal) {
						sum = utils2.SumFloat(sum, utils2.ParseFloat(codeVal), "+")
						// fmt.Println("-------------+++++++---------------", sum)
					}
				}
			}
			if sum == 0.0 {
				sum = 0
			}
			utils2.SetOnlyOneFinalValue(obj, fieldLocationMap[codes[0]][0], utils2.ToString(sum))

		}

	}

	//CSB0108RC0323000
	codes = []string{"fc183", "fc184", "fc185", "fc186", "fc187"}
	for _, code := range codes {
		codeValue := ""
		for ii, loc := range fieldLocationMap[code] {
			if ii == 0 {
				codeValue = utils2.GetFieldValueByLoc(obj, loc, true)
			} else {
				utils2.SetOnlyOneFinalValue(obj, loc, codeValue)
			}
		}
	}

	//CSB0108RC0331000
	//当字段存在，将每个fc040、fc063的结果值相加之和与fc024结果值作对比，不一致时fc024出问题件：发票统筹支付金额与申请书社保报销金额录入不一致
	var locsValue float64 = 0.00
	fc024Locs = fieldLocationMap["fc024"]
	codes = []string{"fc040", "fc063"}
	for _, code := range codes {
		locs := fieldLocationMap[code]
		for _, loc := range locs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					finalValue, _ := strconv.ParseFloat(fieldsArr[loc[2]][loc[3]].FinalValue, 64)
					locsValue += finalValue
				}
			}
		}
	}
	for _, loc := range fc024Locs {
		invoiceMap := obj.Invoice[loc[0]]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
				locsFinal := fmt.Sprintf("%.2f", locsValue)
				if fieldsArr[loc[2]][loc[3]].FinalValue != locsFinal {
					issue := model3.Issue{
						Type:    "",
						Code:    "",
						Message: "发票统筹支付金额与申请书社保报销金额录入不一致",
					}
					fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
				}
			}
		}
	}

	////////////////////////项目报表 - 业务明细 - 单据状态
	obj.Bill.BillType = utils3.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	return nil, obj
}

func setValue(fields []model3.ProjectField, code, value string) []model3.ProjectField {
	for ii, field := range fields {
		if field.Code == code {
			fields[ii].FinalValue = value
			return fields
		}
	}
	return fields
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"yiYuanDaiMaMap", "B0108_太平理赔_医院代码表", "1", "0"},
		{"zhenDuanDaiMaMap", "B0108_太平理赔_诊断代码表", "1", "0"},
		{"shouShuDaiMaMap", "B0108_太平理赔_手术代码表", "1", "0"},
		{"yinHangDaiMaMap", "B0108_太平理赔_银行代码表", "1", "0"},
		{"faPiaoDaXiangMap2", "B0108_太平理赔_门诊发票大项代码表", "0", "0"},
		{"faPiaoDaXiangMap3", "B0108_太平理赔_住院发票大项代码表", "0", "0"},
	}
	for _, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])
				tempMap[strings.TrimSpace(arr[k])] = arr[v]
			}
		}
		constObj[item[0]] = tempMap
	}

	return constObj
}

func IDCard(id string) bool {
	if RegIsMatch(id, `^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\d{3}(\d|X|x)$`) || RegIsMatch(id, `^[1-9]\d{5}\d{2}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\d{2}\d{1}$`) {
		return true
	}
	return false
}
