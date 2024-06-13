package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/report_management/model/request"
	"server/module/report_management/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
)

// GetAbnormalBill
// @Tags Abnormal Bill (异常件数据)
// @Summary 特殊报表--查询异常件数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param types query string false "类型"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/abnormal-bill/list [get]
func GetAbnormalBill(c *gin.Context) {
	var Search request.AbnormalBillSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"proCode":   {utils.NotEmpty()},
		"startTime": {utils.NotEmpty()},
		"endTime":   {utils.NotEmpty()},
		"pageIndex": {utils.Gt("0")},
		"pageSize":  {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetAbnormalBill(Search)
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

// ExportAbnormalBill
// @Tags Abnormal Bill (异常件数据)
// @Summary 特殊报表--导出异常件数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param types query string false "类型"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/abnormal-bill/export [get]
func ExportAbnormalBill(c *gin.Context) {
	var Search request.AbnormalBillSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"proCode":   {utils.NotEmpty()},
		"startTime": {utils.NotEmpty()},
		"endTime":   {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	//uid := c.Request.Header.Get("x-user-id")
	//if uid == "" {
	//	response.FailWithMessage("x-user-id is empty", c)
	//	return
	//}
	//err, a := utils.GetRedisExport("exportAbnormalBill", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("特殊报表--导出异常件数据导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("特殊报表--导出异常件数据导出失败, err: %s", r))
	//			err = utils.DelRedisExport("exportAbnormalBill", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("特殊报表--导出异常件数据删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("特殊报表--导出异常件数据删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("exportAbnormalBill", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("exportAbnormalBill", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("ExportAbnormalBill SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("exportAbnormalBill", uid)
	//		return
	//	}
	//	err, path := service.ExportAbnormalBill(Search)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("exportAbnormalBill", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("ExportAbnormalBill SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("exportAbnormalBill", uid)
	//		return
	//	}
	//	global.GSocketConnMap[uid].Emit("exportAbnormalBill", response.ExportResponse{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("exportAbnormalBill", uid)
	//	return
	//}()

	err, path, name := service.ExportAbnormalBill(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
}
