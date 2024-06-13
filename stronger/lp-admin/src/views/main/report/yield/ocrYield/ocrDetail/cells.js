const headers = [
	{ text: "日期", value: "date" },
	{ text: "字段名称", value: "fieldName" },
	{
		text: "汇总",
		value: "Summary",
		children: [
			{ text: "字符总量", value: "totalChara" },
			{ text: "有效字符总量", value: "effectChara" },
			{ text: "字符准确率", value: "accuracy" },
			{ text: "分块总量", value: "blockNum" },
			{ text: "有效分块数量", value: "blockEffect" },
			{ text: "分块准确率", value: "characterEffect" }
		]
	},
	{
		text: "一码",
		value: "one",
		children: [
			{ text: "字符总量", value: "oneTotalChara" },
			{ text: "有效字符总量", value: "oneEffectChara" },
			{ text: "字符准确率", value: "oneAccuracy" },
			{ text: "分块总量", value: "oneBlockNum" },
			{ text: "有效分块数量", value: "oneBlockEffect" },
			{ text: "分块准确率", value: "oneCharacterEffect" }
		]
	},
	{
		text: "二码",
		value: "two",
		children: [
			{ text: "字符总量", value: "twoTotalChara" },
			{ text: "有效字符总量", value: "twoEffectChara" },
			{ text: "字符准确率", value: "twoAccuracy" },
			{ text: "分块总量", value: "twoBlockNum" },
			{ text: "有效分块数量", value: "twoBlockEffect" },
			{ text: "分块准确率", value: "twoCharacterEffect" }
		]
	},
	{
		text: "问题件",
		value: "problem ",
		children: [
			{ text: "字符总量", value: "proTotalChara" },
			{ text: "有效字符总量", value: "proEffectChara" },
			{ text: "字符准确率", value: "proAccuracy" },
			{ text: "分块总量", value: "proBlockNum" },
			{ text: "有效分块数量", value: "proBlockEffect" },
			{ text: "分块准确率", value: "proCharacterEffect" }
		]
	}
];
const fields = [
	{
		cols: 2,
		formKey: "num1",
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
				label: "明细"
			}
		],
		defaultValue: 1
	},

	{
		cols: 2,
		formKey: "num2",
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
		defaultValue: 1
	},

	{
		cols: 3,
		formKey: "num3",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		range: true,
		defaultValue: undefined
	},

	{
		cols: 2,
		formKey: "num4",
		inputType: "input",
		hideDetails: true,
		label: "字段名称",
		defaultValue: undefined
	}
];

const btns = [
	{
		class: "pr-3",
		color: "primary",
		text: "复制"
	},

	{
		class: "pr-3",
		color: "primary",
		text: "导出"
	}
];

export default {
	headers,
	fields,
	btns
};
