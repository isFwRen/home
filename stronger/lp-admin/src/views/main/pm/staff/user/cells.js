const statusOptions = [
	{
		value: "",
		label: "全部"
	},

	{
		value: true,
		label: "在职"
	},

	{
		value: false,
		label: "离职"
	}
];

const filterFields = [
	{
		formKey: "code",
		inputType: "text",
		class: "mr-4",
		hideDetails: true,
		label: "工号",
		width: 120
	},

	{
		formKey: "name",
		inputType: "text",
		class: "mr-4",
		hideDetails: true,
		label: "姓名",
		width: 120
	},
	{
		formKey: "phone",
		inputType: "text",
		class: "mr-4",
		hideDetails: true,
		label: "手机号",
		width: 120
	},
	// {
	// 	formKey: "role",
	// 	inputType: "text",
	// 	class: "mr-4",
	// 	hideDetails: true,
	// 	label: "角色",
	// 	width: 120
	// },
	//
	// {
	// 	formKey: "status",
	// 	inputType: "select",
	// 	class: "mr-4",
	// 	hideDetails: true,
	// 	label: "状态",
	// 	options: statusOptions,
	// 	width: 120,
	// 	defaultValue: ""
	// },
	//
	// {
	// 	formKey: "date",
	// 	inputType: "date",
	// 	class: "mr-4",
	// 	hideDetails: true,
	// 	label: "上岗日期",
	// 	range: true,
	// 	width: 220
	// }
];

const exportDialogFields = [
	{
		formKey: "proCode",
		inputType: "select",
		label: "项目名称",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "项目名称不能为空." }]
	}
];

const dialogFields = [
	{
		formKey: "roleId",
		inputType: "select",
		label: "角色",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "角色不能为空." }]
	}
];

const headers = [
	{ text: "工号", value: "code", width: 60 },
	{ text: "姓名", value: "name", width: 120 },
	{ text: "手机号", value: "phone", width: 120 },
	// { text: "角色", value: "role", width: 120 },
	// { text: "状态", value: "status", width: 120 },
	// { text: "入职日期", value: "entryDate", date: true },
	// { text: "上岗日期", value: "mountGuardDate", date: true },
	// { text: "离职日期", value: "leaveDate", date: true },
	{ text: "操作", value: "options", width: 195 }
];

export default {
	filterFields,
	exportDialogFields,
	dialogFields,
	headers
};
