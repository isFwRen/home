import CryptoJS from "crypto-js";
import * as XLSX from "xlsx";
import axios from "axios";
import _ from "lodash";
import { localStorage } from "vue-rocket";
const tools = {};
import { TOTP } from "totp-generator";

// 请求头code
tools.GetCode = function (secret, uid) {
	const { otp } = TOTP.generate(secret, {
		digits: 8,
		algorithm: "SHA-512",
		period: 300
	});
	const code = otp.padStart(8, "0");
	console.log(code, "-code");
	return code;
};

tools.getTokenImg = async function (url) {
	const token = localStorage.get("token");
	const user = localStorage.get("user");
	const secret = localStorage.get("secret") || "";
	let code = null;
	if (secret) {
		code = tools.GetCode(secret);
	}

	// 请求图片
	const options = {
		headers: {
			"x-token": token,
			"x-code": String(code),
			"x-user-id": user?.id
		}
	};
	const request = axios.create(options);
	console.log(url, "url");
	const res = await request.get(url, {
		responseType: "blob"
	});
	return res.data;
};
// 将图片转成base64
tools.getBase64 = function (blob) {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.onloadend = () => {
			const base64String = reader.result;
			resolve(base64String);
		};
		reader.onerror = reject;
		reader.readAsDataURL(blob);
	});
};

/**
 * 判断内外网
 */
tools.isIntranet = function () {
	return location.hostname.includes("192.168") ? true : false;
};

/**
 * 基础路径
 */
tools.baseURL = function (devOrigin = "https://www.i-confluence.com:31111") {
	// const baseURL = `${location.hostname !== "localhost" ? location.origin : devOrigin}/`;
	// const baseURLApi = `${baseURL}api/`;
	const baseURL = "http://127.0.0.1:9999";
	const baseURLApi = "http://127.0.0.1:9999";

	return { baseURL, baseURLApi };
};

/**
 * 常量基础路径
 */
tools.constBaseURL = function () {
	const isIntranet = tools.isIntranet();
	const innerUrl = process.env.VUE_APP_CONST_INNER_URL;
	const outerUrl = process.env.VUE_APP_CONST_OUTER_URL;

	const baseURL = `https://admin:Change.Couchdb@${isIntranet ? innerUrl : outerUrl}/`;

	return { baseURL };
};

/**
 * 比较字符串
 */
tools.compareString = function (tValue, cValue, className = "warning--text") {
	if (tValue === "" && cValue === "") {
		return "";
	} else if (tValue === cValue) {
		return {
			targetHtml: tValue,
			diffValue: tValue,
			isDiff: false
		};
	}

	const [tLength, cLength] = [tValue.length, cValue.length];

	const maxLength = tLength >= cLength ? tLength : cLength;

	const tArr = tValue.split("");
	const cArr = cValue.split("");

	let [targetHtml, diffValue, firstDiffIndex] = ["", "", -1];

	for (let i = 0; i < maxLength; i++) {
		if (tArr[i] === cArr[i]) {
			targetHtml = targetHtml + tArr[i];
			diffValue = diffValue + tArr[i];
		} else {
			targetHtml = targetHtml + `<span class="${className}">${tArr[i] || ""}</span>`;

			diffValue = diffValue + "?";

			if (firstDiffIndex === -1) {
				firstDiffIndex = i;
			}
		}
	}

	return {
		targetHtml,
		diffValue,
		firstDiffIndex,
		isDiff: true
	};
};

/**
 * 设置输入框光标位置
 */
tools.setCursorPosition = function (el, index) {
	if (el.setSelectionRange) {
		el.focus();
		el.setSelectionRange(index, index + 1);
	} else if (el.createTextRange) {
		el.focus();
		var rang = el.createTextRange();
		rang.moveStart("character", index);
		rang.moveEnd("character", index + 1);
		rang.collapse(true);
		rang.select();
	}
};

/**
 * 获取某年某月的最后一天
 */
tools.getLastDay = function (year = null, month = null) {
	year = year ? year : new Date().getFullYear();
	month = month ? month : new Date().getMonth() + 1;

	const firstYMD = new Date(year, month, 1);
	const lastDay = new Date(firstYMD.getTime() - 1000 * 60 * 60 * 24).getDate();

	return {
		first: {
			year,
			month,
			day: 1
		},

		last: {
			year,
			month,
			day: lastDay
		}
	};
};

/**
 * 是否为正则表达式
 */
tools.isReg = function (value) {
	let isReg = null;

	try {
		isReg = eval(value) instanceof RegExp;
	} catch (error) {
		isReg = false;
	}

	return isReg;
};

/**
 * 解密数据
 */
tools.aesDecrypt = function (text) {
	let key = CryptoJS.enc.Utf8.parse("xingqiyistronger");
	let decryptedData = CryptoJS.AES.decrypt(text, key, {
		iv: key,
		mode: CryptoJS.mode.CBC,
		padding: CryptoJS.pad.Pkcs7
	});
	return decryptedData.toString(CryptoJS.enc.Utf8);
};

// 节流
tools.throttle = (() => {
	let timer = null;

	return (fn, delay = 100) => {
		if (timer) {
			return;
		}

		timer = setTimeout(() => {
			fn.apply(this, arguments);
			timer = null;
		}, delay);
	};
})();

let enabled = true;
tools.requestAnimationFrameFun = (fn, interval) => {
	if (enabled) {
		window.requestAnimationFrame(fn);
		enabled = false;
		// 每50ms触发一次，提高性能
		setTimeout(() => (enabled = true), interval);
	}
};

// 防抖
tools.debounce = (() => {
	let timer = null;

	return (fn, delay = 300) => {
		if (timer) {
			clearTimeout(timer);
		}

		timer = setTimeout(() => {
			fn();
		}, delay);
	};
})();

// 加载图片
tools.loadImage = function (source, func) {
	const image = new Image();
	image.setAttribute("crossOrigin", "anonymous");

	image.onload = function () {
		func(image.width, image.height);
	};

	image.onerror = function () {
		func(void 0, void 0);
		console.log("image load failed!");
	};

	image.src = source;
};

// 后端导出Excel文件 ，必须在axios请求中加参数 responseType: "blob"
tools.createExcelFun = function (res, name) {
	let blob = new Blob([res], {
		type: "application/vnd.ms-excel"
	});
	let fileName = name + ".xlsx";
	let link = document.createElement("a");
	link.download = fileName;
	link.href = window.URL.createObjectURL(blob);
	document.body.appendChild(link);
	link.click();
	window.URL.revokeObjectURL(link.href);
	return {
		msg: "导出成功"
	};
};

// 前端导出Excel表格
tools.exportExcel = function (name, tableName) {
	const table = document.querySelector(tableName);
	const sheetName = "Sheet1";

	const wb = XLSX.utils.book_new();
	const ws = XLSX.utils.table_to_sheet(table, { raw: true, sheet: sheetName });
	wb.SheetNames.push(sheetName);
	wb.Sheets[sheetName] = ws;

	// 导出 XLSX 文件
	XLSX.writeFile(wb, name);
};

export { tools };
export default tools;
