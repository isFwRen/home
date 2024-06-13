const headers = [
	{ text: "XML代码", value: "xmlNodeCode", width: 125 },
	{ text: "名称", value: "xmlNodeName", width: 125 },
	{ text: "自定义描述", value: "msg", width: 200 },
	{ text: "只能录入", value: "onlyInput", width: 125 },
	{ text: "不能录入", value: "notInput", width: 125 },
	{ text: "最大长度", value: "maxLen", width: 45 },
	{ text: "最小长度", value: "minLen", width: 45 },
	{ text: "最大值", value: "maxVal", width: 45 },
	{ text: "最小值", value: "minVal", width: 45 },
	{ text: "校验", value: "validations" },
	{ text: "操作", value: "options" }
];

const moreOptions = [
	{
		value: "delete",
		label: "删除"
	}
];

const fields = [
	{
		formKey: "xmlNodeCode",
		inputType: "text",
		hideDetails: false,
		label: "XML代码",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "XML代码不能为空." }]
	},
	{
		formKey: "xmlNodeName",
		inputType: "text",
		hideDetails: false,
		label: "xml名称",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "名称不能为空." }]
	},
	{
		formKey: "proName",
		inputType: "text",
		label: "项目名称",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "项目名称名称不能为空." }]
	},

	{
		formKey: "msg",
		inputType: "text",
		label: "自定义描述"
	},

	{
		formKey: "onlyInput",
		inputType: "text",
		label: "只能录入"
	},

	{
		formKey: "notInput",
		inputType: "text",
		label: "不能录入"
	},

	{
		formKey: "maxLen",
		inputType: "text",
		label: "最大长度"
	},

	{
		formKey: "minLen",
		inputType: "text",
		label: "最小长度"
	},

	{
		formKey: "maxVal",
		inputType: "text",
		label: "最大值"
	},

	{
		formKey: "minVal",
		inputType: "text",
		label: "最小值"
	},

	{
		formKey: "validation",
		inputType: "checkboxs",
		label: "校验：",
		options: []
	}
];

export default {
	headers,
	moreOptions,
	fields
};
