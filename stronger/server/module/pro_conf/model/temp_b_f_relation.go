package model

import (
	"server/module/sys_base/model"
)

type TempBFRelation struct {
	model.Model
	TempBId string `json:"tempBId" gorm:"模板分块id"`
	FId     string `json:"fId" gorm:"字段id"`
	FName   string `json:"fName" gorm:"字段名字"`
	FCode   string `json:"fCode" gorm:"字段编码"`
}
