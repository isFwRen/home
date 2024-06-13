package B0113

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc005" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc331 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc331")
		_, fc332 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc332")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		data["fc332"] = fc332.ResultValue
		data["fc331"] = fc331.ResultValue
	}

	return data
}
