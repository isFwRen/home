/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/1/5 11:29 上午
 */

package service

import (
	"server/global"
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
	"server/module/sys_base/model/response"
)

// AddMenu 新增菜单或按钮
func AddMenu(menu model.SysMenu) (err error) {
	err = global.GDb.Model(&model.SysMenu{}).Create(&menu).Error
	return err
}

// DeleteMenuByIds 删除菜单或按钮
func DeleteMenuByIds(ids request.ReqIds) int64 {
	return global.GDb.Model(&model.SysMenu{}).
		Where("id in (?)", ids.Ids).
		Delete(&model.SysMenu{}).RowsAffected
}

// QueryAllMenu 查询所有菜单或按钮
func QueryAllMenu() (err error, menus []response.SysMenuResponse) {
	err = global.GDb.Model(&model.SysMenu{}).Where("is_enable = true").Find(&menus).Error
	return err, menus
}

// EditMenuById 修改菜单或按钮
func EditMenuById(menu request.SysMenuForm) int64 {
	return global.GDb.Model(&model.SysMenu{}).Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"name":       menu.Name,
			"title":      menu.Title,
			"icon":       menu.Icon,
			"path":       menu.Path,
			"menu_type":  menu.MenuType,
			"parent_id":  menu.ParentId,
			"api":        menu.Api,
			"component":  menu.Component,
			"sort":       menu.Sort,
			"is_enable":  menu.IsEnable,
			"is_frame":   menu.IsFrame,
			"action":     menu.Action,
			"api_id":     menu.ApiId,
			"updated_at": menu.UpdatedAt,
		}).RowsAffected
}

// GetMenuById 根据id获取菜单或按钮
func GetMenuById(id string) (err error, menu response.SysMenuResponse) {
	err = global.GDb.Model(&model.SysMenu{}).Where("id = ?", id).First(&menu).Error
	return err, menu
}

// GetMenuByParentId 根据上级id获取菜单或按钮
func GetMenuByParentId(id string) (err error, menus []response.SysMenuResponse) {
	err = global.GDb.Model(&model.SysMenu{}).Where("parent_id = ?", id).Find(&menus).Error
	return err, menus
}

// QueryMenuByTypeAndRole 根据类型查询菜单或按钮
func QueryMenuByTypeAndRole(menuType int, ids []string) (err error, menus []response.SysMenuResponse) {
	err = global.GDb.Model(&model.SysMenu{}).Where("is_enable = true and id in (?)", ids).Find(&menus).Error
	return err, menus
}

// GetApis 获取所有api
func GetApis() (err error, apis []model.SysApi) {
	err = global.GDb.Model(&model.SysApi{}).Find(&apis).Error
	return err, apis
}
