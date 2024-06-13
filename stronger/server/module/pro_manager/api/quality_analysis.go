package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model/request"
	resp "server/module/pro_manager/model/response"
	server "server/module/pro_manager/service"
	"server/utils"
)

// GetQualityAnalysis
// @Tags Quality Management (质量分析)
// @Summary 质量分析--查询
// @accept application/json
// @Produce application/json
// @Param types query string true "类型:1-按项目 2-按字段 3-按人员"
// @Param startTime query string true "开始反馈日期YYYY-MM-DD"
// @Param endTime query string true "结束反馈日期YYYY-MM-DD"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/quality/analysis/list [get]
func GetQualityAnalysis(c *gin.Context) {
	//1-按项目 2-按字段 3-按人员
	var search request.QuaAnaRes
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetQualityAnalysisVerify := utils.Rules{
		"Types":     {utils.NotEmpty()},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	GetQualityAnalysisVerifyErr := utils.Verify(search, GetQualityAnalysisVerify)
	if GetQualityAnalysisVerifyErr != nil {
		response.FailWithMessage(GetQualityAnalysisVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetQualityAnalysis(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询质量分析失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.QualityRes{List: list, Total: total}, c)
	}
}

// ExportQualityAnalysis
// @Tags Quality Management (质量分析)
// @Summary 质量分析--导出
// @accept application/json
// @Produce application/json
// @Param startTime query string true "开始反馈日期YYYY-MM-DD"
// @Param endTime query string true "结束反馈日期YYYY-MM-DD"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/quality/analysis/export [get]
func ExportQualityAnalysis(c *gin.Context) {
	var search request.QuaAnaRes
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ExportQualityAnalysisVerify := utils.Rules{
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	ExportQualityAnalysisVerifyErr := utils.Verify(search, ExportQualityAnalysisVerify)
	if ExportQualityAnalysisVerifyErr != nil {
		response.FailWithMessage(ExportQualityAnalysisVerifyErr.Error(), c)
		return
	}
	//uid := c.Request.Header.Get("x-user-id")
	//if uid == "" {
	//	response.FailWithMessage("x-user-id is empty", c)
	//	return
	//}
	//err, a := utils.GetRedisExport("exportQualityAnalysis", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("质量分析--导出质量分析表导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("质量分析--导出质量分析表导出失败, err: %s", r))
	//			err = utils.DelRedisExport("exportQualityAnalysis", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("质量分析--导出质量分析表删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("质量分析--导出质量分析表删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("exportQualityAnalysis", uid)
	//	if err != nil {
	//		global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("exportQuality SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("exportQualityAnalysis", uid)
	//		return
	//	}
	//	err, p := server.ExportQualityAnalysis(search.StartTime, search.EndTime)
	//	fmt.Println("path", p)
	//	fmt.Println("uid", uid)
	//	if err != nil {
	//		global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("exportQuality err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("exportQualityAnalysis", uid)
	//		return
	//	}
	//	global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
	//		Code: 200,
	//		Data: p,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("exportQualityAnalysis", uid)
	//	fmt.Println("导出完成!")
	//	return
	//}()

	err, p := server.ExportQualityAnalysis(search.StartTime, search.EndTime)
	fmt.Println("path", p)
	fmt.Println("err", err)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		response.GetOkWithData(resp.QualityRes{List: p}, c)
	}
}
