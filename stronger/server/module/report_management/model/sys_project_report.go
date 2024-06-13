package model

import (
	"server/module/sys_base/model"
	"time"
)

// BusinessDetails 业务明细表
type BusinessDetails struct {
	model.Model
	CreateAt            string  `json:"createAt" excel:"日期"`               //日期
	BatchNum            string  `json:"batchNum" excel:"批次号"`              //批次号
	SaleChannel         string  `json:"saleChannel" excel:"销售渠道"`          //销售渠道
	BillName            string  `json:"billName" excel:"案件号"`              //案件号
	PictureNumber       int     `json:"pictureNumber" excel:"影像数量"`        //影像数量
	Agency              string  `json:"agency" excel:"机构"`                 //机构
	InsuranceType       string  `json:"insuranceType" excel:"医保类型"`        //医保类型
	ScanAt              string  `json:"scanAt" excel:"扫描时间"`               //扫描时间
	DownloadAt          string  `json:"downloadAt" excel:"下载时间"`           //下载时间
	ExportAt            string  `json:"exportAt" excel:"导出时间"`             //导出时间
	FirstUploadAt       string  `json:"uploadAt" excel:"初次回传时间"`           //初次回传时间
	LatestUploadAt      string  `json:"latestUploadAt" excel:"最新回传时间"`     //最新回传时间
	AtTheLateUploadAt   string  `json:"atTheLateUploadAt" excel:"最晚回传时间"`  //最晚回传时间
	LateTime            string  `json:"lateTime" excel:"延迟时间"`             //延迟时间
	Status              int     `json:"status" excel:"案件状态"`               //案件状态
	IsTheTimeOut        bool    `json:"isTheTimeOut" excel:"是否超时"`         //是否超时
	WorkTime            string  `json:"workTime" excel:"处理时长"`             //处理时长
	Stage               int     `json:"stage" excel:"录入状态"`                //录入状态
	ClaimType           int     `json:"claimType" excel:"理赔类型"`            //理赔类型
	CountMoney          float32 `json:"countMoney" excel:"账单金额"`           //账单金额
	InvoiceNum          int     `json:"invoiceNum" excel:"发票数量"`           //发票数量
	ListingNum          int     `json:"listingNum" excel:"清单数量"`           //清单数量
	QuestionNum         int     `json:"questionNum" excel:"问题件数量"`         //问题件数量
	DiseaseDiagnosis    string  `json:"diseaseDiagnosis" excel:"疾病诊断"`     //疾病诊断
	FieldCharacter      int     `json:"fieldCharacter" excel:"录入字符数"`      //录入字符数
	SettlementCharacter int     `json:"settlementCharacter" excel:"结算字符数"` //结算字符数
	SettlementMoney     string  `json:"settlementMoney" excel:"结算金额"`      //结算金额
	RequirementOfAging  string  `json:"requirementOfAging" excel:"时效考核要求"` //时效考核要求
	TheNumberOfCase     int     `json:"theNumberOfCase" excel:"案件次数"`      //案件次数  推送的次数，有没有重复单号
	Op0Entry            string  `json:"op0Entry" excel:"初审进入时间"`           //初审进入时间
	Op0End              string  `json:"op0End" excel:"初审结束时间"`             //初审结束时间
	Op0WorkTime         string  `json:"op0WorkTime" excel:"初审处理时间"`        //初审处理时间
	Op1Entry            string  `json:"op1Entry" excel:"一码进入时间"`           //一码进入时间
	Op1End              string  `json:"op1End" excel:"一码结束时间"`             //一码结束时间
	Op1WorkTime         string  `json:"op1WorkTime" excel:"一码处理时间"`        //一码处理时间
	Op2Entry            string  `json:"op2Entry" excel:"二码进入时间"`           //二码进入时间
	Op2End              string  `json:"op2End" excel:"二码结束时间"`             //二码结束时间
	Op2WorkTime         string  `json:"op2WorkTime" excel:"二码处理时间"`        //二码处理时间
	OpQEntry            string  `json:"opQEntry" excel:"问题件进入时间"`          //问题件进入时间
	OpQEnd              string  `json:"opQEnd" excel:"问题件结束时间"`            //问题件结束时间
	OpQWorkTime         string  `json:"opQWorkTime" excel:"问题件处理时间"`       //问题件处理时间
}

type BusinessDetailsExport struct {
	CreateAt    string `json:"createAt" excel:"日期"`      //日期
	BatchNum    string `json:"batchNum" excel:"批次号"`     //批次号
	SaleChannel string `json:"saleChannel" excel:"销售渠道"` //销售渠道
	//BillName            string  `json:"billName" excel:"案件号"`              //案件号
	BillNum             string  `json:"billNum" excel:"案件号"`                   //案件号
	PictureNumber       int     `json:"pictureNumber" excel:"影像数量"`            //影像数量
	Agency              string  `json:"agency" excel:"机构"`                     //机构
	InsuranceType       string  `json:"insuranceType" excel:"单证类型"`            //医保类型
	ScanAt              string  `json:"scanAt" excel:"扫描时间"`                   //扫描时间
	DownloadAt          string  `json:"downloadAt" excel:"下载时间"`               //下载时间
	FirstExportAt       string  `json:"firstExportAt" excel:"初次导出时间"`          //初次导出时间
	ExportAt            string  `json:"exportAt" excel:"导出时间"`                 //导出时间
	FirstUploadAt       string  `json:"uploadAt" excel:"初次回传时间"`               //初次回传时间
	LatestUploadAt      string  `json:"latestUploadAt" excel:"最新回传时间"`         //最新回传时间
	AtTheLateUploadAt   string  `json:"atTheLateUploadAt" excel:"最晚回传时间"`      //最晚回传时间
	LateTime            string  `json:"lateTime" excel:"延迟时间"`                 //延迟时间
	Status              string  `json:"status" excel:"案件状态"`                   //案件状态
	IsTheTimeOut        string  `json:"isTheTimeOut" excel:"是否超时"`             //是否超时
	WorkTime            string  `json:"workTime" excel:"处理时长"`                 //处理时长
	DurationRange       string  `json:"durationRange" excel:"时长范围"`            //时长范围
	Stage               string  `json:"stage" excel:"录入状态"`                    //录入状态
	QualityUserCode     string  `json:"qualityUserCode" excel:"质检人员工号"`        //质检人员工号
	QualityUserName     string  `json:"qualityUserName" excel:"质检人员姓名"`        //质检人员姓名
	ClaimType           string  `json:"claimType" excel:"理赔类型"`                //理赔类型
	CountMoney          float64 `json:"countMoney" excel:"账单金额"`               //账单金额
	InvoiceNum          int     `json:"invoiceNum" excel:"发票数量"`               //发票数量
	ListingNum          int     `json:"listingNum" excel:"清单数量"`               //清单数量
	QuestionNum         int     `json:"questionNum" excel:"问题件数量"`             //问题件数量
	DiseaseDiagnosis    string  `json:"diseaseDiagnosis" excel:"疾病诊断"`         //疾病诊断
	FieldCharacter      int     `json:"fieldCharacter" excel:"录入字符数"`          //录入字符数
	SettlementCharacter int     `json:"settlementCharacter" excel:"结算字符数"`     //结算字符数
	SettlementMoney     string  `json:"settlementMoney" excel:"结算金额"`          //结算金额
	RequirementOfAging  string  `json:"requirementOfAging" excel:"时效考核要求"`     //时效考核要求
	TheNumberOfCase     int     `json:"theNumberOfCase" excel:"案件次数"`          //案件次数  推送的次数，有没有重复单号
	Op0Entry            string  `json:"op0Entry" excel:"初审进入时间"`               //初审进入时间
	Op0End              string  `json:"op0End" excel:"初审结束时间"`                 //初审结束时间
	Op0WorkTime         string  `json:"op0WorkTime" excel:"初审处理时间"`            //初审处理时间
	Op1Entry            string  `json:"op1Entry" excel:"一码进入时间"`               //一码进入时间
	Op1End              string  `json:"op1End" excel:"一码结束时间"`                 //一码结束时间
	Op1WorkTime         string  `json:"op1WorkTime" excel:"一码处理时间"`            //一码处理时间
	Op2Entry            string  `json:"op2Entry" excel:"二码进入时间"`               //二码进入时间
	Op2End              string  `json:"op2End" excel:"二码结束时间"`                 //二码结束时间
	Op2WorkTime         string  `json:"op2WorkTime" excel:"二码处理时间"`            //二码处理时间
	OpQEntry            string  `json:"opQEntry" excel:"问题件进入时间"`              //问题件进入时间
	OpQEnd              string  `json:"opQEnd" excel:"问题件结束时间"`                //问题件结束时间
	OpQWorkTime         string  `json:"opQWorkTime" excel:"问题件处理时间"`           //问题件处理时间
	Remark              string  `json:"remark" excel:"备注"`                     //备注
	ExportStage         string  `json:"exportStage" excel:"首次导出时候的状态"`         //导出时候的状态
	BillType            string  `json:"billType" form:"billType" excel:"单据类型"` //单据类型
	Op1AndOp2Duration   string  `json:"op1AndOp2Duration" excel:"一二码处理时长"`     //一二码处理时长
	ApplicationChar     int     `json:"applicationChar" excel:"申请书录入字符"`       //申请书录入字符
	InvoiceChar         int     `json:"invoiceChar" excel:"发票录入字符"`            //发票录入字符
	ListChar            int     `json:"listChar" excel:"清单录入字符"`               //清单录入字符
	OtherChar           int     `json:"otherChar" excel:"其他录入字符"`              //其他录入字符

	BillTypeOS           int `json:"billTypeOS" excel:"票据类型(门诊)"`          //票据类型(门诊)
	BillTypeIH           int `json:"billTypeIH" excel:"票据类型(住院)"`          //票据类型(住院)
	ElectronBillCount    int `json:"electronBillCount" excel:"电子票据数量"`     //电子票据数量
	NonElectronBillCount int `json:"nonElectronBillCount" excel:"非电子票据数量"` //非电子票据数量
	InvoiceSum           int `json:"invoiceSum" excel:"发票张数汇总"`            //发票张数汇总

}

// CharSum 字符统计表
type CharSum struct {
	model.Model
	ProCode               string    `json:"proCode" excel:"项目编码"`                           //项目编码
	SumDate               time.Time `json:"sumDate" excel:"统计的日期" excelFormat:"2006-01-02"` //统计的日期
	BillCount             int64     `json:"billCount" excel:"案件量"`                          //案件量
	CharCount             int       `json:"charCount" excel:"结算字符数"`                        //结算字符数
	InputCharCount        int       `json:"inputCharCount" excel:"录入字符数"`                   //录入字符数
	AverageCharCount      float64   `json:"averageCharCount" excel:"平均结算字符数"`               //平均结算字符数
	AverageInputCharCount float64   `json:"averageInputCharCount" excel:"平均录入字符数"`          //平均录入字符数
	CharPercent           float64   `json:"charPercent" excel:"结算字符数与录入字符数的比例"`             //结算字符数与录入字符数的比例

	//
	StaffInputCount        int     `json:"staffInputCount" excel:"员工字符"`          //员工字符
	AverageStaffInputCount float64 `json:"averageStaffInputCount" excel:"员工平均字符"` //员工平均字符

	SettleStaffPercent float64 `json:"settleStaffPercent" excel:"结算字符数与员工字符数的比例"` //结算字符数与员工字符数的比例

	//Acc=AverageCharCount His=health insurance
	ClaimsAccHasHis  float64 `json:"claimsAccHasHis" excel:"理赔类型-平均结算字符数:有医保"`  //理赔类型-平均结算字符数:有医保
	ClaimsAccNoneHis float64 `json:"claimsAccNoneHis" excel:"理赔类型-平均结算字符数:无医保"` //理赔类型-平均结算字符数:无医保
	ClaimsAccMixHis  float64 `json:"claimsAccMixHis" excel:"理赔类型-平均结算字符数:混合型"`  //理赔类型-平均结算字符数:混合型

	//Aic = AverageInputCharCount
	ClaimsAicHasHis  float64 `json:"claimsAicHasHis" excel:"理赔类型-平均录入字符：有医保"`  //理赔类型-平均录入字符：有医保
	ClaimsAicNoneHis float64 `json:"claimsAicNoneHis" excel:"理赔类型-平均录入字符：无医保"` //理赔类型-平均录入字符：无医保
	ClaimsAicMixHis  float64 `json:"claimsAicMixHis" excel:"理赔类型-平均录入字符：混合型"`  //理赔类型-平均录入字符：混合型

	// Doc = 单证, Ih = in hospital ,Op = outpatient
	DocAccOp  float64 `json:"docAccOp" excel:"单证类型-平均结算字符数：门诊"`     //单证类型-平均结算字符数：门诊
	DocAccIh  float64 `json:"docAccIh" excel:"单证类型-平均结算字符数：住院"`     //单证类型-平均结算字符数：住院
	DocAccMix float64 `json:"docAccMix" excel:"单证类型-平均结算字符数：门诊+住院"` //单证类型-平均结算字符数：门诊+住院

	DocAicOp  float64 `json:"docAicOp" excel:"单证类型-平均录入字符：门诊"`     //单证类型-平均录入字符：门诊
	DocAicIh  float64 `json:"docAicIh" excel:"单证类型-平均录入字符：住院"`     //单证类型-平均录入字符：住院
	DocAicMix float64 `json:"docAicMix" excel:"单证类型-平均录入字符：门诊+住院"` //单证类型-平均录入字符：门诊+住院
}
