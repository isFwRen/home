package api

import (
	"github.com/gin-gonic/gin"
	"server/global/response"
	reqModel "server/module/report_management/model/request"
	"server/module/report_management/service"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// GetTagList
// @Tags Report Setting (项目报表配置)
// @Summary 项目报表配置--获取选项列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /setting-report-tag/get-tag-list [get]
func GetTagList(c *gin.Context) {
	list := service.GetTagList()
	response.OkDetailed(list, "获取列表成功", c)
	return
}

// GetUserTags
// @Tags Report Setting (项目报表配置)
// @Summary 项目报表配置--获取报表配置列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body reqModel.ReqGetUserTagList true "项目编码"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /setting-report-tag/get-user-tags [post]
func GetUserTags(c *gin.Context) {
	var param reqModel.ReqGetUserTagList
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := api.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := service.GetUserTags(user.ID, param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
	return
}

// SetUserTags
// @Tags Report Setting (项目报表配置)
// @Summary 项目报表配置--设置报表配置列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body reqModel.ReqSettingReport true "项目编码+表头数组"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /setting-report-tag/set-user-tags [post]
func SetUserTags(c *gin.Context) {
	var param reqModel.ReqSettingReport
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := api.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.SetUserTags(user.ID, param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
	return
}
