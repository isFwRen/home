import { request } from "@/api/service";
import { result } from "lodash";
import moment from "moment";

const actions = {
	async GET_HOME_WAIT_MAKE_DATA_LIST({}, { queryDay = moment().format("YYYY-MM-DD") }) {
		const data = {
			queryDay: queryDay
		};
		const result = await request({
			url: "homepage/pro-data/list",
			params: data
		});
		return result;
	},
	async GET_HOME_BUDINESS_RANKING_LIST({}, { rankingType = 1 }) {
		const data = {
			rankingType
		};
		const result = await request({
			url: "homepage/pro-data/business-ranking",
			params: data
		});
		return result;
	},
	async GET_HOME_AGING_TREND_LIST({}, { rankingType = 1 }) {
		const data = {
			rankingType
		};
		const result = await request({
			url: "homepage/pro-data/aging-trend",
			params: data
		});
		return result;
	}
};

export default {
	actions
};
