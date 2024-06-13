// 顶部导航
const entries = [
	{
		index: 0,
		label: "初审",
		value: "op0",
		authKey: "hasOp0",
		class: "mr-3",
		link: "/main/entry/channel/op0"
	},
	{
		index: 1,
		label: "一码",
		value: "op1",
		authKey: "hasOp1",
		class: "mr-3",
		link: "/main/entry/channel/op1"
	},
	{
		index: 2,
		label: "二码",
		value: "op2",
		authKey: "hasOp2",
		class: "mr-3",
		link: "/main/entry/channel/op2"
	},
	{
		index: 3,
		label: "问题件",
		value: "opq",
		authKey: "hasOpq",
		link: "/main/entry/channel/opq"
	}
];

// 删单
const deleteFields = [
	{
		formKey: "password",
		inputType: "text",
		label: "删单密码",
		prependOuter: "*",
		prependOuterClass: "error--text",
		type: "password",
		validation: [{ rule: "required", message: "删单密码不能为空." }]
	},

	{
		formKey: "delRemarks",
		inputType: "textarea",
		hideDetails: false,
		label: "备注",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "备注不能为空." }]
	}
];

export { entries, deleteFields };

export default {
	entries,
	deleteFields
};
