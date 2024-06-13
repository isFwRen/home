import { request } from "@/api/service";

const actions = {
	// 列表
	async GET_PM_QUALITIES_ANALYSIS_LIST({}, params) {
		const data = {
			types: params.types,
			startTime: params.date[0],
			endTime: params.date[1]
		};

		const result = await request({
			url: "pro-manager/quality/analysis/list",
			params: data
		});

		return result;
	},

	// 导出
	async EXPORT_PM_QUALITIES_ANALYSIS({}, params) {
		const data = {
			startTime: params.date[0],
			endTime: params.date[1]
		};

		const result = await request({
			url: "pro-manager/quality/analysis/export",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
