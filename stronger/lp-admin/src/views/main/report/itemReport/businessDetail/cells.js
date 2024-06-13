// import moment from 'moment'

const billType = {
	1: '门诊',
	2: '住院',
	3: '门诊,住院'
}
const headers = [
	{ text: "日期", value: "createAt" },
	{ text: "批次号", value: "batchNum", width: "150px" },
	{ text: "销售渠道", value: "saleChannel" },
	{ text: "案件号", value: "billNum" },
	{ text: "影像数量", value: "pictureNumber" },
	{ text: "机构", value: "agency" },
	{ text: "理赔类型", value: "claimType" },
	// { text: "医保类型", value: "insuranceType" },
	{ text: "单据类型", value: "billType" },
	{ text: "扫描时间", value: "scanAt" },
	{ text: "下载时间", value: "downloadAt" },
	{ text: "初次导出时间", value: "firstExportAt",width: "100px" },
	{ text: "导出时间", value: "exportAt" },
	{ text: "初次回传时间", value: "uploadAt", width: "100px" },
	{ text: "最新回传时间", value: "latestUploadAt", width: "100px" },
	{ text: "最晚回传时间", value: "atTheLateUploadAt", width: "100px" },
	{ text: "处理时长", value: "workTime" },
	{ text: "延迟时间", value: "lateTime" },
	// { text: "备注", value: "remark", width: 100 },
	{ text: "案件状态", value: "status" },
	{ text: "录入状态", value: "stage" },
	{ text: "是否超时", value: "isTheTimeOut" },
	{ text: "时长范围", value: "durationRange" },
	{ text: "质检人工号", value: "qualityUserCode" },
	{ text: "质检姓名", value: "qualityUserName" },
	{ text: "质检状态", value: "exportStage" },
	{ text: "案件次数", value: "theNumberOfCase" },
	{ text: "账单金额", value: "countMoney" },
	{ text: "发票数量", value: "invoiceNum" },
	{ text: "清单数量", value: "listingNum" },
	{ text: "问题件数量", value: "questionNum" },
	{ text: "疾病诊断", value: "diseaseDiagnosis" },
	{ text: "录入字符数", value: "fieldCharacter" },
	{ text: "申请书录入字符", value: "field1" },
	{ text: "发票录入字符", value: "field2" },
	{ text: "清单录入字符", value: "field3" },
	{ text: "结算字符数", value: "settlementCharacter" },
	{ text: "申请书结算字符", value: "settlemen1" },
	{ text: "发票结算字符", value: "settlemen2" },
	{ text: "清单结算录入字符", value: "settlemen3" },
	// { text: "结算金额", value: "settlementMoney" },

	{ text: "时效考核要求", value: "requirementOfAging", width: "100px" },
	{ text: "初审进入时间", value: "op0Entry", width: "100px" },
	{ text: "初审结束时间", value: "op0End", width: "100px" },
	{ text: "初审处理时长", value: "op0WorkTime", width: "100px" },
	{ text: "一码进入时间", value: "op1Entry", width: "100px" },
	{ text: "一码结束时间", value: "op1End", width: "100px" },
	{ text: "一码处理时间", value: "op1WorkTime", width: "100px" },
	{ text: "二码进入时间", value: "op2Entry", width: "100px" },
	{ text: "二码结束时间", value: "op2End", width: "100px" },
	{ text: "二码处理时间", value: "op2WorkTime", width: "100px" },
	{ text: "一二码处理时长", value: "op1AndOp2Duration", width: "120px" },
	{ text: "问题件进入时间", value: "opQEntry", width: "120px" },
	{ text: "问题件结束时间", value: "opQEnd", width: "120px" },
	{ text: "问题件处理时间", value: "opQWorkTime", width: "120" }
	// { text: '质检进入时间', value: 'testEntryTime', width: '100px' },
	// { text: '质检结束时间', value: 'testEndTime', width: '100px' },
	// { text: '质检处理时长', value: 'testHandleTime', width: '100px' },
];

// const fields = [
//   {
//     cols: 2,
//     formKey: 'proCode',
//     inputType: 'select',
//     hideDetails: true,
//     label: '项目',
//     options: []
//   },

//   {
//     cols: 3,
//     formKey: 'date',
//     inputType: 'date',
//     hideDetails: true,
//     label: '日期',
//     clearable: true,
//     range: true,
//     defaultValue: [moment().format('YYYY-MM-DD'), moment().format('YYYY-MM-DD')]
//   },
// ]

export default {
	headers,
	billType
	// fields
};
