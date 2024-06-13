package model

import (
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
	"time"
)

type BusinessPushSend struct {
	model.Model
	PushId       string       `json:"pushId" form:"pushId" db:"push_id"`       // 推送记录表
	UserId       string       `json:"userId" form:"userId" db:"user_id"`       // 推送到的人，用户id
	IsPush       bool         `json:"isPush" form:"isPush" db:"is_push"`       // 是否推送了
	IsRead       bool         `json:"isRead" form:"isRead" db:"isRead"`        // '阅读状态（f未读，t已读）',
	ReadTime     time.Time    `json:"readTime" form:"readTime" db:"read_time"` // 阅读时间
	BusinessPush BusinessPush `gorm:"foreignKey:PushId"`
}

type BusinessPush struct {
	model.Model
	Title   string `json:"title" form:"title" db:"title"`        // 标题
	Msg     string `json:"msg" form:"msg" db:"msg"`              // 内容
	Type    int    `json:"type" form:"type" db:"type"`           // 消息类型1:下载2:上传
	ProCode string `json:"proCode" form:"proCode" db:"pro_code"` // 项目编码
}

// BusinessPushSearchReq 分页请求参数
type BusinessPushSearchReq struct {
	model.BaseTimePageCode
	MsgType int `json:"msgType" form:"msgType"` //消息类型
}

// BusinessPushSendReadReq 已读请求参数
type BusinessPushSendReadReq struct {
	request.ReqIds
	ProCode string `json:"proCode" form:"proCode"` //项目编码
}
