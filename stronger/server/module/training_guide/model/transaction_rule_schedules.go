package model

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
)

// 放在项目下 项目规则进度
type RuleSchedule struct {
	model.Model
	UserId            string         `json:"userId"`
	ProjectCode       string         `json:"projectCode"`
	RequiredLearnList pq.StringArray `json:"requiredLearnList" gorm:"type:text[]"` //必学文件数组
	NeedLearnList     pq.StringArray `json:"needLearnList" gorm:"type:text[]"`     //待学文件数组 数组长度为0时代表全部学习完成
	TrainType         int            `json:"trainType"`                            //类型 1.业务规则 2.报销模板 3.教学视频
}

func (r *RuleSchedule) NoNeedLearn() bool {
	return len(r.NeedLearnList) == 0
}

var TrainTypeMap = map[int]string{
	1: "业务规则",
	2: "报销单模板",
	3: "教学视频",
}
