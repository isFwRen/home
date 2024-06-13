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

// DingtalkNoticeSend
// @Tags msg-manager(消息管理--实时通知)
// @Summary 消息管理--发送钉钉群消息
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model2.DingtalkNoticeSendReq true "钉钉通知实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-notice-msg/send [post]
func DingtalkNoticeSend(c *gin.Context) {
	var dingtalkNoticeSendReq model2.DingtalkNoticeSendReq
	err := c.ShouldBindJSON(&dingtalkNoticeSendReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Msg": {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(dingtalkNoticeSendReq, verify)
	if verifyErr != nil || len(dingtalkNoticeSendReq.GroupId) < 1 {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	//获取钉钉群信息
	err, list := service.GetDingtalkGroupByIds(dingtalkNoticeSendReq.GroupId)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	if list == nil || len(list) == 0 {
		response.FailWithMessage("没有这个群", c)
		return
	}

	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	var errMsg string
	for _, group := range list {
		dingtalkNotice := model.DingtalkNotice{
			GroupName: group.Name,
			GroupId:   group.ID,
			Msg:       dingtalkNoticeSendReq.Msg,
			SendCode:  customClaims.Code,
			SendName:  customClaims.NickName,
		}
		//存消息
		err, noticeId := service.AddDingtalkNotice(dingtalkNotice)
		dingtalkNotice.ID = noticeId
		if err != nil {
			errMsg += group.Name + "发送失败0；"
			global.GLog.Error(group.Name+"的消息新增失败", zap.Error(err))
			continue
		}
		//发消息
		robot := utils.NewRobot(group.AccessToken, group.Secret)
		err = robot.SendTextMessage(dingtalkNoticeSendReq.Msg, []string{}, true)
		//更新消息
		if err != nil {
			errMsg += group.Name + "发送失败1；"
			global.GLog.Error(group.Name+"的消息发送失败", zap.Error(err))
			dingtalkNotice.FailReason = err.Error()
			dingtalkNotice.SendStatus = 1
			err = service.UpdateDingtalkNoticeStatus(dingtalkNotice)
			if err != nil {
				errMsg += group.Name + "发送失败2；"
				global.GLog.Error(group.Name+"的消息更新状态(1:失败)失败", zap.Error(err))
			}
			continue
		}
		dingtalkNotice.SendStatus = 2
		err = service.UpdateDingtalkNoticeStatus(dingtalkNotice)
		if err != nil {
			errMsg += group.Name + "发送失败3；"
			global.GLog.Error(group.Name+"的消息更新状态(2:成功)失败", zap.Error(err))
		}
	}
	if errMsg != "" {
		response.FailWithMessage(errMsg, c)
	} else {
		response.OkWithMessage("发送成功", c)
	}
}

// GetDingtalkNoticeSendMsgPage
// @Tags msg-manager(消息管理--实时通知)
// @Summary 消息管理--群通知信息记录
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/dingtalk-notice-msg/page [get]
func GetDingtalkNoticeSendMsgPage(c *gin.Context) {
	var dingtalkNoticeMsgReqSearch model2.DingtalkNoticeMsgReq
	err := c.ShouldBindQuery(&dingtalkNoticeMsgReqSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service.GetDingtalkNoticeSendMsgPage(dingtalkNoticeMsgReqSearch)
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
