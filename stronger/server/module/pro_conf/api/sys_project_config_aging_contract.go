package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/request"
	"server/global/response"
	"server/module/pro_conf/model"
	proconf4 "server/module/pro_conf/model/response"
	"server/module/pro_conf/service"
	"server/utils"
)

// InsertAgingContractConfig
// @Tags AgingContract(合同时效)
// @Summary 时效配置--新增合同时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProjectConfigAgingContract true "新增合同时效配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /pro-config/project-config-aging-contract/add [post]
func InsertAgingContractConfig(c *gin.Context) {
	var contract model.SysProjectConfigAgingContract
	_ = c.ShouldBindJSON(&contract)
	ContractAgingConfigVerify := utils.Rules{
		"ContractStartTime": {utils.NotEmpty()},
		"ContractEndTime":   {utils.NotEmpty()},
		"ProId":             {utils.NotEmpty()},
		"RequirementsTime":  {utils.NotEmpty()},
	}
	AgingConfigVerifyErr := utils.Verify(contract, ContractAgingConfigVerify)
	if AgingConfigVerifyErr != nil {
		response.FailWithMessage(AgingConfigVerifyErr.Error(), c)
		return
	}
	err := service.InsertAgingContractConfig(contract)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("新增合同时效失败, %v", err), c)
		return

	} else {
		response.OkWithMessage("新增合同时效成功", c)
	}

}

// UpdateAgingContractConfig
// @Tags AgingContract(合同时效)
// @Summary 时效配置--修改合同时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProjectConfigAgingContract true "修改合同时效配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /pro-config/project-config-aging-contract/edit [post]
func UpdateAgingContractConfig(c *gin.Context) {
	var contract model.SysProjectConfigAgingContract
	_ = c.ShouldBindJSON(&contract)
	if contract.ID == "" {
		err := errors.New("id不能为空")
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
		return

	}
	if contract.ProId == "" {
		err := errors.New("项目编码不能为空")
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
		return

	}
	if contract.ClaimType == 0 {
		err := errors.New("理赔类型输入有误")
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
		return

	}
	err := service.UpdateAgingContractConfig(contract)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("修改合同时效失败, %v", err), c)
		return
	} else {
		response.OkWithMessage("修改合同时效成功", c)

	}

}

// DelAgingContractConfig
// @Tags AgingContract(合同时效)
// @Summary 时效配置--删除合同时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.RmById true "合同时效ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/project-config-aging-contract/delete [delete]
func DelAgingContractConfig(c *gin.Context) {
	var rmByIds request.RmById
	_ = c.ShouldBindJSON(&rmByIds)
	IdVerifyErr := utils.Verify(rmByIds, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := service.DelAgingContractConfig(rmByIds.Ids)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除合同时效失败, %v", err), c)
		return
	} else {
		response.OkWithMessage("删除合同时效成功", c)
	}
}

// GetAgingContractConfig
// @Tags AgingContract(合同时效)
// @Summary 时效配置--查询合同时效
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param code query string true "项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pro-config/project-config-aging-contract/list [get]
func GetAgingContractConfig(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage(fmt.Sprintf("查询合同时效失败，项目编码不能为空."), c)
		return
	}
	err, list, total := service.GetAgingContractConfig(code)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败, %v", err), c)
	} else {
		if len(list) > 0 {
			response.OkWithData(proconf4.ProjectConfigAgingResponse{
				List:  list,
				Total: total,
			}, c)
		} else {
			response.FailWithMessage(fmt.Sprintf("没有查询到数据~~~~~~~~"), c)
			return
		}
	}
}
