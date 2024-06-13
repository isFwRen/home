package model

import (
	"github.com/lib/pq"
	modelbase "server/module/sys_base/model"
	"time"
)

type Quality struct {
	modelbase.Model
	Month              string         `json:"month" excel:"月份"`
	ProCode            string         `json:"proCode" form:"proCode" excel:"项目编码"`
	BillName           string         `json:"billName" excel:"案件号"`
	FeedbackDate       time.Time      `json:"feedbackDate" excel:"反馈日期" time:"y"`
	EntryDate          time.Time      `json:"entryDate" excel:"录入日期" time:"y"`
	WrongFieldName     string         `json:"wrongFieldName" excel:"错误字段"`
	Right              string         `json:"right" excel:"正确值"`
	Wrong              string         `json:"wrong" excel:"错误值"`
	Op0ResponsibleCode string         `json:"op0ResponsibleCode" excel:"初审责任人工号"`
	Op0ResponsibleName string         `json:"op0ResponsibleName" excel:"初审责任人姓名"`
	Op1ResponsibleCode string         `json:"op1ResponsibleCode" excel:"一码责任人工号"`
	Op1ResponsibleName string         `json:"op1ResponsibleName" excel:"一码责任人姓名"`
	Op2ResponsibleCode string         `json:"op2ResponsibleCode" excel:"二码责任人工号"`
	Op2ResponsibleName string         `json:"op2ResponsibleName" excel:"二码责任人姓名"`
	OpqResponsibleCode string         `json:"opqResponsibleCode" excel:"问题件责任人工号"`
	OpqResponsibleName string         `json:"opqResponsibleName" excel:"问题件责任人姓名"`
	ImagePath          pq.StringArray `json:"imagePath" gorm:"type:varchar(255)[]" excel:"影像"`
}

type Qualities struct {
	Month              string         `json:"month" excel:"月份"`
	ProCode            string         `json:"proCode" form:"proCode"`
	BillName           string         `json:"billName" excel:"案件号"`
	FeedbackDate       time.Time      `json:"feedbackDate" excel:"反馈日期"`
	EntryDate          string         `json:"entryDate" excel:"录入日期"`
	WrongFieldName     string         `json:"wrongFieldName" excel:"错误字段"`
	Right              string         `json:"right" excel:"正确值"`
	Wrong              string         `json:"wrong" excel:"错误值"`
	Op0ResponsibleCode string         `json:"op0ResponsibleCode" excel:"初审责任人工号"`
	Op0ResponsibleName string         `json:"op0ResponsibleName" excel:"初审责任人姓名"`
	Op1ResponsibleCode string         `json:"op1ResponsibleCode" excel:"一码责任人工号"`
	Op1ResponsibleName string         `json:"op1ResponsibleName" excel:"一码责任人姓名"`
	Op2ResponsibleCode string         `json:"op2ResponsibleCode" excel:"二码责任人工号"`
	Op2ResponsibleName string         `json:"op2ResponsibleName" excel:"二码责任人姓名"`
	OpqResponsibleCode string         `json:"opqResponsibleCode" excel:"问题件责任人工号"`
	OpqResponsibleName string         `json:"opqResponsibleName" excel:"问题件责任人姓名"`
	ImagePath          pq.StringArray `json:"imagePath" gorm:"type:varchar(255)[]"`
}

type QualitiesAddEdit struct {
	ProCode            string         `json:"proCode" form:"proCode"`
	BillName           string         `json:"billName" excel:"案件号"`
	FeedbackDate       time.Time      `json:"feedbackDate" excel:"反馈日期" copier:"-"`
	EntryDate          string         `json:"entryDate" excel:"录入日期"`
	WrongFieldName     string         `json:"wrongFieldName" excel:"错误字段"`
	Right              string         `json:"right" excel:"正确值"`
	Wrong              string         `json:"wrong" excel:"错误值"`
	Op0ResponsibleCode string         `json:"op0ResponsibleCode" excel:"初审责任人工号"`
	Op0ResponsibleName string         `json:"op0ResponsibleName" excel:"初审责任人姓名"`
	Op1ResponsibleCode string         `json:"op1ResponsibleCode" excel:"一码责任人工号"`
	Op1ResponsibleName string         `json:"op1ResponsibleName" excel:"一码责任人姓名"`
	Op2ResponsibleCode string         `json:"op2ResponsibleCode" excel:"二码责任人工号"`
	Op2ResponsibleName string         `json:"op2ResponsibleName" excel:"二码责任人姓名"`
	OpqResponsibleCode string         `json:"opqResponsibleCode" excel:"问题件责任人工号"`
	OpqResponsibleName string         `json:"opqResponsibleName" excel:"问题件责任人姓名"`
	ImagePath          pq.StringArray `json:"imagePath" gorm:"type:varchar(255)[]"`
}
