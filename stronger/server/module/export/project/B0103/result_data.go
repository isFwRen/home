/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年04月23日16:27:47
 */

package B0103

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"server/global"
	model2 "server/module/export/model"
	utils2 "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	unitFunc "server/module/unit"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

//B0103广西贵州国寿理赔
//模板类型字段 : fc003
//账单号: fc005票据号
//总金额: fc011发票总金额
//发票类型: fc009诊疗方式
//
//和发票有对应关系的
//fc018发票属性 MB002-bc001
//fc019清单所属发票 MB002-bc002
//fc020报销单所属发票 MB002-bc003
//fc085诊断书所属发票 MB002-bc004
//fc024第三方支付所属发票 MB002-bc005
//
//和发票没有对应关系的
//无

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:          model2.TypeCode{FieldCode: []string{"fc018"}, BlockCode: []string{"bc001"}},
	QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc019"}, BlockCode: []string{"bc002"}},
	BaoXiaoDanFieldCode:       model2.TypeCode{FieldCode: []string{"fc020"}, BlockCode: []string{"bc003"}},
	ZhenDuanFieldCode:         model2.TypeCode{FieldCode: []string{"fc085"}, BlockCode: []string{"bc004"}},
	ThirdBaoXiaoDan2FieldCode: model2.TypeCode{FieldCode: []string{"fc024"}, BlockCode: []string{"bc005"}},
	OtherTempType:             map[string]string{},
	TempTypeField:             "fc003",
	InvoiceNumField:           []string{"fc005"},
	MoneyField:                []string{"fc011"},
	InvoiceTypeField:          "fc009",
}

// ResultData B0103
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {

	defer func() {
		if err := recover(); err != nil {
			global.GLog.Error("", zap.Any("", err))
			debug.PrintStack()
		}
	}()
	obj = utils2.RelationDealManyInvoiceType(bill, blocks, fieldMap, relationConf)
	//常量
	//constMap := constDeal(bill.ProCode)
	////特殊的常量
	//constSpecialMap := constSpecialDeal(bill.ProCode)

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

	//CSB0103RC0038000
	//当fc005录入内容为A或B或?时，结果值自动生成一个20位数字，且该数字不能与之前的所有fc005内容相同（格式为工号+录入日期+编号）
	for _, fc005Loc := range fieldLocationMap["fc005"] {
		fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc, false)
		nextID := strconv.FormatInt(utils.GWorker.NextId(), 10)
		if utils.RegIsMatch("^(A|B|\\?)$", fc005Val) {
			utils.SetOnlyOneFinalValue(obj, fc005Loc, nextID)
			utils.SetIssue(obj, fc005Loc, "发票号无法辨识", "", "")
		}
	}

	//CSB0103RC0039000
	//fc022录入内容为A、B，fc023结果值为正常日期时，将fc023的结果值复制给fc022
	for _, loc := range fieldLocationMap["fc022"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		loc1 := utils.GetFieldLoc(fieldLocationMap["fc023"], loc[0], loc[1], loc[2], loc[4])
		if len(loc1) < 1 {
			continue
		}
		val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
		_, err := time.Parse("2006-01-02", val1)
		if utils.RegIsMatch("^(A|B)$", val) && err == nil {
			utils.SetOnlyOneFinalValue(obj, loc, val1)
		}
	}

	//CSB0103RC0040000
	//fc023录入内容为A、B，fc022结果值为正常日期时，将fc022的结果值复制给fc023
	for _, loc := range fieldLocationMap["fc023"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		loc1 := utils.GetFieldLoc(fieldLocationMap["fc022"], loc[0], loc[1], loc[2], loc[4])
		if len(loc1) < 1 {
			continue
		}
		val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
		_, err := time.Parse("2006-01-02", val1)
		if utils.RegIsMatch("^(A|B)$", val) && err == nil {
			utils.SetOnlyOneFinalValue(obj, loc, val1)
		}
	}

	//CSB0103RC0041000
	//"fc046为非录入字段时，fc010有录入内容且录入内容不为A不为0时，fc055默认为10；
	//fc046为非录入字段时，fc010为空或为A或为0时，fc055默认为00；
	//fc046为录入字段时，将fc046的结果数据复制给fc055；"
	for _, fc010Loc := range fieldLocationMap["fc010"] {
		fc010Val := utils.GetFieldValueByLoc(obj, fc010Loc, false)
		fc055Loc := utils.GetFieldLoc(fieldLocationMap["fc055"], fc010Loc[0], -1, -1, -1)
		fc046Loc := utils.GetFieldLoc(fieldLocationMap["fc046"], fc010Loc[0], -1, -1, -1)
		if len(fc055Loc) < 1 {
			continue
		}
		if len(fc046Loc) < 1 && (fc010Val != "0" && fc010Val != "A") {
			utils.SetOnlyOneFinalValue(obj, fc055Loc[0], "10")
		}
		if len(fc046Loc) < 1 && (fc010Val == "" || fc010Val == "A" || fc010Val == "0") {
			utils.SetOnlyOneFinalValue(obj, fc055Loc[0], "00")
		}
		if len(fc046Loc) > 0 {
			fc046Val := utils.GetFieldValueByLoc(obj, fc046Loc[0], true)
			utils.SetOnlyOneFinalValue(obj, fc055Loc[0], fc046Val)
		}
	}

	//CSB0103RC0042000
	//"1.将以下右边字段的结果数据放到对应的左边的字段中（右边的字段是报销单中的字段，左边的字段是发票的字段，此校验放在后面）
	//2.当不存在报销单时，将同一个发票中的fc017赋值到同分块中的fc061中
	//3.当fc016值为空时，fc010不需要取fc016的值"
	//myCodes := [][]string{
	//	{"fc057", "fc015"},
	//	{"fc010", "fc016"},
	//	{"fc060", "fc086"},
	//	{"fc061", "fc017"},
	//}
	myCodes := [][]string{
		{"fc015", "fc057"},
		{"fc016", "fc010"},
		{"fc086", "fc060"},
		{"fc017", "fc061"},
	}

	for _, codes := range myCodes {
		for _, loc := range fieldLocationMap[codes[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			loc1 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc[0], -1, -1, -1)
			loc010 := utils.GetFieldLoc(fieldLocationMap["fc010"], loc[0], -1, -1, -1) //fc010
			val010 := utils.GetFieldValueByLoc(obj, loc010[0], false)

			if len(loc1) < 1 {
				continue
			}
			val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)

			if codes[0] == "fc016" && val != "" && val010 == "" {
				utils.SetFinalValue(obj, loc010, val)
			}
			if codes[0] == "fc016" && val1 == "" {
				//3.当fc016值为空时，fc010不需要取fc016的值"
				continue

			}
			utils.SetFinalValue(obj, loc1, val)
		}
	}

	//设置fc010结果值
	for _, fc010Loc := range fieldLocationMap["fc010"] {
		fc010Val := utils.GetFieldValueByLoc(obj, fc010Loc, false) // 10 录入
		val := utils.GetFieldValueByLoc(obj, fc010Loc, true)       // 10结果值
		if val == "" {
			utils.SetOnlyOneFinalValue(obj, fc010Loc, fc010Val)
			fmt.Println("fc010Val=", fc010Val)
		}

	}

	//CSB0103RC0043000
	//fc008的内容根据案件“机构号”去匹配常量表《B0103_广西贵州国寿理赔_数据库编码对应表》中“本级机构代码”（第二列）对应的“医疗目录编码”（第四列）输出结果值
	//bianMaItems := constSpecialMap["shuJuKuBianMaMap"][bill.Agency]
	v, total := utils.FetchConst(bill.ProCode, "B0103_广西贵州国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	global.GLog.Info("B0103_广西贵州国寿理赔_数据库编码对应表,Agency : len(bianMaItems)", zap.Any(bill.Agency, v))
	for _, loc := range fieldLocationMap["fc008"] {
		//if len(bianMaItems) < 4 {
		if total == 0 {
			continue
		}
		utils.SetOnlyOneFinalValue(obj, loc, v)
	}

	//CSB0103RC0044000 20230628新增 ;CSB0103RC0045000 20230628新增 "下列字段中，右边字段的结果数据根据左边字段录入内容与对应医疗目录中“项目名称”（第三列）对应匹配的“自费比例”（第八列）进行输出,如未匹配到自费比例时,默认为1
	//"fc012、fc025、fc026、fc027、fc028、fc029、fc030、fc031根据录入内容匹配对应常量表进行转码输出，当无法转码时不作处理；
	//匹配规则：
	//案件“机构号”匹配常量表《B0103_广西贵州国寿理赔_数据库编码对应表》中“本级机构代码”（第二列）对应的“医疗目录编码”（第四列）的常量表中的“项目名称”（第三列）对应的“医疗项目编码”（第二列）进行转码输出；"
	muLu, total := utils.FetchConst(bill.ProCode, "B0103_广西贵州国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	//bianCode, ok := constSpecialMap["shuJuKuBianMaMap"][bill.Agency]
	//myCode := []string{"fc012", "fc025", "fc026", "fc027", "fc028", "fc029", "fc030", "fc031"}
	myCodes = [][]string{
		{"fc012", "fc047"},
		{"fc025", "fc048"},
		{"fc026", "fc049"},
		{"fc027", "fc050"},
		{"fc028", "fc051"},
		{"fc029", "fc052"},
		{"fc030", "fc053"},
		{"fc031", "fc054"},
	}
	if total != 0 {
		//if ok && len(bianCode) > 3 {
		//yiLiaoMuLu := "医疗目录" + bianCode[3]
		yiLiaoMuLu := "B0103_广西贵州国寿理赔_医疗目录" + muLu
		otherInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(bill.OtherInfo), &otherInfo)
		if err != nil {
			global.GLog.Error("CSB0103RC0044000", zap.Error(err))
		}
		for _, codes := range myCodes {
			for _, twoLoc := range fieldLocationMap[codes[1]] {
				oneLoc := utils.GetFieldLoc(fieldLocationMap[codes[0]], twoLoc[0], twoLoc[1], twoLoc[2], twoLoc[4])
				if len(oneLoc) < 1 {
					continue
				}
				//yiLiaoMuLuConst, _ := constSpecialMap[yiLiaoMuLu]
				oneVal := utils.GetFieldValueByLoc(obj, oneLoc[0], false)
				xiangMu, total1 := utils.FetchConst(bill.ProCode, yiLiaoMuLu, "医疗项目编码", map[string]string{"项目名称": oneVal})
				biLi, _ := utils.FetchConst(bill.ProCode, yiLiaoMuLu, "自费比例", map[string]string{"项目名称": oneVal})
				//medicalCode, isOk := yiLiaoMuLuConst[oneVal]
				//CSB0110RC0045000 项目比例默认1 （左边字段为空时，不执行该校验）
				//if oneVal != "" && isOk == false { //有值 匹配不上就 1
				if oneVal != "" && total1 == 0 { //有值 匹配不上就 1
					utils.SetOnlyOneFinalValue(obj, twoLoc, "1")
				} else if total1 != 0 { //匹配上转码
					utils.SetOnlyOneFinalValue(obj, oneLoc[0], xiangMu)
					utils.SetOnlyOneFinalValue(obj, twoLoc, biLi)
				}
			}
		}
	}

	//CSB0103RC0046000
	//fc062、fc064、fc066、fc068、fc070和fc072的结果值为根据录入值匹配《B0103_广西贵州国寿理赔_第三方出具单位》中的“第三方出具单位名称”（第三列）转换为对应的“第三方出具单位代码”（第二列）
	codeArr := []string{"fc062", "fc064", "fc066", "fc068", "fc070", "fc072"}
	for _, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			v, _ = utils.FetchConst(bill.ProCode, "B0103_广西贵州国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": val})
			//utils.SetOnlyOneFinalValue(obj, loc, constMap["diSanFangChuJuDanWeiMap"][val])
			utils.SetOnlyOneFinalValue(obj, loc, v)
		}
	}
	//CSB0103RC0048000
	//fc006结果值根据《B0103_广西贵州国寿理赔_医疗机构52》的“医院名称”（第一列）转成对应的“医院编号”（第二列）输出，无法转换时默认为“4401008000000003”，并出问题件：医院为数据库以外
	for _, loc := range fieldLocationMap["fc006"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		//v, ok := constMap["yiLiaoJiGouMap"][val]
		v, total = utils.FetchConst(bill.ProCode, "B0103_广西贵州国寿理赔_医疗机构52", "医院编码", map[string]string{"医院名称": val})
		if total != 0 {
			utils.SetOnlyOneFinalValue(obj, loc, v)
		} else {
			utils.SetOnlyOneFinalValue(obj, loc, "4401008000000003")
			utils.SetIssue(obj, loc, "医院为数据库以外", "", "")
		}
	}

	//CSB0103RC0049000 // 20230628新增
	// 录入内容为空、为A、为B时，fc006结果值默认为“4401008000000003”，并出问题件：医院内容模糊无法辨识录入
	for _, fc006Loc := range fieldLocationMap["fc006"] {
		fc005Val := utils.GetFieldValueByLoc(obj, fc006Loc, false)
		//nextID := strconv.FormatInt(utils.GWorker.NextId(), 10)
		if utils.RegIsMatch("^(|A|B|)$", fc005Val) {
			utils.SetOnlyOneFinalValue(obj, fc006Loc, "4401008000000003")
			utils.SetIssue(obj, fc006Loc, "医院内容模糊无法辨识录入", "", "")
		}
	}

	//CSB0103RC0050000 // 20230628新增
	//fc007、fc074、fc075、fc076、fc077、fc078、fc079、fc080、fc081、fc082结果值根据《B0103_广西贵州国寿理赔_ICD10疾病编码》的“疾病名称”（第一列）转成对应的“疾病代码”（第二列）输出
	MCode := []string{"fc007", "fc074", "fc075", "fc076", "fc077", "fc078", "fc079", "fc080", "fc081", "fc082"}
	for _, code := range MCode {
		for _, item := range fieldLocationMap[code] {
			itemVal := utils.GetFieldValueByLoc(obj, item, false)
			v, total = utils.FetchConst(bill.ProCode, "B0103_广西贵州国寿理赔_ICD10疾病编码", "疾病代码", map[string]string{"疾病名称": itemVal})
			//_, ok := constMap["jiBingBianMaMap"][itemVal]
			//if ok {
			if total != 0 {
				//utils.SetOnlyOneFinalValue(obj, item, constMap["jiBingBianMaMap"][itemVal])
				utils.SetOnlyOneFinalValue(obj, item, v)
			}
		}

	}

	//CSB0103RC0051000 // 20230629新增
	//"1、fc007录入内容为A、B、?时，fc007默认为S32，出问题件：疾病诊断无法辨识
	//2、fc074、fc075、fc076、fc077、fc078、fc079、fc080、fc081、fc082以上字段录入内容为A时，默认为空；录入内容为B、?时，默认为S32，出问题件：疾病诊断无法辨识"
	MCodes := []string{"fc074", "fc075", "fc076", "fc077", "fc078", "fc079", "fc080", "fc081", "fc082"}
	for _, loc := range fieldLocationMap["fc007"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		if val == "?" || val == "A" || val == "B" {
			utils.SetOnlyOneFinalValue(obj, loc, "S32")
			utils.SetIssue(obj, loc, "疾病诊断无法辨识", "", "")
		}
	}

	for _, code := range MCodes {
		for _, item := range fieldLocationMap[code] {
			itemVal := utils.GetFieldValueByLoc(obj, item, false)
			if itemVal == "A" {
				utils.SetOnlyOneFinalValue(obj, item, "")
			}

			if itemVal == "B" || itemVal == "?" {
				utils.SetOnlyOneFinalValue(obj, item, "S32")
				utils.SetIssue(obj, item, "疾病诊断无法辨识", "", "")
			}
		}
	}

	//CSB0103RC0047000 // 20230629新增
	//"同一张发票下的MB002-bc003报销单分块中，
	//fc062、fc064、fc066、fc068、fc070和fc072的结果值用“|”符号整合，放到fc062结果值中，
	//fc063、fc065、fc067、fc069、fc071和fc073的结果值用“|”符号整合，放到fc063结果值中；
	//在同一发票属性存在多个fc062及fc063的结果值时，将多个的fc062、fc063的结果值用“|”符号整合，放到最后一个fc062、fc063的结果值中；
	//（结果值为空的字段不进行以上整合操作，注意结果数据的位置进行数据的整合，此条校验在最后进行）"
	//mCodes := []string{"fc062", "fc064", "fc066", "fc068", "fc070", "fc072"}
	//var arrsys []string
	//var emptyString []string
	for iii, invoiceMap := range obj.Invoice {
		if len(invoiceMap.Invoice) > 0 {
			_, fc091 := GetOneField(invoiceMap.Invoice[0], "fc091", false)
			if fc091 == "2" || fc091 == "3" {
				_, fc011 := GetOneField(invoiceMap.Invoice[0], "fc011", true)
				for jjj, field := range invoiceMap.Invoice[0] {
					if field.Code == "fc096" {
						obj.Invoice[iii].Invoice[0][jjj].FinalValue = fc011
					}
					// if RegIsMatch(field.Code, `^(fc091|fc092|fc093|fc094|fc095|fc096|fc097)$`) {
					// 	obj.Invoice[iii].Invoice[0][jjj].FinalValue = ""
					// }
				}
			}
			_, fc090 := GetOneField(invoiceMap.Invoice[0], "fc090", false)
			if fc090 == "2" {
				for jjj, field := range invoiceMap.Invoice[0] {
					if RegIsMatch(field.Code, `^(fc091|fc092|fc093|fc094|fc095|fc096|fc097)$`) {
						obj.Invoice[iii].Invoice[0][jjj].FinalValue = ""
					}
				}
			}
		}

		var arrsys []string
		var emptyString []string
		if len(invoiceMap.BaoXiaoDan) > 0 {
			for _, fields := range invoiceMap.BaoXiaoDan {
				for _, field := range fields {
					if RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) && field.FinalValue != "" {
						arrsys = append(arrsys, field.FinalValue)
					}
					if RegIsMatch(field.Code, `^(fc063|fc065|fc067|fc069|fc071|fc073)$`) && field.FinalValue != "" {
						emptyString = append(emptyString, field.FinalValue)
					}
				}
			}
		}
		if len(invoiceMap.ThirdBaoXiaoDan1) > 0 {
			for _, fields := range invoiceMap.ThirdBaoXiaoDan1 {
				for _, field := range fields {
					if RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) && field.FinalValue != "" {
						arrsys = append(arrsys, field.FinalValue)
					}
					if RegIsMatch(field.Code, `^(fc063|fc065|fc067|fc069|fc071|fc073)$`) && field.FinalValue != "" {
						emptyString = append(emptyString, field.FinalValue)
					}
				}
			}
		}
		if len(invoiceMap.ThirdBaoXiaoDan2) > 0 {
			for _, fields := range invoiceMap.ThirdBaoXiaoDan2 {
				for _, field := range fields {
					if RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) && field.FinalValue != "" {
						arrsys = append(arrsys, field.FinalValue)
					}
					if RegIsMatch(field.Code, `^(fc063|fc065|fc067|fc069|fc071|fc073)$`) && field.FinalValue != "" {
						emptyString = append(emptyString, field.FinalValue)
					}
				}
			}
		}
		if len(invoiceMap.ThirdBaoXiaoDan3) > 0 {
			for _, fields := range invoiceMap.ThirdBaoXiaoDan3 {
				for _, field := range fields {
					if RegIsMatch(field.Code, `^(fc062|fc064|fc066|fc068|fc070|fc072)$`) && field.FinalValue != "" {
						arrsys = append(arrsys, field.FinalValue)
					}
					if RegIsMatch(field.Code, `^(fc063|fc065|fc067|fc069|fc071|fc073)$`) && field.FinalValue != "" {
						emptyString = append(emptyString, field.FinalValue)
					}
				}
			}
		}
		locArr := utils.GetFieldLoc(fieldLocationMap["fc062"], iii, -1, -1, -1)
		fc063locArr := utils.GetFieldLoc(fieldLocationMap["fc063"], iii, -1, -1, -1)
		if len(locArr) < 1 {
			continue
		}
		if len(fc063locArr) < 1 {
			continue
		}
		utils.SetOnlyOneFinalValue(obj, locArr[0], strings.Join(arrsys, "|"))
		utils.SetOnlyOneFinalValue(obj, fc063locArr[0], strings.Join(emptyString, "|"))
	}
	//for _, fc018Loc := range fieldLocationMap["fc018"] {
	//	locArr := utils.GetFieldLoc(fieldLocationMap["fc062"], fc018Loc[0], -1, -1, -1)
	//	fc063locArr := utils.GetFieldLoc(fieldLocationMap["fc063"], fc018Loc[0], -1, -1, -1)
	//	if len(locArr) < 1 {
	//		continue
	//	}
	//	if len(fc063locArr) < 1 {
	//		continue
	//	}
	//	utils.SetOnlyOneFinalValue(obj, locArr[0], strings.Join(arrsys, "|"))
	//	utils.SetOnlyOneFinalValue(obj, fc063locArr[0], strings.Join(emptyString, "|"))
	//
	//}

	mes := ""
	isNotCheck := true
	for ii, invoice := range obj.Invoice {
		fc091 := ""
		fc092 := ""
		// fc093 := ""
		fc094 := ""
		fc095 := ""
		fc096 := ""
		fc005 := ""
		fc097F := ""
		isChange := false
		for jj, fields := range invoice.Invoice {
			for kk, field := range fields {

				// if "fc096" == field.Code {
				// 	_, fc091 := GetOneField(fields, "fc091", false)
				// 	if fc091 == "2" {
				// 		_, fc011 := GetOneField(fields, "fc011", true)
				// 		obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc011
				// 	}
				// }

				// if "fc093" == field.Code {
				// 	_, fc005 := GetOneField(fields, "fc005", true)
				// 	obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc005
				// }
				if field.Code == "fc097" {
					fc097F = field.FinalValue
				}
				if field.Code == "fc005" {
					fc005 = field.FinalValue
				}
				if field.Code == "fc091" {
					fc091 = field.ResultValue
				}
				if field.Code == "fc092" {
					fc092 = field.FinalValue
				}
				// if field.Code == "fc093" {
				// 	fc093 = obj.Invoice[ii].Invoice[jj][kk].FinalValue
				// }
				if field.Code == "fc094" {
					fc094 = field.FinalValue
				}
				if field.Code == "fc095" {
					fc095 = field.FinalValue
				}
				if field.Code == "fc096" {
					fc096 = obj.Invoice[ii].Invoice[jj][kk].FinalValue
				}
				if RegIsMatch(field.Code, `^(fc091|fc092|fc005|fc094|fc095|fc096)$`) && field.IsChange {
					isChange = true
				}

				//CSB0103RC0126000
				if RegIsMatch(field.Code, `^(fc022|fc023)$`) {
					val := obj.Invoice[ii].Invoice[jj][kk].ResultValue
					if val == "?" || val == "？" {
						currentTime := time.Now().Format("2006-01-02")
						obj.Invoice[ii].Invoice[jj][kk].FinalValue = currentTime
						issue := model3.Issue{
							Type:    "",
							Code:    "",
							Message: "【" + fc005 + "】发票日期填写模糊；",
						}
						obj.Invoice[ii].Invoice[jj][kk].Issues = append(obj.Invoice[ii].Invoice[jj][kk].Issues, issue)
					}
				}

			}
		}
		fmt.Println("------------bill.Stage-------------------:", bill.Stage, isChange, bill.Remark)
		if fc097F == "" && fc091 != "" {
			if bill.Stage == 6 || (bill.Stage != 6 && isChange && strings.Contains(bill.Remark, "<<发票查验>>")) {
				bodyData := make(map[string]interface{})
				bodyData["fpdm"] = fc092
				bodyData["fphm"] = fc005
				bodyData["kprq"] = fc094
				bodyData["checkCode"] = fc095
				fmt.Println("------------bodyData-------------------:", invoice.Code, bodyData)
				fc097 := ""
				if fc091 == "1" {
					bodyData["noTaxAmount"] = fc096
					err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData)
					isNotCheck = false
					if err != nil {
						mes += fc005 + fmt.Sprintf("%v", err) + ";"
						fc097 = "99"
						// response.FailWithMessage(fmt.Sprintf("%v", err), c)
						// return
					} else {
						// mes += fc005 + "发票查验成功;"
						del := respData.Data.Del
						if del == "0" {
							fc097 = "01"
						}
						if del == "2" {
							fc097 = "99"
						}
						if del == "3" || del == "7" || del == "8" {
							fc097 = "02"
						}
					}

				} else if fc091 == "2" {
					bodyData["money"] = fc096
					err, respData := unitFunc.Invoice("/v2/eInvoice/query", bodyData)
					isNotCheck = false
					if err != nil {
						mes += fc005 + fmt.Sprintf("%v", err) + ";"
						// if RegIsMatch(fmt.Sprintf("%v", err), `(800|801|1098),`) {
						// fmt.Println("----------------------------------------------------------errerrerrerr-------------------:", err)
						fc097 = "99"
						// }
						// response.FailWithMessage(fmt.Sprintf("%v", err), c)
						// return
					} else {
						// mes += fc005 + "发票查验成功;"
						if respData.Data.IsRed {
							fc097 = "02"
						} else {
							if respData.Data.IsPrint {
								fc097 = "03"
							} else {
								fc097 = "01"
							}
						}
					}

				} else if fc091 == "3" {
					bodyData1 := make(map[string]interface{})
					bodyData1["fphm"] = fc005
					bodyData1["kprq"] = fc094
					bodyData1["jshj"] = fc096
					err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData1)
					isNotCheck = false
					if err != nil {
						mes += "查验失败：" + fmt.Sprintf("%v", err) + ";"
						fc097 = "99"
						// response.FailWithMessage(fmt.Sprintf("%v", err), c)
						//return
					} else {
						// mes += fc005 + "发票查验成功;"
						del := respData.Data.Del
						if del == "0" {
							fc097 = "01"
						}
						if del == "2" {
							fc097 = "99"
						}
						if del == "3" || del == "7" || del == "8" {
							fc097 = "02"
						}
					}
				}
				fmt.Println("------------fc097-------------------:", fc097)
				for jj, fields := range invoice.Invoice {
					for kk, field := range fields {
						if field.Code == "fc097" {
							obj.Invoice[ii].Invoice[jj][kk].ResultValue = fc097
							obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc097
							obj.Invoice[ii].Invoice[jj][kk].IsChange = true
						}
					}
				}
			}
		}

	}

	//CSB0103RC0128000
	for _, loc := range fieldLocationMap["fc098"] {
		val := utils.GetFieldValueByLoc(obj, loc, true)
		result := strings.FieldsFunc(val, func(r rune) bool {
			return r == ';' || r == '；'
		})
		for _, issuesResult := range result {
			var strBuilder strings.Builder
			strBuilder.WriteString(issuesResult)
			issuesMsg := strBuilder.String()
			utils.SetIssue(obj, loc, issuesMsg, "", "")
		}
		utils.SetOnlyOneFinalValue(obj, loc, "")
	}

	if bill.Stage != 6 && isNotCheck && strings.Contains(bill.Remark, "<<发票查验>>") {
		mes += "查验失败：发票查验相关字段的数据未作修改，不执行发票查验功能;"
	}
	// bill.Remark = strings.Replace(bill.Remark, "<<发票查验>>", "", 1)
	fmt.Println("------------mes-------------------:", mes)
	if mes != "" {
		obj.Bill.WrongNote += mes
	}

	obj.Bill.BillType = utils2.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	return nil, obj
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"jiBingBianMaMap", "B0103_广西贵州国寿理赔_ICD10疾病编码", "0", "1"},
		{"yiLiaoJiGouMap", "B0103_广西贵州国寿理赔_医疗机构52", "0", "1"},
		{"diSanFangChuJuDanWeiMap", "B0103_广西贵州国寿理赔_第三方出具单位", "2", "1"},
		{"codeCorrespondence", "B0103_广西贵州国寿理赔_数据库编码对应表", "0", "3"},
		{"yiLiaoMuLu002Map", "B0103_广西贵州国寿理赔_医疗目录522100002", "2", "1"},
	}
	for _, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])
				if v == len(arr) {
					arr = append(arr, strconv.Itoa(k))
				}
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
		{"shuJuKuBianMaMap", "B0103_广西贵州国寿理赔_数据库编码对应表", "1"},
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

	//医疗目录  机构
	for k, items := range global.GProConf[proCode].ConstTable {
		if strings.Index(k, "医疗目录") != -1 {
			tempMap := make(map[string][]string, 0)
			for _, arr := range items {
				tempMap[arr[2]] = arr
			}
			constObj[strings.Replace(k, "B0103_广西贵州国寿理赔_", "", -1)] = tempMap
		}
	}

	return constObj
}

// RemoverEmptyString 删除数组中空字符串
func RemoverEmptyString(arr []string) []string {
	i := 0
	for _, item := range arr {
		if item != "" {
			arr[i] = item
			i++
		}
	}
	return arr[:i]

}
