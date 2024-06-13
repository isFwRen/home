package B0122

import (
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	model2 "server/module/export/model"

	// "server/module/export/service"
	eUtils "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	unitFunc "server/module/unit"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

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
	InvoiceFieldCode:          model2.TypeCode{FieldCode: []string{"fc156", "fc157"}, BlockCode: []string{"bc017", "bc018"}},
	QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc158"}, BlockCode: []string{"bc019"}},
	BaoXiaoDanFieldCode:       model2.TypeCode{FieldCode: []string{"fc159"}, BlockCode: []string{"bc020"}},
	ThirdBaoXiaoDan1FieldCode: model2.TypeCode{FieldCode: []string{"fc160"}, BlockCode: []string{"bc021"}},
	OtherTempType:             map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "10": "10", "11": "11", "12": "12", "13": "13", "14": "14", "15": "15", "16": "16", "17": "17", "18": "18", "19": "19", "20": "20", "21": "21", "22": "22", "23": "23", "24": "24", "25": "25", "26": "26", "31": "31", "32": "32", "33": "33", "34": "34", "35": "35", "36": "36", "37": "37", "38": "38", "39": "39", "40": "40", "9": "9"},
	TempTypeField:             "fc153",
	InvoiceNumField:           []string{"fc079", ""},
	MoneyField:                []string{"fc080", ""},
	InvoiceTypeField:          "fc279",
	//InvoiceTypeField: "fc003",
}

// ResultData B0108
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {
	obj = eUtils.RelationDealManyInvoiceType(bill, blocks, fieldMap, relationConf)
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

	clearFields := false

	// fieldCodeS := []string{"fc003", "fc004", "fc013", "fc014", "fc027", "fc028", "fc054", "fc055", "fc056", "fc081", "fc082", "fc189", "fc190", "fc191", "fc192", "fc193", "fc197", "fc198", "fc210", "fc212", "fc213", "fc217", "fc225", "fc226", "fc235", "fc236", "fc245", "fc246", "fc255", "fc256", "fc278", "fc202", "fc203", "fc019", "fc029", "fc227", "fc237", "fc247", "fc257"}
	for code, fieldLocs := range fieldLocationMap {
		// fieldLocs := fieldLocationMap[fieldCode]
		for _, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					// fieldsArr[loc[2]][loc[3]].ResultValue
					if utils.RegIsMatch(`^(fc003|fc004|fc013|fc014|fc027|fc028|fc054|fc055|fc056|fc081|fc082|fc189|fc190|fc191|fc192|fc193|fc197|fc198|fc210|fc212|fc213|fc217|fc225|fc226|fc235|fc236|fc245|fc246|fc255|fc256|fc278|fc202|fc203|fc019|fc029|fc227|fc237|fc247|fc257|fc473|fc061)$`, code) {
						if !utils.RegIsMatch(`^(A|)$`, fieldsArr[loc[2]][loc[3]].ResultValue) && strings.Index(fieldsArr[loc[2]][loc[3]].ResultValue, "?") == -1 && len(fieldsArr[loc[2]][loc[3]].ResultValue) == 6 {
							fieldsArr[loc[2]][loc[3]].FinalValue = "20" + fieldsArr[loc[2]][loc[3]].FinalValue
						}
					}
					// CSB0122RC0048001 2024-04-28调整 新增字段fc497|fc498|fc499|fc500|fc501|fc502|fc503|fc504|fc505|fc506
					if utils.RegIsMatch(`^(fc164|fc165|fc166|fc167|fc168|fc169|fc170|fc171|fc172|fc173|fc497|fc498|fc499|fc500|fc501|fc502|fc503|fc504|fc505|fc506)$`, code) {
						if utils.HasKey(constMap["jiBingDaiMaMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["jiBingDaiMaMap"][fieldsArr[loc[2]][loc[3]].ResultValue]
						}
						if strings.Index(fieldsArr[loc[2]][loc[3]].ResultValue, "?") != -1 {
							fieldsArr[loc[2]][loc[3]].FinalValue = "R69.X51"
						}
					}

					// if utils.RegIsMatch(`^(fc164|fc165|fc166|fc167|fc168|fc169|fc170|fc171|fc172|fc173)$`, code) {
					// 	if strings.Index(fieldsArr[loc[2]][loc[3]].ResultValue, "?") != -1 {
					// 		fieldsArr[loc[2]][loc[3]].FinalValue = "R69.X51"
					// 	}
					// }

					if utils.RegIsMatch(`^(fc184|fc185|fc186|fc187|fc188)$`, code) {
						if utils.HasKey(constMap["shouShuDaiMaMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["shouShuDaiMaMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						}
					}

					// CSB0122RC0069000
					if utils.RegIsMatch(`^(fc022|fc032|fc230|fc240|fc250|fc260|fc263)$`, code) && fieldsArr[loc[2]][loc[3]].ResultValue != "" {
						if utils.HasKey(constMap["yinHangDaiMaMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["yinHangDaiMaMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						} else {
							fieldsArr[loc[2]][loc[3]].FinalValue = ""
							issue := model3.Issue{
								Type:    "",
								Code:    "C04",
								Message: "开户银行不在列表中",
							}
							obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues, issue)
						}
					}

					// CSB0122RC0068000
					if utils.RegIsMatch(`^(fc023|fc033|fc231|fc241|fc251|fc261)$`, code) {
						if utils.HasKey(constMap["kaiHuChengShiMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["kaiHuChengShiMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						} else if fieldsArr[loc[2]][loc[3]].ResultValue != "" {
							issue := model3.Issue{
								Type:    "",
								Code:    "C05",
								Message: "开户城市不在列表中",
							}
							obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues, issue)
						}
					}

					if utils.RegIsMatch(`^(fc211)$`, code) {
						if utils.HasKey(constMap["zhongJiMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["zhongJiMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						} else {
							fieldsArr[loc[2]][loc[3]].FinalValue = ""
						}
					}

					if utils.RegIsMatch(`^(fc214)$`, code) && !RegIsMatch(fieldsArr[loc[2]][loc[3]].ResultValue, `^(A|)$`) {
						if utils.HasKey(constMap["yiWaiMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["yiWaiMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						} else {
							fieldsArr[loc[2]][loc[3]].FinalValue = "Y10-Y34"
						}
					}
					//调整 ?或录入B时，输出空  编码 CSB0122RC0058000
					if utils.RegIsMatch(`^(fc213)$`, code) {
						if strings.Index(fieldsArr[loc[2]][loc[3]].ResultValue, "?") != -1 {
							obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = ""
							issue := model3.Issue{
								Type:    "",
								Code:    "C29",
								Message: "无法判断意外日期",
							}
							obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues, issue)
						}
					}

					if utils.RegIsMatch(`^(fc216)$`, code) && !RegIsMatch(fieldsArr[loc[2]][loc[3]].ResultValue, `^(A|)$`) {
						if utils.HasKey(constMap["yiWaiDiMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["yiWaiDiMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						} else {
							fieldsArr[loc[2]][loc[3]].FinalValue = "CS31"
						}
					}

					if utils.RegIsMatch(`^(fc218)$`, code) {
						if utils.HasKey(constMap["shenGuMap"], fieldsArr[loc[2]][loc[3]].ResultValue) {
							fieldsArr[loc[2]][loc[3]].FinalValue = constMap["shenGuMap"][fieldsArr[loc[2]][loc[3]].FinalValue]
						} else {
							fieldsArr[loc[2]][loc[3]].FinalValue = ""
						}
					}

					// CSB0122RC0067000
					if utils.RegIsMatch(`^(fc031|fc229|fc239|fc249|fc259)$`, code) {
						if strings.Index(fieldsArr[loc[2]][loc[3]].ResultValue, "?") != -1 {
							issue := model3.Issue{
								Type:    "",
								Code:    "C03",
								Message: "银行卡/存折缺失或不清晰",
							}
							obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues, issue)
							obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].FinalValue = ""
						}
					}
					// CSB0122RC0072000
					if utils.RegIsMatch(`^(fc144|fc339|fc340|fc341|fc342|fc343|fc344|fc345|fc143|fc331|fc332|fc333|fc334|fc335|fc336|fc337)$`, code) {
						if fieldsArr[loc[2]][loc[3]].ResultValue != "" {
							resultValue := utils.ParseFloat(fieldsArr[loc[2]][loc[3]].ResultValue) * 100
							fieldsArr[loc[2]][loc[3]].FinalValue = utils.ToString(resultValue) + "%"
						}
					}

					//CSB0122RC0075000
					if utils.RegIsMatch(`^(fc142|fc324|fc325|fc326|fc327|fc328|fc329|fc330|fc462|fc463|fc464|fc465|fc466|fc467|fc468|fc469|fc140|fc310|fc311|fc312|fc313|fc314|fc315|fc316|fc148|fc367|fc368|fc369|fc370|fc371|fc372|fc373|fc149|fc374|fc375|fc376|fc377|fc378|fc379|fc380|fc151|fc388|fc389|fc390|fc391|fc392|fc393|fc394|fc152|fc395|fc396|fc397|fc398|fc399|fc400|fc401|fc141|fc317|fc318|fc319|fc320|fc321|fc322|fc323)$`, code) {
						value := fieldsArr[loc[2]][loc[3]].FinalValue
						if value == "" {
							value = "0"
						}
						fieldsArr[loc[2]][loc[3]].FinalValue = utils.ToFix2(value)
					}

					// CSB0122RC0139000
					if utils.RegIsMatch(`^(fc139|fc303|fc304|fc305|fc306|fc307|fc308|fc309)$`, code) {
						if utils.RegIsMatch(`(狂犬疫苗|狂犬病疫苗)`, fieldsArr[loc[2]][loc[3]].ResultValue) {
							clearFields = true
						}
					}
					//2024-04-28 CSB0122RC0139001 调整 新增字段：fc497、fc498、fc499、fc500、fc501、fc502、fc503、fc504、fc505、fc506
					if utils.RegIsMatch(`^(fc045|fc164|fc165|fc166|fc167|fc168|fc169|fc170|fc171|fc172|fc173|fc497|fc498|fc499|fc500|fc501|fc502|fc503|fc504|fc505|fc506)$`, code) {
						if utils.RegIsMatch(`(狗咬|猫咬|狗抓|猫抓|鼠咬|鼠抓|狂犬疫苗|狂犬病疫苗)`, fieldsArr[loc[2]][loc[3]].ResultValue) {
							clearFields = true
						}
					}
					if utils.RegIsMatch(`^(fc048)$`, code) {
						if utils.RegIsMatch(`(狂犬|犬伤)`, fieldsArr[loc[2]][loc[3]].ResultValue) {
							clearFields = true
						}
					}

					// CSB0122RC0140000
					// if isFc164 && utils.RegIsMatch(`^(fc144|fc339|fc340|fc341|fc342|fc343|fc344|fc345)$`, code) {
					// 	if fieldsArr[loc[2]][loc[3]].FinalValue != "" && fieldsArr[loc[2]][loc[3]].FinalValue != "0%" && fieldsArr[loc[2]][loc[3]].FinalValue != "0.00%" {
					// 		fieldsArr[loc[2]][loc[3]].FinalValue = "0.00%"
					// 	}
					// }

				}
			}
		}
	}

	// fc110Locs := fieldLocationMap["fc110"]
	// if len(fc110Locs) > 0 {
	// 	fc110Obj := obj.Invoice[fc110Locs[0][0]].Invoice[fc110Locs[0][2]][fc110Locs[0][3]]
	// 	if fc110Obj.ResultValue == "1" {
	// 		fc010Locs := fieldLocationMap["fc010"]
	// 		fc009Locs := fieldLocationMap["fc009"]
	// 		if len(fc010Locs) > 0 && len(fc009Locs) > 0 {
	// 			obj.Invoice[fc009Locs[0][0]].Invoice[fc009Locs[0][2]][fc009Locs[0][3]].FinalValue = obj.Invoice[fc010Locs[0][0]].Invoice[fc010Locs[0][2]][fc010Locs[0][3]].FinalValue
	// 		}
	// 	}
	// }

	// CSB0122RC0075000 CSB0122RC0139001 //调整新增字段 fc497、fc498、fc499、fc500、fc501、fc502、fc503、fc504、fc505、fc506、fc507、fc508、fc509、fc510、fc511、fc512、fc513、fc514、fc515、fc516
	if clearFields {
		cFields := []string{"fc164", "fc165", "fc166", "fc167", "fc168", "fc169", "fc170", "fc171", "fc172", "fc173", "fc174", "fc175", "fc176", "fc177", "fc178", "fc179", "fc180", "fc181", "fc182", "fc183", "fc497", "fc498", "fc499", "fc500", "fc501", "fc502", "fc503", "fc504", "fc505", "fc506", "fc507", "fc508", "fc509", "fc510", "fc511", "fc512", "fc513", "fc514", "fc515", "fc516"}
		for _, code := range cFields {
			utils.SetFinalValue(obj, fieldLocationMap[code], "")
		}
		utils.SetFinalValue(obj, fieldLocationMap["fc164"], "W55.951")
	}

	// fc153 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc153"][0], false)
	fc153s := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc153"], false)
	fc153_4 := false
	for _, fc153 := range fc153s {
		if fc153 == "4" {
			fc153_4 = true
		}
	}
	//调整 销售渠道为“纸质”时  编码 CSB0122RC0052000； CSB0122RC0054000 ；CSB0122RC0061000

	if obj.Bill.SaleChannel == "纸质" {
		if fc153_4 {
			if len(fieldLocationMap["fc015"]) > 0 {
				fc015 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc015"][0], false)
				if fc015 == "2" {
					fc022 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc022"][0], true)
					if fc022 == "" {
						utils.DelIssue(obj, fieldLocationMap["fc022"])
						utils.SetIssue(obj, fieldLocationMap["fc022"][0], "账户信息不完整（纸质申请书）", "C02", "")
						// utils.SetOnlyOneFinalValue(obj, loc1, "")
					}
					fc023 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc023"][0], true)
					if fc023 == "" {
						utils.DelIssue(obj, fieldLocationMap["fc023"])
						utils.SetIssue(obj, fieldLocationMap["fc023"][0], "账户信息不完整（纸质申请书）", "C02", "")
						// utils.SetOnlyOneFinalValue(obj, loc1, "")
					}
				}
			}
			if len(fieldLocationMap["fc011"]) > 0 {
				fc011 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc011"][0], false)
				if utils.RegIsMatch(`^(A|B)$`, fc011) || strings.Index(fc011, "?") != -1 {
					utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc011"][0], "")
					utils.SetIssue(obj, fieldLocationMap["fc011"][0], "申请人联系电话待确认", "C06", "")

				}
			}
		}
	}

	if len(fieldLocationMap["fc217"]) > 0 && len(fieldLocationMap["fc054"]) > 0 {
		fc217 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc217"][0], true)

		fc054 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc054"][0], true)
		if fc054 == "" {
			utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc054"][0], fc217)

		}

	}

	// CSB0122RC0066000
	cards := [][]string{{"fc026", "fc030", "fc029"}, {"fc224", "fc228", "fc227"}, {"fc234", "fc238", "fc237"}, {"fc244", "fc248", "fc247"}, {"fc254", "fc258", "fc257"}, {"fc018", "fc020", "fc019"}}
	for _, card := range cards {
		if len(fieldLocationMap[card[0]]) > 0 {
			f0 := utils.GetFieldValueByLoc(obj, fieldLocationMap[card[0]][0], true)
			if utils.CheckIDCard(f0) {
				sex := utils.Substr(f0, 16, 17)
				aaa, _ := strconv.Atoi(sex)
				if aaa%2 != 0 {
					utils.SetOnlyOneFinalValue(obj, fieldLocationMap[card[1]][0], "M")
				} else {
					utils.SetOnlyOneFinalValue(obj, fieldLocationMap[card[1]][0], "F")
				}
				birthDay := utils.Substr(f0, 6, 14)
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap[card[2]][0], birthDay)

			}
		}

	}
	// 编码 CSB0122RC0062000 ， CSB0122RC0063000 , {"fc039", "fc037", "zbxYiYuanMap", "9999999", "未匹配到中保信医院", "C21"}}
	codes := [][]string{{"fc038", "fc037", "yiYuanMap", "QT2004", "未匹配到内部定点医院", "C20"}, {"fc485", "fc040", "yiYuanMap", "QT2004", "", ""}}
	for _, code := range codes {
		fieldLocs := fieldLocationMap[code[0]]
		for ii, loc := range fieldLocs {
			invoiceMap := obj.Invoice[loc[0]]
			eleLen := reflect.ValueOf(invoiceMap).NumField()
			if eleLen > 0 {
				if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)
					fcxxx := utils.GetFieldValueByLoc(obj, fieldLocationMap[code[1]][ii], true)
					val, is := constMap[code[2]][fcxxx]
					if is {
						fieldsArr[loc[2]][loc[3]].FinalValue = val
					} else {
						// issue := model3.Issue{
						// 	Type:    "",
						// 	Code:    code[5],
						// 	Message: fcxxx + code[4],
						// }
						// fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
						fieldsArr[loc[2]][loc[3]].FinalValue = code[3]
					}

				}
			}
		}
	}

	// //CSB0122RC0070000
	// //"fc161、fc402、fc403、fc404、fc405、fc406、fc407、fc408、fc409、fc410、fc411、
	// //fc412、fc413、fc414、fc415、fc416、fc417、fc418、fc419、fc420、fc421、
	// //fc422、fc423、fc424、fc425、fc426、fc427、fc428、fc429、fc430、fc431
	// //以上字段录入值需根据机构号是否为 北京 来匹配对应的常量表，并转码输出
	// //匹配规则：
	// //北京匹配《B0122_信诚理赔_北京发票大项明细表》的“名称”（第二列）转“代码”（第一列）
	// //非北京匹配《B0122_信诚理赔_发票大项明细表》的“名称”（第二列）转“代码”（第一列）"
	codeArr := []string{
		"fc161", "fc402", "fc403", "fc404", "fc405", "fc406", "fc407", "fc408", "fc409", "fc410", "fc411",
		"fc412", "fc413", "fc414", "fc415", "fc416", "fc417", "fc418", "fc419", "fc420", "fc421",
		"fc422", "fc423", "fc424", "fc425", "fc426", "fc427", "fc428", "fc429", "fc430", "fc431"}
	// for _, code := range codeArr {
	// 	for _, loc := range fieldLocationMap[code] {
	// 		name := utils.GetFieldValueByLoc(obj, loc, false)
	// 		val, _ := constMap["beiJingFaPiaoDaXiangMap"][name]
	// 		if name != "北京" {
	// 			val, _ = constMap["faPiaoDaXiangMap"][name]
	// 		}
	// 		utils.SetOnlyOneFinalValue(obj, loc, val)
	// 	}
	// }

	// CSB0122RC0076000
	codes = [][]string{{"fc402", "fc432"}, {"fc403", "fc433"}, {"fc404", "fc434"}, {"fc405", "fc435"}, {"fc406", "fc436"}, {"fc407", "fc437"}, {"fc408", "fc438"}, {"fc409", "fc439"}, {"fc410", "fc440"}, {"fc411", "fc441"}, {"fc412", "fc442"}, {"fc413", "fc443"}, {"fc414", "fc444"}, {"fc415", "fc445"}, {"fc416", "fc446"}, {"fc417", "fc447"}, {"fc418", "fc448"}, {"fc419", "fc449"}, {"fc420", "fc450"}, {"fc421", "fc451"}, {"fc422", "fc452"}, {"fc423", "fc453"}, {"fc424", "fc454"}, {"fc425", "fc455"}, {"fc426", "fc456"}, {"fc427", "fc457"}, {"fc428", "fc458"}, {"fc429", "fc459"}, {"fc430", "fc460"}, {"fc431", "fc461"}}
	fc402Map := fieldLocationMap["fc402"]
	for ii, _ := range fc402Map {
		// invoiceMap := obj.Invoice[fc402[0]]
		codeSame := map[string]string{}
		for _, code := range codes {
			if len(fieldLocationMap[code[0]]) <= ii {
				continue
			}
			f0 := utils.GetFieldValueByLoc(obj, fieldLocationMap[code[0]][ii], true)
			// fmt.Println("-----------fc402Map-------------", code[0], f0, codeSame)
			if f0 == "" {
				continue
			}
			value, isExist := codeSame[f0]
			if isExist {
				f1 := utils.GetFieldValueByLoc(obj, fieldLocationMap[code[1]][ii], true)
				aa := utils.GetFieldValueByLoc(obj, fieldLocationMap[value][ii], true)
				num := utils.SumFloat(utils.ParseFloat(aa), utils.ParseFloat(f1), "+")
				// fmt.Println("-----------num-------------", num)
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap[value][ii], utils.ToString(num))
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap[code[0]][ii], "")
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap[code[1]][ii], "")
			} else {
				codeSame[f0] = code[1]
			}
		}
		// CSB0122RC0077000
		// fc083 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc083"][ii], true)
		// fc089 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc089"][ii], true)
		// fc093 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc093"][ii], true)
		// fc113 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc113"][ii], true)
		for _, code := range codes {
			if len(fieldLocationMap[code[0]]) <= ii {
				continue
			}
			f0 := utils.GetFieldValueByLoc(obj, fieldLocationMap[code[0]][ii], true)
			f1 := utils.GetFieldValueByLoc(obj, fieldLocationMap[code[1]][ii], true)
			if f0 != "" {
				fc134 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc134"][ii], true)
				fc134 = fc134 + f0 + "," + f1 + "元;"
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc134"][ii], fc134)
			}
			f0 = utils.GetFieldValueByLoc(obj, fieldLocationMap[code[0]][ii], false)
			if f1 == "" {
				continue
			}
			// CSB0122RC0078000
			if f0 == "床位费" {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc083"][ii], f1)
			}
			// CSB0122RC0083000
			if f0 == "护理费" {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc089"][ii], f1)
			}
			// CSB0122RC0087000
			if utils.RegIsMatch(`^(西药费|中草药费|中药费|中成药费|民族药费|自制制剂)$`, f0) {
				fc093 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc093"][ii], true)
				sum := utils.SumFloat(utils.ParseFloat(fc093), utils.ParseFloat(f1), "+")
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc093"][ii], utils.ToString(sum))
			}
			// CSB0122RC0106000
			if f0 == "治疗费" {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc113"][ii], f1)
			}
		}

	}

	// // CSB0122RC0140000
	fc164ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc164"], true)
	isFc164 := false
	for _, fc164 := range fc164ValArr {
		if utils.RegIsMatch(`^(W53\.951|W54\.951|W55\.951|W55\.952|W56\.851|W56\.852|W57\.951|W58\.951|W59\.951|W60\.951|W64\.951|W64\.952|W64\.953|W64\.954)$`, fc164) {
			isFc164 = true
			break
		}
	}

	//CSB0122RC0073000
	//"下列字段为循环分块中的字段，每一行为一组，
	//1.当第2列字段的录入值为1时，第5列字段结果值默认为0%，第6列字段默认为0，第7列字段默认为1；
	//2.当第2列字段的录入值为3时，第5列字段结果值默认为100%；第6列字段默认为第4列字段的值，第7列字段默认为3；
	//3.当第2列字段的录入值为2或5时，将第1列字段的录入内容匹配《B0122_信诚理赔_全国》中“项目名称”，将匹配出来的【自付比例】赋值给第5列，【类别】赋值为第7列，第6列字段默认为“第4列*第5列的值”；
	//4.当第2列字段的录入值为4时，第6列字段结果值默认为“第4列*第5列的值”，第7列字段结果值根据第5列字段的值进行输出，
	//a.当第5列的值不为0%且不为100%时，第7列输出2；
	//b.当第5列的值为0%，第7列输出1；
	//c.当第5列的值为100%，第7列输出3；
	//5.当第2列字段的录入值为6时，第5列字段结果值默认为“第6列/第4列*100%的值”，第7列字段结果值根据第5列字段的值进行输出，
	//a.当第5列的值不为0%且不为100%时，第7列输出2；
	//b.当第5列的值为0%，第7列输出1；
	//c.当第5列的值为100%，第7列输出3；
	//收费名称 清单医保类型 数量 总价 自付比例 清单自费(自付)金额 药品类别
	//①    ②    ③    ④    ⑤    ⑥    ⑦
	codes = [][]string{
		{"fc139", "fc145", "fc141", "fc142", "fc144", "fc462", "fc136"},
		{"fc303", "fc346", "fc317", "fc324", "fc339", "fc463", "fc282"},
		{"fc304", "fc347", "fc318", "fc325", "fc340", "fc464", "fc283"},
		{"fc305", "fc348", "fc319", "fc326", "fc341", "fc465", "fc284"},
		{"fc306", "fc349", "fc320", "fc327", "fc342", "fc466", "fc285"},
		{"fc307", "fc350", "fc321", "fc328", "fc343", "fc467", "fc286"},
		{"fc308", "fc351", "fc322", "fc329", "fc344", "fc468", "fc287"},
		{"fc309", "fc352", "fc323", "fc330", "fc345", "fc469", "fc288"},
	}
	if !isFc164 {
		for _, codeArr := range codes {
			for _, loc1 := range fieldLocationMap[codeArr[0]] {
				loc2 := utils.GetFieldLoc(fieldLocationMap[codeArr[1]], loc1[0], loc1[1], loc1[2], loc1[4])
				loc3 := utils.GetFieldLoc(fieldLocationMap[codeArr[2]], loc1[0], loc1[1], loc1[2], loc1[4])
				loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc1[0], loc1[1], loc1[2], loc1[4])
				loc5 := utils.GetFieldLoc(fieldLocationMap[codeArr[4]], loc1[0], loc1[1], loc1[2], loc1[4])
				loc6 := utils.GetFieldLoc(fieldLocationMap[codeArr[5]], loc1[0], loc1[1], loc1[2], loc1[4])
				loc7 := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], loc1[0], loc1[1], loc1[2], loc1[4])
				if len(loc2) != 1 || len(loc3) != 1 || len(loc4) != 1 || len(loc5) != 1 || len(loc6) != 1 || len(loc7) != 1 {
					continue
				}
				val2 := utils.GetFieldValueByLoc(obj, loc2[0], false)
				val4 := utils.GetFieldValueByLoc(obj, loc4[0], true)
				valDecimal4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
				val5 := utils.GetFieldValueByLoc(obj, loc5[0], true)
				//valDecimal5 := utils.GetFieldDecimalValueByLoc(obj, loc5[0], true)
				val6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
				valDecimal6 := utils.GetFieldDecimalValueByLoc(obj, loc6[0], true)
				val7 := utils.GetFieldValueByLoc(obj, loc7[0], true)
				setVal := []string{val5, val6, val7}

				//1.当第2列字段的录入值为1时，第5列字段结果值默认为0%，第6列字段默认为0，第7列字段默认为1； 20230822 调整为 第7列字段默认为甲
				if val2 == "1" {
					setVal = []string{"0%", "0", "甲"}
				}
				//2.当第2列字段的录入值为3时，第5列字段结果值默认为100%；第6列字段默认为第4列字段的值，第7列字段默认为3；  20230822 调整为 第7列字段默认为丙
				if val2 == "3" {
					setVal = []string{"100%", val4, "丙"}
				}
				//3.当第2列字段的录入值为2或5时，将第1列字段的录入内容匹配《B0122_信诚理赔_全国》中“项目名称”，将匹配出来的【自付比例】（第五列）赋值给第5列，【分类】（第三列）赋值为第7列，第6列字段默认为“第4列*第5列的值”；  20230822 类别调整为分类
				if m, _ := regexp.MatchString("^(2|5)$", val2); m {
					val1 := utils.GetFieldValueByLoc(obj, loc1, true)
					biLi, isExist := constMap["quanGuoMap"][val1]
					if isExist {
						fenLei, _ := constMap["quanGuo1Map"][val1]
						fmt.Println(fenLei)
						valDecimal5 := utils.ParseDecimal(strings.Replace(biLi, "%", "", -1)).Div(utils.ParseDecimal("100"))
						setVal = []string{biLi, valDecimal4.Mul(valDecimal5).StringFixed(2), fenLei}
					} else {
						setVal = []string{"0%", "0", "甲"}
					}

				}

				//4.当第2列字段的录入值为4时，第6列字段结果值默认为“第4列*第5列的值”，第7列字段结果值根据第5列字段的值进行输出，
				//a.当第5列的值不为0%且不为100%时，第7列输出2； 20230822 第7列输出2 调整为 第7列输出乙
				//b.当第5列的值为0%，第7列输出1；	20230822 第7列输出1 调整为 第7列输出甲
				//c.当第5列的值为100%，第7列输出3；	20230822 第7列输出3 调整为 第7列输出丙
				if val2 == "4" {
					valDecimal5 := utils.ParseDecimal(strings.Replace(val5, "%", "", -1)).Div(utils.ParseDecimal("100"))
					setVal = []string{val5, valDecimal4.Mul(valDecimal5).StringFixed(2), ""}
					if val5 != "0%" && val5 != "0.00%" && val5 != "100%" {
						setVal[2] = "乙"
					}
					if val5 == "0%" || val5 == "0.00%" {
						setVal[2] = "甲"
					}
					if val5 == "100%" {
						setVal[2] = "丙"
					}
				}

				//5.当第2列字段的录入值为6时，第5列字段结果值默认为“第6列/第4列*100%的值”，第7列字段结果值根据第5列字段的值进行输出，
				//a.当第5列的值不为0%且不为100%时，第7列输出2；	20230822 第7列输出2 调整为 第7列输出乙
				//b.当第5列的值为0%，第7列输出1；	20230822 第7列输出1 调整为 第7列输出甲
				//c.当第5列的值为100%，第7列输出3；	20230822 第7列输出3 调整为 第7列输出丙
				if val2 == "6" {
					setVal = []string{valDecimal6.Div(valDecimal4).Mul(utils.ParseDecimal("100")).StringFixed(2) + "%", val6, ""}
					if setVal[0] != "0%" && setVal[0] != "0.00%" && setVal[0] != "100%" {
						setVal[2] = "乙"
					}
					if setVal[0] == "0%" && setVal[0] == "0.00%" {
						setVal[2] = "甲"
					}
					if setVal[0] == "100%" {
						setVal[2] = "丙"
					}
				}

				utils.SetFinalValue(obj, loc5, setVal[0])
				utils.SetFinalValue(obj, loc6, setVal[1])
				utils.SetFinalValue(obj, loc7, setVal[2])
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	startTime := time.Now()
	//CSB0122RC0065000
	//"根据fc037的录入值匹配常量表《B0122_信诚理赔_内部医院代码表》及《B0122_信诚理赔_中保信医院代码表》得出对应的“城市”（第三列），
	//将该得到的“城市”与案件“机构”相对比，若一致或为包含关系，则fc132结果值默认1，否则fc133结果值默认2
	//如fc037值为华南理工大学校医院，则根据常量表得出城市为“广州”，案件机构为“广东(广州)”，则fc132默认1"
	for _, fc037Loc := range fieldLocationMap["fc037"] {
		fc037Val := utils.GetFieldValueByLoc(obj, fc037Loc, false)
		//city, ok := constMap["neiBuYiYuanMap"][fc037Val]
		//if !ok {
		city, _ := constMap["zhongBaoXinYiYuanMap"][fc037Val]
		//}
		if strings.Index(obj.Bill.Agency, city) != -1 {
			fc132Loc := utils.GetFieldLoc(fieldLocationMap["fc132"], fc037Loc[0], -1, -1, -1)
			utils.SetFinalValue(obj, fc132Loc, "1")
		} else {
			fc133Loc := utils.GetFieldLoc(fieldLocationMap["fc133"], fc037Loc[0], -1, -1, -1)
			utils.SetFinalValue(obj, fc133Loc, "2")
		}
	}

	//CSB0122RC0074000
	//下列字段为循环分块中的字段，每一行为一组，按顺序执行以下每条校验：
	//1.第5列字段的结果值为：第3列/第2列的值；
	//2.第6列字段的结果值为：100.00%-第4列的值；（如第4列为10.00%，则第6列为90.00%）
	//3.第7列字段的结果值根据第10列字段结果值进行判断：
	//-第10列为“甲”，则第7列结果值为“无自付”
	//-第10列为“乙”，则第7列结果值为“有自付”
	//-第10列为“丙”，则第7列结果值为“全自付”
	//4.第9列字段的结果值为：第3列*第4列的值；
	//5.第8列字段的结果值为：第3列-第9列的值；
	//名称  数量  总价 自付比 单价 报销 有无自付 医保内 医保外 药品类别
	//①    ②    ③    ④    ⑤    ⑥    ⑦    ⑧    ⑨    ⑩
	codes = [][]string{
		{"fc139", "fc141", "fc142", "fc144", "fc140", "fc149", "fc150", "fc151", "fc152", "fc136"},
		{"fc303", "fc317", "fc324", "fc339", "fc310", "fc374", "fc381", "fc388", "fc395", "fc282"},
		{"fc304", "fc318", "fc325", "fc340", "fc311", "fc375", "fc382", "fc389", "fc396", "fc283"},
		{"fc305", "fc319", "fc326", "fc341", "fc312", "fc376", "fc383", "fc390", "fc397", "fc284"},
		{"fc306", "fc320", "fc327", "fc342", "fc313", "fc377", "fc384", "fc391", "fc398", "fc285"},
		{"fc307", "fc321", "fc328", "fc343", "fc314", "fc378", "fc385", "fc392", "fc399", "fc286"},
		{"fc308", "fc322", "fc329", "fc344", "fc315", "fc379", "fc386", "fc393", "fc400", "fc287"},
		{"fc309", "fc323", "fc330", "fc345", "fc316", "fc380", "fc387", "fc394", "fc401", "fc288"},
	}
	if !isFc164 {
		for _, codeArr = range codes {
			for _, loc6 := range fieldLocationMap[codeArr[5]] {
				loc2 := utils.GetFieldLoc(fieldLocationMap[codeArr[1]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc3 := utils.GetFieldLoc(fieldLocationMap[codeArr[2]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc5 := utils.GetFieldLoc(fieldLocationMap[codeArr[4]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc7 := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc8 := utils.GetFieldLoc(fieldLocationMap[codeArr[7]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc9 := utils.GetFieldLoc(fieldLocationMap[codeArr[8]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc10 := utils.GetFieldLoc(fieldLocationMap[codeArr[9]], loc6[0], loc6[1], loc6[2], loc6[4])
				if len(loc2) != 1 || len(loc3) != 1 || len(loc4) != 1 || len(loc5) != 1 || len(loc7) != 1 || len(loc8) != 1 || len(loc9) != 1 || len(loc10) != 1 {
					continue
				}
				//1.第5列字段的结果值为：第3列/第2列的值；
				valDecimal3 := utils.GetFieldDecimalValueByLoc(obj, loc3[0], true)
				valDecimal2 := utils.GetFieldDecimalValueByLoc(obj, loc2[0], true)
				if !valDecimal2.IsZero() {
					utils.SetOnlyOneFinalValue(obj, loc5[0], valDecimal3.Div(valDecimal2).StringFixed(2))
				}

				//2.第6列字段的结果值为：100.00%-第4列的值；（如第4列为10.00%，则第6列为90.00%）
				val4 := strings.ReplaceAll(utils.GetFieldValueByLoc(obj, loc4[0], true), "%", "")
				valDecimal4 := utils.ParseDecimal(val4)
				utils.SetOnlyOneFinalValue(obj, loc6, decimal.NewFromInt(100).Sub(valDecimal4).StringFixed(2)+"%")

				//3.第7列字段的结果值根据第10列字段结果值进行判断：
				//-第10列为“甲”，则第7列结果值为“无自付”
				//-第10列为“乙”，则第7列结果值为“有自付”
				//-第10列为“丙”，则第7列结果值为“全自付”
				myMap := map[string]string{
					"甲": "无自付",
					"乙": "有自付",
					"丙": "全自付",
				}
				for k, v := range myMap {
					if utils.GetFieldValueByLoc(obj, loc10[0], true) == k {
						utils.SetFinalValue(obj, loc7, v)
					}
				}

				//4.第9列字段的结果值为：第3列*第4列的值；
				utils.SetFinalValue(obj, loc9, valDecimal3.Mul(valDecimal4).Div(utils.ParseDecimal("100")).StringFixed(2))

				//5.第8列字段的结果值为：第3列-第9列的值；
				valDecimal9 := utils.GetFieldDecimalValueByLoc(obj, loc9[0], true)
				utils.SetFinalValue(obj, loc8, valDecimal3.Sub(valDecimal9).StringFixed(2))
			}
		}
	}

	// CSB0122RC0079000
	codes = [][]string{
		{"fc139", "fc145", "fc141", "fc142", "fc143", "fc144"},
		{"fc303", "fc346", "fc317", "fc324", "fc331", "fc339"},
		{"fc304", "fc347", "fc318", "fc325", "fc332", "fc340"},
		{"fc305", "fc348", "fc319", "fc326", "fc333", "fc341"},
		{"fc306", "fc349", "fc320", "fc327", "fc334", "fc342"},
		{"fc307", "fc350", "fc321", "fc328", "fc335", "fc343"},
		{"fc308", "fc351", "fc322", "fc329", "fc336", "fc344"},
		{"fc309", "fc352", "fc323", "fc330", "fc337", "fc345"},
	}
	for ii, invoice := range obj.Invoice {
		for jj, fields := range invoice.QingDan {
			isSum := false
			isSum2 := false
			isSum3 := false
			isSum4 := false
			fc161 := ""
			for kk, field := range fields {
				if field.Code == "fc161" {
					fc161 = field.FinalValue
					if field.ResultValue == "床位费" {
						isSum = true
					} else {
						isSum = false
					}
					if field.ResultValue == "护理费" {
						isSum3 = true
					} else {
						isSum3 = false
					}

					if utils.RegIsMatch(`^(西药费|中草药费|中药费|中成药费|民族药费|自制制剂)$`, field.ResultValue) {
						isSum2 = true
					} else {
						isSum2 = false
					}
					if !utils.RegIsMatch(`^(床位费|护理费|西药费|中草药费|中药费|中成药费|民族药费|自制制剂)$`, field.ResultValue) {
						isSum4 = true
					} else {
						isSum4 = false
					}
				}
				// CSB0122RC0071000
				if utils.RegIsMatch(`^(fc138|fc296|fc297|fc298|fc299|fc300|fc301|fc302)$`, field.Code) {
					obj.Invoice[ii].QingDan[jj][kk].FinalValue = fc161
					// utils.SetOnlyOneFinalValue(obj, fieldLocationMap[field.Code][ii], fc161)
				}
				if isSum {
					for _, code := range codes {
						if field.Code == code[0] {
							if utils.RegIsMatch(`(重症监护.*病房|ICU.*病房|CCU.*病房（重症加强护理病房）|NICU.*病房（新生儿重症监护病房）|层流洁净.*病房)`, field.ResultValue) {
								f3 := GetFieldValue(fields, code[3], false, field.BlockIndex)
								fc084 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc084"][ii], true)
								sum := utils.SumFloat(utils.ParseFloat(f3), utils.ParseFloat(fc084), "+")
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc084"][ii], utils.ToFix2(utils.ToString(sum)))

							}
							// CSB0122RC0082000 20230809新增 陪*床|医疗废物处置费|空调*费|降温费|取暖费|*药袋
							if utils.RegIsMatch(`^(观察病房床位费|家庭病床|医疗废物处置费|降温费|取暖费)$`, field.ResultValue) || utils.RegIsMatch(`(陪.*床|空调.*费|药袋)`, field.ResultValue) {
								f3 := GetFieldValue(fields, code[3], false, field.BlockIndex)
								fc087 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc087"][ii], true)
								sum := utils.SumFloat(utils.ParseFloat(f3), utils.ParseFloat(fc087), "+")
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc087"][ii], utils.ToFix2(utils.ToString(sum)))

								fc088 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc088"][ii], true)
								if fc088 != "" {
									fc088 = fc088 + "、"
								}
								fc088 += field.FinalValue
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc088"][ii], fc088)

							}
						}
						// CSB0122RC0080000
						if field.Code == code[5] && utils.RegIsMatch(`^(1\.00|100%)$`, field.FinalValue) {
							f3 := GetFieldValue(fields, code[3], false, field.BlockIndex)
							fc085 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc085"][ii], true)
							sum := utils.SumFloat(utils.ParseFloat(f3), utils.ParseFloat(fc085), "+")
							utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc085"][ii], utils.ToFix2(utils.ToString(sum)))

							// CSB0122RC0081000
							f0 := GetFieldValue(fields, code[0], false, field.BlockIndex)
							fc086 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc086"][ii], true)
							if fc086 != "" {
								fc086 = fc086 + "、"
							}
							fc086 += f0
							utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc086"][ii], fc086)

						}
					}
				} else if isSum2 {
					for _, code := range codes {
						// CSB0122RC0088000
						if field.Code == code[5] && utils.RegIsMatch(`^(1\.00|100%)$`, field.FinalValue) {
							f3 := GetFieldValue(fields, code[3], false, field.BlockIndex)
							if len(fieldLocationMap["fc094"]) > ii {
								fc094 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc094"][ii], true)
								sum := utils.SumFloat(utils.ParseFloat(f3), utils.ParseFloat(fc094), "+")
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc094"][ii], utils.ToFix2(utils.ToString(sum)))
							}

							// CSB0122RC0089000
							if len(fieldLocationMap["fc095"]) > ii {
								f0 := GetFieldValue(fields, code[0], false, field.BlockIndex)
								fc095 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc095"][ii], true)
								if fc095 != "" {
									fc095 = fc095 + "、"
								}
								fc095 += f0
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc095"][ii], fc095)

							}

						}
					}
				} else if isSum3 {
					for _, code := range codes {
						if field.Code == code[0] {
							// CSB0122RC0084000
							if utils.RegIsMatch(`(重症监护.*护|特级.*护理)`, field.ResultValue) {
								f3 := GetFieldValue(fields, code[3], false, field.BlockIndex)
								fc090 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc090"][ii], true)
								sum := utils.SumFloat(utils.ParseFloat(f3), utils.ParseFloat(fc090), "+")
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc090"][ii], utils.ToFix2(utils.ToString(sum)))

							}
						}
						// CSB0122RC0080000
						if field.Code == code[5] && utils.RegIsMatch(`^(1\.00|100%)$`, field.FinalValue) {
							f3 := GetFieldValue(fields, code[3], false, field.BlockIndex)

							//CSB0122RC0085000
							fc091 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc091"][ii], true)
							sum := utils.SumFloat(utils.ParseFloat(f3), utils.ParseFloat(fc091), "+")
							utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc091"][ii], utils.ToFix2(utils.ToString(sum)))

							// CSB0122RC0086000
							f0 := GetFieldValue(fields, code[0], true, field.BlockIndex)
							fc092 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc092"][ii], true)
							if fc092 != "" {
								fc092 = fc092 + "、"
							}
							fc092 += f0
							utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc092"][ii], fc092)

						}
					}
				}
				if isSum4 {
					for _, code := range codes {
						// 编码 CSB0122RC0160000 时间 20230809
						//同一发票下，以fc161（初审）为一组进行判断，当fc161录入不为 【床位费】、【护理费】、【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】时
						//下列字段以一行为一组，当第1列录入内容为【陪*床】、【医疗废物处置费】、【空调*费】、【降温费】、【取暖费】、【*药袋】，将第4列字段的结果值进行汇总，赋值给fc124，数字格式，小数点后保留两位。（如为负数的，无需处理）
						//收费名称 清单医保类型 数量 总价 医保比例 自付比例
						if code[0] == field.Code {
							// matched, _ := regexp.MatchString(`^(陪\*床|医疗废物处置费|空调\*费|降温费|取暖费|\*药袋)$`, field.ResultValue)
							if utils.RegIsMatch(`(陪.*床|空调.*费|药袋)`, field.ResultValue) || utils.RegIsMatch(`^(医疗废物处置费|降温费|取暖费)$`, field.ResultValue) {
								f4 := GetFieldValue(fields, code[3], true, field.BlockIndex)
								fc124 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc124"][ii], true)
								sum := utils.SumFloat(utils.ParseFloat(f4), utils.ParseFloat(fc124), "+")
								if sum > 0 {
									utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc124"][ii], utils.ToFix2(utils.ToString(sum)))
								}
								//--------------------------fc125--------------------------------
								//编码 CSB0122RC0161000
								f1 := GetFieldValue(fields, code[0], true, field.BlockIndex)
								if len(fieldLocationMap["fc125"]) > 0 {
									fc125 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc125"][ii], true)
									// fc125Sum := utils.SumFloat(utils.ParseFloat(f1), utils.ParseFloat(fc125), "+")
									if fc125 != "" {
										fc125 = fc125 + "、"
									}
									fc125 += f1
									// if fc125Sum > 0 {
									utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc125"][ii], fc125)
									// }
								}
							}
						}
					}
				}
			}
		}
	}

	//CSB0122RC0090000
	//"同一发票属性下，判断是否同时存在MB002中的bc020分块、bc021分块，
	//1、如不存在、
	//每一行为一组，当第7列值为中草药费时，将第6列字段进行金额合计，赋值给fc096，数字格式，保留小数点后两位。
	//2、如存在、且fc280、fc281的结果值不为空时；
	//每一行为一组，当第7列值为中草药费时，且6列结果值为1.00或100%时，将第4列字段进行金额合计，赋值给fc096，数字格式，保留小数点后两位。
	//（如存在这2个分块，但fc280、fc281的结果值为空时，参照第1点规则进行处理）

	//CSB0122RC0091000
	//"同一发票属性下，判断是否同时存在MB002中的bc020分块、bc021分块，
	//1、如不存在、
	//每一行为一组，当第7列值为中草药费时，第1列字段的结果值进行汇总，赋值给fc097，例如：陪床费、陪人费、空调费（如为负数的，无需处理）
	//2、如存在、且fc280、fc281的结果值为不空时；
	//每一行为一组，当第7列值为中草药费时，且当6列结果值为1.00或100%时，将第1列字段的结果值进行汇总，赋值给fc097，例如：陪床费、陪人费、空调费（如为负数的，无需处理）
	//（如存在这2个分块，但fc280、fc281的结果值为空时，参照第1点规则进行处理
	codes = [][]string{
		{"fc139", "fc145", "fc141", "fc142", "fc143", "fc144", "fc138"},
		{"fc303", "fc346", "fc317", "fc324", "fc331", "fc339", "fc296"},
		{"fc304", "fc347", "fc318", "fc325", "fc332", "fc340", "fc297"},
		{"fc305", "fc348", "fc319", "fc326", "fc333", "fc341", "fc298"},
		{"fc306", "fc349", "fc320", "fc327", "fc334", "fc342", "fc299"},
		{"fc307", "fc350", "fc321", "fc328", "fc335", "fc343", "fc300"},
		{"fc308", "fc351", "fc322", "fc329", "fc336", "fc344", "fc301"},
		{"fc309", "fc352", "fc323", "fc330", "fc337", "fc345", "fc302"},
	}
	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票
		invoiceMap := obj.Invoice[i]
		fc096Loc := utils.GetFieldLoc(fieldLocationMap["fc096"], i, -1, -1, -1)
		fc097Loc := utils.GetFieldLoc(fieldLocationMap["fc097"], i, -1, -1, -1)
		fc280Loc := utils.GetFieldLoc(fieldLocationMap["fc280"], i, -1, -1, -1)
		fc281Loc := utils.GetFieldLoc(fieldLocationMap["fc281"], i, -1, -1, -1)

		fc280Val := ""
		fc281Val := ""
		if len(fc280Loc) == 1 && len(fc281Loc) == 1 {
			fc280Val = utils.GetFieldValueByLoc(obj, fc280Loc[0], true)
			fc281Val = utils.GetFieldValueByLoc(obj, fc281Loc[0], true)
		}

		//totalVal := decimal.Zero
		totalVal1 := make([]string, 0)
		if (len(invoiceMap.BaoXiaoDan) == 0 && len(invoiceMap.ThirdBaoXiaoDan1) == 0) ||
			(fc280Val == "" && fc281Val == "") {
			for _, codeArr = range codes {
				loc7Arr := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], i, -1, -1, -1)
				for _, loc7 := range loc7Arr {
					val7 := utils.GetFieldValueByLoc(obj, loc7, true)
					//loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc7[0], loc7[1], loc7[2], loc7[4])
					loc1 := utils.GetFieldLoc(fieldLocationMap[codeArr[0]], loc7[0], loc7[1], loc7[2], loc7[4])
					if len(loc1) != 1 {
						continue
					}
					//val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
					val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
					if val7 == "中草药费" {
						//totalVal = totalVal.Add(val4)
						if val1 == "" {
							continue
						}
						totalVal1 = append(totalVal1, val1)
					}
				}
			}
		} else if fc280Val != "" && fc281Val != "" {
			for _, codeArr = range codes {
				loc7Arr := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], i, -1, -1, -1)
				for _, loc7 := range loc7Arr {
					val7 := utils.GetFieldValueByLoc(obj, loc7, true)
					loc6 := utils.GetFieldLoc(fieldLocationMap[codeArr[5]], loc7[0], loc7[1], loc7[2], loc7[4])
					//loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc7[0], loc7[1], loc7[2], loc7[4])
					loc1 := utils.GetFieldLoc(fieldLocationMap[codeArr[0]], loc7[0], loc7[1], loc7[2], loc7[4])
					if len(loc6) != 1 && len(loc1) != 1 {
						continue
					}
					val6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
					//val4 := utils.GetFieldDecimalValueByLoc(obj, loc6[0], true)
					val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
					if val7 == "中草药费" && (val6 == "1.00" || val6 == "100%") {
						//totalVal = totalVal.Add(val4)
						if val1 == "" {
							continue
						}
						totalVal1 = append(totalVal1, val1)
					}
				}
			}
		}
		//调整修改  CSB0122RC0090000 同一发票属性下 第七列结果值为中草药费时，第四列的结果值进行合计，赋值给fc096 数字格式，保留两位小数
		count := decimal.Zero
		for _, codeArr := range codes {
			loc7Arr := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], i, -1, -1, -1)
			for _, loc7 := range loc7Arr {
				val7 := utils.GetFieldValueByLoc(obj, loc7, true)
				loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc7[0], loc7[1], loc7[2], loc7[4])
				if len(loc4) != 1 {
					continue
				}
				val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
				//fmt.Println("-----------val7--------------", val7)
				if val7 == "中草药费" {
					count = count.Add(val4)
				}
			}
		}
		fmt.Println(count)
		utils.SetFinalValue(obj, fc096Loc, count.StringFixed(2))
		utils.SetFinalValue(obj, fc097Loc, strings.Join(totalVal1, "、"))
	}

	//CSB0122RC0092000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值不为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_靶向药清单》的			【通用名】（第二列）、【商品名】（第三列），	 匹配上时，将第4列字段进行金额合计，	赋值给fc099，数字格式，保留小数点后两位。（负数，无需进行处理）
	//CSB0122RC0093000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值不为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_靶向药清单》的			【通用名】（第二列）、【商品名】（第三列），	 匹配上时，将第1列字段的结果值进行汇总，	赋值给fc100，例如：陪床费、陪人费、空调费（如为负数的，无需处理）
	//CSB0122RC0094000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值不为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_化疗药清单》的			【名称】（第二列）、【别名】（第三列），	 	 匹配上时，将第4列字段进行金额合计，	赋值给fc101，数字格式，保留小数点后两位。（负数，无需进行处理
	//CSB0122RC0095000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值不为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_化疗药清单》的			【名称】（第二列）、【别名】（第三列），	 	 匹配上时，将第1列字段的结果值进行汇总，	赋值给fc102，例如：陪床费、陪人费、空调费（如为负数的，无需处理
	//CSB0122RC0096000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值不为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_器官移植抗排斥药》的	【名称】（第二列），				 		 匹配上时，将第4列字段进行金额合计，	赋值给fc103，数字格式，保留小数点后两位。（负数，无需进行处理
	//CSB0122RC0097000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值不为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_器官移植抗排斥药》的	【名称】（第二列），				 		 匹配上时，将第1列字段的结果值进行汇总，	赋值给fc104，例如：陪床费、陪人费、空调费（如为负数的，无需处理
	//CSB0122RC0100000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_靶向药清单》的			【通用名】（第二列）、【商品名】（第三列），	 匹配上时，将第4列字段进行金额合计，	赋值给fc107，数字格式，保留小数点后两位。（负数，无需进行处理
	//CSB0122RC0101000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_靶向药清单》的			【通用名】（第二列）、【商品名】（第三列），	 匹配上时，将第1列字段的结果值进行汇总，	赋值给fc108，例如：陪床费、陪人费、空调费（如为负数的，无需处理
	//CSB0122RC0102000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值为1.00或100%时,第1列录入内容匹配《BB0122_信诚理赔_化疗药清单》的		【名称】（第二列）、【别名】（第三列），	 	 匹配上时，将第4列字段进行金额合计，	赋值给fc109，数字格式，保留小数点后两位。（负数，无需进行处理
	//CSB0122RC0103000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值为1.00或100%时,第1列录入内容匹配《BB0122_信诚理赔_化疗药清单》的		【名称】（第二列）、【别名】（第三列），	 	 匹配上时，将第1列字段的结果值进行汇总，	赋值给fc110，例如：陪床费、陪人费、空调费（如为负数的，无需处理
	//CSB0122RC0104000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_器官移植抗排斥药》的	【名称】（第二列）、				 		 匹配上时，将第4列字段进行金额合计，	赋值给fc111，数字格式，保留小数点后两位。（负数，无需进行处理
	//CSB0122RC0105000	结果数据	"每一行为一组，第7列录入为【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】，6列结果值为1.00或100%时,第1列录入内容匹配《B0122_信诚理赔_器官移植抗排斥药的》的	【名称】（第二列）					 		 匹配上时，将第1列字段的结果值进行汇总，	赋值给fc112，例如：陪床费、陪人费、空调费（如为负数的，无需处理
	myArr := [][]string{
		{"baXiangYao", "fc099", "fc100", "fc107", "fc108"},
		{"huaLiaoYao", "fc101", "fc102", "fc109", "fc110"},
		{"kangPaiCiYao", "fc103", "fc104", "fc111", "fc112"},
	}
	for _, codeFaPiao := range myArr {
		for _, fc099Loc := range fieldLocationMap[codeFaPiao[1]] {
			fc100Loc := utils.GetFieldLoc(fieldLocationMap[codeFaPiao[2]], fc099Loc[0], -1, -1, -1)
			fc107Loc := utils.GetFieldLoc(fieldLocationMap[codeFaPiao[3]], fc099Loc[0], -1, -1, -1)
			fc108Loc := utils.GetFieldLoc(fieldLocationMap[codeFaPiao[4]], fc099Loc[0], -1, -1, -1)
			fc099Val := decimal.Zero
			fc100Val := make([]string, 0)
			fc107Val := decimal.Zero
			fc108Val := make([]string, 0)
			for _, codeArr = range codes {
				loc7Arr := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], fc099Loc[0], -1, -1, -1)
				for _, loc7 := range loc7Arr {
					val7 := utils.GetFieldValueByLoc(obj, loc7, true)
					if strings.Index("西药费,中草药费,中药费,中成药费,民族药费,自制制剂", val7) != -1 {
						loc1 := utils.GetFieldLoc(fieldLocationMap[codeArr[0]], loc7[0], loc7[1], loc7[2], loc7[4])
						loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc7[0], loc7[1], loc7[2], loc7[4])
						loc6 := utils.GetFieldLoc(fieldLocationMap[codeArr[5]], loc7[0], loc7[1], loc7[2], loc7[4])
						if len(loc1) != 1 && len(loc4) != 1 && len(loc6) != 1 {
							continue
						}
						val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
						val6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
						decimalVal4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
						_, ok := constMap[codeFaPiao[0]][val1]
						m, _ := regexp.MatchString("^(1\\.00|100%)$", val6)
						if ok {
							if !m {
								fc099Val = fc099Val.Add(decimalVal4)
								if val1 == "" {
									continue
								}
								fc100Val = append(fc100Val, val1)
							} else {
								fc107Val = fc107Val.Add(decimalVal4)
								if val1 == "" {
									continue
								}
								fc108Val = append(fc108Val, val1)
							}
						}
					}
				}
			}
			utils.SetOnlyOneFinalValue(obj, fc099Loc, fc099Val.StringFixed(2))
			utils.SetFinalValue(obj, fc100Loc, strings.Join(fc100Val, "、"))
			utils.SetFinalValue(obj, fc107Loc, fc107Val.StringFixed(2))
			utils.SetFinalValue(obj, fc108Loc, strings.Join(fc108Val, "、"))
		}
	}

	//CSB0122RC0107000	结果数据	"每一行为一组，第7列录入为治疗费时，下列字段以一行为一组，第6列结果值为100%时，将第4列字段的结果值进行汇总，赋值给fc114，数字格式，保留小数点后两位。（负数，无需进行处理）
	//CSB0122RC0108000	结果数据	"每一行为一组，第7列录入为治疗费时，下列字段以一行为一组，第6列结果值为100%时，将第1列字段的结果值进行汇总，赋值给fc115，例如：陪床费、陪人费、空调费（如为负数的，无需列明）
	for _, loc := range fieldLocationMap["fc114"] {
		fc115Loc := utils.GetFieldLoc(fieldLocationMap["fc115"], loc[0], -1, -1, -1)
		fc114TotalVal := decimal.Zero
		fc115TotalVal := make([]string, 0)
		for _, code := range codes {
			locArr7 := utils.GetFieldLoc(fieldLocationMap[code[6]], loc[0], -1, -1, -1)
			for _, loc7 := range locArr7 {
				val7 := utils.GetFieldValueByLoc(obj, loc7, true)
				loc6 := utils.GetFieldLoc(fieldLocationMap[code[5]], loc7[0], loc7[1], loc7[2], loc7[4])
				loc4 := utils.GetFieldLoc(fieldLocationMap[code[3]], loc7[0], loc7[1], loc7[2], loc7[4])
				loc1 := utils.GetFieldLoc(fieldLocationMap[code[0]], loc7[0], loc7[1], loc7[2], loc7[4])
				if len(loc6) != 1 || len(loc4) != 1 || len(loc1) != 1 {
					continue
				}
				val6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
				val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
				val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
				if val7 == "治疗费" && val6 == "100%" {
					fc114TotalVal = fc114TotalVal.Add(val4)
					fc115TotalVal = append(fc115TotalVal, val1)
				}
			}
		}
		utils.SetOnlyOneFinalValue(obj, loc, fc114TotalVal.StringFixed(2))
		utils.SetFinalValue(obj, fc115Loc, strings.Join(fc115TotalVal, "、"))
	}

	//CSB0122RC0109000	结果数据	"同一发票属性下，判断是否同时存在MB002中的bc020分块、bc021分块，
	//1、如不存在、
	//以fc161（初审）为一组进行判断，当fc161录入为治疗费时，第1列录入内容为【电*疗】、【光*疗】、【磁*疗】、【热*疗】、【声*疗】、【水*疗】、【针灸*】、【推拿*】、【拔罐*】、【刮痧*】、【顺势*疗】、【职业*疗】、【语音*疗】、【音乐*疗】时，将第4列字段的结果值进行汇总，赋值给fc116，数字格式，保留小数点后两位。（负数，无需进行处理）
	//2、如存在、且fc280、fc280的结果值为不空时；
	//以fc161（初审）为一组进行判断，当fc161录入为治疗费时，第1列录入内容为【电*疗】、【光*疗】、【磁*疗】、【热*疗】、【声*疗】、【水*疗】、【针灸*】、【推拿*】、【拔罐*】、【刮痧*】、【顺势*疗】、【职业*疗】、【语音*疗】、【音乐*疗】时，且第6列结果值为100%时，将第4列字段的结果值进行汇总，赋值给fc116，数字格式，保留小数点后两位。（负数，无需进行处理）
	//（如存在这2个分块，但fc280、fc280的结果值为空时，参照第1点规则进行处理）

	//CSB0122RC0110000	结果数据	"同一发票属性下，判断是否同时存在MB002中的bc020分块、bc021分块，
	//1、如不存在、
	//以fc161（初审）为一组进行判断，当fc161录入为治疗费时，第1列录入内容为【电*疗】、【光*疗】、【磁*疗】、【热*疗】、【声*疗】、【水*疗】、【针灸*】、【推拿*】、【拔罐*】、【刮痧*】、【顺势*疗】、【职业*疗】、【语音*疗】、【音乐*疗】时，将第1列字段的结果值进行汇总，赋值给fc117，例如：陪床费、陪人费、空调费（如为负数的，无需处理）
	//2、如存在、且fc280、fc280的结果值为不空时；
	//以ffc161（初审）为一组进行判断，当fc161录入为治疗费时，第1列录入内容为【电*疗】、【光*疗】、【磁*疗】、【热*疗】、【声*疗】、【水*疗】、【针灸*】、【推拿*】、【拔罐*】、【刮痧*】、【顺势*疗】、【职业*疗】、【语音*疗】、【音乐*疗】时，且第6列结果值为100%时，将第1列字段的结果值进行汇总，赋值给fc0117，例如：陪床费、陪人费、空调费（如为负数的，无需处理）
	//（如存在这2个分块，但fc280、fc280的结果值为空时，参照第1点规则进行处理）
	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票
		invoiceMap := obj.Invoice[i]
		fc116Loc := utils.GetFieldLoc(fieldLocationMap["fc116"], i, -1, -1, -1)
		fc117Loc := utils.GetFieldLoc(fieldLocationMap["fc117"], i, -1, -1, -1)

		fc280Loc := utils.GetFieldLoc(fieldLocationMap["fc280"], i, -1, -1, -1)
		fc281Loc := utils.GetFieldLoc(fieldLocationMap["fc281"], i, -1, -1, -1)
		fc280Val := ""
		fc281Val := ""
		if len(fc280Loc) == 1 && len(fc281Loc) == 1 {
			fc280Val = utils.GetFieldValueByLoc(obj, fc280Loc[0], true)
			fc281Val = utils.GetFieldValueByLoc(obj, fc281Loc[0], true)
		}

		fc116TotalVal := decimal.Zero
		fc117TotalVal := make([]string, 0)
		if (len(invoiceMap.BaoXiaoDan) == 0 && len(invoiceMap.ThirdBaoXiaoDan1) == 0) ||
			(fc280Val == "" && fc281Val == "") {
			for _, codeArr = range codes {
				loc7Arr := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], i, -1, -1, -1)
				for _, loc7 := range loc7Arr {
					val7 := utils.GetFieldValueByLoc(obj, loc7, true)
					loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc7[0], loc7[1], loc7[2], loc7[4])
					loc1 := utils.GetFieldLoc(fieldLocationMap[codeArr[0]], loc7[0], loc7[1], loc7[2], loc7[4])
					if len(loc4) != 1 && len(loc1) != 1 {
						continue
					}
					val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
					val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
					//20230809 新增【陪*床】、【医疗废物处置费】、【空调*费】、【降温费】、【取暖费】、【*药袋】
					// m, _ := regexp.MatchString("^(电\\*疗|光\\*疗|磁\\*疗|热\\*疗|声\\*疗|水\\*疗|针灸\\*|推拿\\*|拔罐\\*|刮痧\\*|顺势\\*疗|职业\\*疗|语音\\*疗|音乐\\*疗|陪\\*床|医疗废物处置费|空调\\*费|降温费|取暖费|\\*药袋)$", val1)
					if val7 == "治疗费" && (utils.RegIsMatch(`(电.*疗|光.*疗|磁.*疗|热.*疗|声.*疗|水.*疗|针灸.*|推拿.*|拔罐.*|刮痧.*|顺势.*疗|职业.*疗|语音.*疗|音乐.*疗|陪.*床|医疗废物处置费|空调.*费|降温费|取暖费|.*药袋)`, val1) || utils.RegIsMatch(`^(医疗废物处置费|降温费|取暖费)`, val1)) {
						fc116TotalVal = fc116TotalVal.Add(val4)
						fc117TotalVal = append(fc117TotalVal, val1)
					}
				}
			}
		} else if fc280Val != "" && fc281Val != "" {
			for _, codeArr = range codes {
				loc7Arr := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], i, -1, -1, -1)
				for _, loc7 := range loc7Arr {
					val7 := utils.GetFieldValueByLoc(obj, loc7, true)
					loc6 := utils.GetFieldLoc(fieldLocationMap[codeArr[5]], loc7[0], loc7[1], loc7[2], loc7[4])
					loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc7[0], loc7[1], loc7[2], loc7[4])
					loc1 := utils.GetFieldLoc(fieldLocationMap[codeArr[0]], loc7[0], loc7[1], loc7[2], loc7[4])
					if len(loc6) != 1 && len(loc4) != 1 && len(loc1) != 1 {
						continue
					}
					val6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
					val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
					val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
					// m, _ := regexp.MatchString("^(电\\*疗|光\\*疗|磁\\*疗|热\\*疗|声\\*疗|水\\*疗|针灸\\*|推拿\\*|拔罐\\*|刮痧\\*|顺势\\*疗|职业\\*疗|语音\\*疗|音乐\\*疗|陪\\*床|医疗废物处置费|空调\\*费|降温费|取暖费|\\*药袋)$", val1)
					if val7 == "治疗费" && (utils.RegIsMatch(`(电.*疗|光.*疗|磁.*疗|热.*疗|声.*疗|水.*疗|针灸.*|推拿.*|拔罐.*|刮痧.*|顺势.*疗|职业.*疗|语音.*疗|音乐.*疗|陪.*床|医疗废物处置费|空调.*费|降温费|取暖费|.*药袋)`, val1) || utils.RegIsMatch(`^(医疗废物处置费|降温费|取暖费)`, val1)) && val6 == "100%" {
						fc116TotalVal = fc116TotalVal.Add(val4)
						fc117TotalVal = append(fc117TotalVal, val1)
					}
				}
			}
		}
		utils.SetFinalValue(obj, fc116Loc, fc116TotalVal.StringFixed(2))
		utils.SetFinalValue(obj, fc117Loc, strings.Join(fc117TotalVal, "、"))
	}

	//CSB0122RC0111000	结果数据	"同一发票属性下，判断第一列字段字段录入值是否含有“*透析”，有该关键字时，将第4列字段的结果值进行汇总，赋值给fc119，数字格式，保留小数点后两位。（负数，无需进行处理） //调整 *透析 为 透析
	for _, loc := range fieldLocationMap["fc119"] {
		fc119TotalVal := decimal.Zero
		for _, code := range codes {
			locArr1 := utils.GetFieldLoc(fieldLocationMap[code[0]], loc[0], -1, -1, -1)
			for _, loc1 := range locArr1 {
				val1 := utils.GetFieldValueByLoc(obj, loc1, true)
				loc4 := utils.GetFieldLoc(fieldLocationMap[code[3]], loc1[0], loc1[1], loc1[2], loc1[4])
				if len(loc4) != 1 {
					continue
				}
				val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
				if strings.Index(val1, "透析") != -1 {
					fc119TotalVal = fc119TotalVal.Add(val4)
				}
			}
		}
		utils.SetOnlyOneFinalValue(obj, loc, fc119TotalVal.StringFixed(2))
	}

	//CSB0122RC0112000	结果数据	fc121的结果值=fc080-fc083-fc089-fc093
	codeArr = []string{"fc083", "fc089", "fc093"}
	for _, loc := range fieldLocationMap["fc121"] {
		totalVal := decimal.Zero
		fc080Loc := utils.GetFieldLoc(fieldLocationMap["fc080"], loc[0], -1, -1, -1)
		if len(fc080Loc) != 1 {
			continue
		}
		fc080Val := utils.GetFieldDecimalValueByLoc(obj, fc080Loc[0], true)
		for _, code := range codeArr {
			fc083Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
			if len(fc083Loc) != 1 {
				continue
			}
			fc083Val := utils.GetFieldDecimalValueByLoc(obj, fc083Loc[0], true)
			totalVal = totalVal.Add(fc083Val)
		}
		utils.SetOnlyOneFinalValue(obj, loc, fc080Val.Sub(totalVal).StringFixed(2))
	}

	//CSB0122RC0113000	结果数据	fc122的结果值=fc085+fc091+fc094
	//codeArr = []string{"fc085", "fc091", "fc094"}
	//for _, loc := range fieldLocationMap["fc122"] {
	//	totalVal := decimal.Zero
	//	for _, code := range codeArr {
	//		fc085Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
	//		if len(fc085Loc) != 1 {
	//			continue
	//		}
	//		fc085Val := utils.GetFieldDecimalValueByLoc(obj, fc085Loc[0], true)
	//		totalVal = totalVal.Add(fc085Val)
	//	}
	//	utils.SetOnlyOneFinalValue(obj, loc, totalVal.StringFixed(2))
	//}

	//CSB0122RC0114000	结果数据	"同一发票下，以fc161（初审）为一组进行判断，将fc161录入不为【床位费】、【护理费】、【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】时，下列字段以一行为一组，第6列结果值为100%时，将第4列字段的结果值进行汇总，赋值给fc122，数字格式，保留小数点后两位。（负数，无需进行处理）
	//CSB0122RC0115000	结果数据	"同一发票下，以fc161（初审）为一组进行判断，将fc161录入不为【床位费】、【护理费】、【西药费】、【中草药费】、【中药费】、【中成药费】、【民族药费】、【自制制剂】时，下列字段以一行为一组，第6列结果值为100%时，将第1列字段的结果值进行汇总，赋值给fc123，例如：陪床费、陪人费、空调费（如为负数的，无需处理）
	for _, loc := range fieldLocationMap["fc122"] {
		fc123Loc := utils.GetFieldLoc(fieldLocationMap["fc123"], loc[0], -1, -1, -1)
		fc122TotalVal := decimal.Zero
		fc123TotalVal := make([]string, 0)
		for _, code := range codes {
			locArr7 := utils.GetFieldLoc(fieldLocationMap[code[6]], loc[0], -1, -1, -1)
			for _, loc7 := range locArr7 {
				val7 := utils.GetFieldValueByLoc(obj, loc7, true)
				loc6 := utils.GetFieldLoc(fieldLocationMap[code[5]], loc7[0], loc7[1], loc7[2], loc7[4])
				loc4 := utils.GetFieldLoc(fieldLocationMap[code[3]], loc7[0], loc7[1], loc7[2], loc7[4])
				loc1 := utils.GetFieldLoc(fieldLocationMap[code[0]], loc7[0], loc7[1], loc7[2], loc7[4])
				if len(loc6) != 1 || len(loc4) != 1 || len(loc1) != 1 {
					continue
				}
				val6 := utils.GetFieldValueByLoc(obj, loc6[0], true)
				val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
				val4 := utils.GetFieldDecimalValueByLoc(obj, loc4[0], true)
				m, _ := regexp.MatchString("^(床位费|护理费|西药费|中草药费|中药费|中成药费|民族药费|自制制剂)$", val7)
				if !m && val6 == "100%" {
					fc122TotalVal = fc122TotalVal.Add(val4)
					fc123TotalVal = append(fc123TotalVal, val1)
				}
			}
		}
		utils.SetOnlyOneFinalValue(obj, loc, fc122TotalVal.StringFixed(2))
		utils.SetFinalValue(obj, fc123Loc, strings.Join(fc123TotalVal, "、"))
	}

	//CSB0122RC0135000
	//当fc153（模板类型字段）录入值为10时，校验每个fc040（MB001-bc010诊断）的录入值，
	// 当存在录入值重复时，按顺序将重复的第一个fc040对应的MB002-bc022、bc023、bc024、bc025、bc026
	// 分块下的字段一一对应赋值给其他重复的fc040下的MB002-bc022、bc023、bc024、bc025、bc026
	// 分块的对应字段
	//字段如下：
	//repMap := make(map[string]map[string]string, 0)
	//codeArr = []string{"fc045", "fc046", "fc047", "fc048", "fc044", "fc049", "fc057", "fc050", "fc051", "fc052", "fc053", "fc054", "fc055", "fc056", "fc060", "fc061", "fc062", "fc063", "fc064", "fc065", "fc066", "fc067", "fc068", "fc069", "fc070", "fc074", "fc075", "fc076", "fc077"}
	// "FC040":{
	// 	"FC045":VAL
	// }

	// fc087Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0],LOC[1],loc[2], -1)
	fc040Arrs := fieldLocationMap["fc040"]
	codeArr = []string{"fc045", "fc046", "fc047", "fc048", "fc044", "fc049", "fc057", "fc050", "fc051", "fc052", "fc053", "fc054", "fc055", "fc056", "fc060", "fc061", "fc062", "fc063", "fc064", "fc065", "fc066", "fc067", "fc068", "fc069", "fc070", "fc074", "fc075", "fc076", "fc077"}
	fc040Map := make(map[string]map[string]string, 0)
	for qq, fc040Arr := range fc040Arrs {
		fc040Loc := utils.GetFieldLoc(fieldLocationMap["fc040"], fc040Arr[0], fc040Arr[1], fc040Arr[2], -1)
		fc040Val := utils.GetFieldValueByLoc(obj, fc040Loc[0], false)
		fc055 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc055"][qq], false)
		if fc040Val != "" && fc055 != "" {
			key := fc040Val + "_" + fc055
			_, isExits := fc040Map[key]
			if !isExits {
				fc040Map[key] = map[string]string{}
			}
			_, fc040isExits := fc040Map[fc040Val+"_"]
			if !fc040isExits {
				fc040Map[fc040Val+"_"] = map[string]string{}
			}
			for _, code := range codeArr {
				codeLoc := utils.GetFieldLoc(fieldLocationMap[code], fc040Arr[0], fc040Arr[1], fc040Arr[2], -1)
				if len(codeLoc) > 0 {
					codeVal := utils.GetFieldValueByLoc(obj, codeLoc[0], true)
					if codeVal != "" {
						fc040Map[key][code] = codeVal
						fc040Map[fc040Val+"_"][code] = codeVal
					}
				}

			}
		}
	}
	for qq, fc040Arr := range fc040Arrs {
		fc040Loc := utils.GetFieldLoc(fieldLocationMap["fc040"], fc040Arr[0], fc040Arr[1], fc040Arr[2], -1)
		fc040Val := utils.GetFieldValueByLoc(obj, fc040Loc[0], false)
		fc055 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc055"][qq], false)

		codeValue, isExits := fc040Map[fc040Val+"_"+fc055]
		if isExits {
			for _, code := range codeArr {
				codeLoc := utils.GetFieldLoc(fieldLocationMap[code], fc040Arr[0], fc040Arr[1], fc040Arr[2], -1)
				utils.SetFinalValue(obj, codeLoc, codeValue[code])
			}
		}

	}

	//CSB0122RC0141000
	//当fc164、fc165、fc166、fc167、fc168、fc169、fc170、fc171、fc172、fc173（可能存在多个）
	// 结果值为“W53.951、W54.951、W55.951、W55.952、W56.851、W56.852、W57.951、W58.951、W59.951、
	// W60.951、W64.951、W64.952、W64.953、W64.954”时，
	//fc214结果值默认为X20-X29，fc216结果值默认为CS31
	// CSB0122RC0141001 调整新增字段：fc497、fc498、fc499、fc500、fc501、fc502、fc503、fc504、fc505、fc506
	if len(fieldLocationMap["fc214"]) > 0 || len(fieldLocationMap["fc214"]) > 0 {
		fieldsArrs := []string{"fc164", "fc165", "fc166", "fc167", "fc168", "fc169", "fc170", "fc171", "fc172", "fc173", "fc497", "fc498", "fc499", "fc500", "fc501", "fc502", "fc503", "fc504", "fc505", "fc506"}
		aaa := false
		for _, code := range fieldsArrs {
			codeValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap[code], true)
			for _, value := range codeValArr {
				if utils.RegIsMatch(`^(W53\.951|W54\.951|W55\.951|W55\.952|W56\.851|W56\.852|W57\.951|W58\.951|W59\.951|W60\.951|W64\.951|W64\.952|W64\.953|W64\.954)$`, value) {
					aaa = true
					break
				}
			}
			if aaa {
				break
			}
		}
		if aaa {
			if len(fieldLocationMap["fc214"]) > 0 {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc214"][0], "X20-X29")
			}
			if len(fieldLocationMap["fc216"]) > 0 {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc216"][0], "CS31")
			}
		}
	}

	//CSB0122RC0117000
	//将fc086、fc092、fc095、fc115的结果值按字段顺序汇总赋值给fc127，项目名称中间以“、”相隔。
	//CSB0122RC0119000
	//将fc088、fc097、fc118、fc125的结果值按字段顺序汇总赋值给fc129，项目名称中间以“、”相隔。
	arrRes := []string{"fc127", "fc129"}
	arr := [][]string{
		{"fc086", "fc092", "fc095", "fc115"},
		{"fc088", "fc097", "fc117", "fc125"},
	}
	for i, codeRes := range arrRes {
		for _, loc := range fieldLocationMap[codeRes] {
			totalVal := make([]string, 0)
			for _, code := range arr[i] {
				fc086Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
				if len(fc086Loc) != 1 {
					continue
				}
				fc086Val := utils.GetFieldValueByLoc(obj, fc086Loc[0], true)
				if fc086Val == "" {
					continue
				}
				totalVal = append(totalVal, fc086Val)
			}
			utils.SetOnlyOneFinalValue(obj, loc, strings.Join(totalVal, "、"))
		}
	}

	//CSB0122RC0118000
	//将fc087、fc096、fc117、fc124合计金额赋值给fc128
	codeArr = []string{"fc087", "fc096", "fc116", "fc124"}
	for _, loc := range fieldLocationMap["fc128"] {
		totalVal := decimal.Zero
		for _, code := range codeArr {
			fc087Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
			if len(fc087Loc) != 1 {
				continue
			}
			fc087Val := utils.GetFieldDecimalValueByLoc(obj, fc087Loc[0], true)
			totalVal = totalVal.Add(fc087Val)
		}
		utils.SetOnlyOneFinalValue(obj, loc, totalVal.StringFixed(2))
	}

	var wrong []string
	// 编码CSB0122RC0148000 导出校验 3.当最后一页字段包含?或？时，则出导出校验：最后一页存在?，请检查；
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

	global.GLog.Info("使用时间", zap.Any("", time.Since(startTime)))
	if len(wrong) > 1 {
		obj.Bill.WrongNote += "最后一页存在?，请检查;"
	}

	//编码：CSB0122RC0153000 时间：20230808  当销售渠道为“纸质”时，fc005录入值为A、C时，出问题件，问题件编码：C27，问题件描述：出险人身份无法确认（纸质申请书）
	if obj.Bill.SaleChannel == "纸质" {
		for _, loc := range fieldLocationMap["fc005"] {
			fc005obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
			if fc005obj.ResultValue == "A" || fc005obj.ResultValue == "C" {
				issue := model3.Issue{
					Type:    "",
					Code:    "C27",
					Message: "出险人身份无法确认（纸质申请书）",
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc005obj.Issues, issue)
			}
		}
	}
	//编码：CSB0122RC0154000 时间：20230808  销售渠道为“纸质”时，当fc015录入为2时，fc016结果值为空时，清空该字段已出的问题件，重新出问题件，问题件编码：C02，问题件描述：账户信息不完整（纸质申请书）
	//编码：CSB0122RC0155000 时间：20230808  销售渠道为“纸质”时，当fc015录入为2时，fc017录入值为A时，清空该字段已出的问题件，重新出问题件，问题件编码：C02，问题件描述：账户信息不完整（纸质申请书）
	//编码：CSB0122RC0156000 时间: 20230808  销售渠道为“纸质”时，当fc015录入为2时，fc022录入值为A时，清空该字段已出的问题件，重新出问题件，问题件编码：C02，问题件描述：账户信息不完整（纸质申请书）
	if obj.Bill.SaleChannel == "纸质" {
		//fieldLLocs := fieldLocationMap["fc015"]
		fieldRLocs := fieldLocationMap["fc016"]
		fc017Locs := fieldLocationMap["fc017"]
		fc022Locs := fieldLocationMap["fc022"]
		for i, loc := range fieldLocationMap["fc015"] {
			fc015Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
			fc016Obj := obj.Invoice[fieldRLocs[i][0]].Invoice[fieldRLocs[i][2]][fieldRLocs[i][3]]
			fc017Obj := obj.Invoice[fc017Locs[i][0]].Invoice[fc017Locs[i][2]][fc017Locs[i][3]]
			fc022Obj := obj.Invoice[fc022Locs[i][0]].Invoice[fc022Locs[i][2]][fc022Locs[i][3]]

			if fc015Obj.ResultValue == "2" && fc016Obj.ResultValue == "" && len(fieldRLocs) > 0 {
				utils.DelIssue(obj, fieldLocationMap["fc015"])
				issue := model3.Issue{
					Type:    "",
					Code:    "C02",
					Message: "账户信息不完整（纸质申请书）",
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc015Obj.Issues, issue)
			}
			if fc015Obj.ResultValue == "2" && fc017Obj.ResultValue == "A" && len(fc017Locs) > 0 {
				utils.DelIssue(obj, fieldLocationMap["fc015"])
				issue := model3.Issue{
					Type:    "",
					Code:    "C02",
					Message: "账户信息不完整（纸质申请书）",
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc015Obj.Issues, issue)
			}
			if fc015Obj.ResultValue == "2" && fc022Obj.ResultValue == "A" && len(fc022Locs) > 0 {
				utils.DelIssue(obj, fieldLocationMap["fc015"])
				issue := model3.Issue{
					Type:    "",
					Code:    "C02",
					Message: "账户信息不完整（纸质申请书）",
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc015Obj.Issues, issue)
			}
		}

		// 编码CSB0122RC0163000 时间20230809  销售渠道为“纸质”的案件，所有fc078的结果值默认为3
		if len(fieldLocationMap["fc078"]) > 0 {
			for _, item := range fieldLocationMap["fc078"] {
				utils.SetOnlyOneFinalValue(obj, item, "3")
			}
		}

	}
	// 编码 CSB0122RC0157000 时间：20230808 当fc037录入值不为B且不包含中文时，清空fc038、fc039结果值，fc037出问题件，问题件编码：C25，问题件描述：xxx未匹配到内部定点医院（xxx为fc037的结果值）
	fc038Loc := fieldLocationMap["fc038"]
	fc039Loc := fieldLocationMap["fc039"]
	for i, loc := range fieldLocationMap["fc037"] {
		fc037Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		fmt.Println("fc037Obj.ResultValue=", fc037Obj.ResultValue)
		compile := regexp.MustCompile("([\u4e00-\u9fa5]+)")
		matchString := compile.MatchString(fc037Obj.ResultValue)
		if fc037Obj.ResultValue != "B" && matchString == false {
			if len(fc038Loc) > 0 && len(fc039Loc) > 0 && len(fieldLocationMap["fc037"]) == len(fc038Loc) && len(fieldLocationMap["fc037"]) == len(fc039Loc) {
				//fc038Obj := obj.Invoice[fc038Loc[i][0]].Invoice[fc038Loc[i][2]][fc038Loc[i][3]]
				//fc039Obj := obj.Invoice[fc039Loc[i][0]].Invoice[fc039Loc[i][2]][fc039Loc[i][3]]
				obj.Invoice[fc038Loc[i][0]].Invoice[fc038Loc[i][2]][fc038Loc[i][3]].FinalValue = ""
				obj.Invoice[fc039Loc[i][0]].Invoice[fc039Loc[i][2]][fc039Loc[i][3]].FinalValue = ""
				issue := model3.Issue{
					Type:    "",
					Code:    "C25",
					Message: fc037Obj.FinalValue + "未匹配到内部定点医院",
				}
				obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]].Issues = append(fc037Obj.Issues, issue)
			}
		}
	}
	// //编码 CSB0122RC0159000 时间 20230809 校验每个fc042、fc043的结果值，与fc001作对比，存在不一致时，fc001出问题件，问题件编码：C18，问题件描述：医疗收据患者姓名与事故人不符
	// fc042ValArr := []string{}
	// fc042Arr := fieldLocationMap["fc042"]
	// for _, loc := range fc042Arr {
	// 	fc042Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
	// 	fc042ValArr = append(fc042ValArr, fc042Obj.FinalValue)
	// }
	// for _, loc := range fieldLocationMap["fc043"] {
	// 	fc043Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
	// 	fc042ValArr = append(fc042ValArr, fc043Obj.FinalValue)
	// }
	// for _, loc := range fieldLocationMap["fc001"] {
	// 	if len(fieldLocationMap["fc001"]) > 0 {
	// 		fc001Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
	// 		for _, item := range fc042ValArr {
	// 			if item != fc001Obj.FinalValue {
	// 				utils.SetIssue(obj, fieldLocationMap["fc001"][0], "医疗收据患者姓名与事故人不符", "C18", "")
	// 			}
	// 		}
	// 	}
	// }
	// 编码 CSB0122RC0168000  时间 20230809
	// fc083、fc084、fc085、fc087、fc089、fc090、fc091、fc093、fc094、fc096、fc098、fc099、fc101、fc103、fc105、fc107、fc109、fc111、fc113、fc114、fc116、fc118、fc119、fc120、fc121、fc122、fc124、fc126、fc128、fc130、fc131、
	// fc432、fc433、fc434、fc435、fc436、fc437、fc438、fc439、fc440、fc441、fc442、fc443、fc444、fc445、fc446、fc447、fc448、fc449、fc450、fc451、fc452、fc453、fc454、fc455、fc456、fc457、fc458、fc459、fc460、fc461
	// fc080、fc162、fc163、fc280、fc338、fc281
	// 以上结果值格式默认为保留两位小数（四舍五入），0或空均默认为0.00
	// fields := []string{"fc083", "fc084", "fc085", "fc087", "fc089", "fc090", "fc091", "fc093", "fc094", "fc096", "fc098", "fc099", "fc101", "fc103", "fc105", "fc107", "fc109", "fc111", "fc113", "fc114", "fc116", "fc118", "fc119", "fc120", "fc121", "fc122", "fc124", "fc126", "fc128", "fc130", "fc131",
	// 	"fc432", "fc433", "fc434", "fc435", "fc436", "fc437", "fc438", "fc439", "fc440", "fc441", "fc442", "fc443", "fc444", "fc445", "fc446", "fc447", "fc448", "fc449", "fc450", "fc451", "fc452", "fc453", "fc454", "fc455", "fc456", "fc457", "fc458", "fc459", "fc460", "fc461", "fc080", "fc162", "fc163"}
	// for _, invoiceMap := range obj.Invoice {
	// 	for _, invoice := range invoiceMap.Invoice {
	// 		for _, field := range invoice {
	// 			for _, code := range fields {
	// 				if field.Code == code {
	// 					//f3 := GetFieldValue(invoice, code, false, field.BlockIndex)
	// 					f3 := utils.GetFieldValueByLoc(obj, fieldLocationMap[field.Code][0], true)
	// 					if f3 == "" {
	// 						for _, item := range fieldLocationMap[field.Code] {
	// 							utils.SetOnlyOneFinalValue(obj, item, "0.00")
	// 						}
	// 					} else {
	// 						float, err := strconv.ParseFloat(f3, 64)
	// 						if err != nil {
	// 							fmt.Println(err)
	// 						}
	// 						// round := math.Round(float)
	// 						sprintf := fmt.Sprintf("%.2f", float)
	// 						utils.SetOnlyOneFinalValue(obj, fieldLocationMap[code][0], sprintf)
	// 					}

	// 				}
	// 			}
	// 		}
	// 	}
	// }

	fields := []string{"fc083", "fc084", "fc085", "fc087", "fc089", "fc090", "fc091", "fc093", "fc094", "fc096", "fc098", "fc099", "fc101", "fc103", "fc105", "fc107", "fc109", "fc111", "fc113", "fc114", "fc116", "fc118", "fc119", "fc120", "fc121", "fc122", "fc124", "fc126", "fc128", "fc130", "fc131",
		"fc432", "fc433", "fc434", "fc435", "fc436", "fc437", "fc438", "fc439", "fc440", "fc441", "fc442", "fc443", "fc444", "fc445", "fc446", "fc447", "fc448", "fc449", "fc450", "fc451", "fc452", "fc453", "fc454", "fc455", "fc456", "fc457", "fc458", "fc459", "fc460", "fc461", "fc080", "fc162", "fc163", "fc280", "fc338", "fc281", "fc479", "fc480", "fc481", "fc482", "fc483", "fc492", "fc493", "fc494", "fc495"}
	for _, code := range fields {
		if len(fieldLocationMap[code]) > 0 {
			for _, item := range fieldLocationMap[code] {
				f3 := utils.GetFieldValueByLoc(obj, item, true)
				if f3 == "" {
					// for _, item := range fieldLocationMap[code] {
					utils.SetOnlyOneFinalValue(obj, item, "0.00")
					// }
				} else {
					// float, err := strconv.ParseFloat(f3, 64)
					// if err != nil {
					// 	fmt.Println(err)
					// }
					// round := math.Round(float)
					// sprintf := fmt.Sprintf("%.2f", round)
					utils.SetOnlyOneFinalValue(obj, item, utils.ParseDecimal(f3).StringFixed(2))
				}
			}
		}
	}

	if len(fieldLocationMap["fc072"]) > 0 {
		for ii, item := range fieldLocationMap["fc072"] {
			fc071 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc071"][ii], false)
			fc072 := utils.GetFieldValueByLoc(obj, item, false)
			if fc071 == "1" && fc072 == "2" {
				utils.SetIssue(obj, item, "医疗收据无医院收款章", "C19", "")
			}
		}
	}

	// 编码 CSB0122RC0158000 时间 20230809 fc001默认下载报文accident_name节点值
	nodeRegex := `<accident_name>([^<]*)<\/accident_name>`
	r := regexp.MustCompile(nodeRegex)
	matches := r.FindAllStringSubmatch(obj.Bill.OtherInfo, -1)
	for _, match := range matches {
		nodeVale := match[1]
		if nodeVale != "" && len(fieldLocationMap["fc001"]) > 0 {
			utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc001"][0], nodeVale)
		}
	}

	//编码 CSB0122RC0159000 时间 20230809 校验每个fc042、fc043的结果值，与fc001作对比，存在不一致时，fc001出问题件，问题件编码：C18，问题件描述：医疗收据患者姓名与事故人不符
	fc042ValArr := []string{}
	fc042Arr := fieldLocationMap["fc042"]
	for _, loc := range fc042Arr {
		fc042Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		fc042ValArr = append(fc042ValArr, fc042Obj.FinalValue)
	}
	for _, loc := range fieldLocationMap["fc043"] {
		fc043Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
		fc042ValArr = append(fc042ValArr, fc043Obj.FinalValue)
	}
	for _, loc := range fieldLocationMap["fc001"] {
		if len(fieldLocationMap["fc001"]) > 0 {
			fc001Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
			for _, item := range fc042ValArr {
				if item != fc001Obj.FinalValue {
					utils.SetIssue(obj, fieldLocationMap["fc001"][0], "医疗收据患者姓名与事故人不符", "C18", "")
				}
			}
		}
	}

	// CSB0122RC0174000 时间：20230811 当fc071录入值为1，fc072录入值为2时，fc072出问题件，问题件编码：C19，问题件描述：医疗收据无医院收款章
	// 当fc033结果值与fc009结果值不一致，fc033出问题件：病人姓名与申请表事故者不一致
	fc071Locs := fieldLocationMap["fc071"]
	if len(fc071Locs) > 0 {
		fc071Obj := obj.Invoice[fc071Locs[0][0]].Invoice[fc071Locs[0][2]][fc071Locs[0][3]]
		for _, loc := range fieldLocationMap["fc072"] {
			fc072Obj := obj.Invoice[loc[0]].Invoice[loc[2]][loc[3]]
			if fc071Obj.ResultValue == "1" && fc072Obj.ResultInput == "2" {
				utils.SetIssue(obj, loc, "医疗收据无医院收款章", "C19", "")
			}
		}
	}
	// 编码 CSB0122RC0175000 时间 20230821  当左边字段录入值为4或5时，右边对应字段结果值默认为99999999
	sField := [][]string{
		{"fc025", "fc028"},
		{"fc223", "fc226"},
		{"fc233", "fc236"},
		{"fc243", "fc246"},
		{"fc253", "fc256"},
	}
	for _, codes := range sField {
		for _, loc := range fieldLocationMap[codes[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			loc1 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc[0], -1, -1, -1)
			if val == "4" || val == "5" {
				utils.SetFinalValue(obj, loc1, "99999999")
			}
		}
	}

	//CSB0122RC0116000	结果数据
	//"如存在MB002-bc020分块时，将fc338的结果值赋值为fc126
	//不存在时，将fc085、fc091、fc094、fc114合计金额赋值给fc126"
	for i, loc := range fieldLocationMap["fc126"] {
		invoiceMap := obj.Invoice[i]
		if len(invoiceMap.BaoXiaoDan) > 0 {
			fc338Loc := utils.GetFieldLoc(fieldLocationMap["fc338"], loc[0], -1, -1, -1)
			if len(fc338Loc) != 1 {
				continue
			}
			val := utils.GetFieldValueByLoc(obj, fc338Loc[0], true)
			utils.SetOnlyOneFinalValue(obj, loc, val)
		} else {
			fc085Loc := utils.GetFieldLoc(fieldLocationMap["fc085"], loc[0], -1, -1, -1)
			fc091Loc := utils.GetFieldLoc(fieldLocationMap["fc091"], loc[0], -1, -1, -1)
			fc094Loc := utils.GetFieldLoc(fieldLocationMap["fc094"], loc[0], -1, -1, -1)
			fc114Loc := utils.GetFieldLoc(fieldLocationMap["fc114"], loc[0], -1, -1, -1)
			if len(fc085Loc) != 1 || len(fc091Loc) != 1 || len(fc094Loc) != 1 || len(fc114Loc) != 1 {
				continue
			}
			fc085Val := utils.GetFieldDecimalValueByLoc(obj, fc085Loc[0], true)
			fc091Val := utils.GetFieldDecimalValueByLoc(obj, fc091Loc[0], true)
			fc094Val := utils.GetFieldDecimalValueByLoc(obj, fc094Loc[0], true)
			fc114Val := utils.GetFieldDecimalValueByLoc(obj, fc114Loc[0], true)
			utils.SetOnlyOneFinalValue(obj, loc, fc085Val.Add(fc091Val).Add(fc094Val).Add(fc114Val).StringFixed(2))
		}
	}

	// CSB0122RC0140000
	// 	1.第4列结果值为：第2列/第3列的值
	// 2.第5、6、7列结果值默认为：0.00
	// 3.第8、9、10列结果值默认为：0.00%
	// 4.第11列结果值默认为：无自付
	// 5.第12列结果值默认为：甲
	codes = [][]string{
		{"fc139", "fc142", "fc141", "fc140", "fc462", "fc151", "fc152", "fc143", "fc144", "fc149", "fc150", "fc136"},
		{"fc303", "fc324", "fc317", "fc310", "fc463", "fc388", "fc395", "fc331", "fc339", "fc374", "fc381", "fc282"},
		{"fc304", "fc325", "fc318", "fc311", "fc464", "fc389", "fc396", "fc332", "fc340", "fc375", "fc382", "fc283"},
		{"fc305", "fc326", "fc319", "fc312", "fc465", "fc390", "fc397", "fc333", "fc341", "fc376", "fc383", "fc284"},
		{"fc306", "fc327", "fc320", "fc313", "fc466", "fc391", "fc398", "fc334", "fc342", "fc377", "fc384", "fc285"},
		{"fc307", "fc328", "fc321", "fc314", "fc467", "fc392", "fc399", "fc335", "fc343", "fc378", "fc385", "fc286"},
		{"fc308", "fc329", "fc322", "fc315", "fc468", "fc393", "fc400", "fc336", "fc344", "fc379", "fc386", "fc287"},
		{"fc309", "fc330", "fc323", "fc316", "fc469", "fc394", "fc401", "fc337", "fc345", "fc380", "fc387", "fc288"},
	}
	if isFc164 {
		for _, codeArr = range codes {
			for _, loc6 := range fieldLocationMap[codeArr[5]] {
				loc2 := utils.GetFieldLoc(fieldLocationMap[codeArr[1]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc3 := utils.GetFieldLoc(fieldLocationMap[codeArr[2]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc4 := utils.GetFieldLoc(fieldLocationMap[codeArr[3]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc5 := utils.GetFieldLoc(fieldLocationMap[codeArr[4]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc7 := utils.GetFieldLoc(fieldLocationMap[codeArr[6]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc8 := utils.GetFieldLoc(fieldLocationMap[codeArr[7]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc9 := utils.GetFieldLoc(fieldLocationMap[codeArr[8]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc10 := utils.GetFieldLoc(fieldLocationMap[codeArr[9]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc11 := utils.GetFieldLoc(fieldLocationMap[codeArr[10]], loc6[0], loc6[1], loc6[2], loc6[4])
				loc12 := utils.GetFieldLoc(fieldLocationMap[codeArr[11]], loc6[0], loc6[1], loc6[2], loc6[4])

				if len(loc2) != 1 || len(loc3) != 1 || len(loc4) != 1 || len(loc5) != 1 || len(loc7) != 1 || len(loc8) != 1 || len(loc9) != 1 || len(loc10) != 1 {
					continue
				}
				//1.第5列字段的结果值为：第3列/第2列的值；
				valDecimal3 := utils.GetFieldDecimalValueByLoc(obj, loc3[0], true)
				valDecimal2 := utils.GetFieldDecimalValueByLoc(obj, loc2[0], true)
				if !valDecimal3.IsZero() {
					// utils.SetOnlyOneFinalValue(obj, loc5[0], valDecimal3.Div(valDecimal2).StringFixed(2))
					utils.SetFinalValue(obj, loc4, valDecimal2.Div(valDecimal3).StringFixed(2))
				}

				utils.SetFinalValue(obj, loc5, "0.00")
				utils.SetOnlyOneFinalValue(obj, loc6, "0.00")
				utils.SetFinalValue(obj, loc7, "0.00")

				utils.SetFinalValue(obj, loc8, "0.00%")
				utils.SetFinalValue(obj, loc9, "0.00%")
				utils.SetFinalValue(obj, loc10, "0.00%")

				utils.SetFinalValue(obj, loc11, "无自付")
				utils.SetFinalValue(obj, loc12, "甲")

			}
		}
	}

	// CSB0122RC0140000
	// fc164ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc164"], true)
	// isFc164 := false
	// for _, fc164 := range fc164ValArr {
	// 	if utils.RegIsMatch(`^(W53\.951|W54\.951|W55\.951|W55\.952|W56\.851|W56\.852|W57\.951|W58\.951|W59\.951|W60\.951|W64\.951|W64\.952|W64\.953|W64\.954)$`, fc164) {
	// 		isFc164 = true
	// 		break
	// 	}
	// }

	// for code, fieldLocs := range fieldLocationMap {
	// 	// fieldLocs := fieldLocationMap[fieldCode]
	// 	for _, loc := range fieldLocs {
	// 		invoiceMap := obj.Invoice[loc[0]]
	// 		eleLen := reflect.ValueOf(invoiceMap).NumField()
	// 		if eleLen > 0 {
	// 			if reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(loc[1]).Kind() != reflect.Float64 {
	// 				//每张发票每种类型的字段
	// 				fieldsArr := reflect.ValueOf(invoiceMap).Field(loc[1]).Interface().([][]model3.ProjectField)

	// 				// // CSB0122RC0140000
	// 				if isFc164 && utils.RegIsMatch(`^(fc144|fc339|fc340|fc341|fc342|fc343|fc344|fc345)$`, code) {
	// 					if fieldsArr[loc[2]][loc[3]].FinalValue != "" && fieldsArr[loc[2]][loc[3]].FinalValue != "0%" && fieldsArr[loc[2]][loc[3]].FinalValue != "0.00%" {
	// 						fieldsArr[loc[2]][loc[3]].FinalValue = "0.00%"
	// 					}
	// 				}

	// 			}
	// 		}
	// 	}
	// }

	if len(fieldLocationMap["fc474"]) > 0 {
		fc474 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc474"][0], false)
		if fc474 != "F" && fc474 != "" {
			items := GetAddress(fc474, constMap["dizhiQu"])
			if len(items) == 1 {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc474"][0], items[0])
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc475"][0], Substr(items[0], 0, 4)+"00")
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], Substr(items[0], 0, 2)+"0000")
			} else if len(items) > 1 {
				fc475 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc475"][0], false)
				if fc475 != "F" && fc475 != "" {
					fc475Items := GetAddress(fc475, constMap["dizhiShi"])
					// for _, fc475Item := range fc475Items {
					if len(fc475Items) == 1 {
						for _, item := range items {
							if fc475Items[0] == Substr(item, 0, 4)+"00" {
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc474"][0], item)
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc475"][0], fc475Items[0])
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], Substr(item, 0, 2)+"0000")
							}
						}
					} else {
						utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc474"][0], items[0])
						fc476 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc476"][0], false)
						if fc476 != "F" && fc476 != "" {
							fc476Items := GetAddress(fc476, constMap["dizhiSheng"])
							if len(fc476Items) > 0 {
								utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], fc476Items[0])
							}
						}

					}
				}

			}
		} else {
			fc475 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc475"][0], false)
			if fc475 != "F" && fc475 != "" {
				fc475Items := GetAddress(fc475, constMap["dizhiShi"])
				if len(fc475Items) == 1 {
					utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc475"][0], fc475Items[0])
					utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], Substr(fc475Items[0], 0, 2)+"0000")
				} else {
					fc476 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc476"][0], false)
					if fc476 != "F" && fc476 != "" {
						fc476Items := GetAddress(fc476, constMap["dizhiSheng"])
						if len(fc476Items) > 0 {
							utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], fc476Items[0])
						}
					}
				}
			} else {
				fc476 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc476"][0], false)
				if fc476 != "F" && fc476 != "" {
					fc476Items := GetAddress(fc476, constMap["dizhiSheng"])
					if len(fc476Items) > 0 {
						utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], fc476Items[0])
					}
				}
			}
		}
	}

	if len(fieldLocationMap["fc205"]) > 0 {
		fc205 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc205"][0], false)
		if fc205 != "F" && fc205 != "" {
			fc205Items := GetAddress(fc205, constMap["dizhiShi"])
			if len(fc205Items) == 1 {
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc205"][0], fc205Items[0])
				utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc204"][0], Substr(fc205Items[0], 0, 2)+"0000")
			} else {
				fc204 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc204"][0], false)
				if fc204 != "F" && fc204 != "" {
					fc204Items := GetAddress(fc204, constMap["dizhiSheng"])
					if len(fc204Items) > 0 {
						utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc204"][0], fc204Items[0])
					}
				}
			}
		} else {
			fc204 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc204"][0], false)
			if fc204 != "F" && fc204 != "" {
				fc204Items := GetAddress(fc204, constMap["dizhiSheng"])
				if len(fc204Items) > 0 {
					utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc204"][0], fc204Items[0])
				}
			}
		}
	}

	if len(fieldLocationMap["fc476"]) > 0 {
		fc476 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc476"][0], true)
		// if RegIsMatch(fc476, `^(110000|120000|310000|500000)$`) {
		// 	utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc475"][0], "")
		// }
		fc475 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc475"][0], true)
		fc474 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc474"][0], true)
		if fc475 != "" {
			fc476 = fc476 + "," + fc475
		}
		utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc476"][0], fc476+","+fc474)
	}

	if len(fieldLocationMap["fc021"]) > 0 {
		fc021 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc021"][0], true)
		if fc021 == "" && len(fieldLocationMap["fc264"]) > 0 {
			utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc021"][0], utils.GetFieldValueByLoc(obj, fieldLocationMap["fc264"][0], true))
		}
	}

	if len(fieldLocationMap["fc015"]) > 0 {
		fc015 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc015"][0], false)
		if fc015 == "2" {
			fc021 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc021"][0], true)
			if fc021 == "" {
				utils.SetIssue(obj, fieldLocationMap["fc021"][0], "银行卡/存折缺失或不清晰", "C03", "")

			}
		}
	}

	if len(fieldLocationMap["fc022"]) > 0 {
		fc022 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc022"][0], true)
		if fc022 == "" && len(fieldLocationMap["fc263"]) > 0 {
			utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc022"][0], utils.GetFieldValueByLoc(obj, fieldLocationMap["fc263"][0], true))
		}
	}

	mes := ""
	isNotCheck := true

	mapfc037 := map[string]string{}
	for ii, invoice := range obj.Invoice {
		fc280 := 0.0
		fc281 := 0.0
		for _, baoxiao := range invoice.BaoXiaoDan {
			for _, field := range baoxiao {
				if field.Code == "fc280" && field.FinalValue != "" {
					fc280 = utils.SumFloat(fc280, utils.ParseFloat(field.FinalValue), "+")
				}
			}
		}
		for _, thirdBaoXiaoDan1 := range invoice.ThirdBaoXiaoDan1 {
			for _, field := range thirdBaoXiaoDan1 {
				if field.Code == "fc281" && field.FinalValue != "" {
					fc281 = utils.SumFloat(fc281, utils.ParseFloat(field.FinalValue), "+")
				}
			}
		}
		ifhuli := true
		ischuanwei := true
		fc089 := []int{}
		fc083 := []int{}
		isChange := false
		fc486 := ""
		fc487 := ""
		fc079 := ""
		fc488 := ""
		fc489 := ""
		fc490 := ""
		fc491_f := ""
		for jj, fields := range invoice.Invoice {
			// fc037 :=
			for ff, field := range fields {
				if field.Code == "fc037" && field.FinalValue != "" {
					_, fc081 := GetOneField(fields, "fc081", true)
					_, fc082 := GetOneField(fields, "fc082", true)
					if fc081 != "" || fc082 != "" {
						mapfc037[field.FinalValue+"_"+"_"] = field.FinalValue + "_" + fc081 + "_" + fc082
						obj.Invoice[ii].Invoice[jj][ff].FinalValue = field.FinalValue + "_" + fc081 + "_" + fc082
					}

				}
				if field.Code == "fc040" && field.FinalValue != "" {
					_, fc055 := GetOneField(fields, "fc055", true)
					_, fc056 := GetOneField(fields, "fc056", true)
					if fc055 != "" || fc056 != "" {
						obj.Invoice[ii].Invoice[jj][ff].FinalValue = field.FinalValue + "_" + fc055 + "_" + fc056
					}

				}
				if field.Code == "fc089" {
					fc089 = []int{jj, ff}
				}
				if field.Code == "fc083" {
					fc083 = []int{jj, ff}
				}
				if RegIsMatch(field.Code, `^(fc402|fc403|fc404|fc405|fc406|fc407|fc408|fc409|fc410|fc411|fc412|fc413)$`) && strings.Index(field.ResultValue, "护理费") != -1 {
					ifhuli = false
				}
				if RegIsMatch(field.Code, `^(fc402|fc403|fc404|fc405|fc406|fc407|fc408|fc409|fc410|fc411|fc412|fc413)$`) && strings.Index(field.ResultValue, "床位费") != -1 {
					ischuanwei = false
				}
				if field.Code == "fc130" && fc280 != 0.0 {
					obj.Invoice[ii].Invoice[jj][ff].FinalValue = utils.ToString(fc280)
				}
				if field.Code == "fc131" && fc281 != 0.0 {
					obj.Invoice[ii].Invoice[jj][ff].FinalValue = utils.ToString(fc281)
				}

				if field.Code == "fc490" {
					_, fc486 := GetOneField(fields, "fc486", false)

					if fc486 == "2" {
						_, fc080 := GetOneField(fields, "fc080", true)
						obj.Invoice[ii].Invoice[jj][ff].FinalValue = fc080
					}

				}

				if RegIsMatch(field.Code, `^(fc486|fc487|fc488|fc489|fc490|fc491)$`) {
					_, fc135 := GetOneField(fields, "fc135", false)

					if fc135 == "2" {
						obj.Invoice[ii].Invoice[jj][ff].FinalValue = ""
					}
				}

				if field.Code == "fc486" {
					fc486 = field.FinalValue
				}
				if field.Code == "fc487" {
					fc487 = field.FinalValue
				}
				if field.Code == "fc079" {
					fc079 = field.FinalValue
				}
				if field.Code == "fc488" {
					fc488 = field.FinalValue
				}
				if field.Code == "fc489" {
					fc489 = field.FinalValue
				}
				if field.Code == "fc490" {
					fc490 = obj.Invoice[ii].Invoice[jj][ff].FinalValue
				}
				if field.Code == "fc491" {
					fc491_f = field.FinalValue
				}

				if RegIsMatch(field.Code, `^(fc486|fc487|fc079|fc488|fc489|fc490)$`) && field.IsChange {
					isChange = true
				}
			}
		}
		fmt.Println("------------bill.Stage-------------------:", bill.Stage, isChange)
		if (bill.Stage == 6 && fc491_f == "" && fc486 != "") || (bill.Stage != 6 && isChange && strings.Contains(bill.Remark, "<<发票查验>>")) {
			bodyData := make(map[string]interface{})
			bodyData["fpdm"] = fc487
			bodyData["fphm"] = fc079
			bodyData["kprq"] = fc488
			bodyData["checkCode"] = fc489
			fmt.Println("------------bodyData-------------------:", invoice.Code, bodyData)
			fmt.Println("------------fc490-------------------:", fc490)
			fc491 := ""
			if fc486 == "1" {
				bodyData["noTaxAmount"] = fc490
				err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData)
				isNotCheck = false
				if err != nil {
					mes += invoice.Code + fmt.Sprintf("%v", err) + ";"
					// response.FailWithMessage(fmt.Sprintf("%v", err), c)
					// return
				} else {
					// mes += invoice.Code + "发票查验成功;"
					del := respData.Data.Del
					// if del == "0" {
					// 	fc491 = "01"
					// }
					if del == "0" {
						fc491 = "1"
					}
					if del == "2" || del == "3" || del == "7" || del == "8" {
						fc491 = "2"
					}
				}

			} else if fc486 == "2" {
				bodyData["money"] = fc490
				err, respData := unitFunc.Invoice("/v2/eInvoice/query", bodyData)
				isNotCheck = false
				if err != nil {
					mes += invoice.Code + fmt.Sprintf("%v", err) + ";"
					if RegIsMatch(fmt.Sprintf("%v", err), `(800|801|1098),`) {
						fc491 = "3"
					} else {
						fc491 = "2"
					}
					// response.FailWithMessage(fmt.Sprintf("%v", err), c)
					// return
				} else {
					// mes += invoice.Code + "发票查验成功;"
					if respData.Data.IsRed {
						fc491 = "2"
					} else if respData.Data.IsRed == false {
						fc491 = "1"
					}
					// else {
					// 	if respData.Data.IsPrint {
					// 		fc491 = "03"
					// 	} else {
					// 		fc491 = "01"
					// 	}
					// }
				}

			}

			fmt.Println("------------fc491-------------------:", fc491)
			for jj, fields := range invoice.Invoice {
				for kk, field := range fields {
					if field.Code == "fc491" {
						obj.Invoice[ii].Invoice[jj][kk].ResultValue = fc491
						obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc491
						obj.Invoice[ii].Invoice[jj][kk].IsChange = true
					}
				}
			}
		}
		codes = [][]string{
			{"fc139", "fc142"},
			{"fc303", "fc324"},
			{"fc304", "fc325"},
			{"fc305", "fc326"},
			{"fc306", "fc327"},
			{"fc307", "fc328"},
			{"fc308", "fc329"},
			{"fc309", "fc330"},
		}
		if ifhuli && len(fc089) > 0 {
			sum := 0.0
			for _, fields := range invoice.QingDan {
				for _, code := range codes {
					_, aa := GetOneField(fields, code[0], false)
					if strings.Index(aa, "护理") != -1 {
						_, bb := GetOneField(fields, code[1], true)
						sum = utils.SumFloat(sum, utils.ParseFloat(bb), "+")
					}
				}
			}
			obj.Invoice[ii].Invoice[fc089[0]][fc089[1]].FinalValue = utils.ToString(sum)
			// utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc089"][ii], utils.ToString(sum))
		}
		if ischuanwei && len(fc083) > 0 {
			sum := 0.0
			for _, fields := range invoice.QingDan {
				for _, code := range codes {
					_, aa := GetOneField(fields, code[0], false)
					if RegIsMatch(aa, `(床位|人间|病房|陪护床|走廊)`) && !RegIsMatch(aa, `(取暖|空调)`) {
						_, bb := GetOneField(fields, code[1], true)
						sum = utils.SumFloat(sum, utils.ParseFloat(bb), "+")
					}
				}
			}
			obj.Invoice[ii].Invoice[fc083[0]][fc083[1]].FinalValue = utils.ToString(sum)
			// utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc083"][ii], utils.ToString(sum))
		}

	}

	if bill.Stage != 6 && isNotCheck && strings.Contains(bill.Remark, "<<发票查验>>") {
		mes += "查验失败：发票查验相关字段的数据未作修改，不执行发票查验功能;"
	}
	// bill.Remark = strings.Replace(bill.Remark, "<<发票查验>>", "", 1)
	fmt.Println("------------mes-------------------:", mes)
	if mes != "" {
		obj.Bill.WrongNote += mes
	}

	for _, loc := range fieldLocationMap["fc037"] {
		val := utils.GetFieldValueByLoc(obj, loc, true)
		aa, bb := mapfc037[val+"_"+"_"]
		if bb {
			utils.SetOnlyOneFinalValue(obj, loc, aa)
		}
	}
	for _, loc := range fieldLocationMap["fc040"] {
		val := utils.GetFieldValueByLoc(obj, loc, true)
		aa, bb := mapfc037[val+"_"+"_"]
		if bb {
			utils.SetOnlyOneFinalValue(obj, loc, aa)
		}
	}

	for ff, _ := range fieldLocationMap["fc048"] {
		fc048 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc048"][ff], false)
		if strings.Index(fc048, "特需") != -1 {
			utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc044"][ff], "1")
		} else {
			utils.SetOnlyOneFinalValue(obj, fieldLocationMap["fc044"][ff], "0")
		}
	}

	//CSB0122RC0070000
	//"fc161、fc402、fc403、fc404、fc405、fc406、fc407、fc408、fc409、fc410、fc411、
	//fc412、fc413、fc414、fc415、fc416、fc417、fc418、fc419、fc420、fc421、
	//fc422、fc423、fc424、fc425、fc426、fc427、fc428、fc429、fc430、fc431
	//以上字段录入值需根据机构号是否为 北京 来匹配对应的常量表，并转码输出
	//匹配规则：
	//北京匹配《B0122_信诚理赔_北京发票大项明细表》的“名称”（第二列）转“代码”（第一列）
	//非北京匹配《B0122_信诚理赔_发票大项明细表》的“名称”（第二列）转“代码”（第一列）"
	codeArr = []string{
		"fc161", "fc402", "fc403", "fc404", "fc405", "fc406", "fc407", "fc408", "fc409", "fc410", "fc411",
		"fc412", "fc413", "fc414", "fc415", "fc416", "fc417", "fc418", "fc419", "fc420", "fc421",
		"fc422", "fc423", "fc424", "fc425", "fc426", "fc427", "fc428", "fc429", "fc430", "fc431", "fc138", "fc296", "fc297", "fc298", "fc299", "fc300", "fc301", "fc302"}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			name := utils.GetFieldValueByLoc(obj, loc, false)
			if RegIsMatch(code, `^(fc138|fc296|fc297|fc298|fc299|fc300|fc301|fc302)$`) {
				name = utils.GetFieldValueByLoc(obj, loc, true)
			}
			val, _ := constMap["beiJingFaPiaoDaXiangMap"][name]
			if obj.Bill.Agency != "北京" {
				val, _ = constMap["faPiaoDaXiangMap"][name]
			}
			utils.SetOnlyOneFinalValue(obj, loc, val)
		}
	}

	//CSB0122RC0005001
	// 	同一发票下，（以下字段如存在则发票唯一，字段不存在时不执行该校验）
	// 将fc492的结果值赋值给fc085
	// 将fc493的结果值赋值给fc091
	// 将fc494的结果值赋值给fc094
	// 将fc495的结果值赋值给fc122
	codes = [][]string{
		{"fc085", "fc492"},
		{"fc091", "fc493"},
		{"fc094", "fc494"},
		{"fc122", "fc495"},
	}
	for i, invoices := range obj.Invoice {
		if len(invoices.BaoXiaoDan) == 0 {
			continue
		}
		for j, fields := range invoices.Invoice {
			for k, field := range fields {
				for _, code := range codes {
					if code[0] == field.Code {
						isExit, val := GetOneField(invoices.BaoXiaoDan[0], code[1], true)
						if isExit {
							obj.Invoice[i].Invoice[j][k].FinalValue = val
						}

					}
				}
			}
		}
	}

	//CSB0122RC0099000	结果数据
	//"将fc094的结果值赋值给fc105
	//将fc095的结果值赋值给fc106"
	myArr = [][]string{
		{"fc094", "fc105"},
		{"fc095", "fc106"},
	}
	for _, codeArr = range myArr {
		for _, loc := range fieldLocationMap[codeArr[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			loc1 := utils.GetFieldLoc(fieldLocationMap[codeArr[1]], loc[0], -1, -1, -1)
			utils.SetFinalValue(obj, loc1, val)
		}
	}

	//CSB0122RC0098000	结果数据	同一发票下，fc098的值=fc093-fc094，如fc093的值为0或空时，不进行计算，fc098默认为0.00
	for _, loc := range fieldLocationMap["fc098"] {
		fc093Loc := utils.GetFieldLoc(fieldLocationMap["fc093"], loc[0], -1, -1, -1)
		fc094Loc := utils.GetFieldLoc(fieldLocationMap["fc094"], loc[0], -1, -1, -1)
		if len(fc093Loc) != 1 || len(fc094Loc) != 1 {
			continue
		}
		fc093Val := utils.GetFieldDecimalValueByLoc(obj, fc093Loc[0], true)
		fc094Val := utils.GetFieldDecimalValueByLoc(obj, fc094Loc[0], true)
		if RegIsMatch(utils.GetFieldValueByLoc(obj, fc093Loc[0], true), `^(0|0\.00|)$`) {
			utils.SetOnlyOneFinalValue(obj, loc, "0.00")
		} else {
			utils.SetOnlyOneFinalValue(obj, loc, fc093Val.Sub(fc094Val).StringFixed(2))
		}

	}

	if len(fieldLocationMap["fc079"]) > 0 {
		fc079Map := map[string]int{}
		for _, item := range fieldLocationMap["fc079"] {
			f3 := utils.GetFieldValueByLoc(obj, item, true)
			if f3 != "" {
				fc079Num, isExist := fc079Map[f3]
				if isExist {
					fc079Map[f3] = fc079Num + 1
				} else {
					fc079Map[f3] = 1
				}
			}
		}
		fc079One := ""
		fc079All := ""
		for key, value := range fc079Map {
			if value > 1 {
				if fc079One == "" {
					fc079One = key
				} else {
					fc079One += "、" + key
				}
			}
			// err, num := service.GetCountValueByCode(bill.ProCode, obj.Bill.ID, "fc079", key)
			// if err == nil && num > 0 {
			// 	if fc079All == "" {
			// 		fc079All = key
			// 	} else {
			// 		fc079All += "、" + key
			// 	}
			// }
		}
		if fc079One != "" {
			utils.SetIssue(obj, fieldLocationMap["fc079"][0], "本次案件有重复收据("+fc079One+")", "C24", "")
		}
		if fc079All != "" {
			utils.SetIssue(obj, fieldLocationMap["fc079"][0], "有重复收据号码("+fc079All+")", "C12", "")
		}
	}

	if len(fieldLocationMap["fc496"]) > 0 {
		fc496 := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc496"][0], false)
		arrs := RegSplit(fc496, `(；|;)`)
		for _, arr := range arrs {
			if arr != "" {
				if utils.HasKey(constMap["questionCode"], arr) {
					utils.SetIssue(obj, fieldLocationMap["fc496"][0], constMap["questionMes"][arr], constMap["questionCode"][arr], "")
				}
			}
		}
	}

	////////////////////////项目报表 - 业务明细 - 单据状态
	obj.Bill.BillType = eUtils.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	return nil, obj
}

func GetAddress(fvalue string, tempMap map[string]string) []string {
	items := []string{}
	for key, value := range tempMap {
		if value == fvalue {
			items = append(items, key)
		}
	}
	return items
}

func GetFieldValue(fields []model3.ProjectField, code string, isFinal bool, bidx int) string {
	for _, field := range fields {
		if field.Code == code && (field.BlockIndex == bidx || bidx == -1) {
			if isFinal {
				return field.FinalValue
			} else {
				return field.ResultValue
			}
		}
	}
	return ""
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"yiYuanMap", "B0122_信诚理赔_内部医院代码表", "1", "0"},
		{"zbxYiYuanMap", "B0122_信诚理赔_中保信医院代码表", "1", "0"},
		{"jiBingDaiMaMap", "B0122_信诚理赔_伤病代码表", "1", "0"},
		{"shouShuDaiMaMap", "B0122_信诚理赔_手术代码表", "2", "0"},
		{"yinHangDaiMaMap", "B0122_信诚理赔_银行代码表", "1", "2"},
		{"kaiHuChengShiMap", "B0122_信诚理赔_开户城市表", "2", "0"},
		{"zhongJiMap", "B0122_信诚理赔_重疾原因代码表", "1", "0"},
		{"yiWaiMap", "B0122_信诚理赔_意外原因代码表", "2", "1"},
		{"yiWaiDiMap", "B0122_信诚理赔_意外发生地", "2", "1"},
		{"shenGuMap", "B0122_信诚理赔_身故原因代码表", "1", "0"},
		{"neiBuYiYuanMap", "B0122_信诚理赔_内部医院代码表", "1", "2"},
		{"zhongBaoXinYiYuanMap", "B0122_信诚理赔_中保信医院代码表", "1", "2"},
		{"beiJingFaPiaoDaXiangMap", "B0122_信诚理赔_北京发票大项明细表", "1", "0"},
		{"faPiaoDaXiangMap", "B0122_信诚理赔_发票大项明细表", "1", "0"},
		{"quanGuoMap", "B0122_信诚理赔_全国", "1", "4"},
		{"quanGuo1Map", "B0122_信诚理赔_全国", "1", "2"},

		{"baXiangYao", "B0122_信诚理赔_靶向药清单", "1,2", "0"},
		{"huaLiaoYao", "B0122_信诚理赔_化疗药清单", "1,2", "0"},
		{"kangPaiCiYao", "B0122_信诚理赔_器官移植抗排斥药", "1", "0"},

		{"dizhiQu", "B0122_信诚理赔_行政区划代码", "0", "3"},
		{"dizhiShi", "B0122_信诚理赔_行政区划代码", "0", "2"},
		{"dizhiSheng", "B0122_信诚理赔_行政区划代码", "0", "1"},

		{"questionCode", "B0122_信诚理赔_问题件代码表", "0", "1"},
		{"questionMes", "B0122_信诚理赔_问题件代码表", "0", "2"},
	}
	for _, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				keys := strings.Split(item[2], ",")
				for _, key := range keys {
					k, _ := strconv.Atoi(key)
					v, _ := strconv.Atoi(item[3])
					tempMap[strings.TrimSpace(arr[k])] = arr[v]
				}
			}
		}
		constObj[item[0]] = tempMap
	}

	return constObj
}
