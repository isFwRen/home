package model

import (
	"database/sql/driver"
	"encoding/json"
	"server/module/sys_base/model"

	"github.com/lib/pq"
)

type SysProTemplate struct {
	model.Model
	Name         string         `json:"name" gorm:"comment:'模板名字'"`
	Code         string         `json:"code" gorm:"comment:'模板编码'"`
	Images       pq.StringArray `json:"images" gorm:"type:varchar(255)[] comment:模板图片'"`
	ProName      string         `json:"proName" gorm:"comment:'项目名字'"`
	ProId        string         `json:"proId" gorm:"comment:'项目id'"`
	ProCode      string         `json:"proCode" gorm:"comment:'项目代码'"`
	SysProTempBs []SysProTempB  `json:"sysProTempBs" form:"sysProTempBs" gorm:"foreignKey:ProTempId;references:id;comment:'模板配置'"`
}

type ipList []string

// Value
// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p ipList) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *ipList) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}
