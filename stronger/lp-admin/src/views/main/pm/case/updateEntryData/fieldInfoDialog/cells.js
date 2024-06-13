const opItems = [
	{ label: "初审", value: "op0" },
	{ label: "一码", value: "op1" },
	{ label: "二码", value: "op2" },
	{ label: "问题件", value: "opq" }
];

const fields = [
	{
		formKey: "fieldValue",
		inputType: "text",
		label: "正确数据",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "正确数据不能为空." }]
	},

	{
		formKey: "editDate",
		inputType: "date",
		label: "反馈日期",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "请选择反馈日期." }]
	},

	{
		formKey: "month",
		inputType: "date",
		label: "月份",
		pickerType: "month",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "请选择月份." }]
	},

	{
		formKey: "responsibleCode",
		inputType: "autocomplete",
		label: "责任人工号",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "请选择责任人工号." }]
	},

	{
		formKey: "responsibleName",
		inputType: "text",
		label: "责任人姓名",
		disabled: true,
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "责任人姓名不能为空." }]
	},

	{
		formKey: "op",
		inputType: "select",
		label: "对应工序",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "请选择对应工序." }],
		options: opItems
	}
];

export default {
	fields
};
