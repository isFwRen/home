const headers = [
	{ text: "项目编码", value: "proCode" },
	{ text: "规则名称", value: "ruleName" },
	{ text: "规则类型", value: "ruleType" },
	{ text: "更新日期", value: "UpdatedAt" },
	{ text: "更新人", value: "updatedName" },
	{ text: "操作", value: "options", width: 260 }
];

const ruleTypes = [
	{ label: "项目规则", value: "项目规则" },
	{ label: "易错规则", value: "易错规则" }
];

const fields = [
	{
		formKey: "proCode",
		inputType: "select",
		label: "项目编码",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "项目编码不能为空." }]
	},

	{
		formKey: "ruleName",
		inputType: "text",
		label: "规则名称",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "规则名称不能为空." }]
	},

	{
		formKey: "ruleType",
		inputType: "select",
		label: "规则类型",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "规则类型不能为空." }]
	},
	{
		formKey: "isRequired",
		inputType: "select",
		label: "是否必学",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "是否必学不能为空." }],
		options: [
			{
				label: "是",
				value: 1
			},
			{
				label: "否",
				value: 0
			}
		]
	},
	{
		formKey: "file",
		inputType: "fileInput",
		clearable: true,
		label: "规则文档",
		// accept: 'application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document',
		accept: "application/pdf",
		autoUpload: false,
		maxSize: 1024 * 5,
		prependIcon: "mdi-file-excel-outline",
		prependOuter: "*",
		prependOuterClass: "error--text",
		showUploadList: true,
		validation: [{ rule: "required", message: "规则文档不能为空." }]
	}
];

export default {
	headers,
	ruleTypes,
	fields
};
