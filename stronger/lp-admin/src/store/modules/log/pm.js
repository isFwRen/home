import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
	// Table
	async GET_LOG_PM_LIST({}, params) {
		const startTime = new Date(`${R.isYummy(params.time) ? params.time[0] : TODAY} 00:00:00`);
		const endTime = new Date(`${R.isYummy(params.time) ? params.time[1] : TODAY} 23:59:59`);

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			startTime,
			endTime,
			functionModule: params.functionModule,
			moduleOperation: params.moduleOperation,
			operationPeople: params.operationCodeOrName,
			logType: 2
		};

		const result = await request({
			url: "/sys-logger/list",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
