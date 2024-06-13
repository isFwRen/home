/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 5:34 下午
 */

package B0118

import (
	"fmt"
	"regexp"
	"server/global"
	"server/utils"
	"sort"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/wxnacy/wgo/arrays"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0118:::XmlDeal-----------------------")
	// return nil, xmlValue
	obj := o.(FormatObj)
	xmlValue = SetNodeValue(xmlValue, "claimNum", obj.Bill.BillNum)
	constMap := constDeal(obj.Bill.ProCode)
	fileds := obj.Fields
	errorData := ""
	errMeses := ""
	changeMap := map[string][]string{
		"N0013": {"电子发票（", "）已查验、已打印。"},
		"N0014": {"电子发票（", "）已作废。"},
		"N0015": {"电子发票（", "）已开红票。"},
		//"N0017": {"电子发票（", "）未成功查验，原因：超过该张发票当日查验次数(请于次日再次查验)。"},
		"N0017": {"电子发票（", "）未成功查验，原因："}, //CSB0118RC0340000
	}

	for _, filed := range fileds {
		// if len(filed.Issues ) > 0 {
		for _, issue := range filed.Issues {
			if val, ok := changeMap[issue.Code]; ok {
				changeMap[issue.Code] = append(changeMap[issue.Code], strings.ReplaceAll(strings.ReplaceAll(issue.Message, val[0], ""), val[1], ""))
				continue
			}
			errMes := issue.Code + "_" + issue.Message + ";"
			if strings.Index(errMeses, errMes) == -1 {
				errorData = errorData + "  <errorList>\n   <errorCode>" + issue.Code + "</errorCode>\n   <errorDesc>" + issue.Message + "</errorDesc>\n  </errorList>\n"
				errMeses += errMes
			}

		}
		// error = error + "  <errorList>\n   <errorCode>#{field.final_fields.issue.id}</errorCode>\n   <errorDesc>#{field.final_fields.issue.text}</errorDesc>\n  </errorList>\n"
		// }
	}
	for issueCode, invoiceCode := range changeMap { //CSB0118RC0340000
		if len(invoiceCode) == 2 {
			continue
		}
		if changeMap[issueCode][1] == "）未成功查验，原因：" {

			if len(invoiceCode) >= 3 {
				buffStrList := []string{}
				buffStrList = append(buffStrList, invoiceCode[2:]...)
				N0017Map := make(map[string][]string)
				for _, reasonStr := range buffStrList {
					reg := regexp.MustCompile(`NO\.\d+`)
					numList := reg.FindAllString(reasonStr, -1)
					if len(numList) == 1 {
						sLen := len(numList[0])
						reasonStrLen := len(reasonStr)
						if reasonStrLen > sLen {
							N0017Map[reasonStr[sLen:]] = append(N0017Map[reasonStr[sLen:]], numList[0])
						}
					}
				}
				if len(N0017Map) == 0 {
					continue
				}
				for reason, billNumberList := range N0017Map {
					errorData = errorData +
						"  <errorList>\n   <errorCode>" + issueCode + "</errorCode>\n   <errorDesc>电子发票（" + strings.Join(billNumberList[:], ",") + "）未成功查验，原因：" + reason + "</errorDesc>\n  </errorList>\n"
				}
				continue
			}

		}
		errorData = errorData +
			"  <errorList>\n   <errorCode>" + issueCode + "</errorCode>\n   <errorDesc>电子发票（" + strings.Join(invoiceCode[2:], ",") + changeMap[issueCode][1] + "</errorDesc>\n  </errorList>\n"
	}
	fmt.Println("------------------------error-------------------", errorData)

	xmlValue = strings.Replace(xmlValue, "<errorList></errorList>", errorData, 1)

	reg := regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	// matched, _ := regexp.MatchString(`.zip$`, line)
	items := reg.FindAllString(xmlValue, -1)
	// global.GLog.Info("----------------hospitalList----------------------", zap.Any("", len(items)))
	for _, item := range items {
		new_item := item
		icd10Reg := regexp.MustCompile(`<icd10List>[\s\S]*?<\/icd10List>`)
		icd10s := icd10Reg.FindAllString(xmlValue, -1)
		for ii, icd10 := range icd10s {
			if ii == 0 {
				continue
			}
			if GetNodeValue(icd10, "icdCode") == "" {
				new_item = strings.Replace(new_item, icd10, "", 1)
			}
		}
		// xmlValue = xmlValue.replace(item, new_item)
		xmlValue = strings.Replace(xmlValue, item, new_item, 1)
	}

	//reg = regexp.MustCompile(`<siList>[\s\S]*?<\/siList>`)
	//items = reg.FindAllString(xmlValue, -1)
	//num_inSIAmount := make(map[string]string)
	//for _, item := range items {
	//	obillnumber := GetNodeValue(item, "obillnumber")
	//	inSIAmount := GetNodeValue(item, "inSIAmount")
	//	num_inSIAmount[obillnumber] = inSIAmount
	//}

	//reg = regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	//items = reg.FindAllString(xmlValue, -1)
	//for _, item := range items {
	//	new_item := item
	//	billNumber := GetNodeValue(new_item, "billNumber")
	//	reg = regexp.MustCompile(`<feeList>[\s\S]*?<\/feeList>`)
	//	feeLists := reg.FindAllString(new_item, -1)
	//	xml1299s := []map[string]string{}
	//	initialFee := 0.0
	//	feeList2s := RegMatchAll(new_item, `<feeList2>[\s\S]*?<\/feeList2>`)
	//	if len(feeLists) == 1 && RegIsMatch(feeLists[0], `<feeTypeCode>(199|299)</feeTypeCode>`) {
	//		for _, feeList2 := range feeList2s {
	//			feeTypeCode2 := GetNodeValue(feeList2, "feeTypeCode")
	//			if feeTypeCode2 == "" && GetNodeValue(feeList2, "initialFee") == "" {
	//				new_item = strings.Replace(new_item, feeList2, "", 1)
	//			} else {
	//				if !RegIsMatch(feeTypeCode2, `^(222|224|204)$`) {
	//					xml1299s = append(xml1299s, map[string]string{
	//						"initialFee": GetNodeValue(feeList2, "initialFee"),
	//						"value":      feeList2,
	//					})
	//				}
	//				if GetNodeValue(feeList2, "initialFee") != "" {
	//					aaa, _ := strconv.ParseFloat(GetNodeValue(feeList2, "initialFee"), 64)
	//					initialFee = initialFee + aaa
	//				}
	//			}
	//		}
	//		if len(xml1299s) > 0 {
	//			sort.Slice(xml1299s, func(i, j int) bool {
	//				value, ok := xml1299s[i]["initialFee"]
	//				if !ok {
	//					return false
	//				}
	//				valI, _ := strconv.ParseFloat(value, 64)
	//				value, ok = xml1299s[j]["initialFee"]
	//				if !ok {
	//					return false
	//				}
	//				valJ, _ := strconv.ParseFloat(value, 64)
	//				return valI < valJ
	//			})
	//			selfFee, _ := strconv.ParseFloat(GetNodeValue(feeLists[0], "selfFee"), 64)
	//			bbb, ok := num_inSIAmount[billNumber]
	//			if ok {
	//				fff, _ := strconv.ParseFloat(bbb, 64)
	//				selfFee = initialFee - fff
	//			}
	//			for x, xml1299 := range xml1299s {
	//				new_xml1299, _ := xml1299["value"]
	//				if x == 0 {
	//					selfFeeDesc := GetNodeValue(feeLists[0], "selfFeeDesc")
	//					new_xml1299 = SetNodeValue(new_xml1299, "selfFeeDesc", selfFeeDesc)
	//				}
	//				xxx, _ := xml1299["initialFee"]
	//				fff, _ := strconv.ParseFloat(xxx, 64)
	//				if selfFee <= fff || x == len(xml1299s)-1 {
	//					new_xml1299 = SetNodeValue(new_xml1299, "selfFee", fmt.Sprintf("%v", selfFee))
	//					new_item = strings.Replace(new_item, xml1299["value"], new_xml1299, 1)
	//					break
	//				} else {
	//					new_xml1299 = SetNodeValue(new_xml1299, "selfFee", xml1299["initialFee"])
	//					www, _ := strconv.ParseFloat(xml1299["initialFee"], 64)
	//					selfFee = selfFee - www
	//					new_item = strings.Replace(new_item, xml1299["value"], new_xml1299, 1)
	//				}
	//			}
	//			new_item = strings.Replace(new_item, feeLists[0], "", 1)
	//		} else if strings.Index(new_item, "<feeList2>") == -1 && HasKey(num_inSIAmount, billNumber) {
	//			fff := feeLists[0]
	//			initialFee, _ = strconv.ParseFloat(GetNodeValue(fff, "initialFee"), 64)
	//			www, _ := strconv.ParseFloat(num_inSIAmount[billNumber], 64)
	//			initialFee = initialFee - www
	//			fff = SetNodeValue(fff, "selfFee", fmt.Sprintf("%v", initialFee))
	//			new_item = strings.Replace(new_item, feeLists[0], fff, 1)
	//		}
	//		xmlValue = strings.Replace(xmlValue, item, new_item, 1)
	//	} else {
	//		for _, feeList2 := range feeList2s {
	//			feeTypeCode2 := GetNodeValue(feeList2, "feeTypeCode")
	//			if feeTypeCode2 == "" && GetNodeValue(feeList2, "initialFee") == "" {
	//				new_item = strings.Replace(new_item, feeList2, "", 1)
	//			}
	//		}
	//		xmlValue = strings.Replace(xmlValue, item, new_item, 1)
	//	}
	//}

	// reg = regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	// matched, _ := regexp.MatchString(`.zip$`, line)
	//items = RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	//for _, item := range items {
	//	if strings.Index(item, "") != -1 {
	//		continue
	//	}
	//	new_item := item
	//	feeLists := RegMatchAll(new_item, `<feeList>[\s\S]*?<\/feeList>`)
	//	momey := 0.0
	//	inifeeLists := []map[string]string{}
	//	isone := false
	//	for _, feeList := range feeLists {
	//		new_feeList := feeList
	//		selfFee := GetNodeValue(feeList, "selfFee")
	//		feeTypeCode := GetNodeValue(feeList, "feeTypeCode")
	//		if ParseFloat(selfFee) > 0 {
	//			if isone {
	//				inifeeLists = []map[string]string{}
	//				isone = false
	//				break
	//			} else {
	//				isone = true
	//			}
	//			initialFee := GetNodeValue(feeList, "initialFee")
	//			if RegIsMatch(feeTypeCode, `^(222|224|204)$`) {
	//				momey = ParseFloat(selfFee)
	//				new_feeList = SetNodeValue(new_feeList, "selfFee", "0")
	//				new_item = strings.Replace(new_item, feeList, new_feeList, 1)
	//
	//			} else if ParseFloat(selfFee) > ParseFloat(initialFee) {
	//				momey = ParseFloat(selfFee)
	//				new_feeList = SetNodeValue(new_feeList, "selfFee", "0")
	//				new_item = strings.Replace(new_item, feeList, new_feeList, 1)
	//				inifeeLists = append(inifeeLists, map[string]string{
	//					"initialFee": GetNodeValue(new_feeList, "initialFee"),
	//					"value":      new_feeList,
	//				})
	//			}
	//
	//		} else if !RegIsMatch(feeTypeCode, `^(222|224|204)$`) {
	//			inifeeLists = append(inifeeLists, map[string]string{
	//				"initialFee": GetNodeValue(feeList, "initialFee"),
	//				"value":      feeList,
	//			})
	//		}
	//	}
	//	if !isone {
	//		continue
	//	}
	//	inifeeLists = SortSlice(inifeeLists, "initialFee")
	//	for f, inifeeList := range inifeeLists {
	//		feeList := inifeeList["value"]
	//		if momey <= 0 {
	//			break
	//		}
	//		feeTypeCode := GetNodeValue(feeList, "feeTypeCode")
	//		if RegIsMatch(feeTypeCode, `^(222|224|204)$`) {
	//			continue
	//		}
	//		new_feeList := feeList
	//		selfFee := GetNodeValue(feeList, "selfFee")
	//		if ParseFloat(selfFee) == 0 {
	//			initialFee := GetNodeValue(feeList, "initialFee")
	//			if f == len(inifeeLists)-1 {
	//				momey = momey + ParseFloat(selfFee)
	//				new_feeList = SetNodeValue(new_feeList, "selfFee", fmt.Sprintf("%v", momey))
	//			} else {
	//				if momey > ParseFloat(initialFee) {
	//					new_feeList = SetNodeValue(new_feeList, "selfFee", initialFee)
	//					momey = momey - ParseFloat(initialFee)
	//				} else {
	//					new_feeList = SetNodeValue(new_feeList, "selfFee", fmt.Sprintf("%v", momey))
	//					momey = 0
	//				}
	//			}
	//			new_item = strings.Replace(new_item, feeList, new_feeList, 1)
	//		}
	//	}
	//	xmlValue = strings.Replace(xmlValue, item, new_item, 1)
	//
	//}

	// yiyuan := map[string]string{}
	hospitals := RegMatchAll(xmlValue, `<hospitalCode><\/hospitalCode>\s*<hospitalName>.*<\/hospitalName>`)
	for _, hospital := range hospitals {
		// global.GLog.Info("----------------hospital----------------------" + hospital)
		hospitalName := GetNodeValue(hospital, "hospitalName")
		global.GLog.Info("----------------hospitalName----------------------" + hospitalName)
		code, ok := constMap["yiYuanDaiMaMap"][hospitalName]
		if ok {
			global.GLog.Info("----------------code----------------------" + code)
			new_hospital := hospital
			new_hospital = SetNodeValue(new_hospital, "hospitalCode", code)
			xmlValue = strings.Replace(xmlValue, hospital, new_hospital, 1)
		}

	}

	//一、在opList、hospitalList节点下 存在 feeList2节点时：
	//	1.a、只存在一个feeList节点，节点下的feeTypeCode为299或199时， 当存在siList大节点时，
	//	1.a、存在feeList节点，且每个feeList节点下的feeTypeCode均为299或199时 2022年10月08日11:53:16修改
	//		feeList2.billNumber == siList.obillnumber 获取siList
	//		surplus = feeList2.initialFee全部加起来 - siList.inSIAmount
	//		if （feeTypeCode不为222、224、204) feeList2.selfFee =  if surplus > 0  ? feeList2.initialFee : surplus
	//		删除feeList节点
	//	1.b、只存在一个feeList节点，节点下的feeTypeCode为299或199时，不存在siList，不存在thirdList大节点时
	//	1.b、存在feeList节点，且每个feeList节点下的feeTypeCode均为299或199时 2022年10月08日11:53:16修改
	//		surplus = feeList.selfFee
	//		if （feeTypeCode不为222、224、204) feeList2.initialFee =  if surplus > 0  ? feeList2.selfFee : surplus
	//		删除feeList节点
	//	2、不存在feeList节点，存在siList，不存在thirdList大节点时，
	//		feeList2.billNumber == siList.obillnumber 获取siList
	//		surplus = feeList2.initialFee全部加起来 - siList.inSIAmount
	//		if （feeTypeCode不为222、224、204) feeList2.selfFee =  if surplus > 0  ? feeList2.initialFee : surplus
	//	3、不存在feeList节点，不存在siList，存在thirdList大节点时，
	//		feeList2.billNumber == thirdList.obillnumber 获取thirdList
	//		surplus = feeList2.initialFee全部加起来 - thirdList.thirdInSIAmount
	//		if （feeTypeCode不为222、224、204) feeList2.selfFee =  if surplus > 0  ? feeList2.initialFee : surplus
	//二、在opList、hospitalList节点下 不存在 feeList2节点时：
	//	4、只存在一个feeList节点，节点下的feeTypeCode为299或199时，当存在siList大节点时， 单号：562022080035137
	//		feeList.billNumber == siList.obillnumber 获取siList
	//		feeList.selfFee = feeList.initialFee - siList.inSIAmount
	//	5、存在超过一个feeList节点时 feeList的selfFee>initialFee全部加起来 最终取消2022年10月10日10:55:34
	//		surplus = feeList.initialFee全部加起来
	//		更改为：surplus = feeList.selfFee
	//		if （feeTypeCode不为222、224、204)  feeList.selfFee =  if surplus > 0  ? feeList.initialFee : surplus（最后一个不管都是surplus）

	//将所有siList存起来 key:obillnumber  value:inSIAmount
	reg = regexp.MustCompile(`<siList>[\s\S]*?<\/siList>`)
	siListArr := reg.FindAllString(xmlValue, -1)
	numInSIAmount := make(map[string]decimal.Decimal)
	siListMapArr := make(map[string][]string)
	oBillNumbers := []string{}
	for _, siList := range siListArr {
		oBillNumber := GetNodeValue(siList, "obillnumber")
		oBillNumbers = append(oBillNumbers, oBillNumber)
		inSIAmount := GetNodeDecimalValue(siList, "inSIAmount")
		numInSIAmount[oBillNumber] = inSIAmount
		siListMapArr[oBillNumber] = append(siListMapArr[oBillNumber], siList)
	}
	//将所有thirdList存起来 key:obillnumber  value:inSIAmount
	reg = regexp.MustCompile(`<thirdList>[\s\S]*?<\/thirdList>`)
	thirdListArr := reg.FindAllString(xmlValue, -1)
	numThirdInSIAmount := make(map[string]decimal.Decimal)
	thirdListMapArr := make(map[string][]string)
	for _, thirdList := range thirdListArr {
		oBillNumber := GetNodeValue(thirdList, "obillnumber")
		oBillNumbers = append(oBillNumbers, oBillNumber)
		thirdInSIAmount := GetNodeDecimalValue(thirdList, "thirdInSIAmount")
		numThirdInSIAmount[oBillNumber] = thirdInSIAmount
		thirdListMapArr[oBillNumber] = append(thirdListMapArr[oBillNumber], thirdList)
	}

	reg = regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	invoiceList := reg.FindAllString(xmlValue, -1)
	//先排好序
	for _, invoice := range invoiceList {
		newInvoice := invoice
		feeListArr := RegMatchAll(newInvoice, `<feeList>[\s\S]*?<\/feeList>`)
		feeList2Arr := RegMatchAll(newInvoice, `<feeList2>[\s\S]*?<\/feeList2>`)
		feeListArrAndFeeList2Arr := append(feeList2Arr, feeListArr...)
		newFeeList2Arr := make([]string, len(feeList2Arr))
		copy(newFeeList2Arr, feeList2Arr)
		newFeeListArr := make([]string, len(feeListArr))
		copy(newFeeListArr, feeListArr)
		newFeeListArrAndNewFeeList2Arr := append(newFeeList2Arr, newFeeListArr...)
		sort.Slice(newFeeListArrAndNewFeeList2Arr, func(i, j int) bool {
			return GetNodeDecimalValue(newFeeListArrAndNewFeeList2Arr[i], "initialFee").GreaterThan(GetNodeDecimalValue(newFeeListArrAndNewFeeList2Arr[j], "initialFee"))
		})
		for i, feeListArrAndFeeList2 := range feeListArrAndFeeList2Arr {
			if i == len(feeListArrAndFeeList2Arr)-1 {
				newInvoice = strings.Replace(newInvoice, feeListArrAndFeeList2, strings.Join(newFeeListArrAndNewFeeList2Arr, "\n\t\t"), 1)
			} else {
				newInvoice = strings.Replace(newInvoice, feeListArrAndFeeList2, "", 1)
			}
		}
		xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	}

	reg = regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	invoiceList = reg.FindAllString(xmlValue, -1)
	for _, invoice := range invoiceList {
		newInvoice := invoice
		billNumber := GetNodeValue(newInvoice, "billNumber")
		siListArrInInvoice, isInSIAmount := numInSIAmount[billNumber]
		thirdListArrInInvoice, isThirdInSIAmount := numThirdInSIAmount[billNumber]
		feeListArr := RegMatchAll(newInvoice, `<feeList>[\s\S]*?<\/feeList>`)
		selfFee0 := RegMatchAll(newInvoice, `<selfFee>0<\/selfFee>`)
		if len(feeListArr) != len(selfFee0) {
			continue
		}
		newFeeListArr := make([]string, len(feeListArr))
		copy(newFeeListArr, feeListArr)

		if isInSIAmount {
			totalFeeList2InitialFee := xmlNodeTotalVal(newFeeListArr, "initialFee")
			fmt.Println("---------------numInSIAmount-----------------------", billNumber, totalFeeList2InitialFee, siListArrInInvoice)
			// ---------------numInSIAmount----------------------- 1987654567 0 7600
			surplus := totalFeeList2InitialFee.Sub(siListArrInInvoice)
			newInvoice = setSelfFeeVal(newFeeListArr, surplus, newInvoice, feeListArr)
			// for _, feeListTemp := range feeListArr {
			// 	newInvoice = strings.Replace(newInvoice, feeListTemp, "", 1)
			// }
			xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
		} else if isThirdInSIAmount {
			totalFeeList2InitialFee := xmlNodeTotalVal(newFeeListArr, "initialFee")
			surplus := totalFeeList2InitialFee.Sub(thirdListArrInInvoice)
			newInvoice = setSelfFeeVal(newFeeListArr, surplus, newInvoice, feeListArr)
			// for _, feeListTemp := range feeListArr {
			// 	newInvoice = strings.Replace(newInvoice, feeListTemp, "", 1)
			// }
			xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
		}
	}

	// reg = regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	// invoiceList = reg.FindAllString(xmlValue, -1)
	// for _, invoice := range invoiceList {
	// 	newInvoice := invoice
	// 	feeListArr := RegMatchAll(newInvoice, `<feeList>[\s\S]*?<\/feeList>`)
	// 	feeList2Arr := RegMatchAll(newInvoice, `<feeList2>[\s\S]*?<\/feeList2>`)
	// 	billNumber := GetNodeValue(newInvoice, "billNumber")
	// 	siListArrInInvoice := siListMapArr[billNumber]
	// 	thirdListArrInInvoice := thirdListMapArr[billNumber]
	// 	//第二 在opList、hospitalList节点下 不存在 feeList2节点时：
	// 	if len(feeList2Arr) == 0 {
	// 		//	global.GLog.Info("第二，billNumber" + billNumber)
	// 		//	//4、只存在一个feeList节点，节点下的feeTypeCode为299或199时，当存在siList大节点时
	// 		//	if len(feeListArr) == 1 && RegIsMatch(feeListArr[0], `<feeTypeCode>(199|299)</feeTypeCode>`) && len(numInSIAmount) > 0 {
	// 		//		global.GLog.Info("第4点，billNumber" + billNumber)
	// 		//		newFeeList := feeListArr[0]
	// 		//		initialFee := GetNodeDecimalValue(newFeeList, "initialFee")
	// 		//		subDecimal := initialFee.Sub(numInSIAmount[billNumber])
	// 		//		SetNodeValue(newFeeList, "selfFee", subDecimal.String())
	// 		//		newInvoice = strings.Replace(newInvoice, feeListArr[0], newFeeList, 1)
	// 		//	}
	// 		//	xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)

	// 		//5、存在超过一个feeList节点时 取消2022年10月10日10:55:28
	// 		//if len(feeListArr) > 1 {
	// 		//	global.GLog.Info("第5点，billNumber" + billNumber)
	// 		//	//计算总值
	// 		//	totalFeeListInitialFee := xmlNodeTotalVal(feeListArr, "initialFee")
	// 		//	surplus := totalFeeListInitialFee
	// 		//	fmt.Println("5. totalFeeListInitialFee", totalFeeListInitialFee)
	// 		//	fmt.Println("5. surplus", surplus)
	// 		//	//替换值
	// 		//	for i, feeList := range feeListArr {
	// 		//		newfeeList := feeList
	// 		//		if !RegIsMatch(newfeeList, `<feeTypeCode>(222|224|204)</feeTypeCode>`) {
	// 		//			selfFee := surplus
	// 		//			initialFee := GetNodeDecimalValue(newfeeList, "initialFee")
	// 		//			if i != len(feeListArr)-1 {
	// 		//				if totalFeeListInitialFee.GreaterThan(initialFee) {
	// 		//					selfFee = initialFee
	// 		//					surplus = surplus.Sub(initialFee)
	// 		//				}
	// 		//			}
	// 		//			newfeeList = SetNodeValue(newfeeList, "selfFee", selfFee.String())
	// 		//			newInvoice = strings.Replace(newInvoice, feeList, newfeeList, 1)
	// 		//		}
	// 		//	}
	// 		//}
	// 		//xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	// 	}
	// 	//第一 在opList、hospitalList节点下 存在 feeList2节点时：
	// 	if len(feeList2Arr) != 0 {
	// 		global.GLog.Info("第一，billNumber" + billNumber)
	// 		//newFeeList2Arr := feeList2Arr
	// 		newFeeList2Arr := make([]string, len(feeList2Arr))
	// 		copy(newFeeList2Arr, feeList2Arr)
	// 		//newFeeListArr := feeListArr
	// 		newFeeListArr := make([]string, len(feeListArr))
	// 		copy(newFeeListArr, feeListArr)

	// 		//1、a.只存在一个feeList节点，节点下的feeTypeCode为299或199时，当存在siList大节点时，
	// 		//1、a.存在feeList节点，且每个feeList节点下的feeTypeCode均为299或199时 2022年10月08日11:53:16修改
	// 		flag199 := true
	// 		selfFeeDesc := ""
	// 		for _, feeList := range feeListArr {
	// 			if !RegIsMatch(feeList, `<feeTypeCode>(199|299)</feeTypeCode>`) {
	// 				flag199 = false
	// 			}
	// 			if GetNodeValue(feeList, "selfFeeDesc") != "" {
	// 				selfFeeDesc = GetNodeValue(feeList, "selfFeeDesc")
	// 			}
	// 		}

	// 		//获取最大值的下标
	// 		initialFeeDecimal := decimal.Zero
	// 		maxInitialFeeIndex := 0
	// 		for i, feeList2 := range newFeeList2Arr {
	// 			if !RegIsMatch(feeList2, `<feeTypeCode>(222|224|204)</feeTypeCode>`) {
	// 				if GetNodeDecimalValue(feeList2, "initialFee").GreaterThan(initialFeeDecimal) {
	// 					initialFeeDecimal = GetNodeDecimalValue(feeList2, "initialFee")
	// 					maxInitialFeeIndex = i
	// 				}
	// 			}
	// 		}
	// 		if len(feeListArr) >= 1 && flag199 && len(siListArrInInvoice) > 0 {
	// 			global.GLog.Info("第1.a点，billNumber" + billNumber)

	// 			//赋值在initialFee最大值的节点下的selfFeeDesc
	// 			newFeeList2Arr[maxInitialFeeIndex] = SetNodeValue(newFeeList2Arr[maxInitialFeeIndex], "selfFeeDesc", selfFeeDesc)
	// 			//计算总值
	// 			totalFeeList2InitialFee := xmlNodeTotalVal(newFeeList2Arr, "initialFee")
	// 			fmt.Println("1.a totalFeeList2InitialFee", totalFeeList2InitialFee)
	// 			//计算总值和numInSIAmount的差
	// 			surplus := totalFeeList2InitialFee.Sub(numInSIAmount[billNumber])
	// 			fmt.Println("1.a surplus", surplus)
	// 			newInvoice = setSelfFeeVal(newFeeList2Arr, surplus, newInvoice, feeList2Arr)
	// 			for _, feeListTemp := range feeListArr {
	// 				newInvoice = strings.Replace(newInvoice, feeListTemp, "", 1)
	// 			}
	// 			xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	// 		}

	// 		//1、b.只存在一个feeList节点，节点下的feeTypeCode为299或199时，不存在siList，不存在thirdList大节点时
	// 		if len(feeListArr) >= 1 && flag199 && len(siListArrInInvoice) == 0 && len(thirdListArrInInvoice) == 0 {
	// 			global.GLog.Info("第1.b点，billNumber" + billNumber)
	// 			//赋值在initialFee最大值的节点下的selfFeeDesc
	// 			newFeeList2Arr[maxInitialFeeIndex] = SetNodeValue(newFeeList2Arr[maxInitialFeeIndex], "selfFeeDesc", selfFeeDesc)
	// 			surplus := decimal.Zero
	// 			for _, feeList := range feeListArr {
	// 				if GetNodeDecimalValue(feeList, "selfFee").IsZero() {
	// 					continue
	// 				}
	// 				surplus = GetNodeDecimalValue(feeList, "selfFee")
	// 			}
	// 			fmt.Println("1.b surplus", surplus)
	// 			if !surplus.IsZero() {
	// 				newInvoice = setSelfFeeVal(newFeeList2Arr, surplus, newInvoice, feeList2Arr)
	// 				for _, feeListTemp := range feeListArr {
	// 					newInvoice = strings.Replace(newInvoice, feeListTemp, "", 1)
	// 				}
	// 				xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	// 			}
	// 		}

	// 		//2、不存在feeList节点，存在siList，不存在thirdList大节点时，
	// 		if len(feeListArr) == 0 && len(siListArrInInvoice) > 0 && len(thirdListArrInInvoice) == 0 {
	// 			global.GLog.Info("第2点，billNumber" + billNumber)
	// 			//计算总值
	// 			totalFeeList2InitialFee := xmlNodeTotalVal(newFeeList2Arr, "initialFee")
	// 			fmt.Println("2. totalFeeList2InitialFee", totalFeeList2InitialFee)
	// 			//计算总值和inSIAmount的差
	// 			surplus := totalFeeList2InitialFee.Sub(numInSIAmount[billNumber])
	// 			fmt.Println("2. surplus", surplus)
	// 			newInvoice = setSelfFeeVal(newFeeList2Arr, surplus, newInvoice, feeList2Arr)
	// 			xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	// 		}

	// 		//3、不存在feeList节点，不存在siList，存在thirdList大节点时，
	// 		if len(feeListArr) == 0 && len(siListArrInInvoice) == 0 && len(thirdListArrInInvoice) > 0 {
	// 			global.GLog.Info("第3点，billNumber" + billNumber)
	// 			//计算总值
	// 			totalFeeList2InitialFee := xmlNodeTotalVal(newFeeList2Arr, "initialFee")
	// 			//计算总值和thirdInSIAmount的差
	// 			surplus := totalFeeList2InitialFee.Sub(numThirdInSIAmount[billNumber])
	// 			newInvoice = setSelfFeeVal(newFeeList2Arr, surplus, newInvoice, feeList2Arr)
	// 			xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	// 		}
	// 	}
	// }

	// chaches := map[string]map[string]float64{}
	// items = RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	// for t, item := range items {
	// 	new_item := item
	// 	// billNumber := GetNodeValue(new_item, "billNumber")
	// 	insuredName := GetNodeValue(new_item, "insuredName")
	// 	_, isOK := chaches[insuredName]
	// 	if !isOK {
	// 		chaches[insuredName] = map[string]float64{"num": 0.0, "initialFee": 0.0}
	// 	}
	// 	chaches[insuredName]["num"] = chaches[insuredName]["num"] + 1
	// 	// initialFee := 0
	// 	feeLists := RegMatchAll(xmlValue, `<feeList>[\s\S]*?<\/feeList>`)
	// 	for f, feeList := range feeLists {
	// 		if GetNodeValue(feeList, "initialFee") != "" {
	// 			chaches[insuredName]["initialFee"] = chaches[insuredName]["initialFee"] + ParseFloat(GetNodeValue(feeList, "initialFee"))
	// 		}
	// 		if t == len(items)-1 && f == len(feeLists)-1 {
	// 			if len(chaches) > 1 {
	// 				max_key := map[string]float64{"num": 0.0, "initialFee": 0.0}
	// 				total := 0.0
	// 				for _, chache := range chaches {
	// 					if chache["num"] > max_key["num"] {
	// 						max_key = chache
	// 					}
	// 					total = total + chache["initialFee"]
	// 				}
	// 				total = total - max_key["initialFee"]
	// 				global.GLog.Info("----------------total-----------------", zap.Any("total:", total))
	// 				new_feeList := feeList
	// 				new_feeList = SetNodeValue(new_feeList, "selfPay", fmt.Sprintf("%v", total))
	// 				new_feeList = SetNodeValue(new_feeList, "selfPayDesc", "出险人错误")
	// 				new_item = strings.Replace(new_item, feeList, new_feeList, 1)
	// 				xmlValue = strings.Replace(xmlValue, item, new_item, 1)

	// 			}
	// 		}
	// 	}
	// }

	//items = RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	//for _, item := range items {
	//	new_item := item
	//	if strings.Index(new_item, "<feeTypeCode>199</feeTypeCode>") != -1 || strings.Index(new_item, "<feeTypeCode>299</feeTypeCode>") != -1 {
	//		cache99 := "null"
	//		cacheDatas := []map[string]string{}
	//		feeLists := RegMatchAll(new_item, `<feeList>[\s\S]*?<\/feeList>`)
	//		for _, feeList := range feeLists {
	//			new_feeList := feeList
	//			feeTypeCode := GetNodeValue(new_feeList, "feeTypeCode")
	//			if RegIsMatch(feeTypeCode, `^(199|299)$`) {
	//				if cache99 == "null" {
	//					cache99 = new_feeList
	//				}
	//			} else {
	//				selfFeeDesc := GetNodeValue(new_feeList, "selfFeeDesc")
	//				selfFee := ParseFloat(GetNodeValue(new_feeList, "selfFee"))
	//				if RegIsMatch(selfFeeDesc, `(膳食费|陪床费)`) {
	//					selfFeeDescs := strings.Split(selfFeeDesc, "；")
	//					for _, selfPayDesc := range selfFeeDescs {
	//						selfPayDesc = selfPayDesc + "；"
	//						if RegIsMatch(selfPayDesc, `(膳食费|陪床费)`) {
	//							mms := RegMatchAll(selfPayDesc, `.*自付比例(.+)% 自付金额(.+)元；`)
	//							bl := ParseFloat(mms[1])
	//							selfPay := ParseFloat(mms[2])
	//							selfFee = selfFee - selfPay
	//							selfFeeDesc = strings.Replace(selfFeeDesc, selfPayDesc, "", 1)
	//							cacheDatas = append(cacheDatas, map[string]string{"bl": ToString(bl), "selfPayDesc": selfPayDesc, "selfPay": ToString(selfPay)})
	//						}
	//					}
	//					new_feeList = SetNodeValue(new_feeList, "selfFeeDesc", selfFeeDesc)
	//					new_feeList = SetNodeValue(new_feeList, "selfFee", ToString(selfFee))
	//					new_item = strings.Replace(new_item, feeList, new_feeList, 1)
	//				}
	//			}
	//		}
	//		if cache99 != "null" && len(cacheDatas) > 0 {
	//			new_cache99 := cache99
	//			cacheDatas = SortSlice(cacheDatas, "bl")
	//			selfPay := 0.0
	//			selfPayDesc := ""
	//			for _, cacheData := range cacheDatas {
	//				selfPay = selfPay + ParseFloat(cacheData["selfPay"])
	//				selfPayDesc = selfPayDesc + cacheData["selfPayDesc"]
	//			}
	//			new_cache99 = SetNodeValue(new_cache99, "selfPay", ToString(selfPay))
	//			new_cache99 = SetNodeValue(new_cache99, "selfPayDesc", selfPayDesc)
	//			new_item = strings.Replace(new_item, cache99, new_cache99, 1)
	//			xmlValue = strings.Replace(xmlValue, item, new_item, 1)
	//
	//		}
	//
	//	}
	//}

	//CSB0118RC0293000 xml
	//1、在opList、hospitalList节点下有多个feeList节点中，
	//如部分feeList节点中feeTypeCode、initialFee、selfFee、均有值，
	//部分feeList节点中只有feeTypeCode、initialFee有值时，将部分feeList节点中只有feeTypeCode、initialFee有值对应的feeList节点清空
	//意思是:有多个feeList节点的时候，只有feeTypeCode、initialFee、selfFee三个节点都有值的情况下才保留，如果其中一个节点为空，那对应的那组feeList节点屏蔽
	//2、存在feeList2节点时，且只存在一个feeList节点，节点下的feeTypeCode为299或199时，清空所有的feeList2节点
	reg = regexp.MustCompile(`<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	invoiceList = reg.FindAllString(xmlValue, -1)
	for _, invoice := range invoiceList {
		newInvoice := invoice
		feeListArr := RegMatchAll(newInvoice, `<feeList>[\s\S]*?<\/feeList>`)
		feeList2Arr := RegMatchAll(newInvoice, `<feeList2>[\s\S]*?<\/feeList2>`)
		//第一点
		if len(feeListArr) > 1 {
			for _, feeList := range feeListArr {
				// initialFee := GetNodeValue(feeList, "initialFee")
				feeTypeCode := GetNodeValue(feeList, "feeTypeCode")
				initialFee := GetNodeValue(feeList, "initialFee")
				// selfFee := GetNodeValue(feeList, "selfFee")
				// if initialFee == "" || feeTypeCode == "" || selfFee == "" {
				if feeTypeCode == "" || initialFee == "" {
					newInvoice = strings.Replace(newInvoice, feeList, "", 1)
				}
			}
		}
		//第一点
		if len(feeList2Arr) > 1 {
			for _, feeList2 := range feeList2Arr {
				initialFee := GetNodeValue(feeList2, "initialFee")
				feeTypeCode := GetNodeValue(feeList2, "feeTypeCode")
				selfFee := GetNodeValue(feeList2, "selfFee")
				if initialFee == "" || feeTypeCode == "" || selfFee == "" {
					newInvoice = strings.Replace(newInvoice, feeList2, "", 1)
				}
			}
		}
		//第二点
		//if len(feeList2Arr) > 0 && len(feeListArr) == 1 && RegIsMatch(feeListArr[0], `<feeTypeCode>(199|299)</feeTypeCode>`) {
		//	for _, feeList2 := range feeList2Arr {
		//		newInvoice = strings.Replace(newInvoice, feeList2, "", 1)
		//	}
		//}
		//第二点改成feeList
		// if len(feeList2Arr) > 0 && len(feeListArr) == 1 && RegIsMatch(feeListArr[0], `<feeTypeCode>(199|299)</feeTypeCode>`) {
		// 	for _, feeList := range feeListArr {
		// 		newInvoice = strings.Replace(newInvoice, feeList, "", 1)
		// 	}
		// }
		xmlValue = strings.Replace(xmlValue, invoice, newInvoice, 1)
	}

	invoiceList = RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	for _, invoice := range invoiceList {
		newInvoice := invoice
		feeList2Arr := RegMatchAll(invoice, `<feeList>[\s\S]*?<\/feeList>`)
		hospitalType := GetNodeValue(invoice, "hospitalType")
		// oBillNumbers
		billNumber := GetNodeValue(invoice, "billNumber")
		if arrays.ContainsString(oBillNumbers, billNumber) != -1 {
			for _, feeList2 := range feeList2Arr {
				newfeeList2 := feeList2
				selfFee := GetNodeValue(newfeeList2, "selfFee")
				selfFeeDesc := GetNodeValue(newfeeList2, "selfFeeDesc")
				if selfFeeDesc == "" && selfFee != "0" && selfFee != "0.00" && selfFee != "0.0" {
					feeList2 = SetNodeValue(feeList2, "selfFeeDesc", "按照"+hospitalType+"市结算单扣除自付金额"+selfFee+"元")
				}
				invoice = strings.Replace(invoice, newfeeList2, feeList2, 1)
			}
		}

		invoice = SetNodeValue(invoice, "hospitalType", "")
		xmlValue = strings.Replace(xmlValue, newInvoice, invoice, 1)
	}

	xmlValue = strings.Replace(xmlValue, "<feeList2>", "<feeList>", -1)
	xmlValue = strings.Replace(xmlValue, "</feeList2>", "</feeList>", -1)

	xmlValue = RegReplace(xmlValue, `<obillnumber>AAAA\d+<\/obillnumber>`, "<obillnumber></obillnumber>")
	xmlValue = RegReplace(xmlValue, `<billNumber>AAAA\d+<\/billNumber>`, "<billNumber></billNumber>")

	xmlValue = strings.Replace(xmlValue, "<thirdList2>", "<thirdList>", -1)
	xmlValue = strings.Replace(xmlValue, "</thirdList2>", "</thirdList>", -1)

	items = RegMatchAll(xmlValue, `<(thirdList)>[\s\S]*?<\/(thirdList)>`)
	for _, item := range items {
		orgTypeCodeValue := GetNodeValue(item, "orgTypeCode")
		if orgTypeCodeValue == "" {
			xmlValue = strings.ReplaceAll(xmlValue, item, "")
		}
	}

	//chaches := map[string]map[string]float64{}
	//items = RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	//for t, item := range items {
	//	new_item := item
	//	// billNumber := GetNodeValue(new_item, "billNumber")
	//	insuredName := GetNodeValue(new_item, "insuredName")
	//	_, isOK := chaches[insuredName]
	//	if !isOK {
	//		chaches[insuredName] = map[string]float64{"num": 0.0, "initialFee": 0.0}
	//	}
	//	chaches[insuredName]["num"] = chaches[insuredName]["num"] + 1
	//	// initialFee := 0
	//	feeLists := RegMatchAll(xmlValue, `<feeList>[\s\S]*?<\/feeList>`)
	//	for f, feeList := range feeLists {
	//		if GetNodeValue(feeList, "initialFee") != "" {
	//			chaches[insuredName]["initialFee"] = chaches[insuredName]["initialFee"] + ParseFloat(GetNodeValue(feeList, "initialFee"))
	//		}
	//		if t == len(items)-1 && f == len(feeLists)-1 {
	//			if len(chaches) > 1 {
	//				max_key := map[string]float64{"num": 0.0, "initialFee": 0.0}
	//				total := 0.0
	//				for _, chache := range chaches {
	//					if chache["num"] > max_key["num"] {
	//						max_key = chache
	//					}
	//					total = total + chache["initialFee"]
	//				}
	//				total = total - max_key["initialFee"]
	//				global.GLog.Info("----------------total-----------------", zap.Any("total:", total))
	//				new_feeList := feeList
	//				new_feeList = SetNodeValue(new_feeList, "selfPay", fmt.Sprintf("%v", total))
	//				new_feeList = SetNodeValue(new_feeList, "selfPayDesc", "出险人错误")
	//				new_item = strings.Replace(new_item, feeList, new_feeList, 1)
	//				xmlValue = strings.Replace(xmlValue, item, new_item, 1)
	//
	//			}
	//		}
	//	}
	//}

	//CSB0118RC0306000
	//orgTypeCode节点值为10时转化为1
	//orgTypeCode节点值为20时转化为2
	//校验放在CSB0118RC0298后面

	items = RegMatchAll(xmlValue, `<(thirdList|siList)>[\s\S]*?<\/(thirdList|siList)>`)
	for _, item := range items {
		orgTypeCodeValue := GetNodeValue(item, "orgTypeCode")
		if orgTypeCodeValue == "10" {
			newItem := SetNodeValue(item, "orgTypeCode", "1")
			xmlValue = strings.ReplaceAll(xmlValue, item, newItem)
		}
		if orgTypeCodeValue == "20" {
			newItem := SetNodeValue(item, "orgTypeCode", "2")
			xmlValue = strings.ReplaceAll(xmlValue, item, newItem)
		}
	}

	//CSB0118RC0291002
	//当全部hospitalList、opList大节点下的insuredName的值有存在不一致时（采取少数服从多数的原则），
	//将不一致的insuredName对应的hospitalList或opList下所有的feeList节点下的
	//selfFee默认为0、
	//selfFeeDesc默认为空、
	//selfPay的值取对应节点feeList下initialFee的值、
	//selfPayDesc节点处默认“出险人错误”，需求放在最后
	hospitalMap := make(map[string]int, 0)
	items = utils.RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	flagMaxNumVal := ""
	flagMaxNum := 0
	for _, item := range items {
		insuredName := GetNodeValue(item, "insuredName")
		if _, ok := hospitalMap[insuredName]; ok {
			hospitalMap[insuredName] = hospitalMap[insuredName] + 1
		} else {
			hospitalMap[insuredName] = 1
		}
	}
	for k, v := range hospitalMap {
		if flagMaxNum < v {
			flagMaxNum = v
			flagMaxNumVal = k
		}
	}
	for _, item := range items {
		insuredName := utils.GetNodeValue(item, "insuredName")
		if insuredName == flagMaxNumVal || len(hospitalMap) < 2 {
			continue
		}
		newItem := item
		feeLists := utils.RegMatchAll(newItem, `<feeList>[\s\S]*?<\/feeList>`)
		for _, feeList := range feeLists {
			newFeeList := feeList
			initialFee := utils.GetNodeValue(newFeeList, "initialFee")
			newFeeList = utils.SetNodeValue(newFeeList, "selfPay", initialFee)
			newFeeList = utils.SetNodeValue(newFeeList, "selfFee", "0")
			newFeeList = utils.SetNodeValue(newFeeList, "selfFeeDesc", "")
			newFeeList = utils.SetNodeValue(newFeeList, "selfPayDesc", "出险人错误")
			newItem = strings.ReplaceAll(newItem, feeList, newFeeList)
		}
		xmlValue = strings.ReplaceAll(xmlValue, item, newItem)
	}

	fc277Map := make(map[string]string)

	fc002Map := make(map[string]string, 0)
	invoices := obj.HospitalList
	invoices = append(invoices, obj.OpList...)
	for _, invoice := range invoices {
		fc002 := getOneValue(invoice.Fields, "fc002")
		fc275 := getOneValue(invoice.Fields, "fc275")
		fmt.Println("----------fc002-----fc275-----------------", fc002, fc275)
		fc002Map[fc002] = fc275

		fc277 := getOneValue(invoice.Fields, "fc277")
		fc277Map[fc002] = fc277
	}

	items = utils.RegMatchAll(xmlValue, `<(opList|hospitalList)>[\s\S]*?<\/(opList|hospitalList)>`)
	for _, item := range items {
		billNumber := GetNodeValue(item, "billNumber")
		fc275 := fc002Map[billNumber]
		newitem := item
		if fc275 == "1" {
			newitem = utils.SetNodeValue(newitem, "selfPay", "")
			newitem = utils.SetNodeValue(newitem, "selfFee", "0")
			newitem = utils.SetNodeValue(newitem, "selfFeeDesc", "")
			newitem = utils.SetNodeValue(newitem, "selfPayDesc", "")
			xmlValue = strings.ReplaceAll(xmlValue, item, newitem)
		}

		//CSB0118RC0344000
		//fc277结果值为N时，fc279（billCode）、fc281（billCheckCode）、
		//fc282（billVerification） 节点值默认为空
		fc277 := fc277Map[billNumber]
		if fc277 == "N" {
			newitem = utils.SetNodeValue(newitem, "billCode", "")
			newitem = utils.SetNodeValue(newitem, "billCheckCode", "")
			newitem = utils.SetNodeValue(newitem, "billVerification", "")
			xmlValue = strings.ReplaceAll(xmlValue, item, newitem)
		}

	}

	//CSB0118RC0346000
	//当errorDesc节点值包含“参数错误：”时，屏蔽“参数错误：”这四个字，
	//如：“N0017: 电子发票（NO.xxxx）未成功查验，原因：参数错误：电子票据代码不正确。”，
	//修改为：“N0017: 电子发票（NO.xxxx）未成功查验，原因：电子票据代码不正确。”
	CSB0118RC0346000XmlTmp := ""
	errorList := utils.RegMatchAll(xmlValue, `<(errorList)>[\s\S]*?<\/(errorList)>`)
	for _, item := range errorList {
		newitem := item
		errorDesc := GetNodeValue(item, "errorDesc")
		CSB0118RC0346000XmlTmp = strings.ReplaceAll(errorDesc, "参数错误：", "")
		newitem = utils.SetNodeValue(newitem, "errorDesc", CSB0118RC0346000XmlTmp)
		xmlValue = strings.ReplaceAll(xmlValue, item, newitem)
	}

	newXmlValue = xmlValue
	return err, newXmlValue
}

func ToString(data float64) string {
	// return strconv.FormatFloat(data, 'E', -1, 64)
	return fmt.Sprintf("%v", data)
}

func SortSlice(items []map[string]string, key string) []map[string]string {
	sort.Slice(items, func(i, j int) bool {
		value, ok := items[i][key]
		if !ok {
			return false
		}
		valI, _ := strconv.ParseFloat(value, 64)
		value, ok = items[j][key]
		if !ok {
			return false
		}
		valJ, _ := strconv.ParseFloat(value, 64)
		return valI < valJ
	})
	return items
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

// ParseDecimal string转decimal
func ParseDecimal(v string) decimal.Decimal {
	d, err := decimal.NewFromString(v)
	if err != nil {
		return decimal.Zero
	}
	return d
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

// xmlNodeTotalVal 计算xml里面节点值总和
func xmlNodeTotalVal(xmlArr []string, node string) decimal.Decimal {
	total := decimal.Zero
	for _, xml := range xmlArr {
		total = total.Add(GetNodeDecimalValue(xml, node))
	}
	return total
}

// setSelfFeeVal 设置selfFeeVal的值
func setSelfFeeVal(newFeeArr []string, surplus decimal.Decimal, newInvoice string, feeArr []string) string {
	//最后一个要全部放在那里不在处理
	flagIndex := len(newFeeArr) - 1
	endSurplus := decimal.Zero
	//减到没有了就不减了
	lessZeroFlag := 0

	for i := 0; i < len(newFeeArr); i++ {
		selfFee := surplus
		fmt.Println("selfFee000", selfFee)
		if !RegIsMatch(newFeeArr[i], `<feeTypeCode>(222|224)</feeTypeCode>`) && lessZeroFlag <= 0 {
			initialFee := GetNodeDecimalValue(newFeeArr[i], "initialFee")
			if initialFee.IsZero() {
				continue
			}
			flagIndex = i
			fmt.Println("initialFee", initialFee)
			endSurplus = surplus
			surplus = surplus.Sub(initialFee)
			if surplus.GreaterThan(decimal.Zero) {
				selfFee = initialFee
			} else {
				selfFee = endSurplus
				lessZeroFlag++
			}
			fmt.Println("surplus", surplus)
			fmt.Println("selfFee1111", selfFee)
			newFeeArr[i] = SetNodeValue(newFeeArr[i], "selfFee", selfFee.String())
			newInvoice = strings.Replace(newInvoice, feeArr[i], newFeeArr[i], 1)
		}
	}
	fmt.Println("flagIndex", flagIndex)
	fmt.Println("endSurplus", endSurplus)

	//最后一个要全部放在那里不在处理
	newFeeArrTemp := make([]string, len(newFeeArr))
	copy(newFeeArrTemp, newFeeArr)
	newFeeArrTemp[flagIndex] = SetNodeValue(newFeeArr[flagIndex], "selfFee", endSurplus.String())
	newInvoice = strings.Replace(newInvoice, newFeeArr[flagIndex], newFeeArrTemp[flagIndex], 1)
	return newInvoice
}
