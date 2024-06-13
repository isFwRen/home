const moreOptions = [
	{ label: "大图标", value: "large" },
	{ label: "中图标", value: "medium" },
	{ label: "小图标", value: "small" },
	{ label: "详细信息", value: "detail" }
];

const rightClickOptions = [
	{ label: "重命名", value: "rename" },
	{ label: "删除", value: "delete" },
	{ label: "属性", value: "attr" }
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
		formKey: "file",
		inputType: "fileInput",
		label: "规则文档",
		accept: "image/*",
		autoUpload: false,
		multiple: true,
		prependIcon: "mdi-file-excel-outline",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "规则文档不能为空." }]
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
	}
];

export default {
	moreOptions,
	rightClickOptions,
	fields
};
