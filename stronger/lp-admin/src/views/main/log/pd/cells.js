// 功能模块
const funcOptions = [
	{
		label: "全部",
		value: ""
	},
	{
		label: "配置管理",
		value: "配置管理"
	},
	{
		label: "常量管理",
		value: "常量管理"
	}
];

const fields = [
	{
		formKey: "time",
		inputType: "date",
		hideDetails: true,
		label: "日期"
	},

	{
		formKey: "proCode",
		inputType: "select ",
		hideDetails: true,
		label: "项目",
		options: []
	},

	{
		formKey: "functionModule",
		inputType: "select",
		hideDetails: true,
		label: "功能模块",
		options: funcOptions
	},

	{
		formKey: "moduleOperation",
		inputType: "input",
		hideDetails: true,
		label: "模块操作"
	},

	{
		formKey: "operationCodeOrName",
		inputType: "input",
		hideDetails: true,
		label: "工号/姓名"
	}
];

const headers = [
	{ text: "项目", value: "proCode", sortable: false },
	{ text: "功能模块", value: "functionModule", sortable: false },
	{ text: "模块操作", value: "moduleOperation", sortable: false },
	{ text: "操作人工号", value: "operationCode", sortable: false },
	{ text: "操作人姓名", value: "operationName", sortable: false },
	{ text: "操作时间", value: "CreatedAt", sortable: false }
];

export default {
	fields,
	headers
};
