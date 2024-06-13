package api

import (
	"fmt"
	"server/global"
	"server/global/response"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	rep "server/module/report_management/model/response"
	"server/module/report_management/service"
	modelBase "server/module/sys_base/model"
	"server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateTable(c *gin.Context) {
	err := service.CreateTable()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("建表失败, %v", err), c)
	} else {
		response.OkWithData(nil, c)
	}
}

// GetOutputStatistics
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--查询人员产量统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param code query string false "工号"
// @Param isCheckAll query string true "全部1 or 明细2"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/output-Statistics/list [get]
func GetOutputStatistics(c *gin.Context) {
	var Search request.OutPutStatisticsSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetOutPutStatisticsVerify := utils.Rules{
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	query := c.Query("nickName")
	if query != "" {
		Search.Code = query
	}
	GetOutPutStatisticsVerifyErr := utils.Verify(Search, GetOutPutStatisticsVerify)
	if GetOutPutStatisticsVerifyErr != nil {
		response.FailWithMessage(GetOutPutStatisticsVerifyErr.Error(), c)
		return
	}

	if Search.UpdateTime == "" {
		//查询
		err, list, total, top := service.GetOutputStatistics(Search)

		if err != nil {
			global.GLog.Error(err.Error())
			response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
		} else {
			response.OkWithData(rep.PageResult{
				List:      list,
				Total:     total,
				Top:       top,
				PageIndex: Search.PageInfo.PageIndex,
				PageSize:  Search.PageInfo.PageSize,
			}, c)
		}
	} else {
		times, _ := strconv.Atoi(time.Now().Format("15"))
		if times < 18 || (times >= 0 && times < 8) {
			response.FailWithMessage(fmt.Sprintf("此功能18:00至次日8:00生效，目前无法更新，请知悉"), c)
			return
		}
		R := request.OutPutStatisticsExport{
			PageInfo:   Search.PageInfo,
			ProCode:    Search.ProCode,
			Code:       Search.Code,
			StartTime:  Search.StartTime,
			EndTime:    Search.EndTime,
			UpdateTime: Search.UpdateTime,
		}
		err, list, total, top := service.UpdateOutputStatistics(R)
		if err != nil {
			global.GLog.Error(err.Error())
			response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
		} else {
			response.OkWithData(rep.PageResult{
				List:      list,
				Total:     total,
				Top:       top,
				PageIndex: Search.PageInfo.PageIndex,
				PageSize:  Search.PageInfo.PageSize,
			}, c)
		}
	}

}

// GetOutputStatisticsTask
// @Tags Output Statistics (产量统计(录入系统))
// @Summary 产量统计(录入系统)--查询人员产量统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/output-Statistics-task/list [get]
func GetOutputStatisticsTask(c *gin.Context) {
	var Search request.OutPutStatisticsSearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetOutPutStatisticsVerify := utils.Rules{
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}

	GetOutPutStatisticsVerifyErr := utils.Verify(Search, GetOutPutStatisticsVerify)
	if GetOutPutStatisticsVerifyErr != nil {
		response.FailWithMessage(GetOutPutStatisticsVerifyErr.Error(), c)
		return
	}

	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}

	err, list, total, top := service.GetOutputStatisticsTask(Search, uid)

	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(rep.PageResult{
			List:      list,
			Total:     total,
			Top:       top,
			PageIndex: Search.PageInfo.PageIndex,
			PageSize:  Search.PageInfo.PageSize,
		}, c)
	}
}

// GetOCROutputStatistics
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--查询OCR产量统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/ocr-output-Statistics/list [get]
func GetOCROutputStatistics(c *gin.Context) {
	var Search request.GetOCROutPutStatisticsSearch
	// _ = c.ShouldBindJSON(&Search)
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	// GetOutPutStatisticsVerify := utils.Rules{}
	ProVerify := utils.Rules{
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(Search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetOcrStatistics(Search)

	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(rep.PageResult{
			List:      list,
			Total:     total,
			PageIndex: Search.PageInfo.PageIndex,
			PageSize:  Search.PageInfo.PageSize,
		}, c)
	}

	// GetOutPutStatisticsVerifyErr := utils.Verify(GetOutPutStatistics, GetOutPutStatisticsVerify)
	// if GetOutPutStatisticsVerifyErr != nil {
	// 	response.FailWithMessage(GetOutPutStatisticsVerifyErr.Error(), c)
	// 	return
	// }
}

// UpdateOutputStatistics
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--更新产量统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param updateTime query string true "更新时间"
// @Param code query string false "工号"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/output-Statistics/update [get]
func UpdateOutputStatistics(c *gin.Context) {
	times, _ := strconv.Atoi(time.Now().Format("15"))
	if times < 18 || (times >= 0 && times < 8) {
		response.FailWithMessage(fmt.Sprintf("此功能18:00至次日8:00生效，目前无法更新，请知悉"), c)
	}
	var R request.OutPutStatisticsExport
	if err := c.ShouldBindWith(&R, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
		"StartTime":  {utils.NotEmpty()},
		"EndTime":    {utils.NotEmpty()},
		"UpdateTime": {utils.NotEmpty()},
		"PageIndex":  {utils.Gt("0")},
		"PageSize":   {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(R, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, updateOutput, total, top := service.UpdateOutputStatistics(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData(rep.PageResult{
			List:      updateOutput,
			Total:     total,
			Top:       top,
			PageIndex: R.PageInfo.PageIndex,
			PageSize:  R.PageInfo.PageSize,
		}, c)
	}
}

// ExportOutputStatistics
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--导出人员(全部)产量统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/output-Statistics/export [get]
func ExportOutputStatistics(c *gin.Context) {
	var R request.OutPutStatisticsExport
	if err := c.ShouldBindWith(&R, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ProVerify := utils.Rules{
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
	//err, a := utils.GetRedisExport("outputStatistics", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}

	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("产量统计--导出人员(全部)产量统计导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("产量统计--导出人员(全部)产量统计导出失败, err: %s", r))
	//			err = utils.DelRedisExport("outputStatistics", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("产量统计--导出人员(全部)产量统计删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("产量统计--导出人员(全部)产量统计删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("outputStatistics", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("outputStatistics", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("ExportOutputStatistics SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("outputStatistics", uid)
	//		return
	//	}
	//	err, path := service.ExportOutputStatistics(R)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("outputStatistics", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("ExportOutputStatistics SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("outputStatistics", uid)
	//		return
	//	}
	//	global.GSocketConnMap[uid].Emit("outputStatistics", response.ExportResponse{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("outputStatistics", uid)
	//	return
	//}()

	err, path, name := service.ExportOutputStatistics(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
}

// ExportOutputStatisticsDetail
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--导出人员(明细)产量统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/output-Statistics-detail/export [get]
func ExportOutputStatisticsDetail(c *gin.Context) {
	var R request.OutPutStatisticsDetailExport
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
	//err, a := utils.GetRedisExport("outputStatisticsDetail", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("产量统计--导出人员(明细)产量统计导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("产量统计--导出人员(明细)产量统计导出失败, err: %s", r))
	//			err = utils.DelRedisExport("outputStatisticsDetail", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("产量统计--导出人员(明细)产量统计删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("产量统计--导出人员(明细)产量统计删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("outputStatisticsDetail", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("outputStatisticsDetail", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("outputStatisticsDetail SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("outputStatisticsDetail", uid)
	//		return
	//	}
	//	err, path := service.ExportOutputStatisticsDetail(R)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("outputStatisticsDetail", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("outputStatisticsDetail SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("outputStatisticsDetail", uid)
	//		return
	//	}
	//	global.GSocketConnMap[uid].Emit("outputStatisticsDetail", response.ExportResponse{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	//广播给在这个namespace的人
	//	//global.GSocketIo.BroadcastToNamespace("/global-export", "outputStatisticsDetail", response.ExportResponse{
	//	//	Code: 200,
	//	//	Data: path,
	//	//	Msg:  "导出完成!",
	//	//})
	//	err = utils.DelRedisExport("outputStatisticsDetail", uid)
	//	return
	//}()
	err, path, name := service.ExportOutputStatisticsDetail(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
	//go global.GSocketIo.BroadcastToNamespace("/test", "test1", "000")
}

// ExportOcrOutput
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--OCR产量统计导出报表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/ocr-output-Statistics/export [get]
func ExportOcrOutput(c *gin.Context) {
	var R request.GetOCROutPutStatisticsSearch
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
	err, path, name := service.ExportOcrOutput(R)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
	//go global.GSocketIo.BroadcastToNamespace("/test", "test1", "000")
}

// DeleteOutputStatisticsDetail
// @Tags Output Statistics (产量统计)
// @Summary 产量统计--删除某人某段时间的产量
// @accept application/json
// @Produce application/json
// @Param data body request.OutPutStatisticsSearch true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/output-Statistics/delete  [post]
func DeleteOutputStatisticsDetail(c *gin.Context) {
	var search request.OutPutStatisticsSearch
	_ = c.ShouldBindJSON(&search)
	GVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"Code":      {utils.NotEmpty()},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	GetCorrectedVerifyErr := utils.Verify(search, GVerify)
	if GetCorrectedVerifyErr != nil {
		response.FailWithMessage(GetCorrectedVerifyErr.Error(), c)
		return
	}
	err := service.DeleteOutputStatisticsDetail(search)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("清空失败, %v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetCorrected
// @Tags Corrected (项目折算)
// @Summary 产量统计--查询项目折算比例
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/corrected/list  [get]
func GetCorrected(c *gin.Context) {
	var Search modelBase.BasePageInfo
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetCorrectedVerify := utils.Rules{
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	GetCorrectedVerifyErr := utils.Verify(Search, GetCorrectedVerify)
	if GetCorrectedVerifyErr != nil {
		response.FailWithMessage(GetCorrectedVerifyErr.Error(), c)
		return
	}

	CorrectedArr, count, err := service.GetCorrected(Search)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询项目折算比例失败, %v", err), c)
	} else {
		response.GetOkWithData(rep.CorrectedResponse{
			List:  CorrectedArr,
			Total: count,
		}, c)
	}
}

// InsertCorrected
// @Tags Corrected (项目折算)
// @Summary 产量统计--增加项目折算比例
// @accept application/json
// @Produce application/json
// @Param data body request.InsertCorrected true "项目折算比例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /report-management/corrected/add  [post]
func InsertCorrected(c *gin.Context) {
	var insertCorrected request.InsertCorrected
	_ = c.ShouldBindJSON(&insertCorrected)
	InsertCorrectedVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	InsertCorrectedVerifyErr := utils.Verify(insertCorrected, InsertCorrectedVerify)
	if InsertCorrectedVerifyErr != nil {
		response.FailWithMessage(InsertCorrectedVerifyErr.Error(), c)
		return
	}
	StartTime, _ := time.ParseInLocation("2006-01-02", insertCorrected.StartTime, time.Local)
	InsertCorrect := model.SysCorrected{
		ProCode:              insertCorrected.ProCode,
		Op0AsTheBlock:        Decimal(insertCorrected.Op0AsTheBlock),
		Op0AsTheInvoice:      Decimal(insertCorrected.Op0AsTheInvoice),
		Op1NotExpenseAccount: Decimal(insertCorrected.Op2NotExpenseAccount),
		Op1ExpenseAccount:    Decimal(insertCorrected.Op1ExpenseAccount),
		Op2NotExpenseAccount: Decimal(insertCorrected.Op2NotExpenseAccount),
		Op2ExpenseAccount:    Decimal(insertCorrected.Op2ExpenseAccount),
		Question:             Decimal(insertCorrected.Question),
		StartTime:            StartTime,
	}
	err := service.InsertCorrected(InsertCorrect)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		response.OkWithMessage(fmt.Sprintf("保存成功"), c)
	}
}

// DeleteCorrected
// @Tags Corrected (项目折算)
// @Summary 产量统计--删除项目折算比例
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteCorrectedArr true "项目折算比例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /report-management/corrected/delete  [delete]
func DeleteCorrected(c *gin.Context) {
	var deleteCorrectedArr request.DeleteCorrectedArr
	err := c.ShouldBindJSON(&deleteCorrectedArr)
	if err != nil {
		fmt.Println("123", err)
	}
	DeleteCorrectedArrVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	DeleteCorrectedArrVerifyErr := utils.Verify(deleteCorrectedArr, DeleteCorrectedArrVerify)
	if DeleteCorrectedArrVerifyErr != nil {
		response.FailWithMessage(DeleteCorrectedArrVerifyErr.Error(), c)
		return
	}
	err = service.DeleteCorrected(deleteCorrectedArr)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除项目折算比例失败, %v", err), c)
	} else {
		response.OkWithMessage(fmt.Sprintf("删除项目折算比例成功"), c)
	}
}

// UpdateCorrected
// @Tags Corrected (项目折算)
// @Summary 产量统计--更新项目折算比例
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateCorrected true "项目折算比例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /report-management/corrected/edit  [post]
func UpdateCorrected(c *gin.Context) {
	var updateCorrected request.UpdateCorrected
	_ = c.ShouldBindJSON(&updateCorrected)
	UpdateCorrectedVerify := utils.Rules{
		"Id":      {utils.NotEmpty()},
		"ProCode": {utils.NotEmpty()},
	}
	UpdateCorrectedVerifyErr := utils.Verify(updateCorrected, UpdateCorrectedVerify)
	if UpdateCorrectedVerifyErr != nil {
		response.FailWithMessage(UpdateCorrectedVerifyErr.Error(), c)
		return
	}
	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}

	if updateCorrected.Op0AsTheBlock != 0 && updateCorrected.Op0AsTheInvoice != 0 {
		response.FailWithMessage(fmt.Sprintf("初审（按分块 ）、初审（按发票）只能输入一个"), c)
		return
	}

	var u modelBase.SysUser
	err := global.GDb.Model(&modelBase.SysUser{}).Where("id = ? ", uid).Find(&u).Error
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新项目折算比例失败"), c)
	}
	updateCorrected.Code = u.Code
	updateCorrected.Name = u.NickName
	err = service.UpdateCorrected(updateCorrected)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新项目折算比例失败"), c)
	} else {
		response.OkWithMessage(fmt.Sprintf("更新项目折算比例成功"), c)
	}
}

// GetEditCorrectedLog
// @Tags Corrected (项目折算)
// @Summary 产量统计--获取更改项目折算比例日志
// @accept application/json
// @Produce application/json
// @Param correctedID query string true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/corrected/list-log  [get]
func GetEditCorrectedLog(c *gin.Context) {
	correctedID := c.Query("correctedID")
	fmt.Println(correctedID)
	if correctedID != "" {
		Log, count, err := service.GetUpdateCorrectedLog(correctedID)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("获取更改项目折算比例日志失败, %v", err), c)
		} else {
			response.GetOkWithData(rep.CorrectedResponse{
				List:  Log,
				Total: count,
			}, c)
		}
	} else {
		response.FailWithMessage(fmt.Sprintf("获取更改项目折算比例日志失败, 缺少correctedID"), c)
	}
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
