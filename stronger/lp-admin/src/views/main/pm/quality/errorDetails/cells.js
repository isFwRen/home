import moment from "moment";

const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

const headers = [
	{ text: "日期", value: "num1" },
	{ text: "项目", value: "num2" },
	{ text: "工号", value: "num3" },
	{ text: "姓名", value: "num4" },
	{ text: "案件号", value: "num5" },
	{ text: "字段名称", value: "num6" },
	{ text: "所属分块", value: "num7" },
	{ text: "类型", value: "num8" },
	{ text: "错误数据", value: "num9" },
	{ text: "正确数据", value: "num10" },
	{ text: "抽检人工号", value: "num11" },
	{ text: "抽检人姓名", value: "num12" }
];

export default {
	DEFAULT_DATE,
	headers
};
