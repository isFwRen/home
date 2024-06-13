package model

import (
	"server/module/sys_base/model"
	"time"
)

type DeductionDetail struct {
	model.Model
	ProCode        string    `json:"proCode" form:"proCode"`
	Date           time.Time `json:"date" form:"date"  excel:"日期"`
	BillNum        string    `json:"billNum" form:"billNum"  excel:"账单号"`
	BillId         string    `json:"billId" form:"billId"`
	BatchNum       string    `json:"batchNum" from:"batchNum"` //批次号
	Agency         string    `json:"agency" form:"agency"  excel:"机构号"`
	Province       string    `json:"province" form:"province"  excel:"省份"`
	City           string    `json:"city" form:"city"  excel:"城市"`
	BillCode       string    `json:"billCode" form:"billCode"  excel:"账单号"`
	Hospital       string    `json:"hospital" form:"hospital"  excel:"治疗医院"`
	CostType       string    `json:"costType" form:"costType"  excel:"费用类型"`
	ReimburseType  string    `json:"reimburseType" form:"reimburseType"  excel:"报销类型"`
	InventoryName  string    `json:"inventoryName" form:"inventoryName"  excel:"清单名称"`
	ChargingType   string    `json:"chargingType" form:"chargingType"  excel:"扣费类型"`
	InventoryMoney string    `json:"inventoryMoney" form:"inventoryMoney"  excel:"清单金额"`
	SelfRatio      string    `json:"selfRatio" form:"selfRatio"  excel:"自付比例"`
	SelfMoney      string    `json:"selfMoney" form:"selfMoney"  excel:"自付金额"`
}
