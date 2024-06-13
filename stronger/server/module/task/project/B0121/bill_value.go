package B0121

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc005" || block.Code == "bc016" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc364 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc364")
		_, fc365 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc365")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		data["fc364"] = fc364.ResultValue
		data["fc365"] = fc365.ResultValue
	}

	return data
}
