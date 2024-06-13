/**
 * @Author: 星期一
 * @Description:
 * @Date: 2021/1/5 上午11:20
 */

package request

import (
	"server/module/pro_conf/model"
	modelBase "server/module/sys_base/model"
)

//字段配置搜索
type SysProFieldsSearch struct {
	model.SysProField
	modelBase.BasePageInfo
}

//导出节点配置搜索
type SysExportNodesSearch struct {
	model.SysExportNode
	modelBase.PageInfoSearch
	ProId string `json:"proId" form:"proId" gorm:"comment:项目Id"`
}

//操作记录搜索
type SysOperationRecordSearch struct {
	modelBase.SysOperationRecord
	modelBase.BasePageInfo
}

// SysProjectRecordSearch 项目配置搜索
type SysProjectRecordSearch struct {
	model.SysProject
	modelBase.BasePageInfo
}

//质检配置搜索
type SearchSysQuality struct {
	model.SysQuality
	modelBase.BasePageInfo
}

//审核配置搜索
type SearchSysInspection struct {
	model.SysInspection
	modelBase.BasePageInfo
}
