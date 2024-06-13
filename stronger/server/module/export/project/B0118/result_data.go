/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/24 10:50 上午
 */

package B0118

import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	model2 "server/module/export/model"
	"server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	unitFunc "server/module/unit"
	utils2 "server/utils"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
)

// RelationConf 关系对应
var RelationConf = model2.RelationConf{
	InvoiceFieldCode:          model2.TypeCode{FieldCode: []string{"fc059"}, BlockCode: []string{"bc001"}},
	QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc060"}, BlockCode: []string{"bc002"}},
	BaoXiaoDanFieldCode:       model2.TypeCode{FieldCode: []string{"fc061"}, BlockCode: []string{"bc003"}},
	InvoiceDaXiangFieldCode:   model2.TypeCode{FieldCode: []string{"fc201"}, BlockCode: []string{"bc007"}},
	ThirdBaoXiaoDan1FieldCode: model2.TypeCode{FieldCode: []string{"fc213"}, BlockCode: []string{"bc008"}},
	ThirdBaoXiaoDan2FieldCode: model2.TypeCode{FieldCode: []string{"fc214"}, BlockCode: []string{"bc009"}},
	ThirdBaoXiaoDan3FieldCode: model2.TypeCode{FieldCode: []string{"fc215"}, BlockCode: []string{"bc010"}},
	OtherTempType:             map[string]string{"5": "5", "6": "6", "7": "7"},
	TempTypeField:             "fc056",
	InvoiceNumField:           []string{"fc002"},
	MoneyField:                []string{"fc008"},
	InvoiceTypeField:          "fc003",
	//2-发票;3-清单;4-报销单;5-疾病诊断;6-手术；7-姓名；8-发票大项；9-第三方报销1；10-第三方报销2；11-第三方报销3；
}

//账单号：fc002
//总金额：fc008

// ResultData B0118
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {

	//defer func() {
	//	if r := recover(); r != nil {
	//		global.GLog.Error(fmt.Sprintf("%v", r))
	//	}
	//}()

	global.GLog.Info("B0118:::ResultData")
	obj = utils.RelationDealManyInvoiceType(bill, blocks, fieldMap, RelationConf)

	//常量
	constMap := constDeal(bill.ProCode)
	constSpecialMap := constSpecialDeal(bill.ProCode)
	//fmt.Println(constMap)
	obj.Bill.CountMoney = 0
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

	countErrInvoice := 0
	oneErrIndexArr := []int{-1, -1}
	tt := 1
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

						//CSB0118RC0218000
						//fc053、fc182、fc183、fc184、fc185、fc186、fc187、fc188、fc189、fc190
						//的结果值根据录入内容匹配《中意理赔疾病诊断》的“代码”栏内容
						fieldsCode1 := []string{"fc053", "fc182", "fc183", "fc184", "fc185", "fc186", "fc187", "fc188", "fc189", "fc190"}
						for _, s := range fieldsCode1 {
							if s == fields[l].Code {
								v, ok := constMap["jiBingZhenDuanMap"][fields[l].ResultValue]
								if ok {
									fields[l].FinalValue = v
								}
							}
						}

						//CSB0118RC0217000
						//当左列字段结果值为空时，将右列字段的结果数据赋值给左列字段
						fieldsCode := [][]string{
							{"fc053", "fc191"}, {"fc182", "fc192"}, {"fc183", "fc193"}, {"fc184", "fc194"}, {"fc185", "fc195"},
							{"fc186", "fc196"}, {"fc187", "fc197"}, {"fc188", "fc198"}, {"fc189", "fc199"}, {"fc190", "fc200"}}
						for _, codeArr := range fieldsCode {
							if fields[l].Code == codeArr[0] && fields[l].FinalValue == "" {
								fields[l].FinalValue = getValInFieldsByCodeIndex(fields, codeArr[1], fields[l].BlockIndex)
							}
						}

						//CSB0118RC0219000
						//当fc002录入值包含“?”“？”时，清空结果值，报问题件，问题件编码：N0004，描述账单号模糊
						if fields[l].Code == "fc002" {
							compile, err := regexp.Compile("\\?|？")
							if err != nil {
								global.GLog.Error("CSB0118RC0219000" + err.Error())
								continue
							}
							if compile.MatchString(fields[l].ResultValue) {
								fields[l].FinalValue = ""
								issue := model3.Issue{
									//FieldId: fields[l].ID,
									Type:    "",
									Code:    "N0004",
									Message: "账单号模糊",
								}
								fields[l].Issues = append(fields[l].Issues, issue)
							}

							//CSB0118RC0220000
							//当fc002（会存在多个）的结果值为空时，且存在相同属性的fc146的结果值为空时，
							//将fc002的结果值，默认赋值为一组默认数据，规则为：
							//字母+字母+字母+字母+数字（字母从24位字母的第一位开始，数字从1开始，如：AA1）
							if fields[l].FinalValue == "" {
								//locs := getFieldLoc(fieldLocationMap["fc146"], i, -1, k, -1)
								//if len(locs) > 0 {
								//	//fc146Field := reflect.ValueOf(invoiceMap).Field(locs[0][1]).Interface().([]model3.ProjectField)[locs[0][2]]
								//	fc146Field := invoiceMap.BaoXiaoDan[k][locs[0][3]]
								//	if fc146Field.FinalValue == "" {
								fields[l].FinalValue = string(65) + string(65) + string(65) + string(65) + strconv.Itoa(tt)
								tt++
								//	}
								//}
							}

							//CSB0118RC0221000
							//同一发票属性中，fc002账单号和fc008发票总金额录入内容同时包含?时，且以下字段中录入内容包含?时，
							//将同一发票中的fc003发票费用类型的结果数据清空，
							//并在第一个被清空的fc003字段所对应同一发票属性的fc002账单号出问题件：N张票据模糊；
							//(N为清空的fc003发票费用类型的数量；问题件编码为N0001；被fc003清空的同一属性发票中的所有问题件不需要导出）
							codeArr := []string{"fc009", "fc010", "fc011", "fc012", "fc013", "fc014", "fc015", "fc016", "fc017", "fc018", "fc019", "fc020", "fc021", "fc022", "fc023", "fc024", "fc025", "fc026", "fc027", "fc028", "fc029", "fc030", "fc031", "fc032", "fc033", "fc034", "fc035", "fc036", "fc037", "fc038", "fc039", "fc040", "fc041", "fc042", "fc043", "fc044", "fc045", "fc046", "fc047", "fc048"}
							if strings.Index(fields[l].ResultValue, "?") != -1 {
								locs := getFieldLoc(fieldLocationMap["fc008"], i, -1, -1, -1)
								if len(locs) > 0 &&
									hasQuestionMark(codeArr, fieldLocationMap, i, invoiceMap, locs) {
									fc003locs := getFieldLoc(fieldLocationMap["fc003"], i, -1, -1, -1)
									if len(fc003locs) > 0 {
										invoiceMap.Invoice[0][fc003locs[0][3]].FinalValue = ""

										//清空同一发票问题件
										for o := 0; o < eleLen; o++ {
											if reflect.ValueOf(invoiceMap).Field(o).Kind() == reflect.Slice {
												fieldsO := reflect.ValueOf(invoiceMap).Field(o).Interface().([][]model3.ProjectField)
												for p := 0; p < len(fieldsO); p++ {
													fields1 := fieldsO[p]
													for q := 0; q < len(fields1); q++ {
														fields1[q].Issues = nil
													}
												}
											}
										}

										//模糊数量
										countErrInvoice++
										//记录第一次的i,l
										if oneErrIndexArr[0] == -1 || oneErrIndexArr[1] == -1 {
											oneErrIndexArr[0] = i
											oneErrIndexArr[1] = l
										}
									}
								}
							}
						}

						//CSB0118RC0222000
						//当存在fc101字段时、且当fc101录入为1、2时，将同属性的fc054的结果值赋值给fc079
						rc, err := regexp.Compile("^(1|2)$")
						if err != nil {
							global.GLog.Error("CSB0118RC0222000" + err.Error())
							continue
						}
						if fields[l].Code == "fc101" && rc.MatchString(fields[l].ResultValue) {
							fc054Locs := getFieldLoc(fieldLocationMap["fc054"], i, -1, -1, -1)
							fc079Locs := getFieldLoc(fieldLocationMap["fc079"], i, -1, -1, -1)
							if len(fc054Locs) > 0 && len(fc079Locs) > 0 {
								invoiceMap.BaoXiaoDan[0][fc079Locs[0][3]].FinalValue = invoiceMap.Invoice[0][fc054Locs[0][3]].FinalValue
							}
						}

						//CSB0118RC0223000
						//当fc054的结果值内容不为《中意理赔医院代码表》“名称”栏内容时，出导出校验：医院录入错误，不在数据库中时，应录入“本代码表中不存在的其他医院”
						if fields[l].Code == "fc054" {
							_, ok := constMap["yiYuanDaiMaMap"][fields[l].FinalValue]
							if !ok {
								obj.Bill.WrongNote += "医院录入错误，不在数据库中时，应录入“本代码表中不存在的其他医院”;"
							}
						}

						//CSB0118RC0224000
						//fc005、fc006、fc007录入内容为6位数时，在第一位数前补20，并将格式更改为YYYY-MM-DD，如2019-01-01
						rc, err = regexp.Compile("^(fc005|fc006|fc007|fc280)$")
						if err != nil {
							global.GLog.Error("CSB0118RC0224000" + err.Error())
							continue
						}
						if rc.MatchString(fields[l].Code) && len(fields[l].ResultValue) == 6 {
							fields[l].FinalValue = "20" + fields[l].ResultValue[:2] + "-" + fields[l].ResultValue[2:4] + "-" + fields[l].ResultValue[4:]
						}

						//CSB0118RC0225000
						//1、当同一发票属性的所以的fc153合计金额不等于fc008的结果值时，
						//	fc008报问题件，编码N0005，问题件描述：清单金额与总金额不一致
						//2、当fc153不存在时，同一发票属性的fc008与fc010、fc012、fc014、fc016、fc018、fc020、fc022、fc024、fc026、fc028、fc030、fc032、fc034、fc036、fc038、fc040、fc042、fc044、fc046、fc048汇总金额不一致时，
						//	fc008报问题件，编码N0005，问题件描述：清单金额与总金额不一致
						//3、同时存在bc002、bc007时，同一发票属性的fc008与fc010、fc012、fc014、fc016、fc018、fc020、fc022、fc024、fc026、fc028、fc030、fc032、fc034、fc036、fc038、fc040、fc042、fc044、fc046、fc048汇总金额不一致时，
						//fc008报问题件，编码N0005，问题件描述：清单金额与总金额不一致
						if fields[l].Code == "fc008" {
							isQuestion := false
							fc008, err := decimal.NewFromString(fields[l].FinalValue)
							if err != nil {
								global.GLog.Error("CSB0118RC0225000" + err.Error())
								fc008 = decimal.Zero
								//continue
							}
							// var allFields []model3.ProjectField
							t := decimal.NewFromFloat(0.0)
							fc153s := []string{}
							for _, projectFields := range invoiceMap.QingDan {
								fc153 := GetFieldsFinal(projectFields, "fc153")
								if fc153 == "" || fc153 == "0" {
									continue
								}
								sss := fc153 + "_" + GetFieldsFinal(projectFields, "fc152")
								if arrays.Contains(fc153s, sss) == -1 {
									d, _ := decimal.NewFromString(fc153)
									t = t.Add(d)
								}
								fc153s = append(fc153s, sss)
								// allFields = append(allFields, projectFields...)
							}
							// t := getAllVal(allFields, "fc153")
							//fmt.Println(fc008.Float64())
							//fmt.Println(t.Float64())
							// if !t.Equals(decimal.Zero) && !t.Equals(fc008) && len(invoiceMap.QingDan) > 0 {
							// 	isQuestion = true
							// }

							// if len(invoiceMap.QingDan) < 1 || (len(invoiceMap.InvoiceDaXiang) > 0 && len(invoiceMap.QingDan) > 0) {
							fArr, _ := regexp.Compile("^(fc010|fc012|fc014|fc016|fc018|fc020|fc022|fc024|fc026|fc028|fc030|fc032|fc034|fc036|fc038|fc040|fc042|fc044|fc046|fc048)$")
							t = decimal.NewFromFloat32(0.0)
							for _, projectFields := range invoiceMap.InvoiceDaXiang {
								for _, invoiceDaXiangField := range projectFields {
									if fArr.MatchString(invoiceDaXiangField.Code) {
										d, _ := decimal.NewFromString(invoiceDaXiangField.FinalValue)
										t = t.Add(d)
										//fmt.Println(t)
									}
								}
							}
							//fmt.Println(fc008.Float64())
							//fmt.Println(t.Float64())
							if !t.Equals(fc008) {
								isQuestion = true
							}
							// }

							if isQuestion {
								issue := model3.Issue{
									Type:    "",
									Code:    "N0005",
									Message: "清单金额与总金额不一致",
								}
								fields[l].Issues = append(fields[l].Issues, issue)
							}
						}

						//CSB0118RC0226000
						//当同一发票属性只有bc001跟bc007，不存在bc002时，
						//fc003的录入值为1时，
						//出问题件，编码N0005，问题件描述：门诊发票无结算单和清单明细；
						//fc003的录入值为2时，
						//出问题件，编码N0005，问题件描述：住院发票无结算单和清单明细；
						if fields[l].Code == "fc003" &&
							len(invoiceMap.Invoice) > 0 &&
							len(invoiceMap.InvoiceDaXiang) > 0 &&
							len(invoiceMap.QingDan) < 1 && len(invoiceMap.BaoXiaoDan) < 1 {
							mes := ""
							if fields[l].ResultValue == "1" {
								mes = "门诊发票无结算单和清单明细"
							}
							if fields[l].ResultValue == "2" {
								mes = "住院发票无结算单和清单明细"
							}
							if mes != "" {
								issue := model3.Issue{
									Type:    "",
									Code:    "N0005",
									Message: mes,
								}
								fields[l].Issues = append(fields[l].Issues, issue)
							}
						}

						//CSB0118RC0227000
						//fc055结果值根据录入值对应《中意理赔手术编码》的“名称”匹配《中意理赔手术编码》的“代码”栏内容
						if fields[l].Code == "fc055" {
							fields[l].FinalValue = constMap["shouShuBianMaMap"][fields[l].ResultValue]
						}

						//CSB0118RC0228000
						//fc009、fc011、fc013、fc015、fc017、fc019、fc021、fc023、fc025、fc027、fc029、fc031、fc033、fc035、fc037、fc039、fc041、fc043、fc045、fc047的结果值，
						//根据同一发票属性的fc003的录入值进行匹配转换，
						//当fc003录入为1时，
						//根据录入内容匹配《B0118_中意理赔_住院费用类型》中“费用名称”对应转换为“费用代码”，匹配不上时，转换为199；
						//当fc003录入为2时，
						//根据录入内容匹配《B0118_中意理赔_住院费用类型》中“费用名称”对应转换为“费用代码”，匹配不上时，转换为299
						//2022年04月18日20:07:19 第一个门诊第二个住院
						codeArr := []string{"fc009", "fc011", "fc013", "fc015", "fc017", "fc019", "fc021", "fc023", "fc025", "fc027", "fc029", "fc031", "fc033", "fc035", "fc037", "fc039", "fc041", "fc043", "fc045", "fc047"}
						for _, code := range codeArr {
							if fields[l].Code == code && fields[l].FinalValue != "" {
								fc003locs := getFieldLoc(fieldLocationMap["fc003"], i, -1, -1, -1)
								constName := "zhuYuanFeiYongMap"
								if len(fc003locs) > 0 && invoiceMap.Invoice[0][fc003locs[0][3]].ResultValue == "1" {
									constName = "menZhenFeiYongLeiXingMap"
								}
								val, ok := constMap[constName][fields[l].ResultValue]
								if !ok {
									if len(fc003locs) > 0 && invoiceMap.Invoice[0][fc003locs[0][3]].ResultValue == "1" {
										val = "199"
									}
									if len(fc003locs) > 0 && invoiceMap.Invoice[0][fc003locs[0][3]].ResultValue == "2" {
										val = "299"
									}
								}
								fields[l].FinalValue = val
							}
						}

						//CSB0118RC0254000
						//CSB0118RC0256000
						//CSB0118RC0229000
						//CSB0118RC0230000
						iMap := map[string]string{
							"fc008": "发票总金额模糊",
							//"fc209": "自费金额模糊",
							//"fc210": "自费金额模糊",
							"fc051": "第三方实际赔付金额模糊",
							"fc052": "社保实际赔付金额模糊",
							"fc067": "社保自费金额模糊",
							"fc102": "第三方自费金额模糊",
						}
						if v, ok := iMap[fields[l].Code]; ok {
							compile, err := regexp.Compile("\\?|？")
							if err != nil {
								global.GLog.Error("CSB0118RC0254000" + err.Error())
								//continue
							}
							if compile.MatchString(fields[l].ResultValue) {
								issue := model3.Issue{
									Type:    "",
									Code:    "N0002",
									Message: v,
								}
								fields[l].Issues = append(fields[l].Issues, issue)
								if fields[l].Code == "fc008" {
									fields[l].FinalValue = ""
								}
								if fields[l].Code == "fc051" || fields[l].Code == "fc052" {
									fields[l].FinalValue = "0"
								}
							}
						}

						//CSB0118RC0232000
						//下列字段为循环分块中的字段，
						//1.当第一列字段的录入值为2或6时，第二列对应的字段的结果值默认为1；
						//	若fc181出险地点省为广西壮族且第三列字段匹配不上常量库，且第一列字段录入值为2时，第二列字段结果值默认为20%；
						//2.当第一列字段的录入值为3时，第二列对应的字段的结果值根据第三列字段匹配到 根据同一发票属性的fc181的结果值匹配对应的名称的常量库，
						//  如fc181结果值为北京，则匹配《B0118_中意理赔-北京》常量表中的“中文名称”对应的“自付比例”带出，
						//  当fc181的结果值为空时，则匹配《中意理赔_中意理赔_全国》常量表中的“中文名称”对应的“自付比例”带出，
						//  如项目名称对应多条“自付比例”时，匹配自付比例最小值的那一项（当常量表中的“自付比例”为空时，默认为“0”）；
						//3.当第一列字段的录入值为4时，第二列的字段的结果为第五列除以第四列的值（只保留两位小数点）
						//4.当第一列字段的录入值为5时，第二列对应的字段的结果值根据第三列字段匹配到根据同一发票属性的fc180或fc181的结果值匹配对应的名称的常量库，
						//如fc180结果值为北京，则匹配《B0118_中意理赔-北京》常量表中的“中文名称”对应的“自付比例”带出，
						//当fc180结果值为空，或结果值没有对应的常量表时，则匹配fc181的结果值，
						//如fc181结果值为北京，则匹配《B0118_中意理赔-北京》常量表中的“中文名称”对应的“自付比例”带出，
						//当fc181的结果值为空，或结果值没有对应的常量表时，则匹配《中意理赔全国》常量表中的“中文名称”对应的“自付比例”带出，
						//如项目名称对应多条“自付比例”时，匹配自付比例最小值的那一项；
						//若第三列字段匹配的常量库中，字段对应的自付比例为0或1时，第二列字段结果值默认为5%，
						//否则第二列字段结果值为常量库中对应的自付比例；
						//若第三列字段匹配不上常量库时，第二列字段结果值默认为5%，fc181出险地点省为黑龙江时，第二列字段结果值默认为20%；（
						//默认为5%的字段出问题件：项目名称类型为乙类无扣费比例按照5%进行扣费；默认为20%的字段出问题件：项目名称类型为乙类无扣费比例按照20%进行扣费；）
						//（问题件存在多条时，只需报一条即可）；fc181出险地点省为广西壮族时，第二列字段结果值默认为10%；
						//（该需求放在需求“下列字段为循环分块中的字段，根据fc059-发票属性和fc060-清单所属发票进行匹配…”的前面
						myCode := [][]string{
							{"fc154", "fc162", "fc084", "fc092", "fc172"},
							{"fc155", "fc163", "fc085", "fc093", "fc173"},
							{"fc156", "fc164", "fc086", "fc094", "fc174"},
							{"fc157", "fc165", "fc087", "fc095", "fc175"},
							{"fc158", "fc166", "fc088", "fc096", "fc176"},
							{"fc159", "fc167", "fc089", "fc097", "fc177"},
							{"fc160", "fc168", "fc090", "fc098", "fc178"},
							{"fc161", "fc169", "fc091", "fc099", "fc179"},
						}
						hasIssue := false
						for _, i3 := range myCode {
							if fields[l].Code == i3[0] {
								rc, err = regexp.Compile("^(2|6)$")
								if err != nil {
									global.GLog.Error("CSB0118RC0232000" + err.Error())
									//continue
								}
								fc181Loc := getFieldLoc(fieldLocationMap["fc181"], i, -1, -1, -1)
								fc180Loc := getFieldLoc(fieldLocationMap["fc180"], i, -1, -1, -1)
								fc003Loc := getFieldLoc(fieldLocationMap["fc003"], i, -1, -1, -1)
								fc162Locs := getFieldLoc(fieldLocationMap[i3[1]], i, -1, k, fields[l].BlockIndex)
								for _, fc162Loc := range fc162Locs {
									fc084Loc := getFieldLoc(fieldLocationMap[i3[2]], i, -1, k, fc162Loc[4])
									fc172Loc := getFieldLoc(fieldLocationMap[i3[4]], i, -1, k, fields[l].BlockIndex)
									fc092Loc := getFieldLoc(fieldLocationMap[i3[3]], i, -1, k, fields[l].BlockIndex)
									if len(fc181Loc) < 1 || len(fc084Loc) < 1 || len(fc172Loc) < 1 || len(fc092Loc) < 1 || len(fc180Loc) < 1 {
										continue
									}
									fc084Val := invoiceMap.QingDan[k][fc084Loc[0][3]].FinalValue
									fc172Val := invoiceMap.QingDan[k][fc172Loc[0][3]].FinalValue
									fc092Val := invoiceMap.QingDan[k][fc092Loc[0][3]].FinalValue

									fc181Val := invoiceMap.Invoice[0][fc181Loc[0][3]].ResultValue
									fc180Val := invoiceMap.Invoice[0][fc180Loc[0][3]].ResultValue
									fc003Val := invoiceMap.Invoice[0][fc003Loc[0][3]].ResultValue
									invoiceType := "住院"
									if fc003Val == "1" {
										invoiceType = "门诊"
									}
									if rc.MatchString(fields[l].ResultValue) {
										invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "1"
										if fc181Val == "广西壮族" {
											_, ok := constMap["广西壮族"][fc084Val]
											if !ok && fields[l].ResultValue == "2" {
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "0.2"
											}
										}
									}

									if fields[l].ResultValue == "3" {
										//invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "1"
										ok := false
										addr, ok := constSpecialMap[fc180Val]
										province := fc180Val
										if !ok {
											addr, ok = constSpecialMap[fc181Val]
											province = fc181Val
											if !ok {
												province = "全国"
												addr, ok = constSpecialMap[province]
											}
										}

										if ok {
											valArr, ok1 := addr[fc084Val]
											headArr := constSpecialMap[province+"_head"]["head"]
											ii := arrays.ContainsString(headArr, province+invoiceType)
											global.GLog.Info("3 ii", zap.Any(province+invoiceType, ii))
											if ii > 0 && ii < len(valArr) {
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = valArr[ii]
											}
											if fc180Val != "" {
												ii = arrays.ContainsString(headArr, fc180Val+invoiceType)
												global.GLog.Info("3 ii", zap.Any(fc180Val+invoiceType, ii))
												if ii > 0 && ii < len(valArr) {
													invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = valArr[ii]
												}
											}
											if !ok1 {
												fields[l].FinalValue = ""
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = ""
												invoiceMap.QingDan[k][fc084Loc[0][3]].FinalValue = ""
												invoiceMap.QingDan[k][fc092Loc[0][3]].FinalValue = ""
												invoiceMap.QingDan[k][fc172Loc[0][3]].FinalValue = ""
											}
										}
										//if invoiceMap.QingDan[k][fc162Loc[3]].FinalValue == "" {
										//	invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "0"
										//}
									}

									if fields[l].ResultValue == "4" {
										//第五列除以第四列的
										if fc092Val != "" && fc172Val != "" && fc092Val != "0" {
											fc092, err := decimal.NewFromString(fc092Val)
											fc172, err1 := decimal.NewFromString(fc172Val)
											if err == nil && err1 == nil {
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = fc172.DivRound(fc092, 3).String()
											}
										}
									}

									if fields[l].ResultValue == "5" {
										//2023年04月25日13:41:14 泽如姐说要和第二点一样
										ok := false
										addr, ok := constSpecialMap[fc180Val]
										province := fc180Val
										if !ok {
											addr, ok = constSpecialMap[fc181Val]
											province = fc181Val
											if !ok {
												province = "全国"
												addr, ok = constSpecialMap[province]
											}
										}
										if ok {
											valArr, _ := addr[fc084Val]
											headArr := constSpecialMap[province+"_head"]["head"]
											ii := arrays.ContainsString(headArr, province+invoiceType)
											global.GLog.Info("5 ii", zap.Any(province+invoiceType, ii))
											if ii > 0 && ii < len(valArr) {
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = valArr[ii]
											}
											if fc180Val != "" {
												ii = arrays.ContainsString(headArr, fc180Val+invoiceType)
												global.GLog.Info("5 ii", zap.Any(fc180Val+invoiceType, ii))
												if ii > 0 && ii < len(valArr) {
													invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = valArr[ii]
												}
											}

											if invoiceMap.QingDan[k][fc162Loc[3]].FinalValue == "0" || invoiceMap.QingDan[k][fc162Loc[3]].FinalValue == "1" {
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "0.05"
											}
											if valArr == nil {
												invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "0.05"
												if fc181Val == "黑龙江" {
													invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "0.2"
												}
												if fc181Val == "广西壮族" {
													invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = "0.1"
												}
												//默认为5%的字段出问题件：项目名称类型为乙类无扣费比例按照5%进行扣费；默认为20%的字段出问题件：项目名称类型为乙类无扣费比例按照20%进行扣费；
												if !hasIssue && (invoiceMap.QingDan[k][fc162Loc[3]].FinalValue == "0.05" || invoiceMap.QingDan[k][fc162Loc[3]].FinalValue == "0.2") {
													fc162Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc162Loc[3]].FinalValue)
													issue := model3.Issue{
														Type:    "",
														Code:    "N0005",
														Message: "项目名称类型为乙类无扣费比例按照" + fc162Val.Mul(decimal.NewFromInt(100)).String() + "%进行扣费",
													}
													fields[l].Issues = append(fields[l].Issues, issue)
													hasIssue = true
												}
											}
										}

										//val2, ok2 := addr[fc084Val]
										//invoiceMap.QingDan[k][fc162Loc[3]].FinalValue = val2
										//2023年02月15日09:30:05 泽如删掉
										//2023年02月24日15:18:56 泽如还原

									}

								}
							}
						}

						//CSB0118RC0233000
						//下列字段为循环分块中的字段：
						//若fc181出险地点省为广西壮族，第三列字段录入值匹配《B0118_中意理赔_广西壮族》常量表中的“中文名称”所对应的“统计类别”为材料费时，则
						//当第四列对应的字段录入值处于区间[0,200]时，第二列对应的字段结果值默认为0；
						//当第四列对应的字段录入值处于区间(200,500]时，第二列对应的字段结果值默认为10%；
						//当第四列对应的字段录入值大于500时，第二列对应的字段结果值默认为20%；
						//项目类型，自费比例，项目名称，项目金额，自费金额
						iCodes := [][]string{
							{"fc154", "fc162", "fc084", "fc092", "fc172"},
							{"fc155", "fc163", "fc085", "fc093", "fc173"},
							{"fc156", "fc164", "fc086", "fc094", "fc174"},
							{"fc157", "fc165", "fc087", "fc095", "fc175"},
							{"fc158", "fc166", "fc088", "fc096", "fc176"},
							{"fc159", "fc167", "fc089", "fc097", "fc177"},
							{"fc160", "fc168", "fc090", "fc098", "fc178"},
							{"fc161", "fc169", "fc091", "fc099", "fc179"},
						}
						for _, codes := range iCodes {
							if fields[l].Code == codes[1] {
								fc181Loc := getFieldLoc(fieldLocationMap["fc181"], i, -1, -1, -1)
								fc084Locs := getFieldLoc(fieldLocationMap[codes[2]], i, -1, k, -1)
								if len(fc181Loc) > 0 && invoiceMap.Invoice[0][fc181Loc[0][3]].ResultValue == "广西壮族" {
									for _, fc084Loc := range fc084Locs {
										fc092Loc := getFieldLoc(fieldLocationMap[codes[3]], i, -1, k, fields[l].BlockIndex)
										category, ok := constMap["guangXiZhuangZuMap"][invoiceMap.QingDan[k][fc084Loc[3]].ResultValue]
										if ok && category == "材料费" {
											four := invoiceMap.QingDan[k][fc092Loc[0][3]].ResultValue
											fFour, err := strconv.ParseFloat(four, 64)
											if err != nil {
												global.GLog.Error("CSB0118RC0233000" + err.Error())
												continue
											}
											if fFour >= 0 && fFour <= 200 {
												fields[l].FinalValue = "0"
											}
											if fFour > 200 && fFour <= 500 {
												fields[l].FinalValue = "0.1"
											}
											if fFour > 500 {
												fields[l].FinalValue = "0.2"
											}
										}
									}
								}
							}
						}

						//CSB0118RC0235000
						//fc103为循环分块bc002中的字段，fc103字段的结果值默认为该分块对应的发票的账单号fc002的结果值
						if fields[l].Code == "fc103" {
							fc002locs := getFieldLoc(fieldLocationMap["fc002"], i, -1, -1, -1)
							if len(fc002locs) > 0 {
								fields[l].FinalValue = invoiceMap.Invoice[0][fc002locs[0][3]].FinalValue
							}
						}

						//CSB0118RC0236000
						//fc150是bc002循环分块中的字段，fc150字段的结果值根据同一发票属性的fc152的结果值
						if fields[l].Code == "fc150" {
							fc152locs := getFieldLoc(fieldLocationMap["fc152"], i, -1, -1, -1)
							if len(fc152locs) > 0 {
								fields[l].FinalValue = invoiceMap.QingDan[k][fc152locs[0][3]].FinalValue
							}
						}

						//CSB0118RC0237000
						//fc151是bc002循环分块中的字段，fc151字段的结果值为该清单分块中所以fc092、fc093、fc094、fc095、fc096、fc097、fc098、fc099的合计
						if fields[l].Code == "fc151" {
							fArr := []string{"fc092", "fc093", "fc094", "fc095", "fc096", "fc097", "fc098", "fc099"}
							t := decimal.NewFromFloat32(0.0)
							for _, s := range fArr {
								var qingDanFields []model3.ProjectField
								for _, f := range invoiceMap.QingDan {
									qingDanFields = append(qingDanFields, f...)
								}
								t = t.Add(getAllVal(qingDanFields, s))
							}
							fields[l].FinalValue = t.String()
						}

						//CSB0118RC0238000
						//fc152根据同一发票属性的fc003的结果值为1时，fc152的结果值根据《B0118_中意理赔_门诊费用类型》的“费用名称”转换为“费用代码”，匹配不上时，转换为199
						//当fc152对应同一发票属性的fc003的结果值为2时，fc152的结果值根据《B0118_中意理赔_门诊费用类型》的“费用名称”转换为“费用代码”，匹配不上时，转换为299
						//改为住院费用2022年04月18日16:19:45
						if fields[l].Code == "fc152" {
							fc003locs := getFieldLoc(fieldLocationMap["fc003"], i, -1, -1, -1)
							if len(fc003locs) > 0 {
								constName := "zhuYuanFeiYongMap"
								if invoiceMap.Invoice[0][fc003locs[0][3]].FinalValue == "1" {
									constName = "menZhenFeiYongLeiXingMap"
								}
								val, ok := constMap[constName][fields[l].ResultValue]
								if !ok {
									if invoiceMap.Invoice[0][fc003locs[0][3]].FinalValue == "1" {
										val = "199"
									}
									if invoiceMap.Invoice[0][fc003locs[0][3]].FinalValue == "2" {
										val = "299"
									}
								}
								fields[l].FinalValue = val
							}
						}

						//CSB0118RC0239000
						//当fc152录入值为“工本费、病历费、卡费、复印费、陪护费、陪人费”，将fc152、fc153的录入值赋值给相同属性的
						// iCodes = [][]string{
						// 	{"fc104", "fc105"},
						// }
						// rc, err = regexp.Compile("^(工本费|病历费|卡费|复印费|陪护费|陪人费)$")
						// if err != nil {
						// 	global.GLog.Error("CSB0118RC0239000" + err.Error())
						// 	continue
						// }
						// if fields[l].Code == "fc152" && rc.MatchString(fields[l].ResultValue) {
						// 	for _, code := range iCodes {
						// 		fc152locs := getFieldLoc(fieldLocationMap["fc152"], i, j, k, -1)
						// 		fc153locs := getFieldLoc(fieldLocationMap["fc153"], i, j, k, -1)
						// 		if len(fc152locs) > 0 && len(fc153locs) > 0 {
						// 			fc152Val := invoiceMap.QingDan[k][fc152locs[0][3]].ResultValue
						// 			fc153Val := invoiceMap.QingDan[k][fc153locs[0][3]].ResultValue
						// 			fc104locs := getFieldLoc(fieldLocationMap[code[0]], i, j, k, -1)
						// 			fc105locs := getFieldLoc(fieldLocationMap[code[1]], i, j, k, -1)
						// 			for _, loc := range fc104locs {
						// 				invoiceMap.QingDan[k][loc[3]].FinalValue = fc152Val
						// 			}
						// 			for _, loc := range fc105locs {
						// 				invoiceMap.QingDan[k][loc[3]].FinalValue = fc153Val
						// 			}
						// 		}
						// 	}
						// }

						//CSB0118RC0240000
						//当fc152结果值为199或299时，且左一列字段结果值包含膳食费、陪床费时，将相同属性的左二列的结果值赋值给右二列，每四个字段为一组，只需要根据每组的第一个字段结果值是否包含食费、陪床费进行判断；
						iCodes = [][]string{
							{"fc084", "fc092", "fc104", "fc105"},
							{"fc085", "fc093", "fc106", "fc107"},
							{"fc086", "fc094", "fc108", "fc109"},
							{"fc087", "fc095", "fc110", "fc111"},
							{"fc088", "fc096", "fc112", "fc113"},
							{"fc089", "fc097", "fc114", "fc115"},
							{"fc090", "fc098", "fc116", "fc117"},
							{"fc091", "fc099", "fc118", "fc119"},
						}
						if fields[l].Code == "fc152" && (fields[l].FinalValue == "199" || fields[l].FinalValue == "299") {
							rc, err = regexp.Compile("^(膳食费|陪床费)$")
							if err != nil {
								global.GLog.Error("CSB0118RC0240000" + err.Error())
								continue
							}
							// for _, codes := range iCodes {
							// 	firstLocs := getFieldLoc(fieldLocationMap[codes[0]], i, j, k, -1)
							// 	for _, loc := range firstLocs {
							// 		twoLocs := getFieldLoc(fieldLocationMap[codes[1]], i, j, k, loc[4])
							// 		threeLocs := getFieldLoc(fieldLocationMap[codes[2]], i, j, k, loc[4])
							// 		fourLocs := getFieldLoc(fieldLocationMap[codes[3]], i, j, k, loc[4])
							// 		if len(twoLocs) > 0 && len(threeLocs) > 0 && rc.MatchString(invoiceMap.QingDan[k][loc[3]].FinalValue) {
							// 			invoiceMap.QingDan[k][threeLocs[0][3]].FinalValue = invoiceMap.QingDan[k][firstLocs[0][3]].FinalValue
							// 			invoiceMap.QingDan[k][fourLocs[0][3]].FinalValue = invoiceMap.QingDan[k][twoLocs[0][3]].FinalValue
							// 		}
							// 	}
							// }
						}

						//CSB0118RC0241000,在CSB0118RC0234000需求上
						//下列字段为循环分块中的字段，根据fc059-发票属性和fc060-清单所属发票进行匹配：
						//fc171的值是所有左一列乘以左二列的值的汇总
						//（当清单循环分块中的fc150的结果值存在重复的情况下，将重复的fc150内容该清单属性的内容汇总到第一个fc150对应的清单分块中的fc171；
						//	当一组数据两个字段的值都有空的情况下，则不进行计算，
						//	公式为：fc171=fc092*fc162+fc093*fc163+fc094*fc164+fc095*fc165+fc096*fc166+fc097*fc167+fc098*fc168+fc099*fc169）
						//{"fc092", "fc162", "fc154"},
						//{"fc093", "fc163", "fc155"},
						//{"fc094", "fc164", "fc156"},
						//{"fc095", "fc165", "fc157"},
						//{"fc096", "fc166", "fc158"},
						//{"fc097", "fc167", "fc159"},
						//{"fc098", "fc168", "fc160"},
						//{"fc099", "fc169", "fc161"},

						//CSB0118RC0242000
						//fc084、fc085、fc086、fc087、fc088、fc089、fc090、fc091、fc092、fc093、fc094、fc095、fc096、fc097、fc098、fc099录入为“?”“？”时，清空结果值，
						//报问题件：问题件代码N0001、问题件描述：票据清单模糊（存在多个问题件时，只报一个）
						//2024年3月19日11:10:27
						//yyy票据清单模糊（单张发票存在多个问题件时，只报一个）
						arrCode := []string{"fc084", "fc085", "fc086", "fc087", "fc088", "fc089", "fc090", "fc091", "fc092", "fc093", "fc094", "fc095", "fc096", "fc097", "fc098", "fc099"}
						for _, s := range arrCode {
							if fields[l].Code == s &&
								(fields[l].ResultValue == "?" || fields[l].ResultValue == "？") {
								fc002LocArr := utils2.GetFieldLoc(fieldLocationMap["fc002"], i, -1, -1, -1)
								fc002Val := utils2.GetFieldValueByLocArr(obj, fc002LocArr, true)
								fields[l].FinalValue = ""
								issue := model3.Issue{
									Type:    "",
									Code:    "N0001",
									Message: fc002Val + "票据清单模糊",
								}
								fields[l].Issues = []model3.Issue{issue}
							}
						}

						//CSB0118RC0243000
						//当存在fc205字段时、且当fc205录入为1、2时，将同属性的fc054的结果值赋值给fc208
						if fields[l].Code == "fc205" &&
							(fields[l].ResultValue == "1" || fields[l].ResultValue == "2") {
							fc054locs := getFieldLoc(fieldLocationMap["fc054"], i, -1, -1, -1)
							fc208locs := getFieldLoc(fieldLocationMap["fc208"], i, -1, -1, -1)
							if len(fc054locs) > 0 && len(fc208locs) > 0 {
								invoiceMap.Invoice[0][fc208locs[0][3]].FinalValue = invoiceMap.Invoice[0][fc054locs[0][3]].FinalValue
							}
						}

						//CSB0118RC0246000
						//fc049、fc072、fc206、fc207、fc219、fc220、fc221、fc228、fc229、fc230以上字段为“?”“？”时，结果值默认为“0”
						code := []string{"fc206", "fc207", "fc219", "fc220", "fc221", "fc228", "fc229", "fc230"}
						for _, s := range code {
							rc, err = regexp.Compile("(\\?|？)")
							if err != nil {
								global.GLog.Error("CSB0118RC0246000" + err.Error())
								continue
							}
							if fields[l].Code == s && rc.MatchString(fields[l].ResultValue) {
								fields[l].FinalValue = "0"
							}
						}

						//CSB0118RC0247000
						//fc202、fc203、fc204为bc001中的字段，以下左边字段的结果值默认为该分块对应同一发票属性同分块中bc001的右边字段的结果值
						iCodes = [][]string{
							{"fc202", "fc062"},
							{"fc203", "fc003"},
							{"fc204", "fc002"},
						}
						for _, codes := range iCodes {
							if fields[l].Code == codes[0] {
								rightLocs := getFieldLoc(fieldLocationMap[codes[1]], i, j, k, -1)
								if len(rightLocs) > 0 {
									fields[l].FinalValue = invoiceMap.Invoice[0][rightLocs[0][3]].FinalValue
								}
							}
						}

						//CSB0118RC0248000
						//同一发票属性中，对以下对应的字段进行赋值：
						//bc008第三方报销单1分块：
						//fc216报销单类型1结果数据赋值给fc225发票-报销单类型2
						//fc219第三方报销金额1结果数据赋值给fc228发票-第三方报销金额2
						//fc222服务机构名称1结果数据赋值给fc231发票-服务机构名称2
						//bc009第三方报销单2分块：
						//fc217报销单类型2结果数据赋值给fc226发票-报销单类型3
						//fc220第三方报销金额2结果数据赋值给fc229发票-第三方报销金额3
						//fc223服务机构名称2结果数据赋值给fc232发票-服务机构名称3
						//bc010第三方报销单3分块：
						//fc218报销单类型3结果数据赋值给fc227发票-报销单类型4
						//fc221第三方报销金额3结果数据赋值给fc230发票-第三方报销金额4
						iCodes = [][]string{
							{"fc216", "fc225"},
							{"fc219", "fc228"},
							{"fc222", "fc231"},
							{"fc217", "fc226"},
							{"fc220", "fc229"},
							{"fc223", "fc232"},
							{"fc218", "fc227"},
							{"fc221", "fc230"},
							{"fc224", "fc233"},
						}
						for _, codes := range iCodes {
							if codes[0] == fields[l].Code {
								rightLocs := getFieldLoc(fieldLocationMap[codes[1]], i, -1, -1, -1)
								if len(rightLocs) > 0 {
									invoiceMap.Invoice[0][rightLocs[0][3]].FinalValue = fields[l].FinalValue
								}
							}
						}

						//CSB0118RC0249000
						//fc146的结果值包含“?”“?”时，清空结果值
						if fields[l].Code == "fc146" {
							rc, err = regexp.Compile("(\\?|？)")
							if err != nil {
								global.GLog.Error("CSB0118RC0249000" + err.Error())
								continue
							}
							if rc.MatchString(fields[l].ResultValue) {
								fields[l].FinalValue = ""
							}
						}

						//CSB0118RC0250000
						//fc203结果值为1时，替换为O（大写字母O）
						//fc203结果值为2时，替换为H
						if fields[l].Code == "fc203" {
							if fields[l].FinalValue == "1" {
								fields[l].FinalValue = "O"
							}
							if fields[l].FinalValue == "2" {
								fields[l].FinalValue = "H"
							}
						}

						//CSB0118RC0251000 报销单最多只有一张所以invoiceMap.BaoXiaoDan[0]
						//fc144、fc145、fc146为bc003中的字段，以下左边字段的结果值默认为该分块对应同一发票属性bc001的右边字段的结果值
						iCodes = [][]string{
							{"fc144", "fc062"},
							{"fc145", "fc003"},
							{"fc146", "fc002"},
						}
						for _, codes := range iCodes {
							if fields[l].Code == codes[0] {
								rightLocs := getFieldLoc(fieldLocationMap[codes[1]], i, -1, -1, -1)
								if len(rightLocs) > 0 {
									fields[l].FinalValue = invoiceMap.Invoice[0][rightLocs[0][3]].FinalValue
								}
							}
						}

						//CSB0118RC0255000
						//当同一发票属性的fc154、fc155、fc156、fc157、fc158、fc159、fc160、fc161录入为1或4或2时，同一发票属性的fc171结果值合计金额不等于fc067录入值时，fc067出问题件：清单自费金额合计与结算单自费金额不一致
						if fields[l].Code == "fc067" {
							sFields := []string{"fc154", "fc155", "fc156", "fc157", "fc158", "fc159", "fc160", "fc161"}
							fc171Total := getFc171Total(fieldLocationMap, i, invoiceMap)
							var fieldsVal []string
							for _, field := range sFields {
								locs := getFieldLoc(fieldLocationMap[field], i, -1, -1, -1)
								for _, loc := range locs {
									fieldsVal = append(fieldsVal, invoiceMap.QingDan[loc[2]][loc[3]].ResultValue)
								}
							}
							//包含
							if utils2.HasItem(fieldsVal, "1") || utils2.HasItem(fieldsVal, "4") || utils2.HasItem(fieldsVal, "2") {
								v, err := decimal.NewFromString(fields[l].ResultValue)
								if err != nil {
									global.GLog.Error("CSB0118RC0257000" + err.Error())
									//continue
									v = decimal.Zero
								}
								if !v.Equal(fc171Total) {
									issue := model3.Issue{
										Type:    "",
										Code:    "",
										Message: "清单自费金额合计与结算单自费金额不一致",
									}
									fields[l].Issues = append(fields[l].Issues, issue)
								}
							}

						}

						//CSB0118RC0257000
						//fc067的结果值等于同属性的fc008结果值减去fc067的录入值（当fc067的录入值包含问号或为A时，将fc067的录入值时视为0进行计算）
						if fields[l].Code == "fc067" {
							val := decimal.Zero
							compile, err := regexp.Compile("\\?|？|A")
							if err != nil {
								global.GLog.Error("CSB0118RC0257000" + err.Error())
								continue
							}
							if compile.MatchString(fields[l].ResultValue) {
								val = decimal.Zero
							} else {
								val, err = decimal.NewFromString(fields[l].ResultValue)
								if err != nil {
									global.GLog.Error("CSB0118RC0257000" + err.Error())
									//continue
								}
							}
							locs := getFieldLoc(fieldLocationMap["fc008"], i, -1, -1, -1)
							if len(locs) > 0 {
								fc008vVal, err := decimal.NewFromString(invoiceMap.Invoice[0][locs[0][3]].FinalValue)
								if err != nil {
									global.GLog.Error("CSB0118RC0257000" + err.Error())
									//continue
								}
								fields[l].FinalValue = fc008vVal.Sub(val).String()
							}
						}
						//CSB0118RC0336000 同一发票下，当fc275结果值为1时，报问题件：问题件代码：N0005，问题件描述：票据号xxx执行单病种录入规则（xxx为同一发票下fc002的结果值）
						if fields[l].Code == "fc275" {
							var fc002Val string
							fc002locs := getFieldLoc(fieldLocationMap["fc002"], i, -1, -1, -1)
							if len(fc002locs) > 0 {
								fc002Val = invoiceMap.Invoice[0][fc002locs[0][3]].FinalValue
							}
							issue := model3.Issue{
								Type:    "",
								Code:    "N0005",
								Message: "票据号" + fc002Val + "执行单病种录入规则",
							}
							if fields[l].FinalValue == "1" {
								fields[l].Issues = append(fields[l].Issues, issue)
							}

						}

					}
				}
			}
		}

		//CSB0118RC0221000
		//在发票循环的最后一次去赋值
		if i == len(obj.Invoice)-1 && countErrInvoice > 0 {
			//设置第i张发票的003的问题件
			issue := model3.Issue{
				Type:    "",
				Code:    "N0001",
				Message: fmt.Sprintf("%d张票据模糊", countErrInvoice),
			}
			obj.Invoice[oneErrIndexArr[0]].Invoice[0][oneErrIndexArr[1]].Issues = []model3.Issue{issue}
		}
	}

	for aa, invoice := range obj.Invoice {
		qCodeArr := [][]string{
			{"fc084", "fc162"},
			{"fc085", "fc163"},
			{"fc086", "fc164"},
			{"fc087", "fc165"},
			{"fc088", "fc166"},
			{"fc089", "fc167"},
			{"fc090", "fc168"},
			{"fc091", "fc169"},
		}
		for bb, fields := range invoice.QingDan {
			for cc, field := range fields {
				for _, qCode := range qCodeArr {
					if field.Code == qCode[1] {
						ff1 := getOneResultValue(fields, qCode[0])
						if RegIsMatch(ff1, `^(工本费|病历费|卡费|复印费|陪护费|陪人费|膳食费|陪床费)$`) {
							obj.Invoice[aa].QingDan[bb][cc].FinalValue = "1"
						}
						break
					}
				}
			}
		}

	}

	//CSB0118RC0321000
	//结果数据
	//以下一行为一组，当案件列表的销售渠道为1时，第四列结果值不为0或1时时，将第四、五列结果值赋值为0
	myCode := [][]string{
		{"fc084", "fc154", "fc092", "fc162", "fc172"},
		{"fc085", "fc155", "fc093", "fc163", "fc173"},
		{"fc086", "fc156", "fc094", "fc164", "fc174"},
		{"fc087", "fc157", "fc095", "fc165", "fc175"},
		{"fc088", "fc158", "fc096", "fc166", "fc176"},
		{"fc089", "fc159", "fc097", "fc167", "fc177"},
		{"fc090", "fc160", "fc098", "fc168", "fc178"},
		{"fc091", "fc161", "fc099", "fc169", "fc179"},
	}
	for _, codes := range myCode {
		for _, loc4 := range fieldLocationMap[codes[3]] {
			loc5 := getFieldLoc(fieldLocationMap[codes[4]], loc4[0], loc4[1], loc4[2], loc4[4])
			val4 := utils2.GetFieldValueByLoc(obj, loc4, true)
			if len(loc5) == 1 && val4 != "1" && val4 != "0" && bill.SaleChannel == "1" {
				utils2.SetOnlyOneFinalValue(obj, loc4, "0")
				//utils2.SetOnlyOneFinalValue(obj, loc5[0], "0")
			}
		}
	}

	for i := 0; i < len(obj.Invoice); i++ {
		fc170sInInvoice := make([]map[string]interface{}, 0)
		fc171InInvoice := decimal.Zero
		//同一发票
		invoiceMap := obj.Invoice[i]

		flag := true
		for _, qinDanArr := range invoiceMap.QingDan {
			for _, field := range qinDanArr {
				if field.Code == "fc152" && (field.FinalValue != "199" && field.FinalValue != "299") {
					flag = false
				}
				if field.Code == "fc153" && field.FinalValue != "0" {
					flag = false
				}
			}
		}

		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					for l := 0; l < len(fields); l++ {
						//i:发票index j:发票结构体字段index k:字段二维数组x的index l:字段二维数组y的index

						//CSB0118RC0234000
						//（新华理赔有类似需求）下列字段为循环分块中的字段，根据fc059-发票属性和fc060-清单所属发票进行匹配：
						//1.将该清单属性的所有该分块中第一列录入值为1或3所对应的按以下格式汇总到fc170的结果值中
						//	（当清单循环分块中的fc150、fc152的结果值存在重复的情况下，将重复的fc150内容该清单属性的内容汇总到第一个fc150对应的清单分块中的fc170）：
						//	 AA 自付比例BB% 自付金额C.CC元；DD 自付比例EE% 自付金额F.FF元；...；
						//	（AA为第一列录入值为1所对应的第四列结果值；BB%为对应的第二列结果值，
						// 	如第二列结果值格式为0.1需要转为10%；C.CC为对应的第二列结果值乘第三列结果值）
						//2.将该清单属性的所有该分块中第一列录入值为2或6所对应的按以下格式汇总到fc170的结果值中
						//  （当清单循环分块中的fc150、fc152的结果值存在重复的情况下，将重复的fc150内容该清单属性的内容汇总到第一个fc150对应的清单分块中的fc170）：
						// 	 AA 自付比例100% 自付金额C.CC元；DD 自付比例100% 自付金额F.FF元；...；
						//	 （AA为第一列录入值为2所对应的第四列结果值；C.CC为对应的第三列结果值）
						//3.将该清单属性的所有该分块中第一列录入值为4所对应的按以下格式汇总到fc170的结果值中
						//	（当清单循环分块中的fc150、fc152的结果值存在重复的情况下，将重复的fc150内容该清单属性的内容汇总到第一个fc150对应的清单分块中的fc170）：
						//	AA 自付比例BB% 自付金额C.CC元（AA为第一列录入值为4所对应的第四列结果值，如第二列结果值格式为0.1需要转为10%；C.CC为对应的第五列结果值）
						//	第1和第2点、第3点、第4点需要共同汇总到（第一个）fc170（需按自付比例从大到小进行排列）的结果值中，如以4点规则匹配完之后，（）自付金额为0时，则清空这一组清单5个字段的结果数据
						//	（若与fc059同一分块的fc181-出险地点省为广西壮族：
						//		1、第四列字段匹配不上常量库且第一列字段录入值为2时，不执行第2点改为执行第1点；
						//		2、第四列字段匹配常量表中的“中文名称”所对应的“统计类别”为材料费且第二列字段结果值不为0时，执行第1点）
						iCodes := [][]string{
							{"fc154", "fc162", "fc092", "fc084", "fc172"},
							{"fc155", "fc163", "fc093", "fc085", "fc173"},
							{"fc156", "fc164", "fc094", "fc086", "fc174"},
							{"fc157", "fc165", "fc095", "fc087", "fc175"},
							{"fc158", "fc166", "fc096", "fc088", "fc176"},
							{"fc159", "fc167", "fc097", "fc089", "fc177"},
							{"fc160", "fc168", "fc098", "fc090", "fc178"},
							{"fc161", "fc169", "fc099", "fc091", "fc179"},
						}
						//for _, code := range iCodes {
						//	if fields[l].Code == code[0] {
						//		fc181Loc := getFieldLoc(fieldLocationMap["fc181"], i, -1, -1, -1)
						//		fc1Loc := getFieldLoc(fieldLocationMap[code[1]], i, -1, k, fields[l].BlockIndex)
						//		fc2Loc := getFieldLoc(fieldLocationMap[code[2]], i, -1, k, fields[l].BlockIndex)
						//		fc3Loc := getFieldLoc(fieldLocationMap[code[3]], i, -1, k, fields[l].BlockIndex)
						//		fc4Loc := getFieldLoc(fieldLocationMap[code[4]], i, -1, k, fields[l].BlockIndex)
						//		if len(fc181Loc) < 1 || len(fc1Loc) < 1 || len(fc2Loc) < 1 || len(fc3Loc) < 1 || len(fc4Loc) < 1 {
						//			continue
						//		}
						//		fc181Val := invoiceMap.Invoice[0][fc181Loc[0][3]].ResultValue
						//		fc1Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc1Loc[0][3]].FinalValue)
						//		fc2Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc2Loc[0][3]].FinalValue)
						//		fc3Val := invoiceMap.QingDan[k][fc3Loc[0][3]].FinalValue
						//		fc4Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc4Loc[0][3]].FinalValue)
						//		fc3Input := invoiceMap.QingDan[k][fc3Loc[0][3]].ResultValue
						//		//global.GLog.Warn(fields[l].Code + ":::" + strconv.Itoa(fields[l].BlockIndex) + ":::fc181Val:::" + fc181Val)
						//		//global.GLog.Warn(code[1] + ":::" + strconv.Itoa(fields[l].BlockIndex) + ":::fc1Val:::" + fc1Val.String())
						//		//global.GLog.Warn(code[2] + ":::" + strconv.Itoa(fields[l].BlockIndex) + ":::fc2Val:::" + fc2Val.String())
						//		//global.GLog.Warn(code[3] + ":::" + strconv.Itoa(fields[l].BlockIndex) + ":::fc3Val:::" + fc3Val)
						//		//global.GLog.Warn(code[4] + ":::" + strconv.Itoa(fields[l].BlockIndex) + ":::fc4Val:::" + fc4Val.String())
						//		//global.GLog.Warn(code[3] + ":::" + strconv.Itoa(fields[l].BlockIndex) + ":::fc3Input:::" + fc3Input)
						//
						//		fc1ValMulFc2Val := fc1Val.Mul(fc2Val)
						//		fc1ValPercent := fc1Val.Mul(decimal.NewFromInt(100))
						//
						//		messOrder := 0
						//		if fields[l].ResultValue == "1" || fields[l].ResultValue == "3" {
						//			messOrder = 1
						//		}
						//		if fields[l].ResultValue == "2" || fields[l].ResultValue == "6" {
						//			messOrder = 2
						//		}
						//		if fields[l].ResultValue == "4" {
						//			messOrder = 3
						//		}
						//		if fc181Val == "广西壮族" {
						//			_, ok := constMap["广西壮族"][fc3Input]
						//			if !ok && fields[l].ResultValue == "2" {
						//				messOrder = 1
						//			} else if ok && fields[l].ResultValue == "2" {
						//				messOrder = 2
						//			}
						//		}
						//		if v, ok := constMap["guangXiZhuangZuMap"][fc3Input]; fc181Val == "广西壮族" &&
						//			ok && v == "材料费" && invoiceMap.QingDan[k][fc1Loc[0][3]].FinalValue != "0" {
						//			messOrder = 1
						//		}
						//
						//		switch messOrder {
						//		case 1:
						//			fc170s = append(fc170s, map[string]interface{}{
						//				"mess": fc3Val + " 自付比例" + fc1ValPercent.String() + "% 自付金额" + fc1ValMulFc2Val.String() + "元；", "num": fc1ValPercent,
						//			})
						//			break
						//		case 2:
						//			fc170s = append(fc170s, map[string]interface{}{
						//				"mess": fc3Val + " 自付比例100% 自付金额" + fc2Val.String() + "元；", "num": decimal.NewFromInt(100),
						//			})
						//			break
						//		case 3:
						//			fc170s = append(fc170s, map[string]interface{}{
						//				"mess": fc3Val + " 自付比例" + fc1ValPercent.String() + "% 自付金额" + fc4Val.String() + "元；", "num": fc1ValPercent,
						//			})
						//			break
						//		default:
						//			//global.GLog.Error("CSB0118RC0234000不符合条件")
						//		}
						//
						//		//自付金额为0时，则清空这一组清单5个字段的结果数据
						//		//global.GLog.Info("fc1ValMulFc2Val:::" + fc1ValMulFc2Val.String())
						//		//global.GLog.Info("", zap.Any("", fc1ValMulFc2Val.Equal(decimal.Zero)))
						//		if fc1ValMulFc2Val.Equal(decimal.Zero) {
						//			fields[l].FinalValue = ""
						//			invoiceMap.QingDan[k][fc1Loc[0][3]].FinalValue = ""
						//			invoiceMap.QingDan[k][fc2Loc[0][3]].FinalValue = ""
						//			invoiceMap.QingDan[k][fc3Loc[0][3]].FinalValue = ""
						//			invoiceMap.QingDan[k][fc4Loc[0][3]].FinalValue = ""
						//			if messOrder != 0 {
						//				fc170s = fc170s[:len(fc170s)-1]
						//			}
						//		}
						//
						//		//CSB0118RC0241000
						//		fc171 = fc171.Add(fc1ValMulFc2Val)
						//	}
						//}

						if fields[l].Code == "fc171" && fields[l].BlockIndex == 0 {
							fc171SameMiniBlock := decimal.Zero
							for _, code := range iCodes {
								fc1Locs := getFieldLoc(fieldLocationMap[code[1]], i, -1, k, -1)
								fc2Locs := getFieldLoc(fieldLocationMap[code[2]], i, -1, k, -1)
								for blockIndex, _ := range fc1Locs {
									fc1Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc1Locs[blockIndex][3]].FinalValue)
									fc2Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc2Locs[blockIndex][3]].FinalValue)
									fc1ValMulFc2Val := fc1Val.Mul(fc2Val)
									fc171SameMiniBlock = fc171SameMiniBlock.Add(fc1ValMulFc2Val)
								}
							}
							fields[l].FinalValue = fc171SameMiniBlock.String()
							fc171InInvoice = fc171InInvoice.Add(fc171SameMiniBlock)
							if flag {
								fields[l].FinalValue = ""
							}
						}

						if fields[l].Code == "fc170" && fields[l].BlockIndex == 0 {
							fc170s := make([]map[string]interface{}, 0)
							for _, code := range iCodes {
								fc181Loc := getFieldLoc(fieldLocationMap["fc181"], i, -1, -1, -1)
								fc0Locs := getFieldLoc(fieldLocationMap[code[0]], i, -1, k, -1)
								fc1Locs := getFieldLoc(fieldLocationMap[code[1]], i, -1, k, -1)
								fc2Locs := getFieldLoc(fieldLocationMap[code[2]], i, -1, k, -1)
								fc3Locs := getFieldLoc(fieldLocationMap[code[3]], i, -1, k, -1)
								fc4Locs := getFieldLoc(fieldLocationMap[code[4]], i, -1, k, -1)
								if len(fc181Loc) < 1 || len(fc1Locs) < 1 || len(fc2Locs) < 1 || len(fc3Locs) < 1 || len(fc4Locs) < 1 {
									continue
								}
								fc181Val := invoiceMap.Invoice[0][fc181Loc[0][3]].ResultValue
								for blockIndex, _ := range fc1Locs {
									fc0Val := invoiceMap.QingDan[k][fc0Locs[blockIndex][3]].ResultValue
									fc1Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc1Locs[blockIndex][3]].FinalValue)
									fc2Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc2Locs[blockIndex][3]].FinalValue)
									fc3Val := invoiceMap.QingDan[k][fc3Locs[blockIndex][3]].FinalValue
									fc4Val, _ := decimal.NewFromString(invoiceMap.QingDan[k][fc4Locs[blockIndex][3]].FinalValue)
									fc3Input := invoiceMap.QingDan[k][fc3Locs[blockIndex][3]].ResultValue

									fc1ValMulFc2Val := fc1Val.Mul(fc2Val)
									fc1ValPercent := fc1Val.Mul(decimal.NewFromInt(100))

									messOrder := 0
									if fc0Val == "1" || fc0Val == "3" || fc0Val == "5" {
										messOrder = 1
									}
									if fc0Val == "2" || fc0Val == "6" {
										messOrder = 2
									}
									if fc0Val == "4" {
										messOrder = 3
									}
									if fc181Val == "广西壮族" {
										_, ok := constMap["广西壮族"][fc3Input]
										if !ok && fc0Val == "2" {
											messOrder = 1
										} else if ok && fc0Val == "2" {
											messOrder = 2
										}
									}

									if v, ok := constMap["guangXiZhuangZuMap"][fc3Input]; fc181Val == "广西壮族" &&
										ok && v == "材料费" && invoiceMap.QingDan[k][fc1Locs[blockIndex][3]].FinalValue != "0" {
										messOrder = 1
									}
									switch messOrder {
									case 1:
										if fc1ValMulFc2Val.String() == "" || fc1ValMulFc2Val.String() == "0" {
											break
										}
										fc170s = append(fc170s, map[string]interface{}{
											"mess": fc3Val + " 自付比例" + fc1ValPercent.String() + "% 自付金额" + fc1ValMulFc2Val.String() + "元；", "num": fc1ValPercent,
										})
										break
									case 2:
										if fc2Val.String() == "" || fc2Val.String() == "0" {
											break
										}
										fc170s = append(fc170s, map[string]interface{}{
											"mess": fc3Val + " 自付比例100% 自付金额" + fc2Val.String() + "元；", "num": decimal.NewFromInt(100),
										})
										break
									case 3:
										if fc4Val.String() == "" || fc4Val.String() == "0" {
											break
										}
										fc170s = append(fc170s, map[string]interface{}{
											"mess": fc3Val + " 自付比例" + fc1ValPercent.String() + "% 自付金额" + fc4Val.String() + "元；", "num": fc1ValPercent,
										})
										break
									default:
										//global.GLog.Error("CSB0118RC0234000不符合条件")
									}
								}
							}

							sort.Slice(fc170s, func(i, j int) bool {
								return fc170s[i]["num"].(decimal.Decimal).Cmp(fc170s[j]["num"].(decimal.Decimal)) > 0
							})
							fc170 := ""
							for _, m := range fc170s {
								mess, ok := m["mess"]
								if !ok {
									continue
								}
								fc170 += mess.(string)
							}
							fields[l].FinalValue = fc170
							if flag {
								fields[l].FinalValue = ""
							}
							fc170sInInvoice = append(fc170sInInvoice, fc170s...)
						}

						//CSB0118RC0231000
						//同一发票属性的fc003录入为1时，左边第二列的结果值不为1.00时，清空对应一整列的数据
						//（如fc162的结果值为0.1时，清空fc154、fc162、fc092、fc084、fc172的结果值）
						//2022年04月20日11:42:03 更改不为1.00且不为1
						// iCode := [][]string{
						// 	{"fc154", "fc162", "fc092", "fc084", "fc172"},
						// 	{"fc155", "fc163", "fc093", "fc085", "fc173"},
						// 	{"fc156", "fc164", "fc094", "fc086", "fc174"},
						// 	{"fc157", "fc165", "fc095", "fc087", "fc175"},
						// 	{"fc158", "fc166", "fc096", "fc088", "fc176"},
						// 	{"fc159", "fc167", "fc097", "fc089", "fc177"},
						// 	{"fc160", "fc168", "fc098", "fc090", "fc178"},
						// 	{"fc161", "fc169", "fc099", "fc091", "fc179"},
						// }
						// if fields[l].Code == "fc003" && fields[l].ResultValue == "1" {
						// 	for _, arr := range iCode {
						// 		for i2, _ := range invoiceMap.QingDan {
						// 			fc162Locs := getFieldLoc(fieldLocationMap[arr[1]], i, -1, i2, -1)
						// 			for _, fc162Loc := range fc162Locs {
						// 				if invoiceMap.QingDan[i2][fc162Loc[3]].FinalValue != "1.00" && invoiceMap.QingDan[i2][fc162Loc[3]].FinalValue != "1" {
						// 					for _, s := range arr {
						// 						fc154Loc := getFieldLoc(fieldLocationMap[s], i, -1, i2, fc162Loc[4])
						// 						if len(fc154Loc) > 0 {
						// 							invoiceMap.QingDan[i2][fc154Loc[0][3]].FinalValue = ""
						// 						}
						// 					}
						// 				}
						// 			}
						// 		}
						// 	}
						// }
					}
				}
			}
		}

		if flag {
			sort.Slice(fc170sInInvoice, func(iii, jjj int) bool {
				return fc170sInInvoice[iii]["num"].(decimal.Decimal).Cmp(fc170sInInvoice[jjj]["num"].(decimal.Decimal)) > 0
			})
			fmt.Println("fc170sInInvoice", fc170sInInvoice)
			fmt.Println("fc171InInvoice", fc171InInvoice)
			fc170LocArr := getFieldLoc(fieldLocationMap["fc170"], i, -1, 0, -1)
			fc171LocArr := getFieldLoc(fieldLocationMap["fc171"], i, -1, 0, -1)
			if len(fc170LocArr) > 0 && len(fc171LocArr) > 0 {
				fc170Mess := ""
				for _, m := range fc170sInInvoice {
					mess, ok := m["mess"]
					if !ok {
						continue
					}
					fc170Mess += mess.(string)
				}
				invoiceMap.QingDan[0][fc171LocArr[0][3]].FinalValue = fc171InInvoice.String()
				invoiceMap.QingDan[0][fc170LocArr[0][3]].FinalValue = fc170Mess
			}
		}

		//fc170Loc := getFieldLoc(fieldLocationMap["fc170"], i, -1, 0, -1)
		//fc171Loc := getFieldLoc(fieldLocationMap["fc171"], i, -1, 0, -1)
		//循环到最后一个小分块再赋值
		//fmt.Println(fc170s)
		//fmt.Println(fc171)
		//sort.Slice(fc170s, func(i, j int) bool {
		//	return fc170s[i]["num"].(decimal.Decimal).Cmp(fc170s[j]["num"].(decimal.Decimal)) > 0
		//})
		//fc170 := ""
		//for _, m := range fc170s {
		//	mess, ok := m["mess"]
		//	if !ok {
		//		continue
		//	}
		//	fc170 += mess.(string)
		//}
		//if len(fc170Loc) > 0 {
		//	invoiceMap.QingDan[0][fc170Loc[0][3]].FinalValue = fc170
		//}
		//if len(fc170Loc) > 0 {
		//	fmt.Println(fc171)
		//	invoiceMap.QingDan[0][fc171Loc[0][3]].FinalValue = fc171.String()
		//}

		//CSB0118RC0252000
		//当同一发票属性不存在fc067字段时，且fc205录入为1或2时，将所有fc171的合计结果值赋值给fc209，（放在fc171字段的所有校验之后），fc211的值等于同一发票属性的fc008结果值-fc209结果值
		//当同一发票属性不存在fc067字段时，且fc205录入不为1或2或A时（为空不进行校验），将所有fc171的合计结果值赋值给fc210，（放在fc171字段的所有校验之后），fc212的值等于同一发票属性的fc008结果值-fc210结果值
		fc067Locs := getFieldLoc(fieldLocationMap["fc067"], i, -1, -1, -1)
		fc205Locs := getFieldLoc(fieldLocationMap["fc205"], i, -1, -1, -1)
		if len(fc205Locs) > 0 && len(fc067Locs) < 1 {
			fc171Total := getFc171Total(fieldLocationMap, i, invoiceMap)
			fc008Locs := getFieldLoc(fieldLocationMap["fc008"], i, -1, -1, -1)
			fc205Val := invoiceMap.Invoice[fc205Locs[0][2]][fc205Locs[0][3]].ResultValue
			if fc205Val == "1" || fc205Val == "2" {
				fc209Locs := getFieldLoc(fieldLocationMap["fc209"], i, -1, -1, -1)
				fc211Locs := getFieldLoc(fieldLocationMap["fc211"], i, -1, -1, -1)
				if len(fc209Locs) > 0 && len(fc211Locs) > 0 && len(fc008Locs) > 0 {
					fc008vVal, err := decimal.NewFromString(invoiceMap.Invoice[0][fc008Locs[0][3]].FinalValue)
					if err != nil {
						global.GLog.Error("CSB0118RC0252000" + err.Error())
						//continue
					}
					invoiceMap.Invoice[0][fc209Locs[0][3]].FinalValue = fc171Total.String()
					invoiceMap.Invoice[0][fc211Locs[0][3]].FinalValue = fc008vVal.Sub(fc171Total).String()
				}
			} else if !(fc205Val == "1" || fc205Val == "2" || fc205Val == "A") {
				fc210Locs := getFieldLoc(fieldLocationMap["fc210"], i, -1, -1, -1)
				fc212Locs := getFieldLoc(fieldLocationMap["fc212"], i, -1, -1, -1)
				if len(fc210Locs) > 0 && len(fc212Locs) > 0 && len(fc008Locs) > 0 {
					fc008vVal, err := decimal.NewFromString(invoiceMap.Invoice[0][fc008Locs[0][3]].FinalValue)
					if err != nil {
						global.GLog.Error("CSB0118RC0252000" + err.Error())
						//continue
					}
					invoiceMap.Invoice[0][fc210Locs[0][3]].FinalValue = fc171Total.String()
					invoiceMap.Invoice[0][fc212Locs[0][3]].FinalValue = fc008vVal.Sub(fc171Total).String()
				}
			}
		}

		//CSB0118RC0244000
		//fc067的结果值等于同属性的fc008结果值减去fc067的录入值，当fc067的录入值包含问号或为A为0时，fc067的结果值等于同属性的fc008的结果值减去同属性所有fc171结果值的合计，最后将计算后的计算赋值给fc211，不存在fc067时，不进行该校验
		//fc102的结果值等于同属性的fc008结果值减去fc102的录入值，当fc102的录入值包含问号或为A为0时，将fc102的结果值等于同属性的fc008的结果值减去同属性所有fc171结果值的合计），将计算后的计算赋值给fc212，不存在fc102时，不进行该校验
		iCodes := [][]string{
			{"fc067", "fc211"},
			{"fc102", "fc212"},
		}
		for _, codes := range iCodes {
			fc067Loc := getFieldLoc(fieldLocationMap[codes[0]], i, -1, 0, -1)
			if len(fc067Loc) > 0 {
				compile, err := regexp.Compile("(\\?|？|A)")
				if err != nil {
					global.GLog.Error("CSB0118RC0244000" + err.Error())
					continue
				}
				fc008locs := getFieldLoc(fieldLocationMap["fc008"], i, -1, -1, -1)
				if len(fc008locs) > 0 {
					fc008Str := invoiceMap.Invoice[0][fc008locs[0][3]].FinalValue
					fc067ResultValue := invoiceMap.BaoXiaoDan[fc067Loc[0][2]][fc067Loc[0][3]].ResultValue
					fc008Val, err := decimal.NewFromString(fc008Str)
					if err != nil {
						global.GLog.Error("CSB0118RC0244000" + err.Error())
						continue
					}
					if compile.MatchString(fc067ResultValue) || fc067ResultValue == "0" {
						fc171Total := getFc171Total(fieldLocationMap, i, invoiceMap)
						fmt.Println(fc171Total)
						invoiceMap.BaoXiaoDan[fc067Loc[0][2]][fc067Loc[0][3]].FinalValue = fc008Val.Sub(fc171Total).String()
					} else {
						val, err := decimal.NewFromString(fc067ResultValue)
						if err != nil {
							global.GLog.Error("CSB0118RC0244000" + err.Error())
							continue
						}
						invoiceMap.BaoXiaoDan[fc067Loc[0][2]][fc067Loc[0][3]].FinalValue = fc008Val.Sub(val).String()
					}
					rightLocs := getFieldLoc(fieldLocationMap[codes[1]], i, -1, -1, -1)
					for _, loc := range rightLocs {
						invoiceMap.Invoice[0][loc[3]].FinalValue = invoiceMap.BaoXiaoDan[fc067Loc[0][2]][fc067Loc[0][3]].FinalValue
					}
				}
			}
		}

		//CSB0118RC0253000
		//调整需求：1.当不存在bc007分块时，同一发票属性的fc154、fc155、fc156、fc157、fc158、fc159、fc160、fc161录入为1或4或2时，
		//	同一发票属性的fc171结果值合计金额不等于fc067录入值时，
		//	将fc067的录入值赋值给给同一属性第一个fc171的结果值，
		//	并将其余fc171的结果值默认为0.00（fc171的结果值为空时，不需要默认为0.00），
		//	将最后赋值后的fc171的结果值赋值给fc209，
		//	并fc211的值等于同一发票属性的fc008结果值-fc209结果值，
		//	当同一清单分块里面的fc084、fc085、fc086、fc087、fc088、fc089、fc090、fc091、fc154、fc155、fc156、fc157、fc158、fc159、fc161、fc092、fc093、fc094、fc095、fc096、fc097、fc098、fc099、fc162、fc163、fc164、fc165、fc166、fc167、fc168、fc169、fc172、fc173、fc174、fc175、fc176、fc177、fc178、fc179均为空是时，fc171不需要赋值为0.00
		//2.当不存在bc007分块，同一发票属性的fc154、fc155、fc156、fc157、fc158、fc159、fc160、fc161录入为3，
		//	且同一发票属性的fc171结果值合计金额不等于fc067录入值时，将fc067的录入值赋值给给同一属性第一个fc171的结果值。
		fc067Loc := getFieldLoc(fieldLocationMap["fc067"], i, -1, 0, -1)
		if invoiceMap.InvoiceDaXiang == nil && len(fc067Loc) > 0 {
			sFields := []string{"fc154", "fc155", "fc156", "fc157", "fc158", "fc159", "fc160", "fc161"}
			flag142 := false
			flag3 := false
			for _, item := range sFields {
				var itemValueArr []string
				locs := getFieldLoc(fieldLocationMap[item], i, -1, -1, -1)
				for _, loc := range locs {
					itemValueArr = append(itemValueArr, invoiceMap.QingDan[loc[2]][loc[3]].ResultValue)
				}
				if utils2.HasItem(itemValueArr, "1") || utils2.HasItem(itemValueArr, "4") || utils2.HasItem(itemValueArr, "2") {
					flag142 = true
					break
				}
				if utils2.HasItem(itemValueArr, "3") {
					flag3 = true
					break
				}
			}
			fc171Total := getFc171Total(fieldLocationMap, i, invoiceMap)
			fc067Val, err := decimal.NewFromString(invoiceMap.BaoXiaoDan[fc067Loc[0][2]][fc067Loc[0][3]].ResultValue)
			if err != nil {
				global.GLog.Error("CSB0118RC0253000" + err.Error())
				//continue
			}
			fc171Locs := getFieldLoc(fieldLocationMap["fc171"], i, -1, -1, -1)
			if flag142 && !fc171Total.Equal(fc067Val) {
				flagNotNull := false
				fieldArr := []string{"fc084", "fc085", "fc086", "fc087", "fc088", "fc089", "fc090", "fc091", "fc154", "fc155", "fc156", "fc157", "fc158", "fc159", "fc161", "fc092", "fc093", "fc094", "fc095", "fc096", "fc097", "fc098", "fc099", "fc162", "fc163", "fc164", "fc165", "fc166", "fc167", "fc168", "fc169", "fc172", "fc173", "fc174", "fc175", "fc176", "fc177", "fc178", "fc179"}
				for _, s := range fieldArr {
					locs := getFieldLoc(fieldLocationMap[s], i, -1, -1, -1)
					for _, loc := range locs {
						if invoiceMap.QingDan[loc[2]][loc[3]].ResultValue != "" {
							flagNotNull = true
						}
					}
				}
				for index, loc := range fc171Locs {
					if index == 0 {
						invoiceMap.QingDan[loc[2]][loc[3]].FinalValue = invoiceMap.BaoXiaoDan[fc067Loc[0][2]][fc067Loc[0][3]].ResultValue
					} else if invoiceMap.QingDan[loc[2]][loc[3]].FinalValue == "" && flagNotNull {
						invoiceMap.QingDan[loc[2]][loc[3]].FinalValue = "0.00"
					}
				}
				global.GLog.Info("fc171Total::" + fc171Total.String())
				fc209Locs := getFieldLoc(fieldLocationMap["fc209"], i, -1, -1, -1)
				if len(fc209Locs) > 0 && len(fc171Locs) > 0 {
					invoiceMap.Invoice[0][fc209Locs[0][3]].FinalValue = invoiceMap.QingDan[fc171Locs[0][2]][fc171Locs[0][3]].FinalValue
				}
			}
			if flag3 && !fc171Total.Equal(fc067Val) {
				invoiceMap.QingDan[fc171Locs[0][2]][fc171Locs[0][3]].FinalValue = fc067Val.String()
			}
		}
	}

	//iCodes := [][]string{
	//	{"fc154", "fc162", "fc092", "fc084", "fc172"},
	//	{"fc155", "fc163", "fc093", "fc085", "fc173"},
	//	{"fc156", "fc164", "fc094", "fc086", "fc174"},
	//	{"fc157", "fc165", "fc095", "fc087", "fc175"},
	//	{"fc158", "fc166", "fc096", "fc088", "fc176"},
	//	{"fc159", "fc167", "fc097", "fc089", "fc177"},
	//	{"fc160", "fc168", "fc098", "fc090", "fc178"},
	//	{"fc161", "fc169", "fc099", "fc091", "fc179"},
	//}
	//for _, iCode := range iCodes {
	//	for _, twoLoc := range fieldLocationMap[iCode[1]] {
	//		twoVal := utils2.GetFieldValueByLoc(obj, twoLoc, true)
	//		oneLoc := utils2.GetFieldLoc(fieldLocationMap[iCode[0]], twoLoc[0], twoLoc[1], twoLoc[2], twoLoc[4])
	//		threeLoc := utils2.GetFieldLoc(fieldLocationMap[iCode[2]], twoLoc[0], twoLoc[1], twoLoc[2], twoLoc[4])
	//		fourLoc := utils2.GetFieldLoc(fieldLocationMap[iCode[3]], twoLoc[0], twoLoc[1], twoLoc[2], twoLoc[4])
	//		fiveLoc := utils2.GetFieldLoc(fieldLocationMap[iCode[4]], twoLoc[0], twoLoc[1], twoLoc[2], twoLoc[4])
	//		if len(oneLoc) != 1 || len(threeLoc) != 1 || len(fourLoc) != 1 || len(fiveLoc) != 1 {
	//			continue
	//		}
	//		if twoVal == "0" || twoVal == "" {
	//			utils2.SetFinalValue(obj, oneLoc, "")
	//			utils2.SetFinalValue(obj, fourLoc, "")
	//			utils2.SetFinalValue(obj, threeLoc, "")
	//			utils2.SetFinalValue(obj, fiveLoc, "")
	//			utils2.SetOnlyOneFinalValue(obj, twoLoc, "")
	//		}
	//	}
	//}

	for ii, invoice := range obj.Invoice {
		valueCheck := map[string]int{}
		for qq, fields := range invoice.QingDan {
			fc152 := getOneValue(fields, "fc152")
			// fc153 := getOneValue(fields, "fc153")
			// if fc152 == "" || fc153 == "" {
			// 	continue
			// }
			// val := fc152 + "_" + fc153
			if fc152 == "" {
				continue
			}
			val := fc152
			fmt.Println("-----------valvalvalval---------------", val)
			num, isexit := valueCheck[val]
			if isexit {
				fc170 := getOneValue(fields, "fc170")
				fc171 := getOneValue(fields, "fc171")
				fmt.Println("-----------fc170---------------", fc170, fc171)
				setFc17Value(obj.Invoice[ii].QingDan[num], fc170, fc171)
				clearfcValue(obj.Invoice[ii].QingDan[qq])
				// fc170 := getValInFieldsByCode(fields, "fc170")
				// fc171 := getAllVal(fields, "fc171")
			} else {
				valueCheck[val] = qq
			}
		}
	}

	for ii, invoice := range obj.Invoice {
		for qq, fields := range invoice.QingDan {
			for ff := 0; ff < len(fields); ff++ {
				if fields[ff].Code == "fc170" && fields[ff].FinalValue != "" {
					fc170 := fields[ff].FinalValue
					fc170s := strings.Split(fc170, "；")
					items := make([]map[string]interface{}, 0)
					finalValue := ""
					for _, value := range fc170s {
						if value == "" {
							continue
						}
						aa := RegReplace(value, `.+自付比例`, "")
						aa = RegReplace(value, `% 自付金额.+`, "")
						num, _ := decimal.NewFromString(aa)
						items = append(items, map[string]interface{}{
							"mess": value + "；", "num": num,
						})
					}
					// 自付比例" + fc1ValPercent.String() + "% 自付金额
					sort.Slice(items, func(i, j int) bool {
						return items[i]["num"].(decimal.Decimal).Cmp(items[j]["num"].(decimal.Decimal)) > 0
					})
					for _, m := range items {
						mess, ok := m["mess"]
						if !ok {
							continue
						}
						finalValue += mess.(string)
					}
					obj.Invoice[ii].QingDan[qq][ff].FinalValue = finalValue
				}
			}
		}
	}

	//CSB0118RC0245000
	//当一张单子中，同一发票属性的左、右列字段同时存在时，需将右边字段的结果值赋值给左列字段：同一行的两个字段视为一组数据，当一组字
	//bc001  bc003
	iCodes := [][]string{
		{"fc205", "fc101"},
		{"fc206", "fc052"},
		{"fc207", "fc051"},
		{"fc208", "fc079"},
	}
	for _, code := range iCodes {
		for _, codeRightLoc := range fieldLocationMap[code[1]] {
			fc101Val := utils2.GetFieldValueByLoc(obj, codeRightLoc, true)
			fc205Loc := getFieldLoc(fieldLocationMap[code[0]], codeRightLoc[0], -1, -1, -1)
			utils2.SetFinalValue(obj, fc205Loc, fc101Val)
		}
	}

	//CSB0118RC0306000
	//fc003录入为1时，fc006、fc007默认为fc005的结果值
	for _, codeLoc := range fieldLocationMap["fc003"] {
		fc003Val := utils2.GetFieldValueByLoc(obj, codeLoc, false)
		fc006Loc := getFieldLoc(fieldLocationMap["fc006"], codeLoc[0], -1, -1, -1)
		fc007Loc := getFieldLoc(fieldLocationMap["fc007"], codeLoc[0], -1, -1, -1)
		fc005Loc := getFieldLoc(fieldLocationMap["fc005"], codeLoc[0], -1, -1, -1)
		fc005ValArr := utils2.GetFieldValueArrByLocArr(obj, fc005Loc, true)
		if fc003Val == "1" && len(fc005ValArr) == 1 {
			utils2.SetFinalValue(obj, fc006Loc, fc005ValArr[0])
			utils2.SetFinalValue(obj, fc007Loc, fc005ValArr[0])
		}
	}

	//CSB0118RC0307000导出校验以下一行为一组，
	//当第一列录入内容包含【接种服务】【预防接种】、【疫苗】时，且第二列的结果值或第三列结果值大于0时，
	//出导出校验：xx存在“疫苗接种服务费”，请核实是否为狂犬病相关费用。（xx为对应的对应第一列的字段名称）
	myCodeArr := [][]string{
		{"fc084", "fc162", "fc172"},
		{"fc085", "fc163", "fc173"},
		{"fc086", "fc164", "fc174"},
		{"fc087", "fc165", "fc175"},
		{"fc088", "fc166", "fc176"},
		{"fc089", "fc167", "fc177"},
		{"fc090", "fc168", "fc178"},
		{"fc091", "fc169", "fc179"},
	}
	for _, myCode := range myCodeArr {
		for _, loc1 := range fieldLocationMap[myCode[0]] {
			loc2 := utils2.GetFieldLoc(fieldLocationMap[myCode[1]], loc1[0], loc1[1], loc1[2], loc1[4])
			loc3 := utils2.GetFieldLoc(fieldLocationMap[myCode[2]], loc1[0], loc1[1], loc1[2], loc1[4])
			fc002Loc := utils2.GetFieldLoc(fieldLocationMap["fc002"], loc1[0], -1, -1, -1)
			if len(loc2) < 1 || len(loc3) < 1 {
				continue
			}
			fc002Val := utils2.GetFieldValueByLoc(obj, fc002Loc[0], true)
			err, field1 := utils2.GetFieldByLoc(obj, loc1)
			if err != nil {
				global.GLog.Error("CSB0118RC0307000", zap.Error(errors.New("没有找到字段"+myCode[0])))
				continue
			}
			decimalVal2 := utils2.GetFieldDecimalValueByLoc(obj, loc2[0], true)
			decimalVal3 := utils2.GetFieldDecimalValueByLoc(obj, loc3[0], true)
			r := regexp.MustCompile("(接种服务|预防接种|疫苗)")
			if r.MatchString(field1.ResultValue) &&
				(decimalVal2.GreaterThan(decimal.Zero) || decimalVal3.GreaterThan(decimal.Zero)) {
				msg := fc002Val + "账单号的" + field1.Name + "存在“疫苗接种服务费”，请核实是否为狂犬病相关费用。;"
				if strings.Index(obj.Bill.WrongNote, msg) == -1 {
					obj.Bill.WrongNote += msg
				}
			}
		}
	}

	myCodeArr = [][]string{
		{"fc009", "fc010", "fc104", "fc105", "fc235", "fc255"},
		{"fc011", "fc012", "fc106", "fc107", "fc236", "fc256"},
		{"fc013", "fc014", "fc108", "fc109", "fc237", "fc257"},
		{"fc015", "fc016", "fc110", "fc111", "fc238", "fc258"},
		{"fc017", "fc018", "fc112", "fc113", "fc239", "fc259"},
		{"fc019", "fc020", "fc114", "fc115", "fc240", "fc260"},
		{"fc021", "fc022", "fc116", "fc117", "fc241", "fc261"},
		{"fc023", "fc024", "fc118", "fc119", "fc242", "fc262"},
		{"fc025", "fc026", "fc120", "fc121", "fc243", "fc263"},
		{"fc027", "fc028", "fc122", "fc123", "fc244", "fc264"},
		{"fc029", "fc030", "fc124", "fc125", "fc245", "fc265"},
		{"fc031", "fc032", "fc126", "fc127", "fc246", "fc266"},
		{"fc033", "fc034", "fc128", "fc129", "fc247", "fc267"},
		{"fc035", "fc036", "fc130", "fc131", "fc248", "fc268"},
		{"fc037", "fc038", "fc132", "fc133", "fc249", "fc269"},
		{"fc039", "fc040", "fc134", "fc135", "fc250", "fc270"},
		{"fc041", "fc042", "fc136", "fc137", "fc251", "fc271"},
		{"fc043", "fc044", "fc138", "fc139", "fc252", "fc272"},
		{"fc045", "fc046", "fc140", "fc141", "fc253", "fc273"},
		{"fc047", "fc048", "fc142", "fc143", "fc254", "fc274"},
	}
	nCodes := [][]string{
		{"fc255", "fc235", "fc104", "fc105"},
		{"fc256", "fc236", "fc106", "fc107"},
		{"fc257", "fc237", "fc108", "fc109"},
		{"fc258", "fc238", "fc110", "fc111"},
		{"fc259", "fc239", "fc112", "fc113"},
		{"fc260", "fc240", "fc114", "fc115"},
		{"fc261", "fc241", "fc116", "fc117"},
		{"fc262", "fc242", "fc118", "fc119"},
		{"fc263", "fc243", "fc120", "fc121"},
		{"fc264", "fc244", "fc122", "fc123"},
		{"fc265", "fc245", "fc124", "fc125"},
		{"fc266", "fc246", "fc126", "fc127"},
		{"fc267", "fc247", "fc128", "fc129"},
		{"fc268", "fc248", "fc130", "fc131"},
		{"fc269", "fc249", "fc132", "fc133"},
		{"fc270", "fc250", "fc134", "fc135"},
		{"fc271", "fc251", "fc136", "fc137"},
		{"fc272", "fc252", "fc138", "fc139"},
		{"fc273", "fc253", "fc140", "fc141"},
		{"fc274", "fc254", "fc142", "fc143"},
	}
	for aa, invoice := range obj.Invoice {
		// qCodeArr := [][]string{
		// 	{"fc084", "fc162"},
		// 	{"fc085", "fc163"},
		// 	{"fc086", "fc164"},
		// 	{"fc087", "fc165"},
		// 	{"fc088", "fc166"},
		// 	{"fc089", "fc167"},
		// 	{"fc090", "fc168"},
		// 	{"fc091", "fc169"},
		// }
		// for bb, fields := range invoice.QingDan {
		// 	for cc, field := range fields {
		// 		for _, qCode := range qCodeArr {
		// 			if field.Code == qCode[1] {
		// 				ff1 := getOneResultValue(fields, qCode[0])
		// 				if RegIsMatch(ff1, `^(工本费|病历费|卡费|复印费|陪护费|陪人费|膳食费|陪床费)$`) {
		// 					obj.Invoice[aa].QingDan[bb][cc].FinalValue = "1"
		// 				}
		// 				break
		// 			}
		// 		}
		// 	}
		// }
		// for _, myCode := range myCodeArr {
		cacheValue := make(map[string][]int)
		for bb, projectFields := range invoice.InvoiceDaXiang {
			for cc, field := range projectFields {
				if strings.Index("fc235|fc236|fc237|fc238|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248|fc249|fc250|fc251|fc252|fc253|fc254", field.Code) != -1 && field.FinalValue == "" {
					obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = "0"
				}
			}

			for _, myCode := range myCodeArr {
				field1 := getOneValue(projectFields, myCode[0])
				cacheIdx, isExit := cacheValue[field1]
				if isExit {
					for cc, field := range projectFields {
						if field.Code == myCode[0] {
							obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = ""
						}
						if field.Code == myCode[1] {
							sum := ParseFloat(obj.Invoice[cacheIdx[0]].InvoiceDaXiang[cacheIdx[1]][cacheIdx[2]].FinalValue)
							sum = SumFloat(sum, ParseFloat(field.FinalValue), "+")
							obj.Invoice[cacheIdx[0]].InvoiceDaXiang[cacheIdx[1]][cacheIdx[2]].FinalValue = ToString(sum)
							obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = ""
						}
					}
					field1 = ""
				} else {
					for cc, field := range projectFields {
						if field.Code == myCode[1] {
							idx := []int{aa, bb, cc}
							cacheValue[field1] = idx
							break
						}
					}

					// cacheValue = append(cacheValue, field1)
				}

				// field1R := getOneResultValue(projectFields, myCode[0])
				// if RegIsMatch(field1R, `^(工本费|病历费|卡费|复印费|陪护费|陪人费)$`) {
				// 	field2 := getOneResultValue(projectFields, myCode[1])
				// 	for cc, field := range projectFields {
				// 		if field.Code == myCode[2] {
				// 			obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = field1R
				// 		}
				// 		if field.Code == myCode[3] {
				// 			obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = field2
				// 		}
				// 	}

				// }
				if len(invoice.QingDan) > 0 {
					for _, QingDan := range invoice.QingDan {
						fc152 := getOneValue(QingDan, "fc152")
						fc170 := getOneValue(QingDan, "fc170")
						fc171 := getOneValue(QingDan, "fc171")
						// fmt.Println("----------------invoiceinvoiceinvoice-------------------", fc152, field1, fc170, fc171)
						if fc152 == field1 && fc171 != "" {
							for cc, field := range projectFields {
								if field.Code == myCode[5] {
									obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = fc170
								}
								if field.Code == myCode[4] {
									obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = fc171
								}
							}
						}
					}

				}

			}

			for _, nCode := range nCodes {
				field1R := getOneValue(projectFields, nCode[0])
				// fmt.Println("-------------------------nCodenCodenCode---------------------------------", nCode[0], field1R)
				if RegIsMatch(field1R, `(工本费|病历费|卡费|复印费|陪护费|陪人费|膳食费|陪床费)`) {
					field2 := getOneValue(projectFields, nCode[1])
					items := strings.Split(field1R, "；")
					val1 := ""
					val3 := ""
					val4 := 0.00
					for _, item := range items {
						if item == "" {
							continue
						}
						if RegIsMatch(item, `(工本费|病历费|卡费|复印费|陪护费|陪人费|膳食费|陪床费)`) {
							val3 += item + "；"
							aa := RegReplace(item, `.+自付金额`, "")
							aa = RegReplace(aa, `元.*`, "")
							// fmt.Println("-------------------------aa---------------------------------", item, aa)
							val4 = SumFloat(ParseFloat(aa), val4, "+")
							// fmt.Println("-------------------------val4---------------------------------", val4)
						} else {

							val1 += item + "；"
						}
					}
					val2 := SumFloat(ParseFloat(field2), val4, "-")
					// fmt.Println("-------------------------工本费|病历费|卡费|复印费|陪护费|陪人费|膳食费|陪床费-----------------------------------", val1, val2, val3, val4)
					for cc, field := range projectFields {
						if field.Code == nCode[0] {
							obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = val1
						}
						if field.Code == nCode[1] {
							obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = ToString(val2)
						}
						if field.Code == nCode[2] {
							obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = val3
						}
						if field.Code == nCode[3] {
							obj.Invoice[aa].InvoiceDaXiang[bb][cc].FinalValue = ToString(val4)
						}
					}

				}
			}
			// for _, invoiceDaXiangField := range projectFields {
			// }
		}

	}

	//16:36CSB0118RC0329000
	//结果数据
	//以下字段为1组，如只存在一个字段有结果值，且值为Z04.8时，出问题件，问题件编码：N0005，问题件描述：影像中不存在疾病诊断，
	//fc053、fc182、fc183、fc184、fc185、fc186、fc187、fc188、fc189、fc190
	locFlagIndex := 0
	var locFlag []int
	codeArr := []string{"fc053", "fc182", "fc183", "fc184", "fc185", "fc186", "fc187", "fc188", "fc189", "fc190"}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils2.GetFieldValueByLoc(obj, loc, true)
			if val == "Z04.8" {
				//locFlagIndex++
				locFlag = loc
			}
			if val != "" {
				locFlagIndex++
			}
		}
	}
	if locFlagIndex == 1 && len(locFlag) > 0 {
		utils2.SetIssue(obj, locFlag, "影像中不存在疾病诊断", "N0005", "")
	}

	//CSB0118RC0322000
	//当同一发票属性存在bc003分块时，当fc052的结果值与fc067结果值一致时，
	//出导出校验，账单号[xxx]报销金额加自费金额等于总金额，请修改。(XXXX为对应同一属性的fc002结果值内容)
	for _, codeLoc := range fieldLocationMap["fc052"] {
		fc052Val := utils2.GetFieldDecimalValueByLoc(obj, codeLoc, true)
		fc067Loc := getFieldLoc(fieldLocationMap["fc067"], codeLoc[0], codeLoc[1], codeLoc[2], codeLoc[4])
		fc002Loc := getFieldLoc(fieldLocationMap["fc002"], codeLoc[0], -1, -1, -1)
		if len(fc067Loc) < 1 || len(fc002Loc) < 1 {
			continue
		}
		fc067Val := utils2.GetFieldDecimalValueByLoc(obj, fc067Loc[0], true)
		fc002Val := utils2.GetFieldValueByLoc(obj, fc002Loc[0], true)
		if fc052Val.Equal(fc067Val) {
			obj.Bill.WrongNote += "账单号[" + fc002Val + "]报销金额加自费金额等于总金额，请修改;"
		}
	}

	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		//重新赋值
		//Id               relationConf.InvoiceFieldCode.FieldCode             //发票属性
		//Code             relationConf.InvoiceNumField            			   //账单号
		//Money            relationConf.MoneyField            				   //账单金额
		//InvoiceType      relationConf.InvoiceTypeField      				   //发票类型
		idField := RelationConf.InvoiceFieldCode.FieldCode[0]
		codeField := RelationConf.InvoiceNumField[0]
		moneyField := RelationConf.MoneyField[0]
		invoiceTypeField := RelationConf.InvoiceTypeField
		idFieldLoc := getFieldLoc(fieldLocationMap[idField], i, -1, 0, -1)
		codeFieldLoc := getFieldLoc(fieldLocationMap[codeField], i, -1, 0, -1)
		moneyFieldLoc := getFieldLoc(fieldLocationMap[moneyField], i, -1, 0, -1)
		invoiceTypeFieldLoc := getFieldLoc(fieldLocationMap[invoiceTypeField], i, -1, 0, -1)
		if len(idFieldLoc) > 0 && len(codeFieldLoc) > 0 && len(moneyFieldLoc) > 0 && len(invoiceTypeFieldLoc) > 0 {
			invoiceMap.Id = invoiceMap.Invoice[0][idFieldLoc[0][3]].FinalValue
			invoiceMap.Code = invoiceMap.Invoice[0][codeFieldLoc[0][3]].FinalValue
			invoiceMap.Money, _ = strconv.ParseFloat(invoiceMap.Invoice[0][moneyFieldLoc[0][3]].FinalValue, 64)
			invoiceMap.InvoiceType = invoiceMap.Invoice[0][invoiceTypeFieldLoc[0][3]].FinalValue
			obj.Bill.CountMoney, _ = decimal.NewFromFloat(invoiceMap.Money).Add(decimal.NewFromFloat(obj.Bill.CountMoney)).Float64()
		}
	}

	//CSB0118RC0332000
	//所有字段录入数据只要存在英文问号，不管存在多少个都出一个导出校验：“录入内容存在问号，请核实
	//CSB0118RC0332001 XX发票AA存在?，请修改“(XX取发票号，AA取字段名称)，如无存在发票属性内容时，直接出AA存在?，请修改
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
						msg := fields[l].Name + "存在?，请修改;"
						if invoiceMap.Code != "" {
							msg = invoiceMap.Code + "发票" + msg
						}
						if strings.Index(fields[l].ResultValue, "?") != -1 &&
							strings.Index(obj.Bill.WrongNote, msg) == -1 {
							obj.Bill.WrongNote += msg
						}
					}
				}
			}
		}
	}

	//CSB0118RC0334000
	//当fc062的值存在不一致时，存在不一致时，对一致内容的数量与不一致的数量进行对比，等一致内容的数量与不一致的数量相等时，出导出校验：请联系管理员确认责任外金额是否正确
	valArr := utils2.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc062"], false)
	fc062ValMap := make(map[string]int, 0)
	flagNum := 0  //记录最后一个fc062的数量 为了下面作比较有初始化的值
	flag := false //默认出校验
	for _, val := range valArr {
		if _, ok := fc062ValMap[val]; ok {
			fc062ValMap[val] = fc062ValMap[val] + 1
		} else {
			fc062ValMap[val] = 1
		}
		flagNum = fc062ValMap[val]
	}

	//flagMaxNumVal := "" //记录数量最多的fc062的值 CSB0118RC0335000
	//flagMaxNum := 0     //记录数量最多的fc062的数量 CSB0118RC0335000
	for _, v := range fc062ValMap {
		if flagNum != v || len(fc062ValMap) < 2 {
			//不出校验
			flag = true
		}
		//if flagMaxNum < v {
		//	flagMaxNum = v
		//	flagMaxNumVal = k
		//}
	}
	if !flag {
		obj.Bill.WrongNote += "请联系管理员确认责任外金额是否正确;"
	}

	//CSB0118RC0340000
	//当fc180录入为徐州时，fc067录入值不为0时，出导出校验，该案件为徐州单种病，自费金额录入值必须为0，不可回传
	//for _, loc := range fieldLocationMap["fc180"] {
	//	val := utils2.GetFieldValueByLoc(obj, loc, false)
	//	loc1 := getFieldLoc(fieldLocationMap["fc067"], loc[0], -1, -1, -1)
	//	val1 := utils2.GetFieldValueByLocArr(obj, loc1, false)
	//	if val == "徐州" && val1 != "0" {
	//		obj.Bill.WrongNote += "该案件为徐州单种病，自费金额录入值必须为0，不可回传;"
	//	}
	//}

	//CSB0118RC0343000
	//fc277录入值为1时
	//1、fc279 节点“billCode”的值为空时，出导出校验“票据代码不能为空”。
	//2、fc281 节点“billCheckCode”的值为空时，出导出校验“校验码不能为空”。
	//for _, loc := range fieldLocationMap["fc277"] {
	//	val := utils2.GetFieldValueByLoc(obj, loc, true)
	//	if val == "1" {
	//		locFC279 := getFieldLoc(fieldLocationMap["fc279"], loc[0], -1, -1, -1)
	//		//fc279Field := utils2.GetFieldByLocArr(obj, locFC279)
	//		val1 := utils2.GetFieldValueByLocArr(obj, locFC279, false)
	//		if val1 == "" {
	//			obj.Bill.WrongNote += "票据代码不能为空；"
	//		}
	//		locFC281 := getFieldLoc(fieldLocationMap["fc281"], loc[0], -1, -1, -1)
	//		//fc281Field := utils2.GetFieldByLocArr(obj, locFC281)
	//		valFC281 := utils2.GetFieldValueByLocArr(obj, locFC281, false)
	//		if valFC281 == "" {
	//			obj.Bill.WrongNote += "校验码不能为空"
	//		}
	//	}
	//}

	//CSB0118RC0335000
	//结果数据
	//判断fc062录入值是否存在不一致，存在不一致时，（采取少数服从多数的原则），将不一致的fc062对应发票下的第一列字段的结果值数据赋值给第二列数据、并将第三列字段的结果数据赋值为“出险人错误”。
	//myCode = [][]string{
	//	{"fc010", "fc105", "fc104"},
	//	{"fc012", "fc107", "fc106"},
	//	{"fc014", "fc109", "fc108"},
	//	{"fc016", "fc111", "fc110"},
	//	{"fc018", "fc113", "fc112"},
	//	{"fc020", "fc115", "fc114"},
	//	{"fc022", "fc117", "fc116"},
	//	{"fc024", "fc119", "fc118"},
	//	{"fc026", "fc121", "fc120"},
	//	{"fc028", "fc123", "fc122"},
	//	{"fc030", "fc125", "fc124"},
	//	{"fc032", "fc127", "fc126"},
	//	{"fc034", "fc129", "fc128"},
	//	{"fc036", "fc131", "fc130"},
	//	{"fc038", "fc133", "fc132"},
	//	{"fc040", "fc135", "fc134"},
	//	{"fc042", "fc137", "fc136"},
	//	{"fc044", "fc139", "fc138"},
	//	{"fc046", "fc141", "fc140"},
	//	{"fc048", "fc143", "fc142"},
	//}
	//for _, loc := range fieldLocationMap["fc062"] {
	//	val := utils2.GetFieldValueByLoc(obj, loc, false)
	//	if val == flagMaxNumVal || len(fc062ValMap) < 2 {
	//		continue
	//	}
	//	for _, code := range myCode {
	//		loc1 := getFieldLoc(fieldLocationMap[code[0]], loc[0], -1, -1, -1)
	//		loc2 := getFieldLoc(fieldLocationMap[code[1]], loc[0], -1, -1, -1)
	//		loc3 := getFieldLoc(fieldLocationMap[code[2]], loc[0], -1, -1, -1)
	//		if len(loc1) != 1 || len(loc2) != 1 || len(loc3) != 1 {
	//			continue
	//		}
	//		val1 := utils2.GetFieldValueByLoc(obj, loc1[0], true)
	//		if val1 != "" {
	//			utils2.SetFinalValue(obj, loc2, val1)
	//			utils2.SetFinalValue(obj, loc3, "出险人错误")
	//		}
	//	}
	//}

	// CSB0118RC0335000 同一发票下，当fc275结果值为1时，校验fc003的结果值：（字段不存在时，不执行该校验）
	//1.fc003结果值为1时，fc009结果值默认为100，fc010结果值默认为fc008的结果值，同时清空以下字段的结果值；
	//2.fc003结果值为2时，fc009结果值默认为200，fc010结果值默认为fc008的结果值，同时清空以下字段的结果值；
	//字段：fc011、fc012、fc013、fc014、fc015、fc016、fc017、fc018、fc019、fc020、fc021、fc022、fc023、fc024、fc025、fc026、fc027、fc028、fc029
	//fc030、fc031、fc032、fc033、fc034、fc035、fc036、fc037、fc038、fc039、fc040、fc041、fc042、fc043、fc044、fc045、fc046、fc047、fc048
	codes := []string{"fc011", "fc012", "fc013", "fc014", "fc015", "fc016", "fc017", "fc018", "fc019", "fc020", "fc021", "fc022", "fc023", "fc024", "fc025", "fc026", "fc027", "fc028", "fc029", "fc030", "fc031", "fc032", "fc033", "fc034", "fc035", "fc036", "fc037", "fc038", "fc039", "fc040", "fc041", "fc042", "fc043", "fc044", "fc045", "fc046", "fc047", "fc048"}
	for ii, _ := range fieldLocationMap["fc002"] {
		moneyFieldLoc := getFieldLoc(fieldLocationMap["fc275"], ii, -1, 0, -1)
		for i, codeLoc := range moneyFieldLoc {
			fc275Val := utils2.GetFieldValueByLoc(obj, codeLoc, true)
			if fc275Val == "1" {
				fc002Loc := getFieldLoc(fieldLocationMap["fc002"], i, -1, -1, -1)
				fc002Val := utils2.GetFieldValueByLoc(obj, fc002Loc[0], true)
				fmt.Println("fc002Val=", fc002Val)
				fc003Loc := getFieldLoc(fieldLocationMap["fc003"], ii, -1, -1, -1)
				if len(fc003Loc) < 1 {
					continue
				}
				fc003Val := utils2.GetFieldValueByLoc(obj, fc003Loc[0], true)
				fc008Loc := getFieldLoc(fieldLocationMap["fc008"], ii, -1, -1, -1)
				if len(fc008Loc) < 1 {
					continue
				}
				fc008Val := utils2.GetFieldValueByLoc(obj, fc008Loc[0], true)
				fc009Loc := getFieldLoc(fieldLocationMap["fc009"], ii, -1, -1, -1)
				fc010Loc := getFieldLoc(fieldLocationMap["fc010"], ii, -1, -1, -1)
				if fc003Val == "1" {
					utils2.SetFinalValue(obj, fc009Loc, "100")
					fc009Val := utils2.GetFieldValueByLoc(obj, fc009Loc[0], true)
					fmt.Println("==================fc009Val==========", fc009Val)
					utils2.SetFinalValue(obj, fc010Loc, fc008Val)
					for _, code := range codes {
						field := getFieldLoc(fieldLocationMap[code], ii, -1, -1, -1)
						utils2.SetFinalValue(obj, field, "")
					}
				} else if fc003Val == "2" {
					utils2.SetFinalValue(obj, fc009Loc, "200")
					utils2.SetFinalValue(obj, fc010Loc, fc008Val)
					for _, code := range codes {
						field := getFieldLoc(fieldLocationMap[code], ii, -1, -1, -1)
						utils2.SetFinalValue(obj, field, "")
					}
				}
				//CSB0118RC0337000 同一发票下，当fc275结果值为1时，fc067结果值默认0（如字段不存在则不执行该校验）
				// if len(fc067Loc) > 0 {
				// 	utils2.SetFinalValue(obj, fc067Loc, "0")

				// }
				// if len(fc211Loc) > 0 {
				// 	utils2.SetFinalValue(obj, fc211Loc, fc008Val)

				// }
			}
		}
		//for i, codeLoc := range fieldLocationMap["fc275"] {
		//
		//}
	}
	// CSB0118RC0338000
	//同一发票下，当MB002-bc002清单分块下fc084、fc085、fc086、fc087、fc088、fc089、fc090、fc091清单项目名称字段匹配的常量表为《B0118_中意理赔_省份-全国》时
	//出导出校验：不存在xxx地区目录库，已匹配全国目录库，请联系管理员确认；（xxx为同一发票下fc181的结果值，当fc181结果值为F时，或不存在清单分块时，不执行该校验）
	sFields := []string{"fc084", "fc085", "fc086", "fc087", "fc088", "fc089", "fc090", "fc091"}
	for i, codeLoc := range fieldLocationMap["fc002"] {
		if len(codeLoc) > 0 {
			for _, field := range sFields {
				Loc := getFieldLoc(fieldLocationMap[field], i, -1, -1, -1)
				fc180Loc := getFieldLoc(fieldLocationMap["fc180"], i, -1, -1, -1)
				fc181Loc := getFieldLoc(fieldLocationMap["fc181"], i, -1, -1, -1)
				if len(Loc) > 0 {
					locVal := utils2.GetFieldValueByLoc(obj, Loc[0], false)
					fc180Val := utils2.GetFieldValueByLoc(obj, fc180Loc[0], false)
					fc181Val := utils2.GetFieldValueByLoc(obj, fc181Loc[0], false)
					if locVal != "" {
						ok := false
						addr, ok := constSpecialMap[fc180Val]
						province := fc180Val
						if !ok {
							addr, ok = constSpecialMap[fc181Val]
							province = fc181Val
							if !ok {
								province = "全国"
								addr, ok = constSpecialMap[province]
								val, ok := addr[locVal]
								fmt.Println("==========val=========", val)
								if fc181Val != "F" && ok {
									msg := "不存在" + fc181Val + "地区目录库，已匹配全国目录库，请联系管理员确认；"
									if strings.Index(obj.Bill.WrongNote, msg) == -1 {
										obj.Bill.WrongNote += msg
										fmt.Println("=fc181Val=", fc181Val)
									}
								}
							}
						}
					}

				}
			}
		}
	}

	//CSB0118RC0337000
	//fc276的结果值有内容时，增加该字段的问题件描述，每个分号（;或；）之间作为一个问题件描述，没有分号则视为一个问题件描述，描述内容根据填写内容进行生成，问题件编码均默认为N0005；
	//(如：“缺少发票；”则视为1个问题件；“缺少发票；缺少发票”则视为2个问题件描述)(将该代码放在所有需求的最后面。)
	for _, loc := range fieldLocationMap["fc276"] {
		fc276Val := utils2.GetFieldValueByLoc(obj, loc, true)
		if fc276Val == "" {
			continue
		}
		fc276Val = strings.ReplaceAll(fc276Val, "；", ";")
		for _, msg := range strings.Split(fc276Val, ";") {
			utils2.SetIssue(obj, loc, msg, "N0005", "")
		}
	}

	////////////////////////项目报表 - 业务明细 - 单据状态
	obj.Bill.BillType = utils.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	//CSB0118RC0342000
	//"发票查验功能：结果数据页面新增“发票查验”功能按键
	//一、在对应发票页点击该功能按键，根据当前发票页的fc278结果值将以下对应字段的信息发送到第三方发票查验接口进行查验：
	//1.fc278录入值为1时，发送以下字段结果值到对应参数进行查验得到del的节点值。
	//fc279 fpdm
	//fc002 fphm
	//fc280 kprq
	//fc281 checkCode
	//fc008 noTaxAmount
	//2.fc278录入值为2时，发送以下字段结果值到对应参数进行查验得到isRed和isPrint的节点值。
	//fc279 fpdm
	//fc002 fphm
	//fc280 kprq
	//fc281 checkCode
	//fc008 money
	//二、fc282结果值根据发票查验结果进行输出：
	//1.fc278录入值为1时，根据发票查验del节点值进行转码：
	//del  fc282
	//0      Y
	//2      Y
	//3      Y
	//7      Y
	//8      Y
	//del值为0时出问题件，编码为N0013，问题件描述：电子发票（NO.***********、NO.*******) 已查验、已打印。
	//del值为2时出问题件，编码为N0014，问题件描述：电子发票（NO.***********、NO.*******)已作废
	//del值为3、7、8时出问题件，编码为N0015，问题件描述：电子发票（NO.***********、NO.*******)已开红票
	//2.fc278录入值为2时，根据发票查验isRed和isPrint节点值进行转码：
	//①当isRed为true时，fc282输出N；并出问题件，问题件编码为N0015，问题件描述：电子发票（NO.***********、NO.*******)已开红票
	//②当isRed为false时，则fc282输出Y；并出问题件，编码为N0013，问题件描述：电子发票（NO.***********、NO.*******) 已查验、已打印。
	//3.当查验失败时（即无法拿到【del】或【isRed和isPrint】节点值），fc282输出"N"，
	//校验code的值，
	//当错误码为“1099”时出问题件，编码为N0016，问题件描述：电子发票（NO.***********、NO.*******)票据不存在；
	//当错误码为“800”、“801”、“1009”、“1011”、“2000”“2001”等查验失败时出问题件，编码为N0017，问题件描述：电子发票（NO.***********、NO.*******)未成功查验，原因：******
	//（如：NO.123456789未成功查验，原因：800，税局服务异常，建议 15-20 分钟后重试）
	//其中（NO.***********、NO.*******)为fc002录入内容，假如只有一张单只需要出一个编号，多张单需出多个编号。
	//三、发票查验提示
	//1.成功（即能正常拿到【del】或【isRed和isPrint】节点值）需弹出提示：发票查验成功；
	//2.失败需根据code、message及description三个节点值的内容用“逗号”隔开弹出提示；（如：801，查验失败具体提示，税局服务暂时不可用）"

	//CSB0118RC0343000
	//同一发票下，当fc277录入值为1时，校验fc279、fc002、fc280、fc281、fc008是否存在结果值，任意一个字段结果值为空时，出导出校验：【XX】电子发票五要素未录入齐全，请确认是否为电子发票；（xxx为fc002的结果值）

	//CSB0118RC0345000
	//同一发票下，当fc278录入值为1时，（xxx为fc002的结果值）
	//1、校验fc279的结果值是否为12位数，否则出导出校验：【XX】发票增值税票据代码字段应为12位数，请检查；
	//2、校验fc002的结果值是否为8位数，否则出导出校验：【XX】发票增值税票据号码字段应为8位数，请检查；
	//3、校验fc281的结果值是否为6位数，否则出导出校验：【XX】发票校验码字段应为6位数，请检查；

	//CSB0118RC0343000
	//同一发票下，当fc278录入值为2时，（xxx为fc002的结果值）
	//1、校验fc279的结果值是否为8位数，否则出导出校验：【XX】发票财政部票据代码字段应为8位数，请检查；
	//2、校验fc002的结果值是否为10位数，否则出导出校验：【XX】发票财政部票据号码字段应为10位数，请检查；

	//CSB0118RC0346000
	//同一发票下，当fc277录入值为1，且fc282的结果值不为“Y”时，出导出校验，【XX】发票查验结果不通过，请确认五要素是否录入正确；（xxx取同一发票下fc002的值）

	//CSB0118RC0338000
	//发票查验提示（该需求代码位置放在第一位）
	//失败需根据code、message及description三个节点值的内容用“逗号”隔开出对应导出校验：xxx发票查验失败，原因为：code，message，description；（xxx为同一发票下fc002的结果值）
	//（如：123456发票查验失败，原因为：801，查验失败具体提示，税局服务暂时不可用）
	mes := ""
	isNotCheck := true
	changeMap := map[string][]string{
		"0": {"Y", "N0013", "已查验、已打印。"},
		"2": {"Y", "N0014", "已作废。"},
		"3": {"Y", "N0015", "已开红票。"},
		"7": {"Y", "N0015", "已开红票。"},
		"8": {"Y", "N0015", "已开红票。"},
	}
	changeMap1 := map[bool][]string{
		true:  {"N", "N0015", "已开红票。"},
		false: {"Y", "N0013", "已查验、已打印。"},
	}
	changeErrMap1 := map[int][]string{
		1099: {"N0016", "票据不存在。"},
		800:  {"N0017", "未成功查验，原因："},
		801:  {"N0017", "未成功查验，原因："},
		1009: {"N0017", "未成功查验，原因："},
		1011: {"N0017", "未成功查验，原因："},
		2000: {"N0017", "未成功查验，原因："},
		2001: {"N0017", "未成功查验，原因："},
	}
	//CSB0118RC0346000
	mergeInvoiceList := []string{} //CSB0118RC0346000 定义数组 接收多个结果
	for i, _ := range obj.Invoice {
		fc278Loc := utils2.GetFieldLoc(fieldLocationMap["fc278"], i, -1, -1, -1)
		fc279Loc := utils2.GetFieldLoc(fieldLocationMap["fc279"], i, -1, -1, -1)
		fc002Loc := utils2.GetFieldLoc(fieldLocationMap["fc002"], i, -1, -1, -1)
		fc280Loc := utils2.GetFieldLoc(fieldLocationMap["fc280"], i, -1, -1, -1)
		fc281Loc := utils2.GetFieldLoc(fieldLocationMap["fc281"], i, -1, -1, -1)
		fc008Loc := utils2.GetFieldLoc(fieldLocationMap["fc008"], i, -1, -1, -1)
		fc282Loc := utils2.GetFieldLoc(fieldLocationMap["fc282"], i, -1, -1, -1)
		fc277Loc := utils2.GetFieldLoc(fieldLocationMap["fc277"], i, -1, -1, -1)
		fc278Field := utils2.GetFieldByLocArr(obj, fc278Loc)
		fc279Field := utils2.GetFieldByLocArr(obj, fc279Loc)
		fc002Field := utils2.GetFieldByLocArr(obj, fc002Loc)
		fc280Field := utils2.GetFieldByLocArr(obj, fc280Loc)
		fc281Field := utils2.GetFieldByLocArr(obj, fc281Loc)
		fc008Field := utils2.GetFieldByLocArr(obj, fc008Loc)
		fc282Field := utils2.GetFieldByLocArr(obj, fc282Loc)
		fc277Field := utils2.GetFieldByLocArr(obj, fc277Loc)
		if fc278Field.ResultValue != "" &&
			(bill.Stage == 6 ||
				(bill.Stage != 6 && strings.Contains(bill.Remark, "<<发票查验>>")) && (fc278Field.IsChange || fc279Field.IsChange || fc002Field.IsChange || fc280Field.IsChange || fc281Field.IsChange || fc008Field.IsChange || fc282Field.IsChange)) {
			bodyData := make(map[string]interface{})
			bodyData["fphm"] = fc002Field.FinalValue
			bodyData["kprq"] = fc280Field.FinalValue
			fc282FinalValue := ""
			issue := model3.Issue{}
			if fc278Field.ResultValue == "1" {
				bodyData["noTaxAmount"] = fc008Field.FinalValue
				bodyData["fpdm"] = fc279Field.FinalValue
				bodyData["checkCode"] = fc281Field.FinalValue
				err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData)
				isNotCheck = false
				if err != nil {
					mes += fc002Field.ResultValue + fmt.Sprintf("%v", err) + ";"
					mes = strings.ReplaceAll(mes, "参数错误：", "") //CSB0118RC0350000
					fc282FinalValue = "N"

					if _, ok := changeErrMap1[respData.Code]; ok {
						issue.Code = changeErrMap1[respData.Code][0]
						issue.Message = "电子发票（NO." + fc002Field.ResultValue + "）" + changeErrMap1[respData.Code][1] + respData.Message + "。"
						issue.Message = strings.ReplaceAll(issue.Message, "参数错误：", "") //CSB0118RC0350000
					}
				} else if val, ok := changeMap[respData.Data.Del]; ok {
					fc282FinalValue = val[0]

					//del值为0时出问题件，编码为N0013，问题件描述：电子发票（NO.***********、NO.*******) 已查验、已打印。
					//del值为2时出问题件，编码为N0014，问题件描述：电子发票（NO.***********、NO.*******)已作废
					//del值为3、7、8时出问题件，编码为N0015，问题件描述：电子发票（NO.***********、NO.*******)已开红票
					issue.Code = val[1]
					issue.Message = "电子发票（NO." + fc002Field.ResultValue + "）" + val[2]
				}
			} else if fc278Field.ResultValue == "2" {
				bodyData["money"] = fc008Field.FinalValue
				bodyData["fpdm"] = fc279Field.FinalValue
				bodyData["checkCode"] = fc281Field.FinalValue
				err, respData := unitFunc.Invoice("/v2/eInvoice/query", bodyData)
				isNotCheck = false
				if err != nil {
					mes += fc002Field.ResultValue + fmt.Sprintf("%v", err) + ";"
					mes = strings.ReplaceAll(mes, "参数错误：", "") //CSB0118RC0350000
					fc282FinalValue = "N"

					if _, ok := changeErrMap1[respData.Code]; ok {
						issue.Code = changeErrMap1[respData.Code][0]
						issue.Message = "电子发票（NO." + fc002Field.ResultValue + "）" + changeErrMap1[respData.Code][1] + respData.Message + "。"
						issue.Message = strings.ReplaceAll(issue.Message, "参数错误：", "") //CSB0118RC0350000
					}
				} else if val, ok := changeMap1[respData.Data.IsRed]; ok {
					fc282FinalValue = val[0]

					//①当isRed为true时，fc282输出3；并出问题件，问题件编码为N0015，问题件描述：电子发票（NO.***********、NO.*******)已开红票
					//②当isRed为false时，则fc282输出1；并出问题件，编码为N0013，问题件描述：电子发票（NO.***********、NO.*******) 已查验、已打印。
					issue.Code = val[1]
					issue.Message = "电子发票（NO." + fc002Field.ResultValue + "）" + val[2]
				}
			} else if fc278Field.ResultValue == "3" { //CSB0118RC0342001 2024/05/13
				bodyData["jshj"] = fc008Field.FinalValue
				err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData)
				isNotCheck = false
				if err != nil {
					mes += fc002Field.ResultValue + fmt.Sprintf("%v", err) + ";"
					mes = strings.ReplaceAll(mes, "参数错误：", "") //CSB0118RC0350000
					fc282FinalValue = "N"

					if _, ok := changeErrMap1[respData.Code]; ok {
						issue.Code = changeErrMap1[respData.Code][0]
						issue.Message = "电子发票（NO." + fc002Field.ResultValue + "）" + changeErrMap1[respData.Code][1] + respData.Message + "。"
						issue.Message = strings.ReplaceAll(issue.Message, "参数错误：", "") //CSB0118RC0350000
					}
				} else if val, ok := changeMap[respData.Data.Del]; ok {
					fc282FinalValue = val[0]

					//del值为0时出问题件，编码为N0013，问题件描述：电子发票（NO.***********、NO.*******) 已查验、已打印。
					//del值为2时出问题件，编码为N0014，问题件描述：电子发票（NO.***********、NO.*******)已作废
					//del值为3、7、8时出问题件，编码为N0015，问题件描述：电子发票（NO.***********、NO.*******)已开红票
					issue.Code = val[1]
					issue.Message = "电子发票（NO." + fc002Field.ResultValue + "）" + val[2]
				}
			}

			utils2.SetOnlyOneResultValue(obj, fc282Loc[0], issue.Code+"_"+issue.Message+"_"+fc282FinalValue)
			obj.Invoice[i].Invoice[fc282Loc[0][2]][fc282Loc[0][3]].IsChange = true
		}

		issueArr := strings.Split(utils2.GetFieldValueByLocArr(obj, fc282Loc, false), "_")
		if len(issueArr) == 3 && issueArr[1] != "" && issueArr[0] != "" {
			utils2.SetIssues(obj, fc282Loc, issueArr[1], issueArr[0], "")
			utils2.SetFinalValue(obj, fc282Loc, issueArr[2])
		}

		//CSB0118RC0343000
		//CSB0118RC0349000 2024/05/16
		if fc277Field.ResultValue == "1" &&
			(fc279Field.FinalValue == "" || fc002Field.FinalValue == "" || fc280Field.FinalValue == "" || fc008Field.FinalValue == "") {
			if fc278Field.ResultValue == "3" && (fc279Field.FinalValue == "" && fc281Field.FinalValue == "") {

			} else {
				obj.Bill.WrongNote += "【" + fc002Field.FinalValue + "】电子发票五要素未录入齐全，请确认是否为电子发票；"
			}
		}

		//CSB0118RC0345000
		//CSB0118RC0343000
		//CSB0118RC0347000 2024/05/23 需求编码更正
		iArr := [][]interface{}{
			{"1", fc279Field.FinalValue, 12, fc279Field.Name},
			{"1", fc002Field.FinalValue, 8, fc002Field.Name},
			{"1", fc281Field.FinalValue, 6, fc281Field.Name},
			{"2", fc279Field.FinalValue, 8, fc279Field.Name},
			{"2", fc002Field.FinalValue, 10, fc002Field.Name},
		}
		for _, item := range iArr {
			if fc278Field.ResultValue == item[0].(string) && len(item[1].(string)) != item[2].(int) {
				obj.Bill.WrongNote += fmt.Sprintf("【%v】发票%v字段应为%d位数，请检查；", fc002Field.FinalValue, item[3], item[2])
			}
		}

		//CSB0118RC0346000
		if fc277Field.ResultValue == "1" && utils2.GetFieldValueByLocArr(obj, fc282Loc, true) != "Y" {
			//obj.Bill.WrongNote += fmt.Sprintf("【%v】发票查验结果不通过，请确认五要素是否录入正确；", fc002Field.FinalValue)
			mergeInvoiceList = append(mergeInvoiceList, fc002Field.FinalValue)
		}
	}
	//CSB0118RC0346000
	if len(mergeInvoiceList) > 0 {
		obj.Bill.WrongNote += fmt.Sprintf("【%v】发票查验结果不通过，请确认五要素是否录入正确；", strings.Join(mergeInvoiceList, "，"))
	}

	if bill.Stage != 6 && isNotCheck && strings.Contains(bill.Remark, "<<发票查验>>") {
		mes += "查验失败：发票查验相关字段的数据未作修改，不执行发票查验功能;"
	}
	global.GLog.Info("------------mes-------------------:", zap.Any("", mes))
	if mes != "" {
		obj.Bill.WrongNote += mes
	}

	return nil, obj
}

func clearfcValue(fields []model3.ProjectField) {
	for ff, field := range fields {
		if field.Code == "fc170" || field.Code == "fc171" {
			fields[ff].FinalValue = ""
		}
	}
}

func getOneValue(fields []model3.ProjectField, code string) string {
	for _, field := range fields {
		if field.Code == code {
			return field.FinalValue
		}
	}
	return ""
}

func getOneResultValue(fields []model3.ProjectField, code string) string {
	for _, field := range fields {
		if field.Code == code {
			return field.ResultValue
		}
	}
	return ""
}

func setFc17Value(fields []model3.ProjectField, fc170, fc171 string) string {
	for ff, field := range fields {
		if field.Code == "fc170" {
			fields[ff].FinalValue = fields[ff].FinalValue + "；" + fc170
		}
		if field.Code == "fc171" && fc171 != "" {
			aa, _ := decimal.NewFromString(fc171)
			bb, err := decimal.NewFromString(fields[ff].FinalValue)
			if err != nil {
				bb = decimal.NewFromFloat(0.0)
			}
			fields[ff].FinalValue = bb.Add(aa).String()
		}
	}
	return ""
}

func getFc171Total(fieldLocationMap map[string][][]int, i int, invoiceMap model2.InvoiceMap) decimal.Decimal {
	locs := getFieldLoc(fieldLocationMap["fc171"], i, -1, -1, -1)
	var fc171Total = decimal.Zero
	for _, loc := range locs {
		v, err := decimal.NewFromString(invoiceMap.QingDan[loc[2]][loc[3]].FinalValue)
		if err != nil {
			global.GLog.Error("getFc171Total" + err.Error())
			continue
		}
		fc171Total = fc171Total.Add(v)
	}
	return fc171Total
}

// 获取一个字段所有值的和
func getAllVal(fields []model3.ProjectField, code string) decimal.Decimal {
	t := decimal.NewFromFloat(0.0)
	val := getValInFieldsByCode(fields, code)
	for _, v := range val {
		if v == "" {
			continue
		}
		d, err := decimal.NewFromString(v)
		if err != nil {
			continue
		}
		t = t.Add(d)
	}
	return t
}

// 判断fc008包含问号且那一堆字段是否有一个问号
func hasQuestionMark(codeArr []string, fieldLocationMap map[string][][]int, i int, invoiceMap model2.InvoiceMap, locs [][]int) bool {
	fc008Field := invoiceMap.Invoice[0][locs[0][3]]
	if strings.Index(fc008Field.ResultValue, "?") != -1 {
		for _, s := range codeArr {
			if len(fieldLocationMap[s]) > 0 {
				fc009Locs := getFieldLoc(fieldLocationMap[s], i, -1, -1, -1)
				for _, loc := range fc009Locs {
					fc009Field := invoiceMap.InvoiceDaXiang[loc[2]][loc[3]]
					if strings.Index(fc009Field.ResultValue, "?") != -1 {
						return true
					}
				}

			}
		}
	}
	return false
}

// 字段数组里面获取同一小分块的结果值
func getValInFieldsByCodeIndex(fields []model3.ProjectField, code string, blockIndex int) string {
	for _, f := range fields {
		if f.Code == code && f.BlockIndex == blockIndex {
			return f.FinalValue
		}
	}
	return ""
}

// 根据code获取字段的结果值
func getValInFieldsByCode(fields []model3.ProjectField, code string) []string {
	var val []string
	for _, f := range fields {
		if f.Code == code {
			val = append(val, f.FinalValue)
		}
	}
	return val
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"jiBingZhenDuanMap", "B0118_中意理赔_疾病诊断", "1", "0"},
		{"jiBingZhenDuanDaiMaMap", "B0118_中意理赔_疾病诊断", "0", "1"},
		{"yiYuanDaiMaMap", "B0118_中意理赔_医院代码表", "1", "0"},
		{"shouShuBianMaMap", "B0118_中意理赔_手术编码", "1", "0"},
		{"zhuYuanFeiYongMap", "B0118_中意理赔_住院费用类型", "1", "0"},
		{"zhuYuanFeiYongDaiMaMap", "B0118_中意理赔_住院费用类型", "0", "1"},
		{"guangXiZhuangZuMap", "B0118_中意理赔_省份-广西壮族", "2", "6"},
		{"menZhenFeiYongLeiXingMap", "B0118_中意理赔_门诊费用类型", "1", "0"},
		{"menZhenFeiYongLeiXingDaiMaMap", "B0118_中意理赔_门诊费用类型", "0", "1"},
		{"provinceMap", "B0118_中意理赔_省份-全国", "2", "1"},
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

	//地址
	//for k, v := range global.GProConf[proCode].ConstTable {
	//	if strings.Index(k, "省份") != -1 {
	//		tempMap := make(map[string]string, 0)
	//		for _, arr := range v {
	//			if len(arr) < 2 {
	//				continue
	//			}
	//			//获取最小的一项
	//			if val, ok := tempMap[arr[1]]; ok {
	//				a, err := strconv.ParseFloat(val, 64)
	//				if err != nil {
	//					continue
	//				}
	//				b, err := strconv.ParseFloat(arr[0], 64)
	//				if err != nil {
	//					continue
	//				}
	//				if a < b {
	//					arr[0] = val
	//				}
	//			}
	//			tempMap[arr[1]] = arr[0]
	//		}
	//		constObj[strings.Replace(k, "B0118_中意理赔_省份-", "", -1)] = tempMap
	//	}
	//}
	return constObj
}

// 初始化常量
func constSpecialDeal(proCode string) map[string]map[string][]string {
	constObj := make(map[string]map[string][]string, 0)
	for k, v := range global.GProConf[proCode].ConstTable {
		if strings.Index(k, "省份") != -1 {
			tempMap := make(map[string][]string, 0)
			for _, arr := range v {
				if strings.Index(k, "head") == -1 {
					tempMap[strings.TrimSpace(arr[2])] = arr
				} else {
					tempMap[strings.TrimSpace("head")] = arr
				}
			}
			constObj[strings.Replace(k, "B0118_中意理赔_省份-", "", -1)] = tempMap
		}
	}
	return constObj
}

// 获取同一发票的某一字段的位置 i:第几张发票，j:第几个对象（Invoice，QingDan。。。）
func getFieldLoc(locs [][]int, i, j, k, blockIndex int) [][]int {
	locsNew := make([][]int, 0)
	if len(locs) < 1 {
		return locsNew
	}
	for _, loc := range locs {
		if i != -1 {
			if loc[0] != i {
				continue
			}
		}
		if j != -1 {
			if loc[1] != j {
				continue
			}
		}
		if k != -1 {
			if loc[2] != k {
				continue
			}
		}
		if blockIndex != -1 {
			if loc[4] != blockIndex {
				continue
			}
		}
		locsNew = append(locsNew, loc)
	}
	return locsNew
}
