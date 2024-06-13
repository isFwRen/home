import { request } from "@/api/service";
import { tools as lpTools } from "@/libs/util";
import { storage } from "@/libs/util";

import mock from "./mock.json";
import mock1 from "./mock1.json";
import mock2 from "./mock2.json";
import mock3 from "./mock3.json";

// 后面导出统一用GBK2312
const gb2312ProCodes = ["B0108", "B0121", "B0114"];

const { baseURLApi } = lpTools.baseURL();

const state = {
	cases: {
		rowInfo: {},
		caseInfo: {}
	}
};

const getters = {
	cases: () => state.cases
};

const mutations = {
	UPDATE_CASE(state, data) {
		const { cases } = state;
		state.cases = Object.assign(cases, data);
	}
};

const actions = {
	// 列表
	async GET_CASE_LIST({ }, params) {
		let [startTime, endTime] = params.time;

		const [timeStart, timeEnd] = [
			new Date(`${startTime} 00:00:00`),
			new Date(`${endTime} 23:59:59`)
		];

		const data = {
			orderBy: params.orderBy,
			proCode: params.proCode,
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			timeStart,
			timeEnd,
			billCode: params.billCode,
			status: params.status,
			saleChannel: params.saleChannel,
			batchNum: params.batchNum,
			agency: params.agency,
			insuranceType: params.insuranceType,
			claimType: params.claimType,
			stickLevel: params.stickLevel,
			minCountMoney: params.minCountMoney,
			maxCountMoney: params.maxCountMoney,
			isQuestion: params.isQuestion,
			invoiceNum: params.invoiceNum,
			qualityUser: params.qualityUser,
			stage: params.stage
		};

		if (data.invoiceNum == "") {
			data.invoiceNum = undefined;
		} else if (data.invoiceNum != undefined && !/^\d+$/.test(data.invoiceNum)) {
			data.invoiceNum = -100;
		}

		if (data.billCode != undefined) {
			data.billCode = data.billCode.trim();
		}

		const result = await request({
			url: "pro-manager/bill-list/page",
			params: data
		});

		return result;
	},

	// 删除
	async DELETE_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id,
			delRemarks: form.delRemarks + "-管理"
		};

		const result = await request({
			method: "DELETE",
			url: "pro-manager/bill-list/delete",
			data
		});

		return result;
	},

	// 重加载
	async RELOAD_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/reload",
			data
		});

		return result;
	},

	// 恢复
	async RECOVER_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/recover",
			data
		});

		return result;
	},

	// 导出异常
	async EXPORT_UNUSUAL_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/export-err-bill",
			data
		});

		return result;
	},

	// 强制导出
	async EXPORT_FORCE_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/force-export",
			data
		});

		return result;
	},

	// 判断是否有人在质检
	async CASE_JUDGE_IS_QUALITY({ }, body) {
		const data = {
			billId: body.id,
			proCode: body.proCode
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/is-quality",
			data
		});

		return result;
	},

	// 回传方式
	async AUTO_UPLOAD_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id,
			isAutoUpload: form.isAutoUpload
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/set-upload-type",
			data
		});

		return result;
	},

	// 回传
	async UPLOAD_CASE_ITEM({ }, form) {
		const data = {
			proCode: form.proCode,
			id: form.id
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/upload",
			data
		});

		return result;
	},

	// 获取XML
	GET_CASE_XML({ }, params) {
		const token = storage.get("token");
		const user = storage.get("user");
		const secret = storage.get("secret");
		const code = lpTools.GetCode(secret);
		const { proCode, year, month, date, billNum, types } = params;
		let [xhttp, serialized] = [null, null];

		if (window.XMLHttpRequest) {
			xhttp = new XMLHttpRequest();
		} else {
			xhttp = new window.ActiveXObject("Microsoft.XMLHTTP");
		}

		xhttp.open(
			"GET",
			`${baseURLApi}files/${proCode}/upload_xml/${year}/${month}/${date}/${billNum}.${types}?now=${Date.now()}`,
			false
		);

		xhttp.setRequestHeader("x-token", token);
		xhttp.setRequestHeader("x-user-id", user?.id);
		xhttp.setRequestHeader("x-code", String(code));

		let charset = "";
		if (proCode === "B0113" || proCode === "B0121") {
			charset = "utf-8";
		} else {
			// 后面导出都用gb2312编码
			charset = gb2312ProCodes.includes(proCode) ? "gb2312" : "utf-8";
		}

		xhttp.overrideMimeType(`text/csv;charset=${charset}`);

		// xhttp.setRequestHeader('Content-type', 'application/xml')

		try {
			xhttp.send();

			try {
				const serializer = new XMLSerializer();
				serialized = serializer.serializeToString(xhttp.responseXML);
			} catch (error) {
				serialized = xhttp.responseXML || xhttp.response;
			}
			return serialized;
		} catch (error) {
			return error;
		}
	},

	// 更新XML
	async UPDATE_CASE_XML({ }, body) {
		const data = {
			data: body.xml,
			url: `files/${body.proCode}/upload_xml/${body.year}/${body.month}/${body.date}/${body.billNum}.${body.types}`,
			proCode: body.proCode
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/save-xml",
			data
		});

		return result;
	},

	// 查看结果数据
	async GET_CASE_RESULT_DATA({ }, body) {
		const data = {
			proCode: body.proCode,
			id: body.id
		};

		const result = await request({
			url: "pro-manager/bill-list/see-bill-result-data",
			params: data
		});

		return result;
	},

	// 项目(select)
	async GET_CASE_PRO_LIST() {
		const result = await request({
			url: "pro-config/sys-project/list"
		});

		return result;
	},

	// 常量(select)
	async GET_CASE_CONST_LIST() {
		const result = await request({
			url: "pro-manager/bill-list/dict-const"
		});

		return result;
	},

	// 获取账单报表
	async GET_REPORTS({ }, form) {
		const data = {
			...form
		};

		const result = await request({
			url: "pro-manager/bill-list/qing-dan/get",
			params: data
		});

		return result;
	},

	// 更新备注
	async UPDATE_REMARKS({ }, data) {
		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/remark",
			data
		});

		return result;
	},

	async GET_IMG_API({ }, url) {
		console.log('getimg', url)
		const result = await request({
			method: "GET",
			url: url,
			responseType: "blob"
		});

		return result;
	},

	// 时效简报
	async GET_TIME_BRIEF({ }, params) {
		const result = await request({
			url: "pro-manager/bill-list/get-time-liness-briefing",
			params
		});

		return result;
	},

	// 获取质检信息
	async GET_INSPECT_INFO({ }, data) {
		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/get-qualities-demo",
			data
		});

		return result;
	},

	// 更新质检信息
	async UPDATE_INSPECT_INFO({ }, data) {
		const result = await request({
			method: "POST",
			url: "/pro-manager/bill-list/save-qualities-demo",
			data
		});

		return result;
	},
};

export default {
	state,
	getters,
	mutations,
	actions
};
