package response

import (
	"server/module/pro_conf/model"
)

type SysProTempB struct {
	model.SysProTempB
	OneRelations   []model.TempBlockRelation `json:"oneRelations" gorm:"foreignKey:TempBId;comment:前置分块"`
	TwoRelations   []model.TempBlockRelation `json:"twoRelations" gorm:"foreignKey:TempBId;comment:一码前置分块'"`
	ThreeRelations []model.TempBlockRelation `json:"threeRelations" gorm:"foreignKey:TempBId;comment:参考分块'"`
	TempBFRelation []model.TempBFRelation    `json:"tempBFRelation" gorm:"foreignKey:TempBId;comment:字段'"`
}
