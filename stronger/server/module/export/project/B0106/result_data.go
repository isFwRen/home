/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年05月08日20:44:29
 */

package B0106

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
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

//B0106陕西国寿理赔
//模板类型字段 : fc003
//账单号: fc005票据号
//总金额: fc011发票总金额
//发票类型: fc009诊疗方式
//
//和发票有对应关系的
//fc018发票属性 MB002-bc001
//fc019清单所属发票 MB002-bc002
//fc020报销单所属发票 MB002-bc003
//fc089第三方支付所属发票 MB002-bc005
//
//和发票没有对应关系的
//4-诊断书

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:          model2.TypeCode{FieldCode: []string{"fc018"}, BlockCode: []string{"bc001"}},
	QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc019"}, BlockCode: []string{"bc002"}},
	BaoXiaoDanFieldCode:       model2.TypeCode{FieldCode: []string{"fc020"}, BlockCode: []string{"bc003"}},
	ThirdBaoXiaoDan1FieldCode: model2.TypeCode{FieldCode: []string{"fc089"}, BlockCode: []string{"bc005"}},
	OtherTempType:             map[string]string{"4": "4"},
	TempTypeField:             "fc003",
	InvoiceNumField:           []string{"fc005"},
	MoneyField:                []string{"fc011"},
	InvoiceTypeField:          "fc009",
}

// ResultData B0106
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

	//CSB0106RC0035000	"同一发票下，
	//fc022录入内容为A、B、?，fc023录入内容为正常日期时，将fc023的结果值赋值给fc022结果值；
	//fc023录入内容为A、B、?，fc022录入内容为正常日期时，将fc022的结果值赋值给fc023结果值；"
	myCode := [][]string{
		{"fc022", "fc023"},
		{"fc023", "fc022"},
	}
	for _, code := range myCode {
		for _, loc := range fieldLocationMap[code[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			loc1 := utils.GetFieldLoc(fieldLocationMap[code[1]], loc[0], loc[1], loc[2], loc[4])
			if len(loc1) < 1 {
				continue
			}
			val1 := utils.GetFieldValueByLoc(obj, loc1[0], true)
			_, err := time.Parse("2006-01-02", val1)
			if utils.RegIsMatch("^(A|B|\\?)$", val) && err == nil {
				utils.SetOnlyOneFinalValue(obj, loc, val1)
			}
		}
	}

	//CSB0106RC0036000
	//"以下左边字段为报销单的字段，右边是报销单对应的发票的字段，当fc018与fc020相匹配时，需要将对应的左边的字段的值，赋值到右边的字段中
	myCode = [][]string{
		{"fc046", "fc050"},
		{"fc047", "fc052"},
		{"fc048", "fc053"},
		{"fc049", "fc054"},
	}
	for _, codes := range myCode {
		for _, loc := range fieldLocationMap[codes[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, true)
			loc1 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc[0], -1, -1, -1)
			if len(loc1) < 1 {
				continue
			}
			utils.SetFinalValue(obj, loc1, val)
		}
	}

	//CSB0106RC0037000
	//fc007、fc008结果值根据《B0106_陕西国寿理赔_ICD10疾病编码》中“疾病名称”（第一列）转成对应的“疾病代码”（第二列）输出；
	//CSB0106RC0044000
	//"1、fc007录入与常量库《B0106_陕西国寿理赔_ICD10疾病编码》不匹配时，fc007默认为X59，出问题件：疾病诊断无法辨识；
	//2、fc008录入内容与常量库《B0106_陕西国寿理赔_ICD10疾病编码》不匹配时（为空、为A时不进行校验），fc008默认为X59，出问题件：疾病诊断无法辨识；"
	codeArr := []string{"fc007", "fc008"}
	for i, code := range codeArr {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			//v, ok := constMap["jiBingBianMaMap"][val]
			v, total := utils.FetchConst(bill.ProCode, "B0106_陕西国寿理赔_ICD10疾病编码", "疾病代码", map[string]string{"疾病名称": val})
			utils.SetOnlyOneFinalValue(obj, loc, v)
			//if !ok {
			if total == 0 {
				if i == 1 && !utils.RegIsMatch("^(|A)$", val) {
					utils.SetOnlyOneFinalValue(obj, loc, "X59")
					utils.SetIssue(obj, loc, "疾病诊断无法辨识", "", "")
				} else if i == 0 {
					utils.SetOnlyOneFinalValue(obj, loc, "X59")
					utils.SetIssue(obj, loc, "疾病诊断无法辨识", "", "")
				}
			}
		}
	}

	//CSB0106RC0038000
	//fc087的结果值根据《B0106_陕西国寿理赔_第三方出具单位》中“第三方出具单位名称”（第三列）对应的“第三方出具单位代码”（第二列）进行转码输出
	for _, loc := range fieldLocationMap["fc087"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		v, _ := utils.FetchConst(bill.ProCode, "B0106_陕西国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": val})
		//utils.SetOnlyOneFinalValue(obj, loc, constMap["diSanFangDanWeiMap"][val])
		utils.SetOnlyOneFinalValue(obj, loc, v)
	}

	//CSB0106RC0039000
	//在同一个发票属性存在多个fc087及fc088的结果值时，将多个fc087、fc088的结果值用“|”符号整合，放到最后一个fc087、fc088的结果值中
	//2023年06月09日17:04:07 改成放在第一个
	codes := []string{"fc087", "fc088"}
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

	//CSB0106RC0040000
	//"当fc005录入内容A或B或?时，结果值按顺序默认为999991~n，并出问题件：发票号无法辨识
	//如共有12个fc005为A或B或?时，按顺序默认，第1个fc005为999991，第9个fc005为999999，第10个fc005为9999910，第11个fc005为9999911，第12个fc005为9999912，以此类推……"
	for i, fc005Loc := range fieldLocationMap["fc005"] {
		fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc, false)
		//nextID := strconv.FormatInt(utils.GWorker.NextId(), 10)
		if utils.RegIsMatch("^(A|B|\\?)$", fc005Val) {
			utils.SetOnlyOneFinalValue(obj, fc005Loc, fmt.Sprintf("999999%v", i+1))
			utils.SetIssue(obj, fc005Loc, "发票号无法辨识", "", "")
		}
	}

	//CSB0106RC0041000
	//fc005录入内容不能与之前所有案件的所有fc005重复，若重复，出问题件：发票号与案件xxx中发票重复,请确认（xxx为发票重复的案件单号，fc005为空或前五位数是99999时，不需要报重复的问题件）
	for _, fc005Loc := range fieldLocationMap["fc005"] {
		fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc, false)
		invoiceNums := service.FetchInvoiceNum(bill.ProCode, bill.BillNum, fc005Val)
		if len(invoiceNums) < 1 || fc005Val == "" || utils.RegIsMatch("^999999", fc005Val) {
			continue
		}
		global.GLog.Info("发票号与历史票据重复", zap.Any("", invoiceNums[0].BillNum))
		utils.SetIssue(obj, fc005Loc, "发票号与历史票据重复，请核实", "", "")
	}

	//CSB0106RC0042000
	//fc015为清单循环分块bc002中的字段，fc015字段的结果值默认为该分块对应的发票的票据号fc005的结果值
	for _, loc := range fieldLocationMap["fc015"] {
		fc005Loc := utils.GetFieldLoc(fieldLocationMap["fc005"], loc[0], -1, -1, -1)
		if len(fc005Loc) < 1 {
			continue
		}
		fc005Val := utils.GetFieldValueByLoc(obj, fc005Loc[0], true)
		utils.SetOnlyOneFinalValue(obj, loc, fc005Val)
	}

	//CSB0106RC0043000
	//"fc006根据《B0106_陕西国寿理赔_医疗机构61》中“医院名称”（第一列）转成对应的“医院编码”（第二列）输出。
	//1、机构号是6101开头时，fc006无法转换时默认为“西安市社区卫生服务中心”的“医院编码”输出，不出问题件。
	//2、其他机构号，fc006无法转换时默认为610100000012345，并出问题件：医院为数据库以外。"
	for _, loc := range fieldLocationMap["fc006"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		//yiLiaoJiGouItems, ok := constSpecialMap["医疗机构61"][val]
		v, total := utils.FetchConst(bill.ProCode, "B0106_陕西国寿理赔_医疗机构61", "医院编码", map[string]string{"医院名称": val})
		//if ok && len(yiLiaoJiGouItems) > 1 {
		if total > 0 {
			//utils.SetOnlyOneFinalValue(obj, loc, yiLiaoJiGouItems[1])
			utils.SetOnlyOneFinalValue(obj, loc, v)
		} else if utils.RegIsMatch("^6101", bill.Agency) && total == 0 {
			v, _ = utils.FetchConst(bill.ProCode, "B0106_陕西国寿理赔_医疗机构61", "医院编码", map[string]string{"医院名称": "西安市社区卫生服务中心"})
			//yiLiaoJiGouItemsTemp, _ := constSpecialMap["医疗机构61"]["西安市社区卫生服务中心"]
			//utils.SetOnlyOneFinalValue(obj, loc, yiLiaoJiGouItemsTemp[1])
			utils.SetOnlyOneFinalValue(obj, loc, v)
		} else {
			utils.SetOnlyOneFinalValue(obj, loc, "610100000012345")
			//未匹配要收集
			utils.SetOnlyOneFinalInput(obj, loc)
			utils.SetIssue(obj, loc, "医院为数据库以外", "", "")
		}
	}

	//CSB0106RC0045000
	//"fc012、fc025、fc026、fc027、fc028、fc029、fc030、fc031根据录入内容匹配对应常量表进行转码输出，当无法转码时不作处理；
	//1.匹配常量表《B0106_陕西国寿理赔_数据库编码对应表》中“本级机构代码”（第二列）对应的“医疗目录编码”（第四列）的常量表中的“项目名称”（第三列）对应的“医疗项目编码”（第二列）进行转码输出；
	//2.当下载文件中的bpoSendRemark节点内容包含“生育”时，直接匹配常量表《B0106_陕西国寿理赔_医疗目录M6101002022001》中的“项目名称”（第三列）对应的“医疗项目编码”（第二列）进行转码输出；"
	//CSB0106RC0048000	"下列字段中，右边字段的结果数据根据左边字段录入内容与对应医疗目录中“项目名称”（第三列）对应匹配的“自费比例”（第八列）进行输出,如未匹配到自费比例时,默认为1
	//项目名称 自费比
	//fc012 fc055
	//fc025 fc059
	//fc026 fc063
	//fc027 fc067
	//fc028 fc071
	//fc029 fc075
	//fc030 fc079
	//fc031 fc083"
	v, total := utils.FetchConst(bill.ProCode, "B0106_陕西国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	//shuJuKuBianMaItems, ok := constSpecialMap["shuJuKuBianMa"][bill.Agency]
	//if ok && len(shuJuKuBianMaItems) > 3 {
	if total > 0 {
		yiLiaoMuLu := "B0106_陕西国寿理赔_医疗目录" + v
		otherInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(bill.OtherInfo), &otherInfo)
		if err != nil {
			global.GLog.Error("CSB0106RC0045000", zap.Error(err))
		}
		if strings.Index(otherInfo["bpoSendRemark"].(string), "生育") != -1 {
			yiLiaoMuLu = "医疗目录M6101002022001"
		}
		myCode = [][]string{
			{"fc012", "fc055"},
			{"fc025", "fc059"},
			{"fc026", "fc063"},
			{"fc027", "fc067"},
			{"fc028", "fc071"},
			{"fc029", "fc075"},
			{"fc030", "fc079"},
			{"fc031", "fc083"},
		}
		for _, codes := range myCode {
			for _, loc := range fieldLocationMap[codes[0]] {
				loc1 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc[0], loc[1], loc[2], loc[4])
				val := utils.GetFieldValueByLoc(obj, loc, false)
				v1, total1 := utils.FetchConst(bill.ProCode, yiLiaoMuLu, "医疗项目编码", map[string]string{"项目名称": val})
				v2, _ := utils.FetchConst(bill.ProCode, yiLiaoMuLu, "自费比例", map[string]string{"项目名称": val})
				//yiLiaoMuLuConst, ok1 := constSpecialMap[yiLiaoMuLu]
				//if !ok1 || len(yiLiaoMuLuConst[val]) < 8 || len(loc1) != 1 {
				if total1 == 0 || len(loc1) != 1 {
					//未匹配要收集
					utils.SetOnlyOneFinalInput(obj, loc)
					utils.SetFinalValue(obj, loc1, "1")
					global.GLog.Error("不存在该常量、匹配的值的len()不够或者左右字段不对呀", zap.Any(yiLiaoMuLu, v1))
					continue
				}
				//utils.SetOnlyOneFinalValue(obj, loc, yiLiaoMuLuConst[val][1])
				//utils.SetFinalValue(obj, loc1, yiLiaoMuLuConst[val][7])
				utils.SetOnlyOneFinalValue(obj, loc, v1)
				utils.SetFinalValue(obj, loc1, v2)
			}
		}
	}

	//CSB0106RC0046000
	//fc013、fc032、fc033、fc034、fc035、fc036、fc037、fc038，当以上字段的录入值包含“?”或“？”时，结果值默认为“1”。
	codes = []string{"fc013", "fc032", "fc033", "fc034", "fc035", "fc036", "fc037", "fc038"}
	for _, code := range codes {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			if strings.Index(val, "？") != -1 || strings.Index(val, "?") != -1 {
				utils.SetOnlyOneFinalValue(obj, loc, "1")
			}
		}
	}

	//CSB0106RC0047000
	//"在同一发票中(根据fc018发票属性、fc019清单所属属性判断是否属于同一张发票)对下列分组中左边字段录入内容进行校验，
	//当左边字段的录入值不为1（为?时不进行此校验）时，将右边字段的结果值替换成“右边字段录入值除以左边字段录入值的值”；
	//数量 金额
	myCode = [][]string{
		{"fc013", "fc014"},
		{"fc032", "fc039"},
		{"fc033", "fc040"},
		{"fc034", "fc041"},
		{"fc035", "fc042"},
		{"fc036", "fc043"},
		{"fc037", "fc044"},
		{"fc038", "fc045"},
	}
	for _, code := range myCode {
		for _, loc := range fieldLocationMap[code[0]] {
			loc1 := utils.GetFieldLoc(fieldLocationMap[code[1]], loc[0], loc[1], loc[2], loc[4])
			if len(loc1) < 1 {
				continue
			}
			decimalVal1 := utils.GetFieldDecimalValueByLoc(obj, loc1[0], false)
			decimalVal := utils.GetFieldDecimalValueByLoc(obj, loc, false)
			if !decimalVal.IsZero() && !decimalVal1.Equals(decimal.NewFromInt(1)) {
				utils.SetFinalValue(obj, loc1, decimalVal1.Div(decimalVal).StringFixed(2))
			}
		}
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

	// 编码CSB0106RC0085000 导出校验 最后一页字段录入存在?？ 就导出校验  "最后一页存在?，请检查;" 20230713新增
	var wrong []string
	for p, invoiceMap := range obj.Invoice {
		if p == len(obj.Invoice)-1 {
			for _, fields := range invoiceMap.Invoice {
				if strings.Index(obj.Bill.WrongNote, "最后一页存在?，请检查；") == -1 {
					for _, field := range fields {
						if RegIsMatch(field.ResultValue, `[?？]`) {
							wrong = append(wrong, field.Name)
						}
					}
				}
			}
		}
	}

	mes := ""
	isNotCheck := true
	for ii, invoice := range obj.Invoice {
		fc091 := ""
		fc092 := ""
		fc093 := ""
		fc094 := ""
		fc095 := ""
		fc096 := ""
		fc005 := ""
		fc097F := ""
		isChange := false
		for jj, fields := range invoice.Invoice {
			for kk, field := range fields {

				if "fc096" == field.Code {
					_, fc091 := GetOneField(fields, "fc091", false)
					if fc091 == "2" || fc091 == "3" {
						_, fc011 := GetOneField(fields, "fc011", true)
						obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc011
					}
				}

				if "fc093" == field.Code {
					_, fc005 := GetOneField(fields, "fc005", true)
					obj.Invoice[ii].Invoice[jj][kk].FinalValue = fc005
				}
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
				if field.Code == "fc093" {
					fc093 = obj.Invoice[ii].Invoice[jj][kk].FinalValue
				}
				if field.Code == "fc094" {
					fc094 = field.FinalValue
				}
				if field.Code == "fc095" {
					fc095 = field.FinalValue
				}
				if field.Code == "fc096" {
					fc096 = obj.Invoice[ii].Invoice[jj][kk].FinalValue
				}
				if RegIsMatch(field.Code, `^(fc091|fc092|fc093|fc094|fc095|fc096)$`) && field.IsChange {
					isChange = true
				}

				//CSB0106RC0119000
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
		if fc097F == "" && fc091 != "" {
			if bill.Stage == 6 || (bill.Stage != 6 && isChange && strings.Contains(bill.Remark, "<<发票查验>>")) {
				bodyData := make(map[string]interface{})
				bodyData["fpdm"] = fc092
				bodyData["fphm"] = fc093
				bodyData["kprq"] = fc094
				bodyData["checkCode"] = fc095
				fmt.Println("------------bodyData-------------------:", invoice.Code, bodyData)
				fc097 := ""
				if fc091 == "1" {
					bodyData["noTaxAmount"] = fc096
					fmt.Println("------------发票查验开始-------------------:", time.Now(), invoice.Code, bodyData)
					err, respData := unitFunc.Invoice("/v2/invoice/query", bodyData)
					fmt.Println("------------发票查验结束-------------------:", time.Now())
					fmt.Println("------------发票查验信息-------------------:", respData)
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
					fmt.Println("------------发票查验开始-------------------:", time.Now(), invoice.Code, bodyData)
					err, respData := unitFunc.Invoice("/v2/eInvoice/query", bodyData)
					fmt.Println("------------发票查验结束-------------------:", time.Now())
					fmt.Println("------------发票查验信息-------------------:", respData)

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
	if bill.Stage != 6 && isNotCheck && strings.Contains(bill.Remark, "<<发票查验>>") {
		mes += "查验失败：发票查验相关字段的数据未作修改，不执行发票查验功能;"
	}
	// bill.Remark = strings.Replace(bill.Remark, "<<发票查验>>", "", 1)
	fmt.Println("------------mes-------------------:", mes)
	if mes != "" {
		obj.Bill.WrongNote += mes
	}

	for i, invoices := range obj.Invoice {
		for j, fields := range invoices.Invoice {
			for k, field := range fields {
				// if "fc093" == field.Code {
				// 	_, fc005 := GetOneField(fields, "fc005", true)
				// 	obj.Invoice[i].Invoice[j][k].FinalValue = fc005
				// }
				// if "fc096" == field.Code {
				// 	_, fc091 := GetOneField(fields, "fc091", false)
				// 	if fc091 == "2" {
				// 		_, fc011 := GetOneField(fields, "fc011", true)
				// 		obj.Invoice[i].Invoice[j][k].FinalValue = fc011
				// 	}
				// }
				if RegIsMatch(field.Code, `^(fc091|fc092|fc093|fc094|fc095|fc096|fc097)$`) {
					_, fc090 := GetOneField(fields, "fc090", false)
					if fc090 == "2" {
						obj.Invoice[i].Invoice[j][k].FinalValue = ""
					}
				}
			}
		}
	}

	if len(wrong) > 0 {
		obj.Bill.WrongNote += "最后一页存在?，请检查；"
	}

	for _, loc := range fieldLocationMap["fc098"] {
		val := utils.GetFieldValueByLoc(obj, loc, true)
		if val != "" {
			fc098s := utils.RegMatchAll(val, `[^;；]+`)
			for _, fc098 := range fc098s {
				utils.SetIssue(obj, loc, fc098, "", "")
			}
		}
		// utils.SetOnlyOneFinalValue(obj, loc, "R52.9")

	}

	//CSB0106RC0113000
	//同一发票下，当fc090录入值为1时，fc005的结果值默认为fc092+fc093拼在一起的结果值（该需求代码放在结果数据的最后面）
	//如：fc090为1，fc092结果值为123456，fc093结果值为7890，则fc005结果值为1234567890
	for _, loc := range fieldLocationMap["fc090"] {
		fc090Val := utils.GetFieldValueByLoc(obj, loc, false)
		if fc090Val == "1" {
			fc005Loc := utils.GetFieldLoc(fieldLocationMap["fc005"], loc[0], -1, -1, -1)
			fc092Loc := utils.GetFieldLoc(fieldLocationMap["fc092"], loc[0], -1, -1, -1)
			fc093Loc := utils.GetFieldLoc(fieldLocationMap["fc093"], loc[0], -1, -1, -1)
			if len(fc005Loc) != 1 || len(fc092Loc) != 1 || len(fc093Loc) != 1 {
				continue
			}
			val := utils.GetFieldValueByLoc(obj, fc092Loc[0], true) + utils.GetFieldValueByLoc(obj, fc093Loc[0], true)
			utils.SetFinalValue(obj, fc005Loc, val)
		}
	}

	////////////////////////项目报表 - 业务明细 - 单据状态
	obj.Bill.BillType = utils2.GetBillType(obj)
	////////////////////////项目报表 - 业务明细 - 单据状态

	err = service.UpdateInvoiceNumByBillNum(bill.ProCode, bill.BillNum, invoiceNums)
	return err, obj
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"diSanFangDanWeiMap", "B0106_陕西国寿理赔_第三方出具单位", "2", "1"},
		{"jiBingBianMaMap", "B0106_陕西国寿理赔_ICD10疾病编码", "0", "1"},
		{"jiBingBianMaCodeToNameMap", "B0106_陕西国寿理赔_ICD10疾病编码", "1", "0"},
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
		{"shuJuKuBianMa", "B0106_陕西国寿理赔_数据库编码对应表", "1"},
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
			constObj[strings.Replace(k, "B0106_陕西国寿理赔_", "", -1)] = tempMap
		}
		if strings.Index(k, "医疗机构") != -1 {
			tempMap := make(map[string][]string, 0)
			for _, arr := range items {
				tempMap[arr[0]] = arr
			}
			constObj[strings.Replace(k, "B0106_陕西国寿理赔_", "", -1)] = tempMap
		}
	}
	return constObj
}
