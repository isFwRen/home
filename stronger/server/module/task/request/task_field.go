package request

import (
	"server/module/sys_base/model"
)

type TaskField struct {
	// ID
	model.Model
	Name string `json:"name" form:"name" gorm:"comment:'字段名字'"` //字段名字
	Code string `json:"code" form:"code" gorm:"comment:'字段编码'"` //字段编码

	BillID  string `json:"billID" gorm:"comment:'单ID'"`   //单ID
	BlockID string `json:"blockID" gorm:"comment:'分块ID'"` //分块ID

	BlockIndex int `json:"blockIndex" gorm:"comment:'分块下标'"` //分块下标
	FieldIndex int `json:"fieldIndex" gorm:"comment:'字段下标'"` //字段下标

	// mark string "q11111"

	Op1Value string `json:"op1Value" gorm:"comment:'1码内容'"`                     //1码内容
	Op1Input string `json:"op1Input" gorm:"comment:'1码录入状态(yes|no|no_if|ocr)'"` //1码录入状态(yes|no|no_if|ocr)
	Op2Value string `json:"op2Value" gorm:"comment:'2码内容'"`                     //2码内容
	Op2Input string `json:"op2Input" gorm:"comment:'2码录入状态'"`                   //2码录入状态
	OpQValue string `json:"opQValue" gorm:"comment:'问题件内容'"`                    //问题件内容
	OpQInput string `json:"opQInput" gorm:"comment:'问题件录入状态'"`                  //问题件录入状态
	// OpDValue    string `json:"opDValue" gorm:"comment:'复核录入内容'"`                   //复核录入内容
	// OpDInput    string `json:"opDInput" gorm:"comment:'复核录入状态'"`                   //复核录入状态
	ResultValue string `json:"resultValue" gorm:"comment:'最终录入内容'"` //最终录入内容
	ResultInput string `json:"resultInput" gorm:"comment:'最终录入状态'"` //最终录入状态
	FinalValue  string `json:"finalValue" gorm:"comment:'结果值内容'"`   //结果值内容
	FinalInput  string `json:"finalInput" gorm:"comment:'结果值状态'"`   //结果值状态

}
