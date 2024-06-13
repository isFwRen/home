package response

import (
	"github.com/lib/pq"
	modelBase "server/module/sys_base/model"
)

type WrongRes struct {
	List    interface{} `json:"list"`
	Message string      `json:"message"`
}

type GetWrong struct {
	//Id             string         `json:"id"`
	modelBase.Model
	SubmitDay      string         `json:"submitDay" example:"提交时间"`
	Code           string         `json:"code" example:"工号"`
	NickName       string         `json:"nickName" example:"姓名"`
	Op             string         `json:"op" example:"工序"`
	BillName       string         `json:"billName" example:"单名"`
	BillNum        string         `json:"billNum" example:"单号"`
	CaseNumber     string         `json:"caseNumber" example:"案件号"`
	Agency         string         `json:"agency" example:"机构号"`
	Type           string         `json:"type" example:"模板类型(MB001,MB002)"`
	Types          string         `json:"types" example:"单子类型,比如混合型, 医疗等"`
	Block          string         `json:"block" example:"字段所处分块"`
	FieldId        string         `json:"fieldId" example:"字段id"`
	FieldCode      string         `json:"fieldCode" example:"字段编码"`
	FieldName      string         `json:"fieldName" example:"字段名称"`
	Path           string         `json:"path" example:"单子图片存放路径"`
	Picture        pq.StringArray `json:"picture" gorm:"type:varchar(100)[];comment:'web截图位置'" example:"图片名称"`
	Wrong          string         `json:"wrong" example:"错误内容"`
	Right          string         `json:"right" example:"正确内容"`
	IsComplain     bool           `json:"isComplain"` //是否申诉
	IsWrongConfirm bool           `json:"isWrongConfirm" example:"是否通过"`
	IsOcr          bool           `json:"isOcr" example:"ocr识别的"`
	IsOperationLog string         `json:"isOperationLog"`
	IsAudit        bool           `json:"isAudit"`
}

type GetWrongAnalysis struct {
	StatisticalTime       string  `json:"statisticalTime" excel:"日期"`
	Code                  string  `json:"code" example:"工号" excel:"工号"`
	NickName              string  `json:"nickName" example:"姓名" excel:"姓名"`
	WrongNumber           float64 `json:"wrongNumber" example:"错误数量" excel:"错误数量"`
	TheNumberOfComplaints float64 `json:"theNumberOfComplaints" example:"申诉数量" excel:"申诉数量"`
	TheComplaintRate      string  `json:"theComplaintRate" example:"申诉率" excel:"申诉率"`
	ThroughTheNumber      float64 `json:"throughTheNumber" example:"通过数量" excel:"通过数量"`
	ThePassRate           string  `json:"thePassRate" example:"通过率" excel:"通过率"`
	NonPassingQuantity    float64 `json:"nonPassingQuantity" example:"不通过数量" excel:"不通过数量"`
	UnqualifiedRate       string  `json:"unqualifiedRate" example:"不通过率" excel:"不通过率"`
}
