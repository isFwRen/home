/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/25 14:10
 */

package model

import modelBase "server/module/sys_base/model"

type GroupNoticeOne struct {
	modelBase.Model
	DayOfWeek int    `json:"dayOfWeek"` //星期几0..6 0:星期天 新增必填
	ProCode   string `json:"proCode"`   //项目编码 新增必填
	StartTime string `json:"startTime"` //开始时间 新增必填
	EndTime   string `json:"endTime"`   //结束时间 新增必填
	Interval  int    `json:"interval"`  //间隔时间min 新增必填
	Block     int    `json:"block"`     //分组
	GroupId   string `json:"groupId"`   //群id 新增必填
}

type GroupNoticeTwo struct {
	modelBase.Model
	DayOfWeek int    `json:"dayOfWeek"` //星期几0..6 0:星期天 新增必填
	ProCode   string `json:"proCode"`   //项目编码 新增必填
	SendTime  string `json:"sendTime"`  //开始时间 新增必填
	Block     int    `json:"block"`     //分组
	GroupId   string `json:"groupId"`   //群id 新增必填
}
