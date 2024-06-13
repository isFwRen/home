import moment from "moment";

const headers = [
	{ text: "日期", value: "statisticalTime" },
	{ text: "工号", value: "code" },
	{ text: "姓名", value: "nickName" },
	{ text: "错误数量", value: "wrongNumber" },
	{ text: "申诉数量", value: "theNumberOfComplaints" },
	{ text: "申诉率", value: "theComplaintRate" },
	{ text: "通过数量", value: "throughTheNumber" },
	{ text: "通过率", value: "thePassRate" },
	{ text: "不通过数量", value: "nonPassingQuantity" },
	{ text: "不通过率", value: "unqualifiedRate" }
];

const fields = [
	{
		cols: 2,
		formKey: "proCode",
		inputType: "select",
		hideDetails: true,
		label: "项目"
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
	},

	{
		cols: 2,
		formKey: "code",
		inputType: "input",
		hideDetails: true,
		label: "工号"
	},

	{
		cols: 2,
		formKey: "nickName",
		inputType: "input",
		hideDetails: true,
		label: "姓名"
	}
];

const btns = [
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
