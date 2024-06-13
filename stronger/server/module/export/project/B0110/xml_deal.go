/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年04月23日16:28:53
 */

package B0110

import (
	"encoding/json"
	"fmt"
	"regexp"
	"server/global"
	"server/module/load/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0110:::XmlDeal-----------------------")
	obj := o.(FormatObj)
	fields := obj.Fields
	bill_num := strings.Split(obj.Bill.BillNum, "_")
	xmlValue = SetNodeValue(xmlValue, "claimNo", bill_num[0])
	var otherInfo OtherInfo
	err = json.Unmarshal([]byte(obj.Bill.OtherInfo), &otherInfo)
	if err != nil {
		global.GLog.Error("otherInfo", zap.Error(err))
	}
	//utils.GetNodeValue(obj.Bill.OtherInfo, "claimTpaId")
	xmlValue = SetNodeValue(xmlValue, "claimTpaId", otherInfo.ClaimTpaId)
	// constMap := InitConst(obj.Bill.ProCode)

	// mulubianma := constMap["mulubianma"]
	agency := obj.Bill.Agency
	v, _ := utils.FetchConst(obj.Bill.ProCode, "B0110_新疆国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": agency})
	xmlValue = SetNodeValue(xmlValue, "catalogCode", v)

	fc022s := GetSameField(obj.Fields, "fc022", true)
	eraday := ""
	for _, fc022 := range fc022s {
		if fc022 != "" {
			a, _ := time.Parse("2006-01-02", fc022)
			b, _ := time.Parse("2006-01-02", eraday)
			if eraday == "" {
				eraday = fc022
			} else if a.Before(b) {
				eraday = fc022
			}
		}
	}
	// if fc112 != "" && fc111s[ff] != "" {
	// 	a, _ := time.Parse("2006/01/02", fc111s[ff])
	// 	b, _ := time.Parse("2006/01/02", fc112)
	// 	d := b.Sub(a)
	// 	hospitalAllowance_sum += d.Hours() / 24
	// }

	questionCount := 0
	errMeses := ""
	fc003s := GetSameField(obj.Fields, "fc003", false)
	if arrays.Contains(fc003s, "4") == -1 {
		xmlValue = SetNodeValue(xmlValue, "mIcd10Code", "R52.9")
		errMeses += "案件无疾病诊断;"
		questionCount++
	}

	for _, field := range fields {
		// if len(filed.Issues ) > 0 {
		for _, issue := range field.Issues {
			// fname := field.Name
			// issue.Message = strings.Replace(issue.Message, fname, "", 1)
			// errMes := fname + issue.Message + ";"
			errMes := issue.Message + ";"
			if strings.Index(errMeses, errMes) == -1 {
				errMeses += errMes
				questionCount++
			}
		}
	}

	rcptInfoLists := RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	// xmlValue = SetNodeValue(xmlValue, "pageSum", strconv.Itoa(len(rcptInfoLists)))

	for _, rcptInfoList := range rcptInfoLists {
		old_rcptInfoList := rcptInfoList
		if RegIsMatch(rcptInfoList, `<(rcptList|errorRcptList)>`) {
			rcptInfoList = SetNodeValue(rcptInfoList, "rcptType", "D")
		} else {
			rcptInfoList = SetNodeValue(rcptInfoList, "rcptType", "M")
			// rcptInfoList = RegReplace(rcptInfoList, `<rcptList>[\s\S]*?<\/rcptList>`, "")
			// rcptInfoList = RegReplace(rcptInfoList, `<errorRcptList>[\s\S]*?<\/errorRcptList>`, "")
			// rcptInfoList = SetNodeValue(rcptInfoList, "rcptList", "")
			// rcptInfoList = SetNodeValue(rcptInfoList, "errorRcptList", "")
		}
		expenMode := GetNodeValue(rcptInfoList, "expenMode")
		// fmt.Println("------------expenMode-----------------", expenMode, obj.Bill.Agency)
		if expenMode == "0020" && obj.Bill.Agency == "650192" {
			beginDate := GetNodeValue(rcptInfoList, "beginDate")
			b, _ := time.Parse("2006-01-02", beginDate)
			a, _ := time.Parse("2006-01-02", eraday)
			d := b.Sub(a)
			// fmt.Println("------------b-----------------", b, a, d, d.Hours())
			if (d.Hours() / 24) > 15 {
				rcptInfoList = ""
				if strings.Index(errMeses, "全民门诊，就诊时间超过15日的票据未录入，请处理人员审核;") == -1 {
					errMeses += "全民门诊，就诊时间超过15日的票据未录入，请处理人员审核;"
					questionCount++
				}

			}
		}
		// socialPayAmnt := GetNodeValue(rcptInfoList, "socialPayAmnt")
		// if socialPayAmnt == "0.00" {
		// 	rcptInfoList = SetNodeValue(rcptInfoList, "siHealthType", "00")
		// } else {
		// 	rcptInfoList = SetNodeValue(rcptInfoList, "siHealthType", "10")
		// }
		xmlValue = strings.Replace(xmlValue, old_rcptInfoList, rcptInfoList, 1)
	}

	//CSB0110RC0117000 问题件描述重复时，保留一个即可
	seen := make(map[string]bool)
	var result strings.Builder
	strArr := strings.Split(errMeses, ";")
	for _, str := range strArr {
		if str != "" && !seen[str] {
			seen[str] = true
			result.WriteString(str)
			result.WriteString(";")
		}
	}
	errMeses = result.String()
	xmlValue = SetNodeValue(xmlValue, "bpoAbnormalReason", errMeses)

	if questionCount > 0 {
		xmlValue = SetNodeValue(xmlValue, "bpoIsAbnormal", "0")
	} else {
		xmlValue = SetNodeValue(xmlValue, "bpoIsAbnormal", "1")
	}

	// items := [][]string{
	// 	{"rcptOtherPayDetailList", "uedUnitNo", "otherSpiltAmnt", "\"rcptOtherPayDetailList\": [],"},
	// 	{"costCategoryList", "medicalCategoryName", "classifiedAmount", "\"costCategoryList\": [],"},
	// }
	// for _, item := range items {
	// 	nodes := RegMatchAll(xmlValue, `<`+item[0]+`>[\s\S]*?<\/`+item[0]+`>`)
	// 	for _, node := range nodes {
	// 		val1 := GetNodeValue(node, item[1])
	// 		val2 := GetNodeValue(node, item[2])
	// 		if val1 == "" && val2 == "" {
	// 			xmlValue = strings.Replace(xmlValue, node, item[3], 1)
	// 		}
	// 	}

	// }

	items := [][]string{
		{"rcptList", "listCode"},
		{"errorRcptList", "listName"},
	}
	for _, item := range items {
		nodes := RegMatchAll(xmlValue, `<`+item[0]+`>[\s\S]*?<\/`+item[0]+`>`)
		for _, node := range nodes {
			val1 := GetNodeValue(node, item[1])
			// fmt.Println("------------111111-----------------", item[0], val1)
			if val1 == "" {
				xmlValue = strings.Replace(xmlValue, node, "", 1)
			}
		}

	}

	rcptInfoLists = RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	//xmlValue = SetNodeValue(xmlValue, "pageSum", strconv.Itoa(len(rcptInfoLists)))

	for _, rcptInfoList := range rcptInfoLists {
		old_rcptInfoList := rcptInfoList
		ii := 1
		for _, item := range items {
			nodes := RegMatchAll(rcptInfoList, `<`+item[0]+`>[\s\S]*?<\/`+item[0]+`>`)
			for _, node := range nodes {
				old_node := node
				val1 := GetNodeValue(node, item[1])
				// node = SetNodeValue(node, "rcptListNo", strconv.Itoa(ii+1))
				// fmt.Println("------------22222----------------", item[0], val1)
				if item[0] == "rcptList" && !RegIsMatch(val1, `^\d+$`) {
					node = ""
				}
				if item[0] == "errorRcptList" && RegIsMatch(val1, `^\d+$`) {
					node = ""
				}
				if node != "" {
					node = SetNodeValue(node, "rcptListNo", strconv.Itoa(ii))
					ii++
				}
				rcptInfoList = strings.Replace(rcptInfoList, old_node, node, 1)
			}

		}
		xmlValue = strings.Replace(xmlValue, old_rcptInfoList, rcptInfoList, 1)
	}

	// rcptListNos := RegMatchAll(xmlValue, `<rcptListNo>.*?<\/rcptListNo>`)

	// for ii, rcptListNo := range rcptListNos {
	// 	xmlValue = strings.Replace(xmlValue, rcptListNo, "<rcptListNo>"+strconv.Itoa(ii)+"</rcptListNo>", 1)
	// }

	ids := RegMatchAll(xmlValue, `<id>.*?<\/id>`)
	// fmt.Println("------------ids-----------------", len(ids))
	for _, idItem := range ids {
		old_idItem := idItem
		id := GetNodeValue(idItem, "id")
		// fmt.Println("------------id-----------------", id)
		picIdx, err := strconv.Atoi(id)
		// fmt.Println("------------CmsImageInfoList-----------------", otherInfo.CmsImageInfoList)
		if err == nil {
			pid := otherInfo.CmsImageInfoList[picIdx].ID
			idItem = SetNodeValue(idItem, "id", pid)
			xmlValue = strings.Replace(xmlValue, old_idItem, idItem, 1)
		}
	}

	itemNodes := []string{"bankAccInfoList", "rcptInfoList", "rcptList", "errorRcptList"}

	for _, itemNode := range itemNodes {
		nodeDatas := RegMatchAll(xmlValue, `<`+itemNode+`>[\s\S]*?<\/`+itemNode+`>`)
		if len(nodeDatas) == 1 {
			aaa := "<" + itemNode + "><array_list></array_list></" + itemNode + ">"
			cloneNodeData := nodeDatas[0]
			cloneNodeData = cloneNodeData + aaa
			xmlValue = strings.Replace(xmlValue, nodeDatas[0], cloneNodeData, 1)
		}
	}

	rcptInfoLists = RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	for _, rcptInfoList := range rcptInfoLists {
		old_rcptInfoList := rcptInfoList
		rcptPrintType := GetNodeValue(rcptInfoList, "rcptPrintType")
		if rcptPrintType == "3" || !RegIsMatch(rcptPrintType, `^(3|4|5)$`) {
			aaa := "<vatInvoiceInfoTo><array_list></array_list></vatInvoiceInfoTo>"
			aaa += aaa
			rcptInfoList = RegReplace(rcptInfoList, `<vatInvoiceInfoTo>[\s\S]*?<\/vatInvoiceInfoTo>`, aaa)
		}
		if rcptPrintType == "4" || rcptPrintType == "5" || !RegIsMatch(rcptPrintType, `^(3|4|5)$`) {
			aaa := "<billInfoTo><array_list></array_list></billInfoTo>"
			aaa += aaa
			rcptInfoList = RegReplace(rcptInfoList, `<billInfoTo>[\s\S]*?<\/billInfoTo>`, aaa)
		}
		xmlValue = strings.Replace(xmlValue, old_rcptInfoList, rcptInfoList, 1)
	}

	mapItems := map[string][]string{
		"billInfoTo":       {"billDetailListTo", "billInfoItemListTo", "billInfoExtListTo"},
		"vatInvoiceInfoTo": {"vatInvoiceChargeItenListTo", "vatInvoiceExtListTo"},
	}
	for key, mapItem := range mapItems {
		mLists := RegMatchAll(xmlValue, `<`+key+`>[\s\S]*?<\/`+key+`>`)
		for _, mList := range mLists {
			old_mList := mList
			for _, item := range mapItem {
				mItems := RegMatchAll(xmlValue, `<`+item+`>[\s\S]*?<\/`+item+`>`)
				for _, mItem := range mItems {
					if !RegIsMatch(mItem, `>.+<`) {
						aaa := "<" + item + "><array_list></array_list></" + item + ">"
						aaa += aaa
						mList = strings.Replace(mList, mItem, aaa, 1)
					}
				}
			}
			xmlValue = strings.Replace(xmlValue, old_mList, mList, 1)
		}
	}

	//CSB0110RC0113000 当rcpNo节点为空时，屏蔽rcptInfoList一整个大节点
	rcptInfoList := RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	for _, item := range rcptInfoList {
		rcptNo := GetNodeValue(item, "rcptNo")
		if len(rcptNo) == 0 {
			xmlValue = strings.Replace(xmlValue, item, "", 1)
		}
	}
	rcptInfoLists = RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	xmlValue = SetNodeValue(xmlValue, "pageSum", strconv.Itoa(len(rcptInfoLists)))
	return err, xmlValue
}

type OtherInfo struct {
	BpoSendRemark    string             `json:"bpoSendRemark"`
	ClaimTpaId       string             `json:"claimTpaId"`
	CmsImageInfoList []CmsImageInfoList `json:"cmsImageInfoList"`
}

type CmsImageInfoList struct {
	ID string `json:"id"` // 账单号
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

func GetOneField(fields []model.ProjectField, code string, finalOrResult bool) (bool, string) {
	for _, field := range fields {
		if field.Code == code {
			if finalOrResult {
				return true, field.FinalValue
			} else {
				return true, field.ResultValue
			}
		}
	}
	return false, ""
}

func GetSameField(fields []model.ProjectField, code string, finalOrResult bool) (valuse []string) {
	for _, field := range fields {
		if field.Code == code {
			if finalOrResult {
				valuse = append(valuse, field.FinalValue)
			} else {
				valuse = append(valuse, field.ResultValue)
			}
		}
	}
	return valuse
}

func ToString(data float64) string {
	// return strconv.FormatFloat(data, 'E', -1, 64)
	return fmt.Sprintf("%v", data)
}

func RegReplace(data string, query string, value string) string {
	reg := regexp.MustCompile(query)
	return reg.ReplaceAllString(data, value)
}

func HasKey(data map[string]string, key string) bool {
	_, isOK := data[key]
	return isOK
}

func ParseFloat(value string) float64 {
	val, _ := strconv.ParseFloat(value, 64)
	return val
}

func RegIsMatch(value string, query string) bool {
	// reg := regexp.MustCompile(query)
	matched, _ := regexp.MatchString(query, value)
	return matched
}

func RegMatchAll(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.FindAllString(value, -1)
}

func SetNodeValue(xmlValue, nodeName, value string) string {
	reg := regexp.MustCompile(`>.*</` + nodeName + `>`)
	return reg.ReplaceAllString(xmlValue, ">"+value+"</"+nodeName+">")
}

func GetNodeValue(xmlValue, nodeName string) string {
	beginNode := strings.Index(xmlValue, "<"+nodeName+">") + len(nodeName) + 2
	endNode := strings.Index(xmlValue, "</"+nodeName+">")
	sValue := ""
	if beginNode != -1 && endNode != -1 {
		sValue = Substr(xmlValue, beginNode, endNode)
	}
	return sValue
}

func Substr(str string, start, end int) string {
	if start == -1 || end == -1 {
		return ""
	}
	return string(str[start:end])
}

func InitConst(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"mulubianma", "B0110_新疆国寿理赔_数据库编码对应表", "1", "3"},
		{"sanfang", "B0110_新疆国寿理赔_第三方出具单位", "2", "2"},
		{"ICD10Code", "B0110_新疆国寿理赔_ICD10疾病编码", "0", "1"},
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
