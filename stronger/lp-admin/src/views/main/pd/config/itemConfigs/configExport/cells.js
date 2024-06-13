const headers = [
	{ text: "序号", value: "myOrder" },
	{ text: "代码", value: "name" },
	{ text: "字段1", value: "oneFields" },
	{ text: "字段2", value: "twoFields" },
	{ text: "字段3", value: "threeFields" },
	{ text: "来源", value: "myType" },
	{ text: "固定值", value: "fixedValue" },
	{ text: "备注", value: "remark" },
	{ text: "操作", value: "options", width: 160 }
];

const moreOptions = [
	{
		value: "delete",
		label: "删除"
	}
];

const xmlOptions = [
	{ label: "utf-8", value: "utf-8" },
	{ label: "GBK", value: "gbk" }
	// { label: 'gb2312', value: 'gb2312' }
];

const fields = [
	{
		formKey: "myOrder",
		inputType: "text",
		label: "序号",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "序号不能为空." }]
	},

	{
		formKey: "name",
		inputType: "text",
		label: "代码",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "代码不能为空." }]
	},

	{
		formKey: "oneFields",
		inputType: "autocomplete",
		clearable: true,
		label: "字段一",
		options: []
	},

	{
		formKey: "twoFields",
		inputType: "autocomplete",
		clearable: true,
		label: "字段二",
		options: []
	},

	{
		formKey: "threeFields",
		inputType: "autocomplete",
		clearable: true,
		label: "字段三",
		options: []
	},

	{
		formKey: "myType",
		inputType: "radios",
		label: "来源",
		prepend: "*",
		prependClass: "error--text",
		validation: [{ rule: "required", message: "来源不能为空." }],
		options: [
			{
				value: "1",
				label: "开始"
			},
			{
				value: "2",
				label: "结束"
			},
			{
				value: "3",
				label: "固定值"
			},
			{
				value: "4",
				label: "录入值"
			},
			{
				value: "5",
				label: "代码"
			}
		],
		prependOuter: "*"
	},

	{
		formKey: "fixedValue",
		inputType: "text",
		label: "固定值"
	},

	{
		formKey: "remark",
		inputType: "text",
		label: "备注"
	}
];

export default {
	headers,
	moreOptions,
	xmlOptions,
	fields
};
