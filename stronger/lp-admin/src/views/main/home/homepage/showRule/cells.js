import moment from "moment";

const headers = [
	{ text: "项目编码", value: "proCode" },
	{ text: "标题", value: "title" },
	{
		text: "发布时间",
		value: "CreatedAt",
		output: value => {
			return moment(value).format("YYYY-MM-DD HH:MM:SS");
		}
	},
	{ text: "操作", value: "opitions" }
];

export default {
	headers
};
