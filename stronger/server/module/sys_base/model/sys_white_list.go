package model

import "github.com/lib/pq"

type SysWhiteList struct {
	Model
	ProCode          string         `json:"proCode"`                                     //项目编码
	ProName          string         `json:"proName"`                                     //项目名称
	TempName         string         `json:"tempName"`                                    //模板名称
	UserCode         string         `json:"userCode"`                                    //工号
	UserName         string         `json:"userName"`                                    //姓名
	BlockPermissions pq.StringArray `json:"blockPermissions" gorm:"type:varchar(100)[]"` //分块权限  存编码
}

type SysWhiteListPeopleSum struct {
	BlockSum map[string]string `json:"blockSum"`
}
