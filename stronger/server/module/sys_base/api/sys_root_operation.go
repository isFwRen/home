package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	"server/module/sys_base/model"
	request2 "server/module/sys_base/model/request"
	"server/module/sys_base/service"
	"server/utils"
)

// DeleteUser
// @Tags SysRootOperation
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.GetById true "删除用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /sys-root-operation/root-delete-user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var reqId model.GetById
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := service.DeleteUser(reqId.Id)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// SetUserAuthority
// @Tags SysRootOperation
// @Summary 设置用户权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.SetUserAuth true "设置用户权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /sys-root-operation/root-set-user-authority/setUserAuthority [post]
func SetUserAuthority(c *gin.Context) {
	var sua request2.SetUserAuth
	_ = c.ShouldBindJSON(&sua)
	UserVerify := utils.Rules{
		"UUID":        {utils.NotEmpty()},
		"AuthorityId": {utils.NotEmpty()},
	}
	UserVerifyErr := utils.Verify(sua, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	err := service.SetUserAuthority(sua.UUID, sua.AuthorityId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("修改失败，%v", err), c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// QueryUser
// @Tags SysRootOperation
// @Summary 工号查询
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.QueryUserId true "工号查询"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sys-root-operation/root-query-user/queryUser [post]
func QueryUser(c *gin.Context) {
	var id request2.QueryUserId
	_ = c.ShouldBindJSON(&id)
	UserVerify := utils.Rules{
		"IDCard": {utils.NotEmpty()},
	}
	UserVerifyErr := utils.Verify(id, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	fmt.Println(id)
	err, user := service.QueryUser(request2.RegisterStruct{IDCard: id.IDCard})
	msg := request2.RegisterStruct{Username: user.Username, Reason: user.Reason, State: user.State, Nickname: user.Nickname}
	if err != nil {
		response.FailWithMessage("抱歉,您所查询的身份证还没注册", c)
	} else {
		response.OkWithData(msg, c)
	}
}
