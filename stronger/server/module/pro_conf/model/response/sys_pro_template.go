package response

import (
	"github.com/lib/pq"
	model2 "server/module/pro_conf/model"
	"server/module/sys_base/model"
)

type SysProTemplate struct {
	model.Model
	Name   string         `json:"name" gorm:"comment:'模板名字'"`
	Code   string         `json:"code" gorm:"comment:'模板编码'"`
	Images pq.StringArray `json:"images" gorm:"type:varchar(255)[] comment:模板图片'"`
	ProId  string         `json:"proId" gorm:"comment:'项目id'"`
	//SysProTempB []model.SysProTempB `json:"sysProTempB" gorm:"foreignKey:ProTempId;comment:模板分块"`
}

type NewIdIndexBlock struct {
	NewId         string                     `json:"newId" gorm:"comment:'新的blockID'"`
	OldIndex      []int                      `json:"oldIndex" gorm:"comment:'关系block的Index'"`
	BlockRelaList []model2.TempBlockRelation `json:"blockRelaList" gorm:"comment:旧的Block对象'"`
}
