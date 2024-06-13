package model

type AbnormalBill struct {
	CreateAt       string `json:"createAt" excel:"日期"`         //日期
	BillName       string `json:"billName" excel:"单号"`         //单号
	Agency         string `json:"agency" excel:"机构号"`          //机构号
	AbnormalReason string `json:"abnormalReason" excel:"异常原因"` //异常原因
	UploadAt       string `json:"uploadAt" excel:"回传时间"`       //回传时间
	Stage          int    `json:"stage"`                       //录入状态
	StageStr       string `json:"stageStr" excel:"录入状态"`       //录入状态字符串
}
