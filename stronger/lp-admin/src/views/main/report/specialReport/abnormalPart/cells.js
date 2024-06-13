import moment from "moment";

const headers = [
	{ text: "日期", value: "createAt" },
	{ text: "案件号", value: "billName" },
	{ text: "机构", value: "agency" },
	{ text: "异常原因", value: "abnormalReason" },
	{ text: "回传时间", value: "uploadAt" }
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

	// 邓云峰说泽如让去掉(2022/03/25)
	// {
	//   cols: 2,
	//   formKey: 'types',
	//   inputType: 'select',
	//   hideDetails: true,
	//   label: '类型',
	//   options: [
	//     {
	//       value: '',
	//       label: '全部',
	//     },
	//     {
	//       value: '1',
	//       label: '医疗',
	//     },
	//     {
	//       value: '2',
	//       label: '非医疗',
	//     },
	//   ],
	//   defaultValue: undefined
	// }
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
