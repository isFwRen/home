import moment from "moment";

const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

const typeOptions = [
	{
		value: 1,
		label: "人员抽检"
	},

	{
		value: 2,
		label: "分块抽检"
	}
];

// 人员抽检
const headers1 = [
	{ text: "日期", value: "UpdatedAt" },
	{ text: "工号", value: "code" },
	{ text: "姓名", value: "name" },
	{ text: "抽检总数量", value: "num" },
	{ text: "已抽检数量", value: "doneNum" },
	{ text: "待抽检数量", value: "undoneNum" },
	{ text: "错误数量", value: "wrongNum" },
	{ text: "错误比例", value: "ratio" },
	{ text: "操作", value: "options", width: 170 }
];

// 分块抽检
const headers2 = [
	{ text: "日期", value: "UpdatedAt" },
	{ text: "项目", value: "pro" },
	{ text: "分块名称", value: "name" },
	{ text: "抽检总数量", value: "num" },
	{ text: "已抽检数量", value: "doneNum" },
	{ text: "待抽检数量", value: "undoneNum" },
	{ text: "错误数量", value: "wrongNum" },
	{ text: "错误比例", value: "ratio" },
	{ text: "操作", value: "options", width: 170 }
];

export default {
	DEFAULT_DATE,
	typeOptions,
	headers1,
	headers2
};
