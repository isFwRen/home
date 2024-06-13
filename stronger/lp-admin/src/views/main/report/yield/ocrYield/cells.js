const headers = [
	{
		text: "项目",
		value: "project",
		children: [
			{ text: "合计", value: "total" },
			{ text: "占项目字符比例", value: "proportion" }
		]
	},
	{
		text: "B0101",
		value: "project1",
		children: [
			{ text: "汇总", value: "Summary" },
			{ text: "初审", value: "first" },
			{ text: "一码非报销单", value: "oneNoReim" },
			{ text: "一码报销单", value: "oneReim" },
			{ text: "二码非报销单", value: "twoNoReim" },
			{ text: "二码报销单", value: "twoReim" },
			{ text: "问题件", value: "Problem" }
		]
	},
	{
		text: "B0102",
		value: "project2",
		children: [
			{ text: "汇总", value: "Summary2" },
			{ text: "初审", value: "first2" },
			{ text: "一码非报销单", value: "oneNoReim2" },
			{ text: "一码报销单", value: "oneReim2" },
			{ text: "二码非报销单", value: "twoNoReim2" },
			{ text: "二码报销单", value: "twoReim2" },
			{ text: "问题件", value: "Problem2" }
		]
	}
];

// const fields = [
//   {
//     cols: 2,
//     formKey: 'num1',
//     inputType: 'select',
//     hideDetails: true,
//     label: '项目',
//     options: [
//       {
//         value: '1',
//         label: '全部',
//       },
//       {
//         value: '2',
//         label: '明细',
//       }
//     ],
//     defaultValue: 1
//   },

//   {
//     cols: 3,
//     formKey: 'num2',
//     inputType: 'date',
//     hideDetails: true,
//     label: '时间',
//     defaultValue: undefined
//   },

// ]

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
const typeOptions = [
	{
		label: "全部",
		value: 1
	},

	{
		label: "明细",
		value: 2
	}
];
export { typeOptions };
export default {
	headers,
	typeOptions,
	// fields,
	btns
};
