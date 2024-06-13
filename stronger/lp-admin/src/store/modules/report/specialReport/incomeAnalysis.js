import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
	// 获取来量分析
	async GET_INCOME_ANALYSIS_LIST({}, params) {
		const startTime = R.isYummy(params.time) ? params.time[0] : TODAY;
		const endTime = R.isYummy(params.time) ? params.time[1] : TODAY;
		const isCheckAll = !params.proCode ? true : false;

		const data = {
			proCode: params.proCode,
			startTime,
			endTime,
			types: "download",
			isCheckAll
		};

		const result = await request({
			url: "report-management/business-analysis/download/list",
			params: data
		});

		return result;
	},

	// 导出
	async EXPORT_INCOME_ANALYSIS_DETAIL({}, params) {
		const startTime = R.isYummy(params.time) ? params.time[0] : TODAY;
		const endTime = R.isYummy(params.time) ? params.time[1] : TODAY;
		const isCheckAll = !params.proCode ? true : false;

		const data = {
			proCode: params.proCode,
			startTime,
			endTime,
			types: "download",
			isCheckAll
		};

		const result = await request({
			url: "report-management/business-analysis/download/export",
			params: data,
			responseType: "blob"
		});

		return result;
	}
};

export default {
	actions
};
