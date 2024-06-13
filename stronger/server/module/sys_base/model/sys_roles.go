package model

type SysRoles struct {
	Model
	Name      string `json:"name" gorm:"角色名称"`                //角色名称
	Status    int    `json:"status" gorm:"角色状态"`              //角色状态
	Remark    string `json:"remark" gorm:"角色描述"`              //角色描述
	CountUser int    `json:"countUser" gorm:"用户为当前角色名称的用户数量"` //用户为当前角色名称的用户数量
	CreatedBy string `json:"createdBy" gorm:"创建人"`            //创建人
	UpdatedBy string `json:"updatedBy" gorm:"更新人"`            //更新人
}

//type SysRoles2 struct {
//	Id      string         `json:"id" gorm:"comment:'id'"`
//	Name    string         `json:"name"  gorm:"comment:'权限名字'"`
//	Project string         `json:"project" gorm:"comment:'项目名字'"`
//	Role    pq.StringArray `json:"role" gorm:"type:varchar(100)[] comment:'权限数组'"`
//	ProCode string         `json:"proCode" gorm:"comment:'入职日期'"`
//	UserId  string         `json:"userId" gorm:"comment:'用户id'"`
//}
