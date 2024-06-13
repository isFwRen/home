package b0118

import (
	"fmt"
	"regexp"
	"server/module/load/model"
	"server/module/load/service"
)

func DisableValue(block model.ProjectBlock, fields []model.ProjectField, op string) []model.ProjectField {
	proCode := "B0118_task"
	//CSB0118RC0340000
	//1.当fc277字段录入值为2时，屏蔽fc278（同分块）、fc279票据代码、fc280开票日期、fc281校验码、fc282字段（不同分块）
	fmt.Println("block.Code", op, block.Code)
	if RegIsMatch(block.Code, `^(bc001)$`) {
		_, fc277 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc277")
		fmt.Println("fc277", "s", fc277.ResultValue)
		if fc277.ResultValue == "2" {
			for ii, field := range fields {
				if RegIsMatch(field.Code, `^(fc279|fc280|fc281|fc282)$`) {
					fmt.Println("fc277-------", field.Code)
					fields[ii] = SetOpInput(field, op)
				}
			}
		}
		//CSB0118RC0345000
		//当fc277录入值为1，fc278录入值为3时，屏蔽fc279、fc281字段（不同分块）
		_, fc278 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc278")
		if fc277.ResultValue == "1" && fc278.ResultValue == "3" {
			for ii, field := range fields {
				if RegIsMatch(field.Code, `^(fc279|fc281)$`) {
					fmt.Println("fc277--CSB0118RC0345000", field.Code)
					fields[ii] = SetOpInput(field, op)
				}
			}
		}
	}

	return fields
}

func RegIsMatch(value string, query string) bool {
	// reg := regexp.MustCompile(query)
	matched, _ := regexp.MatchString(query, value)
	return matched
}

func SetOpInput(field model.ProjectField, op string) model.ProjectField {
	if op == "op1" {
		field.Op1Input = "no_if"
	}
	if op == "op2" {
		field.Op2Input = "no_if"
	}
	if op == "opq" {
		field.OpqInput = "no_if"
	}
	return field
}
