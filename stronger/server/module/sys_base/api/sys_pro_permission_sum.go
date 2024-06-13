package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/global/response"
	res "server/module/sys_base/model/response"
	server "server/module/sys_base/service"
)

// GetProPermissionSum
// @Tags Pro Permission Sum (项目权限统计)
// @Summary	人员管理--查询项目权限统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/pro-permission-check/list [get]
func GetProPermissionSum(c *gin.Context) {
	err, total, ProPermissionSum := server.GetProPermissionSum()
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(res.ProPermissionSum{
			List:  ProPermissionSum,
			Total: total,
		}, c)
	}
}
