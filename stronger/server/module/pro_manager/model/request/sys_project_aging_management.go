package request

import (
	modelBase "server/module/sys_base/model"
)

type ProjectAgingManagementSearch struct {
	ProCode    string                 `json:"proCode" form:"proCode" example:"项目编码"`
	StartTime  string                 `json:"startTime" form:"startTime" example:"开始时间"`
	EndTime    string                 `json:"endTime" form:"endTime" example:"结束时间"`
	CaseNumber string                 `json:"caseNumber" form:"caseNumber" example:"案件号"`
	Agency     string                 `json:"agency" form:"agency" example:"机构号"`
	CaseStatus string                 `json:"caseStatus" form:"caseStatus" example:"案件状态"`
	Stage      int                    `json:"stage" form:"stage" example:"录入状态"`
	PageInfo   modelBase.BasePageInfo `json:"page_info"`
}

type ProjectAgingManagementSearchAll struct {
	ProCode    string                 `json:"proCode" form:"proCode" example:"项目编码"`
	StartTime  string                 `json:"startTime" form:"startTime" example:"开始时间"`
	EndTime    string                 `json:"endTime" form:"endTime" example:"结束时间"`
	CaseNumber string                 `json:"caseNumber" form:"caseNumber" example:"案件号"`
	Agency     string                 `json:"agency" form:"agency" example:"机构号"`
	CaseStatus string                 `json:"caseStatus" form:"caseStatus" example:"案件状态"`
	Stage      int                    `json:"stage" form:"stage" example:"录入状态"`
	PageInfo   modelBase.BasePageInfo `json:"page_info"`
}
