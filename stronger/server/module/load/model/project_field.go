package model

import (
	"server/module/sys_base/model"
	"time"
)

type ProjectField struct {
	// ID
	model.Model
	Name string `json:"name" form:"name" gorm:"comment:'字段名字'"` //字段名字
	Code string `json:"code" form:"code" gorm:"comment:'字段编码'"` //字段编码

	BillID  string `json:"billID" gorm:"comment:'单ID'"`   //单ID
	BlockID string `json:"blockID" gorm:"comment:'分块ID'"` //分块ID

	BlockIndex int `json:"blockIndex" gorm:"comment:'分块下标'"` //分块下标
	FieldIndex int `json:"fieldIndex" gorm:"comment:'字段下标'"` //字段下标

	// mark string "q11111"
	Op0Value string `json:"op0Value" gorm:"comment:'初审内容'"`                     //1码内容
	Op0Input string `json:"op0Input" gorm:"comment:'初审录入状态(yes|no|no_if|ocr)'"` //1码录入状态(yes|no|no_if|ocr)
	Op1Value string `json:"op1Value" gorm:"comment:'1码内容'"`                     //1码内容
	Op1Input string `json:"op1Input" gorm:"comment:'1码录入状态(yes|no|no_if|ocr)'"` //1码录入状态(yes|no|no_if|ocr)
	Op2Value string `json:"op2Value" gorm:"comment:'2码内容'"`                     //2码内容
	Op2Input string `json:"op2Input" gorm:"comment:'2码录入状态'"`                   //2码录入状态
	OpqValue string `json:"opqValue" gorm:"comment:'问题件内容'"`                    //问题件内容
	OpqInput string `json:"opqInput" gorm:"comment:'问题件录入状态'"`                  //问题件录入状态
	// OpDInput    string `json:"opDInput" gorm:"comment:'复核录入状态'"`                   //复核录入状态
	ResultValue string `json:"resultValue" gorm:"comment:'最终录入内容'"` //最终录入内容
	ResultInput string `json:"resultInput" gorm:"comment:'最终录入状态'"` //最终录入状态
	FinalValue  string `json:"finalValue" gorm:"comment:'结果值内容'"`   //结果值内容
	FinalInput  string `json:"finalInput" gorm:"comment:'结果值状态'"`   //结果值状态

	RightValue   string    `json:"rightValue" gorm:"comment:'客户反馈对的值'"`  //客户反馈对的值
	FeedbackDate time.Time `json:"feedbackDate" gorm:"comment:'客户反馈日期'"` //客户反馈日期
	IsPractice   bool      `json:"isPractice" gorm:"comment:'练习分块'"`     //练习分块(同步当前分块该字段的值)
	Issues       []Issue   `json:"issues" gorm:"foreignKey:FieldId"`     //问题件
	IsChange     bool      `json:"isChange" gorm:"isChange"`             //是否改了录入值
}

type Issue struct {
	FieldId string `json:"fieldId"` //字段id
	Type    string `json:"type"`    //类型
	Code    string `json:"code"`    //编码
	Message string `json:"message"` //详细信息
}
