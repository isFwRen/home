package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model/request"
	r "server/module/pro_manager/model/response"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
)

// GetProjectAging
// @Tags Aging Management(时效管理)
// @Summary 时效管理--查询所有未回传的单子(单项目)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param caseNumber query string false "案件号"
// @Param agency query string false "机构号"
// @Param caseStatus query string false "案件状态"
// @Param stage query string false "录入状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-aging-management/list/Only [get]
func GetProjectAging(c *gin.Context) {
	//主要查询所有未回传的单子
	var Search request.ProjectAgingManagementSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetProjectAging(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			Total:     total,
			PageIndex: Search.PageInfo.PageIndex,
			PageSize:  Search.PageInfo.PageSize,
		}, c)
	}
}

// GetProjectAgingAll
// @Tags Aging Management(时效管理)
// @Summary 时效管理--查询所有未回传的单子(多项目)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param caseNumber query string false "案件号"
// @Param agency query string false "机构号"
// @Param caseStatus query string false "案件状态"
// @Param stage query string false "录入状态"
// @Param orderBy query   string   false   "排序JSON.stringify([["CreatedAt","desc"]])"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-aging-management/list [get]
func GetProjectAgingAll(c *gin.Context) {
	//主要查询所有未回传的单子
	var Search request.ProjectAgingManagementSearchAll
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetProjectAgingAll(Search)
	response.OkWithData(r.PageResult{
		List:      list,
		Total:     total,
		Msg:       err.Error(),
		PageIndex: Search.PageInfo.PageIndex,
		PageSize:  Search.PageInfo.PageSize,
	}, c)
}
