package service

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	"server/module/report_management/model"
	"strings"
	"unicode"
)

/*
包 unicode 包含了一些针对测试字符的非常有用的函数（其中 ch 代表字符）：
	判断是否为字母： unicode.IsLetter(ch)
	判断是否为数字： unicode.IsDigit(ch)
	判断是否为空白符号： unicode.IsSpace(ch)
*/

// ReWrongSum 申诉通过/不通过, 修改产量
func ReWrongSum(wrong model.Wrong, proCode string, right string, WrongConfirm bool) (string, error) {
	//连接数据库
	db := global.ProDbMap[proCode]
	if db == nil {
		return "nil", global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("err", r)
			tx.Rollback()
		}
	}()
	if WrongConfirm {
		err := tx.Model(&model.Wrong{}).Where("id = ? ", wrong.ID).Updates(map[string]interface{}{
			"is_wrong_confirm": true,
			//"is_complain":      false,
			"is_operation_log": "1",
			"is_audit":         true,
		}).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	} else {
		err := tx.Model(&model.Wrong{}).Where("id = ? ", wrong.ID).Updates(map[string]interface{}{
			"is_wrong_confirm": false,
			"is_complain":      false,
			"is_operation_log": "2",
			"is_audit":         true,
		}).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	//更新产量明细
	var total int64
	var editDetail model.OutputStatistics
	err := tx.Model(&model.OutputStatistics{}).Where("code = ? AND to_char(submit_time,'YYYY-MM-DD') = ?", wrong.Code, wrong.SubmitDay.Format("2006-01-02")).Find(&editDetail).Count(&total).Error
	if err != nil {
		tx.Rollback()
		return "申诉失败", err
	}
	if total == 0 {
		tx.Rollback()
		return "申诉失败", errors.New("产量表没有对应的人员产量信息")
	}

	E := reflect.ValueOf(&editDetail).Elem()

	operation := make([]string, 0)

	if wrong.Op == "Op1" || wrong.Op == "Op2" {
		if wrong.IsBaoXiaoBlock == "1" {
			var t int64
			var wop []model.Wrong
			err = tx.Model(&model.Wrong{}).Where("code = ? AND op = ? AND to_char(submit_day,'YYYY-MM-DD') = ? AND is_wrong_confirm = 'false' AND is_bao_xiao_block = '1'", wrong.Code, wrong.Op, wrong.SubmitDay.Format("2006-01-02")).
				Find(&wop).Count(&t).Error
			if err != nil {
				tx.Rollback()
				return "", err
			}
			wrongSumOp := 0
			for _, v := range wop {
				wrongSumOp += Difference(v.Wrong, "")
			}

			E.FieldByName(wrong.Op + "ExpenseAccountFieldEffectiveCharacter").SetInt(int64(E.FieldByName(wrong.Op+"FieldCharacter").Interface().(int) - wrongSumOp))
			E.FieldByName(wrong.Op + "AccuracyRate").SetFloat(float64(E.FieldByName(wrong.Op+"FieldNum").Interface().(int)-int(t)) / float64(E.FieldByName(wrong.Op+"FieldNum").Interface().(int)))
			operation = append(operation, wrong.Op+"ExpenseAccountFieldEffectiveCharacter", wrong.Op+"AccuracyRate")
		} else if wrong.IsBaoXiaoBlock == "2" {
			var t int64
			var wop []model.Wrong
			err = tx.Model(&model.Wrong{}).Where("code = ? AND op = ? AND to_char(submit_day,'YYYY-MM-DD') = ? AND is_wrong_confirm = 'false' AND is_bao_xiao_block = '2'", wrong.Code, wrong.Op, wrong.SubmitDay.Format("2006-01-02")).
				Find(&wop).Count(&t).Error
			if err != nil {
				tx.Rollback()
				return "", err
			}
			wrongSumOp := 0
			for _, v := range wop {
				wrongSumOp += Difference(v.Wrong, "")
			}

			E.FieldByName(wrong.Op + "NotExpenseAccountFieldEffectiveCharacter").SetInt(int64(E.FieldByName(wrong.Op+"FieldCharacter").Interface().(int) - wrongSumOp))
			E.FieldByName(wrong.Op + "AccuracyRate").SetFloat(float64(E.FieldByName(wrong.Op+"FieldNum").Interface().(int)-int(t)) / float64(E.FieldByName(wrong.Op+"FieldNum").Interface().(int)))
			operation = append(operation, wrong.Op+"NotExpenseAccountFieldEffectiveCharacter", wrong.Op+"AccuracyRate")
		}

	} else {
		var t int64
		var wop []model.Wrong
		err = tx.Model(&model.Wrong{}).Where("code = ? AND op = ? AND to_char(submit_day,'YYYY-MM-DD') = ? AND is_wrong_confirm = 'false' ", wrong.Code, wrong.Op, wrong.SubmitDay.Format("2006-01-02")).
			Find(&wop).Count(&t).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
		wrongSumOp := 0
		for _, v := range wop {
			wrongSumOp += Difference(v.Wrong, "")
		}
		op := wrong.Op
		if wrong.Op == "Opq" {
			op = "OpQ"
		}
		E.FieldByName(op + "FieldEffectiveCharacter").SetInt(int64(E.FieldByName(op+"FieldCharacter").Interface().(int) - wrongSumOp))
		E.FieldByName(op + "AccuracyRate").SetFloat(float64(E.FieldByName(op+"FieldNum").Interface().(int)-int(t)) / float64(E.FieldByName(op+"FieldNum").Interface().(int)))
		operation = append(operation, op+"FieldEffectiveCharacter", op+"AccuracyRate")
	}
	var t int64
	var w []model.Wrong
	err = tx.Model(&model.Wrong{}).Where("code = ? AND to_char(submit_day,'YYYY-MM-DD') = ? AND is_wrong_confirm = 'false'", wrong.Code, wrong.SubmitDay.Format("2006-01-02")).
		Find(&w).Count(&t).Error
	if err != nil {
		tx.Rollback()
		return "", err
	}
	wrongSum := 0
	for _, v := range w {
		a := Difference(v.Wrong, "")
		wrongSum += a
	}
	E.FieldByName("SummaryFieldEffectiveCharacter").SetInt(int64(E.FieldByName("SummaryFieldCharacter").Interface().(int) - wrongSum))
	E.FieldByName("SummaryAccuracyRate").SetFloat(float64(E.FieldByName("SummaryFieldNum").Interface().(int)-int(t)) / float64(E.FieldByName("SummaryFieldNum").Interface().(int)))
	operation = append(operation, "SummaryFieldEffectiveCharacter", "SummaryAccuracyRate")
	err = tx.Select(operation).Model(&model.OutputStatistics{}).Where("code = ? AND to_char(submit_time,'YYYY-MM-DD') = ?", wrong.Code, wrong.SubmitDay.Format("2006-01-02")).Updates(editDetail).Error
	if err != nil {
		tx.Rollback()
		return "申诉失败", err
	}

	var Su model.OutputStatisticsSummary

	//tidy and save
	//tidy Mary ProCode SubmitTime NickName
	Su.Mary = editDetail.Op0FieldEffectiveCharacter + editDetail.Op1NotExpenseAccountFieldEffectiveCharacter + editDetail.Op1ExpenseAccountFieldEffectiveCharacter + editDetail.Op2NotExpenseAccountFieldEffectiveCharacter + editDetail.Op2ExpenseAccountFieldEffectiveCharacter + editDetail.OpQFieldEffectiveCharacter
	Su.Op0 = editDetail.Op0FieldEffectiveCharacter
	Su.Op0InvoiceNum = editDetail.Op0InvoiceNum
	Su.Op1ExpenseAccount = editDetail.Op1ExpenseAccountFieldEffectiveCharacter
	Su.Op1NotExpenseAccount = editDetail.Op1NotExpenseAccountFieldEffectiveCharacter
	Su.Op2ExpenseAccount = editDetail.Op2ExpenseAccountFieldEffectiveCharacter
	Su.Op2NotExpenseAccount = editDetail.Op2NotExpenseAccountFieldEffectiveCharacter
	Su.Question = editDetail.OpQFieldEffectiveCharacter
	Su.ProCode = proCode
	Su.SubmitTime = editDetail.SubmitTime
	Su.NickName = editDetail.NickName
	Su.Code = editDetail.Code

	//save
	var t1 int64
	//fmt.Println("a", "'%"+v.SubmitTime.Format("2006-01-02")+"%'")
	//StartTime, _ := time.ParseInLocation("2006-01-02", v.SubmitTime.Format("2006-01-02"), time.Local)
	err = global.GDb.Model(&model.OutputStatisticsSummary{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? AND pro_code = ? ", Su.SubmitTime.Format("2006-01-02"), Su.Code, Su.ProCode).
		Count(&t1).Error
	if err != nil {
		tx.Rollback()
		return "", err
	}
	if t1 != 0 {
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? AND pro_code = ? ", Su.SubmitTime.Format("2006-01-02"), Su.Code, Su.ProCode).
			Updates(Su).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	} else {
		if Su.Code == "" {
			return "", errors.New("统计报销单, 非报销单, 问题件: 没有对应的工号")
		}
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Create(&Su).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	//更新错误明细
	if WrongConfirm {
		err = tx.Model(&model.Wrong{}).Where("id = ? ", wrong.ID).Updates(map[string]interface{}{"is_wrong_confirm": WrongConfirm, "right": right}).Find(&wrong).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	} else {
		err = tx.Model(&model.Wrong{}).Where("id = ? ", wrong.ID).Updates(map[string]interface{}{"is_wrong_confirm": WrongConfirm}).Find(&wrong).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", err
	}

	return "申诉通过, 修改产量成功", nil
}

func Difference(input, result string) (wrong int) {
	if len(input) > len(result) {
		if result == "" {
			wrong = GetWrongSumVersionTwo(input, result, len(result))
		} else {
			wrong = GetWrongSumVersionTwo(input, result, len(result)) + GetWrongSumVersionTwo(input[len(result):], "", 0)
		}
	} else if len(input) == len(result) {
		wrong = GetWrongSumVersionTwo(input, result, len(input))
	} else if len(input) < len(result) {
		wrong = GetWrongSumVersionTwo(input, result, len(input)) + GetWrongSumVersionTwo("", result[len(input):], 0)
	}
	return wrong
}

func GetWrongSumVersionTwo(input, result string, length int) (wrong int) {
	if input == "" {
		return A(result)
	}
	w := 0
	//匹配中英文?/？
	//reg2 := regexp.MustCompile("^[?|？]$")
	if result != "" {
		if strings.Index(input, "?") != -1 || strings.Index(input, "？") != -1 {
			return A(result)
		} else {
			for i := 0; i < length; i++ {
				if input[i] != result[i] {
					w += A(string(input[i]))
				}
			}
		}
	} else {
		return A(input)
	}
	return w
}

func A(str string) int {
	fieldCharacter := 0
	reg2 := regexp.MustCompile("^[?|？]$")
	for _, rr := range str {
		if reg2.MatchString(string(rr)) {
			continue
		}
		if unicode.Is(unicode.Han, rr) {
			fieldCharacter = fieldCharacter + 2
		} else {
			fieldCharacter = fieldCharacter + 1
		}
	}
	return fieldCharacter
}
