package request

type ProjectConfigAging struct {
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

type ProjectConfigAgingUpdate struct {
	ID                string `json:"id"`
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
	ProId      string `json:"proId"`
	ConfigType string `json:"configType"`
}
