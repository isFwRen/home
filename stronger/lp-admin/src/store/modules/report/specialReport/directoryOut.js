import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
	async GET_DIRECTORY_OUT_LIST({}, params) {
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
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
      type: params.type,
			startTime,
			endTime,
		};  

		const result = await request({
			url: "special-report/new-hospital-catalogue/page",
			params: data
		});

		return result;
	},

	// 导出
	async EXPORT_DIRECTORY_OUT({}, params) {
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
      type: params.type
		};

		const result = await request({
			url: "special-report/new-hospital-catalogue/export",
			params: data,
			responseType: "blob"
		});

		return result;
	}
};

export default {
	actions
};
