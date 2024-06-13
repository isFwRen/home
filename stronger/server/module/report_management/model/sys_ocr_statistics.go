package model

import (
	"server/module/sys_base/model"
	"time"
)

//----------------------------------------------------------------产量明细-----------------------------------------------------------------------------//

type OcrStatistics struct {
	model.Model
	// Code       string `json:"code" gorm:"comment:'工号'" excel:"工号"`
	BillID      string    `json:"bill_id" gorm:"comment:'单ID'" excel:"单ID"`
	BillNum     string    `json:"bill_num" gorm:"comment:'单号'" excel:"单号"`
	NickName    string    `json:"nickName" gorm:"comment:'姓名'" excel:"姓名"`
	JobNumber   string    `json:"jobNumber" gorm:"comment:'工号'" excel:"工号"`
	SubmitTime  time.Time `json:"submit_time" gorm:"comment:'提交时间'" excel:"提交时间" excelFormat:"2006-01-02 15:04:05"`
	ProCode     string    `json:"proCode" gorm:"comment:'项目编码'" excel:"项目编码"`
	Value       string    `json:"value" gorm:"comment:'ocr内容'" excel:"ocr内容"`
	ResultValue string    `json:"result_value" gorm:"comment:'最终值'" excel:"最终值"` //有效字符
	Code        string    `json:"code" gorm:"comment:'字段编码'" excel:"字段编码"`
	Name        string    `json:"name" gorm:"comment:'字段名称'" excel:"字段名称"`    //有效字符
	Disable     string    `json:"disable" gorm:"comment:'是否屏蔽'" excel:"是否屏蔽"` //有效字符
	Compare     string    `json:"compare" gorm:"comment:'数据对比'" excel:"数据对比"` //有效字符
	Rate        float64   `json:"rate" gorm:"comment:'准确率'" excel:"准确率(%)"`   //有效字符
	OcrType     string    `json:"ocr_type" gorm:"comment:'类型'" excel:"类型"`    //类型
	Pic         string    `json:"pic" gorm:"comment:'图片'"`                    //图片
}
