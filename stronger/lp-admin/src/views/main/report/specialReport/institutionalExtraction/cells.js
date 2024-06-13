import moment from "moment";

const headers = [
	{ text: "机构", value: "agency" },
	{ text: "核心立案号", value: "billNum" },
	{ text: "疾病", value: "mIcd10Code" },
	{ text: "是否匹配", value: "isMatch" },
	{ text: "票据类型", value: "expenMode" },
	{ text: "医院名称", value: "hospital" },
	{ text: "总金额", value: "countMoney" },
	{ text: "统筹金额", value: "socialPayMoney" },
	{ text: "范围外金额", value: "outMoney" },
	{ text: "范围内金额", value: "innerMoney" },
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
		cols: 2,
		formKey: "time",
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
