import { request } from "@/api/service";
import moment from "moment";
import mock1 from "./mock.json";
import mock2 from "./mock1.json";

const state = {
	yieldInfo: {
		staffYield: {}
	},
	pro: {
		proOptions: [],
		proMap: {}
	}
};

const getters = {
	yieldInfo: () => state.yieldInfo,
	pro: () => state.pro
};

const mutations = {
	UPDATE_PRO(state, data) {
		const { pro } = state;
		state.pro = Object.assign(pro, data);
	}
};

const actions = {
	async GET_STAFF_YIELD_LIST({}, params) {
		const data = {};

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "report-management/corrected/list",
			params: data
		});

		return result;
	},

	async GET_STAFF_YIELD_LOGS({}, params) {
		const data = {
			correctedID: params.id
		};

		const result = await request({
			url: "report-management/corrected/list-log",
			params: data
		});

		return result;
	},

	// 新增/编辑设置表格
	async UPDATE_REPORT_YIELD_ITEM({}, form) {
		const { status } = form;

		const data = {
			ID: form.id,
			op0AsTheBlock: +form.op0AsTheBlock,
			op0AsTheInvoice: +form.op0AsTheInvoice,
			op1ExpenseAccount: +form.op1ExpenseAccount,
			op1NotExpenseAccount: +form.op1NotExpenseAccount,
			op2ExpenseAccount: +form.op2ExpenseAccount,
			op2NotExpenseAccount: +form.op2NotExpenseAccount,
			proCode: form.proCode,
			question: +form.question,
			startTime: form.startTime
		};

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			method: "POST",
			url: `report-management/corrected/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// Delete
	async DELETE_STAFF_YIELD({}, ids) {
		const result = await request({
			method: "DELETE",
			url: "report-management/corrected/delete",
			data: {
				ids: ids
			}
		});
		return result;
	},

	async GET_STAFF_TOTAL({}, params) {
		var updateTimes = "";

		if (params.updateTime !== undefined && params.updateTime !== "") {
			updateTimes = params.updateTime + " 00:00:00";
		}

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			startTime: (params.date ? params.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
			isCheckAll: params.type || 1,
			updateTime: updateTimes,
			code: params.code
		};
	

		const result = await request({
			url: "report-management/output-Statistics/list",
			params: data
		});

		return result;
	},

	async Delete_STAFF_TOTAL({}, params) {
		const data = {
			proCode: params.proCode,
			startTime: (params.date ? params.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
			code: params.code
		};
		const result = await request({
			method: "POST",
			//baseURL: 'http://127.0.0.1:9999/',
			url: "report-management/output-Statistics/delete",
			data
		});
		return result;
	},

	async GET_OCR_TOTAL({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			startTime: params.date ? params.date[0] : moment().format("yyyy-MM-DD"),
			endTime: params.date ? params.date[0] : moment().format("yyyy-MM-DD"),
			code: params.code,
			isCheckAll: params.currentType || 1
		};

		console.log(data);
		const result = await request({
			url: "report-management/output-Statistics/list",
			params: data
		});

		return result;
	},

	async EXPORT_STAFF_TOTAL({}, params) {
		const data = {
			startTime: (params.date ? params.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59"
		};
		const result = await request({
			url: "/report-management/output-Statistics/export",
			params: data
		});

		return result;
	},

	async EXPORT_STAFF_DETAIL({}, params) {
		const data = {
			startTime: (params.date ? params.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
			proCode: params.proCode
		};
		const result = await request({
			url: "/report-management/output-Statistics-detail/export",
			params: data,
			responseType: "blob"
		});

		return result;
	},

	// 新增/编辑设置表格
	async UPDATE_OUTPUT_STATISTICS({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			startTime: (params.date ? params.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
			code: params.code,
			isCheckAll: params.currentType || 1,
			IsCheckUpdate: params.IsCheckUpdate
		};
		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "report-management/output-Statistics/list",
			params: data
		});
		return result;
	}
};

export default {
	state,
	getters,
	mutations,
	actions
};
