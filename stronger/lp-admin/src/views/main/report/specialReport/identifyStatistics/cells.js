import moment from "moment";

const headers = [
	{ text: "日期", value: "CreatedAt" },
	{ text: "案件号", value: "bill_num" },
	{ text: "字段", value: "name" },
	{ text: "是否OCR", value: "ocr_type" },
	{ text: "姓名", value: "nickName" },
	{ text: "工号", value: "jobNumber" },
	{ text: "ocr内容", value: "value" },
	{ text: "最终回传数据", value: "result_value", type: 'html' },
	{ text: "字段屏蔽", value: "disable" },
	{ text: "数据对比", value: "compare" },
	{ text: "准确率", value: "rate", },
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
				label: "全部"
			},
			{
				value: "2",
				label: "B0101"
			},
			{
				value: "3",
				label: "B0102"
			},
			{
				value: "4",
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
	// {
	// 	cols: 2,
	// 	formKey: "billCode",
	// 	inputType: "input",
	// 	hideDetails: true,
	// 	label: "案件号"
	// },
	// {
	// 	cols: 2,
	// 	formKey: "fieldName",
	// 	inputType: "input",
	// 	hideDetails: true,
	// 	label: "字段名称"
	// },

	// {
	// 	cols: 2,
	// 	formKey: "fieldShield",
	// 	inputType: "select",
	// 	hideDetails: true,
	// 	label: "字段屏蔽",
	// 	options: [
	// 		{
	// 			value: "1",
	// 			label: "请选择"
	// 		},
	// 		{
	// 			value: "2",
	// 			label: "医疗机构"
	// 		}
	// 	],
	// 	defaultValue: undefined
	// },

	// {
	// 	cols: 2,
	// 	formKey: "dataComparison",
	// 	inputType: "select",
	// 	hideDetails: true,
	// 	label: "数据对比",
	// 	options: [
	// 		{
	// 			value: "1",
	// 			label: "请选择"
	// 		},
	// 		{
	// 			value: "2",
	// 			label: "医疗机构"
	// 		}
	// 	],
	// 	defaultValue: undefined
	// }
];


export default {
	fields,
	headers,
};
