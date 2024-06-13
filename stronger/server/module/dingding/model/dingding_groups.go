/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 10:54
 */

package model

import (
	"server/module/sys_base/model"
)

type DingdingGroups struct {
	model.AutoAddIdModel
	Index              string `json:"index" gorm:"comment:'序号'"`
	GroupName          string `json:"groupName" gorm:"comment:'群名称'"`
	ProjectCode        string `json:"projectCode" gorm:"comment:'所属项目编码'"`
	ProjectEnvironment string `json:"projectEnvironment" gorm:"comment:'环境'"`
	AccessToken        string `json:"accessToken" gorm:"comment:'群编号'"`
	Secret             string `json:"secret"`
}
