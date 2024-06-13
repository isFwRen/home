package request

import (
	modelBase "server/module/sys_base/model"
)

type AbnormalBillSearch struct {
	PageInfo  modelBase.BasePageInfo `json:"page_info"`
	ProCode   string                 `json:"proCode" form:"proCode"`     //项目编码
	StartTime string                 `json:"startTime" form:"startTime"` //开始时间
	EndTime   string                 `json:"endTime" form:"endTime"`     //结束时间
	Types     string                 `json:"types" form:"types"`         //类型
}
