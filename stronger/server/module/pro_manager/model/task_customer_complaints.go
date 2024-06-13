package model

import modelBase "server/module/sys_base/model"

type CustomerComplaints struct {
	PageInfo       modelBase.BasePageInfo `json:"page_info"`
	Month          string                 `json:"month" form:"month"`
	ProCode        string                 `json:"proCode" form:"proCode"`
	BillName       string                 `json:"billName" form:"billName"`
	WrongFieldName string                 `json:"wrongFieldName" form:"wrongFieldName"`
}
