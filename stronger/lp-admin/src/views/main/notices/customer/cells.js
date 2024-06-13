const headers = [
	{ text: "所属项目编码", value: "proCode" },
	{ text: "消息发送时间", value: "sendTime" },
	{ text: "消息文件名称", value: "fileName" },
	{ text: "消息类型", value: "msgType" },
	{ text: "消息状态", value: "status" },
	{ text: "消息回复时间", value: "replyTime" },
	{ text: "消息处理人", value: "dealUserName" },
	{ text: "操作", value: "options" }
];

const replyFields = [
	{
		formKey: "proCode",
		inputType: "text",
		label: "项目编码",
		disabled: true
	},

	{
		formKey: "sendTime",
		inputType: "text",
		label: "消息发送时间",
		disabled: true
	},

	{
		formKey: "fileName",
		inputType: "text",
		label: "消息文件名称",
		disabled: true
	},

	{
		formKey: "msgType",
		inputType: "text",
		label: "消息类型",
		disabled: true
	},

	{
		formKey: "content",
		inputType: "text",
		label: "消息内容",
		disabled: true
	},

	{
		formKey: "isReply",
		inputType: "radios",
		label: "回复"
	},

	{
		formKey: "expectNum",
		inputType: "text",
		label: "预计可增加单量",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "预计可增加单量不能为空." }]
	}
];

const viewFields = [
	...replyFields,

	{
		formKey: "replyTime",
		inputType: "text",
		label: "消息回复时间",
		disabled: true
	},

	{
		formKey: "dealUserName",
		inputType: "text",
		label: "消息处理人",
		disabled: true
	}
];

const isReplyItems = [
	{ label: "是", value: true },
	{ label: "否", value: false }
];

export default {
	headers,
	replyFields,
	viewFields,
	isReplyItems
};
