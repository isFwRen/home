import moment from "moment";

const headers = [
	{
		text: "日期",
		value: "payDay"
	},

	{
		text: "工号",
		value: "code"
	},

	{
		text: "姓名",
		value: "nickName"
	},

	{
		text: "产量合计",
		value: "productionStatistics"
	},

	{
		text: "产量工资",
		value: "productionSalary"
	},

	{
		text: "客户投诉",
		value: "customerComplaints"
	},

	{
		text: "推荐奖",
		value: "referralBonus"
	},

	{
		text: "工资合计",
		value: "totalWages"
	},

	{
		text: "税额",
		value: "tax"
	},

	{
		text: "最终工资",
		value: "eventuallyPay"
	}
];

const fields = [
	{
		cols: 3,
		formKey: "date",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		pickerType: "month",
		defaultValue: moment().format("YYYY-MM")
	},

	{
		cols: 2,
		formKey: "code",
		inputType: "input",
		hideDetails: true,
		label: "工号",
		defaultValue: undefined
	},
	{
		cols: 2,
		formKey: "name",
		inputType: "input",
		hideDetails: true,
		label: "姓名",
		defaultValue: undefined
	}
];

const btns = [
	{
		class: "pr-3",
		key: 2,
		color: "primary",
		text: "导出"
	}
];

export default {
	headers,
	fields,
	btns
};
