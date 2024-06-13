/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022年12月16日10:14:22
 */

package B0114

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"reflect"
	"runtime/debug"
	"server/global"
	model2 "server/module/export/model"
	utils2 "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

//华夏理赔：
//fc179模板类型字段
//fc104账单号
//fc132费用金额（总金额）

//3-发票/生育收据;
//4-清单;
//5-报销单;

//0-申请书4（第一栏为被保险人）;
//1-申请书1（有受益人）;
//2-申请书2（第一栏为手写出险人）;
//20-申请书3（第一栏为印刷体出险人）;
//6-伤病诊断;
//8-法院判决日期;
//9-手术术式信息;
//10-救护车津贴;
//11-身故信息;
//61-重疾;
//62-轻症;
//63-中症;
//64-全残;
//65-残疾;
//66-特种病;
//41-抢救室/ICU病房/烧伤病房津贴;

//和发票有对应关系的：
//3-发票/生育收据;
//4-清单;
//5-报销单;

//MB002-bc001  fc094发票属性
//MB002-bc010  fc091清单所属发票
//MB002-bc011  fc092报销单属性

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:    model2.TypeCode{FieldCode: []string{"fc094"}, BlockCode: []string{"bc001"}},
	QingDanFieldCode:    model2.TypeCode{FieldCode: []string{"fc091"}, BlockCode: []string{"bc010"}},
	BaoXiaoDanFieldCode: model2.TypeCode{FieldCode: []string{"fc092"}, BlockCode: []string{"bc011"}},
	OtherTempType:       map[string]string{"0": "0", "1": "1", "2": "2", "20": "20", "6": "6", "8": "8", "9": "9", "10": "10", "11": "11", "61": "61", "62": "62", "63": "63", "64": "64", "65": "65", "66": "66", "41": "41"},
	TempTypeField:       "fc179",
	InvoiceNumField:     []string{"fc104"},
	MoneyField:          []string{"fc132"},
	//InvoiceTypeField: "fc003",
}

// ResultData B0114
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
	specialConstMap := specialConst(bill.ProCode)

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

	//CSB0114RC0067000
	//当fc179不存在1或2或20或0时，且下载文件中的accidentDate、accidentAreaOption、accidentAddr、accidentDesc其中一个值为空时，
	//fc179出问题件：问题件代码02，问题件描述：缺少理赔申请表，且传递信息存在为空情况（fc179为循环字段）
	accidentDate := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentDate")
	accidentAreaOption := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentAreaOption")
	accidentAddr := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentAddr")
	accidentDesc := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentDesc")
	fc179ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc179"], false)
	global.GLog.Info("fc179", zap.Any("val", fc179ValArr))
	if !(utils.HasItem(fc179ValArr, "1") && utils.HasItem(fc179ValArr, "2") &&
		utils.HasItem(fc179ValArr, "20") && utils.HasItem(fc179ValArr, "0")) &&
		(accidentDate == "" || accidentAreaOption == "" || accidentAddr == "" || accidentDesc == "") {
		utils.SetIssues(obj, fieldLocationMap["fc179"], "缺少理赔申请表，且传递信息存在为空情况", "02", "")
	}

	//CSB0114RC0068000
	//"当下载文件中的applyCauses中的causeCode包含01或02或21或22时，符合以下任意条件，出问题件，问题件代码02，问题件描述：缺少账单信息
	//A、fc179不存在3时
	//B、fc179存在3，但不存在对应的4或5时（fc094,fc092,fc091进行判断是否对应，fc094为发票属性，如fc094没有与fc092,fc091其中一个字段录入的内容相同时，则视为无对应的4和5）
	//fc179,fc094,fc092,fc091均为循环字段"
	isIssue := false
	if len(obj.Invoice) < 2 {
		isIssue = true
	}
	for _, invoiceMap := range obj.Invoice {
		if invoiceMap.Id != "other" && len(invoiceMap.BaoXiaoDan) == 0 && len(invoiceMap.QingDan) == 0 {
			isIssue = true
		}
	}
	if isIssue {
		utils.SetIssues(obj, fieldLocationMap["fc179"], "缺少账单信息", "02", "")
	}

	//CSB0114RC0069000
	//fc002录入内容包含?时，结果数据默认为下载文件中的accidentDate的日期；
	//如果fc002录入内容包含?且下载文件中的accidentDate为空，
	//则fc002出问题件：问题件代码：01，问题件描述：事故日期影像不清晰
	//CSB0114RC0070000
	//fc002录入内容为A时，结果数据默认为下载文件中的accidentDate的日期；
	//如果fc002录入内容为A且下载文件中的accidentDate为空，
	//则fc002出问题件：问题件代码：04，问题件描述：事故日期未填写内容。
	issueArr := [][]string{
		{"?", "01", "事故日期影像不清晰"},
		{"A", "04", "事故日期未填写内容"},
	}
	fc002ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc002"], false)
	for i, issue := range issueArr {
		if len(fc002ValArr) == 1 &&
			((i == 0 && strings.Index(fc002ValArr[0], issue[0]) != -1) ||
				(i == 1 && fc002ValArr[0] == issue[0])) {
			utils.SetFinalValue(obj, fieldLocationMap["fc002"], accidentDate)
			if accidentDate == "" {
				utils.SetIssues(obj, fieldLocationMap["fc002"], issue[2], issue[1], "")
			}
		}
	}

	//CSB0114RC0071000
	//fc002的结果日期晚于接收日期时，fc002结果值格式2022/01/01
	//fc002出问题件：问题件代码07，问题件描述：事故日期应该在受理日期（传给BPO的日期）之前。（fc002录入格式为6位数日期，接收日期为8位日期）
	//CSB0114RC0072000
	//fc002的结果日期晚于下载文件中caseReportDate的日期时，
	//fc002出问题件：问题件代码07，问题件描述：事故日期应该在报案记录时间之前。（caseReportDate日期格式为2018/03/13 10:28:39，只需要判断日期即可）
	fc002FinalValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc002"], true)
	if len(fc002FinalValArr) == 1 {
		fc002FinalValTime, _ := time.Parse("2006/01/02", fc002FinalValArr[0])
		if fc002FinalValTime.After(obj.Bill.CreatedAt) {
			utils.SetIssues(obj, fieldLocationMap["fc002"], "事故日期应该在受理日期（传给BPO的日期）之前", "07", "")
		}

		caseReportDate := utils.GetNodeValue(obj.Bill.OtherInfo, "caseReportDate")
		caseReportDateTime, _ := time.Parse("2006/01/02  15:04:05", caseReportDate)
		if fc002FinalValTime.After(caseReportDateTime) {
			utils.SetIssues(obj, fieldLocationMap["fc002"], "事故日期应该在报案记录时间之前", "07", "")
		}
	}

	//CSB0114RC0073000
	//fc008录入内容为A时，默认为下载文件中的accidentAreaOption的值，
	//如下载文件中的值为空时，fc008结果数据默认为空，出问题件：问题件代码04，问题件描述：事故地所属区/县未填写
	//CSB0114RC0074000
	//fc008录入内容为F时，默认为下载文件中的accidentAreaOption的值，
	//如下载文件中的值为空时，fc008结果数据默认为空，出问题件：问题件代码04，问题件描述：事故地所属区/县与代码库不匹配
	//CSB0114RC0075000
	//fc008录入内容包含?时，默认为下载文件中的accidentAreaOption的值，
	//如下载文件中的值为空时，fc008结果数据默认为空，出问题件：问题件代码01，问题件描述：事故地所属区/县影像不清晰
	issueArr = [][]string{
		{"A", "04", "事故地所属区/县未填写"},
		{"F", "04", "事故地所属区/县与代码库不匹配"},
		{"?", "01", "事故地所属区/县影像不清晰"},
	}
	fc008ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc008"], false)
	for _, issue := range issueArr {
		if len(fc008ValArr) == 1 && strings.Index(fc008ValArr[0], issue[0]) != -1 {
			utils.SetFinalValue(obj, fieldLocationMap["fc008"], accidentAreaOption)
			if accidentAreaOption == "" {
				utils.SetIssues(obj, fieldLocationMap["fc008"], issue[2], issue[1], "")
			}
		}
	}

	//CSB0114RC0076000
	//fc008根据以下规则进行匹配转码
	//a.录入内容匹配常量表《B0114_华夏理赔_华夏理赔地址库》的区县（第五列）称对应的事故地所属区/县编码（第二列）进行转码。
	//b.fc008录入内容不唯一时，则根据fc003,fc004的内容匹配到唯一的值进行转码
	//c.以上两个条件均无法转换时,fc008默认为下载XML文件中的accidentAreaOption的值
	//d.fc008录入内容为A或包含?时，fc008默认为默认为下载XML文件中的accidentAreaOption的值"
	fc004ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc004"], false)
	fc003ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc003"], false)
	if len(fc008ValArr) == 1 {
		val, ok := constMap["diZhiMap"][fc008ValArr[0]]
		if ok {
			utils.SetFinalValue(obj, fieldLocationMap["fc008"], val)
			valNumStr, _ := constMap["diZhiShengShiQuMapNum"][fc008ValArr[0]]
			valNum, _ := strconv.Atoi(valNumStr)
			if valNum > 1 {
				str := fc008ValArr[0]
				if fc003ValArr[0] != "" {
					str = fc003ValArr[0] + "_" + str
				}
				if fc004ValArr[0] != "" {
					str = fc004ValArr[0] + "_" + str
				}
				val, _ = constMap["diZhiShengShiQuMap"][str]
				utils.SetFinalValue(obj, fieldLocationMap["fc008"], val)
			}
		} else {
			utils.SetFinalValue(obj, fieldLocationMap["fc008"], accidentAreaOption)
			//if fc008ValArr[0] == "A" || strings.Index(fc008ValArr[0], "?") != -1 {
			//	utils.SetFinalValue(obj, fieldLocationMap["fc008"], accidentAreaOption)
			//}
		}
	}

	//CSB0114RC0077000
	//fc009录入内容为A时，fc009的结果数据默认为下载文件中的accidentAddr的地点；
	//如果fc009录入内容为A且下载文件中的accidentAddr为空，则fc009输出空， fc009出问题件：问题件代码：04，问题件描述：事故地点未填写内容。
	//CSB0114RC0078000
	//fc009录入内容包含“?”时，fc009的结果数据默认为下载文件中的accidentAddr的地点；
	//如果fc009录入内容包含“?”且下载文件中的accidentAddr为空，则fc009输出空， fc009出问题件：问题件代码：01，问题件描述：事故地点影像不清晰
	issueArr = [][]string{
		{"?", "01", "事故地点影像不清晰"},
		{"A", "04", "事故地点未填写内容"},
	}
	fc009ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc009"], false)
	for i, issue := range issueArr {
		if len(fc009ValArr) == 1 &&
			((i == 0 && strings.Index(fc009ValArr[0], issue[0]) != -1) ||
				(i == 1 && fc009ValArr[0] == issue[0])) {
			utils.SetFinalValue(obj, fieldLocationMap["fc009"], accidentAddr)
			if accidentAddr == "" {
				utils.SetFinalValue(obj, fieldLocationMap["fc009"], "")
				utils.SetIssues(obj, fieldLocationMap["fc009"], issue[2], issue[1], "")
			}
		}
	}

	//CSB0114RC0079000
	//fc012录入内容为A时，fc012的结果数据默认为下载文件中的accidentDesc的内容；
	//如果fc012录入内容为A且下载文件中的accidentDesc为空，则fc012输出空，fc012出问题件：问题件代码：04，问题件描述：事故经过未填写内容。
	//CSB0114RC0080000
	//fc012录入内容包含“?”时，fc012的结果数据默认为下载文件中的accidentDesc的内容；
	//如果fc012录入内容包含“?”且下载文件中的accidentDesc为空，则fc012输出空，fc012出问题件：问题件代码：01，问题件描述：事故经过影像不清晰
	issueArr = [][]string{
		{"?", "01", "事故经过影像不清晰"},
		{"A", "04", "事故经过未填写内容"},
	}
	fc012ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc012"], false)
	for i, issue := range issueArr {
		if len(fc012ValArr) == 1 &&
			((i == 0 && strings.Index(fc012ValArr[0], issue[0]) != -1) ||
				(i == 1 && fc012ValArr[0] == issue[0])) {
			utils.SetFinalValue(obj, fieldLocationMap["fc012"], accidentDesc)
			if accidentDesc == "" {
				utils.SetFinalValue(obj, fieldLocationMap["fc012"], "")
				utils.SetIssues(obj, fieldLocationMap["fc012"], issue[2], issue[1], "")
			}
		}
	}

	//CSB0114RC0081000
	//当下载文件XML中的“accidentDesc”的值大于30个字符（标点符号也要算），
	//且fc012的录入值为“A”或“空”时，将下载文件XML中的“accidentDesc”的值赋值到fc012的结果值中
	charNum := utils.CalcChar(accidentDesc)
	fmt.Println(charNum)
	if charNum > 30 && len(fc012ValArr) == 1 && (fc012ValArr[0] == "" || fc012ValArr[0] == "A") {
		utils.SetFinalValue(obj, fieldLocationMap["fc012"], accidentDesc)
	}

	//CSB0114RC0082000
	//"fc014录入内容匹配常量表《B0114_华夏理赔_华夏理赔意外原因表》的意外编码描述（第一列）对应的意外编码（第二列）进行转码；
	//当fc014录入内容不为常量表《B0114_华夏理赔_华夏理赔意外原因表》的意外编码描述（第一列）时（录入内容为A、为空或包含?时不做校验），
	//fc014的结果值默认为空，fc014出问题件：问题件代码07，问题件描述：意外细节均在常量库中无法找到。"
	//CSB0114RC0083000
	//"fc015录入内容匹配常量表《B0114_华夏理赔_华夏理赔损伤外部原因表》的“描述”（第一列）对应的“编码”（第二列）进行转码；
	//当fc015录入内容不为常量表《B0114_华夏理赔_华夏理赔损伤外部原因表》的“描述”（第一列）时（录入内容为A、为空或包含?时不做校验），fc015的结果值默认为空，
	//fc015出问题件：问题件代码07，问题件描述：损伤外部原因均在常量库中无法找到。"
	//CSB0114RC0084000
	//"fc096,fc187,fc190录入内容匹配常量表《B0114_华夏理赔_手术术式编码表》的“手术名称”（第一列）对应的“手术编码”（第二列）进行转码；
	//fc096,fc187,fc190的内容不为常量表《B0114_华夏理赔_手术术式编码表》的手术名称（第一列）时（录入内容为空、为A或包含?时不做校验），fc096,fc187,fc190对应的结果值默认为空，
	//fc096,fc187,fc190对应的字段出问题件：问题件代码07，问题件描述：手术术士编码均在常量库中无法找到。"
	//CSB0114RC0086000
	//"fc101,fc181,fc182录入内容匹配常量表《B0114_华夏理赔_ICD疾病代码表》的疾病名称（第三列）对应的疾病编码（第二列）进行转码；
	//当fc101,fc181,fc182录入内容不为常量表《B0114_华夏理赔_ICD疾病代码表》的疾病名称（第一列）时（录入内容为A、为空或包含?时不做校验），fc101,fc181,fc182对应的结果值默认为空，
	//fc101,fc181,fc182对应的字段出问题件：问题件代码07，问题件描述：伤病疾病代码均在常量库中无法找到。"
	//CSB0114RC0112000
	//fc016录入内容匹配常量表《B0114_华夏理赔_重疾和轻症疾病名称表》的“疾病名称”（第一列）对应的“编码”（第二列）进行转码（录入内容为A、为空或包含?时不做校验），无法转码时fc016的结果值默认为空，fc016出问题件：问题件代码07，描述：重大疾病名称与代码表不匹配
	//CSB0114RC0113000
	//fc019录入内容匹配常量表《B0114_华夏理赔_重疾和轻症疾病名称表》的“疾病名称”（第一列）对应的“编码”（第二列）进行转码，（录入内容为A、为空或包含?时不做校验），无法转码时fc019的结果值默认为空，fc019出问题件：问题件代码07，描述：轻症重疾名称与代码表不匹配
	//CSB0114RC0114000
	//fc021录入内容匹配常量表《B0114_华夏理赔_重疾和轻症疾病名称表》的“疾病名称”（第一列）对应的“编码”（第二列）进行转码，（录入内容为A、为空或包含?时不做校验），无法转码时fc021的结果值默认为空，fc021出问题件：问题件代码07，描述：中症疾病名称与代码表不匹配
	//CSB0114RC0115000
	//fc023录入内容匹配常量表《B0114_华夏理赔_重疾和轻症疾病名称表》的“疾病名称”（第一列）对应的“编码”（第二列）进行转码，（录入内容为A、为空或包含?时不做校验），无法转码时fc023的结果值默认为空，fc023出问题件：问题件代码07，描述：特种病名称与代码表不匹配
	//CSB0114RC0116000
	//fc025录入内容匹配常量表《B0114_华夏理赔_全残信息表》的“伤残名称”（第四列）对应的“伤残代码”（第三列）进行转码（录入内容为A、为空或包含?时不做校验），无法转码时fc025的结果值默认为空，fc025出问题件：问题件代码07，描述：全残项目与代码表不匹配
	constChangeArr := [][][]string{
		{{"fc014"}, {"yiWaiYuanYinMap"}, {"07", "意外细节均在常量库中无法找到"}},
		{{"fc015"}, {"sunShangWaiBuYuanYinMap"}, {"07", "损伤外部原因均在常量库中无法找到"}},
		{{"fc096", "fc187", "fc190"}, {"shouShuBianMaMap"}, {"07", "手术术士编码均在常量库中无法找到"}},
		{{"fc101", "fc181", "fc182"}, {"jiBingDaiMaMap"}, {"07", "伤病疾病代码均在常量库中无法找到"}},
		{{"fc016"}, {"zhongJiQingJiMingChengMap"}, {"07", "重大疾病名称与代码表不匹配"}},
		{{"fc019"}, {"zhongJiQingJiMingChengMap"}, {"07", "轻症重疾名称与代码表不匹配"}},
		{{"fc021"}, {"zhongJiQingJiMingChengMap"}, {"07", "中症疾病名称与代码表不匹配"}},
		{{"fc023"}, {"zhongJiQingJiMingChengMap"}, {"07", "特种病名称与代码表不匹配"}},
		{{"fc025"}, {"quanCanMap"}, {"07", "全残项目与代码表不匹配"}},
	}
	for _, constChange := range constChangeArr {
		for _, code := range constChange[0] {
			fc014ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap[code], false)
			for _, fc014Val := range fc014ValArr {
				v, ok := constMap[constChange[1][0]][fc014Val]
				if ok {
					utils.SetFinalValue(obj, fieldLocationMap[code], v)
				} else {
					if fc014Val != "A" && fc014Val != "" && strings.Index(fc014Val, "?") == -1 {
						utils.SetFinalValue(obj, fieldLocationMap[code], "")
						utils.SetIssues(obj, fieldLocationMap[code], constChange[2][1], constChange[2][0], "")
					}
				}
			}
		}
	}

	//CSB0114RC0085000
	//"将下列左边录入内容复制到的右边字段的结果数据中（左边录入内容为A或F或包含“?”时，右边结果数据默认为空）
	//fc096,fc098
	//fc187,fc189
	//fc190,fc192"
	iCodeArr := [][]string{
		{"fc096", "fc098"},
		{"fc187", "fc189"},
		{"fc190", "fc192"},
	}
	for _, iCode := range iCodeArr {
		fc096ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap[iCode[0]], false)
		if len(fc096ValArr) == 1 {
			if fc096ValArr[0] == "A" || fc096ValArr[0] == "F" || strings.Index(fc096ValArr[0], "?") != -1 {
				utils.SetFinalValue(obj, fieldLocationMap[iCode[1]], "")
			} else {
				utils.SetFinalValue(obj, fieldLocationMap[iCode[1]], fc096ValArr[0])
			}
		}
	}

	//CSB0114RC0087000
	//"1.第二列字段录入内容与常量表中《B0114_华夏理赔_ICD疾病代码》的“疾病名称”（第三列）进行匹配，将匹配对应的“诊断类型”（第一列）的值默认到第一列的结果数据中
	//2.将第二列字段的录入内容复制到第三列字段的结果数据中
	//诊断类型  诊断代码  诊断名称
	//fc099,fc101,fc102
	//fc183,fc181,fc185
	//fc184,fc182,fc186"
	iCodeArr = [][]string{
		{"fc099", "fc101", "fc102"},
		{"fc183", "fc181", "fc185"},
		{"fc184", "fc182", "fc186"},
	}
	for _, iCode := range iCodeArr {
		fc101ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap[iCode[1]], false)
		if len(fc101ValArr) == 1 {
			v, ok := constMap["jiBingDaiMaTypeMap"][fc101ValArr[0]]
			if ok {
				utils.SetFinalValue(obj, fieldLocationMap[iCode[0]], v)
			}
			utils.SetFinalValue(obj, fieldLocationMap[iCode[2]], fc101ValArr[0])
		}
	}

	//CSB0114RC0088000	fc266录入内容为6时,fc105结果数据默认为2（发票信息，fc266,fc105需要相互对应在同一个发票中）
	for _, fc266Loc := range fieldLocationMap["fc266"] {
		fc266Val := utils.GetFieldValueByLoc(obj, fc266Loc, false)
		fc105Loc := utils.GetFieldLoc(fieldLocationMap["fc105"], fc266Loc[0], -1, -1, -1)
		if fc266Val == "6" {
			utils.SetFinalValue(obj, fc105Loc, "2")
		}
	}

	//CSB0114RC0089000	fc175默认为对应的发票或报销单中的fc104的值
	//CSB0114RC0090000	fc176默认为对应的发票或报销单中的fc104的值
	iCodeArr = [][]string{
		{"fc175", "fc104"},
		{"fc176", "fc104"},
	}
	for _, iCode := range iCodeArr {
		for _, fc104Loc := range fieldLocationMap[iCode[1]] {
			fc104Val := utils.GetFieldValueByLoc(obj, fc104Loc, false)
			fc175Loc := utils.GetFieldLoc(fieldLocationMap[iCode[0]], fc104Loc[0], -1, -1, -1)
			utils.SetFinalValue(obj, fc175Loc, fc104Val)
		}
	}

	//CSB0114RC0091000
	//fc106的录入内容匹配常量表《B0114_华夏理赔_医院名称表》的医院名称（第三列）对应的医院编码（第二列）进行转码输出。
	//CSB0114RC0092000
	//fc106的录入值包含“?”时，fc106的结果值转为空并清空fc106的问题件
	//fc106ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc106"], false)
	for _, fc106Loc := range fieldLocationMap["fc106"] {
		fc106Val := utils.GetFieldValueByLoc(obj, fc106Loc, false)
		v, ok := constMap["yiYuanMingChengMap"][fc106Val]
		if ok {
			utils.SetOnlyOneFinalValue(obj, fc106Loc, v)
		}
		if strings.Index(fc106Val, "?") != -1 {
			utils.SetOnlyOneFinalValue(obj, fc106Loc, "")
			utils.DelOnlyOneIssue(obj, fc106Loc)
		}
	}

	//CSB0114RC0094000	fc108的录入日期早于fc002的录入日期时，fc108出问题件：问题件代码07，问题件描述：账单日期不能早于事故日期。（可同一天）
	//CSB0114RC0099000	fc110的录入日期早于fc002的录入日期时，fc110出问题件：问题件代码07，问题件描述：就诊日期不能早于事故日期。（可同一天）
	//CSB0114RC0100000	fc111的录入日期早于fc002的录入日期时，fc111出问题件：问题件代码07，问题件描述：入院日期不能早于事故日期。（可同一天）
	dateMap := [][]string{
		{"fc108", "fc002", "账单日期不能早于事故日期", "07"},
		{"fc110", "fc002", "就诊日期不能早于事故日期", "07"},
		{"fc111", "fc002", "入院日期不能早于事故日期", "07"},
	}
	for _, temp := range dateMap {
		for _, fc108Loc := range fieldLocationMap[temp[0]] {
			fc002Val := ""
			if len(fc002ValArr) == 1 {
				fc002Val = fc002ValArr[0]
			}
			fc108Val := utils.GetFieldValueByLoc(obj, fc108Loc, false)
			if fc108Val == "" {
				continue
			}
			fc108NumVal, _ := strconv.Atoi(fc108Val)
			fc002NumVal, _ := strconv.Atoi(fc002Val)
			if fc108NumVal < fc002NumVal {
				utils.SetIssue(obj, fc108Loc, temp[2], temp[3], "")
			}
		}
	}

	//CSB0114RC0093000	fc108的录入日期早于fc111的录入日期时，fc108出问题件：问题件代码07，问题件描述：账单日期不能早于入院日期。（可同一天）
	//CSB0114RC0101000	fc112的录入日期早于fc111的录入日期时，fc112出问题件：问题件代码07，问题件描述：出院日期不能早于入院日期且不能是同一天。
	dateMap = [][]string{
		{"fc108", "fc111", "账单日期不能早于入院日期", "07"},
		{"fc112", "fc111", "出院日期不能早于入院日期且不能是同一天", "07"},
	}
	for _, temp := range dateMap {
		for _, loc := range fieldLocationMap[temp[0]] {
			fc108Val := utils.GetFieldValueByLoc(obj, loc, false)
			fc111Loc := utils.GetFieldLoc(fieldLocationMap[temp[1]], loc[0], loc[1], -1, -1)
			if len(fc111Loc) != 1 {
				continue
			}
			fc111Val := utils.GetFieldValueByLoc(obj, fc111Loc[0], false)
			if fc111Val == "" || fc108Val == "" {
				continue
			}
			fc108NumVal, _ := strconv.Atoi(fc108Val)
			fc111NumVal, _ := strconv.Atoi(fc111Val)
			if fc108NumVal < fc111NumVal {
				utils.SetIssue(obj, loc, temp[2], temp[3], "")
				utils.SetFinalValue(obj, fc111Loc, "")
			}
		}
	}

	//CSB0114RC0102000	fc266录入内容为6时，fc281的结果值等于fc132的结果值。
	for _, fc266Loc := range fieldLocationMap["fc266"] {
		fc132Loc := utils.GetFieldLoc(fieldLocationMap["fc132"], fc266Loc[0], -1, -1, -1)
		fc281Loc := utils.GetFieldLoc(fieldLocationMap["fc281"], fc266Loc[0], -1, -1, -1)
		if len(fc132Loc) == 1 {
			fc132Val := utils.GetFieldValueByLoc(obj, fc132Loc[0], true)
			utils.SetFinalValue(obj, fc281Loc, fc132Val)
		}
	}

	//CSB0114RC0103000	fc266录入内容为6时，且fc281结果数据为空时，fc281出问题件：问题件代码04，问题件描述：住院发生金额未填写内容。(此校验放在“fc266录入内容为6时，fc281的结果值等于fc132的结果值”的后面)
	for _, fc266Loc := range fieldLocationMap["fc266"] {
		fc266Val := utils.GetFieldValueByLoc(obj, fc266Loc, false)
		fc281Loc := utils.GetFieldLoc(fieldLocationMap["fc281"], fc266Loc[0], -1, -1, -1)
		if len(fc281Loc) == 1 {
			fc281Val := utils.GetFieldValueByLoc(obj, fc281Loc[0], true)
			if fc281Val == "" && fc266Val == "6" {
				utils.SetIssues(obj, fc281Loc, "住院发生金额未填写内容", "04", "")
			}
		}
	}

	//CSB0114RC0104000	fc126为空时结果数据默认为fc120,fc121,fc122,fc124的结果值相加(fc120,fc121,fc122,fc124为空、为A时不相加）（为空、为A或包含?的字段，不相加）
	for _, fc126Loc := range fieldLocationMap["fc126"] {
		fc120Loc := utils.GetFieldLoc(fieldLocationMap["fc120"], fc126Loc[0], -1, -1, -1)
		fc121Loc := utils.GetFieldLoc(fieldLocationMap["fc121"], fc126Loc[0], -1, -1, -1)
		fc122Loc := utils.GetFieldLoc(fieldLocationMap["fc122"], fc126Loc[0], -1, -1, -1)
		fc124Loc := utils.GetFieldLoc(fieldLocationMap["fc124"], fc126Loc[0], -1, -1, -1)
		fc126Val := utils.GetFieldValueByLoc(obj, fc126Loc, true)
		if fc126Val == "" && len(fc120Loc) == 1 && len(fc121Loc) == 1 && len(fc122Loc) == 1 && len(fc124Loc) == 1 {
			fc120Val := utils.GetFieldDecimalValueByLoc(obj, fc120Loc[0], true)
			fc121Val := utils.GetFieldDecimalValueByLoc(obj, fc121Loc[0], true)
			fc122Val := utils.GetFieldDecimalValueByLoc(obj, fc122Loc[0], true)
			fc124Val := utils.GetFieldDecimalValueByLoc(obj, fc124Loc[0], true)
			utils.SetOnlyOneFinalValue(obj, fc126Loc, fc120Val.Add(fc121Val).Add(fc122Val).Add(fc124Val).String())
		}
	}

	//CSB0114RC0095000	fc105录入内容为1时，且fc110录入内容为空或A时，fc110输出空，fc110出问题件：问题件代码：04，问题件描述：就诊日期未填写内容。
	//CSB0114RC0096000	fc105录入内容为2时，且fc111录入内容为空或A时，fc111输出空，fc111出问题件：问题件代码：04，问题件描述：入院日期未填写内容。
	//CSB0114RC0097000	fc105录入内容为2时，且fc112录入内容为空或A时，fc112输出空，fc112出问题件：问题件代码：04，问题件描述：出院日期未填写内容。
	dateMap = [][]string{
		{"fc105", "1", "fc110", "就诊日期未填写内容", "04"},
		{"fc105", "2", "fc111", "入院日期未填写内容", "04"},
		{"fc105", "2", "fc112", "出院日期未填写内容", "04"},
	}
	for _, temp := range dateMap {
		for _, loc := range fieldLocationMap[temp[0]] {
			fc105Val := utils.GetFieldValueByLoc(obj, loc, false)
			fc110Loc := utils.GetFieldLoc(fieldLocationMap[temp[2]], loc[0], loc[1], -1, -1)
			if len(fc110Loc) != 1 {
				continue
			}
			fc110Val := utils.GetFieldValueByLoc(obj, fc110Loc[0], false)
			if fc105Val == temp[1] && (fc110Val == "" || fc110Val == "A") {
				utils.SetIssues(obj, fc110Loc, temp[3], temp[4], "")
				utils.SetFinalValue(obj, fc110Loc, "")
			}
		}
	}

	//CSB0114RC0098000	fc110的结果日期晚于录入当天日期时，fc110出问题件：问题件代码：07，问题件描述：就诊日期不能晚于当前日期。
	for _, loc := range fieldLocationMap["fc110"] {
		fc110Val := utils.GetFieldValueByLoc(obj, loc, false)
		nowStr := time.Now().Format("060102")
		nowVal, _ := strconv.Atoi(nowStr)
		fc110NumVal, _ := strconv.Atoi(fc110Val)
		if fc110NumVal > nowVal {
			utils.SetIssues(obj, fieldLocationMap["fc110"], "就诊日期不能晚于当前日期", "07", "")
		}
	}

	//CSB0114RC0117000
	//fc026的结果数据根据fc025录入内容匹配常量表《B0114_华夏理赔_全残信息表》的“伤残名称”（第四列）自动带出对应的“赔付比例”（第一列）
	fc025ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc025"], false)
	if len(fc025ValArr) == 1 {
		v, ok := constMap["quanCanBiLiMap"][fc025ValArr[0]]
		if ok {
			utils.SetFinalValue(obj, fieldLocationMap["fc026"], v)
		}
	}

	//CSB0114RC0118000
	//下载文件中的applyCauses中的causeCode含05,06时，且fc026的结果值不为1.00时，fc026出问题件：问题件代码：07，问题件描述：鉴定结果表明没有达到全残。
	causeCode := utils.GetNodeData(obj.Bill.OtherInfo, "causeCode")
	fc026ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc026"], true)
	if len(fc026ValArr) == 1 {
		if (utils.HasItem(causeCode, "05") || utils.HasItem(causeCode, "06")) && fc026ValArr[0] != "1.00" {
			utils.SetIssues(obj, fieldLocationMap["fc026"], "鉴定结果表明没有达到全残", "07", "")
		}
	}

	//CSB0114RC0119000
	//fc029录入内容匹配常量表《B0114_华夏理赔_伤残项目表》的“伤残名称”（第三列）对应的“伤残代码”（第二列）进行转码，
	//无法转码时默认为空，出问题件：问题件代码07，描述：伤残项目与代码表不匹配
	//CSB0114RC0120000
	//fc031的结果数据根据fc029录入内容匹配常量表《B0114_华夏理赔_伤残项目表》的“伤残名称”（第三列）自动带出对应的“赔付比例”（第六列）（录入内容为F或包含?时不转码）。
	fc029ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc029"], false)
	if len(fc029ValArr) == 1 {
		v, ok := constMap["shanCanXiangMuMap"][fc029ValArr[0]]
		if ok {
			utils.SetFinalValue(obj, fieldLocationMap["fc029"], v)
		} else {
			utils.SetFinalValue(obj, fieldLocationMap["fc029"], "")
			utils.SetIssues(obj, fieldLocationMap["fc029"], "伤残项目与代码表不匹配", "07", "")
		}

		if fc029ValArr[0] != "F" && strings.Index(fc029ValArr[0], "?") == -1 {
			v, ok = constMap["shanCanXiangMuBiLiMap"][fc029ValArr[0]]
			if ok {
				utils.SetFinalValue(obj, fieldLocationMap["fc031"], v)
			}
		}
	}

	//CSB0114RC0123000
	//下列循环字段的第三列字段的结果值等于第一列字段的结果值乘以第二列字段的结果值：
	//fc231 fc247 fc139
	//fc232 fc248 fc143
	//fc233 fc249 fc147
	//fc234 fc250 fc151
	//fc235 fc251 fc155
	//fc236 fc252 fc159
	//fc237 fc253 fc163
	//fc238 fc254 fc167
	//（该需求放在结果数据需求“1.根据fc008（县）,fc003（市）,fc004（省）的录入内容匹配出事故发生地对应的医疗目录，匹配规则如下...”之前）
	iCodeArr = [][]string{
		{"fc231", "fc247", "fc139"},
		{"fc232", "fc248", "fc143"},
		{"fc233", "fc249", "fc147"},
		{"fc234", "fc250", "fc151"},
		{"fc235", "fc251", "fc155"},
		{"fc236", "fc252", "fc159"},
		{"fc237", "fc253", "fc163"},
		{"fc238", "fc254", "fc167"},
	}
	for _, iCode := range iCodeArr {
		for _, oneLoc := range fieldLocationMap[iCode[0]] {
			oneVal := utils.GetFieldDecimalValueByLoc(obj, oneLoc, false)
			twoLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[1]], oneLoc[0], oneLoc[1], oneLoc[2], oneLoc[4])
			threeLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[2]], oneLoc[0], oneLoc[1], oneLoc[2], oneLoc[4])
			if len(twoLocArr) == 1 && len(threeLocArr) == 1 {
				twoVal := utils.GetFieldDecimalValueByLoc(obj, twoLocArr[0], true)
				utils.SetFinalValue(obj, threeLocArr, oneVal.Mul(twoVal).String())
			}
		}
	}

	//CSB0114RC0124000
	//一、当以下第五列字段录入值为“4”时，对应第三列字段结果值默认为“1”,下列字段中第四列字段内容默认为第二列字段*第三列字段的结果值；
	//二、当以下第五列字段录入值不为“4”（为“5、?或空”）时，对应的字段执行以下六条校验：
	//1.根据fc008（县）,fc003（市）,fc004（省）的录入内容匹配出事故发生地对应的医疗目录，匹配规则如下
	//a.fc008录入内容匹配常量表《B0114_华夏理赔_华夏理赔地址库》的区县（第五列），如匹配成功且匹配的内容是唯一值时，则匹配对应的“所属医疗目录”（第一列）对应的医疗目录的文件名中的常量
	//b.fc008录入内容匹配常量表《B0114_华夏理赔_华夏理赔地址库》的区县（第五列），如匹配成功但匹配的内容不唯一时，则根据与fc003,fc004的内容匹配到对应的“所属医疗目录”（第一列）对应的医疗目录的文件名中的常量
	//c.fc008匹配常量表《B0114_华夏理赔_华夏理赔地址库》的区县（第五列），无法匹配时，则根据fc003录入内容与《B0114_华夏理赔_华夏理赔地址库》的市（第四列）进行匹配，匹配成功时，则匹配对应的“所属医疗目录”（第一列）对应的医疗目录的文件名中的常量
	//d.fc008,fc003录入内容均与代码库无法匹配时，则根据fc004录入内容与《B0114_华夏理赔_华夏理赔地址库》的省（第三列）进行匹配，匹配成功时，则匹配对应的“所属医疗目录”（第一列）对应的医疗目录的文件名中的常量
	//e.如fc008,fc003,fc004均未匹配成功时，则不进行匹配
	//2.根据客户下载文件中的policyBranchCode中的内容与常量表《B0114_华夏理赔_机构编码对应表》的机构编码（第二列）进行匹配，匹配成功后，则匹配对应的“所属医疗目录”（第一列）对应的医疗目录的文件名中的常量
	//下列字段第三列字段根据第一列字段录入内容，与常量表进行匹配，带出常量表中的比例
	//3.条件1和条件2完成匹配后，如下列第一列字段的录入内容只对应一个医疗目录的中文名称（第三列），则根据该医疗目录，将对应的“自费比例”（第一列）放到对应的第三列字段中
	//4.条件1和条件2完成匹配后，如下列第一列字段的录入内容在两个对应的医疗目录的中文名称（第三列）中均有匹配的内容，则根据对应的""自付比例""（第一列）进行对比，取较小的自费比例放到第三列字段中（自付比例为空的默认为0，多个自付比例均为0时取0）
	//5.当无法匹配到自付比例时，第三个字段则默认为0
	//6.下列字段中第四列字段内容默认为第二列字段*第三列字段的结果值
	//项目名称-清单 金额-清单 自费比例-清单 自费金额-清单 项目类型
	//fc138 fc139 fc140 fc141 fc257
	//fc142 fc143 fc144 fc145 fc258
	//fc146 fc147 fc148 fc149 fc259
	//fc150 fc151 fc152 fc153 fc260
	//fc154 fc155 fc156 fc157 fc261
	//fc158 fc159 fc160 fc161 fc262
	//fc162 fc163 fc164 fc165 fc263
	//fc166 fc167 fc168 fc169 fc264

	//获取医疗目录
	policyBranchCode := utils.GetNodeValue(obj.Bill.OtherInfo, "policyBranchCode")
	yiLiaoMuLu := ""
	yiLiaoMuLuInPolicyBranchCode := ""
	valArr := specialConstMap["diZhiMap"][fc008ValArr[0]]
	if len(valArr) == 1 {
		yiLiaoMuLu = valArr[0]
	} else if len(valArr) > 1 {
		str := fc008ValArr[0]
		if fc003ValArr[0] != "" {
			str = fc003ValArr[0] + "_" + str
		}
		if fc004ValArr[0] != "" {
			str = fc004ValArr[0] + "_" + str
		}
		val, ok := constMap["diZhiShengShiQuMap"][str]
		if ok {
			yiLiaoMuLu = val
		}
	} else {
		valArr, ok := specialConstMap["diZhiShiMap"][fc003ValArr[0]]
		if ok {
			yiLiaoMuLu = valArr[0]
		} else {
			valArr, ok = specialConstMap["diZhiShengMap"][fc004ValArr[0]]
			if ok {
				yiLiaoMuLu = valArr[0]
			}
		}
	}
	val, ok := constMap["jiGouBianMaMap"][policyBranchCode]
	if ok {
		yiLiaoMuLuInPolicyBranchCode = val
	}
	global.GLog.Info("yiLiaoMuLu", zap.Any("", yiLiaoMuLu))
	global.GLog.Info("yiLiaoMuLuInPolicyBranchCode", zap.Any("", yiLiaoMuLuInPolicyBranchCode))

	//匹配赋值
	iCodeArr = [][]string{
		{"fc138", "fc139", "fc140", "fc141", "fc257"},
		{"fc142", "fc143", "fc144", "fc145", "fc258"},
		{"fc146", "fc147", "fc148", "fc149", "fc259"},
		{"fc150", "fc151", "fc152", "fc153", "fc260"},
		{"fc154", "fc155", "fc156", "fc157", "fc261"},
		{"fc158", "fc159", "fc160", "fc161", "fc262"},
		{"fc162", "fc163", "fc164", "fc165", "fc263"},
		{"fc166", "fc167", "fc168", "fc169", "fc264"},
	}
	for _, iCode := range iCodeArr {
		for _, fiveLoc := range fieldLocationMap[iCode[4]] {
			fiveVal := utils.GetFieldValueByLoc(obj, fiveLoc, false)
			oneLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[0]], fiveLoc[0], fiveLoc[1], fiveLoc[2], fiveLoc[4])
			twoLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[1]], fiveLoc[0], fiveLoc[1], fiveLoc[2], fiveLoc[4])
			threeLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[2]], fiveLoc[0], fiveLoc[1], fiveLoc[2], fiveLoc[4])
			fourLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[3]], fiveLoc[0], fiveLoc[1], fiveLoc[2], fiveLoc[4])
			if len(twoLocArr) != 1 || len(threeLocArr) != 1 || len(fourLocArr) != 1 {
				continue
			}
			twoVal := utils.GetFieldDecimalValueByLoc(obj, twoLocArr[0], true)
			threeVal := utils.GetFieldDecimalValueByLoc(obj, threeLocArr[0], true)
			if fiveVal == "4" {
				utils.SetFinalValue(obj, threeLocArr, "1")
				utils.SetFinalValue(obj, fourLocArr, twoVal.Sub(threeVal).String())

			} else if (fiveVal == "5" || fiveVal == "?" || fiveVal == "") && len(fc008ValArr) == 1 {
				oneVal := utils.GetFieldValueByLoc(obj, oneLocArr[0], true)
				biLi := constMap[yiLiaoMuLu][oneVal]
				biLiInPolicyBranchCode := constMap[yiLiaoMuLuInPolicyBranchCode][oneVal]
				biLiDecimal := utils.ParseDecimal(biLi)
				biLiDecimalInPolicyBranchCode := utils.ParseDecimal(biLiInPolicyBranchCode)
				newBiLi := decimal.Zero
				if biLi == "" {
					newBiLi = biLiDecimalInPolicyBranchCode
				}
				if biLiInPolicyBranchCode == "" {
					newBiLi = biLiDecimal
				}
				if biLiInPolicyBranchCode != "" && biLi != "" {
					if biLiDecimal.LessThanOrEqual(biLiDecimalInPolicyBranchCode) {
						newBiLi = biLiDecimal
					} else {
						newBiLi = biLiDecimalInPolicyBranchCode
					}
				}
				fmt.Println("newBiLi.String()")
				fmt.Println(newBiLi)
				utils.SetFinalValue(obj, threeLocArr, newBiLi.String())
				utils.SetFinalValue(obj, fourLocArr, newBiLi.Mul(twoVal).String())
			}
		}
	}

	//CSB0114RC0105000
	//fc174结果数据默认为对应清单中（根据fc091,fc094进行匹配）下列左边字段为“1”所对应的右边字段的结果值的合计（清单中的字段均为循环字段，为空、为A或包含?的字段，不相加）
	//比例  自费金额
	iCodeArr = [][]string{
		{"fc140", "fc141"},
		{"fc144", "fc145"},
		{"fc148", "fc149"},
		{"fc152", "fc153"},
		{"fc156", "fc157"},
		{"fc160", "fc161"},
		{"fc164", "fc165"},
		{"fc168", "fc169"},
	}
	for _, fc174Loc := range fieldLocationMap["fc174"] {
		totalRight := decimal.Zero
		for _, iCode := range iCodeArr {
			fc140LocArr := utils.GetFieldLoc(fieldLocationMap[iCode[0]], fc174Loc[0], -1, -1, -1)
			for _, fc140Loc := range fc140LocArr {
				fc140Val := utils.GetFieldValueByLoc(obj, fc140Loc, true)
				fc141LocArr := utils.GetFieldLoc(fieldLocationMap[iCode[1]], fc140Loc[0], fc140Loc[1], fc140Loc[2], fc140Loc[4])
				if len(fc141LocArr) == 1 && fc140Val == "1" {
					fc141Val := utils.GetFieldDecimalValueByLoc(obj, fc141LocArr[0], true)
					totalRight = totalRight.Add(fc141Val)
				}
			}
		}
		utils.SetOnlyOneFinalValue(obj, fc174Loc, totalRight.String())
	}

	//CSB0114RC0121000
	//fc170结果数据默认为所有清单中下列左边字段为“1”所对应的右边字段的结果值的合计（清单中的字段均为循环字段，为空、为A或包含?的字段，不相加）
	//比例  自费金额
	iCodeArr = [][]string{
		{"fc140", "fc141"},
		{"fc144", "fc145"},
		{"fc148", "fc149"},
		{"fc152", "fc153"},
		{"fc156", "fc157"},
		{"fc160", "fc161"},
		{"fc164", "fc165"},
		{"fc168", "fc169"},
	}
	totalRight := decimal.Zero
	for _, iCode := range iCodeArr {
		for _, fc140Loc := range fieldLocationMap[iCode[0]] {
			fc140Val := utils.GetFieldValueByLoc(obj, fc140Loc, true)
			fc141LocArr := utils.GetFieldLoc(fieldLocationMap[iCode[1]], fc140Loc[0], fc140Loc[1], fc140Loc[2], fc140Loc[4])
			if len(fc141LocArr) == 1 && fc140Val == "1" {
				fc141Val := utils.GetFieldDecimalValueByLoc(obj, fc141LocArr[0], true)
				totalRight = totalRight.Add(fc141Val)
			}
		}
	}
	utils.SetFinalValue(obj, fieldLocationMap["fc170"], totalRight.String())

	//CSB0114RC0122000
	//fc171结果数据默认为所有清单中fc140,fc144,fc148,fc152,fc156,fc160,fc164,fc168的结果值为“1”的字段的数量（清单中的字段均为循环字段）
	//（例：fc140为1，fc144为0，fc148/fc152/fc156/fc160/fc164/fc168为空，则fc171默认为1）
	//如没有fc140,fc144,fc148,fc152,fc156,fc160,fc164,fc168或全部内容均不为1时，fc171默认为0"
	num := 0
	codeArr := []string{"fc140", "fc144", "fc148", "fc152", "fc156", "fc160", "fc164", "fc168"}
	for _, code := range codeArr {
		for _, codeLoc := range fieldLocationMap[code] {
			fc140Val := utils.GetFieldValueByLoc(obj, codeLoc, true)
			if fc140Val == "1" {
				num++
			}
		}
	}
	utils.SetFinalValue(obj, fieldLocationMap["fc171"], strconv.Itoa(num))

	//CSB0114RC0128000
	//1、当同一属性的发票fc179的录入值包含5时，将fc273的结果值赋值给fc256
	//2、当同一属性的发票fc179的录入值包含4，不包含5时，将fc174的结果值赋值给fc256
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		fc256LocArr := utils.GetFieldLoc(fieldLocationMap["fc256"], i, -1, -1, -1)
		fc256Val := ""
		if len(invoiceMap.BaoXiaoDan) > 0 {
			fc273LocArr := utils.GetFieldLoc(fieldLocationMap["fc273"], i, -1, -1, -1)
			if len(fc273LocArr) == 1 {
				fc256Val = utils.GetFieldValueByLoc(obj, fc273LocArr[0], true)
			}
		}
		if len(invoiceMap.QingDan) > 0 && len(invoiceMap.BaoXiaoDan) == 0 {
			fc174LocArr := utils.GetFieldLoc(fieldLocationMap["fc174"], i, -1, -1, -1)
			if len(fc174LocArr) == 1 {
				fc256Val = utils.GetFieldValueByLoc(obj, fc174LocArr[0], true)
			}
		}
		utils.SetFinalValue(obj, fc256LocArr, fc256Val)
	}

	//CSB0114RC0129000
	//当fc179的录入值包含5时，将下列第一列字段的结果值赋值给第二列
	iCodeArr = [][]string{
		{"fc265", "fc090"},
		{"fc266", "fc103"},
		{"fc267", "fc109"},
		{"fc268", "fc126"},
		{"fc269", "fc115"},
		{"fc270", "fc128"},
		{"fc271", "fc116"},
		{"fc272", "fc117"},
		{"fc273", "fc118"},
		{"fc274", "fc119"},
		{"fc275", "fc120"},
		{"fc276", "fc121"},
		{"fc277", "fc122"},
		{"fc278", "fc123"},
		{"fc279", "fc124"},
		{"fc282", "fc125"},
		{"fc280", "fc129"},
		{"fc281", "fc114"},
	}
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		if len(invoiceMap.BaoXiaoDan) > 0 {
			for _, iCode := range iCodeArr {
				oneLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[0]], i, -1, -1, -1)
				twoLocArr := utils.GetFieldLoc(fieldLocationMap[iCode[1]], i, -1, -1, -1)
				if len(oneLocArr) == 1 {
					oneVal := utils.GetFieldValueByLoc(obj, oneLocArr[0], true)
					utils.SetFinalValue(obj, twoLocArr, oneVal)
				}
			}
		}
	}

	//CSB0114RC0130000	同一张发票，当fc109的录入值为A时，fc103的结果值默认为01
	for _, fc109Loc := range fieldLocationMap["fc109"] {
		fc109Val := utils.GetFieldValueByLoc(obj, fc109Loc, false)
		if fc109Val == "A" {
			fc103Loc := utils.GetFieldLoc(fieldLocationMap["fc103"], fc109Loc[0], -1, -1, -1)
			utils.SetFinalValue(obj, fc103Loc, "01")
		}
	}

	//CSB0114RC0131000	同一张发票，当fc109录入值不为A也不为空，且fc179录入值不包含5时，fc103的结果值默认为01
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		if len(invoiceMap.BaoXiaoDan) == 0 {
			fc109Loc := utils.GetFieldLoc(fieldLocationMap["fc109"], i, -1, -1, -1)
			fc103Loc := utils.GetFieldLoc(fieldLocationMap["fc103"], i, -1, -1, -1)
			if len(fc109Loc) == 1 {
				fc109Val := utils.GetFieldValueByLoc(obj, fc109Loc[0], false)
				if fc109Val != "" && fc109Val != "A" {
					utils.SetFinalValue(obj, fc103Loc, "01")
				}
			}
		}
	}

	//CSB0114RC0106000
	//1、当fc179的录入值包含5时，fc133的结果值等于同一账单中的fc268、fc270、fc280、fc269、fc256结果数据相加的值。
	//同一账单中fc268、fc270、fc280、fc269可能有多组，这四个字段分别只将第一个结果值不为“空或0.00或包含?”的结果值进行相加；如果只有0.00或者不存在时，则将0.00进行相加。
	//（根据fc094与fc092进行匹配判断fc133与fc268、fc270、fc280、fc269、fc256是否为同一账单）（该需求放在fc268,fc270,fc280,fc269，所有校验的最后面）（为空、为A或包含?的字段，不相加）
	//2、当fc179的录入值不包含5时，fc133的结果值等于同一账单中的fc126、fc128、fc129、fc115、fc256结果数据相加的值。
	//同一账单中fc126、fc128、fc129、fc115可能有多组，这四个字段分别只将第一个结果值不为“空或0.00或包含?”的结果值进行相加；如果只有0.00或者不存在时，则将0.00进行相加。
	//（根据fc094与fc092进行匹配判断fc133与fc126、fc128、fc129、fc115、fc256是否为同一账单）（该需求放在fc126,fc128,fc129,fc115，所有校验的最后面）（为空、为A或包含?的字段，不相加）
	for i := 0; i < len(obj.Invoice); i++ {
		invoiceMap := obj.Invoice[i]
		fc133Val := decimal.Zero
		codeStr := []string{"fc126", "fc128", "fc129", "fc115", "fc256"}
		if len(invoiceMap.BaoXiaoDan) != 0 {
			codeStr = []string{"fc268", "fc270", "fc280", "fc269", "fc256"}
		}
		for _, code := range codeStr {
			codeLocArr := utils.GetFieldLoc(fieldLocationMap[code], i, -1, -1, -1)
			valArr := utils.GetFieldValueArrByLocArr(obj, codeLocArr, true)
			for _, val := range valArr {
				if val != "" && val != "A" && strings.Index(val, "?") == -1 {
					fc133Val = fc133Val.Add(utils.ParseDecimal(val))
				}
			}
		}
		fc133Loc := utils.GetFieldLoc(fieldLocationMap["fc133"], i, -1, -1, -1)
		utils.SetFinalValue(obj, fc133Loc, fc133Val.String())
	}

	//CSB0114RC0107000	fc133的结果值大于fc132的结果值时，fc133出问题件：问题件代码07，问题件描述：核减金额不大于费用金额。
	//CSB0114RC0108000	fc134的结果值等于fc132减去fc133的值。
	for _, fc132Loc := range fieldLocationMap["fc132"] {
		fc133LocArr := utils.GetFieldLoc(fieldLocationMap["fc133"], fc132Loc[0], -1, -1, -1)
		fc134LocArr := utils.GetFieldLoc(fieldLocationMap["fc134"], fc132Loc[0], -1, -1, -1)
		fc132Val := utils.GetFieldDecimalValueByLoc(obj, fc132Loc, true)
		if len(fc133LocArr) == 1 {
			fc133Val := utils.GetFieldDecimalValueByLoc(obj, fc133LocArr[0], true)
			if fc133Val.GreaterThan(fc132Val) {
				utils.SetIssues(obj, fc133LocArr, "核减金额不大于费用金额", "07", "")
			}
			utils.SetFinalValue(obj, fc134LocArr, fc132Val.Sub(fc133Val).String())
		}
	}

	//CSB0114RC0109000	fc135的结果值等于所有fc132的录入值相加的值。（fc132存在多个，为空、为A或包含?的字段，不相加）
	fc132ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc132"], false)
	fc132Total := decimal.Zero
	for _, fc132Val := range fc132ValArr {
		if fc132Val != "" && fc132Val != "A" && strings.Index(fc132Val, "?") == -1 {
			fc132DecimalVal := utils.ParseDecimal(fc132Val)
			fc132Total = fc132Total.Add(fc132DecimalVal)
		}
	}
	utils.SetFinalValue(obj, fieldLocationMap["fc135"], fc132Total.String())

	//CSB0114RC0110000	fc136的结果值等于所有fc133与fc170相加的结果值合计的值（此代码放在fc133,fc170所有校验的后面，为空、为A或包含?的字段，不相加）
	fc170ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc170"], true)
	fc133ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc133"], true)
	fc133Fc170Total := decimal.Zero
	fc170DecimalVal := decimal.Zero
	for _, fc133Val := range fc133ValArr {
		if fc133Val != "" && fc133Val != "A" && strings.Index(fc133Val, "?") == -1 {
			fc133DecimalVal := utils.ParseDecimal(fc133Val)
			fc133Fc170Total = fc133Fc170Total.Add(fc133DecimalVal)
		}
	}
	if len(fc170ValArr) == 1 {
		if fc170ValArr[0] != "" && fc170ValArr[0] != "A" && strings.Index(fc170ValArr[0], "?") == -1 {
			fc170DecimalVal = utils.ParseDecimal(fc170ValArr[0])
		}
	}
	utils.SetFinalValue(obj, fieldLocationMap["fc136"], fc170DecimalVal.Add(fc133Fc170Total).String())

	//CSB0114RC0111000	fc137的结果值等于fc135减去fc136的值。(该代码放在上面两个需求的代码后面)
	if len(fieldLocationMap["fc135"]) == 1 && len(fieldLocationMap["fc136"]) == 1 && len(fieldLocationMap["fc137"]) == 1 {
		fc135Val := utils.GetFieldDecimalValueByLoc(obj, fieldLocationMap["fc135"][0], true)
		fc136Val := utils.GetFieldDecimalValueByLoc(obj, fieldLocationMap["fc136"][0], true)
		utils.SetFinalValue(obj, fieldLocationMap["fc137"], fc135Val.Sub(fc136Val).String())
	}

	//CSB0114RC0125000	第九列字段一二码的录入值不一致时将录入值转换为“?”
	//项目名称-清单 金额-清单 自费比例-清单 自费金额-清单 费用-清单 数量-清单 抢救室 项目大类-清单 项目类型
	//fc138 fc139 fc140 fc141 fc231 fc247 fc207 fc215 fc257
	//fc142 fc143 fc144 fc145 fc232 fc248 fc208 fc216 fc258
	//fc146 fc147 fc148 fc149 fc233 fc249 fc209 fc217 fc259
	//fc150 fc151 fc152 fc153 fc234 fc250 fc210 fc218 fc260
	//fc154 fc155 fc156 fc157 fc235 fc251 fc211 fc219 fc261
	//fc158 fc159 fc160 fc161 fc236 fc252 fc212 fc220 fc262
	//fc162 fc163 fc164 fc165 fc237 fc253 fc213 fc221 fc263
	//fc166 fc167 fc168 fc169 fc238 fc254 fc214 fc222 fc264
	//codeOneStr := "fc257,fc258,fc259,fc260,fc261,fc262,fc263,fc264,"
	//CSB0114RC0127000
	//fc063，fc065，fc067，fc069，fc071，fc073，fc075，fc077，fc079，fc081，fc083，fc085，fc087，fc089，fc114,
	//fc115,fc116,fc117，fc118,fc119,fc120,fc121,fc122,fc123,fc124,fc125,fc126,fc127,fc128,fc129,fc131,
	//fc132，fc133,fc134,fc135,fc136,fc137,fc170，fc174，fc231，fc232，fc233，fc234，fc235，fc236，fc237，
	//fc238，fc139，fc143，fc147，fc151，fc155，fc159，fc163，fc167，fc256，fc281，fc269，fc271，fc272，fc273，
	//fc274，fc275，fc276，fc277，fc278，fc279，fc282，fc280，fc268，fc270结果数据保留两位小数（如结果为空时，则默认为0.00）
	codeTwoStr := "fc063,fc065,fc067,fc069,fc071,fc073,fc075,fc077,fc079,fc081,fc083,fc085,fc087,fc089,fc114,fc115,fc116,fc117,fc118,fc119,fc120,fc121,fc122,fc123,fc124,fc125,fc126,fc127,fc128,fc129,fc131,fc132,fc133,fc134,fc135,fc136,fc137,fc170,fc174,fc231,fc232,fc233,fc234,fc235,fc236,fc237,fc238,fc139,fc143,fc147,fc151,fc155,fc159,fc163,fc167,fc256,fc281,fc269,fc271,fc272,fc273,fc274,fc275,fc276,fc277,fc278,fc279,fc282,fc280,fc268,fc270,"
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
						//i:发票index j:发票结构体字段index k:字段二维数组x的index l:字段二维数组y的index
						//if strings.Index(codeOneStr, fields[l].Code+",") != -1 {
						//	if fields[l].Op1Value != fields[l].Op2Value {
						//		fields[l].ResultValue = "?"
						//	}
						//}

						if strings.Index(codeTwoStr, fields[l].Code+",") != -1 {
							decimalFinalValue := utils.ParseDecimal(fields[l].FinalValue).Round(2)
							fields[l].FinalValue = decimalFinalValue.StringFixed(2)
						}
					}
				}
			}
		}
	}

	//CSB0114RC0126000
	//除了字段“fc138、fc142、fc146、fc150、fc154、fc158、fc162、fc166、fc231、fc232、fc233、fc234、fc235、fc236、fc237、fc238”，
	//fc257、fc258、fc259、fc260、fc261、fc262、fc263、fc264、fc130、fc062、fc064、fc066、fc068、fc070、fc072、fc074、fc076、fc078、fc080、fc082、fc084、fc086、fc088、fc176
	//其他所有字段结果数据包含?时，将结果数据转换为空，该字段出问题件：问题件代码01，问题件描述：（字段名）填写模糊。
	//(将该代码放在所有需求的最后面。)"
	codeThreeStr := "fc138,fc142,fc146,fc150,fc154,fc158,fc162,fc166,fc231,fc232,fc233,fc234,fc235,fc236,fc237,fc238,fc257,fc258,fc259,fc260,fc261,fc262,fc263,fc264,fc130,fc062,fc064,fc066,fc068,fc070,fc072,fc074,fc076,fc078,fc080,fc082,fc084,fc086,fc088,fc176,"
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
						//i:发票index j:发票结构体字段index k:字段二维数组x的index l:字段二维数组y的index
						if strings.Index(codeThreeStr, fields[l].Code+",") == -1 &&
							strings.Index(fields[l].FinalValue, "?") != -1 {
							fields[l].FinalValue = ""
							utils.SetIssue(obj, []int{i, j, k, l}, fields[l].Name+"填写模糊", "01", "")
						}
					}
				}
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
		{"diZhiMap", "B0114_华夏理赔_华夏理赔地址库", "4", "1"},
		{"diZhiShengShiQuMap", "B0114_华夏理赔_华夏理赔地址库", "4", "1"},
		{"yiWaiYuanYinMap", "B0114_华夏理赔_华夏理赔意外原因表", "0", "1"},
		{"sunShangWaiBuYuanYinMap", "B0114_华夏理赔_华夏理赔损伤外部原因表", "0", "1"},
		{"shouShuBianMaMap", "B0114_华夏理赔_手术术式编码表", "0", "1"},
		{"jiBingDaiMaMap", "B0114_华夏理赔_ICD疾病代码表", "2", "1"},
		{"zhongJiQingJiMingChengMap", "B0114_华夏理赔_重疾和轻症疾病名称表", "0", "1"},
		{"quanCanMap", "B0114_华夏理赔_全残信息表", "3", "2"},
		{"jiBingDaiMaTypeMap", "B0114_华夏理赔_ICD疾病代码表", "2", "0"},
		{"yiYuanMingChengMap", "B0114_华夏理赔_医院名称表", "2", "1"},
		{"shanCanXiangMuMap", "B0114_华夏理赔_伤残项目表", "2", "1"},
		{"shanCanXiangMuBiLiMap", "B0114_华夏理赔_伤残项目表", "2", "5"},
		{"quanCanBiLiMap", "B0114_华夏理赔_全残信息表", "3", "0"},
		{"jiGouBianMaMap", "B0114_华夏理赔_机构编码对应表", "1", "0"},
	}
	for i, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		tempNumMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])

				//省市区匹配唯一CSB0114RC0076000
				if i == 1 {
					tempMap[arr[2]+"_"+arr[3]+"_"+arr[4]] = arr[v]
					tempMap[arr[3]+"_"+arr[4]] = arr[v]
					tempMap[arr[2]+"_"+arr[4]] = arr[v]

					//同一个key有多少个值
					num := 1
					numStr, has := tempNumMap[strings.TrimSpace(arr[k])]
					if has {
						num, _ = strconv.Atoi(numStr)
						num++
					}
					tempNumMap[strings.TrimSpace(arr[k])] = fmt.Sprintf("%d", num)
				}
				//else {
				tempMap[strings.TrimSpace(arr[k])] = arr[v]
				//}
			}
		}
		constObj[item[0]] = tempMap
		constObj[item[0]+"Num"] = tempNumMap
	}

	//地址
	for k, v := range global.GProConf[proCode].ConstTable {
		if strings.Index(k, "省份") != -1 {
			tempMap := make(map[string]string, 0)
			for _, arr := range v {
				if len(arr) < 3 {
					continue
				}
				//获取最小的一项
				if val, ok := tempMap[arr[2]]; ok {
					a, err := strconv.ParseFloat(val, 64)
					if err != nil {
						continue
					}
					b, err := strconv.ParseFloat(arr[0], 64)
					if err != nil {
						continue
					}
					if a < b {
						arr[0] = val
					}
				}
				tempMap[arr[2]] = arr[0]
			}
			constObj[strings.Replace(k, "B0114_华夏理赔_省份-", "", -1)] = tempMap
		}
	}

	return constObj
}

// specialConst 特别的常量 需要判断是否唯一的
func specialConst(proCode string) map[string]map[string][]string {
	constObj := make(map[string]map[string][]string, 0)
	nameMap := [][]string{
		{"diZhiMap", "B0114_华夏理赔_华夏理赔地址库", "4", "0"},
		{"diZhiShiMap", "B0114_华夏理赔_华夏理赔地址库", "3", "0"},
		{"diZhiShengMap", "B0114_华夏理赔_华夏理赔地址库", "2", "0"},
	}
	for _, item := range nameMap {
		tempMap := make(map[string][]string, 0)
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])
				tempMap[strings.TrimSpace(arr[k])] = append(tempMap[strings.TrimSpace(arr[k])], arr[v])
			}
		} else {
			global.GLog.Error(item[1], zap.Error(errors.New("没有该常量")))
		}
		constObj[item[0]] = tempMap
	}
	return constObj
}
