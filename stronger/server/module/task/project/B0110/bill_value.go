package B0110

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc001" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc066 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc066")
		_, fc067 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc067")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		data["fc066"] = fc066.ResultValue
		data["fc067"] = fc067.ResultValue
	}

	return data
}
