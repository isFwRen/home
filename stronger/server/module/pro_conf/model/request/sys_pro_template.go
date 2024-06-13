package request

type SysProTemplate struct {
	Name  string `json:"name" from:"name" gorm:"comment:'模板名字'"`
	ProId string `json:"proId" from:"proId" gorm:"comment:'项目id'"`
}
