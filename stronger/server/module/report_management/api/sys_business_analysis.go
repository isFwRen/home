package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/report_management/model"
	R "server/module/report_management/model/response"
	"server/module/report_management/service"
	"server/utils"
)

// GetBusinessDownloadAnalysis
// @Tags Business Analysis (来量/回传分析)
// @Summary 特殊报表--查询来量分析
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param types query string true "类型,2种:来量-download/回传-upload"
// @Param isCheckAll query string true "全部true 或者 明细false"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/business-analysis/download/list [get]
func GetBusinessDownloadAnalysis(c *gin.Context) {
	//在initialize包下的cron.go、report.go、report.go 进行定时统计来量
	var Search model.BusinessAnalysisSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"startTime": {utils.NotEmpty()},
		"endTime":   {utils.NotEmpty()},
		"types":     {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, analysis, total := service.GetBusinessDownloadAnalysis(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(R.BusinessRes{
			List:  analysis,
			Total: total,
		}, c)
	}
}

// GetBusinessUploadAnalysis
// @Tags Business Analysis (来量/回传分析)
// @Summary 特殊报表--查询回传分析
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param types query string true "类型,2种:来量-download/回传-upload"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/business-analysis/upload/list [get]
func GetBusinessUploadAnalysis(c *gin.Context) {
	//在initialize包下的cron.go、report.go、report.go 进行定时统计来量
	var Search model.BusinessAnalysisSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"proCode":   {utils.NotEmpty()},
		"startTime": {utils.NotEmpty()},
		"endTime":   {utils.NotEmpty()},
		"types":     {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, analysis, total := service.GetBusinessUploadAnalysis(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(R.BusinessRes{
			List:  analysis,
			Total: total,
		}, c)
	}
}

// ExportBusinessDownloadAnalysis
// @Tags Error statistics (来量/回传分析)
// @Summary 特殊报表--导出来量分析
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param types query string true "类型,2种:来量-download/回传-upload"
// @Param isCheckAll query string true "全部true 或者 明细false"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/business-analysis/download/export [get]
func ExportBusinessDownloadAnalysis(c *gin.Context) {
	var Search model.BusinessAnalysisSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"startTime": {utils.NotEmpty()},
		"endTime":   {utils.NotEmpty()},
		"types":     {utils.NotEmpty()},
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
	//err, a := utils.GetRedisExport("BusinessDownloadAnalysis", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("特殊报表--导出来量分析导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("特殊报表--导出来量分析导出失败, err: %s", r))
	//			err = utils.DelRedisExport("BusinessDownloadAnalysis", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("特殊报表--导出来量分析删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("特殊报表--导出来量分析删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("BusinessDownloadAnalysis", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("BusinessDownloadAnalysis", response.Response{
	//			Code: 400,
	//			Data: nil,
	//			Msg:  fmt.Sprintf("BusinessDownloadAnalysis SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("BusinessDownloadAnalysis", uid)
	//		return
	//	}
	//	err, path := service.ExportBusinessDownloadAnalysis(Search)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("BusinessDownloadAnalysis", response.Response{
	//			Code: 400,
	//			Data: nil,
	//			Msg:  fmt.Sprintf("ExportBusinessDownloadAnalysis err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("BusinessDownloadAnalysis", uid)
	//		return
	//	}
	//	global.GSocketConnMap[uid].Emit("BusinessDownloadAnalysis", response.Response{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("BusinessDownloadAnalysis", uid)
	//	return
	//}()

	err, path, name := service.ExportBusinessDownloadAnalysis(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
}

// ExportBusinessUploadAnalysis
// @Tags Error statistics (错误统计)
// @Summary 特殊报表--导出回传分析
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param types query string true "类型,2种:来量-download/回传-upload"
// @Param isCheckAll query string true "全部true 或者 明细false"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/business-analysis/upload/export [get]
func ExportBusinessUploadAnalysis(c *gin.Context) {
	var Search model.BusinessAnalysisSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"proCode":   {utils.NotEmpty()},
		"startTime": {utils.NotEmpty()},
		"endTime":   {utils.NotEmpty()},
		"types":     {utils.NotEmpty()},
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
	//err, a := utils.GetRedisExport("ExportBusinessUploadAnalysis", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("特殊报表--导出回传分析导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("特殊报表--导出回传分析导出失败, err: %s", r))
	//			err = utils.DelRedisExport("ExportBusinessUploadAnalysis", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("特殊报表--导出回传分析删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("特殊报表--导出回传分析删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("ExportBusinessUploadAnalysis", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("ExportBusinessUploadAnalysis", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("ExportBusinessUploadAnalysis SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("ExportBusinessUploadAnalysis", uid)
	//		return
	//	}
	//	err, path := service.ExportBusinessUploadAnalysis(Search)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("ExportBusinessUploadAnalysis", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("ExportBusinessUploadAnalysis err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("ExportBusinessUploadAnalysis", uid)
	//		return
	//	}
	//	global.GSocketConnMap[uid].Emit("ExportBusinessUploadAnalysis", response.ExportResponse{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("ExportBusinessUploadAnalysis", uid)
	//	return
	//}()

	err, path, name := service.ExportBusinessUploadAnalysis(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
}
