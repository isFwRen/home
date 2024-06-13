package model

import "server/module/sys_base/model"

type SysFtpMonitor struct {
	model.Model
	Frequency   int    `json:"frequency" form:"frequency"`     //频率
	WrongMsg    string `json:"wrongMsg" form:"wrongMsg"`       //异常描述
	Desc        string `json:"desc" form:"desc"`               //下载说明
	CreatedCode string `json:"createdCode" form:"createdCode"` //创建人工号
	CreatedName string `json:"createdName" form:"createdName"` //创建人姓名
	ProCode     string `json:"proCode" form:"proCode"`         //项目编码
}

type SysFtpMonitorInfoReq struct {
	ProCode string `json:"proCode" form:"proCode" binding:"required"` //项目编码
}

type SysFtpMonitorEditReq struct {
	ID          string `json:"id" form:"id"`
	Frequency   int    `json:"frequency" form:"frequency" binding:"required"` //频率
	WrongMsg    string `json:"wrongMsg" form:"wrongMsg" binding:"required"`   //异常描述
	Desc        string `json:"desc" form:"desc" binding:"required"`           //下载说明
	CreatedName string `json:"createdName" form:"createdName"`                //创建人姓名
	CreatedCode string `json:"createdCode" form:"createdCode"`                //创建人工号
	ProCode     string `json:"proCode" form:"proCode" binding:"required"`     //项目编码
}
