package model

import "server/module/sys_base/model"

type ConstLog struct {
	model.Model
	Type     string `json:"type" form:"type"`         //操作类型
	Content  string `json:"content" form:"content"`   //操作内容
	UserId   string `json:"userId" form:"userId"`     //操作人工号
	UserName string `json:"userName" form:"userName"` //操作人姓名
	Status   int    `json:"status" form:"status"`     //请求状态
}

type ConstLogReq struct {
	model.BasePageInfo
	model.BaseTimeRange
	Type    string `json:"type" form:"type"`       //操作类型
	Content string `json:"content" form:"content"` //操作内容
}
