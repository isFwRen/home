package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/sys_base/model/request"
	responseSysBase "server/module/sys_base/model/response"
	"server/module/sys_base/service"
	"server/utils"
)

// GetUserInformation
// @Tags User Management (用户管理)
// @Summary	人员管理--查询用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param code query string false "工号"
// @Param name query string false "姓名"
// @Param role query string false "角色"
// @Param status query string false "状态"
// @Param startTime query string false "上岗时间start"
// @Param endTime query string false "上岗时间end"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/user-management/list [get]
func GetUserInformation(c *gin.Context) {
	var Search request.UserManagement
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := service.UserPage(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			Total:     total,
			PageIndex: Search.PageInfo.PageIndex,
			PageSize:  Search.PageInfo.PageSize,
		}, c)
	}
}

// AddUserInformation
// @Tags User Management (用户管理)
// @Summary	人员管理--新增用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UserAdd true "新增用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/user-management/add [post]
func AddUserInformation(c *gin.Context) {
	var R request.UserAdd
	_ = c.ShouldBindJSON(&R)
	AddUserVerify := utils.Rules{
		"Code":     {utils.NotEmpty()},
		"NickName": {utils.NotEmpty()},
		"Role":     {utils.NotEmpty()},
		"Status":   {utils.NotEmpty()},
		"Phone":    {utils.NotEmpty()},
	}
	AddUserVerifyVerifyErr := utils.Verify(R, AddUserVerify)
	if AddUserVerifyVerifyErr != nil {
		response.FailWithMessage(AddUserVerifyVerifyErr.Error(), c)
		return
	}
	err := service.AddUser(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("新增用户失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// UpdateProPermissionInUserInformation
// @Tags User Management (用户管理)
// @Summary	人员管理--更新用户项目权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ProPermissionInUserInformation true "更新用户项目权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/user-management/update-user-pro-permission [post]
func UpdateProPermissionInUserInformation(c *gin.Context) {
	var R request.ProPermissionInUserInformation
	_ = c.ShouldBindJSON(&R)
	AddUserVerify := utils.Rules{
		"Code": {utils.NotEmpty()},
	}
	AddUserVerifyVerifyErr := utils.Verify(R, AddUserVerify)
	if AddUserVerifyVerifyErr != nil {
		response.FailWithMessage(AddUserVerifyVerifyErr.Error(), c)
		return
	}
	err := service.UpdateProPermissionInUserInformation(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("更新用户项目权限失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetUserPermissionInformation
// @Tags User Management (用户管理)
// @Summary	人员管理--获取用户项目权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param userId query string true "用户Id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/user-management/user-pro-permission/list [get]
func GetUserPermissionInformation(c *gin.Context) {
	var R request.GetPermissionInformation
	if err := c.ShouldBindWith(&R, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	AddUserVerify := utils.Rules{
		"UserId": {utils.NotEmpty()},
	}
	AddUserVerifyVerifyErr := utils.Verify(R, AddUserVerify)
	if AddUserVerifyVerifyErr != nil {
		response.FailWithMessage(AddUserVerifyVerifyErr.Error(), c)
		return
	}

	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}

	err, list, total := service.GetUserPermissionInformation(R, uid)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询用户项目权限失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}
