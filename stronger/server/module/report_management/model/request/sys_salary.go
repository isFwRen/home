package request

import modelBase "server/module/sys_base/model"

type SysSalarySearch struct {
	PageInfo modelBase.BasePageInfo `json:"page_info"`
	Ym       string                 `json:"ym" form:"ym"`
	Start    string                 `json:"start" form:"start"`
	End      string                 `json:"end" form:"end"`
	Code     string                 `json:"code" form:"code"`
	Name     string                 `json:"name" form:"name"`
}
