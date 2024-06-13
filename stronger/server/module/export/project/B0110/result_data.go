/**
 * @Author: xingqiyi
 * @Description:这个项目没有编码 咨询说不要了 赋值113的
 * @Date: 2023年04月23日16:28:48
 */

package B0110

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"server/global"
	model2 "server/module/export/model"
	"server/module/export/service"
	utils2 "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	unitFunc "server/module/unit"
	"server/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

//B0110新疆国寿理赔
//模板类型字段 : fc003
//账单号: fc005票据号
//总金额: fc011发票总金额
//发票类型: fc009诊疗方式
//
//和发票有对应关系的
//fc018发票属性 MB002-bc001
//fc019清单所属发票 MB002-bc002
//fc020报销单所属发票 MB002-bc003
//fc063第三方支付所属发票 MB002-bc005
//
//和发票没有对应关系的
//4-诊断书

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:          model2.TypeCode{FieldCode: []string{"fc018"}, BlockCode: []string{"bc001"}},
	QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc019"}, BlockCode: []string{"bc002"}},
	BaoXiaoDanFieldCode:       model2.TypeCode{FieldCode: []string{"fc020"}, BlockCode: []string{"bc003"}},
	ThirdBaoXiaoDan1FieldCode: model2.TypeCode{FieldCode: []string{"fc063"}, BlockCode: []string{"bc005"}},
	OtherTempType:             map[string]string{"4": "4"},
	TempTypeField:             "fc003",
	InvoiceNumField:           []string{"fc005"},
	MoneyField:                []string{"fc011"},
	//InvoiceTypeField:          "fc108",
}

// ResultData B0110
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

	//CSB0110RC0033000
	//当fc005录入内容为空或A或B或?时，自动生成一个20位数字，且该数字不能与之前的所有fc005内容相同（格式为工号+录入日期+编号）
	//CSB0110RC0088000
	//fc005录入内容为空、A、B、?时，出问题件：发票号无法辨识；	新增需求
	for _, fc005Loc := range fieldLocationMap["fc005"] {
		fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc, false)
		nextID := strconv.FormatInt(utils.GWorker.NextId(), 10)
		if utils.RegIsMatch("^(A|B|\\?)$", fc005Val) {
			utils.SetOnlyOneFinalValue(obj, fc005Loc, nextID)
			utils.SetIssue(obj, fc005Loc, "发票号无法辨识", "", "")
		}
	}

	//CSB0110RC0034000
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

	//CSB0110RC0035000
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

	//CSB0110RC0036000
	//"fc015为非录入字段时，fc021录入内容为1时，fc016默认为10；
	//fc015为非录入字段时，fc021录入内容为2时，fc016默认为00；
	//fc015为录入字段时，将fc015的结果数据复制给fc016；"
	for _, fc021Loc := range fieldLocationMap["fc021"] {
		fc021Val := utils.GetFieldValueByLoc(obj, fc021Loc, false)
		fc016Loc := utils.GetFieldLoc(fieldLocationMap["fc016"], fc021Loc[0], -1, -1, -1)
		fc015Loc := utils.GetFieldLoc(fieldLocationMap["fc015"], fc021Loc[0], -1, -1, -1)
		if len(fc016Loc) < 1 {
			continue
		}
		if len(fc015Loc) < 1 && fc021Val == "1" {
			utils.SetOnlyOneFinalValue(obj, fc016Loc[0], "10")
		}
		if len(fc015Loc) < 1 && fc021Val == "2" {
			utils.SetOnlyOneFinalValue(obj, fc016Loc[0], "00")
		}
		if len(fc015Loc) > 0 {
			fc015Val := utils.GetFieldValueByLoc(obj, fc015Loc[0], true)
			utils.SetOnlyOneFinalValue(obj, fc016Loc[0], fc015Val)
		}
	}

	//CSB0110RC0037000
	//fc061的结果值为根据录入值匹配《B0110_新疆国寿理赔_第三方出具单位》中的“第三方出具单位名称”（第三列）转换为对应的“第三方出具单位代码”（第二列）
	for _, loc := range fieldLocationMap["fc061"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		v, _ := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": val})
		utils.SetOnlyOneFinalValue(obj, loc, v)
		//utils.SetOnlyOneFinalValue(obj, loc, constMap["diSanFangDanWeiMap"][val])
	}

	//CSB0110RC0038000
	//在同一发票属性存在多个fc061及fc062的结果值时，将多个的fc061、fc062的结果值用“|”符号整合，放到最后一个fc061、fc062的结果值中（注意结果数据的位置进行数据的整合）
	codes := []string{"fc061", "fc062"}
	for _, fc018Loc := range fieldLocationMap["fc018"] {
		for _, code := range codes {
			locArr := utils.GetFieldLoc(fieldLocationMap[code], fc018Loc[0], -1, -1, -1)
			var valArr []string
			for i, loc := range locArr {
				valArr = append(valArr, utils.GetFieldValueByLoc(obj, loc, true))
				if i == len(locArr)-1 {
					utils.SetOnlyOneFinalValue(obj, locArr[0], strings.Join(valArr, "|"))
				}
			}
		}
	}

	//CSB0110RC0039000
	//机构号为650192，下载报文中<bpoSendRemark>节点包含“全民”字样，且该单据中发票fc009录入为1时，同分块的fc021结果值默认为1，fc010结果值默认为0
	if bill.Agency == "650192" {
		otherInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(bill.OtherInfo), &otherInfo)
		if err != nil {
			global.GLog.Error("CSB0110RC0039000", zap.Error(err))
		}
		if strings.Index(otherInfo["bpoSendRemark"].(string), "全民") != -1 {
			for _, fc009Loc := range fieldLocationMap["fc009"] {
				fc009Val := utils.GetFieldValueByLoc(obj, fc009Loc, false)
				fc021Loc := utils.GetFieldLoc(fieldLocationMap["fc021"], fc009Loc[0], fc009Loc[1], fc009Loc[2], fc009Loc[4])
				fc010Loc := utils.GetFieldLoc(fieldLocationMap["fc010"], fc009Loc[0], fc009Loc[1], fc009Loc[2], fc009Loc[4])
				if len(fc021Loc) < 1 || len(fc010Loc) < 1 {
					continue
				}
				if fc009Val == "1" {
					utils.SetFinalValue(obj, fc021Loc, "1")
					utils.SetFinalValue(obj, fc010Loc, "0")
				}
			}
		}
	}

	//CSB0110RC0040000
	//"以下左边字段为报销单的字段，右边是报销单对应的发票的字段，当fc018与fc020相匹配时，需要将对应的左边的字段的值，复制到右边的字段中
	myCode := [][]string{
		{"fc046", "fc051"},
		{"fc047", "fc052"},
	}
	for _, code := range myCode {
		for _, loc := range fieldLocationMap[code[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			loc1 := utils.GetFieldLoc(fieldLocationMap[code[1]], loc[0], -1, -1, -1)
			if len(loc1) < 1 {
				continue
			}
			utils.SetFinalValue(obj, loc1, val)
		}
	}

	//CSB0110RC0041000
	//"机构号匹配常量表《B0110_新疆国寿理赔_数据库编码对应表》中的“本级机构代码”（第二列），匹配上的内容对应的“上级机构代码”（第三列）为650100（其中机构号为650192时，不执行此条校验；）或652300时，同一发票属性下
	//fc046的值大于fc017的值时，fc046的结果数据默认为fc046-fc017的值；
	//fc046的值小于或等于fc017的值时，fc046的结果值默认为空；
	//（fc017、fc046为0或空时，不计算；）"
	//CSB0110RC0042000
	//"机构号匹配常量表《B0110_新疆国寿理赔_数据库编码对应表》中的“本级机构代码”（第二列），匹配上的内容对应的“上级机构代码”（第三列）为652200时，
	//若fc046+fc047的值小于fc017的值，fc024的结果数据默认为fc024-fc017的值"
	//CSB0110RC0043000
	//"fc012、fc025、fc026、fc027、fc028、fc029、fc030、fc031根据录入内容匹配对应常量表进行转码输出，当无法转码时不作处理；
	//匹配规则：
	//案件“机构号”匹配常量表《B0110_新疆国寿理赔_数据库编码对应表》中“本级机构代码”（第二列）对应的“医疗目录编码”（第四列）的常量表中的“项目名称”（第三列）对应的“医疗项目编码”（第二列）进行转码输出；"
	//CSB0110RC0043000
	myCodes := [][]string{
		{"fc012", "fc053"},
		{"fc025", "fc054"},
		{"fc026", "fc055"},
		{"fc027", "fc056"},
		{"fc028", "fc057"},
		{"fc029", "fc058"},
		{"fc030", "fc059"},
		{"fc031", "fc060"},
	}
	for _, codes = range myCodes {
		for _, twoLoc := range fieldLocationMap[codes[1]] {
			oneLoc := utils.GetFieldLoc(fieldLocationMap[codes[0]], twoLoc[0], twoLoc[1], twoLoc[2], twoLoc[4])
			if len(oneLoc) < 1 {
				continue
			}
			oneVal := utils.GetFieldValueByLoc(obj, oneLoc[0], false)
			//CSB0110RC0045000 项目比例默认1 （左边字段为空时，不执行该校验）
			if oneVal != "" {
				utils.SetOnlyOneFinalValue(obj, twoLoc, "1")
			}
		}
	}
	jiGou, total := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_数据库编码对应表", "上级机构代码", map[string]string{"本级机构代码": bill.Agency})
	muLu, _ := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	//shuJuKuBianMaItems, ok := constSpecialMap["shuJuKuBianMa"][bill.Agency]
	//if ok && len(shuJuKuBianMaItems) > 2 {
	if total > 0 {
		//parentAgency := shuJuKuBianMaItems[2]
		//yiLiaoMuLu := "医疗目录" + shuJuKuBianMaItems[3]
		parentAgency := jiGou
		yiLiaoMuLu := "B0110_新疆国寿理赔_医疗目录" + muLu
		for _, fc046Loc := range fieldLocationMap["fc046"] {
			fc046Val := utils.GetFieldDecimalValueByLoc(obj, fc046Loc, true)
			fc017Loc := utils.GetFieldLoc(fieldLocationMap["fc017"], fc046Loc[0], -1, -1, -1)
			fc047Loc := utils.GetFieldLoc(fieldLocationMap["fc047"], fc046Loc[0], -1, -1, -1)
			fc024Loc := utils.GetFieldLoc(fieldLocationMap["fc024"], fc046Loc[0], -1, -1, -1)
			if len(fc017Loc) < 1 || len(fc047Loc) < 1 || len(fc024Loc) < 1 {
				continue
			}
			fc017Val := utils.GetFieldDecimalValueByLoc(obj, fc017Loc[0], true)
			fc047Val := utils.GetFieldDecimalValueByLoc(obj, fc047Loc[0], true)
			fc024Val := utils.GetFieldDecimalValueByLoc(obj, fc024Loc[0], true)
			if fc017Val.IsZero() && fc046Val.IsZero() {
				continue
			}
			//CSB0110RC0041000
			if bill.Agency != "650192" && (parentAgency == "650100" || parentAgency == "652300") {
				utils.SetOnlyOneFinalValue(obj, fc046Loc, "")
				if fc046Val.GreaterThan(fc017Val) {
					utils.SetOnlyOneFinalValue(obj, fc046Loc, fc046Val.Sub(fc017Val).String())
				}
			}
			//CSB0110RC0042000
			if parentAgency == "652200" {
				if fc046Val.Add(fc047Val).LessThan(fc017Val) {
					utils.SetOnlyOneFinalValue(obj, fc024Loc[0], fc024Val.Sub(fc017Val).String())
				}
			}
		}

		for _, codes := range myCodes {
			for _, oneLoc := range fieldLocationMap[codes[0]] {
				oneVal := utils.GetFieldValueByLoc(obj, oneLoc, false)
				twoLoc := utils.GetFieldLoc(fieldLocationMap[codes[1]], oneLoc[0], oneLoc[1], oneLoc[2], oneLoc[4])
				fc006Loc := utils.GetFieldLoc(fieldLocationMap["fc006"], oneLoc[0], -1, -1, -1)
				if len(twoLoc) < 1 || len(fc006Loc) < 1 {
					continue
				}

				xiangMuMingCheng := oneVal
				if strings.Index(oneVal, "床位费") != -1 {
					fc006Val := utils.GetFieldValueByLoc(obj, fc006Loc[0], false)
					dengJi, total1 := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_医疗机构65", "医院等级", map[string]string{"医院名称": fc006Val})
					bianHao, _ := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_医疗机构65", "医疗项目编码", map[string]string{"医院名称": fc006Val})
					//yiLiaoJiGouItems, ok3 := constSpecialMap["医疗机构65"][fc006Val]
					//if !ok3 {
					if total1 == 0 {
						//未匹配要收集
						utils.SetOnlyOneFinalInput(obj, fc006Loc[0])
					}
					//if len(yiLiaoJiGouItems) < 2 {
					//	continue
					//}
					//dengJi := yiLiaoJiGouItems[1]
					//bianHao := yiLiaoJiGouItems[2]
					//1、机构号为6542开头时，匹配医疗机构的医院编号为6542开头或未匹配医院录入内容包含“塔城”字样，则匹配为对应医疗目录“床位费（塔城地区）”，否则匹配为对应医疗目录“床位费（非塔城地区）”；
					if strings.HasPrefix(bill.Agency, "6542") {
						//if strings.HasPrefix(bianHao, "6542") || (!ok3 && strings.Index(oneVal, "塔城") != -1) {
						if strings.HasPrefix(bianHao, "6542") || (total1 == 0 && strings.Index(oneVal, "塔城") != -1) {
							xiangMuMingCheng = "床位费（塔城地区）"
						} else {
							xiangMuMingCheng = "床位费（非塔城地区）"
						}
					}
					//2、机构号为6532开头时，匹配医疗机构的医院编号为6532开头或未匹配医院录入内容包含“和田”字样，则匹配为对应医疗目录“床位费（和田地区）”，否则匹配为对应医疗目录“床位费（非和田地区）”；
					if strings.HasPrefix(bill.Agency, "6532") {
						//if strings.HasPrefix(bianHao, "6532") || (!ok3 && strings.Index(oneVal, "和田") != -1) {
						if strings.HasPrefix(bianHao, "6532") || (total1 == 0 && strings.Index(oneVal, "和田") != -1) {
							xiangMuMingCheng = "床位费（和田地区）"
						} else {
							xiangMuMingCheng = "床位费（非和田地区）"
						}
					}
					//3、机构号为6529、6543、6531、6502、6530、6540开头时，匹配为对应医疗目录的“床位费”；
					if utils.RegIsMatch("^(6529|6543|6531|6502|6530|6540)", bill.Agency) {
						xiangMuMingCheng = "床位费"
					}
					//4、机构不为6542、6532或6529、6543、6531、6502、6530、6540开头时
					//a.医院等级为“01、02、03、04”时，匹配对应医疗目录的“床位费（三级）”；
					//b.医院等级为“05、06、07”时，匹配对应医疗目录的“床位费（二级）”；
					//c.医院等级为“08、09、10、11、12、13”时，匹配对应医疗目录的“床位费（一级）”；
					//d.无法匹配医院等级时，匹配对应医疗目录的“床位费（一级）”；
					if !utils.RegIsMatch("^(6542|6532|6529|6543|6531|6502|6530|6540)", bill.Agency) {
						xiangMuMingCheng = "床位费（一级）"
						if utils.RegIsMatch("^(01|02|03|04)$", dengJi) {
							xiangMuMingCheng = "床位费（三级）"
						}
						if utils.RegIsMatch("^(05|06|07)$", dengJi) {
							xiangMuMingCheng = "床位费（二级）"
						}
						if utils.RegIsMatch("^(08|09|10|11|12|13)$", dengJi) {
							xiangMuMingCheng = "床位费（一级）"
						}
					}
				}

				xiangMu, total2 := utils.FetchConst(bill.ProCode, yiLiaoMuLu, "医疗项目编码", map[string]string{"项目名称": xiangMuMingCheng})
				biLi, _ := utils.FetchConst(bill.ProCode, yiLiaoMuLu, "自费比例", map[string]string{"项目名称": xiangMuMingCheng})
				//yiLiaoMuLuItems, ok1 := constSpecialMap[yiLiaoMuLu]
				//if !ok1 {
				//	//未匹配要收集
				//	utils.SetOnlyOneFinalInput(obj, oneLoc)
				//	global.GLog.Error("yiLiaoMuLu", zap.Any("没有改医疗目录表", yiLiaoMuLu))
				//	continue
				//}
				//item, ok2 := yiLiaoMuLuItems[xiangMuMingCheng]
				//CSB0110RC0045000 项目比例默认1
				//utils.SetOnlyOneFinalValue(obj, twoLoc[0], "1")
				//if ok2 {
				if total2 != 0 {
					//CSB0110RC0043000 匹配目录就直接带出项目编码
					//utils.SetOnlyOneFinalValue(obj, oneLoc, item[1])
					utils.SetOnlyOneFinalValue(obj, oneLoc, xiangMu)

					//CSB0110RC0044000 匹配目录就直接带出项目比例 （左边字段为空时，不执行该校验）
					if oneVal != "" {
						//utils.SetOnlyOneFinalValue(obj, twoLoc[0], item[7])
						utils.SetOnlyOneFinalValue(obj, twoLoc[0], biLi)
					}
				} else {
					//未匹配要收集
					utils.SetOnlyOneFinalInput(obj, oneLoc)
				}
			}
		}

	}

	//CSB0110RC0044000
	//"下列第一列字段(循环分块)录入内容包含“床位费”时，根据同一发票属性下的fc006的录入内容匹配常量《B0110_新疆国寿理赔_医疗机构65》的“医院名称”（第一列）检索出对应的“医院等级”（第二列），根据医院等级与对应医疗目录进行匹配，匹配规则为：
	//1、机构号为6542开头时，匹配医疗机构的医院编号为6542开头或未匹配医院录入内容包含“塔城”字样，则匹配为对应医疗目录“床位费（塔城地区）”，否则匹配为对应医疗目录“床位费（非塔城地区）”；
	//2、机构号为6532开头时，匹配医疗机构的医院编号为6532开头或未匹配医院录入内容包含“和田”字样，则匹配为对应医疗目录“床位费（和田地区）”，否则匹配为对应医疗目录“床位费（非和田地区）”；
	//3、机构号为6529、6543、6531、6502、6530、6540开头时，匹配为对应医疗目录的“床位费”；
	//4、机构不为6542、6532或6529、6543、6531、6502、6530、6540开头时
	//a.医院等级为“01、02、03、04”时，匹配对应医疗目录的“床位费（三级）”；
	//b.医院等级为“05、06、07”时，匹配对应医疗目录的“床位费（二级）”；
	//c.医院等级为“08、09、10、11、12、13”时，匹配对应医疗目录的“床位费（一级）”；
	//d.无法匹配医院等级时，匹配对应医疗目录的“床位费（一级）”；
	//以上所匹配医疗目录的校验需要能兼容查询全角和半角的括号；
	//匹配后检索对应的“医疗项目编码”（第二列）和“自费比例”（第八列），将录入内容包含“床位费”的对应第一列字段(循环分块)结果值默认为对应的“医疗项目编码”、第二列字段(循环分块)结果值默认为对应的“自费比例”；
	//项目名称,自费比
	//fc012,fc053
	//fc025,fc054
	//fc026,fc055
	//fc027,fc056
	//fc028,fc057
	//fc029,fc058
	//fc030,fc059
	//fc031,fc060"
	//CSB0110RC0045000
	//"下列字段中，右边字段的结果数据根据左边字段录入内容与对应医疗目录中“项目名称”（第三列）对应匹配的“自费比例”（第八列）进行输出，默认为1
	//项目名称 自费比

	//CSB0110RC0046000
	//fc005录入内容不能与之前所有案件的所有fc005重复，若重复，出问题件：发票号与历史票据重复，请核实；
	for _, fc005Loc := range fieldLocationMap["fc005"] {
		fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc, false)
		invoiceNums := service.FetchInvoiceNum(bill.ProCode, bill.BillNum, fc005Val)
		if len(invoiceNums) < 1 {
			continue
		}
		global.GLog.Info("发票号与历史票据重复", zap.Any("", invoiceNums[0].BillNum))
		utils.SetIssue(obj, fc005Loc, "发票号与历史票据重复，请核实", "", "")
	}

	//CSB0110RC0047000
	//fc006结果值根据《B0110_新疆国寿理赔_医疗机构65》的“医院名称”（第一列）转成对应的“医院编码”（第三列）输出，如不能转码默认为650100000022，并出问题件：医院为数据库以外
	for _, loc := range fieldLocationMap["fc006"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		v, total := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_医疗机构65", "医院编码", map[string]string{"医院名称": val})
		//yiLiaoJiGouItems, ok := constSpecialMap["医疗机构65"][val]
		//if ok && len(yiLiaoJiGouItems) > 2 {
		if total != 0 {
			//utils.SetOnlyOneFinalValue(obj, loc, yiLiaoJiGouItems[2])
			utils.SetOnlyOneFinalValue(obj, loc, v)
		} else {
			//未匹配要收集
			utils.SetOnlyOneFinalInput(obj, loc)
			utils.SetOnlyOneFinalValue(obj, loc, "650100000022")
			utils.SetIssue(obj, loc, "医院为数据库以外", "", "")
		}
	}

	//CSB0110RC0048000
	//fc006录入内容为空、为A、为B时，fc006默认为6500000000000349，出问题件：医院内容模糊无法辨识录入
	for _, loc := range fieldLocationMap["fc006"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		if utils.RegIsMatch("^(|A|B)$", val) {
			utils.SetOnlyOneFinalValue(obj, loc, "6500000000000349")
			utils.SetIssue(obj, loc, "医院内容模糊无法辨识录入", "", "")
		}
	}

	//CSB0110RC0049000
	//fc007结果值根据《B0110_新疆国寿理赔_ICD10疾病编码》的“疾病名称”（第一列）转成对应的“疾病代码”（第二列）输出，如不能转码默认为R52.9，并出问题件：疾病诊断无法辨识
	for _, loc := range fieldLocationMap["fc007"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		v, total := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_ICD10疾病编码", "疾病代码", map[string]string{"疾病名称": val})
		//v, ok := constMap["jiBingBianMaMap"][val]
		utils.SetOnlyOneFinalValue(obj, loc, v)
		//if !ok {
		if total == 0 {
			utils.SetOnlyOneFinalValue(obj, loc, "R52.9")
			utils.SetIssue(obj, loc, "疾病诊断无法辨识", "", "")
		}
	}

	//CSB0110RC0083000
	//同一发票属性下，fc069的结果值默认为fc005的结果值	新增需求
	for _, loc := range fieldLocationMap["fc005"] {
		val := utils.GetFieldValueByLoc(obj, loc, true)
		fc069Loc := utils.GetFieldLoc(fieldLocationMap["fc069"], loc[0], -1, -1, -1)
		if len(fc069Loc) < 1 {
			continue
		}
		utils.SetFinalValue(obj, fc069Loc, val)
	}

	//CSB0110RC0084000
	//同一发票属性下，当fc067录入2时，fc072的结果值默认为fc011的结果值	新增需求
	// for _, loc := range fieldLocationMap["fc067"] {
	// 	val := utils.GetFieldValueByLoc(obj, loc, true)
	// 	fc072Loc := utils.GetFieldLoc(fieldLocationMap["fc072"], loc[0], -1, -1, -1)
	// 	fc011Loc := utils.GetFieldLoc(fieldLocationMap["fc011"], loc[0], -1, -1, -1)
	// 	if len(fc072Loc) < 1 && len(fc011Loc) < 1 {
	// 		continue
	// 	}
	// 	fc011Val := utils.GetFieldValueByLoc(obj, fc011Loc[0], false)
	// 	if val == "2" {
	// 		utils.SetFinalValue(obj, fc072Loc, fc011Val)
	// 	}
	// }

	//CSB0110RC0085000
	//同一发票属性下，当fc066录入2时，清空fc067、fc068、fc069、fc070、fc071、fc072、fc073字段的结果值	新增需求
	codes = []string{"fc067", "fc068", "fc069", "fc070", "fc071", "fc072", "fc073"}
	for _, loc := range fieldLocationMap["fc066"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		for _, code := range codes {
			loc1 := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
			if len(loc1) < 1 {
				continue
			}
			if val == "2" {
				utils.SetFinalValue(obj, loc1, "")
			}
		}
	}

	//CSB0110RC0086000
	//同一发票属性下的清单分块MB002-bc002，根据fc008清单页码从小到大的顺序进行排序	新增需求
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		sort.Slice(invoiceMap.QingDan, func(i, j int) bool {
			return utils.GetDecimalValByFields(invoiceMap.QingDan[i], "fc008", false).
				LessThan(utils.GetDecimalValByFields(invoiceMap.QingDan[j], "fc008", false))
		})
	}

	//收集发票号
	var invoiceNums []model2.InvoiceNum
	//var fc005ValArr []string
	//for _, fc005Loc := range fieldLocationMap["fc005"] {
	//	var invoiceNum model2.InvoiceNum
	//	fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc, true)
	//	if utils.HasItem(fc005ValArr, fc005Val) {
	//		utils.SetIssue(obj, fc005Loc, "发票号["+fc005Val+"]重复，请核实", "", "")
	//	}
	//	invoiceNum.Num = fc005Val
	//	invoiceNum.BillNum = bill.BillNum
	//	invoiceNums = append(invoiceNums, invoiceNum)
	//	fc005ValArr = append(fc005ValArr, fc005Val)
	//}
	err = service.UpdateInvoiceNumByBillNum(bill.ProCode, bill.BillNum, invoiceNums)

	mes := ""
	isNotCheck := true
	for ii, invoice := range obj.Invoice {
		fc073F := ""
		fc067 := ""
		fc068 := ""
		fc069 := ""
		fc070 := ""
		fc071 := ""
		fc072 := ""
		fc005 := ""
		isChange := false
		for jj, fields := range invoice.Invoice {
			for kk, field := range fields {
				if field.Code == "fc073" {
					fc073F = field.FinalValue
				}
				if field.Code == "fc005" {
					fc005 = field.FinalValue
				}
				if field.Code == "fc067" {
					fc067 = field.ResultValue
				}
				if field.Code == "fc068" {
					fc068 = field.FinalValue
				}
				if field.Code == "fc069" {
					fc069 = field.FinalValue
				}
				if field.Code == "fc070" {
					fc070 = field.FinalValue
				}
				if field.Code == "fc071" {
					fc071 = field.FinalValue
				}

				if "fc072" == field.Code {
					_, fc067 := GetOneField(fields, "fc067", false)
					if fc067 == "2" || fc067 == "3" {
						_, fc011 := GetOneField(fields, "fc011", true)
						obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc011
					}
				}
				if field.Code == "fc072" {
					fc072 = obj.Invoice[ii].Invoice[jj][kk].FinalValue
				}
				if RegIsMatch(field.Code, `^(fc067|fc068|fc069|fc070|fc071|fc072)$`) && field.IsChange {
					isChange = true
				}

				//CSB0110RC0117000
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
		fmt.Println("------------bill.Stage-------------------:", bill.Stage, isChange)
		if fc073F == "" && fc067 != "" {
			if bill.Stage == 6 || (bill.Stage != 6 && isChange && strings.Contains(bill.Remark, "<<发票查验>>")) {
				bodyData := make(map[string]interface{})
				bodyData["fpdm"] = fc068
				bodyData["fphm"] = fc069
				bodyData["kprq"] = fc070
				bodyData["checkCode"] = fc071
				fmt.Println("------------bodyData-------------------:", invoice.Code, bodyData)
				fc073 := ""
				if fc067 == "1" {
					bodyData["noTaxAmount"] = fc072
					err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData)
					isNotCheck = false
					if err != nil {
						mes += fc005 + fmt.Sprintf("%v", err) + ";"
						fc073 = "99"
						// response.FailWithMessage(fmt.Sprintf("%v", err), c)
						// return
					} else {
						// mes += fc005 + "发票查验成功;"
						del := respData.Data.Del
						if del == "0" {
							fc073 = "01"
						}
						if del == "2" {
							fc073 = "99"
						}
						if del == "3" || del == "7" || del == "8" {
							fc073 = "02"
						}
					}

				} else if fc067 == "2" {
					bodyData["money"] = fc072
					err, respData := unitFunc.Invoice("/v2/eInvoice/query", bodyData)
					isNotCheck = false
					if err != nil {
						mes += "查验失败：" + fmt.Sprintf("%v", err) + ";"
						// if RegIsMatch(fmt.Sprintf("%v", err), `(800|801|1098),`) {
						// fmt.Println("----------------------------------------------------------errerrerrerr-------------------:", err)
						fc073 = "99"
						// }
						// response.FailWithMessage(fmt.Sprintf("%v", err), c)
						// return
					} else {
						// mes += fc005 + "发票查验成功;"
						if respData.Data.IsRed {
							fc073 = "02"
						} else {
							if respData.Data.IsPrint {
								fc073 = "03"
							} else {
								fc073 = "01"
							}
						}
					}

				} else if fc067 == "3" {
					bodyData1 := make(map[string]interface{})
					bodyData1["fphm"] = fc069
					bodyData1["kprq"] = fc070
					bodyData1["jshj"] = fc072
					err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData1)
					isNotCheck = false
					if err != nil {
						mes += fc005 + fmt.Sprintf("%v", err) + ";"
						fc073 = "99"
						// response.FailWithMessage(fmt.Sprintf("%v", err), c)
						// return
					} else {
						// mes += fc005 + "发票查验成功;"
						del := respData.Data.Del
						if del == "0" {
							fc073 = "01"
						}
						if del == "2" {
							fc073 = "99"
						}
						if del == "3" || del == "7" || del == "8" {
							fc073 = "02"
						}
					}
				}

				fmt.Println("------------fc073-------------------:", fc073)
				for jj, fields := range invoice.Invoice {
					for kk, field := range fields {
						if field.Code == "fc073" {
							obj.Invoice[ii].Invoice[jj][kk].ResultValue = fc073
							obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc073
							obj.Invoice[ii].Invoice[jj][kk].IsChange = true
						}
					}
				}
			}
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

	// 编码CSB0110RC0091000 导出校验 最后一页字段录入存在?？ 就导出校验  "最后一页存在?，请检查;" 20230713新增
	var fieldArr []string
	var wrong []string
	for p, invoiceMap := range obj.Invoice {
		fc005 := ""
		for _, fields := range invoiceMap.Invoice {
			if p == len(obj.Invoice)-1 {
				if strings.Index(obj.Bill.WrongNote, "最后一页存在?，请检查;") == -1 {
					for _, field := range fields {
						if RegIsMatch(field.ResultValue, `[?？]`) {
							wrong = append(wrong, "最后一页存在?，请检查;")
						}
					}
				}
			}
			if p != len(obj.Invoice)-1 {
				for _, field := range fields {
					if field.Code == "fc005" {
						fc005 = field.FinalValue
					}
					if RegIsMatch(field.ResultValue, `[?？]`) {
						fieldArr = append(fieldArr, field.Name)
					}
				}
			}
		}
		for _, fields := range invoiceMap.QingDan {
			for _, field := range fields {
				if RegIsMatch(field.ResultValue, `[?？]`) {
					fieldArr = append(fieldArr, field.Name)
				}
			}

		}
		for _, fields := range invoiceMap.BaoXiaoDan {
			for _, field := range fields {
				if RegIsMatch(field.ResultValue, `[?？]`) {
					fieldArr = append(fieldArr, field.Name)
				}
			}
		}
		for _, fields := range invoiceMap.ThirdBaoXiaoDan1 {
			for _, field := range fields {
				if RegIsMatch(field.ResultValue, `[?？]`) {
					fieldArr = append(fieldArr, field.Name)
				}
			}
		}
		for _, fields := range invoiceMap.ThirdBaoXiaoDan2 {
			for _, field := range fields {
				if RegIsMatch(field.ResultValue, `[?？]`) {
					fieldArr = append(fieldArr, field.Name)
				}
			}
		}
		for _, fields := range invoiceMap.ThirdBaoXiaoDan3 {
			for _, field := range fields {
				if RegIsMatch(field.ResultValue, `[?？]`) {
					fieldArr = append(fieldArr, field.Name)
				}
			}
		}
		//需求编码 CSB0110RC0091000
		//1.校验所有字段的录入值，当包含?或？时，出导出校验：发票【xxx】的【yyy】存在?号，请核实；
		//2.如同一发票下多个字段包含?或？，则导出校验提示格式为：发票【xxx】的【yyy1、yyy2】存在?号，请核实；（xxx为发票号fc005的值，yyy为包含问号的字段名）
		//3.当最后一页字段包含?或？时，则直接出一条导出校验：最后一页存在?，请检查；
		for _, loc := range fieldLocationMap["fc005"] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			fmt.Println("---------------------------------------AAAAAAAAAAAA")
			fmt.Println("val=", val)
			fmt.Println("---------------------------------------AAAAAAAAAAAA")
			if len(fieldArr) > 0 {
				joinFiled := strings.Join(fieldArr, "、")
				mes := "发票" + fc005 + "的【" + joinFiled + "】存在?号，请核实;"
				obj.Bill.WrongNote += mes
				fieldArr = make([]string, 0)
			}
		}
	}
	if len(wrong) > 1 {
		obj.Bill.WrongNote += "最后一页存在?，请检查;"
	}

	for _, loc := range fieldLocationMap["fc074"] {
		val := utils.GetFieldValueByLoc(obj, loc, true)
		if val != "" {
			fc074s := utils.RegMatchAll(val, `[^;；]+`)
			for _, fc074 := range fc074s {
				utils.SetIssue(obj, loc, fc074, "", "")
			}
		}
		// utils.SetOnlyOneFinalValue(obj, loc, "R52.9")

	}

	////////////////////////项目报表 - 业务明细 - 单据状态
	obj.Bill.BillType = utils2.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	return err, obj
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"diSanFangDanWeiMap", "B0110_新疆国寿理赔_第三方出具单位", "2", "1"},
		{"jiBingBianMaMap", "B0110_新疆国寿理赔_ICD10疾病编码", "0", "1"},
		{"jiBingBianMaCodeToNameMap", "B0110_新疆国寿理赔_ICD10疾病编码", "1", "0"},
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
		{"shuJuKuBianMa", "B0110_新疆国寿理赔_数据库编码对应表", "1"},
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
			constObj[strings.Replace(k, "B0110_新疆国寿理赔_", "", -1)] = tempMap
		}
		if strings.Index(k, "医疗机构") != -1 {
			tempMap := make(map[string][]string, 0)
			for _, arr := range items {
				tempMap[arr[0]] = arr
			}
			constObj[strings.Replace(k, "B0110_新疆国寿理赔_", "", -1)] = tempMap
		}
	}
	return constObj
}
