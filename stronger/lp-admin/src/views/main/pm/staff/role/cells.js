const statusOptions = [
	// {
	//   label: '全部',
	//   value: ''
	// },

	{
		label: "正常",
		value: 1
	},

	{
		label: "停用",
		value: 2
	}
];

const fields = [
	{
		formKey: "name",
		inputType: "text",
		label: "角色名称",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "角色名称不能为空." }]
	},

	{
		formKey: "status",
		inputType: "select",
		label: "状态",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "状态不能为空." }],
		options: statusOptions
	},

	{
		formKey: "remark",
		inputType: "text",
		label: "描述"
	}
];

const headers = [
	{ text: "创建日期", value: "CreatedAt" },
	{ text: "角色名称", value: "name" },
	{ text: "状态", value: "status" },
	// { text: '用户数', value: 'countUser' },
	{ text: "描述", value: "remark" },
	{ text: "操作", value: "options" }
];

export default {
	statusOptions,
	fields,
	headers
};
