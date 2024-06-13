package request

import modelBase "server/module/sys_base/model"

type GetTaskListDetail struct {
	PageInfo         modelBase.BasePageInfo `json:"page_info"`
	ProCode          string                 `json:"proCode" form:"proCode"`
	Op               string                 `json:"op" form:"op"`
	OpStage          int                    `json:"opStage" form:"opStage"`
	IsExpenseAccount int                    `json:"isExpenseAccount" form:"isExpenseAccount"`
}

type GetTaskList struct {
	ProCode string `json:"proCode" form:"proCode"`
}

type GetVariousStateBill struct {
	ProCode    string `json:"proCode" `
	StickLevel int    `json:"stickLevel"`
}

type GetVariousStateBillNum struct {
	BillNum []string `json:"billNum"`
}

type SetVariousStateBill struct {
	ProCode     string   `json:"proCode"`
	StickLevel  int      `json:"stickLevel"`
	CaseNumbers []string `json:"caseNumbers"`
}

type SetPriorityOrganizationNumber struct {
	ProCode            string   `json:"proCode"`
	StickLevel         int      `json:"stickLevel"`
	OrganizationNumber []string `json:"organizationNumber"`
}

type GetCaseDetails struct {
	ProCode     string `json:"proCode" form:"proCode"`
	BillNum     string `json:"billNum" form:"billNum"`
	SaleChannel string `json:"saleChannel" form:"saleChannel"` //销售渠道
}

type GetCaseDetailsBlock struct {
	ProCode string `json:"proCode" form:"proCode"`
	Id      string `json:"id" form:"id"`
}

type GetCaseDetailsField struct {
	ProCode string `json:"proCode" form:"proCode"`
	BillId  string `json:"billId" form:"billId"`
	BlockId string `json:"blockId" form:"blockId"`
}

type Search struct {
	PageInfo modelBase.BasePageInfo `json:"page_info"`
	ProCode  string                 `json:"proCode" form:"proCode"`
	Op       string                 `json:"op" form:"op"`
}
