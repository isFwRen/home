import moment from "moment";

const headers = [
	{ text: "单证号", value: "billNum" },
	{ text: "文件大小", value: "size" },
	{ text: "创建时间", value: "scanAt" },
	{ text: "删除时间", value: "delAt" },
	{ text: "描述文件", value: "describe" },
	{ text: "影像文件", value: "image" },
	{ text: "操作人", value: "name" },
	{ text: "说明", value: "stage" },
	{ text: "备注", value: "remarks" },
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
];

const tableData = [
	{ col1: '销毁文件原因', col2: '按照合同要求汇流保留7天的数据备份，数据备份期满汇流必须销毁数据', },
	{ col1: '文件删除时间范围', col2: '自 ________年____月____日 00 ：00 至自 ________年____月____日 24 ：00', },
	{ col1: '图像文件（单位：任务包）', col2: '' },
	{ col1: '生产数据（单位：表单）', col2: '', },
	{ col1: '图片大小（单位：MB）', col2: '', },
	{ col1: '删除审核日期', col2: '________年____月____日', },
	{ col1: '软件执行人', col2: '', },
	{ col1: '网络确认人', col2: '', },
	{ col1: '运营经理确认', col2: '中意理赔录入项目________年____月份生产环境原始图像及生产数据库数据已经确认全部进行从删除', },
	{ col1: '确认人签名', col2: '', },
]


export default {
	fields,
	headers,
	tableData
};
