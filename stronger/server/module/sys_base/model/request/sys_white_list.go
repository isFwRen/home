package request

import (
	"github.com/lib/pq"
	modelBase "server/module/sys_base/model"
)

type GetWhiteList struct {
	PageInfo modelBase.BasePageInfo `json:"page_info"`
	ProCode  string                 `json:"proCode" form:"proCode" ` //项目编码
	//ProName  string                 `json:"proName" form:"proName"`   //项目名称
	TempName string `json:"tempName" form:"tempName"` //模板名称
	UserCode string `json:"code" form:"code"`         //工号
}

type ExportWhileList struct {
	ProCode   string        `json:"proCode" form:"proCode" `    //项目编码
	TempName  string        `json:"tempName" form:"tempName"`   //模板名称
	BlockCode []interface{} `json:"blockCode" form:"blockCode"` //分块编码
	BlockName []interface{} `json:"blockName" form:"blockName"` //分块名称
}

type CopyWhileList struct {
	ProCode  string   `json:"proCode" form:"proCode"`   //项目编码
	TempName string   `json:"tempName" form:"tempName"` //模板名称
	Code     string   `json:"code" form:"code"`         //工号
	CopyCode []string `json:"copyCode" form:"copyCode"` //复制工号
}

type EditWhiteListArr struct {
	ProCode string `json:"proCode"` //项目编码
	//ProName  string          `json:"proName"`  //项目名称
	TempName string          `json:"tempName"` //模板名称
	Arr      []EditWhiteList `json:"staffs" form:"staffs"`
}

type EditWhiteList struct {
	UserCode         string         `json:"userCode"`         //工号
	UserName         string         `json:"userName"`         //姓名
	BlockPermissions pq.StringArray `json:"blockPermissions"` //分块权限  存编码
}
