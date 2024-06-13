const tab1 = {
	value: "basicInfo",
	label: "基础信息"
};

const tab2 = {
	value: "beneficiaryInfo",
	label: "受益人信息"
};

const tab3 = {
	value: "billInfo",
	label: "账单信息"
};

const tab4 = {
	value: "riskInfo",
	label: "出险信息"
};

const panels = [
	{
		value: 1,
		title: "申请人信息",
		content: ""
	},
	{
		value: 2,
		title: "被保人信息",
		content: ""
	},
	{
		value: 3,
		title: "受托人信息",
		content: ""
	},
	{
		value: 4,
		title: "其它信息",
		content: ""
	}
];

const fieldTypes = new Map([
	[1, "text"],
	[2, "const"],
	[3, "select"]
]);

export { fieldTypes };

export default {
	tab1,
	tab2,
	tab3,
	tab4,
	panels,
	fieldTypes
};
