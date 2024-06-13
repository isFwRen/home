package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/module/practice/model"
	model2 "server/module/pro_manager/model"
	trainStage "server/module/training_guide/service"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/wxnacy/wgo/arrays"
	"gorm.io/gorm"
)

// func SaveSubmitData(proCode string, taskSubmit request.TaskSubmit) error {
// 	op := taskSubmit.Op
// 	block := taskSubmit.Block
// 	mblock := reflect.ValueOf(&block).Elem()
// 	mOp := strings.Replace(op, "o", "O", -1)

// 	subTime := mblock.FieldByName(mOp + "SubmitAt").Interface().(time.Time)
// 	isNew := false
// 	if subTime.Format("2006-01-02") == "0001-01-01" {
// 		mblock.FieldByName(mOp + "SubmitAt").Set(reflect.ValueOf(time.Now()))
// 		isNew = true
// 	}

// 	// mblock.FieldByName("Stage").Set(reflect.ValueOf(op + "Cache"))
// 	opCode := mblock.FieldByName(mOp + "Code").Interface().(string)

// 	return global.ProDbMap[proCode].Transaction(func(tx *gorm.DB) error {
// 		// err := service.UpdateBlock(proCode, block, block.ID)
// 		// query := op + "_code = ? AND id = ? AND " + op + "_submit_at = '0001-01-01 08:05:43+08:05:43'"
// 		query := op + "_code = ? AND id = ? AND " + op + "_stage = ?"
// 		if isNew {
// 			query += " AND to_char(" + op + "_submit_at, 'YYYY-MM-DD') = '0001-01-01'"
// 		}
// 		err := tx.Select(mOp+"SubmitAt").Where(query, opCode, block.ID, op+"Cache").Updates(block).Error
// 		// err := tx.Where("id = ?", block.ID).Updates(block).Error
// 		if err != nil {
// 			return err
// 		}
// 		tFields := taskSubmit.Fields
// 		ids := []string{}
// 		for ii, fields := range tFields {
// 			for jj, field := range fields {
// 				field.BlockIndex = ii
// 				field.FieldIndex = jj
// 				if field.ID != "" {
// 					// err = service.UpdateField(proCode, field, field.ID)
// 					opInput := ""
// 					if op == "op1" {
// 						opInput = field.Op1Input
// 					} else if op == "op2" {
// 						opInput = field.Op2Input
// 					}
// 					if opInput != "no" {
// 						err = tx.Select("BlockIndex", "FieldIndex", "ResultValue", "ResultInput", mOp+"Value", mOp+"Input").Where("id = ?", field.ID).Updates(field).Error
// 					} else {
// 						err = tx.Select("BlockIndex", "FieldIndex", mOp+"Value", mOp+"Input").Where("id = ?", field.ID).Updates(field).Error
// 					}
// 				} else {
// 					// err, _ = service.InsertField(proCode, field)
// 					err = tx.Create(&field).Error
// 				}
// 				if err != nil {
// 					return err
// 				}
// 				ids = append(ids, field.ID)
// 			}
// 		}
// 		if len(ids) > 0 {
// 			err = tx.Where("id not in ? AND block_id = ?", ids, block.ID).Delete(model.ProjectField{}).Error
// 		}
// 		return err
// 	})
// }

func GetTaskBlock(proCode string, code, name string) (err error, block model.PracticeProjectBlock, user model.PracticeUser, pacticeSum model.PracticeSum) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, block, user, pacticeSum
	}
	// user := model.PracticeUser{}
	user.ApplyAt = time.Now()
	user.ProCode = proCode
	user.Code = code
	user.Name = name
	err = db.Model(&model.PracticeUser{}).Where("code = ? and pro_code = ?", code, proCode).FirstOrCreate(&user).Error
	if err != nil {
		return err, block, user, pacticeSum
	}
	if user.SubmitAt.Format("2006-01-02") != "0001-01-01" {
		user.ApplyAt = time.Now()
		user.SubmitAt = time.Time{}
		err = db.Save(&user).Error
	}
	db = db.Model(&model.PracticeProjectBlock{}).Where("is_input = true and status = 1")
	var total int64
	err = db.Count(&total).Error
	if len(user.CacheId) > 0 {
		ids := []string{}
		ids = user.CacheId
		db.Not(ids)
	}
	db.Order("RANDOM()")
	err = db.First(&block).Error
	if len(user.CacheId) > 0 && total > 0 && err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user.CacheId = []string{}
		db = global.ProDbMap[proCode]
		err = db.Save(&user).Error
		if err != nil {
			return err, block, user, pacticeSum
		} else {
			return GetTaskBlock(proCode, code, name)
		}
	}
	fmt.Println("----------err1111-----------", err)
	_, pacticeSum = GetPracticeSum(proCode, user)
	fmt.Println("----------er2222-----------", err)
	return err, block, user, pacticeSum
}

func GetPracticeSum(proCode string, user model.PracticeUser) (err error, pacticeSum model.PracticeSum) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, pacticeSum
	}
	err = db.Model(&model.PracticeSum{}).Where("code = ? and apply_at = ?", user.Code, user.ApplyAt).First(&pacticeSum).Error
	return err, pacticeSum
}

func SubmitTask(proCode string, taskSubmit model.TaskSubmit, code, id, bcode string) (err error) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	user := model.PracticeUser{}
	err = db.Model(&model.PracticeUser{}).Where("code = ? and pro_code = ?", code, proCode).First(&user).Error
	if err != nil {
		return err
	}
	user.CacheId = append(user.CacheId, id)

	if arrays.ContainsString(user.Bcode, bcode) == -1 {
		user.Bcode = append(user.Bcode, bcode)
	}

	err = db.Save(&user).Error
	return Sum_func(proCode, taskSubmit, user)
}

func Exit_func(proCode, code string) (err error, mes string) {
	db := global.ProDbMap[proCode]
	user := model.PracticeUser{}
	err = db.Model(&model.PracticeUser{}).Where("code = ? and pro_code = ?", code, proCode).First(&user).Error
	if err != nil {
		return err, mes
	}
	user.SubmitAt = time.Now()
	user.CacheId = []string{}
	err = db.Save(&user).Error
	if err != nil {
		return err, mes
	}
	pacticeSum := model.PracticeSum{}
	err = db.Model(&model.PracticeSum{}).Where("code = ? and apply_at = ?", user.Code, user.ApplyAt).First(&pacticeSum).Error
	if err != nil {
		return err, mes
	}
	pacticeSum.SubmitAt = time.Now()
	costTime := pacticeSum.SubmitAt.Unix() - pacticeSum.ApplyAt.Unix()
	pacticeSum.SummaryCostTime = costTime
	err = db.Save(&pacticeSum).Error
	if err != nil {
		return err, mes
	}
	err, ask := PracticeAsk(proCode)
	if pacticeSum.SummaryFieldCharacter <= ask.Character {
		mes += "字符数少于合格标准;"
	}
	if pacticeSum.SummaryAccuracyRate*100 <= ask.AccuracyRate {
		mes += "准确率低于合格标准;"
	}
	fmt.Println("-----------------pacticeSumpacticeSum-----------------------", pacticeSum.SummaryFieldCharacter, pacticeSum.SummaryAccuracyRate)
	fmt.Println("-----------------ask-----------------------", ask.Character, ask.AccuracyRate)

	global.GLog.Info("service.Exit_func return---> mes:" + mes + "::")
	if mes == "" {
		err = trainStage.SetTrainingStage4Practice(proCode, code)
	}

	return err, mes
}

func CheckPracticeUser() {
	for true {
		proCode := global.GConfig.System.ProCode
		db := global.ProDbMap[proCode]
		users := []model.PracticeUser{}
		apply_at := time.Now().AddDate(0, 0, -5)
		err := db.Model(&model.PracticeUser{}).Where("pro_code = ? AND to_char(submit_at,'YYYY-MM-DD') = ? AND apply_at <= ? ", proCode, "0001-01-01", apply_at).Find(&users).Error
		if err == nil {
			for _, user := range users {
				Exit_func(proCode, user.Code)
			}
		} else {
			return
		}
		<-time.After(60 * time.Second)
	}

}

func PracticeAsk(proCode string) (err error, ask model.PracticeAsk) {
	err = global.GDb.Model(&model.PracticeAsk{}).Where("pro_code = ?", proCode).First(&ask).Error
	return err, ask
}

func PracticeAskList() (err error, list []model.PracticeAsk) {
	err = global.GDb.Model(&model.PracticeAsk{}).Find(&list).Error
	return err, list
}

func Sum_func(proCode string, taskSubmit model.TaskSubmit, user model.PracticeUser) (err error) {
	db := global.ProDbMap[proCode]
	pacticeSum := model.PracticeSum{}
	pacticeSum.ApplyAt = user.ApplyAt
	pacticeSum.Code = user.Code
	pacticeSum.Name = user.Name
	pacticeSum.ProCode = proCode
	// and to_char(submit_at, 'YYYY-MM-DD') = '0001-01-01'
	err = db.Model(&model.PracticeSum{}).Where("code = ? and apply_at = ?", user.Code, user.ApplyAt).FirstOrCreate(&pacticeSum).Error
	if err != nil {
		return err
	}
	pacticeSum.SummaryBlockNum += 1
	// time.Now()
	costTime := time.Now().Unix() - pacticeSum.ApplyAt.Unix()
	// time.Unix(costTime, 0).Format("15:04:05")
	pacticeSum.SummaryCostTime = costTime
	blockEfficiency := Decimal(3600.0 / (float64(costTime) / float64(pacticeSum.SummaryBlockNum)))
	pacticeSum.SummaryBlockEfficiency = blockEfficiency
	wrongs := []model.PracticeWrong{}
	for _, fields := range taskSubmit.Fields {
		for _, field := range fields {
			if field.Op1Input == "yes" {
				if field.Op1Value != field.ResultValue {
					wrongs = append(wrongs, model.PracticeWrong{
						SubmitDay: time.Now(),
						Code:      user.Code,
						NickName:  user.Name,
						BillNum:   taskSubmit.Bill.BillNum,
						FieldCode: field.Code,
						FieldName: field.Name,
						Path:      taskSubmit.Bill.DownloadPath,
						Picture:   taskSubmit.Block.Picture,
						Wrong:     field.Op1Value,
						Right:     field.ResultValue,
					})
				}
				pacticeSum.SummaryFieldNum += 1
				if strings.Index(field.Op1Value, "?") != -1 {
					pacticeSum.SummaryQuestionMarkNumber += 1
				}
				len := utf8.RuneCountInString(field.Op1Value)
				if len > utf8.RuneCountInString(field.ResultValue) {
					len = utf8.RuneCountInString(field.ResultValue)
				}
				i := -1
				for _, val := range field.Op1Value {
					i++
					fieldCharacter := 1
					if unicode.Is(unicode.Han, val) {
						fieldCharacter = 2
					}
					pacticeSum.SummaryFieldCharacter += fieldCharacter
					if i < len && Substr(field.Op1Value, i, i+1) == Substr(field.ResultValue, i, i+1) {
						pacticeSum.SummaryFieldEffectiveCharacter += fieldCharacter
					}
				}
			}
		}
	}
	fieldEfficiency := Decimal(float64(pacticeSum.SummaryFieldCharacter) / float64(costTime) * 60 * 60)
	pacticeSum.SummaryFieldEfficiency = fieldEfficiency
	questionMarkProportion := 0.0
	if pacticeSum.SummaryFieldNum > 0 {
		questionMarkProportion = Decimal(float64(pacticeSum.SummaryQuestionMarkNumber) / float64(pacticeSum.SummaryFieldNum) * 100)
	}
	pacticeSum.SummaryQuestionMarkProportion = questionMarkProportion
	//错误产量
	if len(wrongs) > 0 {
		err = db.Create(&wrongs).Error
		if err != nil {
			return err
		}
	}

	var total int64
	err = db.Model(&model.PracticeWrong{}).Where("code = ? and submit_day BETWEEN ? AND ? ", user.Code, pacticeSum.ApplyAt, time.Now()).Count(&total).Error
	if err != nil {
		return err
	}
	effectiveFieldNum := pacticeSum.SummaryFieldNum - int(total)
	if pacticeSum.SummaryFieldNum > 0 {
		accuracyRate := Decimal(float64(effectiveFieldNum) / float64(pacticeSum.SummaryFieldNum))
		pacticeSum.SummaryAccuracyRate = accuracyRate
	}

	err = db.Save(&pacticeSum).Error
	if err != nil {
		return err
	}
	return err
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

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

func SelectBillByID(proCode string, id string) (error, model.PracticeProjectBill) {
	var configsRes model.PracticeProjectBill
	err := global.ProDbMap[proCode].Where("id = ?", id).First(&configsRes).Error
	return err, configsRes
}

func SelectOpFieldsByBlockID(proCode string, block_id string) (err error, configs []model.PracticeProjectField) {
	db := global.ProDbMap[proCode].Model(&model.PracticeProjectField{})
	// db.Select("id, created_at, updated_at, name, code, bill_id, block_id, block_index, field_index, op1_value, op1_input, op2_value, op2_input, op_q_value, op_q_input, op_d_value, op_d_input, result_value, result_input, final_value, final_input, right_value")
	if block_id != "" {
		db = db.Where("block_id = ?", block_id)
	}
	// db.Group("block_index")
	db.Order("updated_at asc").Order("block_index asc").Order("field_index asc")
	var configsRes []model.PracticeProjectField
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func FieldsFormat(fields []model.PracticeProjectField) [][]model.PracticeProjectField {
	data := [][]model.PracticeProjectField{}
	blockIndex := -1
	value := []model.PracticeProjectField{}
	for _, field := range fields {
		if field.BlockIndex != blockIndex {
			if len(value) > 0 {
				data = append(data, value)
				value = []model.PracticeProjectField{}
			}
			blockIndex = field.BlockIndex
		}
		value = append(value, field)
	}
	if len(value) > 0 {
		data = append(data, value)
	}
	return data
}

// GetPracticeSumLists
func GetPracticeSumLists(sumSearch model.SumSearch) (err error, total int64, list []model.PracticeSum) {
	limit := sumSearch.PageSize
	offset := sumSearch.PageSize * (sumSearch.PageIndex - 1)

	db := global.ProDbMap[sumSearch.ProCode]
	if db == nil {
		return global.ProDbErr, 0, list
	}
	db = db.Model(&model.PracticeSum{}).Where("apply_at >= ? AND submit_at <= ? ", sumSearch.StartTime, sumSearch.EndTime)
	if sumSearch.Code != "" {
		db.Where("code = ?", sumSearch.Code)
	}
	if sumSearch.Name != "" {
		db.Where("name like ?", "%"+sumSearch.Name+"%")
	}

	err = db.Count(&total).Error
	err = db.Order("apply_at").Limit(limit).Offset(offset).Find(&list).Error

	return err, total, list
}

func GetPracticeWrongLists(sumSearch model.SumSearch) (err error, total int64, list []model.PracticeWrong) {
	limit := sumSearch.PageSize
	offset := sumSearch.PageSize * (sumSearch.PageIndex - 1)

	db := global.ProDbMap[sumSearch.ProCode]
	if db == nil {
		return global.ProDbErr, 0, list
	}
	db = db.Model(&model.PracticeWrong{}).Where("submit_day BETWEEN ? AND ?", sumSearch.StartTime, sumSearch.EndTime)
	if sumSearch.Name != "" {
		db.Where("field_name like ?", "%"+sumSearch.Name+"%")
	}

	if sumSearch.Code != "" {
		db.Where("code = ?", sumSearch.Code)
	}

	err = db.Count(&total).Error
	err = db.Order("submit_day").Limit(limit).Offset(offset).Find(&list).Error

	return err, total, list
}

func GetBlockVideo(code string) []string {
	var teachVideo model2.TeachVideo
	proCode := global.GConfig.System.ProCode
	proId := global.ProCodeId[proCode]
	fmt.Println("---------proId------------:", proId)
	err := global.GDb.Model(&model2.TeachVideo{}).Where("pro_id = ? and sys_block_code = ?", proId, code).First(&teachVideo).Error
	fmt.Println("---------err------------:", err)
	if err != nil {
		return []string{}
	}
	return teachVideo.Video
}
