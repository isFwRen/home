import { request } from "@/api/service";

const actions = {
	// 列表
	async GET_PM_TEACHING_BUSINESS_RULES_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			ruleType: params.ruleType
		};

		const result = await request({
			url: "pro-manager/transactionRule/list",
			params: data
		});

		return result;
	},

	// 新增/修改
	async POST_PM_TEACHING_BUSINESS_RULES_ITEM({}, body) {
		const formData = new FormData();

		for (let key in body) {
			body[key] && formData.append(key, body[key]);
		}

		const result = await request({
			method: "POST",
			url: `pro-manager/transactionRule/${!body.id ? "add" : "edit"}`,
			data: formData
		});

		return result;
	},

	// 删除
	async DELETE_PM_TEACHING_BUSINESS_RULES_ITEM({}, body) {
		const data = {
			proCode: body.proCode,
			ids: body.ids
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/transactionRule/delete",
			data
		});

		return result;
	},

	async GET_PDF({}, url) {
		const result = await request({
			method: "get",
			url,
			responseType: "arraybuffer"
		});

		return result;
	}
};

export default {
	actions
};
