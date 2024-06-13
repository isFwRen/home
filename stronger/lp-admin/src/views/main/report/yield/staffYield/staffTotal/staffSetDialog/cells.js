import moment from "moment";

const headers = [
	{ text: "项目编码", value: "proCode", sortable: false },
	{ text: "初审(按分块)", value: "op0AsTheBlock", sortable: false },
	{ text: "初审(按发票)", value: "op0AsTheInvoice", sortable: false },
	{ text: "一码非报销单", value: "op1NotExpenseAccount", sortable: false },
	{ text: "一码报销单", value: "op1ExpenseAccount", sortable: false },
	{ text: "二码非报销单", value: "op2NotExpenseAccount", sortable: false },
	{ text: "二码报销单", value: "op2ExpenseAccount", sortable: false },
	{ text: "问题件", value: "question", sortable: false },
	{ text: "开始时间", value: "startTime", sortable: false },
	{ text: "更新日期", value: "UpdatedAt", sortable: false },
	{ text: "操作", value: "options", sortable: false }
];

const statusOptions = [
	{
		label: "请选择",
		value: ""
	}
];

const fields = [
	{
		formKey: "proCode",
		inputType: "select",
		label: "项目编码",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "项目编码不能为空." }],
		options: statusOptions,
		defaultValue: undefined
	},

	{
		formKey: "op0AsTheBlock",
		inputType: "text",
		label: "初审(按分块)",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "初审(按分块)不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "op0AsTheInvoice",
		inputType: "text",
		label: "初审(按发票)",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "初审(按发票)不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "op1NotExpenseAccount",
		inputType: "text",
		label: "一码非报销单",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "一码非报销单不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "op1ExpenseAccount",
		inputType: "text",
		label: "一码报销单",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "一码报销单不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "op2NotExpenseAccount",
		inputType: "text",
		label: "二码非报销单",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "二码非报销单不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "op2ExpenseAccount",
		inputType: "text",
		label: "二码报销单",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "二码报销单不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "question",
		inputType: "text",
		label: "问题件",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "问题件不能为空." },
			{ regex: /^\d+(\.\d{0,2})?$/, message: "只能输入2位小数" }
		],
		defaultValue: undefined
	},

	{
		formKey: "startTime",
		inputType: "date",
		hideDetails: true,
		label: "开始时间",
		clearable: true,
		range: false,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "开始时间不能为空." }],
		defaultValue: moment().format("YYYY-MM-DD")
	}
];

export default {
	fields,
	headers
};
