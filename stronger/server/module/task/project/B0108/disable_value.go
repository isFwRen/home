package b0108

import (
	"regexp"
	"server/module/load/model"
	"server/module/load/service"
)

func DisableValue(block model.ProjectBlock, fields []model.ProjectField, op string) []model.ProjectField {
	proCode := "B0108_task"
	if block.Code == "bc021" || block.Code == "bc030" {
		_, fc027 := service.SelectBillFields(proCode, block.BillID, -1, "fc027")
		if fc027.ResultValue == "1" || fc027.ResultValue == "A" {
			for ii, field := range fields {
				if field.Code == "fc029" || field.Code == "fc031" {
					fields[ii] = SetOpInput(field, op)
				}
			}
		}
	}
	if RegIsMatch(block.Code, `^(bc010|bc011|bc014|bc015|bc016|bc017|bc018|bc019|bc020|bc021|bc022|bc023|bc024|bc025|bc026|bc027|bc028|bc029|bc030|bc031|bc032)$`) {
		err, bill := service.SelectBillByID(proCode, block.BillID)
		if err == nil && bill.SaleChannel == "秒赔" {
			for ii, field := range fields {
				if RegIsMatch(field.Code, `^(fc009|fc019|fc099|fc010|fc011|fc095|fc012|fc027|fc028|fc029|fc030|fc031|fc158|fc112|fc169|fc170|fc171|fc172|fc173|fc110|fc190|fc191|fc174|fc175)$`) {
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
