package model

import (
	"github.com/lib/pq"
	modelbase "server/module/sys_base/model"
)

type TeachVideo struct {
	modelbase.Model
	ProId          string         `json:"proId" excel:"序号"`
	SysBlockConfId string         `json:"sysBlockConfId"`
	SysBlockCode   string         `json:"sysBlockCode" excel:"分块编码"`
	SysBlockName   string         `json:"sysBlockName" excel:"分块名字"`
	InputRule      string         `json:"inputRule" excel:"教学视频"`
	Video          pq.StringArray `json:"video" gorm:"type:varchar(255)[]"`
	IsRequired     int            `json:"isRequired"`
}

type TeachVideos struct {
	modelbase.Model `json:"model"`
	ProId           string              `json:"proId" excel:"序号"`
	SysBlockConfId  string              `json:"sysBlockConfId"`
	SysBlockCode    string              `json:"sysBlockCode" excel:"分块编码"`
	SysBlockName    string              `json:"sysBlockName" excel:"分块名字"`
	InputRule       string              `json:"inputRule" excel:"教学视频"`
	IsRequired      int                 `json:"isRequired"`
	Video           []map[string]string `json:"video"`
}
