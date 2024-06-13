package response

import (
	"server/module/pro_conf/model"
)

type SysBlockRelationType struct {
	OneBlockRelation   []model.TempBlockRelation `json:"oneBlockRelation" gorm:"comment:'前置分块'"`
	TwoBlockRelation   []model.TempBlockRelation `json:"twoBlockRelation" gorm:"comment:'一码前置分块'"`
	ThreeBlockRelation []model.TempBlockRelation `json:"threeBlockRelation" gorm:"comment:'参考分块'"`
	TempBFRelation     []model.TempBFRelation    `json:"tempBFRelation" gorm:"comment:'分块字段关系'"`
}
