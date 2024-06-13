import moment from "moment";

const headers = [
	{
		text: "日期",
		value: "sumDate",
		sortable: true,
		width: '100px',
		outPut: val => {
			return val.substr(0, 10);
		}
	},
	{ text: "案件量", value: "billCount", width: '60px'  },
	{
		text: "汇总",
		value: "summary",
		children: [
			{ text: "结算字符数", value: "charCount", width: '90px', },
			{ text: "录入字符数", value: "inputCharCount", width: '90px', },
			{ text: "员工字符", value: "staffInputCount", width: '80px', },
			{ text: "平均结算字符数/件", value: "averageCharCount", width: '130px', },
			{ text: "平均录入字符数/件", value: "averageInputCharCount", width: '130px' },
			{ text: "员工平均字符/件", value: "averageStaffInputCount", width: '130px' },
		]
	},

	{ text: "结算字符数与录入字符数的比例", value: "charPercent", width: '120px' },
	{ text: "结算字符数与员工字符数的比例", value: "settleStaffPercent", width: 120 },

	{
		text: "理赔类型-平均结算字符",

		value: "claimAverageSettlement",
		children: [
			{ text: "有医保", value: "claimsAccHasHis", width: '60px' },
			{ text: "无医保", value: "claimsAccNoneHis", width: '60px' },
			{ text: "混合型", value: "claimsAccMixHis", width: '60px' },
		]
	},

	{
		text: "理赔类型-平均录入字符",
		value: "claimAverageEntry",
		children: [
			{ text: "有医保", value: "claimsAicHasHis", width: '60px' },
			{ text: "无医保", value: "claimsAicNoneHis", width: '60px' },
			{ text: "混合型", value: "claimsAicMixHis", width: '60px' },
		]
	},

	{
		text: "单证类型-平均结算字符",
		value: "documentAverageSettlement",
		children: [
			{ text: "门诊", value: "cv1", width: '50px' },
			{ text: "住院", value: "cv2", width: '50px' },
			{ text: "门诊+住院", value: "cv3", width: '80px' },
		]
	},

	{
		text: "单证类型-平均录入字符",
		value: "documentAverageEntry",
		children: [
			{ text: "门诊", value: "dv1", width: '50px' },
			{ text: "住院", value: "dv2" , width: '50px'},
			{ text: "门诊+住院", value: "dv3", width: '80px'  },
		]
	},



];

const fields = [
	{
		cols: 2,
		formKey: "proCode",
		inputType: "select",
		hideDetails: true,
		label: "项目",
		options: [],
		defaultValue: undefined
	},

	{
		cols: 3,
		formKey: "date",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		range: true,
		defaultValue: [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")]
	}
];

const btns = [
	{
		class: "pr-3",
		color: "primary",
		text: "复制"
	},

	{
		class: "pr-3",
		color: "primary",
		text: "导出"
	}
];

export default {
	headers,
	fields,
	btns
};
