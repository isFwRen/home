package model

type CasbinModel struct {
	ID     uint   `json:"id" gorm:"column:_id"`
	PType  string `json:"pType" gorm:"column:p_type"`
	RoleId string `json:"roleId" gorm:"column:v0"`
	Path   string `json:"path" gorm:"column:v1"`
	Method string `json:"method" gorm:"column:v2"`
}
