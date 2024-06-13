import moment from "moment";

export const REJECT = {
	title: "不通过提示",
	content: "请确认是否不通过!"
};

export const ACCEPT = {
	title: "通过提示",
	content: "请确认是否通过!"
};

export const BATCH_REJECT = {
	title: "批量不通过提示",
	content: "请确认是否批量不通过!"
};

export const BATCH_ACCEPT = {
	title: "批量通过提示",
	content: "请确认是否批量通过!"
};

export const headers = [
	{ text: "日期", value: "submitDay" },
	{ text: "工号", value: "code" },
	{ text: "姓名", value: "nickName" },
	{ text: "案件号", value: "billName" },
	{ text: "机构号", value: "agency" },
	{ text: "字段", value: "fieldName" },
	{ text: "错误数据", value: "wrong" },
	{ text: "正确数据", value: "right" },
	{ text: "解析", value: "analysis" },
	{ text: "差错审核", value: "isWrongConfirm", width: 140 }
];

export const fields = [
	{
		cols: 2,
		formKey: "proCode",
		inputType: "select",
		hideDetails: true,
		label: "项目"
	},

	{
		cols: 3,
		formKey: "date",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		range: true,
		defaultValue: [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")]
	},

	{
		cols: 2,
		formKey: "code",
		inputType: "text",
		hideDetails: true,
		label: "工号"
	},

	{
		cols: 2,
		formKey: "nickName",
		inputType: "text",
		hideDetails: true,
		label: "姓名"
	},

	{
		cols: 2,
		formKey: "fieldName",
		inputType: "text",
		hideDetails: true,
		label: "字段名称"
	},

	{
		cols: 2,
		formKey: "op",
		inputType: "select",
		hideDetails: true,
		label: "工序",
		options: [
			{
				value: "",
				label: ""
			},
			{
				value: "Op0",
				label: "初审"
			},
			{
				value: "Op1",
				label: "一码"
			},
			{
				value: "Op2",
				label: "二码"
			},
			{
				value: "Opq",
				label: "问题件"
			}
		]
	},

	{
		cols: 2,
		formKey: "complaint",
		inputType: "select",
		clearable: true,
		hideDetails: true,
		label: "申诉",
		options: [
			{
				value: "true",
				label: "已申诉"
			},
			{
				value: "false",
				label: "待申诉"
			}
		]
	},

	{
		cols: 2,
		formKey: "confirm",
		inputType: "select",
		hideDetails: true,
		label: "审核",
		options: [
			{
				value: "",
				label: ""
			},
			{
				value: "true",
				label: "审核通过"
			},
			{
				value: "false",
				label: "审核不通过"
			}
		]
	}
];

export const reviewOptions = [
	{
		value: true,
		label: "通过"
	},

	{
		value: false,
		label: "不通过"
	}
];
