package api

import (
	"fmt"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetDeductionDetailsByPage
// @Tags pro-manager/bill-deduction_details(特殊报表)
// @Summary 特殊报表--扣费明细列表
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param timeStart      query   string   true   "开始时间今天开始时间格式'2021-11-05 15:04:05'"
// @Param timeEnd        query   string   true   "结束时间今天结束时间格式'2021-11-05 15:04:05'"
// @Param billCode       query   string   false   "案件号"
// @Param name    		 query   string   false   "名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-deduction_details/page [get]
func GetDeductionDetailsByPage(c *gin.Context) {
	var billListSearch model.DeductionDetailSearch
	err := c.ShouldBindQuery(&billListSearch)
	fmt.Println("--------------billListSearch--------------", billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, bills := service.GetDeductionDetails(billListSearch)
	//utils.DictConvert(bills)
	if err != nil {
		fmt.Println("--------------err--------------", err)
		global.GLog.Error("-------------------------", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      bills,
			Total:     total,
			PageIndex: billListSearch.PageIndex,
			PageSize:  billListSearch.PageSize,
		}, c)
	}
}

// ExportDeductionDetails
// @Tags pro-manager/bill-deduction_details(特殊报表)
// @Summary 特殊报表--扣费明细列表导出excel
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param timeStart      query   string   true   "开始时间今天开始时间格式'2021-11-05 15:04:05'"
// @Param timeEnd        query   string   true   "结束时间今天结束时间格式'2021-11-05 15:04:05'"
// @Param billCode       query   string   false   "案件号"
// @Param name    		 query   string   false   "名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-deduction_details/export [get]
func ExportDeductionDetails(c *gin.Context) {
	var billListSearch model.DeductionDetailSearch
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, path, name := service.ExportDeductionDetailExcel(billListSearch)
	//utils.DictConvert(bills)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
		return
	}
	c.FileAttachment(path+name, name)
}
