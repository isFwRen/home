/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/2/1 11:54
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	model2 "server/module/pro_conf/model"
	"server/module/pro_conf/service"
	responseSysBase "server/module/sys_base/model/response"
)

// GetIssuesByFieldId
// @Tags sys-issue(配置管理--问题件配置)
// @Summary 配置管理--获取字段问题件配置
// @Auth xingqiyi
// @Date 2023/2/1 11:55
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param fId  query   model2.SysIssue      true   "字段id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-issue/list [get]
func GetIssuesByFieldId(c *gin.Context) {
	var sysIssueEditReq model2.SysIssue
	err := c.ShouldBindQuery(&sysIssueEditReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service.GetIssueListByFieldId(sysIssueEditReq.FId)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     0,
			PageIndex: 1,
			PageSize:  0,
		}, c)
	}
}

// EditIssuesByFieldId
// @Tags sys-issue(配置管理--问题件配置)
// @Summary 配置管理--编辑字段问题件配置
// @Auth xingqiyi
// @Date 2023/2/1 11:55
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model2.SysIssueEditReq true "问题件实体list和字段id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /pro-config/sys-issue/edit [post]
func EditIssuesByFieldId(c *gin.Context) {
	var reqParam model2.SysIssueEditReq
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err = service.EditIssuesByFieldId(reqParam.FId, reqParam.List)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("编辑失败"), c)
	} else {
		response.OkWithData("编辑成功", c)
	}
}
