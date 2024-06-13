package request

import "server/module/sys_base/model"

type AddRoleInformation struct {
	RoleName      string `json:"roleName" form:"roleName"`     //角色名称
	RoleStatus    int    `json:"roleStatus" form:"roleStatus"` //角色状态
	RoleDescribe  string `json:"roleRemark" form:"roleRemark"` //角色描述
	RoleCreatedBy string `json:"createdBy" form:"createdBy"`   //创建人
}

type EditRoleInformation struct {
	RoleBeforeName string `json:"beforeName" form:"beforeName"` //角色旧名称
	RoleNewName    string `json:"newName" form:"newName"`       //角色新名称
	RoleStatus     int    `json:"roleStatus" form:"roleStatus"` //角色状态
	RoleDescribe   string `json:"roleRemark" form:"roleRemark"` //角色描述
	RoleUpdatedBy  string `json:"updatedBy" form:"updatedBy"`   //更新人
}

type DeleteRoleInformation struct {
	Ids []string `json:"ids" form:"ids"` //角色名称
}

type RoleAndMenuRelationFormList struct {
	List []RoleAndMenuRelationForm `json:"list" form:"list"`
}

type RoleAndMenuRelationForm struct {
	model.Model
	RoleId   string `json:"roleId" form:"roleId"`              //角色id
	MenuId   string `json:"menuId" form:"menuId"`              //菜单id
	IsSelect bool   `json:"isSelect" form:"isSelect" gorm:"-"` //是否拥有
}
