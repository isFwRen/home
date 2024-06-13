const headers = [
	{ text: "分块名称", value: "name", sortable: false, width: 100 },
	{ text: "初审", value: "op0Code", sortable: false },
	{ text: "初审完成", value: "op0SubmitAt", sortable: false },
	{ text: "一码", value: "op1Code", sortable: false, },
	{ text: "一码完成", value: "op1SubmitAt", sortable: false },
	{ text: "二码", value: "op2Code", sortable: false },
	{ text: "二码完成", value: "op2SubmitAt", sortable: false },
	{ text: "问题件", value: "opqCode", sortable: false },
	{ text: "问题件完成", value: "opqSubmitAt", sortable: false },
	{ text: "分块状态", value: "status", sortable: false },
	{ text: "操作", value: "options", sortable: false }
];

const handleDate = ['op0SubmitAt','op1SubmitAt','op2SubmitAt','opqSubmitAt']

export default {
	headers,
	handleDate
};
