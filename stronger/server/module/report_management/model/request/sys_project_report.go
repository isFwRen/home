package request

import (
	modelBase "server/module/sys_base/model"
	"time"
)

type BusinessDetailsSearch struct {
	ProCode   string `json:"proCode" form:"proCode"`     //项目编码
	StartTime string `json:"startTime" form:"startTime"` //开始时间
	EndTime   string `json:"endTime" form:"endTime"`     //结束时间
	Type      int    `json:"type" form:"type"`           //类型 "1：日，2：周，3：月，4：年"
	PageInfo  modelBase.BasePageInfo
}

type ExportBusinessDetailsSearch struct {
	ProCode   string `json:"proCode" form:"proCode"`     //项目编码
	StartTime string `json:"startTime" form:"startTime"` //开始时间
	EndTime   string `json:"endTime" form:"endTime"`     //结束时间
}

type CopyBusinessDetailsSearch struct {
	ProCode   string                 `json:"proCode"`   //项目编码
	StartTime time.Time              `json:"startTime"` //开始时间
	EndTime   time.Time              `json:"endTime"`   //结束时间
	PageInfo  modelBase.BasePageInfo `json:"page_info"`
}
