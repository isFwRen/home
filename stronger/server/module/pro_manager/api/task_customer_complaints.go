package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	"server/module/pro_manager/model"
	resp "server/module/pro_manager/model/response"
	server "server/module/pro_manager/service"
)

// GetCustomerComplaints
// @Tags Quality Management (客户投诉(录入系统))
// @Summary 客户投诉--查询
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param month query string true "月份"
// @Param proCode query string true "项目编码"
// @Param billName query string true "案件号"
// @Param wrongFieldName query string true "错误字段"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /task/customer_complaints/list [get]
func GetCustomerComplaints(c *gin.Context) {
	var search model.CustomerComplaints
	if err := c.ShouldBind(&search); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}

	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}

	err, list, total := server.GetCustomerComplaints(search, uid)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询客户投诉失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.CC{List: list, Total: total}, c)
	}
}
