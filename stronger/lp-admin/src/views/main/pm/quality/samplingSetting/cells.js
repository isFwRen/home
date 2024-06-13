import moment from "moment";

const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

const typeOptions = [
	{
		value: 1,
		label: "人员抽检"
	},

	{
		value: 2,
		label: "分块抽检"
	}
];

const statusOptions = [
	{
		value: 1,
		label: "启用"
	},

	{
		value: 2,
		label: "停用"
	}
];

const btns = [
	{
		class: "pr-3",
		color: "success",
		icon: "mdi-plus",
		text: "新增",
		value: "new"
	},

	{
		class: "pr-3",
		color: "error",
		icon: "mdi-trash-can-outline",
		text: "删除",
		value: "delete"
	},

	{
		class: "pr-3",
		color: "error",
		icon: "mdi-minus-circle-outline",
		text: "停用",
		value: "stop"
	},

	{
		class: "pr-3",
		color: "success",
		icon: "mdi-check-circle-outline",
		text: "启用",
		value: "enable"
	},

	{
		class: "pr-3",
		color: "primary",
		icon: "mdi-content-copy",
		text: "复制",
		value: "copy"
	},

	{
		color: "primary",
		icon: "mdi-export-variant",
		text: "导出",
		value: "export"
	}
];

// 人员抽检
const headers1 = [
	{ text: "序号", value: "index" },
	{ text: "工号", value: "code" },
	{ text: "姓名", value: "name" },
	{ text: "抽检比例", value: "ratio" },
	{ text: "添加人工号", value: "createdCode" },
	{ text: "添加人姓名", value: "createdName" },
	{ text: "添加时间", value: "CreatedAt" },
	{ text: "操作", value: "options", width: 150 }
];

// 分块抽检
const headers2 = [
	{ text: "序号", value: "index" },
	{ text: "项目", value: "proCode" },
	{ text: "分块名称", value: "name" },
	{ text: "抽检比例", value: "ratio" },
	{ text: "添加人工号", value: "createdCode" },
	{ text: "添加人姓名", value: "createdName" },
	{ text: "添加时间", value: "CreatedAt" },
	{ text: "操作", value: "options", width: 150 }
];

const sameFields = [
	{
		formKey: "proCode",
		inputType: "select",
		label: "项目",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "项目不能为空." }]
	},

	{
		formKey: "ratio",
		inputType: "text",
		hideDetails: false,
		label: "抽检比例",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "抽检比例不能为空." }]
	}
];

const staffFields = [
	...sameFields,

	{
		formKey: "code",
		inputType: "autocomplete",
		label: "工号/姓名",
		multiple: true,
		optionName: "parentXmlNodeOptions",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "工号/姓名不能为空." }]
	}
];

const thunkFields = [
	...sameFields,

	{
		formKey: "name",
		inputType: "autocomplete",
		label: "分块名称",
		optionName: "parentXmlNodeOptions",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "分块名称不能为空." }]
	}
];

export default {
	DEFAULT_DATE,
	typeOptions,
	statusOptions,
	btns,
	headers1,
	headers2,
	staffFields,
	thunkFields
};
