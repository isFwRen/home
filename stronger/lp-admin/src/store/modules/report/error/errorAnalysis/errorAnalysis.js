import { request } from "@/api/service";
import moment from "moment";

const actions = {
	async GET_ERRORANALYSIS_ITEM_LIST({}, data) {
		console.log("asd", data);
		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "/report-management/error-statistics/wrong-analysis/list",
			params: data
		});
		return result;
	},
	async EXPORT_ERRORANALYSIS_DETAIL({}, param) {
		const data = {
			startTime: (param.date ? param.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
			endTime: (param.date ? param.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
			proCode: param.proCode
		};
		const result = await request({
			url: "/report-management/error-statistics/wrong-analysis/export",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
