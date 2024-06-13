const headers = [
	{ text: "排名", value: "rank", width: 90 },
	{ text: "项目", value: "proCode" },
	{ text: "超时数", value: "timeoutCount" },
	{
		text: "保障率",
		value: "billPercent",
		output: val => {
			return val * 100 + "%";
		}
	}
];

const menuList = [
	{ label: "今日", value: "0" },
	{ label: "本月", value: "1" }
];
export default {
	headers,
	menuList
};
