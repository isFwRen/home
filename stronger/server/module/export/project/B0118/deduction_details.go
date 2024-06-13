package B0118

import (
	"fmt"
	"server/module/export/model"
	"server/module/export/service"
	lModel "server/module/load/model"
	"time"
)

func DeductionDetails(obj model.ResultDataBill, xmlValue string) error {
	// obj.Bill.ProCode
	fmt.Println("-----------------DeductionDetails-----------------------------")
	deductionDetails := []model.DeductionDetail{}
	for _, invoice := range obj.Invoice {
		deductionDetail := model.DeductionDetail{}
		deductionDetail.BillId = obj.Bill.ID
		deductionDetail.Date = time.Now()
		deductionDetail.ProCode = obj.Bill.ProCode
		deductionDetail.BillNum = obj.Bill.BillNum
		deductionDetail.BatchNum = obj.Bill.BatchNum
		deductionDetail.Agency = obj.Bill.Agency
		_, fc181 := GetOneField(invoice.Invoice[0], "fc181", false)
		_, fc180 := GetOneField(invoice.Invoice[0], "fc180", false)
		deductionDetail.Province = fc181
		deductionDetail.City = fc180
		_, fc002 := GetOneField(invoice.Invoice[0], "fc002", true)
		deductionDetail.BillCode = fc002
		_, fc054 := GetOneField(invoice.Invoice[0], "fc054", false)
		deductionDetail.Hospital = fc054

		_, fc003 := GetOneField(invoice.Invoice[0], "fc003", false)
		mapData := map[string]string{"1": "门诊", "2": "住院"}
		deductionDetail.CostType = mapData[fc003]

		_, fc205 := GetOneField(invoice.Invoice[0], "fc205", false)
		mapData = map[string]string{"1": "社保", "2": "新农合", "3": "公费医疗", "4": "其他保险公司", "5": "社会福利机构", "6": "其他政府机构", "7": "其他", "8": "补充医疗", "A": "无报销"}
		deductionDetail.ReimburseType = mapData[fc205]

		myCodes := [][]string{
			{"fc154", "fc162", "fc084", "fc092", "fc172"},
			{"fc155", "fc163", "fc085", "fc093", "fc173"},
			{"fc156", "fc164", "fc086", "fc094", "fc174"},
			{"fc157", "fc165", "fc087", "fc095", "fc175"},
			{"fc158", "fc166", "fc088", "fc096", "fc176"},
			{"fc159", "fc167", "fc089", "fc097", "fc177"},
			{"fc160", "fc168", "fc090", "fc098", "fc178"},
			{"fc161", "fc169", "fc091", "fc099", "fc179"},
		}

		// CachedeductionDetails := []model.DeductionDetail{}
		for _, qinDan := range invoice.QingDan {
			for _, myCode := range myCodes {
				_, c1 := GetOneField(qinDan, myCode[0], false)
				if c1 != "3" && c1 != "" {
					deductionDetail_clone := deductionDetail
					_, c1 := GetOneField(qinDan, myCode[0], true)
					mapData = map[string]string{"1": "有写自费/自付比例", "2": "丙（无扣费标准）", "4": "有先付、预付金额", "5": "类型标明乙类但无扣费标准", "6": "全自费"}
					deductionDetail_clone.ChargingType = mapData[c1]
					_, c2 := GetOneField(qinDan, myCode[1], true)
					deductionDetail_clone.SelfRatio = c2
					_, c3 := GetOneField(qinDan, myCode[2], true)
					deductionDetail_clone.InventoryName = c3
					_, c4 := GetOneField(qinDan, myCode[3], true)
					deductionDetail_clone.InventoryMoney = c4
					_, c5 := GetOneField(qinDan, myCode[4], true)
					deductionDetail_clone.SelfMoney = c5
					deductionDetails = append(deductionDetails, deductionDetail_clone)
				}
			}
		}
		// deductionDetails = append(deductionDetails, deductionDetail)
	}

	fmt.Println("-----------------DeductionDetails-----------------------------", len(deductionDetails))
	// if len(deductionDetails) > 0 {
	return service.InsertDeductionDetails(obj.Bill.ProCode, obj.Bill.ID, deductionDetails)
	// }
	// return nil

	// return service.UpdateAgency(obj.Bill.ProCode, agency)
}

func GetOneField(fields []lModel.ProjectField, code string, finalOrResult bool) (bool, string) {
	for _, field := range fields {
		if field.Code == code {
			if finalOrResult {
				return true, field.FinalValue
			} else {
				return true, field.ResultValue
			}
		}
	}
	return false, ""
}
