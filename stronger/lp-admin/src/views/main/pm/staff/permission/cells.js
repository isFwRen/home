const headers = [
	{ text: "项目编码", value: "proCode" },
	{ text: "项目名称", value: "proName" },
	{
		text: "人数",
		value: "NOP",
		children: [
			{ text: "初审", value: "op0" },
			{ text: "一码", value: "op1" },
			{ text: "二码", value: "op2" },
			{ text: "问题件", value: "opq" },
			// { text: "内网", value: "innet" },
			// { text: "外网", value: "outnet" }
		]
	}
];

export default {
	headers
};
