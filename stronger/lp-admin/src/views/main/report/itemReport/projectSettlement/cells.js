import moment from "moment";

const fields = [
	{
		cols: 2,
		formKey: "num1",
		inputType: "select",
		hideDetails: true,
		label: "项目",
		options: [
			{
				value: "1",
				label: "B0114"
			},
			{
				value: "2",
				label: "B0102"
			},
			{
				value: "3",
				label: "B0103"
			}
		],
		defaultValue: undefined
	},

	{
		cols: 3,
		formKey: "num2",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		range: true,
		defaultValue: [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")]
	}
];

const btns = [
	{
		class: "pr-3",
		color: "primary",
		text: "生成并导出"
	}
];

export default {
	fields,
	btns
};
