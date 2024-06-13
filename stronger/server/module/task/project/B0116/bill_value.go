package B0116

import (
	"server/module/load/model"
	"server/module/load/service"
)

func BillValue(proCode string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	data := make(map[string]interface{})
	if block.Code == "bc012" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc090 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc090")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		_, zFeilds := service.GetZeroFieldsValueByCode(proCode, fc090.BlockID, "fc088", fc090.ResultValue)
		for _, field := range zFeilds {
			if field.Code == "fc112" || field.Code == "fc111" {
				data[field.Code] = field.ResultValue
			}
		}
	}

	if block.Code == "bc010" {
		// fmt.Println("------------------------BO114-----------------------:", block.BillID, block.Zero, "fc113")
		_, fc111 := service.SelectBillFields(proCode, block.BillID, block.Zero, "fc111")
		// fmt.Println("------------------------fc113-----------------------:", fc113.ResultValue)
		data["fc111"] = fc111.ResultValue
	}

	if block.Code == "bc029" {
		_, fc013 := service.SelectBillFields(proCode, block.BillID, -1, "fc013")
		data["fc013"] = GetOpValue(fc013, op)
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
