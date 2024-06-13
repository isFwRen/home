package model

import (
	"server/module/sys_base/model"
	"time"
)

// "github.com/lib/pq"
type PracticeProjectBill struct {
	model.Model
	BillName string `json:"billName" form:"billName"` //单据号来源内部
	BillNum  string `json:"billNum" form:"billNum"`   //单据号来源客户
	ProCode  string `json:"proCode" form:"proCode"`   //项目编号
	Agency   string `json:"agency" form:"agency"`     //机构号
	Type     int    `json:"type" form:"type"`         //练习类型
	// IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
	Status                int                    `json:"status" form:"status"`             //录入状态
	DownloadPath          string                 `json:"downloadPath" form:"downloadPath"` //下载路径
	BatchNum              string                 `json:"batchNum" form:"batchNum"`         //批次号
	PracticeProjectBlocks []PracticeProjectBlock `json:"practiceProjectBlocks" form:"practiceProjectBlocks" gorm:"foreignKey:BillID;references:id;comment:'单ID'"`
	PracticeProjectFields []PracticeProjectField `json:"practiceProjectFields" form:"practiceProjectFields" gorm:"foreignKey:BillID;references:id;comment:'单ID'"`
}

type BillUpdate struct {
	BillNum string `json:"billNum" form:"billNum"` //单据号来源客户
	ProCode string `json:"proCode" form:"proCode"` //项目编号
	// IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
}

type BillDel struct {
	ID      string `json:"id" form:"id"`           //单据号来源客户
	ProCode string `json:"proCode" form:"proCode"` //项目编号
	// IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
}

type BillListSearch struct {
	BillNum string `json:"billNum" form:"billNum"` //单据号来源客户
	ProCode string `json:"proCode" form:"proCode"` //项目编号
	model.BasePageInfo
	// IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
}

type BillAdd struct {
	Bill    PracticeProjectBill `json:"bill" form:"bill"`       //单据号来源客户
	ProCode string              `json:"proCode" form:"proCode"` //项目编号
	// IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
}

type OpData struct {
	Bill           PracticeProjectBill      `json:"bill"`
	Block          PracticeProjectBlock     `json:"block"`
	Fields         [][]PracticeProjectField `json:"fields"`
	CodeValues     map[string]interface{}   `json:"codeValues"`
	CacheTime      int                      `json:"cacheTime"`
	ApplyAt        time.Time                `json:"applyAt"`
	FieldCharacter int                      `json:"fieldCharacter"`
	AccuracyRate   float64                  `json:"accuracyRate"`
	Videos         []string                 `json:"videos"`
}

type TaskSubmit struct {
	Bill   PracticeProjectBill      `json:"bill"`
	Block  PracticeProjectBlock     `json:"block"`
	Fields [][]PracticeProjectField `json:"fields"`
	Code   string                   `json:"code"`
	// Op     string                 `json:"op"`
}

type SumSearch struct {
	Code    string `json:"code" form:"code"`       //工号
	Name    string `json:"name" form:"name"`       //姓名
	ProCode string `json:"proCode" form:"proCode"` //项目编号
	model.BasePageInfo
	model.BaseTimeRange
	// IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
}

type TaskGet struct {
	Code string `json:"code" form:"code" gorm:"工号"`
	Name string `json:"name" form:"name" gorm:"姓名"`
	Op   string `json:"op" form:"op" gorm:"工序"`
	Num  int    `json:"num" form:"op" gorm:"返回数字"`
}
