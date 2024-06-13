package task

import (
	"fmt"
	"reflect"
	"server/global"
	"server/module/load/model"
	"server/module/load/service"
	"strings"
	"time"
)

var ReleaseTime float64 = 1800 //  秒
var NextTime float64 = 60      //  秒
var TaskProCode string = global.GConfig.System.ProCode + "_task"
var CheckTaskSecond time.Duration = 30 // 秒
var TaskListSecond time.Duration = 30  // 秒
const ImagePageSize = 20

// func BlockCheck(block model.ProjectBlock, op string) bool {
// 	isOk := true
// 	isOk = PreBlockOK(block, op)
// 	if !isOk {
// 		return isOk
// 	}
// 	return isOk
// }

// CacheFieldConf

//字段判断
func FieldsInputCheck(fields []model.ProjectField, op string) bool {
	isOk := false
	mOp := strings.Replace(op, "o", "O", -1)
	fmt.Println("----------FieldsInputCheck----------", mOp)
	for _, field := range fields {
		rfield := reflect.ValueOf(&field).Elem()
		opInput := rfield.FieldByName(mOp + "Input").Interface().(string)
		if opInput == "yes" {
			isOk = true
			return isOk
		}
	}
	// isOk = PreBlock(block)
	return isOk
}

// 前置分块
func PreBlockOK(block model.ProjectBlock, op string) bool {
	isOk := true
	preBCode := block.PreBCode
	op = strings.Replace(op, "o", "O", -1)
	if len(preBCode) > 0 {
		_, blocks := service.SelectBlockByCodes(TaskProCode, block.BillID, preBCode)
		for _, bk := range blocks {
			mblock := reflect.ValueOf(bk)
			opCode := mblock.FieldByName(op + "Code").Interface().(string)
			if opCode == "0" {
				continue
			}
			submitAt := mblock.FieldByName(op + "SubmitAt").Interface().(time.Time)
			if submitAt.Format("2006-01-02") == "0001-01-01" {
				return false
			}
			subM := time.Now().Sub(submitAt)
			if subM.Seconds() < NextTime {
				return false
			}
		}
	}
	return isOk
}

// 黑名单
func BlackListBlock(block model.ProjectBlock) bool {
	isOk := true
	return isOk
}

// func StructToMapViaJson(data interface{}) map[string]interface{} {
// 	m := make(map[string]interface{})
// 	j, _ := json.Marshal(data)
// 	json.Unmarshal(j, &m)
// 	return m
// }
