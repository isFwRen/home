import { request } from "@/api/service";

const actions = {
	async RULE_GET_TEMPLATE_LIST({}, params) {
		const data = {
			proCode: params.proCode,
			name: params.name
		};
		if (Array.isArray(data["proCode"])) {
			data["proCode"] = data["proCode"].length > 0 ? data["proCode"][0] : "";
		}
		const result = await request({
			url: "pro-manager/reimbursementFormTemplate/task/list",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
