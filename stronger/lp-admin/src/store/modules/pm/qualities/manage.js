import { request } from "@/api/service";

const actions = {
	// 列表
	async GET_PM_QUALITIES_MANAGE_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: String(params.proCode),
			billName: params.billName,
			wrongFieldName: params.wrongFieldName,
			responsibleName: params.responsibleName,
			month: params.month
		};

		const result = await request({
			url: "pro-manager/quality/management/list",
			params: data
		});

		return result;
	},

	// 新增/编辑
	async POST_PM_QUALITIES_MANAGE_ITEM({}, body) {
		const formData = new FormData();

		for (let key in body) {
			formData.append(key, body[key]);
		}

		const result = await request({
			method: "POST",
			url: `pro-manager/quality/management/${!body.id ? "add" : "edit"}`,
			data: formData
		});

		return result;
	},

	// 删除
	async DELETE_PM_QUALITIES_MANAGE_ITEM({}, body) {
		const data = {
			ids: body.ids
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/quality/management/delete",
			data
		});

		return result;
	},

	// 导出
	async EXPORT_PM_QUALITIES_MANAGE_EXCEL({}, params) {
		const data = {
			proCode: String(params.proCode),
			month: params.month
		};

		const result = await request({
			url: "pro-manager/quality/management/export",
			params: data
		});

		return result;
	},

	// 根据工号获取员工
	async GET_PM_QUALITIES_MANAGE_STAFF({}, params) {
		const data = {
			code: params.code
		};

		const result = await request({
			url: "sys-base/user/find",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
