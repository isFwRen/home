import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
	// 获取来量分析
	async GET_INSTITUTION_EXTRACTION_LIST({}, params) {
		const timeStart = moment(R.isYummy(params.time) ? params.time[0] : TODAY).format(
			"YYYY-MM-DD"
		);
		const timeEnd = moment(R.isYummy(params.time) ? params.time[1] : TODAY).format(
			"YYYY-MM-DD"
		);

		const [startTime, endTime] = [
			new Date(`${timeStart} 00:00:00`),
			new Date(`${timeEnd} 23:59:59`)
		];

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,  
			startTime,
			endTime,
		};

		const result = await request({
			url: "special-report/extract-agency/page",
			params: data
		});

		return result;
	},

	// 导出
	async EXPORT_INCOME_ANALYSIS({}, params) {
		const timeStart = moment(R.isYummy(params.time) ? params.time[0] : TODAY).format(
			"YYYY-MM-DD"
		);
		const timeEnd = moment(R.isYummy(params.time) ? params.time[1] : TODAY).format(
			"YYYY-MM-DD"
		);
		const [startTime , endTime ] = [
			new Date(`${timeStart} 00:00:00`),
			new Date(`${timeEnd} 23:59:59`)
		];

		const data = {
			startTime,
			endTime,
			proCode: params.proCode,
		};

		const result = await request({
			url: "special-report/extract-agency/export",
			params: data,
			responseType: "blob"
		});
		
		return result;
	},

	// 获取历史数据
	async EXPORT_HISTORY_ANALYSIS({}, params) {
		const timeStart = moment(R.isYummy(params.time) ? params.time[0] : TODAY).format(
			"YYYY-MM-DD"
		);
		const timeEnd = moment(R.isYummy(params.time) ? params.time[1] : TODAY).format(
			"YYYY-MM-DD"
		);
		const [startTime , endTime ] = [
			new Date(`${timeStart} 00:00:00`),
			new Date(`${timeEnd} 23:59:59`)
		];

		const data = {
			startTime,
			endTime,
			proCode: params.proCode,
		};

		const result = await request({
			url: "special-report/extract-agency/export-real-time",
			params: data,
			responseType: "blob"
		});

		return result;
	},

};

export default {
	actions
};
