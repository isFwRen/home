/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 11:47
 */

package request

import (
	"server/module/sys_base/model"
	"time"
)

type AddGroupStruct struct {
	Index              string `json:"index"`
	GroupName          string `json:"groupName"`
	ProjectCode        string `json:"projectCode"`
	ProjectEnvironment string `json:"projectEnvironment"`
	AccessToken        string `json:"accessToken"`
	Secret             string `json:"secret"`
}

type UpdateDingdingGroupStruct struct {
	Id                 string `json:"id"`
	UpdatedAt          time.Time
	Index              string `json:"index"`
	GroupName          string `json:"groupName"`
	ProjectCode        string `json:"projectCode"`
	ProjectEnvironment string `json:"projectEnvironment"`
	AccessToken        string `json:"accessToken"`
	Secret             string `json:"secret"`
}

type OddsDingdingGroupStruct struct {
	model.PageInfo
	GroupName          string `json:"groupName"`
	ProjectCode        string `json:"projectCode"`
	ProjectEnvironment string `json:"projectEnvironment"`
}
