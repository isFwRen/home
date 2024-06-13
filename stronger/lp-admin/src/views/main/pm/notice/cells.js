const headers = [
	{ text: "标题", value: "title", width: 120 },
	{ text: "发布类型", value: "releaseType", width: 90 },
	{ text: "项目编码", value: "proCode", width: 130 },
	{ text: "发布状态", value: "status", width: 90 },
	{ text: "发布时间", value: "releaseDate", width: 120 },
	{ text: "发布人", value: "releaseUserName", width: 90 },
	{ text: "访问人数", value: "visitCount", width: 90 },
	{ text: "操作", value: "options", width: 150 }
];

const types = [
	{ label: "公告通知", value: "1" },
	{ label: "规则动态", value: "2" }
];

const status = [
	{ label: "待发布", value: "1" },
	{ label: "已发布", value: "2" },
	{ label: "已删除", value: "3" }
];
export default {
	headers,
	types,
	status
};
