package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/sys_base/model/request"
	responseSysBase "server/module/sys_base/model/response"
	"server/module/sys_base/service"
	"server/utils"
)

// GetWhileList
// @Tags While List Management (白名单管理)
// @Summary	人员管理--查询白名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param tempName query string true "模板名字"
// @Param code query string false "工号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/white-list/list [get]
func GetWhileList(c *gin.Context) {
	var Search request.GetWhiteList
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"TempName":  {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetWhileList(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			Total:     total,
			PageIndex: Search.PageInfo.PageIndex,
			PageSize:  Search.PageInfo.PageSize,
		}, c)
	}
}

// EditWhileList
// @Tags While List Management (白名单管理)
// @Summary	人员管理--修改白名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.EditWhiteListArr true "项目修改"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/white-list/edit [post]
func EditWhileList(c *gin.Context) {
	var R request.EditWhiteListArr
	_ = c.ShouldBindJSON(&R)
	ProVerify := utils.Rules{
		"ProCode":  {utils.NotEmpty()},
		"TempName": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(R, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err := service.EditWhiteList(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("修改白名单失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// ExportWhileList
// @Tags While List Management (白名单管理)
// @Summary	人员管理--导出白名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ExportWhileList true "导出"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/white-list/export [post]
func ExportWhileList(c *gin.Context) {
	var Search request.ExportWhileList
	_ = c.ShouldBindJSON(&Search)
	ProVerify := utils.Rules{
		"ProCode":  {utils.NotEmpty()},
		"TempName": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}
	err, a := utils.GetRedisExport("white-list", uid)
	if err == nil && a == "true" {
		response.FailWithMessage("正在导出!", c)
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(fmt.Sprintf("人员管理--导出白名单导出失败, err: %s", r))
				global.GLog.Error(fmt.Sprintf("人员管理--导出白名单导出失败, err: %s", r))
				err = utils.DelRedisExport("white-list", uid)
				if err != nil {
					fmt.Println(fmt.Sprintf("人员管理--导出白名单删除导出缓存失败, err: %s", err.Error()))
					global.GLog.Error(fmt.Sprintf("人员管理--导出白名单删除导出缓存失败, err: %s", err.Error()))
				}
			}
		}()
		//可以广播同一个登录人的客户端的写法
		//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
		err := utils.SetRedisExport("white-list", uid)
		if err != nil {
			global.GSocketConnMap[uid].Emit("whitelist", response.ExportResponse{
				Code: 400,
				Data: "",
				Msg:  fmt.Sprintf("white-list SetRedisExport err: %s", err.Error()),
			})
			err = utils.DelRedisExport("white-list", uid)
			return
		}
		err, path := service.ExportWhileList(Search)
		if err != nil {
			global.GSocketConnMap[uid].Emit("whitelist", response.ExportResponse{
				Code: 400,
				Data: "",
				Msg:  fmt.Sprintf("ExportWhileList err: %s", err.Error()),
			})
			err = utils.DelRedisExport("white-list", uid)
			return
		}
		global.GSocketConnMap[uid].Emit("whitelist", response.ExportResponse{
			Code: 200,
			Data: path,
			Msg:  "导出完成!",
		})
		err = utils.DelRedisExport("white-list", uid)
		return
	}()

	//_, _ = service.ExportWhileList(Search)
	//if err != nil {
	//	global.GLog.Error(err.Error())
	//	response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	//} else {
	response.Ok(c)
	//}
}

// CopyWhileList
// @Tags While List Management (白名单管理)
// @Summary	人员管理--复制白名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CopyWhileList true "复制白名单权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/white-list/copy [post]
func CopyWhileList(c *gin.Context) {
	var Search request.CopyWhileList
	_ = c.ShouldBindJSON(&Search)
	ProVerify := utils.Rules{
		"ProCode":  {utils.NotEmpty()},
		"TempName": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err := service.CopyWhileList(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("复制白名单失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetBlockPeopleSum
// @Tags While List Management (白名单管理)
// @Summary	人员管理--获取白名单每个分块勾选的人数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param tempName query string true "模板名字"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /sys-base/white-list/getBlockPeopleSum [get]
func GetBlockPeopleSum(c *gin.Context) {
	var Search request.GetWhiteList
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"ProCode":  {utils.NotEmpty()},
		"TempName": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list := service.GetBlockPeopleSum(Search.ProCode, Search.TempName)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			PageIndex: Search.PageInfo.PageIndex,
			PageSize:  Search.PageInfo.PageSize,
		}, c)
	}
}
