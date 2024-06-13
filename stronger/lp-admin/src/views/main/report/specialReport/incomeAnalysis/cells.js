import moment from "moment";

const headers = [
	{ text: "项目", value: "proCode" },
	{ text: "00:00-8:30", value: "aHalfPastEight" },
	{ text: "8:30-9:30", value: "aHalfPastNine" },
	{ text: "9:30-10:30", value: "aHalfPastTen" },
	{ text: "10:30-11:30", value: "aHalfPastEleven" },
	{ text: "11:30-12:30", value: "aHalfPastTwelve" },
	{ text: "12:30-13:30", value: "aHalfPastThirteen" },
	{ text: "13:30-14:30", value: "aHalfPastFourteen" },
	{ text: "14:30-15:30", value: "aHalfPastFifteen" },
	{ text: "15:30-16:30", value: "aHalfPastSixteen" },
	{ text: "16:30-17:30", value: "aHalfPastSeventeen" },
	{ text: "17:30-18:30", value: "aHalfPastEighteen" },
	{ text: "18:30-00:00", value: "beforeZeroHour" }
];

const headers2 = [
	{ text: "日期", value: "createAt" },
	{ text: "00:00-8:30", value: "aHalfPastEight" },
	{ text: "8:30-9:30", value: "aHalfPastNine" },
	{ text: "9:30-10:30", value: "aHalfPastTen" },
	{ text: "10:30-11:30", value: "aHalfPastEleven" },
	{ text: "11:30-12:30", value: "aHalfPastTwelve" },
	{ text: "12:30-13:30", value: "aHalfPastThirteen" },
	{ text: "13:30-14:30", value: "aHalfPastFourteen" },
	{ text: "14:30-15:30", value: "aHalfPastFifteen" },
	{ text: "15:30-16:30", value: "aHalfPastSixteen" },
	{ text: "16:30-17:30", value: "aHalfPastSeventeen" },
	{ text: "17:30-18:30", value: "aHalfPastEighteen" },
	{ text: "18:30-00:00", value: "beforeZeroHour" }
];

const fields = [
	{
		cols: 2,
		formKey: "proCode",
		inputType: "select",
		hideDetails: true,
		label: "项目",
		options: [],
		defaultValue: undefined
	},

	{
		cols: 2,
		formKey: "time",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		range: true,
		defaultValue: [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")]
	}
];

const btns = [
	{
		class: "pr-3",
		color: "primary",
		text: "复制"
	},

	{
		class: "pr-3",
		color: "primary",
		text: "导出"
	}
];

export default {
	headers,
	headers2,
	fields,
	btns
};
