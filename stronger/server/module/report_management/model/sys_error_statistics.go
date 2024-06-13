package model

import (
	"github.com/lib/pq"
	modelBase "server/module/sys_base/model"
	"time"
)

type WrongSearch struct {
	PageInfo  modelBase.BasePageInfo `json:"page_info"`
	ProCode   string                 `json:"proCode" form:"proCode"`     //项目编码
	StartTime string                 `json:"startTime" form:"startTime"` //开始时间
	EndTime   string                 `json:"endTime" form:"endTime"`     //结束时间
	Code      string                 `json:"code" form:"code"`           //工号
	NickName  string                 `json:"nickName" form:"nickName"`   //姓名
	FieldName string                 `json:"fieldName" form:"fieldName"` //字段名称
	Op        string                 `json:"op" form:"op"`               //工序
	Complaint string                 `json:"complaint" form:"complaint"` //申诉
	Confirm   string                 `json:"confirm" form:"confirm"`     //审核是否通过
	IsAudit   bool                   `json:"isAudit" form:"isAudit"`
}

type WrongExport struct {
	ProCode   string    `json:"proCode" form:"proCode"  binding:"required"`     //项目编码
	StartTime time.Time `json:"startTime" form:"startTime"  binding:"required"` //开始时间
	EndTime   time.Time `json:"endTime" form:"endTime"  binding:"required"`     //结束时间
}

type Wrong struct {
	//Id             string         `json:"id"`
	modelBase.Model
	SubmitDay      time.Time      `json:"submitDay" example:"提交时间" flag:"time"  excel:"日期"`
	Code           string         `json:"code" example:"工号" excel:"工号"`
	NickName       string         `json:"nickName" example:"姓名"  excel:"姓名"`
	Op             string         `json:"op" example:"工序"`
	BillName       string         `json:"billName" example:"单名" excel:"案件号"`
	BillNum        string         `json:"billNum" example:"单号"`
	CaseNumber     string         `json:"caseNumber" example:"案件号"`
	Agency         string         `json:"agency" example:"机构号"  excel:"机构号"`
	Type           string         `json:"type" example:"模板类型(MB001,MB002)"`
	Types          string         `json:"types" example:"单子类型,比如混合型, 医疗等"`
	Block          string         `json:"block" example:"字段所处分块"`
	FieldId        string         `json:"fieldId" example:"字段id"`
	FieldCode      string         `json:"fieldCode" example:"字段编码"`
	FieldName      string         `json:"fieldName" example:"字段名称"  excel:"字段名称"`
	Path           string         `json:"path" example:"单子图片存放路径"`
	Picture        pq.StringArray `json:"picture" gorm:"type:varchar(100)[];comment:'web截图位置'" example:"图片名称"`
	Wrong          string         `json:"wrong" example:"错误内容"  excel:"错误内容"`
	Right          string         `json:"right" example:"正确内容"  excel:"正确内容"`
	IsComplain     bool           `json:"isComplain"` //是否申诉
	IsWrongConfirm bool           `json:"isWrongConfirm" example:"是否通过"`
	IsOcr          bool           `json:"isOcr" example:"ocr识别的"`
	BillId         string         `json:"billId" example:"billId"`
	IsBaoXiaoBlock string         `json:"isBaoXiaoBlock"` //0:不区分报销与非报销 1:报销 2:非报销
	IsOperationLog string         `json:"isOperationLog"`
	IsAudit        bool           `json:"isAudit"`
}

type WrongExportResp struct {
	SubmitDay time.Time `json:"submitDay" example:"提交时间" flag:"time"  excel:"日期"`
	Code      string    `json:"code" example:"工号" excel:"工号"`
	NickName  string    `json:"nickName" example:"姓名"  excel:"姓名"`
	BillName  string    `json:"billName" example:"单名" excel:"案件号"`
	BillNum   string    `json:"billNum" example:"单号"`
	Agency    string    `json:"agency" example:"机构号"  excel:"机构号"`
	FieldName string    `json:"fieldName" example:"字段名称"  excel:"字段名称"`
	Wrong     string    `json:"wrong" example:"错误内容"  excel:"错误内容"`
	Right     string    `json:"right" example:"正确内容"  excel:"正确内容"`
}

type WrongAnalysis struct {
	modelBase.Model
	//ID                    string    `json:"id"`
	StatisticalTime       time.Time `json:"statisticalTime" excel:"日期" flag:"time"`
	Code                  string    `json:"code" excel:"工号"`
	NickName              string    `json:"nickName" excel:"姓名"`
	WrongNumber           float64   `json:"wrongNumber" excel:"错误数量"`
	TheNumberOfComplaints float64   `json:"theNumberOfComplaints" excel:"申诉数量"`
	TheComplaintRate      float64   `json:"theComplaintRate" excel:"申诉率" flag:"rate"`
	ThroughTheNumber      float64   `json:"throughTheNumber" excel:"通过数量"`
	ThePassRate           float64   `json:"thePassRate" excel:"通过率" flag:"rate"`
	NonPassingQuantity    float64   `json:"nonPassingQuantity" excel:"不通过数量"`
	UnqualifiedRate       float64   `json:"unqualifiedRate" excel:"不通过率" flag:"rate"`
}

type WrongAnalysisExport struct {
	StatisticalTime       time.Time `json:"statisticalTime" excel:"日期" flag:"time"`
	Code                  string    `json:"code" excel:"工号"`
	NickName              string    `json:"nickName" excel:"姓名"`
	WrongNumber           float64   `json:"wrongNumber" excel:"错误数量"`
	TheNumberOfComplaints float64   `json:"theNumberOfComplaints" excel:"申诉数量"`
	TheComplaintRate      string    `json:"theComplaintRate" excel:"申诉率" flag:"rate"`
	ThroughTheNumber      float64   `json:"throughTheNumber" excel:"通过数量"`
	ThePassRate           string    `json:"thePassRate" excel:"通过率" flag:"rate"`
	NonPassingQuantity    float64   `json:"nonPassingQuantity" excel:"不通过数量"`
	UnqualifiedRate       string    `json:"unqualifiedRate" excel:"不通过率" flag:"rate"`
}

type OcrAnalysis struct {
	StatisticalTime time.Time `json:"statisticalTime"` //统计时间(默认格式2021-01-01)
	CaseNumber      string    `json:"caseNumber"`      //案件号
	Agency          string    `json:"agency"`          //机构号
	FieldName       string    `json:"fieldName"`       //字段名称
	Wrong           string    `json:"wrong"`           //错误内容
	Right           string    `json:"right"`           //正确内容
}
