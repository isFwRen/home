package task

import (
	"fmt"
	"reflect"
	"server/global"
	"server/global/response"
	"server/module/load/model"
	"server/module/load/service"
	"server/module/task/project"
	"strings"
	"time"
)

func DoCheckTask() {
	fmt.Println("------------------------------DoCheckTask------------------------------------------")
	for true {
		fmt.Println("DoCheckTask")
		// condition := map[string]interface{}{"stage": "op0"}
		// err, blocks := service.SelectTaskBlocks("B0118", condition) //[]model.ProjectBlock{} // ervice.SelectBlock() op1_ceahe
		err, ids := service.SelectCacheTaskBlocks(TaskProCode)
		fmt.Println("---------------DoCheckTask-------------------------:", err, len(ids))
		for _, id := range ids {
			// fmt.Println("DoCheckTask", ii, id)
			err, block := service.SelectBlockByID(TaskProCode, id)
			block, fields, op := CheckTimeOp(block)
			// service.SaveBlock(TaskProCode, block, block.ID)
			// fmt.Println("-----------UpdateBlockAndFields-------------------", op)
			err = service.UpdateBlockAndFields(TaskProCode, block, fields, op)
			if err != nil {
				fmt.Println("-----------更新失败---------------", block.ID, op, err)
			}
			// fmt.Println("blockblockblock", block)
		}
		fmt.Println("已处理分块数量:", len(ids))
		<-time.After(CheckTaskSecond * time.Second)
		// <-timeAfterTrigger
	}

}

func Check_next_op(block model.ProjectBlock) model.ProjectBlock {
	// if true {
	// 	return nil
	// }
	return block
}

func CheckTimeOp(block model.ProjectBlock) (model.ProjectBlock, []model.ProjectField, string) {
	// json, err := json.Marshal(block)
	if block.Op1Stage == "opCache" {
		err, fields := service.SelectFieldsByBlockID(TaskProCode, block.ID)
		if err != nil {
			return block, fields, ""
		}
		block, fields = NextOp(block, fields, "opCache")
		return block, fields, ""
	}
	// fmt.Println("-----------CheckTimeO-------------------", block.Name)
	nowTime := time.Now()
	mblock := reflect.ValueOf(block)
	// stage := strings.Replace(block.Stage, "Cache", "", -1)
	// stage = strings.Replace(stage, "o", "O", -1)
	op := ""
	if block.Op0Stage == "op0Cache" {
		op = "op0"
	} else if block.Op1Stage == "op1Cache" {
		op = "op1"
	} else if block.Op2Stage == "op2Cache" {
		op = "op2"
	} else if block.OpqStage == "opqCache" {
		op = "opq"
	}
	if op == "" {
		return block, []model.ProjectField{}, op
	}
	mOp := strings.Replace(op, "o", "O", -1)
	// fmt.Println("----------stage---------", mOp)
	applyAt := mblock.FieldByName(mOp + "ApplyAt").Interface().(time.Time)
	submitAt := mblock.FieldByName(mOp + "SubmitAt").Interface().(time.Time)
	// applyAt, submitAt := GetOpTime(block)
	// fmt.Println("applyAt  submitAt ", applyAt, submitAt)
	pFields := []model.ProjectField{}
	if submitAt.Format("2006-01-02") == "0001-01-01" {
		subM := nowTime.Sub(applyAt)
		// fmt.Println("IsZero subM ", subM.Minutes())
		cacheReleaseTime := ReleaseTime
		if block.FreeTime > 0 {
			cacheReleaseTime = float64(block.FreeTime)
		}
		if subM.Seconds() > cacheReleaseTime {
			opCode := mblock.FieldByName(mOp + "Code").Interface().(string)
			block = ClearOp(block, op+"Cache")
			BroadcastRelease("", block.ID, op, opCode)
		}
		return block, pFields, op

	}
	// mm, _ := time.ParseDuration("1m")4e4e4
	// overtime := applyAt.Add(mm)
	subM := nowTime.Sub(submitAt)
	// fmt.Println("-----------------subM!!!!-----------------", block.ID, submitAt, nowTime, subM.Seconds(), NextTime)
	if subM.Seconds() >= NextTime {
		// if op != "op0" {
		_, pFields = service.SelectFieldsByBlockID(TaskProCode, block.ID)
		// }
		block, pFields = NextOp(block, pFields, op+"Cache")
	}
	return block, pFields, op
}

func ClearOp(block model.ProjectBlock, op string) model.ProjectBlock {
	// op := block.Stage
	switch op {
	case "op1Cache":
		block.Op1ApplyAt = time.Time{}
		block.Op1Code = ""
		block.Op1Stage = "op1"
		if block.IsCompetitive {
			block.Op2ApplyAt = time.Time{}
			block.Op2SubmitAt = time.Time{}
			block.Op2Code = ""
			block.Op2Stage = "op2"
		}
		break
	case "op2Cache":
		block.Op2ApplyAt = time.Time{}
		block.Op2Code = ""
		block.Op2Stage = "op2"
		if block.IsCompetitive {
			block.Op1ApplyAt = time.Time{}
			block.Op1SubmitAt = time.Time{}
			block.Op1Code = ""
			block.Op1Stage = "op1"
		}
		break
	case "opqCache":
		block.OpqApplyAt = time.Time{}
		block.OpqCode = ""
		block.OpqStage = "opq"
		break
	case "op0Cache":
		block.Op0ApplyAt = time.Time{}
		block.Op0Code = ""
		block.Op0Stage = "op0"
		break
	}
	return block
}

func NextOp(block model.ProjectBlock, fields []model.ProjectField, op string) (model.ProjectBlock, []model.ProjectField) {
	// op := block.Stage
	switch op {
	case "op0Cache":
		block.Op0Stage = "crop"
		return block, fields
	case "op1Cache":
		isOK := PreBlockOK(block, "op2")
		if !isOK {
			break
		}
		// block.Op1SubmitAt = time.Now()
		block.Op1Stage = "done"
		if block.Op2Stage == "done" {
			return NextOp(block, fields, "op2Cache")
		}
		if block.Op2Stage == "op2" || block.Op2Stage == "op2Cache" || block.Op2Stage == "opCache" {
			break
		}
		block.Op2Stage = "op2"
		fields = project.DisableValue(block, fields, "op2")
		if block.IsLoop || !FieldsInputCheck(fields, "op2") {
			block.Op2Code = "0"
			block.Op2ApplyAt = time.Now()
			block.Op2SubmitAt = time.Now()
			block.Op2Stage = "op2Cache"
			return NextOp(block, fields, "op2Cache")
		}
		break
	case "op2Cache":
		if block.Op2Stage != "done" {
			isOK := PreBlockOK(block, "opq")
			if !isOK {
				break
			}
			// block.Op2SubmitAt = time.Now()
			block.Op2Stage = "done"
		}
		if block.Op2Stage != "done" || block.Op1Stage != "done" {
			break
		}
		block.OpqStage = "opq"
		fields = project.OpQCheck(block, fields, "opq")
		if block.IsLoop || !FieldsInputCheck(fields, "opq") {
			block.OpqCode = "0"
			block.OpqApplyAt = time.Now()
			block.OpqSubmitAt = time.Now()
			block.OpqStage = "done"
			return NextOp(block, fields, "opqCache")
		}
		break
	case "opqCache":
		// block.OpqSubmitAt = time.Now()
		block.OpqStage = "done"
		break
	default:
		if block.Op1Stage == "opCache" {
			isOK := PreBlockOK(block, "op1")
			if !isOK {
				break
			}
			block.Op1Stage = "op1"
			fields = project.DisableValue(block, fields, "op1")
			if !FieldsInputCheck(fields, "op1") {
				block.Op1Code = "0"
				block.Op1ApplyAt = time.Now()
				block.Op1SubmitAt = time.Now()
				block.Op1Stage = "op1Cache"
			}
		}
		if block.Op2Stage == "opCache" {
			block.Op2Stage = "op2"
			fields = project.DisableValue(block, fields, "op2")
			if !FieldsInputCheck(fields, "op2") {
				block.Op2Code = "0"
				block.Op2ApplyAt = time.Now()
				block.Op2SubmitAt = time.Now()
				block.Op2Stage = "op2Cache"
			}
		}
		if block.Op1Stage == "op1Cache" {
			return NextOp(block, fields, "op1Cache")
		}

		// block = project.DisableValue(block, fields)
		// block.OpDSubmitAt = time.Now()
		// block.Stage = "done"
		break
	}
	return block, fields
}

func BroadcastRelease(bill_id string, block_id string, op, userCode string) {
	data := map[string]string{
		"billId":  bill_id,
		"blockId": block_id,
		"op":      op,
	}
	fmt.Println("------------------------BroadcastRelease------userCode-----------------:", data, userCode)
	if userCode != "" {
		userSocket, isExit := global.GSocketConnMap[userCode]
		if isExit {
			userSocket.Emit("release", response.Response{
				Code: 200,
				Data: data,
				Msg:  "释放分块!",
			})
		}
	}

	// global.GSocketIo.BroadcastToNamespace("/global-release", "release", response.Response{
	// 	Code: 200,
	// 	Data: data,
	// 	Msg:  "释放分块",
	// })
}

// func CheckFieldsInput(fields []model.ProjectField) {
// 	for ii, field := range fields {
// 		fmt.Println("CheckFieldsInput:", ii, field)
// 		// field["aaa"]
// 	}
// }

// func StructToMapViaJson(data interface{}) map[string]interface{} {
// 	m := make(map[string]interface{})
// 	// t := time.Now()
// 	j, _ := json.Marshal(data)
// 	json.Unmarshal(j, &m)
// 	// fmt.Println(m)
// 	return m
// 	// fmt.Println(time.Now().Sub(t))
// }

// func BlockMapToStruct(mapInstance map[string]interface{}, data model.ProjectBlock) model.ProjectBlock {
// 	//将 map 转换为指定的结构体
// 	if err := mapstructure.Decode(mapInstance, &data); err != nil {
// 		fmt.Println("aaaaaaaaaa:", err)
// 	}
// 	return data
// }
