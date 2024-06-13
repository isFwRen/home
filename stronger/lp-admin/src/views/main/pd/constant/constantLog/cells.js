const headers = [
	{ text: "操作类型", value: "type" },
	{ text: "操作内容", value: "content" },
	{ text: "操作人", value: "userName" },
	{ text: "操作时间", value: "UpdatedAt", width: 100, sortable: true }
];

const options = [
	{
		label: "查看",
		value: "查看"
	},
	{
		label: "上传",
		value: "上传"
	},
	{
		label: "编辑",
		value: "编辑"
	},
	{
		label: "删除",
		value: "删除"
	},
	{
		label: "导出",
		value: "导出"
	}
];

export default {
	headers,
	options
};
