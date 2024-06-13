package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	"server/module/pro_conf/model"
	proconf3 "server/module/pro_conf/model/request"
	respon "server/module/pro_conf/model/response"
	"server/module/pro_conf/service"
	"server/utils"
)

// CreateProjectConfigAgingHolidayTable
// @Tags AgingHoliday(节假日时效)
// @Summary 建表(前端请勿调用!!!!!!!!)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param configType query string true "配置类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"建表成功"}"
// @Router /pro-config/project-config-aging-holiday/createProjectConfigAgingHolidayTable [post]

func CreateProjectConfigAgingHolidayTable(c *gin.Context) {
	err := service.CreateAgingConfigHoliday()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("建表失败, %v", err), c)
	} else {
		response.OkWithData(nil, c)
	}
}

// InsertProjectConfigAgingHoliday
// @Tags AgingHoliday(节假日时效)
// @Summary 时效配置--添加节假日时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProjectConfigAgingHoliday true "添加节假日时效设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-config-aging-holiday/add [post]
func InsertProjectConfigAgingHoliday(c *gin.Context) {
	var R model.SysProjectConfigAgingHoliday
	_ = c.ShouldBindJSON(&R)
	AgingConfigVerify := utils.Rules{
		"Date": {utils.NotEmpty()},
	}
	AgingConfigVerifyErr := utils.Verify(R, AgingConfigVerify)
	if AgingConfigVerifyErr != nil {
		response.FailWithMessage(AgingConfigVerifyErr.Error(), c)
		return
	}
	err := service.InsertProjectConfigAgingHoliday(R)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("添加节假日时效失败, %v", err), c)
	} else {
		response.OkWithData("添加节假日时效成功", c)
	}
}

// GetProjectConfigAgingHoliday
// @Tags AgingHoliday(节假日时效)
// @Summary 时效配置--获取节假日时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param InquireStartDate query int true "起始时间"
// @Param InquireEndDate query int true "结束时间"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-config-aging-holiday/list [get]
func GetProjectConfigAgingHoliday(c *gin.Context) {
	var R proconf3.ProjectConfigAgingHoliday
	R.InquireStartDate = c.Query("InquireStartDate")
	R.InquireEndDate = c.Query("InquireEndDate")
	fmt.Println(R)
	AgingConfigVerify := utils.Rules{
		"InquireStartDate": {utils.NotEmpty()},
		"InquireEndDate":   {utils.NotEmpty()},
	}
	AgingConfigVerifyErr := utils.Verify(R, AgingConfigVerify)
	if AgingConfigVerifyErr != nil {
		response.FailWithMessage(AgingConfigVerifyErr.Error(), c)
		return
	}
	err, config, count := service.GetAgingConfigHoliday(R)
	fmt.Println(count)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取节假日时效设置失败, err : %v", err), c)
	} else {
		response.OkWithData(respon.ProjectConfigAgingResponse{
			List:  config,
			Total: count,
		}, c)
	}
}

// UpdateProjectConfigAgingHoliday
// @Tags AgingHoliday(节假日时效)
// @Summary 时效配置--更新节假日时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProjectConfigAgingHoliday true "更新节假日时效设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-config-aging-holiday/edit [post]
func UpdateProjectConfigAgingHoliday(c *gin.Context) {
	var R model.SysProjectConfigAgingHoliday
	err := c.ShouldBindJSON(&R)
	if err != nil {
		fmt.Println(err)
	}
	AgingConfigVerify := utils.Rules{
		"Id": {utils.NotEmpty()},
	}
	AgingConfigVerifyErr := utils.Verify(R, AgingConfigVerify)
	if AgingConfigVerifyErr != nil {
		response.FailWithMessage(AgingConfigVerifyErr.Error(), c)
		return
	}
	err = service.UpdateAgingConfigHoliday(R)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新节假日时效设置失败, err : %v", err), c)
	} else {
		response.OkWithData("更新节假日时效设置成功", c)
	}
}
