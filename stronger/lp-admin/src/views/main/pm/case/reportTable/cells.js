const fields = [
	{
		formKey: "invoiceNum",
		inputType: "input",
		hideDetails: true,
		label: "账单号",
		defaultValue: undefined
	},

	{
		formKey: "itemName",
		inputType: "input",
		hideDetails: true,
		label: "项目名称",
		defaultValue: undefined
	}
];

const headers = [
	{ text: "医疗项目分类", value: "medicalType" },
	{ text: "项目名称", value: "name" },
	{ text: "项目类型", value: "type", width: 80 },
	{ text: "项目金额", value: "price", width: 80 },
	{ text: "数量", value: "count", width: 80 },
	{ text: "项目比例", value: "percent", width: 80 },
	{ text: "项目自付（自费、自付）金额", value: "pay" }
];

export default {
	fields,
	headers
};
