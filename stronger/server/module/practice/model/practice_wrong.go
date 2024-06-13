package model

import (
	"server/module/sys_base/model"
	"time"
)

type PracticeWrong struct {
	//Id             string         `json:"id"`
	model.Model
	SubmitDay time.Time `json:"submitDay" example:"提交时间" flag:"time"`
	Code      string    `json:"code" example:"工号"`
	NickName  string    `json:"nickName" example:"姓名"`
	BillNum   string    `json:"billNum" example:"单号"`
	FieldCode string    `json:"fieldCode" example:"字段编码"`
	FieldName string    `json:"fieldName" example:"字段名称"`
	Path      string    `json:"path" example:"单子图片存放路径"`
	Picture   string    `json:"picture" gorm:"type:varchar(100)[];comment:'web截图位置'" example:"图片名称"`
	Wrong     string    `json:"wrong" example:"错误内容"`
	Right     string    `json:"right" example:"正确内容"`
}
