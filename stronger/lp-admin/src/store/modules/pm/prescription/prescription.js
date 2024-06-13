import { request } from "@/api/service";
import moment from "moment";

const actions = {
	async GET_PERSCRIPTION_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			startTime: (params.date ? params.date[0] : "2021-01-01") + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
			caseNumber: params.caseNumber,
			agency: params.agency,
			caseStatus: params.caseStatus,
			stage: params.stage,
			orderBy: params.orderBy
		};

		let result = await request({
			url: "pro-config/project-aging-management/list",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
