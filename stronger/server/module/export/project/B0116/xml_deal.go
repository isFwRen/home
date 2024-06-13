package B0116

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/load/model"
	"strconv"
	"strings"

	"github.com/wxnacy/wgo/arrays"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0116:::XmlDeal-----------------------")
	obj := o.(FormatObj)
	fields := obj.Fields

	issueMessage := RegMatchAll(xmlValue, `<issueMessage>[\s\S]*?<\/issueMessage>`)[0]
	errorData := ""
	errMeses := ""
	count := 0
	for _, filed := range fields {
		// if len(filed.Issues ) > 0 {
		for _, issue := range filed.Issues {
			errMes := issue.Message + ";"
			issueMessage_new := issueMessage
			if strings.Index(errMeses, errMes) == -1 {
				count++
				issueMessage_new = SetNodeValue(issueMessage_new, "issueCode", issue.Code)
				issueMessage_new = SetNodeValue(issueMessage_new, "issueSource", "BPO")
				issueMessage_new = SetNodeValue(issueMessage_new, "issueDescription", issue.Message)
				// errorData = errorData + "  <errorList>\n   <errorCode>" + issue.Code + "</errorCode>\n   <errorDesc>" + issue.Message + "</errorDesc>\n  </errorList>\n"
				errorData += "   " + issueMessage_new + "\n"
				errMeses += errMes
			}
		}
		// error = error + "  <errorList>\n   <errorCode>#{field.final_fields.issue.id}</errorCode>\n   <errorDesc>#{field.final_fields.issue.text}</errorDesc>\n  </errorList>\n"
		// }
	}

	if errorData != "" {
		xmlValue = RegReplace(xmlValue, issueMessage, errorData)
	}
	xmlValue = SetNodeValue(xmlValue, "count", strconv.Itoa(count))

	if obj.Bill.InsuranceType == "0" || obj.Bill.InsuranceType == "1" {
		xmlValue = SetNodeValue(xmlValue, "rgtBpoClass", obj.Bill.InsuranceType)
	}

	xmlValue = SetNodeValue(xmlValue, "grpRgtNo", obj.Bill.BatchNum)
	xmlValue = SetNodeValue(xmlValue, "rgtNo", obj.Bill.BillNum)

	for _, field := range fields {
		if field.Code == "fc003" && field.FinalValue != "" {
			xmlValue = SetNodeValue(xmlValue, "customerName", field.FinalValue)
			break
		}
	}

	for _, field := range fields {
		if field.Code == "fc185" && field.FinalValue != "" {
			xmlValue = SetNodeValue(xmlValue, "payeeName", field.FinalValue)
			break
		}
	}

	_, fc014 := GetOneField(fields, "fc014", true)
	_, fc013 := GetOneField(fields, "fc013", true)
	xmlValue = SetNodeValue(xmlValue, "appAccReason", fc014+fc013)

	fc086 := GetSameField(fields, "fc086", false)
	items := [][]string{{"12", "caseInfo"}, {"11", "operations_G"}, {"14", "operations_N"}, {"15", "operations_H"}}
	// , {"13", "operations_S"}
	for _, item := range items {
		if arrays.Contains(fc086, item[0]) == -1 {
			xmlValue = RegReplace(xmlValue, `<`+item[1]+`>[\s\S]*?<\/`+item[1]+`>`, "")
		}
	}

	if arrays.Contains(fc086, "13") == -1 {
		_, fc013 := GetOneField(fields, "fc013", false)
		if fc013 == "8" {
			nodes := RegMatchAll(xmlValue, `<operations_S>[\s\S]*?<\/operations_S>`)
			if len(nodes) > 0 {
				node := nodes[0]
				_, fc035 := GetOneField(fields, "fc035", true)
				_, fc017 := GetOneField(fields, "fc017", true)
				node_new := node
				node_new = SetNodeValue(node_new, "unitName", fc035)
				node_new = SetNodeValue(node_new, "diagnoseDate", fc017)
				node_new = SetNodeValue(node_new, "feeType", "SP04")
				node_new = SetNodeValue(node_new, "appraiseDate", fc017)

				xmlValue = strings.Replace(xmlValue, node, node_new, 1)
			}
		} else {
			xmlValue = RegReplace(xmlValue, `<operations_S>[\s\S]*?<\/operations_S>`, "")
		}
	}

	items = [][]string{{"viewDesc", "fc044", "fc096", "fc099"}, {"sViewDesc", "fc041", "fc102", "fc105"}}
	for _, item := range items {
		itemVal := ""
		for ii, code := range item {
			if ii == 0 {
				continue
			}
			_, aa := GetOneField(fields, code, false)
			if RegIsMatch(aa, `^(|A|F)$`) || RegIsMatch(aa, `\?`) {
				if itemVal != "" {
					itemVal += ","
				}
				itemVal += aa
			}
		}
		if itemVal != "" {
			xmlValue = SetNodeValue(xmlValue, item[0], itemVal)
		} else {
			xmlValue = SetNodeValue(xmlValue, item[0], "无")
		}
	}

	treatType0 := false
	treatType1 := false
	for _, field := range fields {
		if field.Code == "fc111" {
			if RegIsMatch(field.ResultValue, `^(1|3|6|8|10)$`) {
				treatType0 = true
			}
			if RegIsMatch(field.ResultValue, `^(2|4|5|7|9|11)$`) {
				treatType1 = true
			}
		}
	}
	if treatType1 && treatType0 {
		xmlValue = SetNodeValue(xmlValue, "treatType", "02")
	} else if treatType0 {
		xmlValue = SetNodeValue(xmlValue, "treatType", "00")
	} else if treatType1 {
		xmlValue = SetNodeValue(xmlValue, "treatType", "01")
	} else {
		xmlValue = SetNodeValue(xmlValue, "treatType", "99")
	}

	feemain_As := RegMatchAll(xmlValue, `<feemain_A>[\s\S]*?<\/feemain_A>`)
	feemainMap := map[string]string{"1": "A", "2": "B", "3": "C", "4": "D", "5": "K", "6": "H", "7": "I", "8": "J", "9": "E", "10": "F", "11": "G"}
	for _, feemain_A := range feemain_As {
		feemain_A_old := feemain_A
		billNo := GetNodeValue(feemain_A, "billNo")
		val := feemainMap[billNo]
		feemain_A = strings.ReplaceAll(feemain_A, "feemain_A", "feemain_"+val)
		nodes := RegMatchAll(feemain_A, `.+`)
		items = [][]string{
			{"endDate", "^(B|D|E|G|I|K)$"},
			{"dayCount", "^(B|D|E|G|I|K)$"},
			{"medicallyPayment", "^(C|D|E|J)$"},
			{"insurancePayment", "^(C|D|E|J)$"},
			{"bjPayA", "^(C|D|E|F|G|J)$"},
			{"bjPayB", "^(C|D|E|F|G|J)$"},
			{"bjSelf", "^(C|D|E|F|G|J)$"},
			{"planFundPayment", "^(D|E|F|G|J)$"},
			{"planFundCountPayment", "^(F|G)$"},
			{"hugeCapitalPay", "^(E|F|G|J)$"},
			{"hugeCapitalCountPay", "^(F|G)$"},
			{"hugeOutPatientPay", "^(C)$"},
			{"hugeOutPatientYearCountPay", "^(C)$"},
			{"hgueOutPatientYearBalance", "^(C)$"},
			{"selfCashPayBalance", "^(C|D|E|J)$"},
			{"selfAccountPayBalance", "^(C|D|E|J)$"},
			{"selfAccountBalance", "^(C|D|E|J)$"},
			{"beginPay", "^(C|D|E)$"},
			{"hugeTopPay", "^C|E)$"},
			{"selfPayCount", "^(C|D)$"},
			{"retireSafePay", "^(C|E|J)$"},
			{"govermentSafePay", "^(C|D|E|J)$"},
			{"policySafePay", "^(C|E|J)$"},
			{"safeInrangeCountPay", "^(C|D)$"},
			{"safeOutRangePay", "^(D)$"},
			{"mutualityPayment", "^(D|E|J)$"},
			{"individualPayment", "^(C|E|J)$"},
			{"planFundCountPaymentOne", "^(E)$"},
			{"deliverPay", "^(H|I)$"},
			{"Specialbill", "^(A|B)$"},
		}
		startDate := 1
		for _, node := range nodes {
			if strings.Index(node, "<startDate>") != -1 {
				if startDate == 1 {
					if !RegIsMatch(val, `(A|C|F|H|J)`) {
						feemain_A = strings.Replace(feemain_A, node, "", 1)
					}
				} else {
					if !RegIsMatch(val, `(B|D|E|G|I|K)`) {
						feemain_A = strings.Replace(feemain_A, node, "", 1)
					}
				}
				startDate++
			}
			for _, item := range items {
				if strings.Index(node, "<"+item[0]+">") != -1 {
					if !RegIsMatch(val, item[1]) {
						feemain_A = strings.ReplaceAll(feemain_A, node, "")
					}
				}
			}
		}
		if val == "J" {
			feemain_A = strings.ReplaceAll(feemain_A, "<bjSelf>", "<ownEpense>")
			feemain_A = strings.ReplaceAll(feemain_A, "</bjSelf>", "</ownEpense>")
		}

		xmlValue = strings.Replace(xmlValue, feemain_A_old, feemain_A, 1)
	}

	segmentFlag := GetNodeValue(xmlValue, "segmentFlag")
	if segmentFlag == "" {
		xmlValue = SetNodeValue(xmlValue, "segmentFlag", "0")
	}

	addNodes := []AddNodes{
		{
			Node:  "caseInfo",
			Code:  "fc027",
			Items: [][]string{{"defoType", "fc026"}, {"defoGrade", "fc028"}, {"defoCode", "fc027"}, {"realRate", "fc030"}, {"acceptCom", "fc031"}, {"judgeDate", "fc032"}},
		},
		{
			Node:  "operations_S",
			Code:  "fc109",
			Items: [][]string{{"unitName", "fc035"}, {"diagnoseDate", "fc036"}, {"feeType", "fc109"}, {"appraiseDate", "fc032"}},
		},
		{
			Node:  "operations_N",
			Code:  "fc110",
			Items: [][]string{{"feeType", "fc110"}, {"unitName", "fc035"}, {"diagnoseDate", "fc036"}, {"appraiseDate", "fc032"}, {"disStartDate", "fc216"}, {"disEndDate", "fc217"}, {"MonthlySal", "fc218"}},
		},
		{
			Node:  "operations_H",
			Code:  "fc219",
			Items: [][]string{{"feeType", "fc219"}, {"unitName", "fc035"}, {"appraiseHospital", "fc031"}, {"appraiseDate", "fc032"}, {"StartDate", "fc220"}, {"EndDate", "fc221"}},
		},
	}
	for _, addNode := range addNodes {
		nodes := RegMatchAll(xmlValue, `<`+addNode.Node+`>[\s\S]*?<\/`+addNode.Node+`>`)
		if len(nodes) > 0 {
			node := nodes[0]
			nodeData := ""
			for _, field := range fields {
				if field.Code == addNode.Code && field.FinalValue != "" {
					node_new := node
					for _, item := range addNode.Items {
						node_new = SetNodeValue(node_new, item[0], GetBlockField(fields, item[1], field.BlockID, true))
					}
					nodeData += "       " + node_new + "/n"
				}
			}
			xmlValue = strings.Replace(xmlValue, node, nodeData, 1)
		}
	}

	nodes := []string{"appAmnt", "billMoney", "socialPay", "biPayment", "thirdPay", "selfFee", "medicallyPayment", "insurancePayment", "bjPayA", "bjPayB", "bjSelf", "planFundPayment", "planFundCountPayment", "hugeCapitalPay", "hugeCapitalCountPay", "hugeOutPatientPay", "hugeOutPatientYearCountPay", "hgueOutPatientYearBalance", "selfCashPayBalance", "selfAccountPayBalance", "selfAccountBalance", "beginPay", "hugeTopPay", "selfPayCount", "retireSafePay", "govermentSafePay", "policySafePay", "safeInrangeCountPay", "safeOutRangePay", "mutualityPayment", "individualPayment", "planFundCountPaymentOne", "selfFeeDrug", "selfFeeCheck", "selfFeeTreat", "otSelfFee", "wMoney", "cpMoney", "chMoney", "edMoney", "exaMoney", "exMoney", "dtMoney", "raeMoney", "seMoney", "meMoney", "opeMoney", "nrMoney", "matMoney", "nurMoney", "accMoney", "salMoney", "ambMoney", "hpMoney", "dpMoney", "otMoney"}
	for _, node := range nodes {
		val := GetNodeValue(xmlValue, node)
		if val == "" {
			xmlValue = SetNodeValue(xmlValue, node, "0.00")
		}
	}

	// nodes = RegMatchAll(xmlValue, `<operations_S>[\s\S]*?<\/operations_S>`)
	// if len(nodes) > 0 {
	// 	node := nodes[0]
	// 	nodeData := ""
	// 	for _, field := range fields {
	// 		if field.Code == "fc109" && field.FinalValue != "" {
	// 			node_new := node
	// 			node_new = SetNodeValue(node_new, "unitName", GetBlockField(fields, "fc035", field.BlockID, true))
	// 			node_new = SetNodeValue(node_new, "diagnoseDate", GetBlockField(fields, "fc036", field.BlockID, true))
	// 			node_new = SetNodeValue(node_new, "feeType", field.FinalValue)
	// 			node_new = SetNodeValue(node_new, "appraiseDate", GetBlockField(fields, "fc032", field.BlockID, true))
	// 			nodeData += "       " + node_new + "/n"
	// 		}
	// 	}
	// 	xmlValue = strings.Replace(xmlValue, node, nodeData, 1)
	// }

	return err, xmlValue
}

type AddNodes struct {
	Node  string     `json:"node"`
	Code  string     `json:"code"`
	Items [][]string `json:"items"`
}

func GetBlockField(fields []model.ProjectField, code, blockID string, finalOrResult bool) string {
	for _, field := range fields {
		if field.Code == code && field.BlockID == blockID {
			if finalOrResult {
				return field.FinalValue
			} else {
				return field.ResultValue
			}
		}
	}
	return ""
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

func GetSameField(fields []model.ProjectField, code string, finalOrResult bool) (values []string) {
	for _, field := range fields {
		if field.Code == code {
			if finalOrResult {
				values = append(values, field.FinalValue)
			} else {
				values = append(values, field.ResultValue)
			}
		}
	}
	return values
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
