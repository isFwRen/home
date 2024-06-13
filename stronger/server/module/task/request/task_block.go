package request

import (
	"server/module/sys_base/model"
	"time"
)

type TaskBlock struct {
	// ID
	model.Model
	BillID      string    `json:"billID" gorm:"comment:'单ID'"`             //单ID
	Name        string    `json:"name" gorm:"模板分块名字"`                      //模板分块名字
	Code        string    `json:"code" gorm:"模板分块编码"`                      //模板分块编码
	Op1Code     string    `json:"op1Code" gorm:"comment:'1码人员编号(0为系统录入)'"` //1码人员编号(0为系统录入)
	Op1ApplyAt  time.Time `json:"op1ApplyAt" gorm:"comment:'1码领取时间'"`      //1码领取时间
	Op1SubmitAt time.Time `json:"op1SubmitAt" gorm:"comment:'1码提交时间'"`     //1码提交时间
	Op2Code     string    `json:"op2Code" gorm:"comment:'2码人员编号'"`         //2码人员编号
	Op2ApplyAt  time.Time `json:"op2ApplyAt" gorm:"comment:'2码领取时间'"`      //2码领取时间
	Op2SubmitAt time.Time `json:"op2SubmitAt" gorm:"comment:'2码提交时间'"`     //2码提交时间
	OpQCode     string    `json:"opQCode" gorm:"comment:'问题件人员编号'"`        //问题件人员编号
	OpQApplyAt  time.Time `json:"opQApplyAt" gorm:"comment:'问题件领取时间'"`     //问题件领取时间
	OpQSubmitAt time.Time `json:"opQSubmitAt" gorm:"comment:'问题件提交时间'"`    //问题件提交时间
	OpDCode     string    `json:"opDCode" gorm:"comment:'复核人员编号'"`         //复核人员编号
	OpDApplyAt  time.Time `json:"opDApplyAt" gorm:"comment:'复核领取时间'"`      //复核领取时间
	OpDSubmitAt time.Time `json:"opDSubmitAt" gorm:"comment:'复核提交时间'"`     //复核提交时间

	Op     string      `json:"op" gorm:"工序"`     //工序
	Fields []TaskField `json:"fields" gorm:"字段"` //工序

}
