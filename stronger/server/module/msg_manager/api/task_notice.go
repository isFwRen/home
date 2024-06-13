/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 14:52
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/msg_manager/model"
	model2 "server/module/msg_manager/model/request"
	"server/module/msg_manager/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// TaskNoticeSend
// @Tags msg-manager(消息管理--录入通知)
// @Summary 消息管理--发送录入通知消息
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model2.TaskNoticeSendReq true "录入通知实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/task-notice-msg/send [post]
func TaskNoticeSend(c *gin.Context) {
	var taskNoticeSendReq model2.TaskNoticeSendReq
	err := c.ShouldBindJSON(&taskNoticeSendReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Msg": {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(taskNoticeSendReq, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	if len(taskNoticeSendReq.ProCode) < 1 {
		response.FailWithMessage("项目编码错误", c)
		return
	}
	isSend := global.GSocketIo.BroadcastToNamespace("/global-notice", "notice", taskNoticeSendReq)
	if !isSend {
		response.FailWithMessage("发送失败", c)
		return
	}
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	var taskNoticeArr []model.TaskNotice
	for _, pro := range taskNoticeSendReq.ProCode {
		taskNotice := model.TaskNotice{
			ProCode:    pro,
			Msg:        taskNoticeSendReq.Msg,
			SendCode:   customClaims.Code,
			SendName:   customClaims.NickName,
			SendStatus: 2,
		}
		taskNoticeArr = append(taskNoticeArr, taskNotice)
	}
	err = service.AddTaskNotices(taskNoticeArr)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("发送成功", c)
	}
}

// GetTaskNoticeSendMsgPage
// @Tags msg-manager(消息管理--录入通知)
// @Summary 消息管理--录入通知信息记录
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/task-notice-msg/page [get]
func GetTaskNoticeSendMsgPage(c *gin.Context) {
	var dingtalkNoticeMsgReqSearch model2.DingtalkNoticeMsgReq
	err := c.ShouldBindQuery(&dingtalkNoticeMsgReqSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service.GetTaskNoticeSendMsgPage(dingtalkNoticeMsgReqSearch)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: dingtalkNoticeMsgReqSearch.PageIndex,
			PageSize:  dingtalkNoticeMsgReqSearch.PageSize,
		}, c)
	}
}
