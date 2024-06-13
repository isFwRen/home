package B0122

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/load/model"
	"server/utils"
	"strconv"
	"strings"

	"time"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0122:::XmlDeal-----------------------")
	// return nil, xmlValue
	obj := o.(FormatObj)
	//fields := obj.Fields
	// errorData := ""
	errMeses := ""
	// questionCount := 0
	// medical := obj.Medical

	// for _, field := range obj.Fields {
	// 	// if len(filed.Issues ) > 0 {
	// 	for _, issue := range field.Issues {
	// 		errMes := issue.Code + "_" + issue.Message + ";"
	// 		if strings.Index(errMeses, errMes) == -1 {
	// 			if questionCount == 0 {
	// 				errorData = "	  <question_count></question_count>\n"
	// 			}
	// 			questionCount = questionCount + 1
	// 			errorData = errorData + "	  <question>\n	   <question_code>" + issue.Code + "</question_code>\n	   <question_content>" + issue.Message + "</question_content>\n	  </question>\n"
	// 			errMeses += errMes
	// 		}
	// 	}
	// }
	// if errorData == "" {
	// 	errorData = "	  <question_count>0</question_count>\n"
	// }
	// // fmt.Println("------------------------error-------------------", errorData)
	// errorData = "<question_info>\n	<questions>\n" + errorData + "</questions>\n	</question_info>"
	// xmlValue = strings.Replace(xmlValue, "<question_info></question_info>", errorData, 1)

	//if questionCount != 0 {
	//	xmlValue = utils.SetNodeValue(xmlValue, "question_count", strconv.Itoa(questionCount))
	//}

	//CSB0122RC0127000
	//question_info节点下的question_count节点的值，根据有多少组question节点进行数量赋值
	//出现两个以上的问题件时，生成多一个question节点，有多少个问题件出现多少个question节点
	//question_code节点内容为问题件编码
	//question_content节点内容为问题件描述
	//questionInfos := RegMatchAll(xmlValue, `<question_info>[\s\S]*?<\/question_info>`)
	questions := RegMatchAll(xmlValue, `<questions>[\s\S]*?<\/questions>`)
	errorMessage := "\n"
	errorMessageCount := 0
	for _, field := range obj.Fields {
		for _, issue := range field.Issues {
			errMes := issue.Code + "_" + issue.Message + ";"
			if strings.Index(errMeses, errMes) == -1 {
				// if len(field.Issues) > 0 {
				errorMessageCount++
				fmt.Println("----------------------errorMessageCount----------", field.Code, errorMessageCount, errMes)
				errorMessage = errorMessage + "	  <question>\n	   <question_code>" + issue.Code + "</question_code>\n	   <question_content>" + issue.Message + "</question_content>\n	  </question>\n"
				// }
				errMeses += errMes
			}

			// fmt.Println("field.Issues=len(field.Issues)", field.Issues, len(field.Issues))
		}
	}
	for _, question := range questions {
		if errorMessageCount > 0 {
			xmlValue = strings.Replace(xmlValue, question, "<questions>"+errorMessage+"\t</questions>", 1)
		}
	}

	hospital_names := RegMatchAll(xmlValue, `<hospital_name>.*<\/hospital_name>`)
	for _, hospital_name := range hospital_names {
		hospital_name_clone := hospital_name
		val := GetNodeValue(hospital_name_clone, "hospital_name")
		val = utils.RegReplace(val, `_[^_]*_[^_]*$`, "")
		hospital_name = utils.SetNodeValue(hospital_name, "hospital_name", val)
		xmlValue = strings.Replace(xmlValue, hospital_name_clone, hospital_name, 1)
	}

	// if errorMessageCount > 0 {
	// 	//xmlValue = strings.Replace(xmlValue, "<question_count></question_count>", "<question_count>"+strconv.Itoa(errorMessageCount)+"</question_count>", 1)
	// 	xmlValue = utils.SetNodeValue(xmlValue, "question_count", strconv.Itoa(errorMessageCount))

	// }

	// fmt.Println("errorMessage=", errorMessage)
	// 编码 CSB0122RC0120000 以下节点内容的右边节点的值取下载包中的数据
	//timeliness取下载报文中的timeliness的值
	//task_no取下载报文中的task_no的值
	//claim_id取下载报文中的claim_id的值
	//branch_no取下载报文中的branch_no的值
	//branch_name取下载报文中的branch_name的值
	//send_time取下载报文中的end_time的值
	//case_type取下载报文中的case_type的值
	//accident_name取下载报文中accident_name的值
	//claim_type取下载报文中claim_type的值
	items := []string{"timeliness", "task_no", "claim_id", "branch_no", "branch_name", "send_time", "case_type", "accident_name", "claim_type"}
	for _, item := range items {
		value := utils.GetNodeValue(obj.Bill.OtherInfo, item)
		xmlValue = utils.SetNodeValue(xmlValue, item, value)
	}
	// claim_form取下载报文中claim_type的值
	value := utils.GetNodeValue(obj.Bill.OtherInfo, "claim_type")
	xmlValue = utils.SetNodeValue(xmlValue, "claim_form", value)

	now := time.Now()
	xmlValue = utils.SetNodeValue(xmlValue, "finish_time", now.Format("20060102150405"))
	// accidentAreaOption := utils.GetNodeValue(data, "accidentAreaOption")
	// items = []string{"beneficiary_record", "operation_recordd", "injury_record", "medical_record", "medical_record"}
	// for _, item := range items {
	// 	records := utils.RegMatchAll(xmlValue, `<`+item+`>[\s\S]*?<\/`+item+`>`)
	// 	for _, record := range records {
	// 		if !RegIsMatch(record, `>.+<`) {
	// 			// xmlValue = utils.SetNodeValue(xmlValue, record, "")
	// 			xmlValue = strings.Replace(xmlValue, record, "", 1)
	// 		}
	// 	}

	// }

	//CSB0122RC0136000 CSB0122RC0137000
	// 将每个fc040与fc037字段的结果值作对比，一致则将同一诊断的MB002-bc022、bc023、bc024、bc025、bc026的
	// 字段放置在与fc037同一个medical_record大节点下的已配置对应字段的相关节点处，MB002-bc027分块下的字段放置在与fc037同一个medical_record大节点下的injury_info、injury_info_bj节点
	// 将每个fc041与fc037字段的结果值作对比，一致则将同一手术的MB002-bc034分块下的字段放置在与fc037同一个medical_record大节点下的operation_info、operation_info_bj节点，注：MB002-bc034分块下的字段可能会存在多个不为空且不重复的结果值，将不重复的结果值放置在对应节点内；
	// infoItems := [][]string{{"injury_records", "injury_record", "injury_code"}, {"operation_records", "operation_record", "operation_code"}}
	// for _, infoItem := range infoItems {
	// 	itemList := RegMatchAll(xmlValue, `<`+infoItem[0]+`>[\s\S]*?<\/`+infoItem[0]+`>`)
	// 	codes := []string{}
	// 	for _, medicalRecord := range itemList {
	// 		injuryCodeList := RegMatchAll(medicalRecord, `<`+infoItem[1]+`>[\s\S]*?<\/`+infoItem[1]+`>`)
	// 		for _, injuryCode := range injuryCodeList {
	// 			injuryCodeVal := utils.GetNodeValue(injuryCode, infoItem[2])
	// 			if utils.HasItem(codes, injuryCodeVal) {
	// 				xmlValue = strings.Replace(xmlValue, medicalRecord, "", 1)
	// 			} else {
	// 				codes = append(codes, injuryCodeVal)
	// 			}
	// 		}
	// 	}
	// }

	injury_infos := RegMatchAll(xmlValue, `<injury_info>[\s\S]*?<\/injury_info>`)
	for _, injury_info := range injury_infos {
		injury_info_clone := injury_info
		// fmt.Println("----------injury_infoinjury_info--------", RegIsMatch(injury_info, `<injury_code>.+<\/injury_code>`))
		if strings.Index(injury_info, "<injury_code>") == -1 {
			add := "<injury_records>\n      <injury_record>\n        <injury_code>R69.X51</injury_code>\n        <injury_remark></injury_remark>\n      </injury_record>"
			injury_info = strings.Replace(injury_info, "<injury_records>", add, 1)
			xmlValue = strings.Replace(xmlValue, injury_info_clone, injury_info, 1)
		} else if !RegIsMatch(injury_info, `<injury_code>.+<\/injury_code>`) {
			injury_info = strings.Replace(injury_info, "<injury_code></injury_code>", "<injury_code>R69.X51</injury_code>", 1)
			// fmt.Println("----------injury_info--------", injury_info)
			xmlValue = strings.Replace(xmlValue, injury_info_clone, injury_info, 1)
		}
		// medicalCountArr := RegMatchAll(medicalInfo, `<medical_count>[\s\S]*?<\/medical_count>`)
		// for _, medicalCount := range medicalCountArr {
		// 	medicalCountClone := medicalCount
		// 	medicalCount = utils.SetNodeValue(medicalCount, "medical_count", strconv.Itoa(len(medicalRecordArr)))
		// 	xmlValue = strings.Replace(xmlValue, medicalCountClone, medicalCount, 1)
		// }
	}

	infoItems := [][]string{{"injury_info", "injury_record", "injury_code", "injury_count"}, {"injury_info_bj", "injury_record", "injury_code", "injury_count"}, {"operation_info", "operation_record", "operation_code", "operation_count"}, {"operation_info_bj", "operation_record", "operation_code", "operation_count"}, {"beneficiary_info", "beneficiary_record", "beneficiary", "beneficiary_count"}}
	for _, infoItem := range infoItems {
		infoLists := RegMatchAll(xmlValue, `<`+infoItem[0]+`>[\s\S]*?<\/`+infoItem[0]+`>`)
		for _, infoList := range infoLists {
			codes := []string{}
			infoListOld := infoList
			recordLists := RegMatchAll(infoList, `<`+infoItem[1]+`>[\s\S]*?<\/`+infoItem[1]+`>`)
			count := 0
			for _, recordList := range recordLists {
				// var beneficiaryRecords []string
				if RegIsMatch(infoList, `>.+<`) {
					value := utils.GetNodeValue(recordList, infoItem[2])
					//1.屏蔽一整个beneficiary_record
					if value == "" {
						infoList = strings.Replace(infoList, recordList, "", 1)
					} else {
						if utils.HasItem(codes, value) {
							infoList = strings.Replace(infoList, recordList, "", 1)
						} else {
							codes = append(codes, value)
							count++
						}
					}
				} else {
					// if bb == 0 {
					// 	infoList = strings.Replace(infoList, recordList, "<"+infoItem[1]+"></"+infoItem[1]+">", 1)
					// } else {
					infoList = strings.Replace(infoList, recordList, "", 1)
					// }
				}
			}
			infoList = utils.SetNodeValue(infoList, infoItem[3], strconv.Itoa(count))
			xmlValue = strings.Replace(xmlValue, infoListOld, infoList, 1)
		}
		// codes := []string{}
		// for _, medicalRecord := range itemList {
		// 	injuryCodeList := RegMatchAll(medicalRecord, `<`+infoItem[1]+`>[\s\S]*?<\/`+infoItem[1]+`>`)
		// 	for _, injuryCode := range injuryCodeList {
		// 		injuryCodeVal := utils.GetNodeValue(injuryCode, infoItem[2])
		// 		if utils.HasItem(codes, injuryCodeVal) {
		// 			xmlValue = strings.Replace(xmlValue, medicalRecord, "", 1)
		// 		} else {
		// 			codes = append(codes, injuryCodeVal)
		// 		}
		// 	}
		// }
	}

	//编码 CSB0122RC0124000 1.当beneficiary_record下的beneficiary为空时，屏蔽一整个beneficiary_record大节点
	//2.当5个beneficiary_record节点下的子节点的值均为空时，仅保留1个beneficiary_record父节点，节点下其余子节点进行屏蔽，其余4个beneficiary_record大节点全部屏蔽；
	//3.每一个beneficiary_info大节点下的beneficiary_count节点的值，根据有多少组beneficiary_record节点进行数量赋值；
	// var count int64
	// beneficiaryRecordArr := RegMatchAll(xmlValue, `<beneficiary_record>[\s\S]*?<\/beneficiary_record>`)
	// beneficiaryInfo := RegMatchAll(xmlValue, `<beneficiary_info>[\s\S]*?<\/beneficiary_info>`)[0]
	// beneficiary_count := 0
	// for bb, beneficiaryRecord := range beneficiaryRecordArr {
	// 	// var beneficiaryRecords []string
	// 	if RegIsMatch(beneficiaryInfo, `>.+<`) {
	// 		value := utils.GetNodeValue(beneficiaryRecord, "beneficiary")
	// 		//1.屏蔽一整个beneficiary_record
	// 		if value == "" {
	// 			xmlValue = strings.Replace(xmlValue, beneficiaryRecord, "", 1)
	// 		} else {
	// 			beneficiary_count++
	// 		}
	// 	} else {
	// 		if bb == 0 {
	// 			xmlValue = strings.Replace(xmlValue, beneficiaryRecord, "<beneficiary_record></beneficiary_record>", 1)
	// 		} else {
	// 			xmlValue = strings.Replace(xmlValue, beneficiaryRecord, "", 1)
	// 		}
	// 	}
	// }
	// xmlValue = utils.SetNodeValue(xmlValue, "beneficiary_count", strconv.Itoa(beneficiary_count))

	// CSB0122RC0125000 1.当operation_record下的operation_code为空时，屏蔽一整个operation_record大节点；
	//2.当所有operation_record节点下的子节点的值均为空时，仅保留1个operation_record父节点，且节点下其余子节点进行屏蔽；
	//3.每一个operation_info、operation_info_bj大节点下的operation_count节点的值，根据有多少组operation_record节点进行数量赋值；
	// operationRecordsArr := RegMatchAll(xmlValue, `<operation_records>[\s\S]*?<\/operation_records>`)
	// nodes := []string{"operation_info", "operation_info_bj"}
	// jj := 0
	// for _, operationRecords := range operationRecordsArr {
	// 	//operationRecordsClone := operationRecords
	// 	operationRecordArr := RegMatchAll(operationRecords, `<operation_record>[\s\S]*?<\/operation_record>`)
	// 	operationRecordCount := 0
	// 	var operationRecordValArr []string
	// 	operation_count := len(operationRecordArr)
	// 	for _, operationRecord := range operationRecordArr {
	// 		operationCodeArr := RegMatchAll(operationRecord, `<operation_code>[\s\S]*?<\/operation_code>`)
	// 		for _, operationCode := range operationCodeArr {
	// 			value := utils.GetNodeValue(operationCode, "operation_code")
	// 			if value == "" {
	// 				xmlValue = strings.Replace(xmlValue, operationRecord, "", 1)
	// 				operation_count--
	// 			}
	// 		}
	// 		//---------------------------------------------------------------
	// 		//2.当所有operation_record节点下的子节点的值均为空时，仅保留1个operation_record父节点，且节点下其余子节点进行屏蔽；
	// 		//匹配<xxx>XXX</xxx>
	// 		nodeRegex := `<([^/]+)>([^<]*)<\/[^>]+>`
	// 		r := regexp.MustCompile(nodeRegex)
	// 		matches := r.FindAllStringSubmatch(operationRecord, -1)
	// 		for _, match := range matches {
	// 			//nodeName := match[1]
	// 			nodeValue := match[2]
	// 			if nodeValue != "" {
	// 				fmt.Println("------------nodeValue-------------", nodeValue)
	// 				operationRecordValArr = append(operationRecordValArr, nodeValue)
	// 				continue
	// 			}
	// 		}
	// 		if len(operationRecordValArr) == 0 {
	// 			operationRecordCount++
	// 		}
	// 		for _, node := range nodes {
	// 			combineArr := utils.RegMatchAll(xmlValue, `<`+node+`>[\s\S]*?<\/`+node+`>`)
	// 			for _, combine := range combineArr {
	// 				combineClone := combine
	// 				combine = utils.SetNodeValue(combine, "operation_count", strconv.Itoa(operation_count))
	// 				xmlValue = strings.Replace(xmlValue, combineClone, combine, 1)
	// 			}
	// 		}
	// 	}
	// 	if operationRecordCount == len(operationRecordArr) {
	// 		jj++
	// 		//xmlValue = strings.Replace(xmlValue, operationRecords, "<operation_records>	\n\t\t\t\t\t<operation_record>\n\t\t\t\t\t</operation_record>\n\t\t\t\t</operation_records>", 1)
	// 	}
	// }
	// operationRecordsArrs := RegMatchAll(xmlValue, `<operation_records>[\s\S]*?<\/operation_records>`)
	// for _, operations := range operationRecordsArrs {
	// 	if jj > 0 {
	// 		xmlValue = strings.Replace(xmlValue, operations, "<operation_records>	\n\t\t\t\t\t<operation_record>\n\t\t\t\t\t</operation_record>\n\t\t\t\t</operation_records>", -1)
	// 	}
	// }

	// //编码 CSB0122RC0126000 1.当injury_record下的injury_code为空时，屏蔽一整个injury_record大节点；
	// //2.当所有injury_record节点下的子节点的值均为空时，仅保留1个injury_record父节点，且节点下其余子节点进行屏蔽；
	// //3.每一个injury_info、injury_info_bj大节点下的injury_count节点的值，根据有多少组injury_record节点进行数量赋值；
	// injuryRecordsCount := 0
	// injuryRecordsArr := RegMatchAll(xmlValue, `<injury_records>[\s\S]*?<\/injury_records>`)
	// for _, injuryRecords := range injuryRecordsArr {
	// 	injuryRecordArr := RegMatchAll(injuryRecords, `<injury_record>[\s\S]*?<\/injury_record>`)
	// 	var injuryRecordValArr []string
	// 	var injuryRecordCount int
	// 	for _, injuryRecord := range injuryRecordArr {
	// 		injuryCodeArr := RegMatchAll(injuryRecord, `<injury_code>[\s\S]*?<\/injury_code>`)
	// 		for _, injuryCode := range injuryCodeArr {
	// 			value := utils.GetNodeValue(injuryCode, "injury_code")
	// 			if value == "" {
	// 				xmlValue = strings.Replace(xmlValue, injuryRecord, "", 1)
	// 			}
	// 		}

	// 		//2.当所有operation_record节点下的子节点的值均为空时，仅保留1个operation_record父节点，且节点下其余子节点进行屏蔽；
	// 		//匹配<xxx>XXX</xxx>
	// 		nodeRegex := `<([^/]+)>([^<]*)<\/[^>]+>`
	// 		r := regexp.MustCompile(nodeRegex)
	// 		matches := r.FindAllStringSubmatch(injuryRecord, -1)
	// 		for _, match := range matches {
	// 			nodeValue := match[2]
	// 			if nodeValue != "" {
	// 				injuryRecordValArr = append(injuryRecordValArr, nodeValue)
	// 				continue
	// 			}
	// 		}
	// 		if len(injuryRecordValArr) == 0 {
	// 			injuryRecordCount++
	// 		}
	// 	}
	// 	if injuryRecordCount == len(injuryRecordArr) && len(injuryRecordValArr) == 0 {
	// 		injuryRecordsCount++
	// 		//xmlValue = strings.Replace(xmlValue, injuryRecords, "<injury_records>	\n\t\t\t\t\t<injury_record>\n\t\t\t\t\t</injury_record>\n\t\t\t\t</injury_records>", 1)
	// 	}
	// }
	// injuryRecordsArr = RegMatchAll(xmlValue, `<injury_records>[\s\S]*?<\/injury_records>`)
	// for _, injuryRecords := range injuryRecordsArr {
	// 	if injuryRecordsCount > 0 {
	// 		xmlValue = strings.Replace(xmlValue, injuryRecords, "<injury_records>	\n\t\t\t\t\t<injury_record>\n\t\t\t\t\t</injury_record>\n\t\t\t\t</injury_records>", 1)
	// 	}
	// }

	//编码 CSB0122RC0128000 每一个medical_info大节点下的medical_count节点的值，根据有多少组medical_record节点进行数量赋值
	//medical_record节点下的子节点的值均为空时，保留该节点，节点下其余子节点进行屏蔽
	medicalInfoArr := RegMatchAll(xmlValue, `<medical_info>[\s\S]*?<\/medical_info>`)
	medicalRecordArr := RegMatchAll(xmlValue, `<medical_record>[\s\S]*?<\/medical_record>`)
	for _, medicalInfo := range medicalInfoArr {
		medicalCountArr := RegMatchAll(medicalInfo, `<medical_count>[\s\S]*?<\/medical_count>`)
		for _, medicalCount := range medicalCountArr {
			medicalCountClone := medicalCount
			medicalCount = utils.SetNodeValue(medicalCount, "medical_count", strconv.Itoa(len(medicalRecordArr)))
			xmlValue = strings.Replace(xmlValue, medicalCountClone, medicalCount, 1)
		}
	}
	// medicalRecordVal := []string{}
	// for _, medicalRecord := range medicalRecordArr {
	// 	// 匹配<xxx>XXX</xxx>
	// 	nodeRegex := `<([^/]+)>([^<]*)<\/[^>]+>`
	// 	r := regexp.MustCompile(nodeRegex)
	// 	matches := r.FindAllStringSubmatch(medicalRecord, -1)
	// 	for _, match := range matches {
	// 		//nodeName := match[1]
	// 		nodeValue := match[2]
	// 		if nodeValue != "" {
	// 			medicalRecordVal = append(medicalRecordVal, nodeValue)
	// 			continue
	// 		}
	// 		if len(medicalRecordVal) == 0 {
	// 			xmlValue = strings.Replace(xmlValue, medicalRecord, "<medical_record></medical_record>", 1)
	// 		}
	// 	}
	// }
	//fmt.Println("---------------medicalRecordVal,len(medicalRecordVal)", medicalRecordVal, len(medicalRecordVal))

	// 编码 CSB0122RC0129000
	//1.当fee_record下的fee_name、stringge_name节点值为空时，屏蔽一整个fee_record大节点；
	//2.当所有fee_record节点下的子节点的值均为空时，仅保留1个fee_record节点，且节点下其余子节点进行屏蔽；
	//3.每一个fee_details大节点下的fee_count节点的值，根据有多少组fee_record节点进行数量赋值；
	feeDetails := RegMatchAll(xmlValue, `<fee_details>[\s\S]*?<\/fee_details>`)
	for _, detail := range feeDetails {
		detail_clone := detail
		feeRecords := RegMatchAll(detail, `<fee_record>[\s\S]*?<\/fee_record>`)
		fee_count := len(feeRecords)
		for _, record := range feeRecords {
			record_clone := record
			null_value := ""
			if strings.Index(record, "<fee_name>") != -1 {
				null_value = utils.GetNodeValue(record, "fee_name")
			}
			if strings.Index(record, "<stringge_name>") != -1 {
				null_value = utils.GetNodeValue(record, "stringge_name")
			}
			// fee_name := utils.GetNodeValue(record, "fee_name")
			// stringge_name := utils.GetNodeValue(record, "stringge_name")
			if null_value == "" {
				detail = strings.Replace(detail, record, "", 1)
				fee_count--
			} else {
				record = strings.ReplaceAll(record, "%", "")
				detail = strings.Replace(detail, record_clone, record, 1)
			}

		}
		detail = utils.SetNodeValue(detail, "fee_count", strconv.Itoa(fee_count))
		xmlValue = strings.Replace(xmlValue, detail_clone, detail, 1)
	}
	// nodeArr := []string{"fee_name", "stringge_name"}
	// feeDetails := RegMatchAll(xmlValue, `<fee_details>[\s\S]*?<\/fee_details>`)
	// shieldCount := 0
	// for _, detail := range feeDetails {
	// 	feeCounts := RegMatchAll(detail, `<fee_count>[\s\S]*?<\/fee_count>`)
	// 	feeRecords := RegMatchAll(detail, `<fee_records>[\s\S]*?<\/fee_records>`)
	// 	for _, record := range feeRecords {
	// 		feeRecordArr := RegMatchAll(record, `<fee_record>[\s\S]*?<\/fee_record>`)
	// 		fee_count := len(feeRecordArr)
	// 		feeRecordChild := []string{}
	// 		feeCount := 0
	// 		for _, feeRecord := range feeRecordArr {
	// 			for _, item := range nodeArr {
	// 				itemNames := RegMatchAll(feeRecord, `<`+item+`>[\s\S]*?<\/`+item+`>`)
	// 				for _, NodeName := range itemNames {
	// 					value := utils.GetNodeValue(NodeName, item)
	// 					if value == "" {
	// 						xmlValue = strings.Replace(xmlValue, feeRecord, "", 1)
	// 						fee_count--
	// 					}
	// 				}
	// 			}
	// 			//2.当所有fee_record节点下的子节点的值均为空时，仅保留1个fee_record节点，且节点下其余子节点进行屏蔽；
	// 			nodeRegex := `<([^/]+)>([^<]*)<\/[^>]+>`
	// 			r := regexp.MustCompile(nodeRegex)
	// 			matches := r.FindAllStringSubmatch(feeRecord, -1)
	// 			for _, match := range matches {
	// 				nodeValue := match[2]
	// 				if nodeValue != "" {
	// 					feeRecordChild = append(feeRecordChild, nodeValue)
	// 					continue
	// 				}
	// 			}
	// 			if len(feeRecordChild) == 0 {
	// 				feeCount++
	// 			}
	// 		}
	// 		if feeCount == len(feeRecordArr) && len(feeRecordChild) == 0 {
	// 			shieldCount++
	// 		}

	// 		for _, feeRecordCount := range feeCounts {
	// 			feeRecordCountClone := feeRecordCount
	// 			feeRecordCount = utils.SetNodeValue(feeRecordCount, "fee_count", strconv.Itoa(fee_count))
	// 			xmlValue = strings.Replace(xmlValue, feeRecordCountClone, feeRecordCount, 1)
	// 		}
	// 	}
	// }
	// feeRecords := RegMatchAll(xmlValue, `<fee_records>[\s\S]*?<\/fee_records>`)
	// for _, record := range feeRecords {
	// 	if shieldCount > 0 {
	// 		xmlValue = strings.Replace(xmlValue, record, "<fee_records>	\n\t\t\t\t\t<fee_record>\n\t\t\t\t\t</fee_record>\n\t\t\t\t</fee_records>", 1)
	// 	}
	// }
	// 编码 CSB0122RC0123000 medical_record大节点下的【hospital_bill】、【clinc_bill】、【pharmacy_bill】、【non_medical_bill】对应节点的bill_count节点值为0时，只保留bill_records父节点，该节点下其余子节点进行屏蔽
	//CSB0122RC0133000 medical_record大节点下的【hospital_bill】、【clinc_bill】、【pharmacy_bill】、【non_medical_bill】对应节点的bill_count节点值，根据同一个大节点下的bill_record节点进行数量赋值
	medicalRecordArrays := RegMatchAll(xmlValue, `<medical_record>[\s\S]*?<\/medical_record>`)
	itemArr := []string{"hospital_bill", "clinc_bill", "pharmacy_bill", "non_medical_bill"}
	for _, item := range itemArr {
		for _, medicalRecord := range medicalRecordArrays {
			combineArr := utils.RegMatchAll(medicalRecord, `<`+item+`>[\s\S]*?<\/`+item+`>`)
			for _, combine := range combineArr {
				// billCounts := utils.RegMatchAll(combine, `<bill_count>[\s\S]*?<\/bill_count>`)
				combineClone := combine
				// value := utils.GetNodeValue(combine, "bill_count")
				// if value == "0" {
				// 	billRecordsArr := RegMatchAll(combine, `<bill_records>[\s\S]*?<\/bill_records>`)
				// 	xmlValue = strings.Replace(xmlValue, billRecordsArr[0], "<bill_records></bill_records>", 1)
				// }
				billRecord := RegMatchAll(combine, `<bill_record>[\s\S]*?<\/bill_record>`)
				combine = utils.SetNodeValue(combine, "bill_count", strconv.Itoa(len(billRecord)))
				xmlValue = strings.Replace(xmlValue, combineClone, combine, 1)
				// for _, record := range billRecord {
				// 	billRecordArray := RegMatchAll(record, `<bill_record>[\s\S]*?<\/bill_record>`)
				// 	combine = utils.SetNodeValue(combine, "bill_count", strconv.Itoa(len(billRecordArray)))
				// 	xmlValue = strings.Replace(xmlValue, combineClone, combine, 1)

				// 	for _, billCount := range billCounts {
				// 		billCountClone := billCount
				// 		billCount = utils.SetNodeValue(billCount, "bill_count", strconv.Itoa(len(billRecordArray)))
				// 		xmlValue = strings.Replace(xmlValue, billCountClone, billCount, 1)

				// 	}
				// }

			}
		}
	}

	//编码 CSB0122RC0130000
	//同一个medical_record大节点下，校验每一个fee_name节点值，当其中一个fee_name节点值包含：重症监护*病房、ICU*病房、CCU*病房（重症加强护理病房）、NICU*病房（新生儿重症监护病房）字样时
	//同一个medical_record大节点下的icudate_in、icudate_out两个节点值均默认为99999999
	medicalRecords := RegMatchAll(xmlValue, `<medical_record>[\s\S]*?<\/medical_record>`)
	for _, medicalRecord := range medicalRecords {
		medicalRecordClone := medicalRecord
		feeNameArr := RegMatchAll(medicalRecord, `<fee_name>[\s\S]*?<\/fee_name>`)
		for _, feeName := range feeNameArr {
			feeVal := utils.GetNodeValue(feeName, "fee_name")
			if RegIsMatch(feeVal, `(重症监护\*病房|ICU\*病房|CCU\*病房\(重症加强护理病房\)|NICU\*病房\(新生儿重症监护病房\))`) {
				medicalRecord = utils.SetNodeValue(medicalRecord, "icudate_in", "99999999")
				medicalRecord = utils.SetNodeValue(medicalRecord, "icudate_out", "99999999")
				xmlValue = strings.Replace(xmlValue, medicalRecordClone, medicalRecord, 1)

			}
		}
	}

	// 编码 CSB0122RC0138000 机构为“北京”，屏蔽一整个injury_info、operation_info大节点  机构不为“北京”，屏蔽一整个injury_info_bj、operation_info_bj大节点
	//injuryAndOperation := []string{"injury_info", "operation_info"}
	nodeTwos := []string{"injury_info_bj", "operation_info_bj"}
	//if obj.Bill.Agency == "北京" {
	//	for _, node := range injuryAndOperation {
	//		items := utils.RegMatchAll(xmlValue, `<`+node+`>[\s\S]*?<\/`+node+`>`)
	//		for _, item := range items {
	//			xmlValue = strings.Replace(xmlValue, item, "", 1)
	//		}
	//	}
	//} else
	if obj.Bill.Agency != "北京" {
		for _, nodeTwo := range nodeTwos {
			NotBeijin := utils.RegMatchAll(xmlValue, `<`+nodeTwo+`>[\s\S]*?<\/`+nodeTwo+`>`)
			for _, notBeijing := range NotBeijin {
				xmlValue = strings.Replace(xmlValue, notBeijing, "", 1)
			}
		}
	}
	//编码 CSB0122RC0142000 ;  当injury_code（可能存在多个）节点值包含“W53.951、W54.951、W55.951、W55.952、W56.851、W56.852、W57.951、W58.951、W59.951、W60.951、W64.951、W64.952、W64.953、W64.954”时，
	//屏蔽question_code为C18、C20、C21、C22的question一整个大节点
	//当injury_code（可能存在多个）节点值包含“W53.951、W54.951、W55.951、W55.952、W56.851、W56.852、W57.951、W58.951、W59.951、W60.951、W64.951、W64.952、W64.953、W64.954”时，
	//contingency_reason（节点唯一）节点值默认为X20-X29，contingency_address（节点唯一）节点值默认为CS31
	injuryCodeArr := RegMatchAll(xmlValue, `<injury_code>[\s\S]*?<\/injury_code>`)
	for _, injuryCode := range injuryCodeArr {
		//injuryCodeClone := injuryCode
		injuryCodeVal := utils.GetNodeValue(injuryCode, "injury_code")
		if RegIsMatch(injuryCodeVal, `^(W53.951|W54.951|W55.951|W55.952|W56.851|W56.852|W57.951|W58.951|W59.951|W60.951|W64.951|W64.952|W64.953|W64.954)$`) {
			//contingencyReasonArr := RegMatchAll(xmlValue, `<contingency_reason>[\s\S]*?<\/contingency_reason>`)
			//xmlValue = utils.SetNodeValue(xmlValue, contingencyReasonArr[0], "X20-X29")
			//contingencyAddressArr := RegMatchAll(xmlValue, `<contingency_address>[\s\S]*?<\/contingency_address>`)
			//xmlValue = utils.SetNodeValue(xmlValue, contingencyAddressArr[0], "CS31")
			questionArr := RegMatchAll(xmlValue, `<question>[\s\S]*?<\/question>`)
			for _, question := range questionArr {
				//questionClon := question
				questionVal := utils.GetNodeValue(question, "question_code")
				if RegIsMatch(questionVal, `^(C18|C20|C21|C22)$`) {
					xmlValue = strings.Replace(xmlValue, question, "", 1)
				}
			}
		}
	}

	//CSB0122RC0143000
	//同一个medical_record大节点下，当injury_code（可能存在多个）节点值包含“W53.951、W54.951、W55.951、W55.952、W56.851、W56.852、W57.951、W58.951、W59.951、W60.951、W64.951、W64.952、W64.953、W64.954”时，
	//与其同一个medical_record大节点下contingency_reason（节点唯一）节点值默认为X20-X29，contingency_address（节点唯一）节点值默认为CS31
	//（2） 编码 CSB0122RC0169000  同一个medical_record大节点下，校验hospital_bill大节点下bill_count的值，当bill_count不为0时，清空所有operation_date的节点值
	medicalRecordList := RegMatchAll(xmlValue, `<medical_record>[\s\S]*?<\/medical_record>`)
	//operationDateList := RegMatchAll(xmlValue, `<operation_date>[\s\S]*?<\/operation_date>`)
	for _, medicalRecord := range medicalRecordList {
		medicalRecordClone := medicalRecord

		hospitalBillList := RegMatchAll(medicalRecord, `<hospital_bill>[\s\S]*?<\/hospital_bill>`)
		for _, hospitalBill := range hospitalBillList {
			billCountList := RegMatchAll(hospitalBill, `<bill_count>[\s\S]*?<\/bill_count>`)
			for _, billCount := range billCountList {
				value := utils.GetNodeValue(billCount, "bill_count")
				if value != "0" {
					medicalRecord = utils.SetNodeValue(medicalRecord, "operation_date", "")
					//operation_dates := RegMatchAll(xmlValue, `<operation_date>[\s\S]*?<\/operation_date>`)
				}
			}
		}

		injuryCodeList := RegMatchAll(medicalRecord, `<injury_code>[\s\S]*?<\/injury_code>`)
		for _, injuryCode := range injuryCodeList {
			injuryCodeVal := utils.GetNodeValue(injuryCode, "injury_code")
			if RegIsMatch(injuryCodeVal, `^(W53.951|W54.951|W55.951|W55.952|W56.851|W56.852|W57.951|W58.951|W59.951|W60.951|W64.951|W64.952|W64.953|W64.954)$`) {
				medicalRecord = utils.SetNodeValue(medicalRecord, "contingency_reason", "X20-X29")
				medicalRecord = utils.SetNodeValue(medicalRecord, "contingency_address", "CS31")
			}
		}

		xmlValue = strings.Replace(xmlValue, medicalRecordClone, medicalRecord, 1)

		//（2） 同一个medical_record大节点下，校验hospital_bill大节点下bill_count的值，当bill_count不为0时，清空所有operation_date的节点值

	}

	// medical_records := RegMatchAll(xmlValue, `<medical_record>[\s\S]*?<\/medical_record>`)
	// for _, medical_record := range medical_records {
	// 	medical_record_clone := medical_record
	// 	hospital_bills := RegMatchAll(xmlValue, `<hospital_bill>[\s\S]*?<\/hospital_bill>`)
	// 	if len(hospital_bills) > 0 {
	// 		bill_count := utils.GetNodeValue(hospital_bills[0], "bill_count")
	// 		if bill_count == "0" {
	// 			medical_record = utils.SetNodeValue(medical_record, "medicaldate_in", "")
	// 			medical_record = utils.SetNodeValue(medical_record, "medicaldate_out", "")
	// 		}
	// 	}
	// 	xmlValue = strings.Replace(xmlValue, medical_record_clone, medical_record, 1)
	// }

	if obj.Bill.InsuranceType == "重大疾病" {
		severe_disease_date := utils.GetNodeValue(xmlValue, "severe_disease_date")
		severe_disease_reportdate := utils.GetNodeValue(xmlValue, "severe_disease_reportdate")
		if severe_disease_date == "" && severe_disease_reportdate == "" {
			medical_dates := RegMatchAll(xmlValue, `<medical_date>.+<\/medical_date>`)
			value = ""
			for _, medical_date := range medical_dates {
				aa := utils.GetNodeValue(medical_date, "medical_date")
				a, _ := time.Parse("20060102", aa)
				b, _ := time.Parse("20060102", value)
				if a.Before(b) {
					value = aa
				}
			}
			xmlValue = utils.SetNodeValue(xmlValue, "severe_disease_date", value)
			xmlValue = utils.SetNodeValue(xmlValue, "severe_disease_reportdate", value)
		}
	}

	//CSB0122RC0144000 当一整张单不存在fc196字段时，police_intervention节点值默认0（节点唯一）
	//medicalRecordList := RegMatchAll(xmlValue, `<police_intervention>[\s\S]*?<\/police_intervention>`)
	isfc196 := false
	for _, field := range obj.Fields {
		if field.Code == "fc196" {
			isfc196 = true
			break
		}
	}
	if !isfc196 {
		xmlValue = utils.SetNodeValue(xmlValue, "police_intervention", "0")
	}

	itemNodes := [][]string{
		{"beneficiary_record", "beneficiary_records"},
		{"medical_record", "medical_records"},
		{"bill_record", "bill_records"},
		{"fee_record", "fee_records"},
		{"injury_record", "injury_records"},
		{"operation_record", "operation_records"},
		{"question", "questions"},
		{"pathological_record", "pathological_records"},
	}

	for _, itemNode := range itemNodes {
		nodeDatas := RegMatchAll(xmlValue, `<`+itemNode[1]+`>[\s\S]*?<\/`+itemNode[1]+`>`)
		for _, nodeData := range nodeDatas {
			cloneNodeData := nodeData
			aaa := "<" + itemNode[1] + "><array_list></array_list></" + itemNode[1] + ">"
			if strings.Index(nodeData, "<"+itemNode[0]+">") != -1 {
				nodeData = utils.RegReplace(nodeData, `<`+itemNode[1]+`>`, "")
				nodeData = utils.RegReplace(nodeData, `</`+itemNode[1]+`>`, "")
				nodeData = utils.RegReplace(nodeData, `<`+itemNode[0]+`>`, `<`+itemNode[1]+`>`)
				nodeData = utils.RegReplace(nodeData, `<\/`+itemNode[0]+`>`, `</`+itemNode[1]+`>`)
				nodeData += aaa
			} else {
				nodeData = aaa + aaa
			}
			xmlValue = strings.Replace(xmlValue, cloneNodeData, nodeData, 1)
		}
	}

	medical_records := RegMatchAll(xmlValue, `<medical_records>[\s\S]*?<\/medical_records>`)
	for _, medical_record := range medical_records {
		medical_record_clone := medical_record
		medical_date := utils.GetNodeValue(medical_record, "medical_date")
		// fmt.Println("----------medical_date------------:", medical_date)
		if medical_date == "" {
			value := ""
			// 	a, _ := time.Parse("2006-01-02", endDate)
			// b, _ := time.Parse("2006-01-02", beginDate)
			billdate_ins := RegMatchAll(medical_record, `<(billdate_in|bill_date|medicaldate_in)>.*<\/(billdate_in|bill_date|medicaldate_in)>`)
			for _, billdate_in := range billdate_ins {
				aa := GetNodeValue(billdate_in, "billdate_in")
				if aa == "" {
					aa = GetNodeValue(billdate_in, "bill_date")
				}
				if aa == "" {
					aa = GetNodeValue(billdate_in, "medicaldate_in")
				}
				// fmt.Println("----------billdate_in------------:", billdate_in, aa)
				if aa != "" {
					if value == "" {
						value = aa
					} else {
						a, _ := time.Parse("20060102", aa)
						b, _ := time.Parse("20060102", value)
						if a.Before(b) {
							value = aa
						}
					}
				}
			}
			medical_record = utils.SetNodeValue(medical_record, "medical_date", value)
			xmlValue = strings.Replace(xmlValue, medical_record_clone, medical_record, 1)
		}
	}

	itemNodes = [][]string{
		{"apply_name", "accident_name"},
		{"apply_term_start", "accident_term_start"},
		{"apply_term_close", "accident_term_close"},
	}
	for _, itemNode := range itemNodes {
		aa := GetNodeValue(xmlValue, itemNode[0])
		if aa == "" {
			bb := GetNodeValue(xmlValue, itemNode[1])
			xmlValue = utils.SetNodeValue(xmlValue, itemNode[0], bb)
		}

	}
	// 2024-03-05 取消
	//contingency_reason := GetNodeValue(xmlValue, "contingency_reason")
	//if contingency_reason != "" {
	//	medical_dates := RegMatchAll(xmlValue, `<medical_date>.*<\/medical_date>`)
	//	value := ""
	//	for _, medical_date := range medical_dates {
	//		aa := GetNodeValue(medical_date, "medical_date")
	//		if value == "" {
	//			value = aa
	//		} else {
	//			a, _ := time.Parse("20060102", aa)
	//			b, _ := time.Parse("20060102", value)
	//			if a.Before(b) {
	//				value = aa
	//			}
	//		}
	//	}
	//	xmlValue = utils.SetNodeValue(xmlValue, "contingency_date", value)
	//}

	question_codes := RegMatchAll(xmlValue, `<question_code>.+<\/question_code>`)
	xmlValue = utils.SetNodeValue(xmlValue, "question_count", strconv.Itoa(len(question_codes)))

	newXmlValue = xmlValue
	return err, newXmlValue
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
func RegMatchAll(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.FindAllString(value, -1)
}
func RegIsMatch(value string, query string) bool {
	// reg := regexp.MustCompile(query)
	matched, _ := regexp.MatchString(query, value)
	return matched
}

func RegSplit(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.Split(value, -1)
}
