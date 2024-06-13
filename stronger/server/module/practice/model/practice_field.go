package model

import (
	"server/module/sys_base/model"
)

type PracticeProjectField struct {
	model.Model
	Name        string `json:"name" form:"name" gorm:"comment:'字段名字'"`             //字段名字
	Code        string `json:"code" form:"code" gorm:"comment:'字段编码'"`             //字段编码
	BillID      string `json:"billID" gorm:"comment:'单ID'"`                        //单ID
	BlockID     string `json:"blockID" gorm:"comment:'分块ID'"`                      //分块ID
	BlockIndex  int    `json:"blockIndex" gorm:"comment:'分块下标'"`                   //分块下标
	FieldIndex  int    `json:"fieldIndex" gorm:"comment:'字段下标'"`                   //字段下标
	ResultValue string `json:"resultValue" gorm:"comment:'最终录入内容'"`                //最终录入内容
	ResultInput string `json:"resultInput" gorm:"comment:'最终录入状态'"`                //最终录入状态
	Analysis    string `json:"analysis" gorm:"comment:'解析'"`                       //解析
	Op1Value    string `json:"op1Value" gorm:"comment:'1码内容'"`                     //1码内容
	Op1Input    string `json:"op1Input" gorm:"comment:'1码录入状态(yes|no|no_if|ocr)'"` //1码录入状态(yes|no|no_if|ocr)

}
