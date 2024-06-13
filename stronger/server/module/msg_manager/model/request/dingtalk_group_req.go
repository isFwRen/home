package model

import "server/module/sys_base/model"

type DingtalkGroupReq struct {
	model.BasePageInfo
	Name    string `json:"name" form:"name"`
	ProCode string `json:"proCode" form:"proCode"`
	Env     int    `json:"env" form:"env"`
}
