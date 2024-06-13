import Vue from "vue";
import moment from "moment-timezone";
import { R } from "vue-rocket";

// 保留小数点
const decimalFormat = (value, places = 2) => {
	const numValue = +value;

	if (R.getType(numValue) === "number") {
		return numValue.toFixed(places);
	}

	return value;
};

// 日期格式
const dateFormat = (date, format) => {
	// 时间不存在则不显示
	if (date && date.indexOf("0001") === 0) {
		return "-";
	}

	if (format) {
		return moment.tz(date, "Asia/Shanghai").format(format);
	}

	return date;
};

// 值是否合法
const ifLousyValue = (value, symbol = "-") => {
	if (R.isYummy(value)) {
		return value;
	}
	return symbol;
};

Vue.filter("decimalFormat", decimalFormat);
Vue.filter("dateFormat", dateFormat);
Vue.filter("ifLousyValue", ifLousyValue);
