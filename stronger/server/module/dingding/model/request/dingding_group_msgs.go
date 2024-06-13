/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 11:47
 */

package request

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
	"time"
)

type SendDingdingGroupMsgsStruct struct {
	SendMsg         string        `json:"sendMsg" gorm:"comment:'消息内容'"`
	DingdingGroupId pq.Int64Array `json:"dingdingGroupId" gorm:"type:int4[]; comment:'钉钉群id数组'"`
}

type OddsDingdingGroupMsgsStruct struct {
	model.PageInfo
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
