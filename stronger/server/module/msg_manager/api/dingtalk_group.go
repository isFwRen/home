package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/msg_manager/const_data"
	"server/module/msg_manager/model"
	model2 "server/module/msg_manager/model/request"
	"server/module/msg_manager/service"
	requestSysBase "server/module/sys_base/model/request"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
)

// GetDingtalkGroupByPage
// @Tags msg-manager(消息管理--PT群管理)
// @Summary 消息管理--获取PT群配置分页
// @Auth xingqiyi
// @Date 2022/7/4 16:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param name        	 query   string   false    "群名称"
// @Param proCode        query   string   false    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param env       	 query   int      false    "环境"
// @Param orderBy        query   string   false   "排序JSON.stringify([["CreatedAt","desc"]])"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-group/page [get]
func GetDingtalkGroupByPage(c *gin.Context) {
	var dingtalkGroupSearch model2.DingtalkGroupReq
	err := c.ShouldBindQuery(&dingtalkGroupSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, groups := service.GetDingtalkGroupByPage(dingtalkGroupSearch)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      groups,
			Total:     total,
			PageIndex: dingtalkGroupSearch.PageIndex,
			PageSize:  dingtalkGroupSearch.PageSize,
		}, c)
	}
}

// AddDingtalkGroup
// @Tags msg-manager(消息管理--PT群管理)
// @Summary 消息管理--新增一个PT群配置
// @Auth xingqiyi
// @Date 2022/7/4 16:34
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DingtalkGroup true "PT群配置实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-group/add [post]
func AddDingtalkGroup(c *gin.Context) {
	var dingtalkGroup model.DingtalkGroup
	err := c.ShouldBindJSON(&dingtalkGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Name":        {utils.NotEmpty()},
		"ProCode":     {utils.NotEmpty()},
		"AccessToken": {utils.NotEmpty()},
		"Env":         {utils.NotEmpty()},
		"Secret":      {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(dingtalkGroup, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err = service.AddDingtalkGroup(dingtalkGroup)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// DelDingtalkGroupByIds
// @Tags msg-manager(消息管理--PT群管理)
// @Summary 消息管理--删除PT群配置
// @Auth xingqiyi
// @Date 2022/7/6 14:14
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "id数组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-group/delete [delete]
func DelDingtalkGroupByIds(c *gin.Context) {
	var idsIntReq requestSysBase.ReqIds
	err := c.ShouldBindJSON(&idsIntReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	rows := service.DelDingtalkGroupByIds(idsIntReq)
	if rows < 1 {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", rows), c)
	} else {
		response.OkDetailed(rows, "删除成功", c)
	}
}

// UpdateDingtalkGroupById
// @Tags msg-manager(消息管理--PT群管理)
// @Summary 消息管理--根据id更新PT群配置
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DingtalkGroup true "PT群配置实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-group/edit [post]
func UpdateDingtalkGroupById(c *gin.Context) {
	var dingtalkGroup model.DingtalkGroup
	err := c.ShouldBindJSON(&dingtalkGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"ID":          {utils.NotEmpty()},
		"Name":        {utils.NotEmpty()},
		"ProCode":     {utils.NotEmpty()},
		"AccessToken": {utils.NotEmpty()},
		"Secret":      {utils.NotEmpty()},
		"Env":         {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(dingtalkGroup, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err = service.UpdateDingtalkGroupById(dingtalkGroup)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

// GetDictConst
// @Tags msg-manager(消息管理--PT群管理)
// @Summary 消息管理--获取PT群字典常量
// @Auth xingqiyi
// @Date 2022/7/4 16:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-group/get-dict-const [get]
func GetDictConst(c *gin.Context) {
	obj := map[string]interface{}{
		"env":    const_data.EnvOptions,
		"type":   const_data.CustomerNoticeType,
		"status": const_data.CustomerNoticeStatus,
	}
	response.OkDetailed(obj, "查询成功", c)
}
