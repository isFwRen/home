/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年04月23日16:24:35
 */

package B0103

import (
	"encoding/json"
	"fmt"
	"regexp"
	"server/global"
	"server/module/load/model"
	"server/utils"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0103:::XmlDeal-----------------------")
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
	agency := obj.Bill.Agency
	// constMap := InitConst(obj.Bill.ProCode, agency)

	// mulubianma := constMap["mulubianma"]
	v, _ := utils.FetchConst(obj.Bill.ProCode, "B0103_广西贵州国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": agency})
	xmlValue = SetNodeValue(xmlValue, "catalogCode", v)

	questionCount := 0
	errMeses := ""
	fc003s := GetSameField(obj.Fields, "fc003", false)
	if arrays.Contains(fc003s, "4") == -1 {
		xmlValue = SetNodeValue(xmlValue, "mIcd10Code", "S32")
		errMeses += "案件无疾病诊断;"
		questionCount++
	}

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
		socialPayAmnt := GetNodeValue(rcptInfoList, "socialPayAmnt")
		if RegIsMatch(socialPayAmnt, `^(|0|0\.00)$`) {
			rcptInfoList = SetNodeValue(rcptInfoList, "siHealthType", "00")
		}
		if RegIsMatch(rcptInfoList, `<(rcptList|errorRcptList)>`) {
			rcptInfoList = SetNodeValue(rcptInfoList, "rcptType", "D")
		} else {
			rcptInfoList = SetNodeValue(rcptInfoList, "rcptType", "M")
		}
		if GetNodeValue(rcptInfoList, "siHealthType") == "00" && strings.Index(errMeses, "医保身份未填写内容;") == -1 {
			if !RegIsMatch(GetNodeValue(rcptInfoList, "socialPayAmnt"), `^(|0|0\.00)$`) || !RegIsMatch(GetNodeValue(rcptInfoList, "accountPayAmnt"), `^(|0|0\.00)$`) {
				errMeses += "医保身份未填写内容;"
				questionCount++
			}
		}
		xmlValue = strings.Replace(xmlValue, old_rcptInfoList, rcptInfoList, 1)
	}

	//CSB0103RC0076001 问题件描述重复市，保留一个即可
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

	errorRcpListCache := ""

	items := [][]string{
		{"rcptList", "listCode"},
		{"errorRcptList", "listName"},
	}
	for _, item := range items {
		nodes := RegMatchAll(xmlValue, `<`+item[0]+`>[\s\S]*?<\/`+item[0]+`>`)
		for _, node := range nodes {
			val1 := GetNodeValue(node, item[1])
			if item[0] == "errorRcptList" && errorRcpListCache == "" {
				errorRcpListCache = node
			}
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

	rcptInfoLists = RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)

	mIcd10CodeValue := make(map[string]string, 0)
	for _, invoice := range obj.Invoice {
		_, fc005 := GetOneField(invoice.Fields, "fc005", true)
		mIcd10CodeValue[fc005] = ""
		for _, field := range invoice.Fields {
			if RegIsMatch(field.Code, `^(fc007|fc074|fc075|fc076|fc077|fc078|fc079|fc080|fc081|fc082)$`) && field.FinalValue != "" {
				if mIcd10CodeValue[fc005] != "" {
					mIcd10CodeValue[fc005] += ","
				}
				mIcd10CodeValue[fc005] += field.FinalValue
			}
		}
	}

	fmt.Println("--------mIcd10CodeValuemIcd10CodeValue-----------", mIcd10CodeValue)
	qqq := ""
	for _, rcptInfoList := range rcptInfoLists {
		rcptNo := GetNodeValue(rcptInfoList, "rcptNo")
		if mIcd10CodeValue[rcptNo] != "" {
			qqq = mIcd10CodeValue[rcptNo]
			break
		}
	}
	for _, rcptInfoList := range rcptInfoLists {
		old_rcptInfoList := rcptInfoList
		rcptNo := GetNodeValue(rcptInfoList, "rcptNo")
		fmt.Println("--------rcptNo-----------", rcptNo, mIcd10CodeValue[rcptNo])
		if mIcd10CodeValue[rcptNo] != "" {
			rcptInfoList = SetNodeValue(rcptInfoList, "mIcd10Code", mIcd10CodeValue[rcptNo])
		} else if qqq != "" {
			rcptInfoList = SetNodeValue(rcptInfoList, "mIcd10Code", qqq)
		}

		rcptAmnt := GetNodeValue(rcptInfoList, "rcptAmnt")
		rcptLists := RegMatchAll(rcptInfoList, `<(rcptList|errorRcptList)>[\s\S]*?<\/(rcptList|errorRcptList)>`)
		sum := 0.0
		for rr, rcptList := range rcptLists {
			old_rcptList := rcptList
			quantity := GetNodeValue(rcptList, "quantity")
			price := GetNodeValue(rcptList, "price")
			if quantity != "" && price != "" {
				total := SumFloat(ParseFloat(quantity), ParseFloat(price), "*")
				sum = SumFloat(sum, total, "+")
			}
			if rr == len(rcptLists)-1 {
				fmt.Println("--------------rcptAmnt----------------------", rcptAmnt, sum)
				chae := SumFloat(ParseFloat(rcptAmnt), sum, "-")
				// expenMode := GetNodeValue(rcptInfoList, "expenMode")
				if sum != 0.0 && chae != 0 && chae != 0.00 {
					aaa := SetNodeValue(errorRcpListCache, "rcptListNo", strconv.Itoa(len(rcptLists)+1))
					aaa = SetNodeValue(aaa, "listName", "补差金额")
					aaa = SetNodeValue(aaa, "quantity", "1")
					aaa = SetNodeValue(aaa, "liablePct", "1")
					aaa = SetNodeValue(aaa, "price", ToString(chae))
					rcptInfoList = strings.Replace(rcptInfoList, old_rcptList, rcptList+"\n       "+aaa, 1)
					// rcptInfoList += "\n   <errorRcptList>\n    <listName>补差金额</listName>\n    <price>" + ToString(chae) + "</price>\n   </errorRcptList>\n"
				}
			}
		}

		xmlValue = strings.Replace(xmlValue, old_rcptInfoList, rcptInfoList, 1)
	}

	rcptOtherPayDetailLists := RegMatchAll(xmlValue, `<rcptOtherPayDetailList>[\s\S]*?<\/rcptOtherPayDetailList>`)

	for _, rcptOtherPayDetailList := range rcptOtherPayDetailLists {
		old_rcptOtherPayDetailList := rcptOtherPayDetailList
		uedUnitNos := strings.Split(GetNodeValue(rcptOtherPayDetailList, "uedUnitNo"), "|")
		otherSpiltAmnts := strings.Split(GetNodeValue(rcptOtherPayDetailList, "otherSpiltAmnt"), "|")
		chache := ""
		for ii, value := range uedUnitNos {
			rcptOtherPayDetailList = SetNodeValue(rcptOtherPayDetailList, "uedUnitNo", value)
			rcptOtherPayDetailList = SetNodeValue(rcptOtherPayDetailList, "otherSpiltAmnt", otherSpiltAmnts[ii])
			chache += rcptOtherPayDetailList + "\n"
		}
		xmlValue = strings.Replace(xmlValue, old_rcptOtherPayDetailList, chache, 1)
	}

	// itemNodes := []string{"bankAccInfoList", "rcptInfoList", "rcptList", "errorRcptList"}

	// for _, itemNode := range itemNodes {
	// 	nodeDatas := RegMatchAll(xmlValue, `<`+itemNode+`>[\s\S]*?<\/`+itemNode+`>`)
	// 	if len(nodeDatas) == 1 {
	// 		aaa := "<" + itemNode + "><array_list></array_list></" + itemNode + ">"
	// 		cloneNodeData := nodeDatas[0]
	// 		cloneNodeData = cloneNodeData + aaa
	// 		xmlValue = strings.Replace(xmlValue, nodeDatas[0], cloneNodeData, 1)
	// 	}
	// }

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

	//CSB0103RC0123000 当rcpNo节点为空时，屏蔽rcptInfoList一整个大节点
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
	ID        string `json:"id"`        // 账单号
	ImageType string `json:"imageType"` //下载报文节点
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

func InitConst(proCode string, agency string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"mulubianma", "B0103_广西贵州国寿理赔_数据库编码对应表", "1", "3"},
		{"jibin", "B0103_广西贵州国寿理赔_ICD10疾病编码", "0", "0"},
		{"sanfang", "B0103_广西贵州国寿理赔_第三方出具单位", "1", "1"},
		{"sanfangAll", "B0103_广西贵州国寿理赔_第三方出具单位", "2", "1"},
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
	yiliaomulu, ok := constObj["mulubianma"][agency]
	tempMap := make(map[string]string, 0)
	if ok {
		c, ok := global.GProConf[proCode].ConstTable["B0103_广西贵州国寿理赔_"+yiliaomulu]
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi("2")
				v, _ := strconv.Atoi("1")
				tempMap[strings.TrimSpace(arr[k])] = arr[v]
			}
		}
	}
	constObj["yiliaomulu"] = tempMap

	return constObj
}
