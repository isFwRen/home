/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 5:34 下午
 */

package B0108

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/load/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/wxnacy/wgo/arrays"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0108:::XmlDeal-----------------------")
	// return nil, xmlValue
	obj := o.(FormatObj)
	// xmlValue = SetNodeValue(xmlValue, "claimNum", obj.Bill.BillNum)
	// constMap := constDeal(obj.Bill.ProCode)
	c, ok := global.GProConf[obj.Bill.ProCode].ConstTable["B0108_太平理赔_诊断代码表"]
	tempMap := make(map[string]string, 0)
	if ok {
		for _, arr := range c {
			tempMap[strings.TrimSpace(arr[0])] = arr[2]
		}
	}
	fields := obj.Fields
	errorData := ""
	errMeses := ""
	questionCount := 0
	fcxxx := []string{}
	isFc032 := false
	isFc054 := false
	_, fc264Val := GetOneField(fields, "fc264", true)
	for _, field := range fields {
		if (field.Code == "fc033" || field.Code == "fc059") && field.FinalValue != "" && !utils.HasItem(fcxxx, field.FinalValue) {
			fcxxx = append(fcxxx, field.FinalValue)
		}
		// if len(filed.Issues ) > 0 {
		for _, issue := range field.Issues {
			errMes := issue.Code + "_" + issue.Message + ";"
			if strings.Index(errMeses, errMes) == -1 {
				if questionCount == 0 {
					errorData = "	  <QuestionCount></QuestionCount>\n"
				}
				questionCount = questionCount + 1
				questionType := 1
				if (obj.Bill.Agency == "83301" || obj.Bill.Agency == "83311") && !RegIsMatch(issue.Message, `(模糊|不清晰)`) {
					questionType = 2
				}
				errorData = errorData + "	  <Question>\n	   <QuestionType>" + strconv.Itoa(questionType) + "</QuestionType>\n	   <QuestionContent>" + issue.Message + "</QuestionContent>\n	  </Question>\n"
				errMeses += errMes
			}
		}

		if field.Code == "fc032" {
			if field.FinalValue == fc264Val {
				isFc032 = true
			}
		}
		if field.Code == "fc054" {
			if field.FinalValue == fc264Val {
				isFc054 = true
			}
		}
		// error = error + "  <errorList>\n   <errorCode>#{field.final_fields.issue.id}</errorCode>\n   <errorDesc>#{field.final_fields.issue.text}</errorDesc>\n  </errorList>\n"
		// }
	}
	if len(fcxxx) > 1 {
		questionCount++
		errorData = errorData + "	  <Question>\n	   <QuestionType>1</QuestionType>\n	   <QuestionContent>事故者姓名不一致</QuestionContent>\n	  </Question>\n"
	}
	if errorData == "" {
		errorData = "	  <QuestionCount>0</QuestionCount>\n"
	}
	fmt.Println("------------------------error-------------------", errorData)
	errorData = "<Questions>\n" + errorData + "	</Questions>"
	xmlValue = strings.Replace(xmlValue, "<Questions></Questions>", errorData, 1)

	if questionCount != 0 {
		xmlValue = SetNodeValue(xmlValue, "return_flag", strconv.Itoa(1))
		xmlValue = SetNodeValue(xmlValue, "QuestionCount", strconv.Itoa(questionCount))
	}

	xmlValue = SetNodeValue(xmlValue, "case_no", obj.Bill.BillNum)

	accidentType := GetAccidentTypeNum(fields, tempMap)
	xmlValue = SetNodeValue(xmlValue, "accident_type", accidentType)

	xmlValue = SetNodeValue(xmlValue, "myversion", Substr(obj.Bill.BillName, 8, 16))

	doctor_name_feilds := []string{"fc183", "fc184", "fc185", "fc186", "fc187"}
	for dd, doctor_name_feild := range doctor_name_feilds {
		_, value := GetOneField(fields, doctor_name_feild, true)
		if value != "" {
			xmlMark := "doctor_name"
			if dd != 0 {
				xmlMark = xmlMark + strconv.Itoa(dd+1)
			}
			xmlValue = RegReplace(xmlValue, `<`+xmlMark+`>.*?<\/`+xmlMark+`>`, "<"+xmlMark+">"+value+"</"+xmlMark+">")
		}
	}

	fee_details := RegMatchAll(xmlValue, `<fee_details>[\s\S]*?<\/fee_details>`)
	for _, fee_detail := range fee_details {
		old_fee_detail := fee_detail
		fee_count := len(RegMatchAll(fee_detail, `<item_name>.+?<\/item_name>`))
		// GetNodeValue(fee_detail, "fee_count")
		fee_detail = SetNodeValue(fee_detail, "fee_count", strconv.Itoa(fee_count))
		// if fee_count == 0 {
		re_fee_details := RegMatchAll(fee_detail, `<fee_detail>[\s\S]*?<\/fee_detail>`)
		for ii, re_fee_detail := range re_fee_details {
			if ii == 0 && fee_count == 0 {
				continue
			}
			item_name := GetNodeValue(re_fee_detail, "item_name")
			if item_name == "" {
				fee_detail = strings.Replace(fee_detail, re_fee_detail, "", 1)
			}
		}
		// }
		xmlValue = strings.Replace(xmlValue, old_fee_detail, fee_detail, 1)

	}

	isExit, fc095 := GetOneField(fields, "fc095", true)
	if isExit && !RegIsMatch(fc095, `^1`) {
		xmlValue = SetNodeValue(xmlValue, "apply_cellphone", "")
		if fc095 != "" {
			xmlValue = SetNodeValue(xmlValue, "apply_tel", fc095)
		}
	}

	deduct1s := RegMatchAll(xmlValue, `<deduct1>[\s\S]*?<\/deduct1>`)
	for _, deduct1 := range deduct1s {
		deduct1Value := GetNodeValue(deduct1, "deduct1")
		if deduct1Value != "0.00" && deduct1Value != "" {
			xmlValue = SetNodeValue(xmlValue, "social_security", "Y")
		}
	}

	accident_id := GetAccidentId(fields)
	accident_name_desc := GetAccidentNameDesc(fields)

	clinic_records := RegMatchAll(xmlValue, `<clinic_record>[\s\S]*?<\/clinic_record>`)
	for r, clinic_record := range clinic_records {
		old_clinic_record := clinic_record
		clinic_record = SetNodeValue(clinic_record, "seq", strconv.Itoa(r+1))
		if accident_id != "" {
			clinic_record = SetNodeValue(clinic_record, "accident_id", accident_id)
		}
		if accident_name_desc != "" {
			clinic_record = SetNodeValue(clinic_record, "accident_name_desc", accident_name_desc)
		}
		fee_count := GetNodeValue(clinic_record, "fee_count")
		fee_details := RegMatchAll(clinic_record, `<fee_detail>[\s\S]*?<\/fee_detail>`)
		// fmt.Println("-----------fee_details-------------:", len(fee_details))
		// if fee_count != "0" {
		for ff, fee_detail := range fee_details {
			// fmt.Println("-----------ff-------------:", ff)
			old_fee_detail := fee_detail
			if fee_count != "0" {
				fee_detail = SetNodeValue(fee_detail, "seq", strconv.Itoa(ff+1))
			} else {
				fee_detail = SetNodeValue(fee_detail, "seq", "")
			}
			clinic_record = strings.Replace(clinic_record, old_fee_detail, fee_detail, 1)
			// fmt.Println("-----------fee_detail-------------:", fee_detail)
		}
		// }

		xmlValue = strings.Replace(xmlValue, old_clinic_record, clinic_record, 1)
	}

	operation_type := operationTypeValue(fields)
	xmlValue = RegReplace(xmlValue, `<operation_type>.*?<\/operation_type>`, "<operation_type>"+operation_type+"</operation_type>")

	inpatient_records := RegMatchAll(xmlValue, `<inpatient_record>[\s\S]*?<\/inpatient_record>`)
	for r, inpatient_record := range inpatient_records {
		old_inpatient_record := inpatient_record
		out_date := GetNodeValue(inpatient_record, "out_date")
		in_date := GetNodeValue(inpatient_record, "in_date")
		if out_date != "" && in_date != "" {
			// time.Parse()
			a, _ := time.Parse("20060102", in_date)
			b, _ := time.Parse("20060102", out_date)
			d := b.Sub(a)
			// fmt.Println(d.Hours() / 24)
			inpatient_record = SetNodeValue(inpatient_record, "bed_amt_days", ToString(d.Hours()/24))
		}

		inpatient_record = SetNodeValue(inpatient_record, "seq", strconv.Itoa(r+1))
		if accident_id != "" {
			inpatient_record = SetNodeValue(inpatient_record, "accident_id", accident_id)
		}
		if accident_name_desc != "" {
			inpatient_record = SetNodeValue(inpatient_record, "accident_name_desc", accident_name_desc)
		}
		fee_count := GetNodeValue(inpatient_record, "fee_count")
		// fmt.Println("------------------------fee_count-------------------", fee_count)
		// if fee_count != "0" {
		fee_details := RegMatchAll(inpatient_record, `<fee_detail>[\s\S]*?<\/fee_detail>`)
		for ff, fee_detail := range fee_details {
			old_fee_detail := fee_detail
			// fee_detail = SetNodeValue(fee_detail, "seq", strconv.Itoa(ff+1))
			if fee_count != "0" {
				fee_detail = SetNodeValue(fee_detail, "seq", strconv.Itoa(ff+1))
			} else {
				fee_detail = SetNodeValue(fee_detail, "seq", "")
			}
			inpatient_record = strings.Replace(inpatient_record, old_fee_detail, fee_detail, 1)
		}
		// }
		xmlValue = strings.Replace(xmlValue, old_inpatient_record, inpatient_record, 1)
	}

	if obj.Bill.SaleChannel == "秒赔" {
		xmlValue = SetNodeValue(xmlValue, "apply_certi_type", "")
		xmlValue = SetNodeValue(xmlValue, "apply_certi_validate", "")
	}

	accident_time := GetNodeValue(xmlValue, "accident_time")
	if accident_time == "" {
		early_time_val := GetEarlyTimeOfFdAndid(xmlValue)
		fmt.Println("------------------------early_time_val-------------------", early_time_val)
		xmlValue = SetNodeValue(xmlValue, "accident_time", early_time_val)
	}

	_, fc217 := GetOneField(fields, "fc217", true)
	if fc217 != "" {
		xmlValue = SetNodeValue(xmlValue, "accident_time", fc217)
	}

	isExit, fc036 := GetOneField(fields, "fc036", true)
	baseinfo := GetNodeValue(xmlValue, "base_info")
	accident_id = GetNodeValue(baseinfo, "accident_id")
	baseinfo_clone := baseinfo
	value_ac := ""
	if isExit {
		value_ac = fc036
	} else {
		isExit, fc112 := GetOneField(fields, "fc112", true)
		if isExit {
			value_ac = fc112
		}
	}
	if accident_id == "" {
		baseinfo_clone = RegReplace(baseinfo_clone, `<accident_id\s*\/>|<accident_id\>.*?\<\/accident_id>`, "<accident_id>"+value_ac+"</accident_id>")
		xmlValue = strings.Replace(xmlValue, baseinfo, baseinfo_clone, 1)
	}

	fee_infos := RegMatchAll(xmlValue, `<fee_info>[\s\S]*?<\/fee_info>`)
	for _, fee_info := range fee_infos {
		old_fee_info := fee_info
		sum_amount := GetNodeValue(fee_info, "sum_amount")
		deduct1 := GetNodeValue(fee_info, "deduct1")
		deduct2 := GetNodeValue(fee_info, "deduct2")
		deduct3 := GetNodeValue(fee_info, "deduct3")
		if sum_amount != "" {
			calc_amount := ParseFloat(sum_amount)
			if deduct1 != "" {
				calc_amount = calc_amount - ParseFloat(deduct1)
			}
			if deduct2 != "" {
				calc_amount = calc_amount - ParseFloat(deduct2)
			}
			if deduct3 != "" {
				calc_amount = calc_amount - ParseFloat(deduct3)
			}
			fee_info = RegReplace(fee_info, `<calc_amount\s*\/>|<calc_amount\>.*?\<\/calc_amount>`, "<calc_amount>"+fmt.Sprintf("%.2f", calc_amount)+"</calc_amount>")
			xmlValue = strings.Replace(xmlValue, old_fee_info, fee_info, 1)
		}
	}

	surgery_info0s := RegMatchAll(xmlValue, `<surgery_info>[\s\S]*?<\/surgery_info>`)
	for _, surgery_info := range surgery_info0s {
		old_surgery_info := surgery_info
		surgery_records := RegMatchAll(surgery_info, `<surgery_record>[\s\S]*?<\/surgery_record>`)
		surgery_record_num := 0
		for _, surgery_record := range surgery_records {
			surgery_date := GetNodeValue(surgery_record, "surgery_date")
			surgery_code := GetNodeValue(surgery_record, "surgery_code")
			if surgery_date == "" && surgery_code == "" {
				surgery_record_num++
				if surgery_record_num != len(surgery_records) {
					surgery_info = strings.Replace(surgery_info, surgery_record, "", 1)
				}

			}
		}
		surgery_info = SetNodeValue(surgery_info, "seq", "")
		xmlValue = strings.Replace(xmlValue, old_surgery_info, surgery_info, 1)
	}
	//CSB0108RC0328000
	//"当fc264、fc265、fc266（字段唯一）三个字段同时存在录入值时，将fc264字段结果值与所有fc032字段结果值作对比，存在一致时
	//与fc032相同值的<clinic_no>的同一<clinic_info>节点下<surgery_record>节点下的各节点取对应字段的值为：
	//<surgery_date>取fc265的结果值，<surgery_code>取fc266的结果值
	//（该需求放在CSB0108RC0149000需求的前面）"
	//CSB0108RC0329000
	//"当fc264、fc265、fc266（字段唯一）三个字段同时存在录入值时，将fc264字段结果值与所有fc054字段结果值作对比，存在一致时，"
	//"与fc054相同值的<operation_no>的同一<inpatient_info>节点下<surgery_record>节点下的各节点取对应字段的值为：
	//<surgery_date>取fc265的结果值，<surgery_code>取fc266的结果值
	//该需求放在CSB0108RC0149000需求的前面）"
	_, fc264Input := GetOneField(fields, "fc264", false)
	_, fc264Value := GetOneField(fields, "fc264", true)
	_, fc265Input := GetOneField(fields, "fc265", false)
	_, fc265Val := GetOneField(fields, "fc265", true)
	_, fc266Input := GetOneField(fields, "fc266", false)
	_, fc0266Val := GetOneField(fields, "fc266", true)

	if fc264Input != "" && fc265Input != "" && fc266Input != "" {
		if isFc032 {
			clinicInfo := RegMatchAll(xmlValue, `<(clinic_info)>[\s\S]*?<\/(clinic_info)>`)
			for _, item := range clinicInfo {
				clinicNoValue := GetNodeValue(item, `clinic_no`)
				if clinicNoValue == fc264Value {
					oldItem := item
					item = SetNodeValue(item, "surgery_date", fc265Val)
					item = SetNodeValue(item, "surgery_code", fc0266Val)
					xmlValue = strings.Replace(xmlValue, oldItem, item, 1)
				}
			}
		}
		if isFc054 {
			clinicInfo := RegMatchAll(xmlValue, `<(inpatient_info)>[\s\S]*?<\/(inpatient_info)>`)
			for _, item := range clinicInfo {
				clinicNoValue := GetNodeValue(item, `operation_no`)
				if clinicNoValue == fc264Value {
					oldItem := item
					item = SetNodeValue(item, "surgery_date", fc265Val)
					item = SetNodeValue(item, "surgery_code", fc0266Val)
					xmlValue = strings.Replace(xmlValue, oldItem, item, 1)
				}
			}
		}
	}
	surgery_infos := RegMatchAll(xmlValue, `<surgery_info>[\s\S]*?<\/surgery_info>`)
	for _, surgery_info := range surgery_infos {
		old_surgery_info := surgery_info
		surgery_codes := RegMatchAll(surgery_info, `<surgery_code>.+?<\/surgery_code>`)
		record_counts := strconv.Itoa(len(surgery_codes))
		surgery_info = RegReplace(surgery_info, `<record_counts>.*?<\/record_counts>`, "<record_counts>"+record_counts+"</record_counts>")
		surgery_records := RegMatchAll(old_surgery_info, `<surgery_record>[\s\S]*?<\/surgery_record>`)
		for ss, surgery_record := range surgery_records {
			if ss >= len(surgery_codes) {
				break
			}
			old_surgery_record := surgery_record
			surgery_record = SetNodeValue(surgery_record, "seq", strconv.Itoa(ss+1))
			surgery_info = strings.Replace(surgery_info, old_surgery_record, surgery_record, 1)
		}
		xmlValue = strings.Replace(xmlValue, old_surgery_info, surgery_info, 1)
	}

	item_names := RegMatchAll(xmlValue, `<item_name>[\s\S]*?<\/item_name>`)
	for _, item_name := range item_names {
		value_it := GetNodeValue(item_name, "item_name")
		value_it = RegReplace(value_it, `\<|\>|（|）|\(|\)|\-|※|▪|•|☆|★|&lt;|&gt;`, "")
		new_item_name := SetNodeValue(item_name, "item_name", value_it)
		xmlValue = strings.Replace(xmlValue, item_name, new_item_name, 1)
	}

	// len := utf8.RuneCountInString(RegReplace(accident_id, `[^\,]`, ""))

	// isExit, fc032 := GetOneField(fields, "fc032", true)
	clinic_infos := RegMatchAll(xmlValue, `<clinic_info>[\s\S]*?<\/clinic_info>`)
	for _, clinic_info := range clinic_infos {
		old_clinic_info := clinic_info
		clinic_no := GetNodeValue(old_clinic_info, "clinic_no")
		// if isExit && fc032 == "" {
		// 	clinic_info = RegReplace(clinic_info, `<clinic_record>[\s\S]*?<\/clinic_record>`, "")
		// }
		if clinic_no == "" {
			clinic_info = RegReplace(clinic_info, `\>.*?\<`, "><")
			clinic_info = RegReplace(clinic_info, `<record_counts\s*\/>|<record_counts\>.*?\<\/record_counts>`, "<record_counts>0</record_counts>")
			clinic_info = RegReplace(clinic_info, `<fee_count\s*\/>|<fee_count\>.*?\<\/fee_count>`, "<fee_count>0</fee_count>")
			clinic_info = RegReplace(clinic_info, `<deduct1\s*\/>|<deduct1\>.*?\<\/deduct1>`, "<deduct1>0.00</deduct1>")
			clinic_info = RegReplace(clinic_info, `<deduct2\s*\/>|<deduct2\>.*?\<\/deduct2>`, "<deduct2>0.00</deduct2>")
			clinic_info = RegReplace(clinic_info, `<deduct3\s*\/>|<deduct3\>.*?\<\/deduct3>`, "<deduct3>0.00</deduct3>")
			xmlValue = strings.Replace(xmlValue, old_clinic_info, clinic_info, 1)
		}
	}

	// isExit, fc054 := GetOneField(fields, "fc054", true)
	inpatient_infos := RegMatchAll(xmlValue, `<inpatient_info>[\s\S]*?<\/inpatient_info>`)
	for _, inpatient_info := range inpatient_infos {
		old_inpatient_info := inpatient_info
		operation_no := GetNodeValue(old_inpatient_info, "operation_no")
		// if isExit && fc054 == "" {
		// 	inpatient_info = RegReplace(inpatient_info, `<inpatient_record>[\s\S]*?<\/inpatient_record>`, "")
		// }
		if operation_no == "" {
			inpatient_info = RegReplace(inpatient_info, `\>.*?\<`, "><")
			inpatient_info = RegReplace(inpatient_info, `<record_counts\s*\/>|<record_counts\>.*?\<\/record_counts>`, "<record_counts>0</record_counts>")
			inpatient_info = RegReplace(inpatient_info, `<fee_count\s*\/>|<fee_count\>.*?\<\/fee_count>`, "<fee_count>0</fee_count>")
			inpatient_info = RegReplace(inpatient_info, `<deduct1\s*\/>|<deduct1\>.*?\<\/deduct1>`, "<deduct1>0.00</deduct1>")
			inpatient_info = RegReplace(inpatient_info, `<deduct2\s*\/>|<deduct2\>.*?\<\/deduct2>`, "<deduct2>0.00</deduct2>")
			inpatient_info = RegReplace(inpatient_info, `<deduct3\s*\/>|<deduct3\>.*?\<\/deduct3>`, "<deduct3>0.00</deduct3>")
			xmlValue = strings.Replace(xmlValue, old_inpatient_info, inpatient_info, 1)
		}
	}

	if RegIsMatch(obj.Bill.Agency, `^(00183000|00183010|00183002|00183012|00183300|00183310|00183301|00183311)$`) {
		items := []string{"sum_amount", "calc_amount", "deduct1", "deduct3", "west_medicine", "china_medicine", "herbal_medicine", "examination", "inspection", "laboratory", "special_inspection", "treatment", "surgery", "other", "material", "nursing", "blood_transfusion", "bed", "supplementary", "refund", "self_cash_payment", "account_payment", "prepay", "fee_count"}
		for _, item := range items {
			query := `<` + item + `>.*<\/` + item + `>`
			xmlValue = RegReplace(xmlValue, query, "<"+item+" />")
		}
	}

	xmlValue = utils.RegReplace(xmlValue, `<clinic_no>88888.*<\/clinic_no>`, "<clinic_no></clinic_no>")
	xmlValue = utils.RegReplace(xmlValue, `<operation_no>99999.*<\/operation_no>`, "<operation_no></operation_no>")

	items := RegMatchAll(xmlValue, `<(clinic_info|inpatient_info)>[\s\S]*?<\/(clinic_info|inpatient_info)>`)
	for _, item := range items {
		old_item := item
		surgery_codes := RegMatchAll(item, `<(clinic_no|operation_no)>.+?<\/(clinic_no|operation_no)>`)
		record_counts := strconv.Itoa(len(surgery_codes))
		item = strings.Replace(item, `<record_counts></record_counts>`, "<record_counts>"+record_counts+"</record_counts>", 1)
		// item = RegReplace(item, `<record_counts>.*?<\/record_counts>`, "<record_counts>"+record_counts+"</record_counts>")
		xmlValue = strings.Replace(xmlValue, old_item, item, 1)
	}

	xmlValue = strings.Replace(xmlValue, `encoding="gbk"`, `encoding="gb2312"`, 1)
	xmlValue = RegReplace(xmlValue, `><\/.*>`, " />")
	xmlValue = RegReplace(xmlValue, `[\r]`, "")
	xmlValue = RegReplace(xmlValue, `\n\s*\n`, "\n")

	newXmlValue = xmlValue
	return err, newXmlValue
}

func GetAccidentTypeNum(fields []model.ProjectField, tempMap map[string]string) string {
	val := []string{}
	fmt.Println("------------------------tempMap-------------------", tempMap)
	for _, field := range fields {
		if RegIsMatch(field.Code, `^(fc018|fc112|fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$`) {
			fmt.Println("------------------------field-------------------", field.Code, field.FinalValue)
			if field.FinalValue != "" && HasKey(tempMap, field.FinalValue) {
				val = append(val, tempMap[field.FinalValue])
			}
		}
	}
	fmt.Println("------------------------GetAccidentTypeNum-------------------", val)
	if arrays.ContainsString(val, "1") != -1 && arrays.ContainsString(val, "2") == -1 && arrays.ContainsString(val, "3") == -1 {
		return "1"
	}
	if arrays.ContainsString(val, "1") == -1 && arrays.ContainsString(val, "2") != -1 && arrays.ContainsString(val, "3") == -1 {
		return "2"
	}
	return "3"
}

func operationTypeValue(fields []model.ProjectField) string {
	sum := ""
	for _, field := range fields {
		if RegIsMatch(field.Code, `^(fc177|fc178|fc179|fc180|fc181)$`) {
			if field.FinalValue != "" {
				if sum != "" {
					sum = sum + ","
				}
				sum = sum + field.FinalValue
			}
		}
	}
	return sum
}

func GetAccidentId(fields []model.ProjectField) string {
	val := ""
	for _, field := range fields {
		if RegIsMatch(field.Code, `^(fc036|fc104|fc058|fc161|fc162|fc205|fc206|fc207|fc208|fc209|fc239|fc240|fc241|fc242|fc243|fc244|fc245|fc246|fc247|fc248)$`) {
			if field.FinalValue != "" {
				val = val + field.FinalValue + ","
			}
		}
	}
	return RegReplace(val, `\,$`, "")
}

func GetAccidentNameDesc(fields []model.ProjectField) string {
	val := ""
	for _, field := range fields {
		if RegIsMatch(field.Code, `^(fc228|fc229|fc230|fc231|fc232|fc233|fc234|fc235|fc236|fc237|fc249|fc250|fc251|fc252|fc253|fc254|fc255|fc256|fc257|fc258)$`) {
			if field.FinalValue != "" {
				val = val + field.FinalValue + ","
			}
		}
	}
	return RegReplace(val, `\,$`, "")
}

func GetEarlyTimeOfFdAndid(xmlValue string) string {
	retVal := ""
	arrs := RegMatchAll(xmlValue, `<(first_date|in_date|confirm_date)>.*?<\/(first_date|in_date|confirm_date)>`)
	for _, arr := range arrs {
		date := GetNodeValue(arr, "first_date")
		if date == "" {
			date = GetNodeValue(arr, "in_date")
		}
		if date == "" {
			date = GetNodeValue(arr, "confirm_date")
		}
		if retVal == "" || (date != "" && retVal > date) {
			retVal = date
		}
	}
	return retVal
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
