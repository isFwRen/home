package B0122

import (
	"regexp"
	"server/module/load/model"
)

func DisableValue(block model.ProjectBlock, fields []model.ProjectField, op string) []model.ProjectField {
	// proCode := "B0122_task"
	// if block.Code == "bc021" || block.Code == "bc030" {
	// 	_, fc027 := service.SelectBillFields(proCode, block.BillID, -1, "fc027")
	// 	if fc027.ResultValue == "1" || fc027.ResultValue == "A" {
	// 		for ii, field := range fields {
	// 			if field.Code == "fc029" || field.Code == "fc031" {
	// 				fields[ii] = SetOpInput(field, op)
	// 			}
	// 		}
	// 	}
	// }

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
