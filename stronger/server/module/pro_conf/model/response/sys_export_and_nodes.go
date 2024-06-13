package response

import (
	"server/module/sys_base/model"
)

type SysExport struct {
	model.Model
	ProName  string          `json:"proName" gorm:"comment:项目名称"`
	TempVal  string          `json:"tempVal" gorm:"comment:xml模板数据"`
	XmlType  string          `json:"xmlType" gorm:"comment:xml编码"`
	NodeList []SysExportNode `json:"nodeList" gorm:"foreignKey:ExportId;comment:各个节点"`
}
