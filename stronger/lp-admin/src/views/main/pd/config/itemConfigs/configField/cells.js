const headers = [
	// { text: '序号', value: 'myOrder' },
	{ text: "字段编码", value: "code" },
	{ text: "字段名称", value: "name" },
	{ text: "录入工序", value: "inputProcess" },
	{ text: "问题件配置", value: "opqConfig" },
	{ text: "导出校验配置", value: "exportConfig" },
	{ text: "字段拦截提示配置", value: "interceptionConfig" },
	{ text: "操作", value: "options" }
];

const moreOptions = [
	{
		value: "copy",
		label: "复制"
	},
	{
		value: "delete",
		label: "删除"
	}
];

const fields = [
	{
		cols: 4,
		formKey: "code",
		inputType: "text",
		disabled: true,
		hideDetails: false,
		label: "字段编码",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "字段编码不能为空." }]
	},

	{
		cols: 4,
		formKey: "name",
		inputType: "text",
		hideDetails: false,
		label: "字段名称",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "字段名称不能为空." }]
	},

	{
		cols: 4,
		formKey: "fixValue",
		inputType: "text",
		hideDetails: false,
		label: "固定值"
	},

	{
		cols: 4,
		formKey: "specChar",
		inputType: "text",
		hideDetails: false,
		label: "可通过字符"
	},

	{
		cols: 4,
		formKey: "inputProcess",
		inputType: "select",
		clearable: false,
		hideDetails: false,
		label: "录入工序",
		options: [
			{
				value: 1,
				label: "不录"
			},
			{
				value: 2,
				label: "一码"
			},

			{
				value: 3,
				label: "二码"
			}
		],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "录入工序不能为空." }]
	},

	{
		cols: 4,
		formKey: "checkDate",
		inputType: "select",
		clearable: true,
		hideDetails: false,
		label: "日期时限",
		options: []
	},

	{
		cols: 6,
		formKey: "questionChange",
		inputType: "text",
		hideDetails: false,
		label: "问题件转换"
	},

	{
		cols: 6,
		formKey: "valChange",
		inputType: "text",
		hideDetails: false,
		label: "数值转换"
	},

	{
		cols: 6,
		formKey: "valInsert",
		inputType: "text",
		hideDetails: false,
		label: "数据插入"
	},

	{
		cols: 6,
		formKey: "ignoreIf",
		inputType: "text",
		hideDetails: false,
		label: "不录条件"
	},

	{
		cols: 12,
		formKey: "prompt",
		inputType: "text",
		hideDetails: false,
		label: "录入提示",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "录入提示不能为空." }]
	},

	{
		cols: 4,
		formKey: "maxLen",
		inputType: "text",
		hideDetails: false,
		label: "最大长度"
	},

	{
		cols: 4,
		formKey: "minLen",
		inputType: "text",
		hideDetails: false,
		label: "最小长度"
	},

	{
		cols: 4,
		formKey: "fixLen",
		inputType: "text",
		hideDetails: false,
		label: "固定长度"
	},

	{
		cols: 4,
		formKey: "maxVal",
		inputType: "text",
		hideDetails: false,
		label: "最大值"
	},

	{
		cols: 4,
		formKey: "minVal",
		inputType: "text",
		hideDetails: false,
		label: "最小值"
	},

	{
		cols: 4,
		formKey: "defaultVal",
		inputType: "text",
		hideDetails: false,
		label: "默认值"
	},

	{
		cols: 12,
		formKey: "validations",
		inputType: "checkboxs",
		hideDetails: false,
		options: []
	}
];


const exportObject = [
	{
		items: ['等于', '包含', '不包含'],
		checkType: '等于',
		value: '',
		mark: '',
		icon: 'mdi-delete-empty'
	}
]
const interceptionObject = [
	{
		items: ['等于', '包含', '不包含'],
		item: '等于',
		itemContent: '',
		itemDesc: '',
		isInterception: false,
		icon: 'mdi-delete-empty',
		buttongroups: [
			{
				value: 1,
				label: "是"
			},
			{
				value: 2,
				label: "否"
			},
		],
	}
]

export default {
	headers,
	moreOptions,
	fields,
	exportObject,
	interceptionObject
};
