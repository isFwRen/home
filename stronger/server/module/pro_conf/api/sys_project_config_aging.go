package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/request"
	"server/global/response"
	proconf2 "server/module/pro_conf/model"
	proconf3 "server/module/pro_conf/model/request"
	proconf4 "server/module/pro_conf/model/response"
	"server/module/pro_conf/service"
	"server/utils"
)

// InsertAgingConfig
// @Tags Aging(时效配置)
// @Summary	时效配置--新增时效配置
// @accept application/json
// @Produce application/json
// @Param data body proconf3.ProjectConfigAging true "项目时效配置实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /pro-config/project-config-aging/add [post]
func InsertAgingConfig(c *gin.Context) {
	var R proconf3.ProjectConfigAging
	_ = c.ShouldBindJSON(&R)
	AgingConfigVerify := utils.Rules{
		"AgingStartTime": {utils.NotEmpty()},
		"AgingEndTime":   {utils.NotEmpty()},
		//"AgingOutStartTime": {utils.NotEmpty()},
		//"AgingOutEndTime":   {utils.NotEmpty()},
	}
	AgingConfigVerifyErr := utils.Verify(R, AgingConfigVerify)
	if AgingConfigVerifyErr != nil {
		response.FailWithMessage(AgingConfigVerifyErr.Error(), c)
		return
	}
	agingConfig := &proconf2.SysProjectConfigAgingReq{
		ProId:             R.ProId,
		AgingStartTime:    R.AgingStartTime,
		AgingEndTime:      R.AgingEndTime,
		AgingOutStartTime: R.AgingOutStartTime,
		AgingOutEndTime:   R.AgingOutEndTime,
		RequirementsTime:  R.RequirementsTime,
		ConfigType:        R.ConfigType,
		NodeName:          R.NodeName,
		NodeContent:       R.NodeContent,
		FieldName:         R.FieldName,
		FieldContent:      R.FieldContent,
	}
	err, configReturn := service.InsertAgingConfig(*agingConfig)
	if err != nil {
		response.FailWithDetailed(response.ERROR, proconf4.ProjectConfigAgingResponse{
			List:  configReturn,
			Total: 0,
		}, fmt.Sprintf("%v", err), c)
	} else {
		response.OkDetailed(proconf4.ProjectConfigAgingResponse{
			List:  configReturn,
			Total: 0,
		}, "添加成功", c)
	}
}

// DelAgingConfig
// @Tags  Aging(时效配置)
// @Summary	时效配置--删除时效配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.RmById true "id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-config-aging/delete [delete]
func DelAgingConfig(c *gin.Context) {
	//根据id删除该配置
	var reqId request.RmById
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := service.DelAgingConfig(reqId.Ids)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateAgingConfig
// @Tags Aging(时效配置)
// @Summary 时效配置--修改时效配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body proconf3.ProjectConfigAgingUpdate true "时效实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-config-aging/edit [post]
func UpdateAgingConfig(c *gin.Context) {
	var R proconf3.ProjectConfigAgingUpdate
	_ = c.ShouldBindJSON(&R)
	if R.ID == "" {
		err := errors.New("id不能为空")
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)

	}
	agingConfig := &proconf2.SysProjectConfigAgingReq{
		ProId:             R.ProId,
		AgingStartTime:    R.AgingStartTime,
		AgingEndTime:      R.AgingEndTime,
		AgingOutStartTime: R.AgingOutStartTime,
		AgingOutEndTime:   R.AgingOutEndTime,
		RequirementsTime:  R.RequirementsTime,
		ConfigType:        R.ConfigType,
		NodeName:          R.NodeName,
		NodeContent:       R.NodeContent,
		FieldName:         R.FieldName,
		FieldContent:      R.FieldContent,
	}
	if err := service.UpdateProjectConfigAging(*agingConfig, R.ID); err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// GetAgingConfigByConfigType
// @Tags Aging(时效配置)
// @Summary 时效配置--查询时效配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proId query string true "项目id"
// @Param configType query string true "项目时效类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-config-aging/list [get]
func GetAgingConfigByConfigType(c *gin.Context) {
	//根据配置类型查询时效配置,无分页
	var R proconf3.ReqProjectConfigAgingWithConfigType
	R.ProId = c.Query("proId")
	R.ConfigType = c.Query("configType")
	//if err := c.ShouldBindWith(&R, binding.Query); err != nil {
	//	response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	//	return
	//}
	err, configs := service.SelectAgingConfigByConfigType(R)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败, %v", err), c)
	} else {
		response.OkWithData(proconf4.ProjectConfigAgingResponse{
			List:  configs,
			Total: int64(len(configs)),
		}, c)
	}
}
