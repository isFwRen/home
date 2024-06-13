import moment from "moment";
import { tools } from "@/libs/util";

// const { first, last } = tools.getLastDay()
// const today = new Date().getDate()

const fields = [
	{
		formKey: "date",
		inputType: "date",
		cols: 3,
		hideDetails: true,
		label: "时间",
		range: true,
		defaultValue: [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")]
		//defaultValue: [moment().add(-1, "day").format('YYYY-MM-DD'), moment().add(1, "day").format('YYYY-MM-DD')]
	},

	{
		formKey: "caseNumber",
		inputType: "input",
		cols: 2,
		hideDetails: true,
		label: "案件号",
		defaultValue: undefined
	},

	{
		formKey: "agency",
		inputType: "input",
		cols: 2,
		hideDetails: true,
		label: "机构号",
		defaultValue: undefined
	},

	{
		formKey: "caseStatus",
		inputType: "select",
		cols: 2,
		hideDetails: true,
		label: "案件状态",
		options: [
			{
				value: "",
				label: ""
			},
			{
				value: "2",
				label: "重复"
			},
			{
				value: "3",
				label: "异常"
			},
			{
				value: "4",
				label: "删除"
			}
		],
		defaultValue: undefined
	},

	{
		formKey: "stage",
		inputType: "select",
		cols: 2,
		hideDetails: true,
		label: "录入状态",
		options: [
			{
				value: "",
				label: ""
			},
			{
				value: "1",
				label: "待加载"
			},
			{
				value: "2",
				label: "录入中"
			},
			{
				value: "3",
				label: "已导出"
			},
			{
				value: "4",
				label: "待审核"
			}
		],
		defaultValue: undefined
	}
];

const claimTypes = [
  {
    label: '',
    value: -1
  },
  {
    label: '未定义',
    value: 0
  },
	{
    label: '无发票',
    value: 3
  },
  {
    label: '无报销',
    value: 4
  },
  {
    label: '有报销',
    value: 5
  },
  {
    label: '混合型',
    value: 6
  },
  {
    label: '简易',
    value: 7
  }
]


const headers = [
	{ text: "项目编码", value: "proCode" },
	{ text: "案件号", value: "caseNumber" },
	{ text: "机构", value: "agency" },
	{ text: "理赔类型", value: "claimType" },
	// 扫描时间调整为案件创建时间9-13
	{ text: "扫描时间", value: "createdAt", width: 100, sortable: true },
	{ text: "最晚回传时间", value: "backAtTheLatest" },
	{ text: "剩余时间", value: "timeRemaining" },
	{ text: "录入状态", value: "status" },
	{ text: "案件状态", value: "caseStatus" }
];

export default {
	fields,
	headers,
	claimTypes
};
