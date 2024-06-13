/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/30 3:14 下午
 */

package export

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	model4 "server/module/export/model"
	service2 "server/module/export/service"
	model2 "server/module/load/model"
	"server/module/pro_manager/model"
	model3 "server/module/report_management/model"
	modelBase "server/module/sys_base/model"
	"server/utils"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/wxnacy/wgo/arrays"
	"gorm.io/gorm"
)

// WrongSum 错误统计
func WrongSum(reqParam model.ProCodeAndId, InvoiceLen int, id, flag string, fields []model2.ProjectField) (err error, row int64, tx *gorm.DB) {
	global.GLog.Info("错误统计")
	fmt.Println("错误统计")

	err, obj := getBillObj(reqParam, id, flag, fields)
	if err != nil {
		return err, 0, nil
	}
	var blockMap = make(map[string]model2.ProjectBlock, 0)
	for _, block := range obj.ProjectBlockList {
		blockMap[block.ID] = block
	}
	return WrongSumData(obj, blockMap, InvoiceLen)
}

// WrongSum 错误统计
func WrongSumData(obj model.BillObj, blockMap map[string]model2.ProjectBlock, InvoiceLen int) (err error, row int64, tx *gorm.DB) {

	dayCodeMap := []string{}
	sumDayCodeMap := []string{}

	ops := []string{"Op1", "Op2", "Opq", "Op0"}
	var wrongs []model3.Wrong
	for _, field := range obj.ProjectFieldList {

		f := reflect.ValueOf(field)

		for _, op := range ops {
			inputStatus := fmt.Sprintf("%v", f.FieldByName(op+"Input").Interface())
			block, ok := blockMap[field.BlockID]
			if has, e := regexp.MatchString("yes", inputStatus); e == nil && has && ok {
				b := reflect.ValueOf(block)
				//key = b.FieldByName(op+"Code").Interface().(string) + "|" + b.FieldByName(op+"SubmitAt").Interface().(time.Time).Format("2006-01-02")
				if b.FieldByName(op+"Code").Interface().(string) == "0" {
					continue
				}

				day := b.FieldByName(op + "SubmitAt").Interface().(time.Time).Format("2006-01-02")
				dayCode := day + "_" + b.FieldByName(op+"Code").Interface().(string)
				if arrays.Contains(sumDayCodeMap, dayCode) == -1 {
					sumDayCodeMap = append(sumDayCodeMap, dayCode)
				}

				if f.FieldByName(op+"Value").Interface() != f.FieldByName("ResultValue").Interface() && f.FieldByName(op+"Input").Interface() == "yes" {
					name, okCode := global.UserCodeName[b.FieldByName(op+"Code").Interface().(string)]
					if !okCode {
						name = ""
					}

					// day := b.FieldByName(op + "SubmitAt").Interface().(time.Time).Format("2006-01-02")
					// dayCode := day + "_" + b.FieldByName(op+"Code").Interface().(string)
					if arrays.Contains(dayCodeMap, dayCode) == -1 {
						dayCodeMap = append(dayCodeMap, dayCode)
					}

					// dayCodeMap[day] = b.FieldByName(op + "Code").Interface().(string)
					var wrong = model3.Wrong{}
					if op == "Op1" || op == "Op2" {
						if strings.Index(block.Name, "报销") != -1 {
							wrong.IsBaoXiaoBlock = "1"
						} else {
							wrong.IsBaoXiaoBlock = "2"
						}
					} else {
						wrong.IsBaoXiaoBlock = "0"
					}
					wrong.FieldId = field.ID
					wrong.SubmitDay = b.FieldByName(op + "SubmitAt").Interface().(time.Time)
					wrong.Code = b.FieldByName(op + "Code").Interface().(string)
					wrong.NickName = name
					wrong.Op = op
					wrong.BillName = obj.ProjectBill.BillName
					wrong.BillNum = obj.ProjectBill.BillNum
					wrong.Agency = obj.ProjectBill.Agency
					//wrong.Type = "模板类型"
					//wrong.Types = "单子类型,比如混合型, 医疗等"
					wrong.Block = block.Code
					wrong.FieldCode = field.Code
					wrong.FieldName = field.Name
					wrong.Path = obj.ProjectBill.DownloadPath
					wrong.Picture = []string{block.Picture}
					wrong.Wrong = f.FieldByName(op + "Value").Interface().(string)
					wrong.Right = field.ResultValue
					wrong.IsComplain = false
					wrong.IsWrongConfirm = false
					wrong.IsOcr = inputStatus == "ocr"
					wrong.BillId = obj.ProjectBill.ID
					wrong.IsOperationLog = ""
					wrong.IsAudit = false
					wrongs = append(wrongs, wrong)
				}
			}
		}
	}
	// if len(wrongs) > 0 {
	err, row = service2.InsertWrongs(obj.ProjectBill.ProCode, wrongs, obj.ProjectBill.ID)
	if err != nil {
		return err, 0, nil
	}

	// err, row = OcrSum(obj)
	// if err != nil {
	// 	return err, 0, nil
	// }
	// }

	// err, InvoiceInformation := findInvoice(reqParam.ProCode, id, flag)
	// if err != nil {
	// 	return err, 0, nil
	// }
	// fmt.Println("InvoiceInformation", InvoiceInformation)
	//同一个分块同一个工序只能是一个人录入
	// if flag == "2" {
	// 	err, tx = CalculateOutput(wrongCharsMap, questionMark, key, InvoiceInformation, InvoiceLen, obj)
	// 	return err, row, tx
	// }

	// if flag == "1" {
	// 	err, tx = ReExportBill(wrongs, reqParam.ProCode, InvoiceInformation, InvoiceLen)
	// 	return err, row, tx
	// }

	// if err != nil {
	// 	return err, 0, nil
	// }
	// fmt.Println("-------------dayCodeMap---------------", dayCodeMap)
	err, tx = EffectiveCharacter(dayCodeMap, obj.ProjectBill.ProCode)

	if InvoiceLen > 0 {
		day := ""
		code := ""
		db := global.ProDbMap[obj.ProjectBill.ProCode]
		var Output model3.OutputStatistics
		for _, block := range obj.ProjectBlockList {
			if block.Name == "未定义" && block.Code == "bc001" {
				code = block.Op0Code
				day = block.Op0SubmitAt.Format("2006-01-02")
			}
		}
		err = db.Model(&model3.OutputStatistics{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? ", day, code).Find(&Output).Error
		if err != nil {
			tx.Rollback()
			return err, 0, nil
		}
		Output.Op0InvoiceNum += InvoiceLen
		if err = tx.Select("Op0InvoiceNum").Where("id = ? ", Output.ID).Updates(&Output).Error; err != nil {
			fmt.Println("-------------errerrerr---------------", err)
			tx.Rollback()
			return err, 0, nil
		}

	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	err = ExpenseAccountSum(sumDayCodeMap, obj.ProjectBill.ProCode)
	if err != nil {
		return err, 0, nil
	}

	return err, row, tx
}

func OcrSum(blocks []model2.ProjectBlock, obj model4.ResultDataBill) (error, int64) {
	ocrStatisticses := []model3.OcrStatistics{}
	billID := obj.Bill.ID
	billNum := obj.Bill.BillNum
	proCode := obj.Bill.ProCode
	blockPicMap := map[string]string{}
	if !utils.RegIsMatch(`^(B0114|B0108)$`, proCode) {
		return nil, 0
	}
	//Op1CodeArr := []string{}
	for _, block := range blocks {

		if block.Picture != "" {
			blockPicMap[block.ID] = obj.Bill.DownloadPath + block.Picture
		}
	}
	for _, invoice := range obj.Invoice {
		fc001 := "是"
		for qq, fields := range invoice.QingDan {
			for _, field := range fields {
				if field.Op0Input != "ocr" && field.Op0Input != "ocr_no_if" {
					continue
				}
				if field.Op0Value == "" {
					continue
				}
				// if field.Code == "fc001" {
				// 	if field.FinalValue == "1" {
				// 		fc001 = "是"
				// 	} else {
				// 		fc001 = "否"
				// 	}
				// }
				// if !utils.RegIsMatch(`^项目名称\d*-清单$`, field.Name) {
				// 	continue
				// }

				ocrStatistics := model3.OcrStatistics{}
				ocrStatistics.BillID = billID
				ocrStatistics.BillNum = billNum
				ocrStatistics.Code = field.Code
				ocrStatistics.Name = "清单" + strconv.Itoa(qq+1) + "-" + field.Name
				ocrStatistics.OcrType = fc001
				ocrStatistics.ResultValue = field.ResultValue
				ocrStatistics.SubmitTime = obj.Bill.ScanAt
				ocrStatistics.ProCode = proCode

				if field.Op0Input == "ocr" || field.Op0Input == "ocr_no_if" {
					ocrStatistics.Value = field.Op0Value
					ocrStatistics.Disable = "1"
					ocrStatistics.Compare = "1"
					if field.Op0Value != field.ResultValue {
						ocrStatistics.Compare = "2"
					}
					if field.Op0Input == "ocr" {
						ocrStatistics.Disable = "2"
					}
					min := utf8.RuneCountInString(field.Op0Value)
					max := utf8.RuneCountInString(field.Op0Value)
					wnum := 0
					if min < utf8.RuneCountInString(field.ResultValue) {
						max = utf8.RuneCountInString(field.ResultValue)
					} else {
						min = utf8.RuneCountInString(field.ResultValue)
					}
					for i := 0; i < min; i++ {
						if Substr(field.Op0Value, i, i+1) == Substr(field.ResultValue, i, i+1) {
							wnum++
						}
					}
					// fmt.Println("------------OcrSum--------------", field.Op0Value, field.ResultValue, min, max, wnum)
					// wnum = max - wnum
					if max == 0 {
						ocrStatistics.Rate = 0
					} else {
						ocrStatistics.Rate = Decimal(float64(wnum)/float64(max)) * 100
					}

				} else if field.ResultValue == "" {
					continue
				}

				ocrStatistics.Pic = blockPicMap[field.BlockID]
				split := strings.Split(ocrStatistics.Pic, "/")
				fmt.Println("----------split-----", len(split))
				lastName := split[len(split)-1:]
				for _, block := range blocks {
					if block.Picture == lastName[0] {
						ocrStatistics.JobNumber = block.Op1Code
					}
				}
				tx := global.GDb.Model(&modelBase.SysUser{})
				var user modelBase.SysUser
				tx.Where("code = ?", ocrStatistics.JobNumber).Find(&user)
				ocrStatistics.NickName = user.NickName
				ocrStatisticses = append(ocrStatisticses, ocrStatistics)
			}
		}
	}

	// for _, field := range obj.ProjectFieldList {

	// }
	return service2.InsertOcrs(proCode, ocrStatisticses, billID)
}

func Difference(input, result string) (wrong int) {
	if utf8.RuneCountInString(input) > utf8.RuneCountInString(result) {
		// if result == "" {
		length := utf8.RuneCountInString(result)
		wrong = GetWrongSumCharacter(Substr(input, 0, length), result)
		out_char := Substr(input, length, -1)
		wrong += GetWrongSumCharacter(out_char, "")

		// } else {
		// wrong += GetWrongSumCharacter(input[len(result):], "", 0)
		// }
	} else if utf8.RuneCountInString(input) == utf8.RuneCountInString(result) {
		wrong = GetWrongSumCharacter(input, result)
	} else if utf8.RuneCountInString(input) < utf8.RuneCountInString(result) {
		length := utf8.RuneCountInString(input)
		wrong = GetWrongSumCharacter(input, Substr(result, 0, length))
		out_char := Substr(result, length, -1)
		wrong += GetWrongSumCharacter(out_char, "")
	}
	return wrong
}

// ExpenseAccountSum 统计报销单, 非报销单, 问题件
func ExpenseAccountSum(dayCodeMap []string, proCode string) (err error) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	for _, dayCode := range dayCodeMap {
		day_code := strings.Split(dayCode, "_")
		day := day_code[0]
		code := day_code[1]
		var Outputs []model3.OutputStatistics
		err = db.Model(&model3.OutputStatistics{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? ", day, code).Find(&Outputs).Error
		if err != nil {
			return err
		}
		for _, Output := range Outputs {
			var Su model3.OutputStatisticsSummary

			//tidy and save
			//tidy Mary ProCode SubmitTime NickName
			Su.Mary = Output.Op0FieldEffectiveCharacter + Output.Op1NotExpenseAccountFieldEffectiveCharacter + Output.Op1ExpenseAccountFieldEffectiveCharacter + Output.Op2NotExpenseAccountFieldEffectiveCharacter + Output.Op2ExpenseAccountFieldEffectiveCharacter + Output.OpQFieldEffectiveCharacter
			Su.Op0 = Output.Op0FieldEffectiveCharacter
			Su.Op0InvoiceNum = Output.Op0InvoiceNum
			Su.Op1ExpenseAccount = Output.Op1ExpenseAccountFieldEffectiveCharacter
			Su.Op1NotExpenseAccount = Output.Op1NotExpenseAccountFieldEffectiveCharacter
			Su.Op2ExpenseAccount = Output.Op2ExpenseAccountFieldEffectiveCharacter
			Su.Op2NotExpenseAccount = Output.Op2NotExpenseAccountFieldEffectiveCharacter
			Su.Question = Output.OpQFieldEffectiveCharacter
			Su.ProCode = proCode
			Su.SubmitTime = Output.SubmitTime
			Su.NickName = Output.NickName
			Su.Code = Output.Code
			// fmt.Println("---------------Output---NickName-------------------------------: ", Output.Code, Su.NickName, Su.NickName == "")
			if Su.NickName == "" {
				var user modelBase.SysUser
				err = global.GDb.Model(&modelBase.SysUser{}).Where("code = ? ", Su.Code).Find(&user).Error
				// fmt.Println("--------------err------------------------: ", err, user.NickName)
				if err == nil {
					Su.NickName = user.NickName
				}
			}

			// fmt.Println("---------------Output---NickName-------------------------------: ", Output.Code, Output.NickName)
			// fmt.Println("-------------Su-----NickName-------------------------------: ", Su.NickName)
			//save
			var total int64
			//fmt.Println("a", "'%"+v.SubmitTime.Format("2006-01-02")+"%'")
			//StartTime, _ := time.ParseInLocation("2006-01-02", v.SubmitTime.Format("2006-01-02"), time.Local)
			err = global.GDb.Model(&model3.OutputStatisticsSummary{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? AND pro_code = ? ", Su.SubmitTime.Format("2006-01-02"), Su.Code, Su.ProCode).
				Count(&total).Error
			if err != nil {
				return err
			}
			if total != 0 {
				err = global.GDb.Model(&model3.OutputStatisticsSummary{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? AND pro_code = ? ", Su.SubmitTime.Format("2006-01-02"), Su.Code, Su.ProCode).
					Updates(Su).Error
				if err != nil {
					return err
				}
			} else {
				if Su.Code == "" {
					return errors.New("统计报销单, 非报销单, 问题件: 没有对应的工号")
				}
				err = global.GDb.Model(&model3.OutputStatisticsSummary{}).Create(&Su).Error
				if err != nil {
					return err
				}
			}

		}

	}

	return nil
}

// func DifferenceVersionTwo(input, result string) int {
// 	wrongSum := 0
// 	reg2 := regexp.MustCompile("^[?]$")
// 	for _, v := range input {
// 		if !reg2.MatchString(string(v)) {
// 			if unicode.Is(unicode.Han, v) {
// 				wrongSum = wrongSum + 2
// 			} else {
// 				wrongSum = wrongSum + 1
// 			}
// 		}
// 	}
// 	return wrongSum
// }

func Substr(str string, start, end int) string {
	data := []rune(str)
	if start == -1 && end != -1 {
		return string(data[:end])
	} else if start != -1 && end == -1 {
		return string(data[start:])
	} else if start != -1 && end != -1 {
		return string(data[start:end])
	}
	return str

}

func GetWrongSumCharacter(data1, data2 string) (wrong int) {
	fieldCharacter := 0
	reg2 := regexp.MustCompile("^[?|？]$")
	// fmt.Println("------------GetWrongSumCharacter---------------------", data1, data2)
	// for i := 0; i < length; i++ {
	i := -1
	for _, val := range data1 {
		i++
		// fmt.Println("------------val---------------------", i, fieldCharacter, string(val))
		if reg2.MatchString(Substr(data1, i, i+1)) || (data2 != "" && reg2.MatchString(Substr(data2, i, i+1))) {
			continue
		}
		if data2 != "" && Substr(data1, i, i+1) == Substr(data2, i, i+1) {
			continue
		}
		// fmt.Println("------------data1---------------------", i, fieldCharacter, data1[i], data2[i], val, unicode.Is(unicode.Han, rune(data1[i])))
		if unicode.Is(unicode.Han, val) {
			fieldCharacter = fieldCharacter + 2
		} else {
			fieldCharacter = fieldCharacter + 1
		}
	}
	return fieldCharacter

	// if input == "" {
	// 	return A(result)
	// }
	// w := 0
	// //匹配中英文?/？
	// //reg2 := regexp.MustCompile("[?|？]")
	// if result != "" {
	// 	if strings.Index(input, "?") != -1 || strings.Index(input, "？") != -1 {
	// 		return AHasQuestMa(result)
	// 	} else {
	// 		for i := 0; i < length; i++ {
	// 			if input[i] != result[i] {
	// 				w += A(string(input[i]))
	// 			}
	// 		}
	// 	}
	// } else {
	// 	return A(input)
	// }
	// return w
}

// func A(str string) int {
// 	fieldCharacter := 0
// 	reg2 := regexp.MustCompile("^[?|？]$")
// 	for _, rr := range str {
// 		if reg2.MatchString(string(rr)) {
// 			continue
// 		}
// 		if unicode.Is(unicode.Han, rr) {
// 			fieldCharacter = fieldCharacter + 2
// 		} else {
// 			fieldCharacter = fieldCharacter + 1
// 		}
// 	}
// 	return fieldCharacter
// }

// func AHasQuestMa(str string) int {
// 	fieldCharacter := 0
// 	reg2 := regexp.MustCompile("^[?|？]$")
// 	for _, rr := range str {
// 		if reg2.MatchString(string(rr)) {
// 			continue
// 		}
// 		if unicode.Is(unicode.Han, rr) {
// 			fieldCharacter = fieldCharacter + 2
// 		} else {
// 			fieldCharacter = fieldCharacter + 1
// 		}
// 	}
// 	return fieldCharacter
// }

func EffectiveCharacter(dayCodeMap []string, proCode string) (error, *gorm.DB) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, nil
	}
	// fmt.Println("wrongCharsMap length ", len(wrongCharsMap))
	tx := db.Begin()
	for _, dayCode := range dayCodeMap {
		day_code := strings.Split(dayCode, "_")
		day := day_code[0]
		code := day_code[1]
		var Output model3.OutputStatistics
		var wrongs []model3.Wrong
		err := tx.Model(&model3.Wrong{}).Where("code = ? AND to_char(submit_day,'YYYY-MM-DD') = ? AND is_wrong_confirm = 'false'", code, day).Find(&wrongs).Error
		if err != nil {
			// tx.Rollback()
			return err, nil
		}
		opWrongNum := map[string]int{
			"Op0":     0,
			"Op1":     0,
			"Op2":     0,
			"OpQ":     0,
			"Summary": 0,
		}
		opWrongCharacter := map[string]int{
			"Op0":                  0,
			"Op1":                  0,
			"Op2":                  0,
			"OpQ":                  0,
			"Summary":              0,
			"Op1NotExpenseAccount": 0,
			"Op1ExpenseAccount":    0,
			"Op2NotExpenseAccount": 0,
			"Op2ExpenseAccount":    0,
		}
		// fmt.Println("------------day_code---------------------", day, code, len(wrongs))
		for _, wrong := range wrongs {
			w_op := wrong.Op
			w_op = strings.Replace(w_op, "q", "Q", 1)
			opWrongNum[w_op] += 1
			opWrongNum["Summary"] += 1
			wrongSum := Difference(wrong.Wrong, wrong.Right)
			// fmt.Println("------------Difference---------------------", wrong.Wrong, wrong.Right, wrongSum)
			opWrongCharacter[w_op] += wrongSum
			opWrongCharacter["Summary"] += wrongSum
			if w_op == "Op1" {
				if wrong.IsBaoXiaoBlock == "1" {
					opWrongCharacter["Op1ExpenseAccount"] += wrongSum
				} else if wrong.IsBaoXiaoBlock == "2" {
					opWrongCharacter["Op1NotExpenseAccount"] += wrongSum
				}
			} else if w_op == "Op2" {
				if wrong.IsBaoXiaoBlock == "1" {
					opWrongCharacter["Op2ExpenseAccount"] += wrongSum
				} else if wrong.IsBaoXiaoBlock == "2" {
					opWrongCharacter["Op2NotExpenseAccount"] += wrongSum
				}
			}
			// moutput := reflect.ValueOf(&Output).Elem()
		}
		err = tx.Model(&model3.OutputStatistics{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? ", day, code).Find(&Output).Error
		if err != nil {
			// tx.Rollback()
			return err, nil
		}
		// fmt.Println("-------------opWrongCharacter---------------", opWrongCharacter)
		OperationOp := make([]string, 0)
		for key, value := range opWrongCharacter {
			moutput := reflect.ValueOf(&Output).Elem()
			fieldCharacter := moutput.FieldByName(key + "FieldCharacter").Interface().(int)
			effectiveCharacter := fieldCharacter - value
			if effectiveCharacter < 0 {
				effectiveCharacter = 0
			}
			moutput.FieldByName(key + "FieldEffectiveCharacter").Set(reflect.ValueOf(effectiveCharacter))
			OperationOp = append(OperationOp, key+"FieldEffectiveCharacter")
			if utils.RegIsMatch(`^(Op0|Op1|Op2|OpQ|Summary)$`, key) {
				fieldNum := moutput.FieldByName(key + "FieldNum").Interface().(int)
				if fieldNum > 0 {
					effectiveFieldNum := fieldNum - opWrongNum[key]
					accuracyRate := Decimal(float64(effectiveFieldNum) / float64(fieldNum))
					moutput.FieldByName(key + "AccuracyRate").Set(reflect.ValueOf(accuracyRate))
					OperationOp = append(OperationOp, key+"AccuracyRate")
				}

			}

		}
		// fmt.Println("-------------OperationOp---------------", OperationOp)
		if err = tx.Select(OperationOp).Where("id = ? ", Output.ID).Updates(&Output).Error; err != nil {
			fmt.Println("-------------errerrerr---------------", err)
			tx.Rollback()
			return err, tx
		}
		// fmt.Println("-------------errerrerr---------------", err)
		// var Outputs []model3.OutputStatistics
	}
	return nil, tx
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func getBillObj(reqParam model.ProCodeAndId, id, flag string, fields []model2.ProjectField) (err error, obj model.BillObj) {
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr, obj
	}
	if flag == "1" {
		err = db.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).
			First(&obj.ProjectBill).Error
		err = db.Model(&model2.ProjectBlock{}).Where("bill_id = ?", reqParam.ID).
			Find(&obj.ProjectBlockList).Error
		//err = db.Model(&model2.ProjectField{}).Where("bill_id = ?", reqParam.ID).
		//	Find(&obj.ProjectFieldList).Error
		obj.ProjectFieldList = append(obj.ProjectFieldList, fields...)

	} else if flag == "2" {
		dbs := global.ProDbMap[reqParam.ProCode+"_task"]
		if dbs == nil {
			return global.ProDbErr, obj
		}

		var b model.ProjectBill
		err = db.Model(&model.ProjectBill{}).Where("id = ?", id).Find(&b).Error
		if err != nil {
			return err, obj
		}
		var a model.ProjectBill
		err = dbs.Model(&model.ProjectBill{}).Where("id = ?", id).Find(&a).Error
		if err != nil {
			return err, obj
		}
		if b.Stage == 2 && a.Stage == 6 {
			err = dbs.Model(&model.ProjectBill{}).Where("id = ?", reqParam.ID).
				First(&obj.ProjectBill).Error
			err = dbs.Model(&model2.ProjectBlock{}).Where("bill_id = ?", reqParam.ID).
				Find(&obj.ProjectBlockList).Error
			err = dbs.Model(&model2.ProjectField{}).Where("bill_id = ?", reqParam.ID).
				Find(&obj.ProjectFieldList).Error
		} else {
			err = errors.New("单据状态有误")
		}
	}
	return err, obj
}

func AbnormalityExport(proCode, billId, flag string) error {
	//1. find memory bill
	err, obj, blockMap := getAbnormalityExportBillObj(proCode, billId)
	if err != nil {
		return err
	}
	switch flag {
	case "1":
		err, _, _ = WrongSumData(obj, blockMap, 0)
		// err = forceExportBill(obj, blockMap, proCode, billId)
	case "2":
		err = delBillInTheEntry(obj, blockMap)
	case "3":
		err = exportAbnormalBill(obj, blockMap)
	}
	return err
}

func delBillInTheEntry(obj model.BillObj, blockMap map[string]model2.ProjectBlock) error {
	return nil
}

func exportAbnormalBill(obj model.BillObj, blockMap map[string]model2.ProjectBlock) error {
	return nil
}

func getAbnormalityExportBillObj(proCode, billId string) (err error, obj model.BillObj, s map[string]model2.ProjectBlock) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, obj, nil
	}
	dbs := global.ProDbMap[proCode+"_task"]
	if dbs == nil {
		return global.ProDbErr, obj, nil
	}

	var b model.ProjectBill
	err = db.Model(&model.ProjectBill{}).Where("id = ?", billId).Find(&b).Error
	if err != nil {
		return err, obj, nil
	}
	var a model.ProjectBill
	err = dbs.Model(&model.ProjectBill{}).Where("id = ?", billId).Find(&a).Error
	if err != nil {
		return err, obj, nil
	}
	if b.Stage == 2 && a.Stage == 2 {
		err = dbs.Model(&model.ProjectBill{}).Where("id = ?", billId).
			First(&obj.ProjectBill).Error
		err = dbs.Model(&model2.ProjectBlock{}).Where("bill_id = ?", billId).
			Find(&obj.ProjectBlockList).Error
		err = dbs.Model(&model2.ProjectField{}).Where("bill_id = ?", billId).
			Find(&obj.ProjectFieldList).Error
	} else {
		err = errors.New("单据状态有误")
	}
	var blockMap = make(map[string]model2.ProjectBlock, 0)
	for _, block := range obj.ProjectBlockList {
		blockMap[block.ID] = block
	}
	return nil, obj, blockMap
}
