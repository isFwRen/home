import moment from "moment";

const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

const headers = [
	{ text: "日期", value: "UpdatedAt" },
	{ text: "总数量", value: "sum" },
	{ text: "抽检数量", value: "checkSum" },
	{ text: "抽查率", value: "ratio" },
	{ text: "错误数量", value: "wrongNum" },
	{ text: "错误率", value: "wrongRatio" }
];

export default {
	DEFAULT_DATE,
	headers
};
