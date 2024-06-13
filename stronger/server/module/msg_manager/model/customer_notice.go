/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/18 10:53
 */

package model

import (
	"server/module/sys_base/model"
	"time"
)

//CustomerNoticeSearchReq 分页请求参数
type CustomerNoticeSearchReq struct {
	model.BaseTimePageCode
	MsgType int `json:"msgType" form:"msgType"` //消息类型
	Status  int `json:"status" form:"status"`   //状态
}

//CustomerNotice 实体类
type CustomerNotice struct {
	model.Model
	ProCode      string    `json:"proCode"`      //项目编码 (回复必填)
	SendTime     time.Time `json:"sendTime"`     //发送时间
	FileName     string    `json:"fileName"`     //消息文件名
	MsgType      int       `json:"msgType"`      //消息类型
	Status       int       `json:"status"`       //状态
	ReplyTime    time.Time `json:"replyTime"`    //回复时间
	DealUserCode string    `json:"dealUserCode"` //处理人工号
	DealUserName string    `json:"dealUserName"` //处理人名字
	Content      string    `json:"content"`      //消息内容
	IsReply      bool      `json:"isReply"`      //是否回复 (回复必填)
	ExpectNum    int       `json:"expectNum"`    //预计可增加单量
	DownloadPath string    `json:"downloadPath"` //下载文件路径
	UploadPath   string    `json:"uploadPath"`   //回传文件路径
}
