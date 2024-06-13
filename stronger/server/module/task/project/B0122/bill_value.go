package B0122

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc019" {

		err, fc010 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc010")
		if err == nil {
			data["fc010"] = fc010.ResultValue
		}
		return data
	}

	if block.Code == "bc017" || block.Code == "bc018" {

		err, fc486 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc486")
		if err == nil {
			data["fc486"] = fc486.ResultValue
		}
		err, fc135 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc135")
		if err == nil {
			data["fc135"] = fc135.ResultValue
		}
		return data
	}
	if block.Code == "bc063" {
		err, fc265 := service.SelectBillFields(proCode, block.BillID, -1, "fc265")
		if err == nil {
			data["fc265"] = fc265.ResultValue
		}
		return data
	}

	if block.Code == "bc025" || block.Code == "bc026" {
		err, fcCode := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc484")
		if err == nil {
			data["fc484"] = fcCode.ResultValue
		}
		return data
	}
	return data
}
