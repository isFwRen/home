import { request } from "@/api/service";

const actions = {
	// 业务规则
	async RULE_GET_BUSINESS_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.projectArr,
			ruleType: params.ruleArr
		};
		if (Array.isArray(data["proCode"])) {
			data["proCode"] = data["proCode"].length > 0 ? data["proCode"][0] : "";
		}
		if (Array.isArray(data["ruleType"])) {
			data["ruleType"] = data["ruleType"].length > 0 ? data["ruleType"][0] : "";
		}
		const result = await request({
			url: "pro-manager/transactionRule/task/list",
			params: data
		});
		return result;
	},

	async GET_PDF_FILE({}, url) {
		const result = await request({
			url,
			responseType: "blob"
		});
		return result;
	}
};

export default {
	actions
};
