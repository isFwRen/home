package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/global/response"
	res "server/module/sys_base/model/response"
	"server/module/sys_base/service"
)

// GetUsersInformation
// @Tags SysUser
// @Summary 人员管理--查询员工
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Produce application/json
// @Param code query string false "工号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-base/user/find [get]
func GetUsersInformation(c *gin.Context) {
	code := c.Query("code")
	err, list, total := service.GetUsersInformation(code)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(res.RoleListRes{
			List:  list,
			Total: total,
		}, c)
	}
}
