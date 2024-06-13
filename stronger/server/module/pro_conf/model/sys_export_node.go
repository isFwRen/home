package model

import (
	"server/module/sys_base/model"
	"time"
)

//import "gorm.io/datatypes"

type SysExportNode struct {
	model.Model
	ExportId        string `json:"exportId" form:"exportId" gorm:"comment:导出id"`
	Name            string `json:"name" form:"name" excel:"节点名称" gorm:"comment:节点名称"`
	OneFields       string `json:"oneFields" form:"oneFields" excel:"第一个字段" gorm:"comment:第一个字段"`
	TwoFields       string `json:"twoFields" form:"twoFields" excel:"第二个字段" gorm:"comment:第二个字段"`
	ThreeFields     string `json:"threeFields" form:"threeFields" excel:"第三个字段" gorm:"comment:第三个字段"`
	FixedValue      string `json:"fixedValue" form:"fixedValue" excel:"固定值" gorm:"comment:固定值"`
	Remark          string `json:"remark" form:"remark" excel:"备注" gorm:"comment:备注"`
	MyType          string `json:"myType" form:"myType" excel:"类型" gorm:"comment:类型（开始，结束，录入值，固定值）"`
	MyOrder         int    `json:"myOrder" form:"myOrder" excel:"排序" gorm:"comment:排序"`
	OneFieldsName   string `json:"oneFieldsName" form:"oneFieldsName" excel:"第一个字段" gorm:"comment:第一个字段"`
	TwoFieldsName   string `json:"twoFieldsName" form:"twoFieldsName" excel:"第二个字段" gorm:"comment:第二个字段"`
	ThreeFieldsName string `json:"threeFieldsName" form:"threeFieldsName" excel:"第三个字段" gorm:"comment:第三个字段"`
}

type SysExportNodeExport struct {
	//ID              int `gorm:"primary_key"`
	//CreatedAt       time.Time
	UpdatedAt       time.Time
	ExportId        string `json:"exportId" excel:"导出id" gorm:"comment:导出id"`
	Name            string `json:"name" excel:"节点名称" gorm:"comment:节点名称"`
	OneFields       string `json:"oneFields" excel:"第一个字段id" gorm:"comment:第一个字段id"`
	TwoFields       string `json:"twoFields" excel:"第二个字段id" gorm:"comment:第二个字段id"`
	ThreeFields     string `json:"threeFields" excel:"第三个字段id" gorm:"comment:第三个字段id"`
	FixedValue      string `json:"fixedValue" excel:"固定值" gorm:"comment:固定值"`
	Remark          string `json:"remark" excel:"备注" gorm:"comment:备注"`
	MyType          string `json:"myType" excel:"类型" gorm:"comment:类型（开始，结束，录入值，固定值）"`
	MyOrder         int    `json:"myOrder" excel:"排序" gorm:"comment:排序"`
	NameOneFields   string `json:"nameOneFields" excel:"第一个字段名字" gorm:"comment:第一个字段名字"`
	CodeOneFields   string `json:"codeOneFields" excel:"第一个字段编码" gorm:"comment:第一个字段编码"`
	NameTwoFields   string `json:"nameTwoFields" excel:"第一个字段名字" gorm:"comment:第一个字段名字"`
	CodeTwoFields   string `json:"codeTwoFields" excel:"第一个字段编码" gorm:"comment:第一个字段编码"`
	NameThreeFields string `json:"nameThreeFields" excel:"第一个字段名字" gorm:"comment:第一个字段名字"`
	CodeThreeFields string `json:"codeThreeFields" excel:"第一个字段编码" gorm:"comment:第一个字段编码"`
}
