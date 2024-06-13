import moment from "moment";

const headers = [
	{ text: "日期", value: "CreatedAt" },
	{ text: "案件号", value: "billNum" },
	{ text: "机构", value: "agency" },
	{ text: "目录外医院名称", value: "name" }
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
	},

	{
		cols: 2,
		formKey: "type",
		inputType: "select",
		hideDetails: true,
		label: "类型",
		options: [
			{
				value: "1",
				label: "请选择"
			},
			{
				value: "2",
				label: "医疗机构"
			}
		],
		defaultValue: '1',
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
