package utils

import (
	"reflect"
	"regexp"
	"server/global"
	"server/module/report_management/model"
)

func GetWrongSum(data1, data2 string, length int) int {
	wrong := 0
	//匹配中文字符： [u4e00-u9fa5]
	reg1 := regexp.MustCompile("^[\u4e00-\u9fa5]$")
	//匹配中英文?/？
	reg2 := regexp.MustCompile("^[?|？]$")
	if data2 != "" {
		for i := 0; i < length; i++ {
			if data1[i] != data2[i] && !reg2.MatchString(string(data1[i])) && !reg2.MatchString(string(data2[i])) && reg1.MatchString(string(data1[i])) {
				wrong += 1
			} else {
				wrong += 2
			}
		}
	} else {
		for i := 0; i < length; i++ {
			if !reg2.MatchString(string(data1[i])) && reg1.MatchString(string(data1[i])) {
				wrong += 1
			} else {
				wrong += 2
			}
		}
	}
	return wrong
}

func GetOpWrong(doc model.Wrong, op string, wrongChars map[string]int) {
	if doc.Op == op {
		wrong := 0
		data1 := doc.Wrong
		data2 := doc.Right
		if len(data1) > len(data2) {
			wrong = GetWrongSum(data1, data2, len(data2)) + GetWrongSum(data1[len(data2):], "", len(data1[len(data2):]))
		} else {
			wrong = GetWrongSum(data1, data2, len(data1))
		}

		if _, ok := wrongChars[op]; !ok {
			wrongChars[op] = wrong
		} else {
			wrongChars[op] += wrong
		}
	}
}

func ReWrongSum(W model.Wrong, proName, proCode string) error {
	//匹配中英文?/？
	//reg2 := regexp.MustCompile("^[?|？]$")
	tempOp := W.Op
	wrongChars := make(map[string]int, 0)
	GetOpWrong(W, tempOp, wrongChars)
	//连接数据库
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	var statistics model.OutputStatistics
	err := db.Model(&model.OutputStatistics{}).Where("to_char(submit_time,'YYYY-MM-DD') = ? AND code = ? ", W.SubmitDay.Format("2006-01-02"), W.Code).Find(&statistics).Error
	if err != nil {
		return err
	}
	//statisticsTy := reflect.TypeOf(statistics)
	statisticsVa := reflect.ValueOf(&statistics)
	//有效字符总量
	statisticsVa.FieldByName(W.Op + "FieldEffectiveCharacter").SetInt(int64(statisticsVa.FieldByName(W.Op+"FieldCharacter").Interface().(int) - wrongChars[W.Op]))
	statisticsVa.FieldByName(W.Op + "AccuracyRate").SetFloat(float64(statisticsVa.FieldByName(W.Op+"FieldEffectiveCharacter").Interface().(int) / statisticsVa.FieldByName(W.Op+"FieldCharacter").Interface().(int)))
	err = db.Model(&model.OutputStatistics{}).Where("submit_time = ? AND code = ? ", W.SubmitDay, W.Code).Updates(statistics).Error
	if err != nil {
		return err
	}
	return nil
}
