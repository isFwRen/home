const headers = [
	{ text: "", value: "name" },
	{ text: "初审", value: "op0" },
	{ text: "一码报销单", value: "op1" },
	{ text: "一码非报销单", value: "op1No" },
	{ text: "二码报销单", value: "op2" },
	{ text: "二码非报销单", value: "op2No" },
	{ text: "问题件", value: "opq" }
];

const textareas = [
	{
		formId: "StickLevel1",
		formKey: "stickLevel1",
		val: "urgent",
		title: "紧急件",
		label: "紧急件",
		placeholder: "在此输入紧急处理的案件号，每个案件号为一行"
	},

	{
		formId: "StickLevel2",
		formKey: "stickLevel2",
		val: "priority",
		title: "优先件",
		label: "优先件",
		placeholder: "在此输入优先处理的案件号，每个案件号为一行"
	},

	{
		formId: "OrganizationNumber",
		formKey: "organizationNumber",
		val: "agency",
		title: "机构号",
		label: "机构号",
		placeholder: "在此输入优先处理的机构号，机构号之间用(,)隔开"
	}
];

export default {
	headers,
	textareas
};
