package model

import (
	"server/module/sys_base/model"
)

type SysExport struct {
	model.Model
	ProId          string          `json:"proId"  gorm:"comment:项目id"`
	ProName        string          `json:"proName" gorm:"comment:项目名称"`
	TempVal        string          `json:"tempVal" gorm:"comment:xml模板数据"`
	XmlType        string          `json:"xmlType" gorm:"comment:xml编码"`
	SysExportNodes []SysExportNode `json:"sysExportNodes" form:"sysExportNodes" gorm:"foreignKey:ExportId;references:id;comment:'导出配置'"`
}

//func (SysExport) TableName() string {
//	return "sys_exports"
//}
