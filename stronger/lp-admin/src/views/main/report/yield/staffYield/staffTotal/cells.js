const headers = [
	{
		text: "项目",
		value: "project",
		children: [
			{ text: "工号", value: "code", width: 80, sortable: true },
			{ text: "姓名", value: "nickName", width: 80, sortable: true },
			{ text: "合计", value: "addUpToSomething", width: 80, sortable: true }
		]
	}
	// {
	//   text: 'B0101',
	//   value: 'project1',
	//   children: [
	//     { text: '汇总', value: 'proSummary', },
	//     { text: '初审', value: 'first' },
	//     { text: '一码非报销单', value: 'oneNoReim', },
	//     { text: '一码报销单', value: 'oneReim' },
	//     { text: '二码非报销单', value: 'twoNoReim' },
	//     { text: '二码报销单', value: 'twoReim' },
	//     { text: '问题件', value: 'Problem' }
	//   ]
	// },
];

// const children = [
// 	{ text: "汇总", value: "mary" },
// 	{ text: "初审", value: "op0" },
// 	{ text: "一码非报销单", value: "op1NotExpenseAccount" },
// 	{ text: "一码报销单", value: "op1ExpenseAccount" },
// 	{ text: "二码非报销单", value: "op2NotExpenseAccount" },
// 	{ text: "二码报销单", value: "op2ExpenseAccount" },
// 	{ text: "问题件", value: "question" }
// ];
const children = [
	{ text: "汇总", value: "mary", width: 100 },
	{ text: "初审", value: "op0", width: 100 },
	{ text: "一码", value: "op1", width: 100 },
	{ text: "二码", value: "op2", width: 100 },
	{ text: "问题件", value: "question", width: 100 }
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
//     defaultValue: undefined
//   },

//   {
//     cols: 3,
//     formKey: 'num2',
//     inputType: 'date',
//     hideDetails: true,
//     label: '日期',
//     clearable: true,
//     range: true,
//     defaultValue: undefined
//   },

//   {
//     cols: 2,
//     formKey: 'num3',
//     inputType: 'input',
//     hideDetails: true,
//     label: '工号/姓名',
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
	},

	{
		class: "pr-3",
		color: "primary",
		text: "更新",
		onClick: "onUpdate"
	},

	{
		class: "pr-3",
		color: "primary",
		text: "设置",
		onClick: "onSet"
	}
];

export default {
	headers,
	// fields,
	children,
	btns
};
