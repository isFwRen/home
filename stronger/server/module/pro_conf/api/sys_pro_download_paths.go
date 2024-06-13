/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/30 16:06
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
	"server/global"
	"server/global/response"
	"server/module/pro_conf/model/request"
	"server/module/pro_conf/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	"strings"
)

// SysProDownloadPaths
// @Author 阿先
// @Date 17:20 2020/11/30
// @Tags   pro-config/sys-pro-download-paths(下载路径配置)
// @Summary 配置管理--根据项目id获取下载配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proId query string true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-pro-download-paths/list [get]
func SysProDownloadPaths(c *gin.Context) {
	proId := c.Query("proId")
	err, projectConfigDownloadPaths := service.GetDownloadPathByProjectId(proId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		response.OkDetailed(responseSysBase.BasePageResult{
			List:      projectConfigDownloadPaths,
			Total:     0,
			PageIndex: 0,
			PageSize:  0,
		}, "查询成功", c)
	}
}

// SetDownloadPathAvailable
// @Author 阿先
// @Date 9:26 2020/12/1
// @Tags  pro-config/sys-pro-download-paths(下载路径配置)
// @Summary 配置管理--根据勾选项更改下载路径
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param sysProPathsList body request.SysProPathsList true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-pro-download-paths/set-available [post]
func SetDownloadPathAvailable(c *gin.Context) {
	var R request.SysProPathsList
	_ = c.ShouldBindJSON(&R)
	downloadConfigVerify := utils.Rules{
		"ID": {utils.NotEmpty()},
	}
	if len(R.SysProPathsList) == 0 {
		response.FailWithMessage("参数有误", c)
		return
	}
	for _, v := range R.SysProPathsList {
		downloadConfigVerifyErr := utils.Verify(v, downloadConfigVerify)
		if downloadConfigVerifyErr != nil {
			response.FailWithMessage(fmt.Sprintf("%v", downloadConfigVerifyErr), c)
			return
		}
	}
	err := service.SetDownloadPathAvailable(R)
	if err != nil {
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkDetailed("", "修改成功", c)
	}

}

// GetProcessList
// @Tags pro-config/sys-pro-download-paths(下载路径配置)
// @Summary 配置管理--获取程序进程开启情况
// @Auth xingqiyi
// @Date 2022/3/8 10:45 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-pro-download-paths/process/list [get]
func GetProcessList(c *gin.Context) {
	proCode := c.Query("proCode")
	if _, ok := global.ProCodeId[proCode]; !ok {
		response.FailWithMessage(fmt.Sprintf("查询失败,%v", global.NotPro), c)
		return
	}
	cmd := exec.Command("bash", "-c", "ps aux | grep SCREEN | grep .sh | grep "+proCode)
	output, err := cmd.CombinedOutput()
	obj := map[string]bool{
		"task":     strings.Index(string(output), "task") > -1,
		"download": strings.Index(string(output), "download") > -1,
		"load":     strings.Index(string(output), proCode+"_load") > -1,
		"typeLoad": strings.Index(string(output), "type_load") > -1,
		"export":   strings.Index(string(output), "export") > -1,
		"upload":   strings.Index(string(output), "upload") > -1,
	}
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败,%v", err), c)
	} else {
		response.OkDetailed(responseSysBase.BasePageResult{
			List: []map[string]bool{obj},
		}, "查询成功", c)
	}
}
