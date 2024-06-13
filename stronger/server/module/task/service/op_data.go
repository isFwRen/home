package service

import (
	"fmt"
	"reflect"
	"server/global"
	"server/module/load/model"
	"server/module/task"
	"server/module/task/request"
	"strings"
	"time"

	"gorm.io/gorm"
)

func SaveSubmitData(proCode string, taskSubmit request.TaskSubmit) error {
	op := taskSubmit.Op
	block := taskSubmit.Block
	mblock := reflect.ValueOf(&block).Elem()
	mOp := strings.Replace(op, "o", "O", -1)

	subTime := mblock.FieldByName(mOp + "SubmitAt").Interface().(time.Time)
	isNew := false
	if subTime.Format("2006-01-02") == "0001-01-01" {
		mblock.FieldByName(mOp + "SubmitAt").Set(reflect.ValueOf(time.Now()))
		isNew = true
	}

	// mblock.FieldByName("Stage").Set(reflect.ValueOf(op + "Cache"))
	opCode := mblock.FieldByName(mOp + "Code").Interface().(string)

	return global.ProDbMap[proCode].Transaction(func(tx *gorm.DB) error {
		// err := service.UpdateBlock(proCode, block, block.ID)
		// query := op + "_code = ? AND id = ? AND " + op + "_submit_at = '0001-01-01 08:05:43+08:05:43'"
		query := op + "_code = ? AND id = ? AND " + op + "_stage = ?"
		if isNew {
			query += " AND to_char(" + op + "_submit_at, 'YYYY-MM-DD') = '0001-01-01'"
		}
		err := tx.Select(mOp+"SubmitAt").Where(query, opCode, block.ID, op+"Cache").Updates(block).Error
		// err := tx.Where("id = ?", block.ID).Updates(block).Error
		if err != nil {
			return err
		}
		tFields := taskSubmit.Fields
		ids := []string{}
		for ii, fields := range tFields {
			for jj, field := range fields {
				field.BlockIndex = ii
				field.FieldIndex = jj
				if field.ID != "" {
					// err = service.UpdateField(proCode, field, field.ID)
					opInput := ""
					if op == "op1" {
						opInput = field.Op1Input
					} else if op == "op2" {
						opInput = field.Op2Input
					}
					if opInput != "no" {
						err = tx.Select("BlockIndex", "FieldIndex", "ResultValue", "ResultInput", mOp+"Value", mOp+"Input").Where("id = ?", field.ID).Updates(field).Error
					} else {
						err = tx.Select("BlockIndex", "FieldIndex", mOp+"Value", mOp+"Input").Where("id = ?", field.ID).Updates(field).Error
					}
				} else {
					// err, _ = service.InsertField(proCode, field)
					err = tx.Create(&field).Error
				}
				if err != nil {
					return err
				}
				ids = append(ids, field.ID)
			}
		}
		if len(ids) > 0 {
			err = tx.Where("id not in ? AND block_id = ?", ids, block.ID).Delete(model.ProjectField{}).Error
		}
		return err
	})
}

func GetTaskBlock(proCode string, op string, code string) (err error, configs model.ProjectBlock) {
	db := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
	query := op + "_stage = ? AND " + op + "_code = ?"
	fmt.Println("queryqueryquery:", query)
	db = db.Where(query, op, "")
	if op == "op2" {
		db = db.Where("op1_code != ?", code)
		// 	db.Where("is_competitive = ?", true)
	}
	if op == "op1" {
		db = db.Where("op2_code != ?", code)
		// 	db.Where("is_competitive = ?", true)
	}
	lists, isOK := task.CacheCodeWhiteList[code]
	// fmt.Println("------------------------lists-------------------", lists)
	if isOK && len(lists) > 0 {
		// db.Where( db.Or(lists...) )
		whiteQuery := ""
		// p := global.ProDbMap[proCode].Model(&model.ProjectBlock{})
		values := make([]interface{}, 0)
		where := true
		for _, list := range lists {
			if len(list.Code) == 0 {
				continue
			}
			if where {
				whiteQuery += "( temp = ? AND code IN ? )"
				where = false
			} else {
				whiteQuery += " OR ( temp = ? AND code IN ? )"
			}
			values = append(values, list.Temp)
			values = append(values, list.Code)
		}
		// fmt.Println("------------------------whiteQuery-------------------", whiteQuery)
		// fmt.Println("------------------------values---------------------", values)
		db = db.Where(whiteQuery, values...)
		// aaaa :=
		// db.Where("(temp = ? And code IN ?) || (temp = ? And code IN ?) || (temp = ? And code IN ?) ",list...)
		// db.Where( db.Where(lists[0]).Or(lists[1]).Or(lists[2])   )
		// blackQuery := map[string]interface{}{"code": blackList}
		// map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}
		// db.Not(blackQuery)
	} else {
		db = db.Where(" code = '11111' ")
	}
	if op == "op2" {
		db.Order("is_competitive asc")
	}
	db.Order("level asc")
	db.Order("created_at asc")
	var configsRes model.ProjectBlock
	err = db.First(&configsRes).Error
	return err, configsRes
}
