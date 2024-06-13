// 功能模块
const funcOptions = [
	{
		label: "PT群管理",
		value: "manage"
	},
	{
		label: "固定通知",
		value: "notice"
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
		formKey: "name",
		inputType: "input",
		hideDetails: true,
		label: "名称"
	},

	{
		formKey: "operationCodeOrName",
		inputType: "input",
		hideDetails: true,
		label: "工号/姓名"
	},
	{
		formKey: "modifyContent",
		inputType: "input",
		hideDetails: true,
		label: "修改内容"
	}
];

const headers = [
	{ text: "项目", value: "proCode", sortable: false },
	{ text: "功能模块", value: "functionModule", sortable: false },
	{ text: "名称", value: "name", sortable: false },
	{ text: "修改内容", value: "modificationContent", sortable: false },
	{ text: "修改人工号", value: "modifyLaborCode", sortable: false },
	{ text: "修改人姓名", value: "modifyLaborName", sortable: false },
	{ text: "修改时间", value: "modifyDate", sortable: false, format: "YYYY-MM-DD HH:mm:ss" }
];

export default {
	fields,
	headers
};
