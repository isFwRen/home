/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/25 15:21
 */

package model

import "server/module/msg_manager/model"

type GroupNoticeReq struct {
	//ProCode string `json:"proCode" form:"proCode" binding:"required"`
	GroupId string `json:"groupId" form:"groupId" binding:"required"`
	Type    int    `json:"type" form:"type"` //1:模板一 2:模板二
}

type GroupNoticeAddReq struct {
	Type int `json:"type" form:"type" binding:"required"` //1:模板一 2:模板二
	Ones []model.GroupNoticeOne
	Twos []model.GroupNoticeTwo
}
