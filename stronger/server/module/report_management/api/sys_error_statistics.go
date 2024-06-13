package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	responseReportManagement "server/module/report_management/model/response"
	"server/module/report_management/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	"time"
)

// GetIncorrectList
// @Tags Error statistics (错误查询)
// @Summary 错误查询--查询错误明细
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param code query string false "工号"
// @Param nickName query string false "姓名"
// @Param fieldName query string false "字段名称"
// @Param op query string false "工序"
// @Param complaint query string false "申诉"
// @Param confirm query string false "审核"
// @Param isAudit query bool false "是否已审核 false:待审核 true:全部"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/list [get]
func GetIncorrectList(c *gin.Context) {
	var Search model.WrongSearch
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
	err, list, total := service.GetIncorrectList(Search)
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

// ExportIncorrectList
// @Tags Error statistics (错误查询)
// @Summary 错误查询--导出错误明细
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param export body model.WrongExport true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/export [post]
func ExportIncorrectList(c *gin.Context) {
	var export model.WrongExport
	if err := c.ShouldBindJSON(&export); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	err, list := service.GetExportIncorrectList(export)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
		return
	}
	path := fmt.Sprintf("%v/%v/", global.GConfig.LocalUpload.FilePath, global.PathWrongExport)
	timeNow := time.Now()
	name := export.ProCode + "_错误明细_" + export.StartTime.Format("20060102") + "-" + export.EndTime.Format("20060102") + "_" + timeNow.Format("20060102150405") + ".xlsx"
	err = utils.ExportBigExcel(path, name, "sheet", list)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

// GetIncorrectTaskList
// @Tags Error statistics (录入系统-错误查询)
// @Summary 错误查询(录入系统)--查询错误明细
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param fieldName query string false "字段名称"
// @Param complaint query string false "申诉"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/task/list [get]
func GetIncorrectTaskList(c *gin.Context) {
	var Search model.WrongSearch
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
	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}
	err, list, total := service.GetIncorrectTaskList(Search, uid)
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

// Complain
// @Tags Error statistics (错误查询)
// @Summary 错误查询--错误明细,是否申诉
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Complain true "是否申诉"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/complain [post]
func Complain(c *gin.Context) {
	var complain request.Complain
	_ = c.ShouldBindJSON(&complain)
	confirmVerify := utils.Rules{
		"Id":      {utils.NotEmpty()},
		"ProCode": {utils.NotEmpty()},
	}
	confirmVerifyErr := utils.Verify(complain, confirmVerify)
	if confirmVerifyErr != nil {
		response.FailWithMessage(confirmVerifyErr.Error(), c)
		return
	}
	wrong, err := service.ComplainConfirm(complain)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("申诉失败，%v", err), c)
	} else {
		response.OkWithData(responseReportManagement.WrongRes{List: wrong, Message: "申诉成功"}, c)
	}
}

// ComplainTask @Tags Error statistics (错误查询)
// @Summary 错误查询(录入系统)--错误明细,是否申诉
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ComplainTask true "是否申诉"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/task/complain [post]
func ComplainTask(c *gin.Context) {
	var complain request.ComplainTask
	_ = c.ShouldBindJSON(&complain)
	confirmVerify := utils.Rules{
		"Id":      {utils.NotEmpty()},
		"ProCode": {utils.NotEmpty()},
	}
	confirmVerifyErr := utils.Verify(complain, confirmVerify)
	if confirmVerifyErr != nil {
		response.FailWithMessage(confirmVerifyErr.Error(), c)
		return
	}
	wrong, err := service.ComplainConfirmTask(complain)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("申诉失败，%v", err), c)
	} else {
		response.OkWithData(responseReportManagement.WrongRes{List: wrong, Message: "申诉成功"}, c)
	}
}

// WrongConfirm
// @Tags Error statistics (错误查询)
// @Summary 错误查询--错误明细,申诉是否通过
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.WrongConfirmArray true "申诉是否通过"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/wrong-confirm [post]
func WrongConfirm(c *gin.Context) {
	var confirm request.WrongConfirmArray
	_ = c.ShouldBindJSON(&confirm)
	confirmVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	confirmVerifyErr := utils.Verify(confirm, confirmVerify)
	if confirmVerifyErr != nil {
		response.FailWithMessage(confirmVerifyErr.Error(), c)
		return
	}
	//if confirm.WrongConfirm == "1" {
	//	//通过
	//	confirmVerify := utils.Rules{
	//		"ProCode": {utils.NotEmpty()},
	//	}
	//	confirmVerifyErr := utils.Verify(confirm, confirmVerify)
	//	if confirmVerifyErr != nil {
	//		response.FailWithMessage(confirmVerifyErr.Error(), c)
	//		return
	//	}
	//} else {
	//	confirmVerify := utils.Rules{
	//		"ProCode": {utils.NotEmpty()},
	//	}
	//	confirmVerifyErr := utils.Verify(confirm, confirmVerify)
	//	if confirmVerifyErr != nil {
	//		response.FailWithMessage(confirmVerifyErr.Error(), c)
	//		return
	//	}
	//}

	wrong, err := service.WrongConfirm(confirm)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("申诉失败，%v", err), c)
	} else {
		response.OkWithData(responseReportManagement.WrongRes{List: wrong, Message: "审核成功"}, c)
	}
}

// IncorrectAnalysis
// @Tags Error statistics (错误查询)
// @Summary 错误查询--查询错误分析
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param code query string false "工号"
// @Param nickName query string false "姓名"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/wrong-analysis/list [get]
func IncorrectAnalysis(c *gin.Context) {
	var Search model.WrongSearch
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
	err, list, total := service.IncorrectAnalysis(Search)
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

// OcrAnalysis
// @Tags Error statistics (错误查询)
// @Summary 错误查询--查询OCR错误明细
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param fieldName query string false "字段名称"
// @Param op query string false "工序"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/ocr-analysis/list [get]
func OcrAnalysis(c *gin.Context) {
	var Search model.WrongSearch
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
	err, list, total := service.OcrAnalysis(Search)
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

// ExportIncorrectAnalysis
// @Tags Error statistics (错误查询)
// @Summary 错误查询--导出错误分析
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/error-statistics/wrong-analysis/export [get]
func ExportIncorrectAnalysis(c *gin.Context) {
	var R request.IncorrectAnalysisExport
	if err := c.ShouldBindWith(&R, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(R, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	//uid := c.Request.Header.Get("x-user-id")
	//if uid == "" {
	//	response.FailWithMessage("x-user-id is empty", c)
	//	return
	//}
	//err, a := utils.GetRedisExport("incorrectAnalysis", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("错误查询--导出获取错误分析导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("错误查询--导出获取错误分析导出失败, err: %s", r))
	//			err = utils.DelRedisExport("incorrectAnalysis", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("错误查询--导出获取错误分析删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("错误查询--导出获取错误分析删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("incorrectAnalysis", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("incorrectAnalysis", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("errorStatistics SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("incorrectAnalysis", uid)
	//		return
	//	}
	//	err, path := service.ExportIncorrectAnalysis(R)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("incorrectAnalysis", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("errorStatistics err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("incorrectAnalysis", uid)
	//		return
	//	}
	//	global.GSocketConnMap[uid].Emit("incorrectAnalysis", response.ExportResponse{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("incorrectAnalysis", uid)
	//	return
	//}()

	err, path, name := service.ExportIncorrectAnalysis(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
}
