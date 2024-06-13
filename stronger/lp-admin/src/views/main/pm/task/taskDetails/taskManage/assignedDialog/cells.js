const headers = [
	{ text: "分块名称", value: "blockName", type: 'required' },
	{ text: "初审领取人", value: "op0Code", type: '初审' },
	{ text: "初审领取时间", value: "op0ApplyAt", type: '初审', sortable: true },
	{ text: "一码领取人", value: "op1Code", type: '一码' },
	{ text: "一码领取时间", value: "op1ApplyAt", type: '一码' },
	{ text: "二码领取人", value: "op2Code", type: '二码' },
	{ text: "二码领取时间", value: "op2ApplyAt", type: '二码' },
	{ text: "问题件领取人", value: "opqCode", type: '问题' },
	{ text: "问题件领取时间", value: "opqApplyAt", type: '问题' },
	{ text: "案件号", value: "billNum", type: 'required' },
	{ text: "机构号", value: "agency", type: 'required' },
	{ text: "操作", value: "options", width: 150, type: 'required' }
];

export default {
	headers
};
