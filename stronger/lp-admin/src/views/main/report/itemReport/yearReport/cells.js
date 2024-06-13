import moment from "moment";

const headers1 = [{ text: "项目", value: "project", width: 100 }];

const headers2 = [{ text: "名称", value: "name" }];

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
		defaultValue: [
			moment().subtract(1, "month").format("YYYY-MM-DD"),
			moment().format("YYYY-MM-DD")
		]
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
	headers1,
	headers2,
	fields,
	btns
};
