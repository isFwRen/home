import moment from "moment";

const headers = [
	{ text: "项目", value: "proCode" },
	{ text: "日期", value: "date" },
	{ text: "案件号", value: "billNum" },
	{ text: "批次号", value: "batchNum" },
	{ text: "机构号", value: "agency" },
	{ text: "省份", value: "province" },
	{ text: "城市", value: "city" },
	{ text: "账单号", value: "billCode" },
	{ text: "治疗医院", value: "hospital" },
	{ text: "费用类型", value: "costType" },
	{ text: "报销类型", value: "reimburseType" },
	{ text: "清单名称", value: "inventoryName" },
	{ text: "扣费类型", value: "chargingType" },
	{ text: "清单金额", value: "inventoryMoney" },
	{ text: "自付比例", value: "selfRatio" },
	{ text: "自付金额", value: "selfMoney" },
];


const fields = [
	{
		cols: 2,
		formKey: "proCode",
		inputType: "select",
		hideDetails: true,
		label: "项目",
		options: [
			{
				value: "1",
				label: "B0101"
			},
			{
				value: "2",
				label: "B0102"
			},
			{
				value: "3",
				label: "B0103"
			}
		],
		defaultValue: undefined
	},
	{
		cols: 2,
		formKey: "time",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		range: true,
		defaultValue: [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")]
	},
	{
		cols: 2,
		formKey: "bill_num",
		inputType: "input",
		hideDetails: true,
		label: "案件号",
	},
	{
		cols: 2,
		formKey: "list_name",
		inputType: "input",
		hideDetails: true,
		label: "清单名称",
	},
];


export default {
	fields,
	headers,
};
