package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	dem "server/module/data_entry_management/model/response"
	"server/module/data_entry_management/service"
	sybase "server/module/sys_base/service"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// GetDataEntryChannelInformation
// @Tags Data Entry (录入管理)
// @Summary 录入管理--录入通道
// @accept application/json
// @Produce application/json
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /data-entry/channel/list [get]
func GetDataEntryChannelInformation(c *gin.Context) {
	customClaims, err1 := api.GetUserByToken(c)
	if err1 != nil {
		response.FailWithMessage(fmt.Sprintf("获取登录者失败，%v", err1), c)
	}
	//global.GLog.Info("customClaims.RoleId:::" + customClaims.RoleId)
	//获取项目权限
	err, p := sybase.GetAllPermissionByUId(customClaims.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败%s", err.Error()), c)
		return
	}
	err, ec := service.GetDataEntryChannelInformation(p)
	if err != nil {
		//response.FailWithMessage(fmt.Sprintf("获取数据失败, %v", err), c)
		response.OkWithData(dem.DataEntryChannelRes{
			List:   ec,
			Total:  int64(len(ec)),
			ErrMsg: err.Error(),
		}, c)
	} else {
		response.OkWithData(dem.DataEntryChannelRes{
			List:  ec,
			Total: int64(len(ec)),
		}, c)
	}
}
