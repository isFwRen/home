package model

import (
	"server/module/sys_base/model"
	"time"

	"github.com/lib/pq"
)

// "github.com/lib/pq"
type PracticeUser struct {
	model.Model
	Name     string    `json:"name" from:"name"`               //用户名
	Code     string    `json:"code" from:"code"`               //工号
	ProCode  string    `json:"proCode" from:"proCode"`         //项目编号
	ApplyAt  time.Time `json:"applyAt" gorm:"comment:'领取时间'"`  //1码领取时间
	SubmitAt time.Time `json:"submitAt" gorm:"comment:'提交时间'"` //1码提交时间
	// Status       int            `json:"status" from:"status"`                    //录入状态
	// DownloadPath string         `json:"downloadPath" from:"downloadPath"`        //下载路径
	// BatchNum     string         `json:"batchNum" from:"batchNum"`                //批次号
	CacheId pq.StringArray `json:"cacheId" gorm:"type:varchar(100)[];缓存Id"` //一码前置分块编码}
	Bcode   pq.StringArray `json:"bcode" gorm:"type:varchar(10)[];缓存分块编码"`  //一码前置分块编码}
}

type PracticeAsk struct {
	ID           string  `json:"id" form:"id"`
	ProCode      string  `json:"proCode" form:"proCode"`           //项目编号
	Character    int     `json:"character" form:"character"`       //机构号
	AccuracyRate float64 `json:"accuracyRate" form:"accuracyRate"` //机构号
}
