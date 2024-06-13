const headers = [
	{ text: "项目编码", value: "proCode" },
	{ text: "案件号", value: "billNum" },
	{ text: "机构", value: "agency" },
	{ text: "扫描时间", value: "scanAt", width: 75, sortable: true },
	{ text: "最晚回传时间", value: "last_at", sortable: true },
	{ text: "录入完成时间", value: "appCompleteAt" },
	{ text: "剩余时间", value: "remainderAt", sortable: true },
	{ text: "录入状态", value: "stage" },
	{ text: "操作", value: "options" }
];


const handleDate = ['scanAt','appCompleteAt','last_at']

export default {
	headers,
	handleDate
};
