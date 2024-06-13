/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/9 1:48 下午
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	"server/module/sys_base/model/request"
	responseSysBase "server/module/sys_base/model/response"
	"server/module/sys_base/service"
)

// GetPageByType
// @Tags SysLogger
// @Summary 日志管理--分页获取系统日志
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param startTime      query   time.Time   true    "开始时间"
// @Param endTime        query   time.Time   true    "结束时间"
// @Param functionModule      	 query   string   false   "功能模块"
// @Param moduleOperation         	 query   string  false    "模块操作"
// @Param operationPeople		   	 query   string  false    "工号、姓名"
// @Param logType		 query   int   	  true    "日志类型"
// @Success 200 {object} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-logger/list [get]
func GetPageByType(c *gin.Context) {
	var param request.SysLogger
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, logs := service.GetPageByType(param)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败%s", err.Error()), c)
		return
	}
	response.OkDetailed(responseSysBase.BasePageResult{
		List:      logs,
		Total:     total,
		PageIndex: param.PageIndex,
		PageSize:  param.PageSize,
	}, "获取成功", c)
}
