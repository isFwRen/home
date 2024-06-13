// 功能模块
const funcOptions = [
	{
		label: "全部",
		value: ""
	},
	{
		label: "案件列表",
		value: "案件列表"
	},
	{
		label: "任务管理",
		value: "任务管理"
	},
	{
		label: "时效管理",
		value: "时效管理"
	},
	{
		label: "人员管理",
		value: "人员管理"
	},
	{
		label: "质检管理",
		value: "质检管理"
	}
];

const fields = [
	{
		formKey: "time",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		defaultValue: undefined
	},

	{
		formKey: "proCode",
		inputType: "select ",
		hideDetails: true,
		label: "项目",
		options: [],
		hasDefaultValue: true
	},

	{
		formKey: "functionModule",
		inputType: "select",
		hideDetails: true,
		label: "功能模块",
		options: funcOptions,
		defaultValue: undefined
	},

	{
		formKey: "moduleOperation",
		inputType: "input",
		hideDetails: true,
		label: "模块操作",
		defaultValue: undefined
	},

	{
		formKey: "operationCodeOrName",
		inputType: "input",
		hideDetails: true,
		label: "工号/姓名",
		defaultValue: undefined
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
