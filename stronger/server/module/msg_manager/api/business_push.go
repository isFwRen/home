/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/18 09:52
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/msg_manager/model"
	"server/module/msg_manager/service"
	responseSysBase "server/module/sys_base/model/response"
)

// GetBusinessPushByPage
// @Tags msg-manager(消息管理--业务通知)
// @Summary 消息管理--获取业务通知分页
// @Description 返回参考实体类 BusinessPushSend
// @Date 2023年08月16日15:26:09
// @Security ApiKeyAuth
// @Security ProCode
// @Accept json
// @Produce json
// @Param proCode        query   string   false    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Param msgType        query   int      false    "消息类型"
// @Success 200 {object} responseSysBase.BasePageResult
// @Router /msg-manager/business-push/page [get]
func GetBusinessPushByPage(c *gin.Context) {
	var businessPushSearchReq model.BusinessPushSearchReq
	err := c.ShouldBindQuery(&businessPushSearchReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	uid := c.GetHeader(global.XUserId)
	err, total, list := service.GetBusinessPushByPage(businessPushSearchReq, uid)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: businessPushSearchReq.PageIndex,
			PageSize:  businessPushSearchReq.PageSize,
		}, c)
	}
}

// Read
// @Tags msg-manager(消息管理--业务通知)
// @Summary 消息管理--标志为已读
// @Description
// @Date 2023/8/24
// @Security ApiKeyAuth
// @Security UserID
// @Accept json
// @Produce json
// @Param data body model.BusinessPushSendReadReq true "请求实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /msg-manager/business-push/read [post]
func Read(c *gin.Context) {
	var search model.BusinessPushSendReadReq
	err := c.ShouldBindJSON(&search)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	row := service.Read(search)
	if row == 0 {
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkWithData(row, c)
	}
}
