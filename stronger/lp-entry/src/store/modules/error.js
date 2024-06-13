import { request } from "@/api/service";
import moment from "moment";
import { localStorage } from "vue-rocket";

const actions = {
	// 列表
	async ERROR_GET_LIST({ }, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			startTime: (params.date ? params.date[0] : moment().format("YYYY-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("YYYY-MM-DD")) + " 23:59:59",
			fieldName: params.fieldName,
			complaint: params.complaint
		};

		const result = await request({
			url: "report-management/error-statistics/task/list",
			params: data
		});

		return result;
	},

	// 批量申诉
	async ERROR_APPEAL_ITEMS({ }, body) {
		const data = {
			proCode: body.proCode,
			complainConfirm: body.complainConfirm,
			list: body.list
		};

		const result = await request({
			method: "POST",
			url: "report-management/error-statistics/task/complain",
			data
		});

		return result;
	},

	// 练习错误
	async PRACTICE_ERROR_GET_LIST({ }, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			name: params.name,
			code: localStorage.get('user').code,
			startTime: (params.date ? params.date[0] : moment().format("YYYY-MM-DD")) + "T00:00:00.000Z",
			endTime: (params.date ? params.date[1] : moment().format("YYYY-MM-DD")) + "T23:59:59.000Z"
		};

		const result = await request({
			url: "/practice/wrong",
			params: data
		});

		return result;
	},
};

export default {
	actions
};
