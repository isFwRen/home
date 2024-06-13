const headers = [
	{ text: "项目编码", value: "proCode" },
	{ text: "项目名称", value: "proName" },
	{ text: "权限设置", value: "permission", width: 500 },
	{ text: "内外网", value: "network", width: 200 }
];

// 权限设置 headers checkbox
const permissionHeadersCheckboxs = [
	{
		formKey: "op0",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 0,
		selected: false,
		value: "hasOp0"
	},

	{
		formKey: "op1",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 1,
		selected: false,
		value: "hasOp1"
	},

	{
		formKey: "op2",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 2,
		selected: false,
		value: "hasOp2"
	},

	{
		formKey: "opq",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 3,
		selected: false,
		value: "hasOpq"
	},

	{
		formKey: "pm",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 4,
		selected: false,
		value: "hasPm"
	}
];

// 权限设置 desserts checkbox
const permissionDessertsCheckboxs = [
	{
		initFormKey: "op",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 0,
		label: "初审",
		value: "hasOp0"
	},

	{
		initFormKey: "op",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 1,
		label: "一码",
		value: "hasOp1"
	},

	{
		initFormKey: "op",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 2,
		label: "二码",
		value: "hasOp2"
	},

	{
		initFormKey: "op",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 3,
		label: "问题件",
		value: "hasOpq"
	},

	{
		initFormKey: "pm",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 4,
		label: "PM",
		value: "hasPm"
	}
];

// 内外网 headers checkbox
const netHeadersCheckboxs = [
	{
		formKey: "inNet",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 0,
		selected: false,
		value: "hasInNet"
	},

	{
		formKey: "outNet",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		indeterminate: false,
		index: 1,
		selected: false,
		value: "hasOutNet"
	}
];

// 内外网 headers checkbox
const netDessertsCheckboxs = [
	{
		initFormKey: "net",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 0,
		label: "内网",
		value: "hasInNet"
	},

	{
		initFormKey: "net",
		class: "pa-0 ma-0",
		dense: true,
		hideDetails: true,
		index: 1,
		label: "外网",
		value: "hasOutNet"
	}
];

export { permissionHeadersCheckboxs, netHeadersCheckboxs };

export default {
	headers,
	permissionHeadersCheckboxs,
	permissionDessertsCheckboxs,
	netHeadersCheckboxs,
	netDessertsCheckboxs
};
