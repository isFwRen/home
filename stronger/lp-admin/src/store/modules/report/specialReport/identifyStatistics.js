import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
	async GET_IDENTIFY_STATISTICS_LIST({}, params) {
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
      billCode: params.billCode,
			startTime,
			endTime,
		};

		const result = await request({
			url: "report-management/ocr-output-Statistics/list",
			params: data
		});

		return result;
	},

	// 导出
	async EXPORT_IDENTIFY_STATISTICS({}, params) {
		const timeStart = `${params.time ? params.time[0] : TODAY}`;
		const timeEnd = `${params.time ? params.time[1] : TODAY}`;

		const [startTime , endTime ] = [
			new Date(`${timeStart} 00:00:00`),
			new Date(`${timeEnd} 23:59:59`)
		];


		const data = {
			startTime,
			endTime,
			proCode: params.proCode
		};

		const result = await request({
			url: "report-management/ocr-output-Statistics/export",
			params: data,
			responseType: "blob"
		});

		return result;
	}
};

export default {
	actions
};
