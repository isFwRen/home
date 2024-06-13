/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/28 4:00 下午
 */

package model

type SysMenu struct {
	Model
	Name      string `json:"name" form:"name" gorm:"size:128;"`                    //菜单或按钮名字
	Title     string `json:"title" form:"title" gorm:"size:128;"`                  //菜单或按钮标题
	Icon      string `json:"icon" form:"icon" gorm:"size:128;"`                    //图标
	Path      string `json:"path" form:"path" gorm:"size:128;"`                    //路经(路由url)
	MenuType  int    `json:"menuType" form:"menuType" gorm:"size:1;"`              //菜单类型1：菜单，2：按钮
	Api       string `json:"api" form:"api" gorm:"size:255;"`                      //api
	ParentId  string `json:"parentId" form:"parentId" gorm:"DEFAULT:0"`            //上级id
	Component string `json:"component" form:"component" gorm:"size:255;"`          //组件
	Sort      int    `json:"sort" form:"sort" gorm:"size:4;DEFAULT:0;"`            //排序
	IsEnable  bool   `json:"isEnable" form:"isEnable" gorm:"size:1;DEFAULT:true;"` //是否可用
	IsFrame   bool   `json:"isFrame" form:"isFrame" gorm:"size:1;DEFAULT:false;"`  //是否是弹窗
	Action    string `json:"action" form:"action" gorm:"size:16;"`                 //请求类型
	ApiId     string `json:"apiId" form:"apiId" gorm:"size:20;"`                   //api的id
}

func (SysMenu) TableName() string {
	return "sys_menus"
}

type SysRoleMenuRelations struct {
	Model
	RoleId string `json:"roleId"` //角色id
	MenuId string `json:"menuId"` //菜单id
}
