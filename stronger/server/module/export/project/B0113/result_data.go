/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年03月22日16:35:47
 */

package B0113

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"server/global"
	model2 "server/module/export/model"
	utils2 "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

//B0113百年理赔对应关系
//
//fc002模板类型字段：
//1-申请书;2-发票;3-清单;4-报销单;8-报销单2;9-报销单3;10-报销单4;13-信息补充登记表
//
//和发票有对应关系的：2-发票;3-清单;4-报销单;8-报销单2;9-报销单3;10-报销单4;
//
//MB001-bc002  fc108发票属性
//MB001-bc003  fc109清单所属发票
//MB001-bc004  fc118报销单所属发票
//MB001-bc008  fc328报销单所属发票
//MB001-bc009  fc329报销单所属发票
//MB001-bc010  fc330报销单所属发票

//MB002-bc005 发票
//MB002-bc006 清单
//MB002-bc007 报销单1
//MB002-bc044 报销单2
//MB002-bc045 报销单3
//MB002-bc046 报销单4

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:          model2.TypeCode{FieldCode: []string{"fc108"}, BlockCode: []string{"bc005"}},
	QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc109"}, BlockCode: []string{"bc006"}},
	BaoXiaoDanFieldCode:       model2.TypeCode{FieldCode: []string{"fc118"}, BlockCode: []string{"bc007"}},
	ThirdBaoXiaoDan1FieldCode: model2.TypeCode{FieldCode: []string{"fc328"}, BlockCode: []string{"bc044"}},
	ThirdBaoXiaoDan2FieldCode: model2.TypeCode{FieldCode: []string{"fc329"}, BlockCode: []string{"bc045"}},
	ThirdBaoXiaoDan3FieldCode: model2.TypeCode{FieldCode: []string{"fc330"}, BlockCode: []string{"bc046"}},
	OtherTempType:             map[string]string{"1": "1", "13": "13"},
	TempTypeField:             "fc002",
	InvoiceNumField:           []string{"fc063"},
	MoneyField:                []string{"fc072"},
	//InvoiceTypeField: "fc003",
}

// ResultData B0113
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {
	defer func() {
		if err := recover(); err != nil {
			global.GLog.Error("", zap.Any("", err))
			debug.PrintStack()
		}
	}()
	obj = utils2.RelationDealManyInvoiceType(bill, blocks, fieldMap, relationConf)
	//常量
	constMap := constDeal(bill.ProCode)
	//特殊的常量
	constSpecialMap := constSpecialDeal(bill.ProCode)

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

	//CSB0113RC0038000
	//"录入内容为A或空时，不进行该校验（先执行该需求，再执行地址整合的需求）
	//以下字段：
	//1、判断第一列字段的结果数据最后一位是否为“省”，“直辖市”，“自治区”，“市”，“区”，“县”，“乡”，“村”，“镇”，“旗”，“盟”，“自治州”，不为这些时，在字段的结果数据后增加“区”；
	//2、判断第二列字段的结果数据最后一位是否为“省”，“直辖市”，“自治区”，“市”，“区”，“县”，“乡”，“村”，“镇”，“旗”，“盟”，“自治州”不为这些时，在字段的结果数据后增加“市”；
	//3、第三列字段中，判断字段的结果数据是否为“北京”，“天津”，“上海”或“重庆”，若是，则在结果数据后增加“市”；若不是，判断字段的结果数据最后一位是否为“省”，“直辖市”，“自治区”，“市”，“区”，“县”，“乡”，“村”，“镇”，“旗”，“盟”，“自治州”不为这些时，在字段的结果数据后增加“省”；
	//区     市     省
	//fc024  fc023  fc022"
	myCode := [][]string{
		{"fc024", "区"},
		{"fc023", "市"},
		{"fc022", "省"},
	}
	for i, code := range myCode {
		for _, loc := range fieldLocationMap[code[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			if utils.RegIsMatch("^(|A)$", val) {
				continue
			}
			if i == 2 && utils.RegIsMatch("^(北京|天津|上海|重庆)$", val) {
				utils.SetOnlyOneFinalValue(obj, loc, val+"市")
			} else if !utils.RegIsMatch("(省|直辖市|自治区|市|区|县|乡|村|镇|旗|盟|自治州)$", val) {
				utils.SetOnlyOneFinalValue(obj, loc, val+code[1])
			}
		}
	}

	//CSB0113RC0039000
	//"以下字段中，每一行字段按从左到右的顺序将字段的结果数据进行整合，赋值给最后一列字段后输出：（整张单字段唯一）
	//fc022，fc023，fc024，fc025"
	codeArr := []string{"fc022", "fc023", "fc024", "fc025"}
	locArr0 := fieldLocationMap[codeArr[0]]
	locArr1 := fieldLocationMap[codeArr[1]]
	locArr2 := fieldLocationMap[codeArr[2]]
	locArr3 := fieldLocationMap[codeArr[3]]
	if len(locArr0) == 1 && len(locArr1) == 1 && len(locArr2) == 1 {
		val0 := utils.GetFieldValueByLoc(obj, locArr0[0], true)
		val1 := utils.GetFieldValueByLoc(obj, locArr1[0], true)
		val2 := utils.GetFieldValueByLoc(obj, locArr2[0], true)
		val3 := utils.GetFieldValueByLoc(obj, locArr3[0], true)
		utils.SetFinalValue(obj, locArr3, val0+val1+val2+val3)
	}

	//CSB0113RC0040000
	//fc259、fc258、fc073结果数据四舍五入保留两位小数位。
	codeArr = []string{"fc259", "fc258", "fc073"}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			utils.SetOnlyOneFinalValue(obj, loc, utils.ParseDecimal(val).StringFixed(2))
		}
	}

	//CSB0113RC0042000
	//"当左边字段结果值为空时，将右边字段的结果值赋值为左边的字段，当其中一组字段中的某个字段不存在时，不进行赋值
	myCode = [][]string{
		{"fc268", "fc269"},
		{"fc270", "fc271"},
		{"fc272", "fc273"},
		{"fc274", "fc275"},
	}
	for _, code := range myCode {
		for _, leftLoc := range fieldLocationMap[code[0]] {
			rightLoc := utils.GetFieldLoc(fieldLocationMap[code[1]], leftLoc[0], -1, -1, -1)
			if len(rightLoc) == 0 {
				continue
			}
			leftVal := utils.GetFieldValueByLoc(obj, leftLoc, true)
			if leftVal == "" {
				rightVal := utils.GetFieldValueByLoc(obj, rightLoc[0], true)
				utils.SetOnlyOneFinalValue(obj, leftLoc, rightVal)
			}
		}
	}

	//CSB0113RC0043000
	//"以下字段中，第一列字段录入内容根据《B0113_百年理赔_百年理赔费用项目代码表》中“费用项目名称”（第二列）匹配，将对应的“费用项目编码”（第一列）赋值给第二列字段：
	myCode = [][]string{
		{"fc214", "fc204"},
		{"fc215", "fc205"},
		{"fc216", "fc206"},
		{"fc217", "fc207"},
		{"fc218", "fc208"},
		{"fc219", "fc209"},
		{"fc220", "fc210"},
		{"fc221", "fc211"},
		{"fc222", "fc212"},
		{"fc223", "fc213"},
		{"fc242", "fc247"},
		{"fc243", "fc248"},
		{"fc244", "fc249"},
		{"fc245", "fc250"},
		{"fc246", "fc251"},
	}
	for _, code := range myCode {
		for _, leftLoc := range fieldLocationMap[code[0]] {
			rightLoc := utils.GetFieldLoc(fieldLocationMap[code[1]], leftLoc[0], -1, -1, -1)
			leftVal := utils.GetFieldValueByLoc(obj, leftLoc, true)
			utils.SetFinalValue(obj, rightLoc, constMap["feiYongDaiMaMap"][leftVal])
		}
	}

	//CSB0113RC0044000
	//fc066,fc068,fc069为发票分块MB002-bc005中的字段，若fc066录入值为空，则fc066的结果值默认为fc068的结果值。
	//2023年06月29日09:21:55 取消
	//for _, fc066Loc := range fieldLocationMap["fc066"] {
	//	fc068Loc := utils.GetFieldLoc(fieldLocationMap["fc068"], fc066Loc[0], fc066Loc[1], -1, -1)
	//	if len(fc068Loc) == 0 {
	//		continue
	//	}
	//	fc066Val := utils.GetFieldValueByLoc(obj, fc066Loc, false)
	//	fc068Val := utils.GetFieldValueByLoc(obj, fc068Loc[0], true)
	//	if fc066Val == "" {
	//		utils.SetOnlyOneFinalValue(obj, fc066Loc, fc068Val)
	//	}
	//}

	//CSB0113RC0046000
	//"根据以下第一列字段（全部是循环分块MB002-bc006清单中的字段）的录入值匹配到常量表《B0113_百年理赔_全国》中“项目名称”（第一列）。
	//在校验下若第五列字段的结果值为CZ001或CZ002时，将第五列字段的结果值赋值给第二列字段进行结果值的判断；
	//1、第六列结果值为“2”时，第三列字段结果值默认为对应代码库中的自付比例（第四列），自付比例不足两位小数位时保留两位小数，小数位超过两位时四舍五入保留两位小数位；
	//2、第六列结果值为“3”时，第三列字段结果值默认为1.00；
	//3、第六列结果值为“5”时，执行以下校验：
	//（1）、若第五列字段的结果值为CZ001或CZ002时，将第五列字段的结果值赋值给第二列字段；
	//（2）、若对应的“分类”（第三列）等于“乙”时，第三列字段结果值默认为对应代码库中的自付比例（第四列），自付比例不足两位小数位时保留两位小数，小数位超过两位时四舍五入保留两位小数位；
	//（3）、若对应的“分类”（第三列）等于“丙”时，第三列字段结果值默认为1.00；
	//（4）、若对应的“分类”（第三列）等于“甲”时，清空对应的整条数据
	//（5）、若对应“自付比例”（第四列）等于0或0.00时，清空对应的整条数据
	myCode = [][]string{
		{"fc082", "fc074", "fc098", "fc090", "fc234", "fc260"},
		{"fc083", "fc075", "fc099", "fc091", "fc235", "fc261"},
		{"fc084", "fc076", "fc100", "fc092", "fc236", "fc262"},
		{"fc085", "fc077", "fc101", "fc093", "fc237", "fc263"},
		{"fc086", "fc078", "fc102", "fc094", "fc238", "fc264"},
		{"fc087", "fc079", "fc103", "fc095", "fc239", "fc265"},
		{"fc088", "fc080", "fc104", "fc096", "fc240", "fc266"},
		{"fc089", "fc081", "fc105", "fc097", "fc241", "fc267"},
	}
	for _, codes := range myCode {
		for _, loc1 := range fieldLocationMap[codes[0]] {
			loc2 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc3 := utils.GetFieldLoc(fieldLocationMap[codes[2]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc4 := utils.GetFieldLoc(fieldLocationMap[codes[3]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc5 := utils.GetFieldLoc(fieldLocationMap[codes[4]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc6 := utils.GetFieldLoc(fieldLocationMap[codes[5]], loc1[0], loc1[1], loc1[2], loc1[4])
			if len(loc2) < 1 || len(loc3) < 1 || len(loc4) < 1 || len(loc5) < 1 || len(loc6) < 1 {
				continue
			}
			resultVal1 := utils.GetFieldValueByLoc(obj, loc1, false)
			quanGuoItemArr, ok := constSpecialMap["quanGuoMap"][resultVal1]
			if !ok {
				//global.GLog.Error("CSB0113RC0046000", zap.Error(errors.New("没有匹配全国的常量")))
				continue
			}
			finalVal5 := utils.GetFieldValueByLoc(obj, loc5[0], true)
			if finalVal5 == "CZ001" || finalVal5 == "CZ002" {
				utils.SetOnlyOneFinalValue(obj, loc2[0], finalVal5)
			}
			resultVal6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
			if resultVal6 == "2" {
				utils.SetOnlyOneFinalValue(obj, loc3[0], utils.ParseDecimal(quanGuoItemArr[3]).StringFixed(2))
			} else if resultVal6 == "3" {
				utils.SetOnlyOneFinalValue(obj, loc3[0], "1.00")
			} else if resultVal6 == "5" {
				isDel := false
				if quanGuoItemArr[2] == "乙" {
					utils.SetOnlyOneFinalValue(obj, loc3[0], utils.ParseDecimal(quanGuoItemArr[3]).StringFixed(2))
				} else if quanGuoItemArr[2] == "丙" {
					utils.SetOnlyOneFinalValue(obj, loc3[0], "1.00")
				} else if quanGuoItemArr[2] == "甲" {
					isDel = true
				}
				if utils.ParseDecimal(quanGuoItemArr[3]).IsZero() {
					isDel = true
				}
				if isDel {
					utils.SetOnlyOneFinalValue(obj, loc1, "")
					utils.SetOnlyOneFinalValue(obj, loc2[0], "")
					utils.SetOnlyOneFinalValue(obj, loc3[0], "")
					utils.SetOnlyOneFinalValue(obj, loc4[0], "")
					utils.SetOnlyOneFinalValue(obj, loc5[0], "")
					utils.SetOnlyOneFinalValue(obj, loc6[0], "")
				}
			}

		}
	}

	//CSB0113RC0047000
	//"根据以下第二列字段（全部是循环分块MB002-bc006清单中的字段）的录入值匹配到常量表《B0113_百年理赔_全国》中“项目名称”（第一列）。
	//当第一列字段录入值不为2且不为5，且第二列字段的录入值不在常量库中时，根据第三列字段的录入值内容进行以下转换赋值给第四列字段。
	//1 = CZ002
	//2 = CZ001
	myMap := map[string]string{
		"1": "CZ002",
		"2": "CZ001",
	}
	//医保类型；项目名称；项目明细分类；项目编码；
	myCode = [][]string{
		{"fc260", "fc082", "fc234", "fc074"},
		{"fc261", "fc083", "fc235", "fc075"},
		{"fc262", "fc084", "fc236", "fc076"},
		{"fc263", "fc085", "fc237", "fc077"},
		{"fc264", "fc086", "fc238", "fc078"},
		{"fc265", "fc087", "fc239", "fc079"},
		{"fc266", "fc088", "fc240", "fc080"},
		{"fc267", "fc089", "fc241", "fc081"},
	}
	for _, codes := range myCode {
		for _, loc1 := range fieldLocationMap[codes[0]] {
			loc2 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc3 := utils.GetFieldLoc(fieldLocationMap[codes[2]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc4 := utils.GetFieldLoc(fieldLocationMap[codes[3]], loc1[0], loc1[1], loc1[2], loc1[4])
			if len(loc2) < 1 || len(loc3) < 1 || len(loc4) < 1 {
				continue
			}
			resultVal1 := utils.GetFieldValueByLoc(obj, loc1, false)
			resultVal2 := utils.GetFieldValueByLoc(obj, loc2[0], false)
			resultVal3 := utils.GetFieldValueByLoc(obj, loc3[0], false)
			_, ok := constSpecialMap["quanGuoMap"][resultVal2]
			if resultVal1 != "2" && resultVal1 != "5" && !ok {
				utils.SetOnlyOneFinalValue(obj, loc4[0], myMap[resultVal3])
			}
		}
	}

	//CSB0113RC0048000
	//"以下字段中，当第一列字段的结果值等于CZ001、CZ002，且第三列为数字时，
	//将第三列字段的录入值赋值给第四列字段，赋值后第三列字段的结果数据替换成第四列字段的结果值乘以第五列字段的结果值（在MB002中的bc006循环分块）
	//项目编码，项目名称，项目金额，费用总金额，自费比例
	myCode = [][]string{
		{"fc074", "fc082", "fc090", "fc110", "fc098"},
		{"fc075", "fc083", "fc091", "fc111", "fc099"},
		{"fc076", "fc084", "fc092", "fc112", "fc100"},
		{"fc077", "fc085", "fc093", "fc113", "fc101"},
		{"fc078", "fc086", "fc094", "fc114", "fc102"},
		{"fc079", "fc087", "fc095", "fc115", "fc103"},
		{"fc080", "fc088", "fc096", "fc116", "fc104"},
		{"fc081", "fc089", "fc097", "fc117", "fc105"},
	}
	for _, codes := range myCode {
		for _, loc1 := range fieldLocationMap[codes[0]] {
			loc3 := utils.GetFieldLoc(fieldLocationMap[codes[2]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc4 := utils.GetFieldLoc(fieldLocationMap[codes[3]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc5 := utils.GetFieldLoc(fieldLocationMap[codes[4]], loc1[0], loc1[1], loc1[2], loc1[4])
			if len(loc3) < 1 || len(loc4) < 1 || len(loc5) < 1 {
				continue
			}
			// r := regexp.MustCompile("^[\\d\\.]+$")
			finalVal1 := utils.GetFieldValueByLoc(obj, loc1, true)
			finalVal3 := utils.GetFieldValueByLoc(obj, loc3[0], true)
			if (finalVal1 == "CZ001" || finalVal1 == "CZ002") && RegIsMatch(finalVal3, `^[\d\.]+$`) {
				resultVal3 := utils.GetFieldValueByLoc(obj, loc3[0], false)
				utils.SetOnlyOneFinalValue(obj, loc4[0], resultVal3)

				decimalVal4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
				decimalVal5 := utils.GetFieldDecimalValueByLoc(obj, loc5[0], true)
				utils.SetOnlyOneFinalValue(obj, loc3[0], decimalVal4.Mul(decimalVal5).StringFixed(2))
			}
		}
	}

	//CSB0113RC0049000
	//"以下每一行为一组数据，
	//1、当第三列字段录入值为4时，将第四列字段的录入值赋值给第七列字段；赋值后第四列字段的结果数据替换成“第五列字段的结果值乘以第七列字段的结果值”。
	//2、当第三列字段录入值为1时，将第四列字段的录入值赋值给第七列字段；赋值后第六列字段的录入值赋值给第四列字段；第五列字段的结果值为“第六列字段的结果值除以第七列字段的结果值”。
	//（项目金额、自费比例、费用总金额、清单项目自付（先付、自费）金额的结果值保留两位小数）
	//项目名称；项目明细分类；医保类型；项目金额；自费比例；清单项目自付（先付、自费）金额；费用总金额
	//①     ②     ③     ④     ⑤     ⑥     ⑦
	myCode = [][]string{
		{"fc082", "fc234", "fc260", "fc090", "fc098", "fc308", "fc110"},
		{"fc083", "fc235", "fc261", "fc091", "fc099", "fc309", "fc111"},
		{"fc084", "fc236", "fc262", "fc092", "fc100", "fc310", "fc112"},
		{"fc085", "fc237", "fc263", "fc093", "fc101", "fc311", "fc113"},
		{"fc086", "fc238", "fc264", "fc094", "fc102", "fc312", "fc114"},
		{"fc087", "fc239", "fc265", "fc095", "fc103", "fc313", "fc115"},
		{"fc088", "fc240", "fc266", "fc096", "fc104", "fc314", "fc116"},
		{"fc089", "fc241", "fc267", "fc097", "fc105", "fc315", "fc117"},
	}
	for _, codes := range myCode {
		for _, loc3 := range fieldLocationMap[codes[2]] {
			loc4 := utils.GetFieldLoc(fieldLocationMap[codes[3]], loc3[0], loc3[1], loc3[2], loc3[4])
			loc5 := utils.GetFieldLoc(fieldLocationMap[codes[4]], loc3[0], loc3[1], loc3[2], loc3[4])
			loc6 := utils.GetFieldLoc(fieldLocationMap[codes[5]], loc3[0], loc3[1], loc3[2], loc3[4])
			loc7 := utils.GetFieldLoc(fieldLocationMap[codes[6]], loc3[0], loc3[1], loc3[2], loc3[4])
			if len(loc4) < 1 || len(loc5) < 1 || len(loc6) < 1 || len(loc7) < 1 {
				continue
			}
			resultVal3 := utils.GetFieldValueByLoc(obj, loc3, false)
			resultVal4 := utils.GetFieldValueByLoc(obj, loc4[0], false)
			resultVal6 := utils.GetFieldValueByLoc(obj, loc6[0], false)
			if resultVal3 == "4" || resultVal3 == "1" {
				utils.SetOnlyOneFinalValue(obj, loc7[0], resultVal4)
			}

			decimalVal5 := utils.GetFieldDecimalValueByLoc(obj, loc5[0], true)
			decimalVal7 := utils.GetFieldDecimalValueByLoc(obj, loc7[0], true)
			if resultVal3 == "4" {
				utils.SetOnlyOneFinalValue(obj, loc4[0], decimalVal5.Mul(decimalVal7).StringFixed(2))
			} else if resultVal3 == "1" {
				utils.SetOnlyOneFinalValue(obj, loc4[0], resultVal6)
				decimalVal6 := utils.GetFieldDecimalValueByLoc(obj, loc6[0], true)
				if decimalVal7.IsZero() {
					utils.SetOnlyOneFinalValue(obj, loc5[0], decimalVal7.StringFixed(2))
				} else {
					utils.SetOnlyOneFinalValue(obj, loc5[0], decimalVal6.Div(decimalVal7).StringFixed(2))
				}
			}
		}
	}

	//CSB0113RC0051000
	//"1、以下字段中（全部是循环分块MB002-bc006清单中的字段），若录入值不含小数位或不足两位小数位时，结果值保留两位小数；小数位超过两位时不做处理：
	//fc090,fc091,fc092,fc093,fc094,fc095,fc096,fc097,fc110,fc111,fc112,fc113,fc114,fc115,fc116,fc117
	//2、以下字段中（全部是循环分块MB002-bc005发票中的字段），若录入值不含小数位或不足两位小数位时，结果值保留两位小数；小数位超过两位时不做处理：
	//fc224,fc225,fc226,fc227,fc228,fc229,fc230,fc231,fc232,fc233,fc072,fc073,fc252,fc253,fc254,fc255,fc256
	//(录入内容为数字（包含小数点也是数字）才增加此校验)"
	codeArr = []string{"fc090", "fc091", "fc092", "fc093", "fc094", "fc095", "fc096", "fc097", "fc110", "fc111", "fc112", "fc113", "fc114", "fc115", "fc116", "fc117", "fc224", "fc225", "fc226", "fc227", "fc228", "fc229", "fc230", "fc231", "fc232", "fc233", "fc072", "fc073", "fc252", "fc253", "fc254", "fc255", "fc256"}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			if val == "" {
				continue
			}
			lenDecimal := lenDecimal(val)
			utils.SetOnlyOneFinalValue(obj, loc, val)
			if lenDecimal < 2 {
				utils.SetOnlyOneFinalValue(obj, loc, utils.ParseDecimal(val).StringFixed(2))
			}
		}
	}

	//CSB0113RC0050000
	//"以下每组字段录入内容包含“？”或“?”的时候，对左右两边字段的结果值都进行清空
	myCode = [][]string{
		{"fc082", "fc090"},
		{"fc083", "fc091"},
		{"fc084", "fc092"},
		{"fc085", "fc093"},
		{"fc086", "fc094"},
		{"fc087", "fc095"},
		{"fc088", "fc096"},
		{"fc089", "fc097"},
	}
	r := regexp.MustCompile("(\\?|？)")
	for _, codes := range myCode {
		for _, locLeft := range fieldLocationMap[codes[0]] {
			locRight := utils.GetFieldLoc(fieldLocationMap[codes[1]], locLeft[0], locLeft[1], locLeft[2], locLeft[4])
			if len(locRight) < 1 {
				continue
			}
			leftVal := utils.GetFieldValueByLoc(obj, locLeft, false)
			rightVal := utils.GetFieldValueByLoc(obj, locRight[0], false)
			if r.MatchString(leftVal) || r.MatchString(rightVal) {
				utils.SetOnlyOneFinalValue(obj, locLeft, "")
				utils.SetOnlyOneFinalValue(obj, locRight[0], "")
				fmt.Println("1")
			}
		}
	}

	//CSB0113RC0052000	结果数据
	//fc161录入1时，报问题件，问题件描述：无收据;
	//fc161录入2时，报问题件，问题件描述：收据模糊;
	//fc161录入3时，报问题件，问题件描述：清单不全;
	//fc161录入4时，报问题件，问题件描述：清单模糊;
	//fc161录入5时，报问题件，问题件描述：无清单
	//fc161录入6时，报问题件，问题件描述：案件丢失
	//fc161录入7时，报问题件，问题件描述：身份证过期
	//fc161录入8时，报问题件，问题件描述：缺少身份证
	//fc161录入9时，报问题件，问题件描述：其他
	//fc161录入a时，报问题件，问题件描述：领款人姓名与银行账户名不一致
	myArr := [][]string{
		{"1", "无收据"},
		{"2", "收据模糊"},
		{"3", "清单不全"},
		{"4", "清单模糊"},
		{"5", "无清单"},
		{"6", "案件丢失"},
		{"7", "身份证过期"},
		{"8", "缺少身份证"},
		{"9", "其他"},
		{"a", "领款人姓名与银行账户名不一致"},
	}
	for _, fc161Loc := range fieldLocationMap["fc161"] {
		fc161Val := utils.GetFieldValueByLoc(obj, fc161Loc, false)
		for _, item := range myArr {
			if item[0] == fc161Val {
				utils.SetIssue(obj, fc161Loc, item[1], "", "")
			}
		}
	}

	//CSB0113RC0053000
	//fc022、fc023、fc024和fc025录入内容都为A时，fc025增加问题件：受托人地址未填写内容。
	//CSB0113RC0054000
	//fc022、fc023、fc024、fc025录入内容包含?时，字段对应的结果数据转换为空，fc025清空结果值，fc025增加问题件：受托人地址填写内容无法识别。
	codeArr = []string{"fc022", "fc023", "fc024", "fc025"}
	locArr0 = fieldLocationMap[codeArr[0]]
	locArr1 = fieldLocationMap[codeArr[1]]
	locArr2 = fieldLocationMap[codeArr[2]]
	locArr3 = fieldLocationMap[codeArr[3]]
	if len(locArr0) == 1 && len(locArr1) == 1 && len(locArr2) == 1 && len(locArr3) == 1 {
		val0 := utils.GetFieldValueByLoc(obj, locArr0[0], false)
		val1 := utils.GetFieldValueByLoc(obj, locArr1[0], false)
		val2 := utils.GetFieldValueByLoc(obj, locArr2[0], false)
		val3 := utils.GetFieldValueByLoc(obj, locArr3[0], false)
		if val0 == val1 && val1 == val2 && val2 == val3 && val0 == "A" {
			utils.SetIssue(obj, locArr3[0], "受托人地址未填写内容", "", "")
		}
		if strings.Index(val0+val1+val2+val3, "?") != -1 {
			utils.SetFinalValue(obj, fieldLocationMap["fc025"], "")
			utils.SetIssues(obj, fieldLocationMap["fc025"], "受托人地址填写内容无法识别", "", "")
		}
	}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			if strings.Index(val, "?") != -1 {
				utils.SetOnlyOneFinalValue(obj, loc, "")
			}
		}
	}

	//CSB0113RC0055000
	//fc026录入内容包含?时，fc026输出空并增加问题件：受托人邮编填写内容无法识别。
	//CSB0113RC0057000
	//fc059录入内容包含?时，fc059输出空并增加问题件：银行账户名填写内容无法识别。
	//CSB0113RC0063000
	//"若fc057录入内容包含?时，将结果值转换为空并增加问题件：银行编码无法识别；支付方式无法识别。
	myArr = [][]string{
		{"fc026", "受托人邮编填写内容无法识别"},
		{"fc059", "银行账户名填写内容无法识别"},
		{"fc057", "银行编码无法识别;支付方式无法识别"},
	}
	for _, item := range myArr {
		for _, fc026Loc := range fieldLocationMap[item[0]] {
			fc026Val := utils.GetFieldValueByLoc(obj, fc026Loc, false)
			if strings.Index(fc026Val, "?") != -1 {
				utils.SetOnlyOneFinalValue(obj, fc026Loc, "")
				utils.SetIssue(obj, fc026Loc, item[1], "", "")
			}
		}
	}

	//CSB0113RC0056000
	//fc057录入内容为A时，fc057结果数据转换为空，同时增加问题件：银行编码未填写内容;支付方式未填写内容
	for _, fc057Loc := range fieldLocationMap["fc057"] {
		fc057Val := utils.GetFieldValueByLoc(obj, fc057Loc, false)
		if fc057Val == "A" {
			utils.SetOnlyOneFinalValue(obj, fc057Loc, "")
			utils.SetIssues(obj, fieldLocationMap["fc057"], "银行编码未填写内容;支付方式未填写内容", "", "")
		}
	}

	//CSB0113RC0058000
	//fc063录入内容为B或包含?时，结果数据默认生成一个唯一的账单号，出问题件：账单号XXX填写内容无法识别
	//CSB0113RC0059000
	//fc063录入内容为A时，结果数据默认生成一个唯一的账单号，出问题件：账单号XXX填写内容未填写
	for _, fc063Loc := range fieldLocationMap["fc063"] {
		fc063Val := utils.GetFieldValueByLoc(obj, fc063Loc, false)
		nextID := strconv.FormatInt(utils.GWorker.NextId(), 10)
		if fc063Val == "B" || strings.Index(fc063Val, "?") != -1 {
			utils.SetOnlyOneFinalValue(obj, fc063Loc, nextID)
			utils.SetIssue(obj, fc063Loc, "账单号"+nextID+"填写内容无法识别", "", "")
		}
		if fc063Val == "A" {
			utils.SetOnlyOneFinalValue(obj, fc063Loc, nextID)
			utils.SetIssue(obj, fc063Loc, "账单号"+nextID+"填写内容未填写", "", "")
		}
	}

	//CSB0113RC0041000
	//fc327为循环分块MB002-bc007报销单中的字段，fc327字段的结果值默认为该分块对应的发票的账单号fc063的结果值
	for _, loc := range fieldLocationMap["fc327"] {
		fco63Loc := utils.GetFieldLoc(fieldLocationMap["fc063"], loc[0], -1, -1, -1)
		if len(fco63Loc) != 1 {
			continue
		}
		fco63Val := utils.GetFieldValueByLoc(obj, fco63Loc[0], true)
		utils.SetOnlyOneFinalValue(obj, loc, fco63Val)
	}

	//CSB0113RC0045000
	//fc257为循环分块MB002-bc006清单中的字段，fc257字段的结果值默认为该分块对应的发票的账单号fc063的结果值
	for _, fc257Loc := range fieldLocationMap["fc257"] {
		fc063Loc := utils.GetFieldLoc(fieldLocationMap["fc063"], fc257Loc[0], -1, -1, -1)
		if len(fc063Loc) == 0 {
			continue
		}
		fc063Val := utils.GetFieldValueByLoc(obj, fc063Loc[0], true)
		utils.SetOnlyOneFinalValue(obj, fc257Loc, fc063Val)
	}

	//CSB0113RC0060000
	//将fc065录入值根据《B0113_百年理赔_百年理赔医院代码表》中“医院名称”（第二列）对应的“医院编码”（第一列）进行转码，无法转码时默认为空并增加问题件：医院名称无法识别；医院代码无法识别。
	//CSB0113RC0061000
	//所有字段结果数据包含?时，将结果数据替换为空，出问题件：（字段名）填写内容模糊。fc260、fc261、fc262、fc263、fc264、fc265、fc266、fc267不执行这个校验
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
						if fields[l].Code == "fc065" {
							val, ok := constMap["yiYuanDaiMaMap"][fields[l].ResultValue]
							fields[l].FinalValue = val
							if !ok {
								fields[l].FinalValue = ""
								utils.SetIssue(obj, []int{i, j, k, l}, "医院名称无法识别；医院代码无法识别", "", "")
							}
						}
						if fields[l].Issues == nil && strings.Index(fields[l].ResultValue, "?") != -1 && strings.Index("fc260、fc261、fc262、fc263、fc264、fc265、fc266、fc267", fields[l].Code) == -1 {
							fields[l].FinalValue = ""
							utils.SetIssue(obj, []int{i, j, k, l}, fields[l].Name+"填写内容模糊", "", "")
						}
					}
				}
			}
		}
	}

	//CSB0113RC0062000
	//以下字段录入值含?时，清空字段对应的问题件：
	//fc022,fc023,fc024,fc082,fc090,fc083,fc091,fc084,fc092,fc085,fc093,fc086,fc094,fc087,fc095,fc088,fc096,fc089,fc097"
	//fc270,fc272,fc274,fc276,fc278,fc280,fc214,fc215,fc216,fc217,fc218,fc219,fc220,fc221,fc222,fc223,fc242,fc243,fc244,fc245,fc246,fc224,fc225,fc226,fc227,fc228,fc229,fc230,fc231,fc232,fc233,fc252,fc253,fc254,fc255,fc256,fc234,fc235,fc236,fc237,fc238,fc239,fc240,fc241,fc333,fc334,fc335,fc336
	codeArr = []string{"fc022", "fc023", "fc024", "fc082", "fc090", "fc083", "fc091", "fc084", "fc092", "fc085", "fc093", "fc086", "fc094", "fc087", "fc095", "fc088", "fc096", "fc089", "fc097",
		"fc098", "fc099", "fc100", "fc101", "fc102", "fc103", "fc104", "fc105", "fc308", "fc309", "fc310", "fc311", "fc312", "fc313", "fc314", "fc315",
		"fc270", "fc272", "fc274", "fc276", "fc278", "fc280", "fc214", "fc215", "fc216", "fc217", "fc218", "fc219", "fc220", "fc221", "fc222", "fc223", "fc242", "fc243", "fc244", "fc245", "fc246", "fc224", "fc225", "fc226", "fc227", "fc228", "fc229", "fc230", "fc231", "fc232", "fc233", "fc252", "fc253", "fc254", "fc255", "fc256", "fc234", "fc235", "fc236", "fc237", "fc238", "fc239", "fc240", "fc241", "fc333", "fc334", "fc335", "fc336"}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			if strings.Index(val, "?") != -1 {
				utils.DelOnlyOneIssue(obj, loc)
			}
		}
	}

	//CSB0113RC0113000
	//当fc331录入值为1时，清空fc332、fc333、fc334、fc335、fc336、fc337字段的结果值
	codes := []string{"fc332", "fc333", "fc334", "fc335", "fc336", "fc337"}
	for _, loc := range fieldLocationMap["fc331"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		if val == "1" {
			for _, code := range codes {
				loc1 := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
				utils.SetFinalValue(obj, loc1, "")
			}
		}
	}
	////////////////////////项目报表 - 业务明细 - 单据状态
	obj.Bill.BillType = utils2.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	return nil, obj
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"feiYongDaiMaMap", "B0113_百年理赔_百年理赔费用项目代码表", "1", "0"},
		{"yiYuanDaiMaMap", "B0113_百年理赔_百年理赔医院代码表", "1", "0"},
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

// 初始化常量
func constSpecialDeal(proCode string) map[string]map[string][]string {
	constObj := make(map[string]map[string][]string, 0)
	nameMap := [][]string{
		{"quanGuoMap", "B0113_百年理赔_全国", "0"},
	}
	for _, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string][]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				tempMap[strings.TrimSpace(arr[k])] = arr
			}
		}
		constObj[item[0]] = tempMap
	}
	return constObj
}

// lenDecimal 获取小数长度
func lenDecimal(numStr string) int {
	parts := strings.Split(numStr, ".")
	if len(parts) < 2 {
		return 0 // 整数，返回 0
	}
	return len(parts[1])
}
