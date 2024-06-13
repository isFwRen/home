/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年03月22日16:35:50
 */

package B0113

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/load/model"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0113:::XmlDeal-----------------------")
	obj := o.(FormatObj)
	fields := obj.Fields
	bill_name := obj.Bill.BillName
	bill_name_Arr := strings.Split(bill_name, "_")
	xmlValue = SetNodeValue(xmlValue, "BPOBatNo", bill_name_Arr[0])
	xmlValue = SetNodeValue(xmlValue, "RgtNo", bill_name_Arr[2])
	xmlValue = SetNodeValue(xmlValue, "SeqNo", bill_name_Arr[3])

	xmlValue = strings.Replace(xmlValue, "utf-8", "GBK", 1)

	c, ok := global.GProConf[obj.Bill.ProCode].ConstTable["B0113_百年理赔_百年理赔医院代码表"]
	yiyuan := make(map[string]string, 0)
	if ok {
		fmt.Println("------------------------c-----------------------:", c[0])
		for _, arr := range c {
			yiyuan[strings.TrimSpace(arr[0])] = arr[1]
		}
	}

	questionCount := 0
	errMeses := ""
	for _, field := range fields {
		// if len(filed.Issues ) > 0 {
		for _, issue := range field.Issues {
			fname := field.Name
			issue.Message = strings.Replace(issue.Message, fname, "", 1)
			errMes := fname + issue.Message + ";"
			if strings.Index(errMeses, errMes) == -1 {
				errMeses += errMes
				questionCount++
			}
		}
	}
	xmlValue = SetNodeValue(xmlValue, "StateDesc", errMeses)

	// _, fc161 := GetOneField(fields, "fc161", false)
	// state := questionCount
	// if fc161 == "6" && questionCount > 3 {
	// 	state = 1
	// }
	// xmlValue = SetNodeValue(xmlValue, "State", strconv.Itoa(state))

	// _, fc019 := GetOneField(fields, "fc019", false)

	// if RegIsMatch(fc019, `(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|x|X)$)`) {
	// 	assigneeSex := GetNodeValue(xmlValue, "AssigneeSex")
	// 	if assigneeSex == "0" {
	// 		xmlValue = SetNodeValue(xmlValue, "AssigneeSex", "男")
	// 	} else if assigneeSex == "1" {
	// 		xmlValue = SetNodeValue(xmlValue, "AssigneeSex", "女")
	// 	}

	// }

	_, fc020 := GetOneField(fields, "fc020", false)
	if fc020 == "1" {
		xmlValue = SetNodeValue(xmlValue, "AssigneeIDVailType", "1")
	} else if !RegIsMatch(fc020, `^(1|A|)$`) {
		xmlValue = SetNodeValue(xmlValue, "AssigneeIDVailType", "0")
	}

	lLFeeMainSchemas := RegMatchAll(xmlValue, `<LLFeeMainSchema>[\s\S]*?<\/LLFeeMainSchema>`)
	for ll, lLFeeMainSchema := range lLFeeMainSchemas {
		old_lLFeeMainSchema := lLFeeMainSchema
		// selfItemAmnt := 0.00
		// selfDrugAmnt := 0.00
		// SelfAmntNum := 0.00
		FeeAffixType := "1"

		FeeAmnt := GetNodeValue(lLFeeMainSchema, "FeeAmnt")
		lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "FeeRealAmnt", FeeAmnt)
		MainFeeNo := GetNodeValue(lLFeeMainSchema, "MainFeeNo")
		if MainFeeNo == "" {
			xmlValue = strings.Replace(xmlValue, old_lLFeeMainSchema, "", 1)
			continue
		}
		BPONumPages := 0
		cachellLCaseReceiptSchema := ""
		for _, MedicalBillInfo := range obj.Invoice {
			_, fc063 := GetOneField(MedicalBillInfo.Fields, "fc063", true)
			fmt.Println("-------------MainFeeNo--------------------", ll, MainFeeNo, fc063)
			if fc063 == MainFeeNo {
				fmt.Println("-------------MainFeeNo--------------------", MainFeeNo)
				_, fc062 := GetOneField(MedicalBillInfo.Fields, "fc062", false)
				_, fc068 := GetOneField(MedicalBillInfo.Fields, "fc068", true)
				_, fc069 := GetOneField(MedicalBillInfo.Fields, "fc069", true)
				// RealHospDate := GetNodeValue(lLFeeMainSchema, "RealHospDate")
				// lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "ADJHospDate", RealHospDate)

				if fc068 != "" && fc069 != "" {
					hh := ""
					if fc062 == "2" {
						a, _ := time.Parse("2006-01-02", fc068)
						b, _ := time.Parse("2006-01-02", fc069)
						d := b.Sub(a)
						hh = ToString(d.Hours() / 24)
					}
					if (fc062 == "1" || fc062 == "2") && fc068 == fc069 {
						hh = "1"
					}

					// fmt.Println(d.Hours() / 24)
					fmt.Println("-------------ADJHospDate  RealHospDate--------------------", hh)
					lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "RealHospDate", hh)
					lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "ADJHospDate", hh)
				}
				for _, QingDan := range MedicalBillInfo.QingDan {
					for _, field := range QingDan.Fields {
						if field.Code == "fc109" {
							fmt.Println("-------------fc109fc109fc109--------------------", BPONumPages)
							BPONumPages++
						}
					}
				}
				_, fc065 := GetOneField(MedicalBillInfo.Fields, "fc065", true)
				HospitalCode, isExits := yiyuan[fc065]
				fmt.Println("-------------fc065fc065fc065--------------------", fc065, HospitalCode)
				if isExits {
					lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "HospitalName", HospitalCode)
				}
				// if ll == len(lLFeeMainSchemas)-1 {
				cccs := [][]string{
					{"fc259", "fc258", "fc268"},
					{"fc276", "fc277", "fc270"},
					{"fc278", "fc279", "fc272"},
					{"fc280", "fc281", "fc274"},
				}
				FeeItemName2 := map[string]string{
					"CS00101": "社会医疗保险给付",
					"CS00102": "公费医疗保险给付",
					"CS00103": "商业医疗保险给付",
					"CS00104": "其他",
				}
				for _, ccc := range cccs {
					_, cc0 := GetOneField(MedicalBillInfo.Fields, ccc[0], true)
					_, cc1 := GetOneField(MedicalBillInfo.Fields, ccc[1], true)
					if cc0 == "" || cc0 == "0.00" {
						cc0 = cc1
					}
					if (cc0 != "" && cc0 != "0.00") || (cc1 != "" && cc1 != "0.00") {
						_, cc2 := GetOneField(MedicalBillInfo.Fields, ccc[2], true)
						FeeAffixType = "0"
						cachellLCaseReceiptSchema += "\n\t\t\t\t<LLCaseReceiptSchema>\n\t\t\t\t\t<FeeItemCode>CS001</FeeItemCode>\n\t\t\t\t\t<FeeItemName>第三方给付</FeeItemName>\n\t\t\t\t\t<FeeItemCode2>" + cc2 + "</FeeItemCode2>\n\t\t\t\t\t<FeeItemName2>" + FeeItemName2[cc2] + "</FeeItemName2>\n\t\t\t\t\t<Fee>" + cc0 + "</Fee>\n\t\t\t\t\t<AdjSum>" + cc0 + "</AdjSum>\n\t\t\t\t\t<FeeAmnt></FeeAmnt>\n\t\t\t\t\t<FeeSelfPercent></FeeSelfPercent>\n\t\t\t\t</LLCaseReceiptSchema>\n"

					}
				}
				// }
				lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "BPONumPages", strconv.Itoa(BPONumPages))
				break
			}

		}

		SelfAmnt := GetNodeValue(lLFeeMainSchema, "SelfAmnt")
		idx := strings.LastIndex(lLFeeMainSchema, "</LLCaseReceiptSchema>")
		if idx != -1 {
			idx = idx + len("</LLCaseReceiptSchema>")
			if SelfAmnt != "" {
				cachellLCaseReceiptSchema += "\n\t\t\t\t<LLCaseReceiptSchema>\n\t\t\t\t\t<FeeItemCode>CZ001</FeeItemCode>\n\t\t\t\t\t<FeeItemName>自费项目</FeeItemName>\n\t\t\t\t\t<FeeItemCode2></FeeItemCode2>\n\t\t\t\t\t<FeeItemName2></FeeItemName2>\n\t\t\t\t\t<Fee>" + SelfAmnt + "</Fee>\n\t\t\t\t\t<AdjSum>" + SelfAmnt + "</AdjSum>\n\t\t\t\t\t<FeeAmnt></FeeAmnt>\n\t\t\t\t\t<FeeSelfPercent></FeeSelfPercent>\n\t\t\t\t</LLCaseReceiptSchema>\n"
			}
			lLFeeMainSchema = Substr(lLFeeMainSchema, 0, idx) + cachellLCaseReceiptSchema + string(lLFeeMainSchema[idx:])
		}

		// lLCaseReceiptSchemas := RegMatchAll(lLFeeMainSchema, `<LLCaseReceiptSchema>[\s\S]*?<\/LLCaseReceiptSchema>`)
		// for rr, lLCaseReceiptSchema := range lLCaseReceiptSchemas {
		// 	old_lLCaseReceiptSchema := lLCaseReceiptSchema
		// 	// feeItemCode := GetNodeValue(lLCaseReceiptSchema, "FeeItemCode")
		// 	// fee := GetNodeValue(lLCaseReceiptSchema, "Fee")
		// 	// // SelfAmnt := GetNodeValue(lLFeeMainSchema, "SelfAmnt")
		// 	// if fee != "" {
		// 	// 	if strings.Index(feeItemCode, "CZ001") != -1 {
		// 	// 		selfItemAmnt = SumFloat(selfItemAmnt, ParseFloat(fee), "+")
		// 	// 		SelfAmntNum = SumFloat(SelfAmntNum, ParseFloat(fee), "+")
		// 	// 	} else if strings.Index(feeItemCode, "CZ002") != -1 {
		// 	// 		selfDrugAmnt = SumFloat(selfDrugAmnt, ParseFloat(fee), "+")
		// 	// 		SelfAmntNum = SumFloat(SelfAmntNum, ParseFloat(fee), "+")
		// 	// 	}
		// 	// }
		// 	// if strings.Index(feeItemCode, "CS001") != -1 {
		// 	// 	FeeAffixType = "0"
		// 	// }
		// 	if rr == len(lLCaseReceiptSchemas)-1 {
		// 		lLCaseReceiptSchema += cachellLCaseReceiptSchema
		// 		if SelfAmnt != "" {
		// 			lLCaseReceiptSchema += "\n\t\t\t\t<LLCaseReceiptSchema>\n\t\t\t\t\t<FeeItemCode>CZ001</FeeItemCode>\n\t\t\t\t\t<FeeItemName>自费项目</FeeItemName>\n\t\t\t\t\t<FeeItemCode2></FeeItemCode2>\n\t\t\t\t\t<FeeItemName2></FeeItemName2>\n\t\t\t\t\t<Fee>" + SelfAmnt + "</Fee>\n\t\t\t\t\t<AdjSum>" + SelfAmnt + "</AdjSum>\n\t\t\t\t\t<FeeAmnt></FeeAmnt>\n\t\t\t\t\t<FeeSelfPercent></FeeSelfPercent>\n\t\t\t\t</LLCaseReceiptSchema>\n"
		// 		}
		// 		lLFeeMainSchema = strings.Replace(lLFeeMainSchema, old_lLCaseReceiptSchema, lLCaseReceiptSchema, 1)
		// 	}
		// }
		// if SelfAmnt == "" {
		// 	lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "SelfAmnt", decimal.NewFromFloat(SelfAmntNum).StringFixed(2))
		// }

		lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "FeeAffixType", FeeAffixType)
		// lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "SelfDrugAmnt", decimal.NewFromFloat(selfDrugAmnt).StringFixed(2))
		// lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "SelfItemAmnt", decimal.NewFromFloat(selfItemAmnt).StringFixed(2))

		fmt.Println("-------------llllllll--------------------", ll)

		// if strings.Index(feeItemCode, "CS001") != -1 {
		// 	lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "FeeAffixType", "0")
		// } else {
		// 	lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "FeeAffixType", "1")
		// }
		xmlValue = strings.Replace(xmlValue, old_lLFeeMainSchema, lLFeeMainSchema, 1)
	}

	lLFeeMainSchemas = RegMatchAll(xmlValue, `<LLFeeMainSchema>[\s\S]*?<\/LLFeeMainSchema>`)
	for _, lLFeeMainSchema := range lLFeeMainSchemas {
		old_lLFeeMainSchema := lLFeeMainSchema
		selfItemAmnt := 0.00
		selfDrugAmnt := 0.00
		SelfAmntNum := 0.00
		SelfAmnt := GetNodeValue(lLFeeMainSchema, "SelfAmnt")
		RealHospDate := GetNodeValue(lLFeeMainSchema, "RealHospDate")
		lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "ADJHospDate", RealHospDate)
		HospitalCode := GetNodeValue(lLFeeMainSchema, "HospitalCode")
		if HospitalCode == "" {
			lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "HospitalCode", "370324")
		}

		lLCaseReceiptSchemas := RegMatchAll(lLFeeMainSchema, `<LLCaseReceiptSchema>[\s\S]*?<\/LLCaseReceiptSchema>`)
		for _, lLCaseReceiptSchema := range lLCaseReceiptSchemas {
			feeItemCode := GetNodeValue(lLCaseReceiptSchema, "FeeItemCode")
			fee := GetNodeValue(lLCaseReceiptSchema, "Fee")
			// SelfAmnt := GetNodeValue(lLFeeMainSchema, "SelfAmnt")
			if fee != "" {
				if strings.Index(feeItemCode, "CZ001") != -1 {
					selfItemAmnt = SumFloat(selfItemAmnt, ParseFloat(fee), "+")
					SelfAmntNum = SumFloat(SelfAmntNum, ParseFloat(fee), "+")
				} else if strings.Index(feeItemCode, "CZ002") != -1 {
					selfDrugAmnt = SumFloat(selfDrugAmnt, ParseFloat(fee), "+")
					SelfAmntNum = SumFloat(SelfAmntNum, ParseFloat(fee), "+")
				}
			}
		}
		if SelfAmnt == "" {
			lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "SelfAmnt", decimal.NewFromFloat(SelfAmntNum).StringFixed(2))
		}
		lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "SelfDrugAmnt", decimal.NewFromFloat(selfDrugAmnt).StringFixed(2))
		lLFeeMainSchema = SetNodeValue(lLFeeMainSchema, "SelfItemAmnt", decimal.NewFromFloat(selfItemAmnt).StringFixed(2))
		xmlValue = strings.Replace(xmlValue, old_lLFeeMainSchema, lLFeeMainSchema, 1)
	}

	lLCaseReceiptSchemas := RegMatchAll(xmlValue, `<LLCaseReceiptSchema>[\s\S]*?<\/LLCaseReceiptSchema>`)
	for _, lLCaseReceiptSchema := range lLCaseReceiptSchemas {
		old_lLCaseReceiptSchema := lLCaseReceiptSchema
		Fee := GetNodeValue(lLCaseReceiptSchema, "Fee")
		lLCaseReceiptSchema = SetNodeValue(lLCaseReceiptSchema, "AdjSum", Fee)

		FeeItemCode := GetNodeValue(lLCaseReceiptSchema, "FeeItemCode")
		// FeeItemName := GetNodeValue(lLCaseReceiptSchema, "FeeItemName")  && FeeItemName == ""
		FeeSelfPercent := GetNodeValue(lLCaseReceiptSchema, "FeeSelfPercent")

		if FeeItemCode == "" && (Fee == "" || Fee == "0.00") && FeeSelfPercent == "" {
			xmlValue = strings.Replace(xmlValue, old_lLCaseReceiptSchema, "", 1)
		} else if FeeSelfPercent == "0.00" {
			xmlValue = strings.Replace(xmlValue, old_lLCaseReceiptSchema, "", 1)
		} else {
			xmlValue = strings.Replace(xmlValue, old_lLCaseReceiptSchema, lLCaseReceiptSchema, 1)
		}

	}

	LLBnfSchemas := RegMatchAll(xmlValue, `<LLBnfSchema>[\s\S]*?<\/LLBnfSchema>`)
	for _, LLBnfSchema := range LLBnfSchemas {
		old_LLBnfSchema := LLBnfSchema
		PayeeNonNatural := GetNodeValue(LLBnfSchema, "PayeeNonNatural")
		BankCode := GetNodeValue(LLBnfSchema, "BankCode")
		if PayeeNonNatural == "" {
			LLBnfSchema = SetNodeValue(LLBnfSchema, "PayeeNonNatural", "1")
		}
		if RegIsMatch(BankCode, `^(01|02|03|04|06|17|12|09|14|13|10|05|11|07)$`) {
			LLBnfSchema = SetNodeValue(LLBnfSchema, "CasePayMode", "r")
		} else {
			LLBnfSchema = SetNodeValue(LLBnfSchema, "CasePayMode", "4")
		}
		xmlValue = strings.Replace(xmlValue, old_LLBnfSchema, LLBnfSchema, 1)
	}

	// realHospDate := GetNodeValue(xmlValue, "RealHospDate")
	// xmlValue = SetNodeValue(xmlValue, "RealHospDate", realHospDate)

	return err, xmlValue
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
