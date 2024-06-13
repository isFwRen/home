package b0114

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc010" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc113 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc113")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		data["fc113"] = fc113.ResultValue
		_, fc008 := service.SelectBillFields(proCode, block.BillID, -1, "fc008")
		_, fc003 := service.SelectBillFields(proCode, block.BillID, -1, "fc003")
		_, fc004 := service.SelectBillFields(proCode, block.BillID, -1, "fc004")
		data["fc008"] = fc008.ResultValue
		data["fc003"] = fc003.ResultValue
		data["fc004"] = fc004.ResultValue
		// _, zFeilds := service.GetZeroFieldsValue(proCode, fc060.BlockID, fc060.ResultValue)
		// for _, field := range zFeilds {
		// 	if field.Code == "fc180" || field.Code == "fc181" {
		// 		data[field.Code] = field.ResultValue
		// 	}
		// }
	}

	return data
}
