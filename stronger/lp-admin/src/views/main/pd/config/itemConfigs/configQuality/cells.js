const headers = [
	{ text: "XML节点", value: "parentXmlNodeName" },
	{ text: "XML代码", value: "xmlNodeName" },
	{ text: "字段名称", value: "fieldName" },
	{ text: "输入方式", value: "inputType" },
	{ text: "所属信息", value: "belongType" },
	{ text: "账单", value: "billInfo" },
	{ text: "受益人", value: "beneficiary" },
	{ text: "排版", value: "widthPercent", width: 60 },
	{ text: "排序", value: "myOrder", width: 60 },
	{ text: "操作", value: "options" }
];

const moreOptions = [
	{
		value: "delete",
		label: "删除"
	}
];

// XML节点
const parentXmlNodeNameOptions = [];

// XML代码
const xmlNodeNameOptions = [];

// 字段名称
const fieldNameOptions = [];

// 输入方式
const inputTypeOptions = [
	{
		value: 1,
		label: "输入框"
	},

	{
		value: 2,
		label: "常量表"
	},

	{
		value: 3,
		label: "下拉列表"
	}
];

// 所属信息
const belongTypeOptions = [
	{
		value: 1,
		label: "申请人信息"
	},

	{
		value: 2,
		label: "被保人信息"
	},

	{
		value: 3,
		label: "受托人信息"
	},

	{
		value: 4,
		label: "其它信息"
	},

	{
		value: 5,
		label: "受益人信息"
	},

	{
		value: 6,
		label: "领款人信息"
	},

	{
		value: 7,
		label: "账单信息"
	},

	{
		value: 8,
		label: "出险信息"
	}
];

// 账单
const billInfoOptions = [
	{
		value: 1,
		label: "账单号"
	},

	{
		value: 2,
		label: "账单类型"
	},

	{
		value: 3,
		label: "医院名称"
	},

	{
		value: 4,
		label: "开始时间"
	},

	{
		value: 5,
		label: "结束时间"
	},

	{
		value: 6,
		label: "账单金额"
	},

	{
		value: 7,
		label: "调整金额"
	},

	{
		value: 8,
		label: "自费金额"
	},

	{
		value: 9,
		label: "自付金额"
	},

	{
		value: 10,
		label: "报销金额"
	}
];

// 受益人
const beneficiaryOptions = [
	{
		value: 1,
		label: "受益人姓名"
	},
	{
		value: 2,
		label: "领款人姓名"
	},
	{
		value: 3,
		label: "其他"
	}
];

/**
 * 所属信息为选项为[受益人信息]：显示[受益人] radios
 * 所属信息为选项为[账单信息]：显示[账单] field
 */
const fields = [
	{
		formKey: "parentXmlNodeName",
		inputType: "autocomplete",
		label: "XML节点",
		optionName: "parentXmlNodeOptions",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "XML节点不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "xmlNodeName",
		inputType: "autocomplete",
		label: "XML代码",
		optionName: "xmlNodeOptions",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "XML代码不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "fieldCode",
		inputType: "autocomplete",
		label: "字段名称",
		optionName: "fieldOptions",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "字段名称不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "inputType",
		inputType: "select",
		label: "输入方式",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "输入方式不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "belongType",
		inputType: "select",
		label: "所属信息",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "所属信息不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "billInfo",
		inputType: "select",
		label: "账单",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		// suzerain: {
		//   attackFormKey: 'belongType',
		//   attackValues: [7],
		//   targetFormKey: 'billInfo'
		// },
		validation: [{ rule: "required", message: "账单不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "beneficiary",
		inputType: "radios",
		label: "受益人",
		prepend: "*",
		prependClass: "error--text",
		// suzerain: {
		//   attackFormKey: 'belongType',
		//   attackValues: [5],
		//   targetFormKey: 'beneficiary'
		// },
		options: [],
		validation: [{ rule: "required", message: "受益人不能为空." }],
		defaultValue: undefined
	},

	{
		formKey: "widthPercent",
		inputType: "text",
		label: "排版",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "排版不能为空." },
			{ rule: "numeric", message: "请输入正整数." }
		],
		defaultValue: undefined
	},

	{
		formKey: "myOrder",
		inputType: "text",
		label: "排序",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "排序不能为空." },
			{ rule: "numeric", message: "请输入正整数." }
		],
		defaultValue: undefined
	}
];

export default {
	headers,
	moreOptions,
	fields,
	parentXmlNodeNameOptions,
	xmlNodeNameOptions,
	fieldNameOptions
	// inputTypeOptions,
	// belongTypeOptions,
	// billInfoOptions,
	// beneficiaryOptions
};
