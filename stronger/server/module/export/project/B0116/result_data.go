package B0116

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
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
)

//B0116华夏人寿团险理赔
//模板类型字段 : fc086
//账单号: fc045发票号
//总金额: fc046账单金额
//发票类型: fc111账单类型
//
//和发票有对应关系的
//fc088发票属性 MB002-bc010
//fc089清单所属发票 MB002-bc011
//fc090报销单所属发票 MB002-bc012
//
//和发票没有对应关系的
//2-申请书2（有被保险人）;
//6-伤病诊断;
//8-开户人身份证;
//9-手术术士信息;
//11-重疾/轻症;
//12-伤残/全残;
//13-特种病;
//14-失能;
//15-护理;

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode:    model2.TypeCode{FieldCode: []string{"fc088"}, BlockCode: []string{"bc010"}},
	QingDanFieldCode:    model2.TypeCode{FieldCode: []string{"fc089"}, BlockCode: []string{"bc011"}},
	BaoXiaoDanFieldCode: model2.TypeCode{FieldCode: []string{"fc090"}, BlockCode: []string{"bc012"}},
	OtherTempType:       map[string]string{"2": "2", "6": "6", "8": "8", "9": "9", "11": "11", "12": "12", "13": "13", "14": "14", "15": "15"},
	TempTypeField:       "fc086",
	InvoiceNumField:     []string{"fc045"},
	MoneyField:          []string{"fc046"},
	InvoiceTypeField:    "fc111",
}

// ResultData B0116
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {
	defer func() {
		if err := recover(); err != nil {
			global.GLog.Error("", zap.Any("", err))
			global.GLog.Error(string(debug.Stack()))
		}
	}()
	obj = utils2.RelationDealManyInvoiceType(bill, blocks, fieldMap, relationConf)

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

	//CSB0116RC0045000
	//"当fc008的录入值为“1”时，fc009的结果值默认为“1”；
	//当fc008的录入值为“A”时，fc009的结果值默认为“0”
	//（该需求放在需求“当fc008和fc009的结果值都为空时，将fc193的结果值赋值到fc009的结果值中”前面）"
	for _, loc := range fieldLocationMap["fc008"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		fc009Locs := utils.GetFieldLoc(fieldLocationMap["fc009"], loc[0], -1, -1, -1)
		iMap := map[string]string{
			"1": "1",
			"A": "0",
		}
		utils.SetFinalValue(obj, fc009Locs, iMap[val])
	}

	//CSB0116RC0046000
	//当fc008和fc009的结果值都为空时，将fc193的结果值赋值到fc009的结果值中
	fc008Val := utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc008"], true)
	fc009Val := utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc009"], true)
	fc193Val := utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc193"], true)
	if fc009Val == fc008Val && fc008Val == "" {
		utils.SetFinalValue(obj, fieldLocationMap["fc009"], fc193Val)
	}

	//CSB0116RC0047000
	//fc010的结果值默认为所有fc088录入值不重复的数量，fc010的结果值为空时转换为0（fc088录入值为空的不计算在内，如一个案件有三个fc088，录入值分别为：空、1、2，则fc010的结果值默认为“2”；不存在fc088或者fc088所有录入值为空时，fc010的结果值默认为“0”）
	//发票数量
	utils.SetFinalValue(obj, fieldLocationMap["fc010"], strconv.Itoa(len(obj.Invoice)-1))

	//CSB0116RC0048000
	//fc011的结果值等于所有fc046结果值的和
	valArr := utils.GetFieldDecimalValueArrByLocArr(obj, fieldLocationMap["fc046"], true)
	total := decimal.Zero
	for _, d := range valArr {
		total = total.Add(d)
	}
	utils.SetFinalValue(obj, fieldLocationMap["fc011"], total.StringFixed(2))

	//CSB0116RC0049000
	//"当fc204录入值为“1”时，将下列左边字段第一个不为空的结果值赋值到右边所有字段的结果值中：
	//（该需求放在需求“当fc008和fc009的结果值都为空时，将fc193的结果值赋值到fc009的结果值中”后面）"
	codesArr := [][]string{
		{"fc184", "fc012"},
		{"fc185", "fc003"},
		{"fc186", "fc006"},
		{"fc187", "fc007"},
		{"fc188", "fc005"},
		{"fc189", "fc004"},
		{"fc192", "fc008"},
		{"fc021", "fc196"},
		{"fc022", "fc195"},
		{"fc023", "fc194"},
		{"fc024", "fc197"},
	}
	if len(fieldLocationMap["fc204"]) == 1 {
		fc204Val := utils.GetFieldValueByLoc(obj, fieldLocationMap["fc204"][0], false)
		if fc204Val == "1" {
			for _, codes := range codesArr {
				for _, loc0 := range fieldLocationMap[codes[0]] {
					loc1 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc0[0], -1, -1, -1)
					val0 := utils.GetFieldValueByLoc(obj, loc0, true)
					if val0 != "" {
						utils.SetFinalValue(obj, loc1, val0)
					}
				}
			}
		}
	}

	//CSB0116RC0050000
	//"1.当fc186录入值为“0”或“5”且fc187的结果值不为空时，将fc187结果数据中第七位至第十四位的数字以YYYY-MM-DD格式赋值到fc188的结果值中；
	// 2.当fc006录入值为“0”或“5”且fc007的结果值不为空时，将fc007结果数据中第七位至第十四位的数字以YYYY-MM-DD格式赋值到fc005的结果值中"
	codesArr = [][]string{
		{"fc186", "fc187", "fc188"},
		{"fc006", "fc007", "fc005"},
	}
	for _, codes := range codesArr {
		val0 := utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[0]], false)
		val1 := utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[1]], true)
		if (val0 == "0" || val0 == "5") && val1 != "" && len(val1) >= 14 {
			utils.SetFinalValue(obj, fieldLocationMap[codes[2]], val1[6:10]+"-"+val1[10:12]+"-"+val1[12:14])
		}
	}
	//CSB0116RC0051000
	//当fc192录入值为“1”时，将fc193的结果值默认为“1”
	fc192Val := utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc192"], false)
	if fc192Val == "1" {
		utils.SetFinalValue(obj, fieldLocationMap["fc193"], "1")
	}

	//CSB0116RC0052000
	//"当下列左边字段的录入值不为“A”或“空”时，右边对应字段的结果值默认为“01”
	codesArr = [][]string{
		{"fc044", "fc043"},
		{"fc096", "fc097"},
		{"fc099", "fc100"},
	}
	for _, codes := range codesArr {
		val0 := utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[0]], false)
		if val0 != "A" && val0 != "" {
			utils.SetFinalValue(obj, fieldLocationMap[codes[1]], "01")
		}
	}

	//CSB0116RC0053000
	//"当下列左边字段结果值不为空时，右边对应的字段结果值默认为“06”（如fc102结果值为XXX，则fc103结果值默认为06）
	codesArr = [][]string{
		{"fc041", "fc040"},
		{"fc102", "fc103"},
		{"fc105", "fc106"},
	}
	for _, codes := range codesArr {
		val0 := utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[0]], true)
		if val0 != "" {
			utils.SetFinalValue(obj, fieldLocationMap[codes[1]], "06")
		}
	}

	//CSB0116RC0054000
	//"以下为清单循环字段，
	//1.当第三列字段的录入值为“3”或“4”时，将对应的第二列字段的结果值默认为“1”；
	//2.当第三列字段的录入值不为“3”且不为“4”时，根据fc021（整单唯一）的录入值匹配医疗目录常量表（fc021的录入值为空或为A或包含?时不进行匹配），
	//	如fc021录入值为“北京”，则匹配的医疗目录常量表为《B0116_华夏人寿团险理赔_省份-北京》，
	//	当匹配得出医疗目录常量表时，下列第一列字段的录入值匹配到医疗目录中的“中文名称”（第三列），将与之对应的“自付比例”（第一列）放到第二列对应的字段的结果值中；
	//	匹配规则：
	//	a.无法匹配的情况第二列对应的结果值默认为0（无法匹配有两种情况：1.fc021无法匹配得出医疗目录；2.匹配到医疗目录但上面第一列字段录入内容不在医疗目录内）
	//	b.当上面第一列字段的录入值匹配到的“中文名称”（第三列）不唯一时，取对应的“自付比例”（第一列）最小的值放到第二列对应的字段的结果值中"
	//项目名称 自费比例 项目类型
	codesArr = [][]string{
		{"fc131", "fc147", "fc175"},
		{"fc132", "fc148", "fc176"},
		{"fc133", "fc149", "fc177"},
		{"fc134", "fc150", "fc178"},
		{"fc135", "fc151", "fc179"},
		{"fc136", "fc152", "fc180"},
		{"fc137", "fc153", "fc181"},
		{"fc138", "fc154", "fc182"},
	}
	fc021Val := utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc021"], false)
	for _, codes := range codesArr {
		for _, loc3 := range fieldLocationMap[codes[2]] {
			val3 := utils.GetFieldValueByLoc(obj, loc3, false)
			loc2 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc3[0], loc3[1], loc3[2], loc3[4])
			loc1 := utils.GetFieldLoc(fieldLocationMap[codes[0]], loc3[0], loc3[1], loc3[2], loc3[4])
			if len(loc2) != 1 || len(loc1) != 1 {
				continue
			}
			if m, _ := regexp.MatchString(`^(3|4)$`, val3); m {
				utils.SetFinalValue(obj, loc2, "1")
			} else {
				if m, _ := regexp.MatchString(`^(A|)$`, fc021Val); m || strings.Index(fc021Val, "?") != -1 {
					continue
				}
				val1 := utils.GetFieldValueByLocArr(obj, loc1, false)
				value, row := utils.FetchConst(bill.ProCode, "B0116_华夏人寿团险理赔_省份-"+fc021Val, "自付比例", map[string]string{"中文名称": val1})
				if row == 0 {
					utils.SetFinalValue(obj, loc2, "0")
				} else {
					utils.SetFinalValue(obj, loc2, value)
				}
			}
		}
	}

	//CSB0116RC0055000
	//"同一发票下，以下为清单循环字段，将所有下列右边字段的结果值为“1”所对应的左边字段的结果值合计放到fc051（发票唯一）的结果值中：
	//项目金额 项目类型
	codesArr = [][]string{
		{"fc139", "fc147"},
		{"fc140", "fc148"},
		{"fc141", "fc149"},
		{"fc142", "fc150"},
		{"fc143", "fc151"},
		{"fc144", "fc152"},
		{"fc145", "fc153"},
		{"fc146", "fc154"},
	}
	for _, fc051Loc := range fieldLocationMap["fc051"] {
		total = decimal.Zero
		for _, codes := range codesArr {
			loc1Arr := utils.GetFieldLoc(fieldLocationMap[codes[1]], fc051Loc[0], -1, -1, -1)
			for _, loc1 := range loc1Arr {
				val1 := utils.GetFieldValueByLoc(obj, loc1, true)
				loc0 := utils.GetFieldLoc(fieldLocationMap[codes[0]], loc1[0], loc1[1], loc1[2], loc1[4])
				if val1 == "1" && len(loc0) == 1 {
					total = total.Add(utils.GetFieldDecimalValueByLocArr(obj, loc0, true))
				}
			}
		}
		utils.SetOnlyOneFinalValue(obj, fc051Loc, total.StringFixed(2))
	}

	//CSB0116RC0056000
	//同一发票下，当fc202结果值为非零的数字时，将fc202的结果值赋值到fc051的结果值中
	for _, loc := range fieldLocationMap["fc202"] {
		fc051Loc := utils.GetFieldLoc(fieldLocationMap["fc051"], loc[0], -1, -1, -1)
		fc202Val := utils.GetFieldDecimalValueByLoc(obj, loc, true)
		if !fc202Val.IsZero() {
			utils.SetFinalValue(obj, fc051Loc, fc202Val.StringFixed(2))
		}
	}

	//CSB0116RC0057000
	//同一发票下，fc057的结果值默认为fc056的结果日期减去fc055的结果日期的值
	//（如果fc055和fc056的结果值均为空，则fc057的结果值为“0”；如果fc055和fc056的结果值不为空且相等，则fc057的结果值为“0”）
	for i, fc057Loc := range fieldLocationMap["fc057"] {
		fc056Loc := utils.GetFieldLoc(fieldLocationMap["fc056"], fc057Loc[0], -1, -1, -1)
		fc055Loc := utils.GetFieldLoc(fieldLocationMap["fc055"], fc057Loc[0], -1, -1, -1)
		if len(fc055Loc) != 1 || len(fc056Loc) != 1 {
			continue
		}
		fc056Val := strings.ReplaceAll(utils.GetFieldValueByLoc(obj, fc056Loc[0], true), "-", "")
		fc055Val := strings.ReplaceAll(utils.GetFieldValueByLoc(obj, fc055Loc[0], true), "-", "")
		if fc055Val == fc056Val {
			utils.SetOnlyOneFinalValue(obj, fc057Loc, "0")
		}
		global.GLog.Info(strconv.Itoa(i), zap.Any("fc056Val", fc056Val))
		global.GLog.Info(strconv.Itoa(i), zap.Any("fc055Val", fc055Val))
		utils.SetOnlyOneFinalValue(obj, fc057Loc, utils.ParseDecimal(fc055Val).Sub(utils.ParseDecimal(fc055Val)).StringFixed(0))
	}

	//CSB0116RC0058000
	//"同一发票下，
	//匹配规则，当第一列字段的值为第三列内容时，将第一列对应第二列字段的结果值放在第三列所对应的第四列字段中
	//（如果第一列字段录入值存在重复时，将重复所对应的第二列字段的结果值求和后再根据规则放到下列对应的字段中）：
	//如fc115值为0，fc123值为100，fc116值为0，fc124值为200，则fc155结果值为300"
	//①   ②   ③   ④
	myMap := map[string]string{
		"0":  "fc155",
		"1":  "fc156",
		"2":  "fc157",
		"3":  "fc158",
		"4":  "fc159",
		"5":  "fc160",
		"6":  "fc161",
		"7":  "fc162",
		"8":  "fc163",
		"9":  "fc164",
		"10": "fc165",
		"11": "fc166",
		"12": "fc167",
		"13": "fc168",
		"14": "fc169",
		"15": "fc170",
		"16": "fc171",
		"17": "fc172",
		"18": "fc173",
		"19": "fc174",
	}
	codesArr = [][]string{
		{"fc115", "fc123"},
		{"fc116", "fc124"},
		{"fc117", "fc125"},
		{"fc118", "fc126"},
		{"fc119", "fc127"},
		{"fc120", "fc128"},
		{"fc121", "fc129"},
		{"fc122", "fc130"},
		{"fc222", "fc234"},
		{"fc223", "fc235"},
		{"fc224", "fc236"},
		{"fc225", "fc237"},
		{"fc226", "fc238"},
		{"fc227", "fc239"},
		{"fc228", "fc240"},
		{"fc229", "fc241"},
		{"fc230", "fc242"},
		{"fc231", "fc243"},
		{"fc232", "fc244"},
		{"fc233", "fc245"},
	}

	for _, codes := range codesArr {
		for _, loc1 := range fieldLocationMap[codes[0]] {
			val1 := utils.GetFieldValueByLoc(obj, loc1, false)
			if fc155Code, ok := myMap[val1]; ok {
				loc2 := utils.GetFieldLoc(fieldLocationMap[codes[1]], loc1[0], -1, -1, -1)
				fc155Loc := utils.GetFieldLoc(fieldLocationMap[fc155Code], loc1[0], -1, -1, -1)
				val2 := utils.GetFieldDecimalValueByLocArr(obj, loc2, true)
				fc155Val := utils.GetFieldDecimalValueByLocArr(obj, fc155Loc, true)
				utils.SetFinalValue(obj, fc155Loc, val2.Add(fc155Val).StringFixed(2))
			}
		}
	}

	//CSB0116RC0059000
	//fc026的结果值根据fc027录入内容匹配《B0116_华夏人寿团险理赔_残疾》的“残疾中文名”（第四列）匹配出对应的“残疾类型编码”（第一列）进行输出
	//CSB0116RC0060000
	//fc028的结果值根据fc027录入内容匹配《B0116_华夏人寿团险理赔_残疾》的“残疾中文名”（第四列）匹配出对应的“伤残级别编码”（第二列）进行输出
	//CSB0116RC0061000
	//fc030的结果值根据fc027录入内容匹配《B0116_华夏人寿团险理赔_残疾》的“残疾中文名”（第四列）匹配出对应的“伤残赔付比例”（第五列）进行输出
	//CSB0116RC0062000
	//fc033的结果值根据fc108录入内容匹配《B0116_华夏人寿团险理赔_重疾》的“中文名”（第三列）匹配出对应的“重疾类型编码”（第二列）进行输出
	codesArr = [][]string{
		{"fc026", "fc027", "B0116_华夏人寿团险理赔_残疾", "残疾中文名", "残疾类型编码"},
		{"fc028", "fc027", "B0116_华夏人寿团险理赔_残疾", "残疾中文名", "伤残级别编码"},
		{"fc030", "fc027", "B0116_华夏人寿团险理赔_残疾", "残疾中文名", "伤残赔付比例"},
		{"fc033", "fc108", "B0116_华夏人寿团险理赔_重疾", "残疾中文名", "重疾类型编码"},
	}
	for _, codes := range codesArr {
		val := utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[1]], false)
		value, _ := utils.FetchConst(bill.ProCode, codes[2], codes[4], map[string]string{codes[3]: val})
		utils.SetFinalValue(obj, fieldLocationMap[codes[0]], value)
	}

	//CSB0116RC0063000
	//fc011、fc046、fc048、fc049、fc050、fc051、fc058、fc059、fc060、fc061、fc062、fc063、fc064、fc065、fc066、fc067、fc068、fc069、fc070、fc071、fc072、fc073、fc074、fc075、fc076、fc077、fc078、fc079、fc080、fc081、fc082、fc083、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc139、fc140、fc141、fc142、fc143、fc144、fc145、fc146、fc155、fc156、fc157、fc158、fc159、fc160、fc161、fc162、fc163、fc164、fc165、fc166、fc167、fc168、fc169、fc170、fc171、fc172、fc173、fc174、fc234、fc235、fc236、fc237、fc238、fc239、fc240、fc241、fc242、fc243、fc244、fc245结果数据保留两位小数（如结果为空时，则默认为0.00）
	codes := []string{"fc011", "fc046", "fc048", "fc049", "fc050", "fc051", "fc058", "fc059", "fc060", "fc061", "fc062", "fc063", "fc064", "fc065", "fc066", "fc067", "fc068", "fc069", "fc070", "fc071", "fc072", "fc073", "fc074", "fc075", "fc076", "fc077", "fc078", "fc079", "fc080", "fc081", "fc082", "fc083", "fc123", "fc124", "fc125", "fc126", "fc127", "fc128", "fc129", "fc130", "fc139", "fc140", "fc141", "fc142", "fc143", "fc144", "fc145", "fc146", "fc155", "fc156", "fc157", "fc158", "fc159", "fc160", "fc161", "fc162", "fc163", "fc164", "fc165", "fc166", "fc167", "fc168", "fc169", "fc170", "fc171", "fc172", "fc173", "fc174", "fc234", "fc235", "fc236", "fc237", "fc238", "fc239", "fc240", "fc241", "fc242", "fc243", "fc244", "fc245"}
	for _, code := range codes {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldDecimalValueByLoc(obj, loc, true)
			utils.SetOnlyOneFinalValue(obj, loc, val.StringFixed(2))
		}
	}

	//CSB0116RC0064000
	//当fc086不存在1且不存在2时，第一个fc086出问题件：问题件代码2，问题件描述：缺少理赔申请表（fc086为循环字段）
	fc086ValArr := utils.GetFieldValueArrByLocArr(obj, fieldLocationMap["fc086"], false)
	fc086ValStr := "-" + strings.Join(fc086ValArr, "-") + "-"
	if m, _ := regexp.MatchString(`-(1|2)-`, fc086ValStr); !m {
		utils.SetIssue(obj, fieldLocationMap["fc086"][0], "缺少理赔申请表", "2", "")
	}

	//CSB0116RC0065000
	//当fc086录入值不存在“3”且fc086的录入值不存在“11或12或13或14时（fc086为循环字段），第一个fc086出问题件：问题件代码2，问题件描述：缺少账单信息（fc013录入值为“8”或“9”时不进行该校验）
	fc013Val := utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc013"], false)
	m, _ := regexp.MatchString(`-(3)-`, fc086ValStr)
	m1, _ := regexp.MatchString(`-(11|12|13|14)-`, fc086ValStr)
	if !m && !m1 && fc013Val != "8" && fc013Val != "9" {
		utils.SetIssue(obj, fieldLocationMap["fc086"][0], "缺少账单信息", "2", "")
	}

	//CSB0116RC0066000
	//fc005结果值晚于fc017结果值时（fc005结果值为空时不进行校验），fc005出问题件：问题件代码：7，问题件描述：出生日期不能晚于出险日期！
	//CSB0116RC0067000
	//fc008结果值早于fc017结果值时（fc008结果值为空时不进行校验），fc008出问题件：问题件代码：7，问题件描述：证件有效期不能早于出险日期！
	//CSB0116RC0068000
	//fc018结果值早于fc017结果值时（fc018结果值为空时不进行校验），fc018出问题件：问题件代码：7，问题件描述：认定日期不能早于出险日期！
	//CSB0116RC0072000
	//"1.当fc192的结果值早于fc191的结果值时，fc192出问题件：问题件代码：7，问题件描述：证件有效止期不能早于证件有效起期！
	//2.当fc192的结果值早于fc017的结果值时，fc192出问题件：问题件代码：7，问题件描述：证件有效止期不能早于出险日期！"
	codesArr = [][]string{
		{"fc005", "fc017", "fc005", "出生日期不能晚于出险日期"},
		{"fc017", "fc008", "fc008", "证件有效期不能早于出险日期"},
		{"fc017", "fc018", "fc018", "认定日期不能早于出险日期"},
		{"fc191", "fc192", "fc192", "证件有效止期不能早于证件有效起期"},
		{"fc017", "fc192", "fc192", "证件有效止期不能早于出险日期"},
	}
	for _, codes = range codesArr {
		val0 := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[0]], true), "-", "")
		val1 := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[1]], true), "-", "")
		if utils.ParseDecimal(val0).GreaterThan(utils.ParseDecimal(val1)) {
			utils.SetIssues(obj, fieldLocationMap[codes[2]], codes[3], "7", "")
		}
	}

	//CSB0116RC0069000
	//fc015录入内容匹配《B0116_华夏人寿团险理赔_意外出险细节》的“中文名”（第二列）对应的“编码”（第一列）进行转码；
	//不能转码时结果值默认为空，并出问题件：问题件代码7，问题件描述：意外出险细节不在常量库中（录入内容为A、F、为空或包含?时不做校验）
	//CSB0116RC0070000
	//fc016录入内容匹配《B0116_华夏人寿团险理赔_损伤外部原因》的“中文名”（第二列）对应的“编码”（第一列）进行转码；
	//不能转码时结果值默认为空，并出问题件：问题件代码7，问题件描述：损伤外部原因不在常量库中（录入内容为A、F、为空或包含?时不做校验）
	codesArr = [][]string{
		{"fc015", "B0116_华夏人寿团险理赔_意外出险细节", "意外出险细节不在常量库中"},
		{"fc016", "B0116_华夏人寿团险理赔_损伤外部原因", "损伤外部原因不在常量库中"},
	}
	for _, arr := range codesArr {
		for _, loc := range fieldLocationMap[arr[0]] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			if m, _ := regexp.MatchString(`^(A|F|)$`, val); m || strings.Index(val, "?") != -1 {
				continue
			}
			value, row := utils.FetchConst(bill.ProCode, arr[1], "编码", map[string]string{"中文名": val})
			if row == 0 {
				utils.SetIssue(obj, loc, arr[2], "7", "")
			}
			utils.SetOnlyOneFinalValue(obj, loc, value)
		}
	}

	//CSB0116RC0071000
	//"每一行为一组，判断一组三个字段录入值在常量表《B0116_华夏人寿团险理赔_地址库》的“区、县中文名”（第六列）、“市中文名”（第四列）、“省中文名”（第二列）是否为同一行一一对应，如是则对应转码输出：
	//1.第一列字段根据“省中文名”（第二列）转“省编码”（第一列），无法转码时结果值转为空，并出问题件：问题件代码：7，问题件描述：省、市、区县编码不匹配！；
	//2.第二列字段根据“市中文名”（第四列）转“市编码”（第三列），无法转码时结果值转为空，并出问题件：问题件代码：7，问题件描述：省、市、区县编码不匹配！；
	//3.第三列字段根据“区、县中文名”（第六列）转“区、县编码”（第五列） ，无法转码时结果值转为空，并出问题件：问题件代码：7，问题件描述：省、市、区县编码不匹配！；
	//（录入内容为A、为空、为F或包含?时不做校验）
	//省     市     区
	codesArr = [][]string{
		{"fc021", "fc022", "fc023"},
		{"fc196", "fc195", "fc194"},
		{"fc212", "fc211", "fc210"},
	}
	msg := []string{"省", "市", "区、县"}
	for _, codes := range codesArr {
		for i, code := range codes {
			val := utils.GetFieldValueByLocArr(obj, fieldLocationMap[code], false)
			value, row := utils.FetchConst(bill.ProCode, "B0116_华夏人寿团险理赔_地址库", msg[i]+"编码", map[string]string{msg[i] + "中文名": val})
			if row == 0 {
				utils.SetIssues(obj, fieldLocationMap[code], "省、市、区县编码不匹配", "7", "")
			}
			utils.SetFinalValue(obj, fieldLocationMap[code], value)
		}
	}

	//CSB0116RC0073000
	//fc044，fc096，fc099录入内容匹配《B0116_华夏人寿团险理赔_疾病编码》的“中文名”（第二列）对应的“编码”（第一列）进行转码输出，
	//无法转码时对应字段的结果值默认为空，并出问题件：问题件代码7，问题件描述：疾病编码不在常量库中；（录入内容为A、为空或包含?时不做校验）
	//CSB0116RC0074000
	//fc041，fc102，fc105录入内容匹配《B0116_华夏人寿团险理赔_手术术士编码》的“中文名”（第二列）对应的“编码”（第一列）进行转码输出，
	//无法转码时对应字段结果值默认为空，并出问题件：问题件代码7，问题件描述：手术术式编码不在常量库中；（录入内容为A、为空或包含?时不做校验）
	codesArr = [][]string{
		{"fc044", "fc096", "fc099"},
		{"fc041", "fc102", "fc105"},
	}
	msgArr := [][]string{
		{"B0116_华夏人寿团险理赔_疾病编码", "疾病编码不在常量库中"},
		{"B0116_华夏人寿团险理赔_手术术士编码", "手术术式编码不在常量库中"},
	}
	for i, codes := range codesArr {
		for _, code := range codes {
			for _, loc := range fieldLocationMap[code] {
				val := utils.GetFieldValueByLoc(obj, loc, false)
				if m, _ := regexp.MatchString(`^(A|)$`, val); m || strings.Index(val, "?") != -1 {
					continue
				}
				value, row := utils.FetchConst(bill.ProCode, msgArr[i][0], "编码", map[string]string{"中文名": val})
				if row == 0 {
					utils.SetIssue(obj, loc, msgArr[i][1], "7", "")
				}
				utils.SetOnlyOneFinalValue(obj, loc, value)
			}
		}
	}

	//CSB0116RC0075000
	//fc045的录入值重复时，重复的fc045出问题件：问题件代码7，描述：发票号码不能重复！
	fc045Map := map[string]string{}
	for _, loc := range fieldLocationMap["fc045"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		if _, ok := fc045Map[val]; ok {
			utils.SetIssue(obj, loc, "发票号码不能重复", "7", "")
		} else {
			fc045Map[val] = val
		}
	}

	//CSB0116RC0076000
	//"同一发票下，当以下字段存在时，执行以下校验：
	//1.当fc048、fc049、fc050、fc051的结果值之和大于fc046的结果值时，fc046出问题件：问题件代码7，描述：其他扣减费用、新农合、医保支付金额、其他商保支付之和不得大于账单发生金额！
	//2.当fc060、fc061、fc062、fc059、fc050的结果值之和大于fc046的结果值时，fc046出问题件：问题件代码7，描述：自付一、自付二、自费、医疗保险基金支付金额、其他商保支付金额之和不得大于账单发生金额！
	//3.当fc058的结果值大于fc046的结果值时，fc046出问题件：问题件代码7，描述：医疗保险范围内金额不能大于账单发生金额！
	//4.当fc077的结果值大于fc046的结果值时，fc046出问题件：问题件代码7，描述：筹基金本次支付不能大于账单发生金额！
	//5.当fc080的结果值大于fc046的结果值时，fc046出问题件：问题件代码7，描述：大额互助资金本次支付不能大于账单发生金额！
	//6.当fc063、fc072、fc073、fc074的结果值之和大于fc059的结果值时，fc059出问题件：问题件代码7，描述：门诊大额支付、退休人员补充医疗保险支付、残疾军人补助支付、公务员医疗补助支付之和，不得大于医疗保险基金支付金额！
	//7.当fc077的结果值大于fc082的结果值时，fc082出问题件：问题件代码7，描述：统筹基金本次支付不能大于统筹基金年度内累计支付！
	//8.当fc080的结果值大于fc083的结果值时，fc083出问题件：问题件代码7，描述：大额互助资金本次支付不能大于大额互助资金累计支付金额！"
	codesArr = [][]string{
		{"fc048", "fc049", "fc050", "fc051"},
		{"fc060", "fc061", "fc062", "fc059", "fc050"},
		{"fc058"},
		{"fc077"},
		{"fc080"},
		{"fc063", "fc072", "fc073", "fc074"},
		{"fc077"},
		{"fc080"},
	}
	msgArr = [][]string{
		{"fc046", "其他扣减费用、新农合、医保支付金额、其他商保支付之和不得大于账单发生金额！"},
		{"fc046", "自付一、自付二、自费、医疗保险基金支付金额、其他商保支付金额之和不得大于账单发生金额！"},
		{"fc046", "医疗保险范围内金额不能大于账单发生金额！"},
		{"fc046", "筹基金本次支付不能大于账单发生金额！"},
		{"fc046", "大额互助资金本次支付不能大于账单发生金额！"},
		{"fc059", "门诊大额支付、退休人员补充医疗保险支付、残疾军人补助支付、公务员医疗补助支付之和，不得大于医疗保险基金支付金额！"},
		{"fc082", "统筹基金本次支付不能大于统筹基金年度内累计支付！"},
		{"fc083", "大额互助资金本次支付不能大于大额互助资金累计支付金额！"},
	}
	for i, msg := range msgArr {
		for _, loc := range fieldLocationMap[msg[0]] {
			fc046Val := utils.GetFieldDecimalValueByLoc(obj, loc, true)
			total = decimal.Zero
			for _, code := range codesArr[i] {
				loc1 := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
				if len(loc1) != 1 {
					continue
				}
				val1 := utils.GetFieldDecimalValueByLocArr(obj, loc1, true)
				total = total.Add(val1)
			}
			if total.GreaterThan(fc046Val) {
				utils.SetIssue(obj, loc, msg[1], "7", "")
			}
		}
	}

	//CSB0116RC0077000
	//同一发票下，当fc155、fc156、fc157、fc158、fc159、fc160、fc161、fc162、fc163、fc164、fc165、fc166、fc167、fc168、fc169、fc170、fc171、fc172、fc173、fc174的结果值之和大于fc046的结果值时，fc046出问题件：问题件代码7，描述：所有扣除总金额之和不能大于账单发生金额！
	for _, loc := range fieldLocationMap["fc046"] {
		total = decimal.Zero
		codes = []string{"fc155", "fc156", "fc157", "fc158", "fc159", "fc160", "fc161", "fc162", "fc163", "fc164", "fc165", "fc166", "fc167", "fc168", "fc169", "fc170", "fc171", "fc172", "fc173", "fc174"}
		for _, code := range codes {
			fc155Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
			if len(fc155Loc) != 1 {
				continue
			}
			fc155Val := utils.GetFieldDecimalValueByLocArr(obj, fc155Loc, true)
			total = total.Add(fc155Val)
		}

		fc046Val := utils.GetFieldDecimalValueByLoc(obj, loc, true)
		if total.GreaterThan(fc046Val) {
			utils.SetIssue(obj, loc, "所有扣除总金额之和不能大于账单发生金额！", "7", "")
		}
	}

	//CSB0116RC0078000
	//同一发票下，当fc111录入值为5或8或10或11时，fc050、fc060、fc061、fc062、fc077、fc080的结果值之和不等于fc046的结果值时，fc046出问题件：问题件代码7，描述：自付一、自付二、自费、大额互助基金本次支付、统筹基金本次支付、其他商保支付金额之和必须等于账单金额！（字段不存在不执行该校验）
	for _, loc := range fieldLocationMap["fc046"] {
		fc111Loc := utils.GetFieldLoc(fieldLocationMap["fc111"], loc[0], -1, -1, -1)
		if len(fc111Loc) != 1 {
			continue
		}
		fc111Val := utils.GetFieldValueByLocArr(obj, fc111Loc, true)
		if m, _ := regexp.MatchString(`^(5|8|10|11)$`, fc111Val); !m {
			continue
		}
		total = decimal.Zero
		codes = []string{"fc050", "fc060", "fc061", "fc062", "fc077", "fc080"}
		for _, code := range codes {
			fc050Loc := utils.GetFieldLoc(fieldLocationMap[code], loc[0], -1, -1, -1)
			if len(fc050Loc) != 1 {
				continue
			}
			fc050Val := utils.GetFieldDecimalValueByLocArr(obj, fc050Loc, true)
			total = total.Add(fc050Val)
		}

		fc046Val := utils.GetFieldDecimalValueByLoc(obj, loc, true)
		if !fc046Val.Equals(total) {
			utils.SetIssue(obj, loc, "自付一、自付二、自费、大额互助基金本次支付、统筹基金本次支付、其他商保支付金额之和必须等于账单金额！", "7", "")
		}
	}

	//CSB0116RC0079000
	//同一发票下，当fc051的结果值大于fc046时，fc051出问题件：问题件代码7，描述：所有扣除总金额之和不能大于账单发生金额！
	for _, loc := range fieldLocationMap["fc051"] {
		loc1 := utils.GetFieldLoc(fieldLocationMap["fc046"], loc[0], -1, -1, -1)
		if len(loc1) != 1 {
			continue
		}
		val := utils.GetFieldDecimalValueByLoc(obj, loc, true)
		val1 := utils.GetFieldDecimalValueByLocArr(obj, loc1, true)
		if val.GreaterThan(val1) {
			utils.SetIssue(obj, loc, "所有扣除总金额之和不能大于账单发生金额！", "7", "")
		}
	}

	//CSB0116RC0080000
	//同一发票下，当fc054的结果值早于fc017的结果值时，fc054出问题件：问题件代码7，描述：门诊发生日期不能早于出险日期！
	//CSB0116RC0081000
	//同一发票下，当fc055的结果值早于fc017的结果值时，fc055出问题件：问题件代码7，描述：入院日期不能早于出险日期！
	codesArr = [][]string{
		{"fc054", "门诊发生日期不能早于出险日期！"},
		{"fc055", "入院日期不能早于出险日期！"},
	}
	fc017Val := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, fieldLocationMap["fc017"], true), "-", "")
	for _, code := range codesArr {
		for _, loc := range fieldLocationMap[code[0]] {
			val := strings.ReplaceAll(utils.GetFieldValueByLoc(obj, loc, true), "-", "")
			if utils.ParseDecimal(val).LessThan(utils.ParseDecimal(fc017Val)) {
				utils.SetIssue(obj, loc, codes[1], "7", "")
			}
		}
	}

	//CSB0116RC0082000
	//同一发票下，当fc055的结果值晚于fc056的结果值时，fc056出问题件：问题件代码7，描述：入院日期不能晚于出院日期！
	for _, loc := range fieldLocationMap["fc055"] {
		fc056Loc := utils.GetFieldLoc(fieldLocationMap["fc056"], loc[0], -1, -1, -1)
		if len(fc056Loc) != 1 {
			continue
		}
		val := strings.ReplaceAll(utils.GetFieldValueByLoc(obj, loc, true), "-", "")
		fc056Val := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, fc056Loc, true), "-", "")
		if utils.ParseDecimal(val).GreaterThan(utils.ParseDecimal(fc056Val)) {
			utils.SetIssue(obj, loc, "入院日期不能晚于出院日期！", "7", "")
		}
	}

	//CSB0116RC0083000
	//fc027录入值匹配《B0116_华夏人寿团险理赔_残疾》的“残疾中文名”（第四列）对应的“残疾编码”（第三列）进行转码，
	//无法转码时结果值默认为空，并出问题件：问题件代码7，描述：残疾名称不在常量库中（录入内容为A、为空、为F或包含?时不做校验）
	for _, loc := range fieldLocationMap["fc027"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		if m, _ := regexp.MatchString(`^(A|F|)$`, val); m || strings.Index(val, "?") != -1 {
			continue
		}
		value, row := utils.FetchConst(bill.ProCode, "B0116_华夏人寿团险理赔_残疾", "残疾编码", map[string]string{"残疾中文名": val})
		if row == 0 {
			utils.SetIssue(obj, loc, "残疾名称不在常量库中", "7", "")
		}
		utils.SetOnlyOneFinalValue(obj, loc, value)
	}

	//CSB0116RC0084000
	//fc032结果值早于fc017结果值时（fc032结果值为空时不进行校验），fc032出问题件：问题件代码：7，问题件描述：鉴定日期不能早于出险日期！
	for _, loc := range fieldLocationMap["fc032"] {
		loc1 := utils.GetFieldLoc(fieldLocationMap["fc017"], loc[0], -1, -1, -1)
		if len(loc1) != 1 {
			continue
		}
		val := strings.ReplaceAll(utils.GetFieldValueByLoc(obj, loc, true), "-", "")
		val1 := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, loc1, true), "-", "")
		if utils.ParseDecimal(val).LessThan(utils.ParseDecimal(val1)) {
			utils.SetIssue(obj, loc, "鉴定日期不能早于出险日期！", "7", "")
		}
	}

	//CSB0116RC0085000
	//fc108录入内容匹配《B0116_华夏人寿团险理赔_重疾》的“中文名”（第三列）对应的“编码”（第二列）进行转码输出，
	//无法转码时结果值默认为空，并出问题件：问题件代码7，描述：重疾名称不在常量库中；（录入内容为A、为空、为F或包含?时不做校验）
	for _, loc := range fieldLocationMap["fc108"] {
		val := utils.GetFieldValueByLoc(obj, loc, false)
		if m, _ := regexp.MatchString(`^(A|F|)$`, val); m || strings.Index(val, "?") != -1 {
			continue
		}
		value, row := utils.FetchConst(bill.ProCode, "B0116_华夏人寿团险理赔_重疾", "编码", map[string]string{"中文名": val})
		if row == 0 {
			utils.SetIssue(obj, loc, "重疾名称不在常量库中", "7", "")
		}
		utils.SetOnlyOneFinalValue(obj, loc, value)
	}

	//CSB0116RC0086000
	//fc203、fc206录入内容匹配《B0116_华夏人寿团险理赔_职业编码》的“职业名称”（第二列）对应的“职业编码”（第一列）进行转码输出，
	//当无法转码时结果值默认为空，对应字段出问题件：问题件代码7，问题件描述：申请人职业编码不在常量库中；（录入内容为A、为空或包含?时不做校验）
	codes = []string{"fc203", "fc206"}
	for _, code := range codes {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			if m, _ := regexp.MatchString(`^(A|)$`, val); m || strings.Index(val, "?") != -1 {
				continue
			}
			value, row := utils.FetchConst(bill.ProCode, "B0116_华夏人寿团险理赔_职业编码", "职业编码", map[string]string{"职业名称": val})
			if row == 0 {
				utils.SetIssue(obj, loc, "申请人职业编码不在常量库中", "7", "")
			}
			utils.SetOnlyOneFinalValue(obj, loc, value)
		}
	}

	//CSB0116RC0087000
	//fc036结果值早于fc017结果值时（fc036结果值为空时不进行校验），fc036出问题件：问题件代码：7，问题件描述：确诊日期不能早于出险日期！
	//CSB0116RC0088000
	//fc217结果值早于fc216结果值时（fc217结果值为空时不进行校验），fc217出问题件：问题件代码：7，问题件描述：失能截止日期不能早于失能开始日期！
	//CSB0116RC0089000
	//fc221结果值早于fc220结果值时（fc221结果值为空时不进行校验），fc221出问题件：问题件代码：7，问题件描述：护理截止日期不能早于护理开始日期！
	codesArr = [][]string{
		{"fc036", "fc017", "确诊日期不能早于出险日期！"},
		{"fc217", "fc216", "失能截止日期不能早于失能开始日期！"},
		{"fc221", "fc220", "护理截止日期不能早于护理开始日期！"},
	}
	for _, codes = range codesArr {
		val0 := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[0]], true), "-", "")
		val1 := strings.ReplaceAll(utils.GetFieldValueByLocArr(obj, fieldLocationMap[codes[1]], true), "-", "")
		if utils.ParseDecimal(val0).LessThan(utils.ParseDecimal(val1)) {
			utils.SetIssues(obj, fieldLocationMap[codes[0]], codes[2], "7", "")
		}
	}

	//CSB0116RC0090000
	//"所有字段结果数据包含?时，将结果数据转换为空，该字段出问题件：问题件代码1，问题件描述：（字段名）填写模糊。
	//(将该代码放在所有需求的最后面。)"
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
						if strings.Index(fields[l].ResultValue, "?") != -1 {
							fields[l].FinalValue = ""
							utils.SetIssue(obj, []int{i, j, k, l}, fields[l].Name+"填写模糊", "1", "")
						}
					}
				}
			}
		}
	}

	//CSB0116RC0091000
	//"fc115、fc116、fc117、fc118、fc119、fc120、fc121、fc122、fc123、fc124、fc125、fc126、fc127、fc128、fc129、fc130、fc131、fc132、fc133、fc134、fc135、fc136、fc137、fc138、fc139、fc140、fc141、fc142、fc143、fc144、fc145、fc146
	//的录入值包含?时将上述字段的结果值转为空并清除问题件（上述字段所在分块均为循环分块）
	//该需求需要放在需求“所有字段结果数据包含?时，将结果数据转换为空，该字段出问题件......”后面"
	codes = []string{"fc115", "fc116", "fc117", "fc118", "fc119", "fc120", "fc121", "fc122", "fc123", "fc124", "fc125", "fc126", "fc127", "fc128", "fc129", "fc130", "fc131", "fc132", "fc133", "fc134", "fc135", "fc136", "fc137", "fc138", "fc139", "fc140", "fc141", "fc142", "fc143", "fc144", "fc145", "fc146"}
	for _, code := range codes {
		for _, loc := range fieldLocationMap[code] {
			val := utils.GetFieldValueByLoc(obj, loc, false)
			if strings.Index(val, "?") != -1 {
				utils.SetOnlyOneFinalValue(obj, loc, "")
				utils.DelOnlyOneIssue(obj, loc)
			}
		}
	}

	//CSB0116RC0092000
	//"当单证类型为“1”时，清空以下字段的问题件：
	codes = []string{"fc015", "fc016", "fc018", "fc044", "fc096", "fc099"}
	if obj.Bill.InsuranceType == "" {
		for _, code := range codes {
			utils.DelIssue(obj, fieldLocationMap[code])
		}
	}

	return nil, obj
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"accidental", "B0116_华夏人寿团险理赔_意外出险细节", "1", "0"},
		{"damages", "B0116_华夏人寿团险理赔_损伤外部原因", "1", "0"},
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
