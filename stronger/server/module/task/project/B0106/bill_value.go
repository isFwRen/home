package B0106

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc001" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc090 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc090")
		_, fc091 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc091")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		data["fc090"] = fc090.ResultValue
		data["fc091"] = fc091.ResultValue
	}

	return data
}
