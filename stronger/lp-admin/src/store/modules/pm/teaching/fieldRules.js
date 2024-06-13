import { request } from "@/api/service";

const actions = {
	// 列表
	async GET_PM_TEACHING_FIELD_RULES_LIST({}, params) {
		const data = {
			pageSize: params.pageSize,
			pageIndex: params.pageIndex,
			proCode: params.proCode,
			fieldsName: params.fieldsName,
			rule: params.rule
		};

		const result = await request({
			url: "pro-manager/fieldsRule/list",
			params: data
		});

		return result;
	},

	// 删除
	async DELETE_PM_TEACHING_FIELD_RULES_ITEMS({}, ids) {
		const result = await request({
			method: "POST",
			url: "pro-manager/fieldsRule/delete",
			data: {
				ids
			}
		});

		return result;
	},

	// 导出
	async EXPORT_PM_TEACHING_FIELD_RULES({}, body) {
		const data = {
			proCode: body.proCode,
			rule: body.rule
		};

		const result = await request({
			url: "pro-manager/fieldsRule/export",
			params: data
		});

		return result;
	},

	// 获取excel数据
	async GET_EXPORT_DATA({}, url) {
		const result = await request({
			url: url,
			responseType: "blob"
		});

		return result;
	}
};

export default {
	actions
};
