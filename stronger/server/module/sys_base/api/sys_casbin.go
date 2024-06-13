package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	sys_base2 "server/module/sys_base/model/request"
	sys_base3 "server/module/sys_base/model/response"
	"server/module/sys_base/service"
	"server/utils"
)

// @Tags casbin
// @Summary 更改角色api权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.sys_base.CasbinInReceive true "更改角色api权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /casbin/UpdateCasbin [post]

func UpdateCasbin(c *gin.Context) {
	var cmr sys_base2.CasbinInReceive
	_ = c.ShouldBindJSON(&cmr)
	AuthorityIdVerifyErr := utils.Verify(cmr, utils.CustomizeMap["AuthorityIdVerify"])
	if AuthorityIdVerifyErr != nil {
		response.FailWithMessage(AuthorityIdVerifyErr.Error(), c)
		return
	}
	err := service.UpdateCasbin(cmr.AuthorityId, cmr.CasbinInfos)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("添加规则失败，%v", err), c)
	} else {
		response.OkWithMessage("添加规则成功", c)
	}
}

// @Tags casbin
// @Summary 获取权限列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "获取权限列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /casbin/getPolicyPathByAuthorityId [post]

func GetPolicyPathByAuthorityId(c *gin.Context) {
	var cmr sys_base2.CasbinInReceive
	_ = c.ShouldBindJSON(&cmr)
	AuthorityIdVerifyErr := utils.Verify(cmr, utils.CustomizeMap["AuthorityIdVerify"])
	if AuthorityIdVerifyErr != nil {
		response.FailWithMessage(AuthorityIdVerifyErr.Error(), c)
		return
	}
	paths := service.GetPolicyPathByAuthorityId(cmr.AuthorityId)
	response.OkWithData(sys_base3.PolicyPathResponse{Paths: paths}, c)
}

// @Tags casbin
// @Summary casb RBAC RESTFUL测试路由
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "获取权限列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /casbin/CasbinTest [get]

func CasbinTest(c *gin.Context) {
	// 测试restful以及占位符代码  随意书写
	pathParam := c.Param("pathParam")
	query := c.Query("query")
	response.OkDetailed(gin.H{"pathParam": pathParam, "query": query}, "获取规则成功", c)
}
