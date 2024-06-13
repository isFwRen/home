package response

import (
	"server/module/sys_base/model"
)

//import "gorm.io/datatypes"

type SysExportNode struct {
	model.Model
	ExportId      string `json:"exportId" gorm:"comment:导出id"`
	Name          string `json:"name" gorm:"comment:节点名称"`
	OneFields     string `json:"oneFields" gorm:"comment:第一个字段"`
	TwoFields     string `json:"twoFields" gorm:"comment:第二个字段"`
	ThreeFields   string `json:"threeFields" gorm:"comment:第三个字段"`
	FixedValue    string `json:"fixedValue" gorm:"comment:固定值"`
	Remark        string `json:"remark" gorm:"comment:备注"`
	MyType        string `json:"myType" gorm:"comment:类型（开始，结束，录入值，固定值）"`
	MyOrder       int    `json:"myOrder" gorm:"comment:排序"`
	OneFieldsName string `json:"oneFieldsName"`
	//CodeOneFields   string `json:"codeOneFields"`
	TwoFieldsName string `json:"twoFieldsName"`
	//CodeTwoFields   string `json:"codeTwoFields"`
	ThreeFieldsName string `json:"threeFieldsName"`
	//CodeThreeFields string `json:"codeThreeFields"`
	//OneSysProField   SysProField `json:"oneSysProField" gorm:"ForeignKey:OneFields"`
	//TwoSysProField   SysProField `json:"twoSysProField" gorm:"ForeignKey:TwoFields"`
	//ThreeSysProField SysProField `json:"threeSysProField" gorm:"ForeignKey:ThreeFields"`
}
