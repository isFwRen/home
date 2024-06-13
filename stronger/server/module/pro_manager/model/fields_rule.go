package model

import (
	"github.com/lib/pq"
	modelbase "server/module/sys_base/model"
)

type SysFieldRule struct {
	modelbase.Model
	ProId          string         `json:"proId" excel:"序号"`
	SysFieldConfId string         `json:"sysFieldConfId"`
	SysFieldCode   string         `json:"sysFieldCode" excel:"字段编码"`
	SysFieldName   string         `json:"sysFieldName" excel:"字段名字"`
	InputRule      string         `json:"inputRule" excel:"录入规则"`
	RulePicture    pq.StringArray `json:"rulePicture" gorm:"type:varchar(255)[]" excel:"规则图片"`
}
