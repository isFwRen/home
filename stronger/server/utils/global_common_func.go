package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	model2 "server/module/export/model"
	model3 "server/module/load/model"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func Calculate(start, end string) (startT, endT time.Time, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startT, err = time.ParseInLocation("150405", start, local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endT, err = time.ParseInLocation("150405", end, local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startT, endT, err
}

func StringToTime(t string) (T time.Time, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return T, nil
	}
	T, err = time.ParseInLocation("20060102150405", t, local)
	if err != nil {
		return T, err
	}
	fmt.Println(T.Format("20060102150405"))
	return T, nil
}

// CheckRequired 限 string int 数组 有需再增加
func CheckRequired(str interface{}, Update map[string]interface{}) {
	t := reflect.TypeOf(str)
	v := reflect.ValueOf(str)
	for k := 0; k < t.NumField(); k++ {
		//fieldType := v.Field(k).Kind()
		//if fieldType == reflect.Struct {
		//	CheckRequired(v.Field(k).Interface(), Update)
		//}
		if t.Field(k).Tag.Get("db") != "none" {
			if !v.Field(k).IsZero() {
				if v.Field(k).Kind() == 2 {
					Update[t.Field(k).Tag.Get("db")] = v.Field(k).Int()
				}
				if v.Field(k).Kind() == 23 {
					fmt.Println(t.Field(k).Tag.Get("db"))
					arr := pq.StringArray{}
					for i := 0; i < v.Field(k).Slice(0, v.Field(k).Len()).Len(); i++ {
						arr = append(arr, v.Field(k).Slice(0, v.Field(k).Len()).Index(i).Interface().(string))
					}
					fmt.Println("4", arr)
					Update[t.Field(k).Tag.Get("db")] = arr
				}
				if v.Field(k).Kind() == 24 {
					Update[t.Field(k).Tag.Get("db")] = v.Field(k).String()
				}
			}
		}
	}
}

// GetFieldLoc 获取同一发票的某一字段的位置 i:第几张发票，j:第几个对象数组（Invoice，QingDan。。。）, k:对象数组第几项
func GetFieldLoc(locs [][]int, i, j, k, blockIndex int) [][]int {
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

// GetNodeValue  获取xml节点的第一个值
func GetNodeValue(xmlValue string, nodeName string) string {
	sValue := ""
	beginNode := strings.Index(xmlValue, "<"+nodeName+">") + len(nodeName) + 2
	endNode := strings.Index(xmlValue, "</"+nodeName+">")
	if beginNode != -1 && endNode != -1 {
		sValue = xmlValue[beginNode:endNode]
	}
	return sValue
}

// GetNodeData 取某节点的值
func GetNodeData(xml string, node string) (a []string) {
	reg := regexp.MustCompile("<" + node + ">[\\s\\S]*?</" + node + ">")
	arr := reg.FindAllString(xml, -1)
	for _, s := range arr {
		v := strings.ReplaceAll(s, "<"+node+">", "")
		v = strings.ReplaceAll(v, "</"+node+">", "")
		//fmt.Println(v)
		a = append(a, v)
	}
	return a
}

// GetNode 取整个节点
func GetNode(xml string, node string) []string {
	reg := regexp.MustCompile("<" + node + ">[\\s\\S]*?</" + node + ">")
	arr := reg.FindAllString(xml, -1)
	return arr
}

func SetNodeValue(xmlValue, nodeName, value string) string {
	reg := regexp.MustCompile(`>.*</` + nodeName + `>`)
	return reg.ReplaceAllString(xmlValue, ">"+value+"</"+nodeName+">")
}

func Substr(str string, start, end int) string {
	if start == -1 || end == -1 {
		return ""
	}
	return string(str[start:end])
}

// GetNodeDecimalValue 获取节点值并返回decimal
func GetNodeDecimalValue(xmlValue, nodeName string) decimal.Decimal {
	beginNode := strings.Index(xmlValue, "<"+nodeName+">") + len(nodeName) + 2
	endNode := strings.Index(xmlValue, "</"+nodeName+">")
	sValue := ""
	if beginNode != -1 && endNode != -1 {
		sValue = Substr(xmlValue, beginNode, endNode)
	}
	return ParseDecimal(sValue)
}

// ParseDecimal string转decimal
func ParseDecimal(v string) decimal.Decimal {
	d, err := decimal.NewFromString(v)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// GetFieldByLoc 根据坐标获取字段
func GetFieldByLoc(obj model2.ResultDataBill, loc []int) (err error, field model3.ProjectField) {
	if len(obj.Invoice) < loc[0]+1 {
		global.GLog.Error("loc", zap.Any("", loc))
		return errors.New("发票长度越界"), field
	}
	invoice := reflect.ValueOf(obj.Invoice[loc[0]])
	if invoice.NumField() < loc[1]+1 {
		global.GLog.Error("loc", zap.Any("", loc))
		return errors.New("发票对象字段数量越界"), field
	}
	fieldArr := invoice.Field(loc[1]).Interface().([][]model3.ProjectField)
	if len(fieldArr) < loc[2]+1 {
		global.GLog.Error("loc", zap.Any("", loc))
		return errors.New("没有该下标的字段"), field
	}
	return nil, fieldArr[loc[2]][loc[3]]
}

// GetFieldValueByLoc 根据坐标获取录入值or结果值
func GetFieldValueByLoc(obj model2.ResultDataBill, loc []int, isFinalValue bool) string {
	err, field := GetFieldByLoc(obj, loc)
	if err != nil {
		return ""
	}
	if isFinalValue {
		return field.FinalValue
	}
	return field.ResultValue
}

// GetFieldDecimalValueByLoc 根据坐标获取录入值or结果值 数字
func GetFieldDecimalValueByLoc(obj model2.ResultDataBill, loc []int, isFinalValue bool) decimal.Decimal {
	val := GetFieldValueByLoc(obj, loc, isFinalValue)
	if val == "" || val == "A" || strings.Index(val, "?") != -1 {
		return decimal.Zero
	}
	if isFinalValue {
		return ParseDecimal(val)
	}
	return ParseDecimal(val)
}

// GetFieldDecimalValueArrByLocArr 根据坐标获取录入值or结果值 数字数组
func GetFieldDecimalValueArrByLocArr(obj model2.ResultDataBill, locArr [][]int, isFinalValue bool) []decimal.Decimal {
	iArr := make([]decimal.Decimal, len(locArr))
	for _, loc := range locArr {
		val := GetFieldValueByLoc(obj, loc, isFinalValue)
		if val == "" || val == "A" || strings.Index(val, "?") != -1 {
			iArr = append(iArr, decimal.Zero)
		}
		if isFinalValue {
			iArr = append(iArr, ParseDecimal(val))
		}
	}
	return iArr
}

// GetFieldDecimalValueByLocArr 根据坐标获取录入值or结果值 数字
func GetFieldDecimalValueByLocArr(obj model2.ResultDataBill, locArr [][]int, isFinalValue bool) decimal.Decimal {
	for _, loc := range locArr {
		val := GetFieldValueByLoc(obj, loc, isFinalValue)
		if val == "" || val == "A" || strings.Index(val, "?") != -1 {
			return decimal.Zero
		}
		if isFinalValue {
			return ParseDecimal(val)
		}
	}
	return decimal.Zero
}

// GetFieldValueArrByLocArr 根据坐标数组获取字段录入值or结果值数组
func GetFieldValueArrByLocArr(obj model2.ResultDataBill, locArr [][]int, isFinalValue bool) []string {
	filedArr := make([]string, len(locArr))
	for i, loc := range locArr {
		val := GetFieldValueByLoc(obj, loc, isFinalValue)
		filedArr[i] = val
	}
	return filedArr
}

// GetFieldValueByLocArr 根据坐标数组获取一个字段录入值or结果值数组
func GetFieldValueByLocArr(obj model2.ResultDataBill, locArr [][]int, isFinalValue bool) string {
	for _, loc := range locArr {
		val := GetFieldValueByLoc(obj, loc, isFinalValue)
		return val
	}
	return ""
}

// GetFieldByLocArr 根据坐标数组获取一个字段
func GetFieldByLocArr(obj model2.ResultDataBill, locArr [][]int) (field model3.ProjectField) {
	for _, loc := range locArr {
		_, field = GetFieldByLoc(obj, loc)
		return field
	}
	return field
}

// SetIssues 设置多个问题件
func SetIssues(obj model2.ResultDataBill, locArr [][]int, issueMessage, issueCode, issueType string) {
	for _, loc := range locArr {
		SetIssue(obj, loc, issueMessage, issueCode, issueType)
	}
}

// SetIssue 设置问题件
func SetIssue(obj model2.ResultDataBill, loc []int, issueMessage, issueCode, issueType string) {
	issue := model3.Issue{
		Type:    issueType,
		Code:    issueCode,
		Message: issueMessage,
	}
	eleLen := reflect.ValueOf(obj.Invoice[loc[0]]).NumField()
	if eleLen > 0 {
		if reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.Float64 {
			fieldsArr := reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Interface().([][]model3.ProjectField)
			fieldsArr[loc[2]][loc[3]].Issues = append(fieldsArr[loc[2]][loc[3]].Issues, issue)
		}
	}
}

// DelIssue 清空问题件
func DelIssue(obj model2.ResultDataBill, locArr [][]int) {
	for _, loc := range locArr {
		DelOnlyOneIssue(obj, loc)
	}
}

// DelOnlyOneIssue 清空一个问题件
func DelOnlyOneIssue(obj model2.ResultDataBill, loc []int) {
	eleLen := reflect.ValueOf(obj.Invoice[loc[0]]).NumField()
	if eleLen > 0 {
		if reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.Float64 {
			fieldsArr := reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Interface().([][]model3.ProjectField)
			fieldsArr[loc[2]][loc[3]].Issues = nil
		}
	}

}

// SetFinalValue 设置结果值
func SetFinalValue(obj model2.ResultDataBill, locArr [][]int, value string) {
	for _, loc := range locArr {
		eleLen := reflect.ValueOf(obj.Invoice[loc[0]]).NumField()
		if eleLen > 0 {
			if reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.Float64 {
				fieldsArr := reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Interface().([][]model3.ProjectField)
				fieldsArr[loc[2]][loc[3]].FinalValue = value
			}
		}
	}
}

// SetOnlyOneFinalValue 设置结果值
func SetOnlyOneFinalValue(obj model2.ResultDataBill, loc []int, value string) {
	eleLen := reflect.ValueOf(obj.Invoice[loc[0]]).NumField()
	if eleLen > 0 {
		if reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.Float64 {
			fieldsArr := reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Interface().([][]model3.ProjectField)
			fieldsArr[loc[2]][loc[3]].FinalValue = value
		}
	}
}

// SetOnlyOneFinalInput 设置结果值状态 no:未匹配常量
func SetOnlyOneFinalInput(obj model2.ResultDataBill, loc []int) {
	eleLen := reflect.ValueOf(obj.Invoice[loc[0]]).NumField()
	if eleLen > 0 {
		if reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.Float64 {
			fieldsArr := reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Interface().([][]model3.ProjectField)
			fieldsArr[loc[2]][loc[3]].FinalInput = "no_match_const"
		}
	}
}

// SetOnlyOneResultValue 设置录入值
func SetOnlyOneResultValue(obj model2.ResultDataBill, loc []int, value string) {
	eleLen := reflect.ValueOf(obj.Invoice[loc[0]]).NumField()
	if eleLen > 0 {
		if reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.String && reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Kind() != reflect.Float64 {
			fieldsArr := reflect.ValueOf(obj.Invoice[loc[0]]).Field(loc[1]).Interface().([][]model3.ProjectField)
			fieldsArr[loc[2]][loc[3]].ResultValue = value
			fieldsArr[loc[2]][loc[3]].IsChange = true
		}
	}
}

// CalcChar 计算字符数
func CalcChar(str string) (character int) {
	for _, s := range str {
		v := fmt.Sprintf("%c", s)
		if v == "" {
			continue
		}
		reg3 := regexp.MustCompile("/·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、/")
		for _, k := range v {
			if unicode.Is(unicode.Scripts["Han"], k) || reg3.Match([]byte(string(k))) {
				character += 2
			} else {
				character += 1
			}
		}
	}
	return character
}

// RegIsMatch 正则是否匹配
func RegIsMatch(reg string, value string) bool {
	matched, _ := regexp.MatchString(reg, value)
	return matched
}

func RegMatchAll(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.FindAllString(value, -1)
}

func ToString(data float64) string {
	// return strconv.FormatFloat(data, 'E', -1, 64)
	return fmt.Sprintf("%v", data)
}

func HasKey(data map[string]string, key string) bool {
	_, isOK := data[key]
	return isOK
}

func ParseFloat(value string) float64 {
	val, _ := strconv.ParseFloat(value, 64)
	return val
}

func ToFix2(value string) string {
	decimalFinalValue := ParseDecimal(value).Round(2)
	return decimalFinalValue.StringFixed(2)
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

// GetValByFields 获取字段数组的字段值
func GetValByFields(fields []model3.ProjectField, code string, isFinalValue bool) string {
	for _, field := range fields {
		if field.Code == code {
			if isFinalValue {
				return field.FinalValue
			}
			return field.ResultValue
		}
	}
	return ""
}

// GetDecimalValByFields 获取字段数组的字段值并返回decimal
func GetDecimalValByFields(fields []model3.ProjectField, code string, isFinalValue bool) decimal.Decimal {
	return ParseDecimal(GetValByFields(fields, code, isFinalValue))
}

func RegReplace(data string, query string, value string) string {
	reg := regexp.MustCompile(query)
	return reg.ReplaceAllString(data, value)
}

// FetchFieldValBySaveBlockIndex 获取同一小分块的字段的值
func FetchFieldValBySaveBlockIndex(obj model2.ResultDataBill, locs [][]int, i, j, k, blockIndex int, isFinalValue bool) string {
	locsArr := GetFieldLoc(locs, i, j, k, blockIndex)
	if len(locsArr) != 1 {
		return ""
	}
	return GetFieldValueByLoc(obj, locsArr[0], isFinalValue)
}

func CheckIDCard(id string) bool {
	// 身份证位数不对
	if len(id) != 15 && len(id) != 18 {
		return false
	}

	// 转大写
	id = strings.ToUpper(id)

	if len(id) == 18 {
		// 验证算法
		if !checkValidNo18(id) {
			fmt.Println(id, "身份证算法验证失败！")
			return false
		}

	} else {
		// 转18位
		id = idCard15To18(id)
	}

	// 生日验证
	if !checkBirthdayCode(id[6:14]) {
		fmt.Println(id, "生日验证失败！")
		return false
	}

	// 验证地址
	if !checkAddressCode(id[:2]) {
		fmt.Println(id, "地址验证失败！")
		return false
	}

	return true
}

// 15位身份证转为18位
var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var validValue = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

// 15位转18位
func idCard15To18(id15 string) string {
	nLen := len(id15)
	if nLen != 15 {
		return "身份证不是15位！"
	}
	id18 := make([]byte, 0)
	id18 = append(id18, id15[:6]...)
	id18 = append(id18, '1', '9')
	id18 = append(id18, id15[6:]...)

	sum := 0
	for i, v := range id18 {
		n, _ := strconv.Atoi(string(v))
		sum += n * weight[i]
	}
	mod := sum % 11
	id18 = append(id18, validValue[mod])
	return string(id18)
}

// 18位身份证校验码
func checkValidNo18(id string) bool {
	//string -> []byte
	id18 := []byte(id)
	nSum := 0
	for i := 0; i < len(id18)-1; i++ {
		n, _ := strconv.Atoi(string(id18[i]))
		nSum += n * weight[i]
	}
	//mod得出18位身份证校验码
	mod := nSum % 11
	if validValue[mod] == id18[17] {
		return true
	}

	return false
}

// 验证生日
func checkBirthdayCode(birthday string) bool {
	year, _ := strconv.Atoi(birthday[:4])
	month, _ := strconv.Atoi(birthday[4:6])
	day, _ := strconv.Atoi(birthday[6:])

	curYear, curMonth, curDay := time.Now().Date()
	//出生日期大于现在的日期
	if year < 1900 || year > curYear || month <= 0 || month > 12 || day <= 0 || day > 31 {
		return false
	}

	if year == curYear {
		if month > int(curMonth) {
			return false
		} else if month == int(curMonth) && day > curDay {
			return false
		}
	}

	//出生日期在2月份
	if 2 == month {
		//闰年2月只有29号
		if isLeapYear(year) && day > 29 {
			return false
		} else if day > 28 { //非闰年2月只有28号
			return false
		}
	} else if 4 == month || 6 == month || 9 == month || 11 == month { //小月只有30号
		if day > 30 {
			return false
		}
	}

	return true
}

// 判断是否为闰年
func isLeapYear(year int) bool {
	if year <= 0 {
		return false
	}
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

// 验证地区
// strict: true-验证详细， false-验证省
func checkAddressCode(address string) bool {
	aa, _ := strconv.Atoi(address)
	if aa < 11 || aa > 91 {
		return false
	}

	return true
}
