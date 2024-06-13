const fields = [
	{
		cols: 12,
		formKey: "name",
		inputType: "input",
		label: "项目名称",
		disabled: true,
		defaultValue: undefined
	},

	{
		cols: 12,
		formKey: "code",
		inputType: "input",
		label: "项目编码",
		disabled: true,
		defaultValue: undefined
	},

	{
		cols: 12,
		formKey: "type",
		inputType: "input",
		label: "所属行业",
		defaultValue: undefined
	},

	{
		cols: 12,
		formKey: "cacheTime",
		inputType: "input",
		label: "缓存时间",
		suffix: "秒",
		defaultValue: undefined
	},

	{
		cols: 12,
		formKey: "saveDate",
		inputType: "input",
		label: "数据保存",
		suffix: "天",
		defaultValue: undefined
	},

	{
		cols: 12,
		formKey: "restartAt",
		hideDetails: false,
		hint: "格式：hh:mm:ss，留空表示不自动重启",
		inputType: "input",
		label: "自动重启",
		validation: [
			{
				// regex: /^([0-1][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$/,
				regx: /^(\*|((\*\/)?[1-5]?[0-9])) (\*|((\*\/)?[1-5]?[0-9])) (\*|((\*\/)?(1?[0-9]|2[0-3]))) (\*|((\*\/)?([1-9]|[12][0-9]|3[0-1]))) (\*|((\*\/)?([1-9]|1[0-2]))) (\*|((\*\/)?[0-6]))$/,
				message: "自动重启格式为：hh:mm:ss."
			}
		],
		defaultValue: undefined
	},

	{
		cols: 12,
		formKey: "autoReturn",
		inputType: "radios",
		options: [
			{ value: false, label: "手动回传" },
			{ value: true, label: "自动回传" }
		],
		defaultValue: undefined
	}
];

export default {
	fields
};
