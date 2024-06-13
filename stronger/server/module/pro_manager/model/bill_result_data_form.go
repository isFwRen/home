/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/30 9:42 上午
 */

package model

import (
	"server/module/load/model"
	"time"
)

type EditBillResultData struct {
	FieldId    string `json:"fieldId" form:"fieldId"`                    //字段id
	FieldValue string `json:"fieldValue" form:"fieldValue"`              //字段录入值
	FieldInput string `json:"fieldInput" form:"fieldInput"`              //字段录入状态
	BillId     string `json:"billId" form:"billId" binding:"required"`   //单据id
	ProCode    string `json:"proCode" form:"proCode" binding:"required"` //项目编码
	//BasePageInfo model.BasePageInfo `json:"basePageInfo" form:"basePageInfo"`
}

type Field struct {
	FieldId   string `json:"fieldId" form:"fieldId"`     //字段id
	Code      string `json:"code" form:"code"`           //字段code
	Name      string `json:"name" form:"name"`           //字段name
	BeforeVal string `json:"beforeVal" form:"beforeVal"` //字段修改前的值
	EndVal    string `json:"endVal" form:"endVal"`       //字段修改后的值
}
type EditBillResultDataManyFields struct {
	Fields     []model.ProjectField `json:"fields" form:"fields"`
	EditFields []Field              `json:"editFields" form:"editFields"`
	EditType   int                  `json:"editType" form:"editType"`                  //查看类型
	BillId     string               `json:"billId" form:"billId" binding:"required"`   //单据id
	ProCode    string               `json:"proCode" form:"proCode" binding:"required"` //项目编码
	//BasePageInfo model.BasePageInfo `json:"basePageInfo" form:"basePageInfo"`
}

type SetPracticeForm struct {
	BlockIds   []string `json:"blockIds" form:"blockIds" binding:"required"` //分块id
	ProCode    string   `json:"proCode" form:"proCode" binding:"required"`   //项目编码
	IsPractice bool     `json:"isPractice" form:"isPractice"`                //练习
}

type EditFeedbackVal struct {
	FieldId    string    `json:"fieldId" form:"fieldId" binding:"required"`       //字段id
	FieldValue string    `json:"fieldValue" form:"fieldValue" binding:"required"` //字段正确值
	ProCode    string    `json:"proCode" form:"proCode" binding:"required"`       //项目编码
	EditDate   time.Time `json:"editDate" form:"editDate" binding:"required"`     //反馈日期
}
