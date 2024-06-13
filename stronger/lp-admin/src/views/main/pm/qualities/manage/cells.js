import moment from "moment";

const headers = [
	{ text: "月份", value: "month", width: 65 },
	{ text: "项目编码", value: "proCode", width: 85 },
	{ text: "案件号", value: "billName", width: 120 },
	{ text: "反馈日期", value: "feedbackDate", width: 85 },
	{ text: "录入日期", value: "entryDate", width: 85 },
	{ text: "错误字段", value: "wrongFieldName", width: 85 },
	{ text: "正确值", value: "right", width: 65 },
	{ text: "错误值", value: "wrong", width: 65 },
	{ text: "初审责任人工号", value: "op0ResponsibleCode", width: 65 },
	{ text: "初审责任人姓名", value: "op0ResponsibleName", width: 65 },
	{ text: "一码责任人工号", value: "op1ResponsibleCode", width: 65 },
	{ text: "一码责任人姓名", value: "op1ResponsibleName", width: 65 },
	{ text: "二码责任人工号", value: "op2ResponsibleCode", width: 65 },
	{ text: "二码责任人姓名", value: "op2ResponsibleName", width: 65 },
	{ text: "问题件责任人工号", value: "opqResponsibleCode", width: 65 },
	{ text: "问题件责任人姓名", value: "opqResponsibleName", width: 65 },
	{ text: "影像", value: "imagePath", width: 150 },
	{ text: "操作", value: "options" }
];

const fields = [
	{
		formKey: "month",
		inputType: "date",
		label: "所属月份",
		pickerType: "month",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "所属月份不能为空." }]
	},

	{
		formKey: "proCode",
		inputType: "select",
		label: "项目编码",
		cols: 6,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "项目编码不能为空." }]
	},

	{
		formKey: "billName",
		inputType: "text",
		label: "案件号",
		cols: 6,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "案件号不能为空." }]
	},

	{
		formKey: "feedbackDate",
		inputType: "date",
		label: "反馈日期",
		cols: 6,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "反馈日期不能为空." }]
	},

	{
		formKey: "entryDate",
		inputType: "date",
		label: "录入日期",
		cols: 6,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "录入日期不能为空." }]
	},

	{
		formKey: "wrongFieldName",
		inputType: "text",
		label: "错误字段",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "错误字段不能为空." }]
	},

	{
		formKey: "right",
		inputType: "text",
		label: "正确值",
		cols: 6,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "正确值不能为空." }]
	},

	{
		formKey: "wrong",
		inputType: "text",
		label: "错误值",
		cols: 6,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "错误值不能为空." }]
	},

	{
		formKey: "op0ResponsibleCode",
		inputType: "text",
		label: "初审责任人工号",
		cols: 6
	},

	{
		formKey: "op0ResponsibleName",
		inputType: "text",
		label: "初审责任人姓名",
		disabled: true,
		cols: 6
	},

	{
		formKey: "op1ResponsibleCode",
		inputType: "text",
		label: "一码责任人工号",
		cols: 6
	},

	{
		formKey: "op1ResponsibleName",
		inputType: "text",
		label: "一码责任人姓名",
		disabled: true,
		cols: 6
	},

	{
		formKey: "op2ResponsibleCode",
		inputType: "text",
		label: "二码责任人工号",
		cols: 6
	},

	{
		formKey: "op2ResponsibleName",
		inputType: "text",
		label: "二码责任人姓名",
		disabled: true,
		cols: 6
	},

	{
		formKey: "opqResponsibleCode",
		inputType: "text",
		label: "问题件责任人工号",
		cols: 6
	},

	{
		formKey: "opqResponsibleName",
		inputType: "text",
		label: "问题件责任人姓名",
		disabled: true,
		cols: 6
	},

	{
		formKey: "file",
		inputType: "upload",
		label: "影像",
		action: "",
		autoUpload: false,
		limit: 1,
		prepend: "影像"
	}
];

const DEFAULT_MONTH = (() => {
	const date = new Date();

	const [year, month] = [date.getFullYear(), date.getMonth() + 1];

	return moment(`${year}-${month}`).format("YYYY-MM");
})();

export default {
	headers,
	fields,
	DEFAULT_MONTH
};
