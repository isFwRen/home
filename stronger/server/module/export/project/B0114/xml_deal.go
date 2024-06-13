/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022年12月16日10:15:57
 */

package B0114

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
	global.GLog.Info("------------------B0114:::XmlDeal-----------------------")
	obj := o.(FormatObj)
	fields := obj.Fields

	c, ok := global.GProConf[obj.Bill.ProCode].ConstTable["B0114_华夏理赔_医院名称表"]
	yiyuan_jibie := make(map[string]string, 0)
	if ok {
		fmt.Println("------------------------c-----------------------:", c[0])
		for _, arr := range c {
			yiyuan_jibie[strings.TrimSpace(arr[1])] = arr[0]
		}
	}

	claimInCode := utils.GetNodeValue(obj.Bill.OtherInfo, "claimInCode")
	xmlValue = SetNodeValue(xmlValue, "claimInCode", claimInCode)

	accidentDate := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentDate")
	if accidentDate == "" {
		xmlValue = SetNodeValue(xmlValue, "accidentDate", accidentDate)
	}

	if GetNodeValue(xmlValue, "accidentAreaOption") == "" {
		accidentAreaOption := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentAreaOption")
		xmlValue = SetNodeValue(xmlValue, "accidentAreaOption", accidentAreaOption)
	}

	if GetNodeValue(xmlValue, "accidentAddr") == "" {
		accidentAddr := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentAddr")
		xmlValue = SetNodeValue(xmlValue, "accidentAddr", accidentAddr)
	}

	if GetNodeValue(xmlValue, "accidentDesc") == "" {
		accidentDesc := utils.GetNodeValue(obj.Bill.OtherInfo, "accidentDesc")
		xmlValue = SetNodeValue(xmlValue, "accidentDesc", accidentDesc)
	}

	causeCodes := utils.GetNodeData(obj.Bill.OtherInfo, "causeCode")
	fmt.Println("------------------------causeCodes-----------------------:", causeCodes)
	if utils.HasItem(causeCodes, "03") || utils.HasItem(causeCodes, "04") {
		_, fc018 := GetOneField(fields, "fc018", true)
		xmlValue = SetNodeValue(xmlValue, "deathDate", fc018)
	}

	fc105s := GetSameField(fields, "fc105", false)
	cureType := GetCureType(fc105s)
	xmlValue = SetNodeValue(xmlValue, "cureType", cureType)

	diseaseDescription := GetDiseaseDescription(fields)
	xmlValue = SetNodeValue(xmlValue, "diseaseDescription", diseaseDescription)

	fc106s := GetSameField(fields, "fc106", true)
	fc106_res := []string{}
	target := ""
	for _, fc106 := range fc106s {
		if fc106 != "" && !utils.HasItem(fc106_res, fc106) {
			fc106_res = append(fc106_res, fc106)
		}
	}
	if len(fc106_res) == 0 {
		target += "<hospitalInfo>\n\t\t<hospitalCode></hospitalCode>\n\t\t<hospitalGrade></hospitalGrade>\n\t</hospitalInfo>\n\t"
	} else {
		for _, fc106 := range fc106_res {
			hospitalGrade, _ := yiyuan_jibie[fc106]
			target += "<hospitalInfo>\n\t\t<hospitalCode>" + fc106 + "</hospitalCode>\n\t\t<hospitalGrade>" + hospitalGrade + "</hospitalGrade>\n\t</hospitalInfo>\n\t"
		}
	}
	hospitalInfos := RegMatchAll(xmlValue, `<hospitalInfo>[\s\S]*?<\/hospitalInfo>`)
	xmlValue = strings.Replace(xmlValue, hospitalInfos[0], target, 1)

	hospitalInfos_arr := RegMatchAll(xmlValue, `<hospitalInfos>[\s\S]*?<\/hospitalInfos>`)
	hospitalInfos_value := hospitalInfos_arr[0]
	hospitalInfos_clone := hospitalInfos_arr[0]
	hospitalInfos_value = SetNodeValue(hospitalInfos_value, "count", strconv.Itoa(len(fc106_res)))
	xmlValue = strings.Replace(xmlValue, hospitalInfos_clone, hospitalInfos_value, 1)
	fmt.Println("-----------1111111111--------------")
	_, fc096 := GetOneField(fields, "fc096", true)
	_, fc187 := GetOneField(fields, "fc187", true)
	_, fc190 := GetOneField(fields, "fc190", true)
	surgeryDtos := RegMatchAll(xmlValue, `<surgeryDto>[\s\S]*?<\/surgeryDto>`)
	fmt.Println("-----------surgeryDtos--------------", strconv.Itoa(len(surgeryDtos)))
	surgeryDto_count := 0
	if RegIsMatch(fc096, `^(A|)$`) {
		if RegIsMatch(fc187, `^(A|)$`) && RegIsMatch(fc190, `^(A|)$`) {
			// surgeryDto_value := surgeryDtos[0]
			surgeryDto_clone := surgeryDtos[0]
			// surgeryDto_value = SetNodeValue(surgeryDto_value, "surgeryCode", "")
			xmlValue = strings.Replace(xmlValue, surgeryDto_clone, "", 1)
		}
	} else {
		surgeryDto_count++
	}
	if RegIsMatch(fc187, `^(A|)$`) {
		xmlValue = strings.Replace(xmlValue, surgeryDtos[1], "", 1)
	} else {
		surgeryDto_count++
	}
	if RegIsMatch(fc190, `^(A|)$`) {
		xmlValue = strings.Replace(xmlValue, surgeryDtos[2], "", 1)
	} else {
		surgeryDto_count++
	}

	surgeryDtoses := RegMatchAll(xmlValue, `<surgeryDtos>[\s\S]*?<\/surgeryDtos>`)
	fmt.Println("-----------surgeryDtoses--------------", strconv.Itoa(len(surgeryDtoses)))
	surgeryDtos_value := surgeryDtoses[0]
	surgeryDtos_clone := surgeryDtoses[0]
	surgeryDtos_value = SetNodeValue(surgeryDtos_value, "count", strconv.Itoa(surgeryDto_count))
	xmlValue = strings.Replace(xmlValue, surgeryDtos_clone, surgeryDtos_value, 1)
	fmt.Println("----------AAAAAAAAAAAAAAAAAAAA-------------")
	_, fc101 := GetOneField(fields, "fc101", true)
	_, fc181 := GetOneField(fields, "fc181", true)
	_, fc182 := GetOneField(fields, "fc182", true)
	injureDiagnosisDtos := RegMatchAll(xmlValue, `<injureDiagnosisDto>[\s\S]*?<\/injureDiagnosisDto>`)
	injureDiagnosisDto_count := 0
	if RegIsMatch(fc101, `^(A|)$`) {
		if RegIsMatch(fc181, `^(A|)$`) && RegIsMatch(fc182, `^(A|)$`) {
			injureDiagnosisDto_value := injureDiagnosisDtos[0]
			injureDiagnosisDto_clone := injureDiagnosisDtos[0]
			injureDiagnosisDto_value = SetNodeValue(injureDiagnosisDto_value, "diagnosisModeCode", "")
			xmlValue = strings.Replace(xmlValue, injureDiagnosisDto_clone, injureDiagnosisDto_value, 1)
		}
	} else {
		injureDiagnosisDto_count++
	}
	if RegIsMatch(fc181, `^(A|)$`) {
		xmlValue = strings.Replace(xmlValue, injureDiagnosisDtos[1], "", 1)
	} else {
		injureDiagnosisDto_count++
	}
	if RegIsMatch(fc182, `^(A|)$`) {
		xmlValue = strings.Replace(xmlValue, injureDiagnosisDtos[2], "", 1)
	} else {
		injureDiagnosisDto_count++
	}
	fmt.Println("----------2222222222222222-------------")
	injureDiagnosisDtoses := RegMatchAll(xmlValue, `<injureDiagnosisDtos>[\s\S]*?<\/injureDiagnosisDtos>`)
	injureDiagnosisDtos_value := injureDiagnosisDtoses[0]
	injureDiagnosisDtos_clone := injureDiagnosisDtoses[0]
	injureDiagnosisDtos_value = SetNodeValue(injureDiagnosisDtos_value, "count", strconv.Itoa(injureDiagnosisDto_count))
	xmlValue = strings.Replace(xmlValue, injureDiagnosisDtos_clone, injureDiagnosisDtos_value, 1)

	fc179s := GetSameField(fields, "fc179", false)
	fmt.Println("----------fc179s----------", fc179s)
	medicalBillInfos_count := 0
	for _, fc179 := range fc179s {
		if fc179 == "3" {
			medicalBillInfos_count++
		}
	}
	medicalBillInfos := RegMatchAll(xmlValue, `<medicalBillInfos>[\s\S]*?<\/medicalBillInfos>`)
	medicalBillInfos_value := medicalBillInfos[0]
	medicalBillInfos_clone := medicalBillInfos[0]
	medicalBillInfos_value = SetNodeValue(medicalBillInfos_value, "count", strconv.Itoa(medicalBillInfos_count))
	xmlValue = strings.Replace(xmlValue, medicalBillInfos_clone, medicalBillInfos_value, 1)

	clearItems := [][]string{
		{"09", "10", "seriDiseases", "fc016", "seriDiseaseName", "count", "seriDiseaseConfirmDate", "seriDiseaseDeadDate", "fc018"},
		{"17", "18", "sickDiseases", "fc019", "sickDiseaseName", "count", "sickDiseaseConfirmDate"},
		{"27", "28", "moderateDiseases", "fc021", "moderateDiseaseName", "count", "moderateDiseaseConfirmDate"},
		{"25", "26", "sepcDiseases", "fc023", "sepcDiseaseName", "count", "sepcDiseaseConfirmDate"},
		{"05", "06", "allDiformitis", "fc025", "allDiformityCode", "count", "allDiformityPaymentPerc", "allDiformityConfirmDate", "fc026"},
		{"07", "08", "disableDiformitis", "fc029", "disableDiformityCode", "count", "disableDiformityConfirmDate", "disableDiformityPaymentPerc"},
	}
	for _, clearItem := range clearItems {
		seriDiseases := RegMatchAll(xmlValue, `<`+clearItem[2]+`>[\s\S]*?<\/`+clearItem[2]+`>`)
		seriDiseases_value := seriDiseases[0]
		seriDiseases_clone := seriDiseases[0]
		isFc016, _ := GetOneField(fields, clearItem[3], false)
		if isFc016 {
			seriDiseases_value = SetNodeValue(seriDiseases_value, "count", "1")
		} else {
			seriDiseases_value = SetNodeValue(seriDiseases_value, "count", "0")
		}
		if !utils.HasItem(causeCodes, clearItem[0]) && !utils.HasItem(causeCodes, clearItem[1]) {
			seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[5], "")
			seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[4], "")
			seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[6], "")
			if len(clearItem) >= 8 {
				seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[7], "")
			}
		} else if len(clearItem) == 9 {
			_, fc018 := GetOneField(fields, clearItem[8], true)
			if clearItem[8] == "fc026" {
				if fc018 != "1.00" {
					seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[5], "")
					seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[4], "")
					seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[6], "")
					if len(clearItem) >= 8 {
						seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[7], "")
					}
				}
			} else {
				seriDiseases_value = SetNodeValue(seriDiseases_value, clearItem[7], fc018)
			}

		}
		xmlValue = strings.Replace(xmlValue, seriDiseases_clone, seriDiseases_value, 1)
	}

	fmt.Println("---------333333333333333------------")
	if utils.HasItem(causeCodes, "21") || utils.HasItem(causeCodes, "22") {
		fc112s := GetSameField(fields, "fc112", true)
		fc111s := GetSameField(fields, "fc111", true)
		hospitalAllowance_sum := 0.0
		eraday := "3000/01/01"
		if len(fc112s) > 0 && len(fc112s) == len(fc111s) {
			for ff, fc112 := range fc112s {
				fc111 := fc111s[ff]
				if fc111s[ff] != "" {
					a, _ := time.Parse("2006/01/02", fc111)
					b, _ := time.Parse("2006/01/02", eraday)
					if a.Before(b) {
						eraday = fc111
					}
				}
				if fc112 != "" && fc111s[ff] != "" {
					a, _ := time.Parse("2006/01/02", fc111s[ff])
					b, _ := time.Parse("2006/01/02", fc112)
					d := b.Sub(a)
					hospitalAllowance_sum += d.Hours() / 24
				}
			}
		}
		hospitalAllowances := RegMatchAll(xmlValue, `<hospitalAllowance>[\s\S]*?<\/hospitalAllowance>`)
		for _, hospitalAllowance := range hospitalAllowances {
			hospitalAllowance_value := hospitalAllowance
			hospitalAllowance_clone := hospitalAllowance
			hospitalAllowance_value = SetNodeValue(hospitalAllowance_value, "deductDays", "0")
			hospitalAllowance_value = SetNodeValue(hospitalAllowance_value, "hospitalDays", ToString(hospitalAllowance_sum))
			hospitalAllowance_value = SetNodeValue(hospitalAllowance_value, "approvedDays", ToString(hospitalAllowance_sum))
			xmlValue = strings.Replace(xmlValue, hospitalAllowance_clone, hospitalAllowance_value, 1)
		}

		fc110_1 := []string{}
		fc110_2 := []string{}
		fc110_A := []string{}
		for _, MedicalBillInfo := range obj.MedicalBillInfo {
			_, fc048 := GetOneField(MedicalBillInfo.Fields, "fc048", false)
			_, fc110 := GetOneField(MedicalBillInfo.Fields, "fc110", true)
			if fc110 != "" {
				a, _ := time.Parse("2006/01/02", fc110)
				b, _ := time.Parse("2006/01/02", eraday)
				if fc048 == "1" {
					if !utils.HasItem(fc110_A, fc110) {
						fc110_A = append(fc110_A, fc110)
					}
					if (a.Before(b) || a.Equal(b)) && !utils.HasItem(fc110_1, fc110) {
						fc110_1 = append(fc110_1, fc110)
					}
				} else if fc048 == "2" {
					if (a.Before(b) || a.Equal(b)) && !utils.HasItem(fc110_2, fc110) {
						fc110_2 = append(fc110_2, fc110)
					}
				}
			}
		}

		prehospitalOutpatientAllowances := RegMatchAll(xmlValue, `<prehospitalOutpatientAllowances>[\s\S]*?<\/prehospitalOutpatientAllowances>`)
		prehospitalOutpatientAllowances_value := prehospitalOutpatientAllowances[0]
		prehospitalOutpatientAllowances_clone := prehospitalOutpatientAllowances[0]
		prehospitalOutpatientAllowances_value = SetNodeValue(prehospitalOutpatientAllowances_value, "count", strconv.Itoa(len(fc110_2)))
		xmlValue = strings.Replace(xmlValue, prehospitalOutpatientAllowances_clone, prehospitalOutpatientAllowances_value, 1)
		prehospitalOutpatientAllowance := RegMatchAll(prehospitalOutpatientAllowances_value, ` *<prehospitalOutpatientAllowance>[\s\S]*?<\/prehospitalOutpatientAllowance>`)
		reData := ""
		for kk, data := range fc110_2 {
			prehospitalOutpatientAllowance_value := prehospitalOutpatientAllowance[0]
			prehospitalOutpatientAllowance_value = SetNodeValue(prehospitalOutpatientAllowance_value, "isDeductDate", "0")
			prehospitalOutpatientAllowance_value = SetNodeValue(prehospitalOutpatientAllowance_value, "isApproveDate", "1")
			prehospitalOutpatientAllowance_value = SetNodeValue(prehospitalOutpatientAllowance_value, "prehospitalOutpatientDate", data)
			if kk > 0 {
				reData += "\n"
			}
			reData += prehospitalOutpatientAllowance_value
		}
		xmlValue = strings.Replace(xmlValue, prehospitalOutpatientAllowance[0], reData, 1)
		fmt.Println("----------4444444444444444-----------")
		prehospitalEmergencyAllowances := RegMatchAll(xmlValue, `<prehospitalEmergencyAllowances>[\s\S]*?<\/prehospitalEmergencyAllowances>`)
		prehospitalEmergencyAllowances_value := prehospitalEmergencyAllowances[0]
		prehospitalEmergencyAllowances_clone := prehospitalEmergencyAllowances[0]
		prehospitalEmergencyAllowances_value = SetNodeValue(prehospitalEmergencyAllowances_value, "count", strconv.Itoa(len(fc110_1)))
		xmlValue = strings.Replace(xmlValue, prehospitalEmergencyAllowances_clone, prehospitalEmergencyAllowances_value, 1)
		prehospitalEmergencyAllowance := RegMatchAll(prehospitalEmergencyAllowances_value, ` *<prehospitalEmergencyAllowance>[\s\S]*?<\/prehospitalEmergencyAllowance>`)
		reData = ""
		for kk, data := range fc110_1 {
			prehospitalEmergencyAllowance_value := prehospitalEmergencyAllowance[0]
			prehospitalEmergencyAllowance_value = SetNodeValue(prehospitalEmergencyAllowance_value, "isDeductDate", "0")
			prehospitalEmergencyAllowance_value = SetNodeValue(prehospitalEmergencyAllowance_value, "isApproveDate", "1")
			prehospitalEmergencyAllowance_value = SetNodeValue(prehospitalEmergencyAllowance_value, "prehospitalEmergencyDate", data)
			if kk > 0 {
				reData += "\n"
			}
			reData += prehospitalEmergencyAllowance_value
		}
		xmlValue = strings.Replace(xmlValue, prehospitalEmergencyAllowance[0], reData, 1)

		emergencyAllowances := RegMatchAll(xmlValue, `<emergencyAllowances>[\s\S]*?<\/emergencyAllowances>`)
		emergencyAllowances_value := emergencyAllowances[0]
		emergencyAllowances_clone := emergencyAllowances[0]
		emergencyAllowances_value = SetNodeValue(emergencyAllowances_value, "count", strconv.Itoa(len(fc110_1)))
		xmlValue = strings.Replace(xmlValue, emergencyAllowances_clone, emergencyAllowances_value, 1)
		emergencyAllowance := RegMatchAll(emergencyAllowances_value, ` *<emergencyAllowance>[\s\S]*?<\/emergencyAllowance>`)
		reData = ""
		for kk, data := range fc110_A {
			emergencyAllowance_value := emergencyAllowance[0]
			emergencyAllowance_value = SetNodeValue(emergencyAllowance_value, "isDeductDate", "0")
			emergencyAllowance_value = SetNodeValue(emergencyAllowance_value, "isApproveDate", "1")
			emergencyAllowance_value = SetNodeValue(emergencyAllowance_value, "emergencyDate", data)
			if kk > 0 {
				reData += "\n"
			}
			reData += emergencyAllowance_value
		}
		xmlValue = strings.Replace(xmlValue, emergencyAllowance[0], reData, 1)

	}
	fmt.Println("---------555555555555555555555555----------")
	isAllowance := utils.GetNodeValue(obj.Bill.OtherInfo, "isAllowance")
	fmt.Println("---------isAllowance----------", isAllowance)
	if isAllowance == "Y" {
		sum := 0.0
		sumFileds := []string{"fc207", "fc208", "fc209", "fc210", "fc211", "fc212", "fc213", "fc214", "fc061"}
		for _, field := range fields {
			if utils.HasItem(sumFileds, field.Code) && field.FinalValue != "" {
				fmt.Println("---------sumFileds----------", field.Code, field.FinalValue)
				sum += ParseFloat(field.FinalValue)
			}
		}
		fmt.Println("---------sum----------", sum)
		wardAllowances := RegMatchAll(xmlValue, `<wardAllowance>[\s\S]*?<\/wardAllowance>`)
		for _, wardAllowance := range wardAllowances {
			wardAllowance_value := wardAllowance
			wardAllowance_clone := wardAllowance
			wardAllowance_value = SetNodeValue(wardAllowance_value, "deductDays", "0")
			wardAllowance_value = SetNodeValue(wardAllowance_value, "hospitalDays", ToString(sum))
			wardAllowance_value = SetNodeValue(wardAllowance_value, "approvedDays", ToString(sum))
			xmlValue = strings.Replace(xmlValue, wardAllowance_clone, wardAllowance_value, 1)
		}
	}

	fc044s := GetSameField(fields, "fc044", false)
	approvedDays_sum := 0.0
	for _, fc044 := range fc044s {
		if fc044 != "" {
			approvedDays_sum += ParseFloat(fc044)
		}
	}

	ambulanceAllowance_arr := RegMatchAll(xmlValue, `<ambulanceAllowance>[\s\S]*?<\/ambulanceAllowance>`)
	ambulanceAllowance_value := ambulanceAllowance_arr[0]
	ambulanceAllowance_clone := ambulanceAllowance_arr[0]
	approvedDays := ""
	if approvedDays_sum != 0.0 {
		approvedDays = ToString(approvedDays_sum)
	}
	ambulanceAllowance_value = SetNodeValue(ambulanceAllowance_value, "approvedDays", approvedDays)
	xmlValue = strings.Replace(xmlValue, ambulanceAllowance_clone, ambulanceAllowance_value, 1)

	zeroFields := []string{"hospitalRealAmount", "newAgricultureAmount", "selfPayOne", "selfPayTwo", "selfPayAmount", "insuredPayAmount", "fundPayAmount", "wholeFundAmount", "helpFundHospitalAmount", "wholeFundYearAmount", "helpFundVisitAmount", "helpFundVisitYearAmount", "sociSecuAbateAmount", "ncmsAbateAmount", "commInsuAbateAmount", "otherAbateAmount", "billAmount", "abateAmount", "ratifyAmount"}
	for _, zeroField := range zeroFields {
		xmlValue = RegReplace(xmlValue, `<`+zeroField+`><\/`+zeroField+`>`, "<"+zeroField+">0.00</"+zeroField+">")
	}

	issueMessages := RegMatchAll(xmlValue, `<issueMessage>[\s\S]*?<\/issueMessage>`)
	errorData := ""
	errMeses := ""
	questionCount := 0
	fc175 := true
	fc138 := true
	fc231 := true
	for _, field := range fields {
		// if len(filed.Issues ) > 0 {
		if fc175 && field.Code == "fc175" && len(field.Issues) > 0 {
			questionCount++
			fc175 = false
			errorData = errorData + "<issueMessage>\n\t\t<issueCode>01</issueCode>\n\t\t<issueSource>珠海汇流</issueSource>\n\t\t<issueDescription>账单号-清单填写模糊</issueDescription>\n\t</issueMessage>\n\t"
		}
		if fc138 && RegIsMatch(field.Code, `^(fc138|fc142|fc146|fc150|fc154|fc158|fc162|fc166)$`) && len(field.Issues) > 0 {
			fc138 = false
			questionCount++
			errorData = errorData + "<issueMessage>\n\t\t<issueCode>01</issueCode>\n\t\t<issueSource>珠海汇流</issueSource>\n\t\t<issueDescription>项目名称-清单填写模糊</issueDescription>\n\t</issueMessage>\n\t"
		}
		if fc231 && RegIsMatch(field.Code, `^(fc231|fc232|fc233|fc234|fc235|fc236|fc237|fc238)$`) && len(field.Issues) > 0 {
			fc231 = false
			questionCount++
			errorData = errorData + "<issueMessage>\n\t\t<issueCode>01</issueCode>\n\t\t<issueSource>珠海汇流</issueSource>\n\t\t<issueDescription>费用-清单填写模糊</issueDescription>\n\t</issueMessage>\n\t"
		}
		for _, issue := range field.Issues {
			errMes := issue.Code + "_" + issue.Message + ";"
			if strings.Index(errMeses, errMes) == -1 {
				if !RegIsMatch(field.Code, `^(fc175|fc138|fc142|fc146|fc150|fc154|fc158|fc162|fc166|fc231|fc232|fc233|fc234|fc235|fc236|fc237|fc238)$`) && issue.Message != "" {
					questionCount++
					errorData = errorData + "<issueMessage>\n\t\t<issueCode>" + issue.Code + "</issueCode>\n\t\t<issueSource>珠海汇流</issueSource>\n\t\t<issueDescription>" + issue.Message + "</issueDescription>\n\t</issueMessage>\n\t"
				}
				errMeses += errMes
			}

		}
	}
	fmt.Println("------------------------error-------------------", errorData)
	if utils.HasItem(causeCodes, "21") || utils.HasItem(causeCodes, "22") {
		emergencyAllowances := RegMatchAll(xmlValue, `<emergencyAllowances>[\s\S]*?<\/emergencyAllowances>`)
		isDeductDate := RegMatchAll(emergencyAllowances[0], `<isDeductDate>.+?<\/isDeductDate>`)
		emergencyDate := RegMatchAll(emergencyAllowances[0], `<emergencyDate><\/emergencyDate>`)
		if len(isDeductDate) > 0 && len(emergencyDate) == 0 {
			questionCount++
			errorData = errorData + "<issueMessage>\n\t\t<issueCode>04</issueCode>\n\t\t<issueSource>珠海汇流</issueSource>\n\t\t<issueDescription>急诊日期未填写</issueDescription>\n\t</issueMessage>\n\t"
		}
	}

	if errorData == "" {
		errorData = "<issueMessage>\n\t\t<issueCode></issueCode>\n\t\t<issueSource></issueSource>\n\t\t<issueDescription></issueDescription>\n\t</issueMessage>\n\t"
	}
	xmlValue = strings.Replace(xmlValue, issueMessages[0], errorData, 1)

	issueMessageses := RegMatchAll(xmlValue, `<issueMessages>[\s\S]*?<\/issueMessages>`)[0]
	issueMessageses_value := issueMessageses
	issueMessageses = SetNodeValue(issueMessageses, "count", strconv.Itoa(questionCount))
	xmlValue = strings.Replace(xmlValue, issueMessageses_value, issueMessageses, 1)
	// fmt.Println("----------6666666666666666666666----------")
	if medicalBillInfos_count > 0 {
		medicalDefuctInfo := RegMatchAll(xmlValue, ` *<medicalDefuctInfo>[\s\S]*?<\/medicalDefuctInfo>`)[0]
		drugDetailInfo := RegMatchAll(medicalDefuctInfo, ` *<drugDetailInfo>[\s\S]*?<\/drugDetailInfo>`)[0]
		drugDetailInfoFields := [][]string{
			{"fc215", "fc138", "fc223", "fc231", "fc239", "fc247", "fc139", "fc140", "fc141"},
			{"fc216", "fc142", "fc224", "fc232", "fc240", "fc248", "fc143", "fc144", "fc145"},
			{"fc217", "fc146", "fc225", "fc233", "fc241", "fc249", "fc147", "fc148", "fc149"},
			{"fc218", "fc150", "fc226", "fc234", "fc242", "fc250", "fc151", "fc152", "fc153"},
			{"fc219", "fc154", "fc227", "fc235", "fc243", "fc251", "fc155", "fc156", "fc157"},
			{"fc220", "fc158", "fc228", "fc236", "fc244", "fc252", "fc159", "fc160", "fc161"},
			{"fc221", "fc162", "fc229", "fc237", "fc245", "fc253", "fc163", "fc164", "fc165"},
			{"fc222", "fc166", "fc230", "fc238", "fc246", "fc254", "fc167", "fc168", "fc169"},
		}
		cacheMedicalDefuctInfo := ""
		for _, MedicalBillInfo := range obj.MedicalBillInfo {
			medicalDefuctInfo_value := medicalDefuctInfo
			medicalDefuctInfo_count := 0
			medicalDefuctInfo_value = SetNodeValue(medicalDefuctInfo_value, "billNo", MedicalBillInfo.Code)
			cacheDrugDetailInfo := ""
			for _, QingDan := range MedicalBillInfo.QingDan {
				for _, drugDetailInfoField := range drugDetailInfoFields {
					drugDetailInfo_value := drugDetailInfo
					_, fc140 := GetOneField(QingDan.Fields, drugDetailInfoField[7], true)
					if fc140 == "1" {
						medicalDefuctInfo_count++
						_, fc215 := GetOneField(QingDan.Fields, drugDetailInfoField[0], true)
						_, fc138 := GetOneField(QingDan.Fields, drugDetailInfoField[1], true)
						_, fc231 := GetOneField(QingDan.Fields, drugDetailInfoField[3], true)
						_, fc247 := GetOneField(QingDan.Fields, drugDetailInfoField[5], true)
						_, fc139 := GetOneField(QingDan.Fields, drugDetailInfoField[6], true)
						_, fc141 := GetOneField(QingDan.Fields, drugDetailInfoField[8], true)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "medicalType", fc215)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "drugName", fc138)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "drugUnitPrice", fc231)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "drugAmount", fc247)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "drugSumAmount", fc139)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "selfpaidRatio", fc140)
						drugDetailInfo_value = SetNodeValue(drugDetailInfo_value, "selfpaidAmount", fc141)
						if cacheDrugDetailInfo != "" {
							cacheDrugDetailInfo += "\n"
						}
						cacheDrugDetailInfo += drugDetailInfo_value
					}
				}

			}
			if medicalDefuctInfo_count == 0 {
				continue
			}
			medicalDefuctInfo_value = SetNodeValue(medicalDefuctInfo_value, "count", strconv.Itoa(medicalDefuctInfo_count))

			medicalDefuctInfo_value = strings.Replace(medicalDefuctInfo_value, drugDetailInfo, cacheDrugDetailInfo, 1)
			if cacheMedicalDefuctInfo != "" {
				cacheMedicalDefuctInfo += "\n"
			}
			cacheMedicalDefuctInfo += medicalDefuctInfo_value
		}
		xmlValue = strings.Replace(xmlValue, medicalDefuctInfo, cacheMedicalDefuctInfo, 1)

	}

	medicalBillDetails := RegMatchAll(xmlValue, `<medicalBillDetails>[\s\S]*?<\/medicalBillDetails>`)
	for _, medicalBillDetail := range medicalBillDetails {
		medicalBillDetails_value := medicalBillDetail
		medicalBillDetails_clone := medicalBillDetail
		medicalBillDetails_count := 0
		medicalBillDetail := RegMatchAll(medicalBillDetails_value, `<medicalBillDetail>[\s\S]*?<\/medicalBillDetail>`)
		for _, data := range medicalBillDetail {
			medicalChargeType := GetNodeValue(data, "medicalChargeType")
			medicalBillAmount := GetNodeValue(data, "medicalBillAmount")
			if (medicalBillAmount == "" || medicalBillAmount == "0.00") && (medicalChargeType == "" || medicalChargeType == "0.00") {
				medicalBillDetails_value = strings.Replace(medicalBillDetails_value, data, "", 1)
			} else {
				medicalBillDetails_count++
			}
		}
		medicalBillDetails_value = SetNodeValue(medicalBillDetails_value, "count", strconv.Itoa(medicalBillDetails_count))
		xmlValue = strings.Replace(xmlValue, medicalBillDetails_clone, medicalBillDetails_value, 1)
	}

	if utils.HasItem(causeCodes, "21") || utils.HasItem(causeCodes, "22") {
		xmlValue = SetNodeValue(xmlValue, "operationDeductAmount", "0.00")
		operationMedicalAmount := 0.0
		_, fc045 := GetOneField(fields, "fc045", true)
		if fc045 != "" {
			operationMedicalAmount = ParseFloat(fc045)
		}
		operationMedicalAmountFields := [][]string{
			{"fc062", "fc063"}, {"fc064", "fc065"}, {"fc066", "fc067"}, {"fc068", "fc069"}, {"fc070", "fc071"}, {"fc072", "fc073"}, {"fc074", "fc075"}, {"fc076", "fc077"}, {"fc078", "fc079"}, {"fc080", "fc081"}, {"fc082", "fc083"}, {"fc084", "fc085"}, {"fc086", "fc087"}, {"fc088", "fc089"}, {"fc130", "fc131"},
		}
		for _, MedicalBillInfo := range obj.MedicalBillInfo {
			for _, operationMedicalAmountField := range operationMedicalAmountFields {
				_, fc062 := GetOneField(MedicalBillInfo.Fields, operationMedicalAmountField[0], false)
				if strings.Index(fc062, "10") != -1 {
					_, fc063 := GetOneField(MedicalBillInfo.Fields, operationMedicalAmountField[1], false)
					if fc063 != "" {
						operationMedicalAmount += ParseFloat(fc063)
					}
				}
			}
		}
		xmlValue = SetNodeValue(xmlValue, "operationMedicalAmount", ToString(operationMedicalAmount))
		xmlValue = SetNodeValue(xmlValue, "operationApproveAmount", ToString(operationMedicalAmount))
	}

	//20230609新增 编码：CSB0114RC0125000  medicalDefuctInfos大节点下，存在多少个medicalDefuctInfo大节点，则billQuantity节点的值为多少；当不存在medicalDefuctInfo大节点时，billQuantity默认为0
	medicalInfos := RegMatchAll(xmlValue, `<medicalDefuctInfos>[\s\S]*?<\/medicalDefuctInfos>`)
	medicalInfoList := RegMatchAll(medicalInfos[0], `<medicalDefuctInfo>[\s\S]*?<\/medicalDefuctInfo>`)
	billQuantity := RegMatchAll(medicalInfos[0], `<billQuantity>[\s\S]*?<\/billQuantity>`)
	medicalInfoLength := len(medicalInfoList)
	medicalInfoCount := strconv.Itoa(medicalInfoLength)
	if len(billQuantity) > 0 {
		xmlValue = SetNodeValue(xmlValue, "billQuantity", medicalInfoCount)
	}

	if medicalBillInfos_count == 0 {
		medicalBillInfos := RegMatchAll(xmlValue, `<medicalDefuctInfos>[\s\S]*?<\/medicalDefuctInfos>`)
		for _, medicalBillInfo := range medicalBillInfos {
			medicalBillInfo_clone := medicalBillInfo
			medicalBillInfo = SetNodeValue(medicalBillInfo, "count", "")
			medicalBillInfo = SetNodeValue(medicalBillInfo, "defuctionAmount", "")
			medicalBillInfo = SetNodeValue(medicalBillInfo, "billQuantity", "")
			medicalDefuctInfos := RegMatchAll(medicalBillInfo_clone, `<medicalDefuctInfo>[\s\S]*?<\/medicalDefuctInfo>`)
			for _, medicalDefuctInfo := range medicalDefuctInfos {
				medicalDefuctInfo_clone := medicalDefuctInfo
				clearDefuctInfoItems := []string{"defuctionDesc", "medicalType", "drugName", "drugSpecifications", "drugUnitPrice", "drugUnit", "drugAmount", "drugSumAmount", "selfpaidRatio", "selfpaidAmount"}
				for _, item := range clearDefuctInfoItems {
					medicalDefuctInfo = SetNodeValue(medicalDefuctInfo, item, "")
				}
				medicalBillInfo = strings.Replace(medicalBillInfo, medicalDefuctInfo_clone, medicalDefuctInfo, 1)
			}
			xmlValue = strings.Replace(xmlValue, medicalBillInfo_clone, medicalBillInfo, 1)

		}
	}

	return err, xmlValue
}

func GetDiseaseDescription(fields []model.ProjectField) string {
	_, fc101 := GetOneField(fields, "fc101", false)
	_, fc181 := GetOneField(fields, "fc181", false)
	_, fc182 := GetOneField(fields, "fc182", false)
	value := ""
	if !RegIsMatch(fc101, `^(A|)$`) && strings.Index(fc101, "?") == -1 {
		// _, fc101 = GetOneField(fields, "fc101", true)
		value = fc101
	}
	if !RegIsMatch(fc181, `^(A|)$`) && strings.Index(fc181, "?") == -1 {
		// _, fc181 = GetOneField(fields, "fc181", true)
		if value != "" {
			value = value + "、"
		}
		value += fc181
	}
	if !RegIsMatch(fc182, `^(A|)$`) && strings.Index(fc182, "?") == -1 {
		// _, fc182 = GetOneField(fields, "fc182", true)
		if value != "" {
			value = value + "、"
		}
		value += fc182
	}
	return value
}

func GetCureType(fc105s []string) string {
	num1 := 0
	num2 := 0
	// num3 := 0
	for _, fc105 := range fc105s {
		if fc105 == "1" {
			num1++
		}
		if fc105 == "2" {
			num2++
		}
	}
	if len(fc105s) == num1 {
		return "01"
	} else if len(fc105s) == num2 {
		return "02"
	} else if num1 != 0 && num2 != 0 {
		return "03"
	} else if num1 == 0 && num2 == 0 {
		return "04"
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
