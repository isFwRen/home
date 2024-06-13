/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/25 14:24
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
)

// AddGroupNotice
// @Tags msg-manager(消息管理--固定通知)
// @Summary 消息管理--新增一组模板
// @Auth xingqiyi
// @Date 2022/7/25 14:24
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model2.GroupNoticeAddReq true "通知"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/group-notice/add [post]
func AddGroupNotice(c *gin.Context) {
	var groupNoticeAddReq model2.GroupNoticeAddReq
	err := c.ShouldBindJSON(&groupNoticeAddReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.AddGroupNotices(groupNoticeAddReq)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// GetGroupNoticeByGroupId
// @Tags msg-manager(消息管理--固定通知)
// @Summary 消息管理--获取该项目的模板
// @Auth xingqiyi
// @Date 2022/7/25 15:05
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param groupId        query   string   true    "群id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/group-notice/get-by-group-id [get]
func GetGroupNoticeByGroupId(c *gin.Context) {
	var groupNoticeReq model2.GroupNoticeReq
	err := c.ShouldBindQuery(&groupNoticeReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, ones, twos := service.GetGroupNoticeByGroupId(groupNoticeReq)
	var obj = make(map[int][]model.GroupNoticeOne, 0)
	for _, one := range ones {
		if obj[one.Block] == nil {
			obj[one.Block] = make([]model.GroupNoticeOne, 0)
		}
		obj[one.Block] = append(obj[one.Block], one)
	}

	var objTwo = make(map[int][]model.GroupNoticeTwo, 0)
	for _, two := range twos {
		if objTwo[two.Block] == nil {
			objTwo[two.Block] = make([]model.GroupNoticeTwo, 0)
		}
		objTwo[two.Block] = append(objTwo[two.Block], two)
	}
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(map[string]interface{}{
			"one": obj,
			"two": objTwo,
		}, c)
	}
}

// ReGroupNotice
// @Tags msg-manager(消息管理--固定通知)
// @Summary 消息管理--重置
// @Auth xingqiyi
// @Date 2022/7/25 17:00
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param groupId        query   string   true    "钉钉群id"
// @Param type        query   int   true    "模板类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/group-notice/re [post]
func ReGroupNotice(c *gin.Context) {
	var groupNoticeReq model2.GroupNoticeReq
	err := c.ShouldBindQuery(&groupNoticeReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	row := service.ReGroupNotice(groupNoticeReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkDetailed(row, "重置成功", c)
	}
}

// EditGroupNotice
// @Tags msg-manager(消息管理--固定通知)
// @Summary 消息管理--编辑
// @Auth xingqiyi
// @Date 2022/7/25 17:12
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model2.GroupNoticeAddReq true "通知"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/group-notice/edit [post]
func EditGroupNotice(c *gin.Context) {
	var groupNoticeAddReq model2.GroupNoticeAddReq
	err := c.ShouldBindJSON(&groupNoticeAddReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.EditGroupNotices(groupNoticeAddReq)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}
