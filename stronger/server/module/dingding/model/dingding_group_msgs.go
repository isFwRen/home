/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 10:54
 */

package model

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
	"time"
)

type DingdingGroupMsgs struct {
	model.AutoAddIdModel
	SendMsg         string        `json:"sendMsg" gorm:"comment:'消息内容'"`
	DingdingGroupId pq.Int64Array `json:"dingdingGroupId" gorm:"type:int4[]; comment:'钉钉群id数组'"`
	SendDate        time.Time     `json:"sendDate" gorm:"comment:'发送时间日期'"`
	SendStatus      int8          `json:"sendStatus" gorm:"comment:'发送状态'"`
	SenderId        string        `json:"senderId" gorm:"comment:'发送人id'"`
	FailReason      string        `json:"failReason" gorm:"comment:'发送人id'"`
}
