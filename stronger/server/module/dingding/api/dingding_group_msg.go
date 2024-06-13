/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 11:30
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	dingding2 "server/module/dingding/model"
	dingding3 "server/module/dingding/model/request"
	dingding4 "server/module/dingding/model/response"
	"server/module/dingding/service"
	response2 "server/module/sys_base/model/response"
	"server/utils"
)

/**
 * @Author 阿先
 * @Date 16:34 2020/11/14
 * @Tags
 * @Summary	向钉钉群发送消息
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
 * @Router /dingdingGroups/addGroup [post]
 */
func SendDingdingGroupMsg(c *gin.Context) {
	var R dingding3.SendDingdingGroupMsgsStruct
	fmt.Println("get in here")
	_ = c.ShouldBindJSON(&R)
	DingdingVerify := utils.Rules{
		"SendMsg": {utils.NotEmpty()},
	}
	DingdingVerifyErr := utils.Verify(R, DingdingVerify)
	if DingdingVerifyErr != nil {
		response.FailWithMessage(DingdingVerifyErr.Error(), c)
		return
	}
	//调用第三方发送消息
	//TODO:应该通过消息队列发送
	//通过缓存获取发送者id
	dingdingGroupMsg := &dingding2.DingdingGroupMsgs{SendMsg: R.SendMsg, DingdingGroupId: R.DingdingGroupId, SenderId: ""}
	//dingdingGroupMsg := &model.DingdingGroupMsg{SendMsg:R.SendMsg, SenderId:""}
	err, groupMsgReturn := service.SendMsgToDingdingGroups(*dingdingGroupMsg)
	if err != nil {
		response.FailWithDetailed(response.ERROR, dingding4.DingdingGroupMsgsResponse{DingdingGroupMsgs: groupMsgReturn}, fmt.Sprintf("%v", err), c)
	} else {
		response.OkDetailed(dingding4.DingdingGroupMsgsResponse{DingdingGroupMsgs: groupMsgReturn}, "发送成功", c)
	}
}

/**
 * @Author 阿先
 * @Date 20:40 2020/11/14
 * @Tags
 * @Summary 分页获取钉钉群列表
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /dingdingGroups/GetGroupList [post]
 */
func GetDingdingGroupMsgList(c *gin.Context) {
	var pageInfo dingding3.OddsDingdingGroupMsgsStruct
	_ = c.ShouldBindJSON(&pageInfo)
	fmt.Println(pageInfo)
	PageVerifyErr := utils.Verify(pageInfo, utils.CustomizeMap["PageVerify"])
	fmt.Println(PageVerifyErr)
	if PageVerifyErr != nil {
		response.FailWithMessage(PageVerifyErr.Error(), c)
		return
	}
	err, list, total := service.SelectDingdingGroupMsgListByPage(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(response2.PageResult{
			List:      list,
			Total:     total,
			PageIndex: pageInfo.PageIndex,
			PageSize:  pageInfo.PageSize,
		}, c)
	}
}
