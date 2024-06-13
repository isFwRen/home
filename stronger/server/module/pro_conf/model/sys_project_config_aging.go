package model

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
	"time"
)

type SysProjectConfigAgingReq2 struct {
	model.Model
	ProId             string    `json:"proId" gorm:"comment:'项目名称'"`
	AgingStartTime    time.Time `json:"agingStartTime" gorm:"comment:'时效开始时间'"`
	AgingEndTime      time.Time `json:"agingEndTime" gorm:"comment:'时效结束时间'"`
	AgingOutStartTime time.Time `json:"agingOutStartTime" gorm:"comment:'时效外开始时间'"`
	AgingOutEndTime   time.Time `json:"agingOutEndTime" gorm:"comment:'时效外开始时间'"`
	RequirementsTime  string    `json:"requirementsTime" gorm:"comment:'考核要求(min)'"`
	ConfigType        string    `json:"configType" gorm:"comment:'配置类型'"`
	NodeName          string    `json:"nodeName" gorm:"comment:'节点名称'"`
	NodeContent       string    `json:"nodeContent" gorm:"comment:'节点内容'"`
	FieldName         string    `json:"fieldName" gorm:"comment:'字段名称'"`
	FieldContent      string    `json:"fieldContent" gorm:"comment:'字段内容'"`
}

type SysProjectConfigAgingReq struct {
	model.Model
	ProId             string `json:"proId" gorm:"comment:'项目名称'"`
	AgingStartTime    string `json:"agingStartTime" gorm:"comment:'时效开始时间'"`
	AgingEndTime      string `json:"agingEndTime" gorm:"comment:'时效结束时间'"`
	AgingOutStartTime string `json:"agingOutStartTime" gorm:"comment:'时效外开始时间'"`
	AgingOutEndTime   string `json:"agingOutEndTime" gorm:"comment:'时效外开始时间'"`
	RequirementsTime  string `json:"requirementsTime" gorm:"comment:'考核要求(min)'"`
	ConfigType        string `json:"configType" gorm:"comment:'配置类型'"`
	NodeName          string `json:"nodeName" gorm:"comment:'节点名称'"`
	NodeContent       string `json:"nodeContent" gorm:"comment:'节点内容'"`
	FieldName         string `json:"fieldName" gorm:"comment:'字段名称'"`
	FieldContent      string `json:"fieldContent" gorm:"comment:'字段内容'"`
}

type ReqProjectConfigAgingWithConfigType struct {
	ConfigType string `json:"configType" form:"configType"`
}

type SysProjectConfigAging struct {
	model.Model
	ProId          string `json:"proId" gorm:"comment:'项目名称'"`
	AgingStartTime string `json:"agingStartTime" gorm:"comment:'时效开始时间'"`
	AgingEndTime   string `json:"agingEndTime" gorm:"comment:'时效结束时间'"`

	//AgingOutStartTime string `json:"agingOutStartTime" gorm:"comment:'时效外开始时间'"`
	//AgingOutEndTime   string `json:"agingOutEndTime" gorm:"comment:'时效外开始时间'"`
	//WCoordinate pq.StringArray `json:"wCoordinate" gorm:"type:varchar(100)[];comment:'web截图位置'"` //web截图位置
	AgingOut pq.StringArray `json:"agingOut" gorm:"type:varchar(100)[];comment:'存时效外数据'"`

	RequirementsTime string `json:"requirementsTime" gorm:"comment:'考核要求(min)'"`
	ConfigType       string `json:"configType" gorm:"comment:'配置类型'"`
	NodeName         string `json:"nodeName" gorm:"comment:'节点名称'"`
	NodeContent      string `json:"nodeContent" gorm:"comment:'节点内容'"`
	FieldName        string `json:"fieldName" gorm:"comment:'字段名称'"`
	FieldContent     string `json:"fieldContent" gorm:"comment:'字段内容'"`
}
