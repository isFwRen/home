package model

import (
	"server/module/sys_base/model"
	"time"
)

type ProjectDelBill struct {
	model.Model
	BillName string    `json:"billName" form:"billName"`                                    //
	BillNum  string    `json:"billNum" form:"billNum" gorm:"comment:'单据号'" excel:"单据号"`     //单据号
	Size     float64   `json:"size" form:"size" gorm:"comment:'文件大小'" excel:"文件大小"`         //文件大小
	DelAt    time.Time `json:"delAt" form:"delAt" gorm:"comment:'删除时间'" excel:"删除时间"`       //删除时间
	ScanAt   time.Time `json:"scanAt" form:"scanAt" gorm:"comment:'创建时间'" excel:"创建时间"`     //创建时间
	Remarks  string    `json:"remarks" form:"remarks" gorm:"comment:'备注'" excel:"备注"`       //备注
	Stage    string    `json:"stage" form:"stage" gorm:"comment:'说明'" excel:"说明"`           //说明
	Name     string    `json:"name" form:"name" gorm:"comment:'操作人'" excel:"操作人"`           //操作人
	Describe string    `json:"describe" form:"describe" gorm:"comment:'描述文件'" excel:"描述文件"` //描述文件
	Image    string    `json:"image" form:"image" gorm:"comment:'影像文件'" excel:"影像文件"`       //影像文件
	//单据状态
}

type ProjectDelBillExcel struct {
	model.Model
	Idx      int       `json:"idx" form:"idx" gorm:"comment:'序号'" excel:"序号"`
	BillNum  string    `json:"billNum" form:"billNum" gorm:"comment:'单据号'" excel:"单据号"`     //单据号
	Size     string    `json:"size" form:"size" gorm:"comment:'文件大小'" excel:"文件大小"`         //文件大小
	ScanAt   time.Time `json:"scanAt" form:"scanAt" gorm:"comment:'创建时间'" excel:"创建时间"`     //创建时间
	DelAt    time.Time `json:"delAt" form:"delAt" gorm:"comment:'删除时间'" excel:"删除时间"`       //删除时间
	Describe string    `json:"describe" form:"describe" gorm:"comment:'描述文件'" excel:"描述文件"` //描述文件
	Image    string    `json:"image" form:"image" gorm:"comment:'影像文件'" excel:"影像文件"`       //影像文件
	Name     string    `json:"name" form:"name" gorm:"comment:'操作人'" excel:"操作人"`           //操作人
	Stage    string    `json:"stage" form:"stage" gorm:"comment:'说明'" excel:"说明"`           //说明
	Remarks  string    `json:"remarks" form:"remarks" gorm:"comment:'备注'" excel:"备注"`       //备注
	//单据状态
}
