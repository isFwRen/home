package b0108

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc018" || block.Code == "bc027" {

		err, fc110 := service.SelectBillFields(proCode, block.BillID, -1, "fc110")
		if err == nil {
			data["fc110"] = GetOpValue(fc110, op)
		}
		return data
	}
	return data
}

func GetOpValue(field model.ProjectField, op string) string {
	if op == "op1" {
		return field.Op1Value
	} else if op == "op2" {
		return field.Op2Value
	}
	return field.ResultValue

}
