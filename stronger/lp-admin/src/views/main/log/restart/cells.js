// 功能模块
const funcOptions = [
	{
		label: "全部",
		value: ""
	},
	{
		label: "产量统计",
		value: "产量统计"
	},
	{
		label: "错误查询",
		value: "错误查询"
	},
	{
		label: "工资数据",
		value: "工资数据"
	},
	{
		label: "项目报表",
		value: "项目报表"
	},
	{
		label: "特殊报表",
		value: "特殊报表"
	},
	{
		label: "三大目录",
		value: "三大目录"
	},
	{
		label: "公告管理",
		value: "公告管理"
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
