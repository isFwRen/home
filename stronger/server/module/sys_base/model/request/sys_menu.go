/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/1/5 11:15 上午
 */

package request

import "server/module/sys_base/model"

type SysMenuForm struct {
	model.Model
	Name      string `json:"name" form:"name" gorm:"size:128;"`                           //菜单或按钮名字
	Title     string `json:"title" form:"title" gorm:"size:128;binding:'required'"`       //菜单或按钮标题
	Icon      string `json:"icon" form:"icon" gorm:"size:128;"`                           //图标
	Path      string `json:"path" form:"path" gorm:"size:128;"`                           //路经(路由url)
	MenuType  int    `json:"menuType" form:"menuType" gorm:"size:1;binding:'required'"`   //菜单类型1：菜单，2：按钮
	ParentId  string `json:"parentId" form:"parentId" gorm:"DEFAULT:0"`                   //上级id
	Api       string `json:"api" form:"api" gorm:"size:255;"`                             //api
	Component string `json:"component" form:"component" gorm:"size:255;"`                 //组件
	Sort      int    `json:"sort" form:"sort" gorm:"size:4;DEFAULT:0;binding:'required'"` //排序
	IsEnable  bool   `json:"isEnable" form:"isEnable" gorm:"size:1;DEFAULT:true;"`        //是否可用
	IsFrame   bool   `json:"isFrame" form:"isFrame" gorm:"size:1;DEFAULT:false;"`         //是否是弹窗
	Action    string `json:"action" form:"action" gorm:"size:16;"`                        //请求类型
	ApiId     string `json:"apiId" form:"apiId" gorm:"size:20;"`                          //api的id
}
