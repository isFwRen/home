package b0118

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc002" {

		_, fc060 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc060")
		_, zFeilds := service.GetZeroFieldsValue(proCode, fc060.BlockID, fc060.ResultValue)
		for _, field := range zFeilds {
			if field.Code == "fc180" || field.Code == "fc181" {
				data[field.Code] = field.ResultValue
			}
		}
		return data
	}
	if block.Code == "bc003" {
		_, fc061 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc061")
		_, zFeilds := service.GetZeroFieldsValue(proCode, fc061.BlockID, fc061.ResultValue)
		for _, field := range zFeilds {
			if field.Code == "fc180" {
				data[field.Code] = field.ResultValue
			}
		}
		return data
	}
	return data
}
