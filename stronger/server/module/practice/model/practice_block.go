package model

import (
	"server/module/sys_base/model"
)

//循环分块是同一个分块，用field的BlockIndex区分
//初审切出非初审分块 field.blockIndex = block.zero

type PracticeProjectBlock struct {
	model.Model
	BillID                string                 `json:"billID" gorm:"comment:'单ID'"`                        //单ID
	Temp                  string                 `json:"temp" gorm:"comment:'类型'"`                           //类型
	Name                  string                 `json:"name" gorm:"模板分块名字"`                                 //模板分块名字
	Code                  string                 `json:"code" gorm:"模板分块编码"`                                 //模板分块编码
	FEight                bool                   `json:"fEight" gorm:"是否f8提交"`                               //是否f8提交
	IsLoop                bool                   `json:"isLoop" gorm:"是否循环分块"`                               //是否循环分块
	IsInput               bool                   `json:"isInput" gorm:"是否需要录入"`                              //是否需要录入
	Status                int                    `json:"status" form:"status" gorm:"comment:'状态 1:使用 2:暂停'"` //录入状态
	Picture               string                 `json:"picture" gorm:"图片"`
	PracticeProjectFields []PracticeProjectField `json:"practiceProjectFields" form:"practiceProjectFields" gorm:"foreignKey:BlockID;references:id;comment:'单ID'"`
}
