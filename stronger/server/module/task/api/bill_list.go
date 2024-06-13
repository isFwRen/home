/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 10:34 上午
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	"server/module/pro_manager/model"
	responseSysBase "server/module/sys_base/model/response"
	"server/module/task/service"
)

// GetTaskBillByPage
// @Tags pro-manager/bill-list(项目管理)
// @Summary 获取案件列表
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param timeStart      query   string   false   "开始时间今天开始时间格式'2021-11-05 15:04:05'"
// @Param timeEnd        query   string   false   "结束时间今天结束时间格式'2021-11-05 15:04:05'"
// @Param billCode       query   string   false   "案件号"
// @Param status         query   int      false   "案件状态"
// @Param saleChannel    query   string   false   "案件状态"
// @Param batchNum       query   string   false   "批次号"
// @Param agency         query   string   false   "机构号"
// @Param insuranceType  query   int      false   "医保类型"
// @Param claimType      query   int      false   "理赔类型"
// @Param stickLevel     query   int      false   "加急件"
// @Param minCountMoney  query   float32  false   "最小账单金额"
// @Param maxCountMoney  query   float32  false   "最大账单金额"
// @Param isQuestion     query   int      false   "是否是问题件0:不筛选，1：true，2：false"
// @Param invoiceNum     query   int      false   "发票数量"
// @Param qualityUser    query   string   false   "质检人"
// @Param stage          query   string   false   "录入状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/common/bill-list/page [get]

func GetTaskBillByPage(c *gin.Context) {
	var billListSearch model.BillListSearch
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, bills := service.GetTaskBillByPage(billListSearch)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      bills,
			Total:     total,
			PageIndex: billListSearch.PageIndex,
			PageSize:  billListSearch.PageSize,
		}, c)
	}
}
