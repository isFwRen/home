import moment from "moment";

const headers = [
	{ text: "名称", value: "name" },

	{
		text: "产值",
		value: "ouputValue",
		children: [
			{ text: "预估产值", value: "predictValue" },
			{ text: "实际完成", value: "finishValue" },
			{
				text: "时间比例",
				value: "timePercent",
				output: val => {
					return val * 100 + "%";
				}
			},
			{ text: "完成比例", value: "finishPercent" }
		]
	},

	{
		text: "业务量",
		value: "businessVolume",
		children: [
			{ text: "月业务量", value: "monthCount" },
			{ text: "日业务量", value: "dayCount" }
		]
	},

	{
		text: "时效",
		value: "prescription",
		children: [
			{ text: "月时效保障率", value: "monthAgingPercent" },
			{ text: "月超时数量", value: "monthTimeoutCount" },
			{ text: "日超时数量", value: "dayTimeoutCount" }
		]
	},

	{
		text: "质量",
		value: "quality",
		children: [
			{ text: "月质量准确率", value: "monthRightPercent" },
			{ text: "月差错数量", value: "monthErrorCount" },
			{ text: "日差错数量", value: "dayErrorCount" }
		]
	}
];

const fields = [
	{
		cols: 3,
		formKey: "reportDay",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		range: true,
		defaultValue: [moment().subtract(1, "day").format("YYYY-MM-DD")]
	}
];

const btns = [
	{
		class: "pr-3",
		color: "primary",
		text: "保存"
	},

	{
		class: "pr-3",
		color: "primary",
		text: "导出"
	}
];
const menuList = [
	{ label: "昨日", value: "0" },
	{ label: "前日", value: "1" }
];

export default {
	menuList,
	headers,
	fields,
	btns
};
