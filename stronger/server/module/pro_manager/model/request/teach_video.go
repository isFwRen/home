package request

import modelBase "server/module/sys_base/model"

type TV struct {
	PageInfo  modelBase.BasePageInfo `json:"page_info"`
	ProCode   string                 `json:"proCode" form:"proCode"`
	BlockName string                 `json:"blockName" form:"blockName"`
	Rule      string                 `json:"rule" form:"rule"`
}

type ExportTV struct {
	ProCode string `json:"proCode" form:"proCode"`
}
