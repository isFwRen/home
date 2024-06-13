/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/1/5 11:15 上午
 */

package response

import "server/module/sys_base/model"

type SysMenuResponse struct {
	model.Model
	Name               string            `json:"name" gorm:"size:128;"`                           //菜单或按钮名字
	Title              string            `json:"title" gorm:"size:128;binding:'required'"`        //菜单或按钮标题
	Icon               string            `json:"icon" gorm:"size:128;"`                           //图标
	Path               string            `json:"path" gorm:"size:128;"`                           //路经(路由url)
	MenuType           int               `json:"menuType" gorm:"size:1;binding:'required'"`       //菜单类型0：菜单，1：按钮
	ParentId           string            `json:"parentId" gorm:"size:11;DEFAULT:0"`               //上级id
	Api                string            `json:"api" gorm:"size:255;"`                            //api
	Component          string            `json:"component" gorm:"size:255;"`                      //组件
	Sort               int               `json:"sort" gorm:"size:4;DEFAULT:0;binding:'required'"` //排序
	IsEnable           bool              `json:"isEnable" gorm:"size:1;DEFAULT:true;"`            //是否可用
	IsFrame            bool              `json:"isFrame" gorm:"size:1;DEFAULT:false;"`            //是否是弹窗
	Action             string            `json:"action" gorm:"size:16;"`                          //请求类型
	ApiId              string            `json:"apiId" gorm:"size:20;"`                           //api的id
	Children           []SysMenuResponse `json:"children" gorm:"-"`                               //子节点
	IsSelect           bool              `json:"isSelect" gorm:"-"`                               //有权限
	RoleMenuRelationId string            `json:"roleMenuRelationId" gorm:"-"`                     //角色菜单关系id
}
