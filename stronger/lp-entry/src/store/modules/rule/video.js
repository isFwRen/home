import { request } from "@/api/service";

const actions = {
	async RULE_VIDEO_GET_LIST({}, params) {
		const data = {
			pageSize: params.pageSize,
			pageIndex: params.pageIndex,
			proCode: params.proCode,
			blockName: params.blockName,
			rule: params.ruleArr
		};
		if (Array.isArray(data["rule"])) {
			data["rule"] = data["rule"].length > 0 ? data["rule"][0] : "";
		}
    if (Array.isArray(data["proCode"])) {
			data["proCode"] = data["proCode"].length > 0 ? data["proCode"][0] : "";
		}

		const result = await request({
			url: "pro-manager/teachVideo/task/list",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
