package request

import (
	modelBase "server/module/sys_base/model"
	"time"
)

type OutPutStatisticsDetailExport struct {
	ProCode   string `json:"proCode" form:"proCode"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
	Code      string `json:"code" form:"code"`
}

type OutPutStatisticsExport struct {
	PageInfo   modelBase.BasePageInfo `json:"page_info"`
	ProCode    string                 `json:"proCode" form:"proCode"`
	Code       string                 `json:"code" form:"code"`
	StartTime  string                 `json:"startTime" form:"startTime"`
	EndTime    string                 `json:"endTime" form:"endTime"`
	UpdateTime string                 `json:"updateTime" form:"updateTime"`
}

type OutPutStatisticsSearch struct {
	ProCode    string                 `json:"proCode" form:"proCode"`
	StartTime  string                 `json:"startTime" form:"startTime"`
	EndTime    string                 `json:"endTime" form:"endTime"`
	Code       string                 `json:"code" form:"code"`
	IsCheckAll int                    `json:"isCheckAll" form:"isCheckAll"` //全部true 或者 明细false
	UpdateTime string                 `json:"updateTime" form:"updateTime"`
	PageInfo   modelBase.BasePageInfo `json:"page_info"`
}

type GetOCROutPutStatisticsSearch struct {
	PageInfo    modelBase.BasePageInfo `json:"page_info"`
	ProCode     string                 `json:"proCode" form:"proCode"`
	StartTime   time.Time              `json:"startTime" form:"startTime"`
	EndTime     time.Time              `json:"endTime" form:"endTime"`
	FieldName   string                 `json:"fieldName" form:"fieldName"`
	BillNum     string                 `json:"billNum" form:"billNum"`
	Disable     string                 `json:"disable" form:"disable"`
	Compare     string                 `json:"compare" form:"compare"`
	ResultValue string                 `json:"resultValue" form:"resultValue"`
}

type UpdateCorrected struct {
	Id                   string  `json:"id"`
	ProCode              string  `json:"proCode"`
	Code                 string  `json:"code"` //修改人工号
	Name                 string  `json:"name"` //修改人姓名
	Op0AsTheBlock        float64 `json:"op0AsTheBlock"`
	Op0AsTheInvoice      float64 `json:"op0AsTheInvoice"`
	Op1NotExpenseAccount float64 `json:"op1NotExpenseAccount"`
	Op1ExpenseAccount    float64 `json:"op1ExpenseAccount"`
	Op2NotExpenseAccount float64 `json:"op2NotExpenseAccount"`
	Op2ExpenseAccount    float64 `json:"op2ExpenseAccount"`
	Question             float64 `json:"question"`
	StartTime            string  `json:"startTime"`
}

type InsertCorrected struct {
	Id                   string  `json:"id"`
	ProCode              string  `json:"proCode"`
	Op0AsTheBlock        float64 `json:"op0AsTheBlock"`
	Op0AsTheInvoice      float64 `json:"op0AsTheInvoice"`
	Op1NotExpenseAccount float64 `json:"op1NotExpenseAccount"`
	Op1ExpenseAccount    float64 `json:"op1ExpenseAccount"`
	Op2NotExpenseAccount float64 `json:"op2NotExpenseAccount"`
	Op2ExpenseAccount    float64 `json:"op2ExpenseAccount"`
	Question             float64 `json:"question"`
	StartTime            string  `json:"startTime"`
}

type DeleteCorrectedArr struct {
	Ids []string `json:"ids"`
}
