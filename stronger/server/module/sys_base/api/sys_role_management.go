package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/global/response"
	"server/module/sys_base/model/request"
	res "server/module/sys_base/model/response"
	server "server/module/sys_base/service"
	"server/utils"
)

// GetRoleInformation
// @Tags Role Management (角色管理)
// @Summary	人员管理--查询角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param status query string true "状态 1:正常 2:停用"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/role-management/list [get]
func GetRoleInformation(c *gin.Context) {
	status := c.Query("status")
	err, roles, total := server.GetRoleInformation(status)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(res.RoleListRes{
			List:  roles,
			Total: total,
		}, c)
	}
}

// AddRoleInformation
// @Tags Role Management (角色管理)
// @Summary	人员管理--新增角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.AddRoleInformation true "新增角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/role-management/add [post]
func AddRoleInformation(c *gin.Context) {
	var R request.AddRoleInformation
	_ = c.ShouldBindJSON(&R)
	ProVerify := utils.Rules{
		"RoleName":      {utils.NotEmpty()},
		"RoleStatus":    {utils.NotEmpty()},
		"RoleCreatedBy": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(R, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err := server.AddRoleInformation(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("新增角色失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// EditRoleInformation
// @Tags Role Management (角色管理)
// @Summary	人员管理--修改角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.EditRoleInformation true "修改角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/role-management/edit [post]
func EditRoleInformation(c *gin.Context) {
	var R request.EditRoleInformation
	_ = c.ShouldBindJSON(&R)
	ProVerify := utils.Rules{
		"RoleBeforeName": {utils.NotEmpty()},
		"RoleStatus":     {utils.NotEmpty()},
		"RoleUpdatedBy":  {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(R, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err := server.EditRoleInformation(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("修改角色失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// DeleteRoleInformation
// @Tags Role Management (角色管理)
// @Summary	人员管理--删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteRoleInformation true "修改角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/role-management/delete [delete]
func DeleteRoleInformation(c *gin.Context) {
	var R request.DeleteRoleInformation
	_ = c.ShouldBindJSON(&R)
	ProVerify := utils.Rules{
		"Ids": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(R, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err := server.DeleteRoleInformation(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("删除角色失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetRoleMenuTree
// @Tags Role Management (角色管理)
// @Summary	人员管理--获取当前角色菜单树
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param roleId query string true "角色id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/role-management/sys-menu/tree [get]
func GetRoleMenuTree(c *gin.Context) {
	roleId := c.Query("roleId")
	if roleId == "" {
		response.FailWithMessage("roleId不能为空", c)
		return
	}
	err, r := server.GetRoleMenuRelationsByRoleId(roleId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	var rMap = make(map[string]string)
	for _, relation := range r {
		rMap[relation.MenuId] = relation.ID
	}
	err, treeList := ListToTreeHasSelect(rMap)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(treeList, c)
	}
}

// SetRelation
// @Tags Role Management (角色管理)
// @Summary	人员管理--设置当前角色改菜单的权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param roleAndMenuRelationForms body request.RoleAndMenuRelationFormList true "关系form"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/role-management/sys-menu/relation-set [post]
func SetRelation(c *gin.Context) {
	var reqParam request.RoleAndMenuRelationFormList
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//if !reqParam.IsSelect && reqParam.ID == "" {
	//	response.FailWithMessage("保存失败", c)
	//	return
	//}
	row, err := server.SetRoleMenuRelation(reqParam.List)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("保存失败,%s", err.Error()), c)
		return
	}
	response.OkDetailed(row, "保存成功", c)
}
