const formList1 = [
	{
		name: "select",
		field: "name1",
		title: "理赔类型",
		options: [
			{
				label: "医疗",
				value: "医疗"
			}
		]
	},
	{
		name: "select",
		field: "name2",
		title: "治疗医院",
		options: [
			{
				label: "珠海市妇幼保健院",
				value: "珠海市妇幼保健院"
			}
		]
	},
	{
		name: "input",
		field: "name3",
		title: "入院科别"
	},
	{
		name: "input",
		field: "name4",
		title: "出院科别"
	},
	{
		name: "date",
		field: "name5",
		title: "入院时间"
	},
	{
		name: "date",
		field: "name6",
		title: "出院时间"
	},
	{
		name: "input",
		field: "name8",
		title: "实际住院天数"
	}
];

const formList2 = [
	{
		name: "input",
		field: "name1",
		title: "伤残鉴定书编号"
	},
	{
		name: "input",
		field: "name2",
		title: "伤残鉴定书机构"
	},
	{
		name: "date",
		field: "name3",
		title: "鉴定时间"
	},
	{
		name: "select",
		field: "name4",
		title: "鉴定伤残标准",
		options: [
			{
				label: "7-职业工伤与职业病致残1~10级伤残",
				value: '7-职业工伤与职业病致残1~10级伤残'
			}
		]
	},
	{
		name: "select",
		field: "name5",
		title: "伤残等级",
		options: [
			{
				label: "第一级",
				value: "第一级"
			},
			{
				label: "第二级",
				value: "第二级"
			},
			{
				label: "第三级",
				value: "第三级"
			},
			{
				label: "第四级",
				value: "第四级"
			},
			{
				label: "第五级",
				value: "第五级"
			},
			{
				label: "第六级",
				value: "第六级"
			},
			{
				label: "第七级",
				value: "第七级"
			},
			{
				label: "第八级",
				value: "第八级"
			},
			{
				label: "第九级",
				value: "第九级"
			},
			{
				label: "第十级",
				value: "第十级"
			},
		]
	},
	{
		name: "select",
		field: "name6",
		title: "伤残比例",
		options: [
			{
				label: "10%",
				value: "10%"
			},
			{
				label: "20%",
				value: "20%"
			},
			{
				label: "30%",
				value: "30%"
			},
			{
				label: "40%",
				value: "40%"
			},
			{
				label: "50%",
				value: "50%"
			},
			{
				label: "60%",
				value: "60%"
			},
			{
				label: "70%",
				value: "70%"
			},
			{
				label: "80%",
				value: "80%"
			},
			{
				label: "90%",
				value: "90%"
			},
			{
				label: "100%",
				value: "100%"
			},
		]
	},

	{
		name: "select",
		field: "name8",
		title: "鉴定结果",
		options: [
			{
				label: "下肢缺失",
				value: "下肢缺失"
			}
		]
	},
	{
		name: "input",
		field: "name9",
		title: "备注"
	}
];

const formList3 = [
	{
		name: "date",
		field: "name1",
		title: "身故日期",
		span: 8
	},
	{
		name: "input",
		field: "name2",
		title: "死亡地点",
		span: 8
	},
	{
		name: "input",
		field: "name3",
		title: "身故原因",
		span: 16
	}
];


const formRule1 = {
	name1: [{ required: true, message: "请输入理赔类型" }],
	name2: [{ required: true, message: "请输入治疗医院" }],
	name8: [{ required: true, message: "请输入实际住院天数" }]
};
const formRule2 = {
	name4: [{ required: true, message: "请输入鉴定伤残标准" }],
	name5: [{ required: true, message: "请输入伤残等级" }],
	name6: [{ required: true, message: "请输入伤残比例" }],
	name8: [{ required: true, message: "请输入鉴定结果" }]
};

export default {
	formList1,
	formList2,
	formList3,
	formRule1,
	formRule2
};
