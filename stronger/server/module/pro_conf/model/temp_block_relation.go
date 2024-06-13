package model

import (
	"server/module/sys_base/model"
)

type TempBlockRelation struct {
	model.Model
	MyType   int8   `json:"myType" gorm:"类型（前置分块，一码前置分块，参考分块）"`
	TempBId  string `json:"tempBId" gorm:"模板分块id"`
	PreBId   string `json:"preBId" gorm:"前置or参考分块id"`
	PreBName string `json:"preBName" gorm:"前置or参考分块名字"`
	PreBCode string `json:"preBCode" gorm:"前置or参考分块编码"`
}
