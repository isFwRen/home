package service

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"server/global"
	"server/module/load/model"
	billModel "server/module/pro_manager/model"
	baseModel "server/module/sys_base/model"
	"strconv"
	"strings"
	"time"
	"unicode"

	sumModel "server/module/report_management/model"
	// sumModel "server/module/report_management/model"
	// "strings"

	"gorm.io/gorm"
)

func InsertBlock(proCode string, agingConfig model.ProjectBlock) (err error, configInter model.ProjectBlock) {
	err = global.ProDbMap[proCode].Model(&model.ProjectBlock{}).Create(&agingConfig).Error
	return err, agingConfig
}

func DelBlocksByBillID(proCode string, bill_id string) (err error) {
	var configAging model.ProjectBlock
	err = global.ProDbMap[proCode].Where("bill_id = ?", bill_id).Delete(&configAging).Error
	return err
}

func SaveBlock(proCode string, configAging model.ProjectBlock, id string) (err error) {
	err = global.ProDbMap[proCode].Where("id = ?", id).First(&model.ProjectBlock{}).Save(&configAging).Error
	return err
}

// func UpdateBlock(proCode string, configAging model.ProjectBlock, id string) (err error) {
// 	err = global.ProDbMap[proCode].Where("id = ?", id).Updates(configAging).Error
// 	return err
// }

func UpdateBlock(proCode string, data map[string]interface{}, id string, op string) (err error) {
	// query := "id = ? AND " + op + "_submit_at = '0001-01-01 08:05:43+08:05:43'"
	query := "id = ? AND to_char(" + op + "_submit_at, 'YYYY-MM-DD') = '0001-01-01'"
	err = global.ProDbMap[proCode].Model(&model.ProjectBlock{}).Where(query, id).Updates(data).Error
	return err
}

func SelectBlock(proCode string, name string) (err error, configs []model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	if name != "" {
		db = db.Where("name = ?", name)
	}
	var configsRes []model.ProjectBlock
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func SelectBlockByID(proCode string, id string) (error, model.ProjectBlock) {
	var configsRes model.ProjectBlock
	err := global.ProDbMap[proCode].Where("id = ?", id).First(&configsRes).Error
	return err, configsRes
}
func CountBlock(proCode string, name string) (err error, total int64) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	if name != "" {
		db = db.Where("name = ?", name)
	}
	err = db.Count(&total).Error
	return err, total
}

func SelectTaskBlocks(proCode string, condition map[string]interface{}) (err error, configs []model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	db = db.Where(condition)
	var configsRes []model.ProjectBlock
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func SelectCacheTaskBlocks(proCode string) (err error, configs []string) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	db = db.Select("id").Where("op0_stage = ? OR op1_stage LIKE ? OR op2_stage = ? OR opq_stage = ?", "op0Cache", "%Cache", "op2Cache", "opqCache")
	var configsRes []string
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func GetUserTaskBlock(proCode string, op string, code string) (err error, configs model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	// query := op + "_code = ? AND " + op + "_submit_at like '0001-01-01 08:05:43+08:05:43'"
	query := op + "_code = ? AND to_char(" + op + "_submit_at, 'YYYY-MM-DD') = '0001-01-01'"
	fmt.Println("queryqueryquery:", query)
	db = db.Where(query, code)
	// db.Order("aaa desc")
	var configsRes model.ProjectBlock
	err = db.First(&configsRes).Error
	return err, configsRes
}

func GetUserModifyTaskBlock(proCode string, op string, code string, num int) (err error, configs model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	query := op + "_code = ? AND " + op + "_stage = ?"
	// fmt.Println("queryqueryquery:", query)
	db = db.Where(query, code, op+"Cache")
	db.Order(op + "_submit_at desc")
	var configsRes model.ProjectBlock
	err = db.Limit(1).Offset(num).Find(&configsRes).Error
	return err, configsRes
}

func SelectCropBlocks(proCode string) (err error, configs []model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	db = db.Where("op0_stage = ?", "crop")
	var configsRes []model.ProjectBlock
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func SelectBlockByCodes(proCode, bill_id string, codes []string) (err error, configs []model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	db.Where("bill_id = ? AND code IN ? AND temp != '未定义'", bill_id, codes)
	var configsRes []model.ProjectBlock
	err = db.Find(&configsRes).Error
	return err, configsRes
}

func CountBlockOpNum(proCode string, op string) (err error, total int64) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	query := op + "_stage = ? "
	db = db.Where(query, op)

	err = db.Count(&total).Error
	return err, total
}

func CheckBlockDone(proCode string, bill_id string, block_id string) (err error, total int64) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	db = db.Where("opq_stage <> ? and bill_id = ? and status = 2 and id <> ?", "done", bill_id, block_id)
	err = db.Count(&total).Error
	return err, total
}

func SelectCodeOutput(proCode string, code string, day time.Time) (err error, configs sumModel.OutputStatistics) {
	fmt.Println("---------------SelectCodeOutput-----------------:", proCode)
	db := global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{})
	db = db.Where("code = ? and submit_time = ?", code, day)
	var configsRes sumModel.OutputStatistics
	err = db.First(&configsRes).Error
	return err, configsRes
}

func DayOutput(proCode string, day time.Time) (err error, configs []sumModel.OutputStatistics) {
	fmt.Println("---------------DayOutput-----------------:", proCode)
	proCode = strings.Replace(proCode, "_task", "", -1)
	db := global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{})
	db = db.Where("submit_time = ?", day)
	var configsRes []sumModel.OutputStatistics
	err = db.Select("code", "summary_field_character", "summary_accuracy_rate").Find(&configsRes).Error
	return err, configsRes
}

func SumBlockOutput(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) error {
	err := error(nil)
	if op == "" {
		return err
	}
	mOp := strings.Replace(op, "o", "O", -1)
	// fmt.Println("---------------mOp-----------------:", mOp)
	proCode = strings.Replace(proCode, "_task", "", -1)
	mblock := reflect.ValueOf(&block).Elem()
	stage := mblock.FieldByName(mOp + "Stage").Interface().(string)
	code := mblock.FieldByName(mOp + "Code").Interface().(string)
	// fmt.Println("---------------stage-----------------:", code, stage)
	if code == "0" || code == "" || (stage != "done" && stage != "crop") {
		return err
	}

	applyAt := mblock.FieldByName(mOp + "ApplyAt").Interface().(time.Time)
	submitAt := mblock.FieldByName(mOp + "SubmitAt").Interface().(time.Time)

	fieldNum := 0
	fieldCharacter := 0
	blockNum := 1
	costTime := submitAt.Unix() - applyAt.Unix()
	// fmt.Println("----------------------------costTime----------------------------:", code, applyAt, submitAt, submitAt.Unix(), applyAt.Unix(), costTime)
	submitTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02")) // time.Now().Format("2006-01-02")
	// fmt.Println("---------------SumBlockOutput-----------------:", code, costTime, submitTime)
	err, outputStatistics := SelectCodeOutput(proCode, code, submitTime)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		outputStatistics = sumModel.OutputStatistics{}
		outputStatistics.Code = code
		outputStatistics.SubmitTime = submitTime
	} else if err != nil {
		return err
	}
	questionMarkNumber := 0
	for _, field := range fields {
		mfield := reflect.ValueOf(&field).Elem()
		input := mfield.FieldByName(mOp + "Input").Interface().(string)
		value := mfield.FieldByName(mOp + "Value").Interface().(string)
		// fmt.Println("---------------mfield-----------------:", input, value)
		if input == "no" || input == "no_if" {
			continue
		}
		fieldNum = fieldNum + 1
		if value == "" {
			continue
		}
		if strings.Index(value, "?") != -1 || strings.Index(value, "？") != -1 {
			questionMarkNumber = questionMarkNumber + 1
		}
		aaa := 0
		for _, rr := range value {
			if rr == '?' || rr == '？' {
				continue
			}
			if unicode.Is(unicode.Han, rr) {
				aaa = aaa + 2
			} else {
				aaa = aaa + 1
			}
		}
		if field.Op0Input == "ocr" && field.Op0Value != "" {
			aaa = int(math.Round(float64(aaa) * 0.6))
		}
		fieldCharacter += aaa
		// fieldNum = fieldNum + 1
	}
	mOp = strings.Replace(mOp, "q", "Q", -1)
	moutput := reflect.ValueOf(&outputStatistics).Elem()
	xops := []string{mOp, "Summary"}
	for _, xop := range xops {
		blockNumCache := moutput.FieldByName(xop+"BlockNum").Interface().(int) + blockNum
		fieldNumCache := moutput.FieldByName(xop+"FieldNum").Interface().(int) + fieldNum
		fieldCharacterCache := moutput.FieldByName(xop+"FieldCharacter").Interface().(int) + fieldCharacter
		fieldEffectiveCharacter := moutput.FieldByName(xop+"FieldEffectiveCharacter").Interface().(int) + fieldCharacter
		costTimeCache := moutput.FieldByName(xop+"CostTime").Interface().(int64) + costTime
		fieldEfficiency := Decimal(float64(fieldCharacterCache) / float64(costTimeCache) * 60 * 60)
		blockEfficiency := Decimal(3600.0 / (float64(costTimeCache) / float64(blockNumCache)))
		questionMarkNumberCache := moutput.FieldByName(xop+"QuestionMarkNumber").Interface().(int) + questionMarkNumber
		questionMarkProportion := 0.0
		if fieldNumCache > 0 {
			questionMarkProportion = Decimal(float64(questionMarkNumberCache) / float64(fieldNumCache) * 100)
		}
		fmt.Println("---------------questionMarkProportion------------------:", questionMarkNumberCache, fieldNumCache, questionMarkProportion, outputStatistics.Op1QuestionMarkNumber)
		if xop == "Op1" || xop == "Op2" {
			if strings.Index(block.Name, "报销") != -1 {
				expenseAccountFieldCharacter := moutput.FieldByName(xop+"ExpenseAccountFieldCharacter").Interface().(int) + fieldCharacter
				moutput.FieldByName(xop + "ExpenseAccountFieldCharacter").Set(reflect.ValueOf(expenseAccountFieldCharacter))
				expenseAccountFieldEffectiveCharacter := moutput.FieldByName(xop+"ExpenseAccountFieldEffectiveCharacter").Interface().(int) + fieldCharacter
				moutput.FieldByName(xop + "ExpenseAccountFieldEffectiveCharacter").Set(reflect.ValueOf(expenseAccountFieldEffectiveCharacter))
			} else {
				notExpenseAccountFieldCharacter := moutput.FieldByName(xop+"NotExpenseAccountFieldCharacter").Interface().(int) + fieldCharacter
				moutput.FieldByName(xop + "NotExpenseAccountFieldCharacter").Set(reflect.ValueOf(notExpenseAccountFieldCharacter))
				notExpenseAccountFieldEffectiveCharacter := moutput.FieldByName(xop+"NotExpenseAccountFieldEffectiveCharacter").Interface().(int) + fieldCharacter
				moutput.FieldByName(xop + "NotExpenseAccountFieldEffectiveCharacter").Set(reflect.ValueOf(notExpenseAccountFieldEffectiveCharacter))
			}
		}

		moutput.FieldByName(xop + "BlockNum").Set(reflect.ValueOf(blockNumCache))
		moutput.FieldByName(xop + "FieldNum").Set(reflect.ValueOf(fieldNumCache))
		moutput.FieldByName(xop + "FieldCharacter").Set(reflect.ValueOf(fieldCharacterCache))
		moutput.FieldByName(xop + "FieldEffectiveCharacter").Set(reflect.ValueOf(fieldEffectiveCharacter))
		moutput.FieldByName(xop + "CostTime").Set(reflect.ValueOf(costTimeCache))
		moutput.FieldByName(xop + "FieldEfficiency").Set(reflect.ValueOf(fieldEfficiency))
		moutput.FieldByName(xop + "BlockEfficiency").Set(reflect.ValueOf(blockEfficiency))
		moutput.FieldByName(xop + "QuestionMarkNumber").Set(reflect.ValueOf(questionMarkNumberCache))
		moutput.FieldByName(xop + "QuestionMarkProportion").Set(reflect.ValueOf(questionMarkProportion))
	}
	if outputStatistics.ID != "" {
		err = global.ProDbMap[proCode].Where("id = ?", outputStatistics.ID).First(&sumModel.OutputStatistics{}).Save(&outputStatistics).Error
	} else {
		var user baseModel.SysUser
		var user0p1 baseModel.SysUser
		var user0p2 baseModel.SysUser
		var user0pq baseModel.SysUser
		if block.Op0Code != "" {
			global.GDb.Model(&baseModel.SysUser{}).Where("code = ?", block.Op0Code).Find(&user)
			outputStatistics.NickName = user.NickName
			// err = global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{}).Create(&outputStatistics).Error
		}
		if block.Op1Code != "" {
			global.GDb.Model(&baseModel.SysUser{}).Where("code = ?", block.Op1Code).Find(&user0p1)
			outputStatistics.NickName = user0p1.NickName
			// err = global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{}).Create(&outputStatistics).Error
		}
		if block.Op2Code != "" {
			global.GDb.Model(&baseModel.SysUser{}).Where("code = ?", block.Op2Code).Find(&user0p2)
			outputStatistics.NickName = user0p2.NickName
			// err = global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{}).Create(&outputStatistics).Error
		}
		if block.OpqCode != "" {
			global.GDb.Model(&baseModel.SysUser{}).Where("code = ?", block.OpqCode).Find(&user0pq)
			outputStatistics.NickName = user0pq.NickName
			// err = global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{}).Create(&outputStatistics).Error
		}
		err = global.ProDbMap[proCode].Model(&sumModel.OutputStatistics{}).Create(&outputStatistics).Error
	}
	fmt.Println("------------------!AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Println("outputStatistics=", outputStatistics)
	fmt.Println("block.Op0Code = ? ; block.Op1Code = ?;block.Op2Code = ?;block.OpqCode = ? ", block.Op0Code, block.Op1Code, block.Op2Code, block.OpqCode)
	fmt.Println("------------------!AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	return err
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func UpdateBlockAndFields(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) error {
	return global.ProDbMap[proCode].Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&model.ProjectBlock{}).Where("id = ?", block.ID).Save(block).Error
		if err != nil {
			return err
		}
		for _, field := range fields {
			err = tx.Model(&model.ProjectField{}).Where("id = ?", field.ID).Updates(field).Error
			if err != nil {
				return err
			}
		}
		err = SumBlockOutput(proCode, block, fields, op)
		if err != nil {
			return err
		}
		if block.OpqStage == "done" {
			err, num := CheckBlockDone(proCode, block.BillID, block.ID)
			// fmt.Println("---------------num-----------------:", err, num)
			if err != nil {
				return err
			}
			if num == 0 {
				err = tx.Model(&billModel.ProjectBill{}).Where("id = ?", block.BillID).Update("stage", "6").Error
			}
		}
		// errors.Is(err, gorm.ErrRecordNotFound)
		return err
	})
}

// CountBlockOpNumAndName 统计某一类型分块数量
func CountBlockOpNumAndName(proCode string, op string, blockName string) (err error, total int64) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	db = db.Where(op+"_stage = ? and name like ?", op, "%"+blockName+"%")
	err = db.Count(&total).Error
	return err, total
}
