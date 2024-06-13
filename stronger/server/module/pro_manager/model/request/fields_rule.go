package request

import modelBase "server/module/sys_base/model"

type GFR struct {
	PageInfo   modelBase.BasePageInfo `json:"page_info"`
	ProCode    string                 `json:"proCode" form:"proCode"`
	FieldsName string                 `json:"fieldsName" form:"fieldsName"`
	Rule       string                 `json:"rule" form:"rule"`
}

type ExportFR struct {
	ProCode string `json:"proCode" form:"proCode"`
	Rule    string `json:"rule" form:"rule"`
}
