import moment from "moment";

// 校验规则名
export const mapRules = {
	maxLen: "max_length",
	minLen: "min_length",
	fixLen: "fixed_length",
	fixValue: "included",
	maxVal: "max_value",
	minVal: "min_value",
	checkDate: "check_date",
	specChar: "free_value",

	1: "required",
	2: "spaced",
	3: "show_length",
	4: "digital",
	5: "amount",
	6: "chinese",
	7: "nonnegative",
	8: "alpha",
	9: "integer",
	10: "date",
	11: "id",
	12: "email",
	13: "phone",
	14: "name",

	includes: "included"
};

// 不需要校验的值
export const greenValue = ({ op, value }) => {
	if (!value) {
		return true;
	}

	if (op !== "opq") {
		if (value.includes("?")) {
			return true;
		}
	}

	return false;
};

const MS_DAY = 86400000;
const MS_TODAY = new Date(moment().format("YYYY-MM-DD")).getTime();

// 身份证
const ID_FIRST_TWO = [
	"11",
	"12",
	"13",
	"14",
	"15",
	"21",
	"22",
	"23",
	"31",
	"32",
	"33",
	"34",
	"35",
	"36",
	"37",
	"41",
	"42",
	"43",
	"44",
	"45",
	"46",
	"50",
	"51",
	"52",
	"53",
	"54",
	"61",
	"62",
	"63",
	"64",
	"65"
];
const ID_LAST = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "x", "X"];

// 日期时限
const CHECK_DATE_RULES = ["1", "2", "3", "4"];

export const validations = {
	// 最大长度
	max_length({ op, value, rule }) {
		if (greenValue({ op, value })) {
			return true;
		}

		value = "" + value;
		rule = +rule;

		if (value.length <= rule) {
			return true;
		}

		return `最大长度为${rule}.`;
	},

	// 最小长度
	min_length({ op, value, rule }) {
		if (greenValue({ op, value })) {
			return true;
		}

		value = "" + value;
		rule = +rule;

		if (value.length >= rule) {
			return true;
		}

		return `最小长度为${rule}.`;
	},

	// 固定长度
	fixed_length({ op, value, rule }) {
		if (greenValue({ op, value })) {
			return true;
		}

		value = "" + value;
		rule = +rule;

		if (value.length === rule) {
			return true;
		}

		return `长度只能为${rule}.`;
	},

	// 最大值
	max_value({ op, value, rule }) {
		if (greenValue({ op, value })) {
			return true;
		}

		if (+value <= +rule) {
			return true;
		}

		return `最大值为${rule}.`;
	},

	// 最小值
	min_value({ op, value, rule }) {
		if (greenValue({ op, value })) {
			return true;
		}

		if (+value >= +rule) {
			return true;
		}

		return `最小值为${rule}.`;
	},

	// 日期时限
	check_date({ label, value, rule }) {
		value = value + "";
		rule = rule + "";

		if (!CHECK_DATE_RULES.includes(rule)) {
			return true;
		}

		const result = rules.date({ label, value });

		if (result !== true) {
			return result;
		}

		const year = `20${value.substr(0, 2)}`;
		const month = value.substr(2, 2);
		const day = value.substr(4, 2);
		const customDate = new Date(
			moment(`${year}-${month}-${day}`).format("YYYY-MM-DD")
		).getTime();

		const diff = customDate - MS_TODAY;

		switch (rule) {
			case "1":
				if (diff < 0) {
					return true;
				}
				return `应早于今天.`;

			case "2":
				if (diff >= 0) {
					return true;
				}
				return `应不早于今天`;

			case "3":
				if (diff > MS_DAY) {
					return true;
				}
				return `应晚于今天`;

			case "4":
				if (diff <= MS_DAY) {
					return true;
				}
				return `应不晚于今天`;

			default:
				return true;
		}
	},

	// 跳过校验的值
	free_value({ op, value, rule, message }) {
		if (greenValue({ op, value })) {
			return true;
		}

		if (!rule.includes(value)) {
			return message;
		}

		return true;
	},

	// 必填
	required({ value }) {
		const result = /[\S]+/.test(value);

		if (result) {
			return true;
		}

		return `必填字段.`;
	},

	// 允许空格
	spaced() {
		return true;
	},

	// 显示长度
	show_length() {
		return true;
	},

	// 数字
	digital({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /^(\-|\+)?\d+(\.\d+)?$/.test(value);

		if (result) {
			return true;
		}

		return `只能为数字.`;
	},

	// 金额
	amount({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/.test(
			value
		);

		if (result) {
			return true;
		}

		return `不能小于零，且最多只能保留两位小数.`;
	},

	// 中文
	chinese({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /^[\u4e00-\u9fa5]+$/.test(value);

		if (result) {
			return true;
		}

		return `只能为中文.`;
	},

	// 不能录负数
	nonnegative({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		value = +value;

		if (value >= 0 || isNaN(value)) {
			return true;
		}

		return `不能为负数.`;
	},

	// 字母
	alpha({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /^[A-Za-z]+$/.test(value);

		if (result) {
			return true;
		}

		return `只能为字母.`;
	},

	// 整数
	integer({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /^-?[1-9]\d*$/.test(value);

		if (result) {
			return true;
		}

		return `只能为整数.`;
	},

	// 年月日
	date({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		value = value + "";

		if (value.length !== 6) {
			return `格式为YYMMDD.`;
		}

		const year = `20${value.substr(0, 2)}`;
		const month = value.substr(2, 2);
		const day = value.substr(4, 2);
		const date = `${year}-${month}-${day}`;

		if (moment(date).format("YYYY-MM-DD") === "Invalid date") {
			return "格式不合法.";
		}

		return true;
	},

	// 身份证
	id({ value }) {
		if (greenValue(value)) {
			return true;
		}

		value = value + "";

		const firstTwo = value.substr(0, 2);
		if (!~ID_FIRST_TWO.indexOf(firstTwo)) {
			return `证件不合法.`;
		}

		const sevenToFourteen = value.substr(6, 8);

		const last = value.substr(17, 1);
		if (!~ID_LAST.indexOf(last)) {
			return `证件最后一位只能为x、X或数字.`;
		}

		const length = value.length;
		if (length !== 18) {
			return `证件长度为18位.`;
		}

		return true;
	},

	// 必须包含指定的值
	included({ op, value, rule, message }) {
		if (greenValue({ op, value })) {
			return true;
		}

		if (!rule.includes(value)) {
			return message;
		}

		return true;
	},

	// 邮件
	email({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/.test(value);

		if (result) {
			return true;
		}

		return `邮件格式不正确.`;
	},

	// 手机
	phone({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const result = /^1(?:3\d|4[4-9]|5[0-35-9]|6[67]|7[013-8]|8\d|9\d)\d{8}$/.test(value);

		if (result) {
			return true;
		}

		return `手机格式不正确.`;
	},

	// 姓名
	name({ op, value }) {
		if (greenValue({ op, value })) {
			return true;
		}

		const enResult = /^[A-Za-z]+$/.test(value);
		const cnResult = /^[\u4e00-\u9fa5]+$/.test(value);

		if (enResult || cnResult) {
			return true;
		}

		return `姓名只能为中文或英文.`;
	}
};
