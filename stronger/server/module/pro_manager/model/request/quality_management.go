package request

import (
	"github.com/lib/pq"
	modelBase "server/module/sys_base/model"
)

type QualityRequest struct {
	PageInfo        modelBase.BasePageInfo `json:"page_info"`
	Month           string                 `json:"month" form:"month"`
	ProCode         string                 `json:"proCode" form:"proCode"`
	BillName        string                 `json:"billName" form:"billName"`
	WrongFieldName  string                 `json:"wrongFieldName" form:"wrongFieldName"`
	ResponsibleName string                 `json:"responsibleName" form:"responsibleName"`
	//StartTime       string                 `json:"startTime" form:"startTime"`
	//EndTime         string                 `json:"endTime" form:"endTime"`
}

type QualitiesAeRequest struct {
	Id                 string         `json:"id" form:"id" db:"none"`
	Month              string         `json:"month" form:"month" db:"month"`
	ProCode            string         `json:"proCode" form:"proCode" db:"pro_code"`
	BillName           string         `json:"billName" form:"billName" excel:"案件号" db:"bill_name"`
	FeedbackDate       string         `json:"feedbackDate" form:"feedbackDate" excel:"反馈日期" copier:"-" db:"feedback_date"`
	EntryDate          string         `json:"entryDate" form:"entryDate" excel:"录入日期" db:"entry_date"`
	WrongFieldName     string         `json:"wrongFieldName" form:"wrongFieldName" excel:"错误字段" db:"wrong_field_name"`
	Right              string         `json:"right" form:"right" excel:"正确值" db:"right"`
	Wrong              string         `json:"wrong" form:"wrong" excel:"错误值" db:"wrong"`
	Op0ResponsibleCode string         `json:"op0ResponsibleCode" form:"op0ResponsibleCode" excel:"初审责任人工号" db:"op0_responsible_code"`
	Op0ResponsibleName string         `json:"op0ResponsibleName" form:"op0ResponsibleName" excel:"初审责任人姓名" db:"op0_responsible_name"`
	Op1ResponsibleCode string         `json:"op1ResponsibleCode" form:"op1ResponsibleCode" excel:"一码责任人工号" db:"op1_responsible_code"`
	Op1ResponsibleName string         `json:"op1ResponsibleName" form:"op1ResponsibleName" excel:"一码责任人姓名" db:"op1_responsible_name"`
	Op2ResponsibleCode string         `json:"op2ResponsibleCode" form:"op2ResponsibleCode" excel:"二码责任人工号" db:"op2_responsible_code"`
	Op2ResponsibleName string         `json:"op2ResponsibleName" form:"op2ResponsibleName" excel:"二码责任人姓名" db:"op2_responsible_name"`
	OpqResponsibleCode string         `json:"opqResponsibleCode" form:"opqResponsibleCode" excel:"问题件责任人工号" db:"opq_responsible_code"`
	OpqResponsibleName string         `json:"opqResponsibleName" form:"opqResponsibleName" excel:"问题件责任人姓名" db:"opq_responsible_name"`
	ImagePath          pq.StringArray `json:"imagePath" db:"image_path"`
}

type QualitiesExport struct {
	ProCode string `json:"proCode" form:"proCode"`
	Month   string `json:"month" form:"month"`
}

type QualitiesReqDemo struct {
	Type string `json:"type" binding:"required,oneof=basicInfo benefitInfo billInfo insuranceInfo countInfo"` //类型basicInfo benefitInfo billInfo insuranceInfo countInfo
	ID   string `json:"id" binding:"required"`
	Data string `json:"data" ` //json字符串数据
}
