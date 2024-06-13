package api

import (
	"fmt"
	"math"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetBillDelByPage
// @Tags 删除历史日志
// @Summary 删除历史日志--分页查询
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-del-list/page [get]
func GetBillDelByPage(c *gin.Context) {
	var billListSearch model.BillListSearch
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, bills := service.GetDelBills(billListSearch)
	//utils.DictConvert(bills)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		projectBillExcels := []model.ProjectDelBillExcel{}
		for ii, bill := range bills {
			billExcel := model.ProjectDelBillExcel{}
			billExcel.Idx = ii + 1
			billExcel.BillNum = bill.BillNum
			billExcel.Size = utils.ToString(bill.Size) + "KB"
			billExcel.DelAt = bill.DelAt
			billExcel.ScanAt = bill.ScanAt
			billExcel.Remarks = bill.Remarks
			billExcel.Stage = bill.Stage
			billExcel.Name = bill.Name
			billExcel.Describe = bill.Describe
			billExcel.Image = bill.Image
			projectBillExcels = append(projectBillExcels, billExcel)
		}
		response.OkWithData(responseSysBase.BasePageResult{
			List:      projectBillExcels,
			Total:     total,
			PageIndex: billListSearch.PageIndex,
			PageSize:  billListSearch.PageSize,
		}, c)
	}
}

// HistoryDelList
// @Tags 删除历史日志
// @Summary 删除历史日志--产生数据
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-del-list/history [get]
func HistoryDelList(c *gin.Context) {
	var billListSearch model.HistoryDelProCode
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total := service.HistoryDelList(billListSearch.ProCode)
	// //utils.DictConvert(bills)
	fmt.Println("--------total-------", total)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("操作失败"), c)
	} else {
		response.OkWithData("操作成功", c)
	}

}

// GetBillDelSum
// @Tags 删除历史日志
// @Summary 删除历史日志--统计
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-del-list/sum [get]
func GetBillDelSum(c *gin.Context) {
	var billListSearch model.BillListSearch
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	billListSearch.PageSize = -1
	err, _, bills := service.GetDelBills(billListSearch)
	//utils.DictConvert(bills)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		rojectCache := global.GProConf[billListSearch.ProCode]
		data := map[string]int{
			"total":   0,
			"dNum":    0,
			"uNum":    0,
			"total_2": 0,
			"dNum_2":  0,
			"uNum_2":  0,
			"day":     rojectCache.SaveDate,
		}
		for _, bill := range bills {
			if bill.Stage == "下载文件" {
				if bill.Remarks == "秒赔" {
					data["dNum_2"]++
					data["total_2"] += int(math.Round(bill.Size))
				} else {
					data["dNum"]++
					data["total"] += int(math.Round(bill.Size))
				}
			} else {
				if bill.Remarks == "秒赔" {
					data["uNum_2"]++
				} else {
					data["uNum"]++
				}
			}

		}

		response.OkWithData(data, c)
	}
}

// ExportBillDel
// @Tags 删除历史日志
// @Summary 删除历史日志--导出报表
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-del-list/export [get]
func ExportBillDel(c *gin.Context) {
	var billListSearch model.BillListSearch
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, path, name := service.ExportExcel(billListSearch)
	//utils.DictConvert(bills)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
		return
	}
	c.FileAttachment(path+name, name)
}

// SetDataInRedis
// @Tags 案件列表
// @Summary 案件列表--添加到redis
// @accept application/json
// @Produce application/json
// @param data body request.QualitiesReqDemo true "type，data必填"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/save-qualities-demo [post]
func SetDataInRedis(c *gin.Context) {

	var param request.QualitiesReqDemo
	err := c.BindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	global.GRedis.Set(param.ID+param.Type, param.Data, -1)

	response.OkWithMessage("保存成功", c)
	return

}

// GetDataInRedis
// @Tags 案件列表
// @Summary 案件列表--获取数据
// @accept application/json
// @Produce application/json
// @param data body request.QualitiesReqDemo true "查询参数，type必填,data传空"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/get-qualities-demo [post]
func GetDataInRedis(c *gin.Context) {
	var param request.QualitiesReqDemo
	err := c.BindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res := global.GRedis.Get(param.ID + param.Type)
	strRes := res.Val()

	response.OkWithData(strRes, c)
	return
}
