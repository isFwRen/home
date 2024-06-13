const headers = [
	{ text: "时效起始时间", value: "agingStartTime" },
	{ text: "时效结束时间", value: "agingEndTime" },
	{ text: "时效外开始时间", value: "agingOutStartTime" },
	{ text: "时效外最晚时间", value: "agingOutEndTime" },
	{ text: "考核要求(min)", value: "requirementsTime" },
	{ text: "操作", value: "options", width: 160 }
];

const moreOptions = [
	{
		value: "delete",
		label: "删除"
	}
];

const fields = [
	{
		formKey: "agingStartTime",
		inputType: "date",
		dateFormat: "HH:mm:ss",
		hideDetails: false,
		hint: "格式为hh:mm:ss",
		label: "时效起始时间",
		mode: "time",
		prependOuter: "*",
		prependOuterClass: "error--text",
		timeFormat: "24hr",
		timeUseSeconds: true,
		validation: [{ rule: "required", message: "时效起始时间不能为空." }]
	},

	{
		formKey: "agingEndTime",
		inputType: "date",
		dateFormat: "HH:mm:ss",
		hideDetails: false,
		hint: "格式为hh:mm:ss",
		label: "时效结束时间",
		mode: "time",
		prependOuter: "*",
		prependOuterClass: "error--text",
		timeFormat: "24hr",
		timeUseSeconds: true,
		validation: [{ rule: "required", message: "时效结束时间不能为空." }]
	},

	{
		formKey: "agingOutStartTime",
		inputType: "date",
		dateFormat: "HH:mm:ss",
		hideDetails: false,
		hint: "格式为hh:mm:ss",
		label: "时效外开始时间",
		mode: "time",
		timeFormat: "24hr",
		timeUseSeconds: true
	},

	{
		formKey: "agingOutEndTime",
		inputType: "date",
		dateFormat: "HH:mm:ss",
		hideDetails: false,
		hint: "格式为hh:mm:ss",
		label: "时效外最晚时间",
		mode: "time",
		timeFormat: "24hr",
		timeUseSeconds: true
	},

	{
		formKey: "requirementsTime",
		inputType: "text",
		hideDetails: false,
		hint: "请输入数字",
		label: "考核要求(min)",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [
			{ rule: "required", message: "考核要求不能为空." },
			{ rule: "numeric", message: "考核要求为数字." }
		]
	}
];

export default {
	headers,
	moreOptions,
	fields
};
