package B0121

import (
	"regexp"
	"server/module/load/model"
	"server/module/load/service"
)

func DisableValue(block model.ProjectBlock, fields []model.ProjectField, op string) []model.ProjectField {
	proCode := "B0121_task"
	if RegIsMatch(block.Code, `^(bc002|bc003|bc004)$`) {
		_, fc015 := service.SelectBillFields(proCode, block.BillID, -1, "fc015")
		if fc015.ResultValue == "A" {
			for ii, field := range fields {
				if RegIsMatch(field.Code, `^(fc016|fc017|fc020|fc026|fc024|fc023|fc022|fc025|fc018|fc019|fc021|fc352)$`) {
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
