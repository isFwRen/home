/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/1/5 10:36 上午
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/global/response"
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
	response2 "server/module/sys_base/model/response"
	"server/module/sys_base/service"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// AddMenu
// @Tags SysMenu
// @Summary 人员管理--新增菜单或按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Produce application/json
// @Param data body model.SysMenu true "菜单实体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /sys-menu/add [post]
func AddMenu(c *gin.Context) {
	var reqParam model.SysMenu
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err = service.AddMenu(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("新增失败%s", err.Error()), c)
		return
	}
	response.OkWithMessage("新增成功", c)
}

// DeleteMenu
// @Tags SysMenu
// @Summary 人员管理--删除菜单或按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ReqIds true "id数组"
// @Success 200 {int64} int64 "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sys-menu/delete [delete]
func DeleteMenu(c *gin.Context) {
	var reqParam request.ReqIds
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	rows := service.DeleteMenuByIds(reqParam)
	if rows == 0 {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkDetailed(rows, "删除成功", c)
}

// GetMenuTreeList
// @Tags SysMenu
// @Summary 人员管理--获取菜单或按钮tree
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-menu/tree [get]
func GetMenuTreeList(c *gin.Context) {
	//err, t := ListToTreeHasSelect(nil)
	err, t := service.QueryAllMenu()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败%s", err.Error()), c)
		return
	}
	response.OkDetailed(t, "获取成功", c)
}

// GetMenuTreeAndProByToken
// @Tags SysMenu
// @Summary 人员管理--获取当前登录角色的menu
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-menu/role/get [get]
func GetMenuTreeAndProByToken(c *gin.Context) {
	customClaims, err1 := api.GetUserByToken(c)
	if err1 != nil {
		response.FailWithMessage(fmt.Sprintf("获取登录者失败，%v", err1), c)
	}
	//获取菜单树
	//err, t := getMenuTree(customClaims)
	//if err != nil {
	//	response.FailWithMessage(fmt.Sprintf("获取失败%s", err.Error()), c)
	//	return
	//}
	//global.GLog.Info("customClaims.RoleId:::" + customClaims.RoleId)
	//获取项目权限
	err, p := service.GetAllPermissionByUId(customClaims.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败%s", err.Error()), c)
		return
	}
	err = CachePerm(customClaims.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("权限获取失败,%s", err.Error()), c)
		return
	}
	response.OkDetailed(map[string]interface{}{
		//"Menus":   t,
		"Perm":    p,
		"proCode": global.GConfig.System.ProCode,
	}, "获取成功", c)
}

func CachePerm(uid string) error {
	err, perms := service.FetchPermissionByUId(uid)
	if err != nil {
		return err
	}
	var permMap = make(map[string]interface{})
	for _, perm := range perms {
		permMap[perm.ProCode+"_op0"] = perm.HasOp0
		permMap[perm.ProCode+"_op1"] = perm.HasOp1
		permMap[perm.ProCode+"_op2"] = perm.HasOp2
		permMap[perm.ProCode+"_opq"] = perm.HasOpq
		permMap[perm.ProCode+"_HasInNet"] = perm.HasInNet
		permMap[perm.ProCode+"_HasOutNet"] = perm.HasOutNet
		permMap[perm.ProCode+"_pm"] = perm.HasPm
	}
	if len(permMap) == 0 {
		return nil
	}
	err = service.SetRedisPerm(permMap, uid)
	return err
}

func getMenuTree(customClaims *request.CustomClaims) (error, []response2.SysMenuResponse) {
	err, r := service.GetRoleMenuRelationsByRoleId(customClaims.RoleId)
	if err != nil {
		return err, nil
	}
	var rMap = make(map[string]string)
	var rList []string
	for _, relation := range r {
		rMap[relation.MenuId] = relation.ID
		rList = append(rList, relation.MenuId)
	}
	err, t := ListToTreeWithRole(rMap, rList)
	return err, t
}

// EditMenu
// @Tags SysMenu
// @Summary 人员管理--修改菜单或按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysMenuForm true "菜单实体"
// @Success 200 {string} int64 "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /sys-menu/edit [post]
func EditMenu(c *gin.Context) {
	var reqParam request.SysMenuForm
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	if reqParam.ID == "" {
		response.FailWithMessage("修改失败,id不能为空", c)
		return
	}
	rows := service.EditMenuById(reqParam)
	if rows == 0 {
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkDetailed(rows, "修改成功", c)
}

// ListToTreeHasSelect menu列表转tree
func ListToTreeHasSelect(relationMap map[string]string) (error, []response2.SysMenuResponse) {
	err, list := service.QueryAllMenu()
	if err != nil {
		return err, make([]response2.SysMenuResponse, 0)
	}
	var treeMap = make(map[string][]response2.SysMenuResponse)
	for _, menu := range list {
		if relationMap != nil {
			menu.RoleMenuRelationId, menu.IsSelect = relationMap[menu.ID]
		}
		treeMap[menu.ParentId] = append(treeMap[menu.ParentId], menu)
	}
	var newTree = treeMap["0"]
	setChildNode(&newTree, treeMap)
	return nil, newTree
}

func setChildNode(list *[]response2.SysMenuResponse, treeMap map[string][]response2.SysMenuResponse) {
	for i, menu := range *list {
		v, ok := treeMap[menu.ID]
		if !ok {
			(*list)[i].Children = []response2.SysMenuResponse{}
			continue
		}
		(*list)[i].Children = v
		setChildNode(&(*list)[i].Children, treeMap)
	}
}

// ListToTreeWithRole menu列表转tree
func ListToTreeWithRole(relationMap map[string]string, rList []string) (error, []response2.SysMenuResponse) {
	err, list := service.QueryMenuByTypeAndRole(0, rList)
	if err != nil {
		return err, make([]response2.SysMenuResponse, 0)
	}
	var treeMap = make(map[string][]response2.SysMenuResponse)
	for _, menu := range list {
		if relationMap != nil {
			treeMap[menu.ParentId] = append(treeMap[menu.ParentId], menu)
		}
	}
	var newTree = treeMap["0"]
	setChildNode(&newTree, treeMap)
	return nil, newTree
}

// GetApis
// @Tags SysMenu
// @Summary 人员管理--获取所有api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} model.SysApi "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-menu/api/get [get]
func GetApis(c *gin.Context) {
	err, apis := service.GetApis()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败%s", err.Error()), c)
		return
	}
	response.OkDetailed(apis, "获取成功", c)
}
